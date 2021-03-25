package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

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

	var handler http.Handler

	handler = http.HandlerFunc(indexHandler)
	handler = otelhttp.WithRouteTag("/", handler)
	handler = otelhttp.NewHandler(handler, "index-handler")
	http.HandleFunc("/", handler.ServeHTTP)

	handler = http.HandlerFunc(helloHandler)
	handler = otelhttp.WithRouteTag("/hello/:username", handler)
	handler = otelhttp.NewHandler(handler, "hello-handler")
	http.HandleFunc("/hello/", handler.ServeHTTP)

	srv := &http.Server{
		Addr: ":9999",
	}

	fmt.Println("listening on http://localhost:9999")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
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

	username := strings.Replace(req.URL.Path, "/hello/", "", 1)

	traceURL := uptrace.TraceURL(trace.SpanFromContext(ctx))
	tmpl := `
	<html>
	<h3>Hello %s</h3>
	<p><a href="%s" target="_blank">%s</a></p>
	</html>
	`
	fmt.Fprintf(w, tmpl, username, traceURL, traceURL)
}
