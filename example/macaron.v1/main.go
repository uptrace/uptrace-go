package main

import (
	"context"
	"fmt"

	"go.opentelemetry.io/contrib/instrumentation/gopkg.in/macaron.v1/otelmacaron"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"gopkg.in/macaron.v1"

	"github.com/uptrace/uptrace-go/uptrace"
)

var tracer = otel.Tracer("app_or_package_name")

func main() {
	ctx := context.Background()

	uptrace.ConfigureOpentelemetry(&uptrace.Config{
		// copy your project DSN here or use UPTRACE_DSN env var
		DSN: "",
	})
	defer uptrace.Shutdown(ctx)

	m := macaron.Classic()
	m.Get("/", indexHandler)
	m.Get("/hello/:username", helloHandler)
	m.Use(otelmacaron.Middleware("service-name"))

	m.Run(9999)
}

func indexHandler(c *macaron.Context) string {
	ctx := c.Req.Context()

	traceURL := uptrace.TraceURL(trace.SpanFromContext(ctx))
	tmpl := `
	<html>
	<p>Here are some routes for you:</p>
	<ul>
		<li><a href="/hello/world">Hello world</a></li>
		<li><a href="/hello/foo-bar">Hello foo-bar</a></li>
	</ul>
	<p><a href="%s" target="_blank">%s</a></p>
	</html>
	`
	return fmt.Sprintf(tmpl, traceURL, traceURL)
}

func helloHandler(c *macaron.Context) string {
	ctx := c.Req.Context()

	traceURL := uptrace.TraceURL(trace.SpanFromContext(ctx))
	username := c.Params("username")
	tmpl := `
	<html>
	<h3>Hello %s</h3>
	<p><a href="%s" target="_blank">%s</a></p>
	</html>
	`
	return fmt.Sprintf(tmpl, username, traceURL, traceURL)
}
