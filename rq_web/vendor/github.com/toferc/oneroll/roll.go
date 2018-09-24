package oneroll

import (
	"errors"
	"flag"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// Roll shows all results and variables from an ORE roll
type Roll struct {
	Actor      *Character
	Action     string // type of action act, oppose, maneuver
	NumActions int
	DiePool    *DiePool
	Results    []int
	Matches    []Match
	Loose      []int
	Wiggles    int
	Input      string
}

// DiePool represents a rollable dice set in ORE
type DiePool struct {
	Normal  int
	Hard    int
	Wiggle  int
	Expert  int
	Spray   int
	GoFirst int
}

func (d DiePool) String() string {
	var text string

	if d.Normal > 0 {
		text += fmt.Sprintf("%dd", d.Normal)
	}

	if d.Normal > 0 && d.Hard > 0 {
		text += "+"
	}

	if d.Hard > 0 {
		text += fmt.Sprintf("%dhd", d.Hard)
	}

	if d.Normal > 0 && d.Expert > 0 {
		text += "+"
	}

	if d.Expert > 0 {
		text += fmt.Sprintf("%ded", 1)
	}

	if (d.Hard > 0 && d.Wiggle > 0) || (d.Normal > 0 && d.Wiggle > 0 && d.Hard == 0) {
		text += "+"
	}

	if d.Wiggle > 0 {
		text += fmt.Sprintf("%dwd", d.Wiggle)
	}

	if d.GoFirst > 0 {
		text += fmt.Sprintf(" Go First %d", d.GoFirst)
	}

	if d.Spray > 0 {
		text += fmt.Sprintf(" Spray %d", d.Spray)
	}

	return text
}

// Match shows the height and width of a specific match
type Match struct {
	Actor      *Character
	Height     int
	Width      int
	Initiative int
}

// ByWidthHeight sorts matches in descending order of width then height
type ByWidthHeight []Match

func (a ByWidthHeight) Len() int      { return len(a) }
func (a ByWidthHeight) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByWidthHeight) Less(i, j int) bool {

	if a[i].Initiative > a[j].Initiative {
		return true
	}
	if a[i].Initiative < a[j].Initiative {
		return false
	}

	return a[i].Height > a[j].Height
}

// Resolve ORE dice roll and prints results
func (r *Roll) Resolve(input string) (*Roll, error) {

	r.Input = input

	nd, hd, wd, ed, gf, sp, ac, _, err := r.ParseString(input)

	r.NumActions = ac

	r.DiePool = &DiePool{
		Normal:  nd,
		Hard:    hd,
		Wiggle:  wd,
		Expert:  ed,
		GoFirst: gf,
		Spray:   sp,
	}

	if err != nil {
		return r, err
	}

	r.Wiggles = wd

	actionCount := ac // Disposable counter

	// Check for multiple actions and spray
	// reduce die pool
	if actionCount > 1 && sp == 0 {
		// Remove hd first
		for r.DiePool.Hard > 0 && actionCount > 1 {
			r.DiePool.Hard--
			actionCount--
		}
		for r.DiePool.Normal > 0 && actionCount > 1 {
			r.DiePool.Normal--
			actionCount--
		}
		for r.DiePool.Wiggle > 0 && actionCount > 1 {
			r.DiePool.Wiggle--
			actionCount--
		}
	}

	// Add spray dice to pool as normal dice
	if r.DiePool.Spray > 0 {
		r.DiePool.Normal += r.DiePool.Spray
	}

	// Ensure no more than 10d in pool
	r.verifyLessThan10() // Need to make sure 10d max after multiple actions

	for x := 0; x < r.DiePool.Normal; x++ {
		r.Results = append(r.Results, RollDie(10, 1, 1))
	}

	for x := 0; x < r.DiePool.Hard; x++ {
		r.Results = append(r.Results, 10)
	}

	if r.DiePool.Expert > 0 {
		r.Results = append(r.Results, r.DiePool.Expert)
	}

	r.parseDieRoll()

	// Sort roll by initiative (width+GoFirst) and then height
	sort.Sort(ByWidthHeight(r.Matches))

	return r, nil

}

// ParseString parses string like 5d+1hd+1wd or returns error
func (r *Roll) ParseString(input string) (int, int, int, int, int, int, int, int, error) {

	re := regexp.MustCompile("[0-9]+")

	var sElements []string

	errString := ""

	sElements = strings.SplitN(input, "+", 8)

	var nd, hd, wd, ed, gf, sp int

	ac, nr := 1, 1

	for _, s := range sElements {
		switch {
		case strings.Contains(s, "wd"):
			numString := re.FindString(s)
			wd, _ = strconv.Atoi(numString)

		case strings.Contains(s, "hd"):
			numString := re.FindString(s)
			hd, _ = strconv.Atoi(numString)

		case strings.Contains(s, "ed"):
			numString := re.FindString(s)
			ed, _ = strconv.Atoi(numString)

		case strings.Contains(s, "d"):
			numString := re.FindString(s)
			nd, _ = strconv.Atoi(numString)

		case strings.Contains(s, "gf"):
			numString := re.FindString(s)
			gf, _ = strconv.Atoi(numString)

		case strings.Contains(s, "sp"):
			numString := re.FindString(s)
			sp, _ = strconv.Atoi(numString)

		case strings.Contains(s, "ac"):
			numString := re.FindString(s)
			ac, _ = strconv.Atoi(numString)

		case strings.Contains(s, "nr"):
			numString := re.FindString(s)
			nr, _ = strconv.Atoi(numString)

		default:
			errString = "Error: Not a regular die notation"
		}
	}

	// Ensure at least one roll is made
	if nr < 1 {
		nr = 1
	}

	if errString != "" {
		return 0, 0, 0, 0, 0, 0, 0, 0, errors.New(errString)
	}

	return nd, hd, wd, ed, gf, sp, ac, nr, nil
}

