package main

import (
	"context"
	"log/slog"
	"net/url"
	"os"

	"github.com/teamwork/spacessdkgo/client"
	"github.com/teamwork/spacessdkgo/models"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	baseURL := os.Getenv("SPACES_API_BASE_URL")
	if baseURL == "" {
		baseURL = "https://example.teamwork.com/spaces/api/v1"
	}
	apiKey := os.Getenv("SPACES_API_KEY")

	c := client.NewClient(baseURL,
		client.WithAPIKey(apiKey),
		client.WithLogger(logger),
	)

	ctx := context.Background()

	// List spaces.
	spaces, err := c.Spaces.List(ctx, url.Values{})
	if err != nil {
		logger.ErrorContext(ctx, "failed to list spaces", slog.Any("error", err))
		return
	}
	logger.InfoContext(ctx, "spaces listed", slog.Int("count", len(spaces.Spaces)))

	// Create a space (demo only — requires valid credentials).
	purpose := "Demo space created by the SDK"
	newSpace, err := c.Spaces.Create(ctx, &models.SpaceCreate{
		Title:      "SDK Demo Space",
		Code:       "sdk-demo",
		Purpose:    &purpose,
		SpaceColor: "#3B82F6",
	})
	if err != nil {
		logger.ErrorContext(ctx, "failed to create space", slog.Any("error", err))
		return
	}
	logger.InfoContext(ctx, "space created",
		slog.Int64("id", newSpace.Space.ID),
		slog.String("title", newSpace.Space.Title),
	)

	spaceID := newSpace.Space.ID

	// Create a page in the new space.
	title := "Welcome"
	newPage, err := c.Pages.Create(ctx, spaceID, &models.PageCreate{
		Title:     title,
		Content:   "<p>Hello from the Spaces SDK!</p>",
		IsPublish: true,
	})
	if err != nil {
		logger.ErrorContext(ctx, "failed to create page", slog.Any("error", err))
		return
	}
	logger.InfoContext(ctx, "page created",
		slog.Int64("id", newPage.Page.ID),
		slog.String("title", newPage.Page.Title),
	)

	// Search for pages.
	searchResp, err := c.Search.Search(ctx, models.SearchFilter{
		Query:   "Welcome",
		SpaceID: []int64{spaceID},
	})
	if err != nil {
		logger.ErrorContext(ctx, "search failed", slog.Any("error", err))
		return
	}
	logger.InfoContext(ctx, "search complete",
		slog.Int64("totalResults", searchResp.TotalResults),
	)

	// Clean up: delete the page and space.
	if err := c.Pages.Delete(ctx, spaceID, newPage.Page.ID); err != nil {
		logger.ErrorContext(ctx, "failed to delete page", slog.Any("error", err))
	}
	if err := c.Spaces.Delete(ctx, spaceID); err != nil {
		logger.ErrorContext(ctx, "failed to delete space", slog.Any("error", err))
	}
	logger.InfoContext(ctx, "cleanup complete")
}
