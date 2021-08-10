package trade_knife

func HPIndicator(values []float64, lambda float64, length int) []float64 {
	result := make([]float64, len(values))
	for index := length; index <= len(values); index++ {
		source := values[index-length : index]
		hp := HPFilter(source, lambda)
		result[index-1] = hp[length-1]
	}

	return result
}
