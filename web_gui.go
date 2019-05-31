package main

import (
	"errors"
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gosimple/slug"
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
	Faction         *models.Faction
	Encounter       *models.Encounter
	Campaign        *models.Campaign
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
	Factions          []*models.Faction
	FactionCharacters []*models.CharacterModel
	Encounters        []*models.Encounter
	FactionMap        map[string][]*models.CharacterModel
	Campaigns         []*models.Campaign
	Passions          []string
	CategoryOrder     []string
	WeaponCategories  []string
	Roles             []string
	BaseWeapons       []*runequest.Weapon
	MeleeAttacks      map[string]*runequest.Attack
	RangedAttacks     map[string]*runequest.Attack
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

	Flashes []interface{}
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

func formatStringArray(a []string) string {
	text := ""
	end := len(a)

	for i, t := range a {
		if i+1 == end {
			text += t
		} else {
			text += t + ", "
		}
	}
	return text
}

func formatIntArray(a []int) string {
	text := strconv.Itoa(a[0])
	end := len(a)

	for i, t := range a {
		if i+1 == end {
			str := strconv.Itoa(t)
			text += "-" + str
		}
	}
	return text
}

func indexSpell(str string, spells []runequest.Spell) (int, error) {

	err := errors.New("Spell Not Found")

	for i, spell := range spells {
		if str == spell.CoreString {
			return i, nil
		}
	}

	return 0, err
}

func sortedSkills(skills map[string]*runequest.Skill) []*runequest.Skill {
	skillArray := []*runequest.Skill{}

	for _, v := range skills {
		skillArray = append(skillArray, v)
	}

	total := func(s1, s2 *runequest.Skill) bool {
		return s1.Total > s2.Total
	}

	By(total).Sort(skillArray)

	return skillArray[:9]
}

// By is the type of a "less" function that defines the ordering of its Planet arguments.
type By func(s1, s2 *runequest.Skill) bool

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by By) Sort(skills []*runequest.Skill) {
	ss := &skillSorter{
		skills: skills,
		by:     by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(ss)
}

// skillSorter joins a By function and a slice of Planets to be sorted.
type skillSorter struct {
	skills []*runequest.Skill
	by     func(p1, p2 *runequest.Skill) bool // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *skillSorter) Len() int {
	return len(s.skills)
}

// Swap is part of sort.Interface.
func (s *skillSorter) Swap(i, j int) {
	s.skills[i], s.skills[j] = s.skills[j], s.skills[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *skillSorter) Less(i, j int) bool {
	return s.by(s.skills[i], s.skills[j])
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
	case !cStep["Finishing Touches"]:
		url = "cc8_finishing_touches"
	}
	return url
}

func slugify(st string) string {
	return slug.Make(st)
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
		"formatStringArray":            formatStringArray,
		"formatIntArray":               formatIntArray,
		"sortedSkills":                 sortedSkills,
		"slugify":                      slugify,
	}

	baseTemplate := "templates/layout.html"

	tmpl[filename] = template.Must(template.New("").Funcs(funcMap).ParseFiles(filename, baseTemplate))

	if err := tmpl[filename].ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
