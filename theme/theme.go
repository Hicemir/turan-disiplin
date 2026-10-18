package theme

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"

	"okulevrak/assets"
)

type TurkishTheme struct {
	base      fyne.Theme
	Mode      string  // "light", "cream", "dark"
	BaseSize  float32 // 15, 17, 19
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
	case "cream": // 1:1 Classic Original .NET aesthetic from 1.png
		switch name {
		case theme.ColorNameBackground:
			return color.NRGBA{R: 255, G: 253, B: 231, A: 255} // Cream #fffde7
		case theme.ColorNameOverlayBackground:
			return color.NRGBA{R: 255, G: 250, B: 205, A: 255} // Light cream card #fffacd
		case theme.ColorNameForeground:
			return color.NRGBA{R: 20, G: 20, B: 20, A: 255} // High contrast black
		case theme.ColorNameForegroundOnPrimary:
			return color.NRGBA{R: 255, G: 255, B: 255, A: 255} // White text on buttons
		case theme.ColorNamePrimary:
			return color.NRGBA{R: 2, G: 119, B: 189, A: 255} // Deep blue
		case theme.ColorNameSuccess:
			return color.NRGBA{R: 46, G: 125, B: 50, A: 255} // Deep green
		case theme.ColorNameError:
			return color.NRGBA{R: 198, G: 40, B: 40, A: 255} // Red
		case theme.ColorNameButton:
			return color.NRGBA{R: 240, G: 235, B: 204, A: 255} // Classic Windows warm button
		case theme.ColorNameInputBackground:
			return color.NRGBA{R: 255, G: 255, B: 255, A: 255} // Crisp white
		case theme.ColorNameSeparator:
			return color.NRGBA{R: 218, G: 197, B: 130, A: 255} // Golden border
		default:
			return t.base.Color(name, theme.VariantLight)
		}

	case "dark":
		switch name {
		case theme.ColorNameBackground:
			return color.NRGBA{R: 15, G: 23, B: 42, A: 255} // Deep slate #0f172a
		case theme.ColorNameOverlayBackground:
			return color.NRGBA{R: 30, G: 41, B: 59, A: 255} // Slate #1e293b
		case theme.ColorNameForeground:
			return color.NRGBA{R: 248, G: 250, B: 252, A: 255} // Crisp white
		case theme.ColorNameForegroundOnPrimary:
			return color.NRGBA{R: 255, G: 255, B: 255, A: 255}
		case theme.ColorNamePrimary:
			return color.NRGBA{R: 37, G: 99, B: 235, A: 255} // Royal blue
		case theme.ColorNameSuccess:
			return color.NRGBA{R: 22, G: 163, B: 74, A: 255}
		case theme.ColorNameError:
			return color.NRGBA{R: 220, G: 38, B: 38, A: 255}
		case theme.ColorNameButton:
			return color.NRGBA{R: 51, G: 65, B: 85, A: 255} // Dark slate button (NEVER white on white!)
		case theme.ColorNameInputBackground:
			return color.NRGBA{R: 30, G: 41, B: 59, A: 255}
		case theme.ColorNameSeparator:
			return color.NRGBA{R: 71, G: 85, B: 105, A: 255}
		default:
			return t.base.Color(name, theme.VariantDark)
		}

	default: // "light" (Senior-friendly high contrast daylight)
		switch name {
		case theme.ColorNameBackground:
			return color.NRGBA{R: 248, G: 250, B: 252, A: 255} // Clean soft daylight
		case theme.ColorNameOverlayBackground:
			return color.NRGBA{R: 255, G: 255, B: 255, A: 255} // Pure white card
		case theme.ColorNameForeground:
			return color.NRGBA{R: 15, G: 23, B: 42, A: 255} // Ultra high-contrast black/slate
		case theme.ColorNameForegroundOnPrimary:
			return color.NRGBA{R: 255, G: 255, B: 255, A: 255}
		case theme.ColorNamePrimary:
			return color.NRGBA{R: 29, G: 78, B: 216, A: 255} // Strong blue
		case theme.ColorNameSuccess:
			return color.NRGBA{R: 21, G: 128, B: 61, A: 255} // High contrast green
		case theme.ColorNameError:
			return color.NRGBA{R: 185, G: 28, B: 28, A: 255} // High contrast red
		case theme.ColorNameButton:
			return color.NRGBA{R: 226, G: 232, B: 240, A: 255} // Clear slate-200 button
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
		return t.BaseSize // e.g. 15, 17, 19
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
