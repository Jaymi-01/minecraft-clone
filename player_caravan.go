package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func (p *Player) StartCaravanCycle() {
	ticker := time.NewTicker(60 * time.Minute)
	go func() {
		for range ticker.C {
			if !p.CaravanActive && rand.Float64() < 0.3 {
				p.SpawnCaravan()
			} else if p.CaravanActive {
				p.CaravanActive = false
				p.CaravanInventory = nil
				p.WorldNotice("CARAVAN DEPARTURE: The traveling merchants have moved on.")
			}
		}
	}()
}

func (p *Player) SpawnCaravan() {
	p.CaravanActive = true
	p.CaravanInventory = make(map[string]ShopItem)
	
	rarePool := []ShopItem{
		{ID: "void_essence", Name: "🌌 Void Essence", Price: 5000, Desc: "Pure energy from the void."},
		{ID: "star_matter", Name: "✨ Star Matter", Price: 8000, Desc: "Forged in the heart of a supernova."},
		{ID: "demon_soul", Name: "👹 Demon Lord Soul", Price: 15000, Desc: "The essence of a fallen ruler."},
		{ID: "life_stone", Name: "💎 Life Stone", Price: 25000, Desc: "A stone that defies death."},
		{ID: "abyss_crystal", Name: "🕳️ Abyss Crystal", Price: 10000, Desc: "Radiates dark elemental power."},
	}

	// Pick 3 random rare items
	rand.Shuffle(len(rarePool), func(i, j int) { rarePool[i], rarePool[j] = rarePool[j], rarePool[i] })
	for i := 0; i < 3; i++ {
		p.CaravanInventory[rarePool[i].ID] = rarePool[i]
	}

	p.WorldNotice("CARAVAN ARRIVAL: Traveling merchants have arrived at your domain with rare treasures!")
}

func (p *Player) ListCaravan() {
	if !p.CaravanActive {
		fmt.Println("❌ [SYSTEM]: No active caravans in your sector.")
		return
	}

	fmt.Printf("\n--- 🛒 [TRAVELING MERCHANT CARAVAN] (Your Gold: %d) ---\n", p.Inventory["gold"])
	for id, it := range p.CaravanInventory {
		fmt.Printf("   [%s] %s - 💰 %d Gold\n      📜 %s\n", id, it.Name, it.Price, it.Desc)
	}
	fmt.Println("💡 [TIP]: Use '!buycaravan <id>' to purchase.")
}

func (p *Player) BuyFromCaravan(id string) {
	if !p.CaravanActive {
		fmt.Println("❌ [SYSTEM]: The merchants have already departed.")
		return
	}

	id = strings.ToLower(id)
	it, ok := p.CaravanInventory[id]
	if !ok {
		fmt.Printf("❌ [SYSTEM]: Rare item '%s' is not in the caravan's stock.\n", id)
		return
	}

	if p.Inventory["gold"] < it.Price {
		fmt.Printf("❌ [SYSTEM]: Transaction rejected. Need %d more gold.\n", it.Price-p.Inventory["gold"])
		return
	}

	p.Inventory["gold"] -= it.Price
	p.Inventory[it.ID]++
	p.WorldNotice(fmt.Sprintf("RARE ACQUISITION: %s has been secured.", it.Name))
	p.Save()
}
