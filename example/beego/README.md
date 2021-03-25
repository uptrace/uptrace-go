# Beego OpenTelemetry instrumentation example

[![PkgGoDev](https://pkg.go.dev/badge/go.opentelemetry.io/contrib/instrumentation/github.com/astaxie/beego/otelbeego)](https://pkg.go.dev/go.opentelemetry.io/contrib/instrumentation/github.com/astaxie/beego/otelbeego)

## Quickstart

Install
[otelbeego](https://github.com/open-telemetry/opentelemetry-go-contrib/tree/master/instrumentation/github.com/astaxie/beego/otelbeego)
instrumentation:

```bash
go get go.opentelemetry.io/contrib/instrumentation/github.com/astaxie/beego/otelbeego
```

Then install OpenTelemetry middleware:

```go
// To enable tracing on template rendering, disable autorender and
// call otelbeego.Render manually.
beego.BConfig.WebConfig.AutoRender = false

mware := otelbeego.NewOTelBeegoMiddleWare("service-name")
beego.RunWithMiddleWares(":7777", mware)
```

## Example

To run this example:

```bash
UPTRACE_DSN="https://<token>@api.uptrace.dev/<project_id>" go run main.go
```

Then open http://localhost:9999
