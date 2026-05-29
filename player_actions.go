package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func (p *Player) Mine(locName string) {
	locID := strings.ToLower(locName); loc, ok := Locations[locID]
	if !ok { fmt.Printf("❌ [SYSTEM]: Location '%s' is not mapped.\n", locName); return }
	if p.Level < loc.RequiredLevel { fmt.Printf("🚫 [SYSTEM]: Strength insufficient. Level %d required for %s.\n", loc.RequiredLevel, loc.Name); return }
	if loc.RequiredItem != "" && p.Inventory[loc.RequiredItem] <= 0 { fmt.Printf("🚫 [SYSTEM]: Specialized equipment required: %s.\n", loc.RequiredItem); return }
	if p.Stamina < 10 { p.WorldNotice("EXHAUSTED: Your body refuses to move. Rest required."); return }
	if p.ToolDurability <= 0 { p.WorldNotice("TOOL BROKEN: The extraction instrument has reached its limit."); return }
	
	p.Stamina -= 10; p.ToolDurability -= 1
	fmt.Printf("\n⚒️ [EXTRACTION]: Commencing resource harvest in %s...\n", loc.Name)
	
	if p.HasSkill("miners_instinct") {
		p.SkillUsage["miners_instinct"]++
		if p.SkillUsage["miners_instinct"] >= 10 { p.UpgradeSkill("miners_instinct", true) }
	}
	
	if rand.Float64() <= loc.EncounterChance {
		m := &loc.EncounterTable[rand.Intn(len(loc.EncounterTable))]
		fmt.Printf("⚠️ [AMBUSH]: A wild %s has interrupted your work!\n", m.Name)
		if !p.Combat(m, false) { return }
	}

	pick := 1.0; multis := map[string]float64{"wood_pickaxe":1, "stone_pickaxe":1.2, "iron_pickaxe":1.5, "diamond_pickaxe":2, "abyss_pickaxe":3, "nether_pickaxe":5, "void_pickaxe":10}
	for id, q := range p.Inventory { if q > 0 { if val, ok := multis[id]; ok && val > pick { pick = val } } }
	
	drops := int(float64(1+p.Level/5) * pick); foundItems := make(map[string]int)
	for i := 0; i < drops; i++ {
		r := rand.Float64(); var cum float64; for item, prob := range loc.LootTable { cum += prob; if r <= cum { p.Inventory[item]++; foundItems[item]++; p.TrackQuest("item", item, 1); break } }
	}
	
	if len(foundItems) > 0 {
		fmt.Print("🎁 [SUCCESS]: Gathered materials: ")
		for id, qty := range foundItems { fmt.Printf("[%s x%d] ", strings.ToUpper(id), qty) }
		fmt.Println()
	} else {
		fmt.Println("⚪ [RESULT]: The extraction yielded no valuable resources.")
	}
	p.GainXP(2 + rand.Intn(3)); p.Save()
}

func (p *Player) EnterGate(isAdmin bool) {
	if p.CurrentGate == nil { fmt.Println("📭 [SYSTEM]: No Gates are currently manifest in this sector."); return }
	if !isAdmin && p.Level < p.CurrentGate.MinLevel { fmt.Printf("🚫 [SYSTEM]: Entry denied. Your current power level (%d) is below the threshold (%d).\n", p.Level, p.CurrentGate.MinLevel); return }
	
	if p.CurrentGate.MinLevel >= 10 {
		if !strings.Contains(p.EquippedWeapon, "iron") && !strings.Contains(p.EquippedWeapon, "diamond") && !strings.Contains(p.EquippedWeapon, "void") && !strings.HasPrefix(p.EquippedWeapon, "d_") {
			fmt.Println("🚫 [WARNING]: Gate integrity too high. Iron Sword or better required to breach the interior."); return
		}
	}
	
	p.WorldNotice(fmt.Sprintf("GATE BREACHED: You have entered a Rank %s Gate.", p.CurrentGate.Rank))
	for f := 1; f <= p.CurrentGate.Floors; f++ {
		fmt.Printf("\n🏢 [GATE: FLOOR %d / %d]\n", f, p.CurrentGate.Floors)
		if f == p.CurrentGate.Floors {
			fmt.Printf("👹 [BOSS ROOM]: The core of the Gate has been reached. %s awaits.\n", p.CurrentGate.Boss.Name)
			if p.Combat(&p.CurrentGate.Boss, true) {
				p.Inventory["gold"] += p.CurrentGate.RewardGold
				p.GainHunterXP(p.CurrentGate.RewardXP)
				p.WorldNotice("GATE CLEARED: The manifestation has collapsed.")
				p.CurrentGate = nil; p.Save()
			}
			return
		}
		for i := 0; i < 2; i++ {
			m := Monster{Name: "Gate Sentinel", Health: 20 * p.Level, Damage: 5 * p.Level}
			if !p.Combat(&m, true) { return }
		}
	}
}

