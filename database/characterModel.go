package database

import (
	"fmt"
	"log"

	"github.com/go-pg/pg"
	"github.com/gosimple/slug"
	"github.com/toferc/rq_web/models"
	"github.com/toferc/runequest"
)

// Character represents a generic RPG character
type APICharacter struct {
	Name       string
	Role       string
	Homeland   string
	Occupation string
	Cult       string
	ExtraCults []string
	Age        int
	Clan       string
	Tribe      string
	Abilities  map[string]*runequest.Ability
	// Passions and Reputation
	ElementalRunes map[string]*runequest.Ability
	// Elemental Runes
	PowerRunes       map[string]*runequest.Ability
	ConditionRunes   map[string]*runequest.Ability
	CoreRunes        []*runequest.Ability
	Strength         int
	Dexterity        int
	Constitution     int
	Intelligence     int
	Power            int
	Charisma         int
	Attributes       map[string]*runequest.Attribute
	CurrentHP        int
	CurrentMP        int
	CurrentRP        int
	Movement         []*runequest.Movement
	Skills           map[string]*runequest.Skill
	SkillCategories  map[string]*runequest.SkillCategory
	Advantages       map[string]*runequest.Advantage
	RuneSpells       map[string]*runequest.Spell
	SpiritMagic      map[string]*runequest.Spell
	Powers           map[string]*runequest.Power
	LocationForm     string
	HitLocations     map[string]*runequest.HitLocation
	HitLocationMap   []string
	MeleeAttacks     map[string]*runequest.Attack
	RangedAttacks    map[string]*runequest.Attack
	Equipment        []string
	BoundSpirits     []*runequest.BoundSpirit
	Income           int
	Lunars           int
	Ransom           int
	StandardofLiving string
	InPlay           bool
	Updates          []*runequest.Update
	CreationSteps    map[string]bool

	Tags  []string
	Notes string
}

// SaveCharacterModel saves a Character to the DB
func SaveCharacterModel(db *pg.DB, cm *models.CharacterModel) error {

	cm.Character.UpdateCharacter()

	if cm.Slug == "" {
		cm.Slug = slug.Make(fmt.Sprintf("%s-%s", cm.Author.UserName, cm.Character.Name))
	}

	// Save character in Database
	_, err := db.Model(cm).
		OnConflict("(id) DO UPDATE").
		Set("character = ?character").
		Insert(cm)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(cm.ID, cm.Character.Name)

	updateTSVectors(db, cm.ID)

	return err
}

// UpdateCharacterModel updates a runequest character
func UpdateCharacterModel(db *pg.DB, cm *models.CharacterModel) error {

	cm.Character.UpdateCharacter()

	if cm.Slug == "" {
		cm.Slug = slug.Make(fmt.Sprintf("%s-%s", cm.Author.UserName, cm.Character.Name))
	}

	err := db.Update(cm)
	if err != nil {
		log.Println(err)
	}

	updateTSVectors(db, cm.ID)

	return err
}

// ListAllCharacterModels queries Character names and add to slice
func ListAllCharacterModels(db *pg.DB) ([]*models.CharacterModel, error) {
	var cms []*models.CharacterModel

	_, err := db.Query(&cms, `SELECT * FROM character_models ORDER BY created_at DESC;`)

	if err != nil {
		log.Println(err)
	}

	return cms, nil
}

// ListOpenCharacterModels queries open Character names and add to slice
func ListOpenCharacterModels(db *pg.DB) ([]*models.CharacterModel, error) {
	var cms []*models.CharacterModel

	_, err := db.Query(&cms, `SELECT * FROM character_models 
							WHERE open = true
							ORDER BY character ->> 'Name';`)

	if err != nil {
		log.Println(err)
	}

	return cms, nil
}

// PaginateCharacterModels queries open Character names and add to slice
func PaginateCharacterModels(db *pg.DB, limit, offset int) ([]*models.CharacterModel, error) {

	var cms []*models.CharacterModel

	err := db.Model(&cms).
		Limit(limit).
		Offset(offset * limit).
		Where("open = true").
		Order("created_at DESC").
		Select()

	if err != nil {
		log.Println(err)
	}

	return cms, nil
}

