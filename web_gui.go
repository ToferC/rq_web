package main

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/toferc/rq_web/models"
	"github.com/toferc/runequest"
)

// WebView is a container for Web_gui data
type WebView struct {
	SessionUser string
	IsLoggedIn  string
	IsAdmin     string
	Actor       []*models.CharacterModel
	Actions     []int
	NumRolls    []int
	ErrorString []error
}

// WebUser represents a generic user struct
type WebUser struct {
	IsAuthor    bool
	SessionUser string
	IsLoggedIn  string
	IsAdmin     string
	Users       []*models.User
}

// WebChar is a framework to send objects & data to a Web view
type WebChar struct {
	User            *models.User
	CharacterModel  *models.CharacterModel
	HomelandModel   *models.HomelandModel
	OccupationModel *models.OccupationModel
	CultModel       *models.CultModel
	IsAuthor        bool
	SessionUser     string
	IsLoggedIn      string
	IsAdmin         string
	Wounds          map[string][]int
	// IndexModels
	CharacterModels   []*models.CharacterModel
	HomelandModels    map[string]*models.HomelandModel
	OccupationModels  map[string]*models.OccupationModel
	CultModels        map[string]*models.CultModel
	Passions          []string
	CategoryOrder     []string
	WeaponCategories  []string
	StandardsOfLiving []string
	PowerRunes        []string
	ElementalRunes    []string
	Skills            map[string]runequest.Skill
	SpiritMagic       []runequest.Spell
	RuneSpells        []runequest.Spell
	TotalSpiritMagic  []*runequest.Spell
	TotalRuneSpells   []*runequest.Spell
	NumRunePoints     []int
	NumSpiritMagic    []int
	Counter           []int
	MidCounter        []int
	BigCounter        []int
	RuneArray         []string
}

// SplitLines transfomrs results text string into slice
func SplitLines(s string) []string {
	sli := strings.Split(s, "/n")
	return sli
}

func sliceString(s string, i int) string {

	l := len(s)

	if l > i {
		return s[:i] + "..."
	}
	return s[:l]
}

func skillRoll(id int64, sk *runequest.Skill) string {

	text := ""

	return text
}

func statRoll(id int64, s *runequest.Statistic) string {

	text := ""

	return text
}

func subtract(a, b int) int {
	return a - b
}

func add(a, b int) int {
	return a + b
}

func multiply(a, b int) int {
	return a * b
}

func isIn(s []int, t int) bool {
	for _, n := range s {
		if n == t {
			return true
		}
	}
	return false
}

func isInString(s []string, t string) bool {
	for _, n := range s {
		if n == t {
			return true
		}
	}
	return false
}

// Generate URL for next step of Character creation
func generateCharacterCreationURL(cStep map[string]bool) string {

	url := ""

	switch {
	case !cStep["Personal History"]:
		url = "cc12_personal_history"
	case !cStep["Rune Affinities"]:
		url = "cc2_choose_runes"
	case !cStep["Roll Stats"]:
		url = "cc3_roll_stats"
	case !cStep["Apply Homeland"]:
		url = "cc4_apply_homeland"
	case !cStep["Apply Occupation"]:
		url = "cc5_apply_occupation"
	case !cStep["Apply Cult"]:
		url = "cc6_apply_cult"
	case !cStep["Personal Skills"]:
		url = "cc7_personal_skills"
	}
	return url
}

// Render combines templates, funcs and renders all Web pages in the app
func Render(w http.ResponseWriter, filename string, data interface{}) {

	tmpl := make(map[string]*template.Template)

	// Set up FuncMap
	funcMap := template.FuncMap{
		"skillRoll":                    skillRoll,
		"statRoll":                     statRoll,
		"subtract":                     subtract,
		"add":                          add,
		"multiply":                     multiply,
		"isIn":                         isIn,
		"sliceString":                  sliceString,
		"isInString":                   isInString,
		"generateCharacterCreationURL": generateCharacterCreationURL,
	}

	baseTemplate := "templates/layout.html"

	tmpl[filename] = template.Must(template.New("").Funcs(funcMap).ParseFiles(filename, baseTemplate))

	if err := tmpl[filename].ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
