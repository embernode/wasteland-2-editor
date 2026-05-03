package savegame

import (
	"bytes"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

// realSavePath is the absolute path to the user-supplied test save. It lives
// outside the package directory so tests are skipped automatically when the
// file isn't present (e.g. CI without the asset).
func realSavePath(t *testing.T) string {
	t.Helper()
	// Walk upward from the package dir until we find the project root.
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
	t.Skip("Quicksave 2/Quicksave 2.xml not found — skipping real-savegame test")
	return ""
}

func TestRealSave_Load(t *testing.T) {
	s, err := Load(realSavePath(t))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(s.Characters) != 7 {
		t.Errorf("got %d characters, want 7", len(s.Characters))
	}

	wantNames := []string{
		"Hex", "Pills", "Slick", "Rose", "Lexcanium", "Pizepi Joren", "Big Bert",
	}
	for i, c := range s.Characters {
		if i >= len(wantNames) {
			break
		}
		if c.DisplayName != wantNames[i] {
			t.Errorf("character[%d].DisplayName = %q, want %q",
				i, c.DisplayName, wantNames[i])
		}
		// Every character should have all 7 attributes and at least all skills
		// the game ships with.
		if len(c.Attributes) != len(Attributes) {
			t.Errorf("character %q: %d attributes, want %d",
				c.DisplayName, len(c.Attributes), len(Attributes))
		}
		expectedSkillCount := len(WeaponSkills) + len(GeneralSkills) + len(TechnicalSkills)
		if len(c.Skills) != expectedSkillCount {
			t.Errorf("character %q: %d skills, want %d",
				c.DisplayName, len(c.Skills), expectedSkillCount)
		}
		// Every skill name we know about should be present.
		for _, name := range allSkillNames() {
			if _, ok := c.Skills[name]; !ok {
				t.Errorf("character %q missing skill %q", c.DisplayName, name)
			}
		}
		for _, name := range Attributes {
			if _, ok := c.Attributes[name]; !ok {
				t.Errorf("character %q missing attribute %q", c.DisplayName, name)
			}
		}
	}
}

func TestRealSave_RoundTripByteEqual(t *testing.T) {
	path := realSavePath(t)
	original, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	s, err := Parse(original, path)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	out, err := s.Bytes()
	if err != nil {
		t.Fatalf("Bytes: %v", err)
	}
	if !bytes.Equal(original, out) {
		// Find first byte that differs for a useful diagnostic.
		diff := firstDiff(original, out)
		ctxStart := max0(diff - 40)
		t.Fatalf(
			"round-trip differs at offset %d (file is %d bytes, output is %d)\n"+
				"  original: %q\n"+
				"  output:   %q",
			diff, len(original), len(out),
			string(original[ctxStart:min(diff+40, len(original))]),
			string(out[ctxStart:min(diff+40, len(out))]),
		)
	}
}

func TestRealSave_EditPersists(t *testing.T) {
	path := realSavePath(t)
	original, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	s, err := Parse(original, path)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	// Bump strength to 10, rifle to level 10, give 99 spare points to each.
	for _, c := range s.Characters {
		c.Attributes["strength"] = 10
		c.Skills["rifle"] = SkillXP(10)
		c.AvailableAttributePoints = 7
		c.AvailableSkillPoints = 99
	}

	out, err := s.Bytes()
	if err != nil {
		t.Fatal(err)
	}

	// Re-parse the output and confirm the edits stuck.
	s2, err := Parse(out, path)
	if err != nil {
		t.Fatalf("re-Parse edited save: %v", err)
	}
	if len(s2.Characters) != len(s.Characters) {
		t.Fatalf("character count changed: %d -> %d",
			len(s.Characters), len(s2.Characters))
	}
	for i, c := range s2.Characters {
		if c.Attributes["strength"] != 10 {
			t.Errorf("character[%d] strength = %d, want 10", i, c.Attributes["strength"])
		}
		if SkillLevel(c.Skills["rifle"]) != 10 {
			t.Errorf("character[%d] rifle = %d (level %d), want level 10",
				i, c.Skills["rifle"], SkillLevel(c.Skills["rifle"]))
		}
		if c.AvailableAttributePoints != 7 {
			t.Errorf("character[%d] AvailableAttributePoints = %d, want 7",
				i, c.AvailableAttributePoints)
		}
		if c.AvailableSkillPoints != 99 {
			t.Errorf("character[%d] AvailableSkillPoints = %d, want 99",
				i, c.AvailableSkillPoints)
		}
	}
}

func TestRealSave_AvailablePointsParsed(t *testing.T) {
	s, err := Load(realSavePath(t))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	// Hex (character 0) — observed in raw XML:
	// <availableAttributePoints>0</availableAttributePoints>
	// <availableSkillPoints>6</availableSkillPoints>
	hex := s.Characters[0]
	if hex.DisplayName != "Hex" {
		t.Fatalf("expected character[0] = Hex, got %q", hex.DisplayName)
	}
	if hex.AvailableAttributePoints != 0 {
		t.Errorf("Hex.AvailableAttributePoints = %d, want 0", hex.AvailableAttributePoints)
	}
	if hex.AvailableSkillPoints != 6 {
		t.Errorf("Hex.AvailableSkillPoints = %d, want 6", hex.AvailableSkillPoints)
	}
}

func allSkillNames() []string {
	all := append([]string{}, WeaponSkills...)
	all = append(all, GeneralSkills...)
	all = append(all, TechnicalSkills...)
	sort.Strings(all)
	return all
}

func firstDiff(a, b []byte) int {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		if a[i] != b[i] {
			return i
		}
	}
	return n
}

func max0(n int) int {
	if n < 0 {
		return 0
	}
	return n
}
