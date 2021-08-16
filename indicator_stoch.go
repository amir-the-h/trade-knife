package trade_knife

import "github.com/markcheno/go-talib"

type Stoch struct {
	KTag          IndicatorTag `mapstructure:"kTag"`
	DTag          IndicatorTag `mapstructure:"dTag"`
	InFastKPeriod int          `mapstructure:"kLength"`
	InSlowKPeriod int          `mapstructure:"kSmoothing"`
	InKMaType     talib.MaType `mapstructure:"kMaType"`
	InSlowDPeriod int          `mapstructure:"dSmoothing"`
	InDMaType     talib.MaType `mapstructure:"dMaType"`
}

func (s *Stoch) Add(q *Quote, c *Candle) (ok bool) {
	if c != nil {
		candle, i := q.Find(c.Opentime.Unix())
		if candle == nil {
			goto out
		}

		quote := Quote{
			Market:   q.Market,
			Symbol:   q.Symbol,
			Interval: q.Interval,
			Candles:  q.Candles[:i+1],
		}

		k, d := talib.Stoch(quote.Get(SourceHigh), quote.Get(SourceLow), quote.Get(SourceClose), s.InFastKPeriod, s.InSlowKPeriod, s.InKMaType, s.InSlowDPeriod, s.InDMaType)
		candle.AddIndicator(s.KTag, k[len(k)-1])
		candle.AddIndicator(s.DTag, d[len(d)-1])
		q.Candles[i] = candle

		return
	}

out:
	k, d := talib.Stoch(q.Get(SourceHigh), q.Get(SourceLow), q.Get(SourceClose), s.InFastKPeriod, s.InSlowKPeriod, s.InKMaType, s.InSlowDPeriod, s.InDMaType)
	q.AddIndicator(s.KTag, k)
	q.AddIndicator(s.DTag, d)

	return
}
