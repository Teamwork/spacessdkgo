package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"slices"
)

// checkResponse returns nil if resp.StatusCode matches any of okCodes.
// Otherwise it reads the body, logs the error, and returns a formatted error.
// resp.Request must be non-nil (guaranteed by http.Client and MockRoundTripper).
func (c *Client) checkResponse(ctx context.Context, resp *http.Response, msg string, okCodes ...int) error {
	if slices.Contains(okCodes, resp.StatusCode) {
		return nil
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if c.logger != nil {
		c.logger.LogAttrs(ctx, slog.LevelError, msg,
			slog.String("method", resp.Request.Method),
			slog.String("url", resp.Request.URL.String()),
			slog.Int("status_code", resp.StatusCode),
			slog.String("response_body", string(b)),
		)
	}
	return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(b))
}

// decodeJSON decodes a JSON-encoded value from r into a newly allocated T.
func decodeJSON[T any](r io.Reader) (*T, error) {
	var v T
	if err := json.NewDecoder(r).Decode(&v); err != nil {
		return nil, err
	}
	return &v, nil
}
