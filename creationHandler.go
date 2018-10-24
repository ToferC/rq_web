package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gorilla/mux"
	"github.com/thewhitetulip/Tasks/sessions"
	"github.com/toferc/rq_web/database"
	"github.com/toferc/rq_web/models"
	"github.com/toferc/runequest"
)

// ChooseHomelandHandler allows users to name and select a homeland
func ChooseHomelandHandler(w http.ResponseWriter, req *http.Request) {

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

	cm := models.CharacterModel{}

	c := runequest.NewCharacter("")

	author := database.LoadUser(db, username)
	fmt.Println(author)

	cm = models.CharacterModel{
		Character: c,
		Author:    author,
	}

	homelands, err := database.ListHomelandModels(db)
	if err != nil {
		panic(err)
	}

	wc := WebChar{
		CharacterModel: &cm,
		SessionUser:    username,
		IsLoggedIn:     loggedIn,
		IsAdmin:        isAdmin,
		HomelandModels: homelands,
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/cc1_choose_homeland.html", wc)

	}

	if req.Method == "POST" {

		err := req.ParseMultipartForm(MaxMemory)
		if err != nil {
			panic(err)
		}

		c.Name = req.FormValue("Name")
		c.Description = req.FormValue("Description")

		// Set Homeland
		hlStr := req.FormValue("Homeland")

		fmt.Println("Results: " + hlStr)

		hlID, err := strconv.Atoi(hlStr)
		if err != nil {
			hlID = 0
			fmt.Println(err)
		}

		hl, err := database.PKLoadHomelandModel(db, int64(hlID))
		if err != nil {
			fmt.Println(err)
		}

		c.Homeland = hl.Homeland

		// Upload image to s3
		file, h, err := req.FormFile("image")
		switch err {
		case nil:
			// Prcless image
			defer file.Close()
			// example path media/Major/TestImage/Jason_White.jpg
			path := fmt.Sprintf("/media/%s/%s/%s",
				cm.Author.UserName,
				runequest.ToSnakeCase(cm.Character.Name),
				h.Filename,
			)

			_, err = uploader.Upload(&s3manager.UploadInput{
				Bucket: aws.String(os.Getenv("BUCKET")),
				Key:    aws.String(path),
				Body:   file,
			})
			if err != nil {
				log.Panic(err)
				fmt.Println("Error uploading file ", err)
			}
			fmt.Printf("successfully uploaded %q to %q\n",
				h.Filename, os.Getenv("BUCKET"))

			cm.Image = new(models.Image)
			cm.Image.Path = path

			fmt.Println(path)

		case http.ErrMissingFile:
			log.Println("no file")

		default:
			log.Panic(err)
			fmt.Println("Error getting file ", err)
		}

		//fmt.Println(c)

		err = database.SaveCharacterModel(db, &cm)
		if err != nil {
			log.Panic(err)
		} else {
			fmt.Println("Saved")
		}

		url := fmt.Sprintf("/cc2_choose_runes/%d", cm.ID)

		http.Redirect(w, req, url, http.StatusSeeOther)
	}
}

// PersonalHistoryHandler renders a character in a Web page
func PersonalHistoryHandler(w http.ResponseWriter, req *http.Request) {

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
		fmt.Println(err)
		fmt.Println("Unable to load CharacterModel")
	}

	IsAuthor := false

	if username == cm.Author.UserName {
		IsAuthor = true
	}

	occupations, err := database.ListOccupationModels(db)
	if err != nil {
		panic(err)
	}

	wc := WebChar{
		CharacterModel:   cm,
		OccupationModels: occupations,
		SessionUser:      username,
		IsLoggedIn:       loggedIn,
		IsAdmin:          isAdmin,
		IsAuthor:         IsAuthor,
		Skills:           runequest.Skills,
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/cc12_personal_history.html", wc)

	}

	if req.Method == "POST" {

		err := req.ParseMultipartForm(MaxMemory)
		if err != nil {
			panic(err)
		}

		// Do Stuff

		err = database.UpdateCharacterModel(db, cm)
		if err != nil {
			log.Panic(err)
		} else {
			fmt.Println("Saved")
		}

		url := fmt.Sprintf("/cc2_choose_runes/%d", cm.ID)

		http.Redirect(w, req, url, http.StatusSeeOther)
	}
}

