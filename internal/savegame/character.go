package savegame

import (
	"bytes"
	"fmt"
	"strings"
)

// Character is one party member, with editable attribute / skill maps and
// available-points counters.
type Character struct {
	// Name is the raw <name> field (often an internal ID for recruited NPCs,
	// e.g. "AZ6_Lexcanium_PC").
	Name string
	// DisplayName is the cleaned <displayName>, e.g. "Lexcanium".
	DisplayName string
	// Age and Gender are read-only metadata (parsed but never written back).
	Age    string
	Gender string
	// Level is the character level shown in-game.
	Level int
	// CurrentHP is the persisted current HP. Max HP is not stored — the
	// game derives it at runtime from attributes / level / luckyHitpoints.
	CurrentHP int
	// Attributes maps attribute name -> value (1..10, raw == display).
	Attributes map[string]int
	// Skills maps skill name -> raw XP value. Use SkillLevel/SkillXP to
	// convert to/from the 0..10 displayed level.
	Skills map[string]int
	// AvailableAttributePoints is the unspent attribute-point pool the game
	// shows on the character sheet. Editing skills/attributes via this
	// editor does NOT debit this counter — it's an independent field.
	AvailableAttributePoints int
	// AvailableSkillPoints is the unspent skill-point pool. Same independence
	// rule as above.
	AvailableSkillPoints int

	// Original key order for each dictionary, recorded at parse time.
	// Used by Save() so unedited regions round-trip byte-for-byte.
	attrOrder  []string
	skillOrder []string

	// Byte offsets into Save.raw for every region this character owns.
	attrStart, attrEnd               int
	skillStart, skillEnd             int
	attrPointsStart, attrPointsEnd   int
	skillPointsStart, skillPointsEnd int
	levelStart, levelEnd             int
	curHpStart, curHpEnd             int
}

const (
	pcDataOpen  = "<PcData>"
	pcDataClose = "</PcData>"

	attrOpen   = "<attributes2>"
	attrClose  = "</attributes2>"
	skillOpen  = "<skillXps2>"
	skillClose = "</skillXps2>"

	attrPointsOpen   = "<availableAttributePoints>"
	attrPointsClose  = "</availableAttributePoints>"
	skillPointsOpen  = "<availableSkillPoints>"
	skillPointsClose = "</availableSkillPoints>"

	levelOpen  = "<level>"
	levelClose = "</level>"
	curHpOpen  = "<curHp>"
	curHpClose = "</curHp>"
	ageOpen    = "<age>"
	ageClose   = "</age>"
	genderOpen  = "<gender>"
	genderClose = "</gender>"

	nameOpen         = "<name>"
	nameClose        = "</name>"
	displayNameOpen  = "<displayName>"
	displayNameClose = "</displayName>"

	// XML-encoded "<@>" prefix the game uses on display names.
	atMarker = "&lt;@&gt;"
)

// byteEdit is a contiguous region of the original save buffer to overwrite
// when serializing.
type byteEdit struct {
	start, end int
	body       []byte
}

func parseCharacter(buf []byte, start, end int) (*Character, error) {
	region := buf[start:end]

	name := extractInner(region, nameOpen, nameClose)
	if name == "" {
		// Empty <PcData /> placeholders exist — skip them.
		return nil, nil
	}
	display := extractInner(region, displayNameOpen, displayNameClose)
	if display == "" {
		display = name
	}

	c := &Character{
		Name:        name,
		DisplayName: cleanDisplayName(display),
		Age:         extractInner(region, ageOpen, ageClose),
		Gender:      extractInner(region, genderOpen, genderClose),
		Attributes:  map[string]int{},
		Skills:      map[string]int{},
	}

	if err := locateBlock(region, attrOpen, attrClose, start, &c.attrStart, &c.attrEnd); err != nil {
		return nil, err
	}
	if err := locateBlock(region, skillOpen, skillClose, start, &c.skillStart, &c.skillEnd); err != nil {
		return nil, err
	}
	if err := locateScalar(region, attrPointsOpen, attrPointsClose, start,
		&c.attrPointsStart, &c.attrPointsEnd, &c.AvailableAttributePoints); err != nil {
		return nil, err
	}
	if err := locateScalar(region, skillPointsOpen, skillPointsClose, start,
		&c.skillPointsStart, &c.skillPointsEnd, &c.AvailableSkillPoints); err != nil {
		return nil, err
	}
	if err := locateScalar(region, levelOpen, levelClose, start,
		&c.levelStart, &c.levelEnd, &c.Level); err != nil {
		return nil, err
	}
	if err := locateScalar(region, curHpOpen, curHpClose, start,
		&c.curHpStart, &c.curHpEnd, &c.CurrentHP); err != nil {
		return nil, err
	}

	order, err := parsePairs(buf[c.attrStart:c.attrEnd], c.Attributes)
	if err != nil {
		return nil, fmt.Errorf("attributes2: %w", err)
	}
	c.attrOrder = order

	order, err = parsePairs(buf[c.skillStart:c.skillEnd], c.Skills)
	if err != nil {
		return nil, fmt.Errorf("skillXps2: %w", err)
	}
	c.skillOrder = order

	return c, nil
}

