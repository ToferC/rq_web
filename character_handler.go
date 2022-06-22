package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/thewhitetulip/Tasks/sessions"
	"github.com/toferc/rq_web/database"
	"github.com/toferc/rq_web/models"
	"github.com/toferc/runequest"
)

// CharacterHandler renders a character in a Web page
func CharacterHandler(w http.ResponseWriter, req *http.Request) {

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

	flashes := session.Flashes("message")

	session.Save(req, w)

	vars := mux.Vars(req)
	pk := vars["id"]

	if len(pk) == 0 {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	id, err := strconv.Atoi(pk)
	if err != nil {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	cm, err := database.PKLoadCharacterModel(db, int64(id))
	if err != nil {
		fmt.Println("Unable to load CharacterModel")
		http.Redirect(w, req, "/notfound", http.StatusSeeOther)
		return
	}

	statBlock := cm.Character.StatBlock()

	IsAuthor := false

	if username == cm.Author.UserName {
		IsAuthor = true
	}

	c := cm.Character

	if IsAuthor {
		// Create 4 empty equipment slots for new stuff
		for i := 0; i < 4; i++ {
			c.Equipment = append(c.Equipment, "")
		}
	}

	if cm.Image == nil {
		cm.Image = new(models.Image)
		cm.Image.Path = DefaultCharacterPortrait
	}

	//c.DetermineSkillCategoryValues()

	wc := WebChar{
		CharacterModel: cm,
		IsAuthor:       IsAuthor,
		StatBlock:      statBlock,
		IsLoggedIn:     loggedIn,
		SessionUser:    username,
		IsAdmin:        isAdmin,
		Counter:        numToArray(10),
		Flashes:        flashes,
		StringArray:    runequest.StatMap,
		CategoryOrder:  runequest.CategoryOrder,
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/view_character.html", wc)

	}

	if req.Method == "POST" {

		// Parse Form and redirect
		err := req.ParseForm()
		if err != nil {
			panic(err)
		}

		// Update MP
		str := req.FormValue("MP")
		mp, err := strconv.Atoi(str)
		if err != nil {
			mp = c.CurrentMP
		}

		if mp < 0 {
			mp = 0
		}

		if mp > c.Attributes["MP"].Max {
			mp = c.Attributes["MP"].Max
		}

		c.CurrentMP = mp

		// Update RP
		str = req.FormValue("RP")
		rp, err := strconv.Atoi(str)
		if err != nil {
			rp = c.CurrentMP
		}

		if rp < 0 {
			rp = 0
		}

		c.CurrentRP = rp

		// Update ExtraCults RPs
		for _, ec := range c.ExtraCults {
			str = req.FormValue("RP-" + ec.Name)
			rp, err := strconv.Atoi(str)
			if err != nil {
				rp = ec.CurrentRunePoints
			}

			if rp < 0 {
				rp = 0
			}

			ec.CurrentRunePoints = rp
		}

		// Check skill xp
		for _, v := range c.Skills {
			if req.FormValue(fmt.Sprintf("%s-XP", v.Name)) != "" {
				v.ExperienceCheck = true
			} else {
				v.ExperienceCheck = false
			}
		}

		// Check Elemental Rune xp
		for _, v := range c.ElementalRunes {
			if req.FormValue(fmt.Sprintf("%s-XP", v.Name)) != "" {
				v.ExperienceCheck = true
			} else {
				v.ExperienceCheck = false
			}
		}

		// Check Power Rune XP
		for _, v := range c.PowerRunes {
			if req.FormValue(fmt.Sprintf("%s-XP", v.Name)) != "" {
				v.ExperienceCheck = true
			} else {
				v.ExperienceCheck = false
			}
		}

		// Check Passion xp
		for _, v := range c.Abilities {
			if req.FormValue(fmt.Sprintf("%s-XP", v.Name)) != "" {
				v.ExperienceCheck = true
			} else {
				v.ExperienceCheck = false
			}
		}

		// Check POW xp
		if req.FormValue("POW-XP") != "" {
			c.Statistics["POW"].ExperienceCheck = true
		} else {
			c.Statistics["POW"].ExperienceCheck = false
		}

		// Update HitLocations
		totalDamage := 0

		for k, v := range c.HitLocations {
			str := req.FormValue(fmt.Sprintf("%s-HP", k))
			hp, err := strconv.Atoi(str)
			if err != nil {
				hp = v.Value
			}

			if hp > v.Max {
				hp = v.Max
			}

			if hp < v.Min {
				hp = v.Min
			}

			if hp < 1 {
				v.Disabled = true
			} else {
				v.Disabled = false
			}

			v.Value = hp
			totalDamage += v.Max - hp
		}

		// Determine total damage based on HitLocation HP
		c.CurrentHP = c.Attributes["HP"].Max - totalDamage

		// Read Equipment
		var equipment = []string{}

		for i := 0; i < len(c.Equipment)+1; i++ {
			str := req.FormValue(fmt.Sprintf("Equipment-%d", i))
			if str != "" {
				equipment = append(equipment, str)
			}
		}

		c.Equipment = equipment

		// Track Weapon HP
		for k, v := range c.MeleeAttacks {
			hpString := req.FormValue(fmt.Sprintf("%s-HP", k))
			hp, err := strconv.Atoi(hpString)
			if err != nil {
				hp = v.Weapon.HP
			}

			if hp > v.Weapon.HP {
				hp = v.Weapon.HP
			}

			if hp < -(2 * v.Weapon.HP) {
				hp = -v.Weapon.HP * 2
			}

			v.Weapon.CurrentHP = hp
		}

		// Track Weapon HP
		for k, v := range c.RangedAttacks {
			hpString := req.FormValue(fmt.Sprintf("%s-HP", k))
			hp, err := strconv.Atoi(hpString)
			if err != nil {
				hp = v.Weapon.HP
			}

			if hp > v.Weapon.HP {
				hp = v.Weapon.HP
			}

			if hp < -(2 * v.Weapon.HP) {
				hp = -v.Weapon.HP * 2
			}

			v.Weapon.CurrentHP = hp
		}

		// Track Boundspirit MPs
		for i, bs := range c.BoundSpirits {
			mpString := req.FormValue(fmt.Sprintf("BS-%d-MP", i))
			mp, err := strconv.Atoi(mpString)
			if err != nil {
				mp = bs.Pow
			}
			bs.CurrentMP = mp
		}

		err = database.UpdateCharacterModel(db, cm)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("Saved")
		}

		url := fmt.Sprintf("/view_character/%d#gameplay", cm.ID)

		http.Redirect(w, req, url, http.StatusSeeOther)
		return
	}

}

// ModifyCharacterHandler renders a character in a Web page
func ModifyCharacterHandler(w http.ResponseWriter, req *http.Request) {

	// Get session values or redirect to Login
	session, err := sessions.Store.Get(req, "session")

	if err != nil {
		log.Println("error identifying session")
		http.Redirect(w, req, "/login/", http.StatusFound)
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
		return
	}

	id, err := strconv.Atoi(pk)
	if err != nil {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	// Load CharacterModel
	cm, err := database.PKLoadCharacterModel(db, int64(id))
	if err != nil {
		fmt.Println(err)
	}

	// Validate that User == Author
	IsAuthor := false

	if username == cm.Author.UserName || isAdmin == "true" {
		IsAuthor = true
	} else {
		http.Redirect(w, req, "/", http.StatusFound)
		return
	}

	c := cm.Character

	// Add empty Movement's to character if needed
	if len(c.Movement) == 0 {
		c.Movement = []*runequest.Movement{
			&runequest.Movement{
				Name:  "Ground",
				Value: 8,
			},
		}
	}

	mvLen := len(c.Movement)

	if len(c.Movement) < 3 {
		// Add Movement
		for i := mvLen; i < mvLen+2; i++ {

			c.Movement = append(c.Movement,
				&runequest.Movement{
					Name:  "",
					Value: 0,
				})
		}
	}

	if cm.Image == nil {
		cm.Image = new(models.Image)
		cm.Image.Path = DefaultCharacterPortrait
	}

	wc := WebChar{
		CharacterModel:    cm,
		SessionUser:       username,
		IsAuthor:          IsAuthor,
		IsLoggedIn:        loggedIn,
		IsAdmin:           isAdmin,
		CategoryOrder:     runequest.CategoryOrder,
		StringArray:       runequest.StatMap,
		StandardsOfLiving: runequest.Standards,
	}

	if req.Method == "GET" {

		// Render page

		Render(w, "templates/modify_character.html", wc)

	}

	if req.Method == "POST" {

		err := req.ParseMultipartForm(MaxMemory)
		if err != nil {
			panic(err)
		}

		c.Name = req.FormValue("Name")

		c.Description = req.FormValue("Description")

		eventString := req.FormValue("Event")

		c.Cult.Rank = req.FormValue("Rank")

		rp, err := strconv.Atoi(req.FormValue("RunePoints"))
		if err != nil {
			rp = 3
			fmt.Println("Not a number")
		}

		c.Cult.NumRunePoints = rp

		// Details
		c.StandardofLiving = req.FormValue("Standard")

		income, err := strconv.Atoi(req.FormValue("Income"))
		if err != nil {
			income = 0
		}
		c.Income = income

		ransom, err := strconv.Atoi(req.FormValue("Ransom"))
		if err != nil {
			ransom = 0
		}
		c.Ransom = ransom

		// Update statistics
		for _, st := range runequest.StatMap {

			stat := c.Statistics[st]

			mod, _ := strconv.Atoi(req.FormValue(st))

			if mod != stat.Total {

				modVal := mod - stat.Total

				update := CreateUpdate(eventString, modVal)

				if stat.Updates == nil {
					stat.Updates = []*runequest.Update{}
				}

				stat.Updates = append(stat.Updates, update)

				stat.ExperienceCheck = false
				stat.Max += modVal
			}
			stat.UpdateStatistic()
		}

		// Update Character
		c.SetAttributes()

		// Update Character
		c.CurrentHP = c.Attributes["HP"].Max
		c.CurrentMP = c.Attributes["MP"].Max
		c.CurrentRP = c.Cult.NumRunePoints

		// Hit locations
		for k, v := range c.HitLocations {
			str := req.FormValue(fmt.Sprintf("%s-Armor", k))
			armor, err := strconv.Atoi(str)
			if err != nil {
				armor = v.Armor
			}
			v.Armor = armor
			v.Value = v.Max
		}

		tempMovement := []*runequest.Movement{}

		// Add Movement
		for i := 1; i < mvLen+2; i++ {
			moveName := req.FormValue(fmt.Sprintf("Move-Name-%d", i))

			if moveName != "" {
				mv, err := strconv.Atoi(req.FormValue(fmt.Sprintf("Move-Value-%d", i)))
				if err != nil {
					mv = 8
				}

				tempMovement = append(tempMovement,
					&runequest.Movement{
						Name:  ProcessUserString(moveName),
						Value: mv,
					})
			}
		}

		c.Movement = tempMovement

		for k, s := range c.Skills {
			mod, _ := strconv.Atoi(req.FormValue(s.Name))

			if mod != s.Total {

				modVal := mod - s.Total

				update := CreateUpdate(eventString, modVal)

				if s.Updates == nil {
					s.Updates = []*runequest.Update{}
				}

				s.Updates = append(s.Updates, update)

				s.ExperienceCheck = false
			}
			if s.UserString != "" {
				us := req.FormValue(fmt.Sprintf("%s-UserString", s.Name))
				s.UserString = ProcessUserString(us)
			}
			s.UpdateSkill()

			// Update Weapons Skills
			for _, v := range c.MeleeAttacks {
				if v.Skill.Name == s.Name {
					v.Skill = s
				}
			}

			for _, v := range c.RangedAttacks {
				if v.Skill.Name == s.Name {
					v.Skill = s
				}
			}

			// Remove Character skill is zeroed out
			if s.Total < 1 {
				delete(c.Skills, k)
			}
		}

		// Update Elemental Runes
		for k, v := range c.ElementalRunes {
			mod, err := strconv.Atoi(req.FormValue(k))
			if err != nil {
				mod = v.Total
			}

			if mod != v.Total {

				modVal := mod - v.Total

				update := CreateUpdate(eventString, modVal)

				if v.Updates == nil {
					v.Updates = []*runequest.Update{}
				}

				v.Updates = append(v.Updates, update)
				v.ExperienceCheck = false

			}
			v.UpdateAbility()
		}

		// Create array for watching which updates made to opposed Runes
		triggered := []string{}

		// Update Power Runes
		for k, v := range c.PowerRunes {

			if !isInString(triggered, k) {

				mod, err := strconv.Atoi(req.FormValue(k))
				if err != nil {
					mod = v.Total
				}

				if mod != 0 {

					// Can't have Power rune > 99
					if mod > 99 {
						mod = 99
					}

					if mod != v.Total {

						modVal := mod - v.Total

						update := CreateUpdate(eventString, modVal)

						if v.Updates == nil {
							v.Updates = []*runequest.Update{}
						}

						v.Updates = append(v.Updates, update)

						v.UpdateAbility()
						v.ExperienceCheck = false

						if v.OpposedAbility != "" {

							opposed := c.PowerRunes[v.OpposedAbility]

							// Update opposed Power Rune if needed
							if v.Total+opposed.Total > 99 {

								opposedUpdate := CreateUpdate(eventString, -modVal)

								if opposed.Updates == nil {
									opposed.Updates = []*runequest.Update{}
								}
								opposed.Updates = append(opposed.Updates, opposedUpdate)
								opposed.UpdateAbility()
								triggered = append(triggered, opposed.Name)
							}
						}
					}
				}
			}
		}

		// Update Condition Runes
		for k, v := range c.ConditionRunes {
			mod, err := strconv.Atoi(req.FormValue(k))
			if err != nil {
				mod = v.Total
			}

			if mod != v.Total {

				modVal := mod - v.Total

				update := CreateUpdate(eventString, modVal)

				if v.Updates == nil {
					v.Updates = []*runequest.Update{}
				}

				v.Updates = append(v.Updates, update)
				v.ExperienceCheck = false

			}
			v.UpdateAbility()
		}

		// Update Passions
		for k, a := range c.Abilities {
			mod, _ := strconv.Atoi(req.FormValue(a.Name))

			if mod != a.Total {

				modVal := mod - a.Total

				update := CreateUpdate(eventString, modVal)

				if a.Updates == nil {
					a.Updates = []*runequest.Update{}
				}

				a.Updates = append(a.Updates, update)

				a.ExperienceCheck = false
			}
			if a.UserString != "" {
				us := req.FormValue(fmt.Sprintf("%s-UserString", a.Name))
				a.UserString = ProcessUserString(us)
			}

			a.UpdateAbility()

			// Remove Character Passion if zeroed out
			if a.Total < 1 {
				delete(c.Abilities, k)
			}
		}

		// Set Open to true if user authorizes
		if req.FormValue("Archive") != "" {
			cm.Open = true
		} else {
			cm.Open = false
		}

		c.Clan = req.FormValue("Clan")
		c.Tribe = req.FormValue("Tribe")

		// Upload image to s3
		file, h, err := req.FormFile("image")
		switch err {
		case nil:
			if h.Filename != "" {
				// Process image
				defer file.Close()

				err = ProcessImage(h, file, cm)
				if err != nil {
					log.Printf("Error processing image: %v", err)
				}

			} else {
				fmt.Println("No file provided.")
			}

		case http.ErrMissingFile:
			log.Println("no file")
			fmt.Println("Path: ", cm.Image.Path)

		default:
			log.Panic(err)
			fmt.Println("Error getting file ", err)
			cm.Image.Path = DefaultCharacterPortrait
		}

		cm.Random = false

		err = database.UpdateCharacterModel(db, cm)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("Saved")
		}

		session.AddFlash("Character Updated with "+eventString, "message")
		session.Save(req, w)

		url := fmt.Sprintf("/view_character/%d", cm.ID)

		http.Redirect(w, req, url, http.StatusSeeOther)
		return
	}
}

// AddSkillsHandler renders a character in a Web page
func AddSkillsHandler(w http.ResponseWriter, req *http.Request) {

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
		http.Redirect(w, req, "/", http.StatusFound)
		return
	}

	vars := mux.Vars(req)
	pk := vars["id"]

	if len(pk) == 0 {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	id, err := strconv.Atoi(pk)
	if err != nil {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
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
		SessionUser:    username,
		IsLoggedIn:     loggedIn,
		IsAdmin:        isAdmin,
		IsAuthor:       IsAuthor,
		Counter:        numToArray(6),
		Skills:         runequest.Skills,
		CategoryOrder:  runequest.CategoryOrder,
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/add_skills.html", wc)

	}

	if req.Method == "POST" {

		err := req.ParseMultipartForm(MaxMemory)
		if err != nil {
			panic(err)
		}

		event := req.FormValue("Event")

		// Add Skills
		for i := 1; i < 7; i++ {
			coreString := req.FormValue(fmt.Sprintf("Skill-%d-CoreString", i))
			userString := req.FormValue(fmt.Sprintf("Skill-%d-UserString", i))

			valueStr := req.FormValue(fmt.Sprintf("Skill-%d-Value", i))
			value, err := strconv.Atoi(valueStr)
			if err != nil {
				value = 0
			}

			// Skill exists in Character, modify it via pointer
			if coreString != "" {
				// Determine if skill already exists in c.Skills

				bs := runequest.Skills[coreString]

				sk := &runequest.Skill{
					CoreString: bs.CoreString,
					UserString: bs.UserString,
					Category:   bs.Category,
					Base:       bs.Base,
					UserChoice: bs.UserChoice,
					Updates:    []*runequest.Update{},
				}

				if userString != "" {
					sk.UserString = ProcessUserString(userString)
				}

				sk.GenerateName()

				// Add Skill to Character
				fmt.Println("Add Skill to character: " + sk.Name)
				c.Skills[sk.Name] = sk
				c.Skills[sk.Name].UpdateSkill()

				t := time.Now()
				tString := t.Format("2006-01-02 15:04:05")

				update := &runequest.Update{
					Date:  tString,
					Event: event,
					Value: value - (c.SkillCategories[sk.Category].Value + sk.Base),
				}

				c.Skills[sk.Name].Updates = append(c.Skills[sk.Name].Updates, update)

				c.Skills[sk.Name].UpdateSkill()

				fmt.Println("Updated Character Skill: " + event)
			}
		}

		// Add custom Skills
		for i := 1; i < 7; i++ {
			str := fmt.Sprintf("Custom-Skill-%d-", i)

			coreString := req.FormValue(str + "CoreString")
			val, err := strconv.Atoi(req.FormValue(str + "Value"))
			if err != nil {
				val = 0
			}

			if coreString != "" && val > 0 {

				category := req.FormValue(str + "Category")

				sk := &runequest.Skill{
					CoreString: coreString,
					Category:   category,
					Value:      val - c.SkillCategories[category].Value,
					Total:      val - c.SkillCategories[category].Value,
				}

				userString := req.FormValue(str + "UserString")
				if userString != "" {
					sk.UserString = ProcessUserString(userString)
				}

				sk.GenerateName()
				sk.UpdateSkill()
				c.Skills[sk.Name] = sk
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
		return
	}
}

// AddPassionsHandler renders a character in a Web page
func AddPassionsHandler(w http.ResponseWriter, req *http.Request) {

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
		http.Redirect(w, req, "/", http.StatusFound)
		return
	}

	vars := mux.Vars(req)
	pk := vars["id"]

	if len(pk) == 0 {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	id, err := strconv.Atoi(pk)
	if err != nil {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
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
		SessionUser:    username,
		IsLoggedIn:     loggedIn,
		IsAdmin:        isAdmin,
		IsAuthor:       IsAuthor,
		Counter:        numToArray(6),
		Passions:       runequest.PassionTypes,
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/add_passions.html", wc)

	}

	if req.Method == "POST" {

		err := req.ParseMultipartForm(MaxMemory)
		if err != nil {
			panic(err)
		}

		// Add passions
		for i := 1; i < 7; i++ {

			coreString := req.FormValue(fmt.Sprintf("Passion-%d-CoreString", i))
			userString := req.FormValue(fmt.Sprintf("Passion-%d-UserString", i))

			if coreString != "" {

				p := &runequest.Ability{
					Type:       "Passion",
					CoreString: coreString,
					Updates:    []*runequest.Update{},
				}

				if userString != "" {
					p.UserChoice = true
					p.UserString = ProcessUserString(userString)
				}

				targetString := createName(p.CoreString, p.UserString)

				str := fmt.Sprintf("Passion-%d-Base", i)
				base, err := strconv.Atoi(req.FormValue(str))
				if err != nil {
					base = 0
				}

				update := CreateUpdate("Add Passions", base)
				p.Updates = append(p.Updates, update)

				p.UpdateAbility()
				c.Abilities[targetString] = p
			}
		}

		// Add custom passions
		for i := 1; i < 7; i++ {

			coreString := req.FormValue(fmt.Sprintf("Custom-Passion-%d-CoreString", i))
			userString := req.FormValue(fmt.Sprintf("Custom-Passion-%d-UserString", i))

			if coreString != "" {

				p := &runequest.Ability{
					Type:       "Passion",
					CoreString: ProcessUserString(coreString),
					Updates:    []*runequest.Update{},
				}

				if userString != "" {
					p.UserChoice = true
					p.UserString = ProcessUserString(userString)
				}

				targetString := createName(p.CoreString, p.UserString)

				str := fmt.Sprintf("Custom-Passion-%d-Base", i)
				base, err := strconv.Atoi(req.FormValue(str))
				if err != nil {
					base = 0
				}

				update := CreateUpdate("Add Custom Passions", base)
				p.Updates = append(p.Updates, update)

				p.UpdateAbility()
				c.Abilities[targetString] = p
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
		return
	}
}

// EditMagicHandler renders a character in a Web page
func EditMagicHandler(w http.ResponseWriter, req *http.Request) {

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
		http.Redirect(w, req, "/", http.StatusFound)
		return
	}

	vars := mux.Vars(req)
	pk := vars["id"]

	if len(pk) == 0 {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	id, err := strconv.Atoi(pk)
	if err != nil {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
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

	ecLen := len(c.ExtraCults)

	if ecLen == 0 {
		ecLen = 1
	}

	// Add ExtraCults
	for i := ecLen; i < 4; i++ {
		ec := &runequest.ExtraCult{
			Name:       "",
			RunePoints: 0,
			Rank:       "",
		}
		c.ExtraCults = append(c.ExtraCults, ec)
	}

	// Add extra runespells
	for i := 1; i < 4; i++ {
		trs := &runequest.Spell{
			Name:       "",
			CoreString: "",
			UserString: "",
			Cost:       0,
		}
		c.RuneSpells["zzNewRS-"+string(i)] = trs
	}

	// Add extra spirit magic
	for i := 1; i < 4; i++ {
		tsm := &runequest.Spell{
			Name:       "",
			CoreString: "",
			UserString: "",
			Cost:       0,
		}
		tsm.Name = createName(tsm.CoreString, tsm.UserString)
		c.SpiritMagic["zzNewSM-"+string(i)] = tsm
	}

	// Add extra Powers
	if c.Powers == nil {
		c.Powers = map[string]*runequest.Power{}
	}

	for i := 1; i < 6; i++ {
		tp := &runequest.Power{
			Name:        "",
			Description: "",
		}
		c.Powers[fmt.Sprintf("zzNewPow-%d", i)] = tp
	}

	wc := WebChar{
		CharacterModel: cm,
		SessionUser:    username,
		IsLoggedIn:     loggedIn,
		IsAdmin:        isAdmin,
		IsAuthor:       IsAuthor,
		Counter:        numToArray(3),
		SpiritMagic:    runequest.SpiritMagicSpells,
		RuneSpells:     runequest.RuneSpells,
		CultModels:     cults,
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/edit_magic.html", wc)

	}

	if req.Method == "POST" {

		err := req.ParseMultipartForm(MaxMemory)
		if err != nil {
			panic(err)
		}

		eventString := req.FormValue("Event")

		// Power

		stat := c.Statistics["POW"]

		mod, _ := strconv.Atoi(req.FormValue("Power"))

		if mod != stat.Total {

			modVal := mod - stat.Total

			update := CreateUpdate(eventString, modVal)

			if stat.Updates == nil {
				stat.Updates = []*runequest.Update{}
			}

			stat.Updates = append(stat.Updates, update)

			stat.ExperienceCheck = false
		}

		stat.UpdateStatistic()

		// Primary Cults
		c.Cult.Rank = req.FormValue("Rank")

		rp, err := strconv.Atoi(req.FormValue("RunePoints"))
		if err != nil {
			rp = 3
			fmt.Println("Not a number")
		}

		c.Cult.NumRunePoints = rp

		// Secondary Cults

		tempCults := []*runequest.ExtraCult{}

		for i := 1; i < 4; i++ {
			eCult := req.FormValue(fmt.Sprintf("Cult-%d", i))

			if eCult != "" {

				eRank := req.FormValue(fmt.Sprintf("Rank-%d", i))
				eRP, err := strconv.Atoi(req.FormValue(fmt.Sprintf("RunePoints-%d", i)))
				if err != nil {
					eRP = 1
					fmt.Println("Not a number")
				}

				tempCults = append(tempCults,
					&runequest.ExtraCult{
						Name:       ProcessUserString(eCult),
						Rank:       eRank,
						RunePoints: eRP,
					})
			}
		}

		c.ExtraCults = tempCults

		// Rune Magic

		tempRuneSpells := map[string]*runequest.Spell{}

		for k, v := range c.RuneSpells {
			str := req.FormValue(fmt.Sprintf("RuneSpell-%s", k))
			spec := req.FormValue(fmt.Sprintf("RuneSpell-%s-UserString", k))

			fmt.Println(str, spec)

			switch {
			case str == "":
				fmt.Println("No spell")
				continue

			case v.CoreString == str && v.UserString == spec:
				// Same spell as previous
				fmt.Println("Same spell")
				tempRuneSpells[v.Name] = v

			case v.CoreString != str || v.UserString != spec:
				// Spell has changed or is new
				fmt.Printf("New spell: %s - %s", str, spec)
				// Get info from runequest.RuneSpells
				index, err := indexSpell(str, runequest.RuneSpells)
				if err != nil {
					fmt.Println(err)
					continue
				}

				baseSpell := runequest.RuneSpells[index]

				s := &runequest.Spell{
					Name:       baseSpell.Name,
					CoreString: baseSpell.CoreString,
					UserString: baseSpell.UserString,
					Cost:       baseSpell.Cost,
					Runes:      baseSpell.Runes,
					Domain:     baseSpell.Domain,
				}

				if spec != "" {
					s.UserString = ProcessUserString(spec)
				}

				s.GenerateName()
				tempRuneSpells[s.Name] = s
			}
		}

		c.RuneSpells = tempRuneSpells

		// Spirit Magic

		tempSpiritMagic := map[string]*runequest.Spell{}

		for k, v := range c.SpiritMagic {
			str := req.FormValue(fmt.Sprintf("SpiritMagic-%s", k))
			spec := req.FormValue(fmt.Sprintf("SpiritMagic-%s-UserString", k))
			cString := req.FormValue(fmt.Sprintf("SpiritMagic-%s-Cost", k))

			cost, err := strconv.Atoi(cString)
			if err != nil {
				cost = 1
				fmt.Println("Non-number entered")
			}

			switch {
			case str == "":
				fmt.Println("No spell")
				continue

			case v.CoreString == str && v.UserString == spec && v.Cost == cost:
				fmt.Println("Same spell")

				// Same spell as previous
				tempSpiritMagic[v.Name] = v

			case v.CoreString != str || v.UserString != spec || v.Cost != cost:
				// Spell has changed or is new
				fmt.Printf("NEW Spell: %s - %s", str, spec)

				// Get info from runequest.RuneSpells
				index, err := indexSpell(str, runequest.SpiritMagicSpells)
				if err != nil {
					fmt.Println(err)
					continue
				}

				baseSpell := runequest.SpiritMagicSpells[index]

				s := &runequest.Spell{
					Name:       baseSpell.Name,
					CoreString: baseSpell.CoreString,
					UserString: baseSpell.UserString,
					Cost:       cost,
					Domain:     baseSpell.Domain,
				}

				if spec != "" {
					s.UserString = ProcessUserString(spec)
				}

				s.GenerateName()
				tempSpiritMagic[s.Name] = s
			}
		}

		c.SpiritMagic = tempSpiritMagic

		// Update Powers

		tempPowers := map[string]*runequest.Power{}

		for k := range c.Powers {
			name := req.FormValue(fmt.Sprintf("Power-%s-Name", k))
			desc := req.FormValue(fmt.Sprintf("Power-%s-Description", k))

			if name != "" {
				p := &runequest.Power{
					Name:        ProcessUserString(name),
					Description: desc,
				}
				tempPowers[name] = p
			}
		}

		c.Powers = tempPowers

		// Update CM

		err = database.UpdateCharacterModel(db, cm)
		if err != nil {
			log.Panic(err)
		} else {
			fmt.Println("Saved")
		}

		url := fmt.Sprintf("/view_character/%d", cm.ID)

		http.Redirect(w, req, url, http.StatusSeeOther)
		return
	}
}

// DeleteCharacterHandler renders a character in a Web page
func DeleteCharacterHandler(w http.ResponseWriter, req *http.Request) {

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

	vars := mux.Vars(req)
	pk := vars["id"]

	if len(pk) == 0 {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	id, err := strconv.Atoi(pk)
	if err != nil {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	cm, err := database.PKLoadCharacterModel(db, int64(id))
	if err != nil {
		fmt.Println(err)
	}

	// Validate that User == Author
	IsAuthor := false

	if username == cm.Author.UserName || isAdmin == "true" {
		IsAuthor = true
	} else {
		http.Redirect(w, req, "/", http.StatusFound)
		return
	}

	if cm.Image == nil {
		cm.Image = new(models.Image)
		cm.Image.Path = DefaultCharacterPortrait
	}

	wc := WebChar{
		CharacterModel: cm,
		SessionUser:    username,
		IsAuthor:       IsAuthor,
		IsLoggedIn:     loggedIn,
		IsAdmin:        isAdmin,
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/delete_character.html", wc)

	}

	if req.Method == "POST" {

		database.DeleteCharacterModel(db, cm.ID)

		fmt.Println("Deleted ", cm.Character.Name)
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
}
