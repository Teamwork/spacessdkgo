package client

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/teamwork/spacessdkgo/models"
)

func TestCategoryService_Get(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodGet, "/spaces/api/v1/categories/3.json", http.StatusOK,
			models.CategoryResponse{Category: models.Category{ID: 3, Name: "Engineering"}})

		got, err := newTestClient(mock).Categories.Get(context.Background(), 3)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Category.ID != 3 {
			t.Errorf("got ID %d, want 3", got.Category.ID)
		}
		if got.Category.Name != "Engineering" {
			t.Errorf("got name %q, want %q", got.Category.Name, "Engineering")
		}
		if mock.GetRequests()[0].URL.Path != "/spaces/api/v1/categories/3.json" {
			t.Errorf("unexpected path %s", mock.GetRequests()[0].URL.Path)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		_, err := newTestClient(NewMockRoundTripper()).Categories.Get(context.Background(), 0)
		if err == nil {
			t.Fatal("expected error for id=0")
		}
	})

	t.Run("api error", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodGet, "/spaces/api/v1/categories/3.json", http.StatusNotFound, `{"error":"not found"}`)
		_, err := newTestClient(mock).Categories.Get(context.Background(), 3)
		if err == nil {
			t.Fatal("expected error on 404")
		}
	})
}

func TestCategoryService_List(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodGet, "/spaces/api/v1/categories.json", http.StatusOK,
			models.CategoriesResponse{Categories: []models.Category{
				{ID: 1, Name: "Engineering"},
				{ID: 2, Name: "Product"},
			}})

		got, err := newTestClient(mock).Categories.List(context.Background(), url.Values{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got.Categories) != 2 {
			t.Fatalf("got %d categories, want 2", len(got.Categories))
		}
	})

	t.Run("api error", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodGet, "/spaces/api/v1/categories.json",
			http.StatusInternalServerError, `{"error":"internal"}`)
		_, err := newTestClient(mock).Categories.List(context.Background(), url.Values{})
		if err == nil {
			t.Fatal("expected error on 500")
		}
	})
}

func TestCategoryService_Create(t *testing.T) {
	color := "#4A90E2"

	t.Run("success", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodPost, "/spaces/api/v1/categories.json", http.StatusCreated,
			models.CategoryResponse{Category: models.Category{ID: 7, Name: "Engineering"}})

		got, err := newTestClient(mock).Categories.Create(context.Background(),
			&models.CategoryCreate{Name: "Engineering", Color: &color})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Category.ID != 7 {
			t.Errorf("got ID %d, want 7", got.Category.ID)
		}
		if mock.GetRequests()[0].Method != http.MethodPost {
			t.Errorf("expected POST, got %s", mock.GetRequests()[0].Method)
		}
	})

	t.Run("success 200", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodPost, "/spaces/api/v1/categories.json", http.StatusOK,
			models.CategoryResponse{Category: models.Category{ID: 1}})
		_, err := newTestClient(mock).Categories.Create(context.Background(),
			&models.CategoryCreate{Name: "X"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("nil req", func(t *testing.T) {
		_, err := newTestClient(NewMockRoundTripper()).Categories.Create(context.Background(), nil)
		if err == nil {
			t.Fatal("expected error for nil req")
		}
	})

	t.Run("api error", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodPost, "/spaces/api/v1/categories.json", http.StatusBadRequest, `{"error":"bad"}`)
		_, err := newTestClient(mock).Categories.Create(context.Background(),
			&models.CategoryCreate{Name: "X"})
		if err == nil {
			t.Fatal("expected error on 400")
		}
	})
}

func TestCategoryService_Update(t *testing.T) {
	name := "Updated Engineering"

	t.Run("success", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodPatch, "/spaces/api/v1/categories/3.json", http.StatusOK,
			models.CategoryResponse{Category: models.Category{ID: 3, Name: "Updated Engineering"}})

		got, err := newTestClient(mock).Categories.Update(context.Background(), 3,
			&models.CategoryUpdate{Name: &name})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Category.Name != "Updated Engineering" {
			t.Errorf("got name %q, want %q", got.Category.Name, "Updated Engineering")
		}
		if mock.GetRequests()[0].Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", mock.GetRequests()[0].Method)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		_, err := newTestClient(NewMockRoundTripper()).Categories.Update(context.Background(), 0,
			&models.CategoryUpdate{Name: &name})
		if err == nil {
			t.Fatal("expected error for id=0")
		}
	})

	t.Run("nil req", func(t *testing.T) {
		_, err := newTestClient(NewMockRoundTripper()).Categories.Update(context.Background(), 3, nil)
		if err == nil {
			t.Fatal("expected error for nil req")
		}
	})

	t.Run("api error", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodPatch, "/spaces/api/v1/categories/3.json",
			http.StatusForbidden, `{"error":"forbidden"}`)
		_, err := newTestClient(mock).Categories.Update(context.Background(), 3,
			&models.CategoryUpdate{Name: &name})
		if err == nil {
			t.Fatal("expected error on 403")
		}
	})
}

func TestCategoryService_Delete(t *testing.T) {
	t.Run("success 204", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodDelete, "/spaces/api/v1/categories/3.json", http.StatusNoContent, "")
		if err := newTestClient(mock).Categories.Delete(context.Background(), 3); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if mock.GetRequests()[0].Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", mock.GetRequests()[0].Method)
		}
	})

	t.Run("success 200", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodDelete, "/spaces/api/v1/categories/3.json", http.StatusOK, "")
		if err := newTestClient(mock).Categories.Delete(context.Background(), 3); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		if err := newTestClient(NewMockRoundTripper()).Categories.Delete(context.Background(), 0); err == nil {
			t.Fatal("expected error for id=0")
		}
	})

	t.Run("api error", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodDelete, "/spaces/api/v1/categories/3.json",
			http.StatusNotFound, `{"error":"not found"}`)
		if err := newTestClient(mock).Categories.Delete(context.Background(), 3); err == nil {
			t.Fatal("expected error on 404")
		}
	})
}
