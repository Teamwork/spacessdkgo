package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

// mockResponse holds a queued HTTP response for a specific method+path pair.
type mockResponse struct {
	method string
	path   string
	status int
	body   io.ReadCloser
}

// MockRoundTripper is a fake http.RoundTripper for use in tests.
// It queues responses and records all requests made.
type MockRoundTripper struct {
	mu        sync.Mutex
	responses []mockResponse
	requests  []*http.Request
}

// NewMockRoundTripper creates an empty MockRoundTripper.
func NewMockRoundTripper() *MockRoundTripper {
	return &MockRoundTripper{}
}

// AddResponse queues a response for the given method and path.
// body may be an io.ReadCloser, a string, or any value (marshaled to JSON).
func (m *MockRoundTripper) AddResponse(method, path string, status int, body any) {
	var rc io.ReadCloser
	switch v := body.(type) {
	case io.ReadCloser:
		rc = v
	case string:
		rc = io.NopCloser(bytes.NewBufferString(v))
	default:
		data, err := json.Marshal(v)
		if err != nil {
			panic(fmt.Sprintf("mock_client: failed to marshal response body: %v", err))
		}
		rc = io.NopCloser(bytes.NewBuffer(data))
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	m.responses = append(m.responses, mockResponse{
		method: method,
		path:   path,
		status: status,
		body:   rc,
	})
}

// RoundTrip records the request and returns the next queued response that matches
// the request method and path. If no matching response is found it returns a 500.
func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.requests = append(m.requests, req)

	for i, r := range m.responses {
		if r.method == req.Method && r.path == req.URL.Path {
			m.responses = append(m.responses[:i], m.responses[i+1:]...)
			return &http.Response{
				StatusCode: r.status,
				Body:       r.body,
				Header:     make(http.Header),
				Request:    req,
			}, nil
		}
	}

	return &http.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       io.NopCloser(bytes.NewBufferString("no mock response queued")),
		Header:     make(http.Header),
		Request:    req,
	}, nil
}

// GetRequests returns a copy of all recorded requests.
func (m *MockRoundTripper) GetRequests() []*http.Request {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]*http.Request, len(m.requests))
	copy(out, m.requests)
	return out
}

// Reset clears all queued responses and recorded requests.
func (m *MockRoundTripper) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.responses = nil
	m.requests = nil
}
