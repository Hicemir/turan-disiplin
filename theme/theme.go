package theme

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"

	"okulevrak/assets"
)

type TurkishTheme struct {
	base     fyne.Theme
	Mode     string  // "light", "cream", "dark", "green", "bordeaux", "navy"
	BaseSize float32 // 15, 17, 19
}

type ThemeInfo struct {
	ID          string
	Name        string
	BgColor     color.Color
	AccentColor color.Color
}

func GetAllThemes() []ThemeInfo {
	return []ThemeInfo{
		{ID: "light", Name: "Açık Tema", BgColor: color.NRGBA{R: 248, G: 250, B: 252, A: 255}, AccentColor: color.NRGBA{R: 29, G: 78, B: 216, A: 255}},
		{ID: "cream", Name: "Sarı / Krem", BgColor: color.NRGBA{R: 255, G: 253, B: 231, A: 255}, AccentColor: color.NRGBA{R: 2, G: 119, B: 189, A: 255}},
		{ID: "dark", Name: "Koyu Tema", BgColor: color.NRGBA{R: 15, G: 23, B: 42, A: 255}, AccentColor: color.NRGBA{R: 59, G: 130, B: 246, A: 255}},
		{ID: "green", Name: "Yeşil / Doğa", BgColor: color.NRGBA{R: 240, G: 253, B: 244, A: 255}, AccentColor: color.NRGBA{R: 21, G: 128, B: 61, A: 255}},
		{ID: "bordeaux", Name: "Bordo / Klasik", BgColor: color.NRGBA{R: 250, G: 245, B: 240, A: 255}, AccentColor: color.NRGBA{R: 136, G: 19, B: 55, A: 255}},
		{ID: "navy", Name: "Lacivert / Kurumsal", BgColor: color.NRGBA{R: 241, G: 245, B: 249, A: 255}, AccentColor: color.NRGBA{R: 30, G: 58, B: 138, A: 255}},
	}
}

func NewTurkishTheme(mode string, fontSize float32) fyne.Theme {
	if mode == "" {
		mode = "light"
	}
	if fontSize < 13 {
		fontSize = 15.0
	}
	return &TurkishTheme{
		base:     theme.DefaultTheme(),
		Mode:     mode,
		BaseSize: fontSize,
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
		switch name {
		case theme.ColorNameBackground:
			return color.NRGBA{R: 15, G: 23, B: 42, A: 255} // Slate #0f172a
		case theme.ColorNameOverlayBackground:
			return color.NRGBA{R: 30, G: 41, B: 59, A: 255}
		case theme.ColorNameForeground:
			return color.NRGBA{R: 248, G: 250, B: 252, A: 255}
		case theme.ColorNameForegroundOnPrimary:
			return color.NRGBA{R: 255, G: 255, B: 255, A: 255}
		case theme.ColorNamePrimary:
			return color.NRGBA{R: 59, G: 130, B: 246, A: 255} // Bright royal blue
		case theme.ColorNameSuccess:
			return color.NRGBA{R: 22, G: 163, B: 74, A: 255}
		case theme.ColorNameError:
			return color.NRGBA{R: 220, G: 38, B: 38, A: 255}
		case theme.ColorNameButton:
			return color.NRGBA{R: 51, G: 65, B: 85, A: 255}
		case theme.ColorNameInputBackground:
			return color.NRGBA{R: 30, G: 41, B: 59, A: 255}
		case theme.ColorNameSeparator:
			return color.NRGBA{R: 71, G: 85, B: 105, A: 255}
		default:
			return t.base.Color(name, theme.VariantDark)
		}

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
