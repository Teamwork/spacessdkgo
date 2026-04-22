package client

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/teamwork/spacessdkgo/models"
)

// SpaceService provides methods for interacting with the spaces resource.
type SpaceService struct {
	*Service[models.SpaceResponse, models.SpacesResponse]
	client *Client
}

// NewSpaceService creates a new SpaceService.
func NewSpaceService(client *Client) *SpaceService {
	return &SpaceService{
		Service: NewService[models.SpaceResponse, models.SpacesResponse](
			client,
			NewDefaultPathHandlerWithUpdateMethod("spaces", http.MethodPatch),
		),
		client: client,
	}
}

// Get fetches a single space by ID.
func (s *SpaceService) Get(ctx context.Context, id int64) (*models.SpaceResponse, error) {
	return s.Service.Get(ctx, id)
}

// List fetches a collection of spaces with optional query parameters.
func (s *SpaceService) List(ctx context.Context, params url.Values) (*models.SpacesResponse, error) {
	return s.Service.List(ctx, params)
}

// Create creates a new space.
func (s *SpaceService) Create(ctx context.Context, req *models.SpaceCreate) (*models.SpaceResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("req is required")
	}

	body, err := wrapRequest("space", req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost,
		fmt.Sprintf("%s/spaces.json", s.client.baseURL),
		bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resp, err := s.client.doRequest(ctx, httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := s.client.checkResponse(ctx, resp, "Create space", http.StatusOK, http.StatusCreated); err != nil {
		return nil, err
	}
	return decodeJSON[models.SpaceResponse](resp.Body)
}

// Update updates an existing space by ID.
func (s *SpaceService) Update(ctx context.Context, id int64, req *models.SpaceUpdate) (*models.SpaceResponse, error) {
	if id <= 0 {
		return nil, fmt.Errorf("id must be greater than 0")
	}
	if req == nil {
		return nil, fmt.Errorf("req is required")
	}

	body, err := wrapRequest("space", req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPatch,
		fmt.Sprintf("%s/spaces/%d.json", s.client.baseURL, id),
		bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resp, err := s.client.doRequest(ctx, httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := s.client.checkResponse(ctx, resp, "Update space", http.StatusOK); err != nil {
		return nil, err
	}
	return decodeJSON[models.SpaceResponse](resp.Body)
}

// Delete removes a space by ID. Returns only an error (204 No Content on success).
func (s *SpaceService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("id must be greater than 0")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete,
		fmt.Sprintf("%s/spaces/%d.json", s.client.baseURL, id), nil)
	if err != nil {
		return err
	}

	resp, err := s.client.doRequest(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return s.client.checkResponse(ctx, resp, "Delete space", http.StatusNoContent, http.StatusOK)
}

// Collaborators returns the list of collaborators for a space.
func (s *SpaceService) Collaborators(ctx context.Context, id int64) (*models.SpaceCollaboratorsResponse, error) {
	if id <= 0 {
		return nil, fmt.Errorf("id must be greater than 0")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("%s/spaces/%d/collaborators.json", s.client.baseURL, id), nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := s.client.checkResponse(ctx, resp, "Collaborators", http.StatusOK); err != nil {
		return nil, err
	}
	return decodeJSON[models.SpaceCollaboratorsResponse](resp.Body)
}