// ListCraftedCharacterModels queries open Character names and add to slice
func APICraftedCharacterModels(db *pg.DB) ([]*APICharacter, error) {
	var cms []*models.CharacterModel
	var characters []*APICharacter

	_, err := db.Query(&cms,
		`SELECT * FROM character_models WHERE open = true AND random = false
		UNION
		SELECT * FROM character_models WHERE open = true AND random IS NULL

		ORDER BY created_at DESC;`)

	if err != nil {
		log.Println(err)
	}

	println(len(cms))

	for _, cm := range cms {
		c := APICharacter{
			Name:             cm.Character.Name,
			Role:             cm.Character.Role,
			Homeland:         cm.Character.Homeland.Name,
			Occupation:       cm.Character.Occupation.Name,
			Cult:             cm.Character.Cult.Name,
			Age:              cm.Character.Age,
			Clan:             cm.Character.Clan,
			Tribe:            cm.Character.Tribe,
			Abilities:        cm.Character.Abilities,
			ElementalRunes:   cm.Character.ElementalRunes,
			PowerRunes:       cm.Character.PowerRunes,
			ConditionRunes:   cm.Character.ConditionRunes,
			CoreRunes:        cm.Character.CoreRunes,
			Strength:         cm.Character.Statistics["STR"].Total,
			Dexterity:        cm.Character.Statistics["DEX"].Total,
			Constitution:     cm.Character.Statistics["CON"].Total,
			Intelligence:     cm.Character.Statistics["INT"].Total,
			Power:            cm.Character.Statistics["POW"].Total,
			Charisma:         cm.Character.Statistics["CHA"].Total,
			Attributes:       cm.Character.Attributes,
			CurrentHP:        cm.Character.CurrentHP,
			CurrentMP:        cm.Character.CurrentMP,
			CurrentRP:        cm.Character.CurrentRP,
			Movement:         cm.Character.Movement,
			Skills:           cm.Character.Skills,
			SkillCategories:  cm.Character.SkillCategories,
			Advantages:       cm.Character.Advantages,
			RuneSpells:       cm.Character.RuneSpells,
			SpiritMagic:      cm.Character.SpiritMagic,
			Powers:           cm.Character.Powers,
			LocationForm:     cm.Character.LocationForm,
			HitLocations:     cm.Character.HitLocations,
			HitLocationMap:   cm.Character.HitLocationMap,
			MeleeAttacks:     cm.Character.MeleeAttacks,
			RangedAttacks:    cm.Character.RangedAttacks,
			Equipment:        cm.Character.Equipment,
			BoundSpirits:     cm.Character.BoundSpirits,
			Income:           cm.Character.Income,
			Lunars:           cm.Character.Lunars,
			Ransom:           cm.Character.Ransom,
			StandardofLiving: cm.Character.StandardofLiving,
			InPlay:           cm.Character.InPlay,
			Updates:          cm.Character.Updates,
			CreationSteps:    cm.Character.CreationSteps,
			Tags:             cm.Character.Tags,
			Notes:            cm.Character.Notes,
		}
		characters = append(characters, &c)
	}

	return characters, nil
}

// PaginateCraftedCharacterModels queries open Character names and add to slice
func PaginateCraftedCharacterModels(db *pg.DB, limit, offset int) ([]*models.CharacterModel, error) {
	var cms []*models.CharacterModel

	_, err := db.Query(&cms,
		`SELECT * FROM character_models WHERE open = true AND random = false
		UNION
		SELECT * FROM character_models WHERE open = true AND random IS NULL

		ORDER BY created_at DESC
		LIMIT ?
		OFFSET ?;`, limit, limit*offset)

	if err != nil {
		log.Println(err)
	}

	return cms, nil
}

// ListRandomCharacterModels queries open Character names and add to slice
func ListRandomCharacterModels(db *pg.DB, limit, offset int) ([]*models.CharacterModel, error) {
	var cms []*models.CharacterModel

	_, err := db.Query(&cms,
		`SELECT * FROM character_models WHERE open = true AND random = true
		ORDER BY created_at DESC
		LIMIT ?
		OFFSET ?;`, limit, limit*offset)

	if err != nil {
		log.Println(err)
	}

	return cms, nil
}

// ListUserCharacterModels queries Character names and add to slice
func ListUserCharacterModels(db *pg.DB, username string, limit, offset int) ([]*models.CharacterModel, error) {
	var cms []*models.CharacterModel

	err := db.Model(&cms).
		Limit(limit).
		Offset(offset*limit).
		Where("author ->> 'UserName' = ?", username).
		Order("created_at DESC").
		Select()

	if err != nil {
		log.Println(err)
	}

	return cms, nil
}

// ListOpenUserCharacterModels queries Character names and add to slice
func ListOpenUserCharacterModels(db *pg.DB, username string, limit, offset int) ([]*models.CharacterModel, error) {
	var cms []*models.CharacterModel

	err := db.Model(&cms).
		Limit(limit).
		Offset(offset*limit).
		Where("author ->> 'UserName' = ?", username).
		Where("open = true").
		Order("created_at DESC").
		Select()

	if err != nil {
		log.Println(err)
	}

	return cms, nil
}

func countUserCharacterModels(db *pg.DB, username string) int {

	var count int

	_, err := db.Query(&count,
		`SELECT COUNT(*) FROM character_models WHERE author ->> 'UserName' = ?;`, username)
	if err != nil {
		log.Println(err)
	}

	return count

}

// PKLoadCharacterModel loads a single character from the DB by pk
func PKLoadCharacterModel(db *pg.DB, pk int64) (*models.CharacterModel, error) {
	// Select user by Primary Key
	cm := &models.CharacterModel{ID: pk}
	err := db.Select(cm)

	if err != nil {
		fmt.Println(err)
		return &models.CharacterModel{}, err
	}

	if cm.Character.CreationSteps == nil {
		cm.Character.CreationSteps = map[string]bool{}
		cm.Character.CreationSteps["Complete"] = true
	}

	fmt.Println("Character loaded From DB")
	return cm, nil
}

// DeleteCharacterModel deletes a single character from DB by ID
func DeleteCharacterModel(db *pg.DB, pk int64) error {

	cm := models.CharacterModel{ID: pk}

	fmt.Println("Deleting character...")

	err := db.Delete(&cm)

	return err
}
