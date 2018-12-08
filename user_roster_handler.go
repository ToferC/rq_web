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

	characters, err := database.ListUserCharacterModels(db, username)
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
	}

	Render(w, "templates/user_roster.html", wc)
}

// AddToUserRosterHandler adds an open charactermodel to the individual user roster
func AddToUserRosterHandler(w http.ResponseWriter, req *http.Request) {

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

	author, err := database.LoadUser(db, username)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, req, "/", 302)
	}

	// Get variables from URL
	vars := mux.Vars(req)
	pk := vars["id"]

	if len(pk) == 0 {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}

	id, err := strconv.Atoi(pk)
	if err != nil {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}

	// Load CharacterModel
	cm, err := database.PKLoadCharacterModel(db, int64(id))
	if err != nil {
		fmt.Println(err)
	}

	newCharacterModel := models.CharacterModel{
		Author:    author,
		Character: cm.Character,
		Open:      false,
	}

	err = database.SaveCharacterModel(db, &newCharacterModel)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Saved")
	}

	url := fmt.Sprintf("/view_character/%d", newCharacterModel.ID)

	http.Redirect(w, req, url, 302)
}
