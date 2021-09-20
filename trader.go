package trade_knife

import "github.com/amir-the-h/goex"

// Trader determines how a trader should be implemented.
type Trader interface {
	Open(currency goex.CurrencyPair, position PositionType, quote, entry float64, sl, tp float64, openCandle *Candle) *Trade
	Close(id, exit float64, closeCandle *Candle)
	Start() TradeError
	EntryWatcher()
	ExitWatcher()
	CloseWatcher()
}
