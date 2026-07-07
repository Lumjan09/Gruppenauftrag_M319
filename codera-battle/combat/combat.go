// Package combat enthält den rundenbasierten Kampf-Loop gegen den Drachen.
//
// Vorgegeben sind CombatLoop, buildInitiativeOrder und CalculateDamage.
// Von den Lernenden implementiert werden processHeroTurn, processDragonTurn
// und logAction sowie die zugehörige Effekt-/Skill-Anwendung.
package combat

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"

	"codera-battle/dragon"
	"codera-battle/hero"
	"codera-battle/internal"
	"codera-battle/logging"
)

// inputReader liest die Spielereingaben; in Tests austauschbar.
var inputReader = bufio.NewReader(os.Stdin)

// CalculateDamage berechnet den Schaden eines Angriffs inkl. RNG, Genauigkeit
// und kritischem Treffer.
//
// Rückgabe: (Schaden, kritischerTreffer, verfehlt). Diese Funktion darf laut
// Auftrag nicht verändert werden.
func CalculateDamage(baseMin, baseMax, attackerStat, defenderDef int, accuracy float64) (int, bool, bool) {
	// 1. Genauigkeits-Check (RNG)
	if rand.Float64() > accuracy {
		return 0, false, true // Angriff verfehlt
	}
	// 2. Basisschaden (RNG innerhalb der Spanne)
	baseDamage := rand.Intn(baseMax-baseMin+1) + baseMin
	// 3. Angriffsbonus (Attacker-Stat / 20 als Multiplikator)
	attackMultiplier := 1.0 + float64(attackerStat)/20.0
	// 4. Verteidigungsreduktion
	defenseReduction := 1.0 - float64(defenderDef)/100.0
	if defenseReduction < 0.1 {
		defenseReduction = 0.1 // Minimum 10 % Schaden kommen immer durch
	}
	finalDamage := int(float64(baseDamage) * attackMultiplier * defenseReduction)
	if finalDamage < 1 {
		finalDamage = 1 // Minimum 1 Schaden
	}
	// 5. Kritischer Treffer (10 % Chance, 1.5x Schaden)
	isCrit := rand.Float64() < 0.1
	if isCrit {
		finalDamage = int(float64(finalDamage) * 1.5)
	}
	return finalDamage, isCrit, false
}

// buildInitiativeOrder sortiert alle Teilnehmer absteigend nach Speed.
// Bei Gleichstand mit dem Drachen handeln die Helden zuerst; Gleichstände
// zwischen Helden werden zufällig aufgelöst.
func buildInitiativeOrder(heroes []*hero.Hero, drg *dragon.Dragon) []internal.Combatant {
	order := make([]internal.Combatant, 0, len(heroes)+1)
	for _, h := range heroes {
		order = append(order, h)
	}
	order = append(order, drg)

	rand.Shuffle(len(order), func(i, j int) { order[i], order[j] = order[j], order[i] })
	sort.SliceStable(order, func(i, j int) bool {
		si, sj := order[i].GetStats().Speed, order[j].GetStats().Speed
		if si != sj {
			return si > sj
		}
		// Gleichstand: Held vor Drache
		_, iDragon := order[i].(*dragon.Dragon)
		_, jDragon := order[j].(*dragon.Dragon)
		return !iDragon && jDragon
	})
	return order
}

// CombatLoop führt den kompletten rundenbasierten Kampf aus.
func CombatLoop(heroes []*hero.Hero, drg *dragon.Dragon) {
	logging.Info("Kampf gestartet", "helden", len(heroes), "drache_hp", drg.GetCurrentHP())
	round := 1
	for {
		fmt.Printf("\n══════════════════ Runde %d ══════════════════\n", round)
		order := buildInitiativeOrder(heroes, drg)
		for _, c := range order {
			if !c.IsAlive() {
				continue
			}
			switch actor := c.(type) {
			case *dragon.Dragon:
				processDragonTurn(actor, heroes, round)
			case *hero.Hero:
				processHeroTurn(actor, heroes, drg, round)
			}
			if !drg.IsAlive() {
				printBattleResult(true, heroes, drg)
				return
			}
			if !anyHeroAlive(heroes) {
				printBattleResult(false, heroes, drg)
				return
			}
		}
		for _, h := range heroes {
			h.TickEffects()
		}
		drg.TickEffects()
		round++
	}
}

