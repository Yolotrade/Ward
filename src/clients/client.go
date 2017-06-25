package client

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"../common"
)

const configPath = "./src/config/clients/"

var ClientConstructor = map[string]func(chan *common.Datum) Client{
	"googlefinance": NewGoogleFinanceClient,
	"yahoofinance":  NewYahooFinanceClient,
}

type Client interface {
	ExecuteQuery(symbol string) (*common.Datum, error)
	ExtractData(symbol string, resp *http.Response) (*common.Datum, error)
}

type ClientConfig struct {
	Name   string            `json:"name"`
	Scheme string            `json:"scheme"`
	Host   string            `json:"host"`
	Path   string            `json:"path"`
	Header map[string]string `json:"header"`
}

func loadConfigFromFile(clientname string) ClientConfig {
	raw, err := ioutil.ReadFile(configPath + clientname + ".json")
	if err != nil {
		log.Fatalf("error: could not load client config file (%+v)\n", err)
	}
	var c ClientConfig
	json.Unmarshal(raw, &c)
	return c
}
