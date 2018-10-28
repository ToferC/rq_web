package main

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/toferc/onegc/model"
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
	User            model.User
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
	Skills            map[string]*runequest.Skill
	SpiritMagic       []runequest.Spell
	RuneSpells        []runequest.Spell
	SubCults          []runequest.Cult
	NumRunePoints     []int
	NumSpiritMagic    []int
	Counter           []int
	MidCounter        []int
	BigCounter        []int
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

// Render combines templates, funcs and renders all Web pages in the app
func Render(w http.ResponseWriter, filename string, data interface{}) {

	tmpl := make(map[string]*template.Template)

	// Set up FuncMap
	funcMap := template.FuncMap{
		"skillRoll":   skillRoll,
		"statRoll":    statRoll,
		"subtract":    subtract,
		"add":         add,
		"multiply":    multiply,
		"isIn":        isIn,
		"sliceString": sliceString,
		"isInString":  isInString,
	}

	baseTemplate := "templates/layout.html"

	tmpl[filename] = template.Must(template.New("").Funcs(funcMap).ParseFiles(filename, baseTemplate))

	if err := tmpl[filename].ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
