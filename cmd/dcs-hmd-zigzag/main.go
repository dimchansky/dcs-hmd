package main

import (
	"fmt"
	"log"
	"math"
	"net"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() (err error) {
	udpAddr, err := net.ResolveUDPAddr("udp", ":19089")
	if err != nil {
		return err
	}

	udpConn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return err
	}

	defer func() {
		lErr := udpConn.Close()
		if err == nil {
			err = lErr
		}
	}()

	rotorRPMWave := triangleWave(11.0)
	rotorPitchWave := triangleWave(5.0)
	start := time.Now()
	unixTs := start.Unix()

	for {
		t := time.Since(start).Seconds()

		rotorRPMVal := rotorRPMWave.Value(t)
		rotorPitchVal := rotorPitchWave.Value(t)

		toSend := fmt.Sprintf("%08x*52=%0.4f:53=%0.4f\n", unixTs, rotorRPMVal, rotorPitchVal)
		_, _ = udpConn.Write([]byte(toSend))

		time.Sleep(16 * time.Millisecond)
	}
}

// triangleWave is a triangle wave of period p that spans the range [0,1].
type triangleWave float64

// Value returns value of triangle wave at time t.
func (p triangleWave) Value(t float64) float64 {
	tDivP := t / float64(p)
	return 2.0 * math.Abs(tDivP-math.Floor(tDivP+0.5))
}