// ChooseRunesHandler renders a character in a Web page
func ChooseRunesHandler(w http.ResponseWriter, req *http.Request) {

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
		fmt.Println(err)
		fmt.Println("Unable to load CharacterModel")
	}

	IsAuthor := false

	if username == cm.Author.UserName {
		IsAuthor = true
	}

	c := cm.Character

	wc := WebChar{
		CharacterModel: cm,
		SessionUser:    username,
		IsLoggedIn:     loggedIn,
		IsAdmin:        isAdmin,
		IsAuthor:       IsAuthor,
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/cc2_choose_runes.html", wc)

	}

	if req.Method == "POST" {

		err := req.ParseMultipartForm(MaxMemory)
		if err != nil {
			panic(err)
		}

		// Update Elemental Runes
		for k := range c.ElementalRunes {
			b, err := strconv.Atoi(req.FormValue(k))
			if err != nil {
				fmt.Println(err)
				b = 0
			}

			cbv, err := strconv.Atoi(req.FormValue("Bonus-" + k))
			if err != nil {
				fmt.Println(err)
				cbv = 0
			}

			// Modify Ability
			c.ElementalRunes[k].Base = b
			c.ElementalRunes[k].CreationBonusValue = cbv
		}

		for k, v := range c.PowerRunes {
			b, _ := strconv.Atoi(req.FormValue(k))
			if err != nil {
				b = 0
			}

			cbv, _ := strconv.Atoi(req.FormValue("Bonus-" + k))
			if err != nil {
				cbv = 0
			}
			// Modify Ability
			v.Base = b
			v.CreationBonusValue = cbv
			v.UpdateAbility()
		}

		// Determine opposing Power Rune values
		for k, v := range c.PowerRunes {
			if v.Total > 50 {
				c.UpdateOpposedRune(c.PowerRunes[k])
			}
		}

		fmt.Println(c)

		c.AddRuneModifiers()

		err = database.UpdateCharacterModel(db, cm)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("Saved")
		}

		url := fmt.Sprintf("/cc3_roll_stats/%d", cm.ID)

		http.Redirect(w, req, url, http.StatusSeeOther)
	}
}

// RollStatisticsHandler renders a character in a Web page
func RollStatisticsHandler(w http.ResponseWriter, req *http.Request) {

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
		fmt.Println(err)
		fmt.Println("Unable to load CharacterModel")
	}

	IsAuthor := false

	if username == cm.Author.UserName {
		IsAuthor = true
	}

	c := cm.Character

	occupations, err := database.ListOccupationModels(db)
	if err != nil {
		panic(err)
	}

	wc := WebChar{
		CharacterModel:   cm,
		OccupationModels: occupations,
		SessionUser:      username,
		IsLoggedIn:       loggedIn,
		IsAdmin:          isAdmin,
		IsAuthor:         IsAuthor,
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/cc3_roll_stats.html", wc)

	}

	if req.Method == "POST" {

		err := req.ParseMultipartForm(MaxMemory)
		if err != nil {
			panic(err)
		}

		for _, st := range runequest.StatMap {
			n, _ := strconv.Atoi(req.FormValue(st))
			if err != nil {
				n = 0
			}
			c.Statistics[st].Base = n

			n, _ = strconv.Atoi(req.FormValue("Rune-Bonus-" + st))
			if err != nil {
				n = 0
			}
			c.Statistics[st].RuneBonus = n
		}

		// Set Occupation
		ocStr := req.FormValue("Occupation")

		ocID, err := strconv.Atoi(ocStr)
		if err != nil {
			ocID = 0
		}

		oc, err := database.PKLoadOccupationModel(db, int64(ocID))
		if err != nil {
			fmt.Println("No Occupation Found")
		}

		c.Occupation = oc.Occupation

		fmt.Println(c)

		err = database.UpdateCharacterModel(db, cm)
		if err != nil {
			log.Panic(err)
		} else {
			fmt.Println("Saved")
		}

		url := fmt.Sprintf("/cc4_apply_homeland/%d", cm.ID)

		http.Redirect(w, req, url, http.StatusSeeOther)
	}
}

// ApplyHomelandHandler renders a character in a Web page
func ApplyHomelandHandler(w http.ResponseWriter, req *http.Request) {

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
		fmt.Println(err)
		fmt.Println("Unable to load CharacterModel")
	}

	c := cm.Character

	IsAuthor := false

	if username == cm.Author.UserName {
		IsAuthor = true
	}

	occupations, err := database.ListOccupationModels(db)
	if err != nil {
		panic(err)
	}

	wc := WebChar{
		CharacterModel:   cm,
		OccupationModels: occupations,
		SessionUser:      username,
		IsLoggedIn:       loggedIn,
		IsAdmin:          isAdmin,
		IsAuthor:         IsAuthor,
		Skills:           runequest.Skills,
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/cc4_apply_homeland.html", wc)

	}

	if req.Method == "POST" {

		err := req.ParseMultipartForm(MaxMemory)
		if err != nil {
			panic(err)
		}

		// Remove original homeland
		c.RemoveHomeland()

		// Do Stuff
		for _, s := range c.Homeland.Skills {
			// Or do we just manually modify the skills here?
			str := req.FormValue(fmt.Sprintf("%s-UserString", s.CoreString))
			if str != "" {
				s.UserString = str
			}
			c.ModifySkill(s)
		}

		// Homelands grant 3 base passions
		// Find number of abilities

		for _, selected := range c.Homeland.PassionList {
			c.Homeland.Passions = append(c.Homeland.Passions, selected)
			c.ModifyAbility(selected)
		}

		// Homeland grants a bonus to a rune affinity
		c.ElementalRunes[c.Homeland.RuneBonus].HomelandValue += 10

		for i, sc := range c.Homeland.SkillChoices {
			for m := range sc.Skills {
				str := fmt.Sprintf("SC-%d-%d", i, m)
				if req.FormValue(str) != "" {
					c.ApplySkillChoice(sc, m)
				}
			}
		}

		err = database.UpdateCharacterModel(db, cm)
		if err != nil {
			log.Panic(err)
		} else {
			fmt.Println("Saved")
		}

		url := fmt.Sprintf("/cc5_apply_occupation/%d", cm.ID)

		http.Redirect(w, req, url, http.StatusSeeOther)
	}
}

