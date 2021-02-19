module github.com/uptrace/uptrace-go/example/mongo-driver

go 1.13

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/aws/aws-sdk-go v1.37.14 // indirect
	github.com/golang/snappy v0.0.2 // indirect
	github.com/uptrace/uptrace-go v0.7.7
	go.mongodb.org/mongo-driver v1.4.6
	go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo v0.17.0
	go.opentelemetry.io/otel v0.17.0
	golang.org/x/crypto v0.0.0-20210218145215-b8e89b74b9df // indirect
	golang.org/x/sync v0.0.0-20201207232520-09787c993a3a // indirect
	golang.org/x/text v0.3.5 // indirect
)
