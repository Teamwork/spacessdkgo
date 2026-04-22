package client

import (
	"log/slog"
	"net/http"
	"time"
)

// LoggingTransport is an http.RoundTripper that logs each request and response.
type LoggingTransport struct {
	wrapped http.RoundTripper
	logger  *slog.Logger
}

// NewLoggingTransport creates a LoggingTransport wrapping the given transport.
func NewLoggingTransport(wrapped http.RoundTripper, logger *slog.Logger) *LoggingTransport {
	if wrapped == nil {
		wrapped = http.DefaultTransport
	}
	return &LoggingTransport{
		wrapped: wrapped,
		logger:  logger,
	}
}

// RoundTrip executes the request, logs metadata, and returns the response.
func (t *LoggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()

	if t.logger != nil {
		t.logger.LogAttrs(req.Context(), slog.LevelInfo, "http request",
			slog.String("method", req.Method),
			slog.String("url", req.URL.String()),
		)
	}

	resp, err := t.wrapped.RoundTrip(req)
	duration := time.Since(start)

	if err != nil {
		if t.logger != nil {
			t.logger.LogAttrs(req.Context(), slog.LevelError, "http request failed",
				slog.String("method", req.Method),
				slog.String("url", req.URL.String()),
				slog.Any("error", err),
				slog.Duration("duration", duration),
			)
		}
		return nil, err
	}

	if t.logger != nil {
		t.logger.LogAttrs(req.Context(), slog.LevelInfo, "http response",
			slog.String("method", req.Method),
			slog.String("url", req.URL.String()),
			slog.Int("status_code", resp.StatusCode),
			slog.Duration("duration", duration),
		)
	}

	return resp, nil
}

// NewLoggingClientWithLogger creates an *http.Client with a LoggingTransport.
func NewLoggingClientWithLogger(logger *slog.Logger) *http.Client {
	return &http.Client{
		Transport: NewLoggingTransport(http.DefaultTransport, logger),
	}
}
