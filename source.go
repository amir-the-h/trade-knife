package trade_knife

import "strconv"

// Source is a target field on candle.
type Source string

// Retrieve the value of the target field on the candle.
func (c *Candle) Get(source Source) (float64, bool) {
	switch source {
	// single sources
	case SourceOpen:
		return c.Open, true
	case SourceHigh:
		return c.High, true
	case SourceLow:
		return c.Low, true
	case SourceClose:
		return c.Close, true
	case SourceVolume:
		return c.Volume, true

		// double sources
	case SourceOpenHigh:
		return (c.Open + c.High) / 2, true
	case SourceOpenLow:
		return (c.Open + c.Low) / 2, true
	case SourceOpenClose:
		return (c.Open + c.Close) / 2, true
	case SourceHighLow:
		return (c.High + c.Low) / 2, true
	case SourceHighClose:
		return (c.High + c.Close) / 2, true
	case SourceLowClose:
		return (c.Low + c.Close) / 2, true

		// triple sources
	case SourceOpenHighLow:
		return (c.Open + c.High + c.Low) / 3, true
	case SourceOpenHighClose:
		return (c.Open + c.High + c.Low) / 3, true
	case SourceOpenLowClose:
		return (c.Open + c.Low + c.Close) / 3, true
	case SourceHighLowClose:
		return (c.High + c.Low + c.Close) / 3, true

		// all together
	case SourceOpenHighLowClose:
		return (c.High + c.Low + c.Close) / 3, true
	}

	if value, ok := c.Indicators[IndicatorTag(source)]; ok {
		return value, true
	}

	if value, err := strconv.ParseFloat(string(source), 64); err == nil {
		return value, true
	}

	return 0., false
}

// Retrieve value of target field on all candles.
func (q *Quote) Get(source Source) []float64 {
	quote := *q
	result := make([]float64, len(quote.Candles))
	for i, candle := range quote.Candles {
		result[i], _ = candle.Get(source)
	}

	return result
}
