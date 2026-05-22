package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func NewPlayer(name string) *Player {
	return &Player{
		Name:           name,
		Health:         100,
		MaxHealth:      100,
		Attack:         10,
		Defense:        0,
		Stamina:        50,
		MaxStamina:     50,
		Magic:          100,
		MaxMagic:       100,
		Level:          1,
		XP:             0,
		XPToNext:       100,
		HunterLevel:    1,
		HunterXP:       0,
		HunterXPToNext: 100,
		HunterRank:     "E",
		Inventory:      map[string]int{"wood_pickaxe": 1, "gold": 100},
		ToolDurability: 50,
		Structures:     make(map[string]bool),
		QuestProgress:  make(map[string]int),
		Rank:           "E",
		Kills:          0,
		MonsterKills:   make(map[string]int),
		Taboo:          0,
		SkillPoints:    6,
		Titles:         []string{},
		Skills:         []string{},
		EquippedSkills: []string{},
		SkillSlots:     5,
		SkillLevels:    make(map[string]int),
		SkillUsage:     make(map[string]int),
		SkillCooldowns: make(map[string]int),
		Subordinates:   []Subordinate{},
		SystemOrigin:   "Human",
	}
}

func (p *Player) WorldNotice(msg string) {
	fmt.Printf("\n<< NOTICE: %s >>\n", strings.ToUpper(msg))
}

func (p *Player) UpdateRank() {
	if p.Level >= 150 {
		p.Rank = "SS"
	} else if p.Level >= 100 {
		p.Rank = "S"
	} else if p.Level >= 75 {
		p.Rank = "A"
	} else if p.Level >= 50 {
		p.Rank = "B"
	} else if p.Level >= 30 {
		p.Rank = "C"
	} else if p.Level >= 15 {
		p.Rank = "D"
	} else {
		p.Rank = "E"
	}
}

func (p *Player) UpdateHunterRank() {
	if p.HunterLevel >= 150 {
		p.HunterRank = "SS"
	} else if p.HunterLevel >= 100 {
		p.HunterRank = "S"
	} else if p.HunterLevel >= 75 {
		p.HunterRank = "A"
	} else if p.HunterLevel >= 50 {
		p.HunterRank = "B"
	} else if p.HunterLevel >= 30 {
		p.HunterRank = "C"
	} else if p.HunterLevel >= 15 {
		p.HunterRank = "D"
	} else {
		p.HunterRank = "E"
	}
}

func (p *Player) UpdateSkillSlots() {
	p.SkillSlots = 5 + (p.Level / 5)
}

func (p *Player) GainXP(amount int) {
	if p.Structures["enchanting_table"] {
		amount = int(float64(amount) * 1.5)
	}
	p.XP += amount
	fmt.Printf("[✨ +%d Mine XP]\n", amount)
	for p.XP >= p.XPToNext {
		p.Level++
		p.XP -= p.XPToNext
		p.XPToNext = int(float64(p.XPToNext) * 1.5)
		p.MaxHealth += 10
		p.MaxStamina += 10
		p.MaxMagic += 20
		p.Health = p.MaxHealth
		p.Stamina = p.MaxStamina
		p.Magic = p.MaxMagic
		p.UpdateRank()
		p.UpdateSkillSlots()
		p.TrackQuest("level", "mine", p.Level)
		p.WorldNotice(fmt.Sprintf("Individual '%s' has reached Mine Level %d", p.Name, p.Level))
	}
	p.Save()
}

func (p *Player) GainHunterXP(amount int) {
	p.HunterXP += amount
	fmt.Printf("[🏹 +%d Hunter XP]\n", amount)
	for p.HunterXP >= p.HunterXPToNext {
		oldRank := p.HunterRank
		p.HunterLevel++
		p.HunterXP -= p.HunterXPToNext
		p.HunterXPToNext = int(float64(p.HunterXPToNext) * 1.5)
		if p.HunterLevel%20 == 0 {
			p.SkillSlots++
			p.WorldNotice("Skill Capacity Increased")
		}
		p.UpdateHunterRank()
		p.TrackQuest("level", "hunter", p.HunterLevel)
		p.WorldNotice(fmt.Sprintf("Hunter Level %d achieved", p.HunterLevel))
		if p.HunterRank != oldRank {
			p.WorldNotice(fmt.Sprintf("Hunter Rank promotion to %s successful", p.HunterRank))
		}
	}
	p.Save()
}

func (p *Player) StartGateSpawning() {

	ticker := time.NewTicker(10 * time.Minute)
	go func() {
		for range ticker.C {
			p.SpawnGate()
		}
	}()
	p.SpawnGate()
}

func (p *Player) SpawnGate() {
	ranks := []string{"E", "D", "C", "B", "A", "S", "SS"}
	rank := ranks[rand.Intn(len(ranks))]
	p.ManualSpawnGate(rank)
}

func (p *Player) ManualSpawnGate(rank string) {
	rank = strings.ToUpper(rank)
	gate, ok := Gates[rank]
	if !ok {
		fmt.Printf("❌ Invalid rank: %s.\n", rank)
		return
	}
	possibleBosses := GateBosses[rank]
	if len(possibleBosses) > 0 {
		gate.Boss = possibleBosses[rand.Intn(len(possibleBosses))]
	}
	p.CurrentGate = &gate
	p.WorldNotice(fmt.Sprintf("A %s-Rank Gate has manifested. Boss: %s", rank, gate.Boss.Name))
}

func (p *Player) ChooseOrigin(origin string) {
	origin = strings.ToLower(origin)
	if p.SystemOrigin != "Human" {
		fmt.Printf("❌ Origin already fixed as: %s\n", p.SystemOrigin)
		return
	}
	if p.Level < 5 {
		fmt.Println("🚫 Minimum Level 5 required for System Integration.")
		return
	}
	switch origin {
	case "slime":
		p.SystemOrigin = "Slime"
		p.AddSkill("predator")
		p.AddSkill("great_sage")
		p.WorldNotice("Unique Path: SLIME - Predator and Great Sage acquired")
	case "spider":
		p.SystemOrigin = "Spider"
		p.AddSkill("appraisal")
		p.AddSkill("spider_thread")
		p.WorldNotice("Unique Path: SPIDER - Appraisal and Spider Thread acquired")
	default:
		fmt.Println("❓ Choice unavailable.")
		return
	}
	p.Save()
}

func (p *Player) Evolve() {
	if p.SystemOrigin == "Human" {
		fmt.Println("❌ Origin required.")
		return
	}
	evolved := false
	switch p.SystemOrigin {
	case "Slime":
		if p.Level >= 30 {
			p.SystemOrigin = "Demon Slime"
			p.AddSkill("megiddo"); p.AddSkill("raphael"); p.AddSkill("beelzebuth")
			p.MaxHealth += 500; p.Attack += 50; p.MaxMagic += 500
			p.WorldNotice("Evolution to DEMON SLIME successful.")
			evolved = true
		}
	case "Demon Slime":
		if p.Level >= 60 {
			p.SystemOrigin = "Ultimate Slime (True Dragon)"
			p.MaxHealth += 2000; p.Attack += 200; p.MaxMagic += 2000
			p.WorldNotice("Individual has ascended to TRUE DRAGON status.")
			evolved = true
		}
	case "Spider":
		if p.Level >= 30 {
			p.SystemOrigin = "Arachne"
			p.AddSkill("wisdom"); p.AddSkill("evil_eye"); p.AddSkill("parallel_minds")
			p.MaxHealth += 300; p.Attack += 80; p.MaxMagic += 200
			p.WorldNotice("Evolution to ARACHNE successful.")
			evolved = true
		}
	case "Arachne":
		if p.Level >= 60 {
			p.SystemOrigin = "God (Shiraori)"
			p.MaxHealth += 1500; p.Attack += 300; p.MaxMagic += 5000
			p.WorldNotice("Apotheosis complete. Individual has achieved DIVINITY.")
			evolved = true
		}
	}
	if !evolved {
		fmt.Println("⚠️ Requirements insufficient.")
	} else {
		p.Health = p.MaxHealth; p.Magic = p.MaxMagic; p.Save()
	}
}

