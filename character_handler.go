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
		Counter:        []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
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

		for k, v := range c.HitLocations {
			for i := range v.Wounds {
				v.Wounds[i] = false
				if req.FormValue(fmt.Sprintf("%s-Shock-%d", k, i)) != "" {
					v.Wounds[i] = true
				}
			}
		}

		c.Equipment[0] = req.FormValue("Equipment")

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
	if len(c.HitLocations) < 10 {
		for i := len(c.HitLocations); i < 10; i++ {
			t := runequest.HitLocation{
				Name: "",
			}
			c.HitLocations["z"+string(i)] = &t
		}
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

		Render(w, "templates/modify_character.html", wc)

	}

	if req.Method == "POST" {

		err := req.ParseMultipartForm(MaxMemory)
		if err != nil {
			panic(err)
		}

		c.Name = req.FormValue("Name")

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
