package main

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"

	"github.com/lestrrat-go/dither"
	"github.com/nfnt/resize"
	"github.com/pkg/errors"
)

const (
	maxWidth = 576
)

func printImage(file io.Reader) error {
	wr := printerDev

	// decode jpeg into image.Image
	img, _, err := image.Decode(file)
	if err != nil {
		return errors.Wrap(err, "fail to print image")
	}
	origW, origH := img.Bounds().Dx(), img.Bounds().Dy()

	var w, h int
	if origW < maxWidth {
		w = 576 // maxWidth
		h = ((origH * w) / origW) / 3
		img = resize.Resize(uint(w), uint(h), img, resize.Lanczos3)
	}

	ditheredImg := dither.Monochrome(dither.Burkes, img, 1.18)
	dataBuf := make([]byte, (3*w*h+7)/8)

	// 가로방향 점의 개수: nL + nH x 256
	nH := byte(w / 256)
	nL := byte(w % 256)
	mode := byte(33)
	log.Println(nL, nH, mode)
	cmdBuf := []byte{0x1B, 0x2A, mode, nL, nH}

	dataBufIdx := 0
	for y := 0; y < h; y += 8 * 3 {
		for x := 0; x < w; x++ {
			var dataByte byte
			for yCnt := 0; yCnt < 3; yCnt++ {
				for yi := 0; yi < 8; yi++ {
					currY := y + (8 * yCnt) + yi
					if currY > h {
						continue
					}
					var bit byte
					if r, g, b, _ := ditheredImg.At(x, currY).RGBA(); r == 0 && g == 0 && b == 0 {
						bit = 1 << (7 - yi)
					}
					dataByte |= bit
				}
				dataBuf[dataBufIdx] = dataByte
				dataBufIdx += 1
			}
		}
	}

	// 가운데 정렬
	wr.Write([]byte{0x1B, 0x61, 1})

	// Line spacing
	wr.Write([]byte{0x1B, 0x33, 0})

	for i := 0; i < len(dataBuf); i += (w * 3) {
		printBuf := append(cmdBuf, dataBuf[i:i+w*3]...)
		wr.Write(printBuf)
	}

	// Paper cut
	wr.Write([]byte("\x1B@\x1DVA0"))
	return nil
}