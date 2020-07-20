package main

import (
	"github.com/go-pg/pg/v9/orm"
	"github.com/robinjoseph08/go-pg-migrations/v2"
)

func init() {
	up := func(db orm.DB) error {
		_, err := db.Exec(`CREATE TABLE payments
			(
			    id             serial primary key,
			    created_at     timestamp   not null default now(),
			    account_id     int         not null references accounts (id),
			    transaction_id text unique not null,
			    state          text        not null,
			    amount         numeric     not null,
			    source_type    int         not null references source_types (id),
			    processed      bool        not null default false
			);
		`)

		return err
	}

	down := func(db orm.DB) error {
		_, err := db.Exec("DROP TABLE IF EXISTS payments;")
		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("000003_create_payments_table", up, down, opts)
}
