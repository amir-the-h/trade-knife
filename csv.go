package trade_knife

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// Returns csv row of the candle.
func (c *Candle) Csv(indicators ...string) (csv string) {
	// first basic records
	csv = fmt.Sprintf("%.2f,%.2f,%.2f,%.2f,%.2f,%.2f,%d,%d", c.Open, c.High, c.Low, c.Close, c.Volume, c.Score, c.Opentime.Unix(), c.Closetime.Unix())

	// get indicators values too
	for _, indicator := range indicators {
		if i, ok := c.Indicators[indicator]; ok {
			csv += fmt.Sprintf(",%.2f", i)
		} else {
			csv += ",0.00"
		}
	}

	return
}

// Returns csv formated string of whole quote.
func (q Quote) Csv(indicators ...string) (csv string) {
	if len(q) == 0 {
		return
	}
	// fix the headers
	headers := []string{"Open", "High", "Low", "Close", "Volume", "score", "Open time", "Close time"}

	// and also add the indicators as well
	for _, indicator := range indicators {
		headers = append(headers, strings.Title(string(indicator)))
	}

	csv = strings.Join(headers, ",") + "\n"

	// get each candle csv value
	for _, candle := range q {
		csv += fmt.Sprintln(candle.Csv(indicators...))
	}

	return
}

// Writes down whole quote into a csv file.
func (q Quote) WriteToCsv(filename string, indicators ...string) error {
	if len(q) == 0 {
		return ErrNotEnoughCandles
	}

	// need our file
	if filename == "" {
		filename = fmt.Sprintf("%s-%s.csv", q[0].Symbol, q[0].Interval)
	}

	return ioutil.WriteFile(filename, []byte(q.Csv()), 0644)
}
