package client

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/teamwork/spacessdkgo/models"
)

func TestTagService_Get(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodGet, "/tags/5.json", http.StatusOK,
			models.TagResponse{Tag: models.Tag{ID: 5, Name: "golang", Color: "#00ADD8"}})

		got, err := newTestClient(mock).Tags.Get(context.Background(), 5)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Tag.ID != 5 {
			t.Errorf("got ID %d, want 5", got.Tag.ID)
		}
		if got.Tag.Name != "golang" {
			t.Errorf("got name %q, want %q", got.Tag.Name, "golang")
		}
		if mock.GetRequests()[0].URL.Path != "/tags/5.json" {
			t.Errorf("unexpected path %s", mock.GetRequests()[0].URL.Path)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		_, err := newTestClient(NewMockRoundTripper()).Tags.Get(context.Background(), 0)
		if err == nil {
			t.Fatal("expected error for id=0")
		}
	})

	t.Run("api error", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodGet, "/tags/5.json", http.StatusNotFound, `{"error":"not found"}`)
		_, err := newTestClient(mock).Tags.Get(context.Background(), 5)
		if err == nil {
			t.Fatal("expected error on 404")
		}
	})
}

func TestTagService_List(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodGet, "/tags.json", http.StatusOK,
			models.TagsResponse{Tags: []models.Tag{{ID: 1, Name: "go"}, {ID: 2, Name: "api"}}})

		got, err := newTestClient(mock).Tags.List(context.Background(), url.Values{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got.Tags) != 2 {
			t.Fatalf("got %d tags, want 2", len(got.Tags))
		}
	})

	t.Run("api error", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodGet, "/tags.json", http.StatusInternalServerError, `{"error":"internal"}`)
		_, err := newTestClient(mock).Tags.List(context.Background(), url.Values{})
		if err == nil {
			t.Fatal("expected error on 500")
		}
	})
}

func TestTagService_CreateBatch(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodPost, "/tags.json", http.StatusCreated,
			models.TagsResponse{Tags: []models.Tag{
				{ID: 10, Name: "go", Color: "#00ADD8"},
				{ID: 11, Name: "api", Color: "#FF6B6B"},
			}})

		got, err := newTestClient(mock).Tags.CreateBatch(context.Background(), []models.Tag{
			{Name: "go", Color: "#00ADD8"},
			{Name: "api", Color: "#FF6B6B"},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 2 {
			t.Fatalf("got %d tags, want 2", len(got))
		}
		if got[0].ID != 10 {
			t.Errorf("got ID %d, want 10", got[0].ID)
		}
		if mock.GetRequests()[0].Method != http.MethodPost {
			t.Errorf("expected POST, got %s", mock.GetRequests()[0].Method)
		}
	})

	t.Run("empty tags", func(t *testing.T) {
		_, err := newTestClient(NewMockRoundTripper()).Tags.CreateBatch(context.Background(), []models.Tag{})
		if err == nil {
			t.Fatal("expected error for empty tags")
		}
	})

	t.Run("nil tags", func(t *testing.T) {
		_, err := newTestClient(NewMockRoundTripper()).Tags.CreateBatch(context.Background(), nil)
		if err == nil {
			t.Fatal("expected error for nil tags")
		}
	})

	t.Run("api error", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodPost, "/tags.json", http.StatusBadRequest, `{"error":"bad"}`)
		_, err := newTestClient(mock).Tags.CreateBatch(context.Background(), []models.Tag{{Name: "x"}})
		if err == nil {
			t.Fatal("expected error on 400")
		}
	})
}

func TestTagService_Update(t *testing.T) {
	name := "updated-tag"

	t.Run("success", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodPatch, "/tags/5.json", http.StatusOK,
			models.TagResponse{Tag: models.Tag{ID: 5, Name: "updated-tag"}})

		got, err := newTestClient(mock).Tags.Update(context.Background(), 5,
			&models.TagUpdate{Name: &name})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Tag.Name != "updated-tag" {
			t.Errorf("got name %q, want %q", got.Tag.Name, "updated-tag")
		}
		if mock.GetRequests()[0].Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", mock.GetRequests()[0].Method)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		_, err := newTestClient(NewMockRoundTripper()).Tags.Update(context.Background(), 0,
			&models.TagUpdate{Name: &name})
		if err == nil {
			t.Fatal("expected error for id=0")
		}
	})

	t.Run("nil req", func(t *testing.T) {
		_, err := newTestClient(NewMockRoundTripper()).Tags.Update(context.Background(), 5, nil)
		if err == nil {
			t.Fatal("expected error for nil req")
		}
	})

	t.Run("api error", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodPatch, "/tags/5.json", http.StatusNotFound, `{"error":"not found"}`)
		_, err := newTestClient(mock).Tags.Update(context.Background(), 5, &models.TagUpdate{Name: &name})
		if err == nil {
			t.Fatal("expected error on 404")
		}
	})
}

func TestTagService_Delete(t *testing.T) {
	t.Run("success 204", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodDelete, "/tags/5.json", http.StatusNoContent, "")
		if err := newTestClient(mock).Tags.Delete(context.Background(), 5); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if mock.GetRequests()[0].Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", mock.GetRequests()[0].Method)
		}
	})

	t.Run("success 200", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodDelete, "/tags/5.json", http.StatusOK, "")
		if err := newTestClient(mock).Tags.Delete(context.Background(), 5); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		if err := newTestClient(NewMockRoundTripper()).Tags.Delete(context.Background(), 0); err == nil {
			t.Fatal("expected error for id=0")
		}
	})

	t.Run("api error", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodDelete, "/tags/5.json", http.StatusNotFound, `{"error":"not found"}`)
		if err := newTestClient(mock).Tags.Delete(context.Background(), 5); err == nil {
			t.Fatal("expected error on 404")
		}
	})
}
