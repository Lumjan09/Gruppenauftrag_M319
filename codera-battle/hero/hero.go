// Package hero stellt die Laufzeit-Repräsentation eines Helden bereit.
//
// Jede Rolle (Unterpaket hero/<rolle>) liefert über Definition() ihre
// Charakterdaten. Aus diesen Definitionen werden die Seed-Daten erzeugt und –
// nach dem Laden aus der Datenbank – die spielbaren Hero-Objekte gebaut.
package hero

import (
	"sync"

	"codera-battle/internal"
)

// Strategy ist eine optionale KI-Funktion eines Helden.
//
// Liefert sie einen Skill (und ggf. ein Ziel) zurück, wird die Aktion
// automatisch ausgeführt; bei (nil, nil) wählt die Spielerin/der Spieler
// die Aktion in der CLI selbst.
type Strategy func(self *Hero, allies []*Hero, dragon internal.Combatant, round int) (*internal.Skill, internal.Combatant)

// effect ist ein temporärer, rundenbasierter Buff/Debuff auf einem Helden.
type effect struct {
	attack          int
	defense         int
	damageReduction float64
	rounds          int
}

// Hero ist ein spielbarer Held. Er implementiert internal.Combatant.
type Hero struct {
	Name      string
	Role      string
	BaseStats internal.Stats
	Equipment []internal.Equipment
	Skills    []internal.Skill
	Strategy  Strategy

	currentHP int
	effects   []effect
	mu        sync.Mutex
}

// Definition beschreibt einen Charakter unabhängig von der Laufzeit.
// Sie wird zum Seeden und zum Wiederaufbau aus der DB verwendet.
type Definition struct {
	Name      string
	Role      string
	Base      internal.Stats
	Equipment []internal.Equipment
	Skills    []internal.Skill
}

// New erzeugt einen Helden aus einer Definition. Die Strategie wird – sofern
// für die Rolle registriert – automatisch zugewiesen.
func New(def Definition) *Hero {
	h := &Hero{
		Name:      def.Name,
		Role:      def.Role,
		BaseStats: def.Base,
		Equipment: def.Equipment,
		Skills:    def.Skills,
		Strategy:  StrategyFor(def.Role),
	}
	h.currentHP = h.GetMaxHP()
	return h
}

// GetName liefert den Anzeigenamen.
func (h *Hero) GetName() string { return h.Name }

// GetMaxHP liefert die maximale HP inklusive Ausrüstungs-Boni.
func (h *Hero) GetMaxHP() int {
	hp := h.BaseStats.MaxHP
	for _, e := range h.Equipment {
		hp += e.Bonus.MaxHP
	}
	return hp
}

// GetStats liefert die effektiven Kampfwerte inkl. Ausrüstung und aktiver Effekte.
func (h *Hero) GetStats() internal.Stats {
	s := h.BaseStats
	s.MaxHP = h.GetMaxHP()
	for _, e := range h.Equipment {
		s.Attack += e.Bonus.Attack
		s.Defense += e.Bonus.Defense
		s.Speed += e.Bonus.Speed
	}
	for _, ef := range h.effects {
		s.Attack += ef.attack
		s.Defense += ef.defense
	}
	return s
}

// GetCurrentHP liefert die aktuellen HP (thread-safe).
func (h *Hero) GetCurrentHP() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.currentHP
}

// SetCurrentHP setzt die aktuellen HP und begrenzt sie auf [0, MaxHP].
func (h *Hero) SetCurrentHP(hp int) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if hp < 0 {
		hp = 0
	}
	if max := h.GetMaxHP(); hp > max {
		hp = max
	}
	h.currentHP = hp
}

// IsAlive meldet, ob der Held noch lebt (HP > 0).
func (h *Hero) IsAlive() bool { return h.GetCurrentHP() > 0 }

// HPPercent liefert den HP-Anteil zwischen 0.0 und 1.0.
func (h *Hero) HPPercent() float64 {
	max := h.GetMaxHP()
	if max == 0 {
		return 0
	}
	return float64(h.GetCurrentHP()) / float64(max)
}

// DamageReduction liefert die aktuell aktive Schadensreduktion (0.0–0.9).
func (h *Hero) DamageReduction() float64 {
	r := 0.0
	for _, ef := range h.effects {
		r += ef.damageReduction
	}
	if r > 0.9 {
		r = 0.9
	}
	return r
}

// HasLifeSteal meldet, ob der Held einen Lebensraub-Effekt trägt.
func (h *Hero) HasLifeSteal() bool {
	for _, e := range h.Equipment {
		if e.Special == "life_steal" {
			return true
		}
	}
	return false
}

// AddEffect fügt einen rundenbasierten Buff/Debuff hinzu.
func (h *Hero) AddEffect(attack, defense int, damageReduction float64, rounds int) {
	h.effects = append(h.effects, effect{attack, defense, damageReduction, rounds})
}

// TickEffects altert alle Effekte um eine Runde und entfernt abgelaufene.
func (h *Hero) TickEffects() {
	active := h.effects[:0]
	for _, ef := range h.effects {
		ef.rounds--
		if ef.rounds > 0 {
			active = append(active, ef)
		}
	}
	h.effects = active
}

// --- Strategie-Registry -------------------------------------------------------

var strategies = map[string]Strategy{}

// RegisterStrategy hinterlegt die KI-Funktion einer Rolle.
// Wird typischerweise im init() des jeweiligen Rollenpakets aufgerufen.
func RegisterStrategy(role string, s Strategy) { strategies[role] = s }

// StrategyFor liefert die registrierte Strategie einer Rolle (oder nil).
func StrategyFor(role string) Strategy { return strategies[role] }

// FindSkill liefert einen Pointer auf den Skill mit dem gegebenen Namen (oder nil).
func (h *Hero) FindSkill(name string) *internal.Skill {
	for i := range h.Skills {
		if h.Skills[i].Name == name {
			return &h.Skills[i]
		}
	}
	return nil
}

// LowestHPAlly liefert den lebenden Verbündeten mit den wenigsten HP.
func LowestHPAlly(allies []*Hero) *Hero {
	var lowest *Hero
	for _, a := range allies {
		if !a.IsAlive() {
			continue
		}
		if lowest == nil || a.GetCurrentHP() < lowest.GetCurrentHP() {
			lowest = a
		}
	}
	return lowest
}
