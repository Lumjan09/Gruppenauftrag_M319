DROP TABLE IF EXISTS helden CASCADE;
DROP TABLE IF EXISTS skills CASCADE;
DROP TABLE IF EXISTS equipment CASCADE;

CREATE TABLE equipment (
    id            SERIAL PRIMARY KEY,
    name          VARCHAR(100) NOT NULL UNIQUE,
    typ           VARCHAR(20)  NOT NULL CHECK (typ IN ('weapon','armor','accessory')),
    attack_bonus  INTEGER NOT NULL DEFAULT 0,
    defense_bonus INTEGER NOT NULL DEFAULT 0,
    speed_bonus   INTEGER NOT NULL DEFAULT 0,
    hp_bonus      INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE skills (
    id           SERIAL PRIMARY KEY,
    name         VARCHAR(100) NOT NULL UNIQUE,
    beschreibung TEXT,
    rolle        VARCHAR(20) NOT NULL CHECK (rolle IN ('arkan','druide','kleriker','krieger','schmied','infiltrator')),
    damage_min   INTEGER NOT NULL DEFAULT 0 CHECK (damage_min >= 0),
    damage_max   INTEGER NOT NULL DEFAULT 0 CHECK (damage_max >= damage_min),
    heal         INTEGER NOT NULL DEFAULT 0 CHECK (heal >= 0),
    accuracy     NUMERIC(3,2) NOT NULL CHECK (accuracy >= 0.0 AND accuracy <= 1.0),
    target_type  VARCHAR(20) NOT NULL CHECK (target_type IN ('single_enemy','all_enemies','single_ally','all_allies','self'))
);

CREATE TABLE helden (
    id                    SERIAL PRIMARY KEY,
    name                  VARCHAR(100) NOT NULL,
    rolle                 VARCHAR(20) NOT NULL CHECK (rolle IN ('arkan','druide','kleriker','krieger','schmied','infiltrator')),
    max_hp                INTEGER NOT NULL CHECK (max_hp > 0),
    current_hp            INTEGER NOT NULL CHECK (current_hp >= 0),
    attack                INTEGER NOT NULL DEFAULT 0,
    defense               INTEGER NOT NULL DEFAULT 0,
    speed                 INTEGER NOT NULL DEFAULT 0,
    equipped_weapon_id    INTEGER REFERENCES equipment(id),
    equipped_armor_id     INTEGER REFERENCES equipment(id),
    equipped_accessory_id INTEGER REFERENCES equipment(id)
);