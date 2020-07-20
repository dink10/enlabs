package storage

import (
	"context"
	"fmt"

	"github.com/go-pg/pg/v9"

	"github.com/dink10/enlabs/internal/pkg/payments"
)

// NewPaymentStorage returns a new instance of PaymentStorage.
func NewPaymentStorage(db *pg.DB) *PaymentStorage {
	return &PaymentStorage{db: db}
}

// PaymentStorage provides access to postgres database and
// implements PAYMENT.Storage interface.
type PaymentStorage struct {
	db *pg.DB
}

// Insert inserts user into users table.
func (s *PaymentStorage) SourceTypes(ctx context.Context) ([]payments.SourceType, error) {
	var sourceTypes []payments.SourceType
	if err := s.db.ModelContext(ctx, &sourceTypes).Select(); err != nil {
		return nil, fmt.Errorf("failed to execute insert query: %v", err)
	}

	return sourceTypes, nil
}

// Insert inserts user into users table.
func (s *PaymentStorage) ProceedPayment(ctx context.Context, payment payments.Payment) error {
	err := s.db.RunInTransaction(func(tx *pg.Tx) error {
		_, err := tx.ModelContext(ctx, &payment).Insert()
		if err != nil {
			return err
		}

		var account payments.Account
		query := tx.ModelContext(ctx, &account).
			Where("id=?", ctx.Value("account_id"))
		switch payment.State {
		case "win":
			query.Set("balance=balance+?", payment.Amount)
		case "lost":
			query.Where("balance >= ?", payment.Amount)
			query.Set("balance=balance-?", payment.Amount)
		}

		query.Returning("id").Returning("balance")
		if _, err := query.Update(); err != nil {
			if payment.State == "lost" && err == pg.ErrNoRows {
				return fmt.Errorf("insufficient funds")
			}
			return fmt.Errorf("no such account")
		}

		return nil
	})

	if err != nil {
		pgErr, ok := err.(pg.Error)
		switch {
		case ok && pgErr.IntegrityViolation():
			return fmt.Errorf("transaction_id already processed")
		default:
			payment.Processed = false
			if _, err := s.db.ModelContext(ctx, &payment).Insert(); err != nil {
				return err
			}
		}
	}

	return err
}

func (s *PaymentStorage) Balance(ctx context.Context) (payments.Account, error) {
	var account payments.Account
	err := s.db.ModelContext(ctx, &account).
		Where("id=?", ctx.Value("account_id")).
		Select()

	return account, err
}
