module github.com/uptrace/uptrace-go/example/otlp

go 1.15

require (
	go.opentelemetry.io/otel v1.0.0-RC1
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.0.0-RC1
	go.opentelemetry.io/otel/sdk v1.0.0-RC1
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e // indirect
	golang.org/x/sys v0.0.0-20210616094352-59db8d763f22 // indirect
	google.golang.org/genproto v0.0.0-20210617175327-b9e0b3197ced // indirect
	google.golang.org/grpc v1.38.0
)
