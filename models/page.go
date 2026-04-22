package models

import "time"

// Page represents a wiki page within a space.
type Page struct {
	ID                          int64      `json:"id"`
	ParentID                    *int64     `json:"parentId,omitempty"`
	Title                       string     `json:"title"`
	Slug                        string     `json:"slug"`
	Content                     string     `json:"content"`
	ContentRevision             int64      `json:"contentRevision"`
	IsRequiredReading           bool       `json:"isRequiredReading"`
	IsHomePage                  bool       `json:"isHomePage"`
	IsPublished                 bool       `json:"isPublished"`
	IsPrivate                   bool       `json:"isPrivate"`
	Summary                     string     `json:"summary,omitempty"`
	IsFullWidth                 bool       `json:"isFullWidth"`
	Tags                        []Tag      `json:"tags,omitempty"`
	State                       string     `json:"state"`
	ChangeMessage               *string    `json:"changeMessage,omitempty"`
	DraftVersion                *int64     `json:"draftVersion,omitempty"`
	ReaderInlineCommentsEnabled bool       `json:"readerInlineCommentsEnabled"`
	Space                       EntityRef  `json:"space"`
	CreatedAt                   *time.Time `json:"createdAt,omitempty"`
	UpdatedAt                   *time.Time `json:"updatedAt,omitempty"`
	CreatedBy                   *UserRef   `json:"createdBy,omitempty"`
	UpdatedBy                   *UserRef   `json:"updatedBy,omitempty"`
}

// PageCreate holds fields required to create a new page.
type PageCreate struct {
	ParentID                    *int64 `json:"parentId,omitempty"`
	Title                       string `json:"title"`
	Slug                        string `json:"slug,omitempty"`
	Content                     string `json:"content,omitempty"`
	IsRequiredReading           bool   `json:"isRequiredReading,omitempty"`
	Tags                        []Tag  `json:"tags,omitempty"`
	IsFullWidth                 bool   `json:"isFullWidth,omitempty"`
	IsMinorChange               bool   `json:"isMinorChange,omitempty"`
	IsPublish                   bool   `json:"isPublish,omitempty"`
	ChangeMessage               string `json:"changeMessage,omitempty"`
	ReaderInlineCommentsEnabled bool   `json:"readerInlineCommentsEnabled,omitempty"`
}

// PageUpdate holds fields that can be updated on an existing page.
type PageUpdate struct {
	ParentID                    *int64   `json:"parentId,omitempty"`
	Title                       *string  `json:"title,omitempty"`
	Slug                        *string  `json:"slug,omitempty"`
	IsFullWidth                 *bool    `json:"isFullWidth,omitempty"`
	IsRequiredReading           *bool    `json:"isRequiredReading,omitempty"`
	IsPublish                   *bool    `json:"isPublish,omitempty"`
	Tags                        *[]Tag   `json:"tags,omitempty"`
	Content                     *string  `json:"content,omitempty"`
	IsMinorChange               *bool    `json:"isMinorChange,omitempty"`
	ChangeMessage               *string  `json:"changeMessage,omitempty"`
	ReaderInlineCommentsEnabled *bool    `json:"readerInlineCommentsEnabled,omitempty"`
}

// PageDuplicate holds fields required to duplicate a page.
type PageDuplicate struct {
	Title    string `json:"title"`
	ParentID *int64 `json:"parentId,omitempty"`
	Slug     string `json:"slug,omitempty"`
}

// PageResponse is the single-resource response wrapper for a page.
type PageResponse struct {
	Page     Page         `json:"page"`
	Included IncludedData `json:"included,omitempty"`
}

// PageTreeNode represents a page in the tree returned by the list endpoint.
// The API returns pages as a nested tree, not a flat list.
type PageTreeNode struct {
	ID         int64          `json:"id"`
	Slug       string         `json:"slug"`
	Title      string         `json:"title"`
	UpdatedAt  *time.Time     `json:"updatedAt,omitempty"`
	ChildPages []PageTreeNode `json:"childPages"`
}

// PagesResponse is the list response wrapper for pages.
// The API returns a single root page tree under the "pages" key.
type PagesResponse struct {
	Pages    PageTreeNode `json:"pages"`
	Included IncludedData `json:"included,omitempty"`
	Meta     ResponseMeta `json:"meta"`
}
