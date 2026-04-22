package client

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/teamwork/spacessdkgo/models"
)

// CategoryService provides methods for interacting with the categories resource.
type CategoryService struct {
	*Service[models.CategoryResponse, models.CategoriesResponse]
	client *Client
}

// NewCategoryService creates a new CategoryService.
func NewCategoryService(client *Client) *CategoryService {
	return &CategoryService{
		Service: NewService[models.CategoryResponse, models.CategoriesResponse](
			client,
			NewDefaultPathHandlerWithUpdateMethod("categories", http.MethodPatch),
		),
		client: client,
	}
}

// Get fetches a single category by ID.
func (s *CategoryService) Get(ctx context.Context, id int64) (*models.CategoryResponse, error) {
	return s.Service.Get(ctx, id)
}

// List fetches a collection of categories with optional query parameters.
func (s *CategoryService) List(ctx context.Context, params url.Values) (*models.CategoriesResponse, error) {
	return s.Service.List(ctx, params)
}

// Create creates a new category.
func (s *CategoryService) Create(ctx context.Context, req *models.CategoryCreate) (*models.CategoryResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("req is required")
	}

	body, err := wrapRequest("category", req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost,
		fmt.Sprintf("%s/categories.json", s.client.baseURL),
		bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resp, err := s.client.doRequest(ctx, httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := s.client.checkResponse(ctx, resp, "Create category", http.StatusOK, http.StatusCreated); err != nil {
		return nil, err
	}
	return decodeJSON[models.CategoryResponse](resp.Body)
}

// Update updates an existing category by ID.
func (s *CategoryService) Update(ctx context.Context, id int64, req *models.CategoryUpdate) (*models.CategoryResponse, error) {
	if id <= 0 {
		return nil, fmt.Errorf("id must be greater than 0")
	}
	if req == nil {
		return nil, fmt.Errorf("req is required")
	}

	body, err := wrapRequest("category", req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPatch,
		fmt.Sprintf("%s/categories/%d.json", s.client.baseURL, id),
		bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resp, err := s.client.doRequest(ctx, httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := s.client.checkResponse(ctx, resp, "Update category", http.StatusOK); err != nil {
		return nil, err
	}
	return decodeJSON[models.CategoryResponse](resp.Body)
}

// Delete removes a category by ID. Returns only an error (204 No Content on success).
func (s *CategoryService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("id must be greater than 0")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete,
		fmt.Sprintf("%s/categories/%d.json", s.client.baseURL, id), nil)
	if err != nil {
		return err
	}

	resp, err := s.client.doRequest(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return s.client.checkResponse(ctx, resp, "Delete category", http.StatusNoContent, http.StatusOK)
}
