package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

func (p *Player) Arise() {
	if !p.AriseActive || p.AriseMonster == nil {
		fmt.Println("❌ [SYSTEM]: No extractable mana signature detected.")
		return
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("\n🌌 [SHADOW EXTRACTION]: Defeated entity [%s] detected.\n", p.AriseMonster.Name)
	fmt.Println("⚠️ [WARNING]: You have 3 attempts to extract the shadow. Failure will dissipate the mana.")
	
	for i := 1; i <= 3; i++ {
		fmt.Printf("\n--- [EXTRACTION ATTEMPT %d/3] ---\n", i)
		fmt.Print("Type 'ARISE' to commence: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToUpper(input))

		if input != "ARISE" {
			fmt.Println("🚫 Extraction sequence aborted.")
			continue
		}

		// Success Chance: 30% base + (Magic / 100)% - (Monster Health / 500)%
		chance := 0.3 + (float64(p.MaxMagic) / 5000.0)
		if rand.Float64() <= chance {
			p.SuccessArise()
			return
		}

		fmt.Println("❌ [FAILURE]: The shadow has resisted the call.")
	}

	fmt.Println("💨 [DISSIPATED]: The mana signature has vanished into the abyss.")
	p.AriseActive = false
	p.AriseMonster = nil
	p.Save()
}

func (p *Player) SuccessArise() {
	p.WorldNotice("EXTRACTION SUCCESSFUL: [ARISE]")
	
	shadowName := strings.Replace(p.AriseMonster.Name, "🕷️ ", "", 1)
	shadowName = strings.Replace(shadowName, "👹 ", "", 1)
	shadowName = strings.Replace(shadowName, "🛡️ ", "", 1)
	shadowName = strings.Replace(shadowName, "👑 ", "", 1)
	shadowName = strings.Replace(shadowName, "💀 ", "", 1)
	shadowName = strings.TrimSpace(shadowName)

	newShadow := Subordinate{
		Name:    "Shadow " + shadowName,
		Species: "Shadow",
		Rank:    "Soldier",
		Attack:  p.AriseMonster.Damage / 2,
		Defense: 50,
		Level:   p.Level / 2,
		XP:      0,
		NextXP:  100,
	}
	if newShadow.Level < 1 { newShadow.Level = 1 }
	if newShadow.Attack < 20 { newShadow.Attack = 20 }

	p.Subordinates = append(p.Subordinates, newShadow)
	p.AriseActive = false
	p.AriseMonster = nil
	
	p.WorldNotice(fmt.Sprintf("NEW SHADOW RECRUITED: [%s] joined the Shadow Army.", newShadow.Name))
	p.Save()
}
