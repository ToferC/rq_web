package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

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

	// Create Custom Weapon Slot
	cMelee := &runequest.Attack{
		Skill: &runequest.Skill{},
		Weapon: &runequest.Weapon{
			Custom: true,
			HP:     8,
			SR:     3,
			Damage: "1d6+1",
		},
	}
	c.MeleeAttacks["Custom-Melee"] = cMelee

	cRanged := &runequest.Attack{
		Skill: &runequest.Skill{},
		Weapon: &runequest.Weapon{
			Custom: true,
			Range:  50,
			HP:     5,
			SR:     1,
			Damage: "1d6+1",
		},
	}
	c.RangedAttacks["Custom-Ranged"] = cRanged

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

		damBonus := c.Attributes["DB"]
		dbString := ""
		throwDB := ""

		if c.Attributes["DB"].Text != "-" {
			dbString = damBonus.Text

			if damBonus.Base > 0 {
				throwDB = fmt.Sprintf("+%dD%d", damBonus.Dice, damBonus.Base/2)
			} else {
				throwDB = fmt.Sprintf("-%dD%d", damBonus.Dice, damBonus.Base/2)
			}
		}

		tempMelee := map[string]*runequest.Attack{}

		for k, v := range c.MeleeAttacks {

			if !v.Weapon.Custom {

				// Regular Weapon from game
				weaponString := req.FormValue(fmt.Sprintf("Melee-Weapon-%s", k))
				skillString := req.FormValue(fmt.Sprintf("Melee-Skill-%s", k))

				if weaponString != "" && skillString != "" {

					// Convert weapon name to index
					weaponIndex := weaponsMap[weaponString]

					// Select weapon object from array
					weapon := baseWeapons[weaponIndex]

					tempMelee[weapon.Name] = &runequest.Attack{
						Name:         weapon.Name,
						Skill:        c.Skills[skillString],
						DamageString: weapon.Damage + dbString,
						StrikeRank:   c.Attributes["DEXSR"].Base + c.Attributes["SIZSR"].Base + weapon.SR,
						Weapon:       weapon,
					}
				}
			} else {
				// Custom Weapon
				name := req.FormValue(fmt.Sprintf("Custom-M-%s-Name", k))
				skill := req.FormValue(fmt.Sprintf("Custom-M-%s-Skill", k))

				if name != "" && skill != "" {

					_, ok := weaponsMap[name]
					if ok {
						// Name is already in common weapons
						name += " (custom)"
					}

					damage := req.FormValue(fmt.Sprintf("Custom-M-%s-Damage", k))
					special := req.FormValue(fmt.Sprintf("Custom-M-%s-Special", k))

					str := req.FormValue(fmt.Sprintf("Custom-M-%s-HP", k))
					hp, err := strconv.Atoi(str)
					if err != nil {
						hp = 0
					}

					str = req.FormValue(fmt.Sprintf("Custom-M-%s-SR", k))
					sr, err := strconv.Atoi(str)
					if err != nil {
						sr = 0
					}

					tempMelee[name] = &runequest.Attack{
						Name:         name,
						Skill:        c.Skills[skill],
						DamageString: damage + dbString,
						StrikeRank:   c.Attributes["DEXSR"].Base + c.Attributes["SIZSR"].Base + sr,
						Weapon: &runequest.Weapon{
							Name:      name,
							Type:      "Melee",
							STRDamage: true,
							Damage:    damage,
							HP:        hp,
							CurrentHP: hp,
							SR:        sr,
							Special:   special,
							Custom:    true,
						},
					}
				}
			}
		}

		c.MeleeAttacks = tempMelee

		// Ranged Weapons & Attacks
		tempRanged := map[string]*runequest.Attack{}

		for k, v := range c.RangedAttacks {

			if !v.Weapon.Custom {
				// Regular Weapon
				weaponString := req.FormValue(fmt.Sprintf("Ranged-Weapon-%s", k))
				skillString := req.FormValue(fmt.Sprintf("Ranged-Skill-%s", k))

				if weaponString != "" && skillString != "" {

					weaponIndex := weaponsMap[weaponString]

					weapon := baseWeapons[weaponIndex]

					// Set up for thrown weapons
					throw := false

					if strings.Contains(weapon.Name, "Thrown") {
						throw = true
					}

					damage := ""

					if weapon.Thrown {
						damage = weapon.Damage + throwDB
					} else {
						damage = weapon.Damage
					}

					// Ranged weapon
					tempRanged[weapon.Name] = &runequest.Attack{
						Name:         weapon.Name,
						Skill:        c.Skills[skillString],
						DamageString: damage,
						StrikeRank:   c.Attributes["DEXSR"].Base,
						Weapon:       weapon,
					}
					tempRanged[weapon.Name].Weapon.Thrown = throw
				}
			} else {
				// Custom Weapon
				name := req.FormValue(fmt.Sprintf("Custom-R-%s-Name", k))
				skill := req.FormValue(fmt.Sprintf("Custom-R-%s-Skill", k))

				if name != "" && skill != "" {

					_, ok := weaponsMap[name]
					if ok {
						// Name is already in common weapons
						name += " (custom)"
					}

					damage := req.FormValue(fmt.Sprintf("Custom-R-%s-Damage", k))
					special := req.FormValue(fmt.Sprintf("Custom-R-%s-Special", k))

					str := req.FormValue(fmt.Sprintf("Custom-R-%s-Range", k))
					rangeM, err := strconv.Atoi(str)
					if err != nil {
						rangeM = 0
					}

					str = req.FormValue(fmt.Sprintf("Custom-R-%s-HP", k))
					hp, err := strconv.Atoi(str)
					if err != nil {
						hp = 0
					}

					str = req.FormValue(fmt.Sprintf("Custom-R-%s-SR", k))
					sr, err := strconv.Atoi(str)
					if err != nil {
						sr = 0
					}

					tempRanged[name] = &runequest.Attack{
						Name:         name,
						Skill:        c.Skills[skill],
						DamageString: damage,
						StrikeRank:   c.Attributes["DEXSR"].Base + sr,
						Weapon: &runequest.Weapon{
							Name:      name,
							Type:      "Ranged",
							STRDamage: false,
							Damage:    damage,
							Range:     rangeM,
							HP:        hp,
							CurrentHP: hp,
							SR:        sr,
							Special:   special,
							Custom:    true,
						},
					}
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

		url := fmt.Sprintf("/view_character/%d#gameplay", cm.ID)

		http.Redirect(w, req, url, http.StatusSeeOther)
	}
}
