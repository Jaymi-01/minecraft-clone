package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func NewPlayer(name string) *Player {
	return &Player{
		Name:           name,
		Health:         100,
		MaxHealth:      100,
		Attack:         10,
		Defense:        0,
		Stamina:        50,
		MaxStamina:     50,
		Level:          1,
		XP:             0,
		XPToNext:       100,
		Inventory:      map[string]int{"wood_pickaxe": 1},
		ToolDurability: 50,
		Structures:     make(map[string]bool),
		QuestProgress:  make(map[string]int),
	}
}

func (p *Player) HealFull() {
	p.Health = p.MaxHealth
	p.Stamina = p.MaxStamina
	p.Save()
}

func (p *Player) Save() {
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		fmt.Printf("❌ Error saving data: %v\n", err)
		return
	}
	os.WriteFile("player_data.json", data, 0644)
}

func LoadPlayer() *Player {
	data, err := os.ReadFile("player_data.json")
	if err != nil {
		return NewPlayer("Adventurer")
	}
	var p Player
	if err := json.Unmarshal(data, &p); err != nil {
		return NewPlayer("Adventurer")
	}
	if p.QuestProgress == nil {
		p.QuestProgress = make(map[string]int)
	}
	return &p
}

func (p *Player) GainXP(amount int) {
	p.XP += amount
	fmt.Printf("[✨ +%d XP]\n", amount)
	if p.XP >= p.XPToNext {
		p.Level++
		p.XP -= p.XPToNext
		p.XPToNext = int(float64(p.XPToNext) * 1.5)
		p.MaxHealth += 10
		p.MaxStamina += 10
		p.Health = p.MaxHealth
		p.Stamina = p.MaxStamina
		fmt.Printf("\n🎊 LEVEL UP! You are now level %d! 🎊\n", p.Level)
	}
	p.Save()
}

func (p *Player) TrackQuest(qType, id string, qty int) {
	for _, q := range GlobalQuests {
		if q.TargetType == qType && q.TargetID == id {
			if p.QuestProgress[q.ID] < q.TargetQty {
				p.QuestProgress[q.ID] += qty
				if p.QuestProgress[q.ID] >= q.TargetQty {
					fmt.Printf("\n📜 QUEST COMPLETE: %s! 📜\n", q.Name)
					fmt.Printf("🎁 Rewards: ✨ %d XP, 💰 %d Gold\n", q.RewardXP, q.RewardGold)
					p.Inventory["gold"] += q.RewardGold
					p.GainXP(q.RewardXP)
				}
			}
		}
	}
}

func (p *Player) ListQuests() {
	fmt.Println("\n--- 📜 Active Quests ---")
	for _, q := range GlobalQuests {
		status := "✅ Completed"
		prog := p.QuestProgress[q.ID]
		if prog < q.TargetQty {
			status = fmt.Sprintf("⏳ Progress: %d/%d", prog, q.TargetQty)
		}
		fmt.Printf("[%s] %s\n    📝 %s\n    📊 %s\n", q.ID, q.Name, q.Description, status)
	}
	fmt.Println("-------------------------")
}

func (p *Player) ShowStats() {
	fmt.Printf("\n--- 👤 Player Stats ---\n")
	fmt.Printf("⭐ Level:      %d (XP: %d/%d)\n", p.Level, p.XP, p.XPToNext)
	fmt.Printf("❤️ Health:     %d/%d\n", p.Health, p.MaxHealth)
	fmt.Printf("⚔️ Attack:     %d\n", p.Attack)
	fmt.Printf("🛡️ Defense:    %d\n", p.Defense)
	fmt.Printf("⚡ Stamina:    %d/%d\n", p.Stamina, p.MaxStamina)
	fmt.Printf("🔨 Tool Durability: %d\n", p.ToolDurability)
	if len(p.Structures) > 0 {
		fmt.Printf("🏗️ Structures: ")
		var sList []string
		for s := range p.Structures {
			sList = append(sList, s)
		}
		fmt.Println(strings.Join(sList, ", "))
	}
	fmt.Println("----------------------")
}

