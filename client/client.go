package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

// Client is the top-level SDK client holding configuration and all resource services.
type Client struct {
	baseURL    string
	apiKey     string
	logLevel   slog.Level
	logger     *slog.Logger
	httpClient *http.Client
	middleware []MiddlewareFunc

	Pages      *PageService
	Spaces     *SpaceService
	Comments   *CommentService
	Tags       *TagService
	Categories *CategoryService
	Search     *SearchService
}

// Option is a functional option for configuring a Client.
type Option func(*Client)

// WithAPIKey sets the API key used for Bearer authentication.
func WithAPIKey(apiKey string) Option {
	return func(c *Client) { c.apiKey = apiKey }
}

// WithHTTPClient overrides the default HTTP client.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) { c.httpClient = httpClient }
}

// WithLogLevel sets the minimum log level for the client logger.
func WithLogLevel(level slog.Level) Option {
	return func(c *Client) { c.logLevel = level }
}

// WithLogger sets the structured logger used throughout the client.
func WithLogger(logger *slog.Logger) Option {
	return func(c *Client) { c.logger = logger }
}

// WithMiddleware appends a MiddlewareFunc to the client's middleware chain.
func WithMiddleware(mw MiddlewareFunc) Option {
	return func(c *Client) { c.middleware = append(c.middleware, mw) }
}

// NewClient creates a new Client with the given base URL and options.
func NewClient(baseURL string, opts ...Option) *Client {
	c := &Client{baseURL: baseURL}
	for _, opt := range opts {
		opt(c)
	}
	if c.httpClient == nil {
		c.httpClient = NewLoggingClientWithLogger(c.logger)
	}
	c.Pages = NewPageService(c)
	c.Spaces = NewSpaceService(c)
	c.Comments = NewCommentService(c)
	c.Tags = NewTagService(c)
	c.Categories = NewCategoryService(c)
	c.Search = NewSearchService(c)
	return c
}

// doRequest executes an HTTP request through the middleware chain.
// It sets Authorization, Content-Type, and Accept headers before dispatching.
func (c *Client) doRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	if c.apiKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Build the base handler that calls the underlying http.Client.
	var handler RequestHandler = func(ctx context.Context, r *http.Request) (*http.Response, error) {
		return c.httpClient.Do(r)
	}

	// Wrap with middleware in reverse order so the last-added runs first.
	for i := len(c.middleware) - 1; i >= 0; i-- {
		mw := c.middleware[i]
		next := handler
		handler = func(ctx context.Context, r *http.Request) (*http.Response, error) {
			return mw(ctx, r, next)
		}
	}

	return handler(ctx, req)
}

// wrapRequest serialises data as `{"resource": <data>}` — the format required
// by the Spaces API request body parser.
func wrapRequest(resource string, data any) ([]byte, error) {
	inner, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	// Build {"resource": <inner>}
	out := make([]byte, 0, len(resource)+len(inner)+6)
	out = append(out, '{', '"')
	out = append(out, []byte(resource)...)
	out = append(out, '"', ':')
	out = append(out, inner...)
	out = append(out, '}')
	return out, nil
}
