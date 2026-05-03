package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// wastelandTheme overrides Fyne's default colors with the "Post-Apocalyptic
// Terminal" palette: amber primary, dim brown surfaces, phosphor-green for
// active state. Fonts and icon set fall back to the default theme.
type wastelandTheme struct{}

var _ fyne.Theme = (*wastelandTheme)(nil)

// Theme returns the editor's custom Fyne theme.
func Theme() fyne.Theme { return wastelandTheme{} }

// Palette — values from the Stitch design's "Grimy Desert" spectrum.
var (
	colSurface          = mustColor("#17130f")
	colSurfaceContainer = mustColor("#231f1b")
	colSurfaceLow       = mustColor("#1f1b17")
	colSurfaceHigh      = mustColor("#2e2925")
	colOnSurface        = mustColor("#eae1da")
	colOnSurfaceMuted   = mustColor("#dac2ae")
	colOutline          = mustColor("#544434")
	colPrimary          = mustColor("#ffb86e") // amber
	colPrimaryHover     = mustColor("#e68a00") // hotter amber for hover/active
	colPhosphor         = mustColor("#5dff3b") // active / focus / selection
	colError            = mustColor("#ffb4ab")
)

// Color overrides. Anything we don't handle falls through to the default
// dark theme so we get sensible values for menus, scrollbars, etc.
func (wastelandTheme) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return colSurface
	case theme.ColorNameForeground:
		return colOnSurface
	case theme.ColorNameForegroundOnPrimary:
		return colSurface
	case theme.ColorNamePrimary:
		return colPrimary
	case theme.ColorNameButton:
		return colSurfaceContainer
	case theme.ColorNameDisabledButton:
		return colSurfaceLow
	case theme.ColorNameInputBackground:
		return colSurfaceLow
	case theme.ColorNameInputBorder:
		return colOutline
	case theme.ColorNameSeparator:
		return colOutline
	case theme.ColorNameDisabled:
		return colOnSurfaceMuted
	case theme.ColorNamePlaceHolder:
		return colOnSurfaceMuted
	case theme.ColorNameHover:
		return colSurfaceHigh
	case theme.ColorNameFocus:
		return colPhosphor
	case theme.ColorNameSelection:
		return colPhosphor
	case theme.ColorNamePressed:
		return colPrimaryHover
	case theme.ColorNameMenuBackground:
		return colSurfaceContainer
	case theme.ColorNameOverlayBackground:
		return colSurfaceHigh
	case theme.ColorNameError:
		return colError
	case theme.ColorNameScrollBar, theme.ColorNameShadow:
		return color.NRGBA{R: 0, G: 0, B: 0, A: 96}
	}
	return theme.DefaultTheme().Color(name, theme.VariantDark)
}

func (wastelandTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (wastelandTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (wastelandTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

// mustColor parses a #RRGGBB hex string. Panics if invalid — only used on
// package-level constants, so a typo fails fast at startup.
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
