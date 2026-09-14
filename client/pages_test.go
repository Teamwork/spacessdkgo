package client

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/teamwork/spacessdkgo/models"
)

func TestPageService_Get(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodGet, "/spaces/api/v1/spaces/10/pages/20.json", http.StatusOK,
			models.PageResponse{Page: models.Page{ID: 20, Title: "Getting Started"}})

		got, err := newTestClient(mock).Pages.Get(context.Background(), 10, 20)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Page.ID != 20 {
			t.Errorf("got ID %d, want 20", got.Page.ID)
		}
		if got.Page.Title != "Getting Started" {
			t.Errorf("got title %q, want %q", got.Page.Title, "Getting Started")
		}
		reqs := mock.GetRequests()
		if reqs[0].Method != http.MethodGet {
			t.Errorf("expected GET, got %s", reqs[0].Method)
		}
		if reqs[0].URL.Path != "/spaces/api/v1/spaces/10/pages/20.json" {
			t.Errorf("unexpected path %s", reqs[0].URL.Path)
		}
	})

	t.Run("invalid spaceID", func(t *testing.T) {
		_, err := newTestClient(NewMockRoundTripper()).Pages.Get(context.Background(), 0, 20)
		if err == nil {
			t.Fatal("expected error for spaceID=0")
		}
	})

	t.Run("invalid pageID", func(t *testing.T) {
		_, err := newTestClient(NewMockRoundTripper()).Pages.Get(context.Background(), 10, 0)
		if err == nil {
			t.Fatal("expected error for pageID=0")
		}
	})

	t.Run("api error", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodGet, "/spaces/api/v1/spaces/10/pages/20.json", http.StatusNotFound, `{"error":"not found"}`)
		_, err := newTestClient(mock).Pages.Get(context.Background(), 10, 20)
		if err == nil {
			t.Fatal("expected error on 404")
		}
	})
}

func TestPageService_List(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodGet, "/spaces/api/v1/spaces/10/pages.json", http.StatusOK,
			models.PagesResponse{Pages: models.PageTreeNode{
				ID:    1,
				Title: "Home",
				ChildPages: []models.PageTreeNode{
					{ID: 2, Title: "Intro"},
					{ID: 3, Title: "Setup"},
				},
			}})

		got, err := newTestClient(mock).Pages.List(context.Background(), 10, url.Values{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Pages.ID != 1 {
			t.Fatalf("got root ID %d, want 1", got.Pages.ID)
		}
		if len(got.Pages.ChildPages) != 2 {
			t.Fatalf("got %d child pages, want 2", len(got.Pages.ChildPages))
		}
	})

	t.Run("passes query params", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodGet, "/spaces/api/v1/spaces/10/pages.json", http.StatusOK, models.PagesResponse{})

		_, err := newTestClient(mock).Pages.List(context.Background(), 10,
			url.Values{"page": {"3"}})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if q := mock.GetRequests()[0].URL.Query().Get("page"); q != "3" {
			t.Errorf("expected page=3, got %q", q)
		}
	})

	t.Run("invalid spaceID", func(t *testing.T) {
		_, err := newTestClient(NewMockRoundTripper()).Pages.List(context.Background(), 0, url.Values{})
		if err == nil {
			t.Fatal("expected error for spaceID=0")
		}
	})

	t.Run("api error", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodGet, "/spaces/api/v1/spaces/10/pages.json", http.StatusInternalServerError, `{"error":"internal"}`)
		_, err := newTestClient(mock).Pages.List(context.Background(), 10, url.Values{})
		if err == nil {
			t.Fatal("expected error on 500")
		}
	})
}

func TestPageService_Home(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodGet, "/spaces/api/v1/spaces/10/homepage.json", http.StatusOK,
			models.PageResponse{Page: models.Page{ID: 1, IsHomePage: true}})

		got, err := newTestClient(mock).Pages.Home(context.Background(), 10)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !got.Page.IsHomePage {
			t.Error("expected IsHomePage=true")
		}
		if mock.GetRequests()[0].URL.Path != "/spaces/api/v1/spaces/10/homepage.json" {
			t.Errorf("unexpected path %s", mock.GetRequests()[0].URL.Path)
		}
	})

	t.Run("invalid spaceID", func(t *testing.T) {
		_, err := newTestClient(NewMockRoundTripper()).Pages.Home(context.Background(), 0)
		if err == nil {
			t.Fatal("expected error for spaceID=0")
		}
	})

	t.Run("api error", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodGet, "/spaces/api/v1/spaces/10/homepage.json", http.StatusNotFound, `{"error":"not found"}`)
		_, err := newTestClient(mock).Pages.Home(context.Background(), 10)
		if err == nil {
			t.Fatal("expected error on 404")
		}
	})
}

