package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gosimple/slug"

	"github.com/gorilla/mux"
	"github.com/thewhitetulip/Tasks/sessions"
	"github.com/toferc/rq_web/database"
	"github.com/toferc/rq_web/models"
)

// EncounterIndexHandler renders the basic character roster page
func EncounterIndexHandler(w http.ResponseWriter, req *http.Request) {

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

	encounters, err := database.ListEncounters(db)
	if err != nil {
		panic(err)
	}

	wc := WebChar{
		SessionUser: username,
		IsLoggedIn:  loggedIn,
		IsAdmin:     isAdmin,
		Encounters:  encounters,
	}
	Render(w, "templates/encounter_index.html", wc)
}

// UserEncounterIndexHandler renders the basic character roster page
func UserEncounterIndexHandler(w http.ResponseWriter, req *http.Request) {

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

	encounters, err := database.ListUserEncounters(db, username)
	if err != nil {
		panic(err)
	}

	wc := WebChar{
		SessionUser: username,
		IsLoggedIn:  loggedIn,
		IsAdmin:     isAdmin,
		Encounters:  encounters,
	}
	Render(w, "templates/user_encounter_index.html", wc)
}

// EncounterHandler renders a character in a Web page
func EncounterHandler(w http.ResponseWriter, req *http.Request) {

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

	enc, err := database.SlugLoadEncounter(db, slug)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Unable to load Encounter")
		http.Redirect(w, req, "/notfound", http.StatusSeeOther)
		return
	}

	fmt.Println(enc)

	IsAuthor := false

	if username == enc.Author.UserName {
		IsAuthor = true
	}

	factionMap := map[string][]*models.CharacterModel{}

	for _, fac := range enc.Factions {
		cms, err := database.LoadFactionCharacterModels(db, fac.CharacterModelSlugs)
		if err != nil {
			log.Panic(err)
		}
		factionMap[fac.Name] = cms
	}

	wc := WebChar{
		Encounter:   enc,
		FactionMap:  factionMap,
		IsAuthor:    IsAuthor,
		IsLoggedIn:  loggedIn,
		SessionUser: username,
		IsAdmin:     isAdmin,
	}

	// Render page
	Render(w, "templates/view_encounter.html", wc)

}

// AddEncounterHandler creates a user-generated encounter
func AddEncounterHandler(w http.ResponseWriter, req *http.Request) {

	// Get session values or redirect to Login
	session, err := sessions.Store.Get(req, "session")

	if err != nil {
		log.Println("error identifying session")
		http.Redirect(w, req, "/login/", http.StatusFound)
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
		http.Redirect(w, req, "/", http.StatusFound)
		return
	}

	user, err := database.LoadUser(db, username)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, req, "/", http.StatusFound)
		return
	}

	factions, err := database.ListFactions(db)
	if err != nil {
		panic(err)
	}

	factionMap := map[string]*models.Faction{}

	for _, f := range factions {
		factionMap[f.Slug] = f
	}

	wc := WebChar{
		IsAuthor:    true,
		SessionUser: username,
		IsLoggedIn:  loggedIn,
		IsAdmin:     isAdmin,
		Counter:     numToArray(5),
		Factions:    factions,
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/add_encounter.html", wc)

	}

	if req.Method == "POST" { // POST

		err := req.ParseMultipartForm(MaxMemory)
		if err != nil {
			panic(err)
		}

		enc := &models.Encounter{}

		// Use Schema Decoder to create struct
		err = decoder.Decode(enc, req.PostForm)
		if err != nil {
			panic(err)
		}

		enc.Slug = slug.Make(fmt.Sprintf("%s-%s", username, enc.Name))

		enc.Author = user

		// Get faction Slugs
		for i := 1; i < 5; i++ {
			slug := req.FormValue(fmt.Sprintf("Faction-%d", i))
			if slug != "" {
				enc.Factions = append(enc.Factions, factionMap[slug])
			}
		}

		err = database.SaveEncounter(db, enc)
		if err != nil {
			log.Panic(err)
		} else {
			fmt.Println("Saved Encounter")
		}

		url := fmt.Sprintf("/view_encounter/%s", enc.Slug)

		http.Redirect(w, req, url, http.StatusFound)
	}
}

