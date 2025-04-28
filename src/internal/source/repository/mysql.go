package repository

import (
	"context"
	"xyz_multifinance/src/internal/source/model"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return Repository{db: db}
}

func (r Repository) Create(ctx context.Context, source *model.Source) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO Sources 
			(source_id, secret_hash, category, name, email)
		VALUES
			(?, ?, ?, ?, ?)
	`, source.ID, source.SecretHash, source.Category, source.Name, source.Email)
	if err != nil {
		return err
	}

	return nil
}

func (r Repository) FindByID(ctx context.Context, id string) (*model.Source, error) {
	query := `SELECT source_id, secret_hash, category, name, email FROM Sources WHERE source_id = ?`

	row := r.db.QueryRowContext(ctx, query, id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var source model.Source
	if err := row.Scan(&source.ID, &source.SecretHash, &source.Category, &source.Name, &source.Email); err != nil {
		return nil, err
	}

	return &source, nil
}
