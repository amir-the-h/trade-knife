package trade_knife

// Check if fast indicator crossed over the slow indicator.
//
// fastIndicator > slowIndicator
// PrevfastIndicator <= PrevslowIndicator
func (c *Candle) IndicatorsCrossedOver(fastTag, slowTag IndicatorTag) bool {
	previousCandle := c.Previous
	if previousCandle == nil {
		return false
	}

	fastIndicator := c.Get(Source(fastTag))
	slowIndicator := c.Get(Source(slowTag))
	previousFastIndicator := previousCandle.Get(Source(fastTag))
	previousSlowIndicator := previousCandle.Get(Source(slowTag))

	return fastIndicator > slowIndicator && previousFastIndicator <= previousSlowIndicator
}
