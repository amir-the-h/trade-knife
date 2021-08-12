package trade_knife

import (
	"errors"
	"sort"
	"time"
)

// Quote is the group of candles and make time-series.
type Quote struct {
	Market   MarketType
	Symbol   string
	Interval Interval
	Candles  []*Candle
}

// Search for a candle and it's index amoung Quote by it's symbol and provided timestamp.
func (q *Quote) Find(timestamp int64) (*Candle, int) {
	quote := *q
	for i, candle := range quote.Candles {
		if candle.Opentime.Unix() == timestamp || (candle.Opentime.Unix() < timestamp && candle.Closetime.Unix() > timestamp) {
			return candle, i
		}
	}

	return nil, 0
}

// Run through the quote and reorder candles by the open time.
func (q *Quote) Sort() {
	quote := *q
	sort.Slice(quote.Candles, func(i, j int) bool { return quote.Candles[i].Opentime.Before(quote.Candles[j].Opentime) })
	for i, candle := range quote.Candles {
		if i > 0 {
			candle.Previous = quote.Candles[i-1]
		}
		if i < len(quote.Candles)-1 {
			candle.Next = quote.Candles[i+1]
		}
	}
	*q = quote
}

// Search the quote for provided candle and update it if it exists, otherwise
// it will append to end of the quote.
//
// If you want to update a candle directly then pass sCandle
func (q *Quote) Sync(open, high, low, close, volume float64, openTime, closeTime time.Time, sCandle ...*Candle) (candle *Candle, err CandleError) {
	var lc *Candle
	quote := *q
	checker := func(candle *Candle, openTime, closeTime time.Time) bool {
		return candle.Opentime.Equal(openTime) && candle.Closetime.Equal(closeTime)
	}

	// try last candle first
	if len(quote.Candles) > 0 {
		lc = quote.Candles[len(quote.Candles)-1]
		candle = lc
	}

	// if any suspecious candle provided try it then.
	if len(sCandle) > 0 {
		candle = sCandle[0]
	}

	if candle == nil || !checker(candle, openTime, closeTime) {
		candle, _ = quote.Find(openTime.Unix())
		if candle == nil {
			candle, err = NewCandle(open, high, low, close, volume, openTime, closeTime, lc, nil)
			if err != nil {
				return
			}
			quote.Candles = append(quote.Candles, candle)
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
	quote := *q
	for _, candle := range quote.Candles {
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

// Merge target quote into the current quote, rewrite duplicates and sort it.
func (q *Quote) Merge(target *Quote) {
	quote := *q
	t := *target
	for _, candle := range t.Candles {
		c, i := quote.Find(candle.Opentime.Unix())
		if c != nil {
			quote.Candles[i] = candle
		} else {
			quote.Candles = append(quote.Candles, candle)
		}
	}

	*q = quote
}

// Add indicator values by the given name into the quote.
func (q *Quote) AddIndicator(name string, values []float64) error {
	quote := *q
	if len(values) != len(quote.Candles) {
		return errors.New("count mismatched")
	}

	for i, candle := range quote.Candles {
		candle.Indicators[name] = values[i]
	}

	return nil
}
