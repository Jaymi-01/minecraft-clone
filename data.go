package main

var Locations = map[string]Location{
	"surface": { Name: "🌳 Surface", LootTable: map[string]float64{"wood": 0.6, "stone": 0.4}, EncounterChance: 0.1, EncounterTable: []Monster{{Name: "🟢 Slime", Health: 20, Damage: 5, LootTable: map[string]float64{"gel": 1.0}}}, RequiredLevel: 1 },
	"cave":    { Name: "🕳️ Cave", LootTable: map[string]float64{"iron": 0.5, "coal": 0.5}, EncounterChance: 0.2, EncounterTable: []Monster{{Name: "🧟 Zombie", Health: 40, Damage: 10, LootTable: map[string]float64{"rotten_flesh": 1.0}}}, RequiredLevel: 10, RequiredItem: "stone_pickaxe" },
	"abyss":   { Name: "🌑 Abyss", LootTable: map[string]float64{"gold": 0.5, "diamond": 0.3, "abyss_crystal": 0.2}, EncounterChance: 0.3, EncounterTable: []Monster{{Name: "👻 Shadow Stalker", Health: 60, Damage: 20, LootTable: map[string]float64{"shadow_dust": 1.0}}}, RequiredLevel: 20, RequiredItem: "iron_pickaxe" },
	"nether":  { Name: "🔥 Nether", LootTable: map[string]float64{"quartz": 0.5, "netherite": 0.3, "nether_crystal": 0.2}, EncounterChance: 0.4, EncounterTable: []Monster{{Name: "🔥 Blaze", Health: 80, Damage: 25, LootTable: map[string]float64{"blaze_rod": 1.0}}}, RequiredLevel: 30, RequiredItem: "diamond_pickaxe" },
	"void":    { Name: "🌌 Void", LootTable: map[string]float64{"void_essence": 0.5, "star_matter": 0.3, "void_crystal": 0.2}, EncounterChance: 0.5, EncounterTable: []Monster{{Name: "👾 Void Reaver", Health: 150, Damage: 50, LootTable: map[string]float64{"void_core": 1.0}}}, RequiredLevel: 40, RequiredItem: "nether_pickaxe" },
}

