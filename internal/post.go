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

type statusCodeError struct {
	code int
	msg  string
}

func (e statusCodeError) Error() string {
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

func PostWithRetry(
	ctx context.Context,
	client *http.Client,
	endpoint, token string,
	data []byte,
) error {
	resp, err := postWithRetry(ctx, client, endpoint, token, data)
	if err != nil {
		return err
	}

	defer func() {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		msg := decodeErrorMessage(resp.Body)
		return statusCodeError{
			code: resp.StatusCode,
			msg:  msg,
		}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return statusCodeError{
			code: resp.StatusCode,
		}
	}

	return nil
}

func postWithRetry(
	ctx context.Context,
	client *http.Client,
	endpoint, token string,
	data []byte,
) (resp *http.Response, lastErr error) {
	for attempt := 0; attempt < 3; attempt++ {
		if err := Backoff(ctx, attempt, time.Second, 3*time.Second); err != nil {
			return nil, err
		}

		req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(data))
		if err != nil {
			return nil, err
		}

		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/msgpack")
		req.Header.Set("Content-Encoding", "zstd")

		resp, lastErr = client.Do(req)
		if lastErr == nil {
			return resp, nil
		}
	}
	return nil, lastErr
}
