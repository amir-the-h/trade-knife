package trade_knife

import (
	"fmt"
)

type AutoFibo struct {
	Tag       IndicatorTag `mapstructure:"tag"`
	Ratios    []float64    `mapstructure:"ratios"`
	Deviation float64      `mapstructure:"deviation"`
	Depth     int          `mapstructure:"depth"`
}

func (af *AutoFibo) Add(q *Quote, c *Candle) bool {
	if c != nil {
		candle, i := q.Find(c.Opentime.Unix())
		if candle == nil {
			return false
		}

		quote := Quote{
			Market:   q.Market,
			Currency: q.Currency,
			Interval: q.Interval,
			Candles:  q.Candles[:i+1],
		}

		fibos := AutoFiboRectracement(quote.Get(SourceHigh), quote.Get(SourceLow), quote.Get(SourceClose), af.Ratios, af.Depth, af.Deviation)
		for ratio, fibo := range fibos[len(fibos)-1] {
			c.AddIndicator(IndicatorTag(fmt.Sprintf("%s:%.2f", af.Tag, ratio)), fibo)
		}
		q.Candles[i] = c

		return true
	}

	for _, candle := range q.Candles {
		if !af.Add(q, candle) {
			return false
		}
	}

	return true
}

func (af *AutoFibo) Is(tag IndicatorTag) bool {
	return af.Tag == tag
}
