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

// TextSummaryHandler renders a character in a Web page
func TextSummaryHandler(w http.ResponseWriter, req *http.Request) {

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

	session.Save(req, w)

	vars := mux.Vars(req)
	pk := vars["id"]

	if len(pk) == 0 {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}

	id, err := strconv.Atoi(pk)
	if err != nil {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}

	cm, err := database.PKLoadCharacterModel(db, int64(id))
	if err != nil {
		fmt.Println("Unable to load CharacterModel")
		http.Redirect(w, req, "/notfound", http.StatusSeeOther)
		return
	}

	stringArray := []string{}

	stringArray = append(stringArray, cm.Character.String())

	stringArray = append(stringArray, cm.Character.StatBlock())

	IsAuthor := false

	if username == cm.Author.UserName {
		IsAuthor = true
	}

	if cm.Image == nil {
		cm.Image = new(models.Image)
		cm.Image.Path = DefaultCharacterPortrait
	}

	//c.DetermineSkillCategoryValues()

	wc := WebChar{
		CharacterModel: cm,
		IsAuthor:       IsAuthor,
		StringArray:    stringArray,
		IsLoggedIn:     loggedIn,
		SessionUser:    username,
		IsAdmin:        isAdmin,
	}

	// Render page
	Render(w, "templates/text_summary.html", wc)

}
