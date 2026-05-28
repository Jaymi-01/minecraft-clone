package main

import (
	"fmt"
	"strings"
)

func (p *Player) UnequipSkill(slot int) {
	if slot < 1 || slot > len(p.EquippedSkills) { fmt.Println("❌ [SYSTEM]: Invalid slot selected."); return }
	sID := p.EquippedSkills[slot-1]
	p.EquippedSkills = append(p.EquippedSkills[:slot-1], p.EquippedSkills[slot:]...)
	fmt.Printf("⚪ [SYSTEM]: Skill [%s] has been returned to the library.\n", GlobalSkills[sID].Name)
	p.Save()
}

func (p *Player) AddSkill(id string) {
	id = strings.ToLower(id)
	for _, s := range p.Skills { if s == id { return } }
	p.Skills = append(p.Skills, id)
	if s, ok := GlobalSkills[id]; ok { 
		p.WorldNotice(fmt.Sprintf("SKILL MANIFESTATION: [%s] has been integrated into your neural network.", s.Name)) 
	}
}

func (p *Player) HasSkill(id string) bool { id = strings.ToLower(id); for _, s := range p.Skills { if s == id { return true } }; return false }

func (p *Player) LearnSkill(skillID string) {
	skillID = strings.ToLower(skillID); skill, exists := GlobalSkills[skillID]; if !exists { return }
	for _, s := range p.Skills { if s == skillID { fmt.Println("❌ [SYSTEM]: Mastery over this skill already exists."); return } }
	
	met := false; reason := ""
	if skill.ReqBoss != "" { if p.MonsterKills[strings.ToLower(skill.ReqBoss)] > 0 { met = true } else { reason = "Eliminate " + skill.ReqBoss } }
	if skill.ReqLevel > 0 { if p.Level >= skill.ReqLevel { met = true } else { reason = fmt.Sprintf("Energy Tier %d", skill.ReqLevel) } }
	if skill.Rank == "Forbidden" { if p.Taboo >= 10 { met = true } else { met = false; reason = "Taboo Threshold 10" } }
	
	req := strings.ToLower(skill.UnlockRequirement); origin := strings.ToLower(p.SystemOrigin)
	if strings.Contains(req, "origin") || strings.Contains(req, "evolution") { if strings.Contains(req, origin) { met = true } else { met = false; reason = skill.UnlockRequirement } }
	
	if skill.UnlockRequirement == "Defeat Gate Boss" {
		found := false; for k := range p.MonsterKills { for _, bosses := range GateBosses { for _, b := range bosses { if strings.EqualFold(k, b.Name) { found = true; break } }; if found { break } }; if found { break } }
		if found { met = true } else { met = false; reason = "Extract data from a Gate Boss" }
	}

	if met { 
		p.AddSkill(skillID); p.Save() 
	} else { 
		fmt.Printf("🚫 [SYSTEM]: Access Denied. Requirements: %s\n", reason) 
	}
}

func (p *Player) UpgradeSkill(skillID string, isFree bool) {
	skillID = strings.ToLower(skillID)
	if !isFree && p.SkillPoints < 1 { fmt.Println("❌ [SYSTEM]: Insufficient Skill Points."); return }
	
	if p.SkillLevels[skillID] == 0 { p.SkillLevels[skillID] = 1 }
	p.SkillLevels[skillID]++
	if !isFree { p.SkillPoints-- }
	p.SkillUsage[skillID] = 0
	
	p.WorldNotice(fmt.Sprintf("SKILL ENHANCEMENT: [%s] has reached Level %d. Efficiency increased.", GlobalSkills[skillID].Name, p.SkillLevels[skillID]))
	
	if p.SkillLevels[skillID] >= 10 {
		if evolvedID, canEvolve := SkillEvolutions[skillID]; canEvolve {
			// Atomic replacement
			for i, s := range p.Skills { if s == skillID { p.Skills[i] = evolvedID; break } }
			for i, s := range p.EquippedSkills { if s == skillID { p.EquippedSkills[i] = evolvedID; break } }
			p.SkillLevels[evolvedID] = 1; p.SkillUsage[evolvedID] = 0
			p.WorldNotice(fmt.Sprintf("SKILL EVOLUTION: [%s] has transcended into [%s]!", GlobalSkills[skillID].Name, GlobalSkills[evolvedID].Name))
		}
	}
	p.Save()
}

func (p *Player) EquipSkill(id string) {
	id = strings.ToLower(id)
	owned := false; for _, s := range p.Skills { if s == id { owned = true; break } }
	if !owned { fmt.Println("❌ [SYSTEM]: Skill data not found in local memory."); return }
	
	for i, eq := range p.EquippedSkills { 
		if eq == id { 
			p.EquippedSkills = append(p.EquippedSkills[:i], p.EquippedSkills[i+1:]...)
			fmt.Printf("⚪ [SYSTEM]: Skill [%s] de-synchronized.\n", GlobalSkills[id].Name)
			p.Save(); return 
		} 
	}
	
	if len(p.EquippedSkills) >= p.SkillSlots { fmt.Println("❌ [SYSTEM]: Neural capacity reached. Unequip a skill first."); return }
	p.EquippedSkills = append(p.EquippedSkills, id)
	fmt.Printf("✨ [SYSTEM]: Skill [%s] synchronized to active memory.\n", GlobalSkills[id].Name)
	p.Save()
}

