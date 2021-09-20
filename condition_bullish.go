package trade_knife

// IsBullish checks if candle is bullish
//
// O > C,
func (c *Candle) IsBullish() bool {
	return c.Open > c.Close
}
