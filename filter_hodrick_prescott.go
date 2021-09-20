package trade_knife

import (
	"gonum.org/v1/gonum/mat"
)

// HPFilter filters the given values by Hodrick-Prescott filter.
//
// https://www.google.com/url?sa=t&rct=j&q=&esrc=s&source=web&cd=&cad=rja&uact=8&ved=2ahUKEwjorp3ls6byAhUHDOwKHX27DwgQFnoECAcQAQ&url=https%3A%2F%2Fen.wikipedia.org%2Fwiki%2FHodrick%25E2%2580%2593Prescott_filter&usg=AOvVaw24zh7DousHxEoH5CpSHpIJ
func HPFilter(values []float64, lambda float64) []float64 {
	length := len(values)
	lastIndex := length - 1
	F := mat.NewDense(length, length, nil)
	for x := 0; x < length; x++ {
		F.Set(x, x, 6*lambda+1)
		if x+1 >= 0 && x+1 <= lastIndex {
			F.Set(x+1, x, -4*lambda)
			F.Set(x, x+1, -4*lambda)
		}
		if x+2 >= 0 && x+2 <= lastIndex {
			F.Set(x+2, x, 1*lambda)
			F.Set(x, x+2, 1*lambda)
		}
	}

	F.Set(0, 0, 1*lambda+1)
	F.Set(lastIndex, lastIndex, 1*lambda)
	F.Set(1, 1, 5*lambda+1)
	F.Set(lastIndex-1, lastIndex-1, 5*lambda+1)
	F.Set(1, 0, -2*lambda)
	F.Set(0, 1, -2*lambda)
	F.Set(lastIndex, lastIndex-1, -2*lambda)
	F.Set(lastIndex-1, lastIndex, -2*lambda)

	var Fi mat.Dense
	err := Fi.Inverse(F)
	if err != nil {
		panic(err)
	}

	V := mat.NewDense(length, 1, values)
	var C mat.Dense
	C.Mul(&Fi, V)

	result := make([]float64, length)
	for i := range result {
		result[i] = C.At(i, 0)
	}

	return result
}
