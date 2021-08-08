// Package trade_knife will provide fundamental concept and utils to
// operate with financial time-series data.
package trade_knife

import "errors"

type EnterSignal struct {
	Score      float64
	Cause      EnterCause
	Candle     *Candle
	Quote      float64
	TakeProfit float64
	Stoploss   float64
}
type ExitSignal struct {
	Trade  *Trade
	Candle *Candle
	Cause  ExitCause
}

type PositionType string
type TradeStatus string
type EnterCause string
type ExitCause string

type TradesChannel chan *Trade
type EnterChannel chan EnterSignal
type ExitChannel chan ExitSignal
type CandleChannel chan *Candle

type CandleError error
type SourceError error
type TradeError error

const (
	// Intervals
	Interval1m  = Interval("1m")
	Interval3m  = Interval("3m")
	Interval5m  = Interval("5m")
	Interval15m = Interval("15m")
	Interval30m = Interval("30m")
	Interval1h  = Interval("1h")
	Interval2h  = Interval("2h")
	Interval4h  = Interval("4h")
	Interval6h  = Interval("6h")
	Interval8h  = Interval("8h")
	Interval12h = Interval("12h")
	Interval1d  = Interval("1d")
	Interval3d  = Interval("3d")
	Interval1w  = Interval("1w")
	Interval1M  = Interval("1M")

	// single sources
	SourceOpen   = Source("open")
	SourceHigh   = Source("high")
	SourceLow    = Source("low")
	SourceClose  = Source("close")
	SourceVolume = Source("volume")

	// double sources
	SourceOpenHigh  = Source("oh2")
	SourceOpenLow   = Source("ol2")
	SourceOpenClose = Source("oc2")
	SourceHighLow   = Source("hl2")
	SourceHighClose = Source("hc2")
	SourceLowClose  = Source("lc2")

	// triple sources
	SourceOpenHighLow   = Source("ohl3")
	SourceOpenHighClose = Source("ohc3")
	SourceOpenLowClose  = Source("olc3")
	SourceHighLowClose  = Source("hlc3")

	// all together
	SourceOpenHighLowClose = Source("ohlc4")

	// Position types
	PositionBuy  = PositionType("Buy")
	PositionSell = PositionType("Sell")

	// Trade statuses
	TradeStatusOpen  = TradeStatus("Open")
	TradeStatusClose = TradeStatus("Close")

	// Exit signals
	ExitCauseStopLossTriggered   = ExitCause("Stop loss")
	ExitCauseTakeProfitTriggered = ExitCause("Take profit")
	ExitCauseMarket              = ExitCause("Market")
)

var (
	ErrInvalidCandleData = errors.New("invalid data provided for candle").(CandleError)
	ErrNotEnoughCandles  = errors.New("not enough candles to operate").(CandleError)
	ErrInvalidSource     = errors.New("invalid source provided").(SourceError)
)
