package trade_knife

// Score candles by checking support/resistance line reaction.
func (q *Quote) ScoreBySupportResistance(score float64, source Source) {
	quote := *q
	// loop through quote
	for _, candle := range quote.Candles {
		checker := func(c *Candle) float64 {
			if c.IsBearish() {
				score *= -1
			}
			if c.TouchedDown(source) || c.IsMiddle(source) || c.IsMiddle(source) {
				return score
			}
			return 0
		}

		result := checker(candle)
		if result != 0 {
			candle.Score += result
		}
		if candle.Previous == nil {
			continue
		}
		prevResult := checker(candle.Previous)
		if prevResult*-1 == result {
			candle.Score += result
		}
	}
}
