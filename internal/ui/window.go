package ui

import (
	"fmt"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"

	"github.com/embernode/wasteland-2-editor/internal/savegame"
)

// BuildMainWindow assembles the editor window and wires the open / save /
// character-select controls into the model.
func BuildMainWindow(w fyne.Window) {
	model := &Model{}
	panel := NewCharacterPanel()

	pathLabel := widget.NewLabel("No save loaded.")
	pathLabel.Wrapping = fyne.TextTruncate

	charSelect := widget.NewSelect(nil, nil)
	charSelect.PlaceHolder = "(open a save first)"
	charSelect.Disable()
	charSelect.OnChanged = func(name string) {
		if model.SelectByDisplayName(name) {
			panel.Show(model.Current)
		}
	}

	saveBtn := widget.NewButton("Save", nil)
	saveBtn.Disable()

	loadPath := func(path string) {
		save, err := savegame.Load(path)
		if err != nil {
			dialog.ShowError(fmt.Errorf("load %s: %w", filepath.Base(path), err), w)
			return
		}
		model.SetSave(save)
		pathLabel.SetText(path)
		charSelect.Options = model.CharacterNames()
		charSelect.Refresh()
		if len(charSelect.Options) > 0 {
			charSelect.SetSelected(charSelect.Options[0])
			charSelect.Enable()
		}
		saveBtn.Enable()
		panel.Show(model.Current)
	}

	openBtn := widget.NewButton("Open save…", func() {
		fd := dialog.NewFileOpen(func(rc fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if rc == nil {
				return // user cancelled
			}
			path := rc.URI().Path()
			_ = rc.Close()
			loadPath(path)
		}, w)
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".xml"}))
		if model.Save != nil {
			if u, err := storage.ParseURI("file://" + filepath.Dir(model.Save.Path)); err == nil {
				if dirList, err := storage.ListerForURI(u); err == nil {
					fd.SetLocation(dirList)
				}
			}
		}
		fd.Resize(dialogSize(w))
		fd.Show()
	})

	saveBtn.OnTapped = func() {
		if model.Save == nil {
			return
		}
		err := model.Save.Save(model.Save.Path)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		dialog.ShowInformation("Saved",
			fmt.Sprintf("Wrote %s\n(backup: %s.back)",
				filepath.Base(model.Save.Path), filepath.Base(model.Save.Path)),
			w)
	}

	w.SetOnDropped(func(_ fyne.Position, uris []fyne.URI) {
		if path := pickXMLPath(uris); path != "" {
			loadPath(path)
		}
	})

	header := container.NewBorder(
		nil, nil,
		container.NewHBox(openBtn, saveBtn),
		nil,
		pathLabel,
	)
	top := container.NewVBox(header, charSelect, widget.NewSeparator())

	w.SetContent(container.NewBorder(top, nil, nil, nil, panel.Container()))
	w.Resize(fyne.NewSize(820, 640))
}

// pickXMLPath returns the path of the first dropped URI with a .xml extension,
// falling back to the first URI if none match. Returns "" if uris is empty.
func pickXMLPath(uris []fyne.URI) string {
	for _, u := range uris {
		if u != nil && strings.EqualFold(filepath.Ext(u.Path()), ".xml") {
			return u.Path()
		}
	}
	if len(uris) > 0 && uris[0] != nil {
		return uris[0].Path()
	}
	return ""
}

// dialogSize returns a reasonable default dialog size relative to the parent.
func dialogSize(w fyne.Window) fyne.Size {
	s := w.Canvas().Size()
	return fyne.NewSize(s.Width*0.9, s.Height*0.85)
}
