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

// EquipWeaponsArmorHandler assigns attacks and armor
func EquipWeaponsArmorHandler(w http.ResponseWriter, req *http.Request) {

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

	// Set map for recalling weapons later
	weaponsMap := map[string]int{}

	for i, w := range baseWeapons {
		weaponsMap[w.Name] = i
	}

	// Create maps if needed
	if c.MeleeAttacks == nil {
		c.MeleeAttacks = map[string]*runequest.Attack{}
	}

	if c.RangedAttacks == nil {
		c.RangedAttacks = map[string]*runequest.Attack{}
	}

	// Create empty attack slots if needed
	if len(c.MeleeAttacks) < 6 {
		for i := 1; i < 7-len(c.MeleeAttacks); i++ {
			a := &runequest.Attack{
				Skill:  &runequest.Skill{},
				Weapon: &runequest.Weapon{},
			}
			c.MeleeAttacks[fmt.Sprintf("%d", i)] = a
		}
	} else {
		// Always create at least 3
		for i := 1; i < 4; i++ {
			a := &runequest.Attack{
				Skill:  &runequest.Skill{},
				Weapon: &runequest.Weapon{},
			}
			c.MeleeAttacks[fmt.Sprintf("%d", i)] = a
		}
	}

	if len(c.RangedAttacks) < 6 {
		for i := 1; i < 7-len(c.RangedAttacks); i++ {
			a := &runequest.Attack{
				Skill:  &runequest.Skill{},
				Weapon: &runequest.Weapon{},
			}
			c.RangedAttacks[fmt.Sprintf("%d", i)] = a
		}
	} else {
		// Always create at least 3
		for i := 1; i < 4; i++ {
			a := &runequest.Attack{
				Skill:  &runequest.Skill{},
				Weapon: &runequest.Weapon{},
			}
			c.RangedAttacks[fmt.Sprintf("%d", i)] = a
		}
	}

	wc := WebChar{
		CharacterModel: cm,
		SessionUser:    username,
		IsLoggedIn:     loggedIn,
		IsAdmin:        isAdmin,
		IsAuthor:       IsAuthor,
		Counter:        numToArray(4),
		BigCounter:     numToArray(6),
		Skills:         runequest.Skills,
		MeleeAttacks:   c.MeleeAttacks,
		RangedAttacks:  c.RangedAttacks,
		BaseWeapons:    baseWeapons,
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/weapons_armor.html", wc)

	}

	if req.Method == "POST" {

		err := req.ParseMultipartForm(MaxMemory)
		if err != nil {
			panic(err)
		}

		// Weapons & Attacks
		// Working - pull map info from form or convert attacks to array?

		tempMelee := map[string]*runequest.Attack{}

		for k := range c.MeleeAttacks {
			weaponString := req.FormValue(fmt.Sprintf("Melee-Weapon-%s", k))
			skillString := req.FormValue(fmt.Sprintf("Melee-Skill-%s", k))

			if weaponString != "" && skillString != "" {

				// Convert weapon name to index
				weaponIndex := weaponsMap[weaponString]

				// Select weapon object from array
				weapon := baseWeapons[weaponIndex]

				dbString := ""

				if c.Attributes["DB"].Text != "-" {
					dbString = c.Attributes["DB"].Text
				}

				tempMelee[weapon.Name] = &runequest.Attack{
					Name:         weapon.Name,
					Skill:        c.Skills[skillString],
					DamageString: weapon.Damage + dbString,
					StrikeRank:   c.Attributes["DEXSR"].Base + c.Attributes["SIZSR"].Base + weapon.SR,
					Weapon:       weapon,
				}
			}
		}

		c.MeleeAttacks = tempMelee

		tempRanged := map[string]*runequest.Attack{}

		// Ranged Weapons & Attacks

		for k := range c.RangedAttacks {
			weaponString := req.FormValue(fmt.Sprintf("Ranged-Weapon-%s", k))
			skillString := req.FormValue(fmt.Sprintf("Ranged-Skill-%s", k))

			if weaponString != "" && skillString != "" {

				weaponIndex := weaponsMap[weaponString]

				weapon := baseWeapons[weaponIndex]

				// Ranged weapon
				tempRanged[weapon.Name] = &runequest.Attack{
					Name:         weapon.Name,
					Skill:        c.Skills[skillString],
					DamageString: weapon.Damage,
					StrikeRank:   c.Attributes["DEXSR"].Base,
					Weapon:       weapon,
				}
			}
		}

		c.RangedAttacks = tempRanged

		// Armor
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
