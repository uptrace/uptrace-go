package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/label"
)

const profileTmpl = "profile"

var tracer = otel.Tracer("app_or_package_name")

func main() {
	ctx := context.Background()

	upclient := uptrace.NewClient(&uptrace.Config{
		// copy your project DSN here or use UPTRACE_DSN env var
		DSN: "",
	})
	defer upclient.Close()
	defer upclient.ReportPanic(ctx)

	router := gin.Default()
	router.SetHTMLTemplate(profileTemplate())
	router.Use(otelgin.Middleware("service-name"))
	router.GET("/profiles/:username", userProfileEndpoint)

	if err := router.Run(":9999"); err != nil {
		log.Print(err)
	}
}

func profileTemplate() *template.Template {
	tmpl := `<html><h1>Hello {{ .username }} {{ .name }}</h1></html>` + "\n"
	return template.Must(template.New(profileTmpl).Parse(tmpl))
}

func userProfileEndpoint(c *gin.Context) {
	ctx := c.Request.Context()

	username := c.Param("username")
	name, err := selectUser(ctx, username)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	otelgin.HTML(c, http.StatusOK, profileTmpl, gin.H{
		"username": username,
		"name":     name,
	})
}

func selectUser(ctx context.Context, username string) (string, error) {
	_, span := tracer.Start(ctx, "selectUser")
	defer span.End()

	span.SetAttributes(label.String("username", username))

	if username == "admin" {
		return "Joe", nil
	}

	return "", fmt.Errorf("username=%s not found", username)
}
