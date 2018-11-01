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

// OccupationIndexHandler renders the basic character roster page
func OccupationIndexHandler(w http.ResponseWriter, req *http.Request) {

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

	occupations, err := database.ListOccupationModels(db)
	if err != nil {
		panic(err)
	}

	for _, oc := range occupations {
		if oc.Image == nil {
			oc.Image = new(models.Image)
			oc.Image.Path = DefaultCharacterPortrait
		}
	}

	wc := WebChar{
		SessionUser:      username,
		IsLoggedIn:       loggedIn,
		IsAdmin:          isAdmin,
		OccupationModels: occupations,
	}
	Render(w, "templates/occupation_index.html", wc)
}

// OccupationListHandler applies a Occupation template to a character
func OccupationListHandler(w http.ResponseWriter, req *http.Request) {

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

	occupations, err := database.ListOccupationModels(db)
	if err != nil {
		panic(err)
	}

	wc := WebChar{
		CharacterModel:   cm,
		IsAuthor:         IsAuthor,
		SessionUser:      username,
		IsLoggedIn:       loggedIn,
		IsAdmin:          isAdmin,
		OccupationModels: occupations,
	}

	if req.Method == "GET" {

		// Render page

		Render(w, "templates/add_occupation_from_list.html", wc)

	}

	if req.Method == "POST" {

		err := req.ParseForm()
		if err != nil {
			panic(err)
		}

		occupationName := req.FormValue("Name")

		if occupationName != "" {
			oc := occupations[runequest.ToSnakeCase(occupationName)].Occupation

			fmt.Println(oc)

			c.Occupation = oc

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

// OccupationHandler renders a character in a Web page
func OccupationHandler(w http.ResponseWriter, req *http.Request) {

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

	oc, err := database.PKLoadOccupationModel(db, int64(id))
	if err != nil {
		fmt.Println(err)
		fmt.Println("Unable to load OccupationModel")
	}

	fmt.Println(oc)

	IsAuthor := false

	if username == oc.Author.UserName {
		IsAuthor = true
	}

	if oc.Image == nil {
		oc.Image = new(models.Image)
		oc.Image.Path = DefaultCharacterPortrait
	}

	wc := WebChar{
		OccupationModel: oc,
		IsAuthor:        IsAuthor,
		IsLoggedIn:      loggedIn,
		SessionUser:     username,
		IsAdmin:         isAdmin,
		Skills:          runequest.Skills,
		CategoryOrder:   runequest.CategoryOrder,
	}

	// Render page
	Render(w, "templates/view_occupation.html", wc)

}

// AddOccupationHandler creates a user-generated occupation
func AddOccupationHandler(w http.ResponseWriter, req *http.Request) {

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

	c.Statistics["STR"].Base = 10
	c.Statistics["DEX"].Base = 10
	c.Statistics["INT"].Base = 10
	c.Statistics["POW"].Base = 10
	c.Statistics["CHA"].Base = 10
	c.Statistics["SIZ"].Base = 10

	wc := WebChar{
		CharacterModel:    &cm,
		IsAuthor:          true,
		SessionUser:       username,
		IsLoggedIn:        loggedIn,
		IsAdmin:           isAdmin,
		Counter:           []int{1, 2, 3},
		BigCounter:        []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
		Passions:          runequest.PassionTypes,
		WeaponCategories:  runequest.WeaponCategories,
		CategoryOrder:     runequest.CategoryOrder,
		StandardsOfLiving: runequest.Standards,
		Skills:            runequest.Skills,
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/add_occupation.html", wc)

	}

	if req.Method == "POST" { // POST

		err := req.ParseMultipartForm(MaxMemory)
		if err != nil {
			panic(err)
		}

		user := database.LoadUser(db, username)

		// Map default Occupation to Character.Occupations
		oc := models.OccupationModel{
			Author: user,
			Occupation: &runequest.Occupation{
				Name:             req.FormValue("Name"),
				Description:      req.FormValue("Description"),
				StandardOfLiving: req.FormValue("Standard"),
				//CultChoices:
			},
		}

		income, err := strconv.Atoi(req.FormValue("Income"))
		if err != nil {
			income = 0
		}
		oc.Occupation.Income = income

		ransom, err := strconv.Atoi(req.FormValue("Ransom"))
		if err != nil {
			ransom = 0
		}
		oc.Occupation.Ransom = ransom

		if req.FormValue("Official") != "" {
			oc.Official = true
		} else {
			oc.Official = false
		}

		// Upload image to s3
		file, h, err := req.FormFile("image")
		switch err {
		case nil:
			// Process image
			defer file.Close()
			// example path media/Major/TestImage/Jason_White.jpg
			path := fmt.Sprintf("/media/%s/%s/%s",
				oc.Author.UserName,
				runequest.ToSnakeCase(oc.Occupation.Name),
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

			oc.Image = new(models.Image)
			oc.Image.Path = path

			fmt.Println(path)

		case http.ErrMissingFile:
			log.Println("no file")

		default:
			log.Panic(err)
			fmt.Println("Error getting file ", err)
		}

		var equipment = []string{}

		for i := 1; i < 16; i++ {
			str := req.FormValue(fmt.Sprintf("Equipment-%d", i))
			if str != "" {
				equipment = append(equipment, str)
			}
		}

		oc.Occupation.Equipment = equipment

		// Read Base Skills

		skillArray := []runequest.Skill{}

		for i := 1; i < 16; i++ {

			sk := req.FormValue(fmt.Sprintf("Skill-%d-CoreString", i))

			if sk != "" {

				skbaseSkill := runequest.Skills[sk]
				fmt.Println(skbaseSkill)

				// Skill
				s1 := runequest.Skill{
					CoreString: skbaseSkill.CoreString,
					UserChoice: skbaseSkill.UserChoice,
					Category:   skbaseSkill.Category,
					Base:       skbaseSkill.Base,
				}

				str := fmt.Sprintf("Skill-%d-Value", i)
				v, err := strconv.Atoi(req.FormValue(str))
				if err != nil {
					v = 0
				}
				s1.OccupationValue = v

				if s1.UserChoice {
					userString := fmt.Sprintf("Skill-%d-UserString", i)
					s1.UserString = req.FormValue(userString)
				}
				s1.GenerateName()
				skillArray = append(skillArray, s1)
			}
		}

		oc.Occupation.Skills = skillArray

		// Read Weapons
		for i := 1; i < 4; i++ {

			desc := req.FormValue(fmt.Sprintf("Weapon-%d-Description", i))

			if desc != "" {
				str := fmt.Sprintf("Weapon-%d-Value", i)
				v, err := strconv.Atoi(req.FormValue(str))
				if err != nil {
					v = 0
				}
				ws := runequest.WeaponSelection{
					Description: desc,
					Value:       v,
				}
				oc.Occupation.Weapons = append(oc.Occupation.Weapons, ws)
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
					Base:       s1baseSkill.Base,
				}

				str := fmt.Sprintf("Skill-%d-1-Value", i)
				v, err := strconv.Atoi(req.FormValue(str))
				if err != nil {
					v = 0
				}
				s1.OccupationValue = v

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
					Base:       s2baseSkill.Base,
				}

				str = fmt.Sprintf("Skill-%d-2-Value", i)
				v, err = strconv.Atoi(req.FormValue(str))
				if err != nil {
					v = 0
				}
				s2.OccupationValue = v

				if s2.UserChoice {
					userString := fmt.Sprintf("Skill-%d-2-UserString", i)
					s2.UserString = req.FormValue(userString)
				}

				s1.GenerateName()
				s2.GenerateName()

				// Form SkillChoice
				sc.Skills = []runequest.Skill{
					s1,
					s2,
				}
				// Append skillchoice
				oc.Occupation.SkillChoices = append(oc.Occupation.SkillChoices, sc)
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

				str = fmt.Sprintf("Passion-%d-Value", i)
				v, err := strconv.Atoi(req.FormValue(str))
				if err != nil {
					v = 0
				}
				p.OccupationValue = v

				userString := req.FormValue(fmt.Sprintf("Passion-%d-UserString", i))

				if userString != "" {
					p.UserChoice = true
					p.UserString = userString
				}

				oc.Occupation.PassionList = append(oc.Occupation.PassionList, p)
			}
		}

		// Read New Skills
		for i := 1; i < 4; i++ {

			coreString := req.FormValue(fmt.Sprintf("NewSkill-%d-CoreString", i))

			if coreString != "" {

				sk := runequest.Skill{}

				sk.CoreString = coreString
				sk.Category = req.FormValue(fmt.Sprintf("NewSkill-%d-Category", i))

				userString := req.FormValue(fmt.Sprintf("NewSkill-%d-UserString", i))

				if userString != "" {
					sk.UserChoice = true
					sk.UserString = userString
				}

				str := fmt.Sprintf("NewSkill-%d-Base", i)
				base, err := strconv.Atoi(req.FormValue(str))
				if err != nil {
					base = 0
				}
				sk.Base = base

				str = fmt.Sprintf("NewSkill-%d-Value", i)
				v, err := strconv.Atoi(req.FormValue(str))
				if err != nil {
					v = 0
				}
				sk.OccupationValue = v
				sk.GenerateName()

				oc.Occupation.Skills = append(oc.Occupation.Skills, sk)
			}
		}

		// Add other OccupationModel fields

		author := database.LoadUser(db, username)

		oc.Author = author

		err = database.SaveOccupationModel(db, &oc)
		if err != nil {
			log.Panic(err)
		} else {
			fmt.Println("Saved")
		}

		url := fmt.Sprintf("/view_occupation/%d", oc.ID)

		http.Redirect(w, req, url, http.StatusFound)
	}
}

// ModifyOccupationHandler renders an editable Occupation in a Web page
func ModifyOccupationHandler(w http.ResponseWriter, req *http.Request) {

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

	oc, err := database.PKLoadOccupationModel(db, int64(id))
	if err != nil {
		fmt.Println(err)
	}

	if oc.Author == nil {
		oc.Author = &models.User{
			UserName: "",
		}
	}

	// Validate that User == Author
	IsAuthor := false

	if username == oc.Author.UserName {
		IsAuthor = true
	} else {
		http.Redirect(w, req, "/", 302)
	}

	// Create empty equipment slots if < 16
	if len(oc.Occupation.Equipment) < 16 {
		for i := len(oc.Occupation.Equipment); i < 16; i++ {
			oc.Occupation.Equipment = append(oc.Occupation.Equipment, "")
		}
	}

	// Add extra empty skills if < 20
	if len(oc.Occupation.Skills) < 20 {
		for i := len(oc.Occupation.Skills); i < 20; i++ {
			tempSkill := runequest.Skill{}
			oc.Occupation.Skills = append(oc.Occupation.Skills, tempSkill)
		}
	}

	// Add extra empty skills if < 20
	if len(oc.Occupation.PassionList) < 4 {
		for i := len(oc.Occupation.PassionList); i < 4; i++ {
			p := runequest.Ability{
				Base:            60,
				OccupationValue: 10,
			}
			oc.Occupation.PassionList = append(oc.Occupation.PassionList, p)
		}
	}

	// Add extra empty skillchoices if < 3
	if len(oc.Occupation.SkillChoices) < 4 {
		for i := len(oc.Occupation.SkillChoices); i < 4; i++ {
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
			oc.Occupation.SkillChoices = append(oc.Occupation.SkillChoices, tempSkillChoice)
		}
	}

	// Add empty Weapon options if < 5
	if len(oc.Occupation.Weapons) < 5 {
		for i := len(oc.Occupation.Weapons); i < 5; i++ {
			tempWeapon := runequest.WeaponSelection{
				Description: "",
				Value:       0,
			}
			oc.Occupation.Weapons = append(oc.Occupation.Weapons, tempWeapon)
		}
	}

	wc := WebChar{
		OccupationModel:   oc,
		IsAuthor:          IsAuthor,
		SessionUser:       username,
		IsLoggedIn:        loggedIn,
		IsAdmin:           isAdmin,
		Counter:           []int{1, 2, 3},
		BigCounter:        []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
		Passions:          runequest.PassionTypes,
		WeaponCategories:  runequest.WeaponCategories,
		CategoryOrder:     runequest.CategoryOrder,
		StandardsOfLiving: runequest.Standards,
		Skills:            runequest.Skills,
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/modify_occupation.html", wc)

	}

	if req.Method == "POST" { // POST

		err := req.ParseMultipartForm(MaxMemory)
		if err != nil {
			panic(err)
		}

		// Update Occupation here
		ocName := req.FormValue("Name")
		description := req.FormValue("Description")

		oc.Occupation.Name = ocName
		oc.Occupation.Description = description

		income, err := strconv.Atoi(req.FormValue("Income"))
		if err != nil {
			income = 0
		}
		oc.Occupation.Income = income

		ransom, err := strconv.Atoi(req.FormValue("Ransom"))
		if err != nil {
			ransom = 0
		}
		oc.Occupation.Ransom = ransom

		if req.FormValue("Official") != "" {
			oc.Official = true
		} else {
			oc.Official = false
		}

		// Insert Occupation into App archive if user authorizes
		if req.FormValue("Archive") != "" {
			oc.Open = true
		} else {
			oc.Open = false
		}

		// Open or Official
		if req.FormValue("Official") != "" {
			oc.Official = true
		} else {
			oc.Official = false
		}

		// Upload image to s3
		file, h, err := req.FormFile("image")
		switch err {
		case nil:
			// Process image
			defer file.Close()
			// example path media/Major/TestImage/Jason_White.jpg
			path := fmt.Sprintf("/media/%s/%s/%s",
				oc.Author.UserName,
				runequest.ToSnakeCase(oc.Occupation.Name),
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

			oc.Image = new(models.Image)
			oc.Image.Path = path

			fmt.Println(path)

		case http.ErrMissingFile:
			log.Println("no file")

		default:
			log.Panic(err)
			fmt.Println("Error getting file ", err)
		}

		// Read Equipment
		var equipment = []string{}

		for i := 1; i < 16; i++ {
			str := req.FormValue(fmt.Sprintf("Equipment-%d", i))
			if str != "" {
				equipment = append(equipment, str)
			}
		}

		oc.Occupation.Equipment = equipment

		// Read Skills
		tempSkills := []runequest.Skill{}

		// Read Base Skills from Form
		for i := 1; i < 20; i++ {

			core := req.FormValue(fmt.Sprintf("Skill-%d-CoreString", i))

			str := fmt.Sprintf("Skill-%d-Value", i)
			v, err := strconv.Atoi(req.FormValue(str))
			if err != nil {
				v = 0
			}

			if core != "" && v > 0 {

				sk := runequest.Skill{
					CoreString:      core,
					OccupationValue: v,
				}

				userString := req.FormValue(fmt.Sprintf("Skill-%d-UserString", i))
				if userString != "" {
					sk.UserString = userString
				}
				sk.GenerateName()
				tempSkills = append(tempSkills, sk)
			}
		}

		// Set Occupation.Skills to new array
		oc.Occupation.Skills = tempSkills

		// Read Weapons
		weapons := []runequest.WeaponSelection{}

		for i := 1; i < 5; i++ {

			desc := req.FormValue(fmt.Sprintf("Weapon-%d-Description", i))

			if desc != "" {
				str := fmt.Sprintf("Weapon-%d-Value", i)
				v, err := strconv.Atoi(req.FormValue(str))
				if err != nil {
					v = 0
				}
				ws := runequest.WeaponSelection{
					Description: desc,
					Value:       v,
				}
				weapons = append(weapons, ws)
			}
		}

		oc.Occupation.Weapons = weapons

		// Read passions
		tempPassions := []runequest.Ability{}

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

				str = fmt.Sprintf("Passion-%d-Value", i)
				v, err := strconv.Atoi(req.FormValue(str))
				if err != nil {
					v = 0
				}
				p.OccupationValue = v

				userString := req.FormValue(fmt.Sprintf("Passion-%d-UserString", i))

				if userString != "" {
					p.UserChoice = true
					p.UserString = userString
				}

				tempPassions = append(tempPassions, p)
			}
		}

		// Reset passions to new values
		oc.Occupation.PassionList = tempPassions

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
					Base:       s1baseSkill.Base,
				}

				str := fmt.Sprintf("Skill-%d-1-Value", i)
				v, err := strconv.Atoi(req.FormValue(str))
				if err != nil {
					v = 0
				}
				s1.OccupationValue = v

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
					Base:       s2baseSkill.Base,
				}

				str = fmt.Sprintf("Skill-%d-2-Value", i)
				v, err = strconv.Atoi(req.FormValue(str))
				if err != nil {
					v = 0
				}
				s2.OccupationValue = v

				if s2.UserChoice {
					userString := fmt.Sprintf("Skill-%d-2-UserString", i)
					s2.UserString = req.FormValue(userString)
				}

				s1.GenerateName()
				s2.GenerateName()

				// Form SkillChoice
				sc.Skills = []runequest.Skill{
					s1,
					s2,
				}
				// Append skillchoice
				tempChoices = append(tempChoices, sc)
			}
		}

		oc.Occupation.SkillChoices = tempChoices

		// Read New Skills

		for i := 1; i < 4; i++ {

			coreString := req.FormValue(fmt.Sprintf("NewSkill-%d-CoreString", i))

			if coreString != "" {

				sk := runequest.Skill{}

				sk.CoreString = coreString
				sk.Category = req.FormValue(fmt.Sprintf("NewSkill-%d-Category", i))

				userString := req.FormValue(fmt.Sprintf("NewSkill-%d-UserString", i))

				if userString != "" {
					sk.UserChoice = true
					sk.UserString = userString
				}

				str := fmt.Sprintf("NewSkill-%d-Base", i)
				base, err := strconv.Atoi(req.FormValue(str))
				if err != nil {
					base = 0
				}
				sk.Base = base

				str = fmt.Sprintf("NewSkill-%d-Value", i)
				v, err := strconv.Atoi(req.FormValue(str))
				if err != nil {
					v = 0
				}
				sk.OccupationValue = v

				sk.GenerateName()

				oc.Occupation.Skills = append(oc.Occupation.Skills, sk)
			}
		}

		// Insert Occupation into App archive
		err = database.UpdateOccupationModel(db, oc)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("Saved")
		}

		fmt.Println(oc)

		url := fmt.Sprintf("/view_occupation/%d", oc.ID)

		http.Redirect(w, req, url, http.StatusSeeOther)
	}
}

// DeleteOccupationHandler renders a character in a Web page
func DeleteOccupationHandler(w http.ResponseWriter, req *http.Request) {

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

	oc, err := database.PKLoadOccupationModel(db, int64(id))
	if err != nil {
		fmt.Println(err)
	}

	// Validate that User == Author
	IsAuthor := false

	if username == oc.Author.UserName || isAdmin == "true" {
		IsAuthor = true
	} else {
		http.Redirect(w, req, "/", 302)
	}

	wc := WebChar{
		OccupationModel: oc,
		IsAuthor:        IsAuthor,
		SessionUser:     username,
		IsLoggedIn:      loggedIn,
		IsAdmin:         isAdmin,
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/delete_occupation.html", wc)

	}

	if req.Method == "POST" {

		err := database.DeleteOccupationModel(db, oc.ID)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("Deleted Occupation")
		}

		url := fmt.Sprint("/occupation_index/")

		http.Redirect(w, req, url, http.StatusSeeOther)
	}
}
