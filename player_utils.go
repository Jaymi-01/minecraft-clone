package main

import (
	"fmt"
	"strings"
)

func (p *Player) EquipItem(id string) {
	id = strings.ToLower(id); r, ok := Recipes[id]; if !ok || p.Inventory[id] <= 0 { return }
	if p.Level < r.RequiredLevel { fmt.Printf("🚫 [SYSTEM]: Strength insufficient. Level %d required.\n", r.RequiredLevel); return }
	if r.ResultType == "weapon" { p.EquippedWeapon = id } else if r.ResultType == "armor" { p.EquippedArmor = id }
	p.WorldNotice("EQUIPPED: " + r.Name); p.Save()
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
	p.Taboo += amount; p.WorldNotice(fmt.Sprintf("TABOO INFLUENCE increased to %d", p.Taboo))
	if p.Taboo >= 10 { p.WorldNotice("FORBIDDEN THRESHOLD CROSSED: The System is watching.") }; p.CheckTitles()
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
				if !found { p.Titles = append(p.Titles, id); p.SyncStats(); p.WorldNotice("TITLE GRANTED: " + t.Name) }
			}
		}
	}
}

func (p *Player) ShowHelp() {
	fmt.Println("\n--- 📖 [SYSTEM COMMAND DIRECTORY] ---")
	fmt.Println("   --- CORE ---")
	fmt.Println("   !status / !s      - Display comprehensive identification data.")
	fmt.Println("   !inventory / !i   - Access dimensional storage.")
	fmt.Println("   !skills / !sk     - View integrated skill archive.")
	fmt.Println("   !titles           - List achieved accolades.")
	fmt.Println("   !quests / !q      - Check current mission objectives.")
	
	fmt.Println("\n   --- PROGRESSION ---")
	fmt.Println("   !mine <location>  - Extract resources from mapped sectors.")
	fmt.Println("   !explore          - Enter the Great Labyrinth.")
	fmt.Println("   !enter            - Breach the current Gate manifestation.")
	fmt.Println("   !origin <type>    - Finalize reincarnation logic (Lvl 10).")
	fmt.Println("   !evolve           - Transcend to the next existence tier.")
	
	fmt.Println("\n   --- EQUIPMENT & MASTERY ---")
	fmt.Println("   !equip <id>       - Synchronize gear or skill to active slots.")
	fmt.Println("   !unequip <slot>   - De-synchronize gear or skill.")
	fmt.Println("   !use <item_id>    - Consume material for restoration.")
	fmt.Println("   !upgrade <sk_id>  - Enhance skill level using Skill Points.")
	fmt.Println("   !dupskill <n> <id>- Harvest and duplicate skill from subordinate.")
	fmt.Println("   !create <s1> <s2> - Synthesize a new skill from raw mana.")
	fmt.Println("   !learn <skill_id> - Integrate new skill data.")
	fmt.Println("   !craft <item_id>  - Synthesize gear at the system forge.")
	fmt.Println("   !merge <attr> <sk>- Fuse forbidden attributes with mastered skills.")
	
	fmt.Println("\n   --- ECONOMY & DOMAIN ---")
	fmt.Println("   !shop             - Access merchant network.")
	fmt.Println("   !buy <id>         - Acquire standard instruments.")
	fmt.Println("   !elementalshop    - Contract with fundamental elements.")
	fmt.Println("   !tabooshop        - Access the Forbidden Sanctum.")
	fmt.Println("   !raid <target>    - Commence assault on hostile settlements.")
	fmt.Println("   !squad <list|add|remove> - Manage the Shadow Army vanguard.")
	fmt.Println("   !name <sp> <n>    - Bestow a name upon a creature (Baptism).")
	fmt.Println("   !help             - Re-display this directory.")
}

func (p *Player) Use(id string) { 
	id = strings.ToLower(id)
	if p.Inventory[id] <= 0 { fmt.Printf("❌ [SYSTEM]: Item '%s' not found in storage.\n", id); return }
	
	p.Inventory[id]--
	if p.Inventory[id] == 0 { delete(p.Inventory, id) }
	
	p.HealFull()
	p.WorldNotice(fmt.Sprintf("RESTORED: Existence stabilized via %s.", strings.ToUpper(id)))
	p.Save()
}
