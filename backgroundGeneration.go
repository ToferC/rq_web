package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/mb-14/gomarkov"
	"github.com/toferc/runequest"
)

func buildModel(order int, file string) *gomarkov.Chain {
	chain := gomarkov.NewChain(order)
	for _, data := range getDataset(file) {
		chain.Add(split(data))
	}
	return chain
}

func split(str string) []string {
	return strings.Split(str, "")
}

func getDataset(fileName string) []string {
	file, _ := os.Open(fileName)
	scanner := bufio.NewScanner(file)
	var list []string
	for scanner.Scan() {
		list = append(list, scanner.Text())
	}
	return list
}

func loadModel(jsonFile string) (*gomarkov.Chain, error) {
	var chain gomarkov.Chain
	data, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		return &chain, err
	}
	err = json.Unmarshal(data, &chain)
	if err != nil {
		return &chain, err
	}
	return &chain, nil
}

func saveModel(chain *gomarkov.Chain, jsonFile string) {
	jsonObj, _ := json.Marshal(chain)
	err := ioutil.WriteFile(jsonFile, jsonObj, 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func generateName(chain *gomarkov.Chain) string {
	order := chain.Order
	tokens := make([]string, 0)
	for i := 0; i < order; i++ {
		tokens = append(tokens, gomarkov.StartToken)
	}
	for tokens[len(tokens)-1] != gomarkov.EndToken {
		next, _ := chain.Generate(tokens[(len(tokens) - order):])
		tokens = append(tokens, next)
	}
	return strings.Join(tokens[order:len(tokens)-1], "")
}

func generateBackground(homeland_name string, scale string) (string, string) {

	fmt.Println("Generating Background")
	
	// Select Gender
	var gender, description string
	rr := runequest.RollDice(100, 1, 0, 1)

	homeland := runequest.ToSnakeCase(homeland_name)

	switch homeland {
	case "sartar":
		homeland = "sartar"
	case "lunar_tarsh":
		homeland = "lunar"
	case "balazaring":
		homeland = "balazaring"
	case "grazelands":
		homeland = "grazelands"
	default:
		homeland = "sartar"
	}

	// WORKING BUT REMOVED FOR THE MOMENT
	switch {
	case rr < 45:
		gender = "Male"
		//chainModel = homeland + "MaleModel.json"
	case rr < 91:
		gender = "Female"
		//chainModel = homeland + "FemaleModel.json"
	default:
		gender = "Two Spirited"
		/*
		r2 := runequest.RollDice(10, 1, 0, 1)
		if r2 < 6 {
			chainModel = homeland + "FemaleModel.json"
		} else {
			chainModel = homeland + "MaleModel.json"
		}
		*/
	}

	// ERROR IN HERE
	/*
	// Load MarkovChains
	fmt.Println("Generating Chain")
	chain, err := loadModel(chainModel)
	if err != nil {
		log.Println(err)
	}

	name := generateName(chain)
	fmt.Println("Name: " + name)

	// Traits generation

	*/

	name := "NPC"

	traits := readCSV("traits.csv")

	t1 := traits[ChooseRandom(len(traits))]
	t2 := traits[ChooseRandom(len(traits))]

	description = fmt.Sprintf("%s is %s and %s.\n", name, t1, t2)

	pronoun := []string{}

	switch gender {
	case "Male":
		pronoun = []string{"He", "Him"}
	case "Female":
		pronoun = []string{"She", "Her"}
	case "Two Spirited":
		pronoun = []string{"They", "Them"}
	}

	description += fmt.Sprintf("%s is a %s adventurer. (%s/%s)", pronoun[0], scale, pronoun[0], pronoun[1])

	return name, description
}
