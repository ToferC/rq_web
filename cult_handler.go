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
	pk := vars["id"]

	if len(pk) == 0 {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}

	id, err := strconv.Atoi(pk)
	if err != nil {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}

	cl, err := database.PKLoadCultModel(db, int64(id))
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

	c.Statistics["STR"].Value = 10
	c.Statistics["DEX"].Value = 10
	c.Statistics["INT"].Value = 10
	c.Statistics["POW"].Value = 10
	c.Statistics["CHA"].Value = 10
	c.Statistics["SIZ"].Value = 10

	wc := WebChar{
		CharacterModel:   &cm,
		IsAuthor:         true,
		SessionUser:      username,
		IsLoggedIn:       loggedIn,
		IsAdmin:          isAdmin,
		Counter:          numToArray(3),
		BigCounter:       numToArray(15),
		Passions:         runequest.PassionTypes,
		WeaponCategories: runequest.WeaponCategories,
		CategoryOrder:    runequest.CategoryOrder,
		PowerRunes:       runequest.PowerRuneOrder,
		ElementalRunes:   runequest.ElementalRuneOrder,
		Skills:           runequest.Skills,
		CultModels:       cults,
		SubCults:         []runequest.Cult{},
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

		user := database.LoadUser(db, username)

		// Map default Cult to Character.Cults
		cl := models.CultModel{
			Author: user,
			Cult: &runequest.Cult{
				Name:        req.FormValue("Name"),
				Description: req.FormValue("Description"),
				//CultChoices:
			},
		}

		if req.FormValue("Official") != "" {
			cl.Official = true
		} else {
			cl.Official = false
		}

		// Upload image to s3
		file, h, err := req.FormFile("image")
		switch err {
		case nil:
			// Process image
			defer file.Close()
			// example path media/Major/TestImage/Jason_White.jpg
			path := fmt.Sprintf("/media/%s/%s/%s",
				cl.Author.UserName,
				runequest.ToSnakeCase(cl.Cult.Name),
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

			cl.Image = new(models.Image)
			cl.Image.Path = path

			fmt.Println(path)

		case http.ErrMissingFile:
			log.Println("no file")

		default:
			log.Panic(err)
			fmt.Println("Error getting file ", err)
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

				t := runequest.Spell{
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

				t := runequest.Spell{
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
						Name:       c.Name,
						RuneSpells: c.RuneSpells,
						SubCult:    false,
					}
					cl.Cult.AssociatedCults = append(cl.Cult.AssociatedCults, tempCult)
				}
			}
		}

		// SubCult

		for _, ac := range cults {
			c := ac.Cult

			if c.SubCult {

				str := req.FormValue(fmt.Sprintf("Subcult-%s-Name", c.Name))

				if str != "" {
					tempCult := runequest.Cult{
						Name:       c.Name,
						RuneSpells: c.RuneSpells,
						SubCult:    true,
					}
					cl.Cult.SubCults = append(cl.Cult.SubCults, tempCult)
				}
			}
		}

		// Read Base Skills

		skillArray := []runequest.Skill{}

		// Add common skills

		worship := runequest.Skill{
			CoreString: "Worship",
			UserChoice: true,
			UserString: cl.Cult.Name,
			Category:   "Magic",
			Base:       5,
			CultValue:  20,
		}

		meditate := runequest.Skill{
			CoreString: "Meditate",
			Category:   "Magic",
			Base:       0,
			CultValue:  5,
		}

		lore := runequest.Skill{
			CoreString: "Cult Lore",
			UserChoice: true,
			UserString: cl.Cult.Name,
			Category:   "Knowledge",
			Base:       5,
			CultValue:  15,
		}

		skillArray = append(skillArray, worship)
		skillArray = append(skillArray, meditate)
		skillArray = append(skillArray, lore)

		for i := 1; i < 4; i++ {

			sk := req.FormValue(fmt.Sprintf("Skill-%d-CoreString", i))

			if sk != "" {

				skbaseSkill := runequest.Skills[sk]
				fmt.Println(skbaseSkill)

				// Skill
				s1 := runequest.Skill{
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
					s2.UserString = userString
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
				p.CultValue = v

				userString := req.FormValue(fmt.Sprintf("Passion-%d-UserString", i))

				if userString != "" {
					p.UserChoice = true
					p.UserString = userString
				}

				cl.Cult.PassionList = append(cl.Cult.PassionList, p)
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
				sk.CultValue = v

				cl.Cult.Skills = append(cl.Cult.Skills, sk)
			}
		}

		// Add other CultModel fields

		author := database.LoadUser(db, username)

		cl.Author = author

		err = database.SaveCultModel(db, &cl)
		if err != nil {
			log.Panic(err)
		} else {
			fmt.Println("Saved")
		}

		url := fmt.Sprintf("/view_cult/%d", cl.ID)

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

	if username == cl.Author.UserName {
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

	// Add extra empty skills if < 6
	if len(cl.Cult.Skills) < 8 {
		for i := len(cl.Cult.Skills); i < 8; i++ {
			tempSkill := runequest.Skill{}
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
	if len(cl.Cult.Weapons) < 3 {
		for i := len(cl.Cult.Weapons); i < 3; i++ {
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
		Counter:          []int{1, 2, 3},
		BigCounter:       []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
		Passions:         runequest.PassionTypes,
		WeaponCategories: runequest.WeaponCategories,
		CategoryOrder:    runequest.CategoryOrder,
		PowerRunes:       runequest.PowerRuneOrder,
		ElementalRunes:   runequest.ElementalRuneOrder,
		Skills:           runequest.Skills,
		CultModels:       cults,
		SubCults:         []runequest.Cult{},
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
		clName := req.FormValue("Name")
		description := req.FormValue("Description")

		cl.Cult.Name = clName
		cl.Cult.Description = description

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

		// Upload image to s3
		file, h, err := req.FormFile("image")
		switch err {
		case nil:
			// Prcless image
			defer file.Close()
			// example path media/Major/TestImage/Jason_White.jpg
			path := fmt.Sprintf("/media/%s/%s/%s",
				cl.Author.UserName,
				runequest.ToSnakeCase(cl.Cult.Name),
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

			cl.Image = new(models.Image)
			cl.Image.Path = path

			fmt.Println(path)

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

		tempRuneSpells := []runequest.Spell{}

		for _, rs := range runequest.RuneSpells {
			str := req.FormValue(fmt.Sprintf("RS-%s-CoreString", rs.CoreString))
			if str != "" {

				t := runequest.Spell{
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
		tempSpiritMagic := []runequest.Spell{}

		for _, sm := range runequest.SpiritMagicSpells {
			str := req.FormValue(fmt.Sprintf("SM-%s-CoreString", sm.CoreString))
			if str != "" {

				t := runequest.Spell{
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
						Name:       c.Name,
						RuneSpells: c.RuneSpells,
						SubCult:    false,
					}
					tempAssociatedCults = append(tempAssociatedCults, tempCult)
				}
			}
		}
		cl.Cult.AssociatedCults = tempAssociatedCults

		// SubCult
		tempSubCults := []runequest.Cult{}

		for _, ac := range cults {
			c := ac.Cult

			if c.SubCult {

				str := req.FormValue(fmt.Sprintf("SubCult-%s-Name", c.Name))

				if str != "" {
					tempCult := runequest.Cult{
						Name:       c.Name,
						RuneSpells: c.RuneSpells,
						SubCult:    true,
					}
					tempSubCults = append(tempSubCults, tempCult)
				}
			}
		}
		cl.Cult.SubCults = tempSubCults

		// Read Skills
		tempSkills := []runequest.Skill{}

		// Read Base Skills from Form
		for i := 1; i < 4; i++ {

			core := req.FormValue(fmt.Sprintf("Skill-%d-CoreString", i))

			str := fmt.Sprintf("Skill-%d-Value", i)
			v, err := strconv.Atoi(req.FormValue(str))
			if err != nil {
				v = 0
			}

			if core != "" && v > 0 {

				sk := runequest.Skill{
					CoreString: core,
					CultValue:  v,
				}

				userString := req.FormValue(fmt.Sprintf("Skill-%d-UserString", i))
				if userString != "" {
					sk.UserString = userString
				}
				tempSkills = append(tempSkills, sk)
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
				p.CultValue = v

				userString := req.FormValue(fmt.Sprintf("Passion-%d-UserString", i))

				if userString != "" {
					p.UserChoice = true
					p.UserString = userString
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

		cl.Cult.SkillChoices = tempChoices

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
				sk.CultValue = v

				cl.Cult.Skills = append(cl.Cult.Skills, sk)
			}
		}

		// Insert Cult into App archive
		err = database.UpdateCultModel(db, cl)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("Saved")
		}

		fmt.Println(cl)

		url := fmt.Sprintf("/view_cult/%d", cl.ID)

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