// ModifyEncounterHandler renders an editable Encounter in a Web page
func ModifyEncounterHandler(w http.ResponseWriter, req *http.Request) {

	// Get session values or redirect to Login
	session, err := sessions.Store.Get(req, "session")

	if err != nil {
		log.Println("error identifying session")
		http.Redirect(w, req, "/login/", http.StatusFound)
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

	enc, err := database.SlugLoadEncounter(db, fSlug)
	if err != nil {
		fmt.Println(err)
	}

	if enc.Author == nil {
		enc.Author = &models.User{
			UserName: "",
		}
	}

	// Validate that User == Author
	IsAuthor := false

	if username == enc.Author.UserName || isAdmin == "true" {
		IsAuthor = true
	} else {
		http.Redirect(w, req, "/", http.StatusFound)
		return
	}

	// Load factions for options
	factions, err := database.ListFactions(db)
	if err != nil {
		panic(err)
	}

	// Set mapping for factions
	factionMap := map[string]*models.Faction{}

	for _, f := range factions {
		factionMap[f.Slug] = f
	}

	// Add empty factions options if < 5
	if len(enc.Factions) < 5 {
		for i := len(enc.Factions); i < 5; i++ {
			tf := &models.Faction{
				Name: "",
				Slug: "",
			}
			enc.Factions = append(enc.Factions, tf)
		}
	}

	wc := WebChar{
		Encounter:   enc,
		Factions:    factions,
		IsAuthor:    IsAuthor,
		SessionUser: username,
		IsLoggedIn:  loggedIn,
		IsAdmin:     isAdmin,
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/modify_encounter.html", wc)

	}

	if req.Method == "POST" { // POST

		err := req.ParseMultipartForm(MaxMemory)
		if err != nil {
			panic(err)
		}

		err = decoder.Decode(enc, req.PostForm)
		if err != nil {
			panic(err)
		}

		enc.Slug = slug.Make(fmt.Sprintf("%s-%s", username, enc.Name))

		// reset encounter factions
		enc.Factions = []*models.Faction{}

		// Get faction Slugs
		for i := 1; i < 5; i++ {
			slug := req.FormValue(fmt.Sprintf("Faction-%d", i))
			if slug != "" {
				enc.Factions = append(enc.Factions, factionMap[slug])
			}
		}

		// Insert Encounter into App archive
		err = database.UpdateEncounter(db, enc)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("Saved")
		}

		fmt.Println(enc)

		url := fmt.Sprintf("/view_encounter/%s", enc.Slug)

		http.Redirect(w, req, url, http.StatusSeeOther)
	}
}

// DeleteEncounterHandler renders a character in a Web page
func DeleteEncounterHandler(w http.ResponseWriter, req *http.Request) {

	// Get session values or redirect to Login
	session, err := sessions.Store.Get(req, "session")

	if err != nil {
		log.Println("error identifying session")
		http.Redirect(w, req, "/login/", http.StatusFound)
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

	enc, err := database.SlugLoadEncounter(db, slug)
	if err != nil {
		fmt.Println(err)
	}

	// Validate that User == Author
	IsAuthor := false

	if username == enc.Author.UserName || isAdmin == "true" {
		IsAuthor = true
	} else {
		http.Redirect(w, req, "/", http.StatusFound)
		return
	}

	wc := WebChar{
		Encounter:   enc,
		IsAuthor:    IsAuthor,
		SessionUser: username,
		IsLoggedIn:  loggedIn,
		IsAdmin:     isAdmin,
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/delete_encounter.html", wc)

	}

	if req.Method == "POST" {

		err := database.DeleteEncounter(db, enc.ID)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("Deleted Encounter")
		}

		url := fmt.Sprint("/encounter_index/")

		http.Redirect(w, req, url, http.StatusSeeOther)
	}
}