// processHeroTurn führt den Zug eines Helden aus: zuerst wird eine etwaige
// automatische Strategie geprüft, andernfalls erfolgt die CLI-Auswahl.
func processHeroTurn(h *hero.Hero, heroes []*hero.Hero, drg *dragon.Dragon, round int) {
	if h.Strategy != nil {
		if skill, target := h.Strategy(h, heroes, drg, round); skill != nil {
			fmt.Printf("\n» %s handelt automatisch: %s\n", h.GetName(), skill.Name)
			applyHeroSkill(h, *skill, target, heroes, drg, round)
			return
		}
	}

	printBattleHUD(h, heroes, drg, round)
	skill := readHeroChoice(h)
	target := defaultHeroTarget(skill, h, heroes, drg)
	applyHeroSkill(h, skill, target, heroes, drg, round)
}

// processDragonTurn führt den Zug des Drachen anhand seiner KI aus.
func processDragonTurn(drg *dragon.Dragon, heroes []*hero.Hero, round int) {
	skill := drg.ChooseAction(round)
	status := ""
	if drg.RageActive() {
		status = " [Rage +50%]"
	}
	fmt.Printf("\n🐉 %s setzt %s ein%s\n", drg.GetName(), skill.Name, status)
	logging.Debug("Drachen-Aktion", "skill", skill.Name, "rage", drg.RageActive(), "hp", drg.GetCurrentHP())

	switch skill.Target {
	case internal.TargetSelf: // Corrupted Code – Selbstheilung
		drg.Heal(skill.MaxHeal)
		logAction(fmt.Sprintf("%s heilt sich um %d HP", drg.GetName(), skill.MaxHeal))
	case internal.TargetAllEnemies: // Stack Overflow – AoE auf alle Helden
		for _, h := range heroes {
			if h.IsAlive() {
				dragonHitsHero(drg, h, skill)
			}
		}
	default: // Einzelziel
		if target := randomAliveHero(heroes); target != nil {
			dragonHitsHero(drg, target, skill)
		}
	}
}

// dragonHitsHero wendet einen Drachen-Angriff auf einen Helden an
// (inkl. Rage-Bonus und Schadensreduktion des Helden).
func dragonHitsHero(drg *dragon.Dragon, h *hero.Hero, skill internal.Skill) {
	dmg, crit, missed := CalculateDamage(skill.MinDamage, skill.MaxDamage,
		drg.GetStats().Attack, h.GetStats().Defense, skill.Accuracy)
	if missed {
		fmt.Printf("   → verfehlt %s!\n", h.GetName())
		logAction(fmt.Sprintf("%s verfehlt %s", drg.GetName(), h.GetName()))
		return
	}
	dmg = int(float64(dmg) * drg.RageFactor())
	if red := h.DamageReduction(); red > 0 {
		dmg = int(float64(dmg) * (1.0 - red))
	}
	h.SetCurrentHP(h.GetCurrentHP() - dmg)
	fmt.Printf("   → %d Schaden an %s%s (%d/%d HP)\n", dmg, h.GetName(), critSuffix(crit), h.GetCurrentHP(), h.GetMaxHP())
	logAction(fmt.Sprintf("%s trifft %s für %d Schaden", drg.GetName(), h.GetName(), dmg))
	if !h.IsAlive() {
		fmt.Printf("   ☠ %s wurde besiegt!\n", h.GetName())
		logging.Warn("Held gefallen", "held", h.GetName())
	} else if h.HPPercent() < 0.30 {
		logging.Warn("Held kritisch", "held", h.GetName(), "hp", h.GetCurrentHP())
	}
}

// applyHeroSkill wendet die gewählte Helden-Fähigkeit an.
func applyHeroSkill(h *hero.Hero, skill internal.Skill, target internal.Combatant, heroes []*hero.Hero, drg *dragon.Dragon, _ int) {
	switch {
	case skill.IsHeal():
		applyHeal(h, skill, target, heroes)
	case skill.IsDamage():
		applyHeroDamage(h, skill, drg)
	default:
		applyBuff(h, skill, target, heroes, drg)
	}
}

