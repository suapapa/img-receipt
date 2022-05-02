package main

import (
	"flag"
	"io"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tarm/serial"
)

var (
	printerDev    io.Writer
	flagSerialDev string
	flagUsbDev    string
	flagEnBle     bool
	flagAdvDu     int
)

func main() {
	flag.StringVar(&flagSerialDev, "s", "/dev/ttyACM0", "serial device")
	flag.StringVar(&flagUsbDev, "u", "", "if specify usb lp device -s will be ignored")
	flag.BoolVar(&flagEnBle, "b", false, "enable ble server")
	flag.IntVar(&flagAdvDu, "a", 0, "ble advertisement duration")
	flag.Parse()

	var err error

	if flagUsbDev != "" {
		printerDev, err = os.Open(flagUsbDev)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		c := &serial.Config{Name: flagSerialDev, Baud: 9600}
		printerDev, err = serial.OpenPort(c)
		if err != nil {
			log.Fatal(err)
		}
	}

	r := gin.Default()
	// r.SetTrustedProxies([]string{"192.168.1.2"})
	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	r.POST("/upload", uploadHandler)

	go r.Run(":8080")
	if flagEnBle {
		go bleServer(time.Duration(flagAdvDu) * time.Second)
	}

	stopC := make(chan interface{})
	<-stopC
}

func uploadHandler(c *gin.Context) {
	file, _, _ := c.Request.FormFile("img")
	defer file.Close()

	printImage(file)
}
