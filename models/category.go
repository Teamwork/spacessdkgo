package models

// CategoryMeta holds metadata about a category.
type CategoryMeta struct {
	SpaceCount int64 `json:"spaceCount"`
}

// Category represents a grouping category for spaces.
type Category struct {
	ID    int64        `json:"id"`
	Name  string       `json:"name"`
	Color *string      `json:"color,omitempty"`
	Meta  CategoryMeta `json:"meta"`
}

// CategoryCreate holds fields required to create a new category.
type CategoryCreate struct {
	Name  string  `json:"name"`
	Color *string `json:"color,omitempty"`
}

// CategoryUpdate holds fields that can be updated on an existing category.
type CategoryUpdate struct {
	Name  *string `json:"name,omitempty"`
	Color *string `json:"color,omitempty"`
}

// CategoryResponse is the single-resource response wrapper for a category.
type CategoryResponse struct {
	Category Category     `json:"category"`
	Included IncludedData `json:"included"`
}

// CategoriesResponse is the list response wrapper for categories.
type CategoriesResponse struct {
	Categories []Category   `json:"categories"`
	Included   IncludedData `json:"included"`
	Meta       ResponseMeta `json:"meta"`
}
