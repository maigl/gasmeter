package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Impulse struct {
	Timestamp time.Time `json:"timestamp gorm:"primaryKey"`
	ValueInM3 float64   `json:"value_in_m3"`
	Comment   string    `json:"comment"`
}

func main() {
	close, err := connectDB()
	if err != nil {
		log.Fatal(err)
	}

	defer close()

	lastImpulse, err := lastValueFromDB()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("found last value in db: %v\n", lastImpulse)

	fake := false
	g := NewGasmeter(fake)

	// initialize with last value from db
	g.UpdateValue(lastImpulse.ValueInM3)

	// add a handler to print the current value
	http.HandleFunc("/gasmeter", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "The Current Gasmeter value is %v", g.Reading())
	})

	// add a handler to update the current value
	http.HandleFunc("/gasmeter/update", func(w http.ResponseWriter, r *http.Request) {
		// read json from request body
		decoder := json.NewDecoder(r.Body)
		var t struct {
			Value float64 `json:"value"`
		}
		err := decoder.Decode(&t)
		if err != nil {
			fmt.Fprintf(w, "error: %v", err)
			return
		}
		r.Body.Close()
		// update the value in the database
		err = insertImpulseIntoDB(&Impulse{Timestamp: time.Now(), ValueInM3: g.Reading(), Comment: "manual update"})
		if err != nil {
			fmt.Fprintf(w, "error: %v", err)
			return
		}
		// update the value in the meter
		g.UpdateValue(t.Value)
		fmt.Printf("updated value to %v\n", t.Value)
	})

	go func() {
		// listen and serve
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal(err)
		}
	}()

	// capture exit signals to ensure resources are released on exit.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	defer signal.Stop(quit)
	// wait for exit signal
	<-quit
}
