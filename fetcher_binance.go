package trade_knife

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/delivery"
	"github.com/adshao/go-binance/v2/futures"
)

// Fetches quote from binance spot market.
func NewQuoteFromBinanceSpot(apiKey, secretKey, symbol string, interval Interval) (*Quote, error) {
	var (
		quote                Quote
		tots, tcts, ots, cts time.Time
	)
	diff := interval.Duration()

	client := binance.NewClient(apiKey, secretKey)
	request := client.NewKlinesService().Symbol(symbol).Interval(string(interval))
	klines, err := request.Do(context.Background())
	if err != nil {
		return &quote, err
	}

	for len(klines) > 0 {
		tots = ots
		tcts = cts
		ots = time.Unix(int64(klines[0].OpenTime/1000), 0).Add(time.Hour * -24 * 90).UTC()
		cts = time.Unix(int64(klines[0].CloseTime/1000), 0).Add(diff).UTC()
		if tots == ots || tcts == cts {
			break
		}

		for _, kline := range klines {
			candle, err := createCandleFromKline(symbol, kline.Open, kline.High, kline.Low, kline.Close, kline.Volume, interval, kline.OpenTime, kline.CloseTime)
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

	return &quote, nil
}

// Fetches quote from binace futures market.
func NewQuoteFromBinanceFutures(apiKey, secretKey, symbol string, interval Interval) (*Quote, error) {
	var (
		quote                Quote
		tots, tcts, ots, cts time.Time
	)
	diff := interval.Duration()

	client := futures.NewClient(apiKey, secretKey)
	request := client.NewKlinesService().Symbol(symbol).Interval(string(interval))
	klines, err := request.Do(context.Background())
	if err != nil {
		return &quote, err
	}

	for len(klines) > 0 {
		tots = ots
		tcts = cts
		ots = time.Unix(int64(klines[0].OpenTime/1000), 0).Add(time.Hour * -24 * 90).UTC()
		cts = time.Unix(int64(klines[0].CloseTime/1000), 0).Add(diff).UTC()
		if tots == ots || tcts == cts {
			break
		}

		for _, kline := range klines {
			candle, err := createCandleFromKline(symbol, kline.Open, kline.High, kline.Low, kline.Close, kline.Volume, interval, kline.OpenTime, kline.CloseTime)
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

	return &quote, nil
}

// Fetches quote from binace delivery market.
func NewQuoteFromBinanceDelivery(apiKey, secretKey, symbol string, interval Interval) (*Quote, error) {
	var (
		quote                Quote
		tots, tcts, ots, cts time.Time
	)
	diff := interval.Duration()

	client := delivery.NewClient(apiKey, secretKey)
	request := client.NewKlinesService().Symbol(symbol).Interval(string(interval))
	klines, err := request.Do(context.Background())
	if err != nil {
		return &quote, err
	}

	for len(klines) > 0 {
		tots = ots
		tcts = cts
		ots = time.Unix(int64(klines[0].OpenTime/1000), 0).Add(time.Hour * -24 * 90).UTC()
		cts = time.Unix(int64(klines[0].CloseTime/1000), 0).Add(diff).UTC()
		if tots == ots || tcts == cts {
			break
		}

		for _, kline := range klines {
			candle, err := createCandleFromKline(symbol, kline.Open, kline.High, kline.Low, kline.Close, kline.Volume, interval, kline.OpenTime, kline.CloseTime)
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

	return &quote, nil
}

// Will sync quote with latest binance spot kline info.
func (q *Quote) SyncBinanceSpot(symbol string, interval Interval, update CandleChannel) (doneC chan struct{}, err error) {
	wsKlineHandler := func(event *binance.WsKlineEvent) {
		kline := event.Kline
		o, _ := strconv.ParseFloat(kline.Open, 64)
		h, _ := strconv.ParseFloat(kline.High, 64)
		l, _ := strconv.ParseFloat(kline.Low, 64)
		c, _ := strconv.ParseFloat(kline.Close, 64)
		v, _ := strconv.ParseFloat(kline.Volume, 64)
		ot := time.Unix(int64(kline.StartTime/1000), 0).UTC()
		ct := time.Unix(int64(kline.EndTime/1000), 0).UTC()
		candle, err := q.Sync(symbol, interval, o, h, l, c, v, ot, ct)
		if err != nil {
			return
		}
		update <- candle
	}
	errHandler := func(err error) {
		fmt.Println(err)
	}
	doneC, _, err = binance.WsKlineServe(symbol, string(interval), wsKlineHandler, errHandler)

	return
}

// Will sync quote with latest binance futures kline info.
func (q *Quote) SyncBinanceFutures(symbol string, interval Interval, update CandleChannel) (doneC chan struct{}, err error) {
	wsKlineHandler := func(event *futures.WsKlineEvent) {
		kline := event.Kline
		o, _ := strconv.ParseFloat(kline.Open, 64)
		h, _ := strconv.ParseFloat(kline.High, 64)
		l, _ := strconv.ParseFloat(kline.Low, 64)
		c, _ := strconv.ParseFloat(kline.Close, 64)
		v, _ := strconv.ParseFloat(kline.Volume, 64)
		ot := time.Unix(int64(kline.StartTime/1000), 0).UTC()
		ct := time.Unix(int64(kline.EndTime/1000), 0).UTC()
		candle, err := q.Sync(symbol, interval, o, h, l, c, v, ot, ct)
		if err != nil {
			return
		}
		update <- candle
	}
	errHandler := func(err error) {
		fmt.Println(err)
	}
	doneC, _, err = futures.WsKlineServe(symbol, string(interval), wsKlineHandler, errHandler)

	return
}

// Will sync quote with latest binance spot kline info.
func (q *Quote) SyncBinanceDelivery(symbol string, interval Interval, update CandleChannel) (doneC chan struct{}, err error) {
	wsKlineHandler := func(event *delivery.WsKlineEvent) {
		kline := event.Kline
		o, _ := strconv.ParseFloat(kline.Open, 64)
		h, _ := strconv.ParseFloat(kline.High, 64)
		l, _ := strconv.ParseFloat(kline.Low, 64)
		c, _ := strconv.ParseFloat(kline.Close, 64)
		v, _ := strconv.ParseFloat(kline.Volume, 64)
		ot := time.Unix(int64(kline.StartTime/1000), 0).UTC()
		ct := time.Unix(int64(kline.EndTime/1000), 0).UTC()
		candle, err := q.Sync(symbol, interval, o, h, l, c, v, ot, ct)
		if err != nil {
			return
		}
		update <- candle
	}
	errHandler := func(err error) {
		fmt.Println(err)
	}
	doneC, _, err = delivery.WsKlineServe(symbol, string(interval), wsKlineHandler, errHandler)

	return
}

func createCandleFromKline(symbol string, open, high, low, close, volume string, interval Interval, openTime, closeTime int64) (candle *Candle, err error) {
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

	candle, err = NewCandle(symbol, o, h, l, c, v, ot, ct, interval, nil, nil)
	return
}
