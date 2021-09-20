package trade_knife

// ScoreByBelow scores candles by checking if candles are below source.
func (q *Quote) ScoreByBelow(score float64, source Source) {
	quote := *q
	// loop through quote
	for _, candle := range quote.Candles {
		if candle.IsBelow(source) {
			candle.Score += score
		}
	}
}
