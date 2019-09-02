package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gosimple/slug"

	"github.com/toferc/runequest"

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
	slug := vars["slug"]

	if len(slug) == 0 {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}

	oc, err := database.LoadOccupationModel(db, slug)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Unable to load OccupationModel")
		http.Redirect(w, req, "/notfound", http.StatusSeeOther)
		return
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
		Counter:           numToArray(4),
		BigCounter:        numToArray(15),
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

		user, err := database.LoadUser(db, username)
		if err != nil {
			fmt.Println(err)
			http.Redirect(w, req, "/", 302)
		}

		// Map default Occupation to Character.Occupations
		oc := models.OccupationModel{
			Author: user,
			Occupation: &runequest.Occupation{
				Name:             req.FormValue("Name"),
				Description:      req.FormValue("Description"),
				StandardOfLiving: req.FormValue("Standard"),
				Notes:            req.FormValue("Notes"),
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
			if h.Filename != "" {
				// Process image
				defer file.Close()

				err = ProcessOccupationImage(h, file, &oc)
				if err != nil {
					log.Printf("Error processing image: %v", err)
				}

			} else {
				fmt.Println("No file provided.")
			}

		case http.ErrMissingFile:
			log.Println("no file")
			oc.Image = new(models.Image)
			oc.Image.Path = DefaultCharacterPortrait

		default:
			log.Panic(err)
			fmt.Println("Error getting file ", err)
			oc.Image = new(models.Image)
			oc.Image.Path = DefaultCharacterPortrait
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

		skillArray := []*runequest.Skill{}

		for i := 1; i < 16; i++ {

			sk := req.FormValue(fmt.Sprintf("Skill-%d-CoreString", i))

			if sk != "" {

				skbaseSkill := runequest.Skills[sk]
				fmt.Println(skbaseSkill)

				// Skill
				s1 := &runequest.Skill{
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
		for i := 1; i < 5; i++ {

			coreString := req.FormValue(fmt.Sprintf("Passion-%d-CoreString", i))

			if coreString != "" {

				str := fmt.Sprintf("Passion-%d-Base", i)
				base, err := strconv.Atoi(req.FormValue(str))
				if err != nil {
					base = 0
				}

				str = fmt.Sprintf("Passion-%d-Value", i)
				v, err := strconv.Atoi(req.FormValue(str))
				if err != nil {
					v = 0
				}

				userString := req.FormValue(fmt.Sprintf("Passion-%d-UserString", i))

				p := runequest.Ability{
					Type:            "Passion",
					CoreString:      coreString,
					Base:            base,
					OccupationValue: v,
					UserString:      userString,
				}
				oc.Occupation.PassionList = append(oc.Occupation.PassionList, p)
			}
		}

		// Read New Skills
		for i := 1; i < 5; i++ {

			coreString := req.FormValue(fmt.Sprintf("NewSkill-%d-CoreString", i))

			if coreString != "" {

				sk := &runequest.Skill{}

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

		author, err := database.LoadUser(db, username)
		if err != nil {
			fmt.Println(err)
			http.Redirect(w, req, "/", 302)
		}

		oc.Author = author

		oc.Slug = slug.Make(oc.Occupation.Name)

		err = database.SaveOccupationModel(db, &oc)
		if err != nil {
			log.Panic(err)
		} else {
			fmt.Println("Saved")
		}

		url := fmt.Sprintf("/view_occupation/%s", oc.Slug)

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

	if username == oc.Author.UserName || isAdmin == "true" {
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
			tempSkill := &runequest.Skill{}
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

	// Add extra empty skillchoices if < 5
	if len(oc.Occupation.SkillChoices) < 5 {
		for i := len(oc.Occupation.SkillChoices); i < 5; i++ {
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
		Counter:           numToArray(4),
		BigCounter:        numToArray(16),
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
		oc.Occupation.Notes = req.FormValue("Notes")

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

		oc.Occupation.StandardOfLiving = req.FormValue("Standard")

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
			if h.Filename != "" {
				// Process image
				defer file.Close()

				err = ProcessOccupationImage(h, file, oc)
				if err != nil {
					log.Printf("Error processing image: %v", err)
				}

			} else {
				fmt.Println("No file provided.")
			}

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
		tempSkills := []*runequest.Skill{}

		// Read Base Skills from Form
		for i := 1; i < 20; i++ {

			sk := req.FormValue(fmt.Sprintf("Skill-%d-CoreString", i))

			str := fmt.Sprintf("Skill-%d-Value", i)
			v, err := strconv.Atoi(req.FormValue(str))
			if err != nil {
				v = 0
			}

			userString := req.FormValue(fmt.Sprintf("Skill-%d-UserString", i))

			if sk != "" && v > 0 {

				s1 := &runequest.Skill{}

				skbaseSkill, ok := runequest.Skills[sk]
				if !ok {
					// Skill is new
					for _, ns := range oc.Occupation.Skills {
						// Search for CoreString in Occupation Skills
						if sk == ns.CoreString {
							s1.CoreString = ns.CoreString
							s1.UserChoice = ns.UserChoice
							s1.Base = ns.Base
							s1.Category = ns.Category
						}
					}
				} else {
					// Skill exists in base list
					s1.CoreString = skbaseSkill.CoreString
					s1.UserChoice = skbaseSkill.UserChoice
					s1.Base = skbaseSkill.Base
					s1.Category = skbaseSkill.Category
				}

				if s1.Category == "" {
					s1.Category = req.FormValue(fmt.Sprintf("Skill-%d-Category", i))
				}

				s1.OccupationValue = v

				if userString != "" {
					s1.UserString = userString
				}
				s1.GenerateName()
				tempSkills = append(tempSkills, s1)
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

		for i := 1; i < 5; i++ {

			coreString := req.FormValue(fmt.Sprintf("Passion-%d-CoreString", i))

			if coreString != "" {

				str := fmt.Sprintf("Passion-%d-Base", i)
				base, err := strconv.Atoi(req.FormValue(str))
				if err != nil {
					base = 0
				}

				str = fmt.Sprintf("Passion-%d-Value", i)
				v, err := strconv.Atoi(req.FormValue(str))
				if err != nil {
					v = 0
				}

				userString := req.FormValue(fmt.Sprintf("Passion-%d-UserString", i))

				p := runequest.Ability{
					Type:            "Passion",
					CoreString:      coreString,
					Base:            base,
					OccupationValue: v,
					UserString:      userString,
				}
				tempPassions = append(tempPassions, p)
			}
		}

		// Reset passions to new values
		oc.Occupation.PassionList = tempPassions

		tempChoices := []runequest.SkillChoice{}

		// Read SkillChoices
		for i := 1; i < 5; i++ {

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

		for i := 1; i < 5; i++ {

			coreString := req.FormValue(fmt.Sprintf("NewSkill-%d-CoreString", i))

			if coreString != "" {

				sk := &runequest.Skill{}

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

		oc.Slug = slug.Make(oc.Occupation.Name)

		// Insert Occupation into App archive
		err = database.UpdateOccupationModel(db, oc)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("Saved")
		}

		fmt.Println(oc)

		url := fmt.Sprintf("/view_occupation/%s", oc.Slug)

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

	if oc.Image == nil {
		oc.Image = new(models.Image)
		oc.Image.Path = DefaultCharacterPortrait
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
