package trade_knife

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

// Returns csv row of the candle.
func (c *Candle) Csv(indicators ...string) (csv string) {
	// first basic records
	csv = fmt.Sprintf("%s,%s,%.2f,%.2f,%.2f,%.2f,%.2f,%.2f,%d,%d", c.Symbol, c.Interval, c.Open, c.High, c.Low, c.Close, c.Volume, c.Score, c.Opentime.Unix(), c.Closetime.Unix())

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
	headers := []string{"Symbol", "Interval", "Open", "High", "Low", "Close", "Volume", "Score", "Open time", "Close time"}

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

// Read quote from csv file.
func NewQuoteFromCsv(filename string) (*Quote, error) {
	csvFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	fmt.Println("Successfully Opened CSV file")
	defer csvFile.Close()

	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		fmt.Println(err)
	}
	var quote Quote
	// "Symbol", "Interval", "Open", "High", "Low", "Close", "Volume", "Score", "Open time", "Close time"
	for _, line := range csvLines {
		symbol := line[0]
		interval := Interval(line[1])
		open, _ := strconv.ParseFloat(line[2], 64)
		high, _ := strconv.ParseFloat(line[3], 64)
		low, _ := strconv.ParseFloat(line[4], 64)
		close, _ := strconv.ParseFloat(line[5], 64)
		volume, _ := strconv.ParseFloat(line[6], 64)
		openTimestamp, _ := strconv.ParseFloat(line[8], 64)
		closeTimestamp, _ := strconv.ParseFloat(line[9], 64)
		openTime := time.Unix(int64(openTimestamp), 0)
		closeTime := time.Unix(int64(closeTimestamp), 0)
		candle, err := NewCandle(symbol, open, high, low, close, volume, openTime, closeTime, interval, nil, nil)
		if err != nil {
			return nil, err
		}
		candle.Score, _ = strconv.ParseFloat(line[7], 64)

		quote = append(quote, candle)
	}

	q := &quote
	q.Sort()
	return q, nil
}
