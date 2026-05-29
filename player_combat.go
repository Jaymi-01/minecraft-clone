package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

func (p *Player) Combat(m *Monster, isGate bool) bool {
	p.InCombat = true
	defer func() { p.InCombat = false }()

	fmt.Printf("\n⚔️ [BATTLE COMMENCED]: %s\n", strings.ToUpper(m.Name))
	monsterHealth := m.Health; reader := bufio.NewReader(os.Stdin)
	tempDefense := 0; dimTurns := 0; immActive := p.HasSkill("immortality")
	
	appraisalActive := p.HasSkill("appraisal") || p.HasSkill("sariel")
	parallelActive := p.HasSkill("parallel_minds")

	// Initialize status effects for this combat session
	p.StatusEffects = make(map[string]int)
	m.StatusEffects = make(map[string]int)

	for monsterHealth > 0 && p.Health > 0 {
		// --- STATUS EFFECT PROCESSING (START OF TURN) ---
		if m.StatusEffects["burn"] > 0 {
			burnDmg := int(float64(m.Health) * 0.05)
			if burnDmg < 5 { burnDmg = 5 }
			monsterHealth -= burnDmg
			fmt.Printf("🔥 [STATUS]: %s suffers from BURN! (-%d vitality)\n", m.Name, burnDmg)
			m.StatusEffects["burn"]--
		}
		if p.StatusEffects["burn"] > 0 {
			burnDmg := int(float64(p.MaxHealth) * 0.05)
			p.Health -= burnDmg
			fmt.Printf("🔥 [STATUS]: You are BURNING! (-%d vitality)\n", burnDmg)
			p.StatusEffects["burn"]--
		}

		if monsterHealth <= 0 { break }

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

			if len(m.StatusEffects) > 0 {
				fmt.Print("   🌀 ACTIVE TACTICS: ")
				for eff, dur := range m.StatusEffects { if dur > 0 { fmt.Printf("[%s: %dT] ", strings.ToUpper(eff), dur) } }
				fmt.Println()
			}
		}

		fmt.Printf("\n--- [TURN START] (❤️ %d/%d | 🔮 %d/%d) ---\n", p.Health, p.MaxHealth, p.Magic, p.MaxMagic)
		
		// --- JOB PASSIVES (COMBAT) ---
		if p.Job == "True Demon Lord" {
			regen := int(float64(p.MaxHealth) * 0.05)
			p.Health += regen
			if p.Health > p.MaxHealth { p.Health = p.MaxHealth }
			fmt.Printf("👿 [INFINITE REGENERATION]: Chaos energy restored %d Vitality.\n", regen)
		}

		// --- SQUAD SUPPORT ---
		squadAtkTotal := 0
		if !p.TrialActive {
			for _, name := range p.Squad {
				for i := range p.Subordinates {
					if p.Subordinates[i].Name == name {
						s := &p.Subordinates[i]
						dmg := s.Attack / 2
						// Shadow Monarch Passive: Sovereign's Authority (Subordinates deal +100% damage)
						if p.Job == "Shadow Monarch" {
							dmg *= 2
						}
						squadAtkTotal += dmg
						fmt.Printf("🤝 [%s] unleashes a supporting strike! (-%d vitality)\n", name, dmg)
					}
				}
			}
		} else {
			fmt.Println("🚫 [SYSTEM]: Subordinate support is being suppressed by the Trial's aura.")
		}

		if squadAtkTotal > 0 {
			monsterHealth -= squadAtkTotal
			if monsterHealth <= 0 {
				fmt.Printf("🎯 [SQUAD FINISHER]: Your subordinates have eliminated %s before you could even act!\n", m.Name)
				p.VictoryLogic(m, isGate)
				return true
			}
		}

		// --- PLAYER ACTION ---
		if p.StatusEffects["paralyze"] > 0 {
			fmt.Println("⚡ [STATUS]: Your muscles are locked! Action skipped.")
			p.StatusEffects["paralyze"]--
		} else {
			bonusDef := 0; critChance := 0.05
			// Passive Logic: Check entire library (HasSkill) instead of Equipped
			if p.HasSkill("great_sage") { bonusDef += 10; critChance += 0.05 }
			if p.HasSkill("raphael") { bonusDef += 50; critChance += 0.1 }
			if p.HasSkill("ciel") { bonusDef += 200; critChance += 0.25 }
			if p.HasSkill("critical_eye") { critChance += 0.2 }

			fmt.Print("AVAILABLE ACTIONS: [!fight] ")
			for i, sID := range p.EquippedSkills {
				skill := GlobalSkills[sID]
				if skill.Type != "passive" { fmt.Printf("[!fight%d: %s (MP: %d)] ", i+1, skill.Name, skill.MPCost) }
			}
			
			actionResolved := false
			for !actionResolved {
				fmt.Print("\nCOMMAND: ")
				input, _ := reader.ReadString('\n'); input = strings.TrimSpace(input)
				if input == "!recover" { 
					p.HealFull(); p.Stamina = p.MaxStamina
					fmt.Println("⚡ [CHEAT]: Existence restored to peak state (HP/MP/Stamina).")
					continue 
				}
				
				damage := 0; action := false; usedSkillName := ""

				if input == "!fight" {
					damage = p.Attack + p.GetEquippedWeaponDamage() + rand.Intn(5)
					if parallelActive { damage *= 2; fmt.Println("🧠 [PARALLEL MINDS]: Dual strike logic applied!") }
					if rand.Float64() <= critChance { damage = int(float64(damage) * 1.5); fmt.Println("🎯 [CRITICAL HIT]: Precise strike landed!") }
					usedSkillName = "Basic Attack"; action = true
				} else if strings.HasPrefix(input, "!fight") {
					var idx int; fmt.Sscanf(input, "!fight%d", &idx); idx--
					if idx >= 0 && idx < len(p.EquippedSkills) {
						sID := p.EquippedSkills[idx]; skill := GlobalSkills[sID]; if skill.Type == "passive" { continue }
						if p.Magic < skill.MPCost { fmt.Println("❌ [SYSTEM]: Insufficient Mana!"); continue }
						
						p.Magic -= skill.MPCost; action = true; usedSkillName = skill.Name
						p.SkillUsage[sID]++; if p.SkillUsage[sID] >= 10 { p.UpgradeSkill(sID, true) }
						
						if (sID == "predator" || sID == "gluttony" || sID == "beelzebuth") && float64(monsterHealth)/float64(m.Health) < 0.3 {
							damage = monsterHealth; fmt.Printf("🌀 [SYSTEM]: %s successfully CONSUMED %s!\n", p.Name, m.Name)
							p.Attack += 2; p.MaxHealth += 10; p.GainTaboo(1)
						} else {
							lvl := p.SkillLevels[sID]; if lvl == 0 { lvl = 1 }
							damage = p.Attack + skill.DmgBonus + (skill.DmgBonus * (lvl - 1) / 2) + rand.Intn(10)

							if p.Job == "Abyssal Being" {
								fmt.Println("✨ [SPATIAL DOMINION]: Your strike bypassed the enemy's spatial defenses!")
								damage = int(float64(damage) * 1.25)
							}

							if sID == "dim_maneuver" { dimTurns = 3; fmt.Println("🌌 [SYSTEM]: Dimensional Maneuver active.") }
							if skill.Category == "heal" { p.Health += int(float64(p.MaxHealth) * 0.2); fmt.Printf("💚 [HEAL]: %s restored Vitality.\n", skill.Name) }
							if skill.Category == "defense" { tempDefense = 50; fmt.Printf("🛡️ [DEFENSE]: %s boosted resistance.\n", skill.Name) }
							
							// Apply Status Effect to Monster
							if skill.StatusEffect != "" && rand.Float64() <= skill.StatusChance {
								m.StatusEffects[skill.StatusEffect] = 3
								fmt.Printf("🌀 [TACTICAL]: %s applied %s to %s!\n", skill.Name, strings.ToUpper(skill.StatusEffect), m.Name)
							}
						}
					}
				}

				if action {
					if p.Job == "Abyssal Being" {
						damage = int(float64(damage) * 1.25) 
					}

					monsterHealth -= damage
					fmt.Printf("💥 [%s]: Dealt %d damage to %s.\n", usedSkillName, damage, m.Name)

					// Rune Effect: Lifesteal
					hasLifesteal := false
					for _, rID := range p.ItemRunes[p.EquippedWeapon] { if rID == "lifesteal_rune" { hasLifesteal = true; break } }
					if hasLifesteal {
						heal := int(float64(damage) * 0.05)
						if heal > 0 {
							p.Health += heal
							if p.Health > p.MaxHealth { p.Health = p.MaxHealth }
							fmt.Printf("🩸 [LIFESTEAL]: Absorbed %d vitality from the target.\n", heal)
						}
					}
					actionResolved = true
				} else {
					fmt.Println("❌ [SYSTEM]: Invalid action.")
				}
			}
		}

		if monsterHealth <= 0 {
			p.VictoryLogic(m, isGate)
			return true
		}

		// --- MONSTER TURN ---
		if m.StatusEffects["paralyze"] > 0 {
			fmt.Printf("⚡ [STATUS]: %s is paralyzed! Turn skipped.\n", m.Name)
			m.StatusEffects["paralyze"]--
		} else {
			if dimTurns > 0 && rand.Float64() < 0.5 {
				fmt.Println("🌌 [EVADED]: The attack passed through your current position!"); dimTurns--
			} else {
				if dimTurns > 0 { dimTurns-- }
				currentDef := p.Defense + p.GetEquippedArmorDefense()
				if m.StatusEffects["restrain"] > 0 { currentDef /= 2 }
				
				finalDmg := m.Damage - currentDef
				if tempDefense > 0 { finalDmg = int(float64(finalDmg) * 0.5); tempDefense = 0 }
				if finalDmg < 1 { finalDmg = 1 }; p.Health -= finalDmg
				fmt.Printf("👹 [%s]: Counter-attacked for %d damage. (Vitality: %d left)\n", m.Name, finalDmg, p.Health)
				
				// Monster chance to apply Burn
				if strings.Contains(m.Name, "🔥") && rand.Float64() < 0.2 {
					p.StatusEffects["burn"] = 3
					fmt.Println("🔥 [WARNING]: You have been set ABLAZE!")
				}
			}
		}
		
		if p.Health <= 0 && immActive {
			p.Health = 1; immActive = false
			p.WorldNotice("SKILL [IMMORTALITY] TRIGGERED: Death has been rejected.")
		}

		// --- TITLE PERK: THE ONE WHO OVERCOMES DEATH ---
		hasOvercomerTitle := false
		for _, t := range p.Titles { if t == "the_one_who_overcomes_death" { hasOvercomerTitle = true; break } }
		if hasOvercomerTitle && p.Health < (p.MaxHealth/10) && p.Health > 0 {
			p.HealFull()
			p.WorldNotice("🏆 [TITLE PERK]: THE ONE WHO OVERCOMES DEATH HAS SURPASSED THEIR LIMITS!")
			fmt.Println("✨ [SYSTEM]: Full Vitality and Mana restored via Title Authority.")
		}
	}
	
	if p.Health <= 0 { p.DeathLogic() }
	return false
}

