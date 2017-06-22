package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"./aggregator"
)

const DEBUG = 1

func main() {
	dataPath := "./data/"
	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		err := os.Mkdir(dataPath, 0777)
		if err != nil {
			log.Fatalf("error: could not create data directory (%+v)\n", err)
		}
	}

	log.Printf("Starting aggregator\n")
	a := aggregator.NewAggregator()
	a.Run()

	port := ":8080"
	log.Printf("Starting server on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func debug(o interface{}) {
	fmt.Printf("Debug: %+v\n", o)
}
