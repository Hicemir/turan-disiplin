package theme

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"

	"turandisiplin/assets"
)

type TurkishTheme struct {
	base         fyne.Theme
	Mode         string      // "light", "cream", "dark", ..., "custom"
	BaseSize     float32     // 15, 17, 19
	CustomBg     color.Color // user selected bg
	CustomAccent color.Color // user selected accent
}

type ThemeInfo struct {
	ID          string
	Name        string
	BgColor     color.Color
	AccentColor color.Color
}

func GetAllThemes() []ThemeInfo {
	return []ThemeInfo{
		// 1. Koyu Turkuaz (Gece)
		{ID: "dark_cyan", Name: "Koyu Turkuaz (Gece Modu)", BgColor: color.NRGBA{R: 4, G: 24, B: 36, A: 255}, AccentColor: color.NRGBA{R: 6, G: 182, B: 212, A: 255}},

		// 2. Açık Turkuaz (Gündüz) - Arka Plan: #e2e8f0, Vurgu: #06b6d4
		{ID: "light", Name: "Açık Turkuaz (Gündüz Modu)", BgColor: color.NRGBA{R: 226, G: 232, B: 240, A: 255}, AccentColor: color.NRGBA{R: 6, G: 182, B: 212, A: 255}},
	}
}

func adjustBrightness(c color.Color, delta int) color.Color {
	r, g, b, _ := c.RGBA()
	clamp := func(v int) uint8 {
		if v < 0 {
			return 0
		}
		if v > 255 {
			return 255
		}
		return uint8(v)
	}
	return color.NRGBA{
		R: clamp(int(uint8(r>>8)) + delta),
		G: clamp(int(uint8(g>>8)) + delta),
		B: clamp(int(uint8(b>>8)) + delta),
		A: 255,
	}
}

func buildDarkPalette(name fyne.ThemeColorName, bg, overlay, button, accent, separator color.Color, base fyne.Theme) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return bg
	case theme.ColorNameOverlayBackground, theme.ColorNameInputBackground:
		return overlay
	case theme.ColorNameForeground:
		return color.NRGBA{R: 248, G: 250, B: 252, A: 255}
	case theme.ColorNameForegroundOnPrimary:
		return color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	case theme.ColorNamePrimary:
		return accent
	case theme.ColorNameSuccess:
		return color.NRGBA{R: 34, G: 197, B: 94, A: 255}
	case theme.ColorNameError:
		return color.NRGBA{R: 239, G: 68, B: 68, A: 255}
	case theme.ColorNameButton:
		return button
	case theme.ColorNameSeparator:
		return separator
	case theme.ColorNameHover:
		return color.NRGBA{R: 255, G: 255, B: 255, A: 45}
	default:
		return base.Color(name, theme.VariantDark)
	}
}

func NewTurkishTheme(mode string, fontSize float32, customColors ...color.Color) fyne.Theme {
	if mode == "" {
		mode = "light"
	}
	if fontSize < 13 {
		fontSize = 15.0
	}
	var cBg, cAccent color.Color
	if len(customColors) >= 2 {
		cBg = customColors[0]
		cAccent = customColors[1]
	}
	return &TurkishTheme{
		base:         theme.DefaultTheme(),
		Mode:         mode,
		BaseSize:     fontSize,
		CustomBg:     cBg,
		CustomAccent: cAccent,
	}
}

var (
	resRegular    = fyne.NewStaticResource("LiberationSans-Regular.ttf", assets.FontRegular)
	resBold       = fyne.NewStaticResource("LiberationSans-Bold.ttf", assets.FontBold)
	resItalic     = fyne.NewStaticResource("LiberationSans-Italic.ttf", assets.FontItalic)
	resBoldItalic = fyne.NewStaticResource("LiberationSans-BoldItalic.ttf", assets.FontBoldItalic)
)

func (t *TurkishTheme) Font(style fyne.TextStyle) fyne.Resource {
	if style.Bold && style.Italic {
		return resBoldItalic
	}
	if style.Bold {
		return resBold
	}
	if style.Italic {
		return resItalic
	}
	return resRegular
}

