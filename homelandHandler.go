package main

import (
	"fmt"
	"github.com/toferc/runequest"
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

// HomelandHandler renders a character in a Web page
func HomelandHandler(w http.ResponseWriter, req *http.Request) {

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
	pk := vars["id"]

	if len(pk) == 0 {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}

	id, err := strconv.Atoi(pk)
	if err != nil {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}

	hl, err := database.PKLoadHomelandModel(db, int64(id))
	if err != nil {
		fmt.Println(err)
		fmt.Println("Unable to load HomelandModel")
	}

	fmt.Println(hl)

	IsAuthor := false

	if username == hl.Author.UserName {
		IsAuthor = true
	}

	if hl.Image == nil {
		hl.Image = new(models.Image)
		hl.Image.Path = DefaultCharacterPortrait
	}

	wc := WebChar{
		HomelandModel: hl,
		IsAuthor:      IsAuthor,
		IsLoggedIn:    loggedIn,
		SessionUser:   username,
		IsAdmin:       isAdmin,
		Skills:        runequest.Skills,
		CategoryOrder: runequest.CategoryOrder,
	}

	// Render page
	Render(w, "templates/view_homeland.html", wc)

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

	c := runequest.NewCharacter("name")
	cm := models.CharacterModel{
		Character: c,
	}

	c.Statistics["STR"].Value = 10
	c.Statistics["DEX"].Value = 10
	c.Statistics["INT"].Value = 10
	c.Statistics["POW"].Value = 10
	c.Statistics["CHA"].Value = 10
	c.Statistics["SIZ"].Value = 10

	wc := WebChar{
		CharacterModel: &cm,
		IsAuthor:       true,
		SessionUser:    username,
		IsLoggedIn:     loggedIn,
		IsAdmin:        isAdmin,
		Counter:        []int{1, 2, 3},
		Passions:       runequest.PassionTypes,
		CategoryOrder:  runequest.CategoryOrder,
		Skills:         runequest.Skills,
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/add_homeland.html", wc)

	}

	if req.Method == "POST" { // POST

		err := req.ParseMultipartForm(MaxMemory)
		if err != nil {
			panic(err)
		}

		user := database.LoadUser(db, username)

		// Map default Homeland to Character.Homelands
		hl := models.HomelandModel{
			Author: user,
			Homeland: &runequest.Homeland{
				Name:        req.FormValue("Name"),
				Description: req.FormValue("Description"),
			},
		}

		// Insert power into App archive if user authorizes
		if req.FormValue("Archive") != "" {
			hl.Open = true
		} else {
			hl.Open = false
		}

		// Upload image to s3
		file, h, err := req.FormFile("image")
		switch err {
		case nil:
			// Process image
			defer file.Close()
			// example path media/Major/TestImage/Jason_White.jpg
			path := fmt.Sprintf("/media/%s/%s/%s",
				hl.Author.UserName,
				runequest.ToSnakeCase(hl.Homeland.Name),
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

			hl.Image = new(models.Image)
			hl.Image.Path = path

			fmt.Println(path)

		case http.ErrMissingFile:
			log.Println("no file")

		default:
			log.Panic(err)
			fmt.Println("Error getting file ", err)
		}

		// Read Skills
		for _, s := range c.Skills {

			sk := runequest.Skill{
				CoreString: s.CoreString,
				UserChoice: s.UserChoice,
				Category:   s.Category,
			}

			str := fmt.Sprintf("%s-Base", s.CoreString)
			base, err := strconv.Atoi(req.FormValue(str))
			if err != nil {
				base = 0
			}
			sk.Base = base

			str = fmt.Sprintf("%s-Value", s.CoreString)
			base, err = strconv.Atoi(req.FormValue(str))
			if err != nil {
				base = 0
			}
			sk.HomelandValue = base

			if s.UserChoice {
				sk.UserString = req.FormValue(fmt.Sprintf("%s-UserString", s.CoreString))
			}

			hl.Homeland.Skills = append(hl.Homeland.Skills, sk)
		}

		// Read passions
		for i := 1; i < 4; i++ {

			p := runequest.Ability{
				Type:       "Passion",
				CoreString: req.FormValue(fmt.Sprintf("Passion-%d-CoreString", i)),
			}

			str := fmt.Sprintf("Passion-%d-Base", i)
			base, err := strconv.Atoi(req.FormValue(str))
			if err != nil {
				base = 0
			}
			p.Base = base

			userString := fmt.Sprintf("Passion-%d-UserString", i)

			if userString != "" {
				p.UserChoice = true
				p.UserString = req.FormValue(userString)
			}

			hl.Homeland.PassionList = append(hl.Homeland.PassionList, p)
		}

		// Read SkillChoices
		for i := 1; i < 4; i++ {

			sc := runequest.SkillChoice{}

			s1coreString := req.FormValue(fmt.Sprintf("Skill-%d-1-CoreString", i))
			s2coreString := req.FormValue(fmt.Sprintf("Skill-%d-2-CoreString", i))

			if s1coreString != "" && s2coreString != "" {

				s1baseSkill := runequest.Skills[s1coreString]
				fmt.Println(s1baseSkill)

				// First Skill option
				s1 := runequest.Skill{
					CoreString: s1baseSkill.CoreString,
					UserChoice: s1baseSkill.UserChoice,
					Category:   s1baseSkill.Category,
				}

				str := fmt.Sprintf("Skill-%d-1-Base", i)
				base, err := strconv.Atoi(req.FormValue(str))
				if err != nil {
					base = 0
				}
				s1.Base = base

				userString := fmt.Sprintf("Skill-%d-1-UserString", i)

				if userString != "" {
					s1.UserChoice = true
					s1.UserString = req.FormValue(userString)
				}

				// Second Skill option
				s2baseSkill := runequest.Skills[s2coreString]
				fmt.Println(s2baseSkill)

				// First Skill option
				s2 := runequest.Skill{
					CoreString: s2baseSkill.CoreString,
					UserChoice: s2baseSkill.UserChoice,
					Category:   s2baseSkill.Category,
				}

				str = fmt.Sprintf("Skill-%d-2-Base", i)
				base, err = strconv.Atoi(req.FormValue(str))
				if err != nil {
					base = 0
				}
				s2.Base = base

				userString = fmt.Sprintf("Skill-%d-2-UserString", i)

				if userString != "" {
					s2.UserChoice = true
					s2.UserString = userString
				}

				// Form SkillChoice
				sc.Skills = []runequest.Skill{
					s1,
					s2,
				}
			}

			hl.Homeland.SkillChoices = append(hl.Homeland.SkillChoices, sc)

		}
		// Add other

		// Skill loop based on runequest skills + 10 empty values
		// Create character, show skills, then pull skills in based on changes in base or additions to value

		author := database.LoadUser(db, username)

		hl.Author = author

		err = database.SaveHomelandModel(db, &hl)
		if err != nil {
			log.Panic(err)
		} else {
			fmt.Println("Saved")
		}

		url := fmt.Sprintf("/view_homeland/%d", hl.ID)

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

		hl.Homeland.Name = hlName

		// Insert Homeland into App archive
		err = database.UpdateHomelandModel(db, hl)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("Saved")
		}

		fmt.Println(hl)

		url := fmt.Sprintf("/view_Homeland/%d", hl.ID)

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
