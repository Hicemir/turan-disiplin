package assets

import (
	_ "embed"

	"fyne.io/fyne/v2"
)

//go:embed logo.png
var LogoPNG []byte

var AppIcon = fyne.NewStaticResource("logo.png", LogoPNG)

//go:embed brand_logo_dark.png
var BrandLogoDarkPNG []byte

var BrandLogoDark = fyne.NewStaticResource("brand_logo_dark.png", BrandLogoDarkPNG)

//go:embed brand_logo_light.png
var BrandLogoLightPNG []byte

var BrandLogoLight = fyne.NewStaticResource("brand_logo_light.png", BrandLogoLightPNG)

// GetBrandLogo returns the brand logo matching current theme mode
func GetBrandLogo(isDark bool) fyne.Resource {
	if isDark {
		return BrandLogoDark
	}
	return BrandLogoLight
}

//go:embed yonetmelik.json
var YonetmelikJSON []byte

//go:embed fonts/LiberationSans-Regular.ttf
var FontRegular []byte

//go:embed fonts/LiberationSans-Bold.ttf
var FontBold []byte

//go:embed fonts/LiberationSans-Italic.ttf
var FontItalic []byte

//go:embed fonts/LiberationSans-BoldItalic.ttf
var FontBoldItalic []byte
