package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/tarm/serial"
)

var (
	programName                   = "img-receipt"
	buildStamp, gitHash, buildTag string

	printerDev    *bufio.Writer
	flagSerialDev string
	flagUsbDev    string
	flagEnBle     bool
	flagAdvDu     int
)

func main() {
	log.Printf("%s-%s-%s(%s)", programName, buildTag, gitHash, buildStamp)

	flag.StringVar(&flagSerialDev, "s", "/dev/ttyACM0", "serial device")
	flag.StringVar(&flagUsbDev, "u", "", "if specify usb lp device -s will be ignored")
	flag.BoolVar(&flagEnBle, "b", false, "enable ble server")
	flag.IntVar(&flagAdvDu, "a", 0, "ble advertisement duration")
	flag.Parse()

	if flagUsbDev != "" {
		// /dev/usb/lp0
		dev, err := os.OpenFile(flagUsbDev, os.O_RDWR, 0)
		printerDev = bufio.NewWriter(dev)
		if err != nil {
			log.Fatal(err)
		}
		defer dev.Close()
	} else {
		c := &serial.Config{Name: flagSerialDev, Baud: 9600}
		dev, err := serial.OpenPort(c)
		if err != nil {
			log.Fatal(err)
		}
		printerDev = bufio.NewWriter(dev)
		defer dev.Close()
	}

	r := gin.Default()
	// r.SetTrustedProxies([]string{"192.168.1.2"})
	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	r.POST("/upload", uploadHandler)
	r.POST("/qr", qrHandler)
	r.POST("/cut", cutHandler)

	go r.Run(":8080")
	if flagEnBle {
		go bleServer(time.Duration(flagAdvDu) * time.Second)
	}

	stopC := make(chan interface{})
	<-stopC
}

func uploadHandler(c *gin.Context) {
	file, _, err := c.Request.FormFile("img")
	if err != nil {
		c.Error(errors.Wrap(err, "fail to print"))
	}
	defer file.Close()

	dpi := c.Query("dpi")
	if dpi == "200" {
		if err := printImage24bitDouble(file); err != nil {
			c.Error(errors.Wrap(err, "fail to print"))
		}
	} else {
		if err := printImage8bitDouble(file); err != nil {
			c.Error(errors.Wrap(err, "fail to print"))
		}
	}
	nocut := c.Query("nocut")
	if nocut == "1" || nocut == "true" {
		return
	}
	cutPaper()
}

func cutHandler(c *gin.Context) {
	cutPaper()
}
