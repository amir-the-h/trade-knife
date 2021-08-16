package trade_knife

// Check if candle is below of the indicator.
//
// O,H,L,C < Indicator
func (c *Candle) IsBelowIndicator(source Source) bool {
	indicator := c.Get(source)

	return c.Open < indicator && c.High < indicator && c.Low < indicator && c.Close < indicator
}