var Recipes = map[string]Recipe{
	"wood_pickaxe":    { Name: "🪵⛏️ Wood Pickaxe", Ingredients: map[string]int{"wood": 10}, ResultType: "tool", ResultValue: 50, RequiredLevel: 1 },
	"stone_pickaxe":   { Name: "🪨⛏️ Stone Pickaxe", Ingredients: map[string]int{"wood": 10, "stone": 20}, ResultType: "tool", ResultValue: 100, RequiredLevel: 10 },
	"iron_pickaxe":    { Name: "⛓️⛏️ Iron Pickaxe", Ingredients: map[string]int{"wood": 5, "iron": 20}, ResultType: "tool", ResultValue: 250, RequiredLevel: 20 },
	"diamond_pickaxe": { Name: "💎⛏️ Diamond Pickaxe", Ingredients: map[string]int{"iron": 10, "diamond": 5}, ResultType: "tool", ResultValue: 500, RequiredLevel: 30 },
	"nether_pickaxe":  { Name: "🔥⛏️ Nether Pickaxe", Ingredients: map[string]int{"netherite": 10, "nether_crystal": 5}, ResultType: "tool", ResultValue: 2500, RequiredLevel: 40 },
	"void_pickaxe":    { Name: "🌌⛏️ Void Pickaxe", Ingredients: map[string]int{"star_matter": 10, "void_crystal": 5}, ResultType: "tool", ResultValue: 5000, RequiredLevel: 50 },
	"sword":           { Name: "🗡️ Sword", Ingredients: map[string]int{"wood": 2, "stone": 10}, ResultType: "weapon", ResultValue: 5, RequiredLevel: 1 },
	"iron_sword":      { Name: "⚔️ Iron Sword", Ingredients: map[string]int{"wood": 2, "iron": 15}, ResultType: "weapon", ResultValue: 15, RequiredLevel: 10 },
	"diamond_sword":   { Name: "💎⚔️ Diamond Sword", Ingredients: map[string]int{"iron": 5, "diamond": 10}, ResultType: "weapon", ResultValue: 35, RequiredLevel: 20 },
	"abyss_sword":     { Name: "🌑⚔️ Abyss Sword", Ingredients: map[string]int{"diamond": 10, "abyss_crystal": 10}, ResultType: "weapon", ResultValue: 60, RequiredLevel: 30 },
	"nether_sword":    { Name: "🔥⚔️ Nether Sword", Ingredients: map[string]int{"netherite": 20, "nether_crystal": 10}, ResultType: "weapon", ResultValue: 120, RequiredLevel: 40 },
	"void_sword":      { Name: "🌌⚔️ Void Sword", Ingredients: map[string]int{"star_matter": 20, "void_crystal": 10}, ResultType: "weapon", ResultValue: 250, RequiredLevel: 50 },
	"d_hunter_blade":  { Name: "🗡️ Hunter's Blade", Ingredients: map[string]int{"iron": 20, "goblin_ear": 5}, ResultType: "weapon", ResultValue: 50, RequiredLevel: 10 },
	"d_wolf_slayer":   { Name: "🐺 Wolf Slayer Greatsword", Ingredients: map[string]int{"diamond": 15, "wolf_fang": 10}, ResultType: "weapon", ResultValue: 120, RequiredLevel: 20 },
	"d_naga_trident":  { Name: "🔱 Naga Trident", Ingredients: map[string]int{"abyss_crystal": 20, "naga_scale": 15}, ResultType: "weapon", ResultValue: 300, RequiredLevel: 30 },
	"d_golem_smasher": { Name: "🔨 Golem Smasher", Ingredients: map[string]int{"netherite": 25, "core_fragment": 20}, ResultType: "weapon", ResultValue: 800, RequiredLevel: 40 },
	"d_dragon_slayer": { Name: "🐲 Dragon Slayer", Ingredients: map[string]int{"star_matter": 30, "dragon_heart": 10}, ResultType: "weapon", ResultValue: 2000, RequiredLevel: 50 },
	"d_monarch_sword": { Name: "👑 Monarch's Sword", Ingredients: map[string]int{"void_crystal": 50, "demon_soul": 5}, ResultType: "weapon", ResultValue: 5000, RequiredLevel: 100 },
	"leather_armor":   { Name: "🛡️ Leather Armor", Ingredients: map[string]int{"wood": 20, "stone": 10}, ResultType: "armor", ResultValue: 5, RequiredLevel: 10 },
	"plate_armor":     { Name: "🛡️ Plate Armor", Ingredients: map[string]int{"iron": 50, "gold": 10}, ResultType: "armor", ResultValue: 20, RequiredLevel: 20 },
	"dragon_scale":    { Name: "🛡️ Dragon Scale", Ingredients: map[string]int{"netherite": 20, "dragon_heart": 5}, ResultType: "armor", ResultValue: 100, RequiredLevel: 40 },
	"bread":           { Name: "🥖 Bread", Ingredients: map[string]int{"wood": 5}, ResultType: "food", ResultValue: 20, RequiredLevel: 1 },
	"health_potion":   { Name: "🧪 Health Potion", Ingredients: map[string]int{"gel": 10, "gold": 5}, ResultType: "food", ResultValue: 50, RequiredLevel: 10 },
	"stamina_potion":  { Name: "⚡ Stamina Potion", Ingredients: map[string]int{"gel": 10, "quartz": 10}, ResultType: "stamina_food", ResultValue: 30, RequiredLevel: 20 },
}

var Structures = map[string]Structure{
	"house":            { Name: "🏠 House", Ingredients: map[string]int{"wood": 50, "stone": 50}, RequiredLevel: 10, PerkDesc: "Increases health regeneration (+2 per cycle)" },
	"farm":             { Name: "🌾 Farm", Ingredients: map[string]int{"wood": 100, "gel": 20}, RequiredLevel: 10, PerkDesc: "Increases stamina regeneration (+5 per cycle)" },
	"forge":            { Name: "⚒️ Forge", Ingredients: map[string]int{"stone": 200, "coal": 50, "iron": 20}, RequiredLevel: 20, PerkDesc: "Increases attack power (+10 attack)" },
	"enchanting_table": { Name: "🧪 Enchanting Table", Ingredients: map[string]int{"gold": 50, "abyss_crystal": 10, "diamond": 5}, RequiredLevel: 30, PerkDesc: "Arcane Wisdom (+50% XP gain from all sources)" },
	"vault":            { Name: "🏦 Vault", Ingredients: map[string]int{"iron": 100, "gold": 20}, RequiredLevel: 40, PerkDesc: "Increases max health (+50 HP)" },
	"castle":           { Name: "🏰 Castle", Ingredients: map[string]int{"stone": 1000, "iron": 200, "diamond": 50}, RequiredLevel: 50, PerkDesc: "Global Mastery (+20 attack, +100 HP, +50 Stamina)" },
}

var BotSettlements = map[string]BotSettlement{
	"goblin_camp": { Name: "👺 Goblin Camp", Level: 10, RequiredSword: "iron_sword", Defenders: []Monster{{Name: "👺 Goblin Warrior", Health: 80, Damage: 20}, {Name: "👺 Goblin Archer", Health: 50, Damage: 25}}, LootTable: map[string]int{"wood": 30, "stone": 20, "gold": 25}, Description: "A camp of goblins." },
	"bandit_fort": { Name: "🏴‍☠️ Bandit Fort", Level: 20, RequiredSword: "diamond_sword", Defenders: []Monster{{Name: "🏴‍☠️ Bandit Thug", Health: 250, Damage: 50}, {Name: "🏴‍☠️ Bandit Leader", Health: 500, Damage: 80}}, LootTable: map[string]int{"iron": 100, "gold": 100, "diamond": 5}, Description: "A dangerous fort." },
	"shadow_keep": { Name: "🏰 Shadow Keep", Level: 40, RequiredSword: "void_sword", Defenders: []Monster{{Name: "👻 Shadow Knight", Health: 800, Damage: 150}, {Name: "👻 Shadow Mage", Health: 500, Damage: 200}}, LootTable: map[string]int{"diamond": 50, "netherite": 20, "void_essence": 10}, Description: "The heart of darkness." },
}

var MerchantInventory = map[string]ShopItem{
	"iron_sword":    { ID: "iron_sword", Name: "⚔️ Iron Sword (+15 ATK)", Price: 200, Desc: "Reliable steel. (Lvl 10+)" },
	"diamond_sword": { ID: "diamond_sword", Name: "💎 Diamond Sword (+35 ATK)", Price: 1000, Desc: "Unyielding sharpness. (Lvl 20+)" },
	"void_sword":    { ID: "void_sword", Name: "🌌 Void Sword (+250 ATK)", Price: 5000, Desc: "A blade from beyond. (Lvl 40+)" },
	"leather_armor": { ID: "leather_armor", Name: "🛡️ Leather Armor (+5 DEF)", Price: 150, Desc: "Basic protection. (Lvl 10+)" },
	"plate_armor":   { ID: "plate_armor", Name: "🛡️ Plate Armor (+20 DEF)", Price: 800, Desc: "Solid heavy steel. (Lvl 20+)" },
	"dragon_scale":  { ID: "dragon_scale", Name: "🛡️ Dragon Scale (+100 DEF)", Price: 4000, Desc: "The peak of defense. (Lvl 40+)" },
	"health_potion": { ID: "health_potion", Name: "🧪 Health Potion (+50 HP)", Price: 30, Desc: "Minor healing." },
	"mega_health":   { ID: "mega_health", Name: "🧪 Mega Health (+200 HP)", Price: 150, Desc: "Significant healing. (Lvl 20+)" },
	"magic_potion":  { ID: "magic_potion", Name: "🔮 Magic Potion (+50 MP)", Price: 50, Desc: "Restores mana." },
	"elixir":        { ID: "elixir", Name: "🔮 Divine Elixir (Full Restore)", Price: 500, Desc: "Ultimate restoration. (Lvl 40+)" },
	"repair_kit":    { ID: "repair_kit", Name: "🔧 Repair Kit", Price: 50, Desc: "Restores tool durability." },
	"life_stone":    { ID: "life_stone", Name: "💎 Life Stone", Price: 10000, Desc: "Cheat death once." },
}

var ElementalSkillShopInventory = map[string]struct{ID string; Name string; Price int; Desc string}{
	"fire_bolt":   {ID: "fire_bolt", Name: "🔥 Fire Bolt", Price: 500, Desc: "Basic fire magic."},
	"water_jet":    {ID: "water_jet", Name: "💧 Water Jet", Price: 500, Desc: "Basic water magic."},
	"spark":        {ID: "spark", Name: "⚡ Spark", Price: 500, Desc: "Basic lightning magic."},
	"pebble_shot":  {ID: "pebble_shot", Name: "🪨 Pebble Shot", Price: 500, Desc: "Basic earth magic."},
	"breeze":       {ID: "breeze", Name: "🌬️ Breeze", Price: 500, Desc: "Basic wind magic."},
}

