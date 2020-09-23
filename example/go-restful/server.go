package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	restful "github.com/emicklei/go-restful/v3"
	"github.com/uptrace/uptrace-go/uptrace"
	otelrestful "go.opentelemetry.io/contrib/instrumentation/github.com/emicklei/go-restful"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/label"
)

var (
	upclient *uptrace.Client
	tracer   = global.Tracer("restful-tracer")
)

func main() {
	ctx := context.Background()

	upclient = setupUptrace()
	defer upclient.Close()
	defer upclient.ReportPanic(ctx)

	upclient.ReportError(ctx, errors.New("hello from uptrace-go!"))

	filter := otelrestful.OTelFilter("service-name")
	restful.DefaultContainer.Filter(filter)

	ws := &restful.WebService{}
	ws.Route(ws.GET("/profiles/{username}").To(userProfileHandler))
	restful.Add(ws)
	log.Fatal(http.ListenAndServe(":9999", nil))
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

func userProfileHandler(req *restful.Request, resp *restful.Response) {
	ctx := req.Request.Context()

	username := req.PathParameter("username")
	name, err := selectUser(ctx, username)
	if err != nil {
		resp.WriteError(404, err)
		return
	}

	html := fmt.Sprintf(`<html><h1>Hello %s %s </h1></html>`+"\n", username, name)
	resp.Write([]byte(html))
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
