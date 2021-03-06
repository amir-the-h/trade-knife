package trade_knife

import (
	"fmt"
	"github.com/amir-the-h/goex"
	"time"
)

// Trade represents a real world trade
type Trade struct {
	Currency          goex.CurrencyPair
	Id                string
	Driver            string
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

// Find searches for a trade, and its index among all trades.
func (t Trades) Find(id string) (*Trade, int) {
	for i, trade := range t {
		if trade.Id == id {
			return trade, i
		}
	}

	return nil, 0
}

// NewTrade returns a pointer to a fresh trade.
func NewTrade(id, driver string, currency goex.CurrencyPair, position PositionType, quote, entry float64, sl, tp float64, openCandle *Candle) *Trade {
	var takeProfitPercent, stopLossPercent float64
	now := time.Now().UTC()
	amount := quote / entry
	if tp != 0 || sl != 0 {
		if position == PositionBuy {
			if tp != 0 {
				takeProfitPercent = (tp - entry) / entry * 100
			}
			if sl != 0 {
				stopLossPercent = (entry - sl) / entry * 100
			}
		} else {
			if tp != 0 {
				takeProfitPercent = (entry - tp) / entry * 100
			}
			if sl != 0 {
				stopLossPercent = (sl - entry) / entry * 100
			}
		}
	}

	return &Trade{
		Id:                id,
		Driver:            driver,
		Currency:          currency,
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

// Close closes an active trade.
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

// String Stringify the trade.
func (t Trade) String() string {
	var text string
	text = fmt.Sprintf("#%s\t%s\t%f %s\n", t.Id, t.Position, t.Amount, t.Currency.CurrencyA)
	text += fmt.Sprintf("Quote:\t%f %s\n", t.Quote, t.Currency.CurrencyB)
	text += fmt.Sprintf("Status:\t%s\n", t.Status)
	text += fmt.Sprintf("Entry:\t%.4f\t%s\n", t.Entry, t.OpenAt.Local().Format("06/02/01 15:04:05"))
	if t.Status == TradeStatusClose {
		text += fmt.Sprintf("Exit:\t%.4f\t%s\n", t.Exit, t.CloseAt.Local().Format("06/02/01 15:04:05"))
	}
	if t.TakeProfitPrice != 0 {
		text += fmt.Sprintf("TP:\t%.4f\t%%%.4f\n", t.TakeProfitPrice, t.TakeProfitPercent)
	}
	if t.StopLossPrice != 0 {
		text += fmt.Sprintf("SL:\t%.4f\t%%%.4f\n", t.StopLossPrice, t.StopLossPercent)
	}
	if t.Status == TradeStatusClose {
		text += fmt.Sprintf("\tResult:%.4f\t%%%.4f\n", t.ProfitPrice, t.ProfitPercentage)
	}

	return text
}
