package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/mb-14/gomarkov"
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
