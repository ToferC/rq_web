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
			for i := range v.Value {
				v.Value[i] = false
				if req.FormValue(fmt.Sprintf("%s-Shock-%d", k, i)) != "" {
					v.Value[i] = true
				}
			}
		}

		c.Gear[0] = req.FormValue("Gear")

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

	vars := mux.Vars(req)
	setting := vars["setting"]

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
		Modifiers:      oneroll.Modifiers,
		Counter:        []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		Sources:        oneroll.Sources,
		Permissions:    oneroll.Permissions,
		Intrinsics:     oneroll.Intrinsics,
		Advantages:     nil,
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

		c := &oneroll.Character{}

		setting := req.FormValue("Setting")

		switch setting {
		case "SR":
			c = oneroll.NewSRCharacter(req.FormValue("Name"))
		case "WT":
			c = oneroll.NewWTCharacter(req.FormValue("Name"))
		case "RE":
			c = oneroll.NewReignCharacter(req.FormValue("Name"))
		}

		if setting == "SR" || setting == "WT" {
			c.Archetype = &oneroll.Archetype{
				Type: req.FormValue("Archetype"),
			}
			for _, s := range wc.Counter { // Loop

				sType := req.FormValue(fmt.Sprintf("Source-%d", s))

				pType := req.FormValue(fmt.Sprintf("Permission-%d", s))

				iName := req.FormValue(fmt.Sprintf("Intrinsic-%d-Name", s))

				iInfo := req.FormValue(fmt.Sprintf("Intrinsic-%d-Info", s))

				if iName != "" {
					i := oneroll.Intrinsics[iName]
					l, err := strconv.Atoi(req.FormValue(fmt.Sprintf("Intrinsic-%d-Level", s)))
					if err != nil {
						l = 1
					}
					i.Level = l
					i.Info = iInfo
					c.Archetype.Intrinsics = append(c.Archetype.Intrinsics, &i)
				}

				if sType != "" {
					tS := oneroll.Sources[sType]
					c.Archetype.Sources = append(c.Archetype.Sources, &tS)
				}
				if pType != "" {
					tP := oneroll.Permissions[pType]
					c.Archetype.Permissions = append(c.Archetype.Permissions, &tP)
				}
			}
		}

		c.Description = req.FormValue("Description")

		for _, st := range c.StatMap {
			c.Statistics[st].Dice.Normal, _ = strconv.Atoi(req.FormValue(st))
		}

		for _, sk := range c.Skills {
			sk.Dice.Normal, _ = strconv.Atoi(req.FormValue(sk.Name))
			if sk.ReqSpec {
				sk.Specialization = req.FormValue(fmt.Sprintf("%s-Spec", sk.Name))
			}
		}

		// Hit locations - need to add new map or amend old one

		newHL := map[string]*oneroll.Location{}

		for i := range c.HitLocations {

			name := req.FormValue(fmt.Sprintf("%s-Name", i))

			if name != "" {

				boxes, _ := strconv.Atoi(req.FormValue(fmt.Sprintf("%s-Boxes", i)))
				lar, _ := strconv.Atoi(req.FormValue(fmt.Sprintf("%s-LAR", i)))
				har, _ := strconv.Atoi(req.FormValue(fmt.Sprintf("%s-HAR", i)))

				fmt.Println(name, boxes, lar, har)

				newHL[name] = &oneroll.Location{
					Name:   name,
					Boxes:  boxes,
					LAR:    lar,
					HAR:    har,
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
				oneroll.ToSnakeCase(c.Name),
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

	if c.Setting != "RE" {
		a := c.Archetype

		// Assign additional empty Sources to populate form
		if len(a.Sources) < 4 {
			for i := len(a.Sources); i < 4; i++ {
				tempS := oneroll.Source{
					Type: "",
				}
				a.Sources = append(a.Sources, &tempS)
			}
		}

		// Assign additional empty Permissions to populate form
		if len(a.Permissions) < 4 {
			for i := len(a.Permissions); i < 4; i++ {
				tempP := oneroll.Permission{
					Type: "",
				}
				a.Permissions = append(a.Permissions, &tempP)
			}
		}

		// Assign additional empty Sources to populate form
		if len(a.Intrinsics) < 5 {
			for i := len(a.Intrinsics); i < 5; i++ {
				tempI := oneroll.Intrinsic{
					Name: "",
				}
				a.Intrinsics = append(a.Intrinsics, &tempI)
			}
		}

		// Assign additional empty HitLocations to populate form
		if len(c.HitLocations) < 10 {
			for i := len(c.HitLocations); i < 10; i++ {
				t := oneroll.Location{
					Name: "",
				}
				c.HitLocations["z"+string(i)] = &t
			}
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
		Modifiers:      oneroll.Modifiers,
		Counter:        []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		Sources:        oneroll.Sources,
		Permissions:    oneroll.Permissions,
		Intrinsics:     oneroll.Intrinsics,
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

		if c.Setting != "RE" {

			c.Archetype = &oneroll.Archetype{
				Type: req.FormValue("Archetype"),
			}

			for _, s := range wc.Counter[:3] { // Loop

				sType := req.FormValue(fmt.Sprintf("Source-%d", s))

				pType := req.FormValue(fmt.Sprintf("Permission-%d", s))

				if sType != "" {
					tS := oneroll.Sources[sType]
					c.Archetype.Sources = append(c.Archetype.Sources, &tS)
				}
				if pType != "" {
					tP := oneroll.Permissions[pType]
					c.Archetype.Permissions = append(c.Archetype.Permissions, &tP)
				}
			}

			for _, s := range wc.Counter[:5] {
				iName := req.FormValue(fmt.Sprintf("Intrinsic-%d-Name", s))

				if iName != "" {
					i := oneroll.Intrinsics[iName]

					if i.RequiresLevel {
						l, err := strconv.Atoi(req.FormValue(fmt.Sprintf("Intrinsic-%d-Level", s)))
						if err != nil {
							l = 1
						}
						i.Level = l
					}

					if i.RequiresInfo {
						iInfo := req.FormValue(fmt.Sprintf("Intrinsic-%d-Info", s))
						i.Info = iInfo
					}

					c.Archetype.Intrinsics = append(c.Archetype.Intrinsics, &i)
				}
			}
		}

		c.Description = req.FormValue("Description")

		bw, err := strconv.Atoi(req.FormValue("BaseWill"))
		if err != nil {
			bw = c.BaseWill
		}

		c.BaseWill = bw

		wp, err := strconv.Atoi(req.FormValue("Willpower"))
		if err != nil {
			wp = c.Willpower
		}

		c.Willpower = wp

		for _, st := range c.StatMap {
			c.Statistics[st].Dice.Normal, _ = strconv.Atoi(req.FormValue(st))
		}

		for _, sk := range c.Skills {
			sk.Dice.Normal, _ = strconv.Atoi(req.FormValue(sk.Name))
			if sk.ReqSpec {
				sk.Specialization = req.FormValue(fmt.Sprintf("%s-Spec", sk.Name))
			}
		}

		// Hit locations - need to add new map or amend old one

		newHL := map[string]*oneroll.Location{}

		for i := range c.HitLocations {

			name := req.FormValue(fmt.Sprintf("%s-Name", i))

			if name != "" {
				boxes, _ := strconv.Atoi(req.FormValue(fmt.Sprintf("%s-Boxes", i)))
				lar, _ := strconv.Atoi(req.FormValue(fmt.Sprintf("%s-LAR", i)))
				har, _ := strconv.Atoi(req.FormValue(fmt.Sprintf("%s-HAR", i)))

				fmt.Println(name, boxes, lar, har)

				newHL[name] = &oneroll.Location{
					Name:   name,
					Boxes:  boxes,
					LAR:    lar,
					HAR:    har,
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
					oneroll.ToSnakeCase(c.Name),
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