func (p *Player) GainTaboo(amount int) {
	p.Taboo += amount
	p.WorldNotice(fmt.Sprintf("Taboo Level increased by %d (Current: %d)", amount, p.Taboo))
	if p.Taboo == 10 {
		p.WorldNotice("Individual has crossed the threshold of Forbidden Knowledge.")
	}
	p.CheckTitles()
}

func (p *Player) Combat(m *Monster, isGate bool) bool {
	fmt.Printf("\n⚔️ ENCOUNTER: %s\n", m.Name)
	monsterHealth := m.Health
	reader := bufio.NewReader(os.Stdin)
	p.SkillCooldowns = make(map[string]int)
	tempDefense := 0
	damageMultiplier := 1.0
	rulerPrideTurns := 0
	appraisalActive := false
	parallelMindsActive := false

	for _, sID := range p.EquippedSkills {
		if sID == "appraisal" || sID == "wisdom" { appraisalActive = true }
		if sID == "parallel_minds" { parallelMindsActive = true }
		
		// Passive Usage: Combat Passives
		skill := GlobalSkills[sID]
		passives := map[string]bool{"appraisal":true, "wisdom":true, "great_sage":true, "raphael":true, "critical_eye":true, "battle_hardened":true, "eclipse_affinity":true}
		if skill.Type == "passive" && passives[sID] {
			p.SkillUsage[sID]++
			if p.SkillUsage[sID] >= 10 { p.UpgradeSkill(sID, true) }
		}
	}

	for monsterHealth > 0 && p.Health > 0 {
		tempDefense = 0 
		if rulerPrideTurns > 0 {
			damageMultiplier = 3.0; rulerPrideTurns--
			drain := int(float64(p.MaxHealth) * 0.1); p.Health -= drain
			fmt.Printf("👑 [RULER OF PRIDE] 3x Damage! HP Drained: %d\n", drain)
		} else { damageMultiplier = 1.0 }

		if appraisalActive {
			lvl := 1
			if l, ok := p.SkillLevels["appraisal"]; ok { lvl = l }
			if l, ok := p.SkillLevels["wisdom"]; ok { lvl = l + 10 }

			fmt.Printf("\n--- 👁️ APPRAISAL: %s [HP: %d/%d | DMG: %d] ---\n", m.Name, monsterHealth, m.Health, m.Damage)
			if lvl >= 4 {
				fmt.Printf("   💰 EXPECTED REWARDS: ~%d XP\n", 15+p.Level)
			}
			if lvl >= 7 {
				fmt.Print("   🎁 DROP TABLE: ")
				for item := range m.LootTable { fmt.Printf("[%s] ", strings.Replace(item, "_", " ", -1)) }
				fmt.Println()
			}
			if lvl >= 10 {
				fmt.Print("   🔮 ANALYSIS: ")
				for item, prob := range m.LootTable { fmt.Printf("%s (%.0f%%) ", strings.Replace(item, "_", " ", -1), prob*100) }
				fmt.Println()
			}
		}
		fmt.Printf("--- Your Turn (❤️ %d/%d | 🔮 %d/%d) ---\n", p.Health, p.MaxHealth, p.Magic, p.MaxMagic)
		for sID, cd := range p.SkillCooldowns {
			if cd > 0 { p.SkillCooldowns[sID]-- }
		}

		bonusAtkFromPassives := 0
		critChance := 0.05
		for _, sID := range p.EquippedSkills {
			skill := GlobalSkills[sID]
			if skill.Type == "passive" {
				if sID == "critical_eye" { critChance += 0.15 }
				if sID == "great_sage" { p.Defense += 10; critChance += 0.1 }
				if sID == "raphael" { p.Defense += 50; critChance += 0.3 }
				if sID == "battle_hardened" {
					missingHPPercent := float64(p.MaxHealth-p.Health) / float64(p.MaxHealth)
					bonusAtkFromPassives += int(float64(p.Attack) * missingHPPercent)
				}
			}
		}

		fmt.Print("Actions: [!fight] Attack ")
		for i, sID := range p.EquippedSkills {
			skill := GlobalSkills[sID]
			if skill.Type == "passive" { continue }
			cd := p.SkillCooldowns[sID]
			status := "READY"; if cd > 0 { status = fmt.Sprintf("%d turns", cd) }
			fmt.Printf("[!fight%d] %s (%s, %d MP) ", i+1, skill.Name, status, skill.MPCost)
		}
		fmt.Print("\nChoice: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "!recover" { p.HealFull(); fmt.Println("⚡ [CHEAT] Restored!"); continue }

		damageToMonster := 0; actionTaken := false
		baseAtk := int(float64(p.Attack+p.GetBestSwordDamage()+bonusAtkFromPassives) * damageMultiplier)
		if isGate { baseAtk += p.HunterLevel * 2 } else { baseAtk += p.Level }

		if input == "!fight" {
			damageToMonster = baseAtk + rand.Intn(5)
			if parallelMindsActive { damageToMonster *= 2; fmt.Println("🧠 Parallel Minds hit!") }
			if rand.Float64() <= critChance { damageToMonster = int(float64(damageToMonster) * 1.5); fmt.Println("🎯 CRITICAL!") }
			fmt.Println("🤜 Basic attack."); actionTaken = true
		} else if strings.HasPrefix(input, "!fight") {
			var idx int; fmt.Sscanf(strings.TrimPrefix(input, "!fight"), "%d", &idx); idx--
			if idx >= 0 && idx < len(p.EquippedSkills) {
				sID := p.EquippedSkills[idx]; skill := GlobalSkills[sID]
				if skill.Type == "passive" { continue }
				if p.Magic < skill.MPCost { fmt.Println("❌ Insufficient Magic (MP)!"); continue }
				if p.SkillCooldowns[sID] == 0 {
					lvl := p.SkillLevels[sID]; if lvl == 0 { lvl = 1 }
					p.Magic -= skill.MPCost
					p.SkillCooldowns[sID] = skill.Cooldown; actionTaken = true
					
					if (p.SystemOrigin == "Spider" || p.SystemOrigin == "God (Shiraori)") && (sID == "heresy_magic" || sID == "rot_attack") { 
						p.GainTaboo(1) 
					}

					// Skill Usage Progression
					p.SkillUsage[sID]++
					if p.SkillUsage[sID] >= 10 {
						p.UpgradeSkill(sID, true)
					}

					switch skill.Category {
					case "attack":
						if (sID == "predator" || sID == "beelzebuth") && float64(monsterHealth)/float64(m.Health) < (map[string]float64{"predator": 0.2, "beelzebuth": 0.4}[sID]) {
							damageToMonster = monsterHealth
							fmt.Printf("🌀 CONSUMED! +2 Atk, +10 HP.\n"); p.Attack += 2; p.MaxHealth += 10
							p.GainTaboo(1)
						} else {
							bonusDmg := skill.DmgBonus + (skill.DmgBonus * (lvl - 1) / 2)
							damageToMonster = int(float64(baseAtk+bonusDmg+rand.Intn(10)) * damageMultiplier)
							fmt.Printf("✨ %s unleashed!\n", skill.Name)
						}
					case "heal":
						healAmt := int(float64(p.MaxHealth) * 0.2) + (lvl * 10)
						p.Health += healAmt; if p.Health > p.MaxHealth { p.Health = p.MaxHealth }
						fmt.Printf("💚 Healed %d HP.\n", healAmt)
					case "defense":
						if sID == "thick_skin" { tempDefense = 50 } else { tempDefense = 30 + (lvl * 5) }
						fmt.Printf("🛡️ Defense UP.\n")
					case "utility":
						if sID == "ruler_of_pride" { rulerPrideTurns = 3; fmt.Println("👑 Ruler of Pride Active!") }
					}
				} else { fmt.Printf("❌ Cooldown: %d\n", p.SkillCooldowns[sID]) }
			}
		}
		if !actionTaken { fmt.Println("❓ Unknown."); continue }
		if damageToMonster > 0 {
			monsterHealth -= damageToMonster
			fmt.Printf("💥 Dealt %d damage. (%d HP left)\n", damageToMonster, monsterHealth)
		}
		if monsterHealth <= 0 {
			fmt.Printf("🏆 Defeated %s!\n", m.Name); p.Kills++; p.MonsterKills[strings.ToLower(m.Name)]++
			if p.SystemOrigin == "Spider" && p.Taboo >= 10 {
				p.AddSkill("heresy_magic")
			}
			p.CheckTitles()
			for item, prob := range m.LootTable {
				if rand.Float64() <= prob { p.Inventory[item]++; fmt.Printf("🎁 Dropped: %s\n", item) }
			}
			p.TrackQuest("combat", m.Name, 1)
			if isGate { p.GainHunterXP(20 + rand.Intn(15)) } else { p.GainXP(15 + rand.Intn(10)) }
			return true
		}
		baseDamage := m.Damage + rand.Intn(5); finalDamage := baseDamage - p.Defense
		if tempDefense > 0 { finalDamage = int(float64(finalDamage) * (1.0 - float64(tempDefense)/100.0)) }
		if finalDamage < 1 { finalDamage = 1 }; p.Health -= finalDamage
		fmt.Printf("👹 %s hits for %d. (%d HP left)\n", m.Name, finalDamage, p.Health)
	}

	if p.Health <= 0 {
		if p.Inventory["life_stone"] > 0 {
			p.Inventory["life_stone"]--; if p.Inventory["life_stone"] == 0 { delete(p.Inventory, "life_stone") }
			p.Health = p.MaxHealth; p.Magic = p.MaxMagic; fmt.Println("\n💎 REVIVED!"); p.Save(); return false
		}
		fmt.Println("\n💀 YOU DIED! -50% Gold, -20% XP, -1 Level."); p.Inventory["gold"] = int(float64(p.Inventory["gold"]) * 0.5)
		p.XP = int(float64(p.XP) * 0.8); p.HunterXP = int(float64(p.HunterXP) * 0.8)
		if p.Level > 1 { p.Level--; p.UpdateRank() }
		if p.HunterLevel > 1 { p.HunterLevel--; p.UpdateHunterRank() }
		p.Health = 50; p.Stamina = 10; p.Magic = 0; p.Save(); return false
	}
	return false
}

func (p *Player) EnterGate(isAdmin bool) {
	if p.CurrentGate == nil { fmt.Println("📭 No gate."); return }
	gate := p.CurrentGate
	if !isAdmin && p.Level < gate.MinLevel { fmt.Printf("🚫 Min Level %d required!\n", gate.MinLevel); return }
	if p.Stamina < 20 { fmt.Println("😫 Low stamina!"); return }
	p.Stamina -= 20
	if isAdmin { p.WorldNotice("ADMINISTRATIVE GATE OVERRIDE ACTIVATED") }
	fmt.Printf("\n🌀 Entering %s-Rank Gate...\n", gate.Rank)
	for floor := 1; floor <= gate.Floors; floor++ {
		fmt.Printf("\n🏢 FLOOR %d / %d\n", floor, gate.Floors)
		if floor == gate.Floors {
			fmt.Printf("\n👹 BOSS: %s!\n", gate.Boss.Name)
			if !p.Combat(&gate.Boss, true) { return }
			fmt.Printf("🎊 CLEARED! +%d Gold, +%d Hunter XP\n", gate.RewardGold, gate.RewardXP)
			p.Inventory["gold"] += gate.RewardGold; p.GainHunterXP(gate.RewardXP); p.CurrentGate = nil
			p.Save()
		} else {
			monsterCount := gate.MonsterCount/gate.Floors + 1
			for i := 0; i < monsterCount; i++ {
				monster := Monster{Name: fmt.Sprintf("%s-Rank Beast", gate.Rank), Health: 20 * gate.MinLevel, Damage: 5 * gate.MinLevel}
				if !p.Combat(&monster, true) { return }
			}
			fmt.Printf("\n✅ Floor %d cleared!\n", floor)
		}
	}
}

func (p *Player) NameSubordinate(species, givenName string) {
	species = strings.ToLower(strings.Replace(species, "_", " ", -1))
	if p.MaxMagic < 50 {
		fmt.Println("❌ Insufficient maximum magic capacity to name a subordinate.")
		return
	}
	
	subAtk := 10; subDef := 2
	valid := false
	
	switch species {
	case "slime": valid = true; subAtk = 15; subDef = 5
	case "goblin", "hobgoblin": valid = true; subAtk = 25; subDef = 10
	case "wolf", "alpha wolf": valid = true; subAtk = 40; subDef = 15
	case "small lesser taratect", "lesser taratect": valid = true; subAtk = 30; subDef = 10
	case "greater taratect", "arch taratect": valid = true; subAtk = 100; subDef = 50
	case "naga", "naga warrior": valid = true; subAtk = 60; subDef = 30
	case "orc", "lesser orc lord": valid = true; subAtk = 50; subDef = 40
	}

	if !valid {
		fmt.Printf("❌ Individual of species '%s' cannot be named at this time.\n", species)
		return
	}
	
	p.MaxMagic -= 50
	if p.Magic > p.MaxMagic { p.Magic = p.MaxMagic }
	
	sub := Subordinate{
		Name:    givenName,
		Species: species,
		Attack:  subAtk,
		Defense: subDef,
		Level:   1,
	}
	p.Subordinates = append(p.Subordinates, sub)
	p.WorldNotice(fmt.Sprintf("Individual '%s' (Species: %s) named. Evolution to NAMED STATUS complete.", givenName, species))
	p.Save()
}

func (p *Player) ListSubordinates() {
	fmt.Println("\n   ╔═══════════════════════╗")
	fmt.Println("   ║ 🤝 *SUBORDINATES*     ║")
	fmt.Println("   ╚═══════════════════════╝")
	if len(p.Subordinates) == 0 {
		fmt.Println("   (None identified)")
	} else {
		for _, s := range p.Subordinates {
			fmt.Printf("   🐾 %s [%s] - ATK: %d | DEF: %d\n", s.Name, strings.ToUpper(s.Species), s.Attack, s.Defense)
		}
	}
	fmt.Println("   ━━━━━━━━━━━━━━━━━━━\n   !name <species> <name>")
}

func (p *Player) ListDCraftable() {
	fmt.Println("\n--- 🛠️ Dungeon Crafting ---")
	for id, r := range Recipes {
		if r.RequiredLevel >= 10 { // Only show high-tier dungeon gear
			fmt.Printf("[%s] %s (Req Lvl %d)\n", id, r.Name, r.RequiredLevel)
		}
	}
}

func (p *Player) ListSkills() {
	fmt.Println("\n   ╔═══════════════════════╗")
	fmt.Println("   ║ 🎮 *DUNGEON SKILLS* ║")
	fmt.Println("   ╚═══════════════════════╝")
	fmt.Printf("\n   🎯 *SP:* %d  |  🎮 *Slots:* %d/%d\n", p.SkillPoints, len(p.EquippedSkills), p.SkillSlots)
	fmt.Println("\n   ━━━ *EQUIPPED* ━━━")
	for i, sID := range p.EquippedSkills {
		skill := GlobalSkills[sID]; lvl := p.SkillLevels[sID]; if lvl == 0 { lvl = 1 }
		fmt.Printf("   [%d] %s Lv%d [%s-Rank]\n       _%s_ (%d MP)\n", i+1, skill.Name, lvl, skill.Rank, skill.Desc, skill.MPCost)
	}
	owned := make(map[string]bool); for _, s := range p.Skills { owned[s] = true }
	fmt.Println("\n   ━━━ *LOCKED* ━━━")
	ranks := []string{"E", "D", "C", "B", "A", "S", "Unique", "Ultimate", "Forbidden"}
	for _, r := range ranks {
		for _, skill := range GlobalSkills {
			if skill.Rank == r && !owned[skill.ID] {
				status := "❌"
				if skill.ReqLevel > 0 && p.Level >= skill.ReqLevel { status = "✅" }
				if skill.ReqBoss != "" && p.MonsterKills[strings.ToLower(skill.ReqBoss)] > 0 { status = "✅" }
				if skill.Rank == "Forbidden" && p.Taboo >= 10 { status = "✅" }
				if strings.Contains(skill.UnlockRequirement, "Origin") && strings.Contains(strings.ToLower(skill.UnlockRequirement), strings.ToLower(p.SystemOrigin)) { status = "✅" }
				fmt.Printf("   %s %s [%s-Rank]\n       Req: %s\n", status, skill.Name, skill.Rank, skill.UnlockRequirement)
			}
		}
	}
	fmt.Println("\n   ━━━━━━━━━━━━━━━━━━━\n   !equip <id> · !unequip <slot#> · !dupskill <id>")
}

func (p *Player) UnequipSkill(slot int) {
	if slot < 1 || slot > len(p.EquippedSkills) { return }
	sID := p.EquippedSkills[slot-1]
	p.EquippedSkills = append(p.EquippedSkills[:slot-1], p.EquippedSkills[slot:]...)
	fmt.Printf("⚪ Unequipped %s.\n", GlobalSkills[sID].Name); p.Save()
}

func (p *Player) AddSkill(skillID string) {
	skillID = strings.ToLower(skillID)
	for _, s := range p.Skills { if s == skillID { return } }
	p.Skills = append(p.Skills, skillID)
	if s, ok := GlobalSkills[skillID]; ok {
		p.WorldNotice(fmt.Sprintf("New Skill Acquired: %s", s.Name))
	}
}

func (p *Player) UpgradeSkill(skillID string, isFree bool) {
	skillID = strings.ToLower(skillID)
	if !isFree && p.SkillPoints < 1 { return }
	if p.SkillLevels[skillID] == 0 { p.SkillLevels[skillID] = 1 }
	p.SkillLevels[skillID]++
	if !isFree { p.SkillPoints-- }
	
	p.SkillUsage[skillID] = 0 // Reset usage on level up
	skillName := GlobalSkills[skillID].Name
	fmt.Printf("✨ Upgraded %s to Lv%d!\n", skillName, p.SkillLevels[skillID])
	p.WorldNotice(fmt.Sprintf("Skill Level Increased: %s is now Lv%d", skillName, p.SkillLevels[skillID]))

	// Evolution Logic
	if p.SkillLevels[skillID] >= 10 {
		if evolvedID, canEvolve := SkillEvolutions[skillID]; canEvolve {
			// Remove old skill
			newSkills := []string{}
			for _, s := range p.Skills { if s != skillID { newSkills = append(newSkills, s) } }
			p.Skills = newSkills

			// Remove from equipped
			newEquipped := []string{}
			for _, s := range p.EquippedSkills { if s != skillID { newEquipped = append(newEquipped, s) } }
			p.EquippedSkills = newEquipped

			// Add new skill
			p.AddSkill(evolvedID)
			p.SkillLevels[evolvedID] = 1
			p.SkillUsage[evolvedID] = 0
			p.WorldNotice(fmt.Sprintf("SKILL ASCENSION: %s has evolved into %s!", skillName, GlobalSkills[evolvedID].Name))
		}
	}
	p.Save()
}

func (p *Player) LearnSkill(skillID string) {
	skillID = strings.ToLower(skillID)
	skill, exists := GlobalSkills[skillID]; if !exists { return }
	for _, s := range p.Skills { if s == skillID { return } }
	met := false
	if skill.ReqBoss != "" && p.MonsterKills[strings.ToLower(skill.ReqBoss)] > 0 { met = true }
	if skill.ReqLevel > 0 && p.Level >= skill.ReqLevel { met = true }
	if skill.Rank == "Forbidden" && p.Taboo >= 10 { met = true }
	
	req := strings.ToLower(skill.UnlockRequirement)
	origin := strings.ToLower(p.SystemOrigin)
	if strings.Contains(req, "origin") && strings.Contains(req, origin) { met = true }
	if strings.Contains(req, "evolution") && strings.Contains(req, origin) { met = true }

	if met {
		p.AddSkill(skillID)
		p.Save()
	} else {
		fmt.Printf("🚫 Requirements not met: %s\n", skill.UnlockRequirement)
	}
}

func (p *Player) EquipSkill(skillID string) {
	skillID = strings.ToLower(skillID)
	owned := false; for _, s := range p.Skills { if s == skillID { owned = true; break } }
	if !owned { return }
	for i, eq := range p.EquippedSkills {
		if eq == skillID {
			p.EquippedSkills = append(p.EquippedSkills[:i], p.EquippedSkills[i+1:]...)
			p.WorldNotice(fmt.Sprintf("Skill Unequipped: %s", GlobalSkills[skillID].Name))
			p.Save(); return
		}
	}
	if len(p.EquippedSkills) >= p.SkillSlots { return }
	p.EquippedSkills = append(p.EquippedSkills, skillID)
	p.WorldNotice(fmt.Sprintf("Skill Equipped: %s", GlobalSkills[skillID].Name))
	p.Save()
}

func (p *Player) CheckTitles() {
	for id, t := range GlobalTitles {
		owned := false; for _, o := range p.Titles { if o == id { owned = true; break } }
		if owned { continue }

		met := false
		if t.KillsNeeded > 0 && p.Kills >= t.KillsNeeded { met = true }
		
		// Special Requirements
		switch id {
		case "taboo_master": if p.Taboo >= 10 { met = true }
		case "slime_emperor": if p.SystemOrigin == "Ultimate Slime (True Dragon)" { met = true }
		case "labyrinth_walker": if p.ExplorationDepth >= 100 { met = true }
		case "world_conqueror": if p.Level >= 100 && len(p.Subordinates) >= 5 { met = true }
		case "supreme_hunter": if p.HunterLevel >= 100 { met = true }
		}

		if met {
			p.Titles = append(p.Titles, id)
			p.Attack += t.AttackBonus; p.MaxHealth += t.HPBonus; p.Health += t.HPBonus
			p.Defense += t.DefenseBonus; p.MaxMagic += t.MPBonus; p.Magic += t.MPBonus
			fmt.Printf("\n🏆 [NEW TITLE]: %s UNLOCKED!\n", t.Name)
			p.WorldNotice(fmt.Sprintf("Title '%s' verified. Permanent stat bonus applied.", t.Name))
			p.Save()
		}
	}
}

func (p *Player) ListTitles() {
	fmt.Println("\n--- 🏅 Titles ---")
	for _, tID := range p.Titles { fmt.Printf("%s - %s\n", GlobalTitles[tID].Name, GlobalTitles[tID].PerkDesc) }
}

func (p *Player) ShowStats() {
	fmt.Printf("\n--- 👤 [SYSTEM] STATUS ID ---\n")
	fmt.Printf("   NAME:      %s\n", p.Name)
	fmt.Printf("   ORIGIN:    %s\n", p.SystemOrigin)
	fmt.Printf("   MINE:      [%s] Lvl %d\n", p.Rank, p.Level)
	fmt.Printf("   HUNTER:    [%s] Lvl %d\n", p.HunterRank, p.HunterLevel)
	fmt.Printf("   VITALS:    HP %d/%d | MP %d/%d\n", p.Health, p.MaxHealth, p.Magic, p.MaxMagic)
	fmt.Printf("   COMBAT:    ATK %d | DEF %d\n", p.Attack, p.Defense)
	fmt.Printf("   RECORDS:   KILLS %d | TABOO %d\n", p.Kills, p.Taboo)
	fmt.Println("   ━━━━━━━━━━━━━━━━━━━")
}

func (p *Player) ShowHelp() {
	fmt.Println("\n   ╔═══════════════════════╗")
	fmt.Println("   ║ 📖 *SYSTEM GUIDE*     ║")
	fmt.Println("   ╚═══════════════════════╝")
	fmt.Println("\n   ━━━ *CORE COMMANDS* ━━━")
	fmt.Println("   !mine <loc>   - Gather resources in zones")
	fmt.Println("   !status / !s  - View character profile")
	fmt.Println("   !inventory / !i- Check items")
	fmt.Println("   !quests       - View missions")
	fmt.Println("\n   ━━━ *ANIME SYSTEMS* ━━━")
	fmt.Println("   !origin <type>- Slime or Spider (Lvl 5)")
	fmt.Println("   !evolve       - Species progression")
	fmt.Println("   !name <sp> <n>- Name subordinate (Uses Max MP)")
	fmt.Println("   !subordinates - View your allies")
	fmt.Println("\n   ━━━ *COMBAT & SKILLS* ━━━")
	fmt.Println("   !enter        - Challenge Gate")
	fmt.Println("   !skills       - View collection and requirements")
	fmt.Println("   !equip <id>   - Set skill to an active slot")
	fmt.Println("   !unequip <#>  - Remove skill from slot")
	fmt.Println("   !learn <id>   - Unlock new ability")
	fmt.Println("\n   ━━━ *SKILL USAGE & MANAGEMENT* ━━━")
	fmt.Println("   🔥 ACTIVE: Use !fight1, !fight2, etc. during combat.")
	fmt.Println("   👁️ PASSIVE: (e.g. Great Sage) Active automatically if EQUIPPED.")
	fmt.Println("   🌀 AUTO: Analysis/Appraisal also works during !mine/!explore.")
	fmt.Println("   🧬 !duplicate <sub_name> <skill_id> - Copy sub's skill (Req: Shub-Niggurath)")
	fmt.Println("   ✨ !create <skill_1> <skill_2>      - Birth new skill (Req: Shub-Niggurath)")
}

func (p *Player) HealFull() { p.Health = p.MaxHealth; p.Stamina = p.MaxStamina; p.Magic = p.MaxMagic; p.Save() }

func (p *Player) Save() {
	data, _ := json.MarshalIndent(p, "", "  ")
	os.WriteFile("player_data.json.bak", data, 0644)
	os.WriteFile("player_data.json", data, 0644)
}

func LoadPlayer() *Player {
	data, err := os.ReadFile("player_data.json")
	if err != nil {
		data, err = os.ReadFile("player_data.json.bak")
		if err != nil { return NewPlayer("Adventurer") }
	}
	var p Player; json.Unmarshal(data, &p)
	if p.Inventory == nil { p.Inventory = make(map[string]int) }
	if p.QuestProgress == nil { p.QuestProgress = make(map[string]int) }
	if p.MonsterKills == nil { p.MonsterKills = make(map[string]int) }
	if p.SkillCooldowns == nil { p.SkillCooldowns = make(map[string]int) }
	if p.SkillLevels == nil { p.SkillLevels = make(map[string]int) }
	if p.SkillUsage == nil { p.SkillUsage = make(map[string]int) }
	if p.Subordinates == nil { p.Subordinates = []Subordinate{} }
	if p.Structures == nil { p.Structures = make(map[string]bool) }
	if p.SystemOrigin == "" { p.SystemOrigin = "Human" }
	if p.MaxMagic == 0 { p.MaxMagic = 100; p.Magic = 100 }
	p.UpdateRank(); p.UpdateHunterRank(); p.UpdateSkillSlots()
	return &p
}

func (p *Player) TrackQuest(qType, id string, qty int) {
	for _, q := range GlobalQuests {
		if q.TargetType == qType && q.TargetID == id {
			if p.QuestProgress[q.ID] < q.TargetQty {
				if qType == "level" {
					if qty > p.QuestProgress[q.ID] { p.QuestProgress[q.ID] = qty }
				} else {
					p.QuestProgress[q.ID] += qty
				}
				if p.QuestProgress[q.ID] >= q.TargetQty {
					p.WorldNotice(fmt.Sprintf("Mission '%s' concluded.", q.Name)); p.Inventory["gold"] += q.RewardGold; p.GainXP(q.RewardXP)
				}
			}
		}
	}
}

func (p *Player) ListQuests() {
	for _, q := range GlobalQuests {
		progress := p.QuestProgress[q.ID]
		if progress >= q.TargetQty {
			fmt.Printf("✅ [%s] %s - COMPLETED\n", q.ID, q.Name)
		} else {
			fmt.Printf("🔄 [%s] %s - %d/%d\n", q.ID, q.Name, progress, q.TargetQty)
		}
	}
}

func (p *Player) ShowInventory() {
	fmt.Printf("\n--- 🎒 Inventory ---\n")
	for id, qty := range p.Inventory { if qty > 0 { fmt.Printf("%s: %d\n", id, qty) } }
}

func (p *Player) StartRegeneration() {
	ticker := time.NewTicker(20 * time.Minute)
	go func() { for range ticker.C { p.Regenerate() } }()
}

func (p *Player) LogAction(msg string) {
	timestamp := time.Now().Format("15:04:05")
	p.ActionLog = append([]string{fmt.Sprintf("[%s] %s", timestamp, msg)}, p.ActionLog...)
	if len(p.ActionLog) > 30 { p.ActionLog = p.ActionLog[:30] }
}

func (p *Player) StartSubordinateAutonomy() {
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for range ticker.C {
			for i := range p.Subordinates {
				p.SubordinateAction(&p.Subordinates[i])
			}
			p.Save()
		}
	}()
}

