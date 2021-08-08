package trade_knife

// Check if candle is below of the indicator.
//
// O,H,L,C < Indicator
func (c *Candle) IsBelowIndicator(indicator string) bool {
	source, ok := c.Indicators[indicator]
	if !ok {
		return false
	}

	return c.Open < source && c.High < source && c.Low < source && c.Close < source
}
