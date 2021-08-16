package trade_knife

// Source is a target field on candle.
type Source string

// Retrieve the value of the target field on the candle.
func (c *Candle) Get(source Source) float64 {
	switch source {
	// single sources
	case SourceOpen:
		return c.Open
	case SourceHigh:
		return c.High
	case SourceLow:
		return c.Low
	case SourceClose:
		return c.Close
	case SourceVolume:
		return c.Volume

		// double sources
	case SourceOpenHigh:
		return (c.Open + c.High) / 2
	case SourceOpenLow:
		return (c.Open + c.Low) / 2
	case SourceOpenClose:
		return (c.Open + c.Close) / 2
	case SourceHighLow:
		return (c.High + c.Low) / 2
	case SourceHighClose:
		return (c.High + c.Close) / 2
	case SourceLowClose:
		return (c.Low + c.Close) / 2

		// triple sources
	case SourceOpenHighLow:
		return (c.Open + c.High + c.Low) / 3
	case SourceOpenHighClose:
		return (c.Open + c.High + c.Low) / 3
	case SourceOpenLowClose:
		return (c.Open + c.Low + c.Close) / 3
	case SourceHighLowClose:
		return (c.High + c.Low + c.Close) / 3

		// all together
	case SourceOpenHighLowClose:
		return (c.High + c.Low + c.Close) / 3
	}

	if value, ok := c.Indicators[IndicatorTag(source)]; ok {
		return value
	}

	return 0.
}

// Retrieve value of target field on all candles.
func (q *Quote) Get(source Source) []float64 {
	quote := *q
	result := make([]float64, len(quote.Candles))
	for i, candle := range quote.Candles {
		result[i] = candle.Get(source)
	}

	return result
}
