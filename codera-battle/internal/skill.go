package internal

// TargetType beschreibt das Ziel einer Fähigkeit.
type TargetType string

// Mögliche Zieltypen einer Fähigkeit.
const (
	TargetSingleEnemy TargetType = "single_enemy"
	TargetAllEnemies  TargetType = "all_enemies"
	TargetSingleAlly  TargetType = "single_ally"
	TargetAllAllies   TargetType = "all_allies"
	TargetSelf        TargetType = "self"
)

// EquipmentType unterscheidet die drei Ausrüstungs-Slots.
type EquipmentType string

// Mögliche Ausrüstungstypen.
const (
	EquipWeapon    EquipmentType = "weapon"
	EquipArmor     EquipmentType = "armor"
	EquipAccessory EquipmentType = "accessory"
)

// Equipment ist ein Ausrüstungsgegenstand mit Stat-Boni und optionalem Effekt.
//
// Die Boni werden über einen eingebetteten Stats-Wert abgebildet, sodass jede
// Ausrüstung MaxHP, Attack, Defense und Speed erhöhen kann.
type Equipment struct {
	Name    string
	Type    EquipmentType
	Bonus   Stats
	Special string // z. B. "life_steal"
}

// Skill ist eine Kampffähigkeit eines Helden oder des Drachen.
//
// Schaden und Heilung werden als Spanne (Min/Max) angegeben; konkrete Werte
// entstehen zur Laufzeit per RNG. Buff-/Debuff-Felder sind optional und nur
// gesetzt, wenn die Fähigkeit einen entsprechenden Effekt auslöst.
type Skill struct {
	Name        string
	MinDamage   int
	MaxDamage   int
	MinHeal     int
	MaxHeal     int
	Accuracy    float64
	Target      TargetType
	Description string

	// Optionale Effekte
	SelfDefenseBuff           int     // erhöht eigene Defense (z. B. Schutzschild)
	AllyDefenseBuff           int     // erhöht Defense aller/eines Verbündeten (Schutz-Rune)
	NextTurnAttackBuff        int     // erhöht Attack in der nächsten Runde (Kampfschrei)
	AllyDamageReduction       float64 // -% Schaden für einen Verbündeten (Konstrukt-Schild)
	BuffRounds                int     // Dauer der Buffs/Debuffs in Runden
	EnemyDefenseDebuff        int     // senkt Defense des Drachen (Schwachstelle)
	DoubleDamageEnemyBelowPct float64 // verdoppelt Schaden, wenn Drache unter X% HP
}

// IsHeal meldet, ob die Fähigkeit primär heilt.
func (s Skill) IsHeal() bool { return s.MaxHeal > 0 }

// IsDamage meldet, ob die Fähigkeit Schaden verursacht.
func (s Skill) IsDamage() bool { return s.MaxDamage > 0 }
