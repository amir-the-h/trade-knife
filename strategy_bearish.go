package trade_knife

// Score candles by if candle is bearish
func (q *Quote) ScoreByBearish(score float64) {
	quote := *q
	// loop through quote
	for _, candle := range quote.Candles {
		if candle.IsBearish() {
			candle.Score += score
		}
	}
}
