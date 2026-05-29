package main

import (
	"fmt"
	"math/rand"
	"strings"
)

func (p *Player) ChooseOrigin(origin string) {
	origin = strings.ToLower(origin)
	if p.SystemOrigin != "Human" { fmt.Printf("❌ [SYSTEM]: Origin fixed as %s. Re-selection impossible.\n", p.SystemOrigin); return }
	if p.Level < 10 { fmt.Println("🚫 [SYSTEM]: Reincarnation logic requires Level 10 energy signature."); return }
	
	switch origin {
	case "slime":
		p.SystemOrigin = "Slime"
		p.AddSkill("fire_bolt"); p.AddSkill("predator"); p.AddSkill("great_sage")
		p.WorldNotice("ORIGIN SELECTED: You have reincarnated as a Slime. Predator logic initialized.")
	case "spider":
		p.SystemOrigin = "Spider"
		p.AddSkill("venom_spit"); p.AddSkill("appraisal"); p.AddSkill("spider_thread")
		p.WorldNotice("ORIGIN SELECTED: You have reincarnated as a Small Lesser Taratect. Survival logic initialized.")
	default:
		fmt.Println("❌ [SYSTEM]: Invalid origin specified. Choose 'slime' or 'spider'.")
		return
	}
	p.Save()
}

func (p *Player) Evolve() {
	evolved := false
	oldOrigin := p.SystemOrigin
	switch p.SystemOrigin {
	case "Slime":
		if p.Level >= 30 { p.SystemOrigin = "Demon Slime"; p.AddSkill("raphael"); evolved = true }
	case "Demon Slime":
		if p.Level >= 60 { p.SystemOrigin = "Ultimate Slime (True Dragon)"; evolved = true }
	case "Spider":
		if p.Level >= 20 { p.SystemOrigin = "Small Poison Taratect"; evolved = true }
	case "Small Poison Taratect":
		if p.Level >= 30 { p.SystemOrigin = "Arachne"; p.AddSkill("dim_maneuver"); evolved = true }
	case "Arachne":
		if p.Level >= 60 { p.SystemOrigin = "God (Shiraori)"; p.AddSkill("egg_revival"); evolved = true }
	}
	
	if evolved {
		p.SyncStats(); p.HealFull()
		p.WorldNotice(fmt.Sprintf("EVOLUTION SUCCESSFUL: %s has evolved into %s!", oldOrigin, p.SystemOrigin))
		fmt.Println("🔥 [SYSTEM]: Physical and Magical limits have been recalculated.")
	} else {
		fmt.Println("🚫 [SYSTEM]: Evolution requirements not met. Higher level or specific achievements required.")
	}
}

func (p *Player) StartExploration() {
	if p.Level < 10 { fmt.Println("🚫 [SYSTEM]: Level 10 required to enter the Labyrinth."); return }
	p.Exploring = true; p.ExplorationDepth = 1
	p.WorldNotice("LABYRINTH ENTERED: Reality has shifted. Survival is the only objective.")
}

func (p *Player) Emerge() {
	p.Exploring = false
	p.WorldNotice("LABYRINTH EXITED: You have returned to the surface world.")
}

func (p *Player) Move(dir string) {
	if !p.Exploring || p.Stamina < 2 { return }
	p.Stamina -= 2; p.ExplorationDepth++
	dir = strings.ToUpper(dir)
	fmt.Printf("\n👣 [LABYRINTH]: Moving %s... (Current Depth: %d)\n", dir, p.ExplorationDepth)
	
	hasAnalysis := p.HasSkill("appraisal") || p.HasSkill("great_sage") || p.HasSkill("sariel") || p.HasSkill("raphael") || p.HasSkill("ciel")
	if hasAnalysis {
		if p.HasSkill("sariel") { p.SkillUsage["sariel"]++; if p.SkillUsage["sariel"] >= 10 { p.UpgradeSkill("sariel", true) } } else if p.HasSkill("appraisal") { p.SkillUsage["appraisal"]++; if p.SkillUsage["appraisal"] >= 10 { p.UpgradeSkill("appraisal", true) } }
	}

	e := rand.Intn(100)
	if e < 15 { p.FoundChest(hasAnalysis) } else if e < 30 { p.TriggerTrap(hasAnalysis) } else if e < 55 { p.EncounterMonster() }
}
