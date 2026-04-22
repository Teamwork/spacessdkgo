package client

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/teamwork/spacessdkgo/models"
)

// TagService provides methods for interacting with the tags resource.
type TagService struct {
	*Service[models.TagResponse, models.TagsResponse]
	client *Client
}

// NewTagService creates a new TagService.
func NewTagService(client *Client) *TagService {
	return &TagService{
		Service: NewService[models.TagResponse, models.TagsResponse](
			client,
			NewDefaultPathHandlerWithUpdateMethod("tags", http.MethodPatch),
		),
		client: client,
	}
}

// Get fetches a single tag by ID.
func (s *TagService) Get(ctx context.Context, id int64) (*models.TagResponse, error) {
	return s.Service.Get(ctx, id)
}

// List fetches a collection of tags with optional query parameters.
func (s *TagService) List(ctx context.Context, params url.Values) (*models.TagsResponse, error) {
	return s.Service.List(ctx, params)
}

// CreateBatch creates one or more tags in a single request.
// The API wraps the tags array as {"tags": [...]}.
func (s *TagService) CreateBatch(ctx context.Context, tags []models.Tag) ([]models.Tag, error) {
	if len(tags) == 0 {
		return nil, fmt.Errorf("tags must not be empty")
	}

	body, err := wrapRequest("tags", tags)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		fmt.Sprintf("%s/tags.json", s.client.baseURL),
		bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resp, err := s.client.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := s.client.checkResponse(ctx, resp, "CreateBatch tags", http.StatusOK, http.StatusCreated); err != nil {
		return nil, err
	}
	result, err := decodeJSON[models.TagsResponse](resp.Body)
	if err != nil {
		return nil, err
	}
	return result.Tags, nil
}

// Update updates an existing tag by ID.
func (s *TagService) Update(ctx context.Context, id int64, req *models.TagUpdate) (*models.TagResponse, error) {
	if id <= 0 {
		return nil, fmt.Errorf("id must be greater than 0")
	}
	if req == nil {
		return nil, fmt.Errorf("req is required")
	}

	body, err := wrapRequest("tags", req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPatch,
		fmt.Sprintf("%s/tags/%d.json", s.client.baseURL, id),
		bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resp, err := s.client.doRequest(ctx, httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := s.client.checkResponse(ctx, resp, "Update tag", http.StatusOK); err != nil {
		return nil, err
	}
	return decodeJSON[models.TagResponse](resp.Body)
}

// Delete removes a tag by ID. Returns only an error (204 No Content on success).
func (s *TagService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("id must be greater than 0")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete,
		fmt.Sprintf("%s/tags/%d.json", s.client.baseURL, id), nil)
	if err != nil {
		return err
	}

	resp, err := s.client.doRequest(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return s.client.checkResponse(ctx, resp, "Delete tag", http.StatusNoContent, http.StatusOK)
}
