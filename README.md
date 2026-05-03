# Wasteland 2 Editor

Cross-platform save editor for **Wasteland 2: Director's Cut**. Linux and Windows, single static binary, no system dependencies beyond the OpenGL/X11 libraries every desktop ships with.

A Go + Fyne port of [Tokeshy/Wasteland-2-caracter-editor](https://github.com/Tokeshy/Wasteland-2-caracter-editor). The original Delphi editor corrupts current saves (it writes the legacy `<attributes>...<pair><key>` dictionary format while the live game uses `<attributes2>...<KeyValuePairOfStringInt32>...`); this rewrite targets the modern format and round-trips a real 1.7 MB save **byte-for-byte identical** when nothing is edited.

## Features

- Edit attributes (Coordination / Luck / Awareness / Strength / Speed / Intelligence / Charisma)
- Edit weapon, general, and technical skills (level 0–10)
- Edit unspent attribute and skill points directly (independent of slider changes — moving a slider does **not** debit the point pool)
- Edit character level and current HP; view age and sex
- Drag-and-drop a save XML onto the window, or use **Open save…**
- Backup-on-save: writes `<savename>.xml.back` before overwriting
- Inventory, equipment, position, quest state — everything outside the edited fields is passed through untouched

## Build

```sh
go build ./cmd/wasteland-2-editor
./wasteland-2-editor
```

Tested with Go 1.26. Fyne pulls a handful of indirect dependencies on first build.

## Usage

1. Open or drag your `Quicksave N/Quicksave N.xml` (or any save XML) onto the window.
2. Pick a character from the dropdown.
3. Drag sliders on any of the four tabs.
4. Click **Save**. The original is renamed to `*.xml.back`.

Save files live under:

- Linux (Steam + Proton): `~/.steam/steam/steamapps/compatdata/240760/pfx/drive_c/users/steamuser/Documents/My Games/Wasteland2/Save Games/`
- Windows: `%USERPROFILE%\Documents\My Games\Wasteland2\Save Games\`

## Project layout

```
cmd/wasteland-2-editor/  entry point
internal/savegame/       stdlib XML parser + tests
internal/ui/             Fyne UI (window, character panel, skill tabs/rows)
delphi-original/         (gitignored) original Delphi sources, kept for reference
```

## Tests

```sh
go test ./...
```

The real-savegame tests are skipped automatically if `Quicksave 2/Quicksave 2.xml` isn't present.

## Known caveats

- Slider changes do not debit `availableSkillPoints` / `availableAttributePoints`. Tune those two pools yourself via the entries above the tabs if you want the in-game character sheet to show a particular number.
- Skill XP is set directly via the cumulative-cost table (`0,2,4,6,10,14,18,24,30,36,44`); intermediate values are not preserved.
- **Max HP** is not stored in the save — Wasteland 2 derives it at runtime from attributes / level / `luckyHitpoints`. Editing CON-style attributes is the indirect way to influence it.
- **Perks and quirks** are not present in the XML and likely live in the binary `.bin` sibling, which is not yet decoded. Trait points (`availableTraitPoints`) is also not exposed.

## License

Same spirit as the original: free, no warranty. Make a backup before editing.
