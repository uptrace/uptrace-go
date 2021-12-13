package internal

import (
	"fmt"
	"net/url"
)

type DSN struct {
	original string

	ProjectID string
	Token     string

	Scheme string
	Host   string
}

func (dsn *DSN) String() string {
	return dsn.original
}

func (dsn *DSN) OTLPEndpoint() string {
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

	dsn.Scheme = u.Scheme
	if dsn.Scheme == "" {
		return nil, fmt.Errorf("DSN=%q does not have a scheme", dsnStr)
	}

	dsn.Host = u.Host
	if dsn.Host == "api.uptrace.dev" {
		dsn.Host = "uptrace.dev"
	}
	if dsn.Host == "" {
		return nil, fmt.Errorf("DSN=%q does not have a host", dsnStr)
	}

	return &dsn, nil
}
