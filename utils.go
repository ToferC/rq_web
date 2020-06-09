package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/toferc/runequest"

	"github.com/gorilla/sessions"
)

func getUserSessionValues(s *sessions.Session) map[string]string {

	sessionMap := map[string]string{
		"username": "",
		"loggedin": "false",
		"isAdmin":  "false",
	}

	// Prep for user authentication

	u := s.Values["username"]
	l := s.Values["loggedin"]
	a := s.Values["isAdmin"]

	// Type assertation
	if user, ok := u.(string); !ok {
	} else {
		fmt.Println(user)
		sessionMap["username"] = user
	}

	// Type assertation
	if loggin, ok := l.(string); !ok {
	} else {
		fmt.Println(loggin)
		sessionMap["loggedin"] = loggin
	}

	// Type assertation
	if admin, ok := a.(string); !ok {
	} else {
		fmt.Println(admin)
		sessionMap["isAdmin"] = admin
	}
	return sessionMap
}

// numToArray takes and int and returns an array of [1:int]
func numToArray(m int) []int {

	a := []int{}

	for i := 1; i < m+1; i++ {
		a = append(a, i)
	}

	return a
}

// CreateUpdate creates an update for a runequest Object
func CreateUpdate(text string, value int) *runequest.Update {

	t := time.Now()
	tString := t.Format("2006-01-02")

	update := &runequest.Update{
		Date:  tString,
		Event: fmt.Sprintf("%s", text),
		Value: value,
	}
	fmt.Printf("Add Update: %s %d\n", text, value)
	return update
}

func createName(coreString string, userString string) string {

	targetString := ""

	if userString != "" {
		targetString = fmt.Sprintf("%s (%s)", coreString, userString)
	} else {
		targetString = fmt.Sprintf("%s", coreString)
	}

	return targetString
}

// ProcessUserString returns a simplified
func ProcessUserString(s string) string {
	trimmed := strings.TrimSpace(s)
	lower := strings.ToLower(trimmed)
	title := strings.Title(lower)

	return title
}

func readCSV(f string) []string {

	var a []string

	csvFile, err := os.Open(f)
	if err != nil {
		log.Println("Couldn't open CSV file", err)
	}
	r := csv.NewReader(csvFile)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Panic(err)
		}
		a = append(a, record[0])
	}
	return a
}

// ChooseRandom returns a random number betweeen 0 and l
func ChooseRandom(l int) int {
	return randomSeed.Intn(l)
}

// ChooseFromStringArray takes an array of strings to choose from and an array of strings already chosen
// it returns a slected string and an updated array of chosen stings
func ChooseFromStringArray(stringArray, chosenStrings []string) string {

	fmt.Println("Choose string")

	choice := ChooseRandom(len(stringArray))
	target := stringArray[choice]

	for isInString(chosenStrings, target) {
		fmt.Println("String already chosen")
		choice = ChooseRandom(len(stringArray))
		target = stringArray[choice]
	}

	chosenStrings = append(chosenStrings, target)

	return target
}

// ChooseFromSkillArray takes an array of Skills to choose from and an array of ints already chosen
// it returns a target index and an updated array of chosen indexes
func ChooseFromSkillArray(skillArray []*runequest.Skill, chosenInts []int) int {

	choice := ChooseRandom(len(skillArray))

	for isIn(chosenInts, choice) {
		fmt.Println("Int already chosen")
		choice = ChooseRandom(len(skillArray))
	}

	chosenInts = append(chosenInts, choice)

	return choice
}
