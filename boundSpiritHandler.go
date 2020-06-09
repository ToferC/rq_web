package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/thewhitetulip/Tasks/sessions"
	"github.com/toferc/rq_web/database"
	"github.com/toferc/runequest"
)

// BoundSpiritHandler renders a character in a Web page
func BoundSpiritHandler(w http.ResponseWriter, req *http.Request) {

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

	// Add Spirit Magic Slots to existing bound spirits
	for _, bs := range c.BoundSpirits {

		// Add extra spirit magic
		for i := 1; i < 4; i++ {
			tsm := runequest.Spell{
				Name:       "",
				CoreString: "",
				UserString: "",
				Cost:       0,
			}
			tsm.Name = createName(tsm.CoreString, tsm.UserString)
			bs.SpiritMagicSpells["zzNewSM-"+string(i)] = tsm
		}
	}

	bLen := len(c.BoundSpirits)

	// Add Bound Spirits
	for i := bLen; i < bLen+3; i++ {
		bs := &runequest.BoundSpirit{
			Name:              "",
			Item:              "",
			Pow:               0,
			Cha:               0,
			CurrentMP:         0,
			SpiritMagicSpells: map[string]runequest.Spell{},
		}

		// Add extra spirit magic
		for i := 1; i < 5; i++ {
			tsm := runequest.Spell{
				Name:       "",
				CoreString: "",
				UserString: "",
				Cost:       0,
			}
			tsm.Name = createName(tsm.CoreString, tsm.UserString)
			bs.SpiritMagicSpells["zzNewSM-"+string(i)] = tsm
		}

		c.BoundSpirits = append(c.BoundSpirits, bs)
	}

	wc := WebChar{
		CharacterModel: cm,
		SessionUser:    username,
		IsLoggedIn:     loggedIn,
		IsAdmin:        isAdmin,
		IsAuthor:       IsAuthor,
		SpiritMagic:    runequest.SpiritMagicSpells,
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/bound_spirit.html", wc)

	}

	if req.Method == "POST" {

		err := req.ParseMultipartForm(MaxMemory)
		if err != nil {
			panic(err)
		}

		// Set temp map
		tempSpirits := []*runequest.BoundSpirit{}

		// Read form
		for i, bs := range c.BoundSpirits {
			name := req.FormValue(fmt.Sprintf("BS-%d-Name", i))
			description := req.FormValue(fmt.Sprintf("BS-%d-Description", i))
			item := req.FormValue(fmt.Sprintf("BS-%d-Item", i))

			if name != "" {

				bs.Name = name
				bs.Description = description
				bs.Item = ProcessUserString(item)

				powStr := req.FormValue(fmt.Sprintf("BS-%d-Pow", i))

				pow, err := strconv.Atoi(powStr)
				if err != nil {
					pow = 0
					fmt.Println("Non-number entered")
				}

				chaStr := req.FormValue(fmt.Sprintf("BS-%d-Cha", i))

				cha, err := strconv.Atoi(chaStr)
				if err != nil {
					cha = 0
					fmt.Println("Non-number entered")
				}

				intStr := req.FormValue(fmt.Sprintf("BS-%d-Int", i))

				intelligence, err := strconv.Atoi(intStr)
				if err != nil {
					intelligence = 0
					fmt.Println("Non-number entered")
				}

				bs.Pow = pow
				bs.Int = intelligence
				bs.Cha = cha
				bs.CurrentMP = pow

				// Spirit Magic

				tempSpiritMagic := map[string]runequest.Spell{}

				for k, v := range bs.SpiritMagicSpells {
					str := req.FormValue(fmt.Sprintf("BS-%d-SpiritMagic-%s", i, k))
					spec := req.FormValue(fmt.Sprintf("BS-%d-SpiritMagic-%s-UserString", i, k))
					cString := req.FormValue(fmt.Sprintf("BS-%d-SpiritMagic-%s-Cost", i, k))

					cost, err := strconv.Atoi(cString)
					if err != nil {
						cost = 1
						fmt.Println("Non-number entered")
					}

					switch {
					case str == "":
						fmt.Println("No spell")
						continue

					case str != "" && v.CoreString == str && v.UserString == spec && v.Cost == cost:
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

						if spec != "" {
							spec = ProcessUserString(spec)
						} else {
							spec = baseSpell.UserString
						}

						s := runequest.Spell{
							Name:       baseSpell.Name,
							CoreString: baseSpell.CoreString,
							UserString: spec,
							Cost:       cost,
							Domain:     baseSpell.Domain,
						}

						s.GenerateName()
						tempSpiritMagic[s.Name] = s
					}
				}

				bs.SpiritMagicSpells = tempSpiritMagic

				tempSpirits = append(tempSpirits, bs)
			}
		}

		// Swap map
		c.BoundSpirits = tempSpirits

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
