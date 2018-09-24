package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/thewhitetulip/Tasks/sessions"
	"github.com/toferc/oneroll"
	"github.com/toferc/ore_web_roller/database"
)

// AddSkillHandler renders a character in a Web page
func AddSkillHandler(w http.ResponseWriter, req *http.Request) {

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
	s := vars["stat"]

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

	// Assign basic HyperSkill
	stat := c.Statistics[s]

	skill := &oneroll.Skill{
		Quality: &oneroll.Quality{
			Type: "",
		},
		LinkStat: stat,
		Dice: &oneroll.DiePool{
			Normal: 0,
			Hard:   0,
			Wiggle: 0,
		},
		Narrow:         false,
		Flexible:       false,
		Influence:      false,
		ReqSpec:        false,
		Specialization: "",
	}

	wc := WebChar{
		CharacterModel: cm,
		IsAuthor:       IsAuthor,
		SessionUser:    username,
		IsLoggedIn:     loggedIn,
		IsAdmin:        isAdmin,
		Skill:          skill,
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/add_skill.html", wc)

	}

	if req.Method == "POST" { // POST

		err := req.ParseForm()
		if err != nil {
			panic(err)
		}

		sName := req.FormValue("Name")

		sQuality := req.FormValue("Quality")

		nd, _ := strconv.Atoi(req.FormValue("Normal"))

		skill = new(oneroll.Skill)

		skill.Quality = &oneroll.Quality{
			Type: sQuality,
		}

		skill.Name = sName

		skill.LinkStat = stat

		skill.Dice = &oneroll.DiePool{
			Normal: nd,
		}

		if req.FormValue("Free") != "" {
			skill.Free = true
		}

		if req.FormValue("Narrow") != "" {
			skill.Narrow = true
		}

		if req.FormValue("Flexible") != "" {
			skill.Flexible = true
		}

		if req.FormValue("Influence") != "" {
			skill.Influence = true
		}

		if req.FormValue("ReqSpec") == "Yes" {
			skill.ReqSpec = true
			skill.Specialization = req.FormValue("Specialization")
		}

		c.Skills[sName] = skill

		fmt.Println(c)

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
