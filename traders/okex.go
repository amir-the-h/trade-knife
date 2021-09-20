package traders

import (
	"github.com/amir-the-h/goex"
	"github.com/amir-the-h/goex/builder"
	"github.com/amir-the-h/goex/okex"
	"github.com/amir-the-h/trade-knife"
)

// Okex is an Okay-Exchange trade_knife.Trader
type Okex struct {
	Trades           trade_knife.Trades
	ActiveTrade      *trade_knife.Trade
	Lever            float64
	BuyScoreTrigger  float64
	SellScoreTrigger float64
	CloseOnOpposite  bool
	Cross            bool
	Debug            bool
	Market           trade_knife.MarketType
	candleChannel    trade_knife.CandleChannel
	entryChannel     trade_knife.EnterChannel
	exitChannel      trade_knife.ExitChannel
	openTrades       trade_knife.TradesChannel
	doneTrades       trade_knife.TradesChannel
	logger           *trade_knife.Logger
	Api              *okex.OKExSwap
	apiKey           string
	secretKey        string
	passphrase       string
}

// NewOkex returns a pointer to a fresh Okex's trade_knife.Trader.
func NewOkex(apiKey, secretKey, passphrase string, candleChannel trade_knife.CandleChannel, entryChannel trade_knife.EnterChannel, exitChannel trade_knife.ExitChannel, openTrades trade_knife.TradesChannel, doneTrades trade_knife.TradesChannel, lever, buyScoreTrigger, sellScoreTrigger float64, closeOnOpposite, cross, debug bool, logger *trade_knife.Logger) *Okex {
	defaultBuilder := builder.DefaultAPIBuilder.APIKey(apiKey).APISecretkey(secretKey).ApiPassphrase(passphrase)
	apiConfig := &goex.APIConfig{
		HttpClient:    defaultBuilder.GetHttpClient(),
		ApiKey:        apiKey,
		ApiSecretKey:  secretKey,
		ApiPassphrase: passphrase,
		Lever:         lever,
	}
	api := okex.NewOKEx(apiConfig).OKExSwap
	return &Okex{
		candleChannel:    candleChannel,
		entryChannel:     entryChannel,
		exitChannel:      exitChannel,
		openTrades:       openTrades,
		doneTrades:       doneTrades,
		BuyScoreTrigger:  buyScoreTrigger,
		SellScoreTrigger: sellScoreTrigger,
		CloseOnOpposite:  closeOnOpposite,
		Cross:            cross,
		Debug:            debug,
		Lever:            lever,
		logger:           logger,
		apiKey:           apiKey,
		secretKey:        secretKey,
		passphrase:       passphrase,
		Api:              api,
	}
}

func (ok *Okex) Open(id string, currency goex.CurrencyPair, position trade_knife.PositionType, quote, entry float64, sl, tp float64, openCandle *trade_knife.Candle) *trade_knife.Trade {
	trade := trade_knife.NewTrade(id, "Okex", currency, position, quote, entry, sl, tp, openCandle)
	ok.openTrades <- trade
	ok.Trades = append(ok.Trades, trade)
	return trade
}

func (ok *Okex) Close(id, exit float64, closeCandle *trade_knife.Candle) {
	panic("implement me")
}

func (ok *Okex) Start() trade_knife.TradeError {
	panic("implement me")
}

func (ok *Okex) EntryWatcher() {
	panic("implement me")
}

func (ok *Okex) ExitWatcher() {
	panic("implement me")
}

func (ok *Okex) CloseWatcher() {
	panic("implement me")
}
