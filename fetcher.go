package trade_knife

import (
	"github.com/amir-the-h/goex"
	"time"
)

type Fetcher interface {
	NewQuote(currency goex.CurrencyPair, market MarketType, interval Interval, openTime *time.Time) (*Quote, error)
	Refresh(q *Quote) error
	Sync(q *Quote, update CandleChannel) (err error)
}
