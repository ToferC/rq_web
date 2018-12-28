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

// CreatureIndexHandler renders the basic creature roster page
func CreatureIndexHandler(w http.ResponseWriter, req *http.Request) {

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

	creatures, err := database.ListCreatureModels(db)
	if err != nil {
		panic(err)
	}

	for _, cm := range creatures {
		if cm.Image == nil {
			cm.Image = new(models.Image)
			cm.Image.Path = DefaultCreaturePortrait
		}
	}

	wc := WebChar{
		SessionUser:    username,
		IsLoggedIn:     loggedIn,
		IsAdmin:        isAdmin,
		CreatureModels: creatures,
	}

	Render(w, "templates/roster.html", wc)
}

// CreatureHandler renders a creature in a Web page
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

	cm, err := database.PKLoadCreatureModel(db, int64(id))
	if err != nil {
		fmt.Println(err)
		fmt.Println("Unable to load CreatureModel")
	}

	fmt.Println(cm)

	IsAuthor := false

	if username == cm.Author.UserName {
		IsAuthor = true
	}

	c := cm.Creature

	// Always create 4 empty equipment slots.
	for i := 0; i < 4; i++ {
		c.Equipment = append(c.Equipment, "")
	}

	fmt.Println(c)

	if cm.Image == nil {
		cm.Image = new(models.Image)
		cm.Image.Path = DefaultCreaturePortrait
	}

	//c.DetermineSkillCategoryValues()

	wc := WebChar{
		CreatureModel: cm,
		IsAuthor:      IsAuthor,
		IsLoggedIn:    loggedIn,
		SessionUser:   username,
		IsAdmin:       isAdmin,
		Counter:       numToArray(10),
		Flashes:       flashes,
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/view_creature.html", wc)

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

		err = database.UpdateCreatureModel(db, cm)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("Saved")
		}

		fmt.Println(c)

		url := fmt.Sprintf("/view_creature/%d#gameplay", cm.ID)

		http.Redirect(w, req, url, http.StatusSeeOther)
	}

}

// NewCreatureHandler renders a creature in a Web page
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

	cm := models.CreatureModel{}

	//c := runequest.NewCreature("Default")
	c := &runequest.Creature{}

	//vars := mux.Vars(req)

	// Assign additional empty HitLocations to populate form

	for i := 0; i < 10; i++ {
		t := runequest.HitLocation{
			Name: "",
		}
		c.HitLocations["z"+string(i)] = &t
	}

	author, err := database.LoadUser(db, username)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, req, "/", 302)
	}

	fmt.Println(author)

	cm = models.CreatureModel{
		Creature: c,
		Author:   author,
	}

	wc := WebChar{
		CreatureModel: &cm,
		SessionUser:   username,
		IsLoggedIn:    loggedIn,
		IsAdmin:       isAdmin,
		Counter:       []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
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

		c := &runequest.Creature{}

		c.Description = req.FormValue("Description")

		for _, st := range runequest.StatMap {
			c.Statistics[st].Base, _ = strconv.Atoi(req.FormValue(st))
		}

		for _, sk := range c.Skills {
			sk.Value, _ = strconv.Atoi(req.FormValue(sk.Name))
			if sk.UserChoice {
				sk.UserString = req.FormValue(fmt.Sprintf("%s-Spec", sk.Name))
			}
		}

		// Hit locations - need to add new map or amend old one

		newHL := map[string]*runequest.HitLocation{}

		for i := range c.HitLocations {

			name := req.FormValue(fmt.Sprintf("%s-Name", i))

			if name != "" {

				max, _ := strconv.Atoi(req.FormValue(fmt.Sprintf("%s-Max", i)))
				armor, _ := strconv.Atoi(req.FormValue(fmt.Sprintf("%s-Armor", i)))

				fmt.Println(name, max, armor)

				newHL[name] = &runequest.HitLocation{
					Name:   name,
					Max:    max,
					Armor:  armor,
					HitLoc: []int{},
				}

				newHL[name].FillWounds()

				for j := 1; j < 11; j++ {
					if req.FormValue(fmt.Sprintf("%s-%d-loc", i, j)) != "" {
						newHL[name].HitLoc = append(newHL[name].HitLoc, j)
					}
				}
			}
		}

		fmt.Println(newHL)
		c.HitLocations = newHL

		cm.Creature = c

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

		fmt.Println(c)

		err = database.SaveCreatureModel(db, &cm)
		if err != nil {
			log.Panic(err)
		} else {
			fmt.Println("Saved")
		}

		url := fmt.Sprintf("/view_creature/%d", cm.ID)

		http.Redirect(w, req, url, http.StatusSeeOther)
	}
}

