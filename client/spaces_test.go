package client

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/teamwork/spacessdkgo/models"
)

func newTestClient(transport *MockRoundTripper) *Client {
	return NewClient("https://example.com",
		WithHTTPClient(&http.Client{Transport: transport}),
	)
}

func TestSpaceService_Get(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodGet, "/spaces/api/v1/spaces/42.json", http.StatusOK,
			models.SpaceResponse{Space: models.Space{ID: 42, Title: "Docs"}})

		got, err := newTestClient(mock).Spaces.Get(context.Background(), 42)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Space.ID != 42 {
			t.Errorf("got ID %d, want 42", got.Space.ID)
		}
		if got.Space.Title != "Docs" {
			t.Errorf("got title %q, want %q", got.Space.Title, "Docs")
		}
		reqs := mock.GetRequests()
		if reqs[0].Method != http.MethodGet {
			t.Errorf("expected GET, got %s", reqs[0].Method)
		}
		if reqs[0].URL.Path != "/spaces/api/v1/spaces/42.json" {
			t.Errorf("unexpected path %s", reqs[0].URL.Path)
		}
	})

	for _, id := range []int64{0, -1} {
		t.Run("invalid id", func(t *testing.T) {
			_, err := newTestClient(NewMockRoundTripper()).Spaces.Get(context.Background(), id)
			if err == nil {
				t.Fatalf("expected error for id=%d", id)
			}
		})
	}

	t.Run("api error", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodGet, "/spaces/api/v1/spaces/1.json", http.StatusNotFound, `{"message":"not found"}`)
		_, err := newTestClient(mock).Spaces.Get(context.Background(), 1)
		if err == nil {
			t.Fatal("expected error on 404")
		}
	})
}

func TestSpaceService_List(t *testing.T) {
	t.Run("success no params", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodGet, "/spaces/api/v1/spaces.json", http.StatusOK,
			models.SpacesResponse{Spaces: []models.Space{{ID: 1, Title: "First"}}})

		got, err := newTestClient(mock).Spaces.List(context.Background(), url.Values{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got.Spaces) != 1 {
			t.Fatalf("got %d spaces, want 1", len(got.Spaces))
		}
		if got.Spaces[0].Title != "First" {
			t.Errorf("got title %q, want %q", got.Spaces[0].Title, "First")
		}
	})

	t.Run("passes query params", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodGet, "/spaces/api/v1/spaces.json", http.StatusOK, models.SpacesResponse{})

		params := url.Values{"page": {"2"}, "perPage": {"10"}}
		_, err := newTestClient(mock).Spaces.List(context.Background(), params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if q := mock.GetRequests()[0].URL.Query().Get("page"); q != "2" {
			t.Errorf("expected page=2, got %q", q)
		}
	})

	t.Run("api error", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodGet, "/spaces/api/v1/spaces.json",
			http.StatusInternalServerError, `{"error":"internal"}`)
		_, err := newTestClient(mock).Spaces.List(context.Background(), url.Values{})
		if err == nil {
			t.Fatal("expected error on 500")
		}
	})
}

func TestSpaceService_Create(t *testing.T) {
	t.Run("success 201", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodPost, "/spaces/api/v1/spaces.json", http.StatusCreated,
			models.SpaceResponse{Space: models.Space{ID: 99, Title: "New Space"}})

		got, err := newTestClient(mock).Spaces.Create(context.Background(),
			&models.SpaceCreate{Title: "New Space", Code: "new"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Space.ID != 99 {
			t.Errorf("got ID %d, want 99", got.Space.ID)
		}
		if mock.GetRequests()[0].Method != http.MethodPost {
			t.Errorf("expected POST, got %s", mock.GetRequests()[0].Method)
		}
	})

	t.Run("success 200", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodPost, "/spaces/api/v1/spaces.json", http.StatusOK,
			models.SpaceResponse{Space: models.Space{ID: 1}})
		_, err := newTestClient(mock).Spaces.Create(context.Background(),
			&models.SpaceCreate{Title: "X", Code: "x"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("nil req", func(t *testing.T) {
		_, err := newTestClient(NewMockRoundTripper()).Spaces.Create(context.Background(), nil)
		if err == nil {
			t.Fatal("expected error for nil req")
		}
	})

	t.Run("api error", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodPost, "/spaces/api/v1/spaces.json", http.StatusBadRequest, `{"error":"bad request"}`)
		_, err := newTestClient(mock).Spaces.Create(context.Background(),
			&models.SpaceCreate{Title: "X", Code: "x"})
		if err == nil {
			t.Fatal("expected error on 400")
		}
	})
}

