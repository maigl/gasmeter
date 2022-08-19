package main

import (
	"fmt"
	"log"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

func main() {
	fmt.Println("hallo")

	err := rpio.Open()
	if err != nil {
		log.Fatal(err)
	}

	led := rpio.Pin(17)
	led.Output()

	gPin := rpio.Pin(27)
	gPin.Input()
	gPin.PullDown()

	for {
		led.Toggle()
		time.Sleep(100 * time.Millisecond)

		res := gPin.Read()
		fmt.Printf("gas: %v\n", res)
	}

	rpio.Close()
}
