package models

import (
	"time"

	"github.com/toferc/runequest"
)

// CultModel represents the apps version of runequest.Cult
type CultModel struct {
	ID        int64
	Author    *User
	Cult      *runequest.Cult
	Official  bool
	Open      bool
	Likes     int
	Image     *Image
	Slug      string
	CreatedAt time.Time `sql:"default:now()"`
	UpdatedAt time.Time
}
