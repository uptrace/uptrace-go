# Changelog

## v0.6.0

### Added

- Added span filter and an [example](/example/span-filter/).

## v0.5.0

### Added

- Added default `Config.Resource`.
- Added `Config.TextMapPropagator` with sensible default value.

## v0.4.0

### Added

- `Config.PrettyPrint` for debugging.

### Changed

- `Config.Resource` type is `*resource.Resource`. Before:

```go
Resource: map[string]interface{}{
    "host.name": hostname,
},
```

After:

```go
import "go.opentelemetry.io/otel/sdk/resource"

Resource: resource.New(
    label.String("host.name", hostname),
),
```
