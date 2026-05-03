package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"

	"github.com/embernode/wasteland-2-editor/internal/savegame"
)

// CharacterPanel is the per-character editing area: a stats card row at the
// top, then four tabs (attributes, weapon skills, general skills, technical
// skills). Tabs are driven by the constant skill lists in package savegame.
type CharacterPanel struct {
	body fyne.CanvasObject

	stats *StatsCards
	tabs  *container.AppTabs

	attributes *SkillTab
	weapons    *SkillTab
	general    *SkillTab
	technical  *SkillTab
}

// NewCharacterPanel constructs the editing panel in its empty (no save
// loaded) state.
func NewCharacterPanel() *CharacterPanel {
	p := &CharacterPanel{
		stats:      newStatsCards(),
		attributes: newSkillTab(savegame.Attributes, scaleDirect, 1, 10),
		weapons:    newSkillTab(savegame.WeaponSkills, scaleSkillXP, 0, 10),
		general:    newSkillTab(savegame.GeneralSkills, scaleSkillXP, 0, 10),
		technical:  newSkillTab(savegame.TechnicalSkills, scaleSkillXP, 0, 10),
	}
	p.tabs = container.NewAppTabs(
		container.NewTabItem("Attributes", p.attributes.body),
		container.NewTabItem("Weapons", p.weapons.body),
		container.NewTabItem("General", p.general.body),
		container.NewTabItem("Technical", p.technical.body),
	)
	p.body = container.NewBorder(p.stats.body, nil, nil, nil, p.tabs)
	return p
}

// Container returns the Fyne object to embed in a parent layout.
func (p *CharacterPanel) Container() fyne.CanvasObject { return p.body }

// Show binds the stats and all tabs to a character. Pass nil to clear.
func (p *CharacterPanel) Show(c *savegame.Character) {
	if c == nil {
		p.stats.clear()
		p.attributes.clear()
		p.weapons.clear()
		p.general.clear()
		p.technical.clear()
		return
	}
	p.stats.bind(c)
	p.attributes.bind(c.Attributes)
	p.weapons.bind(c.Skills)
	p.general.bind(c.Skills)
	p.technical.bind(c.Skills)
}
