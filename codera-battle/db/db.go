package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"codera-battle/hero"
	"codera-battle/internal"
)

// Connect stellt die Verbindung zur PostgreSQL-Datenbank her.
func Connect(dsn string) (*gorm.DB, error) {
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("datenbankverbindung fehlgeschlagen: %w", err)
	}
	return conn, nil
}

// Migrate erzeugt bzw. aktualisiert die Tabellen für alle Modelle.
func Migrate(conn *gorm.DB) error {
	return conn.AutoMigrate(&Equipment{}, &Hero{}, &Skill{})
}

// LoadHeroes lädt alle Helden inkl. Ausrüstung und Skills und baut daraus die
// spielbaren Laufzeit-Helden. Die Strategie wird über die Rolle reaktiviert.
func LoadHeroes(conn *gorm.DB) ([]*hero.Hero, error) {
	var rows []Hero
	if err := conn.Preload("Weapon").Preload("Armor").Preload("Accessory").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("helden konnten nicht geladen werden: %w", err)
	}

	heroes := make([]*hero.Hero, 0, len(rows))
	for _, row := range rows {
		var skills []Skill
		if err := conn.Where("role = ?", row.Role).Find(&skills).Error; err != nil {
			return nil, fmt.Errorf("skills für %s konnten nicht geladen werden: %w", row.Role, err)
		}
		heroes = append(heroes, hero.New(toDefinition(row, skills)))
	}
	return heroes, nil
}

// toDefinition wandelt ein DB-Hero-Row in eine hero.Definition um.
func toDefinition(row Hero, skills []Skill) hero.Definition {
	def := hero.Definition{
		Name: row.Name,
		Role: row.Role,
		Base: internal.Stats{MaxHP: row.MaxHP, Attack: row.Attack, Defense: row.Defense, Speed: row.Speed},
	}
	for _, e := range []*Equipment{row.Weapon, row.Armor, row.Accessory} {
		if e != nil {
			def.Equipment = append(def.Equipment, toEquipment(*e))
		}
	}
	for _, s := range skills {
		def.Skills = append(def.Skills, internal.Skill{
			Name: s.Name, MinDamage: s.MinDamage, MaxDamage: s.MaxDamage,
			MinHeal: s.MinHeal, MaxHeal: s.MaxHeal, Accuracy: s.Accuracy,
			Target: internal.TargetType(s.Target), Description: s.Description,
		})
	}
	return def
}

func toEquipment(e Equipment) internal.Equipment {
	return internal.Equipment{
		Name: e.Name,
		Type: internal.EquipmentType(e.Type),
		Bonus: internal.Stats{
			MaxHP: e.BonusHP, Attack: e.BonusAttack, Defense: e.BonusDefense, Speed: e.BonusSpeed,
		},
		Special: e.Special,
	}
}