var TabooShopInventory = map[string]struct{ID string; Name string; Price int; Desc string}{
	"dark_attribute":          {ID: "dark_attribute", Name: "🌑 Dark Attribute", Price: 50, Desc: "Allows merging with high-tier elemental skills."},
	"decay_attribute":         {ID: "decay_attribute", Name: "🍄 Decay Attribute", Price: 50, Desc: "Unlocks the forbidden Rot skill tree."},
	"taboo_resonance":         {ID: "taboo_resonance", Name: "🌌 Taboo Resonance", Price: 50, Desc: "1% elemental dmg per Taboo Level."},
	"shadow_army_expansion":   {ID: "shadow_army_expansion", Name: "👥 Shadow Army Expansion", Price: 75, Desc: "+3 Max Skill Slots."},
	"soul_absorption":         {ID: "soul_absorption", Name: "👻 Soul Absorption", Price: 100, Desc: "Active: Drain 25% target HP."},
	"immortality":             {ID: "immortality", Name: "♾️ Immortality", Price: 150, Desc: "Passive: Survive lethal blow at 1 HP."},
	"shub_niggurath":          {ID: "shub_niggurath", Name: "🧬 Harvest Lord Shub-Niggurath", Price: 200, Desc: "Skill Creation mastery."},
}

var GlobalSkills = map[string]Skill{
	// --- Fire Tree ---
	"fire_bolt":       { ID: "fire_bolt", Name: "🔥 Fire Bolt", DmgBonus: 30, MPCost: 10, Rank: "E", Type: "active", Category: "attack", UnlockRequirement: "Elemental Shop" },
	"flame_lance":     { ID: "flame_lance", Name: "🔥 Flame Lance", DmgBonus: 100, MPCost: 40, Rank: "D", Type: "active", Category: "attack", UnlockRequirement: "Evolve Fire Bolt" },
	"prominence_burn": { ID: "prominence_burn", Name: "☀️ Prominence Burn", DmgBonus: 800, MPCost: 250, Rank: "A", Type: "active", Category: "attack", UnlockRequirement: "Evolve Flame Lance" },
	"inferno":         { ID: "inferno", Name: "🔥💀 Inferno", DmgBonus: 2500, MPCost: 800, Rank: "Ultimate", Type: "active", Category: "attack", UnlockRequirement: "Evolve Prominence Burn" },
	"hellfire":        { ID: "hellfire", Name: "🔥🌑 Hellfire", DmgBonus: 3500, MPCost: 1000, Rank: "Ultimate", Type: "active", Category: "attack", UnlockRequirement: "Merge Dark + Prominence Burn" },

	// --- Water Tree ---
	"water_jet":     { ID: "water_jet", Name: "💧 Water Jet", DmgBonus: 30, MPCost: 10, Rank: "E", Type: "active", Category: "attack", UnlockRequirement: "Elemental Shop" },
	"water_blade":   { ID: "water_blade", Name: "💧 Water Blade", DmgBonus: 90, MPCost: 35, Rank: "D", Type: "active", Category: "attack", UnlockRequirement: "Evolve Water Jet" },
	"tsunami":       { ID: "tsunami", Name: "🌊 Tsunami", DmgBonus: 700, MPCost: 200, Rank: "B", Type: "active", Category: "attack", UnlockRequirement: "Evolve Water Blade" },
	"oceanic_wrath": { ID: "oceanic_wrath", Name: "🔱 Oceanic Wrath", DmgBonus: 1500, MPCost: 450, Rank: "S", Type: "active", Category: "attack", UnlockRequirement: "Evolve Tsunami" },
	"abyss_tide":    { ID: "abyss_tide", Name: "🌊🌑 Abyss Tide", DmgBonus: 4000, MPCost: 1100, Rank: "Ultimate", Type: "active", Category: "attack", UnlockRequirement: "Merge Dark + Oceanic Wrath" },

	// --- Lightning Tree ---
	"spark":           { ID: "spark", Name: "⚡ Spark", DmgBonus: 25, MPCost: 8, Rank: "E", Type: "active", Category: "attack", UnlockRequirement: "Elemental Shop" },
	"chain_lightning": { ID: "chain_lightning", Name: "⚡ Chain Lightning", DmgBonus: 110, MPCost: 50, Rank: "D", Type: "active", Category: "attack", UnlockRequirement: "Evolve Spark" },
	"volt_strike":     { ID: "volt_strike", Name: "⚡ Volt Strike", DmgBonus: 600, MPCost: 180, Rank: "B", Type: "active", Category: "attack", UnlockRequirement: "Evolve Chain Lightning" },
	"indra_judgement": { ID: "indra_judgement", Name: "⚡⚖️ Indra's Judgement", DmgBonus: 2000, MPCost: 600, Rank: "S", Type: "active", Category: "attack", UnlockRequirement: "Evolve Volt Strike" },
	"black_lightning": { ID: "black_lightning", Name: "⚡🌑 Black Lightning", DmgBonus: 5000, MPCost: 1200, Rank: "Ultimate", Type: "active", Category: "attack", UnlockRequirement: "Merge Dark + Indra's Judgement" },

	// --- Earth Tree ---
	"pebble_shot":   { ID: "pebble_shot", Name: "🪨 Pebble Shot", DmgBonus: 35, MPCost: 12, Rank: "E", Type: "active", Category: "attack", UnlockRequirement: "Elemental Shop" },
	"earth_wall":    { ID: "earth_wall", Name: "🧱 Earth Wall", DmgBonus: 0, MPCost: 25, Rank: "D", Type: "active", Category: "defense", UnlockRequirement: "Evolve Pebble Shot" },
	"gravel_storm":  { ID: "gravel_storm", Name: "🌪️ Gravel Storm", DmgBonus: 400, MPCost: 150, Rank: "B", Type: "active", Category: "attack", UnlockRequirement: "Evolve Earth Wall" },
	"terraforming":  { ID: "terraforming", Name: "🌍 Terraforming", DmgBonus: 1200, MPCost: 500, Rank: "S", Type: "active", Category: "attack", UnlockRequirement: "Evolve Gravel Storm" },

	// --- Wind Tree ---
	"breeze":         { ID: "breeze", Name: "🌬️ Breeze", DmgBonus: 20, MPCost: 5, Rank: "E", Type: "active", Category: "attack", UnlockRequirement: "Elemental Shop" },
	"wind_cutter":    { ID: "wind_cutter", Name: "🌬️ Wind Cutter", DmgBonus: 80, MPCost: 30, Rank: "D", Type: "active", Category: "attack", UnlockRequirement: "Evolve Breeze" },
	"cyclone":        { ID: "cyclone", Name: "🌪️ Cyclone", DmgBonus: 500, MPCost: 160, Rank: "B", Type: "active", Category: "attack", UnlockRequirement: "Evolve Wind_cutter" },
	"tempest":        { ID: "tempest", Name: "🌬️🌪️ Tempest", DmgBonus: 1800, MPCost: 600, Rank: "S", Type: "active", Category: "attack", UnlockRequirement: "Evolve Cyclone" },

	// --- Warrior Tree ---
	"power_strike":   { ID: "power_strike", Name: "⚔️ Power Strike", DmgBonus: 50, MPCost: 15, Rank: "E", Type: "active", Category: "attack", UnlockRequirement: "Defeat Gate Boss" },
	"heavy_cleave":   { ID: "heavy_cleave", Name: "🪓 Heavy Cleave", DmgBonus: 120, MPCost: 30, Rank: "D", Type: "active", Category: "attack", UnlockRequirement: "Evolve Power Strike" },
	"armor_break":    { ID: "armor_break", Name: "🛡️ Armor Break", DmgBonus: 250, MPCost: 60, Rank: "C", Type: "active", Category: "attack", UnlockRequirement: "Evolve Heavy Cleave" },
	"earth_shatter":  { ID: "earth_shatter", Name: "🌍 Earth Shatter", DmgBonus: 600, MPCost: 150, Rank: "B", Type: "active", Category: "attack", UnlockRequirement: "Evolve Armor Break" },
	"meteor_strike":  { ID: "meteor_strike", Name: "☄️ Meteor Strike", DmgBonus: 1500, MPCost: 400, Rank: "A", Type: "active", Category: "attack", UnlockRequirement: "Evolve Earth Shatter" },
	"world_severing": { ID: "world_severing", Name: "🗡️ World-Severing Slash", DmgBonus: 4000, MPCost: 1000, Rank: "S", Type: "active", Category: "attack", UnlockRequirement: "Evolve Meteor Strike" },
	"void_slash":     { ID: "void_slash", Name: "🌑 Void Slash", DmgBonus: 8000, MPCost: 2500, Rank: "Ultimate", Type: "active", Category: "attack", UnlockRequirement: "Merge Dark + World-Severing" },

	// --- Spider (Kumoko) ---
	"venom_spit":     { ID: "venom_spit", Name: "🐍 Venom Spit", DmgBonus: 30, Rank: "D", Type: "active", Category: "attack", UnlockRequirement: "Spider Origin" },
	"poison_fang":    { ID: "poison_fang", Name: "🐍 Poison Fang", DmgBonus: 150, Rank: "C", Type: "active", Category: "attack", UnlockRequirement: "Evolve Venom Spit" },
	"deadly_venom":   { ID: "deadly_venom", Name: "☠️ Deadly Venom", DmgBonus: 600, Rank: "A", Type: "active", Category: "attack", UnlockRequirement: "Evolve Poison Fang" },
	"rot_attack":     { ID: "rot_attack", Name: "🍄 Rot Attack", DmgBonus: 6000, MPCost: 1500, Rank: "Forbidden", Type: "active", Category: "attack", UnlockRequirement: "Merge Decay + Deadly Venom" },
	"spider_thread":  { ID: "spider_thread", Name: "🕸️ Spider Thread", MPCost: 10, Rank: "Unique", Type: "active", Category: "defense", UnlockRequirement: "Spider Origin" },
	"steel_thread":   { ID: "steel_thread", Name: "🧶 Steel Thread", DmgBonus: 200, MPCost: 50, Rank: "A", Type: "active", Category: "attack", UnlockRequirement: "Evolve Spider Thread" },
	"divine_thread":  { ID: "divine_thread", Name: "✨ Divine Thread", DmgBonus: 1500, MPCost: 400, Rank: "Ultimate", Type: "active", Category: "attack", UnlockRequirement: "Evolve Steel Thread" },
	"appraisal":      { ID: "appraisal", Name: "👁️ Appraisal", Rank: "Unique", Type: "passive", Category: "utility", UnlockRequirement: "Spider Origin" },
	"evil_eye":       { ID: "evil_eye", Name: "🧿 Evil Eye of Statis", DmgBonus: 200, Rank: "B", Type: "active", Category: "attack", UnlockRequirement: "Spider Origin" },
	"abyss_magic":    { ID: "abyss_magic", Name: "🕳️ Abyss Magic", DmgBonus: 8000, Rank: "Ultimate", Type: "active", Category: "attack", UnlockRequirement: "Evolve Evil Eye" },
	"dim_maneuver":   { ID: "dim_maneuver", Name: "🌌 Dimensional Maneuver", MPCost: 150, Rank: "S", Type: "active", Category: "utility", UnlockRequirement: "Spider Evolution: Arachne" },
	"egg_revival":    { ID: "egg_revival", Name: "🥚 Egg Revivification", Rank: "Ultimate", Type: "passive", Category: "defense", UnlockRequirement: "Spider Evolution: God" },
	"sariel":         { ID: "sariel", Name: "🦉 Wisdom King Sariel", Rank: "Ultimate", Type: "passive", Category: "utility", UnlockRequirement: "Evolve Appraisal" },

	// --- Slime (Rimuru) ---
	"predator":       { ID: "predator", Name: "🌀 Predator", DmgBonus: 100, Rank: "Unique", Type: "active", Category: "attack", UnlockRequirement: "Slime Origin" },
	"gluttony":       { ID: "gluttony", Name: "👿 Gluttony", DmgBonus: 300, Rank: "Unique", Type: "active", Category: "attack", UnlockRequirement: "Evolve Predator" },
	"beelzebuth":     { ID: "beelzebuth", Name: "👹 Gluttonous King Beelzebuth", DmgBonus: 1200, Rank: "Ultimate", Type: "active", Category: "attack", UnlockRequirement: "Evolve Gluttony" },
	"great_sage":     { ID: "great_sage", Name: "🧠 Great Sage", Rank: "Unique", Type: "passive", Category: "utility", UnlockRequirement: "Slime Origin" },
	"raphael":        { ID: "raphael", Name: "📚 Wisdom King Raphael", Rank: "Ultimate", Type: "passive", Category: "utility", UnlockRequirement: "Evolve Great Sage" },
	"shub_niggurath": { ID: "shub_niggurath", Name: "🧬 Harvest Lord Shub-Niggurath", Rank: "Ultimate", Type: "passive", Category: "utility", UnlockRequirement: "Taboo Shop" },

	// --- General ---
	"natures_touch": { ID: "natures_touch", Name: "💚 Nature's Touch", MPCost: 20, Rank: "E", Type: "active", Category: "heal", UnlockRequirement: "Defeat Gate Boss" },
	"bone_armor":    { ID: "bone_armor", Name: "🦴 Bone Armor", MPCost: 20, Rank: "E", Type: "active", Category: "defense", UnlockRequirement: "Defeat Gate Boss" },
}

