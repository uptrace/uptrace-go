package internal

import (
	"fmt"
	"net"
	"net/url"
	"strings"
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
	const subdomain = "otlp."

	endpoint := strings.TrimPrefix(dsn.Host, "api.")

	host, _, err := net.SplitHostPort(endpoint)
	if err != nil {
		host = endpoint
	}

	return subdomain + net.JoinHostPort(host, "4317")
}

func ParseDSN(dsnStr string) (*DSN, error) {
	if dsnStr == "" {
		return nil, fmt.Errorf("either Config.DSN or UPTRACE_DSN required")
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
	if dsn.Host == "" {
		return nil, fmt.Errorf("DSN=%q does not have a host", dsnStr)
	}

	return &dsn, nil
}
