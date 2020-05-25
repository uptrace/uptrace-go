package internal

import (
	"fmt"
	"net/url"
	"os"
)

type DSN struct {
	ProjectID string
	Token     string

	Scheme string
	Host   string
}

func ParseDSN(dsnStr string) (*DSN, error) {
	if dsnStr == "" {
		dsnStr = os.Getenv("UPTRACE_DSN")
	}

	if dsnStr == "" {
		return nil, fmt.Errorf("DSN is empty or missing")
	}

	u, err := url.Parse(dsnStr)
	if err != nil {
		return nil, fmt.Errorf("can't parse DSN=%q: %s", dsnStr, err)
	}

	var dsn DSN

	if len(u.Path) > 0 {
		dsn.ProjectID = u.Path[1:]
	}
	if dsn.ProjectID == "" {
		return nil, fmt.Errorf("project id is empty or missing")
	}

	if u.User != nil {
		dsn.Token = u.User.Username()
	}
	if dsn.Token == "" {
		return nil, fmt.Errorf("project token is empty or missing")
	}

	dsn.Scheme = u.Scheme
	if dsn.Scheme == "" {
		return nil, fmt.Errorf("scheme is empty or missing")
	}

	dsn.Host = u.Host
	if dsn.Host == "" {
		return nil, fmt.Errorf("host is empty or missing")
	}

	return &dsn, nil
}
