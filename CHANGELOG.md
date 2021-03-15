# Changelog

## v0.9.0

- Update OpenTelemetry to
  [v0.18.0](https://github.com/open-telemetry/opentelemetry-go/blob/main/CHANGELOG.md#0180---2020-03-03)

## v0.6.0

### Added

- Added `Config.ServiceName`, `Config.ServiceVersion`, and `Config.ResourceAttributes`.
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
