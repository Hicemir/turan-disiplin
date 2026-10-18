package theme

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"

	"okulevrak/assets"
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

func GetAllThemes(customColors ...color.Color) []ThemeInfo {
	customBg := color.Color(color.NRGBA{R: 30, G: 30, B: 46, A: 255})
	customAccent := color.Color(color.NRGBA{R: 137, G: 180, B: 250, A: 255})
	if len(customColors) >= 1 && customColors[0] != nil {
		customBg = customColors[0]
	}
	if len(customColors) >= 2 && customColors[1] != nil {
		customAccent = customColors[1]
	}

	return []ThemeInfo{
		// Koyu Temalar (Dark Themes)
		{ID: "dark", Name: "Koyu Gri / Mavi", BgColor: color.NRGBA{R: 15, G: 23, B: 42, A: 255}, AccentColor: color.NRGBA{R: 59, G: 130, B: 246, A: 255}},
		{ID: "dark_emerald", Name: "Koyu Gece / Zümrüt Yeşili", BgColor: color.NRGBA{R: 6, G: 26, B: 20, A: 255}, AccentColor: color.NRGBA{R: 16, G: 185, B: 129, A: 255}},
		{ID: "dark_purple", Name: "Koyu Gece / Ametist Mor", BgColor: color.NRGBA{R: 20, G: 9, B: 36, A: 255}, AccentColor: color.NRGBA{R: 168, G: 85, B: 247, A: 255}},
		{ID: "dark_amber", Name: "Koyu Gece / Kehribar Altın", BgColor: color.NRGBA{R: 24, G: 24, B: 27, A: 255}, AccentColor: color.NRGBA{R: 245, G: 158, B: 11, A: 255}},
		{ID: "dark_ruby", Name: "Koyu Gece / Yakut Kırmızı", BgColor: color.NRGBA{R: 31, G: 10, B: 16, A: 255}, AccentColor: color.NRGBA{R: 244, G: 63, B: 94, A: 255}},
		{ID: "dark_cyan", Name: "Koyu Gece / Buz Turkuaz", BgColor: color.NRGBA{R: 4, G: 24, B: 36, A: 255}, AccentColor: color.NRGBA{R: 6, G: 182, B: 212, A: 255}},

		// Açık Temalar (Light Themes)
		{ID: "light", Name: "Açık Tema", BgColor: color.NRGBA{R: 248, G: 250, B: 252, A: 255}, AccentColor: color.NRGBA{R: 29, G: 78, B: 216, A: 255}},
		{ID: "cream", Name: "Sarı / Krem", BgColor: color.NRGBA{R: 255, G: 253, B: 231, A: 255}, AccentColor: color.NRGBA{R: 2, G: 119, B: 189, A: 255}},
		{ID: "green", Name: "Yeşil / Doğa", BgColor: color.NRGBA{R: 240, G: 253, B: 244, A: 255}, AccentColor: color.NRGBA{R: 21, G: 128, B: 61, A: 255}},
		{ID: "bordeaux", Name: "Bordo / Klasik", BgColor: color.NRGBA{R: 250, G: 245, B: 240, A: 255}, AccentColor: color.NRGBA{R: 136, G: 19, B: 55, A: 255}},
		{ID: "navy", Name: "Lacivert / Kurumsal", BgColor: color.NRGBA{R: 241, G: 245, B: 249, A: 255}, AccentColor: color.NRGBA{R: 30, G: 58, B: 138, A: 255}},

		// Özel Kişiselleştirilmiş Tema
		{ID: "custom", Name: "🎨 Özel Tema (Kişiselleştirilmiş)", BgColor: customBg, AccentColor: customAccent},
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
	case "cream": // Warm nostalgia cream
		switch name {
		case theme.ColorNameBackground:
			return color.NRGBA{R: 255, G: 253, B: 231, A: 255} // Cream #fffde7
		case theme.ColorNameOverlayBackground:
			return color.NRGBA{R: 255, G: 250, B: 205, A: 255}
		case theme.ColorNameForeground:
			return color.NRGBA{R: 20, G: 20, B: 20, A: 255}
		case theme.ColorNameForegroundOnPrimary:
			return color.NRGBA{R: 255, G: 255, B: 255, A: 255}
		case theme.ColorNamePrimary:
			return color.NRGBA{R: 2, G: 119, B: 189, A: 255} // Deep blue
		case theme.ColorNameSuccess:
			return color.NRGBA{R: 46, G: 125, B: 50, A: 255}
		case theme.ColorNameError:
			return color.NRGBA{R: 198, G: 40, B: 40, A: 255}
		case theme.ColorNameButton:
			return color.NRGBA{R: 240, G: 235, B: 204, A: 255}
		case theme.ColorNameInputBackground:
			return color.NRGBA{R: 255, G: 255, B: 255, A: 255}
		case theme.ColorNameSeparator:
			return color.NRGBA{R: 218, G: 197, B: 130, A: 255}
		default:
			return t.base.Color(name, theme.VariantLight)
		}

	case "dark": // Deep Slate dark mode
		return buildDarkPalette(name,
			color.NRGBA{R: 15, G: 23, B: 42, A: 255},
			color.NRGBA{R: 30, G: 41, B: 59, A: 255},
			color.NRGBA{R: 51, G: 65, B: 85, A: 255},
			color.NRGBA{R: 59, G: 130, B: 246, A: 255},
			color.NRGBA{R: 71, G: 85, B: 105, A: 255},
			t.base)

	case "dark_emerald": // Deep Night Forest & Emerald
		return buildDarkPalette(name,
			color.NRGBA{R: 6, G: 26, B: 20, A: 255},
			color.NRGBA{R: 13, G: 40, B: 32, A: 255},
			color.NRGBA{R: 19, G: 56, B: 44, A: 255},
			color.NRGBA{R: 16, G: 185, B: 129, A: 255},
			color.NRGBA{R: 27, G: 77, B: 62, A: 255},
			t.base)

	case "dark_purple": // Deep Night Violet & Amethyst
		return buildDarkPalette(name,
			color.NRGBA{R: 20, G: 9, B: 36, A: 255},
			color.NRGBA{R: 33, G: 16, B: 58, A: 255},
			color.NRGBA{R: 48, G: 23, B: 84, A: 255},
			color.NRGBA{R: 168, G: 85, B: 247, A: 255},
			color.NRGBA{R: 67, G: 32, B: 117, A: 255},
			t.base)

	case "dark_amber": // Charcoal Pitch & Warm Amber
		return buildDarkPalette(name,
			color.NRGBA{R: 24, G: 24, B: 27, A: 255},
			color.NRGBA{R: 39, G: 39, B: 42, A: 255},
			color.NRGBA{R: 63, G: 63, B: 70, A: 255},
			color.NRGBA{R: 245, G: 158, B: 11, A: 255},
			color.NRGBA{R: 82, G: 82, B: 91, A: 255},
			t.base)

	case "dark_ruby": // Deep Crimson & Vivid Ruby
		return buildDarkPalette(name,
			color.NRGBA{R: 31, G: 10, B: 16, A: 255},
			color.NRGBA{R: 45, G: 16, B: 24, A: 255},
			color.NRGBA{R: 66, G: 23, B: 35, A: 255},
			color.NRGBA{R: 244, G: 63, B: 94, A: 255},
			color.NRGBA{R: 94, G: 33, B: 50, A: 255},
			t.base)

	case "dark_cyan": // Abyssal Deep Ocean & Cyan
		return buildDarkPalette(name,
			color.NRGBA{R: 4, G: 24, B: 36, A: 255},
			color.NRGBA{R: 8, G: 37, B: 54, A: 255},
			color.NRGBA{R: 12, G: 54, B: 78, A: 255},
			color.NRGBA{R: 6, G: 182, B: 212, A: 255},
			color.NRGBA{R: 17, G: 75, B: 109, A: 255},
			t.base)

	case "green": // Fresh Mint / Forest Emerald
		switch name {
		case theme.ColorNameBackground:
			return color.NRGBA{R: 240, G: 253, B: 244, A: 255} // Mint #f0fdf4
		case theme.ColorNameOverlayBackground:
			return color.NRGBA{R: 255, G: 255, B: 255, A: 255}
		case theme.ColorNameForeground:
			return color.NRGBA{R: 15, G: 23, B: 42, A: 255}
		case theme.ColorNameForegroundOnPrimary:
			return color.NRGBA{R: 255, G: 255, B: 255, A: 255}
		case theme.ColorNamePrimary:
			return color.NRGBA{R: 21, G: 128, B: 61, A: 255} // Emerald
		case theme.ColorNameSuccess:
			return color.NRGBA{R: 22, G: 163, B: 74, A: 255}
		case theme.ColorNameError:
			return color.NRGBA{R: 185, G: 28, B: 28, A: 255}
		case theme.ColorNameButton:
			return color.NRGBA{R: 220, G: 252, B: 231, A: 255}
		case theme.ColorNameInputBackground:
			return color.NRGBA{R: 255, G: 255, B: 255, A: 255}
		case theme.ColorNameSeparator:
			return color.NRGBA{R: 187, G: 247, B: 208, A: 255}
		default:
			return t.base.Color(name, theme.VariantLight)
		}

	case "bordeaux": // Warm Academic Crimson / Burgundy
		switch name {
		case theme.ColorNameBackground:
			return color.NRGBA{R: 250, G: 245, B: 240, A: 255} // Linen #faf5f0
		case theme.ColorNameOverlayBackground:
			return color.NRGBA{R: 255, G: 255, B: 255, A: 255}
		case theme.ColorNameForeground:
			return color.NRGBA{R: 24, G: 24, B: 27, A: 255}
		case theme.ColorNameForegroundOnPrimary:
			return color.NRGBA{R: 255, G: 255, B: 255, A: 255}
		case theme.ColorNamePrimary:
			return color.NRGBA{R: 136, G: 19, B: 55, A: 255} // Bordo #881337
		case theme.ColorNameSuccess:
			return color.NRGBA{R: 21, G: 128, B: 61, A: 255}
		case theme.ColorNameError:
			return color.NRGBA{R: 185, G: 28, B: 28, A: 255}
		case theme.ColorNameButton:
			return color.NRGBA{R: 254, G: 226, B: 226, A: 255}
		case theme.ColorNameInputBackground:
			return color.NRGBA{R: 255, G: 255, B: 255, A: 255}
		case theme.ColorNameSeparator:
			return color.NRGBA{R: 254, G: 205, B: 211, A: 255}
		default:
			return t.base.Color(name, theme.VariantLight)
		}

	case "navy": // Ice cool slate & Corporate Navy
		switch name {
		case theme.ColorNameBackground:
			return color.NRGBA{R: 241, G: 245, B: 249, A: 255} // Ice cool #f1f5f9
		case theme.ColorNameOverlayBackground:
			return color.NRGBA{R: 255, G: 255, B: 255, A: 255}
		case theme.ColorNameForeground:
			return color.NRGBA{R: 15, G: 23, B: 42, A: 255}
		case theme.ColorNameForegroundOnPrimary:
			return color.NRGBA{R: 255, G: 255, B: 255, A: 255}
		case theme.ColorNamePrimary:
			return color.NRGBA{R: 30, G: 58, B: 138, A: 255} // Navy #1e3a8a
		case theme.ColorNameSuccess:
			return color.NRGBA{R: 21, G: 128, B: 61, A: 255}
		case theme.ColorNameError:
			return color.NRGBA{R: 185, G: 28, B: 28, A: 255}
		case theme.ColorNameButton:
			return color.NRGBA{R: 226, G: 232, B: 240, A: 255}
		case theme.ColorNameInputBackground:
			return color.NRGBA{R: 255, G: 255, B: 255, A: 255}
		case theme.ColorNameSeparator:
			return color.NRGBA{R: 203, G: 213, B: 225, A: 255}
		default:
			return t.base.Color(name, theme.VariantLight)
		}

	case "custom":
		bg := t.CustomBg
		if bg == nil {
			bg = color.NRGBA{R: 30, G: 30, B: 46, A: 255}
		}
		accent := t.CustomAccent
		if accent == nil {
			accent = color.NRGBA{R: 137, G: 180, B: 250, A: 255}
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
			default:
				return t.base.Color(name, theme.VariantLight)
			}
		}

	default: // "light" Clean daylight slate
		switch name {
		case theme.ColorNameBackground:
			return color.NRGBA{R: 248, G: 250, B: 252, A: 255}
		case theme.ColorNameOverlayBackground:
			return color.NRGBA{R: 255, G: 255, B: 255, A: 255}
		case theme.ColorNameForeground:
			return color.NRGBA{R: 15, G: 23, B: 42, A: 255}
		case theme.ColorNameForegroundOnPrimary:
			return color.NRGBA{R: 255, G: 255, B: 255, A: 255}
		case theme.ColorNamePrimary:
			return color.NRGBA{R: 29, G: 78, B: 216, A: 255}
		case theme.ColorNameSuccess:
			return color.NRGBA{R: 21, G: 128, B: 61, A: 255}
		case theme.ColorNameError:
			return color.NRGBA{R: 185, G: 28, B: 28, A: 255}
		case theme.ColorNameButton:
			return color.NRGBA{R: 226, G: 232, B: 240, A: 255}
		case theme.ColorNameInputBackground:
			return color.NRGBA{R: 255, G: 255, B: 255, A: 255}
		case theme.ColorNameSeparator:
			return color.NRGBA{R: 203, G: 213, B: 225, A: 255}
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
