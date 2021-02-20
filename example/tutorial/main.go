package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/label"
)

var tracer = otel.Tracer("app_or_package_name")

func main() {
	upclient := uptrace.NewClient(&uptrace.Config{
		// copy your project DSN here or use UPTRACE_DSN env var
		DSN: "",
	})
	defer upclient.Close()

	ctx := context.Background()

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
		label.String("country.code", countryCode),
		label.String("country.name", countryName),
	)

	fmt.Println("trace URL", upclient.TraceURL(span))
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
		label.String("ip", "self"),
		label.Int("resp_len", len(b)),
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
