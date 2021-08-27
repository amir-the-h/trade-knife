package trade_knife

// Check if candle is bullish
//
// O > C,
func (c *Candle) IsBullish() bool {
	return c.Open > c.Close
}
