# Gruppendokumentation – Codera-Battle (M319)

> Repository-Link: _hier den Link zum Git-Repository eintragen_

## 1. Überblick

Codera-Battle ist ein rundenbasiertes CLI-Kampfspiel in Go. Eine Heldengruppe
kämpft gegen den vorgegebenen Entropie-Drachen. Charakterdaten liegen in einer
PostgreSQL-Datenbank (GORM), die Kampflogik läuft in der Kommandozeile.

## 2. C4-Modell

### 2.1 Layer 1 – System-Kontext

```mermaid
flowchart TD
    Spieler["Spieler*in<br/>(Lernende)"]
    System["Codera-Battle<br/>(Go CLI-Anwendung)"]
    DB[("PostgreSQL<br/>Charakterdaten")]

    Spieler -->|Aktionsauswahl über CLI| System
    System -->|Helden, Ausrüstung, Skills laden/seeden| DB
    System -->|Kampfausgabe & Logs| Spieler
```

Der Spieler bedient das System über die CLI. Das System lädt die Charakterdaten
aus PostgreSQL und steuert den Kampf gegen den intern definierten Drachen.

### 2.2 Layer 2 – Container

```mermaid
flowchart TD
    CLI["CLI / main.go<br/>Programmstart & Steuerung"]
    Config["config<br/>.env-Parser"]
    Logging["logging<br/>slog + Rotation → ./logs/"]
    Combat["combat<br/>Kampf-Loop & Schadensberechnung"]
    Hero["hero/*<br/>Heldenrollen & KI"]
    Dragon["dragon<br/>Drachen-Logik (vorgegeben)"]
    DBPkg["db<br/>GORM-Modelle, Migration, Seeds"]
    PG[("PostgreSQL")]

    CLI --> Config
    CLI --> Logging
    CLI --> DBPkg
    CLI --> Combat
    DBPkg --> PG
    DBPkg --> Hero
    Combat --> Hero
    Combat --> Dragon
    Combat --> Logging
```

### 2.3 Layer 3 – Komponenten des `combat`-Pakets (Bonus)

```mermaid
flowchart TD
    Loop["CombatLoop"]
    Init["buildInitiativeOrder"]
    HeroTurn["processHeroTurn"]
    DragonTurn["processDragonTurn"]
    Calc["CalculateDamage"]
    Log["logAction"]
    HUD["printBattleHUD"]

    Loop --> Init
    Loop --> HeroTurn
    Loop --> DragonTurn
    HeroTurn --> HUD
    HeroTurn --> Calc
    HeroTurn --> Log
    DragonTurn --> Calc
    DragonTurn --> Log
```

## 3. Activity-Diagramme der Rollen

> Jedes Gruppenmitglied ergänzt hier das Activity-Diagramm seiner eigenen Rolle.
> Nachfolgend Beispiele für die Auto-KI zweier Rollen.

### 3.1 Arkan-Dokumentar*in (Heil-Strategie)

```mermaid
flowchart TD
    A([Zug beginnt]) --> B{Schwächster<br/>Verbündeter < 40% HP?}
    B -- ja --> C["Klärende Annotation<br/>auf schwächsten Verbündeten"]
    B -- nein --> D[Spieler*in wählt Aktion in CLI]
    C --> E([Zug Ende])
    D --> E
```

### 3.2 Funktions-Krieger*in (Double Strike)

```mermaid
flowchart TD
    A([Zug beginnt]) --> B{Eigene HP < 30%?}
    B -- ja --> C[Schutzschild]
    B -- nein --> D[Spieler*in wählt Aktion]
    D --> E{Präziser Hieb gewählt?}
    E -- ja --> F["2 Goroutines berechnen Schaden<br/>WaitGroup wartet"]
    F --> G["Mutex schützt Drachen-HP<br/>Schaden anwenden"]
    E -- nein --> H[Skill normal anwenden]
    C --> I([Zug Ende])
    G --> I
    H --> I
```

## 4. Datenmodell (GORM)

```mermaid
erDiagram
    HERO ||--o| EQUIPMENT : weapon
    HERO ||--o| EQUIPMENT : armor
    HERO ||--o| EQUIPMENT : accessory
    HERO {
        uint ID
        string Name
        string Role
        int MaxHP
        int Attack
        int Defense
        int Speed
    }
    EQUIPMENT {
        uint ID
        string Name
        string Type
        int BonusAttack
        int BonusDefense
    }
    SKILL {
        uint ID
        string Name
        string Role
        int MinDamage
        int MaxDamage
    }
```

Skills sind über das Feld `Role` einem Helden zugeordnet.

## 5. Aufgabenverteilung

| Person | Rolle | Schwerpunkt M319 |
|--------|-------|------------------|
| Ron | Arkan-Dokumentar*in (Magier*in) | C4-Diagramme, Clean Code, Linter |
| Lumjan | Daten-Druide (Formwandler*in) | GORM-Modelle, DB-Connection |
| Florentin | Runenschmied*in (Architekt) | DB-Migration & Seeds |

> Die Git-History weist pro Person die Commits der eigenen Rolle nach.
