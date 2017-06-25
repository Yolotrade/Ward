package client

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"

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

	u.RawQuery = q.Encode() + "&k=a00,a50,b00,b60,c10,g00,h00,j10,l10,p20,t10,v00,z08,z09"
	fmt.Printf("%s\n", u.String())

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

func (T *YahooFinanceClient) ExecuteQuery(symbol string) (*common.Datum, error) {
	req, err := T.createQuery(symbol)
	if err != nil {
		panic(err.Error())
	}
	resp, err := T.Client.Do(req)
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

type YahooFinanceResponse struct {
	Price string `json:"l10"`
}

func (T *YahooFinanceClient) ExtractData(symbol string, resp *http.Response) (*common.Datum, error) {
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
			fmt.Printf("%s\n", match[1])
		}
	}

	log.Fatalf("Done")
	// datum := &common.Datum{
	// 	Symbol: symbol,
	// 	Price:  price,
	// 	Time:   timestamp.Unix(),
	// }
	return nil, nil
}
