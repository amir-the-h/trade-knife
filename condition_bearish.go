package trade_knife

// Check if candle is bearish
//
// O < C,
func (c *Candle) IsBearish() bool {
	return c.Open < c.Close
}
