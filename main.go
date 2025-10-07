package main

import (
	"context"

	"github.com/Bergion/binance-btc-dca-bot/internal/config"
	"github.com/Bergion/binance-btc-dca-bot/internal/trading/dca"
	"github.com/Bergion/binance-btc-dca-bot/pkg/binance"
	"github.com/robfig/cron/v3"
	"go.uber.org/fx"
)

func main() {
	fx.New(createApp()).Run()
}

func createApp() fx.Option {
	return fx.Options(
		config.Module,
		binance.Module,
		dca.Module,
		fx.Provide(
			NewDCACron,
		),
		fx.Invoke(
			func(lc fx.Lifecycle, cron *cron.Cron) {
				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						go cron.Run()

						return nil
					},
					OnStop: func(ctx context.Context) error {
						cron.Stop()

						return nil
					},
				})
			},
		),
	)
}

func NewDCACron(executor *dca.Executor) *cron.Cron {
	c := cron.New(
		cron.WithChain(
			cron.Recover(cron.DefaultLogger),
		),
	)
	c.AddFunc("30 13 * * *", executor.Execute)

	return c
}
