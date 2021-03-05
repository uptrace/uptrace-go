module github.com/uptrace/uptrace-go/example/mongo-driver

go 1.13

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/aws/aws-sdk-go v1.37.24 // indirect
	github.com/golang/snappy v0.0.3 // indirect
	github.com/uptrace/uptrace-go v0.8.3
	go.mongodb.org/mongo-driver v1.4.6
	go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo v0.18.0
	go.opentelemetry.io/otel v0.18.0
	golang.org/x/crypto v0.0.0-20210220033148-5ea612d1eb83 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/text v0.3.5 // indirect
)