func (p *Player) SubordinateAction(s *Subordinate) {
	if time.Since(s.LastAction) < 5*time.Minute { return }
	s.LastAction = time.Now()

	action := rand.Intn(100)
	if action < 40 { // Mining
		locs := []string{"surface", "cave", "abyss", "nether", "void"}
		locID := locs[rand.Intn(len(locs))]
		loc := Locations[locID]
		if s.Level >= loc.RequiredLevel {
			p.LogAction(fmt.Sprintf("%s is mining in %s...", s.Name, loc.Name))
			found := ""
			for item, prob := range loc.LootTable {
				if rand.Float64() <= prob {
					p.Inventory[item]++
					found += strings.Replace(item, "_", " ", -1) + ", "
				}
			}
			if found != "" { p.LogAction(fmt.Sprintf("%s found: %s", s.Name, strings.TrimSuffix(found, ", "))) }
			p.GainXP(5)
			p.SubordinateGainXPForOne(s, 20)
		}
	} else if action < 70 { // Raiding
		raids := []string{"goblin_camp", "bandit_fort", "shadow_keep"}
		raidID := raids[rand.Intn(len(raids))]
		raid := BotSettlements[raidID]
		if s.Level >= raid.Level {
			p.LogAction(fmt.Sprintf("%s is raiding %s!", s.Name, raid.Name))
			for item, qty := range raid.LootTable {
				p.Inventory[item] += qty
			}
			p.LogAction(fmt.Sprintf("%s completed raid on %s.", s.Name, raid.Name))
			p.GainXP(raid.Level * 10)
			p.SubordinateGainXPForOne(s, raid.Level*50)
		}
	}
}

