package runequest

// Skills is a map of basic common skills in Runequest
var Skills = map[string]*Skill{
	// Agility
	"Boat": &Skill{
		CoreString: "Boat",
		Base:       5,
		Category:   "Agility",
	},
	"Climb": &Skill{
		CoreString: "Climb",
		Base:       40,
		Category:   "Agility",
	},
	"Dodge": &Skill{
		CoreString: "Dodge",
		Base:       20,
		Category:   "Agility",
	},
	"Drive": &Skill{
		CoreString: "Drive",
		UserChoice: true,
		UserString: "Chariot",
		Base:       5,
		Category:   "Agility",
	},
	"Jump": &Skill{
		CoreString: "Jump",
		Base:       30,
		Category:   "Agility",
	},
	"Ride": &Skill{
		CoreString: "Ride",
		UserChoice: true,
		UserString: "Horse",
		Base:       5,
		Category:   "Agility",
	},
	"Swim": &Skill{
		CoreString: "Swim",
		Base:       15,
		Category:   "Agility",
	},

	// Communication
	"Act": &Skill{
		CoreString: "Act",
		Base:       5,
		Category:   "Communication",
	},
	"Art": &Skill{
		CoreString: "Art",
		Base:       5,
		Category:   "Communication",
	},
	"Bargain": &Skill{
		CoreString: "Bargain",
		Base:       5,
		Category:   "Communication",
	},
	"Charm": &Skill{
		CoreString: "Charm",
		Base:       15,
		Category:   "Communication",
	},
	"Dance": &Skill{
		CoreString: "Dance",
		Base:       10,
		Category:   "Communication",
	},
	"Disguise": &Skill{
		CoreString: "Disguise",
		Base:       5,
		Category:   "Communication",
	},
	"Fast Talk": &Skill{
		CoreString: "Fast Talk",
		Base:       5,
		Category:   "Communication",
	},
	"Intimidate": &Skill{
		CoreString: "Intimidate",
		Base:       15,
		Category:   "Communication",
	},
	"Intrigue": &Skill{
		CoreString: "Intrigue",
		Base:       5,
		Category:   "Communication",
	},
	"Orate": &Skill{
		CoreString: "Orate",
		Base:       10,
		Category:   "Communication",
	},
	"Sing": &Skill{
		CoreString: "Sing",
		Base:       10,
		Category:   "Communication",
	},
	"Speak": &Skill{
		CoreString: "Speak",
		UserChoice: true,
		UserString: "Heortling",
		Base:       50,
		Category:   "Communication",
	},
	// Knowledge
	"Alchemy": &Skill{
		CoreString: "Alchemy",
		Base:       0,
		Category:   "Knowledge",
	},
	"Animal Lore": &Skill{
		CoreString: "Animal Lore",
		Base:       5,
		Category:   "Knowledge",
	},
	"Battle": &Skill{
		CoreString: "Battle",
		Base:       10,
		Category:   "Knowledge",
	},
	"Bureacracy": &Skill{
		CoreString: "Bureacracy",
		Base:       0,
		Category:   "Knowledge",
	},
	"Celestial Lore": &Skill{
		CoreString: "Celestial Lore",
		Base:       5,
		Category:   "Knowledge",
	},
	"Cult Lore": &Skill{
		CoreString: "Cult Lore",
		UserChoice: true,
		UserString: "Orlanth",
		Base:       0,
		Category:   "Knowledge",
	},
	"Customs": &Skill{
		CoreString: "Customs",
		UserChoice: true,
		UserString: "Heortling",
		Base:       0,
		Category:   "Knowledge",
	},
	"Elder Race Lore": &Skill{
		CoreString: "Elder Race Lore",
		UserChoice: true,
		UserString: "Elves",
		Base:       5,
		Category:   "Knowledge",
	},
	"Evaluate": &Skill{
		CoreString: "Evaluate",
		Base:       10,
		Category:   "Knowledge",
	},
	"Farm": &Skill{
		CoreString: "Farm",
		Base:       10,
		Category:   "Knowledge",
	},
	"First Aid": &Skill{
		CoreString: "First Aid",
		Base:       10,
		Category:   "Knowledge",
	},
	"Game": &Skill{
		CoreString: "Game",
		Base:       15,
		Category:   "Knowledge",
	},
	"Herd": &Skill{
		CoreString: "Herd",
		Base:       5,
		Category:   "Knowledge",
	},
	"Homeland Lore": &Skill{
		CoreString: "Homeland Lore",
		UserChoice: true,
		UserString: "Local",
		Base:       30,
		Category:   "Knowledge",
	},
	"Library Use": &Skill{
		CoreString: "Library Use",
		Base:       0,
		Category:   "Knowledge",
	},
	"Lore": &Skill{
		CoreString: "Lore",
		UserChoice: true,
		UserString: "Local",
		Base:       0,
		Category:   "Knowledge",
	},
	"Manage Household": &Skill{
		CoreString: "Manage Household",
		Base:       10,
		Category:   "Knowledge",
	},
	"Mineral Lore": &Skill{
		CoreString: "Mineral Lore",
		Base:       5,
		Category:   "Knowledge",
	},
	"Peaceful Cut": &Skill{
		CoreString: "Peaceful Cut",
		Base:       10,
		Category:   "Knowledge",
	},
	"Plant Lore": &Skill{
		CoreString: "Plant Lore",
		Base:       5,
		Category:   "Knowledge",
	},
	"Read/Write": &Skill{
		CoreString: "Read/Write",
		UserChoice: true,
		UserString: "Old Tarsh",
		Base:       0,
		Category:   "Knowledge",
	},
	"Shiphandling": &Skill{
		CoreString: "Shiphandling",
		Base:       0,
		Category:   "Knowledge",
	},
	"Survival": &Skill{
		CoreString: "Survival",
		Base:       15,
		Category:   "Knowledge",
	},
	"Treat Disease": &Skill{
		CoreString: "Treat Disease",
		Base:       5,
		Category:   "Knowledge",
	},
	"Treat Poison": &Skill{
		CoreString: "Treat Poison",
		Base:       5,
		Category:   "Knowledge",
	},
	// Magic
	"Meditate": &Skill{
		CoreString: "Meditate",
		Base:       0,
		Category:   "Magic",
	},
	"Prepare Corpse": &Skill{
		CoreString: "Prepare Corpse",
		Base:       10,
		Category:   "Magic",
	},
	"Sense Assassin": &Skill{
		CoreString: "Sense Assassin",
		Base:       0,
		Category:   "Magic",
	},
	"Sense Chaos": &Skill{
		CoreString: "Sense Chaos",
		Base:       0,
		Category:   "Magic",
	},
	"Sorcery": &Skill{
		CoreString: "Sorcery",
		UserChoice: true,
		UserString: "Spell",
		Base:       0,
		Category:   "Magic",
	},
	"Spirit Combat": &Skill{
		CoreString: "Spirit Combat",
		Base:       20,
		Category:   "Magic",
	},
	"Spirit Dance": &Skill{
		CoreString: "Spirit Dance",
		Base:       0,
		Category:   "Magic",
	},
	"Spirit Lore": &Skill{
		CoreString: "Spirit Lore",
		Base:       0,
		Category:   "Magic",
	},
	"Spirit Travel": &Skill{
		CoreString: "Spirit Travel",
		Base:       0,
		Category:   "Magic",
	},
	"Understand Herd Beast": &Skill{
		CoreString: "Understand Herd Beast",
		Base:       0,
		Category:   "Magic",
	},
	"Worship": &Skill{
		CoreString: "Worship",
		UserChoice: true,
		UserString: "Orlanth",
		Base:       0,
		Category:   "Magic",
	},

	// Manipulation
	"Conceal": &Skill{
		CoreString: "Conceal",
		Base:       5,
		Category:   "Manipulation",
	},
	"Craft": &Skill{
		CoreString: "Craft",
		UserChoice: true,
		UserString: "Arms",
		Base:       10,
		Category:   "Manipulation",
	},
	"Devise": &Skill{
		CoreString: "Devise",
		Base:       5,
		Category:   "Manipulation",
	},
	"Play Instrument": &Skill{
		CoreString: "Play Instrument",
		Base:       5,
		Category:   "Manipulation",
	},
	"Sleight": &Skill{
		CoreString: "Sleight",
		Base:       10,
		Category:   "Manipulation",
	},

	// Melee Melees
	"1H Axe": &Skill{
		CoreString: "1H Axe",
		Base:       10,
		Category:   "Melee",
	},
	"2H Axe": &Skill{
		CoreString: "2H Axe",
		Base:       5,
		Category:   "Melee",
	},
	"Battle Axe": &Skill{
		CoreString: "Battle Axe",
		Base:       10,
		Category:   "Melee",
	},
	"Broadsword": &Skill{
		CoreString: "Broadsword",
		Base:       10,
		Category:   "Melee",
	},
	"Dagger": &Skill{
		CoreString: "Dagger",
		Base:       15,
		Category:   "Melee",
	},
	"Fist": &Skill{
		CoreString: "Fist",
		Base:       25,
		Category:   "Melee",
	},
	"Grapple": &Skill{
		CoreString: "Grapple",
		Base:       25,
		Category:   "Melee",
	},
	"1H Hammer": &Skill{
		CoreString: "1H Hammer",
		Base:       10,
		Category:   "Melee",
	},
	"2H Hammer": &Skill{
		CoreString: "2H Hammer",
		Base:       5,
		Category:   "Melee",
	},
	"Kick": &Skill{
		CoreString: "Kick",
		Base:       15,
		Category:   "Melee",
	},
	"Kopis": &Skill{
		CoreString: "Kopis",
		Base:       10,
		Category:   "Melee",
	},
	"1H Mace": &Skill{
		CoreString: "1H Mace",
		Base:       15,
		Category:   "Melee",
	},
	"2H Mace": &Skill{
		CoreString: "2H Mace",
		Base:       10,
		Category:   "Melee",
	},
	"Pike": &Skill{
		CoreString: "Pike",
		Base:       15,
		Category:   "Melee",
	},
	"Quarterstaff": &Skill{
		CoreString: "Quarterstaff",
		Base:       15,
		Category:   "Melee",
	},
	"Rapier": &Skill{
		CoreString: "Rapier",
		Base:       5,
		Category:   "Melee",
	},
	"Shortsword": &Skill{
		CoreString: "Shortsword",
		Base:       10,
		Category:   "Melee",
	},
	"1H Spear": &Skill{
		CoreString: "1H Spear",
		Base:       05,
		Category:   "Melee",
	},
	"2H Spear": &Skill{
		CoreString: "2H Spear",
		Base:       15,
		Category:   "Melee",
	},

	// Missile Weapons
	"Arbalest": &Skill{
		CoreString: "Arbalest",
		Base:       10,
		Category:   "Ranged",
	},
	"Axe, Throwing": &Skill{
		CoreString: "Axe, Throwing",
		Base:       10,
		Category:   "Ranged",
	},
	"Composite Bow": &Skill{
		CoreString: "Composite Bow",
		Base:       5,
		Category:   "Ranged",
	},
	"Crossbows": &Skill{
		CoreString: "Crossbows",
		Base:       25,
		Category:   "Ranged",
	},
	"Dagger, Throwing": &Skill{
		CoreString: "Dagger, Throwing",
		Base:       5,
		Category:   "Ranged",
	},
	"Elf Bow": &Skill{
		CoreString: "Elf Bow",
		Base:       5,
		Category:   "Ranged",
	},
	"Javelin": &Skill{
		CoreString: "Javelin",
		Base:       10,
		Category:   "Ranged",
	},
	"Pole Lasso": &Skill{
		CoreString: "Pole Lasso",
		Base:       5,
		Category:   "Ranged",
	},
	"Rock": &Skill{
		CoreString: "Rock",
		Base:       15,
		Category:   "Ranged",
	},
	"Self Bow": &Skill{
		CoreString: "Self Bow",
		Base:       5,
		Category:   "Ranged",
	},
	"Sling": &Skill{
		CoreString: "Sling",
		Base:       5,
		Category:   "Ranged",
	},
	"Staff Sling": &Skill{
		CoreString: "Staff Sling",
		Base:       10,
		Category:   "Ranged",
	},
	"Thrown Axe": &Skill{
		CoreString: "Thrown Axe",
		Base:       10,
		Category:   "Ranged",
	},
	"Throwing Dagger": &Skill{
		CoreString: "Throwing Dagger",
		Base:       10,
		Category:   "Ranged",
	},

	// Shields
	"Large Shield": &Skill{
		CoreString: "Large Shield",
		Base:       15,
		Category:   "Shield",
	},
	"Medium Shield": &Skill{
		CoreString: "Medium Shield",
		Base:       15,
		Category:   "Shield",
	},
	"Small Shield": &Skill{
		CoreString: "Small Shield",
		Base:       15,
		Category:   "Shield",
	},

	// Perception
	"Insight": &Skill{
		CoreString: "Insight",
		UserChoice: true,
		UserString: "Species",
		Base:       20,
		Category:   "Perception",
	},
	"Listen": &Skill{
		CoreString: "Listen",
		Base:       25,
		Category:   "Perception",
	},
	"Scan": &Skill{
		CoreString: "Scan",
		Base:       25,
		Category:   "Perception",
	},
	"Search": &Skill{
		CoreString: "Search",
		Base:       25,
		Category:   "Perception",
	},
	"Track": &Skill{
		CoreString: "Track",
		Base:       5,
		Category:   "Perception",
	},

	// Stealth
	"Hide": &Skill{
		CoreString: "Hide",
		Base:       10,
		Category:   "Stealth",
	},
	"Move Quietly": &Skill{
		CoreString: "Move Quietly",
		Base:       10,
		Category:   "Stealth",
	},
}
