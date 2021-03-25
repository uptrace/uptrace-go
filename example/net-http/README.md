# net/http OpenTelemetry instrumentation example

[![PkgGoDev](https://pkg.go.dev/badge/go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp)](https://pkg.go.dev/go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp)

## Quickstart

Install
[otelhttp](https://github.com/open-telemetry/opentelemetry-go-contrib/tree/master/instrumentation/net/http/otelhttp)
instrumentation:

```bash
go get go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp
```

Then Wrap your handlers with `otelhttp.NewHandler`:

```go
func main() {
    var handler http.Handler
    handler = http.HandlerFunc(helloHandler)
    handler = otelhttp.NewHandler(handler, "hello-handler")

    http.Handle("/hello", handler)
}

func helloHandler(w http.ResponseWriter, req *http.Request) { ... }
```

## Example

To run this example:

```bash
UPTRACE_DSN="https://<token>@api.uptrace.dev/<project_id>" go run main.go
```

Then open http://localhost:9999
