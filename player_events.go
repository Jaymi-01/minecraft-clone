package main

import (
	"fmt"
	"math/rand"
)

func (p *Player) TrackQuest(t, id string, q int) {
	for _, qst := range GlobalQuests {
		if qst.TargetType == t && qst.TargetID == id {
			p.QuestProgress[qst.ID] += q
			if p.QuestProgress[qst.ID] >= qst.TargetQty {
				p.WorldNotice("MISSION CONCLUDED: " + qst.Name)
				p.GainXP(qst.RewardXP); p.Inventory["gold"] += qst.RewardGold
			}
		}
	}
}

func (p *Player) ListQuests() {
	fmt.Println("\n--- 📜 Missions ---")
	for _, qst := range GlobalQuests {
		prog := p.QuestProgress[qst.ID]
		status := fmt.Sprintf("%d/%d", prog, qst.TargetQty)
		if prog >= qst.TargetQty { status = "✅ COMPLETED" }
		fmt.Printf("   [%s] %s: %s\n", qst.ID, qst.Name, status)
	}
}

func (p *Player) FoundChest(hasAnalysis bool) {
	loot := []string{"diamond", "void_essence", "life_stone", "health_potion", "star_matter"}
	item := loot[rand.Intn(len(loot))]; qty := 1 + rand.Intn(3)
	p.Inventory[item] += qty; fmt.Printf("🎁 Found a chest! %s x%d acquired.\n", item, qty)
	if hasAnalysis { p.AutoAnalyze(item) }
}

func (p *Player) TriggerTrap(hasAnalysis bool) {
	if hasAnalysis {
		for _, sID := range p.EquippedSkills { if sID == "trap_sense" { p.SkillUsage[sID]++; if p.SkillUsage[sID] >= 10 { p.UpgradeSkill(sID, true) } } }
		if rand.Float64() < 0.7 { fmt.Println("⚠️ [System]: TRAP DODGED."); return }
	}
	dmg := 10 + rand.Intn(15); p.Health -= dmg; fmt.Printf("💥 TRAP! %d DMG (❤️ %d HP left)\n", dmg, p.Health)
}

func (p *Player) EncounterMonster() {
	m := Monster{Name: "Labyrinth Beast", Health: 100 + (p.ExplorationDepth * 20), Damage: 20 + (p.ExplorationDepth * 5)}
	if p.Combat(&m, false) { p.GainXP(20); p.SubordinateGainXP(10) }
}
