package ui

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/embernode/wasteland-2-editor/internal/savegame"
)

// VitalsBar is the top-most header showing the per-character vitals:
// editable Level and CurrentHP, plus read-only Age and Sex.
//
// Max HP is intentionally absent — Wasteland 2 does not persist it; the
// game derives max HP at runtime from attributes / level / luckyHitpoints.
type VitalsBar struct {
	level  *widget.Entry
	curHP  *widget.Entry
	age    *widget.Label
	sex    *widget.Label
	body   fyne.CanvasObject
	target *savegame.Character
}

func newVitalsBar() *VitalsBar {
	v := &VitalsBar{
		age: widget.NewLabel("—"),
		sex: widget.NewLabel("—"),
	}
	v.level = newPointsEntry(func(n int) {
		if v.target != nil {
			v.target.Level = n
		}
	})
	v.curHP = newPointsEntry(func(n int) {
		if v.target != nil {
			v.target.CurrentHP = n
		}
	})
	v.level.Disable()
	v.curHP.Disable()

	v.body = container.NewHBox(
		widget.NewLabel("Level:"), v.level,
		widget.NewLabel("Current HP:"), v.curHP,
		widget.NewSeparator(),
		widget.NewLabel("Age:"), v.age,
		widget.NewLabel("Sex:"), v.sex,
	)
	return v
}

func (v *VitalsBar) bind(c *savegame.Character) {
	v.target = nil // suppress write-back during initial SetText
	v.level.SetText(strconv.Itoa(c.Level))
	v.curHP.SetText(strconv.Itoa(c.CurrentHP))
	v.level.Enable()
	v.curHP.Enable()

	v.age.SetText(blankIfEmpty(c.Age))
	v.sex.SetText(blankIfEmpty(c.Gender))

	v.target = c
}

func (v *VitalsBar) clear() {
	v.target = nil
	v.level.SetText("")
	v.curHP.SetText("")
	v.level.Disable()
	v.curHP.Disable()
	v.age.SetText("—")
	v.sex.SetText("—")
}

func blankIfEmpty(s string) string {
	if s == "" {
		return "—"
	}
	return s
}
