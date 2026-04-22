package models

// Tag represents a tag that can be applied to pages.
type Tag struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Color     string `json:"color"`
	PageCount *int64 `json:"pageCount,omitempty"`
}

// TagCreate holds fields required to create a new tag.
type TagCreate struct {
	Name  string `json:"name"`
	Color string `json:"color,omitempty"`
}

// TagUpdate holds fields that can be updated on an existing tag.
type TagUpdate struct {
	Name  *string `json:"name,omitempty"`
	Color *string `json:"color,omitempty"`
}

// TagResponse is the single-resource response wrapper for a tag.
// Note: the API uses "tags" as the JSON key for both single and list responses.
type TagResponse struct {
	Tag Tag `json:"tags"`
}

// TagsResponse is the list response wrapper for tags.
type TagsResponse struct {
	Tags []Tag `json:"tags"`
}
