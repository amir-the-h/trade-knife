package trade_knife

// Check if indicator is passed through middle of the candle
//
// O >= Indicator,
// H > Indicator,
// L < Indicator,
// C <= Indicator,
func (c *Candle) IsIndicatorMiddle(indicator string) bool {
	source, ok := c.Indicators[indicator]
	if !ok {
		return false
	}

	return c.Open >= source && c.High > source && c.Low < source && c.Close <= source
}
