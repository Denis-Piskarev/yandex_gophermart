package migrations

import (
	"context"
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upInitMigration, downInitMigration)
}

func upInitMigration(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	query := `
		CREATE TABLE IF NOT EXISTS users (
		    id SERIAL,
		    login varchar(255) NOT NULL,
		    password varchar(255) NOT NULL
		);

		CREATE TABLE IF NOT EXISTS orders (
		    number INT NOT NULL,
		    status varchar(255) NOT NULL,
		    accural INT,
		    uploaded_at timestamp NOT NULL default now(),
		    user_id INT NOT NULL
		)
	`

	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}

func downInitMigration(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	query := `DROP TABLE IF EXISTS users, orders;`
	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}
