package ui

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"fyne.io/fyne/v2/test"

	"github.com/embernode/wasteland-2-editor/internal/savegame"
)

// realSavePath walks up from the package dir looking for the test save.
// Skips the test if not found so CI without the asset stays green.
func realSavePath(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 6; i++ {
		candidate := filepath.Join(dir, "Quicksave 2", "Quicksave 2.xml")
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
		dir = filepath.Dir(dir)
	}
	t.Skip("Quicksave 2/Quicksave 2.xml not found")
	return ""
}

// TestCharacterPanel_BindsAndNotifiesEdits is a smoke test covering the most
// load-bearing UI invariants: bind populates entries, bind itself doesn't
// fire onEdit, and a real slider change does.
func TestCharacterPanel_BindsAndNotifiesEdits(t *testing.T) {
	test.NewApp() // initialize Fyne's test driver

	save, err := savegame.Load(realSavePath(t))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	hex := save.Characters[0]
	if hex.DisplayName != "Hex" {
		t.Fatalf("expected Hex first, got %q", hex.DisplayName)
	}

	var dirtyCalls int
	panel := NewCharacterPanel(func() { dirtyCalls++ })
	panel.Show(hex)

	// Stats entries should reflect Hex's vitals.
	if got, want := panel.stats.level.Text, strconv.Itoa(hex.Level); got != want {
		t.Errorf("Level entry = %q, want %q", got, want)
	}
	if got, want := panel.stats.curHP.Text, strconv.Itoa(hex.CurrentHP); got != want {
		t.Errorf("CurrentHP entry = %q, want %q", got, want)
	}
	if got, want := panel.stats.attrPts.Text, strconv.Itoa(hex.AvailableAttributePoints); got != want {
		t.Errorf("attr-points entry = %q, want %q", got, want)
	}

	// Binding should be silent — initial SetText / SetValue must not be
	// counted as user edits.
	if dirtyCalls != 0 {
		t.Errorf("onEdit fired %d times during bind, want 0", dirtyCalls)
	}

	// A user-driven slider change should fire onEdit and update the model.
	// Pick a target that's guaranteed to differ from the current value, so
	// SetValue doesn't early-return.
	row := panel.attributes.rows[0] // coordination
	current := int(row.slider.Value)
	target := current%10 + 1 // 1..10, always != current
	row.slider.SetValue(float64(target))
	if dirtyCalls == 0 {
		t.Errorf("onEdit not fired after slider change")
	}
	if got := hex.Attributes["coordination"]; got != target {
		t.Errorf("coordination = %d after slider change, want %d", got, target)
	}

	// Show(nil) must clear silently — the disable side-effect can't count
	// as an edit either.
	prev := dirtyCalls
	panel.Show(nil)
	if dirtyCalls != prev {
		t.Errorf("onEdit fired during Show(nil): %d -> %d", prev, dirtyCalls)
	}
}

// TestModel_DirtyTransitions covers MarkDirty / ClearDirty edge cases.
func TestModel_DirtyTransitions(t *testing.T) {
	var changes int
	m := &Model{OnDirtyChange: func() { changes++ }}

	if m.IsDirty() {
		t.Error("new Model should not be dirty")
	}

	m.MarkDirty()
	if !m.IsDirty() || changes != 1 {
		t.Errorf("MarkDirty: dirty=%v changes=%d", m.IsDirty(), changes)
	}

	// Idempotent — second MarkDirty shouldn't fire the callback again.
	m.MarkDirty()
	if changes != 1 {
		t.Errorf("MarkDirty fired callback while already dirty: changes=%d", changes)
	}

	m.ClearDirty()
	if m.IsDirty() || changes != 2 {
		t.Errorf("ClearDirty: dirty=%v changes=%d", m.IsDirty(), changes)
	}

	// Idempotent — ClearDirty when already clean is a no-op.
	m.ClearDirty()
	if changes != 2 {
		t.Errorf("ClearDirty fired callback while already clean: changes=%d", changes)
	}
}
