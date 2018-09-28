package oneroll

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Dice is an interface for Statistic & Skills to combine DiePool
type Dice interface {
	getDiePool() *DiePool
}

// Ability is an interface for general ORE object operations
type Ability interface {
	CalculateCost()
}

// ReturnDice implements the Ability to combine DiePool
func ReturnDice(d Dice) *DiePool {
	return d.getDiePool()
}

// UpdateCost implements Ability interface to generate costs
func UpdateCost(a Ability) {
	a.CalculateCost()
}

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

// ParseNumRolls checks how many die rolls are required
func ParseNumRolls(s string) (int, error) {

	re := regexp.MustCompile("[0-9]+")

	var num int
	var numString string

	numString = re.FindString(s)
	num, err := strconv.Atoi(numString)
	if err != nil {
		num = 1
	}
	return num, err
}

// SkillRated returns true if a skill has any points in it
func SkillRated(s *Skill) bool {
	if s.Dice.Normal+s.Dice.Hard+s.Dice.Wiggle > 0 {
		return true
	}
	return false
}

// SumDice sums a DiePool - used in determining BaseWill
func SumDice(d *DiePool) int {

	var r int

	if d.Expert > 0 {
		r = d.Normal + d.Hard + d.Wiggle + 1
	} else {
		r = d.Normal + d.Hard + d.Wiggle
	}
	return r
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
