package trade_knife

// IsMiddle checks if source is passed through middle of the candle
//
// O >= source,
// H > source,
// L < source,
// C <= source,
func (c *Candle) IsMiddle(source Source) bool {
	value, _ := c.Get(source)

	return c.Open >= value && c.High > value && c.Low < value && c.Close <= value
}
