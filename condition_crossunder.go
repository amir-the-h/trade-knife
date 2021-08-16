package trade_knife

// FastIndicator < SlowIndicator
// PrevFastIndicator >= PrevSlowIndicator
func (c *Candle) IndicatorsCrossedUnder(fastTag, slowTag IndicatorTag) bool {
	previousCandle := c.Previous
	if previousCandle == nil {
		return false
	}

	fastIndicator := c.Get(Source(fastTag))
	slowIndicator := c.Get(Source(slowTag))
	previousFastIndicator := previousCandle.Get(Source(fastTag))
	previousSlowIndicator := previousCandle.Get(Source(slowTag))

	return fastIndicator < slowIndicator && previousFastIndicator >= previousSlowIndicator
}
