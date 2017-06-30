package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"./aggregator"
)

const (
	port = ":8080"
)

func main() {
	SetupEnv()
	a := aggregator.NewAggregator()
	a.Run()
	// log.Printf("Starting server on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func debug(o interface{}) {
	fmt.Printf("Debug: %+v\n", o)
}

func SetupEnv() {
	// Make sure file structure exists
	dataPath := []string{
		"./data/",
		"./data/metadata",
		"./data/ticker",
	}
	for _, path := range dataPath {
		err := os.MkdirAll(path, 0777)
		if err != nil {
			log.Fatalf("error: could not create path %s (%+v)\n", path, err)
		}
	}
}
