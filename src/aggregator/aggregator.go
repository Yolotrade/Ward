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
	Clients    []client.Client
	Config     AggregatorConfig
	FileWriter *FileWriter
	Timer      *time.Ticker

	DataStream chan *common.Datum
}

func NewAggregator() Aggregator {
	config := loadConfigFromFile()
	dataStream := make(chan *common.Datum)
	clients := []client.Client{}
	for _, clientName := range config.ClientNames {
		clients = append(clients, client.ClientConstructor[clientName](dataStream))
	}
	fileWriter := NewFileWriter()

	a := Aggregator{
		Clients:    clients,
		Config:     config,
		FileWriter: fileWriter,
		Timer:      time.NewTicker(interval),
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
	// T.pullData()
	T.Clients[0].ExecuteQuery("NVDA")
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
	for {
		select {
		case <-T.Timer.C:
			for _, sym := range T.Config.Stocks {
				for _, client := range T.Clients {
					go func() {
						datum, err := client.ExecuteQuery(sym)
						if err == nil {
							T.DataStream <- datum
						}
					}()
				}
			}
		}
	}
}
