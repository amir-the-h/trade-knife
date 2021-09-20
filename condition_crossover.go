package trade_knife

// CrossedOver checks if fast source crossed over the slow source.
//
// fastSource > slowSource
// prevfastSource <= prevslowSource
func (c *Candle) CrossedOver(fastSource, slowSource Source) bool {
	previousCandle := c.Previous
	if previousCandle == nil {
		return false
	}

	fastValue, ok := c.Get(fastSource)
	if !ok {
		return false
	}
	slowValue, ok := c.Get(slowSource)
	if !ok {
		return false
	}
	previousFastValue, ok := previousCandle.Get(Source(fastSource))
	if !ok {
		return false
	}
	previousSlowValue, ok := previousCandle.Get(Source(slowSource))
	if !ok {
		return false
	}

	return fastValue > slowValue && previousFastValue <= previousSlowValue
}
