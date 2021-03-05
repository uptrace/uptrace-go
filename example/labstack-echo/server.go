package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("app_or_package_name")

func main() {
	ctx := context.Background()

	upclient := uptrace.NewClient(&uptrace.Config{
		// copy your project DSN here or use UPTRACE_DSN env var
		DSN: "",
	})
	defer upclient.Close()
	defer upclient.ReportPanic(ctx)

	e := echo.New()
	e.Use(otelecho.Middleware("service-name"))
	e.Use(middleware.Recover())
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		ctx := c.Request().Context()
		trace.SpanFromContext(ctx).RecordError(err)

		e.DefaultHTTPErrorHandler(err, c)
	}

	e.GET("/profiles/:username", userProfileEndpoint)

	e.Logger.Fatal(e.Start(":9999"))
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

	span.SetAttributes(attribute.String("username", username))

	if username == "admin" {
		return "Joe", nil
	}

	return "", fmt.Errorf("username=%s not found", username)
}
