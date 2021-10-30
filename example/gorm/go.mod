module github.com/uptrace/uptrace-go/example/gorm

go 1.16

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/uptrace/uptrace-go v1.0.5
	go.opentelemetry.io/otel v1.1.0
	gorm.io/driver/sqlite v1.2.3
	gorm.io/gorm v1.22.2
)

require github.com/uptrace/opentelemetry-go-extra/otelgorm v0.1.2