func TestPageService_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodPost, "/spaces/api/v1/spaces/10/pages.json", http.StatusCreated,
			models.PageResponse{Page: models.Page{ID: 55, Title: "New Page"}})

		got, err := newTestClient(mock).Pages.Create(context.Background(), 10,
			&models.PageCreate{Title: "New Page"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Page.ID != 55 {
			t.Errorf("got ID %d, want 55", got.Page.ID)
		}
		if mock.GetRequests()[0].Method != http.MethodPost {
			t.Errorf("expected POST, got %s", mock.GetRequests()[0].Method)
		}
	})

	t.Run("success 200", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodPost, "/spaces/api/v1/spaces/10/pages.json", http.StatusOK,
			models.PageResponse{Page: models.Page{ID: 1}})
		_, err := newTestClient(mock).Pages.Create(context.Background(), 10,
			&models.PageCreate{Title: "X"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("invalid spaceID", func(t *testing.T) {
		_, err := newTestClient(NewMockRoundTripper()).Pages.Create(context.Background(), 0,
			&models.PageCreate{Title: "X"})
		if err == nil {
			t.Fatal("expected error for spaceID=0")
		}
	})

	t.Run("nil req", func(t *testing.T) {
		_, err := newTestClient(NewMockRoundTripper()).Pages.Create(context.Background(), 10, nil)
		if err == nil {
			t.Fatal("expected error for nil req")
		}
	})

	t.Run("api error", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodPost, "/spaces/api/v1/spaces/10/pages.json", http.StatusBadRequest, `{"error":"bad"}`)
		_, err := newTestClient(mock).Pages.Create(context.Background(), 10,
			&models.PageCreate{Title: "X"})
		if err == nil {
			t.Fatal("expected error on 400")
		}
	})
}

func TestPageService_Duplicate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodPost, "/spaces/api/v1/spaces/10/pages/20/duplicate.json", http.StatusCreated,
			models.PageResponse{Page: models.Page{ID: 99, Title: "Copy of Page"}})

		got, err := newTestClient(mock).Pages.Duplicate(context.Background(), 10, 20,
			&models.PageDuplicate{Title: "Copy of Page"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Page.ID != 99 {
			t.Errorf("got ID %d, want 99", got.Page.ID)
		}
		if mock.GetRequests()[0].URL.Path != "/spaces/api/v1/spaces/10/pages/20/duplicate.json" {
			t.Errorf("unexpected path %s", mock.GetRequests()[0].URL.Path)
		}
	})

	t.Run("invalid spaceID", func(t *testing.T) {
		_, err := newTestClient(NewMockRoundTripper()).Pages.Duplicate(context.Background(), 0, 20,
			&models.PageDuplicate{Title: "X"})
		if err == nil {
			t.Fatal("expected error for spaceID=0")
		}
	})

	t.Run("invalid pageID", func(t *testing.T) {
		_, err := newTestClient(NewMockRoundTripper()).Pages.Duplicate(context.Background(), 10, 0,
			&models.PageDuplicate{Title: "X"})
		if err == nil {
			t.Fatal("expected error for pageID=0")
		}
	})

	t.Run("nil req", func(t *testing.T) {
		_, err := newTestClient(NewMockRoundTripper()).Pages.Duplicate(context.Background(), 10, 20, nil)
		if err == nil {
			t.Fatal("expected error for nil req")
		}
	})

	t.Run("api error", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodPost, "/spaces/api/v1/spaces/10/pages/20/duplicate.json", http.StatusBadRequest, `{"error":"bad"}`)
		_, err := newTestClient(mock).Pages.Duplicate(context.Background(), 10, 20,
			&models.PageDuplicate{Title: "X"})
		if err == nil {
			t.Fatal("expected error on 400")
		}
	})
}

func TestPageService_Update(t *testing.T) {
	title := "Updated Title"

	t.Run("success", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodPatch, "/spaces/api/v1/spaces/10/pages/20.json", http.StatusOK,
			models.PageResponse{Page: models.Page{ID: 20, Title: "Updated Title"}})

		got, err := newTestClient(mock).Pages.Update(context.Background(), 10, 20,
			&models.PageUpdate{Title: &title})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Page.Title != "Updated Title" {
			t.Errorf("got title %q, want %q", got.Page.Title, "Updated Title")
		}
		if mock.GetRequests()[0].Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", mock.GetRequests()[0].Method)
		}
	})

	t.Run("invalid spaceID", func(t *testing.T) {
		_, err := newTestClient(NewMockRoundTripper()).Pages.Update(context.Background(), 0, 20,
			&models.PageUpdate{Title: &title})
		if err == nil {
			t.Fatal("expected error for spaceID=0")
		}
	})

	t.Run("invalid pageID", func(t *testing.T) {
		_, err := newTestClient(NewMockRoundTripper()).Pages.Update(context.Background(), 10, 0,
			&models.PageUpdate{Title: &title})
		if err == nil {
			t.Fatal("expected error for pageID=0")
		}
	})

	t.Run("nil req", func(t *testing.T) {
		_, err := newTestClient(NewMockRoundTripper()).Pages.Update(context.Background(), 10, 20, nil)
		if err == nil {
			t.Fatal("expected error for nil req")
		}
	})

	t.Run("api error", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodPatch, "/spaces/api/v1/spaces/10/pages/20.json", http.StatusForbidden, `{"error":"forbidden"}`)
		_, err := newTestClient(mock).Pages.Update(context.Background(), 10, 20,
			&models.PageUpdate{Title: &title})
		if err == nil {
			t.Fatal("expected error on 403")
		}
	})
}

