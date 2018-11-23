package main

import (
	"fmt"
	"github.com/toferc/runequest"
	"time"

	"github.com/gorilla/sessions"
)

func getUserSessionValues(s *sessions.Session) map[string]string {

	sessionMap := map[string]string{
		"username": "",
		"loggedin": "false",
		"isAdmin":  "false",
	}

	// Prep for user authentication

	u := s.Values["username"]
	l := s.Values["loggedin"]
	a := s.Values["isAdmin"]

	// Type assertation
	if user, ok := u.(string); !ok {
	} else {
		fmt.Println(user)
		sessionMap["username"] = user
	}

	// Type assertation
	if loggin, ok := l.(string); !ok {
	} else {
		fmt.Println(loggin)
		sessionMap["loggedin"] = loggin
	}

	// Type assertation
	if admin, ok := a.(string); !ok {
	} else {
		fmt.Println(admin)
		sessionMap["isAdmin"] = admin
	}
	return sessionMap
}

// numToArray takes and int and returns an array of [1:int]
func numToArray(m int) []int {

	a := []int{}

	for i := 1; i < m+1; i++ {
		a = append(a, i)
	}

	return a
}

// CreateUpdate creates an update for a runequest Object
func CreateUpdate(text string, value int) *runequest.Update {

	t := time.Now()
	tString := t.Format("2006-01-02")

	update := &runequest.Update{
		Date:  tString,
		Event: fmt.Sprintf("%s", text),
		Value: value,
	}
	fmt.Printf("Add Update: %s %d\n", text, value)
	return update
}

func createName(coreString string, userString string) string {

	targetString := ""

	if userString != "" {
		targetString = fmt.Sprintf("%s (%s)", coreString, userString)
	} else {
		targetString = fmt.Sprintf("%s", coreString)
	}

	return targetString
}
