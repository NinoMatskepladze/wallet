package models

type WalletID string

// Wallet model
type Wallet struct {
	ID      WalletID `json:"id"`
	Balance int64    `json:"balance"`
}
