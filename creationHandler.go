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

	author, err := database.LoadUser(db, username)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, req, "/", 302)
	}

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

		// Pull form values into cm.Character via gorilla/schema
		err = decoder.Decode(c, req.PostForm)
		if err != nil {
			panic(err)
		}

		// Set Homeland
		hlStr := req.FormValue("HLStr")

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
		ocStr := req.FormValue("OCStr")

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
		cStr := req.FormValue("CStr")

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

		if cultModel.Cult.SubCult {
			// Set base to ParentCult
			// Easier to add to
			parentCult := cultModel.Cult.ParentCult
			if parentCult != nil {
				// Add SubCult skills
				for _, s := range cultModel.Cult.Skills {
					parentCult.Skills = append(parentCult.Skills, s)
				}
				// Add SubCult weapons
				for _, w := range cultModel.Cult.Weapons {
					parentCult.Weapons = append(parentCult.Weapons, w)
				}
				// Add SubCult RuneSpells
				for _, rs := range cultModel.Cult.RuneSpells {
					parentCult.RuneSpells = append(parentCult.RuneSpells, rs)
				}
				// Add SubCult SpiritMagic
				for _, sm := range cultModel.Cult.SpiritMagic {
					parentCult.SpiritMagic = append(parentCult.SpiritMagic, sm)
				}

				// Set Details to ParentCult
				parentCult.Name = cultModel.Cult.Name
				parentCult.Description = cultModel.Cult.Description
				parentCult.Notes += "\n" + cultModel.Cult.Notes
				// Set cult to ParentCult
				c.Cult = parentCult
			}
		} else {
			c.Cult = cultModel.Cult
		}

		c.Cult.Rank = "Initiate"

		fmt.Println("CULT: " + c.Cult.Name)

		if c.Occupation.StandardOfLiving == "Free" {
			c.Abilities["Reputation"].OccupationValue = 5
		}

		if c.Occupation.StandardOfLiving == "Noble" {
			c.Abilities["Reputation"].OccupationValue = 10
		}

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

		// Update CreationSteps
		c.CreationSteps["Base Choices"] = true

		c.Role = "Player Character"

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

	wc := WebChar{
		CharacterModel: cm,
		Counter:        numToArray(5),
		SessionUser:    username,
		IsLoggedIn:     loggedIn,
		IsAdmin:        isAdmin,
		IsAuthor:       IsAuthor,
		Skills:         runequest.Skills,
		Passions:       runequest.PassionTypes,
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

		// Pull form values into cm.Character via gorilla/schema
		err = decoder.Decode(c, req.PostForm)
		if err != nil {
			panic(err)
		}

		str := req.FormValue("Lunars")
		lunars, err := strconv.Atoi(str)
		if err != nil {
			lunars = 0
		}
		if lunars > 0 {
			c.Equipment = append(c.Equipment, fmt.Sprintf("History: %d Lunars", lunars))
		}

		str = fmt.Sprintf("Reputation")
		rep, err := strconv.Atoi(req.FormValue(str))
		if err != nil {
			rep = 0
		}

		c.Abilities["Reputation"].CreationBonusValue = rep

		// Add Skills
		for i := 1; i < 4; i++ {

			sk := req.FormValue(fmt.Sprintf("Skill-%d-CoreString", i))

			if sk != "" {

				skbaseSkill := runequest.Skills[sk]
				fmt.Println(skbaseSkill)

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
				s1.CreationBonusValue = v

				if s1.UserChoice {
					userString := fmt.Sprintf("Skill-%d-UserString", i)
					s1.UserString = req.FormValue(userString)
				}

				targetString := createName(s1.CoreString, skbaseSkill.UserString)

				c.Skills[targetString] = s1
			}
		}

		// Add passions
		for i := 1; i < 6; i++ {

			coreString := req.FormValue(fmt.Sprintf("Passion-%d-CoreString", i))
			userString := req.FormValue(fmt.Sprintf("Passion-%d-UserString", i))

			if coreString != "" {

				targetString := createName(coreString, userString)

				if c.Abilities[targetString] == nil {
					// No ability
					a := &runequest.Ability{
						Type:       "Passion",
						CoreString: coreString,
						Updates:    []*runequest.Update{},
					}

					str := fmt.Sprintf("Passion-%d-Base", i)
					base, err := strconv.Atoi(req.FormValue(str))
					if err != nil {
						base = 0
					}

					if base > 0 {
						update := CreateUpdate("Base", base)
						a.Updates = append(a.Updates, update)
					}

					str = fmt.Sprintf("Passion-%d-Value", i)
					val, err := strconv.Atoi(req.FormValue(str))
					if err != nil {
						val = 0
					}

					if val > 0 {
						update := CreateUpdate("History", val)
						a.Updates = append(a.Updates, update)
					}

					if userString != "" {
						a.UserChoice = true
						a.UserString = userString
					}

					a.UpdateAbility()

					c.Abilities[targetString] = a
				} else {
					// Update existing ability

					update := CreateUpdate("History", 10)
					c.Abilities[targetString].Updates = append(c.Abilities[targetString].Updates, update)

					c.Abilities[targetString].UpdateAbility()
				}

			}
		}

		// Update CreationSteps
		c.CreationSteps["Personal History"] = true

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

		// Set Base Rune values for Power Runes
		baseRunes := []string{
			"Man", "Beast",
			"Fertility", "Death",
			"Harmony", "Disorder",
			"Truth", "Illusion",
			"Movement", "Stasis",
		}

		for _, v := range c.PowerRunes {
			if isInString(baseRunes, v.CoreString) {
				v.Base = 50
			} else {
				v.Base = 0
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

		// Update CreationSteps
		c.CreationSteps["Rune Affinities"] = true

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
		StringArray:    runequest.StatMap,
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

		for k, v := range cm.Character.Homeland.StatisticFrames {
			n, _ := strconv.Atoi(req.FormValue(k))
			if err != nil {
				n = 0
			}
			c.Statistics[k].Base = n
			c.Statistics[k].Max = v.Max
		}

		// Apply Rune Modifiers
		target1 := req.FormValue(fmt.Sprintf("RuneMod-%d", 0))
		c.Statistics[target1].RuneBonus = 2
		c.Statistics[target1].Max += 2

		target2 := req.FormValue(fmt.Sprintf("RuneMod-%d", 1))
		c.Statistics[target2].RuneBonus = 1
		c.Statistics[target2].Max++

		c.LocationForm = c.Homeland.LocationForm
		c.HitLocations = runequest.LocationForms[c.Homeland.LocationForm]
		c.HitLocationMap = runequest.SortLocations(c.HitLocations)

		c.SetAttributes()

		// Apply Move
		for _, m := range c.Homeland.Movement {
			c.Movement = append(c.Movement,
				&runequest.Movement{
					Name:  m.Name,
					Value: m.Value,
				})
		}

		// Update Character
		c.CurrentHP = c.Attributes["HP"].Max
		c.CurrentMP = c.Attributes["MP"].Max
		c.CurrentRP = c.Cult.NumRunePoints

		for _, v := range c.HitLocations {
			v.Value = v.Max
		}

		// Update CreationSteps
		c.CreationSteps["Roll Stats"] = true

		err = database.UpdateCharacterModel(db, cm)
		if err != nil {
			log.Panic(err)
		} else {
			fmt.Println("Saved")
		}

		// Add Flash messages

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

	flashes := session.Flashes("message")

	session.Save(req, w)

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
		Flashes:          flashes,
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

		// Add common skills based on statistics
		dodge := &runequest.Skill{
			CoreString: "Dodge",
			Category:   "Agility",
			Base:       c.Statistics["DEX"].Total * 2,
		}

		jump := &runequest.Skill{
			CoreString: "Jump",
			Category:   "Agility",
			Base:       c.Statistics["DEX"].Total * 3,
		}

		c.Skills["Dodge"] = dodge
		c.Skills["Jump"] = jump

		// Add choices to c.Homeland.Skills
		for i, sc := range c.Homeland.SkillChoices {
			for m := range sc.Skills {
				str := fmt.Sprintf("SC-%d-%d", i, m)
				if req.FormValue(str) != "" {
					c.Homeland.Skills = append(c.Homeland.Skills, &sc.Skills[m])
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

			targetString := createName(s.CoreString, s.UserString)

			if c.Skills[targetString] != nil {
				// Skill exists in Character, modify it via pointer
				c.Skills[targetString].HomelandValue = s.HomelandValue

				if s.Base != c.Skills[targetString].Base {
					c.Skills[targetString].Base = s.Base
				}

				fmt.Println("Updated Character Skill: " + c.Skills[targetString].Name)

			} else {
				// We need to find the base skill from the Master list or create it
				bs, ok := runequest.Skills[s.CoreString]
				if !ok {
					fmt.Println("Skill is new: " + targetString)

					// New Skill
					baseSkill := &runequest.Skill{
						CoreString:    s.CoreString,
						UserString:    s.UserString,
						Category:      s.Category,
						Base:          s.Base,
						UserChoice:    s.UserChoice,
						HomelandValue: s.HomelandValue,
					}

					// Update our new skill

					if baseSkill.Category == "" {
						baseSkill.Category = "Agility"
					}

					sc := c.SkillCategories[baseSkill.Category]

					baseSkill.CategoryValue = sc.Value

					baseSkill.Name = createName(baseSkill.CoreString, baseSkill.UserString)
					baseSkill.UpdateSkill()

					fmt.Println("Add Skill to character: " + baseSkill.Name)
					c.Skills[baseSkill.Name] = baseSkill
				} else {
					// Skill exists in master list
					fmt.Println("Skill in master list: " + targetString)

					fmt.Println(bs)

					baseSkill := &runequest.Skill{
						CoreString: bs.CoreString,
						Category:   bs.Category,
						Base:       s.Base,
						UserChoice: bs.UserChoice,
					}

					baseSkill.HomelandValue = s.HomelandValue
					fmt.Println(s.HomelandValue, baseSkill.HomelandValue)

					if s.Base > baseSkill.Base {
						baseSkill.Base = s.Base
					}
					if s.UserString != "" {
						baseSkill.UserString = s.UserString
					}

					baseSkill.Name = createName(baseSkill.CoreString, baseSkill.UserString)

					// Add Skill to Character
					fmt.Println("Add Skill to character: " + baseSkill.Name)
					c.Skills[baseSkill.Name] = baseSkill
					c.Skills[baseSkill.Name].UpdateSkill()
				}
			}
		}

		// Homelands grant 3 base passions

		for _, a := range c.Homeland.PassionList {
			c.Homeland.Passions = append(c.Homeland.Passions, a)

			targetString := createName(a.CoreString, a.UserString)

			if c.Abilities[targetString] == nil {
				// Need to create a new ability
				a := &runequest.Ability{
					Type:       "Passion",
					CoreString: a.CoreString,
					UserChoice: a.UserChoice,
					UserString: a.UserString,
					Updates:    []*runequest.Update{},
				}

				update := CreateUpdate("Homeland", 60)
				a.Updates = append(a.Updates, update)

				a.UpdateAbility()
				c.Abilities[targetString] = a
			} else {
				// Update existing ability

				update := CreateUpdate("Homeland", 10)
				c.Abilities[targetString].Updates = append(c.Abilities[targetString].Updates, update)

				c.Abilities[targetString].UpdateAbility()
			}
		}

		// Homeland grants a bonus to a rune affinity
		if isInString(runequest.ElementalRuneOrder, c.Homeland.RuneBonus) {
			c.ElementalRunes[c.Homeland.RuneBonus].HomelandValue += 10
		}

		if isInString(runequest.PowerRuneOrder, c.Homeland.RuneBonus) {
			c.PowerRunes[c.Homeland.RuneBonus].HomelandValue += 10
		}

		// Update CreationSteps
		c.CreationSteps["Apply Homeland"] = true

		err = database.UpdateCharacterModel(db, cm)
		if err != nil {
			log.Panic(err)
		} else {
			fmt.Println("Saved")
		}

		// Add Flash messages

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

	flashes := session.Flashes("message")

	session.Save(req, w)

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
		Flashes:          flashes,
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
					c.Occupation.Skills = append(c.Occupation.Skills, &sc.Skills[m])
				}
			}
		}

		// Add choices for Weapon skills
		for i, w := range c.Occupation.Weapons {
			str := fmt.Sprintf("Weapon-%d-CoreString", i)
			fv := req.FormValue(str)
			if fv != "" {

				bs := runequest.Skills[fv]

				ws := &runequest.Skill{
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

			targetString := createName(s.CoreString, s.UserString)

			if c.Skills[targetString] != nil {
				// Skill exists in Character, modify it via pointer
				c.Skills[targetString].OccupationValue = s.OccupationValue

				if s.Base > c.Skills[targetString].Base {
					c.Skills[targetString].Base = s.Base
				}

				fmt.Println("Updated Character Skill: " + c.Skills[targetString].Name)

			} else {
				// We need to find the base skill from the Master list or create it
				bs, ok := runequest.Skills[s.CoreString]
				if !ok {
					fmt.Println("Skill is new: " + targetString)

					// New Skill
					baseSkill := &runequest.Skill{
						CoreString:      s.CoreString,
						UserString:      s.UserString,
						Category:        s.Category,
						Base:            s.Base,
						UserChoice:      s.UserChoice,
						OccupationValue: s.OccupationValue,
					}

					// Update our new skill
					sc := c.SkillCategories[baseSkill.Category]

					baseSkill.CategoryValue = sc.Value

					baseSkill.Name = createName(baseSkill.CoreString, baseSkill.UserString)
					baseSkill.UpdateSkill()

					fmt.Println("Add Skill to character: " + baseSkill.Name)
					c.Skills[baseSkill.Name] = baseSkill
				} else {
					// Skill exists in master list
					fmt.Println("Skill in master list: " + targetString)

					fmt.Println(bs)

					baseSkill := &runequest.Skill{
						CoreString: bs.CoreString,
						Category:   bs.Category,
						Base:       bs.Base,
						UserChoice: bs.UserChoice,
					}

					baseSkill.OccupationValue = s.OccupationValue
					fmt.Println(s.OccupationValue, baseSkill.OccupationValue)

					if s.Base > baseSkill.Base {
						baseSkill.Base = s.Base
					}
					if s.UserString != "" {
						baseSkill.UserString = s.UserString
					}

					baseSkill.Name = createName(baseSkill.CoreString, baseSkill.UserString)

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

			a := c.Occupation.PassionList[n]

			targetString := createName(a.CoreString, a.UserString)

			if c.Abilities[targetString] == nil {
				// Need to create a new ability
				a := &runequest.Ability{
					Type:       "Passion",
					CoreString: a.CoreString,
					UserChoice: a.UserChoice,
					UserString: a.UserString,
					Updates:    []*runequest.Update{},
				}

				update := CreateUpdate("Occupation", 60)
				a.Updates = append(a.Updates, update)

				a.UpdateAbility()
				c.Abilities[targetString] = a
			} else {
				// Update existing ability

				update := CreateUpdate("Occupation", 10)
				c.Abilities[targetString].Updates = append(c.Abilities[targetString].Updates, update)

				c.Abilities[targetString].UpdateAbility()
			}

		}

		// Equipment

		if len(c.Occupation.Equipment) > 0 {
			for _, e := range c.Occupation.Equipment {
				c.Equipment = append(c.Equipment, e)
			}
		}

		// Update CreationSteps
		c.CreationSteps["Apply Occupation"] = true

		err = database.UpdateCharacterModel(db, cm)
		if err != nil {
			log.Panic(err)
		} else {
			fmt.Println("Saved")
		}

		// Add Flash messages

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

	flashes := session.Flashes("message")

	session.Save(req, w)

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

	numRunePoints := numToArray(5)
	numSpiritMagic := numToArray(5)

	// Add Associated Cult spells
	totalSpiritMagic := []*runequest.Spell{}

	totalSpiritMagic = c.Cult.SpiritMagic

	for _, ac := range c.Cult.AssociatedCults {
		for _, spell := range ac.SpiritMagic {
			totalSpiritMagic = append(totalSpiritMagic, spell)
		}
	}

	wc := WebChar{
		CharacterModel:   cm,
		SessionUser:      username,
		IsLoggedIn:       loggedIn,
		IsAdmin:          isAdmin,
		IsAuthor:         IsAuthor,
		Skills:           runequest.Skills,
		NumRunePoints:    numRunePoints,
		NumSpiritMagic:   numSpiritMagic,
		TotalSpiritMagic: totalSpiritMagic,
		WeaponCategories: runequest.WeaponCategories,
		Flashes:          flashes,
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

		rPS, err := strconv.Atoi(req.FormValue("RunePoints"))
		if err != nil {
			rPS = 3
		}

		c.Cult.NumRunePoints = rPS

		if c.Cult.NumRunePoints > 3 {

			update := CreateUpdate("POW pledge to cult", 3-c.Cult.NumRunePoints)

			c.Statistics["POW"].Updates = append(c.Statistics["POW"].Updates, update)
		}

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

				s := baseSpell
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

				baseSpell.Cost = cost

				if spec != "" {
					baseSpell.UserString = spec
				}

				baseSpell.GenerateName()
				c.SpiritMagic[baseSpell.Name] = baseSpell
			}
		}

		// Add choices to c.Cult.Skills
		for i, sc := range c.Cult.SkillChoices {
			for m := range sc.Skills {
				str := fmt.Sprintf("SC-%d-%d", i, m)
				if req.FormValue(str) != "" {
					c.Cult.Skills = append(c.Cult.Skills, &sc.Skills[m])
				}
			}
		}

		// Add choices for Weapon skills
		for i, w := range c.Cult.Weapons {
			str := fmt.Sprintf("Weapon-%d-CoreString", i)
			fv := req.FormValue(str)

			if fv != "" {
				bs := runequest.Skills[fv]

				ws := &runequest.Skill{
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

			targetString := createName(s.CoreString, s.UserString)

			if c.Skills[targetString] != nil {
				// Skill exists in Character, modify it via pointer
				c.Skills[targetString].CultValue = s.CultValue

				if s.Base > c.Skills[targetString].Base {
					c.Skills[targetString].Base = s.Base
				}

				fmt.Println("Updated Character Skill: " + c.Skills[targetString].Name)

			} else {
				// We need to find the base skill from the Master list or create it
				bs, ok := runequest.Skills[s.CoreString]
				if !ok {
					fmt.Println("Skill is new: " + targetString)

					// New Skill
					baseSkill := &runequest.Skill{
						CoreString: s.CoreString,
						UserString: s.UserString,
						Category:   s.Category,
						Base:       s.Base,
						UserChoice: s.UserChoice,
						CultValue:  s.CultValue,
					}

					// Update our new skill
					sc := c.SkillCategories[baseSkill.Category]

					baseSkill.CategoryValue = sc.Value

					baseSkill.Name = createName(baseSkill.CoreString, baseSkill.UserString)
					baseSkill.UpdateSkill()

					fmt.Println("Add Skill to character: " + baseSkill.Name)
					c.Skills[baseSkill.Name] = baseSkill
				} else {
					// Skill exists in master list
					fmt.Println("Skill in master list: " + targetString)

					fmt.Println(bs)

					baseSkill := &runequest.Skill{
						CoreString: bs.CoreString,
						Category:   bs.Category,
						Base:       bs.Base,
						UserChoice: bs.UserChoice,
					}

					baseSkill.CultValue = s.CultValue
					fmt.Println(s.CultValue, baseSkill.CultValue)

					if s.Base > baseSkill.Base {
						baseSkill.Base = s.Base
					}
					if s.UserString != "" {
						baseSkill.UserString = s.UserString
					}

					baseSkill.Name = createName(baseSkill.CoreString, baseSkill.UserString)

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
			a := c.Cult.PassionList[n]

			targetString := createName(a.CoreString, a.UserString)

			if c.Abilities[targetString] == nil {
				// Need to create a new ability
				a := &runequest.Ability{
					Type:       "Passion",
					CoreString: a.CoreString,
					UserChoice: a.UserChoice,
					UserString: a.UserString,
					Updates:    []*runequest.Update{},
				}

				update := CreateUpdate("Cult", 60)
				a.Updates = append(a.Updates, update)

				a.UpdateAbility()
				c.Abilities[targetString] = a
			} else {
				// Update existing ability

				update := CreateUpdate("Cult", 10)
				c.Abilities[targetString].Updates = append(c.Abilities[targetString].Updates, update)

				c.Abilities[targetString].UpdateAbility()
			}
		}

		// Add Spirit Magic
		c.Cult.NumSpiritMagic = 5

		if c.Occupation.Name == "Assistant Shaman" {
			c.Cult.NumSpiritMagic += 5
		}

		// Set Rune Points to Max
		c.CurrentRP = c.Cult.NumRunePoints

		// Update CreationSteps
		c.CreationSteps["Apply Cult"] = true

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
		BigCounter:       numToArray(5),
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

		if req.FormValue("Open") != "" {
			cm.Open = true
		} else {
			cm.Open = false
		}

		c.Clan = req.FormValue("Clan")
		c.Tribe = req.FormValue("Tribe")
		c.Age = 21

		c.StandardofLiving = c.Occupation.StandardOfLiving
		c.Ransom = c.Occupation.Ransom

		// Do Stuff
		// 25% additions
		for i := 1; i < 5; i++ {
			coreString := req.FormValue(fmt.Sprintf("Skill-25-%d", i))
			userString := req.FormValue(fmt.Sprintf("Skill-25-%d-UserString", i))

			targetString := createName(coreString, userString)

			fmt.Println(targetString)

			// Skill exists in Character, modify it via pointer
			if targetString != "" {

				s, ok := c.Skills[targetString]
				if !ok {
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

					sk.Name = createName(sk.CoreString, sk.UserString)

					// Add Skill to Character
					fmt.Println("Add Skill to character: " + sk.Name)
					c.Skills[sk.Name] = sk
					c.Skills[sk.Name].UpdateSkill()
				}

				// Potential error here - Insight (Human)
				s = c.Skills[targetString]

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
		for i := 1; i < 6; i++ {
			coreString := req.FormValue(fmt.Sprintf("Skill-10-%d", i))
			userString := req.FormValue(fmt.Sprintf("Skill-10-%d-UserString", i))

			targetString := createName(coreString, userString)

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

					sk.Name = createName(sk.CoreString, sk.UserString)

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

		// Update CreationSteps
		c.CreationSteps["Personal Skills"] = true

		err = database.UpdateCharacterModel(db, cm)
		if err != nil {
			log.Panic(err)
		} else {
			fmt.Println("Saved")
		}

		url := fmt.Sprintf("/cc8_finishing_touches/%d", cm.ID)

		http.Redirect(w, req, url, http.StatusSeeOther)
	}
}

// FinishingTouchesHandler assigns attacks and armor
func FinishingTouchesHandler(w http.ResponseWriter, req *http.Request) {

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

	baseWeapons := runequest.BaseWeapons

	wc := WebChar{
		CharacterModel: cm,
		SessionUser:    username,
		IsLoggedIn:     loggedIn,
		IsAdmin:        isAdmin,
		IsAuthor:       IsAuthor,
		Counter:        numToArray(4),
		BigCounter:     numToArray(5),
		Skills:         runequest.Skills,
		BaseWeapons:    baseWeapons,
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/cc8_finishing_touches.html", wc)

	}

	if req.Method == "POST" {

		err := req.ParseMultipartForm(MaxMemory)
		if err != nil {
			panic(err)
		}

		if c.MeleeAttacks == nil {
			c.MeleeAttacks = map[string]*runequest.Attack{}
		}

		if c.RangedAttacks == nil {
			c.RangedAttacks = map[string]*runequest.Attack{}
		}

		// Do Stuff
		// Melee Weapons & Attacks
		for i := 1; i < 5; i++ {
			weaponString := req.FormValue(fmt.Sprintf("Melee-Weapon-%d", i))
			skillString := req.FormValue(fmt.Sprintf("Melee-Skill-%d", i))

			if weaponString != "" && skillString != "" {

				weaponIndex, err := strconv.Atoi(weaponString)
				if err != nil {
					weaponIndex = 0
				}

				weapon := baseWeapons[weaponIndex]

				c.MeleeAttacks[weapon.Name] = &runequest.Attack{
					Name:         weapon.Name,
					Skill:        c.Skills[skillString],
					DamageString: weapon.Damage + c.Attributes["DB"].Text,
					StrikeRank:   c.Attributes["DEXSR"].Base + c.Attributes["SIZSR"].Base + weapon.SR,
					Weapon:       weapon,
				}
			}
		}

		// Ranged Weapons & Attacks
		for i := 1; i < 4; i++ {
			weaponString := req.FormValue(fmt.Sprintf("Ranged-Weapon-%d", i))
			skillString := req.FormValue(fmt.Sprintf("Ranged-Skill-%d", i))

			if weaponString != "" && skillString != "" {

				weaponIndex, err := strconv.Atoi(weaponString)
				if err != nil {
					weaponIndex = 0
				}

				weapon := baseWeapons[weaponIndex]

				// Ranged weapon
				c.RangedAttacks[weapon.Name] = &runequest.Attack{
					Name:         weapon.Name,
					Skill:        c.Skills[skillString],
					DamageString: weapon.Damage,
					StrikeRank:   c.Attributes["DEXSR"].Base,
					Weapon:       weapon,
				}
			}
		}

		// Armor
		for k, v := range c.HitLocations {
			str := req.FormValue(fmt.Sprintf("%s-Armor", k))
			armor, err := strconv.Atoi(str)
			if err != nil {
				armor = v.Armor
			}
			v.Armor = armor
		}

		// Update CreationSteps
		c.CreationSteps["Finishing Touches"] = true
		c.CreationSteps["Complete"] = true

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
