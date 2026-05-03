package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/embernode/wasteland-2-editor/internal/savegame"
)

// scaleKind tells the row whether the underlying map stores the displayed
// value directly (attributes) or a raw skill-XP value that needs conversion
// (skills).
type scaleKind int

const (
	scaleDirect scaleKind = iota // attribute: stored == displayed
	scaleSkillXP                  // skill: stored XP, displayed level
)

// SkillRow is a single label + slider + value-readout row.
type SkillRow struct {
	key     string
	scale   scaleKind
	min     int
	max     int
	level   *widget.Label
	slider  *widget.Slider
	target  map[string]int
	enabled bool
	onEdit  func()
}

func newSkillRow(key string, scale scaleKind, min, max int, onEdit func()) *SkillRow {
	r := &SkillRow{key: key, scale: scale, min: min, max: max, onEdit: onEdit}
	r.level = widget.NewLabelWithStyle(
		fmt.Sprintf("0/%d", max),
		fyne.TextAlignTrailing,
		fyne.TextStyle{Monospace: true},
	)
	r.slider = widget.NewSlider(float64(min), float64(max))
	r.slider.Step = 1
	r.slider.OnChanged = func(v float64) {
		level := int(v)
		r.level.SetText(fmt.Sprintf("%d/%d", level, max))
		if r.target == nil {
			return // bind() in progress; don't touch the model
		}
		if r.scale == scaleSkillXP {
			r.target[r.key] = savegame.SkillXP(level)
		} else {
			r.target[r.key] = level
		}
		if r.onEdit != nil {
			r.onEdit()
		}
	}
	r.slider.Disable()
	return r
}

// container returns a single-row layout: [label | slider | level/max].
func (r *SkillRow) container() fyne.CanvasObject {
	label := widget.NewLabel(savegame.SkillLabels[r.key])
	if label.Text == "" {
		label.SetText(r.key)
	}
	// Fixed-width label and value cells; slider takes the rest.
	left := container.NewGridWrap(fyne.NewSize(180, 32), label)
	right := container.NewGridWrap(fyne.NewSize(60, 32), r.level)
	return container.NewBorder(nil, nil, left, right, r.slider)
}

// bind attaches the row to a character's value map. The slider becomes
// interactive and immediately reflects the stored value.
func (r *SkillRow) bind(values map[string]int) {
	// Detach so we don't write to the previous character while we update.
	r.target = nil

	raw := values[r.key]
	level := raw
	if r.scale == scaleSkillXP {
		level = savegame.SkillLevel(raw)
	}
	r.slider.SetValue(float64(level))
	r.level.SetText(fmt.Sprintf("%d/%d", level, r.max))
	r.slider.Enable()
	r.enabled = true

	r.target = values
}

// clear puts the row back into a "no save loaded" state.
func (r *SkillRow) clear() {
	r.target = nil
	r.slider.SetValue(float64(r.min))
	r.slider.Disable()
	r.level.SetText(fmt.Sprintf("0/%d", r.max))
	r.enabled = false
}
