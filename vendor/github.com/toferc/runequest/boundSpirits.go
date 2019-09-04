package runequest

import "fmt"

// BoundSpirit represents spirits that characters can call upon
type BoundSpirit struct {
	Name              string
	Description       string
	Item              string
	Pow               int
	Cha               int
	Int               int
	CurrentMP         int
	SpiritMagicSpells map[string]Spell
}

func (b *BoundSpirit) String() string {
	text := fmt.Sprintf("%s - Pow: %d, Cha: %d, Int: %d, MP: %d", b.Name, b.Pow, b.Cha, b.Int, b.CurrentMP)

	if len(b.SpiritMagicSpells) > 0 {
		text += "\nSpells: " + formatSpellArray(b.SpiritMagicSpells)
	}

	return text
}
