package models

import (
	"time"

	"github.com/toferc/runequest"
)

// HomelandModel represents the apps version of runequest.Homeland
type HomelandModel struct {
	ID        int64
	Author    *User
	Homeland  *runequest.Homeland
	Open      bool
	Likes     int
	Image     *Image
	Slug      string
	CreatedAt time.Time `sql:"default:now()"`
	UpdatedAt time.Time
}
