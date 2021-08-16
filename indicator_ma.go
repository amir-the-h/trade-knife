package trade_knife

import "github.com/markcheno/go-talib"

type Ma struct {
	Tag          IndicatorTag `mapstructure:"tag"`
	Source       Source       `mapstructure:"source"`
	Type         talib.MaType `mapstructure:"type"`
	InTimePeriod int          `mapstructure:"period"`
}

func (ma *Ma) Add(q *Quote, c *Candle) (ok bool) {
	if c != nil {
		candle, i := q.Find(c.Opentime.Unix())
		if candle == nil {
			goto out
		}

		startIndex := i - ma.InTimePeriod
		if startIndex < 0 {
			return
		}

		quote := Quote{
			Market:   q.Market,
			Symbol:   q.Symbol,
			Interval: q.Interval,
			Candles:  q.Candles[startIndex : i+1],
		}

		values := talib.Ma(quote.Get(ma.Source), ma.InTimePeriod, ma.Type)
		candle.AddIndicator(ma.Tag, values[len(values)-1])
		q.Candles[i] = candle

		return
	}

out:
	if len(q.Candles) < ma.InTimePeriod {
		return
	}

	values := talib.Ma(q.Get(ma.Source), ma.InTimePeriod, ma.Type)
	q.AddIndicator(ma.Tag, values)

	return
}

func (ma *Ma) Is(tag IndicatorTag) bool {
	return ma.Tag == tag
}
