package main

import (
	"github.com/go-pg/pg/v9/orm"
	"github.com/robinjoseph08/go-pg-migrations/v2"
)

func init() {
	up := func(db orm.DB) error {
		_, err := db.Exec(`CREATE TABLE accounts
			(
			    id      serial primary key,
			    balance numeric(12,2) default 0
			);
		`)
		if err != nil {
			return err
		}

		_, err = db.Exec(`INSERT INTO accounts (balance) VALUES (0.00);`)

		return err
	}

	down := func(db orm.DB) error {
		_, err := db.Exec("DROP TABLE IF EXISTS accounts;")
		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("000002_create_accounts_table", up, down, opts)
}
