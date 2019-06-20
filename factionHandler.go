package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/toferc/runequest"

	"github.com/gosimple/slug"

	"github.com/gorilla/mux"
	"github.com/thewhitetulip/Tasks/sessions"
	"github.com/toferc/rq_web/database"
	"github.com/toferc/rq_web/models"
)

// FactionIndexHandler renders the basic character roster page
func FactionIndexHandler(w http.ResponseWriter, req *http.Request) {

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

	factions, err := database.ListFactions(db)
	if err != nil {
		panic(err)
	}

	wc := WebChar{
		SessionUser: username,
		IsLoggedIn:  loggedIn,
		IsAdmin:     isAdmin,
		Factions:    factions,
	}
	Render(w, "templates/faction_index.html", wc)
}

// UserFactionIndexHandler renders the basic character roster page
func UserFactionIndexHandler(w http.ResponseWriter, req *http.Request) {

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

	factions, err := database.ListUserFactions(db, username)
	if err != nil {
		panic(err)
	}

	wc := WebChar{
		SessionUser: username,
		IsLoggedIn:  loggedIn,
		IsAdmin:     isAdmin,
		Factions:    factions,
	}
	Render(w, "templates/user_faction_index.html", wc)
}

// FactionHandler renders a character in a Web page
func FactionHandler(w http.ResponseWriter, req *http.Request) {

	session, err := sessions.Store.Get(req, "session")

	if err != nil {
		log.Println("error identifying session")
		// in case of error
	}

	// Prep for user authentication
	sessionMap := getUserSessionValues(session)

	username := sessionMap["username"]
	loggedIn := sessionMap["loggedin"]
	isAdmin := sessionMap["isAdmin"]

	vars := mux.Vars(req)
	slug := vars["slug"]

	fac, err := database.SlugLoadFaction(db, slug)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Unable to load Faction")
	}

	fmt.Println(fac)

	IsAuthor := false

	if username == fac.Author.UserName {
		IsAuthor = true
	}

	cms, err := database.LoadFactionCharacterModels(db, fac.CharacterModelSlugs)
	if err != nil {
		panic(err)
	}

	for _, c := range cms {
		fmt.Println(c.Character.Name)
	}

	wc := WebChar{
		Faction:           fac,
		FactionCharacters: cms,
		IsAuthor:          IsAuthor,
		IsLoggedIn:        loggedIn,
		SessionUser:       username,
		IsAdmin:           isAdmin,
		StringArray:       runequest.StatMap,
	}

	// Render page
	Render(w, "templates/view_faction.html", wc)

}

// AddFactionHandler creates a user-generated faction
func AddFactionHandler(w http.ResponseWriter, req *http.Request) {

	// Get session values or redirect to Login
	session, err := sessions.Store.Get(req, "session")

	if err != nil {
		log.Println("error identifying session")
		http.Redirect(w, req, "/login/", 302)
		return
		// in case of error
	}

	// Prep for user authentication
	sessionMap := getUserSessionValues(session)

	username := sessionMap["username"]
	loggedIn := sessionMap["loggedin"]
	isAdmin := sessionMap["isAdmin"]

	if username == "" {
		// Add user message
		http.Redirect(w, req, "/", 302)
	}

	user, err := database.LoadUser(db, username)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, req, "/", 302)
	}

	cms, err := database.ListUserCharacterModels(db, username)

	factions, err := database.ListFactions(db)
	if err != nil {
		panic(err)
	}

	wc := WebChar{
		IsAuthor:        true,
		SessionUser:     username,
		IsLoggedIn:      loggedIn,
		IsAdmin:         isAdmin,
		Counter:         numToArray(11),
		CharacterModels: cms,
		Factions:        factions,
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/add_faction.html", wc)

	}

	if req.Method == "POST" { // POST

		err := req.ParseMultipartForm(MaxMemory)
		if err != nil {
			panic(err)
		}

		fac := &models.Faction{}

		// Use Schema Decoder to create struct
		err = decoder.Decode(fac, req.PostForm)
		if err != nil {
			panic(err)
		}

		fac.Slug = slug.Make(fmt.Sprintf("%s-%s", username, fac.Name))

		fac.Author = user

		// Get character Slugs
		for i := 1; i < 11; i++ {
			slug := req.FormValue(fmt.Sprintf("Character-%d", i))
			if slug != "" {
				fac.CharacterModelSlugs = append(fac.CharacterModelSlugs, slug)
			}
		}

		err = database.SaveFaction(db, fac)
		if err != nil {
			log.Panic(err)
		} else {
			fmt.Println("Saved Faction")
		}

		url := fmt.Sprintf("/view_faction/%s", fac.Slug)

		http.Redirect(w, req, url, http.StatusFound)
	}
}

