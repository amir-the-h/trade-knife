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

func (af *AutoFibo) Add(q *Quote, c *Candle) (ok bool) {
	if c != nil {
		candle, i := q.Find(c.Opentime.Unix())
		if candle == nil {
			goto out
		}

		quote := Quote{
			Market:   q.Market,
			Symbol:   q.Symbol,
			Interval: q.Interval,
			Candles:  q.Candles[:i+1],
		}

		fibos := AutoFiboRectracement(quote.Get(SourceHigh), quote.Get(SourceLow), quote.Get(SourceClose), af.Ratios, af.Depth, af.Deviation)
		for ratio, fibo := range fibos[len(fibos)-1] {
			candle.AddIndicator(IndicatorTag(fmt.Sprintf("%s:%f", af.Tag, ratio)), fibo)
		}
		q.Candles[i] = candle

		return
	}

out:
	for _, candle := range q.Candles {
		af.Add(q, candle)
	}

	return
}
