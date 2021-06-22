# Changelog

## v0.21.1

- Added back missing resource attributes: `host.name` and `telemetry.sdk.*`.

## v0.21.0

- Updated OpenTelemetry to
  [v1.0.0-RC1](https://github.com/open-telemetry/opentelemetry-go/releases/tag/v1.0.0-RC1).

## v0.20.0

- Updated OpenTelemetry to
  [v0.20.0](https://github.com/open-telemetry/opentelemetry-go/releases/tag/v0.20.0).

## v0.19.0

- Updated OpenTelemetry to
  [v0.19.0](https://github.com/open-telemetry/opentelemetry-go/releases/tag/v0.19.0).
- Changed API and configuration to better indicate that opentelemetry-go can only be configured
  once. Before:

  ```go
  upclient := uptrace.NewClient(&uptrace.Config{...})
  defer upclient.Close()

  fmt.Println(upclient.TraceURL(trace.SpanFromContext(ctx)))
  ```

  Now:

  ```go
  uptrace.ConfigureOpentelemetry(&uptrace.Config{...})
  defer uptrace.Shutdown(ctx)

  fmt.Println(uptrace.TraceURL(trace.SpanFromContext(ctx)))
  ```

- Changed uptrace-go to follow the versioning of opentelemetry-go. For example, uptrace-go v0.19.x
  is compatible with opentelemetry-go v0.19.x.

## v0.9.0

- Updated OpenTelemetry to
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
