module github.com/uptrace/uptrace-go/example/beego

go 1.15

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/astaxie/beego v1.12.3
	github.com/prometheus/client_golang v1.9.0 // indirect
	github.com/prometheus/common v0.18.0 // indirect
	github.com/prometheus/procfs v0.6.0 // indirect
	github.com/shiena/ansicolor v0.0.0-20200904210342-c7312218db18 // indirect
	github.com/uptrace/uptrace-go v0.8.3
	go.opentelemetry.io/contrib/instrumentation/github.com/astaxie/beego/otelbeego v0.18.0
	go.opentelemetry.io/otel v0.18.0
	golang.org/x/crypto v0.0.0-20210220033148-5ea612d1eb83 // indirect
	golang.org/x/net v0.0.0-20210226172049-e18ecbb05110 // indirect
	golang.org/x/sys v0.0.0-20210305034016-7844c3c200c3 // indirect
	golang.org/x/text v0.3.5 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