func (p *Player) ShowInventory() {
	fmt.Printf("\n--- 🎒 Inventory ---\n")
	if len(p.Inventory) == 0 {
		fmt.Println("Empty 📭")
	} else {
		for itemID, qty := range p.Inventory {
			if qty > 0 {
				fmt.Printf("%s: %d\n", itemID, qty)
			}
		}
	}
	fmt.Println("--------------------")
}

func (p *Player) StartRegeneration() {
	ticker := time.NewTicker(20 * time.Minute)
	go func() {
		for range ticker.C {
			p.Regenerate()
		}
	}()
}

func (p *Player) Regenerate() {
	hpRegen := 10
	stRegen := 10

	if p.Structures["house"] {
		hpRegen += 2
	}
	if p.Structures["farm"] {
		stRegen += 5
	}

	if p.Health < p.MaxHealth {
		p.Health += hpRegen
		if p.Health > p.MaxHealth {
			p.Health = p.MaxHealth
		}
	}
	if p.Stamina < p.MaxStamina {
		p.Stamina += stRegen
		if p.Stamina > p.MaxStamina {
			p.Stamina = p.MaxStamina
		}
	}
	p.Save()
}

func (p *Player) StartRaids() {
	ticker := time.NewTicker(30 * time.Minute)
	go func() {
		for range ticker.C {
			if rand.Float64() < 0.3 {
				p.UnderRaid()
			}
		}
	}()
}

func (p *Player) UnderRaid() {
	fmt.Printf("\n🚨 ALERT! Your base is being raided by NPCs! 🚨\n")
	raidStrength := p.Level / 5
	if raidStrength < 1 {
		raidStrength = 1
	}
	raider := Monster{Name: "🏴‍☠️ Raider Party", Health: 50 * raidStrength, Damage: 10 * raidStrength}
	if p.Combat(&raider) {
		fmt.Println("🛡️ You successfully defended your base!")
	} else {
		fmt.Println("📉 The raiders plundered some of your resources!")
		for item, qty := range p.Inventory {
			if qty > 5 {
				lost := rand.Intn(qty / 2)
				p.Inventory[item] -= lost
				if lost > 0 {
					fmt.Printf("💸 Lost %d %s\n", lost, item)
				}
			}
		}
	}
	p.Save()
}

func (p *Player) ListRaids() {
	fmt.Println("\n--- ⚔️ Raid Targets ---")
	for id, s := range BotSettlements {
		fmt.Printf("[%s] %s (⭐ Lvl %d)\n    📝 %s\n", id, s.Name, s.Level, s.Description)
	}
	fmt.Println("-----------------------")
}

func (p *Player) Raid(targetID string) {
	target, ok := BotSettlements[strings.ToLower(targetID)]
	if !ok {
		fmt.Printf("❓ Unknown target: %s. Type !raid to see list.\n", targetID)
		return
	}
	if p.Level < target.Level {
		fmt.Printf("🚫 Your level is too low to raid %s! Required: %d\n", target.Name, target.Level)
		return
	}
	if p.Stamina < 30 {
		fmt.Println("😫 Raiding requires 30 stamina! Wait for regeneration.")
		return
	}
	p.Stamina -= 30
	fmt.Printf("🚀 Starting raid on %s...\n", target.Name)
	for _, defender := range target.Defenders {
		fmt.Printf("⚔️ Facing defender: %s\n", defender.Name)
		if !p.Combat(&defender) {
			fmt.Printf("❌ Raid failed! You were driven back from %s.\n", target.Name)
			return
		}
	}
	fmt.Printf("💰 SUCCESS! You conquered %s and plundered their vault!\n", target.Name)
	for item, qty := range target.LootTable {
		p.Inventory[item] += qty
		fmt.Printf("🎁 Found %d %s\n", qty, item)
	}
	p.GainXP(100 + (target.Level * 10))
	p.Save()
}

