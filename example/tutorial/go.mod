module github.com/uptrace/uptrace-go/example/tutorial

go 1.15

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/klauspost/compress v1.11.1 // indirect
	github.com/sirupsen/logrus v1.7.0 // indirect
	github.com/uptrace/uptrace-go v0.2.0
	github.com/vmihailenco/tagparser v0.1.2 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace v0.13.0 // indirect
	go.opentelemetry.io/otel v0.13.0
	go.opentelemetry.io/otel/sdk v0.13.0 // indirect
	golang.org/x/net v0.0.0-20201009032441-dbdefad45b89 // indirect
	golang.org/x/sys v0.0.0-20201009025420-dfb3f7c4e634 // indirect
	google.golang.org/grpc v1.33.0 // indirect
)
