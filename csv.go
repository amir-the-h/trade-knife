package trade_knife

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Returns csv row of the candle.
func (c *Candle) Csv(indicators ...string) (csv string) {
	// first basic records
	csv = fmt.Sprintf("%.2f,%.2f,%.2f,%.2f,%.2f,%.2f,%s,%s", c.Open, c.High, c.Low, c.Close, c.Volume, c.Score, c.Opentime.UTC().Format(time.RFC3339), c.Closetime.UTC().Format(time.RFC3339))

	// get indicators values too
	if len(indicators) > 0 {
		for _, indicator := range indicators {
			csv += fmt.Sprintf(",%.2f", c.Indicators[indicator])
		}
	} else {
		for _, indicator := range c.Indicators {
			csv += fmt.Sprintf(",%.2f", indicator)
		}
	}

	return
}

// Returns csv formated string of whole quote.
func (q Quote) Csv(indicators ...string) (csv string) {
	if len(q.Candles) == 0 {
		return
	}
	// fix the headers
	headers := []string{"Open", "High", "Low", "Close", "Volume", "Score", "Open time", "Close time"}
	var indicatorNames []string
	if len(indicators) > 0 {
		indicatorNames = indicators
	} else {
		indicatorNames = q.IndicatorNames()
	}
	headers = append(headers, indicatorNames...)
	csv = strings.Join(headers, ",") + "\n"
	// get each candle csv value
	for _, candle := range q.Candles {
		// and also add the indicators as well
		csv += fmt.Sprintln(candle.Csv(indicatorNames...))
	}

	return
}

// Writes down whole quote into a csv file.
func (q Quote) WriteToCsv(filename string, indicators ...string) error {
	if len(q.Candles) == 0 {
		return ErrNotEnoughCandles
	}

	// need our file
	if filename == "" {
		filename = fmt.Sprintf("%s:%s-%s.csv", q.Market, q.Symbol, q.Interval)
	}

	// open or create the file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	// truncate the file
	if err = file.Truncate(0); err != nil {
		return err
	}

	_, err = file.Write([]byte(q.Csv(indicators...)))
	return err
}

// Read quote from csv file.
func NewQuoteFromCsv(filename string, market MarketType, symbol string, interval Interval) (*Quote, error) {
	csvFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer csvFile.Close()

	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		return nil, err
	}
	headers := csvLines[0]
	csvLines = csvLines[1:]
	indexMap := make(map[string]int)
	indicatorsMap := make(map[string]int)
	for i, header := range headers {
		switch header {
		case "Open", "High", "Low", "Close", "Volume", "Score", "Open time", "Close time":
			indexMap[header] = i
		default:
			indicatorsMap[header] = i
		}
	}
	quote := Quote{
		Market:   market,
		Symbol:   symbol,
		Interval: interval,
	}
	for _, line := range csvLines {
		open, _ := strconv.ParseFloat(line[indexMap["Open"]], 64)
		high, _ := strconv.ParseFloat(line[indexMap["High"]], 64)
		low, _ := strconv.ParseFloat(line[indexMap["Low"]], 64)
		close, _ := strconv.ParseFloat(line[indexMap["Close"]], 64)
		volume, _ := strconv.ParseFloat(line[indexMap["Volume"]], 64)
		openTime, _ := time.Parse(time.RFC3339, line[indexMap["Open time"]])
		closeTime, _ := time.Parse(time.RFC3339, line[indexMap["Close time"]])
		candle, err := NewCandle(open, high, low, close, volume, openTime, closeTime, nil, nil)
		if err != nil {
			return nil, err
		}
		candle.Score, _ = strconv.ParseFloat(line[indexMap["Score"]], 64)

		for indicator, index := range indicatorsMap {
			candle.Indicators[indicator], _ = strconv.ParseFloat(line[index], 64)
		}

		quote.Candles = append(quote.Candles, candle)
	}

	q := &quote
	q.Sort()
	return q, nil
}
