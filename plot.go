package trade_knife

import (
	"fmt"
	"image/color"
	"io"
	"log"
	"math/rand"

	"github.com/pplcc/plotext/custplotter"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

// Returns a writer of candlestick plot image.
func Plot(quote Quote, width, height vg.Length, extension string, tags []IndicatorTag) (io.WriterTo, error) {
	data := make(custplotter.TOHLCVs, len(quote.Candles))

	for i, candle := range quote.Candles {
		data[i].T = float64(candle.Opentime.Unix())
		data[i].O = candle.Open
		data[i].H = candle.High
		data[i].L = candle.Low
		data[i].C = candle.Close
	}

	plt := plot.New()

	plt.Title.Text = fmt.Sprintf("%s %s", quote.Symbol, quote.Interval)
	plt.X.Label.Text = "Time"
	plt.Y.Label.Text = "USDT"

	bars, err := custplotter.NewCandlesticks(data)
	if err != nil {
		log.Panic(err)
	}

	grid := plotter.NewGrid()
	plt.Add(bars, grid)
	plt.X.Tick.Marker = plot.TimeTicks{Format: "2006-01-02\n15:04"}
	for _, tag := range tags {
		plotIndicators(plt, quote, tag)
	}
	plt.Legend.ThumbnailWidth = width * 0.1
	plt.Legend.Top = true
	plt.Legend.Left = true
	plt.Legend.TextStyle.Font.Size = font.Centimeter * 1
	plt.Legend.Padding = font.Centimeter * 2

	ioW, err := plt.WriterTo(width, height, extension)
	if err != nil {
		return nil, err
	}

	return ioW, nil
}

// Draw indicators on the main plot.
func plotIndicators(plt *plot.Plot, quote Quote, tag IndicatorTag) {
	if len(quote.Candles) == 0 {
		return
	}
	plot := plotter.NewFunction(func(f float64) float64 {
		candle, _ := quote.Find(int64(f))
		if candle == nil {
			return 0
		}
		v, ok := candle.Indicators[tag]
		if !ok {
			return 0
		}
		return v
	})
	plot.Samples = len(quote.Candles)
	plot.Color = color.RGBA{R: uint8(rand.Intn(255)), G: uint8(rand.Intn(255)), B: uint8(rand.Intn(255)), A: 255}

	plt.Add(plot)
	plt.Legend.Add(string(tag), plot)
}

// TODO: Add signal plotter
