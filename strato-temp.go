package main

import (
	"flag"
	"fmt"
	"github.com/stratoberry/go-bmp085"
	"log"
	"net"
	"os"
	"time"
)

var outFile = flag.String("output", "/data/temperature.csv", "output filename")
var udpPort = flag.Int("udp", 8367, "UDP port for sending data")
var updateFreq = flag.Duration("freq", 10, "update frequency in seconds")

// Feb 1st 2014
const EPOCH = 1391212800

func main() {
	flag.Parse()

	// todo output to a file
	//log.SetOutput("/data/logs/strato-temperature.log")
	var dev *bmp085.Device
	var err error

	if dev, err = bmp085.Init(0x77, 1, bmp085.MODE_STANDARD); err != nil {
		log.Panicln("Failed to init device", err)
	}

	f, err := os.OpenFile(*outFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Panicln("Failed to open output file", err)
	}
	defer f.Close()

	conn, _ := net.Dial("udp", fmt.Sprintf("255.255.255.255:%d", *udpPort))
	defer conn.Close()

	t := time.Tick((*updateFreq) * time.Second)
	for now := range t {
		if temp, pressure, alt, err := dev.GetData(); err != nil {
			log.Panicln("Failed to get data from the device", err)
		} else {
			stanza := fmt.Sprintf("%d;%.2f;%.2f;%.2f\n", now.Unix()-EPOCH, temp, float64(pressure)/100, alt)

			f.WriteString(stanza)
			conn.Write([]byte(stanza))
		}
	}
}
