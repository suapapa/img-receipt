package main

import (
	"bytes"
	"image"
	"log"
	"os"
	"strings"

	_ "image/png"

	"github.com/kenshaw/escpos"
	"github.com/qeesung/image2ascii/convert"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/transform"
)

func printBasictDemo(p *escpos.Escpos) {
	testString := "안녕하세요"

	var bufs bytes.Buffer
	wr := transform.NewWriter(&bufs, korean.EUCKR.NewEncoder())
	wr.Write([]byte(testString))
	wr.Close()

	testBytes := bufs.Bytes()

	p.Init()
	p.SetSmooth(1)
	p.SetFontSize(2, 3)
	p.SetFont("A")
	p.Write("글꼴A-2x3")
	p.Formfeed()

	p.SetFont("B")
	p.WriteRaw(strToEuckrBytes("글꼴B-2x3"))
	p.Formfeed()

	p.SetFont("A")
	p.SetFontSize(1, 1)
	p.Write("123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890")
	p.Formfeed()

	p.SetFont("B")
	p.SetFontSize(1, 1)
	p.Write("123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890")
	p.Formfeed()

	p.SetEmphasize(1)
	p.WriteRaw(strToEuckrBytes("강조"))
	p.Formfeed()

	p.SetUnderline(1)
	p.SetFontSize(3, 3)
	p.WriteRaw(strToEuckrBytes("밑줄-3x3"))
	p.Formfeed()

	p.SetFont("C")
	p.SetFontSize(2, 4)
	p.WriteRaw(strToEuckrBytes("글꼴C-2x4"))
	p.Formfeed()

	p.SetFontSize(8, 8)
	p.WriteRaw(testBytes)
	p.FormfeedN(5)

	convertOptions := convert.Options{
		FixedWidth:  59,
		FixedHeight: 32,
		Reversed:    true,
	}
	p.SetFont("B")
	p.SetFontSize(1, 1)
	p.SetEmphasize(1)
	p.SetUnderline(0)

	file, err := os.Open("Lenna.png")
	if err != nil {
		log.Fatal(err)
	}

	// decode jpeg into image.Image
	imgLenna, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	file.Close()
	converter := convert.NewImageConverter()
	strs := converter.Image2ASCIIString(imgLenna, &convertOptions)
	log.Println(strs)

	for _, s := range strings.Split(strs, "\n") {
		p.Write(s)
		p.Formfeed()
	}

	p.Cut()
	p.End()
}

func strToEuckrBytes(str string) []byte {
	var bufs bytes.Buffer
	wr := transform.NewWriter(&bufs, korean.EUCKR.NewEncoder())
	wr.Write([]byte(str))
	wr.Close()

	return bufs.Bytes()
}
