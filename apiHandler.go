package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/thewhitetulip/Tasks/sessions"

	"github.com/toferc/rq_web/database"
	"github.com/toferc/rq_web/models"
)

// GetCharacterModels handles the basic roster rendering for the app
func GetCharacterModels(w http.ResponseWriter, req *http.Request) {

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

	fmt.Println(loggedIn, isAdmin, username)

	fmt.Println(session)

	/*
		if username == "" {
			http.Redirect(w, req, "/", 302)
			return
		}
	*/

	cms, err := database.ListCharacterModels(db)
	if err != nil {
		log.Println(err)
	}

	json.NewEncoder(w).Encode(cms)
}

// GetUserCharacterModels handles the basic roster rendering for the app
func GetUserCharacterModels(w http.ResponseWriter, req *http.Request) {

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

	fmt.Println(loggedIn, isAdmin, username)

	fmt.Println(session)

	/*
		if username == "" {
			http.Redirect(w, req, "/", 302)
			return
		}
	*/

	vars := mux.Vars(req)
	idString := vars["id"]

	pk, err := strconv.Atoi(idString)
	if err != nil {
		pk = 0
		log.Println(err)
	}

	user, err := database.PKLoadUser(db, int64(pk))
	if err != nil {
		log.Println(err)
		fmt.Println("Unable to load User")
	}

	cms, err := database.ListUserCharacterModels(db, user.UserName)
	if err != nil {
		log.Println(err)
	}

	json.NewEncoder(w).Encode(cms)
}

// GetCharacterModel handles the basic roster rendering for the app
func GetCharacterModel(w http.ResponseWriter, req *http.Request) {

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

	fmt.Println(loggedIn, isAdmin, username)

	fmt.Println(session)

	/*
		if username == "" {
			http.Redirect(w, req, "/", 302)
			return
		}
	*/

	vars := mux.Vars(req)
	idString := vars["id"]

	pk, err := strconv.Atoi(idString)
	if err != nil {
		pk = 0
		log.Println(err)
	}

	cm, err := database.PKLoadCharacterModel(db, int64(pk))
	if err != nil {
		log.Println(err)
	}

	json.NewEncoder(w).Encode(cm)
}

// CreateCharacterModel adds a new character via the API
func CreateCharacterModel(w http.ResponseWriter, req *http.Request) {

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

	fmt.Println(loggedIn, isAdmin, username)

	fmt.Println(session)

	if username == "" {
		http.Redirect(w, req, "/", 302)
		return
	}

	vars := mux.Vars(req)
	idString := vars["id"]

	pk, err := strconv.Atoi(idString)
	if err != nil {
		pk = 0
		log.Println(err)
	}

	cm := &models.CharacterModel{}

	err = json.NewDecoder(req.Body).Decode(cm)
	if err != nil {
		log.Println(err)
		http.Redirect(w, req, "/", 302)
	}

	cm.ID = int64(pk)

	err = database.SaveCharacterModel(db, cm)
	if err != nil {
		log.Println(err)
	}

	json.NewEncoder(w).Encode(cm)
}

// UpdateCharacterModel handles the basic roster rendering for the app
func UpdateCharacterModel(w http.ResponseWriter, req *http.Request) {

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

	fmt.Println(loggedIn, isAdmin, username)

	fmt.Println(session)

	if username == "" {
		http.Redirect(w, req, "/", 302)
		return
	}

	vars := mux.Vars(req)
	idString := vars["id"]

	pk, err := strconv.Atoi(idString)
	if err != nil {
		pk = 0
		log.Println(err)
	}

	cm, err := database.PKLoadCharacterModel(db, int64(pk))

	err = json.NewDecoder(req.Body).Decode(cm)
	if err != nil {
		log.Println(err)
		http.Redirect(w, req, "/", 302)
	}

	err = database.UpdateCharacterModel(db, cm)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Updating Character")

	json.NewEncoder(w).Encode(cm)
}

// CharacterModelLikesHandler handles the basic roster rendering for the app
func CharacterModelLikesHandler(w http.ResponseWriter, req *http.Request) {

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

	fmt.Println(loggedIn, isAdmin, username)

	fmt.Println(session)

	if username == "" {
		http.Redirect(w, req, "/", 302)
		return
	}

	vars := mux.Vars(req)
	idString := vars["id"]

	pk, err := strconv.Atoi(idString)
	if err != nil {
		pk = 0
		log.Println(err)
	}

	cm, err := database.PKLoadCharacterModel(db, int64(pk))

	if req.Method == "PUT" {
		_, ok := cm.LikeData[username]
		if !ok {
			cm.LikeData[username] = &models.Like{
				UserName:  username,
				CreatedAt: time.Now(),
			}
			cm.Likes++
			fmt.Printf("%s liked %s at %s", username, cm.Character.Name, time.Now())
		} else {
			delete(cm.LikeData, username)
			cm.Likes--
			fmt.Printf("%s unliked %s at %s", username, cm.Character.Name, time.Now())
		}

		err = database.UpdateCharacterModel(db, cm)
		if err != nil {
			log.Println(err)
		}
		fmt.Println("Updated Character")
	} else {
		fmt.Println("No request")
	}

	json.NewEncoder(w).Encode(cm)

}

// DeleteCharacterModel deletes a character
func DeleteCharacterModel(w http.ResponseWriter, req *http.Request) {

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

	fmt.Println(loggedIn, isAdmin, username)

	fmt.Println(session)

	if username == "" {
		http.Redirect(w, req, "/", 302)
		return
	}

	vars := mux.Vars(req)
	idString := vars["id"]

	pk, err := strconv.Atoi(idString)
	if err != nil {
		pk = 0
		log.Println(err)
	}

	err = database.DeleteCharacterModel(db, int64(pk))
	if err != nil {
		log.Println(err)
	}

	json.NewEncoder(w).Encode("Deleted")
}
