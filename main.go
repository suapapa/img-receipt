package main

import (
	"log"

	"github.com/tarm/serial"
)

func main() {
	port := "/dev/ttyACM0"
	c := &serial.Config{Name: port, Baud: 115200}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	printImage(s, "Lenna.png")
	s.Flush()

	// p := escpos.New(s)
	// // // printBasictDemo(p)
	// p.Init()
	// p.Cut()

}
