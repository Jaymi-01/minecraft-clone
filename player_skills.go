package main

import (
	"fmt"
	"strings"
)

func (p *Player) UnequipSkill(slot int) {
	if slot < 1 || slot > len(p.EquippedSkills) { fmt.Println("❌ Invalid slot."); return }
	sID := p.EquippedSkills[slot-1]
	p.EquippedSkills = append(p.EquippedSkills[:slot-1], p.EquippedSkills[slot:]...)
	fmt.Printf("⚪ Unequipped %s.\n", GlobalSkills[sID].Name)
	p.Save()
}

func (p *Player) AddSkill(id string) {
	id = strings.ToLower(id); for _, s := range p.Skills { if s == id { return } }
	p.Skills = append(p.Skills, id); if s, ok := GlobalSkills[id]; ok { p.WorldNotice(fmt.Sprintf("SKILL ACQUIRED: %s", s.Name)) }
}

func (p *Player) HasSkill(id string) bool { id = strings.ToLower(id); for _, s := range p.Skills { if s == id { return true } }; return false }

func (p *Player) LearnSkill(skillID string) {
	skillID = strings.ToLower(skillID); skill, exists := GlobalSkills[skillID]; if !exists { return }
	for _, s := range p.Skills { if s == skillID { fmt.Println("❌ Already owned."); return } }
	
	met := false; reason := ""
	if skill.ReqBoss != "" { if p.MonsterKills[strings.ToLower(skill.ReqBoss)] > 0 { met = true } else { reason = "Defeat " + skill.ReqBoss } }
	if skill.ReqLevel > 0 { if p.Level >= skill.ReqLevel { met = true } else { reason = fmt.Sprintf("Lvl %d", skill.ReqLevel) } }
	if skill.Rank == "Forbidden" { if p.Taboo >= 10 { met = true } else { met = false; reason = "Taboo 10" } }
	
	req := strings.ToLower(skill.UnlockRequirement); origin := strings.ToLower(p.SystemOrigin)
	if strings.Contains(req, "origin") || strings.Contains(req, "evolution") { if strings.Contains(req, origin) { met = true } else { met = false; reason = skill.UnlockRequirement } }
	
	if skill.UnlockRequirement == "Defeat Gate Boss" {
		found := false; for k := range p.MonsterKills { for _, bosses := range GateBosses { for _, b := range bosses { if strings.EqualFold(k, b.Name) { found = true; break } }; if found { break } }; if found { break } }
		if found { met = true } else { met = false; reason = "Defeat any Gate Boss" }
	}
	if met { p.AddSkill(skillID); p.Save() } else { fmt.Printf("🚫 Locked: %s\n", reason) }
}

func (p *Player) UpgradeSkill(skillID string, isFree bool) {
	skillID = strings.ToLower(skillID); if !isFree && p.SkillPoints < 1 { return }
	if p.SkillLevels[skillID] == 0 { p.SkillLevels[skillID] = 1 }; p.SkillLevels[skillID]++
	if !isFree { p.SkillPoints-- }; p.SkillUsage[skillID] = 0
	p.WorldNotice(fmt.Sprintf("%s reached Lv%d", GlobalSkills[skillID].Name, p.SkillLevels[skillID]))
	if p.SkillLevels[skillID] >= 10 {
		if evolvedID, canEvolve := SkillEvolutions[skillID]; canEvolve {
			for i, s := range p.Skills { if s == skillID { p.Skills[i] = evolvedID; break } }
			for i, s := range p.EquippedSkills { if s == skillID { p.EquippedSkills[i] = evolvedID; break } }
			p.SkillLevels[evolvedID] = 1; p.SkillUsage[evolvedID] = 0
			p.WorldNotice(fmt.Sprintf("EVOLUTION: %s ascended to %s!", GlobalSkills[skillID].Name, GlobalSkills[evolvedID].Name))
		}
	}
	p.Save()
}

func (p *Player) EquipSkill(id string) {
	id = strings.ToLower(id); owned := false; for _, s := range p.Skills { if s == id { owned = true; break } }
	if !owned { fmt.Println("❌ Not owned."); return }
	for i, eq := range p.EquippedSkills { if eq == id { p.EquippedSkills = append(p.EquippedSkills[:i], p.EquippedSkills[i+1:]...); p.Save(); return } }
	if len(p.EquippedSkills) >= p.SkillSlots { fmt.Println("❌ No slots."); return }
	p.EquippedSkills = append(p.EquippedSkills, id); p.Save()
}

func (p *Player) MergeSkill(attr, target string) {
	attr = strings.ToLower(attr); target = strings.ToLower(target)
	if !p.Attributes[attr] { fmt.Printf("❌ Lack attribute: %s\n", attr); return }
	merges := map[string]map[string]string{"dark_attribute": {"prominence_burn": "hellfire", "indra_judgement": "black_lightning", "oceanic_wrath": "abyss_tide", "world_severing": "void_slash"}, "decay_attribute": {"deadly_venom": "rot_attack"}}
	res, ok := merges[attr][target]; if !ok || !p.HasSkill(target) { fmt.Println("❌ Incompatible merge."); return }
	for i, s := range p.Skills { if s == target { p.Skills[i] = res; break } }; for i, s := range p.EquippedSkills { if s == target { p.EquippedSkills[i] = res; break } }
	p.WorldNotice("MERGE SUCCESSFUL: " + GlobalSkills[res].Name); p.Save()
}

func (p *Player) DuplicateSkill(subName, skillID string) {
	if !p.HasSkill("shub_niggurath") { return }
	for _, s := range p.Subordinates { if strings.EqualFold(s.Name, subName) { for _, sk := range s.Skills { if strings.EqualFold(sk, skillID) { p.AddSkill(sk); p.Save(); return } } } }
}

func (p *Player) CreateSkill(s1, s2 string) {
	if !p.HasSkill("shub_niggurath") || p.Magic < 200 { return }
	p.Magic -= 200; p.AddSkill("fireball"); p.Save()
}
