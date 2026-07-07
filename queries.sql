-- 1) Alle Helden mit Name und Rolle
SELECT name, rolle FROM helden ORDER BY id;

-- 2) Ausrüstung eines Helden via JOIN
SELECT h.name AS held, e.name AS gegenstand, e.typ
FROM helden h
JOIN equipment e ON e.id IN (h.equipped_weapon_id, h.equipped_armor_id, h.equipped_accessory_id)
WHERE h.name = 'Florentin';

-- 3) Alle Skills der Rolle krieger
SELECT name, damage_min, damage_max, heal, accuracy, target_type
FROM skills WHERE rolle = 'krieger';

-- 4) Helden mit allen 3 Slots belegt
SELECT name, rolle FROM helden
WHERE equipped_weapon_id IS NOT NULL
  AND equipped_armor_id IS NOT NULL
  AND equipped_accessory_id IS NOT NULL;

-- 5) Durchschnittlicher Angriff pro Rolle
SELECT rolle, ROUND(AVG(attack), 2) AS avg_attack
FROM helden GROUP BY rolle ORDER BY avg_attack DESC;