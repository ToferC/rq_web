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

// CultIndexHandler renders the basic character roster page
func CultIndexHandler(w http.ResponseWriter, req *http.Request) {

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

	cults, err := database.ListCultModels(db)
	if err != nil {
		panic(err)
	}

	for _, cl := range cults {
		if cl.Image == nil {
			cl.Image = new(models.Image)
			cl.Image.Path = DefaultCharacterPortrait
		}
	}

	wc := WebChar{
		SessionUser: username,
		IsLoggedIn:  loggedIn,
		IsAdmin:     isAdmin,
		CultModels:  cults,
	}
	Render(w, "templates/cult_index.html", wc)
}

// CultListHandler applies a Cult template to a character
func CultListHandler(w http.ResponseWriter, req *http.Request) {

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

	cults, err := database.ListCultModels(db)
	if err != nil {
		panic(err)
	}

	wc := WebChar{
		CharacterModel: cm,
		IsAuthor:       IsAuthor,
		SessionUser:    username,
		IsLoggedIn:     loggedIn,
		IsAdmin:        isAdmin,
		CultModels:     cults,
	}

	if req.Method == "GET" {

		// Render page

		Render(w, "templates/add_cult_from_list.html", wc)

	}

	if req.Method == "POST" {

		err := req.ParseForm()
		if err != nil {
			panic(err)
		}

		cultName := req.FormValue("Name")

		if cultName != "" {
			cl := cults[runequest.ToSnakeCase(cultName)].Cult

			fmt.Println(cl)

			c.Cult = cl

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

// CultHandler renders a character in a Web page
func CultHandler(w http.ResponseWriter, req *http.Request) {

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

	cl, err := database.LoadCultModel(db, slug)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Unable to load CultModel")
	}

	fmt.Println(cl)

	IsAuthor := false

	if username == cl.Author.UserName {
		IsAuthor = true
	}

	if cl.Image == nil {
		cl.Image = new(models.Image)
		cl.Image.Path = DefaultCharacterPortrait
	}

	cults, err := database.ListCultModels(db)
	if err != nil {
		panic(err)
	}

	wc := WebChar{
		CultModel:     cl,
		IsAuthor:      IsAuthor,
		IsLoggedIn:    loggedIn,
		SessionUser:   username,
		IsAdmin:       isAdmin,
		Skills:        runequest.Skills,
		CategoryOrder: runequest.CategoryOrder,
		CultModels:    cults,
	}

	// Render page
	Render(w, "templates/view_cult.html", wc)

}

// AddCultHandler creates a user-generated cult
func AddCultHandler(w http.ResponseWriter, req *http.Request) {

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

	cults, err := database.ListCultModels(db)
	if err != nil {
		panic(err)
	}

	c.Statistics["STR"].Base = 10
	c.Statistics["DEX"].Base = 10
	c.Statistics["INT"].Base = 10
	c.Statistics["POW"].Base = 10
	c.Statistics["CHA"].Base = 10
	c.Statistics["SIZ"].Base = 10

	wc := WebChar{
		CharacterModel:   &cm,
		IsAuthor:         true,
		SessionUser:      username,
		IsLoggedIn:       loggedIn,
		IsAdmin:          isAdmin,
		Counter:          numToArray(7),
		BigCounter:       numToArray(15),
		Passions:         runequest.PassionTypes,
		WeaponCategories: runequest.WeaponCategories,
		CategoryOrder:    runequest.CategoryOrder,
		PowerRunes:       runequest.PowerRuneOrder,
		ElementalRunes:   runequest.ElementalRuneOrder,
		Skills:           runequest.Skills,
		CultModels:       cults,
		RuneSpells:       runequest.RuneSpells,
		SpiritMagic:      runequest.SpiritMagicSpells,
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/add_cult.html", wc)

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

		// Map default Cult to Character.Cults
		cl := models.CultModel{
			Author: user,
			Cult: &runequest.Cult{
				Name:        req.FormValue("Name"),
				Description: req.FormValue("Description"),
				Notes:       req.FormValue("Notes"),
			},
		}

		if req.FormValue("Official") != "" {
			cl.Official = true
		} else {
			cl.Official = false
		}

		if req.FormValue("SubCult") != "" {
			cl.Cult.SubCult = true
		} else {
			cl.Cult.SubCult = false
		}

		parentCultModel := &models.CultModel{}

		if cl.Cult.SubCult {
			// Set ParentCult
			cStr := req.FormValue("ParentCult")
			if cStr != "" {

				cID, err := strconv.Atoi(cStr)
				if err != nil {
					for _, v := range cults {
						// Take first cult in map
						cID = int(v.ID)
						break
					}
				}

				parentCultModel, err = database.PKLoadCultModel(db, int64(cID))
				if err != nil {
					fmt.Println("No Cult Found")
				}

				cl.Cult.ParentCult = parentCultModel.Cult

				fmt.Println("ParentCULT: " + cl.Cult.Name)
			}
		}

		// Upload image to s3
		file, h, err := req.FormFile("image")
		switch err {
		case nil:
			if h.Filename != "" {
				// Process image
				defer file.Close()

				err = ProcessCultImage(h, file, &cl)
				if err != nil {
					log.Printf("Error processing image: %v", err)
				}

			} else {
				fmt.Println("No file provided.")
			}

		case http.ErrMissingFile:
			log.Println("no file")
			cl.Image = new(models.Image)
			cl.Image.Path = DefaultCharacterPortrait

		default:
			log.Panic(err)
			fmt.Println("Error getting file ", err)
			cl.Image = new(models.Image)
			cl.Image.Path = DefaultCharacterPortrait
		}

		// Runes

		for i := 1; i < 4; i++ {

			r := req.FormValue(fmt.Sprintf("Rune-%d", i))

			if r != "" {
				cl.Cult.Runes = append(cl.Cult.Runes, r)
			}
		}

		// Rune Spells

		for _, rs := range runequest.RuneSpells {
			str := req.FormValue(fmt.Sprintf("RS-%s-CoreString", rs.CoreString))
			if str != "" {

				t := &runequest.Spell{
					CoreString: rs.CoreString,
					Points:     rs.Points,
					UserChoice: rs.UserChoice,
					Domain:     rs.Domain,
					Source:     rs.Source,
					Variable:   rs.Variable,
					Cost:       rs.Cost,
				}
				if t.UserChoice {
					t.UserString = req.FormValue(fmt.Sprintf("RS-%s-UserString", rs.CoreString))
				}
				cl.Cult.RuneSpells = append(cl.Cult.RuneSpells, t)
			}
		}

		// Spirit Magic

		for _, sm := range runequest.SpiritMagicSpells {
			str := req.FormValue(fmt.Sprintf("SM-%s-CoreString", sm.CoreString))
			if str != "" {

				t := &runequest.Spell{
					CoreString: sm.CoreString,
					Points:     sm.Points,
					UserChoice: sm.UserChoice,
					Domain:     sm.Domain,
					Source:     sm.Source,
					Variable:   sm.Variable,
					Cost:       sm.Cost,
				}

				if t.UserChoice {
					t.UserString = req.FormValue(fmt.Sprintf("SM-%s-UserString", sm.CoreString))
				}
				cl.Cult.SpiritMagic = append(cl.Cult.SpiritMagic, t)
			}
		}

		// Associated Cults

		for _, ac := range cults {
			c := ac.Cult

			if !c.SubCult {

				str := req.FormValue(fmt.Sprintf("Cult-%s-Name", c.Name))

				if str != "" {
					tempCult := runequest.Cult{
						Name:        c.Name,
						RuneSpells:  c.RuneSpells,
						SpiritMagic: c.SpiritMagic,
						SubCult:     false,
					}
					cl.Cult.AssociatedCults = append(cl.Cult.AssociatedCults, tempCult)
				}
			}
		}

		// Read Base Skills

		skillArray := []*runequest.Skill{}

		// Add common skills for non-SubCults
		if !cl.Cult.SubCult {

			cultName := cl.Cult.Name

			worship := &runequest.Skill{
				CoreString: "Worship",
				UserChoice: true,
				UserString: cultName,
				Category:   "Magic",
				Base:       5,
				CultValue:  20,
			}

			meditate := &runequest.Skill{
				CoreString: "Meditate",
				Category:   "Magic",
				Base:       0,
				CultValue:  5,
			}

			lore := &runequest.Skill{
				CoreString: "Cult Lore",
				UserChoice: true,
				UserString: cultName,
				Category:   "Knowledge",
				Base:       5,
				CultValue:  15,
			}

			skillArray = append(skillArray, worship)
			skillArray = append(skillArray, meditate)
			skillArray = append(skillArray, lore)
		}

		for i := 1; i < 7; i++ {

			sk := req.FormValue(fmt.Sprintf("Skill-%d-CoreString", i))

			if sk != "" {

				skbaseSkill := runequest.Skills[sk]
				fmt.Println(skbaseSkill)

				// Skill
				s1 := &runequest.Skill{
					CoreString: skbaseSkill.CoreString,
					UserChoice: skbaseSkill.UserChoice,
					Category:   skbaseSkill.Category,
				}

				str := fmt.Sprintf("Skill-%d-Value", i)
				v, err := strconv.Atoi(req.FormValue(str))
				if err != nil {
					v = 0
				}
				s1.CultValue = v

				if s1.UserChoice {
					userString := fmt.Sprintf("Skill-%d-UserString", i)
					s1.UserString = req.FormValue(userString)
				}
				skillArray = append(skillArray, s1)
			}
		}

		cl.Cult.Skills = skillArray

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
				cl.Cult.Weapons = append(cl.Cult.Weapons, ws)
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
				s1.CultValue = v

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
				s2.CultValue = v

				if s2.UserChoice {
					userString := fmt.Sprintf("Skill-%d-2-UserString", i)
					s2.UserString = req.FormValue(userString)
				}

				// Form SkillChoice
				sc.Skills = []runequest.Skill{
					s1,
					s2,
				}
				// Append skillchoice
				cl.Cult.SkillChoices = append(cl.Cult.SkillChoices, sc)
			}
		}

		// Read passions
		for i := 1; i < 7; i++ {

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
					Type:       "Passion",
					CoreString: coreString,
					Base:       base,
					CultValue:  v,
					UserString: userString,
				}

				cl.Cult.PassionList = append(cl.Cult.PassionList, p)
			}
		}

		// Read New Skills
		for i := 1; i < 4; i++ {

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
				sk.CultValue = v

				cl.Cult.Skills = append(cl.Cult.Skills, sk)
			}
		}

		// Add other CultModel fields

		author, err := database.LoadUser(db, username)
		if err != nil {
			fmt.Println(err)
			http.Redirect(w, req, "/", 302)
		}

		cl.Author = author

		cl.Slug = slug.Make(cl.Cult.Name)

		err = database.SaveCultModel(db, &cl)
		if err != nil {
			log.Panic(err)
		} else {
			fmt.Println("Saved Cult")
		}

		url := fmt.Sprintf("/view_cult/%s", cl.Slug)

		http.Redirect(w, req, url, http.StatusFound)
	}
}

// ModifyCultHandler renders an editable Cult in a Web page
func ModifyCultHandler(w http.ResponseWriter, req *http.Request) {

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

	cl, err := database.PKLoadCultModel(db, int64(id))
	if err != nil {
		fmt.Println(err)
	}

	if cl.Author == nil {
		cl.Author = &models.User{
			UserName: "",
		}
	}

	// Validate that User == Author
	IsAuthor := false

	if username == cl.Author.UserName || isAdmin == "true" {
		IsAuthor = true
	} else {
		http.Redirect(w, req, "/", 302)
	}

	// Add extra runes
	if len(cl.Cult.Runes) < 4 {
		for i := len(cl.Cult.Runes); i < 4; i++ {
			cl.Cult.Runes = append(cl.Cult.Runes, "")
		}
	}

	// Add extra empty skills if < 10
	if len(cl.Cult.Skills) < 10 {
		for i := len(cl.Cult.Skills); i < 10; i++ {
			tempSkill := &runequest.Skill{}
			cl.Cult.Skills = append(cl.Cult.Skills, tempSkill)
		}
	}

	// Add extra empty skillchoices if < 3
	if len(cl.Cult.SkillChoices) < 4 {
		for i := len(cl.Cult.SkillChoices); i < 4; i++ {
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
			cl.Cult.SkillChoices = append(cl.Cult.SkillChoices, tempSkillChoice)
		}
	}

	// Add empty Weapon options if < 3
	if len(cl.Cult.Weapons) < 4 {
		for i := len(cl.Cult.Weapons); i < 4; i++ {
			tempWeapon := runequest.WeaponSelection{
				Description: "",
				Value:       0,
			}
			cl.Cult.Weapons = append(cl.Cult.Weapons, tempWeapon)
		}
	}

	cults, err := database.ListCultModels(db)
	if err != nil {
		panic(err)
	}

	wc := WebChar{
		CultModel:        cl,
		IsAuthor:         IsAuthor,
		SessionUser:      username,
		IsLoggedIn:       loggedIn,
		IsAdmin:          isAdmin,
		Counter:          numToArray(9),
		BigCounter:       numToArray(15),
		Passions:         runequest.PassionTypes,
		WeaponCategories: runequest.WeaponCategories,
		CategoryOrder:    runequest.CategoryOrder,
		PowerRunes:       runequest.PowerRuneOrder,
		ElementalRunes:   runequest.ElementalRuneOrder,
		Skills:           runequest.Skills,
		CultModels:       cults,
		RuneSpells:       runequest.RuneSpells,
		SpiritMagic:      runequest.SpiritMagicSpells,
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/modify_cult.html", wc)

	}

	if req.Method == "POST" { // POST

		err := req.ParseMultipartForm(MaxMemory)
		if err != nil {
			panic(err)
		}

		// Update Cult here
		cl.Cult.Name = req.FormValue("Name")
		cl.Cult.Description = req.FormValue("Description")
		cl.Cult.Notes = req.FormValue("Notes")

		if req.FormValue("Official") != "" {
			cl.Official = true
		} else {
			cl.Official = false
		}

		// Open or Official
		if req.FormValue("SubCult") != "" {
			cl.Cult.SubCult = true
		} else {
			cl.Cult.SubCult = false
		}

		parentCultModel := &models.CultModel{}

		if cl.Cult.SubCult {
			// Set ParentCult

			cStr := req.FormValue("ParentCult")
			if cStr != cl.Cult.ParentCult.Name && cStr != "" {

				cID, err := strconv.Atoi(cStr)
				if err != nil {
					for _, v := range cults {
						// Take first cult in map
						cID = int(v.ID)
						break
					}
				}

				parentCultModel, err = database.PKLoadCultModel(db, int64(cID))
				if err != nil {
					fmt.Println("No Cult Found")
				}

				cl.Cult.ParentCult = parentCultModel.Cult

				fmt.Println("ParentCULT: " + cl.Cult.Name)
			}
		}

		// Upload image to s3
		file, h, err := req.FormFile("image")
		switch err {
		case nil:
			if h.Filename != "" {
				// Process image
				defer file.Close()

				err = ProcessCultImage(h, file, cl)
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

		// Runes

		tempRunes := []string{}

		for i := 1; i < 4; i++ {

			r := req.FormValue(fmt.Sprintf("Rune-%d", i))

			if r != "" {
				tempRunes = append(tempRunes, r)
			}
		}
		cl.Cult.Runes = tempRunes

		// Rune Spells

		tempRuneSpells := []*runequest.Spell{}

		for _, rs := range runequest.RuneSpells {
			str := req.FormValue(fmt.Sprintf("RS-%s-CoreString", rs.CoreString))
			if str != "" {

				t := &runequest.Spell{
					CoreString: rs.CoreString,
					Points:     rs.Points,
					UserChoice: rs.UserChoice,
					Domain:     rs.Domain,
					Source:     rs.Source,
					Variable:   rs.Variable,
					Cost:       rs.Cost,
				}

				if t.UserChoice {
					t.UserString = req.FormValue(fmt.Sprintf("RS-%s-UserString", rs.CoreString))
				}
				tempRuneSpells = append(tempRuneSpells, t)
			}
		}
		cl.Cult.RuneSpells = tempRuneSpells

		// Spirit Magic
		tempSpiritMagic := []*runequest.Spell{}

		for _, sm := range runequest.SpiritMagicSpells {
			str := req.FormValue(fmt.Sprintf("SM-%s-CoreString", sm.CoreString))
			if str != "" {

				t := &runequest.Spell{
					CoreString: sm.CoreString,
					Points:     sm.Points,
					UserChoice: sm.UserChoice,
					Domain:     sm.Domain,
					Source:     sm.Source,
					Variable:   sm.Variable,
					Cost:       sm.Cost,
				}

				if t.UserChoice {
					t.UserString = req.FormValue(fmt.Sprintf("SM-%s-UserString", sm.CoreString))
				}
				tempSpiritMagic = append(tempSpiritMagic, t)
			}
		}

		cl.Cult.SpiritMagic = tempSpiritMagic

		// Associated Cults
		tempAssociatedCults := []runequest.Cult{}

		for _, ac := range cults {
			c := ac.Cult

			if !c.SubCult {

				str := req.FormValue(fmt.Sprintf("Cult-%s-Name", c.Name))

				if str != "" {
					tempCult := runequest.Cult{
						Name:        c.Name,
						RuneSpells:  c.RuneSpells,
						SpiritMagic: c.SpiritMagic,
						SubCult:     false,
					}
					tempAssociatedCults = append(tempAssociatedCults, tempCult)
				}
			}
		}
		cl.Cult.AssociatedCults = tempAssociatedCults

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
					for _, ns := range cl.Cult.Skills {
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

				s1.CultValue = v

				if userString != "" {
					s1.UserString = userString
				}
				s1.GenerateName()
				tempSkills = append(tempSkills, s1)
			}
		}

		// Set Cult.Skills to new array
		cl.Cult.Skills = tempSkills

		// Read Weapons
		weapons := []runequest.WeaponSelection{}

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
				weapons = append(weapons, ws)
			}
		}

		cl.Cult.Weapons = weapons

		// Read passions
		tempPassions := []runequest.Ability{}

		for i := 1; i < 4; i++ {

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
					Type:       "Passion",
					CoreString: coreString,
					Base:       base,
					CultValue:  v,
					UserString: userString,
				}

				tempPassions = append(tempPassions, p)
			}
		}

		// Reset passions to new values
		cl.Cult.PassionList = tempPassions

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
				s1.CultValue = v

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
				s2.CultValue = v

				if s2.UserChoice {
					userString := fmt.Sprintf("Skill-%d-2-UserString", i)
					s2.UserString = req.FormValue(userString)
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

		cl.Cult.SkillChoices = tempChoices

		// Read New Skills

		for i := 1; i < 4; i++ {

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
				sk.CultValue = v

				cl.Cult.Skills = append(cl.Cult.Skills, sk)
			}
		}

		cl.Slug = slug.Make(cl.Cult.Name)

		// Insert Cult into App archive
		err = database.UpdateCultModel(db, cl)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("Saved")
		}

		fmt.Println(cl)

		url := fmt.Sprintf("/view_cult/%s", cl.Slug)

		http.Redirect(w, req, url, http.StatusSeeOther)
	}
}

// DeleteCultHandler renders a character in a Web page
func DeleteCultHandler(w http.ResponseWriter, req *http.Request) {

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

	cl, err := database.PKLoadCultModel(db, int64(id))
	if err != nil {
		fmt.Println(err)
	}

	if cl.Image == nil {
		cl.Image = new(models.Image)
		cl.Image.Path = DefaultCharacterPortrait
	}

	// Validate that User == Author
	IsAuthor := false

	if username == cl.Author.UserName || isAdmin == "true" {
		IsAuthor = true
	} else {
		http.Redirect(w, req, "/", 302)
	}

	wc := WebChar{
		CultModel:   cl,
		IsAuthor:    IsAuthor,
		SessionUser: username,
		IsLoggedIn:  loggedIn,
		IsAdmin:     isAdmin,
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/delete_cult.html", wc)

	}

	if req.Method == "POST" {

		err := database.DeleteCultModel(db, cl.ID)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("Deleted Cult")
		}

		url := fmt.Sprint("/cult_index/")

		http.Redirect(w, req, url, http.StatusSeeOther)
	}
}
