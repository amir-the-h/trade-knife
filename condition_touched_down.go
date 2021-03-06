package trade_knife

// TouchedDown checks if candle closed above the source but the Low shadow
// touched the source.
//
// O,H,C > source
// L <= source
func (c *Candle) TouchedDown(source Source) bool {
	value, _ := c.Get(source)

	return c.Open > value && c.High > value && c.Low <= value && c.Close > value
}