func (p *Player) SubordinateGainXP(amount int) {
	for i := range p.Subordinates {
		p.SubordinateGainXPForOne(&p.Subordinates[i], amount)
	}
}

func (p *Player) SubordinateGainXPForOne(s *Subordinate, amount int) {
	if s.NextXP == 0 { s.NextXP = 100 }
	s.XP += amount
	if s.XP >= s.NextXP {
		s.Level++
		s.XP -= s.NextXP
		s.NextXP = int(float64(s.NextXP) * 1.5)
		s.Attack += 10
		s.Defense += 10
		p.WorldNotice(fmt.Sprintf("Subordinate '%s' has reached Level %d", s.Name, s.Level))
		p.CheckSubordinateEvolution(s)
		p.CheckSubordinateSkills(s)
	}
}

func (p *Player) CheckSubordinateSkills(s *Subordinate) {
	speciesSkills := map[string][]struct{lvl int; id string}{
		"slime": {{1, "predator"}, {5, "water_blade"}, {15, "gluttony"}, {30, "black_lightning"}},
		"spider": {{1, "appraisal"}, {5, "poison_fang"}, {15, "evil_eye"}, {30, "heresy_magic"}},
		"alpha wolf": {{1, "power_strike"}, {10, "venom_coat"}},
	}
	
	for _, ss := range speciesSkills[s.Species] {
		if s.Level >= ss.lvl {
			owned := false
			for _, sk := range s.Skills { if sk == ss.id { owned = true; break } }
			if !owned {
				s.Skills = append(s.Skills, ss.id)
				p.WorldNotice(fmt.Sprintf("Subordinate '%s' learned: %s", s.Name, GlobalSkills[ss.id].Name))
			}
		}
	}
}