var SkillEvolutions = map[string]string{
	"fire_bolt": "flame_lance", "flame_lance": "prominence_burn", "prominence_burn": "inferno",
	"water_jet": "water_blade", "water_blade": "tsunami", "tsunami": "oceanic_wrath",
	"spark": "chain_lightning", "chain_lightning": "volt_strike", "volt_strike": "indra_judgement",
	"pebble_shot": "earth_wall", "earth_wall": "gravel_storm", "gravel_storm": "terraforming",
	"breeze": "wind_cutter", "wind_cutter": "cyclone", "cyclone": "tempest",
	"power_strike": "heavy_cleave", "heavy_cleave": "armor_break", "armor_break": "earth_shatter", "earth_shatter": "meteor_strike", "meteor_strike": "world_severing",
	"venom_spit": "poison_fang", "poison_fang": "deadly_venom",
	"spider_thread": "steel_thread", "steel_thread": "divine_thread",
	"appraisal": "sariel", "great_sage": "raphael", "predator": "gluttony", "gluttony": "beelzebuth",
	"evil_eye": "abyss_magic",
}

var GlobalQuests = []Quest{
	{ID: "wood_gatherer", Name: "🌲 Wood Gatherer", TargetType: "item", TargetID: "wood", TargetQty: 50, RewardXP: 100, RewardGold: 20, Description: "Collect 50 Wood."},
	{ID: "slime_hunter", Name: "🟢 Slime Hunter", TargetType: "combat", TargetID: "🟢 Slime", TargetQty: 5, RewardXP: 150, RewardGold: 30, Description: "Defeat 5 Slimes."},
	{ID: "iron_miner", Name: "⛓️ Iron Miner", TargetType: "item", TargetID: "iron", TargetQty: 30, RewardXP: 300, RewardGold: 100, Description: "Mine 30 Iron."},
	{ID: "zombie_slayer", Name: "🧟 Zombie Slayer", TargetType: "combat", TargetID: "🧟 Zombie", TargetQty: 10, RewardXP: 500, RewardGold: 200, Description: "Defeat 10 Zombies."},
	{ID: "mine_level_10", Name: "⛏️ Mine Adept", TargetType: "level", TargetID: "mine", TargetQty: 10, RewardXP: 1000, RewardGold: 500, Description: "Reach Mine Level 10."},
	{ID: "hunter_level_10", Name: "🏹 Novice Hunter", TargetType: "level", TargetID: "hunter", TargetQty: 10, RewardXP: 1000, RewardGold: 500, Description: "Reach Hunter Level 10."},
}

