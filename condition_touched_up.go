package trade_knife

// Check if candle close above the indicator but the High shadow
// touched the indicator.
// H >= Indicator
// O,L,C < Indicator
func (c *Candle) TouchedUpIndicator(indicator string) bool {
	source, ok := c.Indicators[indicator]
	if !ok {
		return false
	}

	return c.Open < source && c.High >= source && c.Low < source && c.Close < source
}
