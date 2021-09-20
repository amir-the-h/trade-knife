package trade_knife

import (
	"encoding/csv"
	"fmt"
	"github.com/amir-the-h/goex"
	"os"
	"strconv"
	"strings"
	"time"
)

// Csv returns csv row of the candle.
func (c *Candle) Csv(indicators ...IndicatorTag) (csv string) {
	// first basic records
	csv = fmt.Sprintf("%.2f,%.2f,%.2f,%.2f,%.2f,%.2f,%s,%s", c.Open, c.High, c.Low, c.Close, c.Volume, c.Score, c.Opentime.Local().Format(time.RFC3339), c.Closetime.Local().Format(time.RFC3339))

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

// Csv returns csv formatted string of whole quote.
func (q Quote) Csv(indicators ...IndicatorTag) (csv string) {
	if len(q.Candles) == 0 {
		return
	}
	// fix the headers
	headers := []string{"Open", "High", "Low", "Close", "Volume", "Score", "Open time", "Close time"}
	var indicatorTags []IndicatorTag
	if len(indicators) > 0 {
		indicatorTags = indicators
	} else {
		indicatorTags = q.IndicatorTags()
	}
	for _, indicatorTag := range indicatorTags {
		headers = append(headers, string(indicatorTag))
	}
	csv = strings.Join(headers, ",") + "\n"
	// get each candle csv value
	for _, candle := range q.Candles {
		// and also add the indicators as well
		csv += fmt.Sprintln(candle.Csv(indicatorTags...))
	}

	return
}

// WriteToCsv writes down whole quote into a csv file.
func (q Quote) WriteToCsv(filename string, indicators ...IndicatorTag) error {
	if len(q.Candles) == 0 {
		return ErrNotEnoughCandles
	}

	// need our file
	if filename == "" {
		filename = fmt.Sprintf("%s:%s-%s.csv", q.Market, q.Currency, q.Interval)
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

// NewQuoteFromCsv reads quote from csv file.
func NewQuoteFromCsv(filename string, market MarketType, currency goex.CurrencyPair, interval Interval) (*Quote, error) {
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
		Currency: currency,
		Interval: interval,
	}
	for _, line := range csvLines {
		openPrice, _ := strconv.ParseFloat(line[indexMap["Open"]], 64)
		highPrice, _ := strconv.ParseFloat(line[indexMap["High"]], 64)
		lowPrice, _ := strconv.ParseFloat(line[indexMap["Low"]], 64)
		closePrice, _ := strconv.ParseFloat(line[indexMap["Close"]], 64)
		volume, _ := strconv.ParseFloat(line[indexMap["Volume"]], 64)
		openTime, _ := time.Parse(time.RFC3339, line[indexMap["Open time"]])
		closeTime, _ := time.Parse(time.RFC3339, line[indexMap["Close time"]])
		candle, err := NewCandle(openPrice, highPrice, lowPrice, closePrice, volume, openTime, closeTime, nil, nil)
		if err != nil {
			return nil, err
		}
		candle.Score, _ = strconv.ParseFloat(line[indexMap["Score"]], 64)

		for indicator, index := range indicatorsMap {
			candle.Indicators[IndicatorTag(indicator)], _ = strconv.ParseFloat(line[index], 64)
		}

		quote.Candles = append(quote.Candles, candle)
	}

	q := &quote
	q.Sort()
	return q, nil
}
