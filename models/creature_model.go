package models

import (
	"time"

	"github.com/toferc/runequest"
)

// CreatureModel represents the Web struct for a runequest.creature
type CreatureModel struct {
	ID        int64
	Author    *User
	creature  *runequest.Creature
	Open      bool
	Likes     int
	Image     *Image
	Slug      string
	CreatedAt time.Time `sql:"default:now()"`
	UpdatedAt time.Time
}

// Image  is the image and path for an Image
type Image struct {
	ID   int
	Path string
}
