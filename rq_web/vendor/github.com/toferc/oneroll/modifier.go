package oneroll

import "fmt"

// Modifier models an Extra or Flaw for a Power Quality
type Modifier struct {
	Name          string
	Description   string
	RequiresLevel bool
	Level         int
	RequiresInfo  bool
	RequiresFocus bool
	Info          string
	CostPerLevel  int
	Cost          int
}

func (m Modifier) String() string {
	text := fmt.Sprintf("%s", m.Name)

	if m.RequiresLevel {
		text += fmt.Sprintf(" %d", m.Level)
	}

	if m.RequiresInfo {
		text += fmt.Sprintf(" - %s", m.Info)
	}

	if m.Cost > 0 {
		text += fmt.Sprintf(" (+%d/die)", m.Cost)
	} else {
		text += fmt.Sprintf(" (%d/die)", m.Cost)
	}
	return text
}

// NewModifier returns a new Modifier object
func NewModifier(s string) *Modifier {

	m := new(Modifier)

	m.Name = s
	m.Description = ""
	m.Level = 1
	m.Info = ""
	m.RequiresLevel = false
	m.RequiresInfo = false
	m.CostPerLevel = 1

	return m
}

// CalculateCost updates the cost per level for Modifiers w/ levels
// Called from Power.PowerCost()
func (m *Modifier) CalculateCost(b int) {

	if m.RequiresLevel {
		b = m.CostPerLevel * m.Level
	} else {
		b = m.CostPerLevel
	}
	m.Cost = b
}

