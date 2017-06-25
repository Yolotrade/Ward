package client

import (
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"../common"
)

type GoogleFinanceClient struct {
	Config     ClientConfig
	DataStream chan *common.Datum
}

func NewGoogleFinanceClient(dataStream chan *common.Datum) Client {
	// should be read in from config files
	clientConfig := loadConfigFromFile("googlefinance")
	c := &GoogleFinanceClient{
		Config:     clientConfig,
		DataStream: dataStream,
	}
	return c
}

func (T *GoogleFinanceClient) ExecuteQuery(symbol string) (*common.Datum, error) {
	u := &url.URL{
		Scheme: T.Config.Scheme,
		Host:   T.Config.Host,
		Path:   T.Config.Path,
	}
	q := u.Query()
	q.Set("q", symbol)
	u.RawQuery = q.Encode()
	resp, err := http.Get(u.String())
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}
	datum, err := T.ExtractData(symbol, resp)
	if err != nil {
		return nil, err
	}
	return datum, nil
}

func (T *GoogleFinanceClient) ExtractData(symbol string, resp *http.Response) (*common.Datum, error) {
	timestamp, err := time.Parse(time.RFC1123, resp.Header.Get("Date"))
	if err != nil {
		log.Printf("error: could not parse http header date (%+v)\n", err)
		return nil, err
	}
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	price, err := strconv.ParseFloat(strings.TrimSpace(doc.Find(".pr").Text()), 64)
	if err != nil {
		log.Printf("error: could not parse price value (%+v)", err)
		return nil, err
	}
	datum := &common.Datum{
		Symbol: symbol,
		Price:  price,
		Time:   timestamp.Unix(),
	}
	return datum, nil
}
