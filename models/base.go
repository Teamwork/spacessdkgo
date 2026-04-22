package models

import "time"

// UserRef is a reference to a user by ID and type.
type UserRef struct {
	ID   int64  `json:"id"`
	Type string `json:"type"`
}

// EntityRef is a reference to any entity with optional metadata.
type EntityRef struct {
	ID   int64          `json:"id"`
	Type string         `json:"type"`
	Meta map[string]any `json:"meta,omitempty"`
}

// State represents the lifecycle state of a resource.
type State string

const (
	StateActive   State = "active"
	StateArchived State = "archived"
	StateRemoved  State = "removed"
	StateDeleted  State = "deleted"
)

// TimeRef is a helper for optional timestamp fields.
type TimeRef = *time.Time
