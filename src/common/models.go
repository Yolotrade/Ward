package common

import (
	"strconv"
	"strings"
)

type Datum struct {
	Symbol           string
	Time             int64
	CurrentPrice     *float64
	Ask              *float64
	Bid              *float64
	AskSize          *float64
	BidSize          *float64
	DayLow           *float64
	DayHigh          *float64
	MarketCap        *float64
	Volume           *float64
	PercentageChange *float64
}

var Query = "&k=" + strings.Join([]string{
	"l10",
	"a00",
	"b00",
	"a50",
	"b60",
	"g00",
	"h00",
	"j10",
	"v00",
	"p43",
}, ",")

func (T *Datum) String() string {
	toString := func(f *float64) string {
		if f != nil {
			return strconv.FormatFloat(*f, 'f', 2, 64)
		}
		return "null"
	}
	return "{" +
		"Symbol: " + T.Symbol + ", " +
		"Time: " + strconv.FormatInt(T.Time, 10) + ", " +
		"CurrentPrice: " + toString(T.CurrentPrice) + ", " +
		"Ask: " + toString(T.Ask) + ", " +
		"Bid: " + toString(T.Bid) + ", " +
		"AskSize: " + toString(T.AskSize) + ", " +
		"BidSize: " + toString(T.BidSize) + ", " +
		"DayLow: " + toString(T.DayLow) + ", " +
		"DayHigh: " + toString(T.DayHigh) + ", " +
		"MarketCap: " + toString(T.MarketCap) + ", " +
		"Volume: " + toString(T.Volume) + ", " +
		"PercentageChange: " + toString(T.PercentageChange) + "}"
}
