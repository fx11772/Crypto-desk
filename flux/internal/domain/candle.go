package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

// Candle represents normalized OHLCV market data.
// Exchange uses lowercase (kraken).
// Symbol uses BASE/QUOTE format (BTC/USD).
// Timeframe uses lowercase (15m, 1h).
// Timestamp represents the opening time in UTC.
type Candle struct {
	Exchange  string
	Symbol    string
	Timeframe string
	Timestamp time.Time
	Open      decimal.Decimal
	High      decimal.Decimal
	Low       decimal.Decimal
	Close     decimal.Decimal
	Volume    decimal.Decimal
}
