module github.com/uptrace/uptrace-go/example/mongo-driver

go 1.13

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/uptrace/uptrace-go v1.0.0-RC3
	github.com/youmark/pkcs8 v0.0.0-20201027041543-1326539a0a0a // indirect
	go.mongodb.org/mongo-driver v1.7.1
	go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo v0.22.0
	go.opentelemetry.io/otel v1.0.0-RC2
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97 // indirect
)
