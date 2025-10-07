package binance

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	binancePkg "github.com/adshao/go-binance/v2"
	"github.com/pkg/errors"
)

var (
	binanceBase = "https://api.binance.com"
)

type Client struct {
	client *binancePkg.Client
}

func NewClient(cfg Config) *Client {
	client := binancePkg.NewClient(cfg.APIKey, cfg.APISecret)

	return &Client{client}
}

func (c *Client) GetTickerStat(symbol string) (*TickerStat, error) {
	priceChangeStats, err := c.client.NewListPriceChangeStatsService().Symbol(symbol).Do(context.Background())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	priceChangeStat := priceChangeStats[0]

	priceChangePercentage, err := strconv.ParseFloat(priceChangeStat.PriceChangePercent, 64)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	lastPrice, err := strconv.ParseFloat(priceChangeStat.LastPrice, 64)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return NewTickerStat(priceChangePercentage, lastPrice)
}

func (c *Client) PlaceBuyOrder(symbol string, quantity float64) error {
	quantityStr := fmt.Sprintf("%.5f", quantity)

	res, err := c.client.NewCreateOrderService().
		Symbol(symbol).
		Side(binancePkg.SideTypeBuy).
		Type(binancePkg.OrderTypeMarket).
		Quantity(quantityStr).
		Do(context.Background())
	if err != nil {
		return errors.WithStack(err)
	}

	slog.Info("Buy order placed", slog.Any("response", res))

	return nil
}

func (c *Client) RedeemFlexible(
	productId string,
	amount string,
) error {
	destAccount := "SPOT"
	apiKey := c.client.APIKey
	secKey := c.client.SecretKey

	if apiKey == "" || secKey == "" {
		return errors.New("api key or secret key is empty")
	}

	values := url.Values{}
	values.Set("productId", productId)
	if amount == "" {
		return errors.New("amount must be set when redeemAll is false")
	}

	values.Set("amount", amount)
	values.Set("destAccount", destAccount)

	timestamp := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	values.Set("timestamp", timestamp)

	queryString := values.Encode()
	signature := sign(queryString, secKey)

	fullBody := queryString + "&signature=" + signature

	endpoint := binanceBase + "/sapi/v1/simple-earn/flexible/redeem"
	req, err := http.NewRequest("POST", endpoint, strings.NewReader(fullBody))
	if err != nil {
		return errors.WithStack(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-MBX-APIKEY", apiKey)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return errors.WithStack(err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	var jsonResp map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &jsonResp); err != nil {
		return fmt.Errorf("status: %d, body: %s, jsonErr: %v", resp.StatusCode, string(bodyBytes), err)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("non-200 response: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

func sign(message, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}
