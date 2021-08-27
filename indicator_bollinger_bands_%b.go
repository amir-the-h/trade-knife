package trade_knife

import (
	"fmt"

	"github.com/markcheno/go-talib"
)

type BollingerBandsB struct {
	Tag IndicatorTag      `mapstructure:"tag"`
	Std StandardDeviation `mapstructure:"standardDeviation"`
}

func (bbb *BollingerBandsB) Add(q *Quote, c *Candle) (ok bool) {
	if c != nil {
		candle, i := q.Find(c.Opentime.Unix())
		if candle == nil {
			goto out
		}

		startIndex := i - bbb.Std.InTimePeriod
		if startIndex < 0 {
			return
		}

		quote := Quote{
			Market:   q.Market,
			Symbol:   q.Symbol,
			Interval: q.Interval,
			Candles:  q.Candles[startIndex : i+1],
		}

		deviation := c.Get(Source(bbb.Std.Tag))

		if deviation == 0 {
			if !bbb.Std.Add(&quote, c) {
				return
			}
			deviation = c.Get(Source(bbb.Std.Tag))
		}

		sma := &Ma{
			Tag:          IndicatorTag(fmt.Sprintf("bbb:sma:%s:%d", bbb.Std.Source, bbb.Std.InTimePeriod)),
			Source:       bbb.Std.Source,
			Type:         talib.SMA,
			InTimePeriod: bbb.Std.InTimePeriod,
		}
		basis := c.Get(Source(sma.Tag))
		if basis == 0 {
			if !sma.Add(&quote, c) {
				return
			}
			basis = c.Get(Source(sma.Tag))
		}

		upper := basis + deviation
		lower := basis - deviation
		bbr := (c.Get(bbb.Std.Source) - lower) / (upper - lower)
		candle.AddIndicator(bbb.Tag, bbr)
		q.Candles[i] = candle
		ok = true

		return
	}

out:
	if len(q.Candles) < bbb.Std.InTimePeriod {
		return
	}

	for _, candle := range q.Candles {
		if !bbb.Add(q, candle) {
			return
		}
	}
	ok = true

	return
}

func (bbb *BollingerBandsB) Is(tag IndicatorTag) bool {
	return bbb.Tag == tag
}
