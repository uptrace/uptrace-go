package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/uptrace-go/uptrace"
	otelgomemcache "go.opentelemetry.io/contrib/instrumentation/github.com/bradfitz/gomemcache"
	"go.opentelemetry.io/otel/api/global"
)

const profileTmpl = "profile"

var (
	host = flag.String("host", "127.0.0.1", "memcahe host")
	port = flag.String("port", "11211", "memcache port")
)

var (
	upclient *uptrace.Client
	tracer   = global.Tracer("memcache-tracer")
)

func main() {
	flag.Parse()

	ctx := context.Background()

	upclient = setupUptrace()
	defer upclient.Close()
	defer upclient.ReportPanic(ctx)

	upclient.ReportError(ctx, errors.New("hello from uptrace-go!"))

	// _ = upclient.Tracer("github.com/org/repo")
	mc := otelgomemcache.NewClientWithTracing(
		memcache.New(*host + ":" + *port),
	)

	doMemcacheOperations(ctx, mc)
}

func setupUptrace() *uptrace.Client {
	if os.Getenv("UPTRACE_DSN") == "" {
		log.Printf("using UPTRACE_DSN=%q", os.Getenv("UPTRACE_DSN"))
		panic("UPTRACE_DSN is empty or missing")
	}

	hostname, _ := os.Hostname()
	upclient := uptrace.NewClient(&uptrace.Config{
		// copy your project DSN here or use UPTRACE_DSN enar
		DSN: "",

		Resource: map[string]interface{}{
			"hostname": hostname,
		},
	})

	return upclient
}

func doMemcacheOperations(ctx context.Context, mc *otelgomemcache.Client) {
	err := mc.Add(&memcache.Item{
		Key:   "foo",
		Value: []byte("bar"),
	})
	if err != nil {
		logrus.WithContext(ctx).WithError(err).Error("memcache.Add")
	}

	_, err = mc.Get("foo")
	if err != nil {
		logrus.WithContext(ctx).WithError(err).Error("memcache.Get")

	}

	_, err = mc.Get("hello")
	if err != nil {
		logrus.WithContext(ctx).WithError(err).Error("memcache.Get")
	}

	err = mc.Delete("foo")
	if err != nil {
		logrus.WithContext(ctx).WithError(err).Error("memcache.Delete")
	}
}
