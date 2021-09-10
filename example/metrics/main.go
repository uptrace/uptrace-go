package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"sync/atomic"
	"syscall"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"

	"github.com/uptrace/uptrace-go/uptrace"
)

var meter = metric.Must(global.Meter("app_or_package_name"))

func main() {
	ctx := context.Background()

	// Configure OpenTelemetry with sensible defaults.
	uptrace.ConfigureOpentelemetry(
		// copy your project DSN here or use UPTRACE_DSN env var
		// uptrace.WithDSN("https://<key>@api.uptrace.dev/<project_id>"),

		uptrace.WithServiceName("myservice"),
		uptrace.WithServiceVersion("1.0.0"),
	)
	// Send buffered spans and free resources.
	defer uptrace.Shutdown(ctx)

	// Synchronous instruments.
	go counter(ctx)
	go counterWithLabels(ctx)
	go upDownCounter(ctx)
	go histogram(ctx)

	// Asynchronous instruments.
	go counterObserver(ctx)
	go upDownCounterObserver(ctx)
	go gaugeObserver(ctx)

	fmt.Println("reporting measurements to Uptrace... (press Ctrl+C to stop)")

	ch := make(chan os.Signal, 3)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	<-ch
}

// counter demonstrates how to measure non-decreasing numbers, for example,
// number of requests or connections.
func counter(ctx context.Context) {
	counter := meter.NewInt64Counter("app_or_package_name.component1.requests",
		metric.WithDescription("Number of requests"),
	)

	for {
		counter.Add(ctx, 1)
		time.Sleep(time.Millisecond)
	}
}

// counterWithLabels demonstrates how to add different labels ("hits" and "misses")
// to measurements. Using this simple trick, you can get number of hits, misses,
// sum = hits + misses, and hit_rate = hits / (hits + misses).
func counterWithLabels(ctx context.Context) {
	counter := meter.NewInt64Counter("app_or_package_name.component1.cache",
		metric.WithDescription("Cache hits and misses"),
	)
	// Bind the counter to some labels.
	hits := counter.Bind(attribute.String("type", "hits"))
	misses := counter.Bind(attribute.String("type", "misses"))

	for {
		if rand.Float64() < 0.3 {
			misses.Add(ctx, 1)
		} else {
			hits.Add(ctx, 1)
		}

		time.Sleep(time.Millisecond)
	}
}

// upDownCounter demonstrates how to measure numbers that can go up and down, for example,
// number of goroutines or customers.
//
// See upDownCounterObserver for a better example how to measure number of goroutines.
func upDownCounter(ctx context.Context) {
	counter := meter.NewInt64UpDownCounter("app_or_package_name.component1.goroutines",
		metric.WithDescription("Number of goroutines"),
	)

	for {
		counter.Add(ctx, int64(runtime.NumGoroutine()))

		time.Sleep(time.Second)
	}
}

// histogram demonstrates how to record a distribution of individual values, for example,
// request or query timings. With this instrument you get total number of records,
// avg/min/max values, and heatmaps/percentiles.
func histogram(ctx context.Context) {
	durRecorder := meter.NewInt64Histogram("app_or_package_name.component1.request_duration",
		metric.WithUnit("microseconds"),
		metric.WithDescription("Duration of requests"),
	)

	for {
		dur := time.Duration(rand.NormFloat64()*10000+100000) * time.Microsecond
		durRecorder.Record(ctx, dur.Microseconds())

		time.Sleep(time.Millisecond)
	}
}

// counterObserver demonstrates how to measure monotonic (non-decreasing) numbers,
// for example, number of requests or connections.
func counterObserver(ctx context.Context) {
	// stats is our data source updated by some library.
	var stats struct {
		Hits   int64 // atomic
		Misses int64 // atomic
	}

	var hitsCounter, missesCounter metric.Int64CounterObserver

	batchObserver := meter.NewBatchObserver(
		// SDK periodically calls this function to grab results.
		func(ctx context.Context, result metric.BatchObserverResult) {
			result.Observe(nil,
				hitsCounter.Observation(atomic.LoadInt64(&stats.Hits)),
				missesCounter.Observation(atomic.LoadInt64(&stats.Misses)),
			)
		})

	hitsCounter = batchObserver.NewInt64CounterObserver("app_or_package_name.component2.cache_hits")
	missesCounter = batchObserver.NewInt64CounterObserver("app_or_package_name.component2.cache_misses")

	for {
		if rand.Float64() < 0.3 {
			atomic.AddInt64(&stats.Misses, 1)
		} else {
			atomic.AddInt64(&stats.Hits, 1)
		}

		time.Sleep(time.Millisecond)
	}
}

// upDownCounterObserver demonstrates how to measure numbers that can go up and down,
// for example, number of goroutines or customers.
func upDownCounterObserver(ctx context.Context) {
	_ = meter.NewInt64UpDownCounterObserver("app_or_package_name.component2.goroutines",
		func(ctx context.Context, result metric.Int64ObserverResult) {
			num := runtime.NumGoroutine()
			result.Observe(int64(num))
		},
		metric.WithDescription("Number of goroutines"),
	)
}

// gaugeObserver demonstrates how to measure numbers that can go up and down,
// for example, number of goroutines or customers.
func gaugeObserver(ctx context.Context) {
	_ = meter.NewInt64GaugeObserver("app_or_package_name.component2.goroutines2",
		func(ctx context.Context, result metric.Int64ObserverResult) {
			num := runtime.NumGoroutine()
			result.Observe(int64(num))
		},
		metric.WithDescription("Number of goroutines"),
	)
}
