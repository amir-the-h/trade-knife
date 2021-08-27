package trade_knife

// Check if fast source crossed over the slow source.
//
// fastSource > slowSource
// prevfastSource <= prevslowSource
func (c *Candle) CrossedOver(fastSource, slowSource Source) bool {
	previousCandle := c.Previous
	if previousCandle == nil {
		return false
	}

	fastValue := c.Get(fastSource)
	slowValue := c.Get(slowSource)
	previousFastValue := previousCandle.Get(Source(fastSource))
	previousSlowValue := previousCandle.Get(Source(slowSource))

	return fastValue > slowValue && previousFastValue <= previousSlowValue
}
