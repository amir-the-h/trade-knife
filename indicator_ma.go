package trade_knife

import "github.com/markcheno/go-talib"

type Ma struct {
	Tag          IndicatorTag `mapstructure:"tag"`
	Source       Source       `mapstructure:"source"`
	Type         talib.MaType `mapstructure:"type"`
	InTimePeriod int          `mapstructure:"period"`
}

func (ma *Ma) Add(q *Quote, c *Candle) bool {
	if c != nil {
		candle, i := q.Find(c.Opentime.Unix())
		if candle == nil {
			goto out
		}

		startIndex := i - ma.InTimePeriod
		if startIndex < 0 {
			return false
		}

		quote := Quote{
			Market:   q.Market,
			Symbol:   q.Symbol,
			Interval: q.Interval,
			Candles:  q.Candles[startIndex : i+1],
		}

		values := talib.Ma(quote.Get(ma.Source), ma.InTimePeriod, ma.Type)
		c.AddIndicator(ma.Tag, values[len(values)-1])
		q.Candles[i] = c

		return true
	}

out:
	if len(q.Candles) < ma.InTimePeriod {
		return false
	}

	values := talib.Ma(q.Get(ma.Source), ma.InTimePeriod, ma.Type)
	q.AddIndicator(ma.Tag, values)

	return true
}

func (ma *Ma) Is(tag IndicatorTag) bool {
	return ma.Tag == tag
}
