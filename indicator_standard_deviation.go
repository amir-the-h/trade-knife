package trade_knife

import "github.com/markcheno/go-talib"

type StandardDeviation struct {
	Tag          IndicatorTag `mapstructure:"tag"`
	Source       Source       `mapstructure:"source"`
	InTimePeriod int          `mapstructure:"period"`
	Deviation    float64      `mapstructure:"deviation"`
}

func (sd *StandardDeviation) Add(q *Quote, c *Candle) (ok bool) {
	if c != nil {
		candle, i := q.Find(c.Opentime.Unix())
		if candle == nil {
			goto out
		}

		startIndex := i - sd.InTimePeriod
		if startIndex < 0 {
			return
		}

		quote := Quote{
			Market:   q.Market,
			Symbol:   q.Symbol,
			Interval: q.Interval,
			Candles:  q.Candles[startIndex : i+1],
		}

		values := talib.StdDev(quote.Get(sd.Source), sd.InTimePeriod, sd.Deviation)
		candle.AddIndicator(sd.Tag, values[len(values)-1])
		q.Candles[i] = candle

		return
	}

out:
	if len(q.Candles) < sd.InTimePeriod {
		return
	}

	values := talib.StdDev(q.Get(sd.Source), sd.InTimePeriod, sd.Deviation)
	q.AddIndicator(sd.Tag, values)

	return
}

func (sd *StandardDeviation) Is(tag IndicatorTag) bool {
	return sd.Tag == tag
}
