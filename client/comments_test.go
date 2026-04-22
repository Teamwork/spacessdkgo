package client

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/teamwork/spacessdkgo/models"
)

func TestCommentService_Get(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodGet, "/spaces/10/pages/20/comments/30.json", http.StatusOK,
			models.CommentResponse{Comment: models.Comment{ID: 30, Content: "Hello"}})

		got, err := newTestClient(mock).Comments.Get(context.Background(), 10, 20, 30)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Comment.ID != 30 {
			t.Errorf("got ID %d, want 30", got.Comment.ID)
		}
		if got.Comment.Content != "Hello" {
			t.Errorf("got content %q, want %q", got.Comment.Content, "Hello")
		}
		if mock.GetRequests()[0].URL.Path != "/spaces/10/pages/20/comments/30.json" {
			t.Errorf("unexpected path %s", mock.GetRequests()[0].URL.Path)
		}
	})

	t.Run("invalid spaceID", func(t *testing.T) {
		_, err := newTestClient(NewMockRoundTripper()).Comments.Get(context.Background(), 0, 20, 30)
		if err == nil {
			t.Fatal("expected error for spaceID=0")
		}
	})

	t.Run("invalid pageID", func(t *testing.T) {
		_, err := newTestClient(NewMockRoundTripper()).Comments.Get(context.Background(), 10, 0, 30)
		if err == nil {
			t.Fatal("expected error for pageID=0")
		}
	})

	t.Run("invalid commentID", func(t *testing.T) {
		_, err := newTestClient(NewMockRoundTripper()).Comments.Get(context.Background(), 10, 20, 0)
		if err == nil {
			t.Fatal("expected error for commentID=0")
		}
	})

	t.Run("api error", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodGet, "/spaces/10/pages/20/comments/30.json", http.StatusNotFound, `{"error":"not found"}`)
		_, err := newTestClient(mock).Comments.Get(context.Background(), 10, 20, 30)
		if err == nil {
			t.Fatal("expected error on 404")
		}
	})
}

func TestCommentService_List(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodGet, "/spaces/10/pages/20/comments.json", http.StatusOK,
			models.CommentsResponse{
				Comments: []models.CommentParent{
					{Comment: models.Comment{ID: 1, Content: "First"}},
					{Comment: models.Comment{ID: 2, Content: "Second"}},
				},
			})

		got, err := newTestClient(mock).Comments.List(context.Background(), 10, 20, url.Values{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got.Comments) != 2 {
			t.Fatalf("got %d comments, want 2", len(got.Comments))
		}
	})

	t.Run("passes query params", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodGet, "/spaces/10/pages/20/comments.json", http.StatusOK,
			models.CommentsResponse{})

		_, err := newTestClient(mock).Comments.List(context.Background(), 10, 20,
			url.Values{"page": {"2"}})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if q := mock.GetRequests()[0].URL.Query().Get("page"); q != "2" {
			t.Errorf("expected page=2, got %q", q)
		}
	})

	t.Run("invalid spaceID", func(t *testing.T) {
		_, err := newTestClient(NewMockRoundTripper()).Comments.List(context.Background(), 0, 20, url.Values{})
		if err == nil {
			t.Fatal("expected error for spaceID=0")
		}
	})

	t.Run("invalid pageID", func(t *testing.T) {
		_, err := newTestClient(NewMockRoundTripper()).Comments.List(context.Background(), 10, 0, url.Values{})
		if err == nil {
			t.Fatal("expected error for pageID=0")
		}
	})

	t.Run("api error", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodGet, "/spaces/10/pages/20/comments.json", http.StatusInternalServerError, `{"error":"internal"}`)
		_, err := newTestClient(mock).Comments.List(context.Background(), 10, 20, url.Values{})
		if err == nil {
			t.Fatal("expected error on 500")
		}
	})
}

func TestCommentService_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodPost, "/spaces/10/pages/20/comments.json", http.StatusCreated,
			models.CommentResponse{Comment: models.Comment{ID: 77, Content: "New comment"}})

		got, err := newTestClient(mock).Comments.Create(context.Background(), 10, 20,
			&models.CommentCreate{Content: "New comment"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Comment.ID != 77 {
			t.Errorf("got ID %d, want 77", got.Comment.ID)
		}
		if mock.GetRequests()[0].Method != http.MethodPost {
			t.Errorf("expected POST, got %s", mock.GetRequests()[0].Method)
		}
	})

	t.Run("invalid spaceID", func(t *testing.T) {
		_, err := newTestClient(NewMockRoundTripper()).Comments.Create(context.Background(), 0, 20,
			&models.CommentCreate{Content: "X"})
		if err == nil {
			t.Fatal("expected error for spaceID=0")
		}
	})

	t.Run("invalid pageID", func(t *testing.T) {
		_, err := newTestClient(NewMockRoundTripper()).Comments.Create(context.Background(), 10, 0,
			&models.CommentCreate{Content: "X"})
		if err == nil {
			t.Fatal("expected error for pageID=0")
		}
	})

	t.Run("nil req", func(t *testing.T) {
		_, err := newTestClient(NewMockRoundTripper()).Comments.Create(context.Background(), 10, 20, nil)
		if err == nil {
			t.Fatal("expected error for nil req")
		}
	})

	t.Run("api error", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodPost, "/spaces/10/pages/20/comments.json", http.StatusBadRequest, `{"error":"bad"}`)
		_, err := newTestClient(mock).Comments.Create(context.Background(), 10, 20,
			&models.CommentCreate{Content: "X"})
		if err == nil {
			t.Fatal("expected error on 400")
		}
	})
}