func (p *Player) VictoryLogic(m *Monster, isGate bool) {
	fmt.Printf("\n🏆 [VICTORY]: %s has been eliminated!\n", strings.ToUpper(m.Name))
	p.Kills++; p.MonsterKills[strings.ToLower(m.Name)]++
	p.CheckTitles(); p.TrackQuest("combat", m.Name, 1)
	for item, prob := range m.LootTable {
		if rand.Float64() <= prob {
			p.Inventory[item]++
			fmt.Printf("🎁 [LOOT]: Obtained rare item: %s\n", strings.ToUpper(item))
		}
	}
	if isGate { 
		p.GainHunterXP(20 + rand.Intn(15))
		p.AriseActive = true
		bossCopy := *m
		p.AriseMonster = &bossCopy
		p.WorldNotice("SHADOW EXTRACTION AVAILABLE: Defeated entity has left a mana signature. Use !arise.")
	} else { 
		p.GainXP(15 + rand.Intn(10)) 
	}
}

func (p *Player) DeathLogic() {
	if p.HasSkill("egg_revival") {
		p.Health = p.MaxHealth; p.Magic = p.MaxMagic
		p.WorldNotice("REINCARNATION: A new existence has been born from the abyss.")
	} else if p.Inventory["life_stone"] > 0 {
		p.Inventory["life_stone"]--; p.Health = p.MaxHealth; p.Magic = p.MaxMagic
		fmt.Println("\n💎 [SYSTEM]: LIFE STONE consumed. Existence restored.")
	} else {
		fmt.Println("\n💀 [FATAL]: You have succumbed to your injuries.")
		p.Inventory["gold"] = int(float64(p.Inventory["gold"]) * 0.5)
		p.Health = 50
	}
	p.Save()
}
