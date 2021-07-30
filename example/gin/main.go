package main

import (
	"context"
	"html/template"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/uptrace/uptrace-go/uptrace"
)

const indexTmpl = "index"
const profileTmpl = "profile"

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
	defer uptrace.Shutdown(ctx)

	router := gin.Default()
	router.SetHTMLTemplate(parseTemplates())
	router.Use(otelgin.Middleware("service-name"))
	router.GET("/", indexHandler)
	router.GET("/hello/:username", helloHandler)

	if err := router.Run("localhost:9999"); err != nil {
		log.Print(err)
	}
}

func parseTemplates() *template.Template {
	indexTemplate := `
		<html>
		<p>Here are some routes for you:</p>
		<ul>
			<li><a href="/hello/world">Hello world</a></li>
			<li><a href="/hello/foo-bar">Hello foo-bar</a></li>
		</ul>
		<p><a href="{{ .traceURL }}" target="_blank">{{ .traceURL }}</a></p>
		</html>
	`
	t := template.Must(template.New(indexTmpl).Parse(indexTemplate))

	profileTemplate := `
		<html>
		<h3>Hello {{ .username }}</h3>
		<p><a href="{{ .traceURL }}" target="_blank">{{ .traceURL }}</a></p>
		</html>
	`
	return template.Must(t.New(profileTmpl).Parse(profileTemplate))
}

func indexHandler(c *gin.Context) {
	ctx := c.Request.Context()
	otelgin.HTML(c, http.StatusOK, indexTmpl, gin.H{
		"traceURL": uptrace.TraceURL(trace.SpanFromContext(ctx)),
	})
}

func helloHandler(c *gin.Context) {
	ctx := c.Request.Context()
	otelgin.HTML(c, http.StatusOK, profileTmpl, gin.H{
		"username": c.Param("username"),
		"traceURL": uptrace.TraceURL(trace.SpanFromContext(ctx)),
	})
}
