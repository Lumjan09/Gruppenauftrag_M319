// Package runenschmied implementiert den Charakter der Runenschmied*in (Klassen-Architekt).
package runenschmied

import (
	"codera-battle/hero"
	"codera-battle/internal"
)

const role = "Runenschmied*in"

func init() { hero.RegisterStrategy(role, strategy) }

// Definition liefert die Charakterdaten der Runenschmied*in.
func Definition() hero.Definition {
	return hero.Definition{
		Name: "Florentin",
		Role: role,
		Base: internal.Stats{MaxHP: 130, Attack: 16, Defense: 16, Speed: 10},
		Equipment: []internal.Equipment{
			{Name: "Architekten-Hammer", Type: internal.EquipWeapon, Bonus: internal.Stats{Attack: 7}},
			{Name: "Runen-Plattenpanzer", Type: internal.EquipArmor, Bonus: internal.Stats{Defense: 9}},
			{Name: "Siegelring der Stabilität", Type: internal.EquipAccessory, Bonus: internal.Stats{Speed: 1, MaxHP: 25}},
		},
		Skills: []internal.Skill{
			{Name: "Architekten-Schlag", MinDamage: 14, MaxDamage: 26, Accuracy: 0.85,
				Target: internal.TargetSingleEnemy, Description: "Solider physischer Angriff"},
			{Name: "Schutz-Rune", Accuracy: 1.0, Target: internal.TargetAllAllies,
				AllyDefenseBuff: 3, BuffRounds: 1, Description: "Erhöht Defense aller Helden um 3 für eine Runde"},
			{Name: "Konstrukt-Schild", Accuracy: 1.0, Target: internal.TargetSingleAlly,
				AllyDamageReduction: 0.50, BuffRounds: 1, Description: "Ein Verbündeter erhält -50 % Schaden für 1 Runde"},
		},
	}
}

// strategy schützt das Team automatisch: Konstrukt-Schild für den schwächsten
// Verbündeten unter 25 % HP, sonst Schutz-Rune bei durchschnittlichem Team-HP < 50 %.
func strategy(self *hero.Hero, allies []*hero.Hero, _ internal.Combatant, _ int) (*internal.Skill, internal.Combatant) {
	weakest := hero.LowestHPAlly(allies)
	if weakest != nil && weakest.HPPercent() < 0.25 {
		return self.FindSkill("Konstrukt-Schild"), weakest
	}
	if avgTeamHP(allies) < 0.50 {
		return self.FindSkill("Schutz-Rune"), nil
	}
	return nil, nil
}

func avgTeamHP(allies []*hero.Hero) float64 {
	var sum float64
	var n int
	for _, a := range allies {
		if a.IsAlive() {
			sum += a.HPPercent()
			n++
		}
	}
	if n == 0 {
		return 0
	}
	return sum / float64(n)
}
