package main

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
		Descriptions: []string{
			"The sun shines warmly through the leaves as you gather supplies.",
			"A gentle breeze carries the scent of pine and fresh earth.",
			"Birds chirp overhead while you work on the grassy plains.",
			"You find a sturdy tree and begin harvesting its timber.",
		},
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
		Descriptions: []string{
			"Water droplets echo in the damp silence of the cavern.",
			"The air grows cool and musty as you strike the rocky walls.",
			"Your torch light flickers against a vein of dark coal.",
			"The sound of your pickaxe rings out through the dark tunnels.",
		},
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
		Descriptions: []string{
			"The darkness here feels heavy, almost physical.",
			"Faint whispers seem to drift from the bottomless pits nearby.",
			"A glint of something precious catches your light in the deep gloom.",
			"The rock here is unnaturally hard, singing with a low hum when struck.",
		},
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
		Descriptions: []string{
			"Heat waves distort the air around the bubbling lava lakes.",
			"The ground itself seems to groan under the volcanic pressure.",
			"Cinders fall like glowing snow from the jagged red ceiling.",
			"You chip away at the glowing quartz while shielding your eyes from the glare.",
		},
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
		Descriptions: []string{
			"Reality feels thin here, shimmering like oil on water.",
			"There is no sound, only the vibration of the cosmos in your bones.",
			"Stars seem to drift past you in the endless indigo expanse.",
			"You reach into the nothingness and pull back fragments of pure energy.",
		},
	},
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
	"diamond_sword": {
		Name: "💎 Diamond Sword",
		Ingredients: map[string]int{
			"iron":    5,
			"diamond": 10,
		},
		ResultType:    "weapon",
		ResultValue:   35,
		RequiredLevel: 25,
	},
	"bread": {
		Name: "🥖 Bread",
		Ingredients: map[string]int{
			"wood": 5,
		},
		ResultType:    "food",
		ResultValue:   20,
		RequiredLevel: 2,
	},
	"health_potion": {
		Name: "🧪 Health Potion",
		Ingredients: map[string]int{
			"gel":  10,
			"gold": 5,
		},
		ResultType:    "food",
		ResultValue:   50,
		RequiredLevel: 5,
	},
	"stamina_potion": {
		Name: "⚡ Stamina Potion",
		Ingredients: map[string]int{
			"gel":    10,
			"quartz": 10,
		},
		ResultType:    "stamina_food",
		ResultValue:   30,
		RequiredLevel: 10,
	},
	"leather_armor": {
		Name: "🧥 Leather Armor",
		Ingredients: map[string]int{
			"rotten_flesh": 20,
		},
		ResultType:    "armor",
		ResultValue:   5,
		RequiredLevel: 4,
	},
	"iron_armor": {
		Name: "🛡️ Iron Armor",
		Ingredients: map[string]int{
			"iron": 50,
		},
		ResultType:    "armor",
		ResultValue:   15,
		RequiredLevel: 12,
	},
	"diamond_armor": {
		Name: "💎 Diamond Armor",
		Ingredients: map[string]int{
			"diamond": 25,
		},
		ResultType:    "armor",
		ResultValue:   40,
		RequiredLevel: 30,
	},
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
	"enchanting_table": {
		Name: "🧪 Enchanting Table",
		Ingredients: map[string]int{
			"gold":    50,
			"quartz":  20,
			"diamond": 5,
		},
		RequiredLevel: 15,
		PerkDesc:      "Arcane Wisdom (+50% XP gain from all sources)",
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

var BotSettlements = map[string]BotSettlement{
	"goblin_camp": {
		Name:          "👺 Goblin Camp",
		Level:         5,
		RequiredSword: "sword",
		Defenders: []Monster{
			{Name: "👺 Goblin Warrior", Health: 50, Damage: 10},
			{Name: "👺 Goblin Archer", Health: 30, Damage: 15},
		},
		LootTable:   map[string]int{"wood": 20, "stone": 10, "gold": 5},
		Description: "A small camp of pesky goblins. Easy pickings for a beginner.",
	},
	"bandit_fort": {
		Name:          "🏴‍☠️ Bandit Fort",
		Level:         15,
		RequiredSword: "iron_sword",
		Defenders: []Monster{
			{Name: "🏴‍☠️ Bandit Thug", Health: 100, Damage: 20},
			{Name: "🏴‍☠️ Bandit Leader", Health: 200, Damage: 35},
		},
		LootTable:   map[string]int{"iron": 50, "gold": 25, "diamond": 2},
		Description: "A fortified base filled with ruthless outlaws.",
	},
	"shadow_keep": {
		Name:          "🏰 Shadow Keep",
		Level:         35,
		RequiredSword: "diamond_sword",
		Defenders: []Monster{
			{Name: "👻 Shadow Knight", Health: 300, Damage: 50},
			{Name: "👻 Shadow Mage", Health: 150, Damage: 80},
		},
		LootTable:   map[string]int{"diamond": 20, "netherite": 5, "void_essence": 2},
		Description: "A dark fortress where shadows linger. Extremely dangerous.",
	},
}

var MerchantInventory = map[string]ShopItem{
	"golden_apple": {
		ID:    "golden_apple",
		Name:  "🍎 Golden Apple",
		Price: 50,
		Desc:  "Instantly restores 100 HP.",
	},
	"energy_drink": {
		ID:    "energy_drink",
		Name:  "🥤 Energy Drink",
		Price: 30,
		Desc:  "Instantly restores 50 Stamina.",
	},
	"mystery_box": {
		ID:    "mystery_box",
		Name:  "🎁 Mystery Box",
		Price: 100,
		Desc:  "Contains random high-tier materials.",
	},
	"repair_kit": {
		ID:    "repair_kit",
		Name:  "🔧 Repair Kit",
		Price: 40,
		Desc:  "Fully restores tool durability.",
	},
}

var GlobalQuests = []Quest{
	{ID: "wood_gatherer", Name: "🌲 Wood Gatherer", TargetType: "item", TargetID: "wood", TargetQty: 50, RewardXP: 100, RewardGold: 20, Description: "Collect 50 Wood from the Surface."},
	{ID: "slime_hunter", Name: "🟢 Slime Hunter", TargetType: "combat", TargetID: "🟢 Slime", TargetQty: 5, RewardXP: 150, RewardGold: 30, Description: "Defeat 5 Slimes on the Surface."},
	{ID: "iron_miner", Name: "⛓️ Iron Miner", TargetType: "item", TargetID: "iron", TargetQty: 30, RewardXP: 300, RewardGold: 100, Description: "Mine 30 Iron from the Caves."},
	{ID: "zombie_slayer", Name: "🧟 Zombie Slayer", TargetType: "combat", TargetID: "🧟 Zombie", TargetQty: 10, RewardXP: 500, RewardGold: 200, Description: "Defeat 10 Zombies in the Caves."},
}
