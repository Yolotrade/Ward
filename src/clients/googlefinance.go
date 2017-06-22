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
	Scheme     string
	Host       string
	Path       string
	DataStream chan *common.Datum
}

func NewGoogleFinanceClient(dataStream chan *common.Datum) *GoogleFinanceClient {
	// should be read in from config files
	scheme := "https"
	host := "www.google.com"
	path := "finance"
	c := &GoogleFinanceClient{
		Scheme:     scheme,
		Host:       host,
		Path:       path,
		DataStream: dataStream,
	}
	return c
}

func (T *GoogleFinanceClient) ExecuteQuery(symbol string) error {
	u := &url.URL{
		Scheme: T.Scheme,
		Host:   T.Host,
		Path:   T.Path,
	}
	q := u.Query()
	q.Set("q", symbol)
	u.RawQuery = q.Encode()
	resp, err := http.Get(u.String())
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	datum, err := T.ExtractData(resp)
	if err != nil {
		return err
	}
	T.DataStream <- datum
	return nil
}

func (T *GoogleFinanceClient) ExtractData(resp *http.Response) (*common.Datum, error) {
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
		Symbol: "NVDA",
		Price:  price,
		Time:   timestamp.Unix(),
	}
	return datum, nil
}
