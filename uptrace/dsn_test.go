package uptrace_test

import (
	"fmt"
	"testing"

	"github.com/uptrace/uptrace-go/uptrace"

	"github.com/stretchr/testify/require"
)

func TestParseDSN(t *testing.T) {
	type Test struct {
		dsn     string
		grpc    string
		http    string
		siteURL string
	}

	tests := []Test{
		{
			"https://token@uptrace.dev/1",
			"api.uptrace.dev:4317",
			"api.uptrace.dev:443",
			"https://app.uptrace.dev",
		},
		{
			"https://token@api.uptrace.dev/1",
			"api.uptrace.dev:4317",
			"api.uptrace.dev:443",
			"https://app.uptrace.dev",
		},
		{
			"https://token@demo.uptrace.dev/1?grpc=4317",
			"demo.uptrace.dev:4317",
			"demo.uptrace.dev:443",
			"https://demo.uptrace.dev:443",
		},
		{
			"https://token@localhost:1234/1",
			"localhost:1234",
			"localhost:1234",
			"https://localhost:1234",
		},
		{
			"http://token@localhost:14317/project_id",
			"localhost:14317",
			"localhost:14318",
			"http://localhost:14318",
		},
		{
			"https://AQDan_E_EPe3QAF9fMP0PiVr5UWOu4q5@demo-api.uptrace.dev:4317/1",
			"demo-api.uptrace.dev:4317",
			"demo-api.uptrace.dev:4317",
			"https://demo-api.uptrace.dev:4317",
		},
		{
			"http://Qcn7rcwWO_w0ePo7WmeUtw@localhost:14318?grpc=14317",
			"localhost:14317",
			"localhost:14318",
			"http://localhost:14318",
		},
	}
	for i, test := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			dsn, err := uptrace.ParseDSN(test.dsn)
			require.NoError(t, err)
			require.Equal(t, test.grpc, dsn.OTLPGrpcEndpoint())
			require.Equal(t, test.http, dsn.OTLPHttpEndpoint())
			require.Equal(t, test.siteURL, dsn.SiteURL())
		})
	}
}
