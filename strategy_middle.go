package trade_knife

// ScoreByMiddle scores candles by checking if candles are middle of the source.
func (q *Quote) ScoreByMiddle(score float64, source Source) {
	quote := *q
	// loop through quote
	for _, candle := range quote.Candles {
		if candle.IsMiddle(source) {
			candle.Score += score
		}
	}
}
