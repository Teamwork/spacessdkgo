package client

import (
	"context"
	"net/http"
	"testing"

	"github.com/teamwork/spacessdkgo/models"
)

func TestSearchService_Search(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodGet, "/spaces/api/v1/search.json", http.StatusOK,
			models.SearchResponse{
				TotalResults: 2,
				Results: []models.SearchPage{
					{PageID: 1, Title: "Getting Started", Slug: "getting-started"},
					{PageID: 2, Title: "API Reference", Slug: "api-reference"},
				},
			})

		got, err := newTestClient(mock).Search.Search(context.Background(),
			models.SearchFilter{Query: "getting started"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.TotalResults != 2 {
			t.Errorf("got TotalResults %d, want 2", got.TotalResults)
		}
		if len(got.Results) != 2 {
			t.Fatalf("got %d results, want 2", len(got.Results))
		}
		if got.Results[0].Title != "Getting Started" {
			t.Errorf("got title %q, want %q", got.Results[0].Title, "Getting Started")
		}
	})

	t.Run("encodes query params", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodGet, "/spaces/api/v1/search.json", http.StatusOK, models.SearchResponse{})

		limit := int64(10)
		_, err := newTestClient(mock).Search.Search(context.Background(), models.SearchFilter{
			Query:   "deploy guide",
			SpaceID: []int64{42, 57},
			Limit:   &limit,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		reqs := mock.GetRequests()
		q := reqs[0].URL.Query()
		if q.Get("q") != "deploy guide" {
			t.Errorf("expected q=%q, got %q", "deploy guide", q.Get("q"))
		}
		if q.Get("limit") != "10" {
			t.Errorf("expected limit=10, got %q", q.Get("limit"))
		}
		spaceIDs := q["spaceid"]
		if len(spaceIDs) != 2 {
			t.Errorf("expected 2 spaceid values, got %d", len(spaceIDs))
		}
	})

	t.Run("empty filter", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodGet, "/spaces/api/v1/search.json", http.StatusOK, models.SearchResponse{})

		_, err := newTestClient(mock).Search.Search(context.Background(), models.SearchFilter{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if mock.GetRequests()[0].Method != http.MethodGet {
			t.Errorf("expected GET, got %s", mock.GetRequests()[0].Method)
		}
	})

	t.Run("api error", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodGet, "/spaces/api/v1/search.json",
			http.StatusInternalServerError, `{"error":"internal"}`)
		_, err := newTestClient(mock).Search.Search(context.Background(),
			models.SearchFilter{Query: "test"})
		if err == nil {
			t.Fatal("expected error on 500")
		}
	})
}
