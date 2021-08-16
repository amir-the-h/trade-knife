package trade_knife

type Indicator interface {
	Add(q *Quote, c *Candle) (ok bool)
	Is(tag IndicatorTag) bool
}

type IndicatorTag string
