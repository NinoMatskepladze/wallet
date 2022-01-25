package models

// UpdateWalletRequest type
type UpdateWalletRequest struct {
	Amount int `json:"amount" validate:"required"`
}
