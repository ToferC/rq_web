package runequest

// HomeLandStats provides basic stats for rolling
var HomeLandStats = map[string]*StatisticFrame{
	"STR": &StatisticFrame{
		Dice:     3,
		Modifier: 0,
	},
	"DEX": &StatisticFrame{
		Dice:     3,
		Modifier: 0,
	},
	"CON": &StatisticFrame{
		Dice:     3,
		Modifier: 0,
	},
	"POW": &StatisticFrame{
		Dice:     3,
		Modifier: 0,
	},
	"SIZ": &StatisticFrame{
		Dice:     2,
		Modifier: 6,
	},
	"INT": &StatisticFrame{
		Dice:     2,
		Modifier: 6,
	},
	"CHA": &StatisticFrame{
		Dice:     3,
		Modifier: 0,
	},
}
