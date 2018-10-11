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

// HomelandIndexHandler renders the basic character roster page
func HomelandIndexHandler(w http.ResponseWriter, req *http.Request) {

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

	homelands, err := database.ListHomelandModels(db)
	if err != nil {
		panic(err)
	}

	for _, hl := range homelands {
		if hl.Image == nil {
			hl.Image = new(models.Image)
			hl.Image.Path = DefaultCharacterPortrait
		}
	}

	wc := WebChar{
		SessionUser:    username,
		IsLoggedIn:     loggedIn,
		IsAdmin:        isAdmin,
		HomelandModels: homelands,
	}

	Render(w, "templates/homeland_index.html", wc)
}

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

		// Insert Homeland into App archive if user authorizes
		if req.FormValue("Archive") != "" {
			hl.Open = true
		} else {
			hl.Open = false
		}

		if req.FormValue("Official") != "" {
			hl.Official = true
		} else {
			hl.Official = false
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

		// Read Base Skills
		for _, s := range c.Skills {

			// Build skill based on user input vs. base Skills
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
			v, err := strconv.Atoi(req.FormValue(str))
			if err != nil {
				v = 0
			}
			sk.HomelandValue = v

			if s.UserChoice {
				sk.UserString = req.FormValue(fmt.Sprintf("%s-UserString", s.CoreString))
			}

			if sk.Base > s.Base || sk.HomelandValue > 0 || sk.UserString != s.UserString {
				// If we changed something, add to the Homeland skill list
				hl.Homeland.Skills = append(hl.Homeland.Skills, sk)
			}

		}

		// Read passions
		for i := 1; i < 4; i++ {

			coreString := req.FormValue(fmt.Sprintf("Passion-%d-CoreString", i))

			if coreString != "" {

				p := runequest.Ability{
					Type:       "Passion",
					CoreString: coreString,
				}

				str := fmt.Sprintf("Passion-%d-Base", i)
				base, err := strconv.Atoi(req.FormValue(str))
				if err != nil {
					base = 0
				}
				p.Base = base

				userString := req.FormValue(fmt.Sprintf("Passion-%d-UserString", i))

				if userString != "" {
					p.UserChoice = true
					p.UserString = userString
				}

				hl.Homeland.PassionList = append(hl.Homeland.PassionList, p)
			}
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

				str := fmt.Sprintf("Skill-%d-1-Value", i)
				v, err := strconv.Atoi(req.FormValue(str))
				if err != nil {
					v = 0
				}
				s1.HomelandValue = v

				if s1.UserChoice {
					userString := fmt.Sprintf("Skill-%d-1-UserString", i)
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

				str = fmt.Sprintf("Skill-%d-2-Value", i)
				v, err = strconv.Atoi(req.FormValue(str))
				if err != nil {
					v = 0
				}
				s2.HomelandValue = v

				if s2.UserChoice {
					userString := fmt.Sprintf("Skill-%d-2-UserString", i)
					s2.UserString = userString
				}

				// Form SkillChoice
				sc.Skills = []runequest.Skill{
					s1,
					s2,
				}
				// Append skillchoice
				hl.Homeland.SkillChoices = append(hl.Homeland.SkillChoices, sc)
			}
		}

		// Read New Skills

		for i := 1; i < 4; i++ {

			coreString := req.FormValue(fmt.Sprintf("Skill-%d-CoreString", i))

			if coreString != "" {

				sk := runequest.Skill{}

				sk.CoreString = coreString
				sk.Category = req.FormValue(fmt.Sprintf("Skill-%d-Category", i))

				userString := req.FormValue(fmt.Sprintf("Skill-%d-UserString", i))

				if userString != "" {
					sk.UserChoice = true
					sk.UserString = userString
				}

				str := fmt.Sprintf("Skill-%d-Base", i)
				base, err := strconv.Atoi(req.FormValue(str))
				if err != nil {
					base = 0
				}
				sk.Base = base

				str = fmt.Sprintf("Skill-%d-Value", i)
				v, err := strconv.Atoi(req.FormValue(str))
				if err != nil {
					v = 0
				}
				sk.HomelandValue = v

				hl.Homeland.Skills = append(hl.Homeland.Skills, sk)
			}
		}

		// Add other HomelandModel fields

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

	// Create empty skills for adding to the Homeland

	// Add extra empty skills if < 20
	if len(hl.Homeland.Skills) < 20 {
		for i := len(hl.Homeland.Skills); i < 20; i++ {
			tempSkill := runequest.Skill{}
			hl.Homeland.Skills = append(hl.Homeland.Skills, tempSkill)
		}
	}

	// Add extra empty skillchoices if < 3
	// Something here isn't working
	if len(hl.Homeland.SkillChoices) < 4 {
		for i := len(hl.Homeland.SkillChoices); i < 4; i++ {
			tempSkillChoice := runequest.SkillChoice{
				Skills: []runequest.Skill{
					runequest.Skill{
						Name: "default",
					},
					runequest.Skill{
						Name: "default",
					},
				},
			}
			hl.Homeland.SkillChoices = append(hl.Homeland.SkillChoices, tempSkillChoice)
		}
	}

	wc := WebChar{
		HomelandModel: hl,
		IsAuthor:      IsAuthor,
		SessionUser:   username,
		IsLoggedIn:    loggedIn,
		IsAdmin:       isAdmin,
		Counter:       []int{1, 2, 3},
		CategoryOrder: runequest.CategoryOrder,
		Skills:        runequest.Skills,
		Passions:      runequest.PassionTypes,
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/modify_homeland.html", wc)

	}

	if req.Method == "POST" { // POST

		err := req.ParseMultipartForm(MaxMemory)
		if err != nil {
			panic(err)
		}

		hlName := req.FormValue("Name")

		hl.Homeland.Name = hlName

		// Update Homeland here

		// Insert Homeland into App archive if user authorizes
		if req.FormValue("Archive") != "" {
			hl.Open = true
		} else {
			hl.Open = false
		}

		// Open or Official
		if req.FormValue("Official") != "" {
			hl.Official = true
		} else {
			hl.Official = false
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

		// Skills

		tempSkills := []runequest.Skill{}

		// Read Base Skills from Homeland
		for _, s := range hl.Homeland.Skills {

			// Build skill based on user input vs. base Skills
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
			v, err := strconv.Atoi(req.FormValue(str))
			if err != nil {
				v = 0
			}
			sk.HomelandValue = v

			if s.UserChoice {
				sk.UserString = req.FormValue(fmt.Sprintf("%s-UserString", s.CoreString))
			}

			if sk.Base > s.Base || sk.HomelandValue > 0 || sk.UserString != s.UserString {
				// If we changed something, add to the Homeland skill list
				tempSkills = append(tempSkills, sk)
			}
		}

		// Set Homeland.Skills to new array
		hl.Homeland.Skills = tempSkills

		// Passions

		tempPassions := []runequest.Ability{}

		// Read passions
		for i := 1; i < 4; i++ {

			coreString := req.FormValue(fmt.Sprintf("Passion-%d-CoreString", i))

			if coreString != "" {

				p := runequest.Ability{
					Type:       "Passion",
					CoreString: coreString,
				}

				str := fmt.Sprintf("Passion-%d-Base", i)
				base, err := strconv.Atoi(req.FormValue(str))
				if err != nil {
					base = 0
				}
				p.Base = base

				userString := req.FormValue(fmt.Sprintf("Passion-%d-UserString", i))

				if userString != "" {
					p.UserChoice = true
					p.UserString = userString
				}

				tempPassions = append(tempPassions, p)
			}
		}

		// Reset passions to new values
		hl.Homeland.PassionList = tempPassions

		tempChoices := []runequest.SkillChoice{}

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

				str := fmt.Sprintf("Skill-%d-1-Value", i)
				v, err := strconv.Atoi(req.FormValue(str))
				if err != nil {
					v = 0
				}
				s1.HomelandValue = v

				if s1.UserChoice {
					userString := fmt.Sprintf("Skill-%d-1-UserString", i)
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

				str = fmt.Sprintf("Skill-%d-2-Value", i)
				v, err = strconv.Atoi(req.FormValue(str))
				if err != nil {
					v = 0
				}
				s2.HomelandValue = v

				if s2.UserChoice {
					userString := fmt.Sprintf("Skill-%d-2-UserString", i)
					s2.UserString = userString
				}

				// Form SkillChoice
				sc.Skills = []runequest.Skill{
					s1,
					s2,
				}
				// Append skillchoice
				tempChoices = append(tempChoices, sc)
			}
		}

		hl.Homeland.SkillChoices = tempChoices

		// Read New Skills

		for i := 1; i < 4; i++ {

			coreString := req.FormValue(fmt.Sprintf("Skill-%d-CoreString", i))

			if coreString != "" {

				sk := runequest.Skill{}

				sk.CoreString = coreString
				sk.Category = req.FormValue(fmt.Sprintf("Skill-%d-Category", i))

				userString := req.FormValue(fmt.Sprintf("Skill-%d-UserString", i))

				if userString != "" {
					sk.UserChoice = true
					sk.UserString = userString
				}

				str := fmt.Sprintf("Skill-%d-Base", i)
				base, err := strconv.Atoi(req.FormValue(str))
				if err != nil {
					base = 0
				}
				sk.Base = base

				str = fmt.Sprintf("Skill-%d-Value", i)
				v, err := strconv.Atoi(req.FormValue(str))
				if err != nil {
					v = 0
				}
				sk.HomelandValue = v

				hl.Homeland.Skills = append(hl.Homeland.Skills, sk)
			}
		}

		// Insert Homeland into App archive
		err = database.UpdateHomelandModel(db, hl)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("Saved")
		}

		fmt.Println(hl)

		url := fmt.Sprintf("/view_homeland/%d", hl.ID)

		http.Redirect(w, req, url, http.StatusSeeOther)
	}
}

// DeleteHomelandHandler renders a character in a Web page
func DeleteHomelandHandler(w http.ResponseWriter, req *http.Request) {

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

	if username == hl.Author.UserName || isAdmin == "true" {
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
		Render(w, "templates/delete_homeland.html", wc)

	}

	if req.Method == "POST" {

		err := database.DeleteHomelandModel(db, hl.ID)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("Deleted Homeland")
		}

		url := fmt.Sprint("/homeland_index/")

		http.Redirect(w, req, url, http.StatusSeeOther)
	}
}
