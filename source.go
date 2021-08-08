package trade_knife

// Source is a target field on candle.
type Source string

// Retrieve the value of the target field on the candle.
func (c *Candle) Get(source Source) (float64, error) {
	switch source {
	// single sources
	case SourceOpen:
		return c.Open, nil
	case SourceHigh:
		return c.High, nil
	case SourceLow:
		return c.Low, nil
	case SourceClose:
		return c.Close, nil
	case SourceVolume:
		return c.Volume, nil

		// double sources
	case SourceOpenHigh:
		return (c.Open + c.High) / 2, nil
	case SourceOpenLow:
		return (c.Open + c.Low) / 2, nil
	case SourceOpenClose:
		return (c.Open + c.Close) / 2, nil
	case SourceHighLow:
		return (c.High + c.Low) / 2, nil
	case SourceHighClose:
		return (c.High + c.Close) / 2, nil
	case SourceLowClose:
		return (c.Low + c.Close) / 2, nil

		// triple sources
	case SourceOpenHighLow:
		return (c.Open + c.High + c.Low) / 3, nil
	case SourceOpenHighClose:
		return (c.Open + c.High + c.Low) / 3, nil
	case SourceOpenLowClose:
		return (c.Open + c.Low + c.Close) / 3, nil
	case SourceHighLowClose:
		return (c.High + c.Low + c.Close) / 3, nil

		// all together
	case SourceOpenHighLowClose:
		return (c.High + c.Low + c.Close) / 3, nil

	default:
		return 0, ErrInvalidSource
	}
}

// Retrieve value of target field on all candles.
func (q *Quote) Get(source Source) ([]float64, error) {
	result := make([]float64, len(*q))
	for i, candle := range *q {
		v, err := candle.Get(source)
		if err != nil {
			return result, err
		}
		result[i] = v
	}

	return result, nil
}
