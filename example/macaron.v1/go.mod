module github.com/uptrace/uptrace-go/example/macaron.v1

go 1.14

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/go-macaron/inject v0.0.0-20200308113650-138e5925c53b // indirect
	github.com/unknwon/com v1.0.1 // indirect
	github.com/uptrace/uptrace-go v0.19.1
	go.opentelemetry.io/contrib/instrumentation/gopkg.in/macaron.v1/otelmacaron v0.19.0
	go.opentelemetry.io/otel v0.19.0
	go.opentelemetry.io/otel/trace v0.19.0
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2 // indirect
	gopkg.in/ini.v1 v1.62.0 // indirect
	gopkg.in/macaron.v1 v1.4.0
)
