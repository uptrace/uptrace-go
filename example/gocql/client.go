package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gocql/gocql"
	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gocql/gocql/otelgocql"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var (
	upclient *uptrace.Client
	tracer   = otel.Tracer("gocql-tracer")
)

const keyspace = "gocql_example"

func main() {
	ctx := context.Background()

	upclient = setupUptrace()
	defer upclient.Close()
	defer upclient.ReportPanic(ctx)

	err := initDB()
	if err != nil {
		log.Fatal(err)
	}

	cluster := newCassandraCluster(keyspace)
	session, err := otelgocql.NewSessionWithTracing(ctx, cluster)
	if err != nil {
		log.Fatal(err)
		upclient.ReportError(ctx, err)
	}
	defer session.Close()

	traceGocqlQueries(ctx, session)

	if err := truncateTable(ctx, session); err != nil {
		log.Fatal(err)
	}
}

func traceGocqlQueries(ctx context.Context, session *gocql.Session) {
	ctx, span := tracer.Start(ctx, "test-operations")
	defer span.End()

	insertBooks(ctx, session)
	bookID := selectBook(ctx, session)
	updateBook(ctx, session, bookID)
	deleteBook(ctx, session, bookID)
}

func newCassandraCluster(keyspace string) *gocql.ClusterConfig {
	cluster := gocql.NewCluster("cassandra-server")
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.LocalQuorum
	cluster.ProtoVersion = 3
	cluster.Timeout = 2 * time.Second
	return cluster
}

func setupUptrace() *uptrace.Client {
	upclient := uptrace.NewClient(&uptrace.Config{
		// copy your project DSN here or use UPTRACE_DSN enar
		DSN: "",
	})

	return upclient
}

func insertBooks(ctx context.Context, session *gocql.Session) {
	batch := session.NewBatch(gocql.LoggedBatch)
	for i := 0; i < 5; i++ {
		batch.Query(
			"INSERT INTO book (id, title, author_first_name, author_last_name) VALUES (?, ?, ?, ?)",
			gocql.TimeUUID(),
			fmt.Sprintf("Example Book %d", i),
			fmt.Sprintf("author_last_name %d", i),
			fmt.Sprintf("author_last_name %d", i),
		)
	}
	if err := session.ExecuteBatch(batch.WithContext(ctx)); err != nil {
		trace.SpanFromContext(ctx).RecordError(err)
	}
}

func selectBook(ctx context.Context, session *gocql.Session) string {
	res := session.
		Query(
			"SELECT id from book WHERE author_last_name = ?",
			"author_last_name 1",
		).
		WithContext(ctx).
		Iter()

	var bookID string
	for res.Scan(&bookID) {
		res.Scan(&bookID)
	}

	res.Close()

	return bookID
}

func updateBook(ctx context.Context, session *gocql.Session, bookID string) {
	if err := session.
		Query(
			"UPDATE book SET title = ? WHERE id = ?",
			"Example Book 1 (republished)", bookID,
		).
		WithContext(ctx).
		Exec(); err != nil {
		trace.SpanFromContext(ctx).RecordError(err)
	}
}

func deleteBook(ctx context.Context, session *gocql.Session, bookID string) {
	if err := session.
		Query("DELETE FROM book WHERE id = ?", bookID).
		WithContext(ctx).
		Exec(); err != nil {
		trace.SpanFromContext(ctx).RecordError(err)
	}
}

func initDB() error {
	cluster := newCassandraCluster("system")
	session, err := cluster.CreateSession()
	if err != nil {
		return err
	}

	stmt := fmt.Sprintf(
		"CREATE KEYSPACE IF NOT EXISTS %s WITH replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 }",
		keyspace,
	)
	if err := session.Query(stmt).Exec(); err != nil {
		return err
	}

	session.Close()

	cluster = newCassandraCluster(keyspace)
	session, err = cluster.CreateSession()
	if err != nil {
		return err
	}
	defer session.Close()

	stmt = "CREATE table IF NOT EXISTS book(id UUID, title text, author_first_name text, author_last_name text, PRIMARY KEY(id))"
	if err = session.Query(stmt).Exec(); err != nil {
		return err
	}

	if err := session.Query("CREATE INDEX IF NOT EXISTS ON book(author_last_name)").Exec(); err != nil {
		return err
	}

	return nil
}

func truncateTable(ctx context.Context, session *gocql.Session) error {
	if err := session.Query("TRUNCATE TABLE book").WithContext(ctx).Exec(); err != nil {
		return err
	}

	return nil
}
