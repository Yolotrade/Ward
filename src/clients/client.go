package client

import (
	"net/http"

	"../common"
)

type Client interface {
	ExecuteQuery(symbol string) error
	ExtractData(resp *http.Response) (*common.Datum, error)
}
