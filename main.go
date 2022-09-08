package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	db, err := connectDB()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	ts, lastValueInDB, err := lastValueFromDB()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("found last value in db: %v %v\n", ts, lastValueInDB)

	fake := true
	g := NewGasmeter(fake)

	// initialize with last value from db
	g.UpdateValue(lastValueInDB)

	// add a handler to print the current value
	http.HandleFunc("/gasmeter", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Current Gasmeter value is %v", g.Reading())
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
		err = updateValueInDB(t.Value)
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
