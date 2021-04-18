package main

import (
	"context"
	"fmt"
	"log"

	"github.com/go-pg/pg/extra/pgotel/v10"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"go.opentelemetry.io/otel"

	"github.com/uptrace/uptrace-go/uptrace"
)

var tracer = otel.Tracer("app_or_package_name")

func main() {
	ctx := context.Background()

	uptrace.ConfigureOpentelemetry(&uptrace.Config{
		// copy your project DSN here or use UPTRACE_DSN enar
		DSN: "",
	})
	defer uptrace.Shutdown(ctx)

	db := pg.Connect(&pg.Options{
		Addr:     ":5432",
		User:     "postgres",
		Database: "example",
	})
	defer db.Close()

	db.AddQueryHook(pgotel.NewTracingHook())

	if err := createBookTable(ctx, db); err != nil {
		log.Println(err)
		return
	}

	ctx, span := tracer.Start(ctx, "pg-main-span")
	defer span.End()

	if err := pgQueries(ctx, db); err != nil {
		log.Print(err)
		return
	}

	fmt.Println("trace", uptrace.TraceURL(span))
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
