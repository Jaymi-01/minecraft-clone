package main

import "time"

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
	Name          string             `json:"name"`
	Health        int                `json:"health"`
	Damage        int                `json:"damage"`
	LootTable     map[string]float64 `json:"loot_table"`
	StatusEffects map[string]int     `json:"status_effects"` // effect -> duration
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

type TrainingProgress struct {
	Pushups  int       `json:"pushups"`
	Situps   int       `json:"situps"`
	Squats   int       `json:"squats"`
	Running  int       `json:"running"`
	LastReset time.Time `json:"last_reset"`
}

type ProductionLog struct {
	LastProduced time.Time `json:"last_produced"`
	PendingItems map[string]int `json:"pending_items"`
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
	Magic          int             `json:"magic"`
	MaxMagic       int             `json:"max_magic"`
	Level          int             `json:"level"`
	XP             int             `json:"xp"`
	XPToNext       int             `json:"xp_to_next"`
	Inventory      map[string]int  `json:"inventory"`
	ToolDurability int             `json:"tool_durability"`
	Structures     map[string]bool `json:"structures"`
	QuestProgress  map[string]int  `json:"quest_progress"`
	Rank           string          `json:"rank"`
	HunterLevel    int             `json:"hunter_level"`
	HunterXP       int             `json:"hunter_xp"`
	HunterXPToNext int             `json:"hunter_xp_to_next"`
	HunterRank     string          `json:"hunter_rank"`
	Kills          int             `json:"kills"`
	MonsterKills   map[string]int  `json:"monster_kills"`
	Taboo          int             `json:"taboo"`
	SkillPoints    int             `json:"skill_points"`
	Titles         []string        `json:"titles"`
	Skills         []string        `json:"skills"`
	EquippedSkills []string        `json:"equipped_skills"`
	SkillSlots     int             `json:"skill_slots"`
	SkillLevels    map[string]int  `json:"skill_levels"`
	SkillUsage     map[string]int  `json:"skill_usage"`
	SkillCooldowns map[string]int  `json:"skill_cooldowns"`
	Subordinates   []Subordinate   `json:"subordinates"`
	Squad          []string        `json:"squad"` // Names of subs in party
	EquippedWeapon string          `json:"equipped_weapon"`
	EquippedArmor  string          `json:"equipped_armor"`
	ItemRarities   map[string]string `json:"item_rarities"` // item_id -> rarity
	ItemLevels     map[string]int    `json:"item_levels"`   // item_id -> +1, +2
	Training       TrainingProgress  `json:"training"`
	Production     ProductionLog     `json:"production"`
	CurrentGate    *Gate           `json:"current_gate"`
	SystemOrigin   string          `json:"system_origin"`
	Exploring      bool            `json:"exploring"`
	ExplorationDepth int           `json:"exploration_depth"`
	ActionLog        []string        `json:"action_log"`
	StatusEffects    map[string]int  `json:"status_effects"` // effect -> duration
	Attributes       map[string]bool `json:"attributes"`     // e.g., "dark": true
}

type Subordinate struct {
	Name       string    `json:"name"`
	Species    string    `json:"species"`
	Attack     int       `json:"attack"`
	Defense    int       `json:"defense"`
	Level      int       `json:"level"`
	XP         int       `json:"xp"`
	NextXP     int       `json:"next_xp"`
	Skills     []string  `json:"skills"`
	LastAction time.Time `json:"last_action"`
}

type Title struct {
	ID           string
	Name         string
	KillsNeeded  int
	PerkDesc     string
	AttackBonus  int
	DefenseBonus int
	HPBonus      int
	MPBonus      int
	StaminaBonus int
}

type Skill struct {
	ID                string
	Name              string
	DmgBonus          int
	Cooldown          int
	MPCost            int
	Desc              string
	Rank              string // E, D, C, B, A, S
	Type              string // "active", "passive"
	Category          string // "attack", "defense", "heal"
	Level             int
	UnlockRequirement string
	ReqBoss           string
	ReqLevel          int
	ReqHunterLevel    int
}

type Gate struct {
	Rank          string
	Floors        int
	MinLevel      int
	Boss          Monster
	MonsterCount  int
	RewardXP      int
	RewardGold    int
	Descriptions  []string
}
