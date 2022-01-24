package service

import (
	"context"

	"github.com/NinoMatskepladze/wallet/models"
)

// Service Iface for database related methods
type ServiceIface interface {
	CreateWallet(context.Context) (models.Wallet, error)
	UpdateBalance(context.Context, string, models.AddBalanceRequest) error
	GetWallet(context.Context, string) (models.Wallet, error)
}
