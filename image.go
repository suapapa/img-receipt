package main

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"

	"github.com/kenshaw/escpos"
	"github.com/lestrrat-go/dither"
	"github.com/nfnt/resize"
	"github.com/pkg/errors"
)

func openImg(path string) (*image.Rectangle, []byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, errors.Wrap(err, "fail to load image")
	}

	// decode jpeg into image.Image
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, nil, errors.Wrap(err, "fail to load image")
	}
	file.Close()

	m := resize.Resize(512, 0, img, resize.Lanczos3)

	ditheredImg := dither.Monochrome(dither.Burkes, m, 1.18)
	rect := ditheredImg.Bounds()
	w, h := rect.Dx(), rect.Dy()
	buf := make([]byte, w*h/8)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			bytePos := (w/8)*y + (x / 8)
			pointBit := 7 - (x % 8)
			var bit byte
			if r, g, b, _ := ditheredImg.At(x, y).RGBA(); r == 0 && g == 0 && b == 0 {
				bit = 1
			}
			buf[bytePos] = buf[bytePos] | (bit << pointBit)
			log.Println(x, y, bytePos, pointBit, bit, ditheredImg.At(x, y))
		}
	}

	return &rect, buf, nil
}

func imageDemo(p *escpos.Escpos) {
	r, buf, err := openImg("Lenna.png")
	if err != nil {
		log.Fatal(err)
	}
	// p.Init()
	log.Println(buf)
	log.Println(r.Dx(), r.Dy())
	p.Raster(r.Dx(), r.Dy(), r.Dx()/8, buf)
	p.Cut()
	p.End()
}