func (p *Player) CheckSubordinateEvolution(s *Subordinate) {
	oldSpecies := s.Species
	evolved := false
	if s.Species == "hobgoblin" && s.Level >= 10 {
		s.Species = "ogre"; s.Attack += 20; s.Defense += 10; evolved = true
	} else if s.Species == "ogre" && s.Level >= 25 {
		s.Species = "kijin"; s.Attack += 50; s.Defense += 30; evolved = true
	} else if s.Species == "alpha wolf" && s.Level >= 15 {
		s.Species = "tempest wolf"; s.Attack += 30; s.Defense += 15; evolved = true
	}

	if evolved {
		p.WorldNotice(fmt.Sprintf("Subordinate '%s' has evolved from %s to %s!", s.Name, oldSpecies, s.Species))
	}
}

func (p *Player) AutoAnalyze(itemID string) {
	lvl := 1
	if l, ok := p.SkillLevels["great_sage"]; ok { lvl = l }
	if l, ok := p.SkillLevels["raphael"]; ok { lvl = l + 10 }
	if l, ok := p.SkillLevels["wisdom"]; ok { lvl = l + 10 }

	rarity := "Common"; price := 5
	switch itemID {
	case "diamond", "abyss_crystal", "netherite": rarity = "Epic"; price = 100
	case "void_essence", "star_matter", "void_core": rarity = "Legendary"; price = 500
	case "void_crown", "life_stone", "demon_soul": rarity = "Mythic"; price = 5000
	case "iron", "gold", "quartz": rarity = "Rare"; price = 25
	}

	fmt.Printf("\n[🧠 ANALYSIS]: '%s' identified.", strings.ToUpper(itemID))
	if lvl >= 4 { fmt.Printf(" | Rarity: %s", rarity) }
	if lvl >= 7 { fmt.Printf(" | Market Value: 💰 %d Gold", price) }
	if lvl >= 10 {
		desc := "Standard material."
		if r, ok := Recipes[itemID]; ok { desc = fmt.Sprintf("Used to craft %s.", r.Name) }
		fmt.Printf("\n   📜 %s", desc)
	}
	fmt.Println()
}

