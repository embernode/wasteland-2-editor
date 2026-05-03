package ui

import "github.com/embernode/wasteland-2-editor/internal/savegame"

// Model is the in-memory editor state shared between the window, the
// character sidebar, and the editing panel.
type Model struct {
	Save    *savegame.Save
	Current *savegame.Character

	// OnDirtyChange fires whenever IsDirty transitions; nil-safe.
	OnDirtyChange func()

	dirty bool
}

// SetSave installs a freshly-loaded save and selects its first character.
// Clears the dirty bit since the model is now in sync with disk.
func (m *Model) SetSave(s *savegame.Save) {
	m.Save = s
	if len(s.Characters) > 0 {
		m.Current = s.Characters[0]
	} else {
		m.Current = nil
	}
	m.ClearDirty()
}

// SelectByDisplayName picks the character whose DisplayName matches name and
// returns true if found. Does not affect dirty state — switching views isn't
// an edit.
func (m *Model) SelectByDisplayName(name string) bool {
	if m.Save == nil {
		return false
	}
	for _, c := range m.Save.Characters {
		if c.DisplayName == name {
			m.Current = c
			return true
		}
	}
	return false
}

// CharacterNames returns the display names for the sidebar.
func (m *Model) CharacterNames() []string {
	if m.Save == nil {
		return nil
	}
	out := make([]string, 0, len(m.Save.Characters))
	for _, c := range m.Save.Characters {
		out = append(out, c.DisplayName)
	}
	return out
}

// IsDirty reports whether the model holds unsaved edits.
func (m *Model) IsDirty() bool { return m.dirty }

// MarkDirty flips the dirty bit on (idempotent). Safe to pass as a callback
// to every editing widget.
func (m *Model) MarkDirty() {
	if !m.dirty {
		m.dirty = true
		m.notifyDirty()
	}
}

// ClearDirty flips the dirty bit off. Called after a successful save or load.
func (m *Model) ClearDirty() {
	if m.dirty {
		m.dirty = false
		m.notifyDirty()
	}
}

func (m *Model) notifyDirty() {
	if m.OnDirtyChange != nil {
		m.OnDirtyChange()
	}
}
