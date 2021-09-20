package trade_knife

// ScoreByAbove scores candles by checking if candles are above source.
func (q *Quote) ScoreByAbove(score float64, source Source) {
	quote := *q
	// loop through quote
	for _, candle := range quote.Candles {
		if candle.IsAbove(source) {
			candle.Score += score
		}
	}
}
