package trade_knife

import (
	"sort"
	"time"
)

// Quote is the group of candles and make time-series.
type Quote []*Candle

// Search for a candle and it's index amoung Quote by it's symbol and provided timestamp.
func (q *Quote) Find(symbol string, timestamp int64) (*Candle, int) {
	quote := *q
	for i, candle := range quote {
		if candle.Opentime.Unix() == timestamp || candle.Closetime.Unix() == timestamp || (candle.Opentime.Unix() < timestamp && candle.Closetime.Unix() > timestamp) {
			return candle, i
		}
	}

	return nil, 0
}

// Run through the quote and reorder candles by the open time.
func (q *Quote) Sort() {
	quote := *q
	sort.Slice(quote, func(i, j int) bool { return quote[i].Opentime.Before(quote[j].Opentime) })
	for i, candle := range quote {
		if i > 0 {
			candle.Previous = quote[i-1]
		}
		if i < len(quote)-1 {
			candle.Next = quote[i+1]
		}
	}
	*q = quote
}

// Search the quote for provided candle and update it if it exists, otherwise
// it will append to end of the quote.
//
// If you want to update a candle directly then pass sCandle
func (q *Quote) Sync(symbol string, interval Interval, open, high, low, close, volume float64, openTime, closeTime time.Time, sCandle ...*Candle) (candle *Candle, err CandleError) {
	var lc *Candle
	quote := *q
	checker := func(candle *Candle, openTime, closeTime time.Time) bool {
		return candle.Opentime.Equal(openTime) && candle.Closetime.Equal(closeTime)
	}
	// try last candle first
	if len(quote) > 0 {
		lc = quote[len(quote)-1]
		candle = lc
	}

	// if any suspecious candle provided try it then.
	if len(sCandle) > 0 {
		candle = sCandle[0]
	}

	if candle == nil || !checker(candle, openTime, closeTime) {
		candle, _ = quote.Find(symbol, openTime.Unix())
		if candle == nil {
			candle, err = NewCandle(symbol, open, high, low, close, volume, openTime, closeTime, interval, lc, nil)
			if err != nil {
				return
			}
			quote = append(quote, candle)
			if lc != nil {
				lc.Next = candle
			}
		}
	}
	candle.Open = open
	candle.High = high
	candle.Low = low
	candle.Close = close
	candle.Volume = volume

	*q = quote

	return
}

// Returns a list of indicators used in quote.
func (q *Quote) IndicatorNames() []string {
	indicators := []string{}
	for _, candle := range *q {
		for indicator := range candle.Indicators {
			hasIndicator := false
			for _, name := range indicators {
				if name == indicator {
					hasIndicator = true
					break
				}
			}
			if !hasIndicator {
				indicators = append(indicators, indicator)
			}
		}
	}

	return indicators
}
