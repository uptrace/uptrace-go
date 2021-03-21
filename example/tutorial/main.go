package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/uptrace/uptrace-go/uptrace"
)

var tracer = otel.Tracer("app_or_package_name")

func main() {
	ctx := context.Background()

	uptrace.ConfigureOpentelemetry(&uptrace.Config{
		// copy your project DSN here or use UPTRACE_DSN env var
		DSN: "",
	})
	defer uptrace.Shutdown(ctx)

	ctx, span := tracer.Start(ctx, "fetchCountry")
	defer span.End()

	countryInfo, err := fetchCountryInfo(ctx)
	if err != nil {
		span.RecordError(err)
		return
	}

	countryCode, countryName, err := parseCountryInfo(ctx, countryInfo)
	if err != nil {
		span.RecordError(err)
		return
	}

	span.SetAttributes(
		attribute.String("country.code", countryCode),
		attribute.String("country.name", countryName),
	)

	fmt.Println("trace URL", uptrace.TraceURL(span))
}

func fetchCountryInfo(ctx context.Context) (string, error) {
	ctx, span := tracer.Start(ctx, "fetchCountryInfo")
	defer span.End()

	resp, err := http.Get("https://ip2c.org/self")
	if err != nil {
		span.RecordError(err)
		return "", err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		span.RecordError(err)
		return "", err
	}

	span.SetAttributes(
		attribute.String("ip", "self"),
		attribute.Int("resp_len", len(b)),
	)

	return string(b), nil
}

func parseCountryInfo(ctx context.Context, s string) (code, country string, _ error) {
	ctx, span := tracer.Start(ctx, "parseCountryInfo")
	defer span.End()

	parts := strings.Split(s, ";")
	if len(parts) < 4 {
		err := fmt.Errorf("ip2c: can't parse response: %q", s)
		span.RecordError(err)
		return "", "", err
	}
	return parts[1], parts[3], nil
}
