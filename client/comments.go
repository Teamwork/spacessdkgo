package client

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/teamwork/spacessdkgo/models"
)

// CommentService provides methods for interacting with comments nested under spaces/pages.
type CommentService struct {
	client *Client
}

// NewCommentService creates a new CommentService.
func NewCommentService(client *Client) *CommentService {
	return &CommentService{client: client}
}

// Get fetches a single comment by space, page, and comment ID.
func (s *CommentService) Get(ctx context.Context, spaceID, pageID, commentID int64) (*models.CommentResponse, error) {
	if spaceID <= 0 {
		return nil, fmt.Errorf("spaceID must be greater than 0")
	}
	if pageID <= 0 {
		return nil, fmt.Errorf("pageID must be greater than 0")
	}
	if commentID <= 0 {
		return nil, fmt.Errorf("commentID must be greater than 0")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("%s/spaces/%d/pages/%d/comments/%d.json", s.client.baseURL, spaceID, pageID, commentID), nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := s.client.checkResponse(ctx, resp, "Get comment", http.StatusOK); err != nil {
		return nil, err
	}
	return decodeJSON[models.CommentResponse](resp.Body)
}

// List fetches all comments for a page within a space with optional query parameters.
func (s *CommentService) List(ctx context.Context, spaceID, pageID int64, params url.Values) (*models.CommentsResponse, error) {
	if spaceID <= 0 {
		return nil, fmt.Errorf("spaceID must be greater than 0")
	}
	if pageID <= 0 {
		return nil, fmt.Errorf("pageID must be greater than 0")
	}

	u := fmt.Sprintf("%s/spaces/%d/pages/%d/comments.json", s.client.baseURL, spaceID, pageID)
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

	if err := s.client.checkResponse(ctx, resp, "List comments", http.StatusOK); err != nil {
		return nil, err
	}
	return decodeJSON[models.CommentsResponse](resp.Body)
}

// Create creates a new comment on a page within a space.
func (s *CommentService) Create(ctx context.Context, spaceID, pageID int64, req *models.CommentCreate) (*models.CommentResponse, error) {
	if spaceID <= 0 {
		return nil, fmt.Errorf("spaceID must be greater than 0")
	}
	if pageID <= 0 {
		return nil, fmt.Errorf("pageID must be greater than 0")
	}
	if req == nil {
		return nil, fmt.Errorf("req is required")
	}

	body, err := wrapRequest("comment", req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost,
		fmt.Sprintf("%s/spaces/%d/pages/%d/comments.json", s.client.baseURL, spaceID, pageID),
		bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resp, err := s.client.doRequest(ctx, httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := s.client.checkResponse(ctx, resp, "Create comment", http.StatusOK, http.StatusCreated); err != nil {
		return nil, err
	}
	return decodeJSON[models.CommentResponse](resp.Body)
}

// Update updates an existing comment.
func (s *CommentService) Update(ctx context.Context, spaceID, pageID, commentID int64, req *models.CommentUpdate) (*models.CommentResponse, error) {
	if spaceID <= 0 {
		return nil, fmt.Errorf("spaceID must be greater than 0")
	}
	if pageID <= 0 {
		return nil, fmt.Errorf("pageID must be greater than 0")
	}
	if commentID <= 0 {
		return nil, fmt.Errorf("commentID must be greater than 0")
	}
	if req == nil {
		return nil, fmt.Errorf("req is required")
	}

	body, err := wrapRequest("comment", req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPatch,
		fmt.Sprintf("%s/spaces/%d/pages/%d/comments/%d.json", s.client.baseURL, spaceID, pageID, commentID),
		bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resp, err := s.client.doRequest(ctx, httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := s.client.checkResponse(ctx, resp, "Update comment", http.StatusOK); err != nil {
		return nil, err
	}
	return decodeJSON[models.CommentResponse](resp.Body)
}

// Delete removes a comment. Returns only an error (204 No Content on success).
func (s *CommentService) Delete(ctx context.Context, spaceID, pageID, commentID int64) error {
	if spaceID <= 0 {
		return fmt.Errorf("spaceID must be greater than 0")
	}
	if pageID <= 0 {
		return fmt.Errorf("pageID must be greater than 0")
	}
	if commentID <= 0 {
		return fmt.Errorf("commentID must be greater than 0")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete,
		fmt.Sprintf("%s/spaces/%d/pages/%d/comments/%d.json", s.client.baseURL, spaceID, pageID, commentID), nil)
	if err != nil {
		return err
	}

	resp, err := s.client.doRequest(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return s.client.checkResponse(ctx, resp, "Delete comment", http.StatusNoContent, http.StatusOK)
}
