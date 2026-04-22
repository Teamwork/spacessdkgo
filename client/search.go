package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sonh/qs"
	"github.com/teamwork/spacessdkgo/models"
)

// SearchService provides methods for full-text search across spaces.
type SearchService struct {
	client *Client
}

// NewSearchService creates a new SearchService.
func NewSearchService(client *Client) *SearchService {
	return &SearchService{client: client}
}

// Search performs a full-text search using the provided filter parameters.
func (s *SearchService) Search(ctx context.Context, filter models.SearchFilter) (*models.SearchResponse, error) {
	encoder := qs.NewEncoder()
	values, err := encoder.Values(filter)
	if err != nil {
		return nil, fmt.Errorf("encoding search filter: %w", err)
	}

	u := fmt.Sprintf("%s/search.json", s.client.baseURL)
	if encoded := values.Encode(); encoded != "" {
		u = fmt.Sprintf("%s?%s", u, encoded)
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

	if err := s.client.checkResponse(ctx, resp, "Search", http.StatusOK); err != nil {
		return nil, err
	}
	return decodeJSON[models.SearchResponse](resp.Body)
}
