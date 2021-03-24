package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/uptrace/uptrace-go/uptrace"
)

const indexTmpl = "index"
const profileTmpl = "profile"

var tracer = otel.Tracer("app_or_package_name")

func main() {
	ctx := context.Background()

	uptrace.ConfigureOpentelemetry(&uptrace.Config{
		// copy your project DSN here or use UPTRACE_DSN env var
		DSN: "",
	})
	defer uptrace.Shutdown(ctx)

	router := gin.Default()
	router.SetHTMLTemplate(loadTemplates())
	router.Use(otelgin.Middleware("service-name"))
	router.GET("/", indexEndpoint)
	router.GET("/hello/:username", userProfileEndpoint)

	if err := router.Run("localhost:9999"); err != nil {
		log.Print(err)
	}
}

func loadTemplates() *template.Template {
	indexTemplate := `
		<html>
		<p>Here are some routes for you:</p>

		<ul>
			<li><a href="/hello/admin">Hello admin</a></li>
			<li><a href="/hello/unknown">Hello unknown user</a></li>
		</ul>
	
		<p><a href="{{ .url }}" target="_blank">{{ .url }}</a></p>
		</html>
	`
	t := template.Must(template.New(indexTmpl).Parse(indexTemplate))

	profileTemplate := `
		<html>
		<h1>Hello {{ .username }} {{ .name }}</h1>

		<p><a href="{{ .url }}" target="_blank">{{ .url }}</a></p>
		</html>
	`
	return template.Must(t.New(profileTmpl).Parse(profileTemplate))
}

func indexEndpoint(c *gin.Context) {
	ctx := c.Request.Context()

	traceURL := uptrace.TraceURL(trace.SpanFromContext(ctx))
	otelgin.HTML(c, http.StatusOK, indexTmpl, gin.H{
		"url": traceURL,
	})
}

func userProfileEndpoint(c *gin.Context) {
	ctx := c.Request.Context()

	traceURL := uptrace.TraceURL(trace.SpanFromContext(ctx))

	username := c.Param("username")
	name, err := selectUser(ctx, username)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
			"code":  http.StatusNotFound,
			"url":   traceURL,
		})
		return
	}

	otelgin.HTML(c, http.StatusOK, profileTmpl, gin.H{
		"username": username,
		"name":     name,
		"url":      traceURL,
	})
}

func selectUser(ctx context.Context, username string) (string, error) {
	_, span := tracer.Start(ctx, "selectUser")
	defer span.End()

	span.SetAttributes(attribute.String("username", username))

	if username == "unknown" {
		return "", fmt.Errorf("username=%s not found", username)
	}

	return "Joe", nil
}
