package trade_knife

// ScoreByTouchDown scores candles by checking if candles are touching down source.
func (q *Quote) ScoreByTouchDown(score float64, source Source) {
	quote := *q
	// loop through quote
	for _, candle := range quote.Candles {
		if candle.TouchedDown(source) {
			candle.Score += score
		}
	}
}
