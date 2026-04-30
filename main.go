package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Item represents a game item.
type Item struct {
	ID   string
	Name string
}

type Location struct {
	Name            string
	LootTable       map[string]float64
	EncounterTable  []Monster
	EncounterChance float64
	RequiredLevel   int
	RequiredItem    string
}

type Monster struct {
	Name      string
	Health    int
	Damage    int
	LootTable map[string]float64
}

var Locations = map[string]Location{
	"surface": {
		Name: "🌳 Surface",
		LootTable: map[string]float64{
			"wood":  0.6,
			"stone": 0.4,
		},
		EncounterChance: 0.1,
		EncounterTable: []Monster{
			{Name: "🟢 Slime", Health: 20, Damage: 5, LootTable: map[string]float64{"gel": 1.0}},
		},
		RequiredLevel: 1,
	},
	"cave": {
		Name: "🕳️ Cave",
		LootTable: map[string]float64{
			"iron": 0.5,
			"coal": 0.5,
		},
		EncounterChance: 0.2,
		EncounterTable: []Monster{
			{Name: "🧟 Zombie", Health: 40, Damage: 10, LootTable: map[string]float64{"rotten_flesh": 1.0}},
			{Name: "🕷️ Spider", Health: 30, Damage: 12, LootTable: map[string]float64{"string": 1.0}},
		},
		RequiredLevel: 3,
		RequiredItem:  "stone_pickaxe",
	},
	"abyss": {
		Name: "🕳️ Abyss",
		LootTable: map[string]float64{
			"gold":    0.6,
			"diamond": 0.4,
		},
		EncounterChance: 0.3,
		EncounterTable: []Monster{
			{Name: "👻 Shadow Stalker", Health: 60, Damage: 20, LootTable: map[string]float64{"shadow_dust": 1.0}},
		},
		RequiredLevel: 10,
		RequiredItem:  "iron_pickaxe",
	},
	"nether": {
		Name: "🔥 Nether",
		LootTable: map[string]float64{
			"quartz":    0.7,
			"netherite": 0.3,
		},
		EncounterChance: 0.4,
		EncounterTable: []Monster{
			{Name: "🔥 Blaze", Health: 80, Damage: 25, LootTable: map[string]float64{"blaze_rod": 1.0}},
			{Name: "☁️ Ghast", Health: 50, Damage: 40, LootTable: map[string]float64{"ghast_tear": 1.0}},
		},
		RequiredLevel: 25,
		RequiredItem:  "diamond_pickaxe",
	},
	"void": {
		Name: "🌌 Void",
		LootTable: map[string]float64{
			"void_essence": 0.8,
			"star_matter":  0.2,
		},
		EncounterChance: 0.5,
		EncounterTable: []Monster{
			{Name: "👾 Void Reaver", Health: 150, Damage: 50, LootTable: map[string]float64{"void_core": 1.0}},
		},
		RequiredLevel: 50,
		RequiredItem:  "netherite_pickaxe",
	},
}

type Recipe struct {
	Name          string
	Ingredients   map[string]int
	ResultType    string // "tool", "weapon", "armor"
	ResultValue   int
	RequiredLevel int
}

var Recipes = map[string]Recipe{
	"stone_pickaxe": {
		Name: "🪨 Stone Pickaxe",
		Ingredients: map[string]int{
			"wood":  10,
			"stone": 20,
		},
		ResultType:    "tool",
		ResultValue:   100,
		RequiredLevel: 3,
	},
	"iron_pickaxe": {
		Name: "⛓️ Iron Pickaxe",
		Ingredients: map[string]int{
			"wood": 5,
			"iron": 20,
		},
		ResultType:    "tool",
		ResultValue:   250,
		RequiredLevel: 10,
	},
	"diamond_pickaxe": {
		Name: "💎 Diamond Pickaxe",
		Ingredients: map[string]int{
			"iron":    10,
			"diamond": 5,
		},
		ResultType:    "tool",
		ResultValue:   500,
		RequiredLevel: 25,
	},
	"sword": {
		Name: "🗡️ Sword",
		Ingredients: map[string]int{
			"wood":  2,
			"stone": 10,
		},
		ResultType:    "weapon",
		ResultValue:   5,
		RequiredLevel: 1,
	},
	"iron_sword": {
		Name: "⚔️ Iron Sword",
		Ingredients: map[string]int{
			"wood": 2,
			"iron": 15,
		},
		ResultType:    "weapon",
		ResultValue:   15,
		RequiredLevel: 8,
	},
}

type Structure struct {
	Name          string
	Ingredients   map[string]int
	RequiredLevel int
	PerkDesc      string
}

var Structures = map[string]Structure{
	"house": {
		Name: "🏠 House",
		Ingredients: map[string]int{
			"wood":  50,
			"stone": 50,
		},
		RequiredLevel: 5,
		PerkDesc:      "Increases health regeneration (+2 per cycle)",
	},
	"farm": {
		Name: "🌾 Farm",
		Ingredients: map[string]int{
			"wood": 100,
			"gel":  20,
		},
		RequiredLevel: 8,
		PerkDesc:      "Increases stamina regeneration (+5 per cycle)",
	},
	"forge": {
		Name: "⚒️ Forge",
		Ingredients: map[string]int{
			"stone": 200,
			"coal":  50,
			"iron":  20,
		},
		RequiredLevel: 12,
		PerkDesc:      "Increases attack power (+10 attack)",
	},
	"vault": {
		Name: "🏦 Vault",
		Ingredients: map[string]int{
			"iron": 100,
			"gold": 20,
		},
		RequiredLevel: 20,
		PerkDesc:      "Increases max health (+50 HP)",
	},
	"castle": {
		Name: "🏰 Castle",
		Ingredients: map[string]int{
			"stone":   1000,
			"iron":    200,
			"diamond": 50,
		},
		RequiredLevel: 40,
		PerkDesc:      "Global Mastery (+20 attack, +100 HP, +50 Stamina)",
	},
}

