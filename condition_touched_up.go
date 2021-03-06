package trade_knife

// TouchedUp checks if candle close above the source but the High shadow
// touched the source.
//
// H >= source
// O,L,C < source
func (c *Candle) TouchedUp(source Source) bool {
	value, _ := c.Get(source)

	return c.Open < value && c.High >= value && c.Low < value && c.Close < value
}