func (p *Player) DuplicateSkill(subName, skillID string) {
	if !p.HasSkill("shub_niggurath") {
		fmt.Println("❌ Harvest Lord Shub-Niggurath required.")
		return
	}
	skillID = strings.ToLower(skillID)
	if p.HasSkill(skillID) {
		fmt.Println("❌ Skill already owned.")
		return
	}

	var targetSub *Subordinate
	for i := range p.Subordinates {
		if strings.EqualFold(p.Subordinates[i].Name, subName) {
			targetSub = &p.Subordinates[i]
			break
		}
	}

	if targetSub == nil {
		fmt.Printf("❌ Subordinate '%s' not found.\n", subName)
		return
	}

	found := false
	for _, s := range targetSub.Skills {
		if s == skillID {
			found = true
			break
		}
	}

	if !found {
		fmt.Printf("❌ Subordinate does not possess skill '%s'.\n", skillID)
		return
	}

	p.AddSkill(skillID)
	p.WorldNotice(fmt.Sprintf("SKILL DUPLICATED: '%s' has been harvested from %s.", GlobalSkills[skillID].Name, targetSub.Name))
	p.Save()
}

func (p *Player) CreateSkill(skillA, skillB string) {
	if !p.HasSkill("shub_niggurath") {
		fmt.Println("❌ Harvest Lord Shub-Niggurath required.")
		return
	}
	if !p.HasSkill(skillA) || !p.HasSkill(skillB) {
		fmt.Println("❌ Source skills not owned.")
		return
	}
	if p.Magic < 200 {
		fmt.Println("❌ Insufficient Magic (200 MP needed).")
		return
	}

	p.Magic -= 200
	// Pick a random skill the player doesn't have from high tiers
	possible := []string{}
	for id, s := range GlobalSkills {
		if (s.Rank == "A" || s.Rank == "S" || s.Rank == "Ultimate") && !p.HasSkill(id) {
			possible = append(possible, id)
		}
	}

	if len(possible) == 0 {
		fmt.Println("❌ No new insights attainable at this time.")
		return
	}

	newSkill := possible[rand.Intn(len(possible))]
	p.AddSkill(newSkill)
	p.WorldNotice(fmt.Sprintf("SKILL CREATION: Shub-Niggurath has birthed a new ability: %s!", GlobalSkills[newSkill].Name))
	p.Save()
}

