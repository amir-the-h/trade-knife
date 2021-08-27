package trade_knife

import "github.com/markcheno/go-talib"

type LinearRegression struct {
	Tag          IndicatorTag `mapstructure:"tag"`
	Source       Source       `mapstructure:"source"`
	InTimePeriod int          `mapstructure:"period"`
}

func (lr *LinearRegression) Add(q *Quote, c *Candle) (ok bool) {
	if c != nil {
		candle, i := q.Find(c.Opentime.Unix())
		if candle == nil {
			goto out
		}

		startIndex := i - lr.InTimePeriod
		if startIndex < 0 {
			return
		}

		quote := Quote{
			Market:   q.Market,
			Symbol:   q.Symbol,
			Interval: q.Interval,
			Candles:  q.Candles[startIndex : i+1],
		}

		values := talib.LinearReg(quote.Get(lr.Source), lr.InTimePeriod)
		candle.AddIndicator(lr.Tag, values[len(values)-1])
		q.Candles[i] = candle

		return
	}

out:
	if len(q.Candles) < lr.InTimePeriod {
		return
	}

	values := talib.LinearReg(q.Get(lr.Source), lr.InTimePeriod)
	q.AddIndicator(lr.Tag, values)

	return
}

func (lr *LinearRegression) Is(tag IndicatorTag) bool {
	return lr.Tag == tag
}
