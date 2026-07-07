// Command codera-battle startet das rundenbasierte Kampfsystem gegen den
// Entropie-Drachen.
//
// Ablauf:
//  1. Konfiguration aus .env laden
//  2. Logging initialisieren (panic, falls Logfile nicht beschreibbar)
//  3. Datenbankverbindung aufbauen (panic, falls nicht erreichbar)
//  4. Auto-Migration und Seeds ausführen
//  5. Helden aus der Datenbank laden
//  6. Drachen aus Konstanten erstellen
//  7. Kampf starten (CombatLoop)
package main

import (
	"fmt"

	"codera-battle/combat"
	"codera-battle/config"
	"codera-battle/db"
	"codera-battle/dragon"
	"codera-battle/logging"
)

func main() {
	// 1. Konfiguration laden. Ein fehlendes .env ist kein Fehler – die Werte
	// können auch direkt als Umgebungsvariablen gesetzt sein.
	if err := config.Load(".env"); err != nil {
		fmt.Println("Hinweis: .env konnte nicht geladen werden, nutze Umgebungsvariablen:", err)
	}

	// 2. Logging initialisieren. Ist das Logfile nicht beschreibbar, bricht das
	// Programm laut Vorgabe der Code-Kleriker*in kontrolliert ab.
	logLevel := config.Get("LOG_LEVEL", "info")
	logPath := config.Get("LOG_PATH", "./logs/battle.log")
	if err := logging.Init(logLevel, logPath); err != nil {
		panic(fmt.Sprintf("logging konnte nicht initialisiert werden: %v", err))
	}

	// 3. Datenbankverbindung herstellen. Schlägt sie fehl, ist ein Spielstart
	// nicht möglich – kontrollierter Abbruch.
	conn, err := db.Connect(buildDSN())
	if err != nil {
		logging.Error("datenbankverbindung fehlgeschlagen", "error", err)
		panic(fmt.Sprintf("datenbankverbindung fehlgeschlagen: %v", err))
	}
	logging.Info("datenbankverbindung hergestellt")

	// 4. Auto-Migration und Seeds.
	if err := db.Migrate(conn); err != nil {
		panic(fmt.Sprintf("migration fehlgeschlagen: %v", err))
	}
	if err := db.Seed(conn); err != nil {
		panic(fmt.Sprintf("seeding fehlgeschlagen: %v", err))
	}
	logging.Info("migration und seeds abgeschlossen")

	// 5. Helden laden.
	heroes, err := db.LoadHeroes(conn)
	if err != nil {
		panic(fmt.Sprintf("helden konnten nicht geladen werden: %v", err))
	}
	if len(heroes) == 0 {
		panic("keine helden in der datenbank gefunden")
	}
	logging.Info("helden geladen", "anzahl", len(heroes))

	// 6. Drachen erstellen (aus Konstanten, keine DB nötig).
	drg := dragon.New()

	// 7. Kampf starten.
	logging.Info("kampf beginnt")
	combat.CombatLoop(heroes, drg)
	logging.Info("kampf beendet")
}

// buildDSN baut den PostgreSQL-Verbindungsstring aus Umgebungsvariablen.
func buildDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Get("DB_HOST", "localhost"),
		config.Get("DB_PORT", "5432"),
		config.Get("DB_USER", "postgres"),
		config.Get("DB_PASSWORD", "postgres"),
		config.Get("DB_NAME", "codera"),
		config.Get("DB_SSLMODE", "disable"),
	)
}
