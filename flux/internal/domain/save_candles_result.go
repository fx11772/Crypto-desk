package domain

// SaveCandlesResult represents the outcome of a successful
// candle batch insertion.
type SaveCandlesResult struct {
	InsertedCount  int
	DuplicateCount int
}
