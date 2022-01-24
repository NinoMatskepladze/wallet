package models

type WalletID string

// Wallet type
type Wallet struct {
	ID      WalletID `json:"id"`
	Balance int64    `json:"balance"`
}
