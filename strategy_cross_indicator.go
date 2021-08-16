package trade_knife

// Score candles by checking indicators cross over condition on each of them.
func (q *Quote) ScoreByCrossIndicators(score float64, fastTag, slowTag IndicatorTag) {
	quote := *q
	// loop through quote
	for _, candle := range quote.Candles {
		if candle.IndicatorsCrossedOver(fastTag, slowTag) {
			candle.Score += score
		}
		if candle.IndicatorsCrossedUnder(fastTag, slowTag) {
			candle.Score -= score
		}
	}
}