func (p *Player) ListSkills() {
	fmt.Println("\n--- 🎮 [SYSTEM: SKILL ARCHIVE] ---")
	fmt.Printf("   TARGET: %s | SP: %d | CAPACITY: %d/%d\n", p.Name, p.SkillPoints, len(p.EquippedSkills), p.SkillSlots)
	fmt.Println("   --- ACTIVE SYNC ---")
	for i, sID := range p.EquippedSkills { 
		fmt.Printf("      [%d] %s (Lv.%d)\n", i+1, GlobalSkills[sID].Name, p.SkillLevels[sID]) 
	}
	
	owned := make(map[string]bool); for _, s := range p.Skills { owned[s] = true }
	fmt.Println("\n   --- UNLOCKED SKILLS ---")
	for _, sID := range p.Skills {
		isEquipped := false; for _, eq := range p.EquippedSkills { if eq == sID { isEquipped = true; break } }
		if !isEquipped { fmt.Printf("      [ ] %s (Lv.%d)\n", GlobalSkills[sID].Name, p.SkillLevels[sID]) }
	}
	fmt.Println("\n💡 [TIP]: Use '!allskills' to see the complete System database.")
}

func (p *Player) ListAllSystemSkills() {
	fmt.Println("\n--- 🌌 [GLOBAL SYSTEM ARCHIVE: ALL SKILLS] ---")
	owned := make(map[string]bool); for _, s := range p.Skills { owned[s] = true }
	
	categories := []string{"attack", "defense", "heal", "utility"}
	for _, cat := range categories {
		fmt.Printf("\n   --- %s ---\n", strings.ToUpper(cat))
		for id, s := range GlobalSkills {
			if s.Category == cat {
				status := "🔒 LOCKED"
				if owned[id] { status = "✅ UNLOCKED" }
				fmt.Printf("      - %s [%s] (%s)\n         📜 %s\n", s.Name, s.Rank, status, s.UnlockRequirement)
			}
		}
	}
}

func (p *Player) MergeSkill(attr, target string) {
	attr = strings.ToLower(attr); target = strings.ToLower(target)
	if !p.Attributes[attr] { fmt.Printf("❌ [SYSTEM]: Attribute '%s' not present in soul core.\n", attr); return }
	
	merges := map[string]map[string]string{
		"dark_attribute": {
			"prominence_burn": "hellfire", 
			"indra_judgement": "black_lightning", 
			"oceanic_wrath": "abyss_tide", 
			"world_severing": "void_slash",
		},
		"decay_attribute": {
			"deadly_venom": "rot_attack",
		},
	}
	
	res, ok := merges[attr][target]
	if !ok || !p.HasSkill(target) { fmt.Println("❌ [SYSTEM]: Merging sequence failed. Elements are incompatible."); return }
	
	for i, s := range p.Skills { if s == target { p.Skills[i] = res; break } }
	for i, s := range p.EquippedSkills { if s == target { p.EquippedSkills[i] = res; break } }
	
	p.WorldNotice(fmt.Sprintf("FORBIDDEN MERGE: [%s] and [%s] have fused to birth [%s]!", strings.ToUpper(attr), GlobalSkills[target].Name, GlobalSkills[res].Name))
	p.Save()
}

func (p *Player) DuplicateSkill(subName, skillID string) {
	if !p.HasSkill("shub_niggurath") { fmt.Println("❌ [SYSTEM]: Harvest Lord ability required for duplication."); return }
	for _, s := range p.Subordinates {
		if strings.EqualFold(s.Name, subName) {
			for _, sk := range s.Skills {
				if strings.EqualFold(sk, skillID) {
					p.AddSkill(sk)
					p.WorldNotice(fmt.Sprintf("HARVEST: Skill [%s] extracted from %s and duplicated.", GlobalSkills[sk].Name, subName))
					p.Save(); return
				}
			}
		}
	}
	fmt.Println("❌ [SYSTEM]: Targeted skill not found in subordinate memory.")
}

func (p *Player) CreateSkill(s1, s2 string) {
	if !p.HasSkill("shub_niggurath") || p.Magic < 200 { fmt.Println("❌ [SYSTEM]: Neural creation failed. Check Mana or Lord abilities."); return }
	p.Magic -= 200
	p.AddSkill("fireball")
	p.WorldNotice("CREATION: A new skill has been synthesized from raw mana.")
	p.Save()
}
