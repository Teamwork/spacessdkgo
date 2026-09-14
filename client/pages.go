package client

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/teamwork/spacessdkgo/models"
)

// PageService provides methods for interacting with pages nested under spaces.
type PageService struct {
	client *Client
}

// NewPageService creates a new PageService.
func NewPageService(client *Client) *PageService {
	return &PageService{client: client}
}

// Get fetches a single page by space ID and page ID.
func (s *PageService) Get(ctx context.Context, spaceID, pageID int64) (*models.PageResponse, error) {
	if spaceID <= 0 {
		return nil, fmt.Errorf("spaceID must be greater than 0")
	}
	if pageID <= 0 {
		return nil, fmt.Errorf("pageID must be greater than 0")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("%s/spaces/%d/pages/%d.json", s.client.baseURL, spaceID, pageID), nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if err := s.client.checkResponse(ctx, resp, "Get page", http.StatusOK); err != nil {
		return nil, err
	}
	return decodeJSON[models.PageResponse](resp.Body)
}

// List fetches all pages within a space with optional query parameters.
func (s *PageService) List(ctx context.Context, spaceID int64, params url.Values) (*models.PagesResponse, error) {
	if spaceID <= 0 {
		return nil, fmt.Errorf("spaceID must be greater than 0")
	}

	u := fmt.Sprintf("%s/spaces/%d/pages.json", s.client.baseURL, spaceID)
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
	defer func() { _ = resp.Body.Close() }()

	if err := s.client.checkResponse(ctx, resp, "List pages", http.StatusOK); err != nil {
		return nil, err
	}
	return decodeJSON[models.PagesResponse](resp.Body)
}

// ListWithPrivate fetches a space's pages from
// GET /spaces/api/v2/spaces/{spaceId}/pages.json, which answers with the open
// page tree under "pages" and, under "private", the restricted pages the
// calling user has access to. List reads the v1 route, which returns the open
// tree only.
//
// Both routes run the same query and take the same parameters; the v1 response
// shape simply leaves the private tree out.
func (s *PageService) ListWithPrivate(
	ctx context.Context,
	spaceID int64,
	params url.Values,
) (*models.SpaceContentResponse, error) {
	if spaceID <= 0 {
		return nil, fmt.Errorf("spaceID must be greater than 0")
	}

	u := fmt.Sprintf("%s/spaces/%d/pages.json", s.client.baseURLV2(), spaceID)
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
	defer func() { _ = resp.Body.Close() }()

	if err := s.client.checkResponse(ctx, resp, "List pages", http.StatusOK); err != nil {
		return nil, err
	}
	return decodeJSON[models.SpaceContentResponse](resp.Body)
}

// Home fetches the homepage of a space.
func (s *PageService) Home(ctx context.Context, spaceID int64) (*models.PageResponse, error) {
	if spaceID <= 0 {
		return nil, fmt.Errorf("spaceID must be greater than 0")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("%s/spaces/%d/homepage.json", s.client.baseURL, spaceID), nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if err := s.client.checkResponse(ctx, resp, "Home page", http.StatusOK); err != nil {
		return nil, err
	}
	return decodeJSON[models.PageResponse](resp.Body)
}

// Create creates a new page within a space.
func (s *PageService) Create(ctx context.Context, spaceID int64, req *models.PageCreate) (*models.PageResponse, error) {
	if spaceID <= 0 {
		return nil, fmt.Errorf("spaceID must be greater than 0")
	}
	if req == nil {
		return nil, fmt.Errorf("req is required")
	}

	body, err := wrapRequest("page", req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost,
		fmt.Sprintf("%s/spaces/%d/pages.json", s.client.baseURL, spaceID),
		bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resp, err := s.client.doRequest(ctx, httpReq)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if err := s.client.checkResponse(ctx, resp, "Create page", http.StatusOK, http.StatusCreated); err != nil {
		return nil, err
	}
	return decodeJSON[models.PageResponse](resp.Body)
}

// Duplicate duplicates a page within a space.
func (s *PageService) Duplicate(
	ctx context.Context,
	spaceID, pageID int64,
	req *models.PageDuplicate,
) (*models.PageResponse, error) {
	if spaceID <= 0 {
		return nil, fmt.Errorf("spaceID must be greater than 0")
	}
	if pageID <= 0 {
		return nil, fmt.Errorf("pageID must be greater than 0")
	}
	if req == nil {
		return nil, fmt.Errorf("req is required")
	}

	body, err := wrapRequest("page", req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost,
		fmt.Sprintf("%s/spaces/%d/pages/%d/duplicate.json", s.client.baseURL, spaceID, pageID),
		bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resp, err := s.client.doRequest(ctx, httpReq)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if err := s.client.checkResponse(ctx, resp, "Duplicate page", http.StatusOK, http.StatusCreated); err != nil {
		return nil, err
	}
	return decodeJSON[models.PageResponse](resp.Body)
}

// Update updates an existing page within a space.
func (s *PageService) Update(
	ctx context.Context,
	spaceID, pageID int64,
	req *models.PageUpdate,
) (*models.PageResponse, error) {
	if spaceID <= 0 {
		return nil, fmt.Errorf("spaceID must be greater than 0")
	}
	if pageID <= 0 {
		return nil, fmt.Errorf("pageID must be greater than 0")
	}
	if req == nil {
		return nil, fmt.Errorf("req is required")
	}

	body, err := wrapRequest("page", req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPatch,
		fmt.Sprintf("%s/spaces/%d/pages/%d.json", s.client.baseURL, spaceID, pageID),
		bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resp, err := s.client.doRequest(ctx, httpReq)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if err := s.client.checkResponse(ctx, resp, "Update page", http.StatusOK); err != nil {
		return nil, err
	}
	return decodeJSON[models.PageResponse](resp.Body)
}

// Delete removes a page from a space. Returns only an error (204 No Content on success).
func (s *PageService) Delete(ctx context.Context, spaceID, pageID int64) error {
	if spaceID <= 0 {
		return fmt.Errorf("spaceID must be greater than 0")
	}
	if pageID <= 0 {
		return fmt.Errorf("pageID must be greater than 0")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete,
		fmt.Sprintf("%s/spaces/%d/pages/%d.json", s.client.baseURL, spaceID, pageID), nil)
	if err != nil {
		return err
	}

	resp, err := s.client.doRequest(ctx, req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	return s.client.checkResponse(ctx, resp, "Delete page", http.StatusNoContent, http.StatusOK)
}
