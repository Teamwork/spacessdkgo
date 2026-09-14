package models

import "time"

// Space represents a Teamwork Spaces wiki space.
type Space struct {
	ID              int64      `json:"id"`
	Title           string     `json:"title"`
	Purpose         *string    `json:"purpose,omitempty"`
	Code            string     `json:"code"`
	State           string     `json:"state"`
	SpaceColor      string     `json:"spaceColor"`
	Icon            string     `json:"icon"`
	Banner          string     `json:"banner,omitempty"`
	LinkedProjectID int64      `json:"projectId,omitempty"`
	CreatedAt       *time.Time `json:"createdAt,omitempty"`
	UpdatedAt       *time.Time `json:"updatedAt,omitempty"`
	CreatedBy       *UserRef   `json:"createdBy,omitempty"`
	UpdatedBy       *UserRef   `json:"updatedBy,omitempty"`
}

// SpaceCreate holds fields required to create a new space.
type SpaceCreate struct {
	Title           string  `json:"title"`
	Code            string  `json:"code"`
	Purpose         *string `json:"purpose,omitempty"`
	SpaceColor      string  `json:"spaceColor,omitempty"`
	Icon            string  `json:"icon,omitempty"`
	LinkedProjectID *int64  `json:"projectId,omitempty"`
	CategoryID      *int64  `json:"categoryId,omitempty"`
}

// SpaceUpdate holds fields that can be updated on an existing space.
type SpaceUpdate struct {
	Title           *string `json:"title,omitempty"`
	Code            *string `json:"code,omitempty"`
	Purpose         *string `json:"purpose,omitempty"`
	SpaceColor      *string `json:"spaceColor,omitempty"`
	Icon            *string `json:"icon,omitempty"`
	LinkedProjectID *int64  `json:"projectId,omitempty"`
	State           *string `json:"state,omitempty"`
	CategoryID      *int64  `json:"categoryId,omitempty"`
}

// SpaceCollaborator represents a user or entity collaborating in a space.
type SpaceCollaborator struct {
	ID   int64  `json:"id"`
	Type string `json:"type"`
}

// SpaceCollaboratorsResponse is returned by the collaborators endpoint.
type SpaceCollaboratorsResponse struct {
	Collaborators      []SpaceCollaborator `json:"collaborators"`
	CollaboratorsCount int64               `json:"collaboratorsCount"`
}

// SpaceResponse is the single-resource response wrapper for a space.
type SpaceResponse struct {
	Space    Space        `json:"space"`
	Included IncludedData `json:"included"`
}

// SpacesResponse is the list response wrapper for spaces.
type SpacesResponse struct {
	Spaces   []Space      `json:"spaces"`
	Included IncludedData `json:"included"`
	Meta     ResponseMeta `json:"meta"`
}
