# CLAUDE.md — Go SDK Coding Standards

This file documents the architecture, conventions, and patterns used in this codebase. All new work must follow these standards exactly.

---

## Project Overview

This is a Go SDK for a REST API. It provides a typed, generic client that wraps HTTP calls into service objects per API resource. The SDK is consumed as a library; `main.go` is a demo/CLI entry point only.

**Module:** `github.com/teamwork/desksdkgo`  
**Go version:** `1.24.2`

---

## Directory Structure

```
.
├── api/            # Generic Service interface definition (used for type assertions / documentation)
│   └── api.go
├── client/         # HTTP client, all service implementations, middleware, mocks
│   ├── client.go       # Client struct, Option funcs, ListOptions, doRequest
│   ├── resource.go     # Generic Service[T, L] base with Get/List/Create/Update
│   ├── path.go         # PathHandler interface + DefaultPathHandler
│   ├── middleware.go   # MiddlewareFunc implementations (logging, retry, auth, etc.)
│   ├── transport.go    # LoggingTransport (http.RoundTripper)
│   ├── mock_client.go  # MockRoundTripper for tests
│   ├── filter.go       # MongoDB-style FilterBuilder
│   ├── <resource>.go   # One file per resource (tickets.go, messages.go, etc.)
│   └── <resource>_test.go
├── models/         # All data types: domain models, request/response wrappers
│   ├── base.go         # BaseEntity, EntityRef, UserRef, State
│   ├── response.go     # Pagination, PageMeta, IncludedData, Meta
│   └── <resource>.go   # One file per resource domain
├── util/
│   ├── env.go          # .env loading helpers
│   └── json.go         # MergeJSONData utility
└── main.go             # Demo/CLI only — not part of the library API
```

---

## Dependencies

All dependencies are in `go.mod`. Do not add dependencies without strong justification.

| Package                           | Version | Purpose                                                |
| --------------------------------- | ------- | ------------------------------------------------------ |
| `github.com/brianvoe/gofakeit/v7` | v7.2.1  | Fake data in tests                                     |
| `github.com/gorilla/schema`       | v1.4.1  | Form/query string decoding                             |
| `github.com/joho/godotenv`        | v1.5.1  | Load `.env` files                                      |
| `github.com/sonh/qs`              | v0.6.4  | Encode structs to query strings (use `qs` struct tags) |

Use the **standard library** for: JSON (`encoding/json`), HTTP (`net/http`), logging (`log/slog`), context (`context`), URL building (`net/url`).

---

## Package Layout Rules

- `models/` — pure data. No imports from `client/`. No HTTP logic.
- `client/` — all HTTP logic. Imports `models/`. One service file per resource.
- `util/` — stateless helpers. No imports from `client/` or `models/`.
- `api/` — interface definitions only. Kept minimal.

---

## Client Architecture

### The `Client` Struct

Defined in `client/client.go`. It holds configuration and one exported field per resource service:

```go
type Client struct {
    baseURL    string
    apiKey     string
    logLevel   slog.Level
    logger     *slog.Logger
    httpClient *http.Client
    middleware []MiddlewareFunc

    // One exported field per service, PascalCase, named after the resource (plural)
    Tickets          *TicketService
    Messages         *MessageService
    Companies        *CompanyService
    // ...
}
```

### Constructor

```go
func NewClient(baseURL string, opts ...Option) *Client
```

- Apply all `Option` funcs first.
- Default `httpClient` to `NewLoggingClientWithLogger(...)` if not set via option.
- Initialize all services by calling their `New<Resource>Service(client)` constructors.

### Functional Options

Use the functional options pattern for all client configuration. Define each option as:

```go
type Option func(*Client)

func WithAPIKey(apiKey string) Option {
    return func(c *Client) { c.apiKey = apiKey }
}
```

Available options:

- `WithAPIKey(apiKey string)`
- `WithHTTPClient(httpClient *http.Client)`
- `WithLogLevel(level slog.Level)`
- `WithLogger(logger *slog.Logger)`
- `WithMiddleware(mw MiddlewareFunc)`

### `doRequest`

All HTTP calls go through `(*Client).doRequest(ctx, req)`. It:

1. Sets `Authorization: Bearer <apiKey>` if `apiKey` is non-empty.
2. Sets `Content-Type: application/json` and `Accept: application/json`.
3. Executes the middleware chain in **reverse order** (last added = first executed).
4. Calls `c.httpClient.Do(req)` as the final handler.