func TestCommentService_Update(t *testing.T) {
	content := "Updated content"

	t.Run("success", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodPatch, "/spaces/10/pages/20/comments/30.json", http.StatusOK,
			models.CommentResponse{Comment: models.Comment{ID: 30, Content: "Updated content"}})

		got, err := newTestClient(mock).Comments.Update(context.Background(), 10, 20, 30,
			&models.CommentUpdate{Content: &content})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Comment.Content != "Updated content" {
			t.Errorf("got content %q, want %q", got.Comment.Content, "Updated content")
		}
		if mock.GetRequests()[0].Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", mock.GetRequests()[0].Method)
		}
	})

	t.Run("invalid spaceID", func(t *testing.T) {
		_, err := newTestClient(NewMockRoundTripper()).Comments.Update(context.Background(), 0, 20, 30,
			&models.CommentUpdate{Content: &content})
		if err == nil {
			t.Fatal("expected error for spaceID=0")
		}
	})

	t.Run("invalid pageID", func(t *testing.T) {
		_, err := newTestClient(NewMockRoundTripper()).Comments.Update(context.Background(), 10, 0, 30,
			&models.CommentUpdate{Content: &content})
		if err == nil {
			t.Fatal("expected error for pageID=0")
		}
	})

	t.Run("invalid commentID", func(t *testing.T) {
		_, err := newTestClient(NewMockRoundTripper()).Comments.Update(context.Background(), 10, 20, 0,
			&models.CommentUpdate{Content: &content})
		if err == nil {
			t.Fatal("expected error for commentID=0")
		}
	})

	t.Run("nil req", func(t *testing.T) {
		_, err := newTestClient(NewMockRoundTripper()).Comments.Update(context.Background(), 10, 20, 30, nil)
		if err == nil {
			t.Fatal("expected error for nil req")
		}
	})

	t.Run("api error", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodPatch, "/spaces/10/pages/20/comments/30.json", http.StatusForbidden, `{"error":"forbidden"}`)
		_, err := newTestClient(mock).Comments.Update(context.Background(), 10, 20, 30,
			&models.CommentUpdate{Content: &content})
		if err == nil {
			t.Fatal("expected error on 403")
		}
	})
}

func TestCommentService_Delete(t *testing.T) {
	t.Run("success 204", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodDelete, "/spaces/10/pages/20/comments/30.json", http.StatusNoContent, "")
		if err := newTestClient(mock).Comments.Delete(context.Background(), 10, 20, 30); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if mock.GetRequests()[0].Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", mock.GetRequests()[0].Method)
		}
	})

	t.Run("success 200", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodDelete, "/spaces/10/pages/20/comments/30.json", http.StatusOK, "")
		if err := newTestClient(mock).Comments.Delete(context.Background(), 10, 20, 30); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("invalid spaceID", func(t *testing.T) {
		if err := newTestClient(NewMockRoundTripper()).Comments.Delete(context.Background(), 0, 20, 30); err == nil {
			t.Fatal("expected error for spaceID=0")
		}
	})

	t.Run("invalid pageID", func(t *testing.T) {
		if err := newTestClient(NewMockRoundTripper()).Comments.Delete(context.Background(), 10, 0, 30); err == nil {
			t.Fatal("expected error for pageID=0")
		}
	})

	t.Run("invalid commentID", func(t *testing.T) {
		if err := newTestClient(NewMockRoundTripper()).Comments.Delete(context.Background(), 10, 20, 0); err == nil {
			t.Fatal("expected error for commentID=0")
		}
	})

	t.Run("api error", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodDelete, "/spaces/10/pages/20/comments/30.json", http.StatusNotFound, `{"error":"not found"}`)
		if err := newTestClient(mock).Comments.Delete(context.Background(), 10, 20, 30); err == nil {
			t.Fatal("expected error on 404")
		}
	})
}
