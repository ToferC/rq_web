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

// CharacterIndexHandler renders the basic character roster page
func CharacterIndexHandler(w http.ResponseWriter, req *http.Request) {

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

	characters, err := database.ListCharacterModels(db)
	if err != nil {
		panic(err)
	}

	for _, cm := range characters {
		if cm.Image == nil {
			cm.Image = new(models.Image)
			cm.Image.Path = DefaultCharacterPortrait
		}
	}

	wc := WebChar{
		SessionUser:     username,
		IsLoggedIn:      loggedIn,
		IsAdmin:         isAdmin,
		CharacterModels: characters,
	}

	Render(w, "templates/roster.html", wc)
}

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

		// Update HP

		/*
			str = req.FormValue("HP")
			hp, err := strconv.Atoi(str)
			if err != nil {
				hp = c.CurrentHP
			}

			if hp > c.Attributes["HP"].Max {
				hp = c.Attributes["HP"].Max
			}

			c.CurrentHP = hp
		*/

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

		err = database.UpdateCharacterModel(db, cm)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("Saved")
		}

		fmt.Println(c)
		// Render page
		Render(w, "templates/view_character.html", wc)
	}

}

// NewCharacterHandler renders a character in a Web page
func NewCharacterHandler(w http.ResponseWriter, req *http.Request) {

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

	c := runequest.NewCharacter("Default")

	//vars := mux.Vars(req)

	// Assign additional empty HitLocations to populate form

	for i := 0; i < 10; i++ {
		t := runequest.HitLocation{
			Name: "",
		}
		c.HitLocations["z"+string(i)] = &t
	}

	author := database.LoadUser(db, username)
	fmt.Println(author)

	cm = models.CharacterModel{
		Character: c,
		Author:    author,
	}

	wc := WebChar{
		CharacterModel: &cm,
		SessionUser:    username,
		IsLoggedIn:     loggedIn,
		IsAdmin:        isAdmin,
		Counter:        []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/add_character.html", wc)

	}

	if req.Method == "POST" {

		err := req.ParseMultipartForm(MaxMemory)
		if err != nil {
			panic(err)
		}

		c := &runequest.Character{}

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

		cm.Character = c

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

		// Finalize base Character Cost for play
		if req.FormValue("InPlay") != "" {
			c.InPlay = true
		} else {
			c.InPlay = false
		}

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

// ModifyCharacterHandler renders a character in a Web page
func ModifyCharacterHandler(w http.ResponseWriter, req *http.Request) {

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

		c.Cult.Rank = req.FormValue("Rank")

		rp, err := strconv.Atoi(req.FormValue("RunePoints"))
		if err != nil {
			rp = 3
			fmt.Println("Not a number")
		}

		c.Cult.NumRunePoints = rp

		for _, st := range runequest.StatMap {

			stat := c.Statistics[st]

			mod, _ := strconv.Atoi(req.FormValue(st))

			if mod != stat.Total {

				modVal := mod - stat.Total

				t := time.Now()
				tString := t.Format("2006-01-02")

				update := &runequest.Update{
					Date:  tString,
					Event: fmt.Sprintf("%s", eventString),
					Value: modVal,
				}

				if stat.Updates == nil {
					stat.Updates = []*runequest.Update{}
				}

				stat.Updates = append(stat.Updates, update)

			}

			stat.UpdateStatistic()

		}

		for _, s := range c.Skills {
			mod, _ := strconv.Atoi(req.FormValue(s.Name))

			if mod != s.Total {

				modVal := mod - s.Total

				t := time.Now()
				tString := t.Format("2006-01-02")

				update := &runequest.Update{
					Date:  tString,
					Event: fmt.Sprintf("%s", eventString),
					Value: modVal,
				}

				if s.Updates == nil {
					s.Updates = []*runequest.Update{}
				}

				s.Updates = append(s.Updates, update)

			}
			if s.UserString != "" {
				s.UserString = req.FormValue(fmt.Sprintf("%s-UserString", s.Name))
			}
			s.UpdateSkill()
		}

		// Update Elemental Runes
		for k, v := range c.ElementalRunes {
			mod, err := strconv.Atoi(req.FormValue(k))
			if err != nil {
				mod = v.Total
			}

			if mod != v.Total {

				modVal := mod - v.Total

				t := time.Now()
				tString := t.Format("2006-01-02")

				update := &runequest.Update{
					Date:  tString,
					Event: fmt.Sprintf("%s", eventString),
					Value: modVal,
				}

				if v.Updates == nil {
					v.Updates = []*runequest.Update{}
				}

				v.Updates = append(v.Updates, update)

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

					t := time.Now()
					tString := t.Format("2006-01-02")

					update := &runequest.Update{
						Date:  tString,
						Event: fmt.Sprintf("%s", eventString),
						Value: modVal,
					}

					if v.Updates == nil {
						v.Updates = []*runequest.Update{}
					}

					v.Updates = append(v.Updates, update)

					v.UpdateAbility()

					opposed := c.PowerRunes[v.OpposedAbility]

					// Update opposed Power Rune if needed
					if v.Total+opposed.Total > 100 {

						opposedUpdate := &runequest.Update{
							Date:  tString,
							Event: fmt.Sprintf("%s", eventString),
							Value: -modVal,
						}

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

				t := time.Now()
				tString := t.Format("2006-01-02")

				update := &runequest.Update{
					Date:  tString,
					Event: fmt.Sprintf("%s", eventString),
					Value: modVal,
				}

				if a.Updates == nil {
					a.Updates = []*runequest.Update{}
				}

				a.Updates = append(a.Updates, update)

			}
			if a.UserString != "" {
				a.UserString = req.FormValue(fmt.Sprintf("%s-UserString", a.Name))
			}
			a.UpdateAbility()
		}

		// Hit locations - need to add new map or amend old one

		/*
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

					}
				}
			}

			fmt.Println(newHL)
			c.HitLocations = newHL

		*/

		// Set Open to true if user authorizes
		if req.FormValue("Archive") != "" {
			cm.Open = true
		} else {
			cm.Open = false
		}

		// Finalize base Character Cost for play
		if req.FormValue("InPlay") != "" {
			c.InPlay = true
		} else {
			c.InPlay = false
		}

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

		url := fmt.Sprintf("/view_character/%d", cm.ID)

		http.Redirect(w, req, url, http.StatusSeeOther)
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
	}

	id, err := strconv.Atoi(pk)
	if err != nil {
		http.Redirect(w, req, "/", http.StatusSeeOther)
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
		http.Redirect(w, req, "/", 302)
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
	}
}