Never bypass `doRequest` in service methods.

### `ListOptions`

```go
type ListOptions struct {
    Page    int
    PerPage int
    SortBy  string
    SortDir string
    Embed   string
    Fields  string
    Q       string
}

func (o *ListOptions) Encode() string // returns url-encoded query string
```

Pass `ListOptions.Encode()` output as query parameters for list endpoints. Zero-value fields are omitted.

---

## Generic Service Base

`client/resource.go` defines the reusable CRUD base:

```go
type Service[T any, L any] struct {
    client *Client
    router PathHandler
}

func NewService[T any, L any](client *Client, router PathHandler) *Service[T, L]
```

- `T` = single-resource response type (e.g., `models.TicketResponse`)
- `L` = list response type (e.g., `models.TicketsResponse`)

Methods provided by `Service[T, L]`:

| Method                                         | HTTP         | URL pattern                      | Success codes |
| ---------------------------------------------- | ------------ | -------------------------------- | ------------- |
| `Get(ctx, id int) (*T, error)`                 | GET          | `/<base>/<id>.json?includes=all` | 200           |
| `List(ctx, params url.Values) (*L, error)`     | GET          | `/<base>.json?<params>`          | 200           |
| `Create(ctx, resource *T) (*T, error)`         | POST         | `/<base>.json`                   | 200 or 201    |
| `Update(ctx, id int, resource *T) (*T, error)` | PUT or PATCH | `/<base>/<id>.json`              | 200           |

All methods:

- Use `http.NewRequestWithContext` — never `http.NewRequest`.
- Log errors via `s.logError(msg, attrs...)` before returning.
- Read and close `resp.Body` with `defer resp.Body.Close()`.
- Decode with `json.NewDecoder(resp.Body).Decode(&resource)`.
- On unexpected status: read the body, log it, return `fmt.Errorf("unexpected status code: %d", resp.StatusCode)`.

---

## PathHandler Interface

```go
type PathHandler interface {
    Get(id int) string    // returns "resource/123"
    List() string         // returns "resource"
    Create() string       // returns "resource"
    Update(id int) string // returns "resource/123"
}

type updateMethodProvider interface {
    UpdateMethod() string // return http.MethodPatch or http.MethodPut
}
```

Use `DefaultPathHandler` for standard REST paths:

```go
// Default update method is PUT
NewDefaultPathHandler("tickets")

// Override update method to PATCH
NewDefaultPathHandlerWithUpdateMethod("tickets", http.MethodPatch)
```

`Service.Update` checks if `router` implements `updateMethodProvider` and uses its method if so.

---

## Resource Service Pattern

Every resource gets its own file `client/<resource>s.go` (lowercase plural). The service:

1. Embeds `*Service[models.<Resource>Response, models.<Resource>sResponse]`
2. Holds a `client *Client` for custom operations
3. Has a `New<Resource>Service(client *Client) *<Resource>Service` constructor
4. Wraps every generic method explicitly (do not expose the embedded struct directly)
5. Adds specialized methods for non-standard operations

**Minimal service example:**

```go
type TagService struct {
    *Service[models.TagResponse, models.TagsResponse]
    client *Client
}

func NewTagService(client *Client) *TagService {
    return &TagService{
        Service: NewService[models.TagResponse, models.TagsResponse](
            client,
            NewDefaultPathHandler("tags"),
        ),
        client: client,
    }
}

func (s *TagService) Get(ctx context.Context, id int) (*models.TagResponse, error) {
    return s.Service.Get(ctx, id)
}

func (s *TagService) List(ctx context.Context, params url.Values) (*models.TagsResponse, error) {
    return s.Service.List(ctx, params)
}

func (s *TagService) Create(ctx context.Context, tag *models.TagResponse) (*models.TagResponse, error) {
    return s.Service.Create(ctx, tag)
}

func (s *TagService) Update(ctx context.Context, id int, tag *models.TagResponse) (*models.TagResponse, error) {
    return s.Service.Update(ctx, id, tag)
}
```

**Custom sub-resource method example** (see `messages.go`):

