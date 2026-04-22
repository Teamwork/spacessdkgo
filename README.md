# Teamwork Spaces Go SDK

A typed Go client for the Teamwork Spaces REST API.

**Module:** `github.com/teamwork/spacessdkgo`  
**Go version:** 1.24.2+

---

## Installation

```bash
go get github.com/teamwork/spacessdkgo
```

---

## Quick Start

```go
package main

import (
    "context"
    "log/slog"
    "net/url"
    "os"

    "github.com/teamwork/spacessdkgo/client"
    "github.com/teamwork/spacessdkgo/models"
)

func main() {
    c := client.NewClient(
        "https://your-instance.teamwork.com",
        client.WithAPIKey("your-api-key"),
        client.WithLogLevel(slog.LevelInfo),
    )

    ctx := context.Background()

    // List spaces
    spaces, err := c.Spaces.List(ctx, url.Values{})
    if err != nil {
        slog.Error("failed to list spaces", "error", err)
        os.Exit(1)
    }

    for _, space := range spaces.Spaces {
        slog.Info("space", "id", space.ID, "title", space.Title)
    }
}
```

---

## Configuration

### Client Constructor

```go
c := client.NewClient(baseURL string, opts ...client.Option) *client.Client
```

### Available Options

| Option | Description |
|--------|-------------|
| `WithAPIKey(key string)` | Sets `Authorization: Bearer <key>` on all requests |
| `WithHTTPClient(hc *http.Client)` | Replaces the default HTTP client |
| `WithLogLevel(level slog.Level)` | Sets the minimum log level |
| `WithLogger(logger *slog.Logger)` | Replaces the default structured logger |
| `WithMiddleware(mw MiddlewareFunc)` | Appends middleware to the chain |

### Environment Variables (demo/CLI only)

The `main.go` demo reads:

- `SPACES_API_BASE_URL` — API base URL
- `SPACES_API_KEY` — API key

Use `util.LoadEnv()` to load a `.env` file in your own tooling.

---

## Services

The client exposes one field per resource service:

```go
c.Spaces      // *client.SpaceService
c.Pages       // *client.PageService
c.Comments    // *client.CommentService
c.Tags        // *client.TagService
c.Categories  // *client.CategoryService
c.Search      // *client.SearchService
```

---

### Spaces

```go
// Get a single space
space, err := c.Spaces.Get(ctx, id int64)

// List all spaces
spaces, err := c.Spaces.List(ctx, params url.Values)

// Create a space
space, err := c.Spaces.Create(ctx, &models.SpaceCreate{
    Title:   "Engineering Wiki",
    Purpose: "Internal docs",
})

// Update a space (PATCH)
space, err := c.Spaces.Update(ctx, id int64, &models.SpaceUpdate{
    Title: "Updated Title",
})

// Delete a space
err := c.Spaces.Delete(ctx, id int64)

// List collaborators
collaborators, err := c.Spaces.Collaborators(ctx, id int64)
```

---

### Pages

```go
// Get a single page
page, err := c.Pages.Get(ctx, spaceID, pageID int64)

// List pages in a space
pages, err := c.Pages.List(ctx, spaceID int64, params url.Values)

// Get the home page of a space
home, err := c.Pages.Home(ctx, spaceID int64)

// Create a page
page, err := c.Pages.Create(ctx, spaceID int64, &models.PageCreate{
    Title:   "Getting Started",
    Content: "<p>Welcome!</p>",
})

// Duplicate a page
page, err := c.Pages.Duplicate(ctx, spaceID, pageID int64, &models.PageDuplicate{
    Title: "Copy of Getting Started",
})

// Update a page (PATCH)
page, err := c.Pages.Update(ctx, spaceID, pageID int64, &models.PageUpdate{
    Content: "<p>Updated content</p>",
})

// Delete a page
err := c.Pages.Delete(ctx, spaceID, pageID int64)
```

---

### Comments

