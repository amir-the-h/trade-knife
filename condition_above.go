package trade_knife

// Check if candle is above of the indicator
//
// O,H,L,C > Indicator
func (c *Candle) IsAboveIndicator(source Source) bool {
	indicator := c.Get(source)

	return c.Open > indicator && c.High > indicator && c.Low > indicator && c.Close > indicator
}