func (p *Player) Raid(targetID string) {
	t, ok := BotSettlements[strings.ToLower(targetID)]
	if !ok { fmt.Println("❌ [SYSTEM]: Target settlement not found."); return }
	if p.Level < t.Level { fmt.Printf("🚫 [SYSTEM]: Settlement defenses too strong for your current level (%d < %d).\n", p.Level, t.Level); return }
	if p.Stamina < 30 { fmt.Println("🚫 [SYSTEM]: Stamina too low for a full-scale assault."); return }
	
	p.Stamina -= 30; p.WorldNotice(fmt.Sprintf("RAID COMMENCED: Assaulting %s!", t.Name))
	for _, d := range t.Defenders { if !p.Combat(&d, false) { return } }
	
	fmt.Println("\n💰 [LOOT PLUNDERED]:")
	for id, qty := range t.LootTable { p.Inventory[id] += qty; fmt.Printf("   - %s x%d\n", strings.ToUpper(id), qty) }
	
	p.GainTaboo(1); p.GainXP(100 + t.Level*10)
	p.WorldNotice(fmt.Sprintf("RAID SUCCESSFUL: %s has been completely plundered.", t.Name))
	p.Save()
}

func (p *Player) StartSubordinateAutonomy() {
	ticker := time.NewTicker(5 * time.Minute)
	go func() { for range ticker.C { for i := range p.Subordinates { p.SubordinateAction(&p.Subordinates[i]) }; p.Save() } }()
}

func (p *Player) SubordinateAction(s *Subordinate) {
	if time.Since(s.LastAction) < 5*time.Minute { return }; s.LastAction = time.Now()
	action := rand.Intn(100)
	if action < 40 {
		locs := []string{"surface", "cave", "abyss", "nether", "void"}
		l := Locations[locs[rand.Intn(len(locs))]]
		if s.Level >= l.RequiredLevel {
			p.LogAction(fmt.Sprintf("SCOUTING: %s has successfully navigated the %s.", s.Name, l.Name))
			found := false
			foundMsg := ""
			for it, pr := range l.LootTable { 
				if rand.Float64() <= pr { 
					p.Inventory[it]++
					foundMsg += fmt.Sprintf("[%s] ", strings.ToUpper(it))
					found = true 
				} 
			}
			if found { 
				p.LogAction(fmt.Sprintf("RECOVERY: %s returned with %s.", s.Name, foundMsg))
				p.SubordinateGainXPForOne(s, 20); p.GainXP(10) 
			}
		}
	} else if action < 70 {
		raids := []string{"goblin_camp", "bandit_fort", "shadow_keep"}
		r := BotSettlements[raids[rand.Intn(len(raids))]]
		if s.Level >= r.Level {
			p.LogAction(fmt.Sprintf("INCURSION: %s has breached the defenses of %s.", s.Name, r.Name))
			lootMsg := ""
			for it, q := range r.LootTable { 
				p.Inventory[it] += q
				lootMsg += fmt.Sprintf("[%s x%d] ", strings.ToUpper(it), q)
			}
			p.LogAction(fmt.Sprintf("PLUNDER: %s acquired %s from the raid.", s.Name, lootMsg))
			p.SubordinateGainXPForOne(s, 50); p.GainXP(25)
		}
	}
}

