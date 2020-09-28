package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/label"
)

var (
	upclient *uptrace.Client
	tracer   = global.Tracer("mux-tracer")
)

func main() {
	ctx := context.Background()

	upclient = setupUptrace()
	defer upclient.Close()
	defer upclient.ReportPanic(ctx)

	upclient.ReportError(ctx, errors.New("hello from uptrace-go!"))

	r := mux.NewRouter()
	r.Use(otelmux.Middleware("service-name"))
	r.HandleFunc("/profiles/{username}", userProfileHandler)

	log.Fatal(http.ListenAndServe(":9999", r))
}

func setupUptrace() *uptrace.Client {
	if os.Getenv("UPTRACE_DSN") == "" {
		panic("UPTRACE_DSN is empty or missing")
	}

	hostname, _ := os.Hostname()
	upclient := uptrace.NewClient(&uptrace.Config{
		// copy your project DSN here or use UPTRACE_DSN env var
		DSN: "",

		Resource: map[string]interface{}{
			"hostname": hostname,
		},
	})

	return upclient
}

func userProfileHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	username := mux.Vars(req)["username"]
	name, err := selectUser(ctx, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, `<html><h1>Hello %s %s </h1></html>`+"\n", username, name)
}

func selectUser(ctx context.Context, username string) (string, error) {
	_, span := tracer.Start(ctx, "selectUser")
	defer span.End()

	span.SetAttributes(label.String("username", username))

	if username == "admin" {
		return "Joe", nil
	}

	return "", fmt.Errorf("username=%s not found", username)
}
