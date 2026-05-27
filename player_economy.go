package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func (p *Player) ListShop() {
	fmt.Printf("\n--- 💰 [MERCHANT INVENTORY] (Your Gold: %d) ---\n", p.Inventory["gold"])
	for id, it := range MerchantInventory { fmt.Printf("   [%s] %s - 💰 %d Gold\n      📜 %s\n", id, it.Name, it.Price, it.Desc) }
}

func (p *Player) Buy(id string) {
	id = strings.ToLower(id); it, ok := MerchantInventory[id]
	if !ok {
		// Check if it's in other shops to provide a hint
		if _, ok2 := ElementalSkillShopInventory[id]; ok2 {
			fmt.Printf("❌ [SYSTEM]: Mastery of '%s' requires the '!buyelemental' command.\n", id)
			return
		}
		if _, ok2 := TabooShopInventory[id]; ok2 {
			fmt.Printf("❌ [SYSTEM]: Mastery of '%s' requires the '!buytaboo' command.\n", id)
			return
		}
		fmt.Printf("❌ [SYSTEM]: Item '%s' not found in merchant inventory.\n", id); return
	}
	if p.Inventory["gold"] < it.Price { fmt.Printf("❌ [SYSTEM]: Transaction failed. Insufficient gold (Need %d more).\n", it.Price-p.Inventory["gold"]); return }
	
	p.Inventory["gold"] -= it.Price; p.Inventory[it.ID]++
	p.WorldNotice(fmt.Sprintf("PURCHASE SUCCESSFUL: %s acquired.", it.Name))
	p.Save()
}

func (p *Player) ListTabooShop() {
	fmt.Printf("\n--- 🌌 [TABOO SANCTUM] (Forbidden Insight: %d) ---\n", p.Taboo)
	for id, it := range TabooShopInventory { fmt.Printf("   [%s] %s - %d Taboo\n      📜 %s\n", id, it.Name, it.Price, it.Desc) }
}

func (p *Player) BuyTabooSkill(id string) {
	id = strings.ToLower(id); it, ok := TabooShopInventory[id]
	if !ok { fmt.Printf("❌ [SYSTEM]: Forbidden attribute '%s' not found.\n", id); return }
	if p.Taboo < it.Price { fmt.Printf("❌ [SYSTEM]: Insight rejected. Insufficient Taboo level (Need %d more).\n", it.Price-p.Taboo); return }
	
	p.Taboo -= it.Price; if p.Attributes == nil { p.Attributes = make(map[string]bool) }; p.Attributes[id] = true
	p.WorldNotice(fmt.Sprintf("TABOO MASTERED: Your soul has absorbed the %s.", it.Name))
	p.Save()
}

func (p *Player) ListElementalShop() {
	fmt.Printf("\n--- 🔥 [ELEMENTAL MASTERY] (Your Gold: %d) ---\n", p.Inventory["gold"])
	for id, it := range ElementalSkillShopInventory { fmt.Printf("   [%s] %s - 💰 %d Gold\n      📜 %s\n", id, it.Name, it.Price, it.Desc) }
}

func (p *Player) BuyElementalSkill(id string) {
	id = strings.ToLower(id); it, ok := ElementalSkillShopInventory[id]
	if !ok { fmt.Printf("❌ [SYSTEM]: Elemental skill '%s' not found.\n", id); return }
	if p.Inventory["gold"] < it.Price { fmt.Printf("❌ [SYSTEM]: Mastery failed. Insufficient gold (Need %d more).\n", it.Price-p.Inventory["gold"]); return }
	
	p.Inventory["gold"] -= it.Price; p.AddSkill(id)
	p.WorldNotice(fmt.Sprintf("ELEMENTAL CONTRACT: You have mastered %s.", it.Name))
	p.Save()
}

func (p *Player) ListDCraftable() {
	fmt.Println("\n--- 🛠️ [SYSTEM: FORGE & CONSTRUCTION] ---")
	fmt.Println("   --- WEAPONS & ARMOR ---")
	for id, r := range Recipes { fmt.Printf("      [%s] %s (Req. Lv%d)\n", id, r.Name, r.RequiredLevel) }
	fmt.Println("\n   --- STRUCTURES ---")
	for id, s := range Structures { fmt.Printf("      [%s] %s (Req. Lv%d) - %s\n", id, s.Name, s.RequiredLevel, s.PerkDesc) }
}

