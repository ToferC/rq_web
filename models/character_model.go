package models

import (
	"time"

	"github.com/toferc/runequest"
)

// CharacterModel represents the Web struct for a runequest.Character
type CharacterModel struct {
	ID        int64
	Author    *User
	Character *runequest.Character
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
