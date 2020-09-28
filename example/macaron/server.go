package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/uptrace/uptrace-go/uptrace"
	otelmacaron "go.opentelemetry.io/contrib/instrumentation/macaron"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/label"
	"gopkg.in/macaron.v1"
)

const profileTmpl = "profile"

var (
	upclient *uptrace.Client
	tracer   = global.Tracer("macaron-tracer")
)

func main() {
	ctx := context.Background()

	upclient = setupUptrace()
	defer upclient.Close()
	defer upclient.ReportPanic(ctx)

	upclient.ReportError(ctx, errors.New("hello from uptrace-go!"))

	m := macaron.Classic()
	m.Get("/profiles/:username", userProfileEndpoint)
	m.Use(otelmacaron.Middleware("service-name"))

	if err := router.Run(":9999"); err != nil {
		upclient.ReportError(ctx, err)
	}
}

func userProfileEndpoint() string {
	return "Hello world!"
}

func setupUptrace() *uptrace.Client {
	if os.Getenv("UPTRACE_DSN") == "" {
		panic("UPTRACE_DSN is empty or missing")
	}

	hostname, _ := os.Hostname()
	upclient := uptrace.NewClient(&uptrace.Config{
		// copy your project DSN here or use UPTRACE_DSN env var
		DSN: "",

		Resource: map[string]interface{}{
			"hostname": hostname,
		},
	})

	return upclient
}

func userProfileEndpoint(c *macaron.Context) string {
	ctx := c.Request.Context()

	username := c.Param("username")
	name, err := selectUser(ctx, username)
	if err != nil {
		ctx.Error(http.StatusNotFound, err)
		return
	}

	return fmt.Sprintf(`<html><h1>Hello %s %s </h1></html>`+"\n", username, name)
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
