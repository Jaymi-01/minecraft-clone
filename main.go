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
	
	var player *Player
	reader := bufio.NewReader(os.Stdin)

	// Onboarding Check
	if _, err := os.Stat("player_data.json"); os.IsNotExist(err) {
		fmt.Println("🌟 Welcome to your new adventure! 🌟")
		fmt.Print("Enter your character's name: ")
		name, _ := reader.ReadString('\n')
		name = strings.TrimSpace(name)
		if name == "" { name = "Adventurer" }
		player = NewPlayer(name)
		player.Save()
		fmt.Printf("Hello, %s! Your journey begins now.\n", name)
	} else {
		player = LoadPlayer()
		fmt.Printf("🌟 Welcome back, %s! 🌟\n", player.Name)
	}
	
	// Initialize Background Systems
	player.StartRegeneration()
	player.StartRaids()
	player.StartGateSpawning()
	StartServer(player) // Start the Web Dashboard

	fmt.Println("Available Commands: !mine <location>, !enter, !status, !craft [item], !equip <id>, !unequip <slot#>, !dupskill <id>, !learn <id>, !titles, !build [structure], !shop, !buy <item>, !use <item>, !raid [target], !quests, !subordinates, !name <species> <name>, !stats, !inventory, !origin <slime|spider>, !evolve, !help, !exit")

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
		
		// Handle Movement if Exploring
		if player.Exploring {
			cmdUpper := strings.ToUpper(command)
			if cmdUpper == "W" || cmdUpper == "A" || cmdUpper == "S" || cmdUpper == "D" {
				player.Move(cmdUpper)
				continue
			}
		}

		switch command {
		case "!explore":
			player.StartExploration()
		case "!emerge":
			player.Emerge()
		case "!name":
			if len(parts) < 3 {
				fmt.Println("🤝 Usage: !name <species> <given_name> (Costs 50 Max MP)")
			} else {
				player.NameSubordinate(parts[1], parts[2])
			}
		case "!subordinates":
			player.ListSubordinates()
		case "!origin":
			if len(parts) < 2 {
				fmt.Println("🧬 Choose your System Origin: slime, spider (Unlocked at Level 5)")
			} else {
				player.ChooseOrigin(parts[1])
			}
		case "!evolve":
			player.Evolve()
		case "!mine":
			if len(parts) < 2 {
				fmt.Println("📍 Available Locations: surface, cave, abyss, nether, void")
				fmt.Println("Usage: !mine <location>")
			} else {
				player.Mine(parts[1])
			}
		case "!enter":
			player.EnterGate()
		case "!status", "!stats", "!s", "!id":
			player.ShowStats()
		case "!equip":
			if len(parts) < 2 {
				player.ListSkills()
				fmt.Println("🔮 Usage: !equip <skill_id>")
			} else {
				player.EquipSkill(parts[1])
			}
		case "!unequip":
			if len(parts) < 2 {
				player.ListSkills()
				fmt.Println("🔮 Usage: !unequip <slot_number>")
			} else {
				var slot int
				fmt.Sscanf(parts[1], "%d", &slot)
				player.UnequipSkill(slot)
			}
		case "!dupskill":
			if len(parts) < 2 {
				player.ListSkills()
				fmt.Println("🔮 Usage: !dupskill <skill_id>")
			} else {
				player.UpgradeSkill(parts[1])
			}
		case "!skills":
			player.ListSkills()
		case "!learn":
			if len(parts) < 2 {
				player.ListSkills()
				fmt.Println("🔮 Usage: !learn <skill_id>")
			} else {
				player.LearnSkill(parts[1])
			}
		case "!titles":
			player.ListTitles()
		case "!craft":
			if len(parts) < 2 {
				player.ListDCraftable() // dcraft consolidated into craft
			} else {
				player.Craft(parts[1])
			}
		case "!build":
			if len(parts) < 2 {
				player.ListBuildable()
			} else {
				player.Build(parts[1])
			}
		case "!shop":
			player.ListShop()
		case "!buy":
			if len(parts) < 2 {
				fmt.Println("⚖️ Usage: !buy <item_id>")
			} else {
				player.Buy(parts[1])
			}
		case "!use":
			if len(parts) < 2 {
				fmt.Println("🎒 Usage: !use <item>")
			} else {
				player.Use(parts[1])
			}
		case "!raid":
			if len(parts) < 2 {
				player.ListRaids()
			} else {
				player.Raid(parts[1])
			}
		case "!quests":
			player.ListQuests()
		case "!inventory", "!i":
			player.ShowInventory()
		case "!help", "!h":
			player.ShowHelp()
		case "!recover":
			player.HealFull()
			fmt.Println("⚡ [CHEAT] Health and Stamina fully replenished! ⚡")
		case "!exit":
			player.Save()
			fmt.Println("👋 Goodbye! Your progress has been saved to player_data.json.")
			return
		default:
			if strings.HasPrefix(command, "eval.giveuseritemvar") {
				payload := strings.TrimPrefix(command, "eval.giveuseritemvar")
				if strings.Contains(payload, "=") {
					kv := strings.Split(payload, "=")
					if len(kv) == 2 {
						itemID := kv[0]
						var qty int
						_, err := fmt.Sscanf(kv[1], "%d", &qty)
						if err == nil {
							player.Inventory[itemID] += qty
							fmt.Printf("⚡ [CHEAT] Added %d %s to your inventory! (Total: %d)\n", qty, itemID, player.Inventory[itemID])
							player.Save()
							continue
						}
					}
				}
			} else if strings.HasPrefix(command, "eval.giveusergatevar") {
				rank := strings.TrimPrefix(command, "eval.giveusergatevar")
				if rank != "" {
					player.ManualSpawnGate(rank)
					continue
				}
			}
			fmt.Printf("❓ Unknown command: %s\n", command)
		}
	}
}
