package trade_knife

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/delivery"
	"github.com/adshao/go-binance/v2/futures"
)

// Fetches quote from binance spot market.
func NewQuoteFromBinanceSpot(apiKey, secretKey, symbol string, interval Interval, openTimestamp ...int64) (*Quote, error) {
	var (
		quote                Quote
		tots, tcts, ots, cts time.Time
	)

	client := binance.NewClient(apiKey, secretKey)
	request := client.NewKlinesService().Symbol(symbol).Interval(string(interval))
	direction := -1
	if len(openTimestamp) > 0 {
		direction = 1
		request.StartTime(openTimestamp[0] * 1000)
	}
	klines, err := request.Do(context.Background())
	if err != nil {
		return &quote, err
	}

	for len(klines) > 0 {
		tots = ots
		tcts = cts
		ots = time.Unix(int64(klines[0].OpenTime/1000), 0).Add(time.Hour * time.Duration(direction) * 24 * 90).UTC()
		cts = time.Unix(int64(klines[0].CloseTime/1000), 0).Add(interval.Duration()).UTC()
		if tots == ots || tcts == cts {
			break
		}

		for _, kline := range klines {
			candle, err := createCandleFromKline(MarketSpot, symbol, kline.Open, kline.High, kline.Low, kline.Close, kline.Volume, interval, kline.OpenTime, kline.CloseTime)
			if err != nil {
				return &quote, err
			}
			quote = append(quote, candle)
		}
		klines, err = request.StartTime(ots.Unix() * 1000).EndTime(cts.Unix() * 1000).Do(context.Background())
		if err != nil {
			return &quote, err
		}
	}

	q := &quote
	q.Sort()
	return q, nil
}

// Fetches quote from binace futures market.
func NewQuoteFromBinanceFutures(apiKey, secretKey, symbol string, interval Interval, openTimestamp ...int64) (*Quote, error) {
	var (
		quote                Quote
		tots, tcts, ots, cts time.Time
	)

	client := futures.NewClient(apiKey, secretKey)
	request := client.NewKlinesService().Symbol(symbol).Interval(string(interval))
	direction := -1
	if len(openTimestamp) > 0 {
		direction = 1
		request.StartTime(openTimestamp[0] * 1000)
	}
	klines, err := request.Do(context.Background())
	if err != nil {
		return &quote, err
	}

	for len(klines) > 0 {
		tots = ots
		tcts = cts
		ots = time.Unix(int64(klines[0].OpenTime/1000), 0).Add(time.Hour * time.Duration(direction) * 24 * 90).UTC()
		cts = time.Unix(int64(klines[0].CloseTime/1000), 0).Add(interval.Duration()).UTC()
		if tots == ots || tcts == cts {
			break
		}

		for _, kline := range klines {
			candle, err := createCandleFromKline(MarketFutures, symbol, kline.Open, kline.High, kline.Low, kline.Close, kline.Volume, interval, kline.OpenTime, kline.CloseTime)
			if err != nil {
				return &quote, err
			}
			quote = append(quote, candle)
		}
		klines, err = request.StartTime(ots.Unix() * 1000).EndTime(cts.Unix() * 1000).Do(context.Background())
		if err != nil {
			return &quote, err
		}
	}

	q := &quote
	q.Sort()
	return q, nil
}

