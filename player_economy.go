package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func (p *Player) ListShop() {
	fmt.Printf("\n--- 💰 Shop (Gold: %d) ---\n", p.Inventory["gold"])
	for id, it := range MerchantInventory { fmt.Printf("[%s] %s (💰%d)\n", id, it.Name, it.Price) }
}

func (p *Player) Buy(id string) {
	id = strings.ToLower(id); it, ok := MerchantInventory[id]
	if !ok || p.Inventory["gold"] < it.Price { fmt.Println("❌ Transaction failed."); return }
	p.Inventory["gold"] -= it.Price; p.Inventory[it.ID]++; p.WorldNotice("Bought: " + it.Name); p.Save()
}

func (p *Player) ListTabooShop() {
	fmt.Printf("\n--- 🌌 Taboo Sanctum (Taboo: %d) ---\n", p.Taboo)
	for id, it := range TabooShopInventory { fmt.Printf("[%s] %s (%d Taboo)\n", id, it.Name, it.Price) }
}

func (p *Player) BuyTabooSkill(id string) {
	id = strings.ToLower(id); it, ok := TabooShopInventory[id]
	if !ok || p.Taboo < it.Price { fmt.Println("❌ Insight rejected."); return }
	p.Taboo -= it.Price; if p.Attributes == nil { p.Attributes = make(map[string]bool) }; p.Attributes[id] = true; p.WorldNotice("ATTRIBUTE MASTERED: " + it.Name); p.Save()
}

func (p *Player) ListElementalShop() {
	fmt.Println("\n--- 🔥 Elemental Shop ---")
	for id, it := range ElementalSkillShopInventory { fmt.Printf("[%s] %s (💰%d)\n", id, it.Name, it.Price) }
}

func (p *Player) BuyElementalSkill(id string) {
	id = strings.ToLower(id); it, ok := ElementalSkillShopInventory[id]
	if !ok || p.Inventory["gold"] < it.Price { fmt.Println("❌ Mastery failed."); return }
	p.Inventory["gold"] -= it.Price; p.AddSkill(id); p.Save()
}

func (p *Player) ListDCraftable() {
	fmt.Println("\n--- 🛠️ Crafting & Structures ---")
	for id, r := range Recipes { fmt.Printf("[%s] %s (Req Lvl %d)\n", id, r.Name, r.RequiredLevel) }
	for id, s := range Structures { fmt.Printf("[%s] %s (Req Lvl %d)\n", id, s.Name, s.RequiredLevel) }
}

func (p *Player) ListRaids() {
	fmt.Println("\n--- ⚔️ Bot Settlements ---")
	for id, t := range BotSettlements { fmt.Printf("[%s] %s (Level %d)\n", id, t.Name, t.Level) }
}

func (p *Player) Craft(itemName string) {
	itemName = strings.ToLower(itemName); r, ok := Recipes[itemName]
	if !ok { fmt.Println("❌ No recipe."); return }
	if p.Level < r.RequiredLevel { fmt.Printf("🚫 Need Level %d.\n", r.RequiredLevel); return }
	for k, v := range r.Ingredients { if p.Inventory[k] < v { fmt.Printf("📦 Need %d more %s.\n", v-p.Inventory[k], k); return } }
	for k, v := range r.Ingredients { p.Inventory[k] -= v; if p.Inventory[k] == 0 { delete(p.Inventory, k) } }
	p.Inventory[itemName]++; p.WorldNotice("Crafted: " + r.Name); p.Save()
}

func (p *Player) Build(id string) {
	s, ok := Structures[strings.ToLower(id)]; if !ok || p.Level < s.RequiredLevel || p.Structures[id] { return }
	for k, v := range s.Ingredients { if p.Inventory[k] < v { fmt.Printf("📦 Need %d more %s.\n", v-p.Inventory[k], k); return } }
	for k, v := range s.Ingredients { p.Inventory[k] -= v; if p.Inventory[k] == 0 { delete(p.Inventory, k) } }
	p.Structures[id] = true; p.SyncStats(); p.WorldNotice("CONSTRUCTION COMPLETE: " + s.Name); p.Save()
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
	go func() { for range ticker.C { if rand.Float64() < 0.2 { p.Combat(&Monster{Name: "Raiders", Health: 200, Damage: 50}, false) } } }()
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
	p.CurrentGate = &g; p.WorldNotice("Gate Manifested: " + r)
}
