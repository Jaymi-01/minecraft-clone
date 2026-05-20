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
		SkillPoints:    6,
		Titles:         []string{},
		Skills:         []string{},
		EquippedSkills: []string{},
		SkillSlots:     3,
		SkillLevels:    make(map[string]int),
		SkillCooldowns: make(map[string]int),
	}
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

func (p *Player) GainXP(amount int) {
	if p.Structures["enchanting_table"] {
		amount = int(float64(amount) * 1.5) // +50% XP
	}
	p.XP += amount
	fmt.Printf("[✨ +%d Mine XP]\n", amount)
	for p.XP >= p.XPToNext {
		p.Level++
		p.XP -= p.XPToNext
		p.XPToNext = int(float64(p.XPToNext) * 1.5)
		p.MaxHealth += 10
		p.MaxStamina += 10
		p.Health = p.MaxHealth
		p.Stamina = p.MaxStamina
		
		p.UpdateRank()
		fmt.Printf("\n🎊 MINE LEVEL UP! You are now mine level %d! (Rank: %s) 🎊\n", p.Level, p.Rank)
	}
	p.Save()
}

func (p *Player) GainHunterXP(amount int) {
	p.HunterXP += amount
	fmt.Printf("[🏹 +%d Hunter XP]\n", amount)
	for p.HunterXP >= p.HunterXPToNext {
		p.HunterLevel++
		p.HunterXP -= p.HunterXPToNext
		p.HunterXPToNext = int(float64(p.HunterXPToNext) * 1.5)
		
		// Skill slot increase every 20 hunter levels
		if p.HunterLevel%20 == 0 {
			p.SkillSlots++
			fmt.Printf("🔓 Skill Slot Increased! (Total: %d)\n", p.SkillSlots)
		}

		p.UpdateHunterRank()
		fmt.Printf("\n⚔️ HUNTER LEVEL UP! You are now hunter level %d! (Hunter Rank: %s) ⚔️\n", p.HunterLevel, p.HunterRank)
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
	// Spawn one immediately for testing/first run
	p.SpawnGate()
}

func (p *Player) SpawnGate() {
	ranks := []string{"E", "D", "C", "B", "A", "S", "SS"}
	// Completely random spawn, independent of levels
	rank := ranks[rand.Intn(len(ranks))]
	gate := Gates[rank]
	p.CurrentGate = &gate
	fmt.Printf("\n🌀 [SYSTEM] A %s-Rank Gate has spawned! (Req Mine Lvl: %d) Type !enter to challenge it.\n", rank, gate.MinLevel)
}

func (p *Player) Combat(m *Monster, isGate bool) bool {
	fmt.Printf("\n⚔️ ENCOUNTER: %s (❤️ %d HP, 💥 %d DMG)\n", m.Name, m.Health, m.Damage)
	monsterHealth := m.Health
	reader := bufio.NewReader(os.Stdin)
	
	p.SkillCooldowns = make(map[string]int)
	tempDefense := 0

	for monsterHealth > 0 && p.Health > 0 {
		tempDefense = 0 // Reset every turn
		fmt.Printf("\n--- Your Turn (❤️ %d/%d) ---\n", p.Health, p.MaxHealth)
		
		for sID, cd := range p.SkillCooldowns {
			if cd > 0 {
				p.SkillCooldowns[sID]--
			}
		}

		// Apply Passives
		bonusAtkFromPassives := 0
		critChance := 0.05
		for _, sID := range p.EquippedSkills {
			skill := GlobalSkills[sID]
			if skill.Type == "passive" {
				if sID == "critical_eye" {
					critChance += 0.15
				}
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
			status := "READY"
			if cd > 0 {
				status = fmt.Sprintf("%d turns", cd)
			}
			fmt.Printf("[!fight%d] %s (%s) ", i+1, skill.Name, status)
		}
		fmt.Print("\nChoice: ")
		
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		
		if input == "!recover" {
			p.Health = p.MaxHealth
			p.Stamina = p.MaxStamina
			fmt.Println("⚡ [CHEAT] Health and Stamina fully replenished! ⚡")
			continue
		}

		damageToMonster := 0
		actionTaken := false

		baseAtk := p.Attack + p.GetBestSwordDamage() + bonusAtkFromPassives
		if isGate {
			baseAtk += p.HunterLevel * 2
		} else {
			baseAtk += p.Level
		}

		if input == "!fight" {
			damageToMonster = baseAtk + rand.Intn(5)
			if rand.Float64() <= critChance {
				damageToMonster = int(float64(damageToMonster) * 1.5)
				fmt.Println("🎯 CRITICAL HIT!")
			}
			fmt.Println("🤜 You perform a basic attack.")
			actionTaken = true
		} else if strings.HasPrefix(input, "!fight") {
			skillIdxStr := strings.TrimPrefix(input, "!fight")
			var idx int
			fmt.Sscanf(skillIdxStr, "%d", &idx)
			idx-- // 0-based

			if idx >= 0 && idx < len(p.EquippedSkills) {
				sID := p.EquippedSkills[idx]
				skill := GlobalSkills[sID]
				if skill.Type == "passive" {
					fmt.Printf("ℹ️ %s is a passive skill and is always active!\n", skill.Name)
					continue
				}
				if p.SkillCooldowns[sID] == 0 {
					lvl := p.SkillLevels[sID]
					if lvl == 0 { lvl = 1 }
					
					p.SkillCooldowns[sID] = skill.Cooldown
					actionTaken = true

					switch skill.Category {
					case "attack":
						bonusDmg := skill.DmgBonus + (skill.DmgBonus * (lvl - 1) / 2)
						damageToMonster = baseAtk + bonusDmg + rand.Intn(10)
						fmt.Printf("✨ You unleashed %s (Lv%d)! ✨\n", skill.Name, lvl)
					case "heal":
						healAmt := int(float64(p.MaxHealth) * 0.2) + (lvl * 10)
						p.Health += healAmt
						if p.Health > p.MaxHealth { p.Health = p.MaxHealth }
						fmt.Printf("💚 You used %s! Restored %d HP. (❤️ %d/%d)\n", skill.Name, healAmt, p.Health, p.MaxHealth)
					case "defense":
						if sID == "thick_skin" {
							tempDefense = 50 // 50% reduction
						} else if sID == "bone_armor" || sID == "fortify" {
							tempDefense = 30 + (lvl * 5)
						}
						fmt.Printf("🛡️ You used %s! Defense increased for this turn.\n", skill.Name)
					}
				} else {
					fmt.Printf("❌ %s is still on cooldown (%d turns left).\n", GlobalSkills[sID].Name, p.SkillCooldowns[sID])
				}
			}
		}

		if !actionTaken {
			fmt.Println("❓ Unknown command. Use !fight or !fight1-n.")
			continue
		}

		if damageToMonster > 0 {
			monsterHealth -= damageToMonster
			fmt.Printf("💥 You dealt %d damage to %s. (%d HP left)\n", damageToMonster, m.Name, monsterHealth)
		}
		
		if monsterHealth <= 0 {
			fmt.Printf("🏆 You defeated the %s!\n", m.Name)
			p.Kills++
			p.CheckTitles()
			for item, prob := range m.LootTable {
				if rand.Float64() <= prob {
					p.Inventory[item]++
					fmt.Printf("🎁 Dropped: %s\n", item)
				}
			}
			p.TrackQuest("combat", m.Name, 1)
			if isGate {
				p.GainHunterXP(20 + rand.Intn(15))
			} else {
				p.GainXP(15 + rand.Intn(10))
			}
			return true
		}

		baseDamage := m.Damage + rand.Intn(5)
		finalDamage := baseDamage - p.Defense
		if tempDefense > 0 {
			finalDamage = int(float64(finalDamage) * (1.0 - float64(tempDefense)/100.0))
		}
		if finalDamage < 1 {
			finalDamage = 1
		}
		p.Health -= finalDamage
		fmt.Printf("👹 %s hits you for %d damage (Blocked %d). (%d HP left)\n", m.Name, finalDamage, p.Defense+tempDefense, p.Health)
	}
	
	if p.Health <= 0 {
		if p.Inventory["life_stone"] > 0 {
			p.Inventory["life_stone"]--
			if p.Inventory["life_stone"] == 0 {
				delete(p.Inventory, "life_stone")
			}
			p.Health = p.MaxHealth
			fmt.Println("\n💎 [LIFE STONE] A divine light surrounds you! You have been revived and your penalties prevented.")
			p.Save()
			return false // Saved from penalty, but still "lost" the fight if it was a gate/raid
		}

		fmt.Println("\n💀 YOU DIED!")
		fmt.Println("⚠️ STRICT PENALTY APPLIED: -50% Gold, -20% XP, and -1 Level.")
		p.Inventory["gold"] = int(float64(p.Inventory["gold"]) * 0.5)
		p.XP = int(float64(p.XP) * 0.8)
		p.HunterXP = int(float64(p.HunterXP) * 0.8)
		if p.Level > 1 {
			p.Level--
			p.UpdateRank()
		}
		if p.HunterLevel > 1 {
			p.HunterLevel--
			p.UpdateHunterRank()
		}
		p.Health = 50
		p.Stamina = 10
		p.Save()
		return false
	}
	return false
}

func (p *Player) EnterGate() {
	if p.CurrentGate == nil {
		fmt.Println("📭 No active gate spawned. Wait for the system to detect a rift.")
		return
	}
	gate := p.CurrentGate

	if p.Level < gate.MinLevel {
		fmt.Printf("🚫 Your MINE LEVEL (%d) is too low for this %s-Rank Gate! Required: %d\n", p.Level, gate.Rank, gate.MinLevel)
		return
	}

	if p.Stamina < 20 {
		fmt.Println("😫 You need at least 20 stamina to enter a gate!")
		return
	}
	p.Stamina -= 20

	fmt.Printf("\n🌀 Entering %s-Rank Gate...\n", gate.Rank)
	fmt.Printf("📊 Hunter Status: Level %d (%s-Rank)\n", p.HunterLevel, p.HunterRank)
	
	if len(gate.Descriptions) > 0 {
		fmt.Printf("✨ %s\n", gate.Descriptions[rand.Intn(len(gate.Descriptions))])
	}

	for floor := 1; floor <= gate.Floors; floor++ {
		fmt.Printf("\n🏢 FLOOR %d / %d\n", floor, gate.Floors)
		
		if floor == gate.Floors {
			fmt.Printf("\n👹 BOSS FLOOR REACHED! 👹\n")
			fmt.Printf("The air crackles with immense power. %s awaits you...\n", gate.Boss.Name)
			if !p.Combat(&gate.Boss, true) {
				fmt.Printf("❌ You were defeated by the Boss at the final floor...\n")
				return
			}
			fmt.Printf("🎊 GATE CLEARED! 🎊\n")
			fmt.Printf("💰 Reward: %d Gold, 🏹 %d Hunter XP\n", gate.RewardGold, gate.RewardXP)
			p.Inventory["gold"] += gate.RewardGold
			p.GainHunterXP(gate.RewardXP)

			p.CurrentGate = nil

			// Random skill drop from GlobalSkills
			var eligibleSkills []string
			for id, skill := range GlobalSkills {
				hasSkill := false
				for _, s := range p.Skills {
					if s == id {
						hasSkill = true
						break
					}
				}
				// Boss drops are only E and D rank for now to encourage learning/others
				if !hasSkill && (skill.Rank == "E" || skill.Rank == "D") {
					eligibleSkills = append(eligibleSkills, id)
				}
			}

			if len(eligibleSkills) > 0 {
				newSkill := eligibleSkills[rand.Intn(len(eligibleSkills))]
				p.Skills = append(p.Skills, newSkill)
				fmt.Printf("🔮 You obtained a new skill from the boss: %s!\n", GlobalSkills[newSkill].Name)
				p.Save()
			}
		} else {
			monsterCount := gate.MonsterCount/gate.Floors + 1
			for i := 0; i < monsterCount; i++ {
				monster := Monster{
					Name:   fmt.Sprintf("%s-Rank Dungeon Beast", gate.Rank),
					Health: 20 * gate.MinLevel,
					Damage: 5 * gate.MinLevel,
				}
				if !p.Combat(&monster, true) {
					fmt.Printf("❌ Failed to clear the gate. Driven out from floor %d.\n", floor)
					return
				}
			}
			fmt.Printf("\n✅ Floor %d cleared! Moving deeper...\n", floor)
		}
	}
}

func (p *Player) ListDCraftable() {
	fmt.Println("\n--- 🛠️ Dungeon Crafting Menu ---")
	for id, r := range Recipes {
		fmt.Printf("[%s] %s (Lvl %d)\n    Materials: ", id, r.Name, r.RequiredLevel)
		var ingList []string
		for ing, qty := range r.Ingredients {
			status := "❌"
			if p.Inventory[ing] >= qty {
				status = "✅"
			}
			ingList = append(ingList, fmt.Sprintf("%s %d %s", status, qty, ing))
		}
		fmt.Printf("%s\n", strings.Join(ingList, ", "))
	}
	fmt.Println("-------------------------------")
	fmt.Println("Usage: !dcraft <item_id>")
}

func (p *Player) ListSkills() {
	fmt.Println("\n   ╔═══════════════════════╗")
	fmt.Println("   ║ 🎮 *DUNGEON SKILLS* ║")
	fmt.Println("   ╚═══════════════════════╝")
	fmt.Printf("\n   🎯 *SP:* %d  |  🎮 *Slots:* %d/%d\n", p.SkillPoints, len(p.EquippedSkills), p.SkillSlots)

	fmt.Println("\n   ━━━ *EQUIPPED* ━━━")
	if len(p.EquippedSkills) == 0 {
		fmt.Println("   (None equipped)")
	} else {
		for i, sID := range p.EquippedSkills {
			skill := GlobalSkills[sID]
			lvl := p.SkillLevels[sID]
			if lvl == 0 {
				lvl = 1
			}
			fmt.Printf("   [%d] %s Lv%d [%s-Rank]\n", i+1, skill.Name, lvl, skill.Rank)
			fmt.Printf("       _%s_\n", skill.Desc)
			if skill.Type == "active" {
				fmt.Printf("       CD: %d rounds | Type: %s\n", skill.Cooldown, skill.Type)
			} else {
				fmt.Printf("       Type: %s\n", skill.Type)
			}
		}
	}

	owned := make(map[string]bool)
	for _, s := range p.Skills {
		owned[s] = true
	}

	fmt.Println("\n   ━━━ *LOCKED* ━━━")
	// Sort by Rank for better presentation
	ranks := []string{"E", "D", "C", "B", "A", "S"}
	for _, r := range ranks {
		for _, skill := range GlobalSkills {
			if skill.Rank == r && !owned[skill.ID] {
				fmt.Printf("   🔒 %s [%s-Rank]\n", skill.Name, skill.Rank)
				fmt.Printf("       %s\n", skill.UnlockRequirement)
			}
		}
	}

	fmt.Println("\n   ━━━━━━━━━━━━━━━━━━━")
	fmt.Println("   !dequip <id> · !dunequip <slot#> · !dupskill <id>")
}

func (p *Player) UnequipSkill(slot int) {
	if slot < 1 || slot > len(p.EquippedSkills) {
		fmt.Printf("❌ Invalid slot number. (1-%d)\n", len(p.EquippedSkills))
		return
	}
	sID := p.EquippedSkills[slot-1]
	p.EquippedSkills = append(p.EquippedSkills[:slot-1], p.EquippedSkills[slot:]...)
	fmt.Printf("⚪ Unequipped %s.\n", GlobalSkills[sID].Name)
	p.Save()
}

func (p *Player) UpgradeSkill(skillID string) {
	skillID = strings.ToLower(skillID)
	owned := false
	for _, s := range p.Skills {
		if s == skillID {
			owned = true
			break
		}
	}
	if !owned {
		fmt.Printf("❌ You don't own the skill: %s\n", skillID)
		return
	}

	if p.SkillPoints < 1 {
		fmt.Println("❌ Not enough Skill Points (SP)!")
		return
	}

	if p.SkillLevels[skillID] == 0 {
		p.SkillLevels[skillID] = 1
	}
	p.SkillLevels[skillID]++
	p.SkillPoints--
	fmt.Printf("✨ Upgraded %s to Lv%d! (SP remaining: %d)\n", GlobalSkills[skillID].Name, p.SkillLevels[skillID], p.SkillPoints)
	p.Save()
}

func (p *Player) LearnSkill(skillID string) {
	skillID = strings.ToLower(skillID)
	skill, exists := GlobalSkills[skillID]
	if !exists {
		fmt.Println("❌ Skill not found.")
		return
	}

	for _, s := range p.Skills {
		if s == skillID {
			fmt.Println("❌ You already know this skill!")
			return
		}
	}

	fmt.Printf("🚫 Requirement not met for %s: %s\n", skill.Name, skill.UnlockRequirement)
}

func (p *Player) EquipSkill(skillID string) {
	skillID = strings.ToLower(skillID)
	owned := false
	for _, s := range p.Skills {
		if s == skillID {
			owned = true
			break
		}
	}
	if !owned {
		fmt.Printf("❌ You don't own the skill: %s\n", skillID)
		return
	}

	for i, eq := range p.EquippedSkills {
		if eq == skillID {
			p.EquippedSkills = append(p.EquippedSkills[:i], p.EquippedSkills[i+1:]...)
			fmt.Printf("⚪ Unequipped %s.\n", GlobalSkills[skillID].Name)
			p.Save()
			return
		}
	}

	if len(p.EquippedSkills) >= p.SkillSlots {
		fmt.Printf("🚫 No more skill slots available! (%d/%d)\n", len(p.EquippedSkills), p.SkillSlots)
		return
	}

	p.EquippedSkills = append(p.EquippedSkills, skillID)
	fmt.Printf("🟢 Equipped %s.\n", GlobalSkills[skillID].Name)
	p.Save()
}

func (p *Player) CheckTitles() {
	for id, t := range GlobalTitles {
		if p.Kills >= t.KillsNeeded {
			alreadyOwned := false
			for _, owned := range p.Titles {
				if owned == id {
					alreadyOwned = true
					break
				}
			}
			if !alreadyOwned {
				p.Titles = append(p.Titles, id)
				fmt.Printf("\n🏅 TITLE UNLOCKED: %s! 🏅\n", t.Name)
				fmt.Printf("🎁 Perk: %s\n", t.PerkDesc)
				p.Attack += t.AttackBonus
				p.MaxHealth += t.HPBonus
				p.Health += t.HPBonus
			}
		}
	}
}

func (p *Player) ListTitles() {
	fmt.Println("\n--- 🏅 Earned Titles ---")
	if len(p.Titles) == 0 {
		fmt.Println("No titles earned yet.")
	} else {
		for _, tID := range p.Titles {
			t := GlobalTitles[tID]
			fmt.Printf("%s - %s\n", t.Name, t.PerkDesc)
		}
	}
	fmt.Printf("Total Kills: %d\n", p.Kills)
	fmt.Println("------------------------")
}

func (p *Player) ShowStats() {
	fmt.Printf("\n--- 👤 Player Stats ---\n")
	fmt.Printf("⛏️ Mine Rank:   [%s] (Level: %d)\n", p.Rank, p.Level)
	fmt.Printf("🏹 Hunter Rank: [%s] (Level: %d)\n", p.HunterRank, p.HunterLevel)
	fmt.Printf("✨ Mine XP:     %d/%d\n", p.XP, p.XPToNext)
	fmt.Printf("🏹 Hunter XP:   %d/%d\n", p.HunterXP, p.HunterXPToNext)
	fmt.Printf("❤️ Health:      %d/%d\n", p.Health, p.MaxHealth)
	fmt.Printf("⚔️ Attack:      %d\n", p.Attack)
	fmt.Printf("🛡️ Defense:     %d\n", p.Defense)
	fmt.Printf("⚡ Stamina:     %d/%d\n", p.Stamina, p.MaxStamina)
	fmt.Printf("💀 Kills:       %d\n", p.Kills)
	fmt.Printf("🔨 Durability:  %d\n", p.ToolDurability)
	fmt.Printf("🔮 Skills:      %d/%d slots used\n", len(p.EquippedSkills), p.SkillSlots)
	if len(p.Titles) > 0 {
		fmt.Printf("🏅 Titles: ")
		var tNames []string
		for _, tID := range p.Titles {
			tNames = append(tNames, GlobalTitles[tID].Name)
		}
		fmt.Println(strings.Join(tNames, ", "))
	}
	if len(p.Structures) > 0 {
		fmt.Printf("🏗️ Structures: ")
		var sList []string
		for s := range p.Structures {
			sList = append(sList, s)
		}
		fmt.Println(strings.Join(sList, ", "))
	}
	fmt.Println("----------------------")
}

func (p *Player) HealFull() {
	p.Health = p.MaxHealth
	p.Stamina = p.MaxStamina
	p.Save()
}

func (p *Player) Save() {
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		fmt.Printf("❌ Error saving data: %v\n", err)
		return
	}
	os.WriteFile("player_data.json", data, 0644)
}

func LoadPlayer() *Player {
	data, err := os.ReadFile("player_data.json")
	if err != nil {
		return NewPlayer("Adventurer")
	}
	var p Player
	if err := json.Unmarshal(data, &p); err != nil {
		return NewPlayer("Adventurer")
	}
	if p.QuestProgress == nil {
		p.QuestProgress = make(map[string]int)
	}
	if p.Rank == "" {
		p.UpdateRank()
	}
	if p.HunterRank == "" {
		p.UpdateHunterRank()
	}
	if p.HunterLevel == 0 {
		p.HunterLevel = 1
		p.HunterXPToNext = 100
	}
	if p.SkillSlots == 0 {
		p.SkillSlots = 3
	}
	if p.SkillCooldowns == nil {
		p.SkillCooldowns = make(map[string]int)
	}
	if p.SkillLevels == nil {
		p.SkillLevels = make(map[string]int)
	}
	return &p
}

func (p *Player) TrackQuest(qType, id string, qty int) {
	for _, q := range GlobalQuests {
		if q.TargetType == qType && q.TargetID == id {
			if p.QuestProgress[q.ID] < q.TargetQty {
				p.QuestProgress[q.ID] += qty
				if p.QuestProgress[q.ID] >= q.TargetQty {
					fmt.Printf("\n📜 QUEST COMPLETE: %s! 📜\n", q.Name)
					fmt.Printf("🎁 Rewards: ✨ %d XP, 💰 %d Gold\n", q.RewardXP, q.RewardGold)
					p.Inventory["gold"] += q.RewardGold
					p.GainXP(q.RewardXP)
				}
			}
		}
	}
}

func (p *Player) ListQuests() {
	fmt.Println("\n--- 📜 Active Quests ---")
	for _, q := range GlobalQuests {
		status := "✅ Completed"
		prog := p.QuestProgress[q.ID]
		if prog < q.TargetQty {
			status = fmt.Sprintf("⏳ Progress: %d/%d", prog, q.TargetQty)
		}
		fmt.Printf("[%s] %s\n    📝 %s\n    📊 %s\n", q.ID, q.Name, q.Description, status)
	}
	fmt.Println("-------------------------")
}

func (p *Player) ShowInventory() {
	fmt.Printf("\n--- 🎒 Inventory ---\n")
	if len(p.Inventory) == 0 {
		fmt.Println("Empty 📭")
	} else {
		for itemID, qty := range p.Inventory {
			if qty > 0 {
				fmt.Printf("%s: %d\n", itemID, qty)
			}
		}
	}
	fmt.Println("--------------------")
}

func (p *Player) StartRegeneration() {
	ticker := time.NewTicker(20 * time.Minute)
	go func() {
		for range ticker.C {
			p.Regenerate()
		}
	}()
}

func (p *Player) Regenerate() {
	hpRegen := 10
	stRegen := 10

	if p.Structures["house"] {
		hpRegen += 2
	}
	if p.Structures["farm"] {
		stRegen += 5
	}

	if p.Health < p.MaxHealth {
		p.Health += hpRegen
		if p.Health > p.MaxHealth {
			p.Health = p.MaxHealth
		}
	}
	if p.Stamina < p.MaxStamina {
		p.Stamina += stRegen
		if p.Stamina > p.MaxStamina {
			p.Stamina = p.MaxStamina
		}
	}
	p.Save()
}

func (p *Player) StartRaids() {
	ticker := time.NewTicker(30 * time.Minute)
	go func() {
		for range ticker.C {
			if rand.Float64() < 0.3 {
				p.UnderRaid()
			}
		}
	}()
}

func (p *Player) UnderRaid() {
	fmt.Printf("\n🚨 ALERT! Your base is being raided by NPCs! 🚨\n")
	raidStrength := p.Level / 5
	if raidStrength < 1 {
		raidStrength = 1
	}
	raider := Monster{Name: "🏴‍☠️ Raider Party", Health: 50 * raidStrength, Damage: 10 * raidStrength}
	if p.Combat(&raider, false) {
		fmt.Println("🛡️ You successfully defended your base!")
	} else {
		fmt.Println("📉 The raiders plundered some of your resources!")
		for item, qty := range p.Inventory {
			if qty > 5 {
				lost := rand.Intn(qty / 2)
				p.Inventory[item] -= lost
				if lost > 0 {
					fmt.Printf("💸 Lost %d %s\n", lost, item)
				}
			}
		}
	}
	p.Save()
}

func (p *Player) ListRaids() {
	fmt.Println("\n--- ⚔️ Raid Targets ---")
	for id, s := range BotSettlements {
		fmt.Printf("[%s] %s (⭐ Lvl %d)\n    📝 %s\n", id, s.Name, s.Level, s.Description)
	}
	fmt.Println("-----------------------")
}

func (p *Player) Raid(targetID string) {
	target, ok := BotSettlements[strings.ToLower(targetID)]
	if !ok {
		fmt.Printf("❓ Unknown target: %s. Type !raid to see list.\n", targetID)
		return
	}

	if p.Inventory[target.RequiredSword] <= 0 {
		fmt.Printf("🗡️ A %s is needed to raid %s!\n", target.RequiredSword, target.Name)
		return
	}

	if p.Level < target.Level {
		fmt.Printf("🚫 Your level is too low to raid %s! Required: %d\n", target.Name, target.Level)
		return
	}
	if p.Stamina < 30 {
		fmt.Println("😫 Raiding requires 30 stamina! Wait for regeneration.")
		return
	}
	p.Stamina -= 30
	fmt.Printf("🚀 Starting raid on %s...\n", target.Name)
	for _, defender := range target.Defenders {
		fmt.Printf("⚔️ Facing defender: %s\n", defender.Name)
		if !p.Combat(&defender, false) {
			fmt.Printf("❌ Raid failed! You were driven back from %s.\n", target.Name)
			return
		}
	}
	fmt.Printf("💰 SUCCESS! You conquered %s and plundered their vault!\n", target.Name)
	for item, qty := range target.LootTable {
		p.Inventory[item] += qty
		fmt.Printf("🎁 Found %d %s\n", qty, item)
	}
	p.GainXP(100 + (target.Level * 10))
	p.Save()
}

func (p *Player) ListShop() {
	fmt.Println("\n--- ⚖️ Merchant's Shop ---")
	fmt.Printf("Your Gold: 💰 %d\n", p.Inventory["gold"])
	for id, item := range MerchantInventory {
		fmt.Printf("[%s] %s - 💰 %d\n    📝 %s\n", id, item.Name, item.Price, item.Desc)
	}
	fmt.Println("--------------------------")
}

func (p *Player) Buy(itemID string) {
	item, ok := MerchantInventory[strings.ToLower(itemID)]
	if !ok {
		fmt.Printf("❓ Merchant says: 'I don't have a %s for sale!'\n", itemID)
		return
	}
	if p.Inventory["gold"] < item.Price {
		fmt.Printf("🚫 Merchant says: 'You need 💰 %d gold for that, you only have 💰 %d!'\n", item.Price, p.Inventory["gold"])
		return
	}
	p.Inventory["gold"] -= item.Price
	switch item.ID {
	case "golden_apple":
		p.Health += 100
		if p.Health > p.MaxHealth {
			p.Health = p.MaxHealth
		}
		fmt.Printf("🍎 You bought and ate a Golden Apple! Health restored to %d.\n", p.Health)
	case "energy_drink":
		p.Stamina += 50
		if p.Stamina > p.MaxStamina {
			p.Stamina = p.MaxStamina
		}
		fmt.Printf("🥤 You bought and drank an Energy Drink! Stamina restored to %d.\n", p.Stamina)
	case "repair_kit":
		p.ToolDurability = 500
		fmt.Printf("🔧 You bought a Repair Kit! Your tool is now extremely durable (%d).\n", p.ToolDurability)
	case "mystery_box":
		fmt.Printf("🎁 You opened a Mystery Box and found: ")
		lootPool := []string{"iron", "gold", "diamond", "quartz", "netherite"}
		for i := 0; i < 3; i++ {
			loot := lootPool[rand.Intn(len(lootPool))]
			qty := 5 + rand.Intn(10)
			p.Inventory[loot] += qty
			fmt.Printf("%d %s, ", qty, loot)
		}
		fmt.Println("Not bad!")
	default:
		p.Inventory[item.ID]++
		fmt.Printf("⚖️ You bought 1 %s for 💰 %d gold.\n", item.Name, item.Price)
	}
	p.Save()
}

func (p *Player) ListCraftable() {
	fmt.Println("\n--- 📜 Crafting Menu ---")
	for id, r := range Recipes {
		fmt.Printf("[%s] ⭐ Lvl %d - 📦 Ingredients: ", id, r.RequiredLevel)
		var ingList []string
		for ing, qty := range r.Ingredients {
			ingList = append(ingList, fmt.Sprintf("%d %s", qty, ing))
		}
		fmt.Printf("%s\n", strings.Join(ingList, ", "))
	}
	fmt.Println("------------------------")
}

func (p *Player) Craft(itemName string) {
	recipe, ok := Recipes[strings.ToLower(itemName)]
	if !ok {
		fmt.Printf("❓ Unknown recipe: %s. Type !craft to see options.\n", itemName)
		return
	}
	if p.Level < recipe.RequiredLevel {
		fmt.Printf("🚫 Your level is too low to craft %s! Required: %d\n", recipe.Name, recipe.RequiredLevel)
		return
	}
	for ing, qty := range recipe.Ingredients {
		if p.Inventory[ing] < qty {
			fmt.Printf("❌ Missing ingredients for %s: Need %d %s, have %d\n", recipe.Name, qty, ing, p.Inventory[ing])
			return
		}
	}
	for ing, qty := range recipe.Ingredients {
		p.Inventory[ing] -= qty
		if p.Inventory[ing] == 0 {
			delete(p.Inventory, ing)
		}
	}
	switch recipe.ResultType {
	case "tool":
		p.ToolDurability = recipe.ResultValue
		p.Inventory[strings.ToLower(itemName)] = 1
		fmt.Printf("🛠️ You crafted a %s! Tool durability set to %d.\n", recipe.Name, p.ToolDurability)
	case "weapon":
		p.Inventory[strings.ToLower(itemName)]++
		fmt.Printf("⚔️ You crafted a %s!\n", recipe.Name)
	case "armor":
		p.Defense += recipe.ResultValue
		p.Inventory[strings.ToLower(itemName)] = 1
		fmt.Printf("🛡️ You crafted a %s! Defense increased by %d (Total: %d).\n", recipe.Name, recipe.ResultValue, p.Defense)
	case "food":
		p.Inventory[strings.ToLower(itemName)]++
		fmt.Printf("🍞 You crafted a %s!\n", recipe.Name)
	case "stamina_food":
		p.Inventory[strings.ToLower(itemName)]++
		fmt.Printf("⚡ You crafted a %s!\n", recipe.Name)
	}
	p.GainXP(10 + rand.Intn(5))
	p.Save()
}

func (p *Player) Use(itemName string) {
	itemKey := strings.ToLower(itemName)
	if p.Inventory[itemKey] <= 0 {
		fmt.Printf("❌ You don't have any %s in your inventory.\n", itemName)
		return
	}

	recipe, ok := Recipes[itemKey]
	if !ok || (recipe.ResultType != "food" && recipe.ResultType != "stamina_food") {
		fmt.Printf("❌ %s is not a consumable item.\n", itemName)
		return
	}

	p.Inventory[itemKey]--
	if p.Inventory[itemKey] == 0 {
		delete(p.Inventory, itemKey)
	}

	switch recipe.ResultType {
	case "food":
		oldHP := p.Health
		p.Health += recipe.ResultValue
		if p.Health > p.MaxHealth {
			p.Health = p.MaxHealth
		}
		fmt.Printf("😋 You consumed %s and recovered %d HP! (❤️ %d -> %d)\n", recipe.Name, p.Health-oldHP, oldHP, p.Health)
	case "stamina_food":
		oldStam := p.Stamina
		p.Stamina += recipe.ResultValue
		if p.Stamina > p.MaxStamina {
			p.Stamina = p.MaxStamina
		}
		fmt.Printf("⚡ You consumed %s and recovered %d Stamina! (⚡ %d -> %d)\n", recipe.Name, p.Stamina-oldStam, oldStam, p.Stamina)
	}
	p.Save()
}

func (p *Player) ListBuildable() {
	fmt.Println("\n--- 🏗️ Building Menu ---")
	for id, s := range Structures {
		fmt.Printf("[%s] ⭐ Lvl %d - 📦 Ingredients: ", id, s.RequiredLevel)
		var ingList []string
		for ing, qty := range s.Ingredients {
			ingList = append(ingList, fmt.Sprintf("%d %s", qty, ing))
		}
		fmt.Printf("%s\n    🎁 Perk: %s\n", strings.Join(ingList, ", "), s.PerkDesc)
	}
	fmt.Println("------------------------")
}

func (p *Player) Build(structName string) {
	s, ok := Structures[strings.ToLower(structName)]
	if !ok {
		fmt.Printf("❓ Unknown structure: %s. Type !build to see options.\n", structName)
		return
	}
	if p.Structures[strings.ToLower(structName)] {
		fmt.Printf("🏠 You already built a %s!\n", s.Name)
		return
	}
	if p.Level < s.RequiredLevel {
		fmt.Printf("🚫 Your level is too low to build %s! Required: %d\n", s.Name, s.RequiredLevel)
		return
	}
	for ing, qty := range s.Ingredients {
		if p.Inventory[ing] < qty {
			fmt.Printf("❌ Missing materials for %s: Need %d %s, have %d\n", s.Name, qty, ing, p.Inventory[ing])
			return
		}
	}
	for ing, qty := range s.Ingredients {
		p.Inventory[ing] -= qty
		if p.Inventory[ing] == 0 {
			delete(p.Inventory, ing)
		}
	}
	p.Structures[strings.ToLower(structName)] = true
	fmt.Printf("🔨 You built a %s! Perk Unlocked: %s\n", s.Name, s.PerkDesc)
	switch strings.ToLower(structName) {
	case "forge":
		p.Attack += 10
	case "vault":
		p.MaxHealth += 50
		p.Health += 50
	case "castle":
		p.Attack += 20
		p.MaxHealth += 100
		p.Health += 100
		p.MaxStamina += 50
		p.Stamina += 50
	}
	p.GainXP(50 + rand.Intn(50))
	p.Save()
}

func (p *Player) GetBestSwordDamage() int {
	bestDmg := 0
	for id, qty := range p.Inventory {
		if qty > 0 {
			if r, ok := Recipes[id]; ok && r.ResultType == "weapon" {
				if r.ResultValue > bestDmg {
					bestDmg = r.ResultValue
				}
			}
		}
	}
	return bestDmg
}

func (p *Player) GetBestPickaxeMultiplier() float64 {
	multi := 1.0
	multipliers := map[string]float64{
		"wood_pickaxe":      1.0,
		"stone_pickaxe":     1.2,
		"iron_pickaxe":      1.5,
		"diamond_pickaxe":   2.0,
		"abyss_pickaxe":     3.0,
		"nether_pickaxe":    5.0,
		"void_pickaxe":      10.0,
	}

	for id, qty := range p.Inventory {
		if qty > 0 {
			if m, ok := multipliers[id]; ok {
				if m > multi {
					multi = m
				}
			}
		}
	}
	return multi
}

func (p *Player) Mine(locName string) {
	loc, ok := Locations[strings.ToLower(locName)]
	if !ok {
		fmt.Printf("❓ Unknown location: %s. Type !mine to see available zones.\n", locName)
		return
	}
	if p.Level < loc.RequiredLevel {
		fmt.Printf("🚫 Your level is too low to enter %s! Required: %d\n", loc.Name, loc.RequiredLevel)
		return
	}
	if loc.RequiredItem != "" && p.Inventory[loc.RequiredItem] <= 0 {
		fmt.Printf("🔏 You need a %s to mine in the %s!\n", loc.RequiredItem, loc.Name)
		return
	}
	if p.Stamina < 10 {
		fmt.Println("😫 Not enough stamina! Wait for regeneration.")
		return
	}
	if p.ToolDurability <= 0 {
		fmt.Println("⚠️ Your tool is broken! Craft a new one.")
		return
	}
	p.Stamina -= 10
	p.ToolDurability -= 1

	if len(loc.Descriptions) > 0 {
		desc := loc.Descriptions[rand.Intn(len(loc.Descriptions))]
		fmt.Printf("\n✨ %s\n", desc)
	}

	if rand.Float64() <= loc.EncounterChance {
		monster := loc.EncounterTable[rand.Intn(len(loc.EncounterTable))]
		if !p.Combat(&monster, false) {
			return
		}
	}

	pickMulti := p.GetBestPickaxeMultiplier()
	numDrops := int(float64(1+(p.Level/5)) * pickMulti)
	
	foundSomething := false
	for i := 0; i < numDrops; i++ {
		r := rand.Float64()
		var cumulative float64
		for item, prob := range loc.LootTable {
			cumulative += prob
			if r <= cumulative {
				p.Inventory[item]++
				fmt.Printf("⛏️ You mined in the %s and found: %s!\n", loc.Name, item)
				p.TrackQuest("item", item, 1)
				foundSomething = true
				break
			}
		}
	}
	if !foundSomething {
		fmt.Printf("💨 You mined in the %s but found nothing.\n", loc.Name)
	}
	p.GainXP(2 + rand.Intn(3))
	p.Save()
}
