package main

import (
	"fmt"
	"github.com/toferc/runequest"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/thewhitetulip/Tasks/sessions"
	"github.com/toferc/rq_web/database"
	"github.com/toferc/rq_web/models"
)

const baseDieString string = "1ac+4d+0hd+0wd+0gf+0sp+1nr+0ed"
const blankDieString string = "ac+d+hd+wd+gf+sp+nr+ed"

// RollHandler generates a Web user interface
func RollHandler(w http.ResponseWriter, req *http.Request) {

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

	pk := mux.Vars(req)["id"]

	id, err := strconv.Atoi(pk)
	if err != nil {
		id = 9999
	}

	cm, err := database.PKLoadCharacterModel(db, int64(id))
	if err != nil {
		fmt.Println(err)
		cm = new(models.CharacterModel)
	}

	if cm.Character == nil {
		c := runequest.NewCharacter("Player1")
		cm.Character = c
	}

	wv := WebView{
		Actor:       []*models.CharacterModel{cm},
		SessionUser: username,
		IsLoggedIn:  loggedIn,
		IsAdmin:     isAdmin,
	}

	if req.Method == "GET" {

		Render(w, "templates/roller.html", wv)

		// wv.Rolls = []oneroll.Roll{}

	}

	if req.Method == "POST" {

		http.Redirect(w, req, "/", http.StatusSeeOther)
	}
}

// OpposeHandler generates a Web user interface
func OpposeHandler(w http.ResponseWriter, req *http.Request) {

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

	if req.Method == "GET" {

		wv := WebView{
			SessionUser: username,
			IsLoggedIn:  loggedIn,
			IsAdmin:     isAdmin,
		}

		Render(w, "templates/opposed.html", wv)

		// wv.Rolls = []oneroll.Roll{}

	}

	if req.Method == "POST" {
		// Submit form

		// Encode URL.Query
		//qs := q.Encode()

		http.Redirect(w, req, "/opposed/?", http.StatusSeeOther)
	}
}
