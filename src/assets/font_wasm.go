//go:build wasm

package assets

import (
	"github.com/quasilyte/bitsweetfont"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var (
	Font1, Font2, Font3, Font1_3 font.Face
)

func init() {
	// TODO: more lange
	initFont("ch")
}

func initFont(lang string) {
	var ttfData []byte
	switch lang {
	case "ch":
		ttf, err := gameAssets.ReadFile("_data/ttf/zh_simhei.ttf")
		if err != nil {
			initFont("default")
			return
		}
		ttfData = ttf
	default:
		Font1 = bitsweetfont.New1()
		Font2 = bitsweetfont.Scale(Font1, 2)
		Font3 = bitsweetfont.Scale(Font1, 3)

		Font1_3 = bitsweetfont.New1_3()
		return
	}

	gameFont := parseMustFont(ttfData)
	Font1 = newMustFace(gameFont, &opentype.FaceOptions{
		Size: 12,
		DPI:  72,
		//Hinting: font.HintingFull,
	})
	Font2 = newMustFace(gameFont, &opentype.FaceOptions{
		Size: 24,
		DPI:  72,
		//Hinting: font.HintingFull,
	})
	Font3 = newMustFace(gameFont, &opentype.FaceOptions{
		Size: 32,
		DPI:  72,
		//Hinting: font.HintingFull,
	})

	Font1_3 = newMustFace(gameFont, &opentype.FaceOptions{
		Size: 40,
		DPI:  72,
		//Hinting: font.HintingFull,
	})
}

func parseMustFont(ttf []byte) *opentype.Font {
	font, err := opentype.Parse(ttf)
	if err != nil {
		panic(err)
	}
	return font
}

func newMustFace(f *opentype.Font, opt *opentype.FaceOptions) font.Face {
	face, err := opentype.NewFace(f, opt)
	if err != nil {
		panic(err)
	}
	return face
}
