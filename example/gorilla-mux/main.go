package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/uptrace/uptrace-go/uptrace"
)

var tracer = otel.Tracer("app_or_package_name")

func main() {
	ctx := context.Background()

	// Configure OpenTelemetry with sensible defaults.
	uptrace.ConfigureOpentelemetry(
		// copy your project DSN here or use UPTRACE_DSN env var
		// uptrace.WithDSN("https://<key>@api.uptrace.dev/<project_id>"),

		uptrace.WithServiceName("myservice"),
		uptrace.WithServiceVersion("1.0.0"),
	)
	// Send buffered spans and free resources.
	defer uptrace.Shutdown(ctx)

	r := mux.NewRouter()
	r.Use(otelmux.Middleware("service-name"))
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/hello/{username}", helloHandler)

	fmt.Println("running on http://localhost:9999")
	log.Fatal(http.ListenAndServe(":9999", r))
}

func indexHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

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
	fmt.Fprintf(w, tmpl, traceURL, traceURL)
}

func helloHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	traceURL := uptrace.TraceURL(trace.SpanFromContext(ctx))
	username := mux.Vars(req)["username"]
	tmpl := `
	<html>
	<h3>Hello %s</h3>
	<p><a href="%s" target="_blank">%s</a></p>
	</html>
	`
	fmt.Fprintf(w, tmpl, username, traceURL, traceURL)
}
