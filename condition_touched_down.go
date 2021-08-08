package trade_knife

// Check if candle closed above the indicator but the Low shadow
// touched the indicator.
//
// O,H,C > Indicator
// L <= Indicator
func (c *Candle) TouchedDownIndicator(indicator string) bool {
	source, ok := c.Indicators[indicator]
	if !ok {
		return false
	}

	return c.Open > source && c.High > source && c.Low <= source && c.Close > source
}
