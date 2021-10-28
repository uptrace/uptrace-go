module github.com/uptrace/uptrace-go/example/gorm

go 1.16

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/uptrace/uptrace-go v1.0.5
	github.com/uptrace/uptrace-go/extra/otelgorm v1.0.5
	go.opentelemetry.io/otel v1.0.1
	gorm.io/driver/sqlite v1.1.6
	gorm.io/gorm v1.21.16
)

require (
	github.com/cenkalti/backoff/v4 v4.1.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.2 // indirect
	github.com/mattn/go-sqlite3 v1.14.8 // indirect
	github.com/uptrace/uptrace-go/extra/otelsql v1.0.5 // indirect
	go.opentelemetry.io/contrib/instrumentation/runtime v0.25.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric v0.24.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v0.24.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.0.1 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.0.1 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.0.1 // indirect
	go.opentelemetry.io/otel/internal/metric v0.24.0 // indirect
	go.opentelemetry.io/otel/metric v0.24.0 // indirect
	go.opentelemetry.io/otel/sdk v1.0.1 // indirect
	go.opentelemetry.io/otel/sdk/export/metric v0.24.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v0.24.0 // indirect
	go.opentelemetry.io/otel/trace v1.0.1 // indirect
	go.opentelemetry.io/proto/otlp v0.9.0 // indirect
	golang.org/x/net v0.0.0-20211008194852-3b03d305991f // indirect
	golang.org/x/sys v0.0.0-20211007075335-d3039528d8ac // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20211008145708-270636b82663 // indirect
	google.golang.org/grpc v1.41.0 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
)
