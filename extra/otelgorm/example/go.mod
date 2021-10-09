module github.com/uptrace/uptrace-go/extra/otelgorm/example

go 1.17

replace github.com/uptrace/uptrace-go/extra/otelsql => ../../otelsql

replace github.com/uptrace/uptrace-go/extra/otelgorm => ./..

require (
	github.com/uptrace/uptrace-go/extra/otelgorm v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/otel v1.0.1
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.0.1
	go.opentelemetry.io/otel/sdk v1.0.1
	gorm.io/driver/sqlite v1.1.6
	gorm.io/gorm v1.21.15
)

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.2 // indirect
	github.com/mattn/go-sqlite3 v1.14.8 // indirect
	github.com/uptrace/uptrace-go/extra/otelsql v0.0.0-00010101000000-000000000000 // indirect
	go.opentelemetry.io/otel/internal/metric v0.24.0 // indirect
	go.opentelemetry.io/otel/metric v0.24.0 // indirect
	go.opentelemetry.io/otel/trace v1.0.1 // indirect
	golang.org/x/sys v0.0.0-20210423185535-09eb48e85fd7 // indirect
)
