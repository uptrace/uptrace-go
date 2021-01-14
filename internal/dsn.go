package internal

import (
	"fmt"
	"net/url"
)

type DSN struct {
	ProjectID string
	Token     string

	Scheme string
	Host   string
}

func ParseDSN(dsnStr string) (*DSN, error) {
	if dsnStr == "" {
		return nil, fmt.Errorf("uptrace: either Config.DSN or UPTRACE_DSN env var is required")
	}

	u, err := url.Parse(dsnStr)
	if err != nil {
		return nil, fmt.Errorf("uptrace: can't parse DSN=%q: %s", dsnStr, err)
	}

	var dsn DSN

	if len(u.Path) > 0 {
		dsn.ProjectID = u.Path[1:]
	}
	if dsn.ProjectID == "" {
		return nil, fmt.Errorf("uptrace: DSN does not have project id (DSN=%q)", dsnStr)
	}

	if u.User != nil {
		dsn.Token = u.User.Username()
	}
	if dsn.Token == "" {
		return nil, fmt.Errorf("uptrace: DSN does not have project token (DSN=%q)", dsnStr)
	}

	dsn.Scheme = u.Scheme
	if dsn.Scheme == "" {
		return nil, fmt.Errorf("uptrace: DSN does not have scheme (DSN=%q)", dsnStr)
	}

	dsn.Host = u.Host
	if dsn.Host == "" {
		return nil, fmt.Errorf("uptrace: DSN does not have host (DSN=%q)", dsnStr)
	}

	return &dsn, nil
}