var GateBosses = map[string][]Monster{
	"E": {{Name: "👹 Hobgoblin", Health: 100, Damage: 15, LootTable: map[string]float64{"goblin_ear": 1.0}}, {Name: "🕷️ Small Lesser Taratect", Health: 80, Damage: 10, LootTable: map[string]float64{"string": 1.0}}},
	"D": {{Name: "🐺 Alpha Wolf", Health: 250, Damage: 30, LootTable: map[string]float64{"wolf_fang": 1.0}}, {Name: "🕷️ Spider Queen (Small)", Health: 300, Damage: 35, LootTable: map[string]float64{"spider_eye": 0.5, "string": 1.0}}},
	"C": {{Name: "🦎 Naga Warrior", Health: 600, Damage: 60, LootTable: map[string]float64{"naga_scale": 1.0}}, {Name: "🛡️ Iron Tyrant", Health: 800, Damage: 70, LootTable: map[string]float64{"iron_plate": 0.5}}},
	"B": {{Name: "🦾 Golem Guardian", Health: 1500, Damage: 120, LootTable: map[string]float64{"core_fragment": 1.0}}, {Name: "🎭 Clayman", Health: 2000, Damage: 150, LootTable: map[string]float64{"marionette_string": 0.5}}},
	"A": {{Name: "🔥 Inferno Drake", Health: 4000, Damage: 250, LootTable: map[string]float64{"dragon_heart": 1.0}}, {Name: "⚔️ Hinata Sakaguchi", Health: 5000, Damage: 300, LootTable: map[string]float64{"holy_sword": 0.1}}},
	"S": {{Name: "👑 Demon King", Health: 10000, Damage: 500, LootTable: map[string]float64{"demon_soul": 1.0}}, {Name: "👑 Shadow Monarch", Health: 15000, Damage: 600, LootTable: map[string]float64{"monarch_crown": 0.5}}},
	"SS": {{Name: "🌌 Void Sovereign", Health: 50000, Damage: 1500, LootTable: map[string]float64{"void_crown": 1.0}}, {Name: "🌀 Storm Dragon Veldora", Health: 100000, Damage: 2500, LootTable: map[string]float64{"storm_crest": 1.0}}},
}

