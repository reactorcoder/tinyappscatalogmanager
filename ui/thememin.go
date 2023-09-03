package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"image/color"
)

type customThememin struct {
}

func (customThememin) Color(c fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
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

func (customThememin) Size(s fyne.ThemeSizeName) float32 {
	switch s {
	case theme.SizeNameInnerPadding:
		return 3
	case theme.SizeNamePadding:
		return 0
	case theme.SizeNameInlineIcon:
		return 10
	case theme.SizeNameScrollBar:
		return 10
	case theme.SizeNameScrollBarSmall:
		return 10
	case theme.SizeNameText:
		return 10
	case theme.SizeNameHeadingText:
		return 10
	case theme.SizeNameSubHeadingText:
		return 10
	case theme.SizeNameInputBorder:
		return 1
	case theme.SizeNameSeparatorThickness:
		return 1
	default:
		return 0
	}
}

func (customThememin) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (customThememin) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}

func TinyappsUiRecessed() fyne.Theme {
	return &customThememin{}
}
