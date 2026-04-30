package main

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
	Descriptions    []string
}

type Monster struct {
	Name      string
	Health    int
	Damage    int
	LootTable map[string]float64
}

type Recipe struct {
	Name          string
	Ingredients   map[string]int
	ResultType    string // "tool", "weapon", "armor", "food", "stamina_food"
	ResultValue   int
	RequiredLevel int
}

type Structure struct {
	Name          string
	Ingredients   map[string]int
	RequiredLevel int
	PerkDesc      string
}

type BotSettlement struct {
	Name          string
	Level         int
	RequiredSword string // The specific sword needed to raid this camp
	Defenders     []Monster
	LootTable     map[string]int
	Description   string
}

type ShopItem struct {
	ID    string
	Name  string
	Price int
	Desc  string
}

type Quest struct {
	ID          string
	Name        string
	TargetType  string // "item", "combat", "level"
	TargetID    string // "iron", "Slime", etc.
	TargetQty   int
	RewardXP    int
	RewardGold  int
	Description string
}

// Player represents the player's state.
type Player struct {
	Name           string          `json:"name"`
	Health         int             `json:"health"`
	MaxHealth      int             `json:"max_health"`
	Attack         int             `json:"attack"`
	Defense        int             `json:"defense"`
	Stamina        int             `json:"stamina"`
	MaxStamina     int             `json:"max_stamina"`
	Level          int             `json:"level"`
	XP             int             `json:"xp"`
	XPToNext       int             `json:"xp_to_next"`
	Inventory      map[string]int  `json:"inventory"`
	ToolDurability int             `json:"tool_durability"`
	Structures     map[string]bool `json:"structures"`
	QuestProgress  map[string]int  `json:"quest_progress"`
}
