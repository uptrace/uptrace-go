package uptrace

import (
	"fmt"
	"net"
	"net/url"
)

type DSN struct {
	original string

	Scheme   string
	Host     string
	HTTPPort string
	GRPCPort string
	Token    string
}

func (dsn *DSN) String() string {
	return dsn.original
}

func (dsn *DSN) SiteURL() string {
	if dsn.Host == "uptrace.dev" {
		return "https://app.uptrace.dev"
	}
	return dsn.Scheme + "://" + joinHostPort(dsn.Host, dsn.HTTPPort)
}

func (dsn *DSN) OTLPEndpoint() string {
	if dsn.Host == "uptrace.dev" {
		return "otlp.uptrace.dev:4317"
	}
	return joinHostPort(dsn.Host, dsn.GRPCPort)
}

func ParseDSN(dsnStr string) (*DSN, error) {
	if dsnStr == "" {
		return nil, fmt.Errorf("DSN is empty (use WithDSN or UPTRACE_DSN env var)")
	}

	u, err := url.Parse(dsnStr)
	if err != nil {
		return nil, fmt.Errorf("can't parse DSN=%q: %s", dsnStr, err)
	}

	if u.Scheme == "" {
		return nil, fmt.Errorf("DSN=%q does not have a scheme", dsnStr)
	}
	if u.Host == "" {
		return nil, fmt.Errorf("DSN=%q does not have a host", dsnStr)
	}
	if u.User == nil {
		return nil, fmt.Errorf("DSN=%q does not have a token", dsnStr)
	}

	dsn := DSN{
		original: dsnStr,
		Scheme:   u.Scheme,
		Host:     u.Host,
		Token:    u.User.Username(),
	}

	if host, port, err := net.SplitHostPort(u.Host); err == nil {
		dsn.Host = host
		dsn.HTTPPort = port
	}

	if dsn.Host == "api.uptrace.dev" {
		dsn.Host = "uptrace.dev"
	}

	query := u.Query()
	if grpc := query.Get("grpc"); grpc != "" {
		dsn.GRPCPort = grpc
	}

	if dsn.GRPCPort == "" {
		if dsn.HTTPPort != "" {
			dsn.GRPCPort = dsn.HTTPPort
			if dsn.HTTPPort == "14317" {
				dsn.HTTPPort = "14318"
			}
		} else {
			dsn.GRPCPort = "4317"
		}
	}

	return &dsn, nil
}

func joinHostPort(host, port string) string {
	if port == "" {
		return host
	}
	return net.JoinHostPort(host, port)
}
