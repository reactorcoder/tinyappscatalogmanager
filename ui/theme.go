package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"image/color"
)

type customTheme struct {
}

var (
	tinyappsBackground = &color.NRGBA{R: 0x22, G: 0x2b, B: 0x39, A: 255}
	tinyappsHover      = &color.NRGBA{R: 0x3b, G: 0xa0, B: 0xe6, A: 255}
	tinyappsHover2     = &color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 100}
	tinyappsText       = &color.NRGBA{R: 0xe1, G: 0xdf, B: 0xdc, A: 255}
	tinyappsFocus      = &color.NRGBA{R: 0xcc, G: 0xcc, B: 0xcc, A: 100}
)

func (customTheme) Color(c fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	switch c {
	case theme.ColorNameBackground:
		return tinyappsBackground
	case theme.ColorNameButton, theme.ColorNameDisabled:
		return color.Black
	case theme.ColorNameForeground:
		return tinyappsText
	case theme.ColorNamePlaceHolder, theme.ColorNameScrollBar:
		return color.White
	case theme.ColorNameHover:
		return tinyappsHover2
	case theme.ColorNamePrimary:
		return color.White
	case theme.ColorNameFocus:
		return tinyappsFocus
	case theme.ColorNameShadow: // borders
		return color.Transparent
	case theme.ColorNameSeparator:
		return tinyappsBackground
	default:
		return color.Black
	}
}

func (customTheme) Size(s fyne.ThemeSizeName) float32 {
	switch s {
	case theme.SizeNameInnerPadding:
		return 5
	case theme.SizeNamePadding:
		return 0
	case theme.SizeNameInlineIcon:
		return 11
	case theme.SizeNameScrollBar:
		return 11
	case theme.SizeNameScrollBarSmall:
		return 11
	case theme.SizeNameText:
		return 11
	case theme.SizeNameHeadingText:
		return 11
	case theme.SizeNameSubHeadingText:
		return 11
	case theme.SizeNameInputBorder:
		return 1
	case theme.SizeNameSeparatorThickness:
		return 1
	default:
		return 0
	}
}

func (customTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (customTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}

func TinyappsUi() fyne.Theme {
	return &customTheme{}
}