// ModifyFactionHandler renders an editable Faction in a Web page
func ModifyFactionHandler(w http.ResponseWriter, req *http.Request) {

	// Get session values or redirect to Login
	session, err := sessions.Store.Get(req, "session")

	if err != nil {
		log.Println("error identifying session")
		http.Redirect(w, req, "/login/", 302)
		return
		// in case of error
	}

	// Prep for user authentication
	sessionMap := getUserSessionValues(session)

	username := sessionMap["username"]
	loggedIn := sessionMap["loggedin"]
	isAdmin := sessionMap["isAdmin"]

	vars := mux.Vars(req)
	fSlug := vars["slug"]

	fac, err := database.SlugLoadFaction(db, fSlug)
	if err != nil {
		fmt.Println(err)
	}

	if fac.Author == nil {
		fac.Author = &models.User{
			UserName: "",
		}
	}

	factionCharacters, err := database.LoadFactionCharacterModels(db, fac.CharacterModelSlugs)
	if err != nil {
		log.Panic(err)
	}

	// Validate that User == Author
	IsAuthor := false

	if username == fac.Author.UserName || isAdmin == "true" {
		IsAuthor = true
	} else {
		http.Redirect(w, req, "/", 302)
	}

	cms, err := database.ListUserCharacterModels(db, username)

	// Add empty CharacterModel options if < 10
	if len(factionCharacters) < 10 {
		for i := len(factionCharacters); i < 11; i++ {
			tcm := &models.CharacterModel{
				Character: &runequest.Character{
					Name: "",
				},
				Slug: "",
			}
			factionCharacters = append(factionCharacters, tcm)
		}
	}

	wc := WebChar{
		Faction:           fac,
		FactionCharacters: factionCharacters,
		IsAuthor:          IsAuthor,
		SessionUser:       username,
		IsLoggedIn:        loggedIn,
		IsAdmin:           isAdmin,
		CharacterModels:   cms,
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/modify_faction.html", wc)

	}

	if req.Method == "POST" { // POST

		err := req.ParseMultipartForm(MaxMemory)
		if err != nil {
			panic(err)
		}

		err = decoder.Decode(fac, req.PostForm)
		if err != nil {
			panic(err)
		}

		fac.Slug = slug.Make(fmt.Sprintf("%s-%s", username, fac.Name))

		// Reset string array
		fac.CharacterModelSlugs = []string{}

		// Get character Slugs
		for i := 1; i < 11; i++ {
			slug := req.FormValue(fmt.Sprintf("Character-%d", i))
			if slug != "" {
				fac.CharacterModelSlugs = append(fac.CharacterModelSlugs, slug)
			}
		}

		// Insert Faction into App archive
		err = database.UpdateFaction(db, fac)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("Saved")
		}

		fmt.Println(fac)

		url := fmt.Sprintf("/view_faction/%s", fac.Slug)

		http.Redirect(w, req, url, http.StatusSeeOther)
	}
}

// DeleteFactionHandler renders a character in a Web page
func DeleteFactionHandler(w http.ResponseWriter, req *http.Request) {

	// Get session values or redirect to Login
	session, err := sessions.Store.Get(req, "session")

	if err != nil {
		log.Println("error identifying session")
		http.Redirect(w, req, "/login/", 302)
		return
		// in case of error
	}

	// Prep for user authentication
	sessionMap := getUserSessionValues(session)

	username := sessionMap["username"]
	loggedIn := sessionMap["loggedin"]
	isAdmin := sessionMap["isAdmin"]

	vars := mux.Vars(req)
	slug := vars["slug"]

	fac, err := database.SlugLoadFaction(db, slug)
	if err != nil {
		fmt.Println(err)
	}

	// Validate that User == Author
	IsAuthor := false

	if username == fac.Author.UserName || isAdmin == "true" {
		IsAuthor = true
	} else {
		http.Redirect(w, req, "/", 302)
	}

	wc := WebChar{
		Faction:     fac,
		IsAuthor:    IsAuthor,
		SessionUser: username,
		IsLoggedIn:  loggedIn,
		IsAdmin:     isAdmin,
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/delete_faction.html", wc)

	}

	if req.Method == "POST" {

		err := database.DeleteFaction(db, fac.ID)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("Deleted Faction")
		}

		url := fmt.Sprint("/faction_index/")

		http.Redirect(w, req, url, http.StatusSeeOther)
	}
}
