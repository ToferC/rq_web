package models

import "time"

// Campaign represents a user's campaign framework
type Campaign struct {
	ID          int64
	Author      *User
	Name        string
	Description string
	Open        bool
	Likes       int
	Image       *Image
	Slug        string
	CreatedAt   time.Time `sql:"default:now()"`
	UpdatedAt   time.Time
}
