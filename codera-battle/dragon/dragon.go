// Package dragon enthält die vollständige Implementierung des Entropie-Drachen.
//
// Der Drache benötigt keine Datenbank – seine Werte und Fähigkeiten sind als
// Konstanten/Variablen definiert. Alle HP-Änderungen sind via sync.Mutex
// geschützt, damit parallele Goroutines (z. B. der Doppelangriff des Kriegers)
// keine Race Conditions verursachen.
package dragon

import (
	"math/rand"
	"sync"

	"codera-battle/internal"
)

// Basiswerte des Drachen.
const (
	dragonName       = "Entropie-Drache"
	dragonMaxHP      = 450
	dragonAttack     = 30
	dragonDefense    = 18
	dragonSpeed      = 14
	rageThreshold    = 0.30 // ab <= 30 % HP aktiviert sich Rage
	healThreshold    = 0.20 // unter 20 % HP versucht der Drache zu heilen
	healCooldown     = 4    // Corrupted Code nur alle 4 Runden
	rageDamageFactor = 1.5  // +50 % Schaden im Rage-Modus
)

// Skill-Indizes für die KI-Auswahl.
var (
	skillEntropyClaw = internal.Skill{
		Name: "Entropy Claw", MinDamage: 18, MaxDamage: 32, Accuracy: 0.90,
		Target: internal.TargetSingleEnemy, Description: "Krallenangriff",
	}
	skillNullPointer = internal.Skill{
		Name: "Null Pointer Breath", MinDamage: 24, MaxDamage: 42, Accuracy: 0.75,
		Target: internal.TargetSingleEnemy, Description: "Entropie-Atem",
	}
	skillStackOverflow = internal.Skill{
		Name: "Stack Overflow", MinDamage: 12, MaxDamage: 22, Accuracy: 0.60,
		Target: internal.TargetAllEnemies, Description: "Flächenangriff (alle Helden)",
	}
	skillCorruptedCode = internal.Skill{
		Name: "Corrupted Code", MinHeal: 20, MaxHeal: 20, Accuracy: 1.0,
		Target: internal.TargetSelf, Description: "Drache heilt sich",
	}
)

// Dragon ist der Endgegner. Er implementiert internal.Combatant.
type Dragon struct {
	currentHP int
	Enraged   bool

	defenseDebuff       int
	defenseDebuffRounds int
	lastHealRound       int

	mu sync.Mutex
}

// New erstellt den Drachen mit vollen HP.
func New() *Dragon {
	return &Dragon{currentHP: dragonMaxHP, lastHealRound: -healCooldown}
}

// GetName liefert den Namen des Drachen.
func (d *Dragon) GetName() string { return dragonName }

// GetMaxHP liefert die maximale HP des Drachen.
func (d *Dragon) GetMaxHP() int { return dragonMaxHP }

// GetStats liefert die effektiven Werte inkl. aktiver Defense-Debuffs.
func (d *Dragon) GetStats() internal.Stats {
	d.mu.Lock()
	defer d.mu.Unlock()
	def := dragonDefense - d.defenseDebuff
	if def < 0 {
		def = 0
	}
	return internal.Stats{MaxHP: dragonMaxHP, Attack: dragonAttack, Defense: def, Speed: dragonSpeed}
}

// GetCurrentHP liefert die aktuellen HP (thread-safe).
func (d *Dragon) GetCurrentHP() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.currentHP
}

// SetCurrentHP setzt die HP und begrenzt sie auf [0, MaxHP] (thread-safe).
func (d *Dragon) SetCurrentHP(hp int) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if hp < 0 {
		hp = 0
	}
	if hp > dragonMaxHP {
		hp = dragonMaxHP
	}
	d.currentHP = hp
}

// IsAlive meldet, ob der Drache noch lebt.
func (d *Dragon) IsAlive() bool { return d.GetCurrentHP() > 0 }

// HPPercent liefert den HP-Anteil zwischen 0.0 und 1.0.
func (d *Dragon) HPPercent() float64 {
	return float64(d.GetCurrentHP()) / float64(dragonMaxHP)
}

// RageActive meldet, ob der Rage-Modus aktiv ist (<= 30 % HP).
func (d *Dragon) RageActive() bool { return d.HPPercent() <= rageThreshold }

// RageFactor liefert den Schadensmultiplikator (1.0 oder 1.5 im Rage).
func (d *Dragon) RageFactor() float64 {
	if d.RageActive() {
		return rageDamageFactor
	}
	return 1.0
}

// ApplyDefenseDebuff senkt die Defense des Drachen für mehrere Runden.
func (d *Dragon) ApplyDefenseDebuff(amount, rounds int) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.defenseDebuff = amount
	d.defenseDebuffRounds = rounds
}

// ApplyDamage zieht Schaden thread-sicher von den HP ab (für parallele Goroutines).
func (d *Dragon) ApplyDamage(dmg int) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.currentHP -= dmg
	if d.currentHP < 0 {
		d.currentHP = 0
	}
}

// Heal heilt den Drachen thread-sicher.
func (d *Dragon) Heal(amount int) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.currentHP += amount
	if d.currentHP > dragonMaxHP {
		d.currentHP = dragonMaxHP
	}
}

// DefenseDebuffActive meldet, ob aktuell ein Defense-Debuff aktiv ist.
func (d *Dragon) DefenseDebuffActive() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.defenseDebuffRounds > 0
}

// TickEffects altert Debuffs am Rundenende.
func (d *Dragon) TickEffects() {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.defenseDebuffRounds > 0 {
		d.defenseDebuffRounds--
		if d.defenseDebuffRounds == 0 {
			d.defenseDebuff = 0
		}
	}
}

// ChooseAction wählt anhand der KI-Regeln die Aktion für die aktuelle Runde.
//
//   - Normalmodus (HP > 30 %): zufällige Skill-Auswahl, Corrupted Code nur
//     alle 4 Runden.
//   - Rage-Modus (HP <= 30 %): bevorzugt offensive Skills.
//   - Notheilung (HP < 20 %): 50 % Chance auf Corrupted Code statt Angriff,
//     sofern der Cooldown es zulässt.
func (d *Dragon) ChooseAction(round int) internal.Skill {
	d.Enraged = d.RageActive()
	canHeal := round-d.lastHealRound >= healCooldown

	// Notheilung
	if d.HPPercent() < healThreshold && canHeal && rand.Float64() < 0.50 {
		d.lastHealRound = round
		return skillCorruptedCode
	}

	offensive := []internal.Skill{skillEntropyClaw, skillNullPointer, skillStackOverflow}

	// Rage: nur offensive Skills.
	if d.Enraged {
		return offensive[rand.Intn(len(offensive))]
	}

	// Normalmodus: Corrupted Code in den Pool, wenn Cooldown bereit.
	pool := offensive
	if canHeal && d.HPPercent() < 1.0 {
		pool = append([]internal.Skill{skillCorruptedCode}, offensive...)
	}
	chosen := pool[rand.Intn(len(pool))]
	if chosen.Name == skillCorruptedCode.Name {
		d.lastHealRound = round
	}
	return chosen
}
