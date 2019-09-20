package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/thewhitetulip/Tasks/sessions"
	"github.com/toferc/rq_web/database"
	"github.com/toferc/rq_web/models"
)

// CharacterSearchHandler renders the basic character roster page
func CharacterSearchHandler(w http.ResponseWriter, req *http.Request) {

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

	query := values["query"]

	characters, err := database.SearchCharacterModels(db, query)
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
		CharacterModels: characters,
		Query:           query,
	}

	Render(w, "templates/character_search_results.html", wc)
}
