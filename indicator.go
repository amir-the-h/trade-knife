package trade_knife

type Indicator interface {
	AddToQuote(q *Quote, c *Candle) bool
}

type IndicatorTag string
