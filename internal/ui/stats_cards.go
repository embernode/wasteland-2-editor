package ui

import (
	"image/color"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/embernode/wasteland-2-editor/internal/savegame"
)

// themedBorder is a 1px outline that re-reads the theme's separator color
// on every Refresh(). Without it, swapping themes leaves the card borders
// stuck at the old palette's color.
type themedBorder struct{ widget.BaseWidget }

func newThemedBorder() *themedBorder {
	b := &themedBorder{}
	b.ExtendBaseWidget(b)
	return b
}

func (b *themedBorder) CreateRenderer() fyne.WidgetRenderer {
	rect := canvas.NewRectangle(color.Transparent)
	rect.StrokeWidth = 1
	rect.StrokeColor = theme.Color(theme.ColorNameSeparator)
	return &themedBorderRenderer{rect: rect}
}

type themedBorderRenderer struct{ rect *canvas.Rectangle }

func (r *themedBorderRenderer) Layout(s fyne.Size)        { r.rect.Resize(s) }
// MinSize returns 1×1 so the rectangle always paints — parents (Stack /
// container.NewBorder) override this with the sibling's actual size, but
// the defensive minimum keeps the widget visible if reused elsewhere.
func (r *themedBorderRenderer) MinSize() fyne.Size        { return fyne.NewSize(1, 1) }
func (r *themedBorderRenderer) Objects() []fyne.CanvasObject { return []fyne.CanvasObject{r.rect} }
func (r *themedBorderRenderer) Destroy()                  {}
func (r *themedBorderRenderer) Refresh() {
	r.rect.StrokeColor = theme.Color(theme.ColorNameSeparator)
	r.rect.Refresh()
}

// StatsCards is the per-character header that replaces the older vitals +
// points bars. It renders four editable stats as small bordered cards
// (Level, Current HP, Unspent Attribute Points, Unspent Skill Points) and
// a single read-only summary line below them (display name, sex, age).
type StatsCards struct {
	level    *widget.Entry
	curHP    *widget.Entry
	attrPts  *widget.Entry
	skillPts *widget.Entry
	summary  *widget.Label
	body     fyne.CanvasObject
	target   *savegame.Character
}

func newStatsCards() *StatsCards {
	s := &StatsCards{}
	s.level = newPointsEntry(func(n int) {
		if s.target != nil {
			s.target.Level = n
		}
	})
	s.curHP = newPointsEntry(func(n int) {
		if s.target != nil {
			s.target.CurrentHP = n
		}
	})
	s.attrPts = newPointsEntry(func(n int) {
		if s.target != nil {
			s.target.AvailableAttributePoints = n
		}
	})
	s.skillPts = newPointsEntry(func(n int) {
		if s.target != nil {
			s.target.AvailableSkillPoints = n
		}
	})
	s.disableAll()

	s.summary = widget.NewLabel("")

	row := container.NewGridWithColumns(4,
		card("LEVEL", s.level),
		card("CURRENT HP", s.curHP),
		card("ATTR POINTS", s.attrPts),
		card("SKILL POINTS", s.skillPts),
	)
	s.body = container.NewPadded(container.NewVBox(
		row,
		container.NewPadded(s.summary),
	))
	return s
}

// bind attaches the cards to a character's data.
func (s *StatsCards) bind(c *savegame.Character) {
	s.target = nil
	s.level.SetText(strconv.Itoa(c.Level))
	s.curHP.SetText(strconv.Itoa(c.CurrentHP))
	s.attrPts.SetText(strconv.Itoa(c.AvailableAttributePoints))
	s.skillPts.SetText(strconv.Itoa(c.AvailableSkillPoints))
	s.enableAll()
	s.summary.SetText(summaryLine(c))
	s.target = c
}

func (s *StatsCards) clear() {
	s.target = nil
	s.level.SetText("")
	s.curHP.SetText("")
	s.attrPts.SetText("")
	s.skillPts.SetText("")
	s.disableAll()
	s.summary.SetText("")
}

func (s *StatsCards) enableAll() {
	s.level.Enable()
	s.curHP.Enable()
	s.attrPts.Enable()
	s.skillPts.Enable()
}

func (s *StatsCards) disableAll() {
	s.level.Disable()
	s.curHP.Disable()
	s.attrPts.Disable()
	s.skillPts.Disable()
}

func summaryLine(c *savegame.Character) string {
	parts := []string{c.DisplayName}
	if c.Gender != "" {
		parts = append(parts, c.Gender)
	}
	if c.Age != "" {
		parts = append(parts, "Age "+c.Age)
	}
	out := ""
	for i, p := range parts {
		if i > 0 {
			out += "  •  "
		}
		out += p
	}
	return out
}

// card wraps content in a small bordered box with an uppercase caption.
func card(caption string, content fyne.CanvasObject) fyne.CanvasObject {
	label := widget.NewLabelWithStyle(caption, fyne.TextAlignCenter, fyne.TextStyle{Bold: true, Monospace: true})

	// Double-pad: NewPadded once gives us breathing room from the border,
	// the second wrap stops the entry from touching the caption.
	body := container.NewPadded(container.NewPadded(container.NewVBox(label, content)))
	return container.NewStack(newThemedBorder(), body)
}

// newPointsEntry returns an Entry that calls onCommit when the text parses
// as a non-negative integer. Garbage / empty input is ignored.
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
