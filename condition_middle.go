package trade_knife

// Check if indicator is passed through middle of the candle
//
// O >= Indicator,
// H > Indicator,
// L < Indicator,
// C <= Indicator,
func (c *Candle) IsIndicatorMiddle(source Source) bool {
	indicator := c.Get(source)

	return c.Open >= indicator && c.High > indicator && c.Low < indicator && c.Close <= indicator
}
