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
		Name:           name,
		Health:         100, MaxHealth: 100, Attack: 10, Defense: 0, Stamina: 50, MaxStamina: 50, Magic: 100, MaxMagic: 100,
		Level:          1, XP: 0, XPToNext: 100, HunterLevel: 1, HunterRank: "E", Inventory: map[string]int{"wood_pickaxe": 1, "gold": 100},
		ToolDurability: 50, Structures: make(map[string]bool), QuestProgress: make(map[string]int), Rank: "E", MonsterKills: make(map[string]int),
		SkillLevels: make(map[string]int), SkillUsage: make(map[string]int), SkillCooldowns: make(map[string]int),
		Subordinates: []Subordinate{}, Squad: []string{}, ItemRarities: make(map[string]string), ItemLevels: make(map[string]int),
		Training: TrainingProgress{LastReset: time.Now()}, Production: ProductionLog{LastProduced: time.Now(), PendingItems: make(map[string]int)},
		SystemOrigin: "Human", Attributes: make(map[string]bool), StatusEffects: make(map[string]int),
	}
}

func (p *Player) Save() { d, _ := json.MarshalIndent(p, "", "  "); os.WriteFile("player_data.json.bak", d, 0644); os.WriteFile("player_data.json", d, 0644) }

func LoadPlayer() *Player {
	data, err := os.ReadFile("player_data.json")
	if err != nil { data, err = os.ReadFile("player_data.json.bak"); if err != nil { return NewPlayer("Adventurer") } }
	var p Player; json.Unmarshal(data, &p)
	
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
	
	p.UpdateRank(); p.UpdateHunterRank(); p.UpdateSkillSlots(); p.SyncStats()
	return &p
}

func (p *Player) WorldNotice(msg string) { fmt.Printf("\n<< NOTICE: %s >>\n", strings.ToUpper(msg)); p.LogAction(msg) }
func (p *Player) LogAction(msg string) { t := time.Now().Format("15:04:05"); p.ActionLog = append([]string{fmt.Sprintf("[%s] %s", t, msg)}, p.ActionLog...); if len(p.ActionLog) > 30 { p.ActionLog = p.ActionLog[:30] } }

func (p *Player) UpdateRank() { if p.Level >= 150 { p.Rank = "SS" } else if p.Level >= 100 { p.Rank = "S" } else if p.Level >= 75 { p.Rank = "A" } else if p.Level >= 50 { p.Rank = "B" } else if p.Level >= 30 { p.Rank = "C" } else if p.Level >= 15 { p.Rank = "D" } else { p.Rank = "E" } }
func (p *Player) UpdateHunterRank() { if p.HunterLevel >= 150 { p.HunterRank = "SS" } else if p.HunterLevel >= 100 { p.HunterRank = "S" } else if p.HunterLevel >= 75 { p.HunterRank = "A" } else if p.HunterLevel >= 50 { p.HunterRank = "B" } else if p.HunterLevel >= 30 { p.HunterRank = "C" } else if p.HunterLevel >= 15 { p.HunterRank = "D" } else { p.HunterRank = "E" } }
func (p *Player) UpdateSkillSlots() { p.SkillSlots = 5 + (p.Level / 5); if p.Attributes["shadow_army_expansion"] { p.SkillSlots += 3 } }

func (p *Player) GainXP(amount int) {
	if p.Structures["enchanting_table"] { amount = int(float64(amount) * 1.5) }; p.XP += amount
	for p.XP >= p.XPToNext { p.Level++; p.XP -= p.XPToNext; p.XPToNext = int(float64(p.XPToNext) * 1.5); p.SyncStats(); p.WorldNotice(fmt.Sprintf("Level %d reached.", p.Level)) }; p.Save()
}

func (p *Player) GainHunterXP(amount int) {
	p.HunterXP += amount; for p.HunterXP >= p.HunterXPToNext { p.HunterLevel++; p.HunterXP -= p.HunterXPToNext; p.HunterXPToNext = int(float64(p.HunterXPToNext) * 1.5); p.UpdateHunterRank(); p.WorldNotice(fmt.Sprintf("Hunter Lv%d achieved.", p.HunterLevel)) }; p.Save()
}

func (p *Player) SyncStats() {
	p.MaxHealth = 100 + ((p.Level - 1) * 10); p.MaxStamina = 50 + ((p.Level - 1) * 10); p.MaxMagic = 100 + ((p.Level - 1) * 20); p.Attack = 10; p.Defense = 0
	for _, tID := range p.Titles { if t, ok := GlobalTitles[tID]; ok { p.Attack += t.AttackBonus; p.MaxHealth += t.HPBonus; p.Defense += t.DefenseBonus; p.MaxMagic += t.MPBonus; p.MaxStamina += t.StaminaBonus } }
	if p.Structures["vault"] { p.MaxHealth += 50 }; if p.Structures["castle"] { p.Attack += 20; p.MaxHealth += 100; p.MaxStamina += 50 }; if p.Structures["forge"] { p.Attack += 10 }
	if p.Health > p.MaxHealth { p.Health = p.MaxHealth }; if p.Stamina > p.MaxStamina { p.Stamina = p.MaxStamina }; if p.Magic > p.MaxMagic { p.Magic = p.MaxMagic }
}

func (p *Player) HealFull() { p.Health = p.MaxHealth; p.Magic = p.MaxMagic; p.Stamina = p.MaxStamina; p.Save() }
