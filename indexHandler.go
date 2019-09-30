package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/thewhitetulip/Tasks/sessions"
	"github.com/toferc/rq_web/database"
	"github.com/toferc/rq_web/models"
)

// CraftedCharacterIndexHandler renders the basic character roster page
func CraftedCharacterIndexHandler(w http.ResponseWriter, req *http.Request) {

	session, err := sessions.Store.Get(req, "session")

	if err != nil {
		log.Println("error identifying session")
		Render(w, "templates/login.html", nil)
		return
		// in case of error
	}

	// Prep for user authentication
	sessionMap := getUserSessionValues(session)

	username := sessionMap["username"]
	loggedIn := sessionMap["loggedin"]
	isAdmin := sessionMap["isAdmin"]

	values := mux.Vars(req)

	l := values["limit"]
	limit, err := strconv.Atoi(l)
	if err != nil {
		limit = 66
	}

	o := values["offset"]
	offset, err := strconv.Atoi(o)
	if err != nil {
		offset = 0
	}

	characters, err := database.ListCraftedCharacterModels(db, limit, offset)
	if err != nil {
		panic(err)
	}

	for _, cm := range characters {

		if cm.Image == nil {
			cm.Image = new(models.Image)
			cm.Image.Path = DefaultCharacterPortrait
		}
	}

	wc := WebChar{
		SessionUser:     username,
		IsLoggedIn:      loggedIn,
		IsAdmin:         isAdmin,
		CharacterModels: characters,
		Limit:           limit,
		Offset:          offset,
		Index:           "true",
	}

	Render(w, "templates/crafted_roster.html", wc)
}

// AllCharacterIndexHandler renders the basic character roster page
func AllCharacterIndexHandler(w http.ResponseWriter, req *http.Request) {

	session, err := sessions.Store.Get(req, "session")

	if err != nil {
		log.Println("error identifying session")
		Render(w, "templates/login.html", nil)
		return
		// in case of error
	}

	// Prep for user authentication
	sessionMap := getUserSessionValues(session)

	username := sessionMap["username"]
	loggedIn := sessionMap["loggedin"]
	isAdmin := sessionMap["isAdmin"]

	values := mux.Vars(req)

	l := values["limit"]
	limit, err := strconv.Atoi(l)
	if err != nil {
		limit = 66
	}

	o := values["offset"]
	offset, err := strconv.Atoi(o)
	if err != nil {
		offset = 0
	}

	characters, err := database.PaginateCharacterModels(db, limit, offset)
	if err != nil {
		panic(err)
	}

	for _, cm := range characters {

		if cm.Image == nil {
			cm.Image = new(models.Image)
			cm.Image.Path = DefaultCharacterPortrait
		}
	}

	wc := WebChar{
		SessionUser:     username,
		IsLoggedIn:      loggedIn,
		IsAdmin:         isAdmin,
		CharacterModels: characters,
		Offset:          offset,
		Limit:           limit,
		Index:           "true",
	}

	Render(w, "templates/roster.html", wc)
}

// RandomCharacterIndexHandler renders the basic character roster page
func RandomCharacterIndexHandler(w http.ResponseWriter, req *http.Request) {

	session, err := sessions.Store.Get(req, "session")

	if err != nil {
		log.Println("error identifying session")
		Render(w, "templates/login.html", nil)
		return
		// in case of error
	}

	// Prep for user authentication
	sessionMap := getUserSessionValues(session)

	username := sessionMap["username"]
	loggedIn := sessionMap["loggedin"]
	isAdmin := sessionMap["isAdmin"]

	values := mux.Vars(req)

	l := values["limit"]
	limit, err := strconv.Atoi(l)
	if err != nil {
		limit = 66
	}

	o := values["offset"]
	offset, err := strconv.Atoi(o)
	if err != nil {
		offset = 0
	}

	characters, err := database.ListRandomCharacterModels(db, limit, offset)
	if err != nil {
		panic(err)
	}

	for _, cm := range characters {

		if cm.Image == nil {
			cm.Image = new(models.Image)
			cm.Image.Path = DefaultCharacterPortrait
		}
	}

	wc := WebChar{
		SessionUser:     username,
		IsLoggedIn:      loggedIn,
		IsAdmin:         isAdmin,
		CharacterModels: characters,
		Limit:           limit,
		Offset:          offset,
		Index:           "true",
	}

	Render(w, "templates/random_roster.html", wc)
}

// UserCharacterRosterHandler handles user-specific rosters
func UserCharacterRosterHandler(w http.ResponseWriter, req *http.Request) {

	session, err := sessions.Store.Get(req, "session")

	if err != nil {
		log.Println("error identifying session")
		Render(w, "templates/login.html", nil)
		return
		// in case of error
	}

	// Prep for user authentication
	sessionMap := getUserSessionValues(session)

	username := sessionMap["username"]
	loggedIn := sessionMap["loggedin"]
	isAdmin := sessionMap["isAdmin"]

	if username == "" {
		http.Redirect(w, req, "/", 302)
	}

	values := mux.Vars(req)

	l := values["limit"]
	limit, err := strconv.Atoi(l)
	if err != nil {
		limit = 66
	}

	o := values["offset"]
	offset, err := strconv.Atoi(o)
	if err != nil {
		offset = 0
	}

	characters, err := database.ListUserCharacterModels(db, username, limit, offset)
	if err != nil {
		log.Println(err)
	}

	for _, cm := range characters {
		if cm.Image == nil {
			cm.Image = new(models.Image)
			cm.Image.Path = DefaultCharacterPortrait
		}
	}

	u := &models.User{
		UserName: username,
	}

	wc := WebChar{
		SessionUser:     username,
		User:            u,
		IsLoggedIn:      loggedIn,
		IsAdmin:         isAdmin,
		CharacterModels: characters,
		Limit:           limit,
		Offset:          offset,
		Index:           "true",
	}

	Render(w, "templates/user_roster.html", wc)

}

// UserOpenCharacterRosterHandler handles user-specific rosters
func UserOpenCharacterRosterHandler(w http.ResponseWriter, req *http.Request) {

	session, err := sessions.Store.Get(req, "session")

	if err != nil {
		log.Println("error identifying session")
		Render(w, "templates/login.html", nil)
		return
		// in case of error
	}

	// Prep for user authentication
	sessionMap := getUserSessionValues(session)

	username := sessionMap["username"]
	loggedIn := sessionMap["loggedin"]
	isAdmin := sessionMap["isAdmin"]

	values := mux.Vars(req)

	pk := values["id"]

	if len(pk) == 0 {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}

	id, err := strconv.Atoi(pk)
	if err != nil {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}

	u, err := database.PKLoadUser(db, int64(id))
	if err != nil {
		fmt.Println(err)
	}

	l := values["limit"]
	limit, err := strconv.Atoi(l)
	if err != nil {
		limit = 66
	}

	o := values["offset"]
	offset, err := strconv.Atoi(o)
	if err != nil {
		offset = 0
	}

	characters, err := database.ListOpenUserCharacterModels(db, u.UserName, limit, offset)
	if err != nil {
		log.Println(err)
	}

	for _, cm := range characters {
		if cm.Image == nil {
			cm.Image = new(models.Image)
			cm.Image.Path = DefaultCharacterPortrait
		}
	}

	wc := WebChar{
		SessionUser:     username,
		IsLoggedIn:      loggedIn,
		IsAdmin:         isAdmin,
		User:            u,
		CharacterModels: characters,
		Limit:           limit,
		Offset:          offset,
		Index:           "true",
	}

	Render(w, "templates/user_open_roster.html", wc)

}
