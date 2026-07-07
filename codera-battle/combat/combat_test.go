package combat

import (
	"testing"

	"codera-battle/dragon"
	"codera-battle/hero"
	"codera-battle/internal"
)

func testHero(role string, base internal.Stats, skills []internal.Skill, equip []internal.Equipment) *hero.Hero {
	h := hero.New(hero.Definition{Name: "Test-" + role, Role: role, Base: base, Skills: skills, Equipment: equip})
	h.Strategy = nil // Tests steuern die Aktionen selbst
	return h
}

func TestCalculateDamage_Miss(t *testing.T) {
	_, _, missed := CalculateDamage(10, 20, 20, 10, 0.0)
	if !missed {
		t.Fatal("bei accuracy 0.0 muss der Angriff verfehlen")
	}
}

func TestCalculateDamage_HitMinimum(t *testing.T) {
	for i := 0; i < 1000; i++ {
		dmg, _, missed := CalculateDamage(10, 20, 20, 200, 1.0)
		if missed {
			t.Fatal("bei accuracy 1.0 darf nie verfehlt werden")
		}
		if dmg < 1 {
			t.Fatalf("Schaden muss mindestens 1 sein, war %d", dmg)
		}
	}
}

func TestHeroEquipmentBonuses(t *testing.T) {
	h := testHero("X",
		internal.Stats{MaxHP: 100, Attack: 10, Defense: 5, Speed: 10},
		nil,
		[]internal.Equipment{{Name: "W", Type: internal.EquipWeapon, Bonus: internal.Stats{Attack: 8, MaxHP: 20}}},
	)
	if got := h.GetMaxHP(); got != 120 {
		t.Errorf("MaxHP erwartet 120, war %d", got)
	}
	if got := h.GetStats().Attack; got != 18 {
		t.Errorf("Attack erwartet 18, war %d", got)
	}
}

func TestHeroHealClampsToMax(t *testing.T) {
	h := testHero("H", internal.Stats{MaxHP: 100, Attack: 10, Defense: 5, Speed: 10}, nil, nil)
	h.SetCurrentHP(90)
	h.SetCurrentHP(h.GetCurrentHP() + 50)
	if h.GetCurrentHP() != 100 {
		t.Errorf("HP darf MaxHP nicht überschreiten, war %d", h.GetCurrentHP())
	}
}

func TestHeroDeath(t *testing.T) {
	h := testHero("D", internal.Stats{MaxHP: 50, Attack: 10, Defense: 5, Speed: 10}, nil, nil)
	h.SetCurrentHP(0)
	if h.IsAlive() {
		t.Error("Held mit 0 HP darf nicht mehr leben")
	}
}

func TestEffectExpiry(t *testing.T) {
	h := testHero("E", internal.Stats{MaxHP: 100, Attack: 10, Defense: 5, Speed: 10}, nil, nil)
	h.AddEffect(0, 5, 0, 1)
	if h.GetStats().Defense != 10 {
		t.Fatalf("Defense-Buff nicht aktiv, war %d", h.GetStats().Defense)
	}
	h.TickEffects()
	if h.GetStats().Defense != 5 {
		t.Errorf("Defense-Buff nach Tick noch aktiv, war %d", h.GetStats().Defense)
	}
}

func TestInitiativeOrder_DragonTieGoesToHero(t *testing.T) {
	// Held mit Speed 14 (= Drachen-Speed) muss vor dem Drachen handeln.
	h := testHero("Tie", internal.Stats{MaxHP: 100, Attack: 10, Defense: 5, Speed: 14}, nil, nil)
	drg := dragon.New()
	order := buildInitiativeOrder([]*hero.Hero{h}, drg)
	if _, ok := order[0].(*hero.Hero); !ok {
		t.Error("bei Gleichstand mit dem Drachen muss der Held zuerst handeln")
	}
}

func TestInitiativeOrder_SortedBySpeed(t *testing.T) {
	slow := testHero("Slow", internal.Stats{Speed: 5}, nil, nil)
	fast := testHero("Fast", internal.Stats{Speed: 30}, nil, nil)
	drg := dragon.New()
	order := buildInitiativeOrder([]*hero.Hero{slow, fast}, drg)
	if order[0].GetName() != "Test-Fast" {
		t.Errorf("schnellster Kämpfer muss zuerst handeln, war %s", order[0].GetName())
	}
}

func TestApplyHeroDamage_ReducesDragonHP(t *testing.T) {
	skill := internal.Skill{Name: "Hit", MinDamage: 10, MaxDamage: 20, Accuracy: 1.0, Target: internal.TargetSingleEnemy}
	h := testHero("Atk", internal.Stats{MaxHP: 100, Attack: 20, Defense: 5, Speed: 10}, []internal.Skill{skill}, nil)
	drg := dragon.New()
	start := drg.GetCurrentHP()
	applyHeroDamage(h, skill, drg)
	if drg.GetCurrentHP() >= start {
		t.Errorf("Drachen-HP muss sinken: vorher %d, nachher %d", start, drg.GetCurrentHP())
	}
}

func TestApplyHeal_RestoresAlly(t *testing.T) {
	heal := internal.Skill{Name: "Heal", MinHeal: 20, MaxHeal: 20, Accuracy: 1.0, Target: internal.TargetSingleAlly}
	healer := testHero("Healer", internal.Stats{MaxHP: 100, Speed: 10}, []internal.Skill{heal}, nil)
	wounded := testHero("Wounded", internal.Stats{MaxHP: 100, Speed: 10}, nil, nil)
	wounded.SetCurrentHP(50)
	applyHeal(healer, heal, wounded, []*hero.Hero{healer, wounded})
	if wounded.GetCurrentHP() != 70 {
		t.Errorf("Verbündeter sollte auf 70 HP geheilt werden, war %d", wounded.GetCurrentHP())
	}
}

func TestLifeSteal_HealsAttacker(t *testing.T) {
	skill := internal.Skill{Name: "Steal", MinDamage: 30, MaxDamage: 40, Accuracy: 1.0, Target: internal.TargetSingleEnemy}
	h := testHero("System-Infiltrator*in",
		internal.Stats{MaxHP: 120, Attack: 30, Defense: 10, Speed: 20},
		[]internal.Skill{skill},
		[]internal.Equipment{{Name: "Dolch", Type: internal.EquipWeapon, Bonus: internal.Stats{Attack: 14}, Special: "life_steal"}},
	)
	h.SetCurrentHP(50)
	drg := dragon.New()
	applyHeroDamage(h, skill, drg)
	if h.GetCurrentHP() <= 50 {
		t.Errorf("Life Steal sollte den Angreifer heilen, HP war %d", h.GetCurrentHP())
	}
}
