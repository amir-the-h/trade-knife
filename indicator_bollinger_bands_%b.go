package trade_knife

import (
	"fmt"

	"github.com/markcheno/go-talib"
)

type BollingerBandsB struct {
	Tag IndicatorTag      `mapstructure:"tag"`
	Std StandardDeviation `mapstructure:"standardDeviation"`
}

func (bbb *BollingerBandsB) Add(q *Quote, c *Candle) bool {
	if c != nil {
		candle, i := q.Find(c.Opentime.Unix())
		if candle == nil {
			return false
		}

		startIndex := i - bbb.Std.InTimePeriod
		if startIndex < 0 {
			return false
		}

		quote := Quote{
			Market:   q.Market,
			Currency:   q.Currency,
			Interval: q.Interval,
			Candles:  q.Candles[startIndex : i+1],
		}

		deviation, ok := c.Get(Source(bbb.Std.Tag))
		if !ok {
			if !bbb.Std.Add(&quote, c) {
				return false
			}

			deviation, ok = c.Get(Source(bbb.Std.Tag))
			if !ok {
				return false
			}
		}

		sma := &Ma{
			Tag:          IndicatorTag(fmt.Sprintf("bbb:sma:%s:%d", bbb.Std.Source, bbb.Std.InTimePeriod)),
			Source:       bbb.Std.Source,
			Type:         talib.SMA,
			InTimePeriod: bbb.Std.InTimePeriod,
		}
		basis, ok := c.Get(Source(sma.Tag))
		if !ok {
			if !sma.Add(&quote, c) {
				return false
			}
			basis, ok = c.Get(Source(sma.Tag))
			if !ok {
				return false
			}
		}

		upper := basis + deviation
		lower := basis - deviation
		bbr, ok := c.Get(bbb.Std.Source)
		if !ok {
			return false
		}
		bbr = (bbr - lower) / (upper - lower)
		candle.AddIndicator(bbb.Tag, bbr)
		q.Candles[i] = candle

		return true
	}

	if len(q.Candles) < bbb.Std.InTimePeriod {
		return false
	}

	for _, candle := range q.Candles {
		if !bbb.Add(q, candle) {
			return false
		}
	}

	return true
}

func (bbb *BollingerBandsB) Is(tag IndicatorTag) bool {
	return bbb.Tag == tag
}
