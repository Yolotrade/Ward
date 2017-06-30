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
	interval   = time.Second * 1
	configPath = "./src/config/aggregator.json"
)

type AggregatorConfig struct {
	Stocks      []string `json:"stocks"`
	ClientNames []string `json:"clients"`
}

type Aggregator struct {
	Client     client.Client
	Config     AggregatorConfig
	FileWriter *FileWriter

	DataStream chan *common.Datum
}

func NewAggregator() Aggregator {
	dataStream := make(chan *common.Datum)

	a := Aggregator{
		Client:     client.NewYahooFinanceClient(dataStream),
		Config:     loadConfigFromFile(),
		FileWriter: NewFileWriter(),
		DataStream: dataStream,
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
	T.pullData()
}

func (T *Aggregator) flushToFile() {
	for {
		select {
		case x := <-T.DataStream:
			T.FileWriter.WriteData(x)
		}
	}
}

func (T *Aggregator) pullData() {
	for _, sym := range T.Config.Stocks {
		go func(s string) {
			err := T.Client.Connect(s)
			if err != nil {
				log.Printf("error: connection terminated (%+v)\n", err)
			}
		}(sym)
	}
}
