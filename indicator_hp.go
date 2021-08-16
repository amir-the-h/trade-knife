package trade_knife

type Hp struct {
	Tag    IndicatorTag `mapstructure:"tag"`
	Source Source       `mapstructure:"source"`
	Lambda float64      `mapstructure:"lambda"`
	Length int          `mapstructure:"length"`
}

func (hp *Hp) Add(q *Quote, c *Candle) (ok bool) {
	if c != nil {
		candle, i := q.Find(c.Opentime.Unix())
		if candle == nil {
			goto out
		}

		startIndex := i - hp.Length
		if startIndex < 0 {
			return
		}

		quote := Quote{
			Market:   q.Market,
			Symbol:   q.Symbol,
			Interval: q.Interval,
			Candles:  q.Candles[startIndex : i+1],
		}

		values := HPFilter(quote.Get(hp.Source), hp.Lambda)
		candle.AddIndicator(hp.Tag, values[len(values)-1])
		q.Candles[i] = candle

		return
	}

out:
	if len(q.Candles) < hp.Length {
		return
	}

	for _, candle := range q.Candles {
		hp.Add(q, candle)
	}

	return
}

func (hp *Hp) Is(tag IndicatorTag) bool {
	return hp.Tag == tag
}
