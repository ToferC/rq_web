package runequest

// Race represents a genetically similar group of beings
type Race struct {
	Name       string
	StatBounds map[string]*StatisticBound
}

// StatisticBound are the limits for a stat based on Race
type StatisticBound struct {
	Stat  *Statistic
	Min   int
	Max   int
	Dice  int
	Bonus int
}
