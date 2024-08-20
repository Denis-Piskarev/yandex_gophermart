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
		    password varchar(255) NOT NULL,
		    current FLOAT DEFAULT 0,
		    withdrawn INT  DEFAULT 0
		);

		CREATE TABLE IF NOT EXISTS orders (
		    id SERIAL,
		    number BIGINT NOT NULL,
		    status varchar(255) NOT NULL,
		    accrual INT,
		    uploaded_at timestamptz NOT NULL default now(),
		    user_id INT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS withdrawals (
		    id SERIAL,
		    user_id INT NOT NULL,
		    number BIGINT NOT NULL,
		    sum INT NOT NULL,
		    processed_at timestamptz NOT NULL default now()
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
