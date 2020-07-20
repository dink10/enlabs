package processing

import (
	"context"
	"fmt"

	"github.com/go-pg/pg/v9"
	"github.com/jasonlvhit/gocron"
	"github.com/sirupsen/logrus"

	"github.com/dink10/enlabs/internal/pkg/config"
	"github.com/dink10/enlabs/internal/pkg/database"
	"github.com/dink10/enlabs/internal/pkg/logger"
	"github.com/dink10/enlabs/internal/pkg/payments"
)

// Run runs application.
func Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var cfg Config
	if err := config.LoadConfig(&cfg); err != nil {
		return fmt.Errorf("failed to parse config: %v", err)
	}

	if err := logger.Init(&cfg.Logger); err != nil {
		return fmt.Errorf("failed to initialize logger: %v", err)
	}

	db, err := database.Connect(ctx, &cfg.Database)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	defer database.Close(db)

	err = gocron.Every(cfg.Processing.CancellationTime).Minute().Do(func() {
		logrus.Info("Start of post processing")
		defer logrus.Info("End of post processing")
		var pays []payments.Payment
		err := db.ModelContext(ctx, &pays).
			Where("(id % 2) = 1").
			Where("processed = true").
			Order("id DESC").
			Limit(10).
			Select()
		if err != nil {
			logrus.Errorf("query error: %s", err)
			return
		}

		for _, p := range pays {
			err := db.RunInTransaction(func(tx *pg.Tx) error {
				var account payments.Account
				query := tx.ModelContext(ctx, &account).
					Where("id=?", p.AccountID)
				switch p.State {
				case "win":
					query.Where("balance >= ?", p.Amount)
					query.Set("balance=balance-?", p.Amount)
				case "lost":
					query.Set("balance=balance+?", p.Amount)
				}

				query.Returning("id").Returning("balance")

				if _, err := query.Update(); err != nil {
					if p.State == "win" && err == pg.ErrNoRows {
						return fmt.Errorf("insufficient funds")
					}
					return fmt.Errorf("no such account")
				}

				_, err := tx.Model(&p).WherePK().Set("processed = ?", false).Update()
				if err != nil {
					return fmt.Errorf("query error: %s", err)
				}

				return nil
			})
			switch {
			case err != nil:
				logrus.Errorf("Payment with transaction_id [%s] can't be cancelled due to: %s", p.TransactionID, err)
			default:
				logrus.Errorf("Payment with transaction_id [%s] successfully cancelled", p.TransactionID)
			}
		}
	})
	if err != nil {
		return err
	}

	<-gocron.Start()

	return nil
}
