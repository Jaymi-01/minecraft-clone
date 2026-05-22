# 🌟 Mine & Exploration RPG 🌟

A feature-rich, text-based RPG built in Go with a real-time web dashboard. Combine mining mechanics with deep anime-inspired progression systems.

---

## 🎮 How to Play

### Just want to jump in?
1. Navigate to the [**`release/`**](./release) folder.
2. Download the `.zip` file for your OS.
3. Extract and run the executable.
4. **📊 Visual Dashboard:** While playing, open [**`http://localhost:8080`**](http://localhost:8080) in your browser to see your real-time stats and inventory!

### Developers
- **Prerequisites:** [Go 1.20+](https://go.dev/dl/)
- **Run:** `go run .`
- **Build:** `go build -o mine-rpg .`

---

## 📜 Core Game Systems

### ⛏️ Mining & Resources
Explore 5 distinct zones from the **Surface** to the **Void**. Each zone requires a minimum level and specific tools (Pickaxes) to harvest rare ores and crystals.

### 🌀 The Gate System (Dungeons)
Gates from **Rank E to SS** manifest throughout the world. Enter them to clear floors of monsters and defeat a **Gate Boss** for massive Gold and Hunter XP rewards.

### 🗺️ The Great Elroe Labyrinth
A deep exploration mode where you navigate through dangerous tunnels. The further you move (`ExplorationDepth`), the stronger the monsters become, but the better the hidden loot.

### 🛡️ Building & Settlement
Construct structures to gain permanent account-wide buffs:
- **Forge:** +10 Attack
- **Vault:** +50 Max HP
- **Enchanting Table:** +50% XP Gain
- **Castle:** Massive stats across the board.

---

## 🧬 Anime-Inspired Origins (Lvl 5+)

Once you reach Level 5, you can choose a **System Origin** to unlock unique skill trees and evolution paths:

- **💧 Slime Path:** Focuses on consumption and analysis.
  - *Skills:* Predator, Great Sage, Raphael, Beelzebuth.
  - *Evolution:* Demon Slime -> True Dragon.
- **🕷️ Spider Path:** Focuses on survival and status effects.
  - *Skills:* Appraisal, Spider Thread, Evil Eye, Parallel Minds.
  - *Evolution:* Arachne -> God (Shiraori).

---

## ⚔️ Combat & Progression
- **Active Skills:** Equip up to 3 active skills (expandable) to use in turn-based combat.
- **Titles:** Defeat large quantities of monsters to unlock Titles (e.g., *Wolf Slayer*) that grant permanent stat bonuses.
- **Subordinates:** Tame and **Name** monsters (Slimes, Goblins, Wolves) to have them assist you in combat and base defense.
- **Raiding:** Launch raids on NPC settlements like Goblin Camps or Bandit Forts to plunder resources.

---

## 💾 System Integrity
- **Real-time Sync:** The web dashboard updates every second with your latest stats.
- **Save Safety:** The system maintains a `player_data.json.bak` file to prevent progress loss in case of crashes.
- **Auto-Save:** Every major action (combat, mining, crafting) triggers an automatic save.

---

## ⌨️ Common Commands
- `!mine <zone>` - Start resource gathering.
- `!status` / `!s` - View full character profile.
- `!enter` - Challenge the currently manifested Gate.
- `!explore` - Enter the Great Elroe Labyrinth.
- `!origin <slime|spider>` - Choose your path at Level 5.
- `!evolve` - Advance to the next species tier.
- `!skills` - View your collection and unlock requirements.
- `!quests` - Check your mission progress.
- `!build <name>` - Construct a new settlement structure.

Enjoy the journey to the top of the System! 🗡️💎👑
