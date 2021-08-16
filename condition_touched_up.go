package trade_knife

// Check if candle close above the indicator but the High shadow
// touched the indicator.
// H >= Indicator
// O,L,C < Indicator
func (c *Candle) TouchedUpIndicator(source Source) bool {
	indicator := c.Get(source)

	return c.Open < indicator && c.High >= indicator && c.Low < indicator && c.Close < indicator
}
