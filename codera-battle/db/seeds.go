package db

import (
	"fmt"

	"gorm.io/gorm"

	"codera-battle/hero"
	"codera-battle/internal"

	"codera-battle/hero/arkan"
	"codera-battle/hero/druide"
	"codera-battle/hero/runenschmied"
)

// seedDefinitions liefert alle Charakterdefinitionen, die in die Datenbank
// geschrieben werden. Pro Gruppenmitglied genau eine Rolle:
//
//	Ron       -> Arkan-Dokumentar*in (Magier*in)
//	Lumjan    -> Daten-Druide (Formwandler*in)
//	Florentin -> Runenschmied*in (Architekt)
func seedDefinitions() []hero.Definition {
	return []hero.Definition{
		arkan.Definition(),
		druide.Definition(),
		runenschmied.Definition(),
	}
}

// Seed befüllt die Datenbank mit allen Helden, Ausrüstungsgegenständen und
// Skills. Der Vorgang ist idempotent: Sind bereits Helden vorhanden, passiert
// nichts.
func Seed(conn *gorm.DB) error {
	var count int64
	if err := conn.Model(&Hero{}).Count(&count).Error; err != nil {
		return fmt.Errorf("seed-prüfung fehlgeschlagen: %w", err)
	}
	if count > 0 {
		return nil
	}

	return conn.Transaction(func(tx *gorm.DB) error {
		for _, def := range seedDefinitions() {
			if err := seedDefinition(tx, def); err != nil {
				return err
			}
		}
		return nil
	})
}

// seedDefinition schreibt eine einzelne Charakterdefinition (Ausrüstung, Held,
// Skills) in die Datenbank.
func seedDefinition(tx *gorm.DB, def hero.Definition) error {
	hr := Hero{
		Name:      def.Name,
		Role:      def.Role,
		MaxHP:     def.Base.MaxHP,
		CurrentHP: def.Base.MaxHP,
		Attack:    def.Base.Attack,
		Defense:   def.Base.Defense,
		Speed:     def.Base.Speed,
	}

	for _, e := range def.Equipment {
		item := Equipment{
			Name:         e.Name,
			Type:         string(e.Type),
			BonusAttack:  e.Bonus.Attack,
			BonusDefense: e.Bonus.Defense,
			BonusSpeed:   e.Bonus.Speed,
			BonusHP:      e.Bonus.MaxHP,
			Special:      e.Special,
		}
		if err := tx.Create(&item).Error; err != nil {
			return fmt.Errorf("ausrüstung %q konnte nicht angelegt werden: %w", e.Name, err)
		}
		id := item.ID
		switch e.Type {
		case internal.EquipWeapon:
			hr.WeaponID = &id
		case internal.EquipArmor:
			hr.ArmorID = &id
		case internal.EquipAccessory:
			hr.AccessoryID = &id
		}
	}

	if err := tx.Create(&hr).Error; err != nil {
		return fmt.Errorf("held %q konnte nicht angelegt werden: %w", def.Name, err)
	}

	for _, s := range def.Skills {
		skill := Skill{
			Name:        s.Name,
			Role:        def.Role,
			MinDamage:   s.MinDamage,
			MaxDamage:   s.MaxDamage,
			MinHeal:     s.MinHeal,
			MaxHeal:     s.MaxHeal,
			Accuracy:    s.Accuracy,
			Target:      string(s.Target),
			Description: s.Description,
		}
		if err := tx.Create(&skill).Error; err != nil {
			return fmt.Errorf("skill %q konnte nicht angelegt werden: %w", s.Name, err)
		}
	}

	return nil
}
