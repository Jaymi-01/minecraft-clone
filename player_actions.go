package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func (p *Player) Mine(locName string) {
	locID := strings.ToLower(locName); loc, ok := Locations[locID]; if !ok || p.Level < loc.RequiredLevel { return }
	if loc.RequiredItem != "" && p.Inventory[loc.RequiredItem] <= 0 { fmt.Printf("🚫 Need %s.\n", loc.RequiredItem); return }
	if p.Stamina < 10 { p.WorldNotice("EXHAUSTED"); return }; if p.ToolDurability <= 0 { p.WorldNotice("TOOL BROKEN"); return }
	p.Stamina -= 10; p.ToolDurability -= 1
	if rand.Float64() <= loc.EncounterChance { if !p.Combat(&loc.EncounterTable[rand.Intn(len(loc.EncounterTable))], false) { return } }
	pick := 1.0; multis := map[string]float64{"wood_pickaxe":1, "stone_pickaxe":1.2, "iron_pickaxe":1.5, "diamond_pickaxe":2, "abyss_pickaxe":3, "nether_pickaxe":5, "void_pickaxe":10}
	for id, q := range p.Inventory { if q > 0 { if val, ok := multis[id]; ok && val > pick { pick = val } } }
	drops := int(float64(1+p.Level/5) * pick); foundItems := make(map[string]int)
	for i := 0; i < drops; i++ {
		r := rand.Float64(); var cum float64; for item, prob := range loc.LootTable { cum += prob; if r <= cum { p.Inventory[item]++; foundItems[item]++; p.TrackQuest("item", item, 1); break } }
	}
	if len(foundItems) > 0 { for id, qty := range foundItems { fmt.Printf("🎁 [GATHERED] %s x%d\n", id, qty) } }
	p.GainXP(2 + rand.Intn(3)); p.Save()
}

func (p *Player) EnterGate(isAdmin bool) {
	if p.CurrentGate == nil { fmt.Println("📭 No gate manifested."); return }
	if !isAdmin && p.Level < p.CurrentGate.MinLevel { fmt.Printf("🚫 Min Level %d required!\n", p.CurrentGate.MinLevel); return }
	if p.CurrentGate.MinLevel >= 10 {
		if !strings.Contains(p.EquippedWeapon, "iron") && !strings.Contains(p.EquippedWeapon, "diamond") && !strings.Contains(p.EquippedWeapon, "void") && !strings.HasPrefix(p.EquippedWeapon, "d_") {
			fmt.Println("🚫 DANGEROUS: Iron Sword or better required for Level 10+ Gates."); return
		}
	}
	fmt.Printf("🌀 Entering %s Gate...\n", p.CurrentGate.Rank)
	for f := 1; f <= p.CurrentGate.Floors; f++ {
		fmt.Printf("\n🏢 FLOOR %d / %d\n", f, p.CurrentGate.Floors)
		if f == p.CurrentGate.Floors {
			if p.Combat(&p.CurrentGate.Boss, true) { p.Inventory["gold"] += p.CurrentGate.RewardGold; p.GainHunterXP(p.CurrentGate.RewardXP); p.CurrentGate = nil; p.Save() }
			return
		}
		for i := 0; i < 3; i++ { m := Monster{Name: "Gate Beast", Health: 20 * p.Level, Damage: 5 * p.Level}; if !p.Combat(&m, true) { return } }
	}
}

func (p *Player) Raid(targetID string) {
	t, ok := BotSettlements[strings.ToLower(targetID)]; if !ok { return }
	if p.Level < t.Level || p.Stamina < 30 { fmt.Println("🚫 Insufficient preparation."); return }
	p.Stamina -= 30; p.WorldNotice(fmt.Sprintf("Commencing Raid on %s", t.Name))
	for _, d := range t.Defenders { if !p.Combat(&d, false) { return } }
	for id, qty := range t.LootTable { p.Inventory[id] += qty; fmt.Printf("   💰 Plundered: %s x%d\n", id, qty) }
	p.GainTaboo(1); p.GainXP(100 + t.Level*10); p.Save(); p.WorldNotice(fmt.Sprintf("Raid on %s successful.", t.Name))
}

func (p *Player) StartSubordinateAutonomy() {
	ticker := time.NewTicker(5 * time.Minute)
	go func() { for range ticker.C { for i := range p.Subordinates { p.SubordinateAction(&p.Subordinates[i]) }; p.Save() } }()
}

func (p *Player) SubordinateAction(s *Subordinate) {
	if time.Since(s.LastAction) < 5*time.Minute { return }; s.LastAction = time.Now(); action := rand.Intn(100)
	if action < 40 {
		locs := []string{"surface", "cave", "abyss", "nether", "void"}; l := Locations[locs[rand.Intn(len(locs))]]
		if s.Level >= l.RequiredLevel {
			p.LogAction(s.Name + " mining in " + l.Name); for it, pr := range l.LootTable { if rand.Float64() <= pr { p.Inventory[it]++; p.LogAction(s.Name + " found " + it) } }
			p.SubordinateGainXPForOne(s, 20); p.GainXP(10)
		}
	} else if action < 70 {
		raids := []string{"goblin_camp", "bandit_fort", "shadow_keep"}; r := BotSettlements[raids[rand.Intn(len(raids))]]
		if s.Level >= r.Level {
			p.LogAction(s.Name + " raiding " + r.Name); for it, q := range r.LootTable { p.Inventory[it] += q }
			p.SubordinateGainXPForOne(s, 50); p.GainXP(25)
		}
	}
}

func (p *Player) SubordinateGainXP(amount int) { for i := range p.Subordinates { p.SubordinateGainXPForOne(&p.Subordinates[i], amount) } }

func (p *Player) SubordinateGainXPForOne(s *Subordinate, amount int) {
	if s.NextXP == 0 { s.NextXP = 100 }; s.XP += amount
	if s.XP >= s.NextXP {
		s.Level++; s.XP -= s.NextXP; s.NextXP = int(float64(s.NextXP) * 1.5); s.Attack += 10; s.Defense += 10
		p.WorldNotice(fmt.Sprintf("%s reached Lv%d", s.Name, s.Level)); p.CheckSubordinateEvolution(s); p.CheckSubordinateSkills(s)
	}
}

func (p *Player) CheckSubordinateSkills(s *Subordinate) {
	skills := map[string][]struct{lvl int; id string}{"slime": {{1, "predator"}, {5, "water_jet"}, {15, "gluttony"}, {30, "lightning"}}, "spider": {{1, "appraisal"}, {5, "venom_spit"}, {15, "evil_eye"}, {30, "heresy_magic"}}}
	for _, ss := range skills[s.Species] {
		if s.Level >= ss.lvl {
			owned := false; for _, sk := range s.Skills { if sk == ss.id { owned = true; break } }
			if !owned { s.Skills = append(s.Skills, ss.id); p.WorldNotice(s.Name + " learned " + GlobalSkills[ss.id].Name) }
		}
	}
}

func (p *Player) CheckSubordinateEvolution(s *Subordinate) {
	if s.Species == "hobgoblin" && s.Level >= 10 { s.Species = "ogre" } else if s.Species == "ogre" && s.Level >= 25 { s.Species = "kijin" } else if s.Species == "alpha wolf" && s.Level >= 15 { s.Species = "tempest wolf" }
}
