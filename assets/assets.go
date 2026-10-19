package assets

import (
	_ "embed"

	"fyne.io/fyne/v2"
)

//go:embed logo.png
var LogoPNG []byte

var AppIcon = fyne.NewStaticResource("logo.png", LogoPNG)

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
