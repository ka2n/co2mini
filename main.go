package co2mini

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/zserge/hid"
)

const (
	vendorID = "04d9:a052:0100:00"
	co2op    = 0x50
	tempop   = 0x42
)

var (
	key = []byte{0x86, 0x41, 0xc9, 0xa8, 0x7f, 0x41, 0x3c, 0xac}
)

// Value represents mesurements result from device
type Value struct {
	Temp *float64 `json:"temp"`
	CO2  *int     `json:"co2"`
}

// FindDevice from connected devices
func FindDevice() *hid.Device {
	var dev *hid.Device
	hid.UsbWalk(func(device hid.Device) {
		info := device.Info()
		id := fmt.Sprintf("%04x:%04x:%04x:%02x", info.Vendor, info.Product, info.Revision, info.Interface)
		if id != vendorID {
			return
		}
		dev = &device
	})
	return dev
}

// Oneshot read
func Oneshot(device hid.Device, output OutputWriter) error {
	ctx, cancel := context.WithCancel(context.Background())
	recv := make(chan Value, 1)
	var v Value

	go func() {
		for vv := range recv {
			if vv.CO2 != nil {
				v.CO2 = vv.CO2
			}
			if vv.Temp != nil {
				v.Temp = vv.Temp
			}
			if v.Temp != nil && v.CO2 != nil {
				break
			}
		}
		cancel()
	}()

	if err := readMeterValue(ctx, device, recv); err != nil {
		return err
	}

	return output.Write(v)
}

// Watch device values
func Watch(device hid.Device, output OutputWriter) error {
	ctx := context.Background()
	recv := make(chan Value, 1)

	go func() {
		for v := range recv {
			output.Write(v)
		}
	}()

	if err := readMeterValue(ctx, device, recv); err != nil {
		return err
	}
	return nil
}

func readMeterValue(ctx context.Context, device hid.Device, result chan<- Value) error {
	if err := device.Open(); err != nil {
		log.Println("Open error: ", err)
		return err
	}
	defer device.Close()
	if err := device.SetReport(0, key); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if buf, err := device.Read(-1, 1*time.Second); err == nil {
				dec := decrypt(buf, key)
				if len(dec) == 0 {
					continue
				}
				val := int(dec[1])<<8 | int(dec[2])
				if dec[0] == co2op {
					result <- Value{CO2: &val}
				}
				if dec[0] == tempop {
					temp := float64(val)/16.0 - 273.15
					result <- Value{Temp: &temp}
				}
			}
		}
	}
}
