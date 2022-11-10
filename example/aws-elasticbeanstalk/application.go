package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		fmt.Fprintf(w, "Hello World! Append a name to the URL to say hello. For example, use %s/Mary to say hello to Mary.", r.Host)
	} else {
		fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
	}
	fmt.Fprintf(w, "\n%s", uptrace.TraceURL(trace.SpanFromContext(r.Context())))
}

func main() {
	uptrace.ConfigureOpentelemetry(
		uptrace.WithDSN("https://<token>@uptrace.dev/<project_id>"),

		uptrace.WithServiceName("myservice"),
		uptrace.WithServiceVersion("v1.0.0"),
		uptrace.WithDeploymentEnvironment("production"),
	)
	defer uptrace.Shutdown(context.TODO())

	handler := http.Handler(http.HandlerFunc(handlerFunc))
	handler = otelhttp.NewHandler(handler, "")

	httpServer := &http.Server{
		Addr:         ":5000",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
		Handler:      handler,
	}
	httpServer.ListenAndServe()
}
