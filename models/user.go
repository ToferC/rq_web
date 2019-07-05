package models

import (
	"fmt"
	"time"
)

//User implements a generic user model
type User struct {
	ID         int64
	UserName   string    `sql:",unique"`
	Email      string    `json:"-"`
	Password   string    `json:"-"`
	IsAdmin    bool      `json:"-"`
	CreatedAt  time.Time `sql:"default:now()"`
	Characters int
	UpdatedAt  time.Time
}

func (u User) String() string {
	text := fmt.Sprintf("%s %s %s %d", u.UserName, u.Email, u.Password, u.Characters)
	return text
}
