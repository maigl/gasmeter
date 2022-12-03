package main

import (
	"fmt"
	"time"

	rpio "github.com/stianeikeland/go-rpio/v4"
)

type PulseMeter interface {
	CountPulse()
	Reading() float64
	UpdateValue(float64)
	Close()
}

type Gasmeter struct {
	data100 int
	pin     rpio.Pin
}

func NewGasmeter(fake bool) *Gasmeter {
	g := &Gasmeter{}

	if fake {
		go func() {
			for {
				fmt.Println("fake pulse")
				g.CountPulse()
				// wait random time between 1 and 5 seconds
				time.Sleep(time.Duration(1+int64(4)*time.Second.Nanoseconds()) * time.Nanosecond)
			}
		}()
	} else {
		err := rpio.Open()
		if err != nil {
			panic(err)
		}
		g.pin = rpio.Pin(27)
		// of course, we actually want to read the input
		g.pin.Input()
		g.pin.PullDown()

		var last rpio.State
		go func() {
			for {
				current := g.pin.Read()
				// detect falling edge
				if current == rpio.Low && last == rpio.High {
					g.CountPulse()
					fmt.Println("pulse .. wait a bit")
					time.Sleep(5 * time.Second)
					fmt.Println(".. continue detecting")
				}
				last = current
				time.Sleep(500 * time.Millisecond)
			}
		}()
	}

	return g
}

func (g *Gasmeter) Close() {
	g.pin.Detect(rpio.NoEdge)
	rpio.Close()
}

func (g *Gasmeter) CountPulse() {
	g.data100++

	im := &Impulse{
		Timestamp: time.Now(),
		ValueInM3: g.Reading(),
		Comment:   "impulse",
	}

	err := insertImpulseIntoDB(im)
	if err != nil {
		fmt.Printf("error counting pulse: %v", err)
	}
}

func (g *Gasmeter) Reading() float64 {
	return float64(g.data100) / float64(100)
}

func (g *Gasmeter) UpdateValue(value float64) {
	g.data100 = int(value * 100)
}
