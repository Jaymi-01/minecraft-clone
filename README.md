# 🌟 Mine & Exploration RPG 🌟

A feature-rich, text-based RPG built in Go. Explore vast locations, build a thriving settlement, craft powerful gear, and raid NPC fortresses!

---

## 🎮 For Players (Just want to play!)

If you just want to jump into the adventure, you **don't** need to install Go or clone this repository. 

1. Navigate to the [**`release/`**](./release) folder in this repository.
2. Download the `.zip` file that matches your operating system:
   - `MineRPG_Windows.zip` (for PC)
   - `MineRPG_Linux.zip` (for Linux)
   - `MineRPG_Mac_M1_M2.zip` (for Apple Silicon Macs)
   - `MineRPG_Mac_Intel.zip` (for older Macs)
3. Extract the ZIP to a new folder.
4. Run the executable (`mine-system.exe` on Windows or `./mine-system` on Linux/Mac).
5. **Start your journey!** Your progress is saved automatically to `player_data.json`.
6. **📊 Visual Dashboard:** While the game is running, open your browser to [**`http://localhost:8080`**](http://localhost:8080) to see your real-time stats, XP, and inventory!

---

## 🛠️ For Developers (Want to contribute!)

We welcome contributors! To set up the development environment:

### Prerequisites
- [Go 1.20+](https://go.dev/dl/)

### Setup
1. **Clone the repository:**
   ```bash
   git clone https://github.com/Jaymi-01/minecraft-clone.git
   cd minecraft-clone
   ```
2. **Run the game locally:**
   ```bash
   go run .
   ```
3. **Build for your specific platform:**
   ```bash
   go build -o mine-game .
   ```

### Project Structure
- `main.go`: CLI Entry point and command loop.
- `player.go`: Core game logic (Mining, Combat, Crafting, Saving).
- `data.go`: Static definitions for locations, items, and quests.
- `types.go`: Struct and model definitions.
- `build.ps1`: Automation script for multi-platform packaging.

---

## 📜 Core Features
- **🌍 5 Exploration Zones:** From the Surface to the Cosmic Void.
- **🏗️ Building System:** Construct Houses, Farms, and Castles for passive perks.
- **⚔️ Bot Raiding:** Attack NPC settlements or defend your base from incoming raids.
- **📜 Quest System:** Track objectives and earn massive gold/XP rewards.
- **⚖️ Merchant Shop:** Spend your gold on rare consumables and repair kits.
- **💾 Persistence:** Seamless JSON-based save system.

Enjoy the grind! ⛏️💎🏰
