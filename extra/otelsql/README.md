[![PkgGoDev](https://pkg.go.dev/badge/github.com/uptrace/uptrace-go/extra/otelsql)](https://pkg.go.dev/github.com/uptrace/uptrace-go/extra/otelsql)

# OpenTelemetry Go instrumentation for database/sql package

This OpenTelemetry instrumentation records database queries (including `Tx` and `Stmt` queries) and
`DBStats` metrics.

## Installation

```shell
go get github.com/uptrace/uptrace-go/extra/otelsql
```

## Usage

To instrument database/sql client, you need to connect to a database using the API provided by this
package:

- `sql.Open(driverName, dsn)` becomes `otelsql.Open(driverName, dsn)`.
- `sql.OpenDB(connector)` becomes `otelsql.OpenDB(connector)`.

```go
import (
    "github.com/uptrace/uptrace-go/extra/otelsql"
    semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

db, err := otelsql.Open("sqlite", "file::memory:?cache=shared",
	otelsql.WithAttributes(semconv.DBSystemSqlite))
if err != nil {
    panic(err)
}
```

See [example](/example/) for details.
