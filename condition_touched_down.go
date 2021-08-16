package trade_knife

// Check if candle closed above the indicator but the Low shadow
// touched the indicator.
//
// O,H,C > Indicator
// L <= Indicator
func (c *Candle) TouchedDownIndicator(source Source) bool {
	indicator := c.Get(source)

	return c.Open > indicator && c.High > indicator && c.Low <= indicator && c.Close > indicator
}