// Modifiers creates map of standard WT extras & Flaws
var Modifiers = map[string]Modifier{

	// Extras
	"Custom Extra": Modifier{
		Name:          "Custom Extra",
		Description:   "",
		RequiresLevel: true,
		Level:         1,
		RequiresInfo:  true,
		Info:          "",
		CostPerLevel:  1,
	},
	"Custom Flaw": Modifier{
		Name:          "Custom Flaw",
		Description:   "",
		RequiresLevel: true,
		Level:         1,
		RequiresInfo:  true,
		Info:          "",
		CostPerLevel:  -1,
	},

	// Shadowrun modifiers
	"Drain": Modifier{
		Name:          "Drain",
		Description:   "One Shock damage once threshold passed.",
		RequiresLevel: false,
		Level:         1,
		RequiresInfo:  false,
		Info:          "",
		CostPerLevel:  -1,
	},
	"Essence Cost": Modifier{
		Name:          "Essence Cost",
		Description:   "",
		RequiresLevel: true,
		Level:         1,
		RequiresInfo:  false,
		Info:          "",
		CostPerLevel:  -1,
	},

	// Regular modifiers
	"Area": Modifier{
		Name:          "Area",
		Description:   "",
		RequiresLevel: true,
		Level:         1,
		RequiresInfo:  false,
		Info:          "",
		CostPerLevel:  1,
	},
	"Augment": Modifier{
		Name:          "Augment",
		Description:   "",
		RequiresLevel: false,
		RequiresInfo:  false,
		Info:          "",
		CostPerLevel:  4,
	},
	"Booster": Modifier{
		Name:          "Booster",
		Description:   "",
		RequiresLevel: true,
		Level:         1,
		CostPerLevel:  1,
	},
	"Burn": Modifier{
		Name:         "Burn",
		Description:  "",
		CostPerLevel: 2,
	},
	"Controlled Effect": Modifier{
		Name:         "Controlled Effect",
		Description:  "",
		CostPerLevel: 1,
	},
	"Daze": Modifier{
		Name:         "Daze",
		Description:  "",
		CostPerLevel: 1,
	},
	"Deadly": Modifier{
		Name:          "Deadly",
		Description:   "1: Killing, 2: Shock & Killing",
		RequiresLevel: true,
		Level:         1,
		CostPerLevel:  1,
	},
	"Disintigrate": Modifier{
		Name:         "Disintigrate",
		Description:  "",
		CostPerLevel: 2,
	},
	"Duration": Modifier{
		Name:         "Duration",
		Description:  "",
		CostPerLevel: 2,
	},
	"Electrocuting": Modifier{
		Name:         "Electrocuting",
		Description:  "",
		CostPerLevel: 1,
	},
	"Endless": Modifier{
		Name:         "Endless",
		Description:  "",
		CostPerLevel: 3,
	},
	"Engulf": Modifier{
		Name:         "Engulf",
		Description:  "",
		CostPerLevel: 2,
	},
	"Go First": Modifier{
		Name:          "Go First",
		Description:   "",
		RequiresLevel: true,
		Level:         1,
		CostPerLevel:  1,
	},
	"Hardened Defense": Modifier{
		Name:         "Hardened Defense",
		Description:  "",
		CostPerLevel: 2,
	},
	"High Capacity": Modifier{
		Name:         "High Capacity",
		Description:  "",
		RequiresInfo: true,
		Info:         "",
		CostPerLevel: 1,
	},
	"Interference": Modifier{
		Name:         "Interference",
		Description:  "",
		CostPerLevel: 3,
	},
	"Native Power": Modifier{
		Name:         "Native Power",
		Description:  "",
		CostPerLevel: 1,
	},
	"No Physics": Modifier{
		Name:         "No Physics",
		Description:  "",
		CostPerLevel: 1,
	},
	"No Upward Limit": Modifier{
		Name:         "No Upward Limit",
		Description:  "",
		CostPerLevel: 2,
	},
	"Non-Physical": Modifier{
		Name:         "Non-Physical",
		Description:  "",
		CostPerLevel: 2,
	},
	"On Sight": Modifier{
		Name:         "On Sight",
		Description:  "",
		CostPerLevel: 1,
	},
	"Penetration": Modifier{
		Name:          "Penetration",
		Description:   "",
		RequiresLevel: true,
		Level:         1,
		CostPerLevel:  1,
	},
	"Permanent": Modifier{
		Name:         "Permanent",
		Description:  "",
		CostPerLevel: 4,
	},
	"Radius": Modifier{
		Name:          "Radius",
		Description:   "10m x2/level",
		RequiresLevel: true,
		Level:         1,
		CostPerLevel:  2,
	},
	"Power Capacity": Modifier{
		Name:          "Power Capacity",
		Description:   "Power Capacity Type",
		RequiresLevel: true,
		Level:         2,
		RequiresInfo:  false,
		Info:          "Mass, Range, Speed or Touch",
		CostPerLevel:  1,
	},
	"Speeding Bullet": Modifier{
		Name:         "Speeding Bullet",
		Description:  "",
		CostPerLevel: 2,
	},
	"Spray": Modifier{
		Name:          "Spray",
		Description:   "",
		RequiresLevel: true,
		Level:         1,
		CostPerLevel:  1,
	},
	"Subtle": Modifier{
		Name:         "Subtle",
		Description:  "",
		CostPerLevel: 1,
	},
	"Traumatic": Modifier{
		Name:         "Traumatic",
		Description:  "",
		CostPerLevel: 1,
	},
	"Variable Effect": Modifier{
		Name:         "Variable Effect",
		Description:  "",
		CostPerLevel: 4,
	},

	// Flaws
	"Always On": Modifier{
		Name:          "Always On",
		Description:   "Combines with Permanent Extra",
		RequiresLevel: false,
		Level:         1,
		RequiresInfo:  false,
		Info:          "",
		CostPerLevel:  -1,
	},
	"Armored Defense": Modifier{
		Name:          "Armored Defense",
		Description:   "",
		RequiresLevel: false,
		Level:         1,
		RequiresInfo:  false,
		Info:          "",
		CostPerLevel:  -2,
	},
	"Attached": Modifier{
		Name: "Attached",
		Description: `Attached is worth –2 if it applies only when you use a specific Miracle or Skill. If Attached applies when you use a
particular Stat (which can be used with multiple Skills), it’s worth –1.`,
		RequiresLevel: true,
		Level:         1,
		RequiresInfo:  true,
		Info:          "",
		CostPerLevel:  -1,
	},
	"Automatic": Modifier{
		Name:          "Automatic",
		Description:   "",
		RequiresLevel: false,
		Level:         1,
		RequiresInfo:  false,
		Info:          "",
		CostPerLevel:  -1,
	},
	"Backfires": Modifier{
		Name:          "Backfires",
		Description:   "",
		RequiresLevel: false,
		Level:         1,
		RequiresInfo:  false,
		Info:          "",
		CostPerLevel:  -2,
	},
	"Base Will Cost": Modifier{
		Name:          "Base Will Cost",
		Description:   "",
		RequiresLevel: false,
		Level:         1,
		RequiresInfo:  false,
		Info:          "",
		CostPerLevel:  -4,
	},
	"Delayed Effect": Modifier{
		Name:          "Delayed Effect",
		Description:   "",
		RequiresLevel: false,
		Level:         1,
		RequiresInfo:  false,
		Info:          "",
		CostPerLevel:  -2,
	},
	"Depleted": Modifier{
		Name:          "Depleted",
		Description:   "",
		RequiresLevel: false,
		Level:         1,
		RequiresInfo:  true,
		Info:          "",
		CostPerLevel:  -1,
	},
	"Direct Feed": Modifier{
		Name:          "Direct Feed",
		Description:   "",
		RequiresLevel: false,
		Level:         1,
		RequiresInfo:  false,
		Info:          "",
		CostPerLevel:  -2,
	},
	"Exhausted": Modifier{
		Name:          "Exhausted",
		Description:   "",
		RequiresLevel: false,
		Level:         1,
		RequiresInfo:  false,
		Info:          "",
		CostPerLevel:  -3,
	},
	"Focus": Modifier{
		Name:          "Focus",
		Description:   "",
		RequiresLevel: false,
		Level:         1,
		RequiresInfo:  true,
		Info:          "",
		CostPerLevel:  -1,
	},
	"Full Power Only": Modifier{
		Name:          "Full Power Only",
		Description:   "",
		RequiresLevel: false,
		Level:         1,
		RequiresInfo:  false,
		Info:          "",
		CostPerLevel:  -1,
	},
	"Fragile": Modifier{
		Name:          "Fragile",
		Description:   "",
		RequiresLevel: false,
		Level:         1,
		RequiresInfo:  false,
		Info:          "",
		CostPerLevel:  -1,
	},
	"Go Last": Modifier{
		Name:          "Go Last",
		Description:   "",
		RequiresLevel: false,
		Level:         1,
		RequiresInfo:  false,
		Info:          "",
		CostPerLevel:  -1,
	},
	"Horrifying": Modifier{
		Name:          "Horrifying",
		Description:   "",
		RequiresLevel: false,
		Level:         1,
		RequiresInfo:  false,
		Info:          "",
		CostPerLevel:  -1,
	},
	"If/Then": Modifier{
		Name:          "If/Then",
		Description:   "",
		RequiresLevel: true,
		Level:         1,
		RequiresInfo:  true,
		Info:          "",
		CostPerLevel:  -1,
	},
	"Limited Damage": Modifier{
		Name:          "Limited Damage",
		Description:   "",
		RequiresLevel: false,
		Level:         1,
		RequiresInfo:  false,
		Info:          "",
		CostPerLevel:  -1,
	},
	"Limited Width": Modifier{
		Name:          "Limited Width",
		Description:   "",
		RequiresLevel: false,
		Level:         1,
		RequiresInfo:  false,
		Info:          "",
		CostPerLevel:  -1,
	},
	"Locational": Modifier{
		Name:          "Locational",
		Description:   "",
		RequiresLevel: false,
		Level:         1,
		RequiresInfo:  true,
		Info:          "",
		CostPerLevel:  -1,
	},
	"Loopy": Modifier{
		Name:          "Loopy",
		Description:   "",
		RequiresLevel: false,
		Level:         1,
		RequiresInfo:  false,
		Info:          "",
		CostPerLevel:  -1,
	},
	"Mental Strain": Modifier{
		Name:          "Mental Strain",
		Description:   "",
		RequiresLevel: false,
		Level:         1,
		RequiresInfo:  false,
		Info:          "",
		CostPerLevel:  -2,
	},
	"No Physical Change": Modifier{
		Name:          "No Physical Change",
		Description:   "",
		RequiresLevel: false,
		Level:         1,
		RequiresInfo:  false,
		Info:          "",
		CostPerLevel:  -1,
	},
	"Obvious": Modifier{
		Name:          "Obvious",
		Description:   "",
		RequiresLevel: false,
		Level:         1,
		RequiresInfo:  true,
		Info:          "",
		CostPerLevel:  -1,
	},
	"One Use": Modifier{
		Name:          "One Use",
		Description:   "",
		RequiresLevel: false,
		Level:         1,
		RequiresInfo:  false,
		Info:          "",
		CostPerLevel:  -4,
	},
	"Reduced Capacities": Modifier{
		Name:          "Reduced Capacities",
		Description:   "",
		RequiresLevel: false,
		Level:         1,
		RequiresInfo:  true,
		Info:          "",
		CostPerLevel:  -1,
	},
	"Scattered Damage": Modifier{
		Name:          "Scattered Damage",
		Description:   "",
		RequiresLevel: false,
		Level:         1,
		RequiresInfo:  false,
		Info:          "",
		CostPerLevel:  -1,
	},
	"Self Only": Modifier{
		Name:          "Self Only",
		Description:   "",
		RequiresLevel: false,
		Level:         1,
		RequiresInfo:  false,
		Info:          "",
		CostPerLevel:  -3,
	},
	"Slow": Modifier{
		Name:          "Slow",
		Description:   "",
		RequiresLevel: false,
		Level:         1,
		RequiresInfo:  false,
		Info:          "",
		CostPerLevel:  -2,
	},
	"Touch Only": Modifier{
		Name:          "Touch Only",
		Description:   "",
		RequiresLevel: false,
		Level:         1,
		RequiresInfo:  false,
		Info:          "",
		CostPerLevel:  -2,
	},
	"Uncontrollable": Modifier{
		Name:          "Uncontrollable",
		Description:   "",
		RequiresLevel: false,
		Level:         1,
		RequiresInfo:  false,
		Info:          "",
		CostPerLevel:  -2,
	},
	"Willpower Bid": Modifier{
		Name:          "Willpower Bid",
		Description:   "",
		RequiresLevel: false,
		Level:         1,
		RequiresInfo:  false,
		Info:          "",
		CostPerLevel:  -1,
	},
	"Willpower Cost": Modifier{
		Name:          "Willpower Cost",
		Description:   "",
		RequiresLevel: false,
		Level:         1,
		RequiresInfo:  false,
		Info:          "",
		CostPerLevel:  -2,
	},
	"Willpower Investment": Modifier{
		Name:          "Willpower Investment",
		Description:   "",
		RequiresLevel: false,
		Level:         1,
		RequiresInfo:  false,
		Info:          "",
		CostPerLevel:  -1,
	},

	// Focus Flaws
	"Accessible": Modifier{
		Name:          "Accessible",
		Description:   "",
		RequiresLevel: false,
		Level:         1,
		RequiresInfo:  false,
		RequiresFocus: true,
		Info:          "",
		CostPerLevel:  -1,
	},
	"Adaptation": Modifier{
		Name:          "Adaptation",
		Description:   "",
		RequiresLevel: false,
		Level:         1,
		RequiresInfo:  false,
		RequiresFocus: true,
		Info:          "",
		CostPerLevel:  -2,
	},
	"Booby-Trapped": Modifier{
		Name:          "Booby-Trapped",
		Description:   "",
		RequiresLevel: false,
		Level:         1,
		RequiresInfo:  false,
		RequiresFocus: true,
		Info:          "",
		CostPerLevel:  1,
	},
	"Bulky": Modifier{
		Name:          "Bulky",
		Description:   "",
		RequiresLevel: false,
		Level:         1,
		RequiresInfo:  false,
		RequiresFocus: true,
		Info:          "",
		CostPerLevel:  -1,
	},
	"Crew": Modifier{
		Name:          "Crew",
		Description:   "",
		RequiresLevel: false,
		Level:         1,
		RequiresInfo:  false,
		RequiresFocus: true,
		Info:          "",
		CostPerLevel:  -1,
	},
	"Delicate": Modifier{
		Name:          "Delicate",
		Description:   "",
		RequiresLevel: false,
		Level:         1,
		RequiresInfo:  false,
		RequiresFocus: true,
		Info:          "",
		CostPerLevel:  -1,
	},
	"Durable": Modifier{
		Name:          "Durable",
		Description:   "",
		RequiresLevel: false,
		Level:         1,
		RequiresInfo:  false,
		RequiresFocus: true,
		Info:          "",
		CostPerLevel:  1,
	},
	"Environment Bound": Modifier{
		Name:          "Environment Bound",
		Description:   "",
		RequiresLevel: false,
		Level:         1,
		RequiresInfo:  false,
		RequiresFocus: true,
		Info:          "",
		CostPerLevel:  -1,
	},
	"Friends Only": Modifier{
		Name:          "Friends Only",
		Description:   "",
		RequiresLevel: false,
		Level:         1,
		RequiresInfo:  false,
		RequiresFocus: true,
		Info:          "",
		CostPerLevel:  2,
	},
	"Immutable": Modifier{
		Name:          "Immutable",
		Description:   "",
		RequiresLevel: false,
		Level:         1,
		RequiresInfo:  false,
		RequiresFocus: true,
		Info:          "",
		CostPerLevel:  -1,
	},
	"Indestructible": Modifier{
		Name:          "Indestructible",
		Description:   "",
		RequiresLevel: false,
		Level:         1,
		RequiresInfo:  false,
		RequiresFocus: true,
		Info:          "",
		CostPerLevel:  2,
	},
	"Irreplaceable": Modifier{
		Name:          "Irreplaceable",
		Description:   "",
		RequiresLevel: false,
		Level:         1,
		RequiresInfo:  false,
		RequiresFocus: true,
		Info:          "",
		CostPerLevel:  -2,
	},
	"Manufacturable": Modifier{
		Name:          "Manufacturable",
		Description:   "",
		RequiresLevel: false,
		Level:         1,
		RequiresInfo:  false,
		RequiresFocus: true,
		Info:          "",
		CostPerLevel:  2,
	},
	"Operational Skill": Modifier{
		Name:          "Operational Skill",
		Description:   "",
		RequiresLevel: false,
		Level:         1,
		RequiresInfo:  true,
		RequiresFocus: true,
		Info:          "",
		CostPerLevel:  -1,
	},
	"Secret": Modifier{
		Name:          "Secret",
		Description:   "",
		RequiresLevel: false,
		Level:         1,
		RequiresInfo:  false,
		RequiresFocus: true,
		Info:          "",
		CostPerLevel:  1,
	},
	"Unwieldy": Modifier{
		Name:          "Unwieldy",
		Description:   "",
		RequiresLevel: true,
		Level:         1,
		RequiresInfo:  false,
		RequiresFocus: true,
		Info:          "",
		CostPerLevel:  -1,
	},
}
