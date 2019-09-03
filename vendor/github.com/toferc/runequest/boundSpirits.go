package runequest

import "fmt"

// BoundSpirit represents spirits that characters can call upon
type BoundSpirit struct {
	Name              string
	Description       string
	Item              string
	Pow               int
	Cha               int
	CurrentMP         int
	SpiritMagicSpells map[string]Spell
}

func (b *BoundSpirit) String() string {
	text := fmt.Sprintf("%s - Pow: %d, Cha: %d, MP: %d", b.Name, b.Pow, b.Cha, b.CurrentMP)

	if len(b.SpiritMagicSpells) > 0 {
		text += ", Spells: "
		for _, s := range b.SpiritMagicSpells {
			text += fmt.Sprintf("%v, ", s)
		}
	}

	return text
}
