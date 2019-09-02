package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gosimple/slug"

	"github.com/toferc/runequest"

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
	slug := vars["slug"]

	if len(slug) == 0 {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}

	hl, err := database.LoadHomelandModel(db, slug)
	if err != nil {
		fmt.Println("Unable to load HomelandModel")
		http.Redirect(w, req, "/notfound", http.StatusSeeOther)
		return
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
		StringArray:   runequest.StatMap,
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

	c.Statistics["STR"].Base = 10
	c.Statistics["DEX"].Base = 10
	c.Statistics["INT"].Base = 10
	c.Statistics["POW"].Base = 10
	c.Statistics["CHA"].Base = 10
	c.Statistics["SIZ"].Base = 10

	user, err := database.LoadUser(db, username)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, req, "/", 302)
	}

	// Map default Homeland to Character.Homelands
	hl := models.HomelandModel{
		Author: user,
		Homeland: &runequest.Homeland{
			StatisticFrames: runequest.HomeLandStats,
		},
	}

	// Add empty Movement's to homeland
	if len(hl.Homeland.Movement) == 0 {
		hl.Homeland.Movement = []runequest.Movement{
			runequest.Movement{
				Name:  "Ground",
				Value: 8,
			},
			runequest.Movement{
				Name:  "",
				Value: 0,
			},
			runequest.Movement{
				Name:  "",
				Value: 0,
			},
		}
	}

	wc := WebChar{
		CharacterModel:   &cm,
		HomelandModel:    &hl,
		IsAuthor:         true,
		SessionUser:      username,
		IsLoggedIn:       loggedIn,
		IsAdmin:          isAdmin,
		Counter:          numToArray(3),
		BigCounter:       numToArray(20),
		Passions:         runequest.PassionTypes,
		CategoryOrder:    runequest.CategoryOrder,
		Skills:           runequest.Skills,
		PowerRunes:       runequest.PowerRuneOrder,
		ElementalRunes:   runequest.ElementalRuneOrder,
		HitLocationForms: runequest.LocationForms,
		StringArray:      runequest.StatMap,
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

		hl.Homeland.Name = req.FormValue("Name")
		hl.Homeland.Description = req.FormValue("Description")
		hl.Homeland.Notes = req.FormValue("Notes")

		tempMovement := []runequest.Movement{}

		// Add Movement
		for i := 1; i < 4; i++ {
			moveName := req.FormValue(fmt.Sprintf("Move-Name-%d", i))

			if moveName != "" {
				mv, err := strconv.Atoi(req.FormValue(fmt.Sprintf("Move-Value-%d", i)))
				if err != nil {
					mv = 8
				}

				tempMovement = append(tempMovement,
					runequest.Movement{
						Name:  moveName,
						Value: mv,
					})
			}
		}
		hl.Homeland.Movement = tempMovement

		hl.Homeland.LocationForm = req.FormValue("Hit-Location-Form")

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
			if h.Filename != "" {
				// Process image
				defer file.Close()

				err = ProcessHomelandImage(h, file, &hl)
				if err != nil {
					log.Printf("Error processing image: %v", err)
				}

			} else {
				fmt.Println("No file provided.")
			}

		case http.ErrMissingFile:
			log.Println("no file")
			hl.Image = new(models.Image)
			hl.Image.Path = DefaultCharacterPortrait

		default:
			log.Panic(err)
			fmt.Println("Error getting file ", err)
			hl.Image = new(models.Image)
			hl.Image.Path = DefaultCharacterPortrait
		}

		// Read Rune

		hl.Homeland.RuneBonus = req.FormValue("Rune")

		// Read Stat Framework

		for k, v := range hl.Homeland.StatisticFrames {

			dString := fmt.Sprintf("Stat-%s-Dice", k)
			dice, err := strconv.Atoi(req.FormValue(dString))
			if err != nil {
				dice = 3
			}

			mString := fmt.Sprintf("Stat-%s-Modifier", k)
			mod, err := strconv.Atoi(req.FormValue(mString))
			if err != nil {
				mod = 0
			}

			v.Dice = dice
			v.Modifier = mod
			v.Max = (dice * 6) + mod + dice

			if mod > 0 {
				v.Max++
			}

			fmt.Println(dice, mod)
		}

		// Read Base Skills

		skillArray := []*runequest.Skill{}

		for i := 1; i < 20; i++ {

			sk := req.FormValue(fmt.Sprintf("Skill-%d-CoreString", i))

			if sk != "" {

				skbaseSkill := runequest.Skills[sk]
				fmt.Println(skbaseSkill)

				// Skill
				s1 := &runequest.Skill{
					CoreString: skbaseSkill.CoreString,
					UserChoice: skbaseSkill.UserChoice,
					Base:       skbaseSkill.Base,
					Category:   skbaseSkill.Category,
				}

				str := fmt.Sprintf("Skill-%d-Value", i)
				v, err := strconv.Atoi(req.FormValue(str))
				if err != nil {
					v = 0
				}
				s1.HomelandValue = v

				str = fmt.Sprintf("Skill-%d-Base", i)
				b, err := strconv.Atoi(req.FormValue(str))
				if err != nil {
					b = 0
				}
				if b > s1.Base {
					s1.Base = b
				}

				if s1.UserChoice {
					userString := fmt.Sprintf("Skill-%d-UserString", i)
					s1.UserString = strings.TrimSpace(req.FormValue(userString))
				}
				skillArray = append(skillArray, s1)
			}
		}

		hl.Homeland.Skills = skillArray

		// Read passions
		for i := 1; i < 4; i++ {

			coreString := req.FormValue(fmt.Sprintf("Passion-%d-CoreString", i))

			if coreString != "" {

				str := fmt.Sprintf("Passion-%d-Base", i)
				base, err := strconv.Atoi(req.FormValue(str))
				if err != nil {
					fmt.Println(err)
					base = 60
				}

				u := req.FormValue(fmt.Sprintf("Passion-%d-UserString", i))
				userString := strings.TrimSpace(u)

				p := runequest.Ability{
					Type:          "Passion",
					CoreString:    coreString,
					Base:          base,
					HomelandValue: 10,
					UserString:    userString,
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
					Base:       s1baseSkill.Base,
				}

				str := fmt.Sprintf("Skill-%d-1-Value", i)
				v, err := strconv.Atoi(req.FormValue(str))
				if err != nil {
					v = 0
				}
				s1.HomelandValue = v

				if s1.UserChoice {
					userString := fmt.Sprintf("Skill-%d-1-UserString", i)
					s1.UserString = strings.TrimSpace(req.FormValue(userString))
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
				s2.HomelandValue = v

				if s2.UserChoice {
					userString := fmt.Sprintf("Skill-%d-2-UserString", i)
					s2.UserString = strings.TrimSpace(req.FormValue(userString))
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

		// Read Custom Skills

		for i := 1; i < 4; i++ {

			c := req.FormValue(fmt.Sprintf("NewSkill-%d-CoreString", i))
			coreString := strings.TrimSpace(c)

			if coreString != "" {

				sk := &runequest.Skill{}

				sk.CoreString = coreString
				sk.Category = req.FormValue(fmt.Sprintf("NewSkill-%d-Category", i))

				u := req.FormValue(fmt.Sprintf("NewSkill-%d-UserString", i))
				userString := strings.TrimSpace(u)

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
				sk.HomelandValue = v

				hl.Homeland.Skills = append(hl.Homeland.Skills, sk)
			}
		}

		// Add other HomelandModel fields

		author, err := database.LoadUser(db, username)
		if err != nil {
			fmt.Println(err)
			http.Redirect(w, req, "/", 302)
		}

		hl.Author = author

		hl.Slug = slug.Make(hl.Homeland.Name)

		err = database.SaveHomelandModel(db, &hl)
		if err != nil {
			log.Panic(err)
		} else {
			fmt.Println("Saved")
		}

		url := fmt.Sprintf("/view_homeland/%s", hl.Slug)

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

	if username == hl.Author.UserName || isAdmin == "true" {
		IsAuthor = true
	} else {
		http.Redirect(w, req, "/", 302)
	}

	// Add Homelandstats if not already there
	if hl.Homeland.StatisticFrames == nil {
		hl.Homeland.StatisticFrames = runequest.HomeLandStats
	}

	// Create empty passions for adding to the Homeland

	if len(hl.Homeland.PassionList) < 4 {
		for i := len(hl.Homeland.PassionList); i < 4; i++ {
			passion := runequest.Ability{}
			hl.Homeland.PassionList = append(hl.Homeland.PassionList, passion)
		}
	}

	// Create empty skills for adding to the Homeland

	// Add extra empty skills if < 20
	if len(hl.Homeland.Skills) < 20 {
		for i := len(hl.Homeland.Skills); i < 20; i++ {
			tempSkill := &runequest.Skill{}
			hl.Homeland.Skills = append(hl.Homeland.Skills, tempSkill)
		}
	}

	// Add extra empty skillchoices if < 3
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

	// Add empty Movement's to homeland
	if len(hl.Homeland.Movement) == 0 {
		hl.Homeland.Movement = []runequest.Movement{
			runequest.Movement{
				Name:  "Ground",
				Value: 8,
			},
		}
	}

	mvLen := len(hl.Homeland.Movement)

	if mvLen < 3 {

		// Add Movement
		for i := mvLen; i < mvLen+2; i++ {

			hl.Homeland.Movement = append(hl.Homeland.Movement,
				runequest.Movement{
					Name:  "",
					Value: 0,
				})
		}
	}

	if hl.Image == nil {
		hl.Image = new(models.Image)
		hl.Image.Path = DefaultCharacterPortrait
	}

	wc := WebChar{
		HomelandModel:    hl,
		IsAuthor:         IsAuthor,
		SessionUser:      username,
		IsLoggedIn:       loggedIn,
		IsAdmin:          isAdmin,
		Counter:          []int{1, 2, 3},
		CategoryOrder:    runequest.CategoryOrder,
		Skills:           runequest.Skills,
		Passions:         runequest.PassionTypes,
		PowerRunes:       runequest.PowerRuneOrder,
		ElementalRunes:   runequest.ElementalRuneOrder,
		HitLocationForms: runequest.LocationForms,
		StringArray:      runequest.StatMap,
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

		// Update Homeland here
		hl.Homeland.Name = req.FormValue("Name")
		hl.Homeland.Description = req.FormValue("Description")
		hl.Homeland.Notes = req.FormValue("Notes")

		tempMovement := []runequest.Movement{}

		// Add Movement
		for i := 1; i < mvLen+2; i++ {
			m := req.FormValue(fmt.Sprintf("Move-Name-%d", i))
			moveName := strings.TrimSpace(m)

			if moveName != "" {
				mv, err := strconv.Atoi(req.FormValue(fmt.Sprintf("Move-Value-%d", i)))
				if err != nil {
					mv = 8
				}

				tempMovement = append(tempMovement,
					runequest.Movement{
						Name:  moveName,
						Value: mv,
					})
			}
		}
		hl.Homeland.Movement = tempMovement

		hl.Homeland.LocationForm = req.FormValue("Hit-Location-Form")

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

		// Read Stat Framework

		for k, v := range hl.Homeland.StatisticFrames {

			dString := fmt.Sprintf("Stat-%s-Dice", k)
			dice, err := strconv.Atoi(req.FormValue(dString))
			if err != nil {
				dice = 3
			}

			mString := fmt.Sprintf("Stat-%s-Modifier", k)
			mod, err := strconv.Atoi(req.FormValue(mString))
			if err != nil {
				mod = 0
			}

			v.Dice = dice
			v.Modifier = mod
			v.Max = (dice * 6) + mod + dice

			if mod > 0 {
				v.Max++
			}

			fmt.Println(dice, mod)
		}

		// Upload image to s3
		file, h, err := req.FormFile("image")
		switch err {
		case nil:
			if h.Filename != "" {
				// Process image
				defer file.Close()

				err = ProcessHomelandImage(h, file, hl)
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

		hl.Homeland.RuneBonus = req.FormValue("Rune")

		// Read Base Skills

		skillArray := []*runequest.Skill{}

		for i := 1; i < 20; i++ {

			sk := req.FormValue(fmt.Sprintf("Skill-%d-CoreString", i))
			u := req.FormValue(fmt.Sprintf("Skill-%d-UserString", i))
			userString := strings.TrimSpace(u)

			if sk != "" {

				s1 := &runequest.Skill{}

				skbaseSkill, ok := runequest.Skills[sk]
				if !ok {
					// Skill is new
					for _, ns := range hl.Homeland.Skills {
						// Search for CoreString in Homeland Skills
						if sk == ns.CoreString {
							s1.CoreString = ns.CoreString
							s1.UserChoice = ns.UserChoice
							s1.Base = ns.Base
							s1.Category = ns.Category
						}
					}
				} else {
					// Skill
					s1.CoreString = skbaseSkill.CoreString
					s1.UserChoice = skbaseSkill.UserChoice
					s1.Base = skbaseSkill.Base
					s1.Category = skbaseSkill.Category
				}

				fmt.Println(skbaseSkill)

				str := fmt.Sprintf("Skill-%d-Value", i)
				v, err := strconv.Atoi(req.FormValue(str))
				if err != nil {
					v = 0
				}
				s1.HomelandValue = v

				str = fmt.Sprintf("Skill-%d-Base", i)
				b, err := strconv.Atoi(req.FormValue(str))
				if err != nil {
					b = 0
				}

				s1.Base = b

				if s1.UserChoice && userString != "" {
					s1.UserString = userString
				}

				s1.GenerateName()

				skillArray = append(skillArray, s1)
			}
		}

		hl.Homeland.Skills = skillArray

		// Passions

		tempPassions := []runequest.Ability{}

		// Read passions
		for i := 1; i < 4; i++ {

			coreString := req.FormValue(fmt.Sprintf("Passion-%d-CoreString", i))

			if coreString != "" {

				str := fmt.Sprintf("Passion-%d-Base", i)
				base, err := strconv.Atoi(req.FormValue(str))
				if err != nil {
					base = 0
				}

				u := req.FormValue(fmt.Sprintf("Passion-%d-UserString", i))
				userString := strings.TrimSpace(u)

				p := runequest.Ability{
					Type:          "Passion",
					CoreString:    coreString,
					Base:          base,
					HomelandValue: 10,
					UserString:    userString,
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
					Base:       s1baseSkill.Base,
				}

				str := fmt.Sprintf("Skill-%d-1-Value", i)
				v, err := strconv.Atoi(req.FormValue(str))
				if err != nil {
					v = 0
				}
				s1.HomelandValue = v

				if s1.UserChoice {
					userString := fmt.Sprintf("Skill-%d-1-UserString", i)
					s1.UserString = strings.TrimSpace(req.FormValue(userString))
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
				s2.HomelandValue = v

				if s2.UserChoice {
					userString := fmt.Sprintf("Skill-%d-2-UserString", i)
					s2.UserString = strings.TrimSpace(req.FormValue(userString))
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

		hl.Homeland.SkillChoices = tempChoices

		// Read New Skills

		for i := 1; i < 4; i++ {

			coreString := req.FormValue(fmt.Sprintf("NewSkill-%d-CoreString", i))

			if coreString != "" {

				sk := runequest.Skill{}

				sk.CoreString = coreString
				sk.Category = req.FormValue(fmt.Sprintf("NewSkill-%d-Category", i))

				u := req.FormValue(fmt.Sprintf("NewSkill-%d-UserString", i))
				userString := strings.TrimSpace(u)

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
				sk.HomelandValue = v

				hl.Homeland.Skills = append(hl.Homeland.Skills, &sk)
			}
		}

		hl.Slug = slug.Make(hl.Homeland.Name)

		// Insert Homeland into App archive
		err = database.UpdateHomelandModel(db, hl)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("Saved")
		}

		fmt.Println(hl)

		url := fmt.Sprintf("/view_homeland/%s", hl.Slug)

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

	if hl.Image == nil {
		hl.Image = new(models.Image)
		hl.Image.Path = DefaultCharacterPortrait
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
