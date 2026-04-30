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
	StartServer(player) // Start the Web Dashboard

	fmt.Println("Available Commands: !mine <location>, !craft [item], !build [structure], !shop, !buy <item>, !use <item>, !raid [target], !quests, !stats, !inventory, !exit")

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
		case "!stats":
			player.ShowStats()
		case "!inventory":
			player.ShowInventory()
		case "!recover":
			player.HealFull()
			fmt.Println("⚡ [CHEAT] Health and Stamina fully replenished! ⚡")
		case "!exit":
			player.Save()
			fmt.Println("👋 Goodbye! Your progress has been saved to player_data.json.")
			return
		default:
			fmt.Printf("❓ Unknown command: %s\n", command)
		}
	}
}