func (p *Player) ListShop() {
	fmt.Println("\n--- ⚖️ Merchant's Shop ---")
	fmt.Printf("Your Gold: 💰 %d\n", p.Inventory["gold"])
	for id, item := range MerchantInventory {
		fmt.Printf("[%s] %s - 💰 %d\n    📝 %s\n", id, item.Name, item.Price, item.Desc)
	}
	fmt.Println("--------------------------")
}

func (p *Player) Buy(itemID string) {
	item, ok := MerchantInventory[strings.ToLower(itemID)]
	if !ok {
		fmt.Printf("❓ Merchant says: 'I don't have a %s for sale!'\n", itemID)
		return
	}
	if p.Inventory["gold"] < item.Price {
		fmt.Printf("🚫 Merchant says: 'You need 💰 %d gold for that, you only have 💰 %d!'\n", item.Price, p.Inventory["gold"])
		return
	}
	p.Inventory["gold"] -= item.Price
	switch item.ID {
	case "golden_apple":
		p.Health += 100
		if p.Health > p.MaxHealth {
			p.Health = p.MaxHealth
		}
		fmt.Printf("🍎 You bought and ate a Golden Apple! Health restored to %d.\n", p.Health)
	case "energy_drink":
		p.Stamina += 50
		if p.Stamina > p.MaxStamina {
			p.Stamina = p.MaxStamina
		}
		fmt.Printf("🥤 You bought and drank an Energy Drink! Stamina restored to %d.\n", p.Stamina)
	case "repair_kit":
		p.ToolDurability = 500
		fmt.Printf("🔧 You bought a Repair Kit! Your tool is now extremely durable (%d).\n", p.ToolDurability)
	case "mystery_box":
		fmt.Printf("🎁 You opened a Mystery Box and found: ")
		lootPool := []string{"iron", "gold", "diamond", "quartz", "netherite"}
		for i := 0; i < 3; i++ {
			loot := lootPool[rand.Intn(len(lootPool))]
			qty := 5 + rand.Intn(10)
			p.Inventory[loot] += qty
			fmt.Printf("%d %s, ", qty, loot)
		}
		fmt.Println("Not bad!")
	default:
		p.Inventory[item.ID]++
		fmt.Printf("⚖️ You bought 1 %s for 💰 %d gold.\n", item.Name, item.Price)
	}
	p.Save()
}

func (p *Player) ListCraftable() {
	fmt.Println("\n--- 📜 Crafting Menu ---")
	for id, r := range Recipes {
		fmt.Printf("[%s] ⭐ Lvl %d - 📦 Ingredients: ", id, r.RequiredLevel)
		var ingList []string
		for ing, qty := range r.Ingredients {
			ingList = append(ingList, fmt.Sprintf("%d %s", qty, ing))
		}
		fmt.Printf("%s\n", strings.Join(ingList, ", "))
	}
	fmt.Println("------------------------")
}

func (p *Player) Craft(itemName string) {
	recipe, ok := Recipes[strings.ToLower(itemName)]
	if !ok {
		fmt.Printf("❓ Unknown recipe: %s. Type !craft to see options.\n", itemName)
		return
	}
	if p.Level < recipe.RequiredLevel {
		fmt.Printf("🚫 Your level is too low to craft %s! Required: %d\n", recipe.Name, recipe.RequiredLevel)
		return
	}
	for ing, qty := range recipe.Ingredients {
		if p.Inventory[ing] < qty {
			fmt.Printf("❌ Missing ingredients for %s: Need %d %s, have %d\n", recipe.Name, qty, ing, p.Inventory[ing])
			return
		}
	}
	for ing, qty := range recipe.Ingredients {
		p.Inventory[ing] -= qty
	}
	switch recipe.ResultType {
	case "tool":
		p.ToolDurability = recipe.ResultValue
		p.Inventory[strings.ToLower(itemName)] = 1
		fmt.Printf("🛠️ You crafted a %s! Tool durability set to %d.\n", recipe.Name, p.ToolDurability)
	case "weapon":
		p.Attack += recipe.ResultValue
		fmt.Printf("⚔️ You crafted a %s! Attack increased by %d (Total: %d).\n", recipe.Name, recipe.ResultValue, p.Attack)
	case "armor":
		p.Defense += recipe.ResultValue
		fmt.Printf("🛡️ You crafted a %s! Defense increased by %d (Total: %d).\n", recipe.Name, recipe.ResultValue, p.Defense)
	case "food":
		p.Inventory[strings.ToLower(itemName)]++
		fmt.Printf("🍞 You crafted a %s!\n", recipe.Name)
	case "stamina_food":
		p.Inventory[strings.ToLower(itemName)]++
		fmt.Printf("⚡ You crafted a %s!\n", recipe.Name)
	}
	p.GainXP(10 + rand.Intn(5))
	p.Save()
}

