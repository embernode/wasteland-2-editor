package ui

import "github.com/embernode/wasteland-2-editor/internal/savegame"

// Model is the in-memory editor state shared between the window, the
// character dropdown, and the editing panel.
type Model struct {
	Save    *savegame.Save
	Current *savegame.Character
}

// SetSave installs a freshly-loaded save and selects its first character.
func (m *Model) SetSave(s *savegame.Save) {
	m.Save = s
	if len(s.Characters) > 0 {
		m.Current = s.Characters[0]
	} else {
		m.Current = nil
	}
}

// SelectByDisplayName picks the character whose DisplayName matches name and
// returns true if found.
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

// CharacterNames returns the display names for the dropdown.
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
