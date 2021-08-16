package trade_knife

type Indicator interface {
	Add(q *Quote, c *Candle) bool
}

type IndicatorTag string
