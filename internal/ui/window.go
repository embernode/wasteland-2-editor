package ui

import (
	"fmt"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"

	"github.com/embernode/wasteland-2-editor/internal/savegame"
)

// PreferenceKeyTheme is the fyne.Preferences key used to remember the
// active palette name across runs.
const PreferenceKeyTheme = "theme"

// BuildMainWindow assembles the editor window: a top action bar, a left
// character sidebar, and the per-character editing panel on the right.
func BuildMainWindow(w fyne.Window) {
	app := fyne.CurrentApp()
	model := &Model{}

	pathLabel := widget.NewLabel("No save loaded.")
	pathLabel.Wrapping = fyne.TextTruncate

	updatePathLabel := func() {
		if model.Save == nil {
			pathLabel.SetText("No save loaded.")
			return
		}
		text := filepath.Base(model.Save.Path)
		if model.IsDirty() {
			text += " •"
		}
		pathLabel.SetText(text)
	}
	model.OnDirtyChange = updatePathLabel

	panel := NewCharacterPanel(model.MarkDirty)

	themeSelect := widget.NewSelect(ThemeNames(), nil)
	themeSelect.SetSelected(app.Preferences().StringWithFallback(PreferenceKeyTheme, WastelandPalette.Name))
	themeSelect.OnChanged = func(name string) {
		app.Settings().SetTheme(ThemeByName(name))
		app.Preferences().SetString(PreferenceKeyTheme, name)
	}

	sidebar := newCharacterSidebar(func(name string) {
		if model.SelectByDisplayName(name) {
			panel.Show(model.Current)
		}
	})

	saveBtn := widget.NewButton("Save", nil)
	saveBtn.Disable()

	loadPath := func(path string) {
		save, err := savegame.Load(path)
		if err != nil {
			dialog.ShowError(fmt.Errorf("load %s: %w", filepath.Base(path), err), w)
			return
		}
		model.SetSave(save) // also clears dirty + fires updatePathLabel
		sidebar.SetOptions(model.CharacterNames())
		saveBtn.Enable()
		panel.Show(model.Current)
	}

	// confirmIfDirty wraps any operation that would discard unsaved edits.
	confirmIfDirty := func(action string, then func()) {
		if !model.IsDirty() {
			then()
			return
		}
		dialog.ShowConfirm(
			"Discard unsaved changes?",
			fmt.Sprintf("You have unsaved edits. %s anyway?", action),
			func(ok bool) {
				if ok {
					then()
				}
			},
			w,
		)
	}

	showOpenDialog := func() {
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
	}

	openBtn := widget.NewButton("Open save…", func() {
		confirmIfDirty("Open a different save", showOpenDialog)
	})

	saveBtn.OnTapped = func() {
		if model.Save == nil {
			return
		}
		if err := model.Save.WriteTo(model.Save.Path); err != nil {
			dialog.ShowError(err, w)
			return
		}
		model.ClearDirty() // fires updatePathLabel; no modal toast
	}

	w.SetOnDropped(func(_ fyne.Position, uris []fyne.URI) {
		path := pickXMLPath(uris)
		if path == "" {
			return
		}
		confirmIfDirty("Open the dropped file", func() { loadPath(path) })
	})

	w.SetCloseIntercept(func() {
		if !model.IsDirty() {
			w.Close()
			return
		}
		dialog.ShowConfirm(
			"Unsaved changes",
			"You have unsaved edits. Quit anyway?",
			func(ok bool) {
				if ok {
					w.Close()
				}
			},
			w,
		)
	})

	// Cmd/Ctrl-O = Open save…  ;  Cmd/Ctrl-S = Save (when enabled).
	w.Canvas().AddShortcut(
		&desktop.CustomShortcut{KeyName: fyne.KeyO, Modifier: fyne.KeyModifierShortcutDefault},
		func(fyne.Shortcut) { openBtn.OnTapped() },
	)
	w.Canvas().AddShortcut(
		&desktop.CustomShortcut{KeyName: fyne.KeyS, Modifier: fyne.KeyModifierShortcutDefault},
		func(fyne.Shortcut) {
			if !saveBtn.Disabled() {
				saveBtn.OnTapped()
			}
		},
	)

	header := container.NewBorder(
		nil, nil,
		container.NewHBox(openBtn, saveBtn),
		container.NewHBox(widget.NewLabel("Theme:"), themeSelect),
		pathLabel,
	)

	w.SetContent(container.NewBorder(
		container.NewPadded(header), nil,
		container.NewPadded(sidebar.Container()), nil,
		container.NewPadded(panel.Container()),
	))
	w.Resize(fyne.NewSize(960, 640))
}

// pickXMLPath returns the path of the first dropped URI with a .xml extension.
// Non-XML drops are ignored — silently dropping is friendlier than parsing
// random files and showing a confusing error.
func pickXMLPath(uris []fyne.URI) string {
	for _, u := range uris {
		if u != nil && strings.EqualFold(filepath.Ext(u.Path()), ".xml") {
			return u.Path()
		}
	}
	return ""
}

// dialogSize returns a reasonable default dialog size relative to the parent.
func dialogSize(w fyne.Window) fyne.Size {
	s := w.Canvas().Size()
	return fyne.NewSize(s.Width*0.9, s.Height*0.85)
}
