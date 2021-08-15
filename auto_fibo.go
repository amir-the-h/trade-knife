package trade_knife

import (
	"math"

	"github.com/markcheno/go-talib"
)

func AutoFibo(inHigh, inLow, inClose, ratios []float64, depth int, deviation float64) []map[float64]float64 {
	result := make([]map[float64]float64, len(inClose))

	for i := len(inClose) - 1; i > depth*2; i-- {
		outHigh := inHigh[:i]
		outLow := inLow[:i]
		outClose := inClose[:i]
		atrs := talib.Atr(outHigh, outLow, outClose, 10)
		atr := atrs[len(atrs)-1]
		deviationThreshold := atr / inClose[i] * 100. * deviation
		pivots := pivots(outHigh, outLow, deviationThreshold, depth/2)
		if len(pivots) != 2 {
			result[i] = map[float64]float64{}
			continue
		}

		row := make(map[float64]float64)
		row[0] = pivots[0]
		row[100] = pivots[1]
		diff := row[100] - row[0]
		for _, ratio := range ratios {
			row[ratio] = row[0] + (diff * ratio / 100)
			row[ratio*-1] = row[0] - (diff * ratio / 100)
		}

		result[i] = row
	}

	return result
}

func calcDev(base_price, price float64) float64 {
	return 100 * (price - base_price) / price
}

func pivots(inHigh, inLow []float64, deviationThreshold float64, length int) (pivots []float64) {
	if len(inHigh) != len(inLow) || len(inHigh) <= length*2 {
		return
	}

	var (
		highPivot, lowPivot float64
		highFixed, lowFixed bool
	)

	for i := len(inHigh) - length - 1; i > length; i-- {
		// initialize compare point
		highTarget := inHigh[i]
		lowTarget := inLow[i]

		// prepare targets batches
		hrTargets := inHigh[i+1 : i+length+1]
		hlTargets := inHigh[i-length : i]
		lrTargets := inLow[i+1 : i+length+1]
		llTargets := inLow[i-length : i]
		isHigh := true
		isLow := true

		for o := range hrTargets {
			if highPivot == 0 && isHigh && (hrTargets[o] > highTarget || hlTargets[o] > highTarget) {
				isHigh = false
			}

			if lowPivot == 0 && isLow && (lrTargets[o] < lowTarget || llTargets[o] < lowTarget) {
				isLow = false
			}
		}

		if isHigh && !highFixed {
			highPivot = highTarget
			pivots = append(pivots, highPivot)
			highFixed = true
		}

		if isLow && !lowFixed {
			lowPivot = lowTarget
			pivots = append(pivots, lowPivot)
			lowFixed = true
		}

		if len(pivots) == 2 {
			deviation := calcDev(pivots[0], pivots[1])
			if math.Abs(deviation) < deviationThreshold {
				if pivots[0] == highPivot {
					lowPivot = 0
					lowFixed = false
				} else {
					highPivot = 0
					highFixed = false
				}
				pivots = pivots[:1]
				continue
			}

			return
		}
	}

	return
}
