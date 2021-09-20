package trade_knife

import "github.com/markcheno/go-talib"

type Stoch struct {
	Tag           IndicatorTag `mapstructure:"tag"`
	KTag          IndicatorTag `mapstructure:"kTag"`
	DTag          IndicatorTag `mapstructure:"dTag"`
	InFastKPeriod int          `mapstructure:"kLength"`
	InSlowKPeriod int          `mapstructure:"kSmoothing"`
	InKMaType     talib.MaType `mapstructure:"kMaType"`
	InSlowDPeriod int          `mapstructure:"dSmoothing"`
	InDMaType     talib.MaType `mapstructure:"dMaType"`
}

func (s *Stoch) Add(q *Quote, c *Candle) bool {
	if c != nil {
		candle, i := q.Find(c.Opentime.Unix())
		if candle == nil {
			return false
		}

		quote := Quote{
			Market:   q.Market,
			Currency:   q.Currency,
			Interval: q.Interval,
			Candles:  q.Candles[:i+1],
		}

		k, d := talib.Stoch(quote.Get(SourceHigh), quote.Get(SourceLow), quote.Get(SourceClose), s.InFastKPeriod, s.InSlowKPeriod, s.InKMaType, s.InSlowDPeriod, s.InDMaType)
		c.AddIndicator(s.KTag, k[len(k)-1])
		c.AddIndicator(s.DTag, d[len(d)-1])
		q.Candles[i] = c

		return true
	}

	k, d := talib.Stoch(q.Get(SourceHigh), q.Get(SourceLow), q.Get(SourceClose), s.InFastKPeriod, s.InSlowKPeriod, s.InKMaType, s.InSlowDPeriod, s.InDMaType)
	q.AddIndicator(s.KTag, k)
	q.AddIndicator(s.DTag, d)

	return true
}

func (s *Stoch) Is(tag IndicatorTag) bool {
	return s.Tag == tag
}
