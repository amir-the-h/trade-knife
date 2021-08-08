package trade_knife

// Check if candle is above of the indicator
//
// O,H,L,C > Indicator
func (c *Candle) IsAboveIndicator(indicator string) bool {
	source, ok := c.Indicators[indicator]
	if !ok {
		return false
	}

	return c.Open > source && c.High > source && c.Low > source && c.Close > source
}
