package binance

type TickerStat struct {
	priceChangePercentage float64
	lastPrice             float64
}

func NewTickerStat(priceChangePercentage float64, lastPrice float64) (*TickerStat, error) {
	return &TickerStat{priceChangePercentage: priceChangePercentage, lastPrice: lastPrice}, nil
}

func (t *TickerStat) PriceChangePercentage() float64 {
	return t.priceChangePercentage
}

func (t *TickerStat) LastPrice() float64 {
	return t.lastPrice
}
