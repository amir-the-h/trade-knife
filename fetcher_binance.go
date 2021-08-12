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

// Fetches quote from binance market.
func NewQuoteFromBinance(apiKey, secretKey, symbol string, market MarketType, interval Interval, openTimestamp ...int64) (*Quote, error) {
	var (
		tots, tcts, ots, cts time.Time
		quote                = Quote{
			Symbol:   symbol,
			Market:   market,
			Interval: interval,
		}
	)

	direction := -1
	switch market {
	case MarketSpot:
		client := binance.NewClient(apiKey, secretKey)
		request := client.NewKlinesService().Symbol(symbol).Interval(string(interval))
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
				candle, err := createCandleFromKline(kline.Open, kline.High, kline.Low, kline.Close, kline.Volume, kline.OpenTime, kline.CloseTime)
				if err != nil {
					return &quote, err
				}
				quote.Candles = append(quote.Candles, candle)
			}
			klines, err = request.StartTime(ots.Unix() * 1000).EndTime(cts.Unix() * 1000).Do(context.Background())
			if err != nil {
				return &quote, err
			}
		}
	case MarketFutures:
		client := futures.NewClient(apiKey, secretKey)
		request := client.NewKlinesService().Symbol(symbol).Interval(string(interval))
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
				candle, err := createCandleFromKline(kline.Open, kline.High, kline.Low, kline.Close, kline.Volume, kline.OpenTime, kline.CloseTime)
				if err != nil {
					return &quote, err
				}
				quote.Candles = append(quote.Candles, candle)
			}
			klines, err = request.StartTime(ots.Unix() * 1000).EndTime(cts.Unix() * 1000).Do(context.Background())
			if err != nil {
				return &quote, err
			}
		}
	case MarketDelivery:
		client := delivery.NewClient(apiKey, secretKey)
		request := client.NewKlinesService().Symbol(symbol).Interval(string(interval))
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
				candle, err := createCandleFromKline(kline.Open, kline.High, kline.Low, kline.Close, kline.Volume, kline.OpenTime, kline.CloseTime)
				if err != nil {
					return &quote, err
				}
				quote.Candles = append(quote.Candles, candle)
			}
			klines, err = request.StartTime(ots.Unix() * 1000).EndTime(cts.Unix() * 1000).Do(context.Background())
			if err != nil {
				return &quote, err
			}
		}
	}

	q := &quote
	q.Sort()
	return q, nil
}

// Fetch all candles after last candle including itself.
func (q *Quote) RefreshBinance(apiKey, secretKey string) error {
	quote := *q
	if len(quote.Candles) == 0 {
		return errors.New("won't be able to refresh an empty quote")
	}

	var (
		lastCandle   = quote.Candles[len(quote.Candles)-1]
		openTime     = lastCandle.Opentime.Unix()
		fetchedQuote *Quote
		err          error
	)
	fetchedQuote, err = NewQuoteFromBinance(apiKey, secretKey, quote.Symbol, quote.Market, quote.Interval, openTime)
	if err != nil {
		return err
	}

	q.Merge(fetchedQuote)

	return nil
}

// Will sync quote with latest binance kline info.
func (q *Quote) SyncBinance(update CandleChannel) (doneC chan struct{}, err error) {
	quote := *q
	if len(quote.Candles) == 0 {
		return nil, errors.New("won't be able to sync an empty quote")
	}

	errHandler := func(err error) {
		fmt.Println(err)
	}

	switch quote.Market {
	case MarketSpot:
		wsKlineHandler := func(event *binance.WsKlineEvent) {
			kline := event.Kline
			o, _ := strconv.ParseFloat(kline.Open, 64)
			h, _ := strconv.ParseFloat(kline.High, 64)
			l, _ := strconv.ParseFloat(kline.Low, 64)
			c, _ := strconv.ParseFloat(kline.Close, 64)
			v, _ := strconv.ParseFloat(kline.Volume, 64)
			ot := time.Unix(int64(kline.StartTime/1000), 0).UTC()
			ct := time.Unix(int64(kline.EndTime/1000), 0).UTC()
			candle, err := q.Sync(o, h, l, c, v, ot, ct)
			if err != nil {
				return
			}
			update <- candle
		}
		doneC, _, err = binance.WsKlineServe(quote.Symbol, string(quote.Interval), wsKlineHandler, errHandler)
	case MarketFutures:
		wsKlineHandler := func(event *futures.WsKlineEvent) {
			kline := event.Kline
			o, _ := strconv.ParseFloat(kline.Open, 64)
			h, _ := strconv.ParseFloat(kline.High, 64)
			l, _ := strconv.ParseFloat(kline.Low, 64)
			c, _ := strconv.ParseFloat(kline.Close, 64)
			v, _ := strconv.ParseFloat(kline.Volume, 64)
			ot := time.Unix(int64(kline.StartTime/1000), 0).UTC()
			ct := time.Unix(int64(kline.EndTime/1000), 0).UTC()
			candle, err := q.Sync(o, h, l, c, v, ot, ct)
			if err != nil {
				return
			}
			update <- candle
		}
		doneC, _, err = futures.WsKlineServe(quote.Symbol, string(quote.Interval), wsKlineHandler, errHandler)
	case MarketDelivery:
		wsKlineHandler := func(event *delivery.WsKlineEvent) {
			kline := event.Kline
			o, _ := strconv.ParseFloat(kline.Open, 64)
			h, _ := strconv.ParseFloat(kline.High, 64)
			l, _ := strconv.ParseFloat(kline.Low, 64)
			c, _ := strconv.ParseFloat(kline.Close, 64)
			v, _ := strconv.ParseFloat(kline.Volume, 64)
			ot := time.Unix(int64(kline.StartTime/1000), 0).UTC()
			ct := time.Unix(int64(kline.EndTime/1000), 0).UTC()
			candle, err := q.Sync(o, h, l, c, v, ot, ct)
			if err != nil {
				return
			}
			update <- candle
		}
		doneC, _, err = delivery.WsKlineServe(quote.Symbol, string(quote.Interval), wsKlineHandler, errHandler)
	}

	return
}

func createCandleFromKline(open, high, low, close, volume string, openTime, closeTime int64) (candle *Candle, err error) {
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

	candle, err = NewCandle(o, h, l, c, v, ot, ct, nil, nil)
	return
}
