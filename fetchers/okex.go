package fetchers

import (
	"errors"
	"github.com/amir-the-h/goex"
	"github.com/amir-the-h/goex/builder"
	"github.com/amir-the-h/goex/okex"
	"github.com/amir-the-h/trade-knife"
	"time"
)

// Okex is an Okay-Exchange trade_knife.Fetcher
type Okex struct {
}

// NewOkex returns a pointer to a fresh Okex trade_knife.Fetcher.
func NewOkex() *Okex {
	return &Okex{}
}

// NewQuote fetches quote from okex market.
func (ok *Okex) NewQuote(currency goex.CurrencyPair, market trade_knife.MarketType, interval trade_knife.Interval, openTime *time.Time) (*trade_knife.Quote, error) {
	var (
		tOpenTime, tEndTime time.Time
		quote               = trade_knife.Quote{
			Currency: currency,
			Market:   market,
			Interval: interval,
		}
	)

	direction := -1
	size := 200
	optionalParameter := goex.OptionalParameter{}
	if openTime != nil {
		tOpenTime = *openTime
		tEndTime = tOpenTime.Add(interval.Duration() * time.Duration(size))
		direction = 1
		optionalParameter["start"] = openTime.Format(time.RFC3339)
		optionalParameter["end"] = tEndTime.Format(time.RFC3339)
	}
	switch market {
	case trade_knife.MarketSpot:
		client := builder.DefaultAPIBuilder.Build(goex.OKEX)
		klines, err := client.GetKlineRecords(currency, interval.KlinePeriod(), size, optionalParameter)
		if err != nil {
			return &quote, err
		}

		for len(klines) > 0 {
			lastTOpenTime := tOpenTime
			if direction == 1 {
				tOpenTime = time.Unix(klines[0].Timestamp, 0).Add(interval.Duration()).UTC()
			} else {
				tOpenTime = time.Unix(klines[len(klines)-1].Timestamp, 0).Add(interval.Duration() * time.Duration(direction*size)).UTC()
			}
			if tOpenTime.After(time.Now()) {
				break
			}
			if openTime != nil && openTime.After(tOpenTime) {
				break
			}
			tEndTime = tOpenTime.Add(interval.Duration() * time.Duration(size))
			if lastTOpenTime.Equal(tOpenTime) {
				break
			}

			for _, kline := range klines {
				candle, err := createCandleFromOkexKline(kline.Open, kline.High, kline.Low, kline.Close, kline.Vol, kline.Timestamp, interval)
				if err != nil {
					return &quote, err
				}
				quote.Candles = append(quote.Candles, candle)
			}
			optionalParameter["start"] = tOpenTime.Format(time.RFC3339)
			optionalParameter["end"] = tEndTime.Format(time.RFC3339)
			klines, err = client.GetKlineRecords(currency, interval.KlinePeriod(), size, optionalParameter)
			if err != nil {
				return &quote, err
			}
		}
	case trade_knife.MarketFutures:
		client := builder.DefaultAPIBuilder.BuildFuture(goex.OKEX_SWAP)
		klines, err := client.GetKlineRecords("", currency, interval.KlinePeriod(), size, optionalParameter)
		if err != nil {
			return &quote, err
		}

		for len(klines) > 0 {
			lastTOpenTime := tOpenTime
			if direction == 1 {
				tOpenTime = time.Unix(klines[0].Timestamp, 0).Add(interval.Duration()).UTC()
			} else {
				tOpenTime = time.Unix(klines[len(klines)-1].Timestamp, 0).Add(interval.Duration() * time.Duration(direction*size)).UTC()
			}
			if tOpenTime.After(time.Now()) {
				break
			}
			if openTime != nil && openTime.After(tOpenTime) {
				break
			}
			tEndTime = tOpenTime.Add(interval.Duration() * time.Duration(size))
			if lastTOpenTime.Equal(tOpenTime) {
				break
			}

			for _, kline := range klines {
				candle, err := createCandleFromOkexKline(kline.Open, kline.High, kline.Low, kline.Close, kline.Vol, kline.Timestamp, interval)
				if err != nil {
					return &quote, err
				}
				quote.Candles = append(quote.Candles, candle)
			}
			optionalParameter["start"] = tOpenTime.Format(time.RFC3339)
			optionalParameter["end"] = tEndTime.Format(time.RFC3339)
			klines, err = client.GetKlineRecords("", currency, interval.KlinePeriod(), size, optionalParameter)
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
func (ok *Okex) Refresh(q *trade_knife.Quote) error {
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
	fetchedQuote, err = ok.NewQuote(q.Currency, quote.Market, quote.Interval, &openTime)
	if err != nil {
		return err
	}

	q.Merge(fetchedQuote)

	return nil
}

// Sync syncs quote with latest okex kline info.
func (ok *Okex) Sync(q *trade_knife.Quote, update trade_knife.CandleChannel) (err error) {
	quote := *q
	if len(quote.Candles) == 0 {
		return errors.New("won't be able to sync an empty quote")
	}
	api := okex.NewOKExV3SwapWs(okex.NewOKEx(&goex.APIConfig{
		HttpClient: builder.DefaultAPIBuilder.GetHttpClient(),
	}))
	callback := func(k *goex.FutureKline, n int) {
		ot := time.Unix(k.Timestamp, 0).UTC()
		ct := ot.Add(q.Interval.Duration()).UTC()
		candle, err := q.Sync(k.Open, k.High, k.Low, k.Close, k.Vol, ot, ct)
		if err != nil {
			return
		}
		update <- candle
	}
	api.KlineCallback(callback)
	return api.SubscribeKline(quote.Currency, goex.SWAP_CONTRACT, int(q.Interval.KlinePeriod()))
}

func createCandleFromOkexKline(open, high, low, close, volume float64, timestamp int64, interval trade_knife.Interval) (*trade_knife.Candle, error) {
	ot := time.Unix(timestamp, 0).UTC()
	ct := ot.Add(interval.Duration())
	return trade_knife.NewCandle(open, high, low, close, volume, ot, ct, nil, nil)
}