// ModifyCreatureHandler renders a creature in a Web page
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

	// Load CreatureModel
	cm, err := database.PKLoadCreatureModel(db, int64(id))
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

	c := cm.Creature

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
		cm.Image.Path = DefaultCreaturePortrait
	}

	wc := WebChar{
		CreatureModel: cm,
		SessionUser:   username,
		IsAuthor:      IsAuthor,
		IsLoggedIn:    loggedIn,
		IsAdmin:       isAdmin,
	}

	if req.Method == "GET" {

		// Render page

		Render(w, "templates/modify_creature.html", wc)

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

			// Remove Creature skill is zeroed out
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
			cm.Image.Path = DefaultCreaturePortrait
		}

		err = database.UpdateCreatureModel(db, cm)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("Saved")
		}

		session.AddFlash("Creature Updated with "+eventString, "message")
		session.Save(req, w)

		url := fmt.Sprintf("/view_creature/%d", cm.ID)

		http.Redirect(w, req, url, http.StatusSeeOther)
	}
}

// AddCreatureContentHandler renders a creature in a Web page
func AddCreatureContentHandler(w http.ResponseWriter, req *http.Request) {

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

	cm, err := database.PKLoadCreatureModel(db, int64(id))
	if err != nil {
		fmt.Println(err)
		fmt.Println("Unable to load CreatureModel")
	}

	c := cm.Creature

	IsAuthor := false

	if username == cm.Author.UserName {
		IsAuthor = true
	}

	wc := WebChar{
		CreatureModel: cm,
		SessionUser:   username,
		IsLoggedIn:    loggedIn,
		IsAdmin:       isAdmin,
		IsAuthor:      IsAuthor,
		Counter:       numToArray(3),
		Passions:      runequest.PassionTypes,
		Skills:        runequest.Skills,
		SpiritMagic:   runequest.SpiritMagicSpells,
		RuneSpells:    runequest.RuneSpells,
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/add_creature_content.html", wc)

	}

	if req.Method == "POST" {

		err := req.ParseMultipartForm(MaxMemory)
		if err != nil {
			panic(err)
		}

		event := req.FormValue("Event")

		// Add Skills
		for i := 1; i < 4; i++ {
			coreString := req.FormValue(fmt.Sprintf("Skill-%d-CoreString", i))
			userString := req.FormValue(fmt.Sprintf("Skill-%d-UserString", i))

			valueStr := req.FormValue(fmt.Sprintf("Skill-%d-Value", i))
			value, err := strconv.Atoi(valueStr)
			if err != nil {
				value = 0
			}

			// Skill exists in Creature, modify it via pointer
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
					sk.UserString = userString
				}

				sk.GenerateName()

				// Add Skill to Creature
				fmt.Println("Add Skill to creature: " + sk.Name)
				c.Skills[sk.Name] = sk
				c.Skills[sk.Name].UpdateSkill()

				t := time.Now()
				tString := t.Format("2006-01-02 15:04:05")

				update := &runequest.Update{
					Date:  tString,
					Event: event,
					Value: value,
				}

				c.Skills[sk.Name].Updates = append(c.Skills[sk.Name].Updates, update)

				c.Skills[sk.Name].UpdateSkill()

				fmt.Println("Updated Creature Skill: " + event)
			}
		}

		// Add passions
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

				userString := req.FormValue(fmt.Sprintf("Passion-%d-UserString", i))

				if userString != "" {
					p.UserChoice = true
					p.UserString = userString
				}

				//c.ModifyAbility(p)
				p.UpdateAbility()
			}
		}

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

				s := baseSpell
				if spec != "" {
					s.UserString = spec
				}
				s.GenerateName()
				c.RuneSpells[s.Name] = &s
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
				c.SpiritMagic[s.Name] = s
			}
		}

		err = database.UpdateCreatureModel(db, cm)
		if err != nil {
			log.Panic(err)
		} else {
			fmt.Println("Saved")
		}

		url := fmt.Sprintf("/view_creature/%d", cm.ID)

		http.Redirect(w, req, url, http.StatusSeeOther)
	}
}

// DeleteCreatureHandler renders a creature in a Web page
func DeleteCreatureHandler(w http.ResponseWriter, req *http.Request) {

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
	}

	id, err := strconv.Atoi(pk)
	if err != nil {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}

	cm, err := database.PKLoadCreatureModel(db, int64(id))
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

	if cm.Image == nil {
		cm.Image = new(models.Image)
		cm.Image.Path = DefaultCreaturePortrait
	}

	wc := WebChar{
		CreatureModel: cm,
		SessionUser:   username,
		IsAuthor:      IsAuthor,
		IsLoggedIn:    loggedIn,
		IsAdmin:       isAdmin,
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/delete_creature.html", wc)

	}

	if req.Method == "POST" {

		database.DeleteCreatureModel(db, cm.ID)

		fmt.Println("Deleted ", cm.Creature.Name)
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}
}