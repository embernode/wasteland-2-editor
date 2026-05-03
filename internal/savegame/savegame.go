// Package savegame loads, edits, and writes Wasteland 2 (Director's Cut)
// player save files.
//
// The save XML is a serialized .NET object graph; party members live under
// /UserData/data/pcs/PcData. Each PcData carries an <attributes2> and a
// <skillXps2> block — both are dictionaries serialized as a sequence of
// <KeyValuePairOfStringInt32><Key>name</Key><Value>n</Value>...</…> entries.
//
// To stay safe on a save format we don't fully model, the parser leaves the
// original byte buffer untouched everywhere except inside those two dictionary
// blocks for each character. On Save() we regenerate just those regions and
// splice them back in.
package savegame

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"sort"
)

// Save is a parsed save game. The original file bytes are preserved so we can
// write back with minimal changes.
type Save struct {
	Path       string
	raw        []byte
	Characters []*Character
}

// Load reads and parses the save game at path.
func Load(path string) (*Save, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read save: %w", err)
	}
	return Parse(raw, path)
}

// Parse parses an in-memory save buffer.
func Parse(raw []byte, path string) (*Save, error) {
	s := &Save{Path: path, raw: raw}

	cursor := 0
	for {
		open := bytes.Index(s.raw[cursor:], []byte(pcDataOpen))
		if open < 0 {
			break
		}
		open += cursor
		end := bytes.Index(s.raw[open:], []byte(pcDataClose))
		if end < 0 {
			return nil, fmt.Errorf("savegame: unterminated %s at offset %d", pcDataOpen, open)
		}
		end += open + len(pcDataClose)

		c, err := parseCharacter(s.raw, open, end)
		if err != nil {
			return nil, fmt.Errorf("character at offset %d: %w", open, err)
		}
		if c != nil {
			s.Characters = append(s.Characters, c)
		}
		cursor = end
	}
	if len(s.Characters) == 0 {
		return nil, errors.New("savegame: no characters found")
	}
	return s, nil
}

// Save writes the (possibly edited) save to dst. If dst exists it is renamed
// to dst+".back" first; any pre-existing .back is removed.
func (s *Save) Save(dst string) error {
	out, err := s.Bytes()
	if err != nil {
		return err
	}
	if _, err := os.Stat(dst); err == nil {
		backup := dst + ".back"
		_ = os.Remove(backup)
		if err := os.Rename(dst, backup); err != nil {
			return fmt.Errorf("backup: %w", err)
		}
	}
	if err := os.WriteFile(dst, out, 0o644); err != nil {
		return fmt.Errorf("write save: %w", err)
	}
	return nil
}

// Bytes returns the serialized save with each character's edited regions
// spliced into the original byte buffer. Everything outside those regions
// passes through unchanged.
func (s *Save) Bytes() ([]byte, error) {
	var edits []byteEdit
	for _, c := range s.Characters {
		edits = append(edits, c.edits()...)
	}
	sort.Slice(edits, func(i, j int) bool { return edits[i].start < edits[j].start })

	var out bytes.Buffer
	out.Grow(len(s.raw))
	cursor := 0
	for _, e := range edits {
		if e.start < cursor {
			return nil, fmt.Errorf("savegame: overlapping edits at %d", e.start)
		}
		out.Write(s.raw[cursor:e.start])
		out.Write(e.body)
		cursor = e.end
	}
	out.Write(s.raw[cursor:])
	return out.Bytes(), nil
}
