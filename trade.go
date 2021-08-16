package trade_knife

import (
	"fmt"
	"time"
)

type Trade struct {
	Id                string
	Driver            string
	Symbol            string
	Base              string
	Quote             float64
	Amount            float64
	Entry             float64
	Exit              float64
	ProfitPrice       float64
	ProfitPercentage  float64
	StopLossPercent   float64
	StopLossPrice     float64
	TakeProfitPercent float64
	TakeProfitPrice   float64
	Position          PositionType
	Status            TradeStatus
	OpenCandle        *Candle
	CloseCandle       *Candle
	OpenAt            *time.Time
	CloseAt           *time.Time
}
type Trades []*Trade

// Search for a trade and it's index amoung all trades.
func (t Trades) Find(id string) (*Trade, int) {
	for i, trade := range t {
		if trade.Id == id {
			return trade, i
		}
	}

	return nil, 0
}

// Interface to determine how a trader should implements.
type Trader interface {
	Open(id, symbol, base string, position PositionType, quote, entry float64, sl, tp float64, openCandle *Candle) *Trade
	Close(id, exit float64, closeCandle *Candle)
	EntryWatcher()
	ExitWatcher()
	CloseWatcher()
}

// returns a pointer to a fresh trade.
func NewTrade(id, driver, symbol, base string, position PositionType, quote, entry float64, sl, tp float64, openCandle *Candle) *Trade {
	var takeProfitPercent, stopLossPercent float64
	now := time.Now().UTC()
	amount := quote / entry
	if id == "" {
		id = fmt.Sprint(now.Unix())
	}
	if position == PositionBuy {
		takeProfitPercent = (100 * tp / entry) - 100
		stopLossPercent = 100 - (100 * sl / entry)
	} else {
		takeProfitPercent = 100 - (100 * tp / entry)
		stopLossPercent = (100 * sl / entry) - 100
	}

	return &Trade{
		Id:                id,
		Driver:            driver,
		Symbol:            symbol,
		Base:              base,
		Quote:             quote,
		Amount:            amount,
		Entry:             entry,
		Position:          position,
		StopLossPercent:   stopLossPercent,
		StopLossPrice:     sl,
		TakeProfitPercent: takeProfitPercent,
		TakeProfitPrice:   tp,
		OpenAt:            &now,
		OpenCandle:        openCandle,
		Status:            TradeStatusOpen,
	}
}

// Close an active trade.
func (t *Trade) Close(price float64, candle *Candle) {
	if t.Status == TradeStatusClose {
		return
	}
	now := time.Now().UTC()
	t.CloseCandle = candle
	t.CloseAt = &now
	t.Exit = price
	t.Status = TradeStatusClose

	if t.Position == PositionBuy {
		t.ProfitPrice = (t.Exit - t.Entry) * t.Amount
	} else {
		t.ProfitPrice = (t.Entry - t.Exit) * t.Amount
	}
	t.ProfitPercentage = t.ProfitPrice * 100 / t.Quote
}

// Stringify the trade.
func (t Trade) String() string {
	var text string
	text = fmt.Sprintf("#%s\t%s\t%f\t%s\n", t.Id, t.Position, t.Amount, t.Symbol)
	text += fmt.Sprintf("Quote:\t%f %s\n", t.Quote, t.Base)
	text += fmt.Sprintf("Status:\t%s\n", t.Status)
	text += fmt.Sprintf("Entry:\t%.2f\t%s\n", t.Entry, t.OpenAt.Local().Format("06/02/01 15:04:05"))
	if t.Status == TradeStatusClose {
		text += fmt.Sprintf("Exit:\t%.2f\t%s\n", t.Exit, t.CloseAt.Local().Format("06/02/01 15:04:05"))
	}
	if t.TakeProfitPrice != 0 {
		text += fmt.Sprintf("TP:\t%.2f\t%%%.2f\n", t.TakeProfitPrice, t.TakeProfitPercent)
	}
	if t.StopLossPrice != 0 {
		text += fmt.Sprintf("SL:\t%.2f\t%%%.2f\n", t.StopLossPrice, t.StopLossPercent)
	}
	if t.Status == TradeStatusClose {
		text += fmt.Sprintf("\tResult:%.2f\t%%%.2f\n", t.ProfitPrice, t.ProfitPercentage)
	}

	return text
}