// applyHeroDamage berechnet und wendet den Schaden eines Helden auf den Drachen an.
// Der Funktions-Krieger nutzt bei "Präziser Hieb" einen parallelen Doppelangriff
// (Goroutines + WaitGroup), dessen Schaden via Mutex thread-sicher angewendet wird.
func applyHeroDamage(h *hero.Hero, skill internal.Skill, drg *dragon.Dragon) {
	atkStat := h.GetStats().Attack
	defStat := drg.GetStats().Defense
	belowThreshold := skill.DoubleDamageEnemyBelowPct > 0 && drg.HPPercent() < skill.DoubleDamageEnemyBelowPct

	doubleStrike := h.Role == "Funktions-Krieger*in" && skill.Name == "Präziser Hieb"
	hits := 1
	if doubleStrike {
		hits = 2
		fmt.Printf("   ⚔ Doppelangriff (parallel)!\n")
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	total := 0
	anyCrit := false
	for i := 0; i < hits; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			dmg, crit, missed := CalculateDamage(skill.MinDamage, skill.MaxDamage, atkStat, defStat, skill.Accuracy)
			if missed {
				return
			}
			if belowThreshold {
				dmg *= 2
			}
			mu.Lock()
			total += dmg
			anyCrit = anyCrit || crit
			mu.Unlock()
			drg.ApplyDamage(dmg) // thread-sicher per Mutex im Dragon
		}()
	}
	wg.Wait()

	if total == 0 {
		fmt.Printf("   → %s verfehlt den Drachen!\n", h.GetName())
		logAction(fmt.Sprintf("%s verfehlt mit %s", h.GetName(), skill.Name))
		return
	}
	if skill.EnemyDefenseDebuff > 0 {
		drg.ApplyDefenseDebuff(skill.EnemyDefenseDebuff, skill.BuffRounds)
	}
	if h.HasLifeSteal() {
		steal := total / 10
		if steal > 0 {
			h.SetCurrentHP(h.GetCurrentHP() + steal)
			fmt.Printf("   ♥ %s saugt %d HP ab\n", h.GetName(), steal)
		}
	}
	if skill.NextTurnAttackBuff > 0 {
		h.AddEffect(skill.NextTurnAttackBuff, 0, 0, skill.BuffRounds)
	}
	fmt.Printf("   → %s trifft den Drachen für %d Schaden%s (%d/%d HP)\n",
		h.GetName(), total, critSuffix(anyCrit), drg.GetCurrentHP(), drg.GetMaxHP())
	logAction(fmt.Sprintf("%s trifft Drache für %d Schaden mit %s", h.GetName(), total, skill.Name))
}

// applyHeal wendet eine Heilfähigkeit an (Selbst, Verbündeter oder alle).
func applyHeal(h *hero.Hero, skill internal.Skill, target internal.Combatant, heroes []*hero.Hero) {
	heal := rand.Intn(skill.MaxHeal-skill.MinHeal+1) + skill.MinHeal
	switch skill.Target {
	case internal.TargetAllAllies:
		for _, ally := range heroes {
			if ally.IsAlive() {
				ally.SetCurrentHP(ally.GetCurrentHP() + heal)
			}
		}
		fmt.Printf("   ✚ %s heilt das ganze Team um %d HP\n", h.GetName(), heal)
		logAction(fmt.Sprintf("%s heilt alle Helden um %d", h.GetName(), heal))
	default:
		t, ok := target.(*hero.Hero)
		if !ok {
			t = h
		}
		t.SetCurrentHP(t.GetCurrentHP() + heal)
		fmt.Printf("   ✚ %s heilt %s um %d HP (%d/%d)\n", h.GetName(), t.GetName(), heal, t.GetCurrentHP(), t.GetMaxHP())
		logAction(fmt.Sprintf("%s heilt %s um %d", h.GetName(), t.GetName(), heal))
	}
}

// applyBuff wendet reine Buff-/Schutzfähigkeiten an (ohne Schaden/Heilung).
func applyBuff(h *hero.Hero, skill internal.Skill, target internal.Combatant, heroes []*hero.Hero, _ *dragon.Dragon) {
	switch {
	case skill.SelfDefenseBuff > 0:
		h.AddEffect(0, skill.SelfDefenseBuff, 0, skill.BuffRounds)
		fmt.Printf("   🛡 %s erhöht die eigene Defense um %d\n", h.GetName(), skill.SelfDefenseBuff)
	case skill.AllyDefenseBuff > 0:
		for _, ally := range heroes {
			if ally.IsAlive() {
				ally.AddEffect(0, skill.AllyDefenseBuff, 0, skill.BuffRounds)
			}
		}
		fmt.Printf("   🛡 %s erhöht die Defense aller Helden um %d\n", h.GetName(), skill.AllyDefenseBuff)
	case skill.AllyDamageReduction > 0:
		t, ok := target.(*hero.Hero)
		if !ok {
			t = h
		}
		t.AddEffect(0, 0, skill.AllyDamageReduction, skill.BuffRounds)
		fmt.Printf("   🛡 %s schützt %s (-%.0f%% Schaden)\n", h.GetName(), t.GetName(), skill.AllyDamageReduction*100)
	}
	logAction(fmt.Sprintf("%s nutzt %s", h.GetName(), skill.Name))
}

// logAction protokolliert eine Kampfaktion auf Info-Stufe.
func logAction(msg string) { logging.Info(msg) }

// --- CLI-Hilfsfunktionen ------------------------------------------------------

