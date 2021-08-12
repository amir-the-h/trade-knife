package trade_knife

import (
	"log"
	"strings"
)

// PaperTrader exchange papertrade driver.
type PaperTrader struct {
	Trades           Trades
	candleChannel    CandleChannel
	enterChannel     EnterChannel
	exitChannel      ExitChannel
	BuyScoreTrigger  float64
	SellScoreTrigger float64
	CloseOnOpposite  bool
	Cross            bool
	Debug            bool
	activeTrade      *Trade
}

// Returns a pointer to fresh binance papertrade driver.
func NewPaperTrader(candleChannel CandleChannel, entryChannel EnterChannel, buyscoreTrigger, sellscoreTrigger float64, closeOnOpposite, cross bool) *PaperTrader {
	exitChannel := make(ExitChannel)
	return &PaperTrader{
		candleChannel:    candleChannel,
		enterChannel:     entryChannel,
		exitChannel:      exitChannel,
		BuyScoreTrigger:  buyscoreTrigger,
		SellScoreTrigger: sellscoreTrigger,
		CloseOnOpposite:  closeOnOpposite,
		Cross:            cross,
	}
}

// Launch all watchers of the driver.
func (pt *PaperTrader) Start() TradeError {
	if pt.Debug {
		log.Println("⏳", "Starting paper trade")
	}
	// setup watchers threads
	go pt.EntryWatcher()
	go pt.ExitWatcher()
	go pt.CloseWatcher()

	if pt.Debug {
		log.Println("✅", "Paper trade started")
	}
	return nil
}

// Create a new trade immediately.
func (pt *PaperTrader) Open(id, symbol, base string, quote, entry float64, position PositionType, sl, tp float64, openCandle *Candle) *Trade {
	trade := NewTrade(id, "PaperTrader Papertrade", symbol, base, quote, entry, position, sl, tp, openCandle)
	pt.Trades = append(pt.Trades, trade)
	return trade
}

// Close an open trade immediately.
func (pt *PaperTrader) Close(id, exit float64, closeCandle *Candle) {

}

// Watch for entry signals and open proper positions.
func (pt *PaperTrader) EntryWatcher() {
	if pt.Debug {
		log.Println("✅", "Entry watcher started")
	}
	for enterSignal := range pt.enterChannel {
		var position PositionType
		if enterSignal.Candle.Score >= pt.BuyScoreTrigger {
			position = PositionBuy
		} else if enterSignal.Candle.Score <= pt.SellScoreTrigger {
			position = PositionSell
		} else {
			continue
		}

		// check for crossed positions
		if pt.activeTrade != nil {
			if pt.activeTrade.Position != position {
				if pt.Cross || pt.CloseOnOpposite {
					// fire exit signal
					trade := *pt.activeTrade
					pt.exitChannel <- ExitSignal{
						Trade:  &trade,
						Candle: enterSignal.Candle,
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
		trade := pt.Open("", symbol, base, enterSignal.Quote, enterSignal.Candle.Close, position, enterSignal.Stoploss, enterSignal.TakeProfit, enterSignal.Candle)
		pt.activeTrade = trade
		if pt.Debug {
			log.Printf("💰 Trade started by score %f casued %s\n%s", enterSignal.Candle.Score, enterSignal.Cause, *trade)
		}
	}
}

// Watch for exit signals and and fire proper close signals.
func (pt *PaperTrader) ExitWatcher() {
	if pt.Debug {
		log.Println("✅", "Exit watcher started")
	}
	for candle := range pt.candleChannel {
		if pt.activeTrade != nil && pt.activeTrade.Position == PositionBuy {
			// check for stop loss first
			if candle.Close <= pt.activeTrade.StopLossPercent {
				pt.exitChannel <- ExitSignal{
					Trade:  pt.activeTrade,
					Candle: candle,
					Cause:  ExitCauseStopLossTriggered,
				}
				continue
			}
			// and take profit as well
			if candle.Close >= pt.activeTrade.TakeProfitPrice {
				pt.exitChannel <- ExitSignal{
					Trade:  pt.activeTrade,
					Candle: candle,
					Cause:  ExitCauseTakeProfitTriggered,
				}
				continue
			}
		} else if pt.activeTrade != nil && pt.activeTrade.Position == PositionSell {
			// same rules here
			if candle.Close >= pt.activeTrade.StopLossPrice {
				pt.exitChannel <- ExitSignal{
					Trade:  pt.activeTrade,
					Candle: candle,
					Cause:  ExitCauseStopLossTriggered,
				}
				continue
			}
			if candle.Close <= pt.activeTrade.TakeProfitPrice {
				pt.exitChannel <- ExitSignal{
					Trade:  pt.activeTrade,
					Candle: candle,
					Cause:  ExitCauseTakeProfitTriggered,
				}
				continue
			}
		}
	}
}

// Watch for close signals and close the trade immediately.
func (pt *PaperTrader) CloseWatcher() {
	if pt.Debug {
		log.Println("✅", "Close watcher started")
	}
	for exitSignal := range pt.exitChannel {
		// close the trade
		exitSignal.Trade.Close(exitSignal.Candle.Close, exitSignal.Candle)
		// and remove it from active trades
		if pt.Debug {
			icon := "📈"
			if exitSignal.Trade.ProfitPrice < 0 {
				icon = "📉"
			}
			log.Printf("%s Trade finished by %s\n%s", icon, exitSignal.Cause, *exitSignal.Trade)
		}
		if pt.activeTrade != nil && pt.activeTrade.Id == exitSignal.Trade.Id {
			pt.activeTrade = nil
		}
	}
}
