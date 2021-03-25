# Zap OpenTelemetry instrumentation example

[![PkgGoDev](https://pkg.go.dev/badge/github.com/uptrace/uptrace-go/extra/otelzap)](https://pkg.go.dev/github.com/uptrace/uptrace-go/extra/otelzap)

Install [otelzap](https://github.com/uptrace/uptrace-go/tree/master/extra/otelzap) instrumentation:

```bash
go get github.com/uptrace/uptrace-go/extra/otelzap
go mod edit -replace go.uber.org/zap=github.com/uptrace/zap@master
```

Then wrap zap logger:

```go
logger, err := zap.NewDevelopment()
if err != nil {
	panic(err)
}
defer logger.Sync()

logger = otelzap.Wrap(logger, otelzap.WithLevel(zap.NewAtomicLevelAt(zap.ErrorLevel)))
```

And use `zap.Ctx` to propagate the active span:

```go
logger.Ctx(ctx).Error("hello from zap",
	zap.Error(errors.New("hello world")),
	zap.String("foo", "bar"))
```

## Example

To run this example:

```bash
UPTRACE_DSN="https://<token>@api.uptrace.dev/<project_id>" go run main.go
```

**Note** that this example requires patching the zap package:

```go
go mod edit -replace go.uber.org/zap=github.com/uptrace/zap@master
```
