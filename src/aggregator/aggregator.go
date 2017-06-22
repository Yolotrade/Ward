package aggregator

import (
	"time"

	"../clients"
	"../common"
)

const (
	interval = time.Second * 10
)

type Aggregator struct {
	Clients    []client.Client
	Symbols    []string
	FileWriter *FileWriter
	Timer      *time.Ticker

	DataStream chan *common.Datum
}

func NewAggregator() Aggregator {
	dataStream := make(chan *common.Datum)
	// maybe load from config files and use reflection?
	clients := []client.Client{
		client.NewGoogleFinanceClient(dataStream),
	}
	fileWriter := NewFileWriter()
	symbols := []string{
		"NVDA",
	}

	a := Aggregator{
		Clients:    clients,
		Symbols:    symbols,
		FileWriter: fileWriter,
		Timer:      time.NewTicker(interval),
		DataStream: dataStream,
	}
	return a
}

func (T *Aggregator) Run() {
	go func() {
		for {
			select {
			case <-T.Timer.C:
				T.execute()
			}
		}
	}()
	go func() {
		for {
			select {
			case x := <-T.DataStream:
				T.FileWriter.WriteData(x)
			}
		}
	}()
}

func (T *Aggregator) execute() {
	for _, sym := range T.Symbols {
		for _, client := range T.Clients {
			go client.ExecuteQuery(sym)
		}
	}
}
