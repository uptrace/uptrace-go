package main

import (
	"context"
	"log"
	"os"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/go-pg/pgext"
	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/otel/api/global"
)

var tracer = global.Tracer("go-pg-tracer")

func main() {
	ctx := context.Background()

	upclient := setupUptrace()
	defer upclient.Close()
	defer upclient.ReportPanic(ctx)

	options := &pg.Options{
		Addr:     "postgresql-server:5432",
		User:     "postgres",
		Database: "example",
	}
	db := pg.Connect(options)
	db.AddQueryHook(pgext.OpenTelemetryHook{})
	defer db.Close()

	if err := createBookTable(ctx, db); err != nil {
		upclient.ReportError(ctx, err)
		log.Println(err.Error())
		return
	}

	ctx, span := tracer.Start(ctx, "pg-main-span")
	defer span.End()

	if err := pgQueries(ctx, db); err != nil {
		upclient.ReportError(ctx, err)
		log.Println(err.Error())
		return
	}

	log.Println("trace", upclient.TraceURL(span))
}

func setupUptrace() *uptrace.Client {
	if os.Getenv("UPTRACE_DSN") == "" {
		panic("UPTRACE_DSN is empty or missing")
	}

	hostname, _ := os.Hostname()
	upclient := uptrace.NewClient(&uptrace.Config{
		// copy your project DSN here or use UPTRACE_DSN enar
		DSN: "",

		Resource: map[string]interface{}{
			"host.name": hostname,
		},
	})

	return upclient
}

type Book struct {
	ID              int
	Title           string
	AuthorFirstName string
	AuthorLastName  string
}

func pgQueries(ctx context.Context, db *pg.DB) error {
	book := &Book{
		Title:           "Harry Potter",
		AuthorFirstName: "Rowling",
		AuthorLastName:  "Joanne",
	}
	_, err := db.ModelContext(ctx, book).Insert()
	if err != nil {
		return err
	}

	_, err = db.ModelContext(ctx, book).
		Set("title = ?", "Harry Potter and the Deathly Hallows").
		Where("id = ?", book.ID).
		Update()
	if err != nil {
		return err
	}

	_, err = db.ModelContext(ctx, book).
		Where("id = ?", book.ID).
		Delete()
	if err != nil {
		return err
	}

	return nil
}

func createBookTable(ctx context.Context, db *pg.DB) error {
	if err := db.ModelContext(ctx, (*Book)(nil)).DropTable(&orm.DropTableOptions{
		IfExists: true,
	}); err != nil {
		return err
	}

	if err := db.ModelContext(ctx, (*Book)(nil)).CreateTable(nil); err != nil {
		return err
	}

	return nil
}