func TestSpaceService_Update(t *testing.T) {
	title := "Updated"

	t.Run("success", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodPatch, "/spaces/api/v1/spaces/1.json", http.StatusOK,
			models.SpaceResponse{Space: models.Space{ID: 1, Title: "Updated"}})

		got, err := newTestClient(mock).Spaces.Update(context.Background(), 1,
			&models.SpaceUpdate{Title: &title})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Space.Title != "Updated" {
			t.Errorf("got title %q, want %q", got.Space.Title, "Updated")
		}
		if mock.GetRequests()[0].Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", mock.GetRequests()[0].Method)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		_, err := newTestClient(NewMockRoundTripper()).Spaces.Update(context.Background(), 0,
			&models.SpaceUpdate{Title: &title})
		if err == nil {
			t.Fatal("expected error for id=0")
		}
	})

	t.Run("nil req", func(t *testing.T) {
		_, err := newTestClient(NewMockRoundTripper()).Spaces.Update(context.Background(), 1, nil)
		if err == nil {
			t.Fatal("expected error for nil req")
		}
	})

	t.Run("api error", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodPatch, "/spaces/api/v1/spaces/1.json", http.StatusForbidden, `{"error":"forbidden"}`)
		_, err := newTestClient(mock).Spaces.Update(context.Background(), 1,
			&models.SpaceUpdate{Title: &title})
		if err == nil {
			t.Fatal("expected error on 403")
		}
	})
}

func TestSpaceService_Delete(t *testing.T) {
	t.Run("success 204", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodDelete, "/spaces/api/v1/spaces/1.json", http.StatusNoContent, "")
		if err := newTestClient(mock).Spaces.Delete(context.Background(), 1); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if mock.GetRequests()[0].Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", mock.GetRequests()[0].Method)
		}
	})

	t.Run("success 200", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodDelete, "/spaces/api/v1/spaces/2.json", http.StatusOK, "")
		if err := newTestClient(mock).Spaces.Delete(context.Background(), 2); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		if err := newTestClient(NewMockRoundTripper()).Spaces.Delete(context.Background(), 0); err == nil {
			t.Fatal("expected error for id=0")
		}
	})

	t.Run("api error", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodDelete, "/spaces/api/v1/spaces/1.json", http.StatusNotFound, `{"error":"not found"}`)
		if err := newTestClient(mock).Spaces.Delete(context.Background(), 1); err == nil {
			t.Fatal("expected error on 404")
		}
	})
}

func TestSpaceService_Collaborators(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodGet, "/spaces/api/v1/spaces/5/collaborators.json", http.StatusOK,
			models.SpaceCollaboratorsResponse{
				CollaboratorsCount: 2,
				Collaborators: []models.SpaceCollaborator{
					{ID: 10, Type: "user"},
					{ID: 11, Type: "user"},
				},
			})

		got, err := newTestClient(mock).Spaces.Collaborators(context.Background(), 5)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.CollaboratorsCount != 2 {
			t.Errorf("got count %d, want 2", got.CollaboratorsCount)
		}
		if len(got.Collaborators) != 2 {
			t.Errorf("got %d collaborators, want 2", len(got.Collaborators))
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		_, err := newTestClient(NewMockRoundTripper()).Spaces.Collaborators(context.Background(), 0)
		if err == nil {
			t.Fatal("expected error for id=0")
		}
	})

	t.Run("api error", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodGet, "/spaces/api/v1/spaces/5/collaborators.json",
			http.StatusNotFound, `{"error":"not found"}`)
		_, err := newTestClient(mock).Spaces.Collaborators(context.Background(), 5)
		if err == nil {
			t.Fatal("expected error on 404")
		}
	})
}
