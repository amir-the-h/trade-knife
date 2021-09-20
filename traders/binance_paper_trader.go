package traders

import (
	"fmt"
	"github.com/amir-the-h/goex"
	"github.com/amir-the-h/trade-knife"
	"time"
)

// BinancePaperTrader is a Binance-Exchange-PaperTrader trade_knife.Trader
type BinancePaperTrader struct {
	Trades           trade_knife.Trades
	candleChannel    trade_knife.CandleChannel
	entryChannel     trade_knife.EnterChannel
	exitChannel      trade_knife.ExitChannel
	openTrades       trade_knife.TradesChannel
	doneTrades       trade_knife.TradesChannel
	Wallet           float64
	BuyScoreTrigger  float64
	SellScoreTrigger float64
	CloseOnOpposite  bool
	Cross            bool
	Debug            bool
	ActiveTrade      *trade_knife.Trade
	logger           trade_knife.Logger
}

// NewBinancePaperTrader returns a pointer to fresh BinancePaperTrader's trade_knife.Trader.
func NewBinancePaperTrader(candleChannel trade_knife.CandleChannel, entryChannel trade_knife.EnterChannel, exitChannel trade_knife.ExitChannel, openTrades trade_knife.TradesChannel, doneTrades trade_knife.TradesChannel, wallet, buyScoreTrigger, sellScoreTrigger float64, closeOnOpposite, cross, debug bool, logger trade_knife.Logger) *BinancePaperTrader {
	return &BinancePaperTrader{
		candleChannel:    candleChannel,
		entryChannel:     entryChannel,
		exitChannel:      exitChannel,
		openTrades:       openTrades,
		doneTrades:       doneTrades,
		Wallet:           wallet,
		BuyScoreTrigger:  buyScoreTrigger,
		SellScoreTrigger: sellScoreTrigger,
		CloseOnOpposite:  closeOnOpposite,
		Cross:            cross,
		Debug:            debug,
		logger:           logger,
	}
}

// Open creates a new trade immediately.
func (pt *BinancePaperTrader) Open(currency goex.CurrencyPair, position trade_knife.PositionType, quote, entry, sl, tp float64, openCandle *trade_knife.Candle) *trade_knife.Trade {
	id := fmt.Sprint(time.Now().Unix())
	trade := trade_knife.NewTrade(id, "BinancePaperTrader Papertrade", currency, position, quote, entry, sl, tp, openCandle)
	pt.openTrades <- trade
	pt.Trades = append(pt.Trades, trade)
	return trade
}

// Close closes the chosen trade
func (pt *BinancePaperTrader) Close(id string, exit float64, closeCandle *trade_knife.Candle) {
	for _, trade := range pt.Trades {
		if trade.Id == id {
			trade.Close(exit, closeCandle)
			pt.Wallet += trade.ProfitPrice
			pt.doneTrades <- trade
		}
	}
}

// Start launches all watchers of the driver.
func (pt *BinancePaperTrader) Start() trade_knife.TradeError {
	if pt.Debug {
		pt.logger.Info.Println("Starting paper trade")
	}
	// setup watchers threads
	go pt.EntryWatcher()
	go pt.ExitWatcher()
	go pt.CloseWatcher()

	if pt.Debug {
		pt.logger.Success.Println("Paper trade started")
	}
	return nil
}

// EntryWatcher watches for entry signals and open proper positions.
func (pt *BinancePaperTrader) EntryWatcher() {
	if pt.Debug {
		pt.logger.Success.Println("Entry watcher started")
	}
	for enterSignal := range pt.entryChannel {
		var position trade_knife.PositionType
		if enterSignal.Score >= pt.BuyScoreTrigger {
			position = trade_knife.PositionBuy
		} else if enterSignal.Score <= pt.SellScoreTrigger {
			position = trade_knife.PositionSell
		} else {
			continue
		}

		// check for crossed positions
		if pt.ActiveTrade != nil {
			if pt.ActiveTrade.Position != position {
				if pt.CloseOnOpposite {
					// fire exit signal
					pt.exitChannel <- trade_knife.ExitSignal{
						Trade:  pt.ActiveTrade,
						Candle: &enterSignal.Candle,
						Cause:  "Crossed position",
					}

					// shall we pass?
					if !pt.Cross {
						continue
					}
				}
			} else {
				continue
			}
		}

		trade := pt.Open(enterSignal.Currency, position, pt.Wallet, enterSignal.Candle.Close, enterSignal.Stoploss, enterSignal.TakeProfit, &enterSignal.Candle)
		pt.ActiveTrade = trade
		if pt.Debug {
			pt.logger.Info.Printf("Trade started by score %f casued %s\n%s", enterSignal.Candle.Score, enterSignal.Cause, *trade)
		}
	}
}

// ExitWatcher watches for exit signals and fire proper close signals.
func (pt *BinancePaperTrader) ExitWatcher() {
	if pt.Debug {
		pt.logger.Success.Println("Exit watcher started")
	}
	for candle := range pt.candleChannel {
		if pt.ActiveTrade != nil && (pt.ActiveTrade.StopLossPercent != 0 || pt.ActiveTrade.TakeProfitPercent != 0) {
			if pt.ActiveTrade.Position == trade_knife.PositionBuy {
				// check for stop loss first
				if pt.ActiveTrade.StopLossPercent != 0 && candle.Close <= pt.ActiveTrade.StopLossPercent {
					pt.exitChannel <- trade_knife.ExitSignal{
						Trade:  pt.ActiveTrade,
						Candle: candle,
						Cause:  trade_knife.ExitCauseStopLossTriggered,
					}
					continue
				}
				// and take profit as well
				if pt.ActiveTrade.TakeProfitPercent != 0 && candle.Close >= pt.ActiveTrade.TakeProfitPrice {
					pt.exitChannel <- trade_knife.ExitSignal{
						Trade:  pt.ActiveTrade,
						Candle: candle,
						Cause:  trade_knife.ExitCauseTakeProfitTriggered,
					}
					continue
				}
			} else {
				// same rules here
				if pt.ActiveTrade.StopLossPercent != 0 && candle.Close >= pt.ActiveTrade.StopLossPrice {
					pt.exitChannel <- trade_knife.ExitSignal{
						Trade:  pt.ActiveTrade,
						Candle: candle,
						Cause:  trade_knife.ExitCauseStopLossTriggered,
					}
					continue
				}
				if pt.ActiveTrade.TakeProfitPercent != 0 && candle.Close <= pt.ActiveTrade.TakeProfitPrice {
					pt.exitChannel <- trade_knife.ExitSignal{
						Trade:  pt.ActiveTrade,
						Candle: candle,
						Cause:  trade_knife.ExitCauseTakeProfitTriggered,
					}
					continue
				}
			}
		}
	}
}

// CloseWatcher watches for close signals and close the trade immediately.
func (pt *BinancePaperTrader) CloseWatcher() {
	if pt.Debug {
		pt.logger.Success.Println("Close watcher started")
	}
	for exitSignal := range pt.exitChannel {
		// close the trade
		pt.Close(exitSignal.Trade.Id, exitSignal.Candle.Close, exitSignal.Candle)
		// and remove it from active trades
		if pt.Debug {
			icon := "📈"
			if exitSignal.Trade.ProfitPrice < 0 {
				icon = "📉"
			}
			pt.logger.Info.Printf("%s Trade finished by %s\n%s", icon, exitSignal.Cause, *exitSignal.Trade)
		}
		if pt.ActiveTrade != nil && pt.ActiveTrade.Id == exitSignal.Trade.Id {
			pt.ActiveTrade = nil
		}
	}
}