func (t *TurkishTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch t.Mode {
	case "dark_cyan", "dark", "dark_emerald", "dark_purple", "dark_amber", "dark_ruby": // Koyu Turkuaz (Gece Modu)
		return buildDarkPalette(name,
			color.NRGBA{R: 4, G: 24, B: 36, A: 255},
			color.NRGBA{R: 8, G: 37, B: 54, A: 255},
			color.NRGBA{R: 12, G: 54, B: 78, A: 255},
			color.NRGBA{R: 6, G: 182, B: 212, A: 255},
			color.NRGBA{R: 17, G: 75, B: 109, A: 255},
			t.base)

	case "custom":
		bg := t.CustomBg
		if bg == nil {
			bg = color.NRGBA{R: 226, G: 232, B: 240, A: 255}
		}
		accent := t.CustomAccent
		if accent == nil {
			accent = color.NRGBA{R: 6, G: 182, B: 212, A: 255}
		}

		r, g, b, _ := bg.RGBA()
		r8, g8, b8 := uint8(r>>8), uint8(g>>8), uint8(b>>8)
		lum := 0.299*float64(r8) + 0.587*float64(g8) + 0.114*float64(b8)
		isDark := lum < 140

		if isDark {
			overlay := adjustBrightness(bg, 18)
			button := adjustBrightness(bg, 32)
			separator := adjustBrightness(bg, 48)
			return buildDarkPalette(name, bg, overlay, button, accent, separator, t.base)
		} else {
			overlay := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
			button := adjustBrightness(bg, -18)
			separator := adjustBrightness(bg, -35)
			switch name {
			case theme.ColorNameBackground:
				return bg
			case theme.ColorNameOverlayBackground, theme.ColorNameInputBackground:
				return overlay
			case theme.ColorNameForeground:
				return color.NRGBA{R: 15, G: 23, B: 42, A: 255}
			case theme.ColorNameForegroundOnPrimary:
				return color.NRGBA{R: 255, G: 255, B: 255, A: 255}
			case theme.ColorNamePrimary:
				return accent
			case theme.ColorNameSuccess:
				return color.NRGBA{R: 21, G: 128, B: 61, A: 255}
			case theme.ColorNameError:
				return color.NRGBA{R: 185, G: 28, B: 28, A: 255}
			case theme.ColorNameButton:
				return button
			case theme.ColorNameSeparator:
				return separator
			case theme.ColorNameHover:
				return color.NRGBA{R: 15, G: 23, B: 42, A: 40}
			default:
				return t.base.Color(name, theme.VariantLight)
			}
		}

	default: // "light" Açık Turkuaz (Gündüz Modu: #e2e8f0 Arka Plan, #06b6d4 Vurgu)
		switch name {
		case theme.ColorNameBackground:
			return color.NRGBA{R: 226, G: 232, B: 240, A: 255} // Soft light slate #e2e8f0
		case theme.ColorNameOverlayBackground:
			return color.NRGBA{R: 255, G: 255, B: 255, A: 255} // Crisp white cards
		case theme.ColorNameForeground:
			return color.NRGBA{R: 15, G: 23, B: 42, A: 255} // Dark slate text
		case theme.ColorNameForegroundOnPrimary:
			return color.NRGBA{R: 255, G: 255, B: 255, A: 255}
		case theme.ColorNamePrimary:
			return color.NRGBA{R: 6, G: 182, B: 212, A: 255} // Vibrant cyan #06b6d4
		case theme.ColorNameSuccess:
			return color.NRGBA{R: 22, G: 163, B: 74, A: 255}
		case theme.ColorNameError:
			return color.NRGBA{R: 220, G: 38, B: 38, A: 255}
		case theme.ColorNameButton:
			return color.NRGBA{R: 203, G: 213, B: 225, A: 255} // Soft button #cbd5e1
		case theme.ColorNameInputBackground:
			return color.NRGBA{R: 255, G: 255, B: 255, A: 255}
		case theme.ColorNameSeparator:
			return color.NRGBA{R: 148, G: 163, B: 184, A: 255}
		case theme.ColorNameHover:
			return color.NRGBA{R: 15, G: 23, B: 42, A: 40}
		default:
			return t.base.Color(name, theme.VariantLight)
		}
	}
}

func (t *TurkishTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return t.base.Icon(name)
}

func (t *TurkishTheme) Size(name fyne.ThemeSizeName) float32 {
	ratio := t.BaseSize / 15.0

	switch name {
	case theme.SizeNameText:
		return t.BaseSize
	case theme.SizeNameHeadingText:
		return 22 * ratio
	case theme.SizeNameSubHeadingText:
		return 18 * ratio
	case theme.SizeNameCaptionText:
		return 13 * ratio
	case theme.SizeNamePadding:
		return 8 * ratio
	case theme.SizeNameInlineIcon:
		return 22 * ratio
	default:
		return t.base.Size(name)
	}
}