```go
func (s *MessageService) CreateForTicket(ctx context.Context, ticketID int, message *models.MessageResponse) (*models.MessageResponse, error) {
    if ticketID <= 0 {
        return nil, fmt.Errorf("ticketID must be greater than 0")
    }
    if message == nil {
        return nil, fmt.Errorf("message is required")
    }

    body, err := json.Marshal(message.Message)
    if err != nil {
        return nil, err
    }

    req, err := http.NewRequestWithContext(ctx, http.MethodPost,
        fmt.Sprintf("%s/tickets/%d/messages.json", s.client.baseURL, ticketID),
        bytes.NewBuffer(body))
    if err != nil {
        return nil, err
    }

    resp, err := s.client.doRequest(ctx, req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
        b, err := io.ReadAll(resp.Body)
        if err != nil {
            return nil, err
        }
        return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(b))
    }

    var result models.MessageResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }
    return &result, nil
}
```

---

## Models

### File Organization

- `models/base.go` — `BaseEntity`, `EntityRef`, `UserRef`, `State`
- `models/response.go` — `Pagination`, `PageMeta`, `IncludedData`, `Meta`
- `models/<resource>.go` — domain model + its single/list response wrappers

All three types (domain model, single response, list response) for a resource live in the **same file**.

### `BaseEntity`

All primary resource types embed `BaseEntity`:

```go
type BaseEntity struct {
    ID        int        `json:"id"`
    Type      any        `json:"type"` // string or object depending on API
    CreatedAt *time.Time `json:"createdAt,omitempty"`
    UpdatedAt *time.Time `json:"updatedAt,omitempty"`
    CreatedBy *UserRef   `json:"createdBy,omitempty"`
    UpdatedBy *UserRef   `json:"updatedBy,omitempty"`
    State     State      `json:"state"`
}
```

### `EntityRef`

Use `*EntityRef` for optional references to related objects (agent, inbox, status, type, priority, etc.):

```go
type EntityRef struct {
    ID   int            `json:"id"`
    Type string         `json:"type"`
    Meta map[string]any `json:"meta"`
}
```

Use `[]EntityRef` (not pointer) for collections of references (tags, files, messages, activities).

### Response Wrapper Pattern

Every resource must have exactly these two response types:

```go
// Single resource — used by Get, Create, Update
type TicketResponse struct {
    Ticket   Ticket       `json:"ticket"`
    Included IncludedData `json:"included"`
}

// List — used by List, Search
type TicketsResponse struct {
    Tickets    []Ticket     `json:"tickets"`
    Included   IncludedData `json:"included"`
    Pagination Pagination   `json:"pagination"`
    Meta       Meta         `json:"meta"`
}
```

Naming: `<Resource>Response` (singular) and `<Resource>sResponse` (plural). The JSON key in the single response is the lowercase resource name; in the list response it is the lowercase plural.

### `IncludedData`

All sideloaded/embedded resources are decoded into `IncludedData` in `models/response.go`. When adding a new resource that can appear as included data, add a field here:

```go
type IncludedData struct {
    Companies []Company `json:"companies"`
    // add new resource here
}
```

### JSON Struct Tags

- Always include `json:"fieldName"` on every exported field.
- Use `omitempty` for optional fields: `json:"field,omitempty"`.
- Use `*time.Time` (pointer) for optional timestamps; use `time.Time` (non-pointer) only when the field is always present.
- Use `*EntityRef` for optional single references.
- Use `[]EntityRef` (nil-able slice, with `omitempty`) for optional collections.
- Use `any` only when the API genuinely returns heterogeneous types for a field.

### `State` Enum

```go
type State string

const (
    StateActive  State = "active"
    StateDeleted State = "deleted"
)
```

Define domain-specific enums the same way: named string type + exported constants.

### Custom Unmarshalers

Use the alias pattern when a field must be remapped:

```go
func (m *Message) UnmarshalJSON(data []byte) error {
    type Alias Message
    aux := &struct {
        HtmlBody string `json:"htmlBody"`
        *Alias
    }{Alias: (*Alias)(m)}
    if err := json.Unmarshal(data, aux); err != nil {
        return err
    }
    m.Body = aux.HtmlBody
    return nil
}
```

### Query String Filter Structs

For complex search/filter params, define a dedicated struct with `qs` tags (not `json`):

```go
type SearchTicketsFilter struct {
    Search    string     `qs:"search"`
    Tags      []int64    `qs:"tags"`
    StartDate *time.Time `qs:"startDate,omitempty"`
    // ...
}
```

Encode with `github.com/sonh/qs`:

