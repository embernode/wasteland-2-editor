package ui

import (
	_ "embed"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

//go:embed fonts/SpaceGrotesk-Variable.ttf
var spaceGroteskTTF []byte

var spaceGrotesk = &fyne.StaticResource{
	StaticName:    "SpaceGrotesk-Variable.ttf",
	StaticContent: spaceGroteskTTF,
}

// Palette is a complete color scheme. Add a new theme by defining another
// Palette value and registering it in themesByName.
type Palette struct {
	Name string

	Surface          color.Color
	SurfaceLow       color.Color
	SurfaceContainer color.Color
	SurfaceHigh      color.Color
	OnSurface        color.Color
	OnSurfaceMuted   color.Color
	Outline          color.Color

	Primary      color.Color
	PrimaryHover color.Color
	Selection    color.Color
	Focus        color.Color
	Error        color.Color
}

// WastelandPalette — Stitch "Post-Apocalyptic Terminal": amber primary,
// grimy desert browns. Focus / selection are dark amber so foreground text
// stays legible (Fyne uses both as fill backgrounds, not outlines).
var WastelandPalette = Palette{
	Name:             "Wasteland",
	Surface:          mustColor("#17130f"),
	SurfaceLow:       mustColor("#1f1b17"),
	SurfaceContainer: mustColor("#231f1b"),
	SurfaceHigh:      mustColor("#2e2925"),
	OnSurface:        mustColor("#eae1da"),
	OnSurfaceMuted:   mustColor("#dac2ae"),
	Outline:          mustColor("#544434"),
	Primary:          mustColor("#ffb86e"),
	PrimaryHover:     mustColor("#e68a00"),
	Selection:        mustColor("#492900"),
	Focus:            mustColor("#492900"),
	Error:            mustColor("#ffb4ab"),
}

// PrecisionPalette — Stitch "Precision Utility": deep-charcoal background,
// electric cyan accent, professional dashboard feel.
var PrecisionPalette = Palette{
	Name:             "Precision",
	Surface:          mustColor("#121414"),
	SurfaceLow:       mustColor("#1a1c1c"),
	SurfaceContainer: mustColor("#1e2020"),
	SurfaceHigh:      mustColor("#282a2b"),
	OnSurface:        mustColor("#e2e2e2"),
	OnSurfaceMuted:   mustColor("#b9cacb"),
	Outline:          mustColor("#3b494b"),
	Primary:          mustColor("#00dbe9"),
	PrimaryHover:     mustColor("#00f0ff"),
	Selection:        mustColor("#004f54"),
	Focus:            mustColor("#004f54"),
	Error:            mustColor("#ffb4ab"),
}

var themesByName = map[string]Palette{
	WastelandPalette.Name: WastelandPalette,
	PrecisionPalette.Name: PrecisionPalette,
}

// ThemeNames returns the registered palette names in display order.
func ThemeNames() []string {
	return []string{WastelandPalette.Name, PrecisionPalette.Name}
}

// ThemeByName looks up a palette wrapped in a fyne.Theme. Falls back to
// Wasteland if name is unknown.
func ThemeByName(name string) fyne.Theme {
	p, ok := themesByName[name]
	if !ok {
		p = WastelandPalette
	}
	return paletteTheme{p}
}

// paletteTheme is the fyne.Theme implementation backed by a Palette.
type paletteTheme struct{ p Palette }

var _ fyne.Theme = paletteTheme{}

func (t paletteTheme) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return t.p.Surface
	case theme.ColorNameForeground:
		return t.p.OnSurface
	case theme.ColorNameForegroundOnPrimary:
		return t.p.Surface
	case theme.ColorNamePrimary:
		return t.p.Primary
	case theme.ColorNameButton:
		return t.p.SurfaceContainer
	case theme.ColorNameDisabledButton:
		return t.p.SurfaceLow
	case theme.ColorNameInputBackground:
		return t.p.SurfaceLow
	case theme.ColorNameInputBorder:
		return t.p.Outline
	case theme.ColorNameSeparator:
		return t.p.Outline
	case theme.ColorNameDisabled:
		return t.p.OnSurfaceMuted
	case theme.ColorNamePlaceHolder:
		return t.p.OnSurfaceMuted
	case theme.ColorNameHover:
		return t.p.SurfaceHigh
	case theme.ColorNameFocus:
		return t.p.Focus
	case theme.ColorNameSelection:
		return t.p.Selection
	case theme.ColorNamePressed:
		return t.p.PrimaryHover
	case theme.ColorNameMenuBackground:
		return t.p.SurfaceContainer
	case theme.ColorNameOverlayBackground:
		return t.p.SurfaceHigh
	case theme.ColorNameError:
		return t.p.Error
	case theme.ColorNameScrollBar, theme.ColorNameShadow:
		return color.NRGBA{R: 0, G: 0, B: 0, A: 96}
	}
	return theme.DefaultTheme().Color(name, theme.VariantDark)
}

func (paletteTheme) Font(style fyne.TextStyle) fyne.Resource {
	if style.Monospace {
		return spaceGrotesk
	}
	return theme.DefaultTheme().Font(style)
}

func (paletteTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (paletteTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

// mustColor parses a #RRGGBB hex string. Panics on a typo at startup.
func mustColor(hex string) color.NRGBA {
	if len(hex) != 7 || hex[0] != '#' {
		panic("invalid color literal: " + hex)
	}
	var r, g, b uint8
	for i, dst := range []*uint8{&r, &g, &b} {
		hi := hexNibble(hex[1+i*2])
		lo := hexNibble(hex[2+i*2])
		*dst = hi<<4 | lo
	}
	return color.NRGBA{R: r, G: g, B: b, A: 0xff}
}

func hexNibble(c byte) uint8 {
	switch {
	case c >= '0' && c <= '9':
		return c - '0'
	case c >= 'a' && c <= 'f':
		return c - 'a' + 10
	case c >= 'A' && c <= 'F':
		return c - 'A' + 10
	}
	panic("invalid hex digit: " + string(c))
}
