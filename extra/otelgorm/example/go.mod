module github.com/uptrace/uptrace-go/extra/otelgorm/example

go 1.17

replace github.com/uptrace/uptrace-go/extra/otelsql => ../../otelsql

replace github.com/uptrace/uptrace-go/extra/otelgorm => ./..

require (
	github.com/uptrace/uptrace-go/extra/otelgorm v1.1.0
	go.opentelemetry.io/otel v1.2.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.2.0
	go.opentelemetry.io/otel/sdk v1.2.0
	gorm.io/driver/sqlite v1.2.4
	gorm.io/gorm v1.22.3
)

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.2 // indirect
	github.com/mattn/go-sqlite3 v1.14.9 // indirect
	github.com/uptrace/uptrace-go/extra/otelsql v1.1.0 // indirect
	go.opentelemetry.io/otel/internal/metric v0.25.0 // indirect
	go.opentelemetry.io/otel/metric v0.25.0 // indirect
	go.opentelemetry.io/otel/trace v1.2.0 // indirect
	golang.org/x/sys v0.0.0-20211116061358-0a5406a5449c // indirect
)
