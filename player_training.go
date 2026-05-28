package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func (p *Player) CheckTrainingReset() {
	now := time.Now()
	// If it's a new day since last reset
	if now.Year() > p.Training.LastReset.Year() || now.YearDay() > p.Training.LastReset.YearDay() {
		// If they didn't complete training yesterday, trigger penalty
		if !p.TrainingCompletedToday && (p.Training.Pushups > 0 || p.Training.Situps > 0 || p.Training.Squats > 0 || p.Training.Running > 0) {
			p.TriggerPenalty()
		}
		
		// Reset for the new day
		p.Training.Pushups = 0
		p.Training.Situps = 0
		p.Training.Squats = 0
		p.Training.Running = 0
		p.Training.LastReset = now
		p.TrainingCompletedToday = false
		p.Save()
		p.WorldNotice("DAILY QUEST: [PREPARATION TO BECOME STRONG] HAS ARRIVED.")
	}
}

func (p *Player) Train(exercise string) {
	if p.PenaltyActive { fmt.Println("🚫 [SYSTEM]: Training unavailable while in the Penalty Zone."); return }
	if p.TrainingCompletedToday { fmt.Println("✅ [SYSTEM]: Daily training already completed. Limits reached for today."); return }

	exercise = strings.ToLower(exercise)
	target := 100
	msg := ""

	switch exercise {
	case "pushups":
		if p.Training.Pushups >= target { fmt.Println("✅ Pushups complete."); return }
		p.Training.Pushups += 10
		msg = fmt.Sprintf("💪 Pushups: %d/%d", p.Training.Pushups, target)
	case "situps":
		if p.Training.Situps >= target { fmt.Println("✅ Situps complete."); return }
		p.Training.Situps += 10
		msg = fmt.Sprintf("🧘 Situps: %d/%d", p.Training.Situps, target)
	case "squats":
		if p.Training.Squats >= target { fmt.Println("✅ Squats complete."); return }
		p.Training.Squats += 10
		msg = fmt.Sprintf("🦵 Squats: %d/%d", p.Training.Squats, target)
	case "running":
		if p.Training.Running >= 10 { fmt.Println("✅ Running complete."); return }
		p.Training.Running += 1
		msg = fmt.Sprintf("🏃 Running: %d/10 km", p.Training.Running)
	default:
		fmt.Println("❓ Usage: !train <pushups|situps|squats|running>")
		return
	}

	fmt.Println(msg)
	if p.Training.Pushups >= target && p.Training.Situps >= target && p.Training.Squats >= target && p.Training.Running >= 10 {
		p.CompleteDailyTraining()
	}
	p.Save()
}

func (p *Player) CompleteDailyTraining() {
	p.TrainingCompletedToday = true
	p.WorldNotice("DAILY QUEST COMPLETED: Limits Surpassed.")
	
	// Rewards
	stats := []string{"atk", "def", "hp", "mp", "stamina"}
	reward := stats[rand.Intn(len(stats))]
	switch reward {
	case "atk": p.Attack += 1; fmt.Println("✨ Reward: Permanent +1 ATK")
	case "def": p.Defense += 1; fmt.Println("✨ Reward: Permanent +1 DEF")
	case "hp": p.MaxHealth += 10; fmt.Println("✨ Reward: Permanent +10 Max HP")
	case "mp": p.MaxMagic += 10; fmt.Println("✨ Reward: Permanent +10 Max MP")
	case "stamina": p.MaxStamina += 5; fmt.Println("✨ Reward: Permanent +5 Max Stamina")
	}

	p.Inventory["hidden_box"]++
	fmt.Println("🎁 Reward: [Hidden Box] acquired. Use !use hidden_box to open.")
	p.SyncStats()
}

func (p *Player) TriggerPenalty() {
	p.PenaltyActive = true
	p.WorldNotice("PENALTY QUEST: SURVIVAL OF THE UNPREPARED.")
	fmt.Println("⚠️ [SYSTEM]: You failed to complete your daily quest. You are being force-transported to the Penalty Zone.")
	
	m := Monster{
		Name: "🦂 Poisonous Giant Centipede",
		Health: 500,
		Damage: 50,
		LootTable: map[string]float64{"sand_core": 0.5},
	}

	fmt.Println("\n🌵 [PENALTY ZONE]: You must survive for 4 combat rounds against the desert swarm.")
	for i := 1; i <= 4; i++ {
		fmt.Printf("\n--- [SURVIVAL WAVE %d/4] ---\n", i)
		if !p.Combat(&m, false) {
			p.PenaltyActive = false
			return // They died, combat handle death
		}
	}

	p.PenaltyActive = false
	p.WorldNotice("PENALTY CONCLUDED: You have returned from the desert abyss.")
	p.HealFull()
}

func (p *Player) OpenHiddenBox() {
	if p.Inventory["hidden_box"] <= 0 { return }
	p.Inventory["hidden_box"]--
	
	loot := rand.Intn(100)
	if loot < 50 {
		p.Inventory["gold"] += 500
		fmt.Println("💰 [HIDDEN BOX]: Found 500 Gold!")
	} else if loot < 80 {
		p.SkillPoints += 1
		fmt.Println("🔮 [HIDDEN BOX]: Found 1 Skill Point!")
	} else if loot < 95 {
		p.Inventory["life_stone"]++
		fmt.Println("💎 [HIDDEN BOX]: RARE! Found a Life Stone!")
	} else {
		p.AddSkill("natures_touch")
		fmt.Println("✨ [HIDDEN BOX]: LEGENDARY! Mastered [Nature's Touch]!")
	}
	p.Save()
}
