package trade_knife

type Hp struct {
	Tag    IndicatorTag `mapstructure:"tag"`
	Source Source       `mapstructure:"source"`
	Lambda float64      `mapstructure:"lambda"`
	Length int          `mapstructure:"length"`
}

func (hp *Hp) Add(q *Quote, c *Candle) bool {
	if c != nil {
		candle, i := q.Find(c.Opentime.Unix())
		if candle == nil {
			return false
		}

		startIndex := i - hp.Length
		if startIndex < 0 {
			return false
		}

		quote := Quote{
			Market:   q.Market,
			Symbol:   q.Symbol,
			Interval: q.Interval,
			Candles:  q.Candles[startIndex : i+1],
		}

		values := HPFilter(quote.Get(hp.Source), hp.Lambda)
		c.AddIndicator(hp.Tag, values[len(values)-1])
		q.Candles[i] = c

		return true
	}

	if len(q.Candles) < hp.Length {
		return false
	}

	for _, candle := range q.Candles {
		hp.Add(q, candle)
	}

	return true
}

func (hp *Hp) Is(tag IndicatorTag) bool {
	return hp.Tag == tag
}