func (p *Player) ListRaids() {
	fmt.Println("\n--- ⚔️ [THREAT ANALYSIS: SETTLEMENTS] ---")
	for id, t := range BotSettlements { fmt.Printf("   [%s] %s - THREAT LV.%d\n      📜 %s\n", id, t.Name, t.Level, t.Description) }
}

func (p *Player) Craft(itemName string) {
	itemName = strings.ToLower(itemName); r, ok := Recipes[itemName]
	if !ok { fmt.Printf("❌ [SYSTEM]: Recipe for '%s' not found.\n", itemName); return }
	if p.Level < r.RequiredLevel { fmt.Printf("🚫 [SYSTEM]: Crafting failed. Level %d required.\n", r.RequiredLevel); return }
	
	for k, v := range r.Ingredients { if p.Inventory[k] < v { fmt.Printf("📦 [SYSTEM]: Missing materials. Need %d more %s.\n", v-p.Inventory[k], k); return } }
	for k, v := range r.Ingredients { p.Inventory[k] -= v; if p.Inventory[k] == 0 { delete(p.Inventory, k) } }
	
	p.Inventory[itemName]++
	p.WorldNotice(fmt.Sprintf("FORGE SUCCESSFUL: %s has been created.", r.Name))
	p.Save()
}

func (p *Player) Build(id string) {
	s, ok := Structures[strings.ToLower(id)]
	if !ok { fmt.Println("❌ [SYSTEM]: Blueprint not found."); return }
	if p.Level < s.RequiredLevel { fmt.Printf("❌ [SYSTEM]: Level %d required for this construction.\n", s.RequiredLevel); return }
	if p.Structures[id] { fmt.Println("❌ [SYSTEM]: Structure already exists."); return }
	
	for k, v := range s.Ingredients { if p.Inventory[k] < v { fmt.Printf("❌ [SYSTEM]: Construction halted. Need %d more %s.\n", v-p.Inventory[k], k); return } }
	for k, v := range s.Ingredients { p.Inventory[k] -= v; if p.Inventory[k] == 0 { delete(p.Inventory, k) } }
	
	p.Structures[id] = true; p.SyncStats()
	p.WorldNotice(fmt.Sprintf("CONSTRUCTION COMPLETE: The %s has been established.", s.Name))
	p.Save()
}

func (p *Player) StartRegeneration() {
	ticker := time.NewTicker(2 * time.Minute)
	go func() { for range ticker.C { p.Regenerate(); p.Save() } }()
}

func (p *Player) Regenerate() {
	p.Health += 10; p.Stamina += 10; p.Magic += 20
	if p.Health > p.MaxHealth { p.Health = p.MaxHealth }; if p.Stamina > p.MaxStamina { p.Stamina = p.MaxStamina }; if p.Magic > p.MaxMagic { p.Magic = p.MaxMagic }
}

func (p *Player) StartRaids() {
	ticker := time.NewTicker(30 * time.Minute)
	go func() { for range ticker.C { if rand.Float64() < 0.2 { p.UnderRaid() } } }()
}

func (p *Player) UnderRaid() {
	m := Monster{Name: "Raiders", Health: 200, Damage: 50}
	p.WorldNotice("ALERT: YOUR DOMAIN IS UNDER ATTACK BY RAIDERS!")
	p.Combat(&m, false)
}

func (p *Player) StartGateSpawning() {
	ticker := time.NewTicker(10 * time.Minute)
	go func() { for range ticker.C { p.SpawnGate(); p.Save() } }()
	p.SpawnGate()
}

func (p *Player) SpawnGate() {
	ranks := []string{"E", "D", "C", "B", "A", "S", "SS"}
	r := ranks[rand.Intn(len(ranks))]; g := Gates[r]
	bosses := GateBosses[r]; if len(bosses) > 0 { g.Boss = bosses[rand.Intn(len(bosses))] }
	p.CurrentGate = &g
	p.WorldNotice(fmt.Sprintf("A Rank %s Gate has manifested in the world! Recommended Level: %d", r, g.MinLevel))
}
