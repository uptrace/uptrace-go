package main

import (
	"context"

	"github.com/astaxie/beego"
	"go.opentelemetry.io/contrib/instrumentation/github.com/astaxie/beego/otelbeego"
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

	// To enable tracing on template rendering, disable autorender and
	// call otelbeego.Render manually.
	beego.BConfig.WebConfig.AutoRender = false

	beego.Router("/", &IndexController{})
	beego.Router("/hello/:username", &HelloController{})

	mware := otelbeego.NewOTelBeegoMiddleWare("service-name")
	beego.RunWithMiddleWares("localhost:9999", mware)
}

type IndexController struct {
	beego.Controller
}

func (c *IndexController) Get() {
	ctx := c.Ctx.Request.Context()

	c.Data["traceURL"] = uptrace.TraceURL(trace.SpanFromContext(ctx))
	c.TplName = "index.tpl"

	if err := otelbeego.Render(&c.Controller); err != nil {
		c.Abort("500")
	}
}

type HelloController struct {
	beego.Controller
}

func (c *HelloController) Get() {
	ctx := c.Ctx.Request.Context()

	c.Data["username"] = c.Ctx.Input.Param(":username")
	c.Data["traceURL"] = uptrace.TraceURL(trace.SpanFromContext(ctx))
	c.TplName = "hello.tpl"

	if err := otelbeego.Render(&c.Controller); err != nil {
		c.Abort("500")
	}
}