var Gates = map[string]Gate{
	"E": { Rank: "E", Floors: 3, MinLevel: 10, MonsterCount: 5, RewardXP: 500, RewardGold: 100, Descriptions: []string{"A weak crack."} },
	"D": { Rank: "D", Floors: 4, MinLevel: 20, MonsterCount: 8, RewardXP: 1500, RewardGold: 500, Descriptions: []string{"A blue portal."} },
	"C": { Rank: "C", Floors: 5, MinLevel: 30, MonsterCount: 12, RewardXP: 5000, RewardGold: 2000, Descriptions: []string{"A green gate."} },
	"B": { Rank: "B", Floors: 6, MinLevel: 40, MonsterCount: 15, RewardXP: 15000, RewardGold: 8000, Descriptions: []string{"A purple portal."} },
	"A": { Rank: "A", Floors: 8, MinLevel: 50, MonsterCount: 20, RewardXP: 50000, RewardGold: 25000, Descriptions: []string{"A red gate."} },
	"S": { Rank: "S", Floors: 10, MinLevel: 100, MonsterCount: 30, RewardXP: 200000, RewardGold: 100000, Descriptions: []string{"A black gate."} },
	"SS": { Rank: "SS", Floors: 12, MinLevel: 150, MonsterCount: 50, RewardXP: 1000000, RewardGold: 500000, Descriptions: []string{"A god-like rift."} },
}