// ApplyOccupationHandler renders a character in a Web page
func ApplyOccupationHandler(w http.ResponseWriter, req *http.Request) {

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
		fmt.Println(err)
		fmt.Println("Unable to load CharacterModel")
	}

	IsAuthor := false

	if username == cm.Author.UserName {
		IsAuthor = true
	}

	occupations, err := database.ListOccupationModels(db)
	if err != nil {
		panic(err)
	}

	wc := WebChar{
		CharacterModel:   cm,
		OccupationModels: occupations,
		SessionUser:      username,
		IsLoggedIn:       loggedIn,
		IsAdmin:          isAdmin,
		IsAuthor:         IsAuthor,
		Skills:           runequest.Skills,
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/cc5_apply_occupation.html", wc)

	}

	if req.Method == "POST" {

		err := req.ParseMultipartForm(MaxMemory)
		if err != nil {
			panic(err)
		}

		// Do Stuff

		err = database.UpdateCharacterModel(db, cm)
		if err != nil {
			log.Panic(err)
		} else {
			fmt.Println("Saved")
		}

		url := fmt.Sprintf("/cc6_apply_cult/%d", cm.ID)

		http.Redirect(w, req, url, http.StatusSeeOther)
	}
}

// ApplyCultHandler renders a character in a Web page
func ApplyCultHandler(w http.ResponseWriter, req *http.Request) {

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
		fmt.Println(err)
		fmt.Println("Unable to load CharacterModel")
	}

	IsAuthor := false

	if username == cm.Author.UserName {
		IsAuthor = true
	}

	occupations, err := database.ListOccupationModels(db)
	if err != nil {
		panic(err)
	}

	wc := WebChar{
		CharacterModel:   cm,
		OccupationModels: occupations,
		SessionUser:      username,
		IsLoggedIn:       loggedIn,
		IsAdmin:          isAdmin,
		IsAuthor:         IsAuthor,
		Skills:           runequest.Skills,
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/cc6_apply_cult.html", wc)

	}

	if req.Method == "POST" {

		err := req.ParseMultipartForm(MaxMemory)
		if err != nil {
			panic(err)
		}

		// Do Stuff

		err = database.UpdateCharacterModel(db, cm)
		if err != nil {
			log.Panic(err)
		} else {
			fmt.Println("Saved")
		}

		url := fmt.Sprintf("/cc7_personal_skills/%d", cm.ID)

		http.Redirect(w, req, url, http.StatusSeeOther)
	}
}

// PersonalSkillsHandler renders a character in a Web page
func PersonalSkillsHandler(w http.ResponseWriter, req *http.Request) {

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
		fmt.Println(err)
		fmt.Println("Unable to load CharacterModel")
	}

	IsAuthor := false

	if username == cm.Author.UserName {
		IsAuthor = true
	}

	occupations, err := database.ListOccupationModels(db)
	if err != nil {
		panic(err)
	}

	wc := WebChar{
		CharacterModel:   cm,
		OccupationModels: occupations,
		SessionUser:      username,
		IsLoggedIn:       loggedIn,
		IsAdmin:          isAdmin,
		IsAuthor:         IsAuthor,
		Skills:           runequest.Skills,
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/cc7_personal_skills.html", wc)

	}

	if req.Method == "POST" {

		err := req.ParseMultipartForm(MaxMemory)
		if err != nil {
			panic(err)
		}

		// Do Stuff

		err = database.UpdateCharacterModel(db, cm)
		if err != nil {
			log.Panic(err)
		} else {
			fmt.Println("Saved")
		}

		url := fmt.Sprintf("/view_character/%d", cm.ID)

		http.Redirect(w, req, url, http.StatusSeeOther)
	}
}
