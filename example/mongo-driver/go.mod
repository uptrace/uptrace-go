module github.com/uptrace/uptrace-go/example/mongo-driver

go 1.13

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/aws/aws-sdk-go v1.38.3 // indirect
	github.com/golang/snappy v0.0.3 // indirect
	github.com/uptrace/uptrace-go v0.19.1
	github.com/youmark/pkcs8 v0.0.0-20201027041543-1326539a0a0a // indirect
	go.mongodb.org/mongo-driver v1.5.0
	go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo v0.19.0
	go.opentelemetry.io/otel v0.19.0
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/text v0.3.5 // indirect
)
