package trade_knife

// Score candles by checking indicators cross over condition on each of them.
func (q *Quote) ScoreByCrossIndicators(score float64, fastIndicator, slowIndicator string, fastSource, slowSource Source) {
	// loop through quote
	for _, candle := range *q {
		if candle.IndicatorsCrossedOver(fastIndicator, slowIndicator) {
			candle.Score += score
		}
		if candle.IndicatorsCrossedUnder(fastIndicator, slowIndicator) {
			candle.Score -= score
		}
	}
}
