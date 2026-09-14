package models

import "time"

// Comment represents a comment on a wiki page.
type Comment struct {
	ID         int64      `json:"id"`
	ParentID   int64      `json:"parentId,omitempty"`
	Content    string     `json:"content"`
	State      string     `json:"state"`
	Identifier *string    `json:"identifier,omitempty"`
	Selection  *string    `json:"selection,omitempty"`
	IsPrivate  bool       `json:"isPrivate"`
	Page       EntityRef  `json:"page"`
	Space      EntityRef  `json:"space"`
	CreatedAt  *time.Time `json:"createdAt,omitempty"`
	UpdatedAt  *time.Time `json:"updatedAt,omitempty"`
	CreatedBy  *UserRef   `json:"createdBy,omitempty"`
	UpdatedBy  *UserRef   `json:"updatedBy,omitempty"`
}

// CommentCreate holds fields required to create a new comment.
type CommentCreate struct {
	ParentID  *int64 `json:"parentID,omitempty"`
	Content   string `json:"content"`
	IsPrivate bool   `json:"isPrivate,omitempty"`
}

// CommentUpdate holds fields that can be updated on an existing comment.
type CommentUpdate struct {
	Content   *string `json:"content,omitempty"`
	State     *string `json:"state,omitempty"`
	IsPrivate *bool   `json:"isPrivate,omitempty"`
}

// CommentParent is a top-level comment that may contain replies.
type CommentParent struct {
	Comment
	Replies []Comment `json:"replies,omitempty"`
}

// CommentResponse is the single-resource response wrapper for a comment.
type CommentResponse struct {
	Comment  Comment      `json:"comment"`
	Included IncludedData `json:"included"`
}

// CommentsResponse is the list response wrapper for comments.
type CommentsResponse struct {
	Comments []CommentParent `json:"comments"`
	Included IncludedData    `json:"included"`
}
