package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

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

	occupations, err := database.ListOccupationModels(db)
	if err != nil {
		panic(err)
	}

	cults, err := database.ListCultModels(db)
	if err != nil {
		panic(err)
	}

	wc := WebChar{
		CharacterModel:   &cm,
		SessionUser:      username,
		IsLoggedIn:       loggedIn,
		IsAdmin:          isAdmin,
		HomelandModels:   homelands,
		OccupationModels: occupations,
		CultModels:       cults,
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
			for _, k := range homelands {
				// Get first homeland
				hlID = int(k.ID)
			}
			fmt.Println(err)
		}

		hl, err := database.PKLoadHomelandModel(db, int64(hlID))
		if err != nil {
			fmt.Println(err)
		}

		c.Homeland = hl.Homeland
		fmt.Println("HOMELAND: " + c.Homeland.Name)

		// Set Occupation
		ocStr := req.FormValue("Occupation")

		ocID, err := strconv.Atoi(ocStr)
		if err != nil {
			for _, v := range occupations {
				// Take first occupation in map
				ocID = int(v.ID)
				break
			}
		}

		oc, err := database.PKLoadOccupationModel(db, int64(ocID))
		if err != nil {
			fmt.Println("No Occupation Found")
		}

		c.Occupation = oc.Occupation
		fmt.Println("OCCUPATION: " + c.Occupation.Name)

		// Set Cult
		cStr := req.FormValue("Cult")

		cID, err := strconv.Atoi(cStr)
		if err != nil {
			for _, v := range cults {
				// Take first cult in map
				cID = int(v.ID)
				break
			}
		}

		cultModel, err := database.PKLoadCultModel(db, int64(cID))
		if err != nil {
			fmt.Println("No Cult Found")
		}

		c.Cult = cultModel.Cult
		fmt.Println("CULT: " + c.Cult.Name)

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
			cm.Image = new(models.Image)
			cm.Image.Path = DefaultCharacterPortrait

		default:
			log.Panic(err)
			fmt.Println("Error getting file ", err)
			cm.Image = new(models.Image)
			cm.Image.Path = DefaultCharacterPortrait
		}

		cm.Open = true
		//fmt.Println(c)

		err = database.SaveCharacterModel(db, &cm)
		if err != nil {
			log.Panic(err)
		} else {
			fmt.Println("Saved")
		}

		url := fmt.Sprintf("/cc12_personal_history/%d", cm.ID)

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
		Counter:          numToArray(5),
		SessionUser:      username,
		IsLoggedIn:       loggedIn,
		IsAdmin:          isAdmin,
		IsAuthor:         IsAuthor,
		Skills:           runequest.Skills,
		Passions:         runequest.PassionTypes,
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

		// Read passions
		for i := 1; i < 6; i++ {

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

				c.ModifyAbility(p)
			}
		}

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

		// Set Elemental Runes
		for i := 1; i < 4; i++ {
			eRune := req.FormValue(fmt.Sprintf("Rune%d", i))

			switch i {
			case 1:
				c.ElementalRunes[eRune].Base = 60
			case 2:
				c.ElementalRunes[eRune].Base = 40
			case 3:
				c.ElementalRunes[eRune].Base = 20
			}
		}

		// Set Elemental Runes
		for i := 1; i < 3; i++ {
			pRune := req.FormValue(fmt.Sprintf("PowerRune%d", i))
			c.PowerRunes[pRune].Base = 75
		}

		// Update Elemental Rune Bonus
		for k := range c.ElementalRunes {

			cbv, err := strconv.Atoi(req.FormValue("Bonus-" + k))
			if err != nil {
				fmt.Println(err)
				cbv = 0
			}

			// Modify Ability
			c.ElementalRunes[k].CreationBonusValue = cbv
		}

		for k, v := range c.PowerRunes {

			cbv, _ := strconv.Atoi(req.FormValue("Bonus-" + k))
			if err != nil {
				cbv = 0
			}
			// Modify Ability
			c.PowerRunes[k].CreationBonusValue = cbv
			v.UpdateAbility()
		}

		// Determine opposing Power Rune values
		for k, v := range c.PowerRunes {
			if v.Total > 50 {
				c.UpdateOpposedRune(c.PowerRunes[k])
			}
		}

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

	// Determine Rune bonuses (returns array of 2 strings)
	runeArray := c.DetermineRuneModifiers()

	wc := WebChar{
		CharacterModel: cm,
		SessionUser:    username,
		IsLoggedIn:     loggedIn,
		IsAdmin:        isAdmin,
		IsAuthor:       IsAuthor,
		RuneArray:      runeArray,
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
		}

		// Apply Rune Modifiers
		target1 := req.FormValue(fmt.Sprintf("RuneMod-%d", 0))
		c.Statistics[target1].RuneBonus = 2

		target2 := req.FormValue(fmt.Sprintf("RuneMod-%d", 1))
		c.Statistics[target2].RuneBonus = 1

		c.SetAttributes()

		// Update Character

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

		// Add choices to c.Homeland.Skills
		for i, sc := range c.Homeland.SkillChoices {
			for m := range sc.Skills {
				str := fmt.Sprintf("SC-%d-%d", i, m)
				if req.FormValue(str) != "" {
					c.Homeland.Skills = append(c.Homeland.Skills, sc.Skills[m])
				}
			}
		}

		// Add Skills to Character
		for i, s := range c.Homeland.Skills {
			// Modify Skill
			if s.UserString == "any" {
				// User Chooses a new specialization
				str := req.FormValue(fmt.Sprintf("Skill-%d-UserString", i))
				if str != "" {
					s.UserString = str
				}
			}

			var targetString string

			// Find target string for Skill
			if s.UserString != "" {
				targetString = fmt.Sprintf("%s (%s)", s.CoreString, s.UserString)
			} else {
				targetString = fmt.Sprintf("%s", s.CoreString)
			}

			if c.Skills[targetString] != nil {
				// Skill exists in Character, modify it via pointer
				c.Skills[targetString].HomelandValue = s.HomelandValue

				fmt.Println("Updated Character Skill: " + c.Skills[targetString].Name)

			} else {
				// We need to find the base skill from the Master list or create it
				if runequest.Skills[s.CoreString] == nil {
					fmt.Println("Skill is new: " + targetString)

					// New Skill
					baseSkill := &s
					// Update our new skill
					sc := c.SkillCategories[baseSkill.Category]

					baseSkill.CategoryValue = sc.Value

					baseSkill.UpdateSkill()

					fmt.Println("Add Skill to character: " + baseSkill.Name)
					c.Skills[baseSkill.Name] = baseSkill
				} else {
					// Skill exists in master list
					fmt.Println("Skill in master list: " + targetString)

					bs := runequest.Skills[s.CoreString]

					baseSkill := &runequest.Skill{
						CoreString: bs.CoreString,
						UserString: bs.UserString,
						Category:   bs.Category,
						Base:       bs.Base,
						UserChoice: bs.UserChoice,
					}

					fmt.Println(baseSkill)
					baseSkill.HomelandValue = s.HomelandValue
					fmt.Println(s.HomelandValue, baseSkill.HomelandValue)

					if s.Base > baseSkill.Base {
						baseSkill.Base = s.Base
					}
					if s.UserString != "" {
						baseSkill.UserString = s.UserString
					}

					baseSkill.GenerateName()

					// Add Skill to Character
					fmt.Println("Add Skill to character: " + baseSkill.Name)
					c.Skills[baseSkill.Name] = baseSkill
					c.Skills[baseSkill.Name].UpdateSkill()
				}
			}
		}

		// Homelands grant 3 base passions
		// Find number of abilities

		for _, selected := range c.Homeland.PassionList {
			c.Homeland.Passions = append(c.Homeland.Passions, selected)
			c.ModifyAbility(selected)
		}

		// Homeland grants a bonus to a rune affinity
		c.ElementalRunes[c.Homeland.RuneBonus].HomelandValue += 10

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

	c := cm.Character

	IsAuthor := false

	if username == cm.Author.UserName {
		IsAuthor = true
	}

	wc := WebChar{
		CharacterModel:   cm,
		SessionUser:      username,
		IsLoggedIn:       loggedIn,
		IsAdmin:          isAdmin,
		IsAuthor:         IsAuthor,
		Skills:           runequest.Skills,
		WeaponCategories: runequest.WeaponCategories,
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

		// Add choices to c.Occupation.Skills
		for i, sc := range c.Occupation.SkillChoices {
			for m := range sc.Skills {
				str := fmt.Sprintf("SC-%d-%d", i, m)
				if req.FormValue(str) != "" {
					c.Occupation.Skills = append(c.Occupation.Skills, sc.Skills[m])
				}
			}
		}

		// Add choices for Weapon skills
		for i, w := range c.Occupation.Weapons {
			str := fmt.Sprintf("Weapon-%d-CoreString", i)
			fv := req.FormValue(str)
			if fv != "" {

				bs := runequest.Skills[fv]

				ws := runequest.Skill{
					CoreString:      bs.CoreString,
					Category:        bs.Category,
					Base:            bs.Base,
					OccupationValue: w.Value,
				}
				c.Occupation.Skills = append(c.Occupation.Skills, ws)
			}
		}

		// Add Skills to Character
		for i, s := range c.Occupation.Skills {
			// Modify Skill
			if s.UserString == "any" {
				// User Chooses a new specialization
				str := req.FormValue(fmt.Sprintf("Skill-%d-UserString", i))
				if str != "" {
					s.UserString = str
				}
			}

			var targetString string

			// Find target string for Skill
			if s.UserString != "" {
				targetString = fmt.Sprintf("%s (%s)", s.CoreString, s.UserString)
			} else {
				targetString = fmt.Sprintf("%s", s.CoreString)
			}

			if c.Skills[targetString] != nil {
				// Skill exists in Character, modify it via pointer
				c.Skills[targetString].OccupationValue = s.OccupationValue

				if s.Base > c.Skills[targetString].Base {
					c.Skills[targetString].Base = s.Base
				}

				fmt.Println("Updated Character Skill: " + c.Skills[targetString].Name)

			} else {
				// We need to find the base skill from the Master list or create it
				if runequest.Skills[s.CoreString] == nil {
					fmt.Println("Skill is new: " + targetString)

					// New Skill
					baseSkill := &s
					// Update our new skill
					sc := c.SkillCategories[baseSkill.Category]

					baseSkill.CategoryValue = sc.Value

					baseSkill.UpdateSkill()

					fmt.Println("Add Skill to character: " + baseSkill.Name)
					c.Skills[baseSkill.Name] = baseSkill
				} else {
					// Skill exists in master list
					fmt.Println("Skill in master list: " + targetString)

					bs := runequest.Skills[s.CoreString]

					baseSkill := &runequest.Skill{
						CoreString: bs.CoreString,
						UserString: bs.UserString,
						Category:   bs.Category,
						Base:       bs.Base,
						UserChoice: bs.UserChoice,
					}

					fmt.Println(baseSkill)
					baseSkill.OccupationValue = s.OccupationValue
					fmt.Println(s.OccupationValue, baseSkill.OccupationValue)

					if s.Base > baseSkill.Base {
						baseSkill.Base = s.Base
					}
					if s.UserString != "" {
						baseSkill.UserString = s.UserString
					}

					baseSkill.GenerateName()

					// Add Skill to Character
					fmt.Println("Add Skill to character: " + baseSkill.Name)
					c.Skills[baseSkill.Name] = baseSkill
					c.Skills[baseSkill.Name].UpdateSkill()
				}
			}
		}

		// Occupations grant a bonus to one Passion
		if len(c.Occupation.PassionList) > 0 {
			str := req.FormValue("Passion")
			n, err := strconv.Atoi(str)
			if err != nil {
				n = 0
			}

			c.ModifyAbility(c.Occupation.PassionList[n])
		}

		// Equipment

		if len(c.Occupation.Equipment) > 0 {
			for _, e := range c.Occupation.Equipment {
				c.Equipment = append(c.Equipment, e)
			}
		}

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

	c := cm.Character

	IsAuthor := false

	if username == cm.Author.UserName {
		IsAuthor = true
	}

	numRunePoints := numToArray(3)
	numSpiritMagic := numToArray(5)

	// Add Associated Cult spells
	totalSpiritMagic := []runequest.Spell{}

	totalSpiritMagic = c.Cult.SpiritMagic

	for _, ac := range c.Cult.AssociatedCults {
		for _, spell := range ac.SpiritMagic {
			totalSpiritMagic = append(totalSpiritMagic, spell)
		}
	}

	wc := WebChar{
		CharacterModel: cm,
		SessionUser:    username,
		IsLoggedIn:     loggedIn,
		IsAdmin:        isAdmin,
		IsAuthor:       IsAuthor,
		Skills:         runequest.Skills,
		NumRunePoints:  numRunePoints,
		NumSpiritMagic: numSpiritMagic,
		SpiritMagic:    totalSpiritMagic,
	}
	// Test
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

		c.Cult.Rank = req.FormValue("Rank")

		c.RuneSpells = map[string]*runequest.Spell{}
		c.SpiritMagic = map[string]*runequest.Spell{}

		// Rune Magic
		for i := 1; i < 6; i++ {
			str := req.FormValue(fmt.Sprintf("RuneSpell-%d", i))
			spec := req.FormValue(fmt.Sprintf("RuneSpell-%d-UserString", i))
			if str != "" {
				index, err := strconv.Atoi(str)
				if err != nil {
					index = 0
					fmt.Println("Spell Not found")
				}
				baseSpell := c.Cult.RuneSpells[index]

				s := &baseSpell
				if spec != "" {
					s.UserString = spec
				}
				s.GenerateName()
				c.RuneSpells[s.Name] = s
			}
		}

		// Spirit Magic
		for i := 1; i < 6; i++ {
			str := req.FormValue(fmt.Sprintf("SpiritMagic-%d", i))
			spec := req.FormValue(fmt.Sprintf("SpiritMagic-%d-UserString", i))
			cString := req.FormValue(fmt.Sprintf("SpiritMagic-%d-Cost", i))

			if str != "" {

				index, err := strconv.Atoi(str)
				if err != nil {
					index = 0
					fmt.Println("Spell Not found")
				}

				cost, err := strconv.Atoi(cString)
				if err != nil {
					cost = 1
					fmt.Println("Non-number entered")
				}

				baseSpell := totalSpiritMagic[index]

				s := &runequest.Spell{
					Name:       baseSpell.Name,
					CoreString: baseSpell.CoreString,
					UserString: baseSpell.UserString,
					Cost:       cost,
					Domain:     baseSpell.Domain,
				}

				if spec != "" {
					s.UserString = spec
				}

				s.GenerateName()
				c.SpiritMagic[s.Name] = s
			}
		}

		// Add choices to c.Cult.Skills
		for i, sc := range c.Cult.SkillChoices {
			for m := range sc.Skills {
				str := fmt.Sprintf("SC-%d-%d", i, m)
				if req.FormValue(str) != "" {
					c.Cult.Skills = append(c.Cult.Skills, sc.Skills[m])
				}
			}
		}

		// Add choices for Weapon skills
		for i, w := range c.Cult.Weapons {
			str := fmt.Sprintf("Weapon-%d-CoreString", i)
			fv := req.FormValue(str)

			if fv != "" {
				bs := runequest.Skills[fv]

				ws := runequest.Skill{
					CoreString: bs.CoreString,
					Category:   bs.Category,
					Base:       bs.Base,
					CultValue:  w.Value,
				}
				c.Cult.Skills = append(c.Cult.Skills, ws)
			}
		}

		// Add Skills to Character
		for i, s := range c.Cult.Skills {
			// Modify Skill
			if s.UserString == "any" {
				// User Chooses a new specialization
				str := req.FormValue(fmt.Sprintf("Skill-%d-UserString", i))
				if str != "" {
					s.UserString = str
				}
			}

			var targetString string

			// Find target string for Skill
			if s.UserString != "" {
				targetString = fmt.Sprintf("%s (%s)", s.CoreString, s.UserString)
			} else {
				targetString = fmt.Sprintf("%s", s.CoreString)
			}

			if c.Skills[targetString] != nil {
				// Skill exists in Character, modify it via pointer
				c.Skills[targetString].CultValue = s.CultValue

				if s.Base > c.Skills[targetString].Base {
					c.Skills[targetString].Base = s.Base
				}

				fmt.Println("Updated Character Skill: " + c.Skills[targetString].Name)

			} else {
				// We need to find the base skill from the Master list or create it
				if runequest.Skills[s.CoreString] == nil {
					// Skill doesn't exist in master list
					fmt.Println("Skill is new: " + targetString)

					// New Skill
					baseSkill := &s
					// Update our new skill
					sc := c.SkillCategories[baseSkill.Category]

					baseSkill.CategoryValue = sc.Value
					baseSkill.UserString = s.UserString

					baseSkill.UpdateSkill()

					fmt.Println("Add Skill to character: " + baseSkill.Name)
					c.Skills[baseSkill.Name] = baseSkill
				} else {
					// Skill exists in master list
					fmt.Println("Skill in master list: " + targetString)

					bs := runequest.Skills[s.CoreString]

					baseSkill := &runequest.Skill{
						CoreString: bs.CoreString,
						UserString: bs.UserString,
						Category:   bs.Category,
						Base:       bs.Base,
						UserChoice: bs.UserChoice,
					}

					fmt.Println(baseSkill)
					baseSkill.CultValue = s.CultValue
					fmt.Println(s.CultValue, baseSkill.CultValue)

					if s.Base > baseSkill.Base {
						baseSkill.Base = s.Base
					}
					if s.UserString != "" {
						baseSkill.UserString = s.UserString
					}

					baseSkill.GenerateName()

					// Add Skill to Character
					fmt.Println("Add Skill to character: " + baseSkill.Name)
					c.Skills[baseSkill.Name] = baseSkill
					c.Skills[baseSkill.Name].UpdateSkill()
				}
			}
		}

		// Add 20 to one cult skill
		targetString20 := req.FormValue("Skill-20")

		// Skill exists in Character, modify it via pointer
		if targetString20 != "" {

			index, err := strconv.Atoi(targetString20)
			if err != nil {
				index = 0
				fmt.Println("Skill Not found")
			}

			baseSkill := c.Cult.Skills[index]
			targetSkill := &runequest.Skill{
				CoreString: baseSkill.CoreString,
				UserString: baseSkill.UserString,
			}

			targetSkill.GenerateName()

			t := time.Now()
			tString := t.Format("2006-01-02 15:04:05")

			update := &runequest.Update{
				Date:  tString,
				Event: "Cult Skill (20)",
				Value: 20,
			}

			s := c.Skills[targetSkill.Name]

			if s.Updates == nil {
				s.Updates = []*runequest.Update{}
			}

			s.Updates = append(s.Updates, update)

			s.UpdateSkill()

			fmt.Println("Updated Character Skill 20%: " + s.Name)
		}

		// Add 15 to one Cult Skill
		targetString15 := req.FormValue("Skill-15")

		// Skill exists in Character, modify it via pointer
		if targetString15 != "" {

			index, err := strconv.Atoi(targetString15)
			if err != nil {
				index = 0
				fmt.Println("Skill Not found")
			}

			baseSkill := c.Cult.Skills[index]
			targetSkill := &runequest.Skill{
				CoreString: baseSkill.CoreString,
				UserString: baseSkill.UserString,
			}

			targetSkill.GenerateName()

			t := time.Now()
			tString := t.Format("2006-01-02 15:04:05")

			update := &runequest.Update{
				Date:  tString,
				Event: "Cult Skill (15)",
				Value: 15,
			}

			s := c.Skills[targetSkill.Name]

			if s.Updates == nil {
				s.Updates = []*runequest.Update{}
			}

			s.Updates = append(s.Updates, update)

			s.UpdateSkill()

			fmt.Println("Updated Character Skill 10%: " + s.Name)
		}

		// Cults grant a bonus to one Passion
		str := req.FormValue("Passion")
		if str != "" {
			n, err := strconv.Atoi(str)
			if err != nil {
				n = 0
			}
			c.ModifyAbility(c.Cult.PassionList[n])
		}

		// Add base runepoints
		c.Cult.NumRunePoints = 3
		c.Cult.NumSpiritMagic = 5

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
		Counter:          numToArray(4),
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
		// 25% additions
		for i := 1; i < 5; i++ {
			coreString := req.FormValue(fmt.Sprintf("Skill-25-%d", i))
			userString := req.FormValue(fmt.Sprintf("Skill-25-%d-UserString", i))

			var targetString string

			// Find target string for Skill
			if userString != "" {
				targetString = fmt.Sprintf("%s (%s)", coreString, userString)
			} else {
				targetString = fmt.Sprintf("%s", coreString)
			}

			fmt.Println(targetString)

			// Skill exists in Character, modify it via pointer
			if targetString != "" {
				// Determine if skill already exists in c.Skills

				if c.Skills[targetString] == nil {
					// Skill isn't in c.Skills, so create new skill
					fmt.Println("Unable to find skill: " + targetString)

					bs := runequest.Skills[coreString]

					sk := &runequest.Skill{
						CoreString: bs.CoreString,
						UserString: bs.UserString,
						Category:   bs.Category,
						Base:       bs.Base,
						UserChoice: bs.UserChoice,
					}

					if userString != "" {
						sk.UserString = userString
					}

					sk.GenerateName()

					// Add Skill to Character
					fmt.Println("Add Skill to character: " + sk.Name)
					c.Skills[sk.Name] = sk
					c.Skills[sk.Name].UpdateSkill()
				}

				s := c.Skills[targetString]

				t := time.Now()
				tString := t.Format("2006-01-02 15:04:05")

				update := &runequest.Update{
					Date:  tString,
					Event: "Personal Skills (25)",
					Value: 25,
				}

				if s.Updates == nil {
					s.Updates = []*runequest.Update{}
				}

				s.Updates = append(s.Updates, update)

				s.UpdateSkill()

				fmt.Println("Updated Character Skill (25): " + s.Name)
			}
		}

		// 10% additions
		for i := 1; i < 5; i++ {
			coreString := req.FormValue(fmt.Sprintf("Skill-10-%d", i))
			userString := req.FormValue(fmt.Sprintf("Skill-10-%d-UserString", i))

			var targetString string

			// Find target string for Skill
			if userString != "" {
				targetString = fmt.Sprintf("%s (%s)", coreString, userString)
			} else {
				targetString = fmt.Sprintf("%s", coreString)
			}

			// Skill exists in Character, modify it via pointer
			if targetString != "" {
				// Skill isn't in c.Skills, so create new skill
				fmt.Println("Unable to find skill: " + targetString)

				s := c.Skills[targetString]

				if c.Skills[targetString] == nil {
					// Skill isn't in c.Skills, so create new skill
					fmt.Println("Unable to find skill: " + targetString)

					bs := runequest.Skills[coreString]

					sk := &runequest.Skill{
						CoreString: bs.CoreString,
						UserString: bs.UserString,
						Category:   bs.Category,
						Base:       bs.Base,
						UserChoice: bs.UserChoice,
					}

					if userString != "" {
						sk.UserString = userString
					}

					sk.GenerateName()

					// Add Skill to Character
					fmt.Println("Add Skill to character: " + sk.Name)
					c.Skills[sk.Name] = sk
					c.Skills[sk.Name].UpdateSkill()
				}

				s = c.Skills[targetString]

				t := time.Now()
				tString := t.Format("2006-01-02 15:04:05")

				update := &runequest.Update{
					Date:  tString,
					Event: "Personal Skills (10)",
					Value: 10,
				}

				s.Updates = append(s.Updates, update)

				s.UpdateSkill()

				fmt.Println("Updated Character Skill 10%: " + s.Name)
			}
		}

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
