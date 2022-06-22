package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/gosimple/slug"
	"github.com/thewhitetulip/Tasks/sessions"
	"github.com/toferc/rq_web/database"
	"github.com/toferc/rq_web/models"
	"github.com/toferc/runequest"
)

// NoteIndexHandler renders a list of notes under a character
func NoteIndexHandler(w http.ResponseWriter, req *http.Request) {

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

	// Get variables from URL
	vars := mux.Vars(req)
	pk := vars["id"]

	if len(pk) == 0 {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	id, err := strconv.Atoi(pk)
	if err != nil {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	cm, err := database.PKLoadCharacterModel(db, int64(id))
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(cm.Character.Name, cm.ID)

	notes, err := database.ListNotes(db, cm.ID)
	if err != nil {
		fmt.Println(err)
	}

	IsAuthor := false

	if username == cm.Author.UserName {
		IsAuthor = true
	}

	wv := WebView{
		SessionUser:    username,
		CharacterModel: cm,
		IsLoggedIn:     loggedIn,
		IsAuthor:       IsAuthor,
		IsAdmin:        isAdmin,
		Notes:          notes,
	}

	Render(w, "templates/notes_index.html", wv)
}

// NoteHandler renders a note in a Web page
func NoteHandler(w http.ResponseWriter, req *http.Request) {

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

	nt, err := database.SlugLoadNote(db, slug)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Unable to load Note")
	}

	cm, err := database.PKLoadCharacterModel(db, nt.CharacterModelID)
	if err != nil {
		fmt.Println(err)
	}

	IsAuthor := false

	if username == nt.AuthorUserName {
		IsAuthor = true
	}

	wv := WebView{
		Note:           nt,
		IsAuthor:       IsAuthor,
		CharacterModel: cm,
		IsLoggedIn:     loggedIn,
		SessionUser:    username,
		IsAdmin:        isAdmin,
	}

	// Render page
	Render(w, "templates/view_note.html", wv)

}

// AddNoteHandler creates a user-generated note
func AddNoteHandler(w http.ResponseWriter, req *http.Request) {

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

	// Get variables from URL
	vars := mux.Vars(req)
	pk := vars["id"]

	if len(pk) == 0 {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	id, err := strconv.Atoi(pk)
	if err != nil {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	cm, err := database.PKLoadCharacterModel(db, int64(id))
	if err != nil {
		fmt.Println(err)
	}

	nt := models.Note{
		Title:            "New",
		Year:             1625,
		CharacterModelID: cm.ID,
		AuthorUserName:   username,
		Tags:             []string{"", "", ""},
	}

	wv := WebView{
		Note:           &nt,
		CharacterModel: cm,
		IsAuthor:       true,
		SessionUser:    username,
		IsLoggedIn:     loggedIn,
		IsAdmin:        isAdmin,
		Seasons:        models.Seasons,
		Weeks:          models.Weeks,
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/add_note.html", wv)

	}

	if req.Method == "POST" { // POST

		err := req.ParseForm()
		if err != nil {
			panic(err)
		}

		// Map default Note

		// Pull form values into Note via gorilla/schema
		err = decoder.Decode(&nt, req.PostForm)
		if err != nil {
			panic(err)
		}

		nt.PublishedOn = time.Now()

		fmt.Println(nt.Tags)

		nt.Slug = slug.Make(fmt.Sprintf("%s-%d-%s-%s-%s",
			cm.Character.Name, nt.Year, nt.Season, nt.Week, nt.Title))

		err = database.SaveNote(db, &nt)
		if err != nil {
			log.Panic(err)
		} else {
			fmt.Println("Saved Note")
			cm.Notes++
			err = database.UpdateCharacterModel(db, cm)
			if err != nil {
				log.Println(err)
			}
			fmt.Println("Updated CharacterModel")
		}

		url := fmt.Sprintf("/view_note/%s", nt.Slug)

		http.Redirect(w, req, url, http.StatusFound)
	}
}

// ModifyNoteHandler renders an editable Note in a Web page
func ModifyNoteHandler(w http.ResponseWriter, req *http.Request) {

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

	nt, err := database.SlugLoadNote(db, slug)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, req, "/", http.StatusFound)
		return
	}

	if len(nt.Tags) < 3 {
		for i := 0; i < 3; i++ {
			nt.Tags = append(nt.Tags, "")
		}
	}

	cm, err := database.PKLoadCharacterModel(db, nt.CharacterModelID)
	if err != nil {
		fmt.Println(err)
		cm = &models.CharacterModel{
			Character: &runequest.Character{},
			Author:    &models.User{},
		}
	}

	// Validate that User == Author
	IsAuthor := false

	if username == nt.AuthorUserName || isAdmin == "true" {
		IsAuthor = true
	} else {
		http.Redirect(w, req, "/", http.StatusFound)
		return
	}

	wv := WebView{
		Note:           nt,
		IsAuthor:       IsAuthor,
		CharacterModel: cm,
		SessionUser:    username,
		IsLoggedIn:     loggedIn,
		IsAdmin:        isAdmin,
		Seasons:        models.Seasons,
		Weeks:          models.Weeks,
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/modify_note.html", wv)

	}

	if req.Method == "POST" { // POST

		err := req.ParseForm()
		if err != nil {
			panic(err)
		}

		// Pull form values into Note via gorilla/schema
		err = decoder.Decode(nt, req.PostForm)
		if err != nil {
			panic(err)
		}

		nt.CharacterModelID = cm.ID
		nt.AuthorID = cm.Author.ID
		nt.AuthorUserName = username

		// Insert Note into App archive
		err = database.UpdateNote(db, nt)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("Saved")
		}

		fmt.Println(nt)

		url := fmt.Sprintf("/view_note/%s", nt.Slug)

		http.Redirect(w, req, url, http.StatusSeeOther)
	}
}

// DeleteNoteHandler renders a character in a Web page
func DeleteNoteHandler(w http.ResponseWriter, req *http.Request) {

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
	pk := vars["id"]

	id, err := strconv.Atoi(pk)
	if err != nil {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	nt, err := database.PKLoadNote(db, int64(id))
	if err != nil {
		fmt.Println(err)
	}

	cm, err := database.PKLoadCharacterModel(db, nt.CharacterModelID)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(nt)

	// Validate that User == Author
	IsAuthor := false

	if username == nt.AuthorUserName || isAdmin == "true" {
		IsAuthor = true
	} else {
		http.Redirect(w, req, "/", http.StatusFound)
		return
	}

	wv := WebView{
		Note:           nt,
		CharacterModel: cm,
		IsAuthor:       IsAuthor,
		SessionUser:    username,
		IsLoggedIn:     loggedIn,
		IsAdmin:        isAdmin,
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/delete_note.html", wv)

	}

	if req.Method == "POST" {

		err := database.DeleteNote(db, nt.ID)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("Deleted Note")
			cm.Notes--
			err = database.UpdateCharacterModel(db, cm)
			if err != nil {
				log.Println(err)
			}
			fmt.Println("Updated CharacterModel")
		}

		url := fmt.Sprint("/")

		http.Redirect(w, req, url, http.StatusSeeOther)
	}
}
