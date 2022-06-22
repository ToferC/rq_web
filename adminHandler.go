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

// AdminUserRosterViewHandler handles user-specific rosters
func AdminUserRosterViewHandler(w http.ResponseWriter, req *http.Request) {

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
		http.Redirect(w, req, "/", http.StatusFound)
		return
	}

	if isAdmin != "true" {
		http.Redirect(w, req, "/", http.StatusFound)
		return
	}

	values := mux.Vars(req)
	pk := values["id"]

	if len(pk) == 0 {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	id, err := strconv.Atoi(pk)
	if err != nil {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
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

	characters, err := database.ListUserCharacterModels(db, u.UserName, limit, offset)
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
		User:            u,
		SessionUser:     username,
		IsLoggedIn:      loggedIn,
		IsAdmin:         isAdmin,
		CharacterModels: characters,
		Limit:           limit,
		Offset:          offset,
	}

	Render(w, "templates/admin_view_user_roster.html", wc)
}

// MakeAdminHandler handles the basic roster rendering for the app
func MakeAdminHandler(w http.ResponseWriter, req *http.Request) {

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

	isAdmin := sessionMap["isAdmin"]

	vars := mux.Vars(req)
	idString := vars["id"]

	pk, err := strconv.Atoi(idString)
	if err != nil {
		pk = 0
		log.Println(err)
	}

	fmt.Println(session)

	if isAdmin != "true" {
		http.Redirect(w, req, "/", http.StatusFound)
		return
	}

	user, err := database.PKLoadUser(db, int64(pk))
	if err != nil {
		log.Fatal(err)
		fmt.Println("Unable to load User")
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	if !user.IsAdmin {
		user.IsAdmin = true
	} else {
		user.IsAdmin = false
	}

	err = database.UpdateUser(db, user)
	if err != nil {
		log.Println(err)
	}

	url := "/user_index/"

	http.Redirect(w, req, url, http.StatusFound)
	return
}

// DeleteUserHandler handles the basic roster rendering for the app
func DeleteUserHandler(w http.ResponseWriter, req *http.Request) {

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

	values := mux.Vars(req)
	idString := values["id"]

	pk, err := strconv.Atoi(idString)
	if err != nil {
		pk = 0
		log.Println(err)
	}

	fmt.Println(session)

	if isAdmin != "true" {
		http.Redirect(w, req, "/", http.StatusFound)
		return
	}

	user, err := database.PKLoadUser(db, int64(pk))
	if err != nil {
		log.Fatal(err)
		fmt.Println("Unable to load User")
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	limit := 9999

	offset := 0

	cms, err := database.ListUserCharacterModels(db, user.UserName, limit, offset)
	if err != nil {
		log.Println(err)
	}

	wv := WebView{
		User:        user,
		IsLoggedIn:  loggedIn,
		SessionUser: username,
		IsAdmin:     isAdmin,
		Characters:  cms,
	}

	if req.Method == "GET" {
		Render(w, "templates/delete_user.html", wv)
	}

	if req.Method == "POST" {

		for _, cm := range cms {
			err = database.DeleteCharacterModel(db, cm.ID)
			if err != nil {
				log.Println(err)
			}
		}

		err = database.DeleteUser(db, user.ID)
		if err != nil {
			log.Println(err)
		}

		url := "/user_index/"

		http.Redirect(w, req, url, http.StatusFound)
		return
	}

}
