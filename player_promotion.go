package main

import (
	"fmt"
	"strings"
)

var ShadowRanks = []string{"Soldier", "Elite", "Knight", "Elite Knight", "Commander", "Grand Marshal"}

func (p *Player) PromoteShadow(nameInput string) {
	foundIdx := -1
	nameInput = strings.ToLower(nameInput)

	// Phase 1: Try exact match (case-insensitive)
	for i := range p.Subordinates {
		if strings.EqualFold(p.Subordinates[i].Name, nameInput) {
			foundIdx = i; break
		}
	}

	// Phase 2: Try partial match (if no exact match found)
	if foundIdx == -1 {
		for i := range p.Subordinates {
			// Check if input is in name OR if name (without 'Shadow ') is in input
			subName := strings.ToLower(p.Subordinates[i].Name)
			cleanSubName := strings.Replace(subName, "shadow ", "", 1)
			if strings.Contains(subName, nameInput) || strings.Contains(nameInput, cleanSubName) {
				foundIdx = i; break
			}
		}
	}

	if foundIdx == -1 {
		fmt.Printf("❌ [SYSTEM]: Shadow '%s' not found in your army. Use !subordinates to check names.\n", nameInput)
		return
	}

	s := &p.Subordinates[foundIdx]
	if s.Species != "Shadow" {
		fmt.Println("❌ [SYSTEM]: Only entities of the Shadow Species can undergo Rank Ascension.")
		return
	}

	currentRankIdx := -1
	for i, r := range ShadowRanks {
		if s.Rank == r {
			currentRankIdx = i
			break
		}
	}

	if currentRankIdx == len(ShadowRanks)-1 {
		fmt.Println("❌ [SYSTEM]: This shadow has already reached the peak rank: Grand Marshal.")
		return
	}

	nextRank := ShadowRanks[currentRankIdx+1]
	reqLevel := 0
	reqSouls := 0

	switch nextRank {
	case "Elite": reqLevel = 20; reqSouls = 1
	case "Knight": reqLevel = 30; reqSouls = 3
	case "Elite Knight": reqLevel = 45; reqSouls = 7
	case "Commander": reqLevel = 60; reqSouls = 15
	case "Grand Marshal": reqLevel = 80; reqSouls = 30
	}

	if s.Level < reqLevel {
		fmt.Printf("🚫 [SYSTEM]: Ascension denied. %s requires Level %d (Current: %d).\n", nextRank, reqLevel, s.Level)
		return
	}

	if p.Inventory["demon_soul"] < reqSouls {
		fmt.Printf("📦 [SYSTEM]: Insufficient materials. Need %d more Demon Souls.\n", reqSouls-p.Inventory["demon_soul"])
		return
	}

	// Perform Ascension
	p.Inventory["demon_soul"] -= reqSouls
	if p.Inventory["demon_soul"] == 0 { delete(p.Inventory, "demon_soul") }

	oldRank := s.Rank
	s.Rank = nextRank
	s.Attack = int(float64(s.Attack) * 1.5)
	s.Defense = int(float64(s.Defense) * 1.5)
	
	p.WorldNotice(fmt.Sprintf("RANK ASCENSION: [%s] has evolved from %s to %s!", s.Name, oldRank, nextRank))
	fmt.Printf("🔥 [SYSTEM]: %s's combat specs have been significantly enhanced.\n", s.Name)
	p.Save()
}
