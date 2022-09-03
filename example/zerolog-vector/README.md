# Zerolog and Vector example

This example demonstrates how to collect [Zerolog](https://github.com/rs/zerolog) logs with
[Vector](https://vector.dev/).

Because Zerolog does not support
[context](https://uptrace.dev/opentelemetry/go-tracing.html#context), we have to propagate
`trace_id` and `span_id` in the log message:

```shell
logger.Error().
	// trace_id and span_id are needed to properly link the log message with the span.
	Str("trace_id", child1.SpanContext().TraceID().String()).
	Str("span_id", child1.SpanContext().SpanID().String()).
	Str("foo", "bar").
	Msg("message from zerolog")
```

That is not very convenient, so we recommend to use
[Zap](https://github.com/uptrace/opentelemetry-go-extra/tree/main/otelzap) or
[logrus](https://github.com/uptrace/opentelemetry-go-extra/tree/main/otellogrus) instead.

## Running this example

**Step 1**. Replace `headers.uptrace-dsn` in [vector.toml](vector.toml) with your Uptrace DSN.

**Step 2**. Start Vector:

```shell
vector --config vector.toml
```

**Step 3**. Run the example:

```shell
UPTRACE_DSN="https://<token>@uptrace.dev/<project_id>" go run .
```

Click the link in your terminal to open Uptrace. See
[Structured Logging](https://uptrace.dev/opentelemetry/structured-logging.html) for details.
