package client

import (
	"fmt"
	"net/http"
)

// PathHandler generates URL path segments for a resource.
type PathHandler interface {
	Get(id int64) string    // returns "resource/123"
	List() string           // returns "resource"
	Create() string         // returns "resource"
	Update(id int64) string // returns "resource/123"
}

// updateMethodProvider optionally specifies the HTTP method used for updates.
type updateMethodProvider interface {
	UpdateMethod() string // returns http.MethodPatch or http.MethodPut
}

// DefaultPathHandler is the standard implementation of PathHandler for flat REST resources.
type DefaultPathHandler struct {
	resource     string
	updateMethod string
}

// NewDefaultPathHandler creates a PathHandler that uses PUT for updates.
func NewDefaultPathHandler(resource string) *DefaultPathHandler {
	return &DefaultPathHandler{
		resource:     resource,
		updateMethod: http.MethodPut,
	}
}

// NewDefaultPathHandlerWithUpdateMethod creates a PathHandler with a custom update HTTP method.
func NewDefaultPathHandlerWithUpdateMethod(resource, method string) *DefaultPathHandler {
	return &DefaultPathHandler{
		resource:     resource,
		updateMethod: method,
	}
}

// Get returns the path for fetching a single resource by ID.
func (h *DefaultPathHandler) Get(id int64) string {
	return fmt.Sprintf("%s/%d", h.resource, id)
}

// List returns the path for listing resources.
func (h *DefaultPathHandler) List() string {
	return h.resource
}

// Create returns the path for creating a new resource.
func (h *DefaultPathHandler) Create() string {
	return h.resource
}

// Update returns the path for updating a resource by ID.
func (h *DefaultPathHandler) Update(id int64) string {
	return fmt.Sprintf("%s/%d", h.resource, id)
}

// UpdateMethod returns the HTTP method used for update operations.
func (h *DefaultPathHandler) UpdateMethod() string {
	return h.updateMethod
}
