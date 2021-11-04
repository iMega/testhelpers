package db

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
)

func Test_createDB(t *testing.T) {
	_, close, err := Create("", func(ctx context.Context, tx *sql.Tx) error {
		if err := createEmailTable(ctx, tx); err != nil {
			return err
		}

		if err := addEmail(ctx, tx, "info@example.com"); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		t.Errorf("failed to create db, %s", err)
	}
	defer close()
}

func createEmailTable(ctx context.Context, tx *sql.Tx) error {
	q := `CREATE TABLE IF NOT EXISTS email (
        email VARCHAR(16) NOT NULL
    )`

	if _, err := tx.ExecContext(ctx, q); err != nil {
		return fmt.Errorf("failed to execute query, %w", err)
	}

	return nil
}

func addEmail(ctx context.Context, tx *sql.Tx, email string) error {
	q := `insert into email (email) values (?)`

	if _, err := tx.ExecContext(ctx, q, email); err != nil {
		return fmt.Errorf("failed to execute query, %w", err)
	}

	return nil
}
