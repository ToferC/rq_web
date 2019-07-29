package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gorilla/mux"
	"github.com/gosimple/slug"
	"github.com/thewhitetulip/Tasks/sessions"
	"github.com/toferc/rq_web/database"
	"github.com/toferc/rq_web/models"
	"github.com/toferc/runequest"
)

// CreatureHandler renders a character in a Web page
func CreatureHandler(w http.ResponseWriter, req *http.Request) {

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

	fmt.Println(cm)

	IsAuthor := false

	if username == cm.Author.UserName {
		IsAuthor = true
	}

	c := cm.Character

	// Always create 4 empty equipment slots.
	for i := 0; i < 4; i++ {
		c.Equipment = append(c.Equipment, "")
	}

	fmt.Println(c)

	if cm.Image == nil {
		cm.Image = new(models.Image)
		cm.Image.Path = DefaultCharacterPortrait
	}

	//c.DetermineSkillCategoryValues()

	wc := WebChar{
		CharacterModel: cm,
		IsAuthor:       IsAuthor,
		IsLoggedIn:     loggedIn,
		SessionUser:    username,
		IsAdmin:        isAdmin,
		Counter:        numToArray(10),
		Flashes:        flashes,
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

		// Update MP
		str = req.FormValue("RP")
		rp, err := strconv.Atoi(str)
		if err != nil {
			rp = c.CurrentMP
		}

		if rp < 0 {
			rp = 0
		}

		c.CurrentRP = rp

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

		// Always create 4 empty equipment slots.
		for i := 0; i < 4; i++ {
			c.Equipment = append(c.Equipment, "")
		}

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

		err = database.UpdateCharacterModel(db, cm)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("Saved")
		}

		fmt.Println(c)

		url := fmt.Sprintf("/view_character/%d#gameplay", cm.ID)

		http.Redirect(w, req, url, http.StatusSeeOther)
	}

}

