package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/dghubble/gologin/google"
	"github.com/gorilla/mux"

	"github.com/thewhitetulip/Tasks/sessions"
	"github.com/toferc/rq_web/database"
	"github.com/toferc/rq_web/models"
)

// UserIndexHandler handles the basic roster rendering for the app
func UserIndexHandler(w http.ResponseWriter, req *http.Request) {

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

	fmt.Println(session)

	if isAdmin != "true" {
		http.Redirect(w, req, "/", 302)
		return
	}

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

	users, err := database.PaginateUsers(db, limit, offset)
	if err != nil {
		panic(err)
	}

	characters := 0

	for _, v := range users {
		characters += v.Characters
	}

	wu := WebUser{
		SessionUser:    username,
		IsLoggedIn:     loggedIn,
		IsAdmin:        isAdmin,
		Users:          users,
		UserCount:      len(users),
		CharacterCount: characters,
		Limit:          limit,
		Offset:         offset,
	}

	Render(w, "templates/admin_view_users.html", wu)
}

func googleLoginFunc() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {

		session, err := sessions.Store.Get(req, "session")

		ctx := req.Context()
		googleUser, err := google.UserFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Implement a success handler
		session.Values["loggedin"] = "true"
		session.Values["username"] = googleUser.Name

		user, err := database.LoadUser(db, googleUser.Name)
		if err != nil {
			// No user exist for Google ID - Create one
			user, err = database.CreateGoogleUser(db, googleUser.Name, googleUser.Email)
			if err != nil {
				fmt.Println(err)
			}
		}

		if user.IsAdmin {
			session.Values["isAdmin"] = "true"
		} else {
			session.Values["isAdmin"] = "false"
		}

		session.Save(req, w)
		log.Print("user ", googleUser.Name, " is authenticated")
		fmt.Println(session.Values)
		http.Redirect(w, req, "/", 302)
		return
	}
	return http.HandlerFunc(fn)
}

//LogoutFunc Implements the logout functionality
//Will delete the session information from the cookie Store
func LogoutFunc(w http.ResponseWriter, req *http.Request) {
	session, err := sessions.Store.Get(req, "session")
	if err == nil {
		if session.Values["loggedin"] != false {
			session.Values["loggedin"] = "false"
			session.Values["username"] = ""
			session.Values["isAdmin"] = "false"
			session.Save(req, w)
			fmt.Println("Logged Out")
		}
	}
	http.Redirect(w, req, "/", 302)
	// Redirect to main page
}

//LoginFunc implements the login functionality, will add a cookie to cookie Store
//to manage authentication
func LoginFunc(w http.ResponseWriter, req *http.Request) {
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

	wc := WebChar{
		SessionUser: username,
		IsLoggedIn:  loggedIn,
		IsAdmin:     isAdmin,
	}

	switch req.Method {
	case "GET":
		Render(w, "templates/login.html", wc)
	case "POST":
		log.Print("Inside POST")
		req.ParseForm()
		username := req.Form.Get("username")
		password := req.Form.Get("password")

		if (username != "" && password != "") && database.ValidUser(db, username, password) {
			session.Values["loggedin"] = "true"
			session.Values["username"] = username

			user, err := database.LoadUser(db, username)
			if err != nil {
				fmt.Println(err)
				http.Redirect(w, req, "/", 302)
				return
			}

			if user.IsAdmin {
				session.Values["isAdmin"] = "true"
			} else {
				session.Values["isAdmin"] = "false"
			}

			session.Save(req, w)
			log.Print("user ", username, " is authenticated")
			fmt.Println(session.Values)
			http.Redirect(w, req, "/", 302)
			return
		} else {
			log.Print("Invalid user " + username)
			Render(w, "templates/login.html", wc)
		}
	}
}

//SignUpFunc implements sign-up functionality
func SignUpFunc(w http.ResponseWriter, req *http.Request) {
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

	wc := WebChar{
		SessionUser: username,
		IsLoggedIn:  loggedIn,
		IsAdmin:     isAdmin,
	}

	if req.Method == "POST" {
		req.ParseForm()

		username := req.Form.Get("username")
		username = strings.TrimSpace(username)
		rawpassword := req.Form.Get("password")
		email := req.Form.Get("email")

		if len(username) < 2 || len(username) > 24 || len(rawpassword) < 2 {
			http.Redirect(w, req, "/signup/", 302)
			return
		}

		password, err := database.HashPassword(rawpassword)
		if err != nil {
			fmt.Println(err)
		}

		log.Println(username, password, email)

		u := models.User{
			UserName: username,
			Password: password,
			Email:    email,
		}

		err = database.SaveUser(db, &u)
		if err != nil {
			http.Error(w, "Unable to sign user up", http.StatusInternalServerError)
		} else {
			//Log in user and continue
			session.Values["loggedin"] = "true"
			session.Values["username"] = username
			session.Save(req, w)
			log.Print("user ", username, " is authenticated")
			fmt.Println(session.Values)
			http.Redirect(w, req, "/", 302)
			return
		}
	} else if req.Method == "GET" {
		Render(w, "templates/signup.html", wc)
	}
}

func notFound(w http.ResponseWriter, req *http.Request) {

	session, err := sessions.Store.Get(req, "session")

	if err != nil {
		log.Println("error identifying session")
		// in case of error
	}

	w.WriteHeader(http.StatusNotFound)

	// Prep for user authentication
	sessionMap := getUserSessionValues(session)

	username := sessionMap["username"]
	loggedIn := sessionMap["loggedin"]
	isAdmin := sessionMap["isAdmin"]

	wc := WebChar{
		SessionUser: username,
		IsLoggedIn:  loggedIn,
		IsAdmin:     isAdmin,
	}

	Render(w, "templates/404.html", wc)
}