// edits returns every byte region this character contributes to the output.
func (c *Character) edits() []byteEdit {
	return []byteEdit{
		{c.attrStart, c.attrEnd, renderPairs(attrOpen, attrClose, c.Attributes, c.attrOrder)},
		{c.skillStart, c.skillEnd, renderPairs(skillOpen, skillClose, c.Skills, c.skillOrder)},
		{c.attrPointsStart, c.attrPointsEnd, renderScalar(attrPointsOpen, attrPointsClose, c.AvailableAttributePoints)},
		{c.skillPointsStart, c.skillPointsEnd, renderScalar(skillPointsOpen, skillPointsClose, c.AvailableSkillPoints)},
		{c.levelStart, c.levelEnd, renderScalar(levelOpen, levelClose, c.Level)},
		{c.curHpStart, c.curHpEnd, renderScalar(curHpOpen, curHpClose, c.CurrentHP)},
	}
}

func locateBlock(region []byte, open, close string, base int, gotStart, gotEnd *int) error {
	o := bytes.Index(region, []byte(open))
	if o < 0 {
		return fmt.Errorf("missing %s", open)
	}
	c := bytes.Index(region[o:], []byte(close))
	if c < 0 {
		return fmt.Errorf("unterminated %s", open)
	}
	*gotStart = base + o
	*gotEnd = base + o + c + len(close)
	return nil
}

// locateScalar finds <open>N</close>, records the byte range that covers
// the whole element, and parses N into *gotValue.
func locateScalar(region []byte, open, close string, base int, gotStart, gotEnd, gotValue *int) error {
	o := bytes.Index(region, []byte(open))
	if o < 0 {
		return fmt.Errorf("missing %s", open)
	}
	contentStart := o + len(open)
	c := bytes.Index(region[contentStart:], []byte(close))
	if c < 0 {
		return fmt.Errorf("unterminated %s", open)
	}
	*gotStart = base + o
	*gotEnd = base + contentStart + c + len(close)

	n, err := parseInt(string(region[contentStart : contentStart+c]))
	if err != nil {
		return fmt.Errorf("%s: %w", open, err)
	}
	*gotValue = n
	return nil
}

func renderScalar(open, close string, value int) []byte {
	return []byte(fmt.Sprintf("%s%d%s", open, value, close))
}

func extractInner(region []byte, open, close string) string {
	o := bytes.Index(region, []byte(open))
	if o < 0 {
		return ""
	}
	o += len(open)
	c := bytes.Index(region[o:], []byte(close))
	if c < 0 {
		return ""
	}
	return string(region[o : o+c])
}

// cleanDisplayName turns "<@>Hex{F}" into "Hex".
func cleanDisplayName(raw string) string {
	s := strings.TrimPrefix(raw, "<@>")
	s = strings.TrimPrefix(s, atMarker)
	if i := strings.LastIndex(s, "{"); i >= 0 {
		if j := strings.LastIndex(s, "}"); j > i {
			s = s[:i] + s[j+1:]
		}
	}
	return strings.TrimSpace(s)
}
