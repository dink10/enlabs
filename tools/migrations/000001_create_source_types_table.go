package main

import (
	"github.com/go-pg/pg/v9/orm"
	"github.com/robinjoseph08/go-pg-migrations/v2"
)

func init() {
	up := func(db orm.DB) error {
		_, err := db.Exec(`CREATE TABLE source_types
			(
			    id    serial primary key,
			    value text not null
			);
		`)
		if err != nil {
			return err
		}

		_, err = db.Exec(`INSERT INTO source_types (value) 
			VALUES ('game'), ('server'), ('payment');
		`)

		return err
	}

	down := func(db orm.DB) error {
		_, err := db.Exec("DROP TABLE IF EXISTS source_types;")
		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("000002_create_accounts_table", up, down, opts)
}
