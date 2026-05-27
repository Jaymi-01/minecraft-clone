package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	
	player := LoadPlayer()
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("🌟 SYSTEM INITIALIZED 🌟")
	player.WorldNotice("System Online. Welcome, " + player.Name)

	// Start Game Cycles
	player.StartRegeneration()
	player.StartRaids()
	player.StartGateSpawning()
	player.StartSubordinateAutonomy()
	StartServer(player)

	for {
		if player.Exploring {
			fmt.Printf("\n📍 [LABYRINTH DEPTH %d] Actions: W (Forward), A (Left), D (Right), S (Backward), !emerge\nChoice: ", player.ExplorationDepth)
		} else {
			fmt.Print("\nChoice: ")
		}

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		parts := strings.Fields(input)
		if len(parts) == 0 { continue }

		command := parts[0]

		if player.Exploring {
			switch strings.ToUpper(command) {
			case "W", "A", "S", "D": player.Move(command)
			case "!EMERGE": player.Emerge()
			default: fmt.Println("🚫 Use W/A/S/D or !emerge.")
			}
			continue
		}

		switch command {
		case "!mine":
			if len(parts) < 2 { fmt.Println("📍 Usage: !mine <location>") } else { player.Mine(parts[1]) }
		case "!enter":
			player.EnterGate(false)
		case "!allowadminenter":
			player.EnterGate(true)
		case "!status", "!stats", "!s", "!id":
			player.ShowStats()
		case "!inventory", "!i":
			player.ShowInventory()
		case "!equip":
			if len(parts) < 2 { fmt.Println("🔮 Usage: !equip <id>") } else {
				id := parts[1]
				isSkill := false
				for id2 := range GlobalSkills { if id2 == id { isSkill = true; break } }
				if isSkill { player.EquipSkill(id) } else { player.EquipItem(id) }
			}
		case "!unequip":
			if len(parts) < 2 { fmt.Println("🔮 Usage: !unequip <slot_number|weapon|armor>") } else {
				var slot int
				if n, err := fmt.Sscanf(parts[1], "%d", &slot); err == nil && n == 1 {
					player.UnequipSkill(slot)
				} else {
					player.UnequipItem(parts[1])
				}
			}
		case "!dupskill":
			if len(parts) < 2 { fmt.Println("🔮 Usage: !dupskill <id>") } else { player.UpgradeSkill(parts[1], false) }
		case "!learn":
			if len(parts) < 2 { fmt.Println("📖 Usage: !learn <id>") } else { player.LearnSkill(parts[1]) }
		case "!tabooshop":
			player.ListTabooShop()
		case "!buytaboo":
			if len(parts) < 2 { fmt.Println("🌌 Usage: !buytaboo <id>") } else { player.BuyTabooSkill(parts[1]) }
		case "!shop":
			player.ListShop()
		case "!buy":
			if len(parts) < 2 { fmt.Println("💰 Usage: !buy <id>") } else { player.Buy(parts[1]) }
		case "!elementalshop":
			player.ListElementalShop()
		case "!buyelemental":
			if len(parts) < 2 { fmt.Println("🔥 Usage: !buyelemental <id>") } else { player.BuyElementalSkill(parts[1]) }
		case "!merge":
			if len(parts) < 3 { fmt.Println("🧬 Usage: !merge <attr> <skill>") } else { player.MergeSkill(parts[1], parts[2]) }
		case "!origin":
			if len(parts) < 2 { fmt.Println("🧬 Usage: !origin <slime|spider>") } else { player.ChooseOrigin(parts[1]) }
		case "!evolve":
			player.Evolve()
		case "!explore":
			player.StartExploration()
		case "!craft":
			if len(parts) < 2 { player.ListDCraftable() } else { player.Craft(parts[1]) }
		case "!raid":
			if len(parts) < 2 { player.ListRaids() } else { player.Raid(parts[1]) }
		case "!squad":
			if len(parts) < 2 { fmt.Println("👥 !squad add <n>, !squad list") } else {
				if parts[1] == "add" && len(parts) > 2 { player.AddToSquad(parts[2]) } else if parts[1] == "list" { player.ListSquad() }
			}
		case "!help", "!h", "?":
			player.ShowHelp()
		case "!exit", "!quit":
			player.Save()
			fmt.Println("System Offline. Progress Saved.")
			return
		default:
			fmt.Printf("❓ Unknown command: %s. Type !help.\n", command)
		}
	}
}
