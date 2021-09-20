package trade_knife

import (
	"errors"
	"github.com/amir-the-h/goex"
	"sort"
	"time"
)

// Quote is the group of candles and make time-series.
type Quote struct {
	Currency goex.CurrencyPair
	Market   MarketType
	Interval Interval
	Candles  []*Candle
}

// Find searches for a candle, and its index among Quote by its symbol and provided timestamp.
func (q *Quote) Find(timestamp int64) (*Candle, int) {
	quote := *q
	for i, candle := range quote.Candles {
		if candle.Opentime.Unix() == timestamp || (candle.Opentime.Unix() < timestamp && candle.Closetime.Unix() > timestamp) {
			return candle, i
		}
	}

	return nil, 0
}

// Sort runs through the quote and reorder candles by the open time.
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

// Sync searches the quote for provided candle and update it if it exists, otherwise
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

// IndicatorTags returns a list of indicators used in quote.
func (q *Quote) IndicatorTags() []IndicatorTag {
	tags := []IndicatorTag{}
	quote := *q
	for _, candle := range quote.Candles {
		for indicator := range candle.Indicators {
			hasIndicator := false
			for _, tag := range tags {
				if tag == indicator {
					hasIndicator = true
					break
				}
			}
			if !hasIndicator {
				tags = append(tags, indicator)
			}
		}
	}

	return tags
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

// AddIndicator adds indicator values by the given tag into the quote.
func (q *Quote) AddIndicator(tag IndicatorTag, values []float64) error {
	quote := *q
	if len(values) != len(quote.Candles) {
		return errors.New("count mismatched")
	}

	for i := range values {
		q.Candles[i].Indicators[tag] = values[i]
	}

	return nil
}
