package aggregator

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"

	"../clients"
	"../common"
)

const (
	interval   = 1
	configPath = "./src/config/aggregator.json"
)
const (
	Pre MarketState = iota
	Open
	After
	Closed
	Null
)

type MarketState int

type AggregatorConfig struct {
	Stocks []string `json:"stocks"`
}

type Aggregator struct {
	Client      client.Client
	Config      AggregatorConfig
	Timer       *time.Ticker
	MarketState MarketState
	FileWriter  *FileWriter
	DataStream  chan *common.Datum
}

func NewAggregator() Aggregator {
	dataStream := make(chan *common.Datum)

	a := Aggregator{
		Client:      client.NewYahooFinanceClient(dataStream),
		Config:      loadConfigFromFile(),
		Timer:       time.NewTicker(interval * time.Second),
		MarketState: Null,
		FileWriter:  NewFileWriter(),
		DataStream:  dataStream,
	}
	return a
}

func loadConfigFromFile() AggregatorConfig {
	raw, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalf("error: could not load client config file (%+v)\n", err)
	}
	var c AggregatorConfig
	json.Unmarshal(raw, &c)
	return c
}

func (T *Aggregator) Run() {
	log.Printf("Starting aggregator\n")
	go T.flushToFile()
	go T.Cron()
}

func (T *Aggregator) Cron() {
	for {
		t := <-T.Timer.C
		state := getMarketState()
		// if a state transition occurs
		if T.MarketState != state {
			T.MarketState = state
			switch state {
			case Pre:
				T.closeConnections()
			case Open:
				T.establishConnection()
			case After:
				T.closeConnections()
			case Closed:
				T.closeConnections()
			}
		}
		// every 10 seconds if the market is open...
		if t.Second()%10 == 0 && state == Open {
			// hot reloading config :D
			T.Config = loadConfigFromFile()
			T.establishConnection()
		}
	}
}

func getMarketState() MarketState {
	now := time.Now()
	if now.Weekday() != time.Saturday && now.Weekday() != time.Sunday {
		if now.Hour() >= 20 {
			return Closed
		} else if now.Hour() >= 16 {
			return After
		} else if (now.Hour() == 9 && now.Minute() >= 30) || now.Hour() > 9 {
			return Open
		} else if now.Hour() >= 4 {
			return Pre
		} else {
			return Closed
		}
	}
	return Closed
}

func (T *Aggregator) flushToFile() {
	for {
		select {
		case x := <-T.DataStream:
			T.FileWriter.WriteData(x)
		}
	}
}

func (T *Aggregator) establishConnection() {
	for _, sym := range T.Config.Stocks {
		err := T.Client.Connect(sym)
		if err != nil {
			log.Printf("error: connection terminated (%+v)\n", err)
		}
	}
}

func (T *Aggregator) closeConnections() {
	for _, sym := range T.Config.Stocks {
		T.Client.Disconnect(sym)
	}
}
