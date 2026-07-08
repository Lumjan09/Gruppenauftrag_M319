// Package druide implementiert den Charakter des Daten-Druiden (Formwandler*in).
package druide

import (
	"codera-battle/hero"
	"codera-battle/internal"
)

const role = "Daten-Druide"

func init() { hero.RegisterStrategy(role, strategy) }

// Definition liefert die Charakterdaten des Daten-Druiden.
func Definition() hero.Definition {
	return hero.Definition{
		Name: "Lumjan",
		Role: role,
		Base: internal.Stats{MaxHP: 100, Attack: 14, Defense: 10, Speed: 16},
		Equipment: []internal.Equipment{
			{Name: "Transformations-Kristall", Type: internal.EquipWeapon, Bonus: internal.Stats{Attack: 6}},
			{Name: "Datenstrom-Mantel", Type: internal.EquipArmor, Bonus: internal.Stats{Defense: 4}},
			{Name: "Schema-Ring", Type: internal.EquipAccessory, Bonus: internal.Stats{Speed: 5, MaxHP: 10}},
		},
		Skills: []internal.Skill{
			{Name: "Datenklinge", MinDamage: 10, MaxDamage: 20, Accuracy: 0.85,
				Target: internal.TargetSingleEnemy, Description: "Transformierte Daten als Klinge"},
			{Name: "Strukturwandel", MinDamage: 14, MaxDamage: 28, Accuracy: 0.70,
				Target: internal.TargetSingleEnemy, Description: "Hoher Schaden, niedrige Genauigkeit"},
			{Name: "Transformative Regeneration", MinHeal: 12, MaxHeal: 20, Accuracy: 1.0,
				Target: internal.TargetSelf, Description: "Heilt sich selbst"},
		},
	}
}

// strategy heilt sich selbst automatisch, sobald die eigenen HP unter 40 % fallen.
func strategy(self *hero.Hero, _ []*hero.Hero, _ internal.Combatant, _ int) (*internal.Skill, internal.Combatant) {
	if self.HPPercent() < 0.40 {
		return self.FindSkill("Transformative Regeneration"), self
	}
	return nil, nil
}