// Fetches quote from binace delivery market.
func NewQuoteFromBinanceDelivery(apiKey, secretKey, symbol string, interval Interval, openTimestamp ...int64) (*Quote, error) {
	var (
		quote                Quote
		tots, tcts, ots, cts time.Time
	)

	client := delivery.NewClient(apiKey, secretKey)
	request := client.NewKlinesService().Symbol(symbol).Interval(string(interval))
	direction := -1
	if len(openTimestamp) > 0 {
		direction = 1
		request.StartTime(openTimestamp[0] * 1000)
	}
	klines, err := request.Do(context.Background())
	if err != nil {
		return &quote, err
	}

	for len(klines) > 0 {
		tots = ots
		tcts = cts
		ots = time.Unix(int64(klines[0].OpenTime/1000), 0).Add(time.Hour * time.Duration(direction) * -24 * 90).UTC()
		cts = time.Unix(int64(klines[0].CloseTime/1000), 0).Add(interval.Duration()).UTC()
		if tots == ots || tcts == cts {
			break
		}

		for _, kline := range klines {
			candle, err := createCandleFromKline(MarketDelivery, symbol, kline.Open, kline.High, kline.Low, kline.Close, kline.Volume, interval, kline.OpenTime, kline.CloseTime)
			if err != nil {
				return &quote, err
			}
			quote = append(quote, candle)
		}

		klines, err = request.StartTime(ots.Unix() * 1000).EndTime(cts.Unix() * 1000).Do(context.Background())
		if err != nil {
			return &quote, err
		}
	}

	q := &quote
	q.Sort()
	return q, nil
}

// Fetch all spot candles after last candle including itself.
func (q *Quote) RefreshBinanceSpot(apiKey, secretKey string) error {
	quote := *q
	if len(quote) == 0 {
		return errors.New("won't be able to refresh an empty quote")
	}

	lastCandle := quote[len(quote)-1]
	openTime := lastCandle.Opentime.Unix()
	fetchedQuote, err := NewQuoteFromBinanceSpot(apiKey, secretKey, lastCandle.Symbol, lastCandle.Interval, openTime)
	if err != nil {
		return err
	}

	q.Merge(fetchedQuote)

	return nil
}

// Fetch all futures candles after last candle including itself.
func (q *Quote) RefreshBinanceFutures(apiKey, secretKey string) error {
	quote := *q
	if len(quote) == 0 {
		return errors.New("won't be able to refresh an empty quote")
	}

	lastCandle := quote[len(quote)-1]
	openTime := lastCandle.Opentime.Unix()
	fetchedQuote, err := NewQuoteFromBinanceFutures(apiKey, secretKey, lastCandle.Symbol, lastCandle.Interval, openTime)
	if err != nil {
		return err
	}

	q.Merge(fetchedQuote)

	return nil
}

// Fetch all delivery candles after last candle including itself.
func (q *Quote) RefreshBinanceDelivery(apiKey, secretKey string) error {
	quote := *q
	if len(quote) == 0 {
		return errors.New("won't be able to refresh an empty quote")
	}

	lastCandle := quote[len(quote)-1]
	openTime := lastCandle.Opentime.Unix()
	fetchedQuote, err := NewQuoteFromBinanceDelivery(apiKey, secretKey, lastCandle.Symbol, lastCandle.Interval, openTime)
	if err != nil {
		return err
	}

	q.Merge(fetchedQuote)

	return nil
}

// Will sync quote with latest binance spot kline info.
func (q *Quote) SyncBinanceSpot(update CandleChannel) (doneC chan struct{}, err error) {
	quote := *q
	if len(quote) == 0 {
		return nil, errors.New("won't be able to sync an empty quote")
	}

	lastCandle := quote[len(quote)-1]
	wsKlineHandler := func(event *binance.WsKlineEvent) {
		kline := event.Kline
		o, _ := strconv.ParseFloat(kline.Open, 64)
		h, _ := strconv.ParseFloat(kline.High, 64)
		l, _ := strconv.ParseFloat(kline.Low, 64)
		c, _ := strconv.ParseFloat(kline.Close, 64)
		v, _ := strconv.ParseFloat(kline.Volume, 64)
		ot := time.Unix(int64(kline.StartTime/1000), 0).UTC()
		ct := time.Unix(int64(kline.EndTime/1000), 0).UTC()
		candle, err := q.Sync(lastCandle.Market, lastCandle.Symbol, lastCandle.Interval, o, h, l, c, v, ot, ct)
		if err != nil {
			return
		}
		update <- candle
	}
	errHandler := func(err error) {
		fmt.Println(err)
	}
	doneC, _, err = binance.WsKlineServe(lastCandle.Symbol, string(lastCandle.Interval), wsKlineHandler, errHandler)

	return
}

