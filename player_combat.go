package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

func (p *Player) Combat(m *Monster, isGate bool) bool {
	fmt.Printf("\n⚔️ [BATTLE COMMENCED]: %s\n", strings.ToUpper(m.Name))
	monsterHealth := m.Health; reader := bufio.NewReader(os.Stdin)
	tempDefense := 0; dimTurns := 0; immActive := p.HasSkill("immortality")
	
	appraisalActive := p.HasSkill("appraisal") || p.HasSkill("sariel")
	parallelActive := p.HasSkill("parallel_minds")

	for monsterHealth > 0 && p.Health > 0 {
		if appraisalActive {
			lvl := 1; if l, ok := p.SkillLevels["appraisal"]; ok { lvl = l }; if l, ok := p.SkillLevels["sariel"]; ok { lvl = l + 10 }
			fmt.Printf("\n--- 👁️ [ANALYSIS]: %s ---\n", strings.ToUpper(m.Name))
			fmt.Printf("   ❤️ VITALITY: %d/%d\n", monsterHealth, m.Health)
			fmt.Printf("   💥 THREAT LEVEL: %d DMG\n", m.Damage)
			if lvl >= 10 {
				fmt.Print("   🔮 LOOT PROBABILITY: ")
				for item, prob := range m.LootTable { fmt.Printf("[%s: %.0f%%] ", item, prob*100) }
				fmt.Println()
			}
		}

		fmt.Printf("\n--- [TURN START] (❤️ %d/%d | 🔮 %d/%d) ---\n", p.Health, p.MaxHealth, p.Magic, p.MaxMagic)
		
		squadAtk := 0
		for _, name := range p.Squad {
			for i := range p.Subordinates {
				if p.Subordinates[i].Name == name {
					squadAtk += p.Subordinates[i].Attack / 2
					fmt.Printf("🤝 [%s] provides fire support! (+%d support dmg)\n", name, p.Subordinates[i].Attack/2)
				}
			}
		}

		bonusDef := 0; critChance := 0.05
		for _, sID := range p.EquippedSkills {
			if sID == "great_sage" { bonusDef += 10 }; if sID == "raphael" { bonusDef += 50 }
			if sID == "critical_eye" { critChance += 0.2 }
		}

		fmt.Print("AVAILABLE ACTIONS: [!fight] ")
		for i, sID := range p.EquippedSkills {
			skill := GlobalSkills[sID]
			if skill.Type != "passive" {
				fmt.Printf("[!fight%d: %s (MP: %d)] ", i+1, skill.Name, skill.MPCost)
			}
		}
		fmt.Print("COMMAND: ")
		input, _ := reader.ReadString('\n'); input = strings.TrimSpace(input)
		if input == "!recover" { 
			p.HealFull()
			p.Stamina = p.MaxStamina
			fmt.Println("⚡ [CHEAT]: Existence restored to peak state (HP/MP/Stamina).")
			continue 
		}
		damage := 0; action := false; usedSkillName := ""

		if input == "!fight" {

			damage = p.Attack + p.GetEquippedWeaponDamage() + squadAtk + rand.Intn(5)
			if parallelActive { damage *= 2; fmt.Println("🧠 [PARALLEL MINDS]: Dual strike logic applied!") }
			if rand.Float64() <= critChance { damage = int(float64(damage) * 1.5); fmt.Println("🎯 [CRITICAL HIT]: Precise strike landed!") }
			usedSkillName = "Basic Attack"
			action = true
		} else if strings.HasPrefix(input, "!fight") {
			var idx int; fmt.Sscanf(input, "!fight%d", &idx); idx--
			if idx >= 0 && idx < len(p.EquippedSkills) {
				sID := p.EquippedSkills[idx]; skill := GlobalSkills[sID]; if skill.Type == "passive" { continue }
				if p.Magic < skill.MPCost { fmt.Println("❌ [SYSTEM]: Insufficient Mana for " + skill.Name); continue }
				
				p.Magic -= skill.MPCost; action = true; usedSkillName = skill.Name
				p.SkillUsage[sID]++; if p.SkillUsage[sID] >= 10 { p.UpgradeSkill(sID, true) }
				
				if (sID == "predator" || sID == "gluttony" || sID == "beelzebuth") && float64(monsterHealth)/float64(m.Health) < 0.3 {
					damage = monsterHealth; fmt.Printf("🌀 [SYSTEM]: %s successfully CONSUMED %s!\n", p.Name, m.Name)
					p.Attack += 2; p.MaxHealth += 10; p.GainTaboo(1)
				} else {
					lvl := p.SkillLevels[sID]; if lvl == 0 { lvl = 1 }
					damage = p.Attack + skill.DmgBonus + (skill.DmgBonus * (lvl - 1) / 2) + rand.Intn(10)
					if sID == "dim_maneuver" { dimTurns = 3; fmt.Println("🌌 [SYSTEM]: Dimensional Maneuver active. Evasion increased significantly.") }
					if skill.Category == "heal" { p.Health += int(float64(p.MaxHealth) * 0.2); fmt.Printf("💚 [HEAL]: %s restored Vitality.\n", skill.Name) }
					if skill.Category == "defense" { tempDefense = 50; fmt.Printf("🛡️ [DEFENSE]: %s boosted physical resistance.\n", skill.Name) }
				}
			}
		}

		if !action { continue }; monsterHealth -= damage
		fmt.Printf("💥 [%s]: Dealt %d damage to %s.\n", usedSkillName, damage, m.Name)
		
		if monsterHealth <= 0 {
			fmt.Printf("\n🏆 [VICTORY]: %s has been eliminated!\n", strings.ToUpper(m.Name))
			p.Kills++; p.MonsterKills[strings.ToLower(m.Name)]++
			p.CheckTitles(); p.TrackQuest("combat", m.Name, 1)
			for item, prob := range m.LootTable {
				if rand.Float64() <= prob {
					p.Inventory[item]++
					fmt.Printf("🎁 [LOOT]: Obtained rare item: %s\n", strings.ToUpper(item))
				}
			}
			if isGate { p.GainHunterXP(20 + rand.Intn(15)) } else { p.GainXP(15 + rand.Intn(10)) }
			return true
		}

		// Monster Turn
		if dimTurns > 0 && rand.Float64() < 0.5 {
			fmt.Println("🌌 [EVADED]: The attack passed through your current position!"); dimTurns--
		} else {
			if dimTurns > 0 { dimTurns-- }
			finalDmg := m.Damage - (p.Defense + p.GetEquippedArmorDefense() + bonusDef)
			if tempDefense > 0 { finalDmg = int(float64(finalDmg) * 0.5); tempDefense = 0 }
			if finalDmg < 1 { finalDmg = 1 }; p.Health -= finalDmg
			fmt.Printf("👹 [%s]: Counter-attacked for %d damage. (Vitality: %d left)\n", m.Name, finalDmg, p.Health)
		}
		
		if p.Health <= 0 && immActive {
			p.Health = 1; immActive = false
			p.WorldNotice("SKILL [IMMORTALITY] TRIGGERED: Death has been rejected.")
		}
	}
	
	if p.Health <= 0 {
		if p.HasSkill("egg_revival") {
			p.Health = p.MaxHealth; p.Magic = p.MaxMagic; p.WorldNotice("REINCARNATION: A new existence has been born from the abyss."); p.Save(); return false
		}
		if p.Inventory["life_stone"] > 0 {
			p.Inventory["life_stone"]--; p.Health = p.MaxHealth; p.Magic = p.MaxMagic; fmt.Println("\n💎 [SYSTEM]: LIFE STONE consumed. Existence restored."); p.Save(); return false
		}
		fmt.Println("\n💀 [FATAL]: You have succumbed to your injuries."); p.Inventory["gold"] = int(float64(p.Inventory["gold"]) * 0.5); p.Health = 50; p.Save()
	}
	return false
}
