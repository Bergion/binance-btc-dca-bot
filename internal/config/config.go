package config

import (
	"log"

	"github.com/Bergion/binance-btc-dca-bot/internal/trading/dca"
	"github.com/Bergion/binance-btc-dca-bot/pkg/binance"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

type AppConfig struct {
	Binance binance.Config `mapstructure:"binance"`
	DCA     dca.Config     `mapstructure:"dca"`
}

type Result struct {
	fx.Out

	Binance binance.Config
	DCA     dca.Config
}

func NewAppConfig() (Result, error) {
	config := AppConfig{}

	viper.SetConfigName("config")
	viper.SetConfigType("toml")

	viper.AddConfigPath("./")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Failed to read config: ", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatal("Failed to unmarshal config: ", err)
	}

	return Result{
		Binance: config.Binance,
		DCA:     config.DCA,
	}, nil
}
