package apiserver

import (
	"encoding/base64"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"log"
	"math"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

var (
	dpi      float64 = 72
	fontfile string  = "3952.ttf"
	hinting  string  = "none"
	size     float64 = 12
	spacing  float64 = 1.5
	wonb     bool    = false
)

const title = "Jabberwocky"

func renderImage(text string) (string, error) {

	// Read the font data.
	fontBytes, err := ioutil.ReadFile(fontfile)
	if err != nil {
		log.Println(err)
		return "", err
	}
	f, err := truetype.Parse(fontBytes)
	if err != nil {
		log.Println(err)
		return "", err
	}

	// Draw the background and the guidelines.
	fg, bg := image.Black, image.White
	ruler := color.RGBA{0xdd, 0xdd, 0xdd, 0xff}
	if wonb {
		fg, bg = image.White, image.Black
		ruler = color.RGBA{0x22, 0x22, 0x22, 0xff}
	}
	const imgW, imgH = 640, 480
	rgba := image.NewRGBA(image.Rect(0, 0, imgW, imgH))
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
	for i := 0; i < 200; i++ {
		rgba.Set(10, 10+i, ruler)
		rgba.Set(10+i, 10, ruler)
	}

	// Draw the text.
	h := font.HintingNone
	switch hinting {
	case "full":
		h = font.HintingFull
	}
	d := &font.Drawer{
		Dst: rgba,
		Src: fg,
		Face: truetype.NewFace(f, &truetype.Options{
			Size:    size,
			DPI:     dpi,
			Hinting: h,
		}),
	}
	y := 10 + int(math.Ceil(size*dpi/72))
	// dy := int(math.Ceil(size * spacing * dpi / 72))
	d.Dot = fixed.Point26_6{
		X: (fixed.I(imgW) - d.MeasureString(title)) / 2,
		Y: fixed.I(y),
	}
	d.DrawString(title)

	d.Dot = fixed.P(10, y)
	d.DrawString(text)

	encoded := base64.StdEncoding.EncodeToString([]byte(rgba.Pix))
	return encoded, nil

}

/*
// Save that RGBA image to disk.
outFile, err := os.Create("out.png")
if err != nil {
	log.Println(err)
	os.Exit(1)
}
defer outFile.Close()
b := bufio.NewWriter(outFile)
err = png.Encode(b, rgba)
if err != nil {
	log.Println(err)
	os.Exit(1)
}
err = b.Flush()
if err != nil {
	log.Println(err)
	os.Exit(1)
}
fmt.Println("Wrote out.png OK.")
*/
