package uptrace

import (
	"fmt"
	"net"
	"net/url"
)

type DSN struct {
	original string

	Scheme string
	Host   string

	ProjectID string
	Token     string
}

func (dsn *DSN) String() string {
	return dsn.original
}

func (dsn *DSN) AppAddr() string {
	if dsn.Host == "uptrace.dev" {
		return "https://app.uptrace.dev"
	}
	host, _, err := net.SplitHostPort(dsn.Host)
	if err != nil {
		return dsn.Host
	}
	return dsn.Scheme + "://" + net.JoinHostPort(host, "14318")
}

func (dsn *DSN) OTLPHost() string {
	if dsn.Host == "uptrace.dev" {
		return "otlp.uptrace.dev:4317"
	}
	return dsn.Host
}

func ParseDSN(dsnStr string) (*DSN, error) {
	if dsnStr == "" {
		return nil, fmt.Errorf("DSN is empty (use WithDSN or UPTRACE_DSN env var)")
	}

	u, err := url.Parse(dsnStr)
	if err != nil {
		return nil, fmt.Errorf("can't parse DSN=%q: %s", dsnStr, err)
	}

	dsn := DSN{
		original: dsnStr,
	}

	dsn.Scheme = u.Scheme
	if dsn.Scheme == "" {
		return nil, fmt.Errorf("DSN=%q does not have a scheme", dsnStr)
	}

	dsn.Host = u.Host
	if dsn.Host == "" {
		return nil, fmt.Errorf("DSN=%q does not have a host", dsnStr)
	}
	if dsn.Host == "api.uptrace.dev" {
		dsn.Host = "uptrace.dev"
	}

	if len(u.Path) > 0 {
		dsn.ProjectID = u.Path[1:]
	}
	if dsn.ProjectID == "" {
		return nil, fmt.Errorf("DSN=%q does not have a project id", dsnStr)
	}

	if u.User != nil {
		dsn.Token = u.User.Username()
	}
	if dsn.Token == "" {
		return nil, fmt.Errorf("DSN=%q does not have a token", dsnStr)
	}

	return &dsn, nil
}
