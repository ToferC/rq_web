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

		cm.Open = true
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

		c.AddRuneModifiers()

		c.SetAttributes()

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
		for _, s := range c.Homeland.Skills {
			// Modify Skill
			if s.UserString == "any" {
				// User Chooses a new specialization
				str := req.FormValue(fmt.Sprintf("%s-UserString", s.CoreString))
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

					baseSkill := runequest.Skills[s.CoreString]
					fmt.Println(baseSkill)
					baseSkill.HomelandValue = s.HomelandValue
					fmt.Println(s.HomelandValue, baseSkill.HomelandValue)

					if s.Base > baseSkill.Base {
						baseSkill.Base = s.Base
					}
					if s.UserString != "" {
						baseSkill.UserString = s.UserString
					}
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

	cults, err := database.ListCultModels(db)
	if err != nil {
		panic(err)
	}

	wc := WebChar{
		CharacterModel:   cm,
		CultModels:       cults,
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
			bs := runequest.Skills[fv]

			ws := runequest.Skill{
				CoreString:      bs.CoreString,
				Category:        bs.Category,
				Base:            bs.Base,
				OccupationValue: w.Value,
			}
			c.Occupation.Skills = append(c.Occupation.Skills, ws)
		}

		// Add Skills to Character
		for _, s := range c.Occupation.Skills {
			// Modify Skill
			if s.UserString == "any" {
				// User Chooses a new specialization
				str := req.FormValue(fmt.Sprintf("%s-UserString", s.CoreString))
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

					baseSkill := runequest.Skills[s.CoreString]
					fmt.Println(baseSkill)
					baseSkill.OccupationValue = s.OccupationValue
					fmt.Println(s.OccupationValue, baseSkill.OccupationValue)

					if s.Base > baseSkill.Base {
						baseSkill.Base = s.Base
					}
					if s.UserString != "" {
						baseSkill.UserString = s.UserString
					}
					// Add Skill to Character
					fmt.Println("Add Skill to character: " + baseSkill.Name)
					c.Skills[baseSkill.Name] = baseSkill
					c.Skills[baseSkill.Name].UpdateSkill()
				}
			}
		}

		// Occupations grant a bonus to one Passion
		str := req.FormValue("Passion")
		n, err := strconv.Atoi(str)
		if err != nil {
			n = 0
		}

		c.ModifyAbility(c.Occupation.PassionList[n])

		// Equipment

		for _, e := range c.Occupation.Equipment {
			c.Equipment = append(c.Equipment, e)
		}

		// Set Cult
		cStr := req.FormValue("Cult")

		cID, err := strconv.Atoi(cStr)
		if err != nil {
			cID = 0
		}

		cultModel, err := database.PKLoadCultModel(db, int64(cID))
		if err != nil {
			fmt.Println("No Cult Found")
		}

		c.Cult = cultModel.Cult

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

	wc := WebChar{
		CharacterModel: cm,
		SessionUser:    username,
		IsLoggedIn:     loggedIn,
		IsAdmin:        isAdmin,
		IsAuthor:       IsAuthor,
		Skills:         runequest.Skills,
		NumRunePoints:  numRunePoints,
		NumSpiritMagic: numSpiritMagic,
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

		c.RuneSpells = map[string]*runequest.Spell{}
		c.SpiritMagic = map[string]*runequest.Spell{}

		// Rune Magic
		for i := 1; i < 4; i++ {
			str := req.FormValue(fmt.Sprintf("RuneSpell-%d", i))
			spec := req.FormValue(fmt.Sprintf("RuneSpell-%d-UserString", i))
			if str != "" {
				index, err := strconv.Atoi(str)
				if err != nil {
					index = 0
					fmt.Println("Spell Not found")
				}
				baseSpell := runequest.RuneSpells[index]

				s := &baseSpell
				if spec != "" {
					s.UserString = spec
				}
				s.GenerateName()
				c.RuneSpells[s.Name] = s
			}
		}

		// Spirit Magic
		for i := 1; i < 4; i++ {
			str := req.FormValue(fmt.Sprintf("SpiritMagic-%d", i))
			spec := req.FormValue(fmt.Sprintf("SpiritMagic-%d-UserString", i))
			pString := req.FormValue(fmt.Sprintf("SpiritMagic-%d-Points", i))

			if str != "" {

				index, err := strconv.Atoi(str)
				if err != nil {
					index = 0
					fmt.Println("Spell Not found")
				}

				pts, err := strconv.Atoi(pString)
				if err != nil {
					pts = 1
					fmt.Println("Non-number entered")
				}

				baseSpell := runequest.RuneSpells[index]

				s := &baseSpell
				s.Cost = pts

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
		for _, s := range c.Cult.Skills {
			// Modify Skill
			if s.UserString == "any" {
				// User Chooses a new specialization
				str := req.FormValue(fmt.Sprintf("%s-UserString", s.CoreString))
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

					baseSkill := runequest.Skills[s.CoreString]
					fmt.Println(baseSkill)
					baseSkill.CultValue = s.CultValue
					fmt.Println(s.CultValue, baseSkill.CultValue)

					if s.Base > baseSkill.Base {
						baseSkill.Base = s.Base
					}
					if s.UserString != "" {
						baseSkill.UserString = s.UserString
					}
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
			targetString := req.FormValue(fmt.Sprintf("Skill-25-%d", i))
			userString := req.FormValue(fmt.Sprintf("Skill-25-%d-UserString", i))

			fmt.Println(targetString)

			// Skill exists in Character, modify it via pointer
			if targetString != "" {
				// Determine if skill already exists in c.Skills

				s := c.Skills[targetString]

				if s == nil {
					s = runequest.Skills[targetString]
					if userString != "" {
						s.UserString = userString
					}

					s.GenerateName()

					// Add Skill to Character
					fmt.Println("Add Skill to character: " + s.Name)
					c.Skills[s.Name] = s
					c.Skills[s.Name].UpdateSkill()
				}

				t := time.Now()
				tString := t.Format("2006-01-02 15:04:05")

				update := &runequest.Update{
					Date:  tString,
					Event: "Personal Skills 25%",
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
			targetString := req.FormValue(fmt.Sprintf("Skill-10-%d", i))
			userString := req.FormValue(fmt.Sprintf("Skill-10-%d-UserString", i))

			// Skill exists in Character, modify it via pointer
			if targetString != "" {

				s := c.Skills[targetString]

				if s == nil {
					s = runequest.Skills[targetString]
					if userString != "" {
						s.UserString = userString
					}

					s.GenerateName()

					// Add Skill to Character
					fmt.Println("Add Skill to character: " + s.Name)
					c.Skills[s.Name] = s
					c.Skills[s.Name].UpdateSkill()
				}

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
