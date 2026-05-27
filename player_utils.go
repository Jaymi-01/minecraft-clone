package main

import (
	"fmt"
	"strings"
)

func (p *Player) EquipItem(id string) {
	id = strings.ToLower(id); r, ok := Recipes[id]; if !ok || p.Inventory[id] <= 0 { return }
	if p.Level < r.RequiredLevel { fmt.Printf("🚫 Need Level %d.\n", r.RequiredLevel); return }
	if r.ResultType == "weapon" { p.EquippedWeapon = id } else if r.ResultType == "armor" { p.EquippedArmor = id }
	p.WorldNotice("Equipped: " + r.Name); p.Save()
}

func (p *Player) UnequipItem(slot string) {
	if slot == "weapon" { p.EquippedWeapon = "" } else if slot == "armor" { p.EquippedArmor = "" }; p.Save()
}

func (p *Player) GetEquippedWeaponDamage() int {
	if p.EquippedWeapon == "" { return 0 }; r, ok := Recipes[p.EquippedWeapon]; if ok { return r.ResultValue }; return 0
}

func (p *Player) GetEquippedArmorDefense() int {
	if p.EquippedArmor == "" { return 0 }; r, ok := Recipes[p.EquippedArmor]; if ok { return r.ResultValue }; return 0
}

func (p *Player) GainTaboo(amount int) {
	p.Taboo += amount; p.WorldNotice(fmt.Sprintf("Taboo Level increased to %d", p.Taboo))
	if p.Taboo == 10 { p.WorldNotice("Forbidden threshold crossed.") }; p.CheckTitles()
}

func (p *Player) CheckTitles() {
	for id, t := range GlobalTitles {
		if p.Kills >= t.KillsNeeded {
			met := false; if t.KillsNeeded > 0 { met = true }
			switch id {
			case "taboo_master": if p.Taboo >= 10 { met = true }
			case "taboo_priest": if p.Taboo >= 30 { met = true }
			case "taboo_prophet": if p.Taboo >= 50 { met = true }
			case "labyrinth_walker": if p.ExplorationDepth >= 100 { met = true }
			}
			if met {
				found := false; for _, owned := range p.Titles { if owned == id { found = true; break } }
				if !found { p.Titles = append(p.Titles, id); p.SyncStats(); p.WorldNotice("New Title: " + t.Name) }
			}
		}
	}
}

func (p *Player) ShowHelp() {
	fmt.Println("\n--- 📖 SYSTEM GUIDE ---")
	fmt.Println("   !mine <loc>   - Gather resources")
	fmt.Println("   !status / !s  - View profile")
	fmt.Println("   !inventory / !i- Check items")
	fmt.Println("   !shop / !buy  - Purchase gear")
	fmt.Println("   !elementalshop- Master elements")
	fmt.Println("   !tabooshop    - Forbidden sanctum")
	fmt.Println("   !merge <attr> <skill> - Evolve skills")
	fmt.Println("   !squad add <n>- Manage combat party")
}

func (p *Player) Use(id string) { if p.Inventory[id] > 0 { p.Inventory[id]--; p.HealFull() } }
