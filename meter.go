package main

import (
	"fmt"
	"log"
	"time"

	"github.com/warthog618/gpio"
)

type PulseMeter interface {
	CountPulse()
	Reading() float64
	UpdateValue(float64)
	Close()
}

type Gasmeter struct {
	data100 int
	pin     *gpio.Pin
}

func NewGasmeter(fake bool) *Gasmeter {
	g := &Gasmeter{}

	ts, lastValueInDB, err := lastValueFromDB()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("found last value in db: %v %v\n", ts, lastValueInDB)
	g.UpdateValue(lastValueInDB)

	// somehow on startup the directly get a falling edge event
	// so we need to ignore the first one
	ignoreFirst := true

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
		err := gpio.Open()
		if err != nil {
			panic(err)
		}
		g.pin = gpio.NewPin(gpio.J8p13)
		g.pin.Input()
		g.pin.PullDown()

		err = g.pin.Watch(gpio.EdgeFalling, func(pin *gpio.Pin) {
			if !ignoreFirst {
				g.CountPulse()
			}
			ignoreFirst = false
			fmt.Printf("Current Gasmeter value is %v", g.Reading())
		})
		fmt.Println("waiting for pulse")

		if err != nil {
			panic(err)
		}
	}

	// FIXME: do we need to do this ??
	// defer pin.Unwatch()

	return g
}

func (g *Gasmeter) Close() {
	gpio.Close()
}

func (g *Gasmeter) CountPulse() {
	g.data100++

	err := insertImpulseIntoDB(g.Reading())
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
