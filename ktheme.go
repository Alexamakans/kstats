package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type kTheme struct {
	font fyne.Resource
}

var _ fyne.Theme = (*kTheme)(nil)

func (k kTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameBackground {
		return color.Black
	}

	if name == theme.ColorNameForeground {
		return color.White
	}

	if name == theme.ColorNameSuccess {
		return color.RGBA{110, 110, 110, 255}
	}

	if name == theme.ColorNameError {
		return color.RGBA{130, 20, 20, 255}
	}

	return theme.DefaultTheme().Color(name, variant)
}

func (k kTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (k kTheme) Font(style fyne.TextStyle) fyne.Resource {
	if style.Monospace {
		return theme.DefaultTextMonospaceFont()
	}
	return k.font
}

func (k kTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
