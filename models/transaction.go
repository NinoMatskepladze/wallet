package models

import "time"

// Transaction model
type Transaction struct {
	ID          string    `json:"id"`
	WalletID    string    `json:"wallet_id"`
	Balance     int       `json:"balance"`
	Amount      int       `json:"amount"`
	TimeCreated time.Time `json:"time_created"`
}
