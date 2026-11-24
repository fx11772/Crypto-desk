package domain

import "time"

// CollectionTarget represents a configured market data collection target.
// ID is zero when the target has not yet been persisted.
// Exchange uses lowercase (kraken).
// Symbol uses BASE/QUOTE format (BTC/USD).
// Timeframe uses lowercase (15m, 1h).
// CreatedAt and UpdatedAt are expressed in UTC.
type CollectionTarget struct {
	ID        int64
	Exchange  string
	Symbol    string
	Timeframe string
	Enabled   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
