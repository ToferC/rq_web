package models

import (
	"fmt"
	"time"
)

//User implements a generic user model
type User struct {
	ID        int64
	UserName  string `sql:",unique"`
	Email     string
	Password  string
	IsAdmin   bool
	CreatedAt time.Time `sql:"default:now()"`
	UpdatedAt time.Time
}

func (u User) String() string {
	text := fmt.Sprintf("%s %s %s", u.UserName, u.Email, u.Password)
	return text
}
