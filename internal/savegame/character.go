package savegame

import (
	"bytes"
	"fmt"
	"strings"
)

// Character is one party member, with editable attribute and skill maps.
type Character struct {
	// Name is the raw <name> field (often an internal ID for recruited NPCs,
	// e.g. "AZ6_Lexcanium_PC").
	Name string
	// DisplayName is the cleaned <displayName>, e.g. "Lexcanium".
	DisplayName string
	// Attributes maps attribute name -> value (1..10, raw == display).
	Attributes map[string]int
	// Skills maps skill name -> raw XP value. Use SkillLevel/SkillXP to
	// convert to/from the 0..10 displayed level.
	Skills map[string]int

	// Original key order for each dictionary, recorded at parse time.
	// Used by Save() so unedited regions round-trip byte-for-byte.
	attrOrder  []string
	skillOrder []string

	// Byte offsets into Save.raw of the dictionary blocks.
	attrStart, attrEnd   int
	skillStart, skillEnd int
}

const (
	pcDataOpen  = "<PcData>"
	pcDataClose = "</PcData>"

	attrOpen   = "<attributes2>"
	attrClose  = "</attributes2>"
	skillOpen  = "<skillXps2>"
	skillClose = "</skillXps2>"

	nameOpen         = "<name>"
	nameClose        = "</name>"
	displayNameOpen  = "<displayName>"
	displayNameClose = "</displayName>"

	// XML-encoded "<@>" prefix the game uses on display names.
	atMarker = "&lt;@&gt;"
)

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
		Attributes:  map[string]int{},
		Skills:      map[string]int{},
	}

	if err := locateBlock(region, attrOpen, attrClose, start, &c.attrStart, &c.attrEnd); err != nil {
		return nil, err
	}
	if err := locateBlock(region, skillOpen, skillClose, start, &c.skillStart, &c.skillEnd); err != nil {
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
