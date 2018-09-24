package oneroll

import "fmt"

// Capacity is Range, Mass, Touch or Speed
type Capacity struct {
	Type    string
	Level   int
	Value   string
	Booster *Booster
}

var powerCapacities = []byte(`{
  "mass": {
    "base": 25,
    "measure": "kg",
  },
  "range": {
    "base": 10,
    "measure": "m",
  },
  "speed": {
    "base": 2,
    "measure": "m",
  },
	"self": {
		"base": 0,
		"measure": nil,
	}
}`)

func (c Capacity) String() string {

	var text string

	if c.Type == "Self" {
		text = fmt.Sprintf("%s", c.Type)
	} else {
		text = fmt.Sprintf("%s (%s)",
			c.Type,
			c.Value)
	}

	return text
}

// Booster multiplies a Capacity or Statistic
type Booster struct {
	Level int
}