// NewCreatureHandler renders a character in a Web page
func NewCreatureHandler(w http.ResponseWriter, req *http.Request) {

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

	author, err := database.LoadUser(db, username)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, req, "/", 302)
	}

	fmt.Println(author)

	cm := models.CharacterModel{
		Author: author,
	}

	c := runequest.NewCharacter("")

	cm.Character = c

	c.Abilities = map[string]*runequest.Ability{
		"Reputation": &runequest.Ability{
			CoreString: "Reputation",
			Type:       "Base",
			Value:      5,
		},
	}

	// Add Movement
	for i := 0; i < 2; i++ {
		var moveName string
		var mv int

		c.Movement = append(c.Movement,
			&runequest.Movement{
				Name:  moveName,
				Value: mv,
			})
	}

	c.Skills = map[string]*runequest.Skill{}

	// Set 6 skills
	for i := 0; i < 6; i++ {
		s := runequest.Skill{
			Name:       "",
			CoreString: "",
			UserString: "",
			Base:       0,
			Value:      0,
			Category:   "Agility",
		}
		c.Skills[fmt.Sprintf("z-%d", i)] = &s
	}

	cults, err := database.ListCultModels(db)
	if err != nil {
		panic(err)
	}

	// Add extra runespells
	for i := 1; i < 6; i++ {
		trs := &runequest.Spell{
			Name:       "",
			CoreString: "",
			UserString: "",
			Cost:       0,
		}
		c.RuneSpells["zzNewRS-"+string(i)] = trs
	}

	// Add extra spirit magic
	for i := 1; i < 6; i++ {
		tsm := &runequest.Spell{
			Name:       "",
			CoreString: "",
			UserString: "",
			Cost:       0,
		}
		tsm.Name = createName(tsm.CoreString, tsm.UserString)
		c.SpiritMagic["zzNewSM-"+string(i)] = tsm
	}

	wc := WebChar{
		CharacterModel:   &cm,
		SessionUser:      username,
		Counter:          numToArray(4),
		BigCounter:       numToArray(7),
		IsLoggedIn:       loggedIn,
		IsAdmin:          isAdmin,
		Passions:         runequest.PassionTypes,
		Skills:           runequest.Skills,
		SpiritMagic:      runequest.SpiritMagicSpells,
		RuneSpells:       runequest.RuneSpells,
		CultModels:       cults,
		CategoryOrder:    runequest.CategoryOrder,
		HitLocationForms: runequest.LocationForms,
		StringArray:      runequest.StatMap,
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/add_creature.html", wc)

	}

	if req.Method == "POST" {

		err := req.ParseMultipartForm(MaxMemory)
		if err != nil {
			panic(err)
		}

		err = decoder.Decode(c, req.PostForm)
		if err != nil {
			panic(err)
		}

		cm.Slug = slug.Make(fmt.Sprintf("%s-%s", username, c.Name))

		for _, st := range runequest.StatMap {
			num, err := strconv.Atoi(req.FormValue(st))
			if err != nil {
				num = 10
			}
			c.Statistics[st].Base = num
			c.Statistics[st].Max = num
		}

		c.DetermineSkillCategoryValues()
		c.SetAttributes()

		tempMovement := []*runequest.Movement{}

		// Add Movement
		for i := 1; i < 4; i++ {
			moveName := req.FormValue(fmt.Sprintf("Move-Name-%d", i))

			if moveName != "" {
				mv, err := strconv.Atoi(req.FormValue(fmt.Sprintf("Move-Value-%d", i)))
				if err != nil {
					mv = 8
				}

				tempMovement = append(tempMovement,
					&runequest.Movement{
						Name:  moveName,
						Value: mv,
					})
			}
		}

		c.Movement = tempMovement

		// Add Skills
		newSkills := map[string]*runequest.Skill{}

		for k := range c.Skills {
			str := fmt.Sprintf("Skill-%s-", k)

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
					sk.UserString = userString
				}

				sk.GenerateName()
				sk.UpdateSkill()
				newSkills[sk.Name] = sk
			}
		}

		// Update skills
		c.Skills = newSkills

		// Add passions
		for i := 1; i < 6; i++ {

			coreString := req.FormValue(fmt.Sprintf("Passion-%d-CoreString", i))
			userString := req.FormValue(fmt.Sprintf("Passion-%d-UserString", i))

			str := fmt.Sprintf("Passion-%d-Base", i)
			base, err := strconv.Atoi(req.FormValue(str))
			if err != nil {
				base = 0
			}

			if coreString != "" && base > 0 {

				targetString := createName(coreString, userString)

				// No ability
				a := &runequest.Ability{
					Type:       "Passion",
					CoreString: coreString,
					Base:       base,
					Total:      base,
				}

				if userString != "" {
					a.UserString = userString
				}

				a.UpdateAbility()

				c.Abilities[targetString] = a

			}
		}

		// Update Elemental Runes
		for k, v := range c.ElementalRunes {
			val, err := strconv.Atoi(req.FormValue(k))
			if err != nil {
				val = 0
			}

			if val > 0 {
				v.Base = val
				v.Total = val
				v.UpdateAbility()
			}
		}

		// Update Power Runes
		for k, v := range c.PowerRunes {
			val, err := strconv.Atoi(req.FormValue(k))
			if err != nil {
				val = 0
			}

			if val > 0 {
				v.Base = val
				v.Total = val
				v.UpdateAbility()
			}
		}

		// Update Condition Runes
		for k, v := range c.ConditionRunes {
			val, err := strconv.Atoi(req.FormValue(k))
			if err != nil {
				val = 0
			}

			if val > 0 {
				v.Base = val
				v.Total = val
				v.UpdateAbility()
			}
		}

		// Update Character

		form := req.FormValue("Hit-Location-Form")

		c.LocationForm = form
		c.HitLocations = runequest.LocationForms[form]
		c.HitLocationMap = runequest.SortLocations(c.HitLocations)

		c.CurrentHP = c.Attributes["HP"].Max
		c.CurrentMP = c.Attributes["MP"].Max
		c.CurrentRP = c.Cult.NumRunePoints

		armorStr := req.FormValue("Armor")
		armor, err := strconv.Atoi(armorStr)
		if err != nil {
			armor = 0
		}
		for _, v := range c.HitLocations {
			v.Armor = armor
			v.Value = v.Max
		}

		fmt.Println(c.HitLocations)

		c.UpdateCharacter()

		// Update Powers
		c.Powers = map[string]*runequest.Power{}

		for _, n := range numToArray(7) {
			str := fmt.Sprintf("Power-%d-", n)
			name := req.FormValue(str + "Name")
			if name != "" {
				c.Powers[name] = &runequest.Power{
					Name:        name,
					Description: req.FormValue(str + "Description"),
				}
			}
		}

		meleeAttacks := map[string]*runequest.Attack{}
		rangedAttacks := map[string]*runequest.Attack{}

		// Melee Weapons & Attacks
		for _, i := range numToArray(7) {
			attackName := req.FormValue(fmt.Sprintf("Attack-%d-Name", i))
			attackType := req.FormValue(fmt.Sprintf("Attack-%d-Type", i))
			special := req.FormValue(fmt.Sprintf("Attack-%d-Special", i))
			damage := req.FormValue(fmt.Sprintf("Attack-%d-Damage", i))

			skillString := req.FormValue(fmt.Sprintf("Attack-%d-Skill", i))

			skillVal, err := strconv.Atoi(skillString)
			if err != nil {
				skillVal = 0
			}

			hpStr := req.FormValue(fmt.Sprintf("Attack-%d-HP", i))
			hpVal, err := strconv.Atoi(hpStr)
			if err != nil {
				hpVal = 0
			}

			srStr := req.FormValue(fmt.Sprintf("Attack-%d-SR", i))
			srVal, err := strconv.Atoi(srStr)
			if err != nil {
				srVal = 0
			}

			rangeStr := req.FormValue(fmt.Sprintf("Attack-%d-Range", i))
			rangeVal, err := strconv.Atoi(rangeStr)
			if err != nil {
				rangeVal = 0
			}

			if attackName != "" && skillVal > 0 {

				// Create Weapon Skill
				sk := &runequest.Skill{
					CoreString: attackName,
					Value:      skillVal - c.SkillCategories[attackType].Value,
					Category:   attackType,
				}

				c.Skills[attackName] = sk

				if attackType == "Ranged" {

					rangedAttacks[attackName] = &runequest.Attack{
						Name:         attackName,
						Skill:        c.Skills[attackName],
						DamageString: damage,
						StrikeRank:   c.Attributes["DEXSR"].Base,
						Weapon: &runequest.Weapon{
							Name:      attackName,
							Type:      attackType,
							SR:        srVal,
							STRDamage: false,
							Range:     rangeVal,
							Damage:    damage,
							HP:        hpVal,
							CurrentHP: hpVal,
							Special:   special,
							Custom:    true,
						},
					}
				} else {

					dbString := ""

					if c.Attributes["DB"].Text != "-" {
						dbString = c.Attributes["DB"].Text
					}

					meleeAttacks[attackName] = &runequest.Attack{
						Name:         attackName,
						Skill:        c.Skills[attackName],
						DamageString: damage + dbString,
						StrikeRank:   c.Attributes["DEXSR"].Base + c.Attributes["SIZSR"].Base + srVal,
						Weapon: &runequest.Weapon{
							Name:      attackName,
							Type:      attackType,
							SR:        srVal,
							STRDamage: true,
							Damage:    damage,
							HP:        hpVal,
							CurrentHP: hpVal,
							Special:   special,
							Custom:    true,
						},
					}
				}
			}
		}

		c.MeleeAttacks = meleeAttacks
		c.RangedAttacks = rangedAttacks

		tempRuneSpells := map[string]*runequest.Spell{}

		for k, v := range c.RuneSpells {
			str := req.FormValue(fmt.Sprintf("RuneSpell-%s", k))
			spec := req.FormValue(fmt.Sprintf("RuneSpell-%s-UserString", k))

			fmt.Println(str, spec)

			switch {
			case str == "":
				fmt.Println("No spell")
				continue

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
					Domain:     baseSpell.Domain,
				}

				if spec != "" {
					s.UserString = spec
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
					s.UserString = spec
				}

				s.GenerateName()
				tempSpiritMagic[s.Name] = s
			}
		}

		c.SpiritMagic = tempSpiritMagic

		// Insert power into App archive if user authorizes
		if req.FormValue("Archive") != "" {
			cm.Open = true
		} else {
			cm.Open = false
		}

		// Upload image to s3
		file, h, err := req.FormFile("image")
		switch err {
		case nil:
			// Process image
			defer file.Close()
			// example path media/Major/TestImage/Jason_White.jpg
			path := fmt.Sprintf("/media/%s/%s/%s",
				cm.Author.UserName,
				runequest.ToSnakeCase(c.Name),
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

		c.CreationSteps["Complete"] = true

		fmt.Println(c)

		err = database.SaveCharacterModel(db, &cm)
		if err != nil {
			log.Panic(err)
		} else {
			fmt.Println("Saved")
		}

		url := fmt.Sprintf("/view_character/%d", cm.ID)

		http.Redirect(w, req, url, http.StatusSeeOther)
	}
}