func (p *Player) SubordinateGainXP(amount int) { for i := range p.Subordinates { p.SubordinateGainXPForOne(&p.Subordinates[i], amount) } }

func (p *Player) SubordinateGainXPForOne(s *Subordinate, amount int) {
	// Domain Bonus: Training Ground
	if tgLvl := p.DomainStructures["training_ground"]; tgLvl > 0 {
		bonus := float64(tgLvl) * 0.1
		amount = int(float64(amount) * (1.0 + bonus))
	}

	if s.NextXP == 0 { s.NextXP = 100 }; s.XP += amount
	if s.XP >= s.NextXP {
		s.Level++; s.XP -= s.NextXP; s.NextXP = int(float64(s.NextXP) * 1.5); s.Attack += 10; s.Defense += 10
		p.WorldNotice(fmt.Sprintf("SUBORDINATE GROWTH: %s has achieved Level %d. Potential increased.", s.Name, s.Level))
		p.CheckSubordinateEvolution(s); p.CheckSubordinateSkills(s)
	}
}

func (p *Player) CheckSubordinateSkills(s *Subordinate) {
	// Level-based skill mapping for all species
	skills := map[string][]struct{lvl int; id string}{
		"slime": {
			{1, "predator"}, {5, "water_jet"}, {10, "gluttony"}, 
			{20, "great_sage"}, {30, "beelzebuth"}, {50, "raphael"},
		},
		"spider": {
			{1, "appraisal"}, {3, "venom_spit"}, {7, "spider_thread"}, 
			{15, "evil_eye"}, {25, "parallel_minds"}, {40, "annihilating_eye"},
		},
		"taratect": {
			{1, "appraisal"}, {5, "poison_fang"}, {10, "steel_thread"}, 
			{20, "parallel_minds"}, {35, "heresy_magic"},
		},
		"wolf": {
			{1, "power_strike"}, {10, "heavy_cleave"}, {20, "spark"},
		},
		"alpha wolf": {
			{1, "heavy_cleave"}, {10, "armor_break"}, {25, "chain_lightning"},
		},
		"goblin": {
			{1, "power_strike"}, {5, "natures_touch"}, {15, "bone_armor"},
		},
		"hobgoblin": {
			{1, "heavy_cleave"}, {10, "armor_break"}, {20, "earth_shatter"},
		},
		"ogre": {
			{1, "armor_break"}, {15, "earth_shatter"}, {30, "meteor_strike"},
		},
		"kijin": {
			{1, "meteor_strike"}, {20, "world_severing"}, {50, "flame_lance"},
		},
	}

	if list, ok := skills[s.Species]; ok {
		for _, ss := range list {
			if s.Level >= ss.lvl {
				owned := false
				for _, sk := range s.Skills { if sk == ss.id { owned = true; break } }
				if !owned { 
					s.Skills = append(s.Skills, ss.id)
					p.WorldNotice(fmt.Sprintf("LEGACY ACQUIRED: %s has mastered the art of %s.", s.Name, GlobalSkills[ss.id].Name)) 
				}
			}
		}
	}
}

func (p *Player) CheckSubordinateEvolution(s *Subordinate) {
	oldSp := s.Species
	if s.Species == "hobgoblin" && s.Level >= 10 { s.Species = "ogre" } else if s.Species == "ogre" && s.Level >= 25 { s.Species = "kijin" } else if s.Species == "alpha wolf" && s.Level >= 15 { s.Species = "tempest wolf" }
	if oldSp != s.Species { p.WorldNotice(fmt.Sprintf("ASCENSION: %s has evolved from %s into a %s!", s.Name, oldSp, strings.ToUpper(s.Species))) }
}
