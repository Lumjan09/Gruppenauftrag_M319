// Package internal stellt die geteilten Basistypen des Kampfsystems bereit.
//
// Hier liegen die von der Lehrperson vorgegebenen Typen (Stats, Combatant)
// sowie die von der Gruppe definierten Domänentypen (Skill, Equipment),
// damit sowohl Helden als auch der Drache dieselben Strukturen nutzen können.
package internal

// Stats enthält die Kampfwerte eines Kämpfers.
type Stats struct {
	MaxHP   int
	Attack  int
	Defense int
	Speed   int
}

// Combatant ist das gemeinsame Interface von Helden und Drache.
// Es wird sowohl von den Helden-Strukturen als auch vom Drachen implementiert.
type Combatant interface {
	GetName() string
	GetStats() Stats
	GetCurrentHP() int
	SetCurrentHP(hp int)
	GetMaxHP() int
	IsAlive() bool
}
