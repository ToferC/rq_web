package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/thewhitetulip/Tasks/sessions"
	"github.com/toferc/rq_web/database"
	"github.com/toferc/rq_web/models"
	"github.com/toferc/runequest"
)

// RandomCharacterHandler allows users to name and select a homeland
func RandomCharacterHandler(w http.ResponseWriter, req *http.Request) {

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

	c := runequest.NewCharacter("")

	author, err := database.LoadUser(db, username)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, req, "/", 302)
	}

	fmt.Println(author)

	cm = models.CharacterModel{
		Character: c,
		Author:    author,
	}

	homelands, err := database.ListHomelandModels(db)
	if err != nil {
		panic(err)
	}

	occupations, err := database.ListOccupationModels(db)
	if err != nil {
		panic(err)
	}

	cults, err := database.ListCultModels(db)
	if err != nil {
		panic(err)
	}

	wc := WebChar{
		CharacterModel:   &cm,
		SessionUser:      username,
		IsLoggedIn:       loggedIn,
		IsAdmin:          isAdmin,
		HomelandModels:   homelands,
		OccupationModels: occupations,
		CultModels:       cults,
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/random_character.html", wc)

	}

	if req.Method == "POST" {

		err := req.ParseMultipartForm(MaxMemory)
		if err != nil {
			panic(err)
		}

		// Get character scale
		scale := req.FormValue("Scale")

		// Set Homeland
		hlStr := req.FormValue("HLStr")

		fmt.Println("Results: " + hlStr)

		hlID, err := strconv.Atoi(hlStr)
		if err != nil {
			for _, k := range homelands {
				// Get first homeland
				hlID = int(k.ID)
			}
			fmt.Println(err)
		}

		hl, err := database.PKLoadHomelandModel(db, int64(hlID))
		if err != nil {
			fmt.Println(err)
		}

		c.Homeland = hl.Homeland
		fmt.Println("HOMELAND: " + c.Homeland.Name)

		// Load MarkovChains
		sartarChain, err := loadModel("sartarModel.json")
		if err != nil {
			log.Println(err)
		}

		// Name generation
		c.Name = generateName(sartarChain)

		// Traits generation

		traits := readCSV("traits.csv")

		t1 := traits[ChooseRandom(len(traits))]
		t2 := traits[ChooseRandom(len(traits))]

		c.Description = fmt.Sprintf("%s is %s and %s.", c.Name, t1, t2)

		// Set Occupation
		ocStr := req.FormValue("OCStr")

		ocID, err := strconv.Atoi(ocStr)
		if err != nil {
			for _, v := range occupations {
				// Take first occupation in map
				ocID = int(v.ID)
				break
			}
		}

		oc, err := database.PKLoadOccupationModel(db, int64(ocID))
		if err != nil {
			fmt.Println("No Occupation Found")
		}

		c.Occupation = oc.Occupation
		fmt.Println("OCCUPATION: " + c.Occupation.Name)

		// Set Cult
		cStr := req.FormValue("CStr")

		cID, err := strconv.Atoi(cStr)
		if err != nil {
			for _, v := range cults {
				// Take first cult in map
				cID = int(v.ID)
				break
			}
		}

		cultModel, err := database.PKLoadCultModel(db, int64(cID))
		if err != nil {
			fmt.Println("No Cult Found")
		}

		if cultModel.Cult.SubCult {
			// Set base to ParentCult
			// Easier to add to
			parentCult := cultModel.Cult.ParentCult
			if parentCult != nil {
				// Add SubCult skills
				for _, s := range cultModel.Cult.Skills {
					parentCult.Skills = append(parentCult.Skills, s)
				}
				// Add SubCult weapons
				for _, w := range cultModel.Cult.Weapons {
					parentCult.Weapons = append(parentCult.Weapons, w)
				}
				// Add SubCult RuneSpells
				for _, rs := range cultModel.Cult.RuneSpells {
					parentCult.RuneSpells = append(parentCult.RuneSpells, rs)
				}
				// Add SubCult SpiritMagic
				for _, sm := range cultModel.Cult.SpiritMagic {
					parentCult.SpiritMagic = append(parentCult.SpiritMagic, sm)
				}

				// Set Details to ParentCult
				parentCult.Name = cultModel.Cult.Name
				parentCult.Description = cultModel.Cult.Description
				parentCult.Notes += "\n" + cultModel.Cult.Notes
				// Set cult to ParentCult
				c.Cult = parentCult
			}
		} else {
			c.Cult = cultModel.Cult
		}

		c.Cult.Rank = "Initiate"

		fmt.Println("CULT: " + c.Cult.Name)

		if c.Occupation.StandardOfLiving == "Free" {
			c.Abilities["Reputation"].OccupationValue = 5
		}

		if c.Occupation.StandardOfLiving == "Noble" {
			c.Abilities["Reputation"].OccupationValue = 10
		}

		cm.Image = new(models.Image)
		cm.Image.Path = DefaultCharacterPortrait

		// Update CreationSteps
		c.CreationSteps["Base Choices"] = true
		c.CreationSteps["Personal History"] = true

		c.Role = "Player Character"

		// Choose Runes
		runeValueArray := []int{}
		boostRunes := 0

		runesChosen := []string{}

		switch scale {
		case "Heroic":
			runeValueArray = []int{80, 60, 40, 85}
			boostRunes = 1
		case "Epic":
			runeValueArray = []int{90, 70, 60, 95}
			boostRunes = 2
		default:
			runeValueArray = []int{60, 40, 20, 75}
			boostRunes = 0
		}

		// Allocate one Rune bonuses to Cult runes first
		for i, r := range c.Cult.Runes {
			if i <= boostRunes {
				if isInString(runequest.ElementalRuneOrder, r) {
					c.ElementalRunes[r].Base = runeValueArray[0]
					runesChosen = append(runesChosen, r)
				}
			}
		}

		if len(runesChosen) == 0 {
			// No cult rune in Elemental Runes
			ru := runequest.ElementalRuneOrder[ChooseRandom(
				len(runequest.ElementalRuneOrder))]
			c.ElementalRunes[ru].Base = runeValueArray[0]
			runesChosen = append(runesChosen, ru)
		}

		// Id remaining runes
		remainingRunes := []string{}
		for _, r := range runequest.ElementalRuneOrder {
			if !isInString(runesChosen, r) {
				remainingRunes = append(remainingRunes, r)
			}
		}

		// Allocate 40 rune
		r := ChooseRandom(4)
		r40 := c.ElementalRunes[remainingRunes[r]]
		r40.Base = runeValueArray[1]
		runesChosen = append(runesChosen, r40.Name)

		// Id remaining runes
		remainingRunes = []string{}
		for _, r := range runequest.ElementalRuneOrder {
			if !isInString(runesChosen, r) {
				remainingRunes = append(remainingRunes, r)
			}
		}

		// Allocate 20 rune
		r = ChooseRandom(3)
		r20 := c.ElementalRunes[remainingRunes[r]]
		r20.Base = runeValueArray[2]

		// Elemental Runes
		// Allocate one Rune bonuses to Cult runes first
		runesChosen = []string{}
		for _, r := range c.Cult.Runes {
			if isInString(runequest.PowerRuneOrder[:9], r) {
				c.PowerRunes[r].Base = runeValueArray[3]
				c.PowerRunes[r].UpdateAbility()
				c.UpdateOpposedRune(c.PowerRunes[r])
				runesChosen = append(runesChosen, r)
				runesChosen = append(runesChosen, c.PowerRunes[r].OpposedAbility)
				break
			}
		}

		// Id remaining runes
		remainingRunes = []string{}
		for _, r := range runequest.PowerRuneOrder[:9] {
			if !isInString(runesChosen, r) {
				remainingRunes = append(remainingRunes, r)
			}
		}

		// Add second rune
		r = ChooseRandom(len(remainingRunes))
		r75 := c.PowerRunes[remainingRunes[r]]
		r75.Base = runeValueArray[3]

		c.CreationSteps["Rune Affinities"] = true

		// Set Stats
		minStat := 1

		// Define minimum die roll value
		switch scale {
		case "Common":
			minStat = 1
		case "Heroic":
			minStat = 2
		case "Epic":
			minStat = 3
		}

		// Roll stats
		for k, v := range cm.Character.Homeland.StatisticFrames {
			c.Statistics[k].Base = runequest.RollDice(6, minStat, v.Modifier, v.Dice)
			c.Statistics[k].Max = v.Max
		}

		// Determine Rune bonuses (returns array of 2 strings)
		runeArray := c.DetermineRuneModifiers()

		c.Statistics[runeArray[0]].RuneBonus = 2
		c.Statistics[runeArray[0]].Max += 2

		c.Statistics[runeArray[1]].RuneBonus = 1
		c.Statistics[runeArray[1]].Max++

		c.LocationForm = c.Homeland.LocationForm
		c.HitLocations = runequest.LocationForms[c.Homeland.LocationForm]
		c.HitLocationMap = runequest.SortLocations(c.HitLocations)

		c.SetAttributes()

		// Apply Move
		for _, m := range c.Homeland.Movement {
			c.Movement = append(c.Movement,
				&runequest.Movement{
					Name:  m.Name,
					Value: m.Value,
				})
		}

		// Update Character
		c.CurrentHP = c.Attributes["HP"].Max
		c.CurrentMP = c.Attributes["MP"].Max
		c.CurrentRP = c.Cult.NumRunePoints

		for _, v := range c.HitLocations {
			v.Value = v.Max
		}

		c.CreationSteps["Roll Stats"] = true

		// Update skills now that we have stats
		c.Skills["Dodge"].Base = c.Statistics["DEX"].Total * 2
		c.Skills["Jump"].Base = c.Statistics["DEX"].Total * 3

		// Add choices to c.Homeland.Skills
		for _, sc := range c.Homeland.SkillChoices {
			choice := ChooseRandom(2)
			c.Homeland.Skills = append(c.Homeland.Skills, &sc.Skills[choice])

		}

		// Add Skills to Character
		fmt.Println("**Homeland Skills**")
		for _, s := range c.Homeland.Skills {
			// Modify Skill
			if s.UserString == "any" {
				// use default specialization
				// Could add list and choose random here

				s.UserString = runequest.Skills[s.CoreString].UserString
			}

			targetString := createName(s.CoreString, s.UserString)

			sk, ok := c.Skills[targetString]

			if ok {
				// Skill exists in Character, modify it via pointer
				sk.HomelandValue = s.HomelandValue

				if s.Base != sk.Base {
					sk.Base = s.Base
				}

				sk.Name = createName(sk.CoreString, sk.UserString)
				sk.UpdateSkill()

				fmt.Printf("Updated Character Skill: %s: %d\n", sk.Name, s.Base)

			} else {
				// We need to find the base skill from the Master list or create it
				bs, ok := runequest.Skills[s.CoreString]
				if !ok {
					fmt.Println("Skill is new: " + targetString)

					// New Skill
					baseSkill := &runequest.Skill{
						CoreString:    s.CoreString,
						UserString:    s.UserString,
						Category:      s.Category,
						Base:          s.Base,
						UserChoice:    s.UserChoice,
						HomelandValue: s.HomelandValue,
					}

					// Update our new skill

					if baseSkill.Category == "" {
						baseSkill.Category = "Agility"
					}

					sc, ok := c.SkillCategories[baseSkill.Category]
					if !ok {
						// Figure this one out
						sc = c.SkillCategories["Perception"]
					}

					baseSkill.CategoryValue = sc.Value

					baseSkill.Name = createName(baseSkill.CoreString, baseSkill.UserString)
					baseSkill.UpdateSkill()

					fmt.Println("Add Skill to character: " + baseSkill.Name)
					c.Skills[baseSkill.Name] = baseSkill
				} else {
					// Skill exists in master list create new skill based off template
					fmt.Println("Skill in master list: " + targetString)

					fmt.Println(bs)

					baseSkill := &runequest.Skill{
						CoreString: bs.CoreString,
						Category:   bs.Category,
						Base:       s.Base,
						UserChoice: bs.UserChoice,
					}

					baseSkill.HomelandValue = s.HomelandValue
					fmt.Println(s.HomelandValue, baseSkill.HomelandValue)

					if s.Base > baseSkill.Base {
						baseSkill.Base = s.Base
					}
					if s.UserString != "" {
						baseSkill.UserString = s.UserString
					}

					baseSkill.Name = createName(baseSkill.CoreString, baseSkill.UserString)

					// Add Skill to Character
					fmt.Println("Add Skill to character: " + baseSkill.Name)
					c.Skills[baseSkill.Name] = baseSkill
					c.Skills[baseSkill.Name].UpdateSkill()
				}
			}
		}

		// Homelands grant 3 base passions
		fmt.Println("**Homeland Passions**")
		for _, a := range c.Homeland.PassionList {
			c.Homeland.Passions = append(c.Homeland.Passions, a)

			targetString := createName(a.CoreString, a.UserString)

			if c.Abilities[targetString] == nil {
				// Need to create a new ability
				a := &runequest.Ability{
					Type:       "Passion",
					CoreString: a.CoreString,
					UserChoice: a.UserChoice,
					UserString: a.UserString,
					Updates:    []*runequest.Update{},
				}

				update := CreateUpdate("Homeland", 60)
				a.Updates = append(a.Updates, update)

				a.UpdateAbility()
				c.Abilities[targetString] = a
			} else {
				// Update existing ability

				update := CreateUpdate("Homeland", 10)
				c.Abilities[targetString].Updates = append(c.Abilities[targetString].Updates, update)

				c.Abilities[targetString].UpdateAbility()
			}
		}

		// Homeland grants a bonus to a rune affinity
		if isInString(runequest.ElementalRuneOrder, c.Homeland.RuneBonus) {
			c.ElementalRunes[c.Homeland.RuneBonus].HomelandValue += 10
		}

		if isInString(runequest.PowerRuneOrder, c.Homeland.RuneBonus) {
			c.PowerRunes[c.Homeland.RuneBonus].HomelandValue += 10
		}

		c.CreationSteps["Apply Homeland"] = true

		// Occupation

		fmt.Println("**Occupation**")

		// Add choices to c.Occupation.Skills
		for _, sc := range c.Occupation.SkillChoices {
			choice := ChooseRandom(2)
			c.Occupation.Skills = append(c.Occupation.Skills, &sc.Skills[choice])
		}

		// Add choices for Weapon skills

		meleeSkills := []string{}
		rangedSkills := []string{}
		allSkills := []string{"Small Shield", "Medium Shield", "Large Shield"}
		shieldSkills := []string{"Small Shield", "Medium Shield", "Large Shield"}

		for k, v := range runequest.Skills {

			switch {
			case v.Category == "Melee":
				meleeSkills = append(meleeSkills, k)
				allSkills = append(allSkills, k)
			case v.Category == "Ranged":
				rangedSkills = append(rangedSkills, k)
				allSkills = append(allSkills, k)
			}
		}

		// Choose weapon skills
		var choice int

		for _, w := range c.Occupation.Weapons {

			target := ""

			switch {
			case w.Description == "Melee":
				choice = ChooseRandom(len(meleeSkills))
				target = meleeSkills[choice]
			case w.Description == "Ranged":
				choice = ChooseRandom(len(rangedSkills))
				target = rangedSkills[choice]
			case w.Description == "Shield":
				choice = ChooseRandom(len(shieldSkills))
				target = shieldSkills[choice]
			default:
				choice = ChooseRandom(len(allSkills))
				target = allSkills[choice]
			}

			bs := runequest.Skills[target]

			ws := &runequest.Skill{
				CoreString:      bs.CoreString,
				Category:        bs.Category,
				Base:            bs.Base,
				OccupationValue: w.Value,
			}

			c.Occupation.Skills = append(c.Occupation.Skills, ws)
		}

		// Add Skills to Character
		for _, s := range c.Occupation.Skills {
			// Modify Skill
			if s.UserString == "any" {
				// User Chooses a new specialization
				s.UserString = runequest.Skills[s.CoreString].UserString
			}

			targetString := createName(s.CoreString, s.UserString)

			sk, ok := c.Skills[targetString]

			if ok {
				// Skill exists in Character, modify it via pointer
				sk.OccupationValue = s.OccupationValue

				sk.Name = createName(sk.CoreString, sk.UserString)
				sk.UpdateSkill()

				fmt.Printf("Updated Character Skill: %s: %d\n", sk.Name, s.Base)

			} else {
				// We need to find the base skill from the Master list or create it
				bs, ok := runequest.Skills[s.CoreString]
				if !ok {
					fmt.Println("Skill is new: " + targetString)

					// New Skill
					baseSkill := &runequest.Skill{
						CoreString:      s.CoreString,
						UserString:      s.UserString,
						Category:        s.Category,
						Base:            s.Base,
						UserChoice:      s.UserChoice,
						OccupationValue: s.OccupationValue,
					}

					// Update our new skill
					sc := c.SkillCategories[baseSkill.Category]

					baseSkill.CategoryValue = sc.Value

					baseSkill.Name = createName(baseSkill.CoreString, baseSkill.UserString)
					baseSkill.UpdateSkill()

					fmt.Println("Add Skill to character: " + baseSkill.Name)
					c.Skills[baseSkill.Name] = baseSkill
				} else {
					// Skill exists in master list
					fmt.Println("Skill in master list: " + targetString)

					fmt.Println(bs)

					baseSkill := &runequest.Skill{
						CoreString: bs.CoreString,
						Category:   bs.Category,
						Base:       bs.Base,
						UserChoice: bs.UserChoice,
					}

					baseSkill.OccupationValue = s.OccupationValue
					fmt.Println(s.OccupationValue, baseSkill.OccupationValue)

					if s.Base > baseSkill.Base {
						baseSkill.Base = s.Base
					}
					if s.UserString != "" {
						baseSkill.UserString = s.UserString
					}

					baseSkill.Name = createName(baseSkill.CoreString, baseSkill.UserString)

					// Add Skill to Character
					fmt.Println("Add Skill to character: " + baseSkill.Name)
					c.Skills[baseSkill.Name] = baseSkill
					c.Skills[baseSkill.Name].UpdateSkill()
				}
			}
		}

		// Occupations grant a bonus to one Passion
		fmt.Println("**Occupations Passion**")
		if len(c.Occupation.PassionList) > 0 {

			n := ChooseRandom(len(c.Occupation.PassionList))

			a := c.Occupation.PassionList[n]

			targetString := createName(a.CoreString, a.UserString)

			if c.Abilities[targetString] == nil {
				// Need to create a new ability
				a := &runequest.Ability{
					Type:       "Passion",
					CoreString: a.CoreString,
					UserChoice: a.UserChoice,
					UserString: a.UserString,
					Updates:    []*runequest.Update{},
				}

				update := CreateUpdate("Occupation", 60)
				a.Updates = append(a.Updates, update)

				a.UpdateAbility()
				c.Abilities[targetString] = a
			} else {
				// Update existing ability

				update := CreateUpdate("Occupation", 10)
				c.Abilities[targetString].Updates = append(c.Abilities[targetString].Updates, update)

				c.Abilities[targetString].UpdateAbility()
			}
		}

		// Set generic armor value
		//armor := 0

		// Equipment
		if len(c.Occupation.Equipment) > 0 {
			for _, e := range c.Occupation.Equipment {
				c.Equipment = append(c.Equipment, e)
			}
		}

		c.Income = c.Occupation.Income

		// Update CreationSteps
		c.CreationSteps["Apply Occupation"] = true

		// Cults
		fmt.Println("**Cults**")

		switch scale {
		case "Heroic":
			c.Cult.Rank = "Initiate"
			c.Cult.NumRunePoints = 3
		case "Epic":
			c.Cult.Rank = "Rune Lord"
			c.Cult.NumRunePoints = 8
		default:
			c.Cult.Rank = "Initiate"
			c.Cult.NumRunePoints = 1
		}

		c.RuneSpells = map[string]*runequest.Spell{}
		c.SpiritMagic = map[string]*runequest.Spell{}

		// Rune Magic
		fmt.Println("**Choose Rune Spells**")

		numRuneSpells := 1

		switch scale {
		case "Heroic":
			numRuneSpells = 3
		case "Epic":
			numRuneSpells = 5
		}

		chosenInt := []int{}
		index := ChooseRandom(len(c.Cult.RuneSpells))

		for i := 1; i <= numRuneSpells; i++ {

			for isIn(chosenInt, index) {
				fmt.Println("Spell already chosen")
				index = ChooseRandom(len(c.Cult.RuneSpells))
			}

			baseSpell := c.Cult.RuneSpells[index]

			s := baseSpell

			s.GenerateName()
			c.RuneSpells[s.Name] = s

			chosenInt = append(chosenInt, index)

		}

		// Spirit Magic

		// Add Associated Cult spells
		totalSpiritMagic := []*runequest.Spell{}

		totalSpiritMagic = c.Cult.SpiritMagic

		for _, ac := range c.Cult.AssociatedCults {
			for _, spell := range ac.SpiritMagic {
				totalSpiritMagic = append(totalSpiritMagic, spell)
			}
		}

		mSP := 5

		switch scale {
		case "Heroic":
			mSP = 10
		case "Epic":
			mSP = c.Statistics["CHA"].Total
		default:
			mSP = 5
		}

		cSP := 0

		for {
			index := ChooseRandom(len(c.Cult.SpiritMagic))

			baseSpell := totalSpiritMagic[index]

			if baseSpell.Variable {
				pts := runequest.RollDice(mSP-cSP, 1, 0, 1)
				baseSpell.Cost = pts
			}

			baseSpell.GenerateName()
			c.SpiritMagic[baseSpell.Name] = baseSpell
			cSP += baseSpell.Cost
			if cSP >= mSP {
				break
			}
		}

		// Add choices to c.Cult.Skills
		for _, sc := range c.Cult.SkillChoices {
			m := ChooseRandom(2)
			c.Cult.Skills = append(c.Cult.Skills, &sc.Skills[m])
		}

		for _, w := range c.Cult.Weapons {

			target := ""

			switch {
			case w.Description == "Melee":
				choice = ChooseRandom(len(meleeSkills))
				target = meleeSkills[choice]
			case w.Description == "Ranged":
				choice = ChooseRandom(len(rangedSkills))
				target = rangedSkills[choice]
			case w.Description == "Shield":
				choice = ChooseRandom(len(shieldSkills))
				target = shieldSkills[choice]
			default:
				choice = ChooseRandom(len(allSkills))
				target = allSkills[choice]
			}

			bs := runequest.Skills[target]

			ws := &runequest.Skill{
				CoreString:      bs.CoreString,
				Category:        bs.Category,
				Base:            bs.Base,
				OccupationValue: w.Value,
			}

			c.Cult.Skills = append(c.Cult.Skills, ws)
		}

		// Add Skills to Character
		for _, s := range c.Cult.Skills {
			// Modify Skill
			if s.UserString == "any" {
				// User Chooses a new specialization

				s.UserString = runequest.Skills[s.CoreString].UserString
			}

			targetString := createName(s.CoreString, s.UserString)

			sk, ok := c.Skills[targetString]

			if ok {
				// Skill exists in Character, modify it via pointer
				sk.CultValue = s.CultValue

				sk.Name = createName(sk.CoreString, sk.UserString)
				sk.UpdateSkill()

				fmt.Printf("Updated Character Skill: %s: %d\n", sk.Name, s.Base)

			} else {
				// We need to find the base skill from the Master list or create it
				bs, ok := runequest.Skills[s.CoreString]
				if !ok {
					fmt.Println("Skill is new: " + targetString)

					// New Skill
					baseSkill := &runequest.Skill{
						CoreString: s.CoreString,
						UserString: s.UserString,
						Category:   s.Category,
						Base:       s.Base,
						UserChoice: s.UserChoice,
						CultValue:  s.CultValue,
					}

					// Update our new skill
					sc := c.SkillCategories[baseSkill.Category]

					baseSkill.CategoryValue = sc.Value

					baseSkill.Name = createName(baseSkill.CoreString, baseSkill.UserString)
					baseSkill.UpdateSkill()

					fmt.Println("Add Skill to character: " + baseSkill.Name)
					c.Skills[baseSkill.Name] = baseSkill
				} else {
					// Skill exists in master list
					fmt.Println("Skill in master list: " + targetString)

					fmt.Println(bs)

					baseSkill := &runequest.Skill{
						CoreString: bs.CoreString,
						Category:   bs.Category,
						Base:       bs.Base,
						UserChoice: bs.UserChoice,
					}

					baseSkill.CultValue = s.CultValue
					fmt.Println(s.CultValue, baseSkill.CultValue)

					if s.Base > baseSkill.Base {
						baseSkill.Base = s.Base
					}
					if s.UserString != "" {
						baseSkill.UserString = s.UserString
					}

					baseSkill.Name = createName(baseSkill.CoreString, baseSkill.UserString)

					// Add Skill to Character
					fmt.Println("Add Skill to character: " + baseSkill.Name)
					c.Skills[baseSkill.Name] = baseSkill
					c.Skills[baseSkill.Name].UpdateSkill()
				}
			}
		}

		fmt.Println("**Cult Bonuses**")
		// Add 20 to one cult skill
		// Skill exists in Character, modify it via pointer

		index = ChooseRandom(len(c.Cult.Skills))

		baseSkill := c.Cult.Skills[index]
		targetSkill := &runequest.Skill{
			CoreString: baseSkill.CoreString,
			UserString: baseSkill.UserString,
		}

		targetSkill.GenerateName()

		t := time.Now()
		tString := t.Format("2006-01-02 15:04:05")

		update := &runequest.Update{
			Date:  tString,
			Event: "Cult Skill (20)",
			Value: 20,
		}

		sk := c.Skills[targetSkill.Name]

		if sk.Updates == nil {
			sk.Updates = []*runequest.Update{}
		}

		sk.Updates = append(sk.Updates, update)

		sk.UpdateSkill()

		fmt.Println("Updated Character Skill 20%: " + sk.Name)

		// Add 15 to one cult skill
		// Skill exists in Character, modify it via pointer

		index = ChooseRandom(len(c.Cult.Skills))

		baseSkill2 := c.Cult.Skills[index]
		targetSkill2 := &runequest.Skill{
			CoreString: baseSkill2.CoreString,
			UserString: baseSkill2.UserString,
		}

		targetSkill2.GenerateName()

		update2 := &runequest.Update{
			Date:  tString,
			Event: "Cult Skill (15)",
			Value: 15,
		}

		sk2 := c.Skills[targetSkill2.Name]

		if sk2.Updates == nil {
			sk2.Updates = []*runequest.Update{}
		}

		sk2.Updates = append(sk2.Updates, update2)

		sk2.UpdateSkill()

		fmt.Println("Updated Character Skill 15%: " + sk2.Name)

		// Cults grant a bonus to one Passion
		if len(c.Cult.PassionList) > 0 {
			n := ChooseRandom(len(c.Cult.PassionList))
			a := c.Cult.PassionList[n]
			targetString := createName(a.CoreString, a.UserString)

			if c.Abilities[targetString] == nil {
				// Need to create a new ability
				a := &runequest.Ability{
					Type:       "Passion",
					CoreString: a.CoreString,
					UserChoice: a.UserChoice,
					UserString: a.UserString,
					Updates:    []*runequest.Update{},
				}

				update := CreateUpdate("Cult", 60)
				a.Updates = append(a.Updates, update)

				a.UpdateAbility()
				c.Abilities[targetString] = a
			} else {
				// Update existing ability

				update := CreateUpdate("Cult", 10)
				c.Abilities[targetString].Updates = append(c.Abilities[targetString].Updates, update)

				c.Abilities[targetString].UpdateAbility()
			}
		}

		// Add Spirit Magic
		c.Cult.NumSpiritMagic = 5

		if c.Occupation.Name == "Assistant Shaman" {
			c.Cult.NumSpiritMagic += 5
		}

		// Set Rune Points to Max
		c.CurrentRP = c.Cult.NumRunePoints

		// Update CreationSteps
		c.CreationSteps["Apply Cult"] = true

		// Personal Skills
		fmt.Println("**Personal Skills**")

		sortedSkillArray := sortedSkills(c.Skills)

		switch scale {
		case "Heroic":
			for i, ss := range sortedSkillArray {
				switch {
				case i <= 3:
					s := c.Skills[ss.Name]
					if s.CoreString == "Speak" {
						s = sortedSkillArray[ChooseRandom(len(sortedSkillArray))]
					}
					s.AddSkillUpdate("Heroic Skill 25", 25)
				case i <= 8:
					s := c.Skills[ss.Name]
					if s.CoreString == "Speak" {
						s = sortedSkillArray[ChooseRandom(len(sortedSkillArray))]
					}
					s.AddSkillUpdate("Heroic Skill 10", 10)
				default:
					break
				}
			}
		case "Epic":
			for i, ss := range sortedSkillArray {
				switch {
				case i <= 7:
					s := c.Skills[ss.Name]
					if s.CoreString == "Speak" {
						s = sortedSkillArray[ChooseRandom(len(sortedSkillArray))]
					}
					s.AddSkillUpdate("Epic Skill 40", 40)
				case i <= 14:
					s := c.Skills[ss.Name]
					if s.CoreString == "Speak" {
						s = sortedSkillArray[ChooseRandom(len(sortedSkillArray))]
					}
					s.AddSkillUpdate("Epic Skill 25", 25)
				case i <= 21:
					s := c.Skills[ss.Name]
					if s.CoreString == "Speak" {
						s = sortedSkillArray[ChooseRandom(len(sortedSkillArray))]
					}
					s.AddSkillUpdate("Epic Skill 10", 10)
				default:
					break
				}
			}
		default:
			// Common Character
			for i, ss := range sortedSkillArray {
				if i <= 5 {
					s := c.Skills[ss.Name]
					if s.CoreString == "Speak" {
						s = sortedSkillArray[ChooseRandom(len(sortedSkillArray))]
					}
					s.AddSkillUpdate("Heroic Skill 10", 10)
				}
			}
		}

		// Weapons

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

		// Get best weapons skill
		for _, ss := range sortedSkillArray {
			if ss.Category == "Melee" {

				weapon := &runequest.Weapon{}

				for _, bw := range baseWeapons {
					if bw.MainSkill == ss.Name {
						weapon = bw
						break
					}
				}

				tempMelee[weapon.Name] = &runequest.Attack{
					Name:         weapon.Name,
					Skill:        ss,
					DamageString: weapon.Damage + dbString,
					StrikeRank:   c.Attributes["DEXSR"].Base + c.Attributes["SIZSR"].Base + weapon.SR,
					Weapon:       weapon,
				}
				break
			}
		}

		// Get best weapons skill
		for _, ss := range sortedSkillArray {
			if ss.Category == "Shield" {

				weapon := &runequest.Weapon{}

				for _, bw := range baseWeapons {
					if bw.MainSkill == ss.Name {
						weapon = bw
						break
					}
				}

				tempMelee[weapon.Name] = &runequest.Attack{
					Name:         weapon.Name,
					Skill:        ss,
					DamageString: weapon.Damage + dbString,
					StrikeRank:   c.Attributes["DEXSR"].Base + c.Attributes["SIZSR"].Base + weapon.SR,
					Weapon:       weapon,
				}
				break
			}
		}

		c.MeleeAttacks = tempMelee

		tempRanged := map[string]*runequest.Attack{}

		// Get best weapons skill
		for _, ss := range sortedSkillArray {
			if ss.Category == "Ranged" {

				weapon := &runequest.Weapon{}

				for _, bw := range baseWeapons {
					if bw.MainSkill == ss.Name {
						weapon = bw
						break
					}
				}

				damage := ""

				if weapon.Thrown {
					damage = weapon.Damage + throwDB
				} else {
					damage = weapon.Damage
				}

				tempRanged[weapon.Name] = &runequest.Attack{
					Name:         weapon.Name,
					Skill:        ss,
					DamageString: damage,
					StrikeRank:   c.Attributes["DEXSR"].Base + weapon.SR,
					Weapon:       weapon,
				}
				break
			}
		}

		c.RangedAttacks = tempRanged

		// Armor

		// Update CreationSteps
		c.CreationSteps["Finishing Touches"] = true
		c.CreationSteps["Complete"] = true

		// Save Character

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
