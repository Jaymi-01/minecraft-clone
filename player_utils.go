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
	if p.EquippedWeapon == "" { return 0 }
	r, ok := Recipes[p.EquippedWeapon]
	if !ok { return 0 }
	
	base := r.ResultValue
	lvl := p.ItemLevels[p.EquippedWeapon]
	// Each enhancement level adds 20% of base value
	boost := int(float64(base) * (float64(lvl) * 0.2))
	return base + boost
}

func (p *Player) GetEquippedArmorDefense() int {
	if p.EquippedArmor == "" { return 0 }
	r, ok := Recipes[p.EquippedArmor]
	if !ok { return 0 }
	
	base := r.ResultValue
	lvl := p.ItemLevels[p.EquippedArmor]
	boost := int(float64(base) * (float64(lvl) * 0.2))
	return base + boost
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
	
	fmt.Println("\n   --- CORE & IDENTIFICATION ---")
	fmt.Println("   !status / !s      - Display comprehensive identification data.")
	fmt.Println("   !inventory / !i   - Access dimensional storage.")
	fmt.Println("   !skills / !sk     - View your active syncs and unlocked skills.")
	fmt.Println("   !allskills        - Pull the complete System skill database.")
	fmt.Println("   !subordinates     - Inspect your recruited Shadow Army.")
	fmt.Println("   !titles           - List achieved accolades and their buffs.")
	fmt.Println("   !quests / !q      - Check current mission objectives.")

	fmt.Println("\n   --- WORLD & PROGRESSION ---")
	fmt.Println("   !mine <location>  - Extract resources from mapped sectors.")
	fmt.Println("   !train <type>     - Perform daily physical training (Pushups/Situps/Squats/Running).")
	fmt.Println("   !explore          - Enter the Great Labyrinth (requires Level 10).")
	fmt.Println("   !enter            - Breach the current Gate manifestation.")
	fmt.Println("   !arise            - Attempt to extract a Shadow from a defeated Gate Boss.")
	fmt.Println("   !origin <type>    - Finalize reincarnation logic (Slime or Spider).")
	fmt.Println("   !evolve           - Transcend to the next existence tier.")
	fmt.Println("   !jobtrial         - Challenge your shadow reflection (Level 40).")

	fmt.Println("\n   --- EQUIPMENT & MASTERY ---")
	fmt.Println("   !equip <id>       - Synchronize gear or skill to active slots.")
	fmt.Println("   !unequip <slot>   - De-synchronize gear or skill (ID or Slot #).")
	fmt.Println("   !use <item_id>    - Consume material for restoration or extraction.")
	fmt.Println("   !upgrade <sk_id>  - Enhance skill level using Skill Points.")
	fmt.Println("   !enhance <it_id>  - Force-upgrade weapon/armor at the forge (+1 to +10).")
	fmt.Println("   !socket <it> <rn>- Insert a combat rune into gear slots.")
	fmt.Println("   !dupskill <n> <id>- Harvest Lord: Duplicate skill from subordinate.")
	fmt.Println("   !create <s1> <s2> - Harvest Lord: Synthesize a new skill from raw mana.")
	fmt.Println("   !learn <skill_id> - Integrate new skill data from the archive.")
	fmt.Println("   !craft <item_id>  - Synthesize gear at the System forge.")
	fmt.Println("   !merge <attr> <sk>- Fuse forbidden attributes with mastered skills.")

	fmt.Println("\n   --- ECONOMY & DOMAIN ---")
	fmt.Println("   !shop             - Access the static merchant network.")
	fmt.Println("   !buy <item_id>    - Acquire standard instruments or gear.")
	fmt.Println("   !elementalshop    - Contract with fundamental elements.")
	fmt.Println("   !buyelemental <id>- Master a specific elemental logic.")
	fmt.Println("   !tabooshop        - Access the Forbidden Sanctum (requires Taboo).")
	fmt.Println("   !buytaboo <id>    - Master a forbidden attribute or ultimate skill.")
	fmt.Println("   !caravan          - Access the Traveling Merchant Caravan (Rare treasures).")
	fmt.Println("   !buycaravan <id>  - Secure rare treasures before the merchants depart.")
	fmt.Println("   !raid <target>    - Commence full-scale assault on hostile settlements.")
	fmt.Println("   !domain           - View Jura Tempest Federation status and structures.")
	fmt.Println("   !domain build <id>- Construct specialized dominion structures.")
	fmt.Println("   !domain claim     - Collect passive income from the Federation Treasury.")
	fmt.Println("   !squad <list|add|remove> - Manage your active combat vanguard.")
	fmt.Println("   !shadowexchange   - Monarch: Swap squad members in real-time.")
	fmt.Println("   !name <sp> <name> - Bestow a name upon a creature (Baptism).")
	fmt.Println("   !promote <shadow> - Ascend a Shadow's rank using Demon Souls.")
	
	fmt.Println("\n   --- SYSTEM ---")
	fmt.Println("   !help / !h        - Re-display this System Command Directory.")
	fmt.Println("   !exit             - Terminate system link and save progress.")
}

func (p *Player) Use(id string) { 
	id = strings.ToLower(id)
	if p.Inventory[id] <= 0 { fmt.Printf("❌ [SYSTEM]: Item '%s' not found in storage.\n", id); return }
	
	if id == "hidden_box" {
		p.OpenHiddenBox()
		return
	}

	p.Inventory[id]--
	if p.Inventory[id] == 0 { delete(p.Inventory, id) }
	
	p.HealFull()
	p.WorldNotice(fmt.Sprintf("RESTORED: Existence stabilized via %s.", strings.ToUpper(id)))
	p.Save()
}
