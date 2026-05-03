package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// CharacterSidebar is a fixed-width vertical list of party members. Selecting
// a row fires the onChanged callback with the chosen DisplayName.
type CharacterSidebar struct {
	list      *widget.List
	options   []string
	selected  int // -1 = none
	onChanged func(string)
	body      fyne.CanvasObject
}

const sidebarWidth = 180

func newCharacterSidebar(onChanged func(string)) *CharacterSidebar {
	s := &CharacterSidebar{selected: -1, onChanged: onChanged}
	s.list = widget.NewList(
		func() int { return len(s.options) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(s.options[i])
		},
	)
	s.list.OnSelected = func(i widget.ListItemID) {
		if i == s.selected {
			return
		}
		s.selected = i
		if s.onChanged != nil && i >= 0 && i < len(s.options) {
			s.onChanged(s.options[i])
		}
	}
	// Constrain width; let the list fill height in its parent.
	s.body = container.New(&fixedWidth{width: sidebarWidth}, s.list)
	return s
}

// Container returns the embeddable widget.
func (s *CharacterSidebar) Container() fyne.CanvasObject { return s.body }

// SetOptions replaces the option list and selects the first entry.
func (s *CharacterSidebar) SetOptions(opts []string) {
	s.options = opts
	s.selected = -1
	s.list.Refresh()
	if len(opts) > 0 {
		s.selected = 0
		s.list.Select(0)
	}
}

// SetSelected highlights the row whose text matches name.
func (s *CharacterSidebar) SetSelected(name string) {
	for i, opt := range s.options {
		if opt == name {
			if i != s.selected {
				s.selected = i
				s.list.Select(i)
			}
			return
		}
	}
}

// fixedWidth is a layout that locks its only child to a fixed width and the
// available height. Used to keep the sidebar a stable column.
type fixedWidth struct{ width float32 }

func (f *fixedWidth) MinSize(objs []fyne.CanvasObject) fyne.Size {
	if len(objs) == 0 {
		return fyne.NewSize(f.width, 0)
	}
	min := objs[0].MinSize()
	return fyne.NewSize(f.width, min.Height)
}

func (f *fixedWidth) Layout(objs []fyne.CanvasObject, size fyne.Size) {
	for _, o := range objs {
		o.Resize(fyne.NewSize(f.width, size.Height))
		o.Move(fyne.NewPos(0, 0))
	}
}
