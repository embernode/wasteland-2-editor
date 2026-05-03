package savegame

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
)

// The .NET XmlSerializer emits dictionaries as a sequence of
//   <KeyValuePairOfStringInt32><Key>name</Key><Value>n</Value></KeyValuePairOfStringInt32>
// inside the dictionary's container element.

const (
	pairElement  = "KeyValuePairOfStringInt32"
	keyElement   = "Key"
	valueElement = "Value"
)

// parsePairs decodes a <container>...key/value pairs...</container> block
// into the supplied map and returns the keys in the order they appeared in
// the source. Preserving order matters: <skillXps2> is alphabetical but
// <attributes2> follows a game-defined order, and the editor must round-trip
// either layout byte-for-byte.
func parsePairs(block []byte, into map[string]int) ([]string, error) {
	dec := xml.NewDecoder(bytes.NewReader(block))
	var (
		curKey   string
		curValue int
		inKey    bool
		inValue  bool
		haveKV   bool
		order    []string
	)
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			return order, nil
		}
		if err != nil {
			return nil, err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case pairElement:
				curKey = ""
				curValue = 0
				haveKV = false
			case keyElement:
				inKey = true
			case valueElement:
				inValue = true
			}
		case xml.CharData:
			if inKey {
				curKey += string(t)
			}
			if inValue {
				v, err := parseInt(string(t))
				if err != nil {
					return nil, fmt.Errorf("value for %q: %w", curKey, err)
				}
				curValue = v
				haveKV = true
			}
		case xml.EndElement:
			switch t.Name.Local {
			case keyElement:
				inKey = false
			case valueElement:
				inValue = false
			case pairElement:
				key := strings.TrimSpace(curKey)
				if key != "" && haveKV {
					if _, exists := into[key]; !exists {
						order = append(order, key)
					}
					into[key] = curValue
				}
			}
		}
	}
}

// renderPairs serializes a map back into the .NET XmlSerializer dictionary
// format. Keys appear first in the order recorded during parse; any keys
// added since (not in `order`) are appended alphabetically.
func renderPairs(open, close string, m map[string]int, order []string) []byte {
	seen := make(map[string]bool, len(order))
	var b bytes.Buffer
	b.Grow(len(open) + len(close) + len(m)*80)
	b.WriteString(open)
	for _, k := range order {
		if _, ok := m[k]; !ok {
			continue // key was deleted
		}
		writePair(&b, k, m[k])
		seen[k] = true
	}
	// Anything new is sorted to keep the output deterministic.
	var extras []string
	for k := range m {
		if !seen[k] {
			extras = append(extras, k)
		}
	}
	sort.Strings(extras)
	for _, k := range extras {
		writePair(&b, k, m[k])
	}
	b.WriteString(close)
	return b.Bytes()
}

func writePair(b *bytes.Buffer, key string, value int) {
	fmt.Fprintf(b,
		"<%s><%s>%s</%s><%s>%d</%s></%s>",
		pairElement,
		keyElement, escape(key), keyElement,
		valueElement, value, valueElement,
		pairElement,
	)
}

func escape(s string) string {
	var b strings.Builder
	_ = xml.EscapeText(&b, []byte(s))
	return b.String()
}

// parseInt is a thin wrapper around strconv.Atoi that tolerates surrounding
// whitespace (CharData between pretty-printed XML tags) and treats an empty
// string as 0.
func parseInt(s string) (int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, nil
	}
	return strconv.Atoi(s)
}
