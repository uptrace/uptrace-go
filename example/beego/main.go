package main

import (
	"context"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/contrib/instrumentation/github.com/astaxie/beego/otelbeego"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
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

	beego.Router("/profiles/:username", &ProfileController{})

	mware := otelbeego.NewOTelBeegoMiddleWare("service-name")
	beego.RunWithMiddleWares("localhost:9999", mware)
}

type ProfileController struct {
	beego.Controller
}

func (c *ProfileController) Get() {
	ctx := c.Ctx.Request.Context()

	username := c.Ctx.Input.Param(":username")
	name, err := selectUser(ctx, username)
	if err != nil {
		c.Abort("404")
		return
	}

	c.Data["username"] = username
	c.Data["name"] = name
	c.TplName = "hello.tpl"

	// Don't forget to call render manually.
	if err := otelbeego.Render(&c.Controller); err != nil {
		c.Abort("500")
	}
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
