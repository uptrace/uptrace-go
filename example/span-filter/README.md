# Using filters to filter and change spans

This example demonstrates how to use a filter function with uptrace-go:

```go
func newUptraceClient() *uptrace.Client {
	upclient := uptrace.NewClient(&uptrace.Config{
		// copy your project DSN here or use UPTRACE_DSN env var
		DSN: "",

		ServiceName:    "test",
		ServiceVersion: "v1.0.0",
	}, uptrace.WithFilter(spanFilter))

	return upclient
}

func spanFilter(span *spanexp.Span) bool {
	span.Name += " [filter]"

	return true // true keeps the span
}
```

To run this example:

```bash
UPTRACE_DSN="https://<key>@uptrace.dev/<project_id>" go run main.go
```
