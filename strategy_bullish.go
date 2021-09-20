package trade_knife

// ScoreByBullish scores candles by checking if candles are bullish.
func (q *Quote) ScoreByBullish(score float64) {
	quote := *q
	// loop through quote
	for _, candle := range quote.Candles {
		if candle.IsBullish() {
			candle.Score += score
		}
	}
}
