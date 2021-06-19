module github.com/uptrace/uptrace-go/example/beego

go 1.15

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/astaxie/beego v1.12.3
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/prometheus/common v0.29.0 // indirect
	github.com/shiena/ansicolor v0.0.0-20200904210342-c7312218db18 // indirect
	github.com/uptrace/uptrace-go v0.20.0
	go.opentelemetry.io/contrib/instrumentation/github.com/astaxie/beego/otelbeego v0.21.0
	go.opentelemetry.io/otel v1.0.0-RC1
	go.opentelemetry.io/otel/trace v1.0.0-RC1
	golang.org/x/crypto v0.0.0-20210616213533-5ff15b29337e // indirect
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e // indirect
	golang.org/x/sys v0.0.0-20210616094352-59db8d763f22 // indirect
)
