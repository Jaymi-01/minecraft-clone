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
		case "!skills", "!sk":
			player.ListSkills()
		case "!allskills":
			player.ListAllSystemSkills()
		case "!subordinates", "!subs":
			player.ListSubordinates()
		case "!titles":
			player.ListTitles()
		case "!quests", "!q":
			player.ListQuests()
		case "!use":
			if len(parts) < 2 { fmt.Println("🧪 Usage: !use <item_id>") } else { player.Use(parts[1]) }
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
			if len(parts) < 2 { fmt.Println("🔮 Usage: !unequip <slot_number|skill_id|weapon|armor>") } else {
				id := parts[1]
				var slot int
				if n, err := fmt.Sscanf(id, "%d", &slot); err == nil && n == 1 {
					player.UnequipSkill(slot)
				} else if id == "weapon" || id == "armor" {
					player.UnequipItem(id)
				} else {
					// Try to unequip by skill ID
					found := false
					for i, eq := range player.EquippedSkills {
						if eq == id {
							player.UnequipSkill(i + 1)
							found = true
							break
						}
					}
					if !found {
						// Maybe it's a weapon/armor ID they typed?
						if player.EquippedWeapon == id { player.UnequipItem("weapon") } else if player.EquippedArmor == id { player.UnequipItem("armor") } else {
							fmt.Printf("❌ [SYSTEM]: '%s' is not currently equipped in any slot.\n", id)
						}
					}
				}
			}
		case "!upgrade":
			if len(parts) < 2 { fmt.Println("🔮 Usage: !upgrade <skill_id>") } else { player.UpgradeSkill(parts[1], false) }
		case "!dupskill":
			if len(parts) < 3 { fmt.Println("🔮 Usage: !dupskill <subordinate_name> <skill_id>") } else { player.DuplicateSkill(parts[1], parts[2]) }
		case "!create":
			if len(parts) < 3 { fmt.Println("🧬 Usage: !create <skill_1> <skill_2>") } else { player.CreateSkill(parts[1], parts[2]) }
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
			if len(parts) < 2 { fmt.Println("👥 !squad <add|remove|list>") } else {
				if parts[1] == "add" && len(parts) > 2 { 
					player.AddToSquad(parts[2]) 
				} else if parts[1] == "remove" && len(parts) > 2 {
					player.RemoveFromSquad(parts[2])
				} else if parts[1] == "list" { 
					player.ListSquad() 
				}
			}
		case "!name":
			if len(parts) < 3 { fmt.Println("🤝 Usage: !name <species> <given_name>") } else { player.NameSubordinate(parts[1], parts[2]) }
		case "!help", "!h", "?":
			player.ShowHelp()
		case "!recover":
			player.HealFull()
			player.Stamina = player.MaxStamina
			player.WorldNotice("RESTORED: Vitals and Energy stabilized.")
		case "!exit", "!quit":
			player.Save()
			fmt.Println("System Offline. Progress Saved.")
			return
		default:
			if strings.HasPrefix(command, "eval.giveuseritemvar") {
				cmd := strings.TrimPrefix(command, "eval.giveuseritemvar")
				parts2 := strings.Split(cmd, "=")
				if len(parts2) == 2 {
					itemName := parts2[0]
					qty := 0; fmt.Sscanf(parts2[1], "%d", &qty)
					player.Inventory[itemName] += qty
					player.WorldNotice(fmt.Sprintf("EVAL: itemvar.%s = %d", itemName, player.Inventory[itemName]))
				}
			} else if strings.HasPrefix(command, "eval.giveusergatevar") {
				gateType := strings.ToUpper(strings.TrimPrefix(command, "eval.giveusergatevar"))
				if g, ok := Gates[gateType]; ok {
					newGate := g
					bosses := GateBosses[gateType]
					if len(bosses) > 0 {
						newGate.Boss = bosses[rand.Intn(len(bosses))]
					}
					player.CurrentGate = &newGate
					player.WorldNotice(fmt.Sprintf("EVAL: gatevar.manifest = %s (Boss: %s)", gateType, newGate.Boss.Name))
				}
			} else if strings.HasPrefix(command, "eval.giveusertaboovar") {
				qtyStr := strings.TrimPrefix(command, "eval.giveusertaboovar")
				qty := 0; fmt.Sscanf(qtyStr, "%d", &qty)
				player.GainTaboo(qty)
			} else {
				fmt.Printf("❓ Unknown command: %s. Type !help.\n", command)
			}

		}
	}
}
