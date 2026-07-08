// Package arkan implementiert den Charakter der Arkan-Dokumentar*in (Magier*in).
package arkan

import (
	"codera-battle/hero"
	"codera-battle/internal"
)

// role ist der eindeutige Rollenschlüssel dieses Charakters.
const role = "Arkan-Dokumentar*in"

func init() { hero.RegisterStrategy(role, strategy) }

// Definition liefert die Charakterdaten der Magier*in (Ron).
func Definition() hero.Definition {
	return hero.Definition{
		Name: "Ron",
		Role: role,
		Base: internal.Stats{MaxHP: 120, Attack: 18, Defense: 8, Speed: 14},
		Equipment: []internal.Equipment{
			{Name: "Pergament-Stab", Type: internal.EquipWeapon, Bonus: internal.Stats{Attack: 8}},
			{Name: "Runen-Gewand", Type: internal.EquipArmor, Bonus: internal.Stats{Defense: 5}},
			{Name: "Tintenfass-Amulett", Type: internal.EquipAccessory, Bonus: internal.Stats{Speed: 3, MaxHP: 20}},
		},
		Skills: []internal.Skill{
			{Name: "Runen-Geschoss", MinDamage: 12, MaxDamage: 24, Accuracy: 0.90,
				Target: internal.TargetSingleEnemy, Description: "Magischer Runenangriff"},
			{Name: "Arkaner Bann", MinDamage: 8, MaxDamage: 16, Accuracy: 0.85,
				Target: internal.TargetAllEnemies, Description: "Schwacher Flächenangriff"},
			{Name: "Klärende Annotation", MinHeal: 15, MaxHeal: 25, Accuracy: 1.0,
				Target: internal.TargetSingleAlly, Description: "Heilt einen Verbündeten"},
		},
	}
}

// strategy wählt automatisch den Heilzauber, wenn ein Verbündeter unter 40 % HP
// liegt; ansonsten überlässt sie die Wahl der Spielerin/dem Spieler.
func strategy(self *hero.Hero, allies []*hero.Hero, _ internal.Combatant, _ int) (*internal.Skill, internal.Combatant) {
	target := hero.LowestHPAlly(allies)
	if target != nil && target.HPPercent() < 0.40 {
		return self.FindSkill("Klärende Annotation"), target
	}
	return nil, nil
}
