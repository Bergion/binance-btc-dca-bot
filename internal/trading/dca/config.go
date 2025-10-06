package dca

type Config struct {
	QuantityUSDT float64 `mapstructure:"quantity_usdt"`
	Symbol       string  `mapstructure:"symbol"`
}