// Determine matches including width, height and initiative for a roll
func (r *Roll) parseDieRoll() *Roll {

	matches := make(map[int]int)
	for _, d := range r.Results {
		matches[d]++
	}

	goFirst := 0
	if r.DiePool.GoFirst != 0 {
		goFirst = r.DiePool.GoFirst
	}

	for k, v := range matches {
		switch {
		case v == 1:
			r.Loose = append(r.Loose, k)
		case v > 1:
			r.Matches = append(r.Matches, Match{
				Actor:      r.Actor,
				Height:     k,
				Width:      v,
				Initiative: v + goFirst,
			})
		}
	}
	return r
}

// VerifyLessThan10 checks and reduces die pools to less than 10d
func (r *Roll) verifyLessThan10() {

	if SumDice(r.DiePool) > 10 {

		fmt.Println("Error: Can't roll more than 10 dice. Reducing to less than 10.")
		fmt.Printf(fmt.Sprintf("Current Dice: %dd+%dhd+%dwd+default%ded.\n",
			r.DiePool.Normal,
			r.DiePool.Hard,
			r.DiePool.Wiggle,
			r.DiePool.Expert,
		))

		// Remove normal dice first
		for r.DiePool.Normal > 0 && SumDice(r.DiePool) > 10 {
			fmt.Printf("reduced Normal dice from %d to %d. \n",
				r.DiePool.Normal,
				r.DiePool.Normal-1,
			)
			r.DiePool.Normal--
			fmt.Printf(fmt.Sprintf("Current Dice: %dd+%dhd+%dwd+default%ded.\n",
				r.DiePool.Normal,
				r.DiePool.Hard,
				r.DiePool.Wiggle,
				r.DiePool.Expert))
		}

		// Reduce hard dice next
		for r.DiePool.Hard > 0 && SumDice(r.DiePool) > 10 {
			fmt.Printf("reduced Hard dice from %d to %d. \n",
				r.DiePool.Hard,
				r.DiePool.Hard-1,
			)
			r.DiePool.Hard--
			fmt.Printf(fmt.Sprintf("Current Dice: %dd+%dhd+%dwd+default%ded.\n",
				r.DiePool.Normal,
				r.DiePool.Hard,
				r.DiePool.Wiggle,
				r.DiePool.Expert))
		}

		// Reduce expert dice next
		for r.DiePool.Expert > 0 && SumDice(r.DiePool) > 10 {
			fmt.Printf("reduced Expert dice from %d to %d. \n",
				1,
				0,
			)
			r.DiePool.Expert = 0
			fmt.Printf(fmt.Sprintf("Current Dice: %dd+%dhd+%dwd.\n",
				r.DiePool.Normal,
				r.DiePool.Hard,
				r.DiePool.Wiggle))
		}

		// Reduce wiggle dice last
		for r.DiePool.Wiggle > 0 && SumDice(r.DiePool) > 10 {
			fmt.Printf("reduced Wiggle dice from %d to %d. \n",
				r.DiePool.Wiggle,
				r.DiePool.Wiggle-1,
			)
			r.DiePool.Wiggle--

			fmt.Printf(fmt.Sprintf("Current Dice: %dd+%dhd+%dwd.\n",
				r.DiePool.Normal,
				r.DiePool.Hard,
				r.DiePool.Wiggle))
		}
	}
}

// Provides standard string formatting for roll
func (r Roll) String() string {

	text := ""
	var results []Match

	text += fmt.Sprintf("Actor: %s, Action: %s, Go First: %d, Spray: %d\n\n",
		r.Actor.Name,
		r.Action,
		r.DiePool.GoFirst,
		r.DiePool.Spray,
	)

	if len(r.Matches) > 0 {

		text += "Matches:\n"

		for _, m := range r.Matches {
			results = append(results, m)
		}
		sort.Sort(ByWidthHeight(results))
	}

	text += fmt.Sprintln("***Resolution***")
	text += fmt.Sprintf("%s Actions: %d\n", r.Actor.Name, r.NumActions)

	rs := TrimSliceBrackets(r.Results)

	text += fmt.Sprintf("Dice show: %s\n\n", rs)

	for _, m := range results {
		text += fmt.Sprintf("Match: %dx%d, Initiative: %dx%d\n",
			m.Width, m.Height,
			m.Initiative, m.Height,
		)
	}
	if r.Wiggles > 0 {
		text += fmt.Sprintf("+%d wiggle dice\n", r.Wiggles)
	}

	ls := TrimSliceBrackets(r.Loose)

	if len(r.Loose) > 0 {
		text += fmt.Sprintf("\nLoose dice %s\n", ls)
	}

	return text + "\n"
}

func main() {

	diePool := flag.String("d", "4d", "a die string separated by + like 4d+2hd+1wd")
	numRolls := flag.Int("n", 1, "an int that represents the number of rolls to make")

	flag.Parse()

	for x := 0; x < *numRolls; x++ {
		roll := Roll{}
		roll.Resolve(*diePool)
		fmt.Println(roll)
	}
}
