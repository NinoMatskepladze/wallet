package test

import (
	"context"

	"github.com/NinoMatskepladze/wallet/models"
)

// FakeService for mocking Wallet service Iface
type FakeService struct {
}

// CreateWallet dummy mock func
func (f FakeService) CreateWallet(context.Context) (models.Wallet, error) {
	return models.Wallet{
		ID:      "fakeID",
		Balance: 0,
	}, nil
}

// UpdateBalance dummy mock func
func (f FakeService) UpdateBalance(context.Context, string, models.UpdateWalletRequest) error {
	return nil
}

// GetWallet dummy mock func
func (f FakeService) GetWallet(context.Context, string) (models.GetWalletResponse, error) {
	return models.GetWalletResponse{
		Wallet: models.Wallet{
			ID:      "FakeID",
			Balance: 100,
		},
	}, nil
}