// Will sync quote with latest binance futures kline info.
func (q *Quote) SyncBinanceFutures(update CandleChannel) (doneC chan struct{}, err error) {
	quote := *q
	if len(quote) == 0 {
		return nil, errors.New("won't be able to sync an empty quote")
	}

	lastCandle := quote[len(quote)-1]
	wsKlineHandler := func(event *futures.WsKlineEvent) {
		kline := event.Kline
		o, _ := strconv.ParseFloat(kline.Open, 64)
		h, _ := strconv.ParseFloat(kline.High, 64)
		l, _ := strconv.ParseFloat(kline.Low, 64)
		c, _ := strconv.ParseFloat(kline.Close, 64)
		v, _ := strconv.ParseFloat(kline.Volume, 64)
		ot := time.Unix(int64(kline.StartTime/1000), 0).UTC()
		ct := time.Unix(int64(kline.EndTime/1000), 0).UTC()
		candle, err := q.Sync(lastCandle.Market, lastCandle.Symbol, lastCandle.Interval, o, h, l, c, v, ot, ct)
		if err != nil {
			return
		}
		update <- candle
	}
	errHandler := func(err error) {
		fmt.Println(err)
	}
	doneC, _, err = futures.WsKlineServe(lastCandle.Symbol, string(lastCandle.Interval), wsKlineHandler, errHandler)

	return
}

// Will sync quote with latest binance spot kline info.
func (q *Quote) SyncBinanceDelivery(update CandleChannel) (doneC chan struct{}, err error) {
	quote := *q
	if len(quote) == 0 {
		return nil, errors.New("won't be able to sync an empty quote")
	}

	lastCandle := quote[len(quote)-1]
	wsKlineHandler := func(event *delivery.WsKlineEvent) {
		kline := event.Kline
		o, _ := strconv.ParseFloat(kline.Open, 64)
		h, _ := strconv.ParseFloat(kline.High, 64)
		l, _ := strconv.ParseFloat(kline.Low, 64)
		c, _ := strconv.ParseFloat(kline.Close, 64)
		v, _ := strconv.ParseFloat(kline.Volume, 64)
		ot := time.Unix(int64(kline.StartTime/1000), 0).UTC()
		ct := time.Unix(int64(kline.EndTime/1000), 0).UTC()
		candle, err := q.Sync(lastCandle.Market, lastCandle.Symbol, lastCandle.Interval, o, h, l, c, v, ot, ct)
		if err != nil {
			return
		}
		update <- candle
	}
	errHandler := func(err error) {
		fmt.Println(err)
	}
	doneC, _, err = delivery.WsKlineServe(lastCandle.Symbol, string(lastCandle.Interval), wsKlineHandler, errHandler)

	return
}

func createCandleFromKline(market MarketType, symbol string, open, high, low, close, volume string, interval Interval, openTime, closeTime int64) (candle *Candle, err error) {
	ot := time.Unix(int64(openTime/1000), 0).UTC()
	ct := time.Unix(int64(closeTime/1000), 0).UTC()
	o, err := strconv.ParseFloat(open, 64)
	if err != nil {
		return
	}

	h, err := strconv.ParseFloat(high, 64)
	if err != nil {
		return
	}

	l, err := strconv.ParseFloat(low, 64)
	if err != nil {
		return
	}

	c, err := strconv.ParseFloat(close, 64)
	if err != nil {
		return
	}

	v, err := strconv.ParseFloat(volume, 64)
	if err != nil {
		return
	}

	candle, err = NewCandle(market, symbol, o, h, l, c, v, ot, ct, interval, nil, nil)
	return
}
