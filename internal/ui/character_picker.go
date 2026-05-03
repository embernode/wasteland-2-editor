package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// CharacterPicker is a dropdown-style control whose popup highlights the
// currently-selected option. widget.Select doesn't do that, so we roll our
// own as a Button-that-pops-a-List.
type CharacterPicker struct {
	button    *widget.Button
	options   []string
	selected  int // -1 = none
	onChanged func(string)
}

const pickerPlaceholder = "(open a save first)"

func newCharacterPicker(onChanged func(string)) *CharacterPicker {
	p := &CharacterPicker{selected: -1, onChanged: onChanged}
	p.button = widget.NewButtonWithIcon(pickerPlaceholder, theme.MenuDropDownIcon(), p.showPopup)
	p.button.IconPlacement = widget.ButtonIconTrailingText
	p.button.Alignment = widget.ButtonAlignLeading
	p.button.Disable()
	return p
}

// Container returns the embeddable widget.
func (p *CharacterPicker) Container() fyne.CanvasObject { return p.button }

// SetOptions replaces the option list. The current selection is cleared.
func (p *CharacterPicker) SetOptions(opts []string) {
	p.options = opts
	p.selected = -1
	if len(opts) == 0 {
		p.button.SetText(pickerPlaceholder)
		p.button.Disable()
		return
	}
	p.button.SetText(opts[0])
	p.selected = 0
}

// SetSelected picks the option whose text matches name. No-op if not found.
func (p *CharacterPicker) SetSelected(name string) {
	for i, opt := range p.options {
		if opt == name {
			p.selected = i
			p.button.SetText(opt)
			return
		}
	}
}

// Enable / Disable mirror widget.Select.
func (p *CharacterPicker) Enable()  { p.button.Enable() }
func (p *CharacterPicker) Disable() { p.button.Disable() }

func (p *CharacterPicker) showPopup() {
	if len(p.options) == 0 {
		return
	}
	canvas := fyne.CurrentApp().Driver().CanvasForObject(p.button)
	if canvas == nil {
		return
	}

	list := widget.NewList(
		func() int { return len(p.options) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(p.options[i])
		},
	)
	// Pre-select before wiring OnSelected so the highlight appears without
	// triggering the change handler.
	if p.selected >= 0 {
		list.Select(p.selected)
	}

	// Size the popup to the button's width and roughly one row per option,
	// capped so a large party doesn't overflow the window.
	btnSize := p.button.Size()
	const rowHeight = float32(36)
	const maxRows = 10
	rows := len(p.options)
	if rows > maxRows {
		rows = maxRows
	}
	popupSize := fyne.NewSize(btnSize.Width, rowHeight*float32(rows)+8)

	btnPos := fyne.CurrentApp().Driver().AbsolutePositionForObject(p.button)
	popupPos := fyne.NewPos(btnPos.X, btnPos.Y+btnSize.Height)

	popup := widget.NewPopUp(list, canvas)
	popup.Resize(popupSize)
	popup.ShowAtPosition(popupPos)

	list.OnSelected = func(i widget.ListItemID) {
		p.selected = i
		p.button.SetText(p.options[i])
		popup.Hide()
		if p.onChanged != nil {
			p.onChanged(p.options[i])
		}
	}
}
