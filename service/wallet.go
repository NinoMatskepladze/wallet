package service

import (
	"context"
	"log"

	"github.com/NinoMatskepladze/wallet/db"
	"github.com/NinoMatskepladze/wallet/models"
	"github.com/google/uuid"
)

type Service struct {
	db *db.Datastore
}

// NewService defines new service for wallet
func NewService(db *db.Datastore) *Service {
	return &Service{
		db: db,
	}
}

// CreateWallet creates a new wallet with default balance 0
func (s *Service) CreateWallet(ctx context.Context) (models.Wallet, error) {
	newWalletID := uuid.New().String()
	wallet := &models.Wallet{
		ID:      models.WalletID(newWalletID),
		Balance: 0,
	}

	_, err := s.db.DB.ExecContext(
		ctx,
		"insert into wallets (id, balance, created_at, updated_at) values ($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);",
		wallet.ID, wallet.Balance,
	)

	if err != nil {
		return models.Wallet{}, err
	}
	return *wallet, nil
}

// UpdateBalance updates balance using amount which can be both negative and positive
// based on that balance decreases or increases. Insufficient balance cant be subtracted
func (s *Service) UpdateBalance(ctx context.Context, walletID string, addBalanceRequest models.AddBalanceRequest) error {
	// Create a new context, and begin a transaction
	tx, err := s.db.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	newTransactionID := uuid.New().String()
	finalBalance := new(int)
	// Query for updating wallet
	// If balance is not sifficient it will throw an error
	err = tx.QueryRow(`UPDATE wallets
	SET balance = balance + $2, updated_at = CURRENT_TIMESTAMP
	WHERE id = $1
	RETURNING balance`, walletID, addBalanceRequest.Amount).Scan(finalBalance)
	if err != nil {
		// In case there is an error in the query execution, rollback the transaction
		tx.Rollback()
		return err
	}

	// The next query is handled similarly
	_, err = tx.ExecContext(ctx, `INSERT INTO transactions (id, wallet_id, balance, amount, created_at)
	values ($1, $2, $3, $4, CURRENT_TIMESTAMP)
	`, newTransactionID, walletID, finalBalance, addBalanceRequest.Amount)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction if there were no Query execution errors
	return tx.Commit()
}

// GetWallet returns current state of wallet
func (s *Service) GetWallet(ctx context.Context, walletID string) (models.Wallet, error) {
	wallet := &models.Wallet{}

	row := s.db.DB.QueryRow("SELECT id, balance FROM wallets WHERE id=$1;", walletID)
	err := row.Scan(&wallet.ID, &wallet.Balance)
	return *wallet, err
}