```go
encoder := qs.NewEncoder()
values, err := encoder.Values(filter)
// use values.Encode() as query string
```

---

## Middleware

Middleware lives in `client/middleware.go`. The type signatures are:

```go
type MiddlewareFunc func(ctx context.Context, req *http.Request, next RequestHandler) (*http.Response, error)
type RequestHandler  func(ctx context.Context, req *http.Request) (*http.Response, error)
```

Each middleware is a factory function returning a `MiddlewareFunc`:

```go
func LoggingMiddleware(logger *slog.Logger) MiddlewareFunc {
    return func(ctx context.Context, req *http.Request, next RequestHandler) (*http.Response, error) {
        // pre-request work
        resp, err := next(ctx, req)
        // post-request work
        return resp, err
    }
}
```

Available middleware (do not duplicate):

- `LoggingMiddleware(logger)` — logs method, URL, status, duration
- `RetryMiddleware(maxRetries, retryDelay)` — retries on error, clones request per attempt
- `AuthMiddleware(token)` — sets `Authorization: Bearer <token>`
- `UserAgentMiddleware(userAgent)` — sets `User-Agent`
- `RateLimitMiddleware(requestsPerSecond)` — token-bucket rate limiting
- `RequestIDMiddleware()` — adds `X-Request-ID` header
- `TimeoutMiddleware(timeout)` — wraps context with deadline
- `HeaderMiddleware(headers map[string]string)` — adds arbitrary headers
- `ConditionalMiddleware(condition, middleware)` — conditional application

The chain is applied in **reverse append order**: `middleware[last]` runs first.

---

## Error Handling

- Return `(*T, error)` — never panic in library code.
- Validate inputs at method entry; return `fmt.Errorf("fieldName is required")` or `fmt.Errorf("fieldName must be > 0")`.
- On unexpected HTTP status: read body, include status code and body in the error string.
- Use `s.logError(msg, slog.Attr...)` to log before returning, only when a `logger` is configured.
- Use `fmt.Errorf(...)` — no custom error types, no `errors.New` with sentinel errors.

Error message formats:

```go
fmt.Errorf("unexpected status code: %d", resp.StatusCode)                     // when no body read
fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, body)    // when body is read
```

---

## Logging

Use `log/slog` (standard library) throughout. Never use `fmt.Println`, `log.Printf`, or third-party loggers.

```go
s.client.logger.LogAttrs(ctx, slog.LevelError, "message",
    slog.Any("error", err),
    slog.String("method", http.MethodGet),
    slog.String("url", req.URL.String()),
    slog.Int("status_code", resp.StatusCode),
    slog.String("response_body", string(body)),
)
```

Always guard logger calls: `if s.client != nil && s.client.logger != nil`.

Middleware uses `logger.InfoContext` / `logger.ErrorContext` with context.

---

## Testing

### Test Files

Place test files next to what they test: `client/tickets_test.go` for `client/tickets.go`.

### Mock HTTP Transport

Use `MockRoundTripper` from `client/mock_client.go`. Never use external mocking libraries.

```go
func TestMyService(t *testing.T) {
    mockTransport := NewMockRoundTripper()
    mockTransport.AddResponse(
        http.MethodPost,
        "/tickets/123/messages.json",
        http.StatusCreated,
        expectedResponseStruct, // marshaled to JSON automatically
    )

    c := NewClient("https://example.com",
        WithHTTPClient(&http.Client{Transport: mockTransport}),
    )

    resp, err := c.Messages.CreateForTicket(context.Background(), 123, &inputMessage)
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }

    requests := mockTransport.GetRequests()
    if len(requests) != 1 {
        t.Fatalf("expected 1 request, got %d", len(requests))
    }
    if requests[0].URL.Path != "/tickets/123/messages.json" {
        t.Fatalf("unexpected path: %s", requests[0].URL.Path)
    }
}
```

`AddResponse` accepts: `io.ReadCloser`, `string`, or any value (marshaled to JSON).

### Test Data

Use `github.com/brianvoe/gofakeit/v7` for realistic fake values:

```go
import "github.com/brianvoe/gofakeit/v7"

subject := gofakeit.Sentence(5)
email   := gofakeit.Email()
body    := gofakeit.Paragraph(1, 3, 10, " ")
```

### Assertions

Use only the standard `testing` package. No testify or other assertion libraries.

