package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/warthog618/gpio"
)

type gasmeter struct {
	data100 int
}

func (g *gasmeter) CountPulse() {
	g.data100++
}

func (g *gasmeter) Reading() float64 {
	return float64(g.data100) / float64(100)
}

func main() {
	err := gpio.Open()
	if err != nil {
		panic(err)
	}
	defer gpio.Close()
	pin := gpio.NewPin(gpio.J8p13)
	pin.Input()
	pin.PullDown()

	gmeter := &gasmeter{}

	// capture exit signals to ensure resources are released on exit.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	defer signal.Stop(quit)

	err = pin.Watch(gpio.EdgeFalling, func(pin *gpio.Pin) {
		gmeter.CountPulse()
		fmt.Printf("Current Gasmeter value is %v", gmeter.Reading())
	})
	if err != nil {
		panic(err)
	}
	defer pin.Unwatch()

	// In a real application the main thread would do something useful here.
	// But we'll just run for a minute then exit.
	fmt.Println("Watching Pin 4...")
	select {
	case <-time.After(time.Minute):
	case <-quit:
	}
}