func (p *Player) Use(itemName string) {
	itemKey := strings.ToLower(itemName)
	if p.Inventory[itemKey] <= 0 {
		fmt.Printf("❌ You don't have any %s in your inventory.\n", itemName)
		return
	}
	recipe, ok := Recipes[itemKey]
	if !ok || (recipe.ResultType != "food" && recipe.ResultType != "stamina_food") {
		fmt.Printf("❌ %s is not a consumable item.\n", itemName)
		return
	}
	p.Inventory[itemKey]--
	if p.Inventory[itemKey] == 0 {
		delete(p.Inventory, itemKey)
	}
	switch recipe.ResultType {
	case "food":
		oldHP := p.Health
		p.Health += recipe.ResultValue
		if p.Health > p.MaxHealth {
			p.Health = p.MaxHealth
		}
		fmt.Printf("😋 You consumed %s and recovered %d HP! (❤️ %d -> %d)\n", recipe.Name, p.Health-oldHP, oldHP, p.Health)
	case "stamina_food":
		oldStam := p.Stamina
		p.Stamina += recipe.ResultValue
		if p.Stamina > p.MaxStamina {
			p.Stamina = p.MaxStamina
		}
		fmt.Printf("⚡ You consumed %s and recovered %d Stamina! (⚡ %d -> %d)\n", recipe.Name, p.Stamina-oldStam, oldStam, p.Stamina)
	}
	p.Save()
}

func (p *Player) ListBuildable() {
	fmt.Println("\n--- 🏗️ Building Menu ---")
	for id, s := range Structures {
		fmt.Printf("[%s] ⭐ Lvl %d - 📦 Ingredients: ", id, s.RequiredLevel)
		var ingList []string
		for ing, qty := range s.Ingredients {
			ingList = append(ingList, fmt.Sprintf("%d %s", qty, ing))
		}
		fmt.Printf("%s\n    🎁 Perk: %s\n", strings.Join(ingList, ", "), s.PerkDesc)
	}
	fmt.Println("------------------------")
}

func (p *Player) Build(structName string) {
	s, ok := Structures[strings.ToLower(structName)]
	if !ok {
		fmt.Printf("❓ Unknown structure: %s. Type !build to see options.\n", structName)
		return
	}
	if p.Structures[strings.ToLower(structName)] {
		fmt.Printf("🏠 You already built a %s!\n", s.Name)
		return
	}
	if p.Level < s.RequiredLevel {
		fmt.Printf("🚫 Your level is too low to build %s! Required: %d\n", s.Name, s.RequiredLevel)
		return
	}
	for ing, qty := range s.Ingredients {
		if p.Inventory[ing] < qty {
			fmt.Printf("❌ Missing materials for %s: Need %d %s, have %d\n", s.Name, qty, ing, p.Inventory[ing])
			return
		}
	}
	for ing, qty := range s.Ingredients {
		p.Inventory[ing] -= qty
	}
	p.Structures[strings.ToLower(structName)] = true
	fmt.Printf("🔨 You built a %s! Perk Unlocked: %s\n", s.Name, s.PerkDesc)
	switch strings.ToLower(structName) {
	case "forge":
		p.Attack += 10
	case "vault":
		p.MaxHealth += 50
		p.Health += 50
	case "castle":
		p.Attack += 20
		p.MaxHealth += 100
		p.Health += 100
		p.MaxStamina += 50
		p.Stamina += 50
	}
	p.GainXP(50 + rand.Intn(50))
	p.Save()
}

