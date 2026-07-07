// Package config lädt die Konfiguration aus einer .env-Datei und stellt
// typisierte Zugriffsfunktionen bereit. Es werden nur Pakete der
// Standardbibliothek verwendet.
package config

import (
	"bufio"
	"os"
	"strings"
)

// Load liest die angegebene .env-Datei und setzt die Werte als
// Umgebungsvariablen (bereits vorhandene Variablen werden nicht überschrieben).
// Eine fehlende Datei ist kein Fehler – dann gelten die OS-Variablen/Defaults.
func Load(path string) error {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.Trim(strings.TrimSpace(value), `"'`)
		if _, exists := os.LookupEnv(key); !exists {
			_ = os.Setenv(key, value)
		}
	}
	return sc.Err()
}

// Get liefert den Wert einer Variable oder den Fallback.
func Get(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
