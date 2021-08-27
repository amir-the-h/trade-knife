package trade_knife

// Check if fast source crossed under the slow source.
//
// fastSource < slowSource
// prevFastSource >= prevSlowSource
func (c *Candle) CrossedUnder(fastSource, slowSource Source) bool {
	previousCandle := c.Previous
	if previousCandle == nil {
		return false
	}

	fastValue := c.Get(Source(fastSource))
	slowValue := c.Get(Source(slowSource))
	previousFastValue := previousCandle.Get(Source(fastSource))
	previousSlowValue := previousCandle.Get(Source(slowSource))

	return fastValue < slowValue && previousFastValue >= previousSlowValue
}