```go
if err != nil {
    t.Fatalf("expected no error, got %v", err)
}
if got != want {
    t.Errorf("got %v, want %v", got, want)
}
```

---

## Naming Conventions

### Files

| Location           | Pattern               | Examples                                         |
| ------------------ | --------------------- | ------------------------------------------------ |
| `models/`          | lowercase snake_case  | `ticket.go`, `ticket_priority.go`, `base.go`     |
| `client/` services | lowercase plural      | `tickets.go`, `companies.go`, `slas.go`          |
| `client/` infra    | descriptive lowercase | `resource.go`, `middleware.go`, `mock_client.go` |
| tests              | `<file>_test.go`      | `tickets_test.go`, `filter_test.go`              |

### Types

| Thing                | Convention                  | Example                       |
| -------------------- | --------------------------- | ----------------------------- |
| Domain model         | PascalCase                  | `Ticket`, `Company`           |
| Single response      | `<Model>Response`           | `TicketResponse`              |
| List response        | `<Model>sResponse`          | `TicketsResponse`             |
| Service type         | `<Model>Service`            | `TicketService`               |
| Filter/search struct | `<Action><Model>Filter`     | `SearchTicketsFilter`         |
| Enum type            | named string type           | `type State string`           |
| Enum constants       | PascalCase with type prefix | `StateActive`, `StateDeleted` |

### Functions and Methods

| Thing                | Convention                        | Example                                |
| -------------------- | --------------------------------- | -------------------------------------- |
| Service constructor  | `New<Service>`                    | `NewTicketService`                     |
| CRUD methods         | `Get`, `List`, `Create`, `Update` | (always these four names)              |
| Sub-resource methods | `<Verb>For<Parent>`               | `CreateForTicket`                      |
| Option funcs         | `With<Field>`                     | `WithAPIKey`, `WithLogLevel`           |
| Middleware factories | `<Behavior>Middleware`            | `LoggingMiddleware`, `RetryMiddleware` |

### Variables and Receivers

- Receiver name: single lowercase letter or two-letter abbreviation. `s` for services, `c` for client, `m` for models.
- Private fields: camelCase (`apiKey`, `httpClient`, `baseURL`).
- Service fields on `Client`: PascalCase plural (`Tickets`, `Messages`, `Companies`).

---

## URL Construction

All URLs follow the pattern: `<baseURL>/<resource>/<id>.json[?<query>]`

```go
// Get
fmt.Sprintf("%s/%s.json?includes=all", s.client.baseURL, s.router.Get(id))

// List
fmt.Sprintf("%s/%s.json?%s", s.client.baseURL, s.router.List(), params.Encode())

// Create
fmt.Sprintf("%s/%s.json", s.client.baseURL, s.router.Create())

// Update
fmt.Sprintf("%s/%s.json", s.client.baseURL, s.router.Update(id))

// Sub-resource
fmt.Sprintf("%s/tickets/%d/messages.json", s.client.baseURL, ticketID)

// Search
fmt.Sprintf("%s/search/tickets.json?%s", s.client.baseURL, values.Encode())
```

Always append `.json` to resource paths. Always use `http.NewRequestWithContext`.

---

## Input Validation

Validate at the service method boundary only — not inside `doRequest` or model constructors. Check:

- Required IDs: `if id <= 0 { return nil, fmt.Errorf("id must be greater than 0") }`
- Required pointers: `if resource == nil { return nil, fmt.Errorf("resource is required") }`
- Required nested IDs: `if resource.Message.Ticket.ID <= 0 { return nil, fmt.Errorf("...") }`

Do not validate fields that the API will validate (formats, lengths, enums). Only guard against panics and obviously broken calls.

---

## Patterns to Avoid

- Do not add global state or package-level variables.
- Do not use `init()`.
- Do not use `log.Fatal`, `os.Exit`, or `panic` in library code.
- Do not introduce interfaces beyond `PathHandler`, `updateMethodProvider`, and the `api.Service` interface unless there is a concrete need for multiple implementations.
- Do not add caching inside the SDK — leave that to callers.
- Do not retry inside individual service methods — use `RetryMiddleware` instead.
- Do not embed `http.Client` directly in service types — always use the shared `Client.httpClient` via `doRequest`.
- Do not use `gorilla/schema` for encoding (only decoding). Use `sonh/qs` for encoding structs to query strings.