type BotSettlement struct {
	Name        string
	Level       int
	Defenders   []Monster
	LootTable   map[string]int
	Description string
}

var BotSettlements = map[string]BotSettlement{
	"goblin_camp": {
		Name:  "👺 Goblin Camp",
		Level: 5,
		Defenders: []Monster{
			{Name: "👺 Goblin Warrior", Health: 50, Damage: 10},
			{Name: "👺 Goblin Archer", Health: 30, Damage: 15},
		},
		LootTable:   map[string]int{"wood": 20, "stone": 10, "gold": 5},
		Description: "A small camp of pesky goblins. Easy pickings for a beginner.",
	},
	"bandit_fort": {
		Name:  "🏴‍☠️ Bandit Fort",
		Level: 15,
		Defenders: []Monster{
			{Name: "🏴‍☠️ Bandit Thug", Health: 100, Damage: 20},
			{Name: "🏴‍☠️ Bandit Leader", Health: 200, Damage: 35},
		},
		LootTable:   map[string]int{"iron": 50, "gold": 25, "diamond": 2},
		Description: "A fortified base filled with ruthless outlaws.",
	},
	"shadow_keep": {
		Name:  "🏰 Shadow Keep",
		Level: 35,
		Defenders: []Monster{
			{Name: "👻 Shadow Knight", Health: 300, Damage: 50},
			{Name: "👻 Shadow Mage", Health: 150, Damage: 80},
		},
		LootTable:   map[string]int{"diamond": 20, "netherite": 5, "void_essence": 2},
		Description: "A dark fortress where shadows linger. Extremely dangerous.",
	},
}

// Player represents the player's state.
type Player struct {
	Health         int             `json:"health"`
	MaxHealth      int             `json:"max_health"`
	Attack         int             `json:"attack"`
	Stamina        int             `json:"stamina"`
	MaxStamina     int             `json:"max_stamina"`
	Level          int             `json:"level"`
	XP             int             `json:"xp"`
	XPToNext       int             `json:"xp_to_next"`
	Inventory      map[string]int  `json:"inventory"`
	ToolDurability int             `json:"tool_durability"`
	Structures     map[string]bool `json:"structures"`
}

func NewPlayer() *Player {
	return &Player{
		Health:         100,
		MaxHealth:      100,
		Attack:         10,
		Stamina:        50,
		MaxStamina:     50,
		Level:          1,
		XP:             0,
		XPToNext:       100,
		Inventory:      map[string]int{"wood_pickaxe": 1},
		ToolDurability: 50,
		Structures:     make(map[string]bool),
	}
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
		return NewPlayer()
	}
	var p Player
	if err := json.Unmarshal(data, &p); err != nil {
		return NewPlayer()
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

func (p *Player) ShowStats() {
	fmt.Printf("\n--- 👤 Player Stats ---\n")
	fmt.Printf("⭐ Level:      %d (XP: %d/%d)\n", p.Level, p.XP, p.XPToNext)
	fmt.Printf("❤️ Health:     %d/%d\n", p.Health, p.MaxHealth)
	fmt.Printf("⚔️ Attack:     %d\n", p.Attack)
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
	}
	p.GainXP(10 + rand.Intn(5))
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
			p.GainXP(15 + rand.Intn(10))
			return true
		}
		damageToPlayer := m.Damage + rand.Intn(5)
		p.Health -= damageToPlayer
		fmt.Printf("💥 %s hits you for %d damage. (%d HP left)\n", m.Name, damageToPlayer, p.Health)
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

func main() {
	rand.Seed(time.Now().UnixNano())
	player := LoadPlayer()
	player.StartRegeneration()
	player.StartRaids()
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("🌟 Welcome back to the Mine & Exploration System! 🌟")
	fmt.Println("Available Commands: !mine <location>, !craft [item], !build [structure], !raid [target], !stats, !inventory, !exit")

	for {
		fmt.Print("\n> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			continue
		}
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}
		parts := strings.Fields(input)
		command := parts[0]
		switch command {
		case "!mine":
			if len(parts) < 2 {
				fmt.Println("📍 Available Locations: surface, cave, abyss, nether, void")
				fmt.Println("Usage: !mine <location>")
			} else {
				player.Mine(parts[1])
			}
		case "!craft":
			if len(parts) < 2 {
				player.ListCraftable()
			} else {
				player.Craft(parts[1])
			}
		case "!build":
			if len(parts) < 2 {
				player.ListBuildable()
			} else {
				player.Build(parts[1])
			}
		case "!raid":
			if len(parts) < 2 {
				player.ListRaids()
			} else {
				player.Raid(parts[1])
			}
		case "!stats":
			player.ShowStats()
		case "!inventory":
			player.ShowInventory()
		case "!exit":
			player.Save()
			fmt.Println("👋 Goodbye! Your progress has been saved to player_data.json.")
			return
		default:
			fmt.Printf("❓ Unknown command: %s\n", command)
		}
	}
}
