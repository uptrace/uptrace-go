# Changelog

## v1.6.0

- Updated OpenTelemetry to
  [v1.6.0](https://github.com/open-telemetry/opentelemetry-go/releases/tag/v1.6.0).

## v1.5.0

- Updated OpenTelemetry to
  [v1.5.0](https://github.com/open-telemetry/opentelemetry-go/releases/tag/v1.5.0).

## v1.4.0

- Updated OpenTelemetry to
  [v1.4.0](https://github.com/open-telemetry/opentelemetry-go/releases/tag/v1.4.0).

## v1.3.0

- Updated OpenTelemetry to
  [v1.3.0](https://github.com/open-telemetry/opentelemetry-go/releases/tag/v1.3.0).

- Added `WithResourceDetectors` to configure resource detectors, for example:

```go
import (
	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/contrib/detectors/aws/ec2"
)

uptrace.ConfigureOpentelemetry(
	uptrace.WithResourceDetectors(ec2.NewResourceDetector()),
)
```

See [documentation](https://docs.uptrace.dev/guide/go.html#resource-detectors) for details.

## v1.2.0

- Updated OpenTelemetry to
  [v1.2.0](https://github.com/open-telemetry/opentelemetry-go/releases/tag/v1.2.0).

## v1.1.0

- Updated OpenTelemetry to
  [v1.1.0](https://github.com/open-telemetry/opentelemetry-go/releases/tag/v1.1.0).
- Moved instrumentations to
  [opentelemetry-go-extra](https://github.com/uptrace/opentelemetry-go-extra)

## v1.0.5

- otelzap: added support for sugared loggers.

## v1.0.4

- Added [otelgorm](/extra/otelgorm/) instrumentation for GORM.
- Changed otelsql to not set `error` status on sql.ErrNoRows errors.

## v1.0.3

- Updated OpenTelemetry to
  [v1.0.1](https://github.com/open-telemetry/opentelemetry-go/releases/tag/v1.0.1).
- Added [otelsql](/extra/otelsql/) instrumentation to instrument database/sql client. The
  instrumentation records processed queries and reports `sql.DBStats` metrics.
- Changed [otelzap](/extra/otelzap/) instrumentation to work with a standard unpatched version of
  Zap. The logging API is compatible, but you now have to wrap a `zap.Logger` with a
  `otelzap.Logger` to add OpenTelemetry instrumentation.

## v1.0.2

- Updated OpenTelemetry to
  [v1.0.0](https://github.com/open-telemetry/opentelemetry-go/releases/tag/v1.0.0).

## v1.0.1

- Updated OpenTelemetry to
  [v1.0.0-RC3](https://github.com/open-telemetry/opentelemetry-go/releases/tag/v1.0.0-RC3).

## v1.0.0

- No changes. The purpose of this release is to avoid confusion with Go installing v0.21.1 by
  default.

## v1.0.0-RC3

- Fully switched to using OpenTelemetry Protocol (OTLP) for exporting spans and metrics. This is
  fully backwards compatible and should not cause any disruptive changes.

## v1.0.0-RC2

- Updated OpenTelemetry to
  [v1.0.0-RC2](https://github.com/open-telemetry/opentelemetry-go/releases/tag/v1.0.0-RC2).
- Changed configuration to use options instead of a single `Config` struct. All the previous
  configuration options are fully supported.

  There are 3 types of options:

  - [Option](https://pkg.go.dev/github.com/uptrace/uptrace-go@v1.0.0-RC2/uptrace#Option) for common
    options that configure tracing and metrics.
  - [TracingOption](https://pkg.go.dev/github.com/uptrace/uptrace-go@v1.0.0-RC2/uptrace#TracingOption)
    for options specific to tracing.
  - [MetricsOption](https://pkg.go.dev/github.com/uptrace/uptrace-go@v1.0.0-RC2/uptrace#MetricsOption)
    for options specific to metrics.

  For example, to configure tracing and metrics:

  ```go
  uptrace.ConfigureOpentelemetry(
      uptrace.WithDSN("https://<token>@api.uptrace.dev/<project_id>"),
      uptrace.WithServiceName("myservice"),
      uptrace.WithServiceVersion("1.0.0"),
  )
  ```

  To configure only tracing, use `WithMetricsDisabled` option:

  ```go
  uptrace.ConfigureOpentelemetry(
      uptrace.WithMetricsDisabled(),

      uptrace.WithDSN("https://<token>@api.uptrace.dev/<project_id>"),
      uptrace.WithServiceName("myservice"),
      uptrace.WithServiceVersion("1.0.0"),
  )
  ```

- Added support for OpenTelemetry Metrics using standard OTLP exporter.
- Enabled metrics by default. `WithMetricsDisabled` option can be used to disable metrics.

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
