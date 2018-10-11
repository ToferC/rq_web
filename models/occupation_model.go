package models

import (
	"time"

	"github.com/toferc/runequest"
)

// OccupationModel represents the apps version of runequest.Occupation
type OccupationModel struct {
	ID         int64
	Author     *User
	Occupation *runequest.Occupation
	Official   bool
	Open       bool
	Likes      int
	Image      *Image
	Slug       string
	CreatedAt  time.Time `sql:"default:now()"`
	UpdatedAt  time.Time
}