var GlobalTitles = map[string]Title{
	"wolf_slayer": { ID: "wolf_slayer", Name: "🐺 Wolf Slayer", KillsNeeded: 50, PerkDesc: "+10 Attack", AttackBonus: 10 },
	"goblin_bane": { ID: "goblin_bane", Name: "👺 Goblin Bane", KillsNeeded: 100, PerkDesc: "+20 Attack, +50 HP", AttackBonus: 20, HPBonus: 50 },
	"spider_crusher": { ID: "spider_crusher", Name: "🕷️ Spider Crusher", KillsNeeded: 150, PerkDesc: "+30 Defense, +100 MP", DefenseBonus: 30, MPBonus: 100 },
	"demon_hunter": { ID: "demon_hunter", Name: "😈 Demon Hunter", KillsNeeded: 250, PerkDesc: "+50 Attack, +500 Max HP", AttackBonus: 50, HPBonus: 500 },
	"dragon_slayer": { ID: "dragon_slayer", Name: "🐲 Dragon Slayer", KillsNeeded: 500, PerkDesc: "+100 Attack, +1000 Max HP", AttackBonus: 100, HPBonus: 1000 },
	"god_slayer": { ID: "god_slayer", Name: "✨ God Slayer", KillsNeeded: 1000, PerkDesc: "+250 Attack, +5000 Max HP", AttackBonus: 250, HPBonus: 5000 },
	"taboo_master": { ID: "taboo_master", Name: "🌌 Master of Taboo", KillsNeeded: 0, PerkDesc: "+500 MP, +50 Defense", DefenseBonus: 50, MPBonus: 500 },
	"taboo_priest": { ID: "taboo_priest", Name: "🕯️ High Priest of Taboo", KillsNeeded: 0, PerkDesc: "+1500 MP, +100 Attack", AttackBonus: 100, MPBonus: 1500 },
	"taboo_prophet": { ID: "taboo_prophet", Name: "🌑 Abyssal Prophet", KillsNeeded: 0, PerkDesc: "+5000 MP, +300 Attack", AttackBonus: 300, MPBonus: 5000 },
	"slime_emperor": { ID: "slime_emperor", Name: "💧 Slime Emperor", KillsNeeded: 0, PerkDesc: "+100 Attack, +1000 MP", AttackBonus: 100, MPBonus: 1000 },
	"labyrinth_walker": { ID: "labyrinth_walker", Name: "🕵️ Labyrinth Walker", KillsNeeded: 0, PerkDesc: "+50 Defense, +200 Stamina", DefenseBonus: 50, StaminaBonus: 200 },
	"world_conqueror": { ID: "world_conqueror", Name: "🌍 World Conqueror", KillsNeeded: 0, PerkDesc: "+500 Attack, +10000 HP", AttackBonus: 500, HPBonus: 10000 },
	"supreme_hunter": { ID: "supreme_hunter", Name: "🏹 Supreme Hunter", KillsNeeded: 0, PerkDesc: "+200 Attack, +500 Defense", AttackBonus: 200, DefenseBonus: 500 },
}
