package trade_knife

// CrossedUnder checks if fast source crossed under the slow source.
//
// fastSource < slowSource
// prevFastSource >= prevSlowSource
func (c *Candle) CrossedUnder(fastSource, slowSource Source) bool {
	previousCandle := c.Previous
	if previousCandle == nil {
		return false
	}

	fastValue, ok := c.Get(Source(fastSource))
	if !ok {
		return false
	}
	slowValue, ok := c.Get(Source(slowSource))
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

	return fastValue < slowValue && previousFastValue >= previousSlowValue
}
