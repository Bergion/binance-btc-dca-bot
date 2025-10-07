package binance

import (
	"fmt"
	"strconv"

	binancePkg "github.com/eranyanay/binance-api"
	"github.com/pkg/errors"
)

type Client struct {
	client *binancePkg.BinanceClient
}

func NewClient(cfg Config) *Client {
	client := binancePkg.NewBinanceClient(cfg.APIKey, cfg.APISecret)

	return &Client{client}
}

func (c *Client) GetTickerStat(symbol string) (*TickerStat, error) {
	priceChange, err := c.client.Ticker(&binancePkg.TickerOpts{Symbol: symbol})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	priceChangePercentage, err := strconv.ParseFloat(priceChange.PriceChangePercentage, 64)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	lastPrice, err := strconv.ParseFloat(priceChange.LastPrice, 64)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return NewTickerStat(priceChangePercentage, lastPrice)
}

func (c *Client) PlaceBuyOrder(symbol string, quantity float64) error {
	quantityStr := fmt.Sprintf("%.2f", quantity)

	_, err := c.client.NewOrder(&binancePkg.NewOrderOpts{
		Symbol:      symbol,
		Side:        binancePkg.OrderSideBuy,
		Type:        binancePkg.OrderTypeMarket,
		TimeInForce: binancePkg.TimeInForceGTC,
		Quantity:    quantityStr,
	})
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
