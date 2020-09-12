package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/uptrace/uptrace-go/uptrace"
	otelecho "go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/label"
)

const profileTmpl = "profile"

var (
	upclient *uptrace.Client
	tracer   = global.Tracer("echo-tracer")
)

func main() {
	ctx := context.Background()

	upclient = setupUptrace()

	defer upclient.Close()
	defer upclient.ReportPanic(ctx)

	upclient.ReportError(ctx, errors.New("hello from echo-go!"))

	e := echo.New()
	e.Use(otelecho.Middleware("service-name"))
	e.Use(middleware.Recover())
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		ctx := c.Request().Context()
		trace.SpanFromContext(ctx).RecordError(ctx, err)

		e.DefaultHTTPErrorHandler(err, c)
	}

	e.GET("/profiles/:username", userProfileEndpoint)

	e.Logger.Fatal(e.Start(":9999"))
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

func userProfileEndpoint(c echo.Context) error {
	ctx := c.Request().Context()

	username := c.Param("username")
	name, err := selectUser(ctx, username)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	html := fmt.Sprintf(`<html><h1>Hello %s %s </h1></html>`+"\n", username, name)
	return c.HTML(http.StatusOK, html)
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