```go
// Get a single comment
comment, err := c.Comments.Get(ctx, spaceID, pageID, commentID int64)

// List comments on a page
comments, err := c.Comments.List(ctx, spaceID, pageID int64, params url.Values)

// Create a comment
comment, err := c.Comments.Create(ctx, spaceID, pageID int64, &models.CommentCreate{
    Content: "Great page!",
})

// Update a comment (PATCH)
comment, err := c.Comments.Update(ctx, spaceID, pageID, commentID int64, &models.CommentUpdate{
    Content: "Edited: Great page!",
})

// Delete a comment
err := c.Comments.Delete(ctx, spaceID, pageID, commentID int64)
```

---

### Tags

```go
// Get a single tag
tag, err := c.Tags.Get(ctx, id int64)

// List all tags
tags, err := c.Tags.List(ctx, params url.Values)

// Create multiple tags in one request
tags, err := c.Tags.CreateBatch(ctx, []models.Tag{
    {Name: "go", Color: "#00ADD8"},
    {Name: "api", Color: "#FF6B6B"},
})

// Update a tag (PATCH)
tag, err := c.Tags.Update(ctx, id int64, &models.TagUpdate{
    Name: "golang",
})

// Delete a tag
err := c.Tags.Delete(ctx, id int64)
```

---

### Categories

```go
// Get a single category
cat, err := c.Categories.Get(ctx, id int64)

// List all categories
cats, err := c.Categories.List(ctx, params url.Values)

// Create a category
cat, err := c.Categories.Create(ctx, &models.CategoryCreate{
    Name:  "Engineering",
    Color: "#4A90E2",
})

// Update a category (PATCH)
cat, err := c.Categories.Update(ctx, id int64, &models.CategoryUpdate{
    Name: "Engineering Teams",
})

// Delete a category
err := c.Categories.Delete(ctx, id int64)
```

---

### Search

```go
results, err := c.Search.Search(ctx, models.SearchFilter{
    Query:   "deployment guide",
    SpaceID: []int64{42, 57},
    Limit:   func() *int64 { v := int64(20); return &v }(),
})

for _, r := range results.Results {
    fmt.Printf("page %d: %s (space %d)\n", r.PageID, r.Title, r.Space.ID)
    for field, snippet := range r.MatchedText {
        fmt.Printf("  matched in %s: %s\n", field, snippet)
    }
}
```

---

## Middleware

Middleware wraps every HTTP request. Add it via `WithMiddleware`; multiple calls are allowed. The chain executes in **reverse append order** (last-added runs first).

```go
c := client.NewClient(baseURL,
    client.WithAPIKey(apiKey),
    client.WithMiddleware(client.RetryMiddleware(3, 500*time.Millisecond)),
    client.WithMiddleware(client.RequestIDMiddleware()),
    client.WithMiddleware(client.LoggingMiddleware(logger)),
)
```

### Available Middleware

| Middleware | Description |
|------------|-------------|
| `LoggingMiddleware(logger)` | Logs method, URL, status, and duration |
| `RetryMiddleware(maxRetries, retryDelay)` | Retries failed requests |
| `AuthMiddleware(token)` | Sets `Authorization: Bearer` (alternative to `WithAPIKey`) |
| `UserAgentMiddleware(userAgent)` | Sets the `User-Agent` header |
| `RateLimitMiddleware(requestsPerSecond)` | Token-bucket rate limiting |
| `RequestIDMiddleware()` | Adds a unique `X-Request-ID` to each request |
| `TimeoutMiddleware(timeout)` | Wraps request context with a deadline |
| `HeaderMiddleware(headers)` | Adds static headers to all requests |
| `ConditionalMiddleware(condition, mw)` | Applies middleware only when `condition` returns true |

### Writing Custom Middleware

```go
func TracingMiddleware(traceID string) client.MiddlewareFunc {
    return func(ctx context.Context, req *http.Request, next client.RequestHandler) (*http.Response, error) {
        req.Header.Set("X-Trace-ID", traceID)
        return next(ctx, req)
    }
}
```

---

## Pagination

List responses include a `Meta` field with page information:

