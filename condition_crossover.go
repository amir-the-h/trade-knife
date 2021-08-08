package trade_knife

// Check if fast indicator crossed over the slow indicator.
//
// FastIndicator > SlowIndicator
// PrevFastIndicator <= PrevSlowIndicator
func (c *Candle) IndicatorsCrossedOver(fastIndicator, slowIndicator string) bool {
	previousCandle := c.Previous
	if previousCandle == nil {
		return false
	}

	fastSource, ok := c.Indicators[fastIndicator]
	if !ok {
		return false
	}

	slowSource, ok := c.Indicators[slowIndicator]
	if !ok {
		return false
	}

	previousFastSource, ok := previousCandle.Indicators[fastIndicator]
	if !ok {
		return false
	}
	previousSlowSource, ok := previousCandle.Indicators[slowIndicator]
	if !ok {
		return false
	}

	return fastSource > slowSource && previousFastSource <= previousSlowSource
}
