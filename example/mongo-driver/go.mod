module github.com/uptrace/uptrace-go/example/mongo-driver

go 1.13

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/aws/aws-sdk-go v1.38.64 // indirect
	github.com/uptrace/uptrace-go v0.20.0
	github.com/youmark/pkcs8 v0.0.0-20201027041543-1326539a0a0a // indirect
	go.mongodb.org/mongo-driver v1.5.3
	go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo v0.21.0
	go.opentelemetry.io/otel v1.0.0-RC1
	golang.org/x/crypto v0.0.0-20210616213533-5ff15b29337e // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/text v0.3.6 // indirect
)
