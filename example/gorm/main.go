package main

import (
	"context"
	"fmt"

	"github.com/uptrace/uptrace-go/extra/otelgorm"
	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/otel"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {
	ctx := context.Background()

	// Configure OpenTelemetry with sensible defaults.
	uptrace.ConfigureOpentelemetry(
		// copy your project DSN here or use UPTRACE_DSN env var
		// uptrace.WithDSN("https://<key>@api.uptrace.dev/<project_id>"),

		uptrace.WithServiceName("myservice"),
		uptrace.WithServiceVersion("1.0.0"),
	)
	// Send buffered spans and free resources.
	defer uptrace.Shutdown(ctx)

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Install otelgorm plugin.
	if err := db.Use(otelgorm.NewPlugin()); err != nil {
		panic(err)
	}

	// Create a tracer. Usually, tracer is a global variable.
	tracer := otel.Tracer("app_or_package_name")

	// Create a root span (a trace) to measure some operation.
	ctx, main := tracer.Start(ctx, "main-operation")
	// End the span when the operation we are measuring is done.
	defer main.End()

	gormExample(ctx, db)

	fmt.Printf("trace: %s\n", uptrace.TraceURL(main))
}

func gormExample(ctx context.Context, db *gorm.DB) {
	// Migrate the schema
	db.WithContext(ctx).AutoMigrate(&Product{})

	// Create
	db.WithContext(ctx).Create(&Product{Code: "D42", Price: 100})

	// Read
	var product Product
	db.WithContext(ctx).First(&product, 1)                 // find product with integer primary key
	db.WithContext(ctx).First(&product, "code = ?", "D42") // find product with code D42

	// Update - update product's price to 200
	db.WithContext(ctx).Model(&product).Update("Price", 200)
	// Update - update multiple fields
	db.WithContext(ctx).Model(&product).Updates(Product{Price: 200, Code: "F42"}) // non-zero fields
	db.WithContext(ctx).Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// Delete - delete product
	db.WithContext(ctx).Delete(&product, 1)
}
