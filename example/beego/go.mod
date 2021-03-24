module github.com/uptrace/uptrace-go/example/beego

go 1.15

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/astaxie/beego v1.12.3
	github.com/golang/protobuf v1.5.1 // indirect
	github.com/prometheus/client_golang v1.10.0 // indirect
	github.com/prometheus/common v0.20.0 // indirect
	github.com/shiena/ansicolor v0.0.0-20200904210342-c7312218db18 // indirect
	github.com/uptrace/uptrace-go v0.19.1
	go.opentelemetry.io/contrib/instrumentation/github.com/astaxie/beego/otelbeego v0.19.0
	go.opentelemetry.io/otel v0.19.0
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2 // indirect
	golang.org/x/net v0.0.0-20210324051636-2c4c8ecb7826 // indirect
	golang.org/x/sys v0.0.0-20210324051608-47abb6519492 // indirect
	golang.org/x/text v0.3.5 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
