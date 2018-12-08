package main

import (
	"os"
)

func configGoogleOAUTH() {
	os.Getenv("GOOGLE_CLIENT_ID")
	os.Getenv("GOOGLE_CLIENT_SECRET")
}

// Config is a base configuration for Goole Oauth
type Config struct {
	ClientID     string
	ClientSecret string
}
