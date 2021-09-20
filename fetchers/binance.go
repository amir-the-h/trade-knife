package fetchers

import (
	"context"
	"errors"
	"fmt"
	"github.com/amir-the-h/goex"
	"github.com/amir-the-h/trade-knife"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/delivery"
	"github.com/adshao/go-binance/v2/futures"
)

// Binance is a Binance-Exchange trade_knife.Fetcher
type Binance struct {
}

// NewBinance returns a pointer to a fresh Binance trader.
func NewBinance() *Binance {
	return &Binance{}
}

// NewQuote fetches quote from binance market.
func (b *Binance) NewQuote(currency goex.CurrencyPair, market trade_knife.MarketType, interval trade_knife.Interval, openTime *time.Time) (*trade_knife.Quote, error) {
	var (
		tots, tcts, ots, cts time.Time
		quote                = trade_knife.Quote{
			Currency: currency,
			Market:   market,
			Interval: interval,
		}
	)

	direction := -1
	switch market {
	case trade_knife.MarketSpot:
		client := binance.NewClient("", "")
		request := client.NewKlinesService().Symbol(currency.ToSymbol("")).Interval(string(interval))
		if openTime != nil {
			direction = 1
			request.StartTime(openTime.Unix() * 1000)
		}
		klines, err := request.Do(context.Background())
		if err != nil {
			return &quote, err
		}

		for len(klines) > 0 {
			tots = ots
			tcts = cts
			ots = time.Unix(klines[0].OpenTime/1000, 0).Add(time.Hour * time.Duration(direction) * 24 * 90).UTC()
			cts = time.Unix(klines[0].CloseTime/1000, 0).Add(interval.Duration()).UTC()
			if tots == ots || tcts == cts {
				break
			}

			for _, kline := range klines {
				candle, err := createCandleFromBinanceKline(kline.Open, kline.High, kline.Low, kline.Close, kline.Volume, kline.OpenTime, kline.CloseTime)
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
	case trade_knife.MarketFutures:
		client := futures.NewClient("","")
		request := client.NewKlinesService().Symbol(currency.ToSymbol("")).Interval(string(interval))
		if openTime != nil {
			direction = 1
			request.StartTime(openTime.Unix() * 1000)
		}
		klines, err := request.Do(context.Background())
		if err != nil {
			return &quote, err
		}

		for len(klines) > 0 {
			tots = ots
			tcts = cts
			ots = time.Unix(klines[0].OpenTime/1000, 0).Add(time.Hour * time.Duration(direction) * 24 * 90).UTC()
			cts = time.Unix(klines[0].CloseTime/1000, 0).Add(interval.Duration()).UTC()
			if tots == ots || tcts == cts {
				break
			}

			for _, kline := range klines {
				candle, err := createCandleFromBinanceKline(kline.Open, kline.High, kline.Low, kline.Close, kline.Volume, kline.OpenTime, kline.CloseTime)
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
	case trade_knife.MarketDelivery:
		client := delivery.NewClient("","")
		request := client.NewKlinesService().Symbol(currency.ToSymbol("")).Interval(string(interval))
		if openTime != nil {
			direction = 1
			request.StartTime(openTime.Unix() * 1000)
		}
		klines, err := request.Do(context.Background())
		if err != nil {
			return &quote, err
		}

		for len(klines) > 0 {
			tots = ots
			tcts = cts
			ots = time.Unix(klines[0].OpenTime/1000, 0).Add(time.Hour * time.Duration(direction) * 24 * 90).UTC()
			cts = time.Unix(klines[0].CloseTime/1000, 0).Add(interval.Duration()).UTC()
			if tots == ots || tcts == cts {
				break
			}

			for _, kline := range klines {
				candle, err := createCandleFromBinanceKline(kline.Open, kline.High, kline.Low, kline.Close, kline.Volume, kline.OpenTime, kline.CloseTime)
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

// Refresh fetches all candles after last candle including itself.
func (b *Binance) Refresh(q *trade_knife.Quote) error {
	quote := *q
	if len(quote.Candles) == 0 {
		return errors.New("won't be able to refresh an empty quote")
	}

	var (
		lastCandle   = quote.Candles[len(quote.Candles)-1]
		openTime     = lastCandle.Opentime
		fetchedQuote *trade_knife.Quote
		err          error
	)
	fetchedQuote, err = b.NewQuote(quote.Currency, quote.Market, quote.Interval, &openTime)
	if err != nil {
		return err
	}

	q.Merge(fetchedQuote)

	return nil
}

// Sync syncs quote with latest binance kline info.
func (b *Binance) Sync(q *trade_knife.Quote, update trade_knife.CandleChannel) (err error) {
	quote := *q
	if len(quote.Candles) == 0 {
		return errors.New("won't be able to sync an empty quote")
	}

	errHandler := func(err error) {
		fmt.Println(err)
	}

	switch quote.Market {
	case trade_knife.MarketSpot:
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
		_, _, err = binance.WsKlineServe(quote.Currency.ToSymbol(""), string(quote.Interval), wsKlineHandler, errHandler)
	case trade_knife.MarketFutures:
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
		_, _, err = futures.WsKlineServe(quote.Currency.ToSymbol(""), string(quote.Interval), wsKlineHandler, errHandler)
	case trade_knife.MarketDelivery:
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
		_, _, err = delivery.WsKlineServe(quote.Currency.ToSymbol(""), string(quote.Interval), wsKlineHandler, errHandler)
	}

	return
}

func createCandleFromBinanceKline(open, high, low, close, volume string, openTime, closeTime int64) (candle *trade_knife.Candle, err error) {
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

	candle, err = trade_knife.NewCandle(o, h, l, c, v, ot, ct, nil, nil)
	return
}