// printBattleHUD zeichnet den Drachen-Status, die Team-HP und das Aktionsmenü.
func printBattleHUD(active *hero.Hero, heroes []*hero.Hero, drg *dragon.Dragon, round int) {
	status := ""
	if drg.RageActive() {
		status = "Status: Rage aktiv (+50% Schaden)"
	}
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Printf("║ %-40s ║\n", fmt.Sprintf("%s - HP: %d/%d", drg.GetName(), drg.GetCurrentHP(), drg.GetMaxHP()))
	if status != "" {
		fmt.Printf("║ %-40s ║\n", status)
	}
	fmt.Println("╚══════════════════════════════════════════╝")
	fmt.Printf("\nRunde %d - Zug von %s (%s)\n", round, active.GetName(), active.Role)
	fmt.Println("─── Team HP ───")
	for _, h := range heroes {
		fmt.Printf("  %-22s %3d/%-3d %s\n", h.GetName(), h.GetCurrentHP(), h.GetMaxHP(), hpBar(h))
	}
	fmt.Println("─── Aktionen ───")
	for i, s := range active.Skills {
		fmt.Printf("  %d. %s\n", i+1, skillLabel(s))
	}
	fmt.Print("Deine Wahl: ")
}

func skillLabel(s internal.Skill) string {
	switch {
	case s.IsHeal():
		return fmt.Sprintf("%s (heilt %d-%d)", s.Name, s.MinHeal, s.MaxHeal)
	case s.IsDamage():
		return fmt.Sprintf("%s (%d-%d Schaden, %.0f%%)", s.Name, s.MinDamage, s.MaxDamage, s.Accuracy*100)
	default:
		return fmt.Sprintf("%s (%s)", s.Name, s.Description)
	}
}

func hpBar(h *hero.Hero) string {
	switch p := h.HPPercent(); {
	case p < 0.30:
		return "▼▼"
	case p < 0.70:
		return "▼"
	default:
		return ""
	}
}

// readHeroChoice liest eine gültige Skill-Auswahl von der Standardeingabe.
func readHeroChoice(h *hero.Hero) internal.Skill {
	for {
		line, err := inputReader.ReadString('\n')
		if err != nil && line == "" {
			return h.Skills[0] // Fallback (z. B. EOF in nicht-interaktiver Umgebung)
		}
		choice, convErr := strconv.Atoi(strings.TrimSpace(line))
		if convErr == nil && choice >= 1 && choice <= len(h.Skills) {
			return h.Skills[choice-1]
		}
		fmt.Printf("Ungültige Eingabe. Bitte 1-%d wählen: ", len(h.Skills))
	}
}

// defaultHeroTarget bestimmt das Ziel für vom Spieler gewählte Skills.
func defaultHeroTarget(skill internal.Skill, self *hero.Hero, heroes []*hero.Hero, drg *dragon.Dragon) internal.Combatant {
	switch skill.Target {
	case internal.TargetSelf:
		return self
	case internal.TargetSingleAlly:
		if t := hero.LowestHPAlly(heroes); t != nil {
			return t
		}
		return self
	default:
		return drg
	}
}

func randomAliveHero(heroes []*hero.Hero) *hero.Hero {
	alive := make([]*hero.Hero, 0, len(heroes))
	for _, h := range heroes {
		if h.IsAlive() {
			alive = append(alive, h)
		}
	}
	if len(alive) == 0 {
		return nil
	}
	return alive[rand.Intn(len(alive))]
}

func anyHeroAlive(heroes []*hero.Hero) bool {
	for _, h := range heroes {
		if h.IsAlive() {
			return true
		}
	}
	return false
}

// critSuffix liefert eine Kennzeichnung für kritische Treffer.
func critSuffix(crit bool) string {
	if crit {
		return " (KRITISCH!)"
	}
	return ""
}

// printBattleResult gibt das Endergebnis des Kampfes aus.
func printBattleResult(victory bool, heroes []*hero.Hero, drg *dragon.Dragon) {
	fmt.Println("\n═══════════════════════════════════════════════")
	if victory {
		fmt.Println("🏆 SIEG! Der Entropie-Drache wurde besiegt!")
		logging.Info("Kampf beendet", "ergebnis", "sieg")
	} else {
		fmt.Printf("💀 NIEDERLAGE! Der Drache hat noch %d HP.\n", drg.GetCurrentHP())
		logging.Info("Kampf beendet", "ergebnis", "niederlage", "drache_hp", drg.GetCurrentHP())
	}
	fmt.Println("─── Überlebende ───")
	for _, h := range heroes {
		state := "gefallen"
		if h.IsAlive() {
			state = fmt.Sprintf("%d/%d HP", h.GetCurrentHP(), h.GetMaxHP())
		}
		fmt.Printf("  %-22s %s\n", h.GetName(), state)
	}
	fmt.Println("═══════════════════════════════════════════════")
}
