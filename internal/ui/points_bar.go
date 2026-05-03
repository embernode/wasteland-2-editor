package ui

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/embernode/wasteland-2-editor/internal/savegame"
)

// PointsBar is the header row that exposes the two unspent-point counters
// (attribute and skill) as plain numeric entries. Editing a slider does NOT
// adjust these — they're independent fields the player can tune freely.
type PointsBar struct {
	attr   *widget.Entry
	skill  *widget.Entry
	body   fyne.CanvasObject
	target *savegame.Character
}

func newPointsBar() *PointsBar {
	b := &PointsBar{}
	b.attr = newPointsEntry(func(n int) {
		if b.target != nil {
			b.target.AvailableAttributePoints = n
		}
	})
	b.skill = newPointsEntry(func(n int) {
		if b.target != nil {
			b.target.AvailableSkillPoints = n
		}
	})
	b.attr.Disable()
	b.skill.Disable()

	b.body = container.NewHBox(
		widget.NewLabel("Unspent attribute points:"), b.attr,
		widget.NewLabel("Unspent skill points:"), b.skill,
	)
	return b
}

// bind attaches the entries to a character. Setting the entry text to the
// character's current value will trigger OnChanged, which writes back the
// same value — harmless.
func (b *PointsBar) bind(c *savegame.Character) {
	b.target = nil // suppress write-back during initial SetText
	b.attr.SetText(strconv.Itoa(c.AvailableAttributePoints))
	b.skill.SetText(strconv.Itoa(c.AvailableSkillPoints))
	b.attr.Enable()
	b.skill.Enable()
	b.target = c
}

func (b *PointsBar) clear() {
	b.target = nil
	b.attr.SetText("")
	b.skill.SetText("")
	b.attr.Disable()
	b.skill.Disable()
}

// newPointsEntry returns a fixed-width Entry that calls onCommit when the
// text parses as a non-negative integer. Garbage input is ignored — the
// previous good value stays in effect.
func newPointsEntry(onCommit func(int)) *widget.Entry {
	e := widget.NewEntry()
	e.OnChanged = func(text string) {
		if text == "" {
			return
		}
		n, err := strconv.Atoi(text)
		if err != nil || n < 0 {
			return
		}
		onCommit(n)
	}
	return e
}
