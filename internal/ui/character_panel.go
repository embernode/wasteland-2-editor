package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/embernode/wasteland-2-editor/internal/savegame"
)

// CharacterPanel is the per-character editing area: a vitals header and
// an unspent-points header above four tabs (attributes, weapon skills,
// general skills, technical skills). Each tab is driven by the constant
// skill lists in package savegame.
type CharacterPanel struct {
	body fyne.CanvasObject

	vitals *VitalsBar
	points *PointsBar
	tabs   *container.AppTabs

	attributes *SkillTab
	weapons    *SkillTab
	general    *SkillTab
	technical  *SkillTab
}

// NewCharacterPanel constructs the editing panel in its empty (no save
// loaded) state.
func NewCharacterPanel() *CharacterPanel {
	p := &CharacterPanel{
		vitals:     newVitalsBar(),
		points:     newPointsBar(),
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
	header := container.NewVBox(
		p.vitals.body,
		p.points.body,
		widget.NewSeparator(),
	)
	p.body = container.NewBorder(header, nil, nil, nil, p.tabs)
	return p
}

// Container returns the Fyne object to embed in a parent layout.
func (p *CharacterPanel) Container() fyne.CanvasObject { return p.body }

// Show binds the headers and all tabs to a character. Pass nil to clear.
func (p *CharacterPanel) Show(c *savegame.Character) {
	if c == nil {
		p.vitals.clear()
		p.points.clear()
		p.attributes.clear()
		p.weapons.clear()
		p.general.clear()
		p.technical.clear()
		return
	}
	p.vitals.bind(c)
	p.points.bind(c)
	p.attributes.bind(c.Attributes)
	p.weapons.bind(c.Skills)
	p.general.bind(c.Skills)
	p.technical.bind(c.Skills)
}
