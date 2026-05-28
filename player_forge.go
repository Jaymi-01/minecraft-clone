package main

import (
	"fmt"
	"math/rand"
	"strings"
)

func (p *Player) EnhanceItem(itemID string) {
	itemID = strings.ToLower(itemID)
	if p.Inventory[itemID] <= 0 { fmt.Printf("❌ [SYSTEM]: Item '%s' not found in storage.\n", itemID); return }
	
	r, ok := Recipes[itemID]
	if !ok || (r.ResultType != "weapon" && r.ResultType != "armor") {
		fmt.Println("❌ [SYSTEM]: Only Weapons and Armor can be enhanced via the Forge.")
		return
	}

	currentLvl := p.ItemLevels[itemID]
	if currentLvl >= 10 {
		fmt.Println("❌ [SYSTEM]: This item has already reached the maximum enhancement tier (+10).")
		return
	}

	nextLvl := currentLvl + 1
	reqMat := ""
	reqQty := 0

	if nextLvl <= 3 { reqMat = "void_essence"; reqQty = 1 }
	if nextLvl >= 4 && nextLvl <= 7 { reqMat = "star_matter"; reqQty = 1 }
	if nextLvl >= 8 { reqMat = "star_matter"; reqQty = 3 }

	if p.Inventory[reqMat] < reqQty {
		fmt.Printf("📦 [SYSTEM]: Enhancement failed. Need %d more %s.\n", reqQty-p.Inventory[reqMat], reqMat)
		return
	}

	// Chance of success
	chance := 1.0 - (float64(currentLvl) * 0.08)
	p.Inventory[reqMat] -= reqQty
	if p.Inventory[reqMat] == 0 { delete(p.Inventory, reqMat) }

	if rand.Float64() <= chance {
		p.ItemLevels[itemID]++
		p.WorldNotice(fmt.Sprintf("ENHANCEMENT SUCCESS: [%s] is now +%d!", r.Name, p.ItemLevels[itemID]))
		p.SyncStats()
	} else {
		fmt.Printf("💥 [FAILURE]: The enhancement sequence failed. The %s was consumed.\n", reqMat)
	}
	p.Save()
}

func (p *Player) SocketRune(itemID, runeID string) {
	itemID = strings.ToLower(itemID)
	runeID = strings.ToLower(runeID)
	
	if p.Inventory[itemID] <= 0 || p.Inventory[runeID] <= 0 {
		fmt.Println("❌ [SYSTEM]: Item or Rune not found in storage.")
		return
	}

	r, ok := Recipes[itemID]
	if !ok || (r.ResultType != "weapon" && r.ResultType != "armor") {
		fmt.Println("❌ [SYSTEM]: Only combat gear supports socketing.")
		return
	}

	// Max 2 runes per item
	if len(p.ItemRunes[itemID]) >= 2 {
		fmt.Println("❌ [SYSTEM]: All socket slots are currently occupied.")
		return
	}

	validRunes := map[string]string{
		"lifesteal_rune": "Lifesteal (5% DMG -> HP)",
		"mana_rune":      "Mana Flux (+50 Max MP)",
		"defense_rune":   "Iron Skin (+20 Defense)",
	}

	if _, valid := validRunes[runeID]; !valid {
		fmt.Println("❌ [SYSTEM]: The specified material is not a valid combat rune.")
		return
	}

	p.Inventory[runeID]--
	if p.Inventory[runeID] == 0 { delete(p.Inventory, runeID) }

	if p.ItemRunes == nil { p.ItemRunes = make(map[string][]string) }
	p.ItemRunes[itemID] = append(p.ItemRunes[itemID], runeID)
	
	p.WorldNotice(fmt.Sprintf("RUNE SYNCHRONIZATION: [%s] has been socketed into [%s].", strings.ToUpper(runeID), r.Name))
	p.SyncStats()
	p.Save()
}
