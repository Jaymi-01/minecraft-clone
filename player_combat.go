package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

func (p *Player) Combat(m *Monster, isGate bool) bool {
	fmt.Printf("\n⚔️ ENCOUNTER: %s\n", m.Name); monsterHealth := m.Health; reader := bufio.NewReader(os.Stdin)
	tempDefense := 0; dimTurns := 0; immActive := p.HasSkill("immortality")
	
	appraisalActive := p.HasSkill("appraisal") || p.HasSkill("sariel")
	parallelActive := p.HasSkill("parallel_minds")

	for monsterHealth > 0 && p.Health > 0 {
		if appraisalActive {
			lvl := 1; if l, ok := p.SkillLevels["appraisal"]; ok { lvl = l }; if l, ok := p.SkillLevels["sariel"]; ok { lvl = l + 10 }
			fmt.Printf("\n--- 👁️ APPRAISAL: %s [HP: %d/%d | DMG: %d] ---\n", m.Name, monsterHealth, m.Health, m.Damage)
			if lvl >= 10 { fmt.Print("   🔮 ANALYSIS: "); for item, prob := range m.LootTable { fmt.Printf("%s (%.0f%%) ", item, prob*100) }; fmt.Println() }
		}

		fmt.Printf("Turn (❤️ %d/%d | 🔮 %d/%d)\n", p.Health, p.MaxHealth, p.Magic, p.MaxMagic)
		
		squadAtk := 0
		for _, name := range p.Squad {
			for i := range p.Subordinates { if p.Subordinates[i].Name == name { squadAtk += p.Subordinates[i].Attack / 2; fmt.Printf("🤝 %s support!\n", name) } }
		}

		bonusDef := 0; critChance := 0.05
		for _, sID := range p.EquippedSkills {
			if sID == "great_sage" { bonusDef += 10 }; if sID == "raphael" { bonusDef += 50 }
			if sID == "critical_eye" { critChance += 0.2 }
		}

		fmt.Print("Choice (!fight, !fight1...): "); input, _ := reader.ReadString('\n'); input = strings.TrimSpace(input)
		damage := 0; action := false

		if input == "!fight" {
			damage = p.Attack + p.GetEquippedWeaponDamage() + squadAtk + rand.Intn(5)
			if parallelActive { damage *= 2; fmt.Println("🧠 Parallel Minds hit!") }
			if rand.Float64() <= critChance { damage = int(float64(damage) * 1.5); fmt.Println("🎯 CRITICAL!") }
			action = true
		} else if strings.HasPrefix(input, "!fight") {
			var idx int; fmt.Sscanf(input, "!fight%d", &idx); idx--
			if idx >= 0 && idx < len(p.EquippedSkills) {
				sID := p.EquippedSkills[idx]; skill := GlobalSkills[sID]; if skill.Type == "passive" { continue }
				if p.Magic < skill.MPCost { fmt.Println("❌ Low Magic!"); continue }
				p.Magic -= skill.MPCost; action = true
				p.SkillUsage[sID]++; if p.SkillUsage[sID] >= 10 { p.UpgradeSkill(sID, true) }
				
				if (sID == "predator" || sID == "gluttony" || sID == "beelzebuth") && float64(monsterHealth)/float64(m.Health) < 0.3 {
					damage = monsterHealth; fmt.Println("🌀 CONSUMED!"); p.Attack += 2; p.MaxHealth += 10; p.GainTaboo(1)
				} else {
					lvl := p.SkillLevels[sID]; if lvl == 0 { lvl = 1 }
					damage = p.Attack + skill.DmgBonus + (skill.DmgBonus * (lvl - 1) / 2) + rand.Intn(10)
					if sID == "dim_maneuver" { dimTurns = 3; fmt.Println("🌌 Spatial dodge active!") }
					if skill.Category == "heal" { p.Health += int(float64(p.MaxHealth) * 0.2); fmt.Println("💚 Healed.") }
					if skill.Category == "defense" { tempDefense = 50 }
				}
			}
		}

		if !action { continue }; monsterHealth -= damage; fmt.Printf("💥 Dealt %d damage.\n", damage)
		if monsterHealth <= 0 {
			fmt.Printf("🏆 Defeated %s!\n", m.Name); p.Kills++; p.MonsterKills[strings.ToLower(m.Name)]++
			p.CheckTitles(); p.TrackQuest("combat", m.Name, 1)
			for item, prob := range m.LootTable { if rand.Float64() <= prob { p.Inventory[item]++; fmt.Printf("🎁 Dropped: %s\n", item) } }
			if isGate { p.GainHunterXP(20 + rand.Intn(15)) } else { p.GainXP(15 + rand.Intn(10)) }
			return true
		}

		if dimTurns > 0 && rand.Float64() < 0.5 { fmt.Println("🌌 EVADED!"); dimTurns-- } else {
			if dimTurns > 0 { dimTurns-- }
			finalDmg := m.Damage - (p.Defense + p.GetEquippedArmorDefense() + bonusDef)
			if tempDefense > 0 { finalDmg = int(float64(finalDmg) * 0.5) }
			if finalDmg < 1 { finalDmg = 1 }; p.Health -= finalDmg
			fmt.Printf("👹 %s hits for %d. (%d HP left)\n", m.Name, finalDmg, p.Health)
		}
		if p.Health <= 0 && immActive { p.Health = 1; immActive = false; p.WorldNotice("IMMORTALITY TRIGGERED") }
	}
	
	if p.Health <= 0 {
		if p.HasSkill("egg_revival") { p.Health = p.MaxHealth; p.Magic = p.MaxMagic; p.WorldNotice("REINCARNATED."); p.Save(); return false }
		if p.Inventory["life_stone"] > 0 { p.Inventory["life_stone"]--; p.Health = p.MaxHealth; p.Magic = p.MaxMagic; fmt.Println("\n💎 REVIVED!"); p.Save(); return false }
		fmt.Println("\n💀 DIED."); p.Inventory["gold"] = int(float64(p.Inventory["gold"]) * 0.5); p.Health = 50; p.Save()
	}
	return false
}
