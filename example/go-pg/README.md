# go-pg OpenTelemetry instrumentation example

[![PkgGoDev](https://pkg.go.dev/badge/github.com/go-pg/pg/extra/pgotel)](https://pkg.go.dev/github.com/go-pg/pg/extra/pgotel)

## Quickstart

To install [pgotel](https://github.com/go-pg/pg/tree/v10/extra/pgotel) instrumentation:

```bash
go get github.com/go-pg/pg/extra/pgotel/v10
```

Then add OpenTelemetry hook:

```go
db := pg.Connect(&pg.Options{
    Addr:     "postgresql-server:5432",
    User:     "postgres",
    Database: "example",
})

db.AddQueryHook(pgotel.NewTracingHook())
```

## Example

To run this example you need a PostgreSQL server. You can start one with Docker:

```bash
make up
```

Then run the example:

```bash
UPTRACE_DSN="https://<token>@api.uptrace.dev/<project_id>" go run main.go
```
