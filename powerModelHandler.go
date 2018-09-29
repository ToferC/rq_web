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
	"github.com/toferc/runequest"
)

// HomelandListHandler applies a Homeland template to a character
func HomelandListHandler(w http.ResponseWriter, req *http.Request) {

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

	// Validate that User == Author
	IsAuthor := false

	if username == cm.Author.UserName {
		IsAuthor = true
	} else {
		http.Redirect(w, req, "/", 302)
	}

	c := cm.Character

	homelands, err := database.ListHomelandModels(db)
	if err != nil {
		panic(err)
	}

	wc := WebChar{
		CharacterModel: cm,
		IsAuthor:       IsAuthor,
		SessionUser:    username,
		IsLoggedIn:     loggedIn,
		IsAdmin:        isAdmin,
		HomelandModels: homelands,
	}

	if req.Method == "GET" {

		// Render page

		Render(w, "templates/add_homeland_from_list.html", wc)

	}

	if req.Method == "POST" {

		err := req.ParseForm()
		if err != nil {
			panic(err)
		}

		homelandName := req.FormValue("Name")

		if homelandName != "" {
			hl := homelands[runequest.ToSnakeCase(homelandName)].Homeland

			fmt.Println(hl)

			c.Homeland = hl

			c.UpdateCharacter()
		}

		cm.Character.Homeland = 

		err = database.UpdateCharacterModel(db, cm)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("Saved")
		}

		fmt.Println(c)

		url := fmt.Sprintf("/view_character/%d", cm.ID)

		http.Redirect(w, req, url, http.StatusSeeOther)
	}
}

// AddHomelandHandler creates a user-generated homeland
func AddHomelandHandler(w http.ResponseWriter, req *http.Request) {

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

	// Create default Homeland to populate page
	defaultHomeland := &runequest.Homeland{
		Name: "",
	}

	// Map default Homeland to Character.Homelands
	hl := models.HomelandModel{
		Homeland: defaultHomeland,
	}

	wc := WebChar{
		HomelandModel: &hl,
		IsAuthor:      true,
		SessionUser:   username,
		IsLoggedIn:    loggedIn,
		IsAdmin:       isAdmin,
		Counter:       []int{1, 2, 3, 4, 5, 6, 7, 8},
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/add_homeland.html", wc)

	}

	if req.Method == "POST" { // POST

		err := req.ParseForm()
		if err != nil {
			panic(err)
		}

		homelandName := req.FormValue("Name")

		homeland := runequest.Homeland{
			Name: homelandName,
		}

		// Skill loop based on runequest skills + 10 empty values

		author := database.LoadUser(db, username)
		fmt.Println(author)

		hl.Author = author

		// Insert Homeland into App archive
		p.DetermineHomelandCapacities()
		p.CalculateCost()

		hl.Homeland = &p
		database.SaveHomelandModel(db, &hl)

		url := fmt.Sprintf("/view_Homeland/%d", hl.ID)

		http.Redirect(w, req, url, http.StatusFound)
	}
}

// ModifyHomelandHandler renders an editable Homeland in a Web page
func ModifyHomelandHandler(w http.ResponseWriter, req *http.Request) {

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
	pk := vars["id"]

	id, err := strconv.Atoi(pk)
	if err != nil {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}

	hl, err := database.PKLoadHomelandModel(db, int64(id))
	if err != nil {
		fmt.Println(err)
	}

	if hl.Author == nil {
		hl.Author = &models.User{
			UserName: "",
		}
	}

	// Validate that User == Author
	IsAuthor := false

	if username == hl.Author.UserName {
		IsAuthor = true
	} else {
		http.Redirect(w, req, "/", 302)
	}

	h := hl.Homeland

	// Assign additional empty skills & passions to populate form

	wc := WebChar{
		HomelandModel: hl,
		IsAuthor:      IsAuthor,
		SessionUser:   username,
		IsLoggedIn:    loggedIn,
		IsAdmin:       isAdmin,
		Counter:       []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/modify_homeland.html", wc)

	}

	if req.Method == "POST" { // POST

		err := req.ParseForm()
		if err != nil {
			panic(err)
		}

		hlName := req.FormValue("Name")

		hl.Name = hlName

		for _, qLoop := range wc.Counter[:qCount] { // Quality Loop

			qType := req.FormValue(fmt.Sprintf("Q%d-Type", qLoop))

			if qType != "" {
				l, err := strconv.Atoi(req.FormValue(fmt.Sprintf("Q%d-Level", qLoop)))
				if err != nil {
					l = 0
				}
				q := &runequest.Quality{
					Type:  req.FormValue(fmt.Sprintf("Q%d-Type", qLoop)),
					Level: l,
					Name:  req.FormValue(fmt.Sprintf("Q%d-Name", qLoop)),
				}

				for _, cLoop := range wc.Counter[:3] {
					cType := req.FormValue(fmt.Sprintf("Q%d-C%d-Type", qLoop, cLoop))
					if cType != "" {
						cap := &runequest.Capacity{
							Type: cType,
						}
						q.Capacities = append(q.Capacities, cap)
					}
				}

				// Append Quality to Homeland Qualities
				p.Qualities = append(p.Qualities, q)
			}
		}

		// Insert Homeland into App archive

		hlModel.Homeland = hl

		err = database.UpdateHomelandModel(db, hlModel)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("Saved")
		}

		fmt.Println(p)

		url := fmt.Sprintf("/view_Homeland/%d", hlModel.ID)

		http.Redirect(w, req, url, http.StatusSeeOther)
	}
}

// DeleteStandaloneHomelandHandler renders a character in a Web page
func DeleteStandaloneHomelandHandler(w http.ResponseWriter, req *http.Request) {

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
	pk := vars["id"]

	id, err := strconv.Atoi(pk)
	if err != nil {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}

	hl, err := database.PKLoadHomelandModel(db, int64(id))
	if err != nil {
		fmt.Println(err)
	}

	// Validate that User == Author
	IsAuthor := false

	if username == hl.Author.UserName {
		IsAuthor = true
	} else {
		http.Redirect(w, req, "/", 302)
	}

	wc := WebChar{
		HomelandModel: hl,
		IsAuthor:      IsAuthor,
		SessionUser:   username,
		IsLoggedIn:    loggedIn,
		IsAdmin:       isAdmin,
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/delete_standalone_Homeland.html", wc)

	}

	if req.Method == "POST" {

		err := database.DeleteHomelandModel(db, hl.ID)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("Deleted Homeland")
		}

		url := fmt.Sprint("/index_Homelands/")

		http.Redirect(w, req, url, http.StatusSeeOther)
	}
}
