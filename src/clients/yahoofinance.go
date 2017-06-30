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
	Config     ClientConfig
	DataStream chan *common.Datum
	Client     *http.Client
}

func NewYahooFinanceClient(dataStream chan *common.Datum) Client {
	clientConfig := loadConfigFromFile("yahoofinance")
	c := &YahooFinanceClient{
		Config:     clientConfig,
		DataStream: dataStream,
		Client:     &http.Client{},
	}
	return c
}

func (T *YahooFinanceClient) Connect(symbol string) error {
	req, err := T.createQuery(symbol)
	if err != nil {
		panic(err.Error())
	}
	resp, err := T.Client.Do(req)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	return T.ExtractData(symbol, resp)
}

type YahooFinanceResponse struct {
	Price string `json:"l10"`
}

func (T *YahooFinanceClient) ExtractData(symbol string, resp *http.Response) error {
	reader := bufio.NewReader(resp.Body)
	regex := regexp.MustCompile(`yfs_.*?\((?P<json>.*?)\)`)
	for {
		tok, err := reader.ReadBytes('>')
		if err != nil {
			log.Printf("error: %+v", err)
			break
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
	return nil
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
