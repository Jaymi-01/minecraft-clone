package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

func NewPlayer(name string) *Player {
	return &Player{
		Name:             name,
		Health:           100, MaxHealth: 100, Attack: 10, Defense: 0, Stamina: 50, MaxStamina: 50, Magic: 100, MaxMagic: 100,
		Level:            1, XP: 0, XPToNext: 100, HunterLevel: 1, HunterRank: "E", 
		Inventory:        map[string]int{"wood_pickaxe": 1, "gold": 100},
		ToolDurability:   50, 
		Structures:       make(map[string]bool), 
		QuestProgress:    make(map[string]int), 
		Rank:             "E", 
		MonsterKills:     make(map[string]int),
		SkillLevels:      make(map[string]int), 
		SkillUsage:       make(map[string]int), 
		SkillCooldowns:   make(map[string]int),
		Subordinates:     []Subordinate{}, 
		Squad:            []string{}, 
		ItemRarities:     make(map[string]string), 
		ItemLevels:       make(map[string]int),
		ItemRunes:        make(map[string][]string),
		Training:         TrainingProgress{LastReset: time.Now()}, 
		Production:       ProductionLog{LastProduced: time.Now(), PendingItems: make(map[string]int)},
		SystemOrigin:     "Human", 
		Attributes:       make(map[string]bool), 
		StatusEffects:    make(map[string]int),
		DomainStructures: make(map[string]int),
	}
}

func (p *Player) Save() { d, _ := json.MarshalIndent(p, "", "  "); os.WriteFile("player_data.json.bak", d, 0644); os.WriteFile("player_data.json", d, 0644) }

func LoadPlayer() *Player {
	data, err := os.ReadFile("player_data.json")
	if err != nil { data, err = os.ReadFile("player_data.json.bak"); if err != nil { return NewPlayer("Adventurer") } }
	var p Player; json.Unmarshal(data, &p)
	
	// Rigorous map and slice initialization
	if p.Inventory == nil { p.Inventory = make(map[string]int) }
	if p.QuestProgress == nil { p.QuestProgress = make(map[string]int) }
	if p.MonsterKills == nil { p.MonsterKills = make(map[string]int) }
	if p.SkillLevels == nil { p.SkillLevels = make(map[string]int) }
	if p.SkillUsage == nil { p.SkillUsage = make(map[string]int) }
	if p.SkillCooldowns == nil { p.SkillCooldowns = make(map[string]int) }
	if p.Subordinates == nil { p.Subordinates = []Subordinate{} }
	if p.Squad == nil { p.Squad = []string{} }
	if p.Attributes == nil { p.Attributes = make(map[string]bool) }
	if p.StatusEffects == nil { p.StatusEffects = make(map[string]int) }
	if p.Structures == nil { p.Structures = make(map[string]bool) }
	if p.ItemRarities == nil { p.ItemRarities = make(map[string]string) }
	if p.ItemLevels == nil { p.ItemLevels = make(map[string]int) }
	if p.ItemRunes == nil { p.ItemRunes = make(map[string][]string) }
	if p.DomainStructures == nil { p.DomainStructures = make(map[string]int) }
	
	p.InCombat = false // Reset state on load
	p.TrialActive = false

	p.UpdateRank(); p.UpdateHunterRank(); p.UpdateSkillSlots(); p.SyncStats()
	return &p
}

func (p *Player) WorldNotice(msg string) { fmt.Printf("\n<< SYSTEM: %s >>\n", strings.ToUpper(msg)); p.LogAction(msg) }
func (p *Player) LogAction(msg string) { t := time.Now().Format("15:04:05"); p.ActionLog = append([]string{fmt.Sprintf("[%s] %s", t, msg)}, p.ActionLog...); if len(p.ActionLog) > 30 { p.ActionLog = p.ActionLog[:30] } }

func (p *Player) UpdateRank() { if p.Level >= 150 { p.Rank = "SS" } else if p.Level >= 100 { p.Rank = "S" } else if p.Level >= 75 { p.Rank = "A" } else if p.Level >= 50 { p.Rank = "B" } else if p.Level >= 30 { p.Rank = "C" } else if p.Level >= 15 { p.Rank = "D" } else { p.Rank = "E" } }
func (p *Player) UpdateHunterRank() { if p.HunterLevel >= 150 { p.HunterRank = "SS" } else if p.HunterLevel >= 100 { p.HunterRank = "S" } else if p.HunterLevel >= 75 { p.HunterRank = "A" } else if p.HunterLevel >= 50 { p.HunterRank = "B" } else if p.HunterLevel >= 30 { p.HunterRank = "C" } else if p.HunterLevel >= 15 { p.HunterRank = "D" } else { p.HunterRank = "E" } }
func (p *Player) UpdateSkillSlots() { p.SkillSlots = 5 + (p.Level / 5); if p.Attributes["shadow_army_expansion"] { p.SkillSlots += 3 } }

