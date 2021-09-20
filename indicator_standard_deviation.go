package trade_knife

import "github.com/markcheno/go-talib"

type StandardDeviation struct {
	Tag          IndicatorTag `mapstructure:"tag"`
	Source       Source       `mapstructure:"source"`
	InTimePeriod int          `mapstructure:"period"`
	Deviation    float64      `mapstructure:"deviation"`
}

func (sd *StandardDeviation) Add(q *Quote, c *Candle) bool {
	if c != nil {
		candle, i := q.Find(c.Opentime.Unix())
		if candle == nil {
			return false
		}

		startIndex := i - sd.InTimePeriod
		if startIndex < 0 {
			return false
		}

		quote := Quote{
			Market:   q.Market,
			Currency:   q.Currency,
			Interval: q.Interval,
			Candles:  q.Candles[startIndex : i+1],
		}

		values := talib.StdDev(quote.Get(sd.Source), sd.InTimePeriod, sd.Deviation)
		c.AddIndicator(sd.Tag, values[len(values)-1])
		q.Candles[i] = c

		return true
	}

	if len(q.Candles) < sd.InTimePeriod {
		return false
	}

	values := talib.StdDev(q.Get(sd.Source), sd.InTimePeriod, sd.Deviation)
	q.AddIndicator(sd.Tag, values)

	return true
}

func (sd *StandardDeviation) Is(tag IndicatorTag) bool {
	return sd.Tag == tag
}
