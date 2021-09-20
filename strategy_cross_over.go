package trade_knife

// ScoreByCrossOver scores candles by checking sources cross over condition on each of them.
func (q *Quote) ScoreByCrossOver(score float64, fastSource, slowSource Source) {
	quote := *q
	// loop through quote
	for _, candle := range quote.Candles {
		if candle.CrossedOver(fastSource, slowSource) {
			candle.Score += score
		}
	}
}
