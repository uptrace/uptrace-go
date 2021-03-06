package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var HTTPClient = &http.Client{
	Timeout: 10 * time.Second,
}

type StatusCodeError struct {
	code int
	msg  string
}

func NewStatusCodeError(code int, msg string) *StatusCodeError {
	return &StatusCodeError{
		code: code,
		msg:  msg,
	}
}

func (e *StatusCodeError) Code() int {
	return e.code
}

// Temporary reports whether the error is temporary and safe to retry.
func (e *StatusCodeError) Temporary() bool {
	return e.code >= 500
}

func (e *StatusCodeError) Error() string {
	if e.msg != "" {
		return fmt.Sprintf("status=%d: %s", e.code, e.msg)
	}
	return "got status=" + strconv.Itoa(e.code) + ", wanted 200 OK"
}

func decodeErrorMessage(r io.Reader) string {
	m := make(map[string]interface{})
	if err := json.NewDecoder(r).Decode(&m); err != nil {
		return err.Error()
	}
	msg, _ := m["message"].(string)
	return msg
}

type SimpleClient struct {
	Client     *http.Client
	Token      string
	MaxRetries int
}

func (c *SimpleClient) Post(
	ctx context.Context, endpoint string, data []byte,
) error {
	resp, err := c.postWithRetry(ctx, endpoint, data)
	if err != nil {
		return err
	}

	defer func() {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		msg := decodeErrorMessage(resp.Body)
		return NewStatusCodeError(resp.StatusCode, msg)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return NewStatusCodeError(resp.StatusCode, "")
	}

	return nil
}

func (c *SimpleClient) postWithRetry(
	ctx context.Context, endpoint string, data []byte,
) (resp *http.Response, lastErr error) {
	for attempt := 0; attempt <= c.MaxRetries; attempt++ {
		if err := Backoff(ctx, attempt, time.Second, 3*time.Second); err != nil {
			return nil, err
		}

		resp, lastErr = c.post(ctx, endpoint, data)
		if lastErr != nil || resp.StatusCode >= 500 {
			continue
		}

		return resp, nil
	}
	return resp, lastErr
}

func (c *SimpleClient) post(
	ctx context.Context, endpoint string, data []byte,
) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/msgpack")
	req.Header.Set("Content-Encoding", "zstd")

	return c.Client.Do(req)
}
