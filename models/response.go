package models

// Pagination holds page metadata returned by list endpoints.
type Pagination struct {
	PageOffset int64 `json:"pageOffset"`
	PageSize   int64 `json:"pageSize"`
	Count      int64 `json:"count"`
}

// ResponseMeta wraps the pagination object returned in list responses.
type ResponseMeta struct {
	Page Pagination `json:"page"`
}

// IncludedData holds sideloaded resources returned alongside primary data.
type IncludedData struct {
	Spaces []Space `json:"spaces,omitempty"`
	Pages  []Page  `json:"pages,omitempty"`
	Users  []any   `json:"users,omitempty"`
}
