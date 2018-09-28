package model

// User defines a basic User for a Web app
type User struct {
	ID       int64
	Email    string
	Password string
	Profile  *Profile
}
