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
	tickerStat, err := d.binanceClient.GetTickerStat(d.config.Symbol)
	if err != nil {
		slog.Error("Failed to get ticker stat: ", slog.Any("error", err))
	}

	quantityUSDT := d.config.QuantityUSDT

	if tickerStat.PriceChangePercentage() < -5 {
		quantityUSDT *= 2
	}

	quantityBTC := quantityUSDT / tickerStat.LastPrice()

	err = d.binanceClient.PlaceBuyOrder(d.config.Symbol, quantityBTC)
	if err != nil {
		slog.Error("Failed to place buy order: ", slog.Any("error", err))
	}
}