func TestPageService_Delete(t *testing.T) {
	t.Run("success 204", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodDelete, "/spaces/api/v1/spaces/10/pages/20.json", http.StatusNoContent, "")
		if err := newTestClient(mock).Pages.Delete(context.Background(), 10, 20); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if mock.GetRequests()[0].Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", mock.GetRequests()[0].Method)
		}
	})

	t.Run("success 200", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodDelete, "/spaces/api/v1/spaces/10/pages/20.json", http.StatusOK, "")
		if err := newTestClient(mock).Pages.Delete(context.Background(), 10, 20); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("invalid spaceID", func(t *testing.T) {
		if err := newTestClient(NewMockRoundTripper()).Pages.Delete(context.Background(), 0, 20); err == nil {
			t.Fatal("expected error for spaceID=0")
		}
	})

	t.Run("invalid pageID", func(t *testing.T) {
		if err := newTestClient(NewMockRoundTripper()).Pages.Delete(context.Background(), 10, 0); err == nil {
			t.Fatal("expected error for pageID=0")
		}
	})

	t.Run("api error", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodDelete, "/spaces/api/v1/spaces/10/pages/20.json", http.StatusNotFound, `{"error":"not found"}`)
		if err := newTestClient(mock).Pages.Delete(context.Background(), 10, 20); err == nil {
			t.Fatal("expected error on 404")
		}
	})
}

func TestPageService_ListWithPrivate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodGet, "/spaces/api/v2/spaces/10/pages.json", http.StatusOK,
			models.SpaceContentResponse{SpaceContent: models.SpaceContentTree{
				Pages: models.PageTreeNode{
					ID:         1,
					Title:      "Home",
					ChildPages: []models.PageTreeNode{{ID: 2, Title: "Intro"}},
				},
				Private: []models.PageTreeNode{{ID: 3, Title: "Salaries"}},
			}})

		got, err := newTestClient(mock).Pages.ListWithPrivate(context.Background(), 10, url.Values{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.SpaceContent.Pages.ChildPages[0].ID != 2 {
			t.Errorf("got open page ID %d, want 2", got.SpaceContent.Pages.ChildPages[0].ID)
		}
		if len(got.SpaceContent.Private) != 1 || got.SpaceContent.Private[0].ID != 3 {
			t.Errorf("expected the private tree to survive decoding, got %+v", got.SpaceContent.Private)
		}

		// The v1 default must not be what answered: the two routes return
		// different shapes from the same query, and a v1 body decodes into the
		// v2 type as an empty tree rather than an error.
		reqs := mock.GetRequests()
		if reqs[0].URL.Path != "/spaces/api/v2/spaces/10/pages.json" {
			t.Errorf("unexpected path %s", reqs[0].URL.Path)
		}
	})

	t.Run("forwards query parameters", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodGet, "/spaces/api/v2/spaces/10/pages.json", http.StatusOK,
			models.SpaceContentResponse{})

		params := url.Values{}
		params.Set("pageSize", "25")
		if _, err := newTestClient(mock).Pages.ListWithPrivate(context.Background(), 10, params); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got := mock.GetRequests()[0].URL.Query().Get("pageSize"); got != "25" {
			t.Errorf("got pageSize %q, want %q", got, "25")
		}
	})

	t.Run("invalid spaceID", func(t *testing.T) {
		_, err := newTestClient(NewMockRoundTripper()).Pages.ListWithPrivate(context.Background(), 0, url.Values{})
		if err == nil {
			t.Fatal("expected error for spaceID=0")
		}
	})

	t.Run("api error", func(t *testing.T) {
		mock := NewMockRoundTripper()
		mock.AddResponse(http.MethodGet, "/spaces/api/v2/spaces/10/pages.json", http.StatusNotFound,
			`{"error":"not found"}`)
		_, err := newTestClient(mock).Pages.ListWithPrivate(context.Background(), 10, url.Values{})
		if err == nil {
			t.Fatal("expected error on 404")
		}
	})
}

// TestBaseURLV2 pins that only the version segment moves. normalizeBaseURL
// repairs a wrongly-configured base URL by truncating it at the v1 suffix, so a
// builder that appended rather than swapped would answer a v1 path with "/v2"
// stuck on the end.
func TestBaseURLV2(t *testing.T) {
	for _, tc := range []struct {
		input string
		want  string
	}{
		{"https://test.teamwork.com", "https://test.teamwork.com/spaces/api/v2"},
		{"https://test.teamwork.com/spaces/api/v1", "https://test.teamwork.com/spaces/api/v2"},
		{"https://test.teamwork.com/spaces/api/v1/extra", "https://test.teamwork.com/spaces/api/v2"},
		{"", ""},
	} {
		t.Run(tc.input, func(t *testing.T) {
			c := &Client{baseURL: normalizeBaseURL(tc.input)}
			if got := c.baseURLV2(); got != tc.want {
				t.Errorf("baseURLV2() = %q, want %q", got, tc.want)
			}
		})
	}
}
