package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func (p *Player) StartJobTrial() {
	if p.Level < 40 {
		fmt.Printf("🚫 [SYSTEM]: Job Advancement requires Level 40. Current Level: %d\n", p.Level)
		return
	}
	if p.Job != "" {
		fmt.Printf("❌ [SYSTEM]: You have already obtained the Job: %s\n", p.Job)
		return
	}

	p.WorldNotice("JOB ADVANCEMENT QUEST: THE FINAL TRIAL.")
	fmt.Println("⚠️ [WARNING]: This is a solo trial. Subordinate support is blocked by the System.")
	
	// Create a reflection boss
	boss := Monster{
		Name:   "👤 Shadow Reflection of " + p.Name,
		Health: p.MaxHealth * 2,
		Damage: p.Attack,
		LootTable: map[string]float64{"job_token": 1.0},
	}

	p.TrialActive = true
	// Force solo combat by temporarily clearing squad or handling it in combat logic
	// For simplicity, we just trigger combat and let the user know they are solo
	if p.Combat(&boss, false) {
		p.SelectJob()
	}
	p.TrialActive = false
	p.Save()
}

func (p *Player) SelectJob() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("\n🌟 [QUEST COMPLETED]: You have defeated your reflection.")
	fmt.Println("📜 [SYSTEM]: Choose your permanent path of power:")
	fmt.Println("   1. [SHADOW MONARCH] - Master of the Shadow Army. Unlocks !shadowexchange.")
	fmt.Println("   2. [TRUE DEMON LORD] - Master of Chaos. Unlocks Passive: Infinite Regeneration.")
	fmt.Println("   3. [ABYSSAL BEING] - Master of Space. Unlocks Passive: Spatial Dominion.")
	
	for {
		fmt.Print("\nSELECT PATH (1-3): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		
		switch input {
		case "1":
			p.Job = "Shadow Monarch"
			p.Attack += 100
			p.MaxMagic += 500
			p.WorldNotice("CLASS ADVANCEMENT: YOU ARE THE SHADOW MONARCH.")
			return
		case "2":
			p.Job = "True Demon Lord"
			p.MaxHealth += 2000
			p.MaxMagic += 1000
			p.WorldNotice("CLASS ADVANCEMENT: YOU ARE THE TRUE DEMON LORD.")
			return
		case "3":
			p.Job = "Abyssal Being"
			p.Attack += 200
			p.Defense += 100
			p.WorldNotice("CLASS ADVANCEMENT: YOU ARE THE ABYSSAL BEING.")
			return
		default:
			fmt.Println("❌ Invalid choice.")
		}
	}
}

func (p *Player) ShadowExchange(current, replacement string) {
	if p.Job != "Shadow Monarch" { return }
	
	foundIdx := -1
	for i, n := range p.Squad {
		if strings.EqualFold(n, current) {
			foundIdx = i
			break
		}
	}

	if foundIdx == -1 {
		fmt.Printf("❌ [SYSTEM]: '%s' is not in the active squad.\n", current)
		return
	}

	subExists := false
	for _, s := range p.Subordinates {
		if strings.EqualFold(s.Name, replacement) {
			subExists = true
			break
		}
	}

	if !subExists {
		fmt.Printf("❌ [SYSTEM]: '%s' not found in your Shadow Army.\n", replacement)
		return
	}

	// Swap
	p.Squad[foundIdx] = replacement
	p.WorldNotice(fmt.Sprintf("SHADOW EXCHANGE: [%s] has been swapped with [%s].", current, replacement))
	p.Save()
}

func (p *Player) ApplyJobPassives() {
	if p.Job == "True Demon Lord" {
		// Passive: Infinite Regeneration (5% Max HP every turn in combat handled in combat.go)
		// Standard regen handled in economy.go
	}
	if p.Job == "Abyssal Being" {
		// Passive: Spatial Dominion (Ignore 25% monster defense handled in combat.go)
	}
}