func (p *Player) HasSkill(id string) bool {
	id = strings.ToLower(id)
	for _, s := range p.Skills { if s == id { return true } }
	return false
}

func (p *Player) Regenerate() {
	h := 10; s := 10; m := 20
	if p.Structures["house"] { h += 2 }; if p.Structures["farm"] { s += 5 }
	p.Health += h; if p.Health > p.MaxHealth { p.Health = p.MaxHealth }
	p.Stamina += s; if p.Stamina > p.MaxStamina { p.Stamina = p.MaxStamina }
	p.Magic += m; if p.Magic > p.MaxMagic { p.Magic = p.MaxMagic }
	if p.Health < p.MaxHealth/4 { p.WorldNotice("Warning: Vitality low.") }
	if p.Magic < p.MaxMagic/5 { p.WorldNotice("Warning: Magicule levels low.") }
	p.Save()
}

func (p *Player) StartRaids() {
	ticker := time.NewTicker(30 * time.Minute)
	go func() { for range ticker.C { if rand.Float64() < 0.3 { p.UnderRaid() } } }()
}

func (p *Player) UnderRaid() {
	fmt.Printf("\n🚨 ALERT: HOSTILE FORCES DETECTED! 🚨\n")
	raidStrength := (p.Level / 5) + 1
	raider := Monster{Name: "🏴‍☠️ Raider Battalion", Health: 100 * raidStrength, Damage: 20 * raidStrength}
	totalSubDefense := 0
	for _, s := range p.Subordinates {
		totalSubDefense += s.Defense
		fmt.Printf("🤝 Subordinate '%s' is assisting in defense!\n", s.Name)
	}
	p.Defense += totalSubDefense
	success := p.Combat(&raider, false)
	p.Defense -= totalSubDefense
	if success { fmt.Println("🛡️ RAID REPELLED.") } else { fmt.Println("📉 ASSETS PLUNDERED.") }
	p.Save()
}

func (p *Player) ListRaids() {
	for id, s := range BotSettlements { fmt.Printf("[%s] %s (Lvl %d)\n", id, s.Name, s.Level) }
}

func (p *Player) Raid(targetID string) {
	t, ok := BotSettlements[strings.ToLower(targetID)]; if !ok { return }
	if p.Inventory[t.RequiredSword] <= 0 || p.Level < t.Level || p.Stamina < 30 { return }
	p.Stamina -= 30
	p.WorldNotice(fmt.Sprintf("Commencing Raid on %s", t.Name))
	for _, d := range t.Defenders { if !p.Combat(&d, false) { return } }
	fmt.Println("\n💰 [LOOT ACQUIRED]:")
	for id, qty := range t.LootTable {
		p.Inventory[id] += qty
		fmt.Printf("   - %s x%d\n", strings.Replace(id, "_", " ", -1), qty)
	}
	p.GainTaboo(1)
	p.GainXP(100 + t.Level*10); p.Save()
	p.WorldNotice(fmt.Sprintf("Raid on %s successful.", t.Name))
}

func (p *Player) ListShop() {
	fmt.Printf("Gold: 💰 %d\n", p.Inventory["gold"])
	for id, it := range MerchantInventory { fmt.Printf("[%s] %s - 💰 %d\n", id, it.Name, it.Price) }
}

func (p *Player) Buy(itemID string) {
	it, ok := MerchantInventory[strings.ToLower(itemID)]; if !ok || p.Inventory["gold"] < it.Price { return }
	p.Inventory["gold"] -= it.Price
	p.WorldNotice(fmt.Sprintf("Purchased %s for 💰 %d gold.", it.Name, it.Price))
	switch it.ID {
	case "golden_apple": p.Health += 100; if p.Health > p.MaxHealth { p.Health = p.MaxHealth }
	case "energy_drink": p.Stamina += 50; if p.Stamina > p.MaxStamina { p.Stamina = p.MaxStamina }
	case "repair_kit": p.ToolDurability = 500
	case "mystery_box":
		loot := []string{"iron", "gold", "diamond", "quartz", "netherite"}
		for i := 0; i < 3; i++ { p.Inventory[loot[rand.Intn(len(loot))]] += 5 + rand.Intn(10) }
	default: p.Inventory[it.ID]++
	}
	p.Save()
}

func (p *Player) ListCraftable() {
	for id, r := range Recipes { fmt.Printf("[%s] Lvl %d\n", id, r.RequiredLevel) }
}

func (p *Player) Craft(itemName string) {
	r, ok := Recipes[strings.ToLower(itemName)]; if !ok || p.Level < r.RequiredLevel { return }
	for id, qty := range r.Ingredients { if p.Inventory[id] < qty { return } }
	for id, qty := range r.Ingredients { p.Inventory[id] -= qty; if p.Inventory[id] == 0 { delete(p.Inventory, id) } }
	switch r.ResultType {
	case "tool": p.ToolDurability = r.ResultValue; p.Inventory[strings.ToLower(itemName)] = 1
	case "weapon", "food", "stamina_food": p.Inventory[strings.ToLower(itemName)]++
	case "armor": p.Defense += r.ResultValue; p.Inventory[strings.ToLower(itemName)] = 1
	}
	p.WorldNotice(fmt.Sprintf("Successfully crafted: %s", r.Name))
	p.GainXP(10 + rand.Intn(5)); p.Save()
}

func (p *Player) Use(itemName string) {
	k := strings.ToLower(itemName); r, ok := Recipes[k]; if !ok || p.Inventory[k] <= 0 { return }
	p.Inventory[k]--; if p.Inventory[k] == 0 { delete(p.Inventory, k) }
	if r.ResultType == "food" { p.Health += r.ResultValue; if p.Health > p.MaxHealth { p.Health = p.MaxHealth } }
	if r.ResultType == "stamina_food" { p.Stamina += r.ResultValue; if p.Stamina > p.MaxStamina { p.Stamina = p.MaxStamina } }
	p.WorldNotice(fmt.Sprintf("Used item: %s", r.Name))
	p.Save()
}

func (p *Player) ListBuildable() {
	for id, s := range Structures { fmt.Printf("[%s] Lvl %d - %s\n", id, s.RequiredLevel, s.PerkDesc) }
}

func (p *Player) Build(structName string) {
	s, ok := Structures[strings.ToLower(structName)]; if !ok || p.Structures[strings.ToLower(structName)] || p.Level < s.RequiredLevel { return }
	for id, qty := range s.Ingredients { if p.Inventory[id] < qty { return } }
	for id, qty := range s.Ingredients { p.Inventory[id] -= qty; if p.Inventory[id] == 0 { delete(p.Inventory, id) } }
	p.Structures[strings.ToLower(structName)] = true
	p.WorldNotice(fmt.Sprintf("Construction complete: %s", s.Name))
	switch strings.ToLower(structName) {
	case "forge": p.Attack += 10
	case "vault": p.MaxHealth += 50; p.Health += 50
	case "castle": p.Attack += 20; p.MaxHealth += 100; p.Health += 100; p.MaxStamina += 50; p.Stamina += 50
	}
	p.GainXP(50 + rand.Intn(50)); p.Save()
}

