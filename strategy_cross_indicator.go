package trade_knife

// Score candles by checking indicators cross over condition on each of them.
func (q *Quote) ScoreByCrossIndicators(score float64, fastIndicator, slowIndicator string, fastSource, slowSource Source) {
	quote := *q
	// loop through quote
	for _, candle := range quote.Candles {
		if candle.IndicatorsCrossedOver(fastIndicator, slowIndicator) {
			candle.Score += score
		}
		if candle.IndicatorsCrossedUnder(fastIndicator, slowIndicator) {
			candle.Score -= score
		}
	}
}
