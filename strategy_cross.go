package trade_knife

// ScoreByCross scores candles by checking sources cross over/under condition on each of them.
func (q *Quote) ScoreByCross(score float64, fastSource, slowSource Source) {
	quote := *q
	// loop through quote
	for _, candle := range quote.Candles {
		if candle.CrossedOver(fastSource, slowSource) {
			candle.Score += score
		}
		if candle.CrossedUnder(fastSource, slowSource) {
			candle.Score -= score
		}
	}
}
