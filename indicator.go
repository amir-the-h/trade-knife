package trade_knife

// Indicator indicates how an indicator should be implemented
type Indicator interface {
	Add(q *Quote, c *Candle) (ok bool)
	Is(tag IndicatorTag) bool
}

type IndicatorTag string
