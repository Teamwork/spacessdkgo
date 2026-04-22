package client

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

// MiddlewareFunc is a function that wraps an HTTP request handler.
type MiddlewareFunc func(ctx context.Context, req *http.Request, next RequestHandler) (*http.Response, error)

// RequestHandler is the function signature for the next handler in the chain.
type RequestHandler func(ctx context.Context, req *http.Request) (*http.Response, error)

// LoggingMiddleware logs method, URL, status code, and duration for each request.
func LoggingMiddleware(logger *slog.Logger) MiddlewareFunc {
	return func(ctx context.Context, req *http.Request, next RequestHandler) (*http.Response, error) {
		start := time.Now()
		logger.InfoContext(ctx, "outgoing request",
			slog.String("method", req.Method),
			slog.String("url", req.URL.String()),
		)

		resp, err := next(ctx, req)
		duration := time.Since(start)

		if err != nil {
			logger.ErrorContext(ctx, "request error",
				slog.String("method", req.Method),
				slog.String("url", req.URL.String()),
				slog.Any("error", err),
				slog.Duration("duration", duration),
			)
			return nil, err
		}

		logger.InfoContext(ctx, "request complete",
			slog.String("method", req.Method),
			slog.String("url", req.URL.String()),
			slog.Int("status_code", resp.StatusCode),
			slog.Duration("duration", duration),
		)
		return resp, nil
	}
}

// RetryMiddleware retries failed requests up to maxRetries times with a delay between attempts.
func RetryMiddleware(maxRetries int, retryDelay time.Duration) MiddlewareFunc {
	return func(ctx context.Context, req *http.Request, next RequestHandler) (*http.Response, error) {
		var (
			resp *http.Response
			err  error
		)
		for attempt := 0; attempt <= maxRetries; attempt++ {
			if attempt > 0 {
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				case <-time.After(retryDelay):
				}
			}

			// Clone the request body for each attempt since it may have been read.
			cloned := req.Clone(ctx)
			resp, err = next(ctx, cloned)
			if err == nil {
				return resp, nil
			}
		}
		return resp, err
	}
}

// AuthMiddleware sets the Authorization header to "Bearer <token>" on every request.
func AuthMiddleware(token string) MiddlewareFunc {
	return func(ctx context.Context, req *http.Request, next RequestHandler) (*http.Response, error) {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		return next(ctx, req)
	}
}

// UserAgentMiddleware sets the User-Agent header on every request.
func UserAgentMiddleware(userAgent string) MiddlewareFunc {
	return func(ctx context.Context, req *http.Request, next RequestHandler) (*http.Response, error) {
		req.Header.Set("User-Agent", userAgent)
		return next(ctx, req)
	}
}

// RateLimitMiddleware limits outgoing requests to the given rate (requests per second).
func RateLimitMiddleware(requestsPerSecond float64) MiddlewareFunc {
	var (
		mu        sync.Mutex
		tokens    = requestsPerSecond
		lastCheck = time.Now()
	)
	return func(ctx context.Context, req *http.Request, next RequestHandler) (*http.Response, error) {
		mu.Lock()
		now := time.Now()
		elapsed := now.Sub(lastCheck).Seconds()
		lastCheck = now
		tokens += elapsed * requestsPerSecond
		if tokens > requestsPerSecond {
			tokens = requestsPerSecond
		}
		if tokens < 1 {
			wait := time.Duration((1-tokens)/requestsPerSecond*1000) * time.Millisecond
			mu.Unlock()
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(wait):
			}
			mu.Lock()
			tokens = 0
		} else {
			tokens--
		}
		mu.Unlock()
		return next(ctx, req)
	}
}

// RequestIDMiddleware adds a unique X-Request-ID header to every request.
func RequestIDMiddleware() MiddlewareFunc {
	return func(ctx context.Context, req *http.Request, next RequestHandler) (*http.Response, error) {
		req.Header.Set("X-Request-ID", generateRequestID())
		return next(ctx, req)
	}
}

// generateRequestID returns a random hex request ID.
func generateRequestID() string {
	return fmt.Sprintf("%016x", rand.Int63()) //nolint:gosec
}

// TimeoutMiddleware wraps the request context with a deadline.
func TimeoutMiddleware(timeout time.Duration) MiddlewareFunc {
	return func(ctx context.Context, req *http.Request, next RequestHandler) (*http.Response, error) {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		return next(ctx, req.WithContext(ctx))
	}
}

// HeaderMiddleware adds a set of static headers to every request.
func HeaderMiddleware(headers map[string]string) MiddlewareFunc {
	return func(ctx context.Context, req *http.Request, next RequestHandler) (*http.Response, error) {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
		return next(ctx, req)
	}
}

// ConditionalMiddleware applies mw only when condition returns true.
func ConditionalMiddleware(condition func(*http.Request) bool, mw MiddlewareFunc) MiddlewareFunc {
	return func(ctx context.Context, req *http.Request, next RequestHandler) (*http.Response, error) {
		if condition(req) {
			return mw(ctx, req, next)
		}
		return next(ctx, req)
	}
}