func (p *Player) StartExploration() {
	if p.Level < 10 { fmt.Println("🚫 Lvl 10 Required."); return }
	p.Exploring = true; p.ExplorationDepth = 1
	p.WorldNotice("ENTERING THE GREAT ELROE LABYRINTH")
	fmt.Println("📍 POSITION: [Entrance]\nACTIONS: [W] Forward | [A] Left | [D] Right | [S] Backward | !emerge")
}

func (p *Player) Emerge() {
	if !p.Exploring {
		fmt.Println("🕵️ You are not presently exploring in any known labyrinth.")
		return
	}
	p.Exploring = false
	p.WorldNotice("EMERGED")
}

func (p *Player) Move(dir string) {
	if !p.Exploring { return }
	if p.Stamina < 2 { fmt.Println("😫 Exhausted!"); return }
	p.Stamina -= 2
	p.ExplorationDepth++
	dir = strings.ToUpper(dir)
	fmt.Printf("\n👣 You move %s... (Depth: %d)\n", map[string]string{"W":"Forward","A":"Left","D":"Right","S":"Backward"}[dir], p.ExplorationDepth)
	hasAnalysis := false
	for _, sID := range p.EquippedSkills { if sID == "great_sage" || sID == "wisdom" { hasAnalysis = true; break } }
	event := rand.Intn(100)
	if event < 15 { p.FoundChest(hasAnalysis) } else if event < 30 { p.TriggerTrap(hasAnalysis) } else if event < 55 { p.EncounterMonster() } else { fmt.Println("🌫️ Path clear.") }
	p.Save()
}

func (p *Player) FoundChest(hasAnalysis bool) {
	fmt.Println("🎁 [HIDDEN CHEST]!")
	loot := []string{"diamond", "void_essence", "life_stone", "health_potion", "star_matter"}
	item := loot[rand.Intn(len(loot))]; qty := 1 + rand.Intn(3)
	p.Inventory[item] += qty; fmt.Printf("   - Obtained: %s x%d\n", item, qty)
	if hasAnalysis { p.AutoAnalyze(item) }
}

func (p *Player) TriggerTrap(hasAnalysis bool) {
	if hasAnalysis {
		for _, sID := range p.EquippedSkills { if sID == "trap_sense" { p.SkillUsage[sID]++; if p.SkillUsage[sID] >= 10 { p.UpgradeSkill(sID, true) } } }
		if rand.Float64() < 0.7 { fmt.Println("⚠️ [System]: TRAP DODGED."); return }
	}
	fmt.Println("⚠️ [TRAP] Pressure plate!"); dmg := 10 + rand.Intn(15); p.Health -= dmg
	fmt.Printf("💥 Arrows hit! (❤️ %d HP left)\n", p.Health)
}

func (p *Player) EncounterMonster() {
	monster := Monster{Name: "🕸️ Labyrinth Stalker", Health: 100 + (p.ExplorationDepth * 20), Damage: 15 + (p.ExplorationDepth * 5), LootTable: map[string]float64{"string": 1.0, "quartz": 0.2}}
	if p.Combat(&monster, false) { p.GainXP(20); p.SubordinateGainXP(10) }
}

func (p *Player) GetBestSwordDamage() int {
	b := 0; for id, q := range p.Inventory { if q > 0 { if r, ok := Recipes[id]; ok && r.ResultType == "weapon" { if r.ResultValue > b { b = r.ResultValue } } } }
	return b
}

func (p *Player) GetBestPickaxeMultiplier() float64 {
	m := 1.0; multis := map[string]float64{"wood_pickaxe":1, "stone_pickaxe":1.2, "iron_pickaxe":1.5, "diamond_pickaxe":2, "abyss_pickaxe":3, "nether_pickaxe":5, "void_pickaxe":10}
	for id, q := range p.Inventory { if q > 0 { if val, ok := multis[id]; ok && val > m { m = val } } }
	return m
}

func (p *Player) Mine(locName string) {
	locID := strings.ToLower(locName)
	loc, ok := Locations[locID]
	
	// Exhaustive Validation
	if !ok {
		fmt.Println("❌ Unknown location. Valid: surface, cave, abyss, nether, void")
		return
	}
	if p.Level < loc.RequiredLevel {
		fmt.Printf("🚫 Insufficient Mine Level! Req: %d (You: %d)\n", loc.RequiredLevel, p.Level)
		return
	}
	if loc.RequiredItem != "" && p.Inventory[loc.RequiredItem] <= 0 {
		fmt.Printf("⛏️ You need a %s to mine in the %s!\n", strings.Replace(loc.RequiredItem, "_", " ", -1), loc.Name)
		return
	}
	if p.Stamina < 10 {
		fmt.Println("😫 You are too exhausted to mine! (Stamina < 10)")
		return
	}
	if p.ToolDurability <= 0 {
		fmt.Println("❌ Your tool is broken! Craft a new one or use a Repair Kit.")
		return
	}

	// Start Mining Action
	p.Stamina -= 10
	p.ToolDurability -= 1

	// Passive Usage: Miner's Instinct
	for _, sID := range p.EquippedSkills { if sID == "miners_instinct" { p.SkillUsage[sID]++; if p.SkillUsage[sID] >= 10 { p.UpgradeSkill(sID, true) } } }
	
	if rand.Float64() <= loc.EncounterChance {
		if !p.Combat(&loc.EncounterTable[rand.Intn(len(loc.EncounterTable))], false) {
			return // Combat failed or player died/escaped
		}
	}
	
	pick := p.GetBestPickaxeMultiplier()
	drops := int(float64(1+p.Level/5) * pick)
	if drops < 1 { drops = 1 }
	
	foundItems := make(map[string]int)
	for i := 0; i < drops; i++ {
		r := rand.Float64()
		var cum float64
		itemFound := false
		for item, prob := range loc.LootTable {
			cum += prob
			if r <= cum {
				p.Inventory[item]++
				foundItems[item]++
				p.TrackQuest("item", item, 1)
				
				// Auto Analysis
				hasAnalysis := false
				for _, sID := range p.EquippedSkills { 
					if sID == "great_sage" || sID == "wisdom" { 
						hasAnalysis = true
						break 
					} 
				}
				if hasAnalysis { p.AutoAnalyze(item) }
				itemFound = true
				break
			}
		}
		// Fallback for rare cases where r is exactly 1.0 (though rand.Float64 is [0,1))
		if !itemFound {
			// Pick first item in map as fallback if nothing found
			for item := range loc.LootTable {
				p.Inventory[item]++
				foundItems[item]++
				break
			}
		}
	}

	if len(foundItems) > 0 {
		for id, qty := range foundItems {
			fmt.Printf("🎁 [GATHERED] %s x%d\n", strings.Replace(id, "_", " ", -1), qty)
		}
	} else {
		fmt.Printf("💨 You mined in the %s but found nothing.\n", loc.Name)
	}

	p.GainXP(2 + rand.Intn(3))
	p.Save()
}