// ModifyCreatureHandler renders a character in a Web page
func ModifyCreatureHandler(w http.ResponseWriter, req *http.Request) {

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

	if username == cm.Author.UserName || isAdmin == "true" {
		IsAuthor = true
	} else {
		http.Redirect(w, req, "/", 302)
	}

	c := cm.Character

	// Assign additional empty HitLocations to populate form
	/*
		if len(c.HitLocations) < 10 {
			for i := len(c.HitLocations); i < 10; i++ {
				t := runequest.HitLocation{
					Name: "",
				}
				c.HitLocations["z"+string(i)] = &t
			}
		}
	*/

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

		//c.Cult.Rank = req.FormValue("Rank")

		rp, err := strconv.Atoi(req.FormValue("RunePoints"))
		if err != nil {
			rp = 3
			fmt.Println("Not a number")
		}

		//c.Cult.NumRunePoints = rp
		fmt.Println(rp)

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
			}

			stat.UpdateStatistic()

		}

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
				s.UserString = req.FormValue(fmt.Sprintf("%s-UserString", s.Name))
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

					opposed := c.PowerRunes[v.OpposedAbility]

					// Update opposed Power Rune if needed
					if v.Total+opposed.Total > 100 {

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

		// Update Abilities
		for _, a := range c.Abilities {
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
				a.UserString = req.FormValue(fmt.Sprintf("%s-UserString", a.Name))
			}
			a.UpdateAbility()
		}

		// Hit locations - need to add new map or amend old one

		for k, v := range c.HitLocations {
			str := req.FormValue(fmt.Sprintf("%s-Armor", k))
			armor, err := strconv.Atoi(str)
			if err != nil {
				armor = v.Armor
			}
			v.Armor = armor

			str = req.FormValue(fmt.Sprintf("%s-HP-Max", k))
			max, err := strconv.Atoi(str)
			if err != nil {
				max = v.Max
			}
			v.Max = max
		}

		// Set Open to true if user authorizes
		if req.FormValue("Archive") != "" {
			cm.Open = true
		} else {
			cm.Open = false
		}

		//c.Clan = req.FormValue("Clan")
		//c.Tribe = req.FormValue("Tribe")

		// Upload image to s3
		file, h, err := req.FormFile("image")
		switch err {
		case nil:
			if h.Filename != "" {
				// Process image
				defer file.Close()
				// example path media/Major/TestImage/Jason_White.jpg
				path := fmt.Sprintf("/media/%s/%s/%s",
					cm.Author.UserName,
					runequest.ToSnakeCase(c.Name),
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

				if cm.Image == nil {
					cm.Image = new(models.Image)
				}
				cm.Image.Path = path

				fmt.Println(path)
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
	}
}
