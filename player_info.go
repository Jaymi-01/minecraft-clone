package main

import (
	"fmt"
	"strings"
)

func (p *Player) ShowStats() {
	fmt.Printf("\n--- 👤 STATUS ---\n   NAME: %s\n   ORIGIN: %s\n   MINE: [%s] Lv%d\n   HUNTER: [%s] Lv%d\n   COMBAT: ATK %d | DEF %d\n   RECORDS: KILLS %d | TABOO %d\n", p.Name, p.SystemOrigin, p.Rank, p.Level, p.HunterRank, p.HunterLevel, p.Attack, p.Defense, p.Kills, p.Taboo)
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

	fmt.Printf("\n[🧠 ANALYSIS]: '%s' identified.", strings.ToUpper(itemID))
	if lvl >= 4 { fmt.Printf(" | Rarity: %s", rarity) }
	if lvl >= 7 { fmt.Printf(" | Value: 💰 %d", price) }
	if lvl >= 10 {
		desc := "Standard material."
		if r, ok := Recipes[itemID]; ok { desc = fmt.Sprintf("Used to craft %s.", r.Name) }
		fmt.Printf("\n   📜 %s", desc)
	}
	fmt.Println()
}

func (p *Player) ShowInventory() {
	fmt.Println("\n--- 🎒 Inventory ---")
	for k, v := range p.Inventory { if v > 0 { fmt.Printf("   %s: %d\n", k, v) } }
}

func (p *Player) ListSubordinates() {
	fmt.Println("\n--- 🤝 Subordinates ---")
	for _, s := range p.Subordinates { fmt.Printf("   🐾 %s [%s] - LV.%d | ATK: %d\n", s.Name, s.Species, s.Level, s.Attack) }
}

func (p *Player) AddToSquad(name string) {
	if len(p.Squad) >= 3 { fmt.Println("❌ Squad full."); return }
	for _, s := range p.Subordinates {
		if strings.EqualFold(s.Name, name) {
			for _, n := range p.Squad { if strings.EqualFold(n, name) { return } }
			p.Squad = append(p.Squad, s.Name); p.WorldNotice(s.Name + " joined squad."); p.Save(); return
		}
	}
}

func (p *Player) ListSquad() {
	fmt.Println("\n--- 👥 Combat Squad ---")
	for _, n := range p.Squad { fmt.Printf("   🤝 %s\n", n) }
}

func (p *Player) ListTitles() {
	fmt.Println("\n--- 🏆 Titles ---")
	for _, tid := range p.Titles { fmt.Printf("   🏅 %s\n", GlobalTitles[tid].Name) }
}

func (p *Player) NameSubordinate(species, givenName string) {
	species = strings.ToLower(strings.Replace(species, "_", " ", -1))
	if p.MaxMagic < 50 { fmt.Println("❌ Need 50 Max MP."); return }
	p.MaxMagic -= 50; if p.Magic > p.MaxMagic { p.Magic = p.MaxMagic }
	sub := Subordinate{Name: givenName, Species: species, Attack: 20, Defense: 10, Level: 1}
	p.Subordinates = append(p.Subordinates, sub); p.WorldNotice("NAMED: " + givenName); p.Save()
}
