package oneroll

type CostFramework struct {
	Setting    string
	Stat       int
	Skill      int
	HyperSkill int
	HyperStat  int
	Quality    int
	HardMult   int
	WiggleMult int
	ExpertMult int
}

var Settings = map[string]CostFramework{
	"WT": CostFramework{
		Setting:    "Wild Talents",
		Stat:       5,
		Skill:      2,
		HyperSkill: 1,
		HyperStat:  4,
		Quality:    2,
		HardMult:   2,
		WiggleMult: 4,
		ExpertMult: 2,
	},
	"SR": CostFramework{
		Setting:    "Shadowrun",
		Stat:       5,
		Skill:      2,
		HyperSkill: 1,
		HyperStat:  4,
		Quality:    2,
		HardMult:   2,
		WiggleMult: 4,
		ExpertMult: 2,
	},
	"RE": CostFramework{
		Setting:    "Reign",
		Stat:       5,
		Skill:      1,
		HardMult:   2,
		WiggleMult: 6,
		ExpertMult: 2,
	},
}
