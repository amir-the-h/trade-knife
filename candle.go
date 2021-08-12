package trade_knife

import (
	"time"
)

// Candle is the main structure which contains a group of useful candlestick data.
type Candle struct {
	Symbol     string             `json:"symbol"`
	Open       float64            `json:"open"`
	High       float64            `json:"high"`
	Low        float64            `json:"low"`
	Close      float64            `json:"close"`
	Volume     float64            `json:"volume"`
	Score      float64            `json:"score"`
	Indicators map[string]float64 `json:"indicators"`
	Interval   Interval           `json:"interval"`
	Opentime   time.Time          `json:"open_time"`
	Closetime  time.Time          `json:"close_time"`
	Next       *Candle            `json:"-"`
	Previous   *Candle            `json:"-"`
}

// Returns a pointer to a fresh candle with provided data.
func NewCandle(symbol string, open, high, low, close, volume float64, openTime, closeTime time.Time, interval Interval, previous, next *Candle) (candle *Candle, err CandleError) {
	if high < low || high < open || high < close || low > open || low > close {
		err = ErrInvalidCandleData
		return
	}

	candle = &Candle{
		Symbol:     symbol,
		Open:       open,
		High:       high,
		Low:        low,
		Close:      close,
		Volume:     volume,
		Opentime:   openTime,
		Closetime:  closeTime,
		Interval:   interval,
		Indicators: make(map[string]float64),
		Previous:   previous,
		Next:       next,
	}

	return
}

// Add indicator value by the given name into the candle.
func (c *Candle) AddIndicator(name string, value float64) {
	c.Indicators[name] = value
}