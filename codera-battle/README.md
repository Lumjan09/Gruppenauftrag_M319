# Codera-Battle – Der finale Kampf gegen den Entropie-Drachen

Rundenbasiertes CLI-Kampfsystem in Go (Modul M319). Eine Heldengruppe tritt
gegen den vollständig vorgegebenen Entropie-Drachen an. Die Helden, ihre
Ausrüstung und Skills werden via GORM in einer PostgreSQL-Datenbank verwaltet;
der Drache ist in Konstanten definiert.

## Features

- Rundenbasierter Kampf-Loop mit Initiative nach Speed
- RNG-basierte Schadensberechnung (Streuung, Genauigkeit, kritische Treffer)
- Sechs Helden-Rollen mit eigenen Stats, Ausrüstung, Skills und Auto-KI
- Drachen-KI mit Normal-, Rage- und Notheilungs-Modus
- Goroutine-basierter Double Strike der Funktions-Krieger*in (Mutex-geschützt)
- Logging (Debug/Info/Warn/Error) mit täglicher Log-Rotation in `./logs/`
- Konfiguration über `.env`
- Unit-Tests für Schadensberechnung, Heilung und Kampf-Logik

## Projektstruktur

```
codera-battle/
├── main.go                 # Programmstart: Config, Logging, DB, Kampf
├── internal/               # Combatant-Interface, Stats, Skill, Equipment
├── dragon/                 # Vorgegebene Drachen-Implementierung (unverändert)
├── hero/                   # Basis-Heldenlogik + ein Paket pro Rolle
│   ├── arkan/              # Arkan-Dokumentar*in (Magier*in)
│   ├── druide/             # Daten-Druide (Formwandler*in)
│   ├── kleriker/           # Code-Kleriker*in (Heiler*in)
│   ├── krieger/            # Funktions-Krieger*in (Warrior)
│   ├── runenschmied/       # Runenschmied*in (Architekt)
│   └── rogue/              # System-Infiltrator*in (optional)
├── combat/                 # Kampf-Loop, Schadensberechnung, CLI-Anzeige
├── db/                     # GORM-Modelle, Migration, Seeds, Queries
├── config/                 # .env-Parser
├── logging/                # Logging-Framework mit Rotation
├── .env-example            # Konfigurationsvorlage
└── .gitignore
```

## Voraussetzungen

- Go 1.22 oder neuer
- Eine lokale PostgreSQL-Instanz (eigener Container pro Lernende*r, siehe M164)

## Setup

1. Repository klonen.
2. Abhängigkeiten laden (benötigt Internet):
   ```
   go mod tidy
   ```
   Dadurch werden `gorm.io/gorm` und `gorm.io/driver/postgres` geladen und
   `go.sum` erzeugt.
3. Konfiguration anlegen:
   ```
   cp .env-example .env
   ```
   Anschliessend die Zugangsdaten zur eigenen PostgreSQL-Instanz eintragen.
4. Programm starten:
   ```
   go run .
   ```
   Beim ersten Start werden Tabellen migriert und die Seed-Daten eingespielt.

## Tests

```
go test ./...
```

## Dokumentation (Godoc)

```
godoc -http :8080
```
Danach im Browser `http://localhost:8080/pkg/codera-battle/` öffnen.

## Hinweise

- **Drachen-HP:** Der Drache startet laut Spezifikation mit 450 HP. Der
  Beispiel-Screenshot im Auftrag zeigt `187/300` – das ist nur ein
  Anzeigebeispiel; massgeblich sind die Werte aus `dragon/dragon.go`.
- **Logfile:** Ist `./logs/battle.log` nicht beschreibbar, bricht das Programm
  beim Start kontrolliert mit einer aussagekräftigen Meldung ab.
- **Seed-Namen:** In `db/seeds.go` bzw. den jeweiligen Rollen-Paketen müssen die
  Platzhalternamen durch die echten Namen der Gruppenmitglieder ersetzt werden
  (siehe `// TODO`-Markierungen).
