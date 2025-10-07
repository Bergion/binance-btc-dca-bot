package dca

import (
	"log/slog"

	"github.com/Bergion/binance-btc-dca-bot/pkg/binance"
)

type Executor struct {
	config        Config
	binanceClient *binance.Client
}

func NewExecutor(cfg Config, binanceClient *binance.Client) *Executor {
	return &Executor{config: cfg, binanceClient: binanceClient}
}

func (d *Executor) Execute() {
	slog.Info("Executing DCA")
	tickerStat, err := d.binanceClient.GetTickerStat(d.config.Symbol)
	if err != nil {
		slog.Error("Failed to get ticker stat: ", slog.Any("error", err))
	}

	slog.Info(
		"Ticker stat",
		slog.Float64("price_change_percentage", tickerStat.PriceChangePercentage()),
		slog.Float64("last_price", tickerStat.LastPrice()),
	)

	quantityUSDT := d.config.QuantityUSDT

	if tickerStat.PriceChangePercentage() < -5 {
		quantityUSDT *= 2
	}

	quantityBTC := quantityUSDT / tickerStat.LastPrice()

	slog.Info("Placing buy order", slog.Float64("quantity_btc", quantityBTC))

	err = d.binanceClient.PlaceBuyOrder(d.config.Symbol, quantityBTC)
	if err != nil {
		slog.Error("Failed to place buy order: ", slog.Any("error", err))
	}
}