func (p *Player) Combat(m *Monster) bool {
	fmt.Printf("\n⚔️ A wild %s appeared!\n", m.Name)
	monsterHealth := m.Health
	for monsterHealth > 0 && p.Health > 0 {
		damageToMonster := p.Attack + rand.Intn(5)
		monsterHealth -= damageToMonster
		fmt.Printf("🤜 You hit %s for %d damage. (%d HP left)\n", m.Name, damageToMonster, monsterHealth)
		if monsterHealth <= 0 {
			fmt.Printf("🏆 You defeated the %s!\n", m.Name)
			for item, prob := range m.LootTable {
				if rand.Float64() <= prob {
					p.Inventory[item]++
					fmt.Printf("🎁 Dropped: %s\n", item)
				}
			}
			p.TrackQuest("combat", m.Name, 1)
			p.GainXP(15 + rand.Intn(10))
			return true
		}
		baseDamage := m.Damage + rand.Intn(5)
		finalDamage := baseDamage - p.Defense
		if finalDamage < 1 {
			finalDamage = 1
		}
		p.Health -= finalDamage
		fmt.Printf("💥 %s hits you for %d damage (Blocked %d). (%d HP left)\n", m.Name, finalDamage, p.Defense, p.Health)
	}
	if p.Health <= 0 {
		fmt.Println("💀 You were defeated...")
		p.Health = 20
		fmt.Println("🩹 You limped back to safety and recovered a bit of health.")
		p.Save()
		return false
	}
	return false
}

func (p *Player) Mine(locName string) {
	loc, ok := Locations[strings.ToLower(locName)]
	if !ok {
		fmt.Printf("❓ Unknown location: %s. Type !mine to see available zones.\n", locName)
		return
	}
	if p.Level < loc.RequiredLevel {
		fmt.Printf("🚫 Your level is too low to enter %s! Required: %d\n", loc.Name, loc.RequiredLevel)
		return
	}
	if loc.RequiredItem != "" && p.Inventory[loc.RequiredItem] <= 0 {
		fmt.Printf("🔏 You need a %s to mine in the %s!\n", loc.RequiredItem, loc.Name)
		return
	}
	if p.Stamina < 10 {
		fmt.Println("😫 Not enough stamina! Wait for regeneration.")
		return
	}
	if p.ToolDurability <= 0 {
		fmt.Println("⚠️ Your tool is broken! Craft a new one.")
		return
	}
	p.Stamina -= 10
	p.ToolDurability -= 1
	if len(loc.Descriptions) > 0 {
		desc := loc.Descriptions[rand.Intn(len(loc.Descriptions))]
		fmt.Printf("\n✨ %s\n", desc)
	}
	if rand.Float64() <= loc.EncounterChance {
		monster := loc.EncounterTable[rand.Intn(len(loc.EncounterTable))]
		if !p.Combat(&monster) {
			return
		}
	}
	numDrops := 1 + (p.Level / 5)
	foundSomething := false
	for i := 0; i < numDrops; i++ {
		r := rand.Float64()
		var cumulative float64
		for item, prob := range loc.LootTable {
			cumulative += prob
			if r <= cumulative {
				p.Inventory[item]++
				fmt.Printf("⛏️ You mined in the %s and found: %s!\n", loc.Name, item)
				p.TrackQuest("item", item, 1)
				foundSomething = true
				break
			}
		}
	}
	if !foundSomething {
		fmt.Printf("💨 You mined in the %s but found nothing.\n", loc.Name)
	}
	p.GainXP(2 + rand.Intn(3))
	p.Save()
}
