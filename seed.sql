INSERT INTO equipment (name, typ, attack_bonus, defense_bonus, speed_bonus, hp_bonus) VALUES
('Pergament-Stab',             'weapon',    8, 0, 0,  0),
('Runen-Gewand',               'armor',     0, 5, 0,  0),
('Tintenfass-Amulett',         'accessory', 0, 0, 3, 20),
('Datenstrom-Mantel',          'armor',     0, 4, 0,  0),
('Transformations-Kristall',   'weapon',    6, 0, 0,  0),
('Schema-Ring',                'accessory', 0, 0, 5, 10),
('Architekten-Hammer',         'weapon',    7, 0, 0,  0),
('Runen-Plattenpanzer',        'armor',     0, 9, 0,  0),
('Siegelring-der-Stabilität',  'accessory', 0, 0, 1, 25);

INSERT INTO skills (name, beschreibung, rolle, damage_min, damage_max, heal, accuracy, target_type) VALUES
('Runen-Geschoss',             'Magisches Geschoss aus reiner Rune.',     'arkan',   12, 24,  0, 0.90, 'single_enemy'),
('Arkaner Bann',               'Flaechenschaden auf alle Gegner.',        'arkan',    8, 16,  0, 0.85, 'all_enemies'),
('Klärende-Annotation',        'Heilt einen Verbuendeten.',               'arkan',    0,  0, 20, 1.00, 'single_ally'),
('Datenklinge',                'Schnitt mit gebuendelten Datenstroemen.', 'druide',  10, 20,  0, 0.85, 'single_enemy'),
('Strukturwandel',             'Hoher Schaden durch Transformation.',     'druide',  14, 28,  0, 0.70, 'single_enemy'),
('Transformative-Regeneration','Heilt sich selbst.',                      'druide',   0,  0, 16, 1.00, 'self'),
('Architekten-Schlag',         'Solider Hammerschlag.',                   'schmied', 14, 26,  0, 0.85, 'single_enemy'),
('Schutz-Rune',                'Buff: +3 Verteidigung fuer alle.',        'schmied',  0,  0,  0, 1.00, 'all_allies'),
('Konstrukt-Schild',           'Buff: -50% Schaden fuer einen Ally.',     'schmied',  0,  0,  0, 1.00, 'single_ally');

INSERT INTO helden (name, rolle, max_hp, current_hp, attack, defense, speed, equipped_weapon_id, equipped_armor_id, equipped_accessory_id) VALUES
('Ron',      'arkan',   120, 120, 18,  8, 14,
  (SELECT id FROM equipment WHERE name='Pergament-Stab'),
  (SELECT id FROM equipment WHERE name='Runen-Gewand'),
  (SELECT id FROM equipment WHERE name='Tintenfass-Amulett')),
('Lumjan',   'druide',  100, 100, 14, 10, 16,
  (SELECT id FROM equipment WHERE name='Transformations-Kristall'),
  (SELECT id FROM equipment WHERE name='Datenstrom-Mantel'),
  (SELECT id FROM equipment WHERE name='Schema-Ring')),
('Florentin','schmied', 130, 130, 16, 16, 10,
  (SELECT id FROM equipment WHERE name='Architekten-Hammer'),
  (SELECT id FROM equipment WHERE name='Runen-Plattenpanzer'),
  (SELECT id FROM equipment WHERE name='Siegelring-der-Stabilität'));