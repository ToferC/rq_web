package runequest

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Max returns the larger of two ints
func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

// RollDie rolls and sum dice
func RollDie(max, min, numDice int) int {

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	result := 0
	for i := 1; i < numDice+1; i++ {
		roll := r1.Intn(max+1-min) + min
		result += roll
	}
	return result
}

// TrimSliceBrackets trims the brackets from a slice and return ints as a string
func TrimSliceBrackets(s []int) string {
	rs := fmt.Sprintf("%d", s)
	rs = strings.Trim(rs, "[]")
	return rs
}

// UserQuery creates and question and returns the User's input as a string
func UserQuery(q string) string {
	question := bufio.NewReader(os.Stdin)
	fmt.Print(q)
	r, _ := question.ReadString('\n')

	input := strings.Trim(r, " \n")

	return input
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

// ToSnakeCase transforms a string to snake case
func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// ChooseRandom chooses a random element in a slice
// when given l as len(slice)
func ChooseRandom(l int) int {
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	return r.Intn(l)
}

func readCSVFromURL(url string) ([][]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	reader := csv.NewReader(resp.Body)
	data, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func formatIntArray(a []int) string {
	text := strconv.Itoa(a[0])
	end := len(a)

	if len(a) > 1 {
		for i, t := range a {
			if i+1 == end {
				str := strconv.Itoa(t)
				text += "-" + str
			}
		}
	}
	return text
}

// IsInIntArray finds an int in an array of ints
func IsInIntArray(a []int, t int) bool {
	for _, i := range a {
		if t == i {
			return true
		}
	}
	return false
}
