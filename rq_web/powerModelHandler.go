package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/thewhitetulip/Tasks/sessions"
	"github.com/toferc/oneroll"
	"github.com/toferc/ore_web_roller/database"
	"github.com/toferc/ore_web_roller/models"
)

// PowerListHandler renders a character in a Web page
func PowerListHandler(w http.ResponseWriter, req *http.Request) {

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

	pows, err := database.ListPowerModels(db)
	if err != nil {
		panic(err)
	}

	wc := WebChar{
		CharacterModel: cm,
		IsAuthor:       IsAuthor,
		SessionUser:    username,
		IsLoggedIn:     loggedIn,
		IsAdmin:        isAdmin,
		PowerModels:    pows,
	}

	if req.Method == "GET" {

		// Render page

		Render(w, "templates/add_power_from_list.html", wc)

	}

	if req.Method == "POST" {

		err := req.ParseForm()
		if err != nil {
			panic(err)
		}

		if c.Powers == nil {
			c.Powers = map[string]*oneroll.Power{}
		}

		pName := req.FormValue("Name")

		nd, _ := strconv.Atoi(req.FormValue("Normal"))
		hd, _ := strconv.Atoi(req.FormValue("Hard"))
		wd, _ := strconv.Atoi(req.FormValue("Wiggle"))

		if pName != "" {
			p := pows[oneroll.ToSnakeCase(pName)].Power

			fmt.Println(p)

			p.Dice.Normal = nd
			p.Dice.Hard = hd
			p.Dice.Wiggle = wd

			p.Slug = oneroll.ToSnakeCase(pName)

			oneroll.UpdateCost(p)

			c.Powers[p.Slug] = p
		}
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

// AddStandalonePowerHandler renders a character in a Web page
func AddStandalonePowerHandler(w http.ResponseWriter, req *http.Request) {

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

	// Create default Power to populate page
	defaultPower := &oneroll.Power{
		Name: "",
		Dice: &oneroll.DiePool{
			Normal: 0,
			Hard:   0,
			Wiggle: 0,
		},
		Effect:    "",
		Qualities: []*oneroll.Quality{},
		Slug:      "",
	}

	// Map default Power to Character.Powers
	pm := models.PowerModel{
		Power: defaultPower,
	}

	wc := WebChar{
		PowerModel:  &pm,
		IsAuthor:    true,
		SessionUser: username,
		IsLoggedIn:  loggedIn,
		IsAdmin:     isAdmin,
		Modifiers:   oneroll.Modifiers,
		Counter:     []int{1, 2, 3, 4, 5, 6, 7, 8},
		Capacities: map[string]float32{
			"Mass":  25.0,
			"Range": 10.0,
			"Speed": 2.5,
			"Self":  0.0,
		},
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/add_standalone_power.html", wc)

	}

	if req.Method == "POST" { // POST

		err := req.ParseForm()
		if err != nil {
			panic(err)
		}

		pName := req.FormValue("Name")

		nd, _ := strconv.Atoi(req.FormValue("Normal"))
		hd, _ := strconv.Atoi(req.FormValue("Hard"))
		wd, _ := strconv.Atoi(req.FormValue("Wiggle"))

		p := oneroll.Power{
			Name: pName,
			Dice: &oneroll.DiePool{
				Normal: nd,
				Hard:   hd,
				Wiggle: wd,
			},
			Effect:    req.FormValue("Effect"),
			Qualities: []*oneroll.Quality{},
			Slug:      oneroll.ToSnakeCase(pName),
		}

		for _, qLoop := range wc.Counter[:4] { // Quality Loop

			qType := req.FormValue(fmt.Sprintf("Q%d-Type", qLoop))

			if qType != "" {
				l, err := strconv.Atoi(req.FormValue(fmt.Sprintf("Q%d-Level", qLoop)))
				if err != nil {
					l = 0
				}
				q := &oneroll.Quality{
					Type:  req.FormValue(fmt.Sprintf("Q%d-Type", qLoop)),
					Level: l,
					Name:  req.FormValue(fmt.Sprintf("Q%d-Name", qLoop)),
				}

				for _, cLoop := range wc.Counter[:3] {
					cType := req.FormValue(fmt.Sprintf("Q%d-C%d-Type", qLoop, cLoop))
					if cType != "" {
						cap := &oneroll.Capacity{
							Type: cType,
						}
						q.Capacities = append(q.Capacities, cap)
					}
				}

				m := new(oneroll.Modifier)

				for _, mLoop := range wc.Counter { // Modifier Loop
					mName := req.FormValue(fmt.Sprintf("Q%d-M%d-Name", qLoop, mLoop))
					if mName != "" {

						tM := oneroll.Modifiers[mName]

						m = &tM

						if m.RequiresLevel {
							l, err := strconv.Atoi(req.FormValue(fmt.Sprintf("Q%d-M%d-Level", qLoop, mLoop)))
							if err != nil {
								l = 1
							}
							m.Level = l
						}

						if m.RequiresInfo {
							m.Info = req.FormValue(fmt.Sprintf("Q%d-M%d-Info", qLoop, mLoop))
						}
						q.Modifiers = append(q.Modifiers, m)
					}
				}
				p.Qualities = append(p.Qualities, q)
			}
		}

		author := database.LoadUser(db, username)
		fmt.Println(author)

		pm.Author = author

		// Insert power into App archive
		p.DeterminePowerCapacities()
		p.CalculateCost()

		pm.Power = &p
		database.SavePowerModel(db, &pm)

		url := fmt.Sprintf("/view_power/%d", pm.ID)

		http.Redirect(w, req, url, http.StatusFound)
	}
}

// ModifyStandalonePowerHandler renders an editable Power in a Web page
func ModifyStandalonePowerHandler(w http.ResponseWriter, req *http.Request) {

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

	pm, err := database.PKLoadPowerModel(db, int64(id))
	if err != nil {
		fmt.Println(err)
	}

	if pm.Author == nil {
		pm.Author = &models.User{
			UserName: "",
		}
	}

	// Validate that User == Author
	IsAuthor := false

	if username == pm.Author.UserName {
		IsAuthor = true
	} else {
		http.Redirect(w, req, "/", 302)
	}

	p := pm.Power

	// Assign additional empty Qualities to populate form
	if len(p.Qualities) < 4 {
		for i := len(p.Qualities); i < 4; i++ {
			tempQ := oneroll.NewQuality("")
			p.Qualities = append(p.Qualities, tempQ)
		}
	} else {
		// Always create at least 2 Qualities
		for i := 0; i < 2; i++ {
			tempQ := oneroll.NewQuality("")
			p.Qualities = append(p.Qualities, tempQ)
		}
	}

	// Counter for form reading
	qCount := len(p.Qualities)

	// Assign additional empty Capacities to populate form
	for _, q := range p.Qualities {
		if len(q.Capacities) < 3 {
			for i := len(q.Capacities); i < 3; i++ {
				tempC := oneroll.Capacity{
					Type: "",
				}
				q.Capacities = append(q.Capacities, &tempC)
			}
		}
		if len(q.Modifiers) < 8 {
			for i := len(q.Modifiers); i < 8; i++ {
				tempM := oneroll.NewModifier("")
				q.Modifiers = append(q.Modifiers, tempM)
			}
		}
	}

	wc := WebChar{
		PowerModel:  pm,
		IsAuthor:    IsAuthor,
		SessionUser: username,
		IsLoggedIn:  loggedIn,
		IsAdmin:     isAdmin,
		Modifiers:   oneroll.Modifiers,
		Counter:     []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		Capacities: map[string]float32{
			"Mass":  25.0,
			"Range": 10.0,
			"Speed": 2.5,
			"Self":  0.0,
		},
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/modify_standalone_power.html", wc)

	}

	if req.Method == "POST" { // POST

		err := req.ParseForm()
		if err != nil {
			panic(err)
		}

		pName := req.FormValue("Name")

		nd, _ := strconv.Atoi(req.FormValue("Normal"))
		hd, _ := strconv.Atoi(req.FormValue("Hard"))
		wd, _ := strconv.Atoi(req.FormValue("Wiggle"))

		p.Name = pName
		p.Dice.Normal = nd
		p.Dice.Hard = hd
		p.Dice.Wiggle = wd
		p.Effect = req.FormValue("Effect")
		p.Qualities = []*oneroll.Quality{}
		p.Slug = oneroll.ToSnakeCase(pName)

		for _, qLoop := range wc.Counter[:qCount] { // Quality Loop

			qType := req.FormValue(fmt.Sprintf("Q%d-Type", qLoop))

			if qType != "" {
				l, err := strconv.Atoi(req.FormValue(fmt.Sprintf("Q%d-Level", qLoop)))
				if err != nil {
					l = 0
				}
				q := &oneroll.Quality{
					Type:  req.FormValue(fmt.Sprintf("Q%d-Type", qLoop)),
					Level: l,
					Name:  req.FormValue(fmt.Sprintf("Q%d-Name", qLoop)),
				}

				for _, cLoop := range wc.Counter[:3] {
					cType := req.FormValue(fmt.Sprintf("Q%d-C%d-Type", qLoop, cLoop))
					if cType != "" {
						cap := &oneroll.Capacity{
							Type: cType,
						}
						q.Capacities = append(q.Capacities, cap)
					}
				}

				m := new(oneroll.Modifier)

				for _, mLoop := range wc.Counter { // Modifier Loop
					mName := req.FormValue(fmt.Sprintf("Q%d-M%d-Name", qLoop, mLoop))
					if mName != "" {

						// Take base modifier struct from Modifiers
						tM := oneroll.Modifiers[mName]

						m = &tM

						if m.RequiresLevel {
							// Ensure level is a number or set to 1
							l, err := strconv.Atoi(req.FormValue(fmt.Sprintf("Q%d-M%d-Level", qLoop, mLoop)))
							if err != nil {
								l = 1
							}
							m.Level = l
						}

						if m.RequiresInfo {
							m.Info = req.FormValue(fmt.Sprintf("Q%d-M%d-Info", qLoop, mLoop))
						}
						// Append new modifier to Quality Modifiers
						q.Modifiers = append(q.Modifiers, m)

					}
				}
				// Append Quality to Power Qualities
				p.Qualities = append(p.Qualities, q)
			}
		}

		// Insert power into App archive
		p.DeterminePowerCapacities()
		p.CalculateCost()

		pm.Power = p

		err = database.UpdatePowerModel(db, pm)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("Saved")
		}

		fmt.Println(p)

		url := fmt.Sprintf("/view_power/%d", pm.ID)

		http.Redirect(w, req, url, http.StatusSeeOther)
	}
}

// DeleteStandalonePowerHandler renders a character in a Web page
func DeleteStandalonePowerHandler(w http.ResponseWriter, req *http.Request) {

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

	pm, err := database.PKLoadPowerModel(db, int64(id))
	if err != nil {
		fmt.Println(err)
	}

	// Validate that User == Author
	IsAuthor := false

	if username == pm.Author.UserName {
		IsAuthor = true
	} else {
		http.Redirect(w, req, "/", 302)
	}

	wc := WebChar{
		PowerModel:  pm,
		IsAuthor:    IsAuthor,
		SessionUser: username,
		IsLoggedIn:  loggedIn,
		IsAdmin:     isAdmin,
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/delete_standalone_power.html", wc)

	}

	if req.Method == "POST" {

		err := database.DeletePowerModel(db, pm.ID)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("Deleted Power")
		}

		url := fmt.Sprint("/index_powers/")

		http.Redirect(w, req, url, http.StatusSeeOther)
	}
}