```go
resp, err := c.Spaces.List(ctx, url.Values{
    "page":    []string{"2"},
    "perPage": []string{"25"},
})

fmt.Println(resp.Meta.Page.Count)  // total items
fmt.Println(resp.Meta.Page.Size)   // page size
fmt.Println(resp.Meta.Page.Offset) // current offset
```

---

## Models

### Core Types

| Model | Key Fields |
|-------|-----------|
| `Space` | `ID`, `Title`, `Code`, `Purpose`, `SpaceColor`, `Icon`, `State` |
| `Page` | `ID`, `Title`, `Slug`, `Content`, `ContentRevision`, `Tags`, `IsPublished`, `IsPrivate` |
| `Comment` | `ID`, `Content`, `ParentID`, `IsPrivate`, `State` |
| `Tag` | `ID`, `Name`, `Color`, `PageCount` |
| `Category` | `ID`, `Name`, `Color` |

### State Constants

```go
models.StateActive    // "active"
models.StateArchived  // "archived"
models.StateRemoved   // "removed"
models.StateDeleted   // "deleted"
```

### Included Data (Sideloading)

Get and List responses include an `Included` field with sideloaded related resources:

```go
page, err := c.Pages.Get(ctx, spaceID, pageID)
// page.Included contains related entities embedded in the response
```

---

## Error Handling

All methods return `(*T, error)`. Errors are plain `error` values — no custom types or sentinel errors.

```go
space, err := c.Spaces.Get(ctx, 123)
if err != nil {
    // err contains the HTTP status code and response body on API errors:
    // "unexpected status code: 404, body: ..."
    log.Fatal(err)
}
```

---

## Testing

The SDK ships a `MockRoundTripper` for unit tests. It intercepts HTTP calls without a real server.

```go
import "github.com/teamwork/spacessdkgo/client"

func TestMyCode(t *testing.T) {
    mock := client.NewMockRoundTripper()
    mock.AddResponse(
        http.MethodGet,
        "/spaces/42.json",
        http.StatusOK,
        models.SpaceResponse{Space: models.Space{ID: 42, Title: "Docs"}},
    )

    c := client.NewClient("https://example.com",
        client.WithHTTPClient(&http.Client{Transport: mock}),
    )

    space, err := c.Spaces.Get(context.Background(), 42)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if space.Space.Title != "Docs" {
        t.Errorf("got %q, want %q", space.Space.Title, "Docs")
    }

    reqs := mock.GetRequests()
    if len(reqs) != 1 {
        t.Fatalf("expected 1 request, got %d", len(reqs))
    }
}
```

`AddResponse` accepts: `io.ReadCloser`, `string`, or any value (marshaled to JSON automatically).

---

## Project Layout

```
.
├── api/            # Service interface definitions
├── client/         # HTTP client, all service implementations, middleware
│   ├── client.go       # Client struct and constructor
│   ├── resource.go     # Generic Service[T, L] base (Get/List/Create/Update)
│   ├── middleware.go   # All middleware implementations
│   ├── transport.go    # LoggingTransport (http.RoundTripper)
│   ├── mock_client.go  # MockRoundTripper for tests
│   ├── filter.go       # MongoDB-style FilterBuilder
│   ├── spaces.go
│   ├── pages.go
│   ├── comments.go
│   ├── tags.go
│   ├── categories.go
│   └── search.go
├── models/         # All data types: domain models and response wrappers
│   ├── base.go         # BaseEntity, EntityRef, UserRef, State
│   ├── response.go     # Pagination, Meta, IncludedData
│   └── <resource>.go   # Domain model + request/response types per resource
├── util/
│   ├── env.go          # .env loading helpers
│   └── json.go         # MergeJSONData utility
└── main.go             # Demo/CLI — not part of the library API
```

---

## Contributing

Follow the patterns in [CLAUDE.md](CLAUDE.md) for all new code. Key rules:

- One file per resource in `client/` and `models/`
- All HTTP calls go through `doRequest` — never bypass it
- Use `http.NewRequestWithContext` — never `http.NewRequest`
- Return `(*T, error)` — no panics in library code
- Use `log/slog` — never `fmt.Println` or `log.Printf`
- Tests use only the standard `testing` package and `MockRoundTripper`
