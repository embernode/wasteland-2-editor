package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

// SkillTab is a tab body holding a column of SkillRow widgets driven by a
// fixed list of skill keys.
type SkillTab struct {
	rows []*SkillRow
	body fyne.CanvasObject
}

func newSkillTab(keys []string, scale scaleKind, min, max int, onEdit func()) *SkillTab {
	t := &SkillTab{rows: make([]*SkillRow, 0, len(keys))}
	rowObjs := make([]fyne.CanvasObject, 0, len(keys))
	for _, k := range keys {
		r := newSkillRow(k, scale, min, max, onEdit)
		t.rows = append(t.rows, r)
		rowObjs = append(rowObjs, r.container())
	}
	t.body = container.NewVScroll(container.NewVBox(rowObjs...))
	return t
}

func (t *SkillTab) bind(values map[string]int) {
	for _, r := range t.rows {
		r.bind(values)
	}
}

func (t *SkillTab) clear() {
	for _, r := range t.rows {
		r.clear()
	}
}
