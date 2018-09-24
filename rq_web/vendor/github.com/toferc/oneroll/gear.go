package oneroll

// Weapon models a weapon in ORE
type Weapon struct {
	ID          int64
	Name        string
	Material    Material
	Volume      int // In litres
	Wounds      int
	HAR         int
	LAR         int
	Shock       string
	Kill        string
	Area        int
	Penetration int
}

// Armor models physical defenses in ORE
type Armor struct {
	ID        int64
	Name      string
	Material  Material
	HAR       int
	LAR       int
	Locations []int
	Volume    int // In litres
}

// Material models a material type in ORE
type Material struct {
	ID       int64
	Name     string
	HAR      int
	LAR      int
	Mass     int // Per litre
	Melt     int // in C
	Freeze   int // in C
	Vaporize int // in C
	Burn     int // in C
}
