package trade_knife

// IsBelow checks if candle is below of the source.
//
// O,H,L,C < source
func (c *Candle) IsBelow(source Source) bool {
	value, _ := c.Get(source)

	return c.Open < value && c.High < value && c.Low < value && c.Close < value
}
