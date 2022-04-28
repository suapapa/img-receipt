package main

import (
	"log"
	"os"

	"github.com/tarm/serial"
)

func main() {
	port := "/dev/ttyACM0"
	c := &serial.Config{Name: port, Baud: 115200}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	f, _ := os.Open("_img/Lenna.png")
	defer f.Close()
	printImage(s, f)
	s.Flush()
}
