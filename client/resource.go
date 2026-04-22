package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Service is the generic base for all resource services.
// T is the single-resource response type; L is the list response type.
type Service[T any, L any] struct {
	client *Client
	router PathHandler
}

// NewService creates a new Service with the given client and router.
func NewService[T any, L any](client *Client, router PathHandler) *Service[T, L] {
	return &Service[T, L]{
		client: client,
		router: router,
	}
}

// Get fetches a single resource by ID.
func (s *Service[T, L]) Get(ctx context.Context, id int64) (*T, error) {
	if id <= 0 {
		return nil, fmt.Errorf("id must be greater than 0")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("%s/%s.json", s.client.baseURL, s.router.Get(id)), nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := s.client.checkResponse(ctx, resp, "Get", http.StatusOK); err != nil {
		return nil, err
	}
	return decodeJSON[T](resp.Body)
}

// List fetches a collection of resources with optional query parameters.
func (s *Service[T, L]) List(ctx context.Context, params url.Values) (*L, error) {
	u := fmt.Sprintf("%s/%s.json", s.client.baseURL, s.router.List())
	if len(params) > 0 {
		u = fmt.Sprintf("%s?%s", u, params.Encode())
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := s.client.checkResponse(ctx, resp, "List", http.StatusOK); err != nil {
		return nil, err
	}
	return decodeJSON[L](resp.Body)
}

// Create posts a new resource and returns the created resource response.
func (s *Service[T, L]) Create(ctx context.Context, resource *T) (*T, error) {
	if resource == nil {
		return nil, fmt.Errorf("resource is required")
	}

	body, err := json.Marshal(resource)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		fmt.Sprintf("%s/%s.json", s.client.baseURL, s.router.Create()),
		bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resp, err := s.client.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := s.client.checkResponse(ctx, resp, "Create", http.StatusOK, http.StatusCreated); err != nil {
		return nil, err
	}
	return decodeJSON[T](resp.Body)
}

// Update sends a PUT or PATCH request to update a resource by ID.
func (s *Service[T, L]) Update(ctx context.Context, id int64, resource *T) (*T, error) {
	if id <= 0 {
		return nil, fmt.Errorf("id must be greater than 0")
	}
	if resource == nil {
		return nil, fmt.Errorf("resource is required")
	}

	body, err := json.Marshal(resource)
	if err != nil {
		return nil, err
	}

	method := http.MethodPut
	if ump, ok := s.router.(updateMethodProvider); ok {
		method = ump.UpdateMethod()
	}

	req, err := http.NewRequestWithContext(ctx, method,
		fmt.Sprintf("%s/%s.json", s.client.baseURL, s.router.Update(id)),
		bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resp, err := s.client.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := s.client.checkResponse(ctx, resp, "Update", http.StatusOK); err != nil {
		return nil, err
	}
	return decodeJSON[T](resp.Body)
}
