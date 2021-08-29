package trade_knife

import (
	"strings"
)

// PaperTrader exchange papertrade driver.
type PaperTrader struct {
	Trades           Trades
	candleChannel    CandleChannel
	entryChannel     EnterChannel
	exitChannel      ExitChannel
	tradeChannel     TradesChannel
	Wallet           float64
	BuyScoreTrigger  float64
	SellScoreTrigger float64
	CloseOnOpposite  bool
	Cross            bool
	Debug            bool
	ActiveTrade      *Trade
	logger           Logger
}

// Returns a pointer to fresh binance papertrade driver.
func NewPaperTrader(candleChannel CandleChannel, entryChannel EnterChannel, exitChannel ExitChannel, tradeChannel TradesChannel, wallet, buyscoreTrigger, sellscoreTrigger float64, closeOnOpposite, cross, debug bool, logger Logger) *PaperTrader {
	return &PaperTrader{
		candleChannel:    candleChannel,
		entryChannel:     entryChannel,
		exitChannel:      exitChannel,
		tradeChannel:     tradeChannel,
		Wallet:           wallet,
		BuyScoreTrigger:  buyscoreTrigger,
		SellScoreTrigger: sellscoreTrigger,
		CloseOnOpposite:  closeOnOpposite,
		Cross:            cross,
		Debug:            debug,
		logger:           logger,
	}
}

// Launch all watchers of the driver.
func (pt *PaperTrader) Start() TradeError {
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

// Create a new trade immediately.
func (pt *PaperTrader) Open(id, symbol, base string, position PositionType, quote, entry, sl, tp float64, openCandle *Candle) *Trade {
	trade := NewTrade(id, "PaperTrader Papertrade", symbol, base, position, quote, entry, sl, tp, openCandle)
	pt.Trades = append(pt.Trades, trade)
	return trade
}

func (pt *PaperTrader) Close(id string, exit float64, closeCandle *Candle) {
	for _, trade := range pt.Trades {
		if trade.Id == id {
			trade.Close(exit, closeCandle)
			pt.Wallet += trade.ProfitPrice
		}
	}
}

// Watch for entry signals and open proper positions.
func (pt *PaperTrader) EntryWatcher() {
	if pt.Debug {
		pt.logger.Success.Println("Entry watcher started")
	}
	for enterSignal := range pt.entryChannel {
		var position PositionType
		if enterSignal.Score >= pt.BuyScoreTrigger {
			position = PositionBuy
		} else if enterSignal.Score <= pt.SellScoreTrigger {
			position = PositionSell
		} else {
			continue
		}

		// check for crossed positions
		if pt.ActiveTrade != nil {
			if pt.ActiveTrade.Position != position {
				if pt.CloseOnOpposite {
					// fire exit signal
					pt.exitChannel <- ExitSignal{
						Trade:  pt.ActiveTrade,
						Candle: &enterSignal.Candle,
						Cause:  ExitCause("Crossed position"),
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

		base := "USDT"
		symbol := strings.ReplaceAll(enterSignal.Symbol, base, "")
		trade := pt.Open("", symbol, base, position, pt.Wallet, enterSignal.Candle.Close, enterSignal.Stoploss, enterSignal.TakeProfit, &enterSignal.Candle)
		pt.ActiveTrade = trade
		if pt.Debug {
			pt.logger.Info.Printf("Trade started by score %f casued %s\n%s", enterSignal.Candle.Score, enterSignal.Cause, *trade)
		}
		pt.tradeChannel <- trade
	}
}

// Watch for exit signals and and fire proper close signals.
func (pt *PaperTrader) ExitWatcher() {
	if pt.Debug {
		pt.logger.Success.Println("Exit watcher started")
	}
	for candle := range pt.candleChannel {
		if pt.ActiveTrade != nil && (pt.ActiveTrade.StopLossPercent != 0 || pt.ActiveTrade.TakeProfitPercent != 0) {
			if pt.ActiveTrade.Position == PositionBuy {
				// check for stop loss first
				if pt.ActiveTrade.StopLossPercent != 0 && candle.Close <= pt.ActiveTrade.StopLossPercent {
					pt.exitChannel <- ExitSignal{
						Trade:  pt.ActiveTrade,
						Candle: candle,
						Cause:  ExitCauseStopLossTriggered,
					}
					continue
				}
				// and take profit as well
				if pt.ActiveTrade.TakeProfitPercent != 0 && candle.Close >= pt.ActiveTrade.TakeProfitPrice {
					pt.exitChannel <- ExitSignal{
						Trade:  pt.ActiveTrade,
						Candle: candle,
						Cause:  ExitCauseTakeProfitTriggered,
					}
					continue
				}
			} else {
				// same rules here
				if pt.ActiveTrade.StopLossPercent != 0 && candle.Close >= pt.ActiveTrade.StopLossPrice {
					pt.exitChannel <- ExitSignal{
						Trade:  pt.ActiveTrade,
						Candle: candle,
						Cause:  ExitCauseStopLossTriggered,
					}
					continue
				}
				if pt.ActiveTrade.TakeProfitPercent != 0 && candle.Close <= pt.ActiveTrade.TakeProfitPrice {
					pt.exitChannel <- ExitSignal{
						Trade:  pt.ActiveTrade,
						Candle: candle,
						Cause:  ExitCauseTakeProfitTriggered,
					}
					continue
				}
			}
		}
	}
}

// Watch for close signals and close the trade immediately.
func (pt *PaperTrader) CloseWatcher() {
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
