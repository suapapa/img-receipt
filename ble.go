package main

import (
	"bytes"
	"context"
	"log"
	"time"

	"github.com/go-ble/ble"
	ble_linux "github.com/go-ble/ble/linux"
	"github.com/pkg/errors"
)

var (
	imgReceiptSvcUUID           = ble.MustParse("b6c5fd61-1702-462d-beb6-d63ac3611e50")
	imgReceiptDataWriteCharUUID = ble.MustParse("b6c5fd62-1702-462d-beb6-d63ac3611e50")
)

func bleServer(advDu time.Duration) error {
	d, err := ble_linux.NewDeviceWithName("imgReceipt")
	if err != nil {
		return err
	}

	ble.SetDefaultDevice(d)
	if err := ble.RemoveAllServices(); err != nil {
		return err
	}

	imgReceiptSvc := ble.NewService(imgReceiptSvcUUID)
	imgReceiptSvc.AddCharacteristic(imgReceiptDataWriteChar())

	if err := ble.AddService(imgReceiptSvc); err != nil {
		return err
	}

	log.Printf("advertising for %v secs...\n", advDu)
	var bleCtx context.Context
	if advDu > 0 {
		bleCtx = ble.WithSigHandler(context.WithTimeout(context.Background(), advDu))
	} else {
		bleCtx = ble.WithSigHandler(context.WithCancel(context.Background()))
	}
	chkErr(ble.AdvertiseNameAndServices(bleCtx, "imgReceipt", imgReceiptSvc.UUID))

	return nil
}

func imgReceiptDataWriteChar() *ble.Characteristic {
	c := ble.NewCharacteristic(imgReceiptDataWriteCharUUID)
	c.HandleWrite(ble.WriteHandlerFunc(func(req ble.Request, rsp ble.ResponseWriter) {
		data := req.Data()
		log.Printf("onDataWrite: receivced %v bytes\n", len(data))
		byteBuff := bytes.NewBuffer(data)

		err := printImage(byteBuff)
		chkErr(err)
	}))
	return c
}

func chkErr(err error) {
	switch errors.Cause(err) {
	case nil:
	case context.DeadlineExceeded:
		log.Printf("done\n")
	case context.Canceled:
		log.Printf("canceled\n")
	default:
		log.Fatalf(err.Error())
	}
}
