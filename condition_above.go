package trade_knife

// IsAbove checks if candle is above of the source
//
// O,H,L,C > source
func (c *Candle) IsAbove(source Source) bool {
	value, _ := c.Get(source)

	return c.Open > value && c.High > value && c.Low > value && c.Close > value
}
