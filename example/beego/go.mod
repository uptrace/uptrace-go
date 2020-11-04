module github.com/uptrace/uptrace-go/example/beego

go 1.15

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/astaxie/beego v1.12.2
	github.com/prometheus/common v0.14.0 // indirect
	github.com/prometheus/procfs v0.2.0 // indirect
	github.com/shiena/ansicolor v0.0.0-20200904210342-c7312218db18 // indirect
	github.com/uptrace/uptrace-go v0.4.2
	go.opentelemetry.io/contrib/instrumentation/github.com/astaxie/beego/otelbeego v0.13.0
	go.opentelemetry.io/otel v0.13.0
	golang.org/x/crypto v0.0.0-20201002170205-7f63de1d35b0 // indirect
)
