package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var tracer = otel.Tracer("app_or_package_name")

func main() {
	ctx := context.Background()

	upclient := uptrace.NewClient(&uptrace.Config{
		// copy your project DSN here or use UPTRACE_DSN env var
		DSN: "",
	})
	defer upclient.Close()
	defer upclient.ReportPanic(ctx)

	r := mux.NewRouter()
	r.Use(otelmux.Middleware("service-name"))
	r.HandleFunc("/profiles/{username}", userProfileHandler)

	log.Fatal(http.ListenAndServe(":9999", r))
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

	span.SetAttributes(attribute.String("username", username))

	if username == "admin" {
		return "Joe", nil
	}

	return "", fmt.Errorf("username=%s not found", username)
}
