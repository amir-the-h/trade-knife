package trade_knife

// IsBearish checks if candle is bearish
//
// O < C,
func (c *Candle) IsBearish() bool {
	return c.Open < c.Close
}
