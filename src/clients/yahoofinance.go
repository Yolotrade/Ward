package client

import (
	"bufio"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"../common"
)

type YahooFinanceClient struct {
	Config      ClientConfig
	DataStream  chan *common.Datum
	Client      *http.Client
	Connections map[string]*http.Response
}

func NewYahooFinanceClient(dataStream chan *common.Datum) Client {
	clientConfig := loadConfigFromFile("yahoofinance")
	c := &YahooFinanceClient{
		Config:      clientConfig,
		DataStream:  dataStream,
		Client:      &http.Client{},
		Connections: make(map[string]*http.Response),
	}
	return c
}

func (T *YahooFinanceClient) Connect(symbol string) error {
	// check if connection already exists
	if _, ok := T.Connections[symbol]; ok {
		return nil
	}

	// construct connection url
	req, err := T.createQuery(symbol)
	if err != nil {
		panic(err.Error())
	}
	resp, err := T.Client.Do(req)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	// store connection handle
	T.Connections[symbol] = resp

	go T.ExtractData(symbol, resp)

	return nil
}

func (T *YahooFinanceClient) Disconnect(symbol string) {
	if resp, ok := T.Connections[symbol]; ok {
		resp.Body.Close()
		delete(T.Connections, symbol)
	}
}

func (T *YahooFinanceClient) ExtractData(symbol string, resp *http.Response) error {
	reader := bufio.NewReader(resp.Body)
	regex := regexp.MustCompile(`yfs_.*?\((?P<json>.*?)\)`)
	for {
		tok, err := reader.ReadBytes('>')
		if err != nil {
			return err
		}
		match := regex.FindSubmatch(tok)
		if len(match) > 1 {
			datum, err := deserialize(symbol, match[1])
			if err != nil {
				log.Printf("error: could not deserialize data payload (%+v)\n", err)
				continue
			}
			T.DataStream <- datum
		}
	}
}

func (T *YahooFinanceClient) createQuery(symbol string) (*http.Request, error) {
	// create request url string
	u := &url.URL{
		Scheme: T.Config.Scheme,
		Host:   T.Config.Host,
		Path:   T.Config.Path,
	}
	q := u.Query()
	q.Add("s", symbol)
	q.Add("marketid", "us_market")
	q.Add("callback", "parent.yfs_u1f")
	q.Add("mktmcb", "parent.yfs_mktmcb")
	q.Add("gencallback", "parent.yfs_gencb")

	u.RawQuery = q.Encode() + common.Query

	// create request
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		log.Printf("error: could not create get request to %s\n", u.String())
		return nil, err
	}

	// set request headers
	for key, value := range T.Config.Header {
		req.Header.Add(key, value)
	}
	return req, nil
}

func deserialize(symbol string, json []byte) (*common.Datum, error) {
	extract := func(prop string) *float64 {
		regex := regexp.MustCompile(prop + `:"(.*?)"`)
		match := regex.FindSubmatch(json)
		if len(match) > 1 {
			rawVal := strings.Replace(string(match[1]), ",", "", -1)
			val, err := strconv.ParseFloat(rawVal, 64)
			if err != nil {
				return nil
			}
			return &val
		}
		return nil
	}
	d := &common.Datum{
		Symbol:           symbol,
		Time:             time.Now().Unix(),
		CurrentPrice:     extract("l10"),
		Ask:              extract("a00"),
		Bid:              extract("b00"),
		AskSize:          extract("a50"),
		BidSize:          extract("b60"),
		DayLow:           extract("g00"),
		DayHigh:          extract("h00"),
		MarketCap:        extract("j10"),
		Volume:           extract("v00"),
		PercentageChange: extract("p43"),
	}
	return d, nil
}
