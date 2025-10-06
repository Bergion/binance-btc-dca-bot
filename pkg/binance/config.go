package binance

type Config struct {
	APIKey    string `mapstructure:"api_key"`
	APISecret string `mapstructure:"api_secret"`
}
