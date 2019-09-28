package main

import (
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

	homelands, err := database.ListHomelandModels(db)
	if err != nil {
		panic(err)
	}

	occupations, err := database.ListOccupationModels(db)
	if err != nil {
		panic(err)
	}

	cults, err := database.ListCultModels(db)
	if err != nil {
		panic(err)
	}

	wc := WebChar{
		SessionUser:      username,
		IsLoggedIn:       loggedIn,
		IsAdmin:          isAdmin,
		CharacterModels:  characters,
		HomelandModels:   homelands,
		OccupationModels: occupations,
		CultModels:       cults,
		Limit:            limit,
		Offset:           offset,
		Index:            "true",
	}

	if req.Method == "GET" {
		Render(w, "templates/crafted_roster.html", wc)
	}

	if req.Method == "POST" {

		// Parse Form and redirect
		err := req.ParseForm()
		if err != nil {
			panic(err)
		}

		query := &database.QueryArgs{
			Homeland:   req.FormValue("Homeland"),
			Occupation: req.FormValue("Occupation"),
			Cult:       req.FormValue("Cult"),
		}

		wc.CharacterModels, err = query.GetFilteredCharacterModels(db)
		if err != nil {
			log.Println(err)
		}
		Render(w, "templates/crafted_roster.html", wc)
	}

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

	homelands, err := database.ListHomelandModels(db)
	if err != nil {
		panic(err)
	}

	occupations, err := database.ListOccupationModels(db)
	if err != nil {
		panic(err)
	}

	cults, err := database.ListCultModels(db)
	if err != nil {
		panic(err)
	}

	wc := WebChar{
		SessionUser:      username,
		IsLoggedIn:       loggedIn,
		IsAdmin:          isAdmin,
		CharacterModels:  characters,
		HomelandModels:   homelands,
		OccupationModels: occupations,
		CultModels:       cults,
		Offset:           offset,
		Limit:            limit,
		Index:            "true",
	}

	if req.Method == "GET" {
		Render(w, "templates/roster.html", wc)
	}

	if req.Method == "POST" {

		// Parse Form and redirect
		err := req.ParseForm()
		if err != nil {
			panic(err)
		}

		query := &database.QueryArgs{
			Homeland:   req.FormValue("Homeland"),
			Occupation: req.FormValue("Occupation"),
			Cult:       req.FormValue("Cult"),
		}

		wc.CharacterModels, err = query.GetFilteredCharacterModels(db)
		if err != nil {
			log.Println(err)
		}
		Render(w, "templates/query_roster.html", wc)
	}
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

	homelands, err := database.ListHomelandModels(db)
	if err != nil {
		panic(err)
	}

	occupations, err := database.ListOccupationModels(db)
	if err != nil {
		panic(err)
	}

	cults, err := database.ListCultModels(db)
	if err != nil {
		panic(err)
	}

	wc := WebChar{
		SessionUser:      username,
		IsLoggedIn:       loggedIn,
		IsAdmin:          isAdmin,
		CharacterModels:  characters,
		HomelandModels:   homelands,
		OccupationModels: occupations,
		CultModels:       cults,
		Limit:            limit,
		Offset:           offset,
		Index:            "true",
	}

	if req.Method == "GET" {
		Render(w, "templates/random_roster.html", wc)
	}

	if req.Method == "POST" {

		// Parse Form and redirect
		err := req.ParseForm()
		if err != nil {
			panic(err)
		}

		query := &database.QueryArgs{
			Homeland:   req.FormValue("Homeland"),
			Occupation: req.FormValue("Occupation"),
			Cult:       req.FormValue("Cult"),
		}

		wc.CharacterModels, err = query.GetFilteredCharacterModels(db)
		if err != nil {
			log.Println(err)
		}
		Render(w, "templates/random_roster.html", wc)
	}

}
