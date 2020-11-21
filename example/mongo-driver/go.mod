module github.com/uptrace/uptrace-go/example/mongo-driver

go 1.13

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/aws/aws-sdk-go v1.34.32 // indirect
	github.com/golang/snappy v0.0.2 // indirect
	github.com/uptrace/uptrace-go v0.3.0
	go.mongodb.org/mongo-driver v1.4.2
	go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo v0.14.0
	go.opentelemetry.io/otel v0.14.0
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a // indirect
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208 // indirect
)
