package main

import (
	"fmt"
	"strings"
	"time"
)

func (p *Player) ListDomain() {
	if p.DomainLevel == 0 { p.DomainLevel = 1 }
	if p.DomainStructures == nil { p.DomainStructures = make(map[string]int) }

	fmt.Println("\n--- 🏯 [DOMINION: JURA TEMPEST FEDERATION] ---")
	fmt.Printf("   FEDERATION LEVEL: %d\n", p.DomainLevel)
	fmt.Printf("   CITIZEN COUNT: %d Subordinates\n", len(p.Subordinates))
	
	fmt.Println("\n   --- STRUCTURES ---")
	if len(p.DomainStructures) == 0 {
		fmt.Println("      (No specialized structures established)")
	} else {
		for id, lvl := range p.DomainStructures {
			name := strings.Title(strings.Replace(id, "_", " ", -1))
			fmt.Printf("      - %s (Lv.%d)\n", name, lvl)
		}
	}

	fmt.Println("\n   --- AVAILABLE FOR CONSTRUCTION ---")
	fmt.Println("      [research_lab]   - Increases Skill XP gain (+10%/Lv)")
	fmt.Println("      [training_ground]- Increases Subordinate XP gain (+10%/Lv)")
	fmt.Println("      [treasury]       - Generates passive income (💰 100/Lv/Hr)")
	
	fmt.Println("\n💡 [TIP]: Use '!domain build <id>' or '!domain claim'.")
}

func (p *Player) BuildDomainStructure(id string) {
	id = strings.ToLower(id)
	valid := map[string]bool{"research_lab": true, "training_ground": true, "treasury": true}
	if !valid[id] { fmt.Println("❌ [SYSTEM]: Invalid blueprint ID."); return }

	currentLvl := p.DomainStructures[id]
	cost := (currentLvl + 1) * 2000
	
	if p.Inventory["gold"] < cost {
		fmt.Printf("📦 [SYSTEM]: Construction failed. Need %d more Gold.\n", cost-p.Inventory["gold"])
		return
	}

	p.Inventory["gold"] -= cost
	p.DomainStructures[id]++
	p.WorldNotice(fmt.Sprintf("CONSTRUCTION COMPLETE: %s is now Level %d.", strings.ToUpper(id), p.DomainStructures[id]))
	p.Save()
}

func (p *Player) ClaimPassiveIncome() {
	if p.LastIncomeClaim.IsZero() { p.LastIncomeClaim = time.Now() }
	
	hours := int(time.Since(p.LastIncomeClaim).Hours())
	if hours < 1 {
		fmt.Println("❌ [SYSTEM]: Treasury reserves are currently accumulating. Wait at least 1 hour.")
		return
	}

	treasuryLvl := p.DomainStructures["treasury"]
	goldGained := treasuryLvl * 100 * hours
	
	if goldGained > 0 {
		p.Inventory["gold"] += goldGained
		p.WorldNotice(fmt.Sprintf("INCOME SECURED: Collected 💰 %d Gold from the Federation Treasury.", goldGained))
	} else {
		fmt.Println("⚪ [SYSTEM]: Treasury is empty. Established level 1+ required.")
	}

	p.LastIncomeClaim = time.Now()
	p.Save()
}
