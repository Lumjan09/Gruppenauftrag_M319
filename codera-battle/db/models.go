// Package db enthält die GORM-Datenbankmodelle, die Migration, das Seeding
// und das Laden der Helden aus der PostgreSQL-Datenbank.
//
// Nur die Charakterdaten der Lernenden (Helden, Ausrüstung, Skills) werden
// persistiert. Der Drache ist in Konstanten definiert und benötigt keine DB.
package db

import "gorm.io/gorm"

// Equipment ist das GORM-Modell eines Ausrüstungsgegenstands.
type Equipment struct {
	gorm.Model
	Name         string
	Type         string // weapon | armor | accessory
	BonusAttack  int
	BonusDefense int
	BonusSpeed   int
	BonusHP      int
	Special      string // optionaler Spezialeffekt, z. B. life_steal
}

// Hero ist das GORM-Modell eines Helden inklusive der drei Ausrüstungs-Slots.
//
// Die Beziehung zu Equipment wird über drei einzelne Fremdschlüssel
// (Waffe, Rüstung, Accessoire) abgebildet.
type Hero struct {
	gorm.Model
	Name      string
	Role      string
	MaxHP     int
	CurrentHP int
	Attack    int
	Defense   int
	Speed     int

	WeaponID    *uint
	Weapon      *Equipment `gorm:"foreignKey:WeaponID"`
	ArmorID     *uint
	Armor       *Equipment `gorm:"foreignKey:ArmorID"`
	AccessoryID *uint
	Accessory   *Equipment `gorm:"foreignKey:AccessoryID"`
}

// Skill ist das GORM-Modell einer Fähigkeit, zugeordnet über die Rolle.
type Skill struct {
	gorm.Model
	Name        string
	Role        string
	MinDamage   int
	MaxDamage   int
	MinHeal     int
	MaxHeal     int
	Accuracy    float64
	Target      string
	Description string
}
