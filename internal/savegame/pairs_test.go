package savegame

import (
	"reflect"
	"testing"
)

func TestParsePairs_simple(t *testing.T) {
	in := []byte(`<attributes2>` +
		`<KeyValuePairOfStringInt32><Key>coordination</Key><Value>5</Value></KeyValuePairOfStringInt32>` +
		`<KeyValuePairOfStringInt32><Key>luck</Key><Value>2</Value></KeyValuePairOfStringInt32>` +
		`</attributes2>`)
	got := map[string]int{}
	order, err := parsePairs(in, got)
	if err != nil {
		t.Fatalf("parsePairs: %v", err)
	}
	wantMap := map[string]int{"coordination": 5, "luck": 2}
	wantOrder := []string{"coordination", "luck"}
	if !reflect.DeepEqual(got, wantMap) {
		t.Errorf("map: got %v, want %v", got, wantMap)
	}
	if !reflect.DeepEqual(order, wantOrder) {
		t.Errorf("order: got %v, want %v", order, wantOrder)
	}
}

func TestRenderPairs_preservesRecordedOrder(t *testing.T) {
	// Game-defined attribute order — not alphabetical.
	order := []string{"coordination", "luck", "awareness", "strength",
		"speed", "intelligence", "charisma"}
	m := map[string]int{
		"coordination": 5, "luck": 2, "awareness": 6, "strength": 2,
		"speed": 4, "intelligence": 10, "charisma": 2,
	}
	got := string(renderPairs("<attributes2>", "</attributes2>", m, order))
	want := `<attributes2>` +
		`<KeyValuePairOfStringInt32><Key>coordination</Key><Value>5</Value></KeyValuePairOfStringInt32>` +
		`<KeyValuePairOfStringInt32><Key>luck</Key><Value>2</Value></KeyValuePairOfStringInt32>` +
		`<KeyValuePairOfStringInt32><Key>awareness</Key><Value>6</Value></KeyValuePairOfStringInt32>` +
		`<KeyValuePairOfStringInt32><Key>strength</Key><Value>2</Value></KeyValuePairOfStringInt32>` +
		`<KeyValuePairOfStringInt32><Key>speed</Key><Value>4</Value></KeyValuePairOfStringInt32>` +
		`<KeyValuePairOfStringInt32><Key>intelligence</Key><Value>10</Value></KeyValuePairOfStringInt32>` +
		`<KeyValuePairOfStringInt32><Key>charisma</Key><Value>2</Value></KeyValuePairOfStringInt32>` +
		`</attributes2>`
	if got != want {
		t.Errorf("\n got: %s\nwant: %s", got, want)
	}
}

func TestRenderPairs_unknownKeysAppendedAlphabetically(t *testing.T) {
	// Recorded order has only one entry — extras should follow, sorted.
	order := []string{"luck"}
	m := map[string]int{"luck": 2, "zebra": 9, "apple": 1}
	got := string(renderPairs("<x>", "</x>", m, order))
	want := `<x>` +
		`<KeyValuePairOfStringInt32><Key>luck</Key><Value>2</Value></KeyValuePairOfStringInt32>` +
		`<KeyValuePairOfStringInt32><Key>apple</Key><Value>1</Value></KeyValuePairOfStringInt32>` +
		`<KeyValuePairOfStringInt32><Key>zebra</Key><Value>9</Value></KeyValuePairOfStringInt32>` +
		`</x>`
	if got != want {
		t.Errorf("\n got: %s\nwant: %s", got, want)
	}
}

func TestRenderPairs_skipsDeletedKeys(t *testing.T) {
	order := []string{"a", "b", "c"}
	m := map[string]int{"a": 1, "c": 3} // "b" was removed from the map
	got := string(renderPairs("<x>", "</x>", m, order))
	want := `<x>` +
		`<KeyValuePairOfStringInt32><Key>a</Key><Value>1</Value></KeyValuePairOfStringInt32>` +
		`<KeyValuePairOfStringInt32><Key>c</Key><Value>3</Value></KeyValuePairOfStringInt32>` +
		`</x>`
	if got != want {
		t.Errorf("\n got: %s\nwant: %s", got, want)
	}
}

func TestRoundTrip_orderedPreserved(t *testing.T) {
	original := []byte(`<skillXps2>` +
		`<KeyValuePairOfStringInt32><Key>handgun</Key><Value>0</Value></KeyValuePairOfStringInt32>` +
		`<KeyValuePairOfStringInt32><Key>perception</Key><Value>36</Value></KeyValuePairOfStringInt32>` +
		`<KeyValuePairOfStringInt32><Key>rifle</Key><Value>44</Value></KeyValuePairOfStringInt32>` +
		`</skillXps2>`)
	parsed := map[string]int{}
	order, err := parsePairs(original, parsed)
	if err != nil {
		t.Fatalf("parsePairs: %v", err)
	}
	rerendered := renderPairs("<skillXps2>", "</skillXps2>", parsed, order)
	if string(rerendered) != string(original) {
		t.Errorf("round-trip mismatch:\n got: %s\nwant: %s", rerendered, original)
	}
}

func TestSkillLevelXP(t *testing.T) {
	cases := []struct {
		xp, level int
	}{
		{0, 0}, {2, 1}, {4, 2}, {6, 3}, {10, 4}, {14, 5},
		{18, 6}, {24, 7}, {30, 8}, {36, 9}, {44, 10},
		// in-between values round down
		{3, 1}, {5, 2}, {15, 5}, {43, 9},
		// out-of-range rounds to top of table
		{100, 10},
	}
	for _, c := range cases {
		if got := SkillLevel(c.xp); got != c.level {
			t.Errorf("SkillLevel(%d)=%d, want %d", c.xp, got, c.level)
		}
	}
	for level := 0; level <= 10; level++ {
		xp := SkillXP(level)
		if SkillLevel(xp) != level {
			t.Errorf("SkillXP(%d)=%d, but SkillLevel(%d)=%d (want %d)",
				level, xp, xp, SkillLevel(xp), level)
		}
	}
}

func TestCleanDisplayName(t *testing.T) {
	cases := map[string]string{
		"<@>Hex{F}":          "Hex",
		"<@>Pizepi Joren{F}": "Pizepi Joren",
		"<@>Big Bert{M}":     "Big Bert",
		"Plain":              "Plain",
	}
	for in, want := range cases {
		if got := cleanDisplayName(in); got != want {
			t.Errorf("cleanDisplayName(%q)=%q, want %q", in, got, want)
		}
	}
}
