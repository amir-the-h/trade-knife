package trade_knife

import "github.com/markcheno/go-talib"

type LinearRegression struct {
	Tag          IndicatorTag `mapstructure:"tag"`
	Source       Source       `mapstructure:"source"`
	InTimePeriod int          `mapstructure:"period"`
}

func (lr *LinearRegression) Add(q *Quote, c *Candle) bool {
	if c != nil {
		candle, i := q.Find(c.Opentime.Unix())
		if candle == nil {
			return false
		}

		startIndex := i - lr.InTimePeriod
		if startIndex < 0 {
			return false
		}

		quote := Quote{
			Market:   q.Market,
			Currency:   q.Currency,
			Interval: q.Interval,
			Candles:  q.Candles[startIndex : i+1],
		}

		values := talib.LinearReg(quote.Get(lr.Source), lr.InTimePeriod)
		c.AddIndicator(lr.Tag, values[len(values)-1])
		q.Candles[i] = c

		return true
	}

	if len(q.Candles) < lr.InTimePeriod {
		return false
	}

	values := talib.LinearReg(q.Get(lr.Source), lr.InTimePeriod)
	q.AddIndicator(lr.Tag, values)

	return true
}

func (lr *LinearRegression) Is(tag IndicatorTag) bool {
	return lr.Tag == tag
}
