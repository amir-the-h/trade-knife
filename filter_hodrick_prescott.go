package trade_knife

import (
	"gonum.org/v1/gonum/mat"
)

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
