package trade_knife

// Score candles by checking if candles are touching up source.
func (q *Quote) ScoreByTouchUp(score float64, source Source) {
	quote := *q
	// loop through quote
	for _, candle := range quote.Candles {
		if candle.TouchedUp(source) {
			candle.Score += score
		}
	}
}
