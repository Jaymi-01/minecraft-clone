package main

import (
	"fmt"
	"strings"
	"time"
)

func (p *Player) ShowStats() {
	fmt.Printf("\n--- 👤 [STATUS IDENTIFICATION] ---\n")
	fmt.Printf("   NAME: %s\n", p.Name)
	fmt.Printf("   ORIGIN: %s\n", p.SystemOrigin)
	fmt.Printf("   MINE RANK: [%s] Lv%d\n", p.Rank, p.Level)
	fmt.Printf("   HUNTER RANK: [%s] Lv%d\n", p.HunterRank, p.HunterLevel)
	fmt.Printf("   COMBAT SPECS: ATK %d | DEF %d\n", p.Attack, p.Defense)
	fmt.Printf("   SYSTEM RECORDS: KILLS %d | TABOO %d\n", p.Kills, p.Taboo)
}

func (p *Player) AutoAnalyze(itemID string) {
	lvl := 1
	if l, ok := p.SkillLevels["great_sage"]; ok { lvl = l }
	if l, ok := p.SkillLevels["raphael"]; ok { lvl = l + 10 }
	if l, ok := p.SkillLevels["sariel"]; ok { lvl = l + 10 }

	rarity := "Common"; price := 5
	switch itemID {
	case "diamond", "abyss_crystal", "netherite": rarity = "Epic"; price = 100
	case "void_essence", "star_matter", "void_core": rarity = "Legendary"; price = 500
	case "void_crown", "life_stone", "demon_soul": rarity = "Mythic"; price = 5000
	case "iron", "gold", "quartz": rarity = "Rare"; price = 25
	}

	fmt.Printf("\n[🧠 ANALYZE]: '%s' identified.", strings.ToUpper(itemID))
	if lvl >= 4 { fmt.Printf(" | RARITY: %s", rarity) }
	if lvl >= 7 { fmt.Printf(" | MARKET VALUE: 💰 %d Gold", price) }
	if lvl >= 10 {
		desc := "Standard material."
		if r, ok := Recipes[itemID]; ok { desc = fmt.Sprintf("Essential for crafting %s.", r.Name) }
		fmt.Printf("\n   📜 DATA: %s", desc)
	}
	fmt.Println()
}

func (p *Player) ShowInventory() {
	fmt.Println("\n--- 🎒 [DIMENSIONAL STORAGE] ---")
	found := false
	for k, v := range p.Inventory { 
		if v > 0 { 
			fmt.Printf("   - %s: %d units\n", strings.Title(strings.Replace(k, "_", " ", -1)), v) 
			found = true
		} 
	}
	if !found { fmt.Println("   (Inventory is empty)") }
}

func (p *Player) ListSubordinates() {
	fmt.Println("\n--- 🤝 [SHADOW ARMY: SUBORDINATES] ---")
	if len(p.Subordinates) == 0 { fmt.Println("   (No subordinates recruited)"); return }
	for _, s := range p.Subordinates { 
		rankText := ""
		if s.Rank != "" { rankText = fmt.Sprintf("[%s] ", strings.ToUpper(s.Rank)) }
		fmt.Printf("   🐾 %s %s[%s] - LV.%d\n", rankText, s.Name, strings.ToUpper(s.Species), s.Level)
		fmt.Printf("      ⚔️ ATK: %d | 🛡️ DEF: %d\n", s.Attack, s.Defense)
		if len(s.Skills) > 0 {
			fmt.Print("      ✨ SKILLS: ")
			for _, skID := range s.Skills { fmt.Printf("[%s] ", GlobalSkills[skID].Name) }
			fmt.Println()
		} else {
			fmt.Println("      ✨ SKILLS: (None mastered)")
		}
	}
}

func (p *Player) AddToSquad(name string) {
	if len(p.Squad) >= 3 { fmt.Println("❌ [SYSTEM]: Squad capacity reached. Maximum 3 units allowed."); return }
	for _, s := range p.Subordinates {
		if strings.EqualFold(s.Name, name) {
			for _, n := range p.Squad { 
				if strings.EqualFold(n, name) { 
					fmt.Println("❌ [SYSTEM]: Unit is already active in the current squad.")
					return 
				} 
			}
			p.Squad = append(p.Squad, s.Name)
			p.WorldNotice(fmt.Sprintf("SHADOW EXTRACTION: [%s] has been integrated into the Combat Squad.", s.Name))
			p.Save()
			return
		}
	}
	fmt.Printf("❌ [SYSTEM]: Subordinate '%s' not found in records.\n", name)
}

func (p *Player) ListSquad() {
	fmt.Println("\n--- 👥 [ACTIVE COMBAT SQUAD] ---")
	if len(p.Squad) == 0 { fmt.Println("   (Squad is currently empty)"); return }
	for _, n := range p.Squad { fmt.Printf("   🤝 UNIT: %s\n", n) }
}

func (p *Player) RemoveFromSquad(name string) {
	for i, n := range p.Squad {
		if strings.EqualFold(n, name) {
			p.Squad = append(p.Squad[:i], p.Squad[i+1:]...)
			p.WorldNotice(fmt.Sprintf("DISMISSAL: [%s] has been removed from active duty.", n))
			p.Save(); return
		}
	}
	fmt.Printf("❌ [SYSTEM]: Unit '%s' is not in the active squad.\n", name)
}

func (p *Player) ListTitles() {
	fmt.Println("\n--- 🏆 [ACHIEVED TITLES] ---")
	if len(p.Titles) == 0 { fmt.Println("   (No titles earned yet)"); return }
	for _, tid := range p.Titles { fmt.Printf("   🏅 %s\n", GlobalTitles[tid].Name) }
}

func (p *Player) NameSubordinate(species, givenName string) {
	speciesID := strings.ToLower(strings.Replace(species, "_", " ", -1))
	if p.MaxMagic < 50 { fmt.Println("❌ [SYSTEM]: Bestowal of name requires 50 Max Mana. Current reserves insufficient."); return }
	
	valid := false
	switch speciesID {
	case "slime", "goblin", "hobgoblin", "wolf", "alpha wolf", "spider", "taratect", "ogre", "kijin": valid = true
	}
	if !valid { fmt.Printf("❌ [SYSTEM]: Entity of species '%s' cannot be named through standard protocols.\n", species); return }

	p.MaxMagic -= 50; if p.Magic > p.MaxMagic { p.Magic = p.MaxMagic }
	sub := Subordinate{
		Name: givenName, 
		Species: speciesID, 
		Attack: 20, 
		Defense: 10, 
		Level: 1,
		LastAction: time.Now(),
	}
	p.Subordinates = append(p.Subordinates, sub)
	p.WorldNotice(fmt.Sprintf("BAPTISM: The name [%s] has been bestowed upon the [%s]. Physical evolution likely.", givenName, strings.ToUpper(speciesID)))
	p.Save()
}
