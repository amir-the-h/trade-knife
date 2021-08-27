package trade_knife

// Check if candle is above of the source
//
// O,H,L,C > source
func (c *Candle) IsAbove(source Source) bool {
	value := c.Get(source)

	return c.Open > value && c.High > value && c.Low > value && c.Close > value
}
