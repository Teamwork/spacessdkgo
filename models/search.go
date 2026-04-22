package models

// SearchFilter holds parameters for the search endpoint.
type SearchFilter struct {
	Query      string   `qs:"q"`
	SpaceID    []int64  `qs:"spaceid"`
	Limit      *int64   `qs:"limit,omitempty"`
	Offset     *int64   `qs:"offset,omitempty"`
	IncludeDel bool     `qs:"deleted,omitempty"`
	Fields     []string `qs:"fields,omitempty"`
}

// SearchPageSpace is a compact space reference within search results.
type SearchPageSpace struct {
	ID   int64  `json:"id"`
	Type string `json:"type"`
}

// ShortTag is a compact tag representation within search results.
type ShortTag struct {
	TagID int64  `json:"tagId"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

// SearchPage represents a single search result item.
type SearchPage struct {
	PageID      int64               `json:"pageId"`
	Title       string              `json:"title"`
	Slug        string              `json:"slug"`
	MatchedText map[string][]string `json:"matched"`
	Space       SearchPageSpace     `json:"space"`
	Tags        []ShortTag          `json:"tags,omitempty"`
}

// SearchSpaceIncluded holds spaces included in search results.
type SearchSpaceIncluded struct {
	Spaces map[string]Space `json:"spaces"`
}

// SearchResponse is the response returned by the search endpoint.
type SearchResponse struct {
	TotalResults int64               `json:"totalResults"`
	Results      []SearchPage        `json:"results"`
	Included     SearchSpaceIncluded `json:"included"`
}