func (p *Player) GainXP(amount int) {
	if p.Structures["enchanting_table"] { amount = int(float64(amount) * 1.5) }
	
	// Domain Bonus: Research Lab
	if labLvl := p.DomainStructures["research_lab"]; labLvl > 0 {
		bonus := float64(labLvl) * 0.1
		amount = int(float64(amount) * (1.0 + bonus))
	}

	p.XP += amount
	fmt.Printf("[✨ +%d EXPERIENCE POINTS]\n", amount)
	for p.XP >= p.XPToNext {
		p.Level++; p.XP -= p.XPToNext; p.XPToNext = int(float64(p.XPToNext) * 1.5)
		p.UpdateRank(); p.UpdateSkillSlots(); p.SyncStats()
		p.WorldNotice(fmt.Sprintf("CONGRATULATIONS: You have reached Level %d. Limits broken.", p.Level))
		if p.Level%10 == 0 { fmt.Println("🔥 [SYSTEM]: Stat synchronization complete. Significant growth detected.") }
	}
	p.Save()
}

func (p *Player) GainHunterXP(amount int) {
	p.HunterXP += amount
	fmt.Printf("[🏹 +%d HUNTER XP]\n", amount)
	for p.HunterXP >= p.HunterXPToNext {
		p.HunterLevel++; p.HunterXP -= p.HunterXPToNext; p.HunterXPToNext = int(float64(p.HunterXPToNext) * 1.5)
		p.UpdateHunterRank()
		p.WorldNotice(fmt.Sprintf("HUNTER PROMOTION: Rank %s achieved (Level %d).", p.HunterRank, p.HunterLevel))
	}
	p.Save()
}

func (p *Player) SyncStats() {
	// Base Stats from Level
	baseHP := 100 + ((p.Level - 1) * 20)
	baseMP := 100 + ((p.Level - 1) * 30)
	baseStamina := 50 + ((p.Level - 1) * 15)
	baseAtk := 10 + ((p.Level - 1) * 5)
	baseDef := 0 + ((p.Level - 1) * 2)

	// Evolution Multipliers & Flat Bonuses
	multiplier := 1.0
	flatHP, flatMP, flatAtk, flatDef := 0, 0, 0, 0

	switch p.SystemOrigin {
	case "Slime", "Spider":
		multiplier = 1.2; flatHP = 100; flatMP = 100; flatAtk = 20
	case "Small Poison Taratect":
		multiplier = 1.5; flatHP = 250; flatMP = 250; flatAtk = 50; flatDef = 20
	case "Demon Slime", "Arachne":
		multiplier = 2.5; flatHP = 1000; flatMP = 1000; flatAtk = 150; flatDef = 100
	case "Ultimate Slime (True Dragon)", "God (Shiraori)":
		multiplier = 5.0; flatHP = 5000; flatMP = 5000; flatAtk = 500; flatDef = 300
	}

	// Apply Evolution
	p.MaxHealth = int(float64(baseHP+flatHP) * multiplier)
	p.MaxMagic = int(float64(baseMP+flatMP) * multiplier)
	p.MaxStamina = int(float64(baseStamina) * multiplier)
	p.Attack = int(float64(baseAtk+flatAtk) * multiplier)
	p.Defense = int(float64(baseDef+flatDef) * multiplier)

	// Title Bonuses
	for _, tID := range p.Titles {
		if t, ok := GlobalTitles[tID]; ok {
			p.Attack += t.AttackBonus; p.MaxHealth += t.HPBonus; p.Defense += t.DefenseBonus; p.MaxMagic += t.MPBonus; p.MaxStamina += t.StaminaBonus
		}
	}

	// Rune Bonuses
	for _, runes := range p.ItemRunes {
		for _, runeID := range runes {
			if runeID == "mana_rune" { p.MaxMagic += 50 }
			if runeID == "defense_rune" { p.Defense += 20 }
		}
	}

	// Structure Bonuses
	if p.Structures["vault"] { p.MaxHealth += 50 }
	if p.Structures["castle"] { p.Attack += 20; p.MaxHealth += 100; p.MaxStamina += 50 }
	if p.Structures["forge"] { p.Attack += 10 }

	// Cap current values
	if p.Health > p.MaxHealth { p.Health = p.MaxHealth }
	if p.Stamina > p.MaxStamina { p.Stamina = p.MaxStamina }
	if p.Magic > p.MaxMagic { p.Magic = p.MaxMagic }
}

func (p *Player) HealFull() { p.Health = p.MaxHealth; p.Magic = p.MaxMagic; p.Stamina = p.MaxStamina; p.Save() }
