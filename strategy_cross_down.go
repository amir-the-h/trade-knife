package trade_knife

// Score candles by checking sources cross under condition on each of them.
func (q *Quote) ScoreByCrossDown(score float64, fastSource, slowSource Source) {
	quote := *q
	// loop through quote
	for _, candle := range quote.Candles {
		if candle.CrossedUnder(fastSource, slowSource) {
			candle.Score += score
		}
	}
}