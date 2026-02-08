package repository

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	"r2manager/domain"
)

type SettingsRepository struct {
	db *sql.DB
}

func NewSettingsRepository(db *sql.DB) *SettingsRepository {
	return &SettingsRepository{db: db}
}

func (r *SettingsRepository) GetAllBucketSettings(ctx context.Context) ([]domain.BucketSettings, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT bucket_name, public_url FROM bucket_settings`)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query bucket settings")
	}
	defer rows.Close()

	var settings []domain.BucketSettings
	for rows.Next() {
		var s domain.BucketSettings
		if err := rows.Scan(&s.BucketName, &s.PublicUrl); err != nil {
			return nil, errors.Wrap(err, "failed to scan bucket settings")
		}
		settings = append(settings, s)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "failed to iterate bucket settings")
	}

	return settings, nil
}

func (r *SettingsRepository) GetBucketSettings(ctx context.Context, bucketName string) (*domain.BucketSettings, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT bucket_name, public_url FROM bucket_settings WHERE bucket_name = ?`,
		bucketName,
	)

	var s domain.BucketSettings
	err := row.Scan(&s.BucketName, &s.PublicUrl)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to query bucket settings")
	}

	return &s, nil
}

func (r *SettingsRepository) UpsertBucketSettings(ctx context.Context, bucketName, publicUrl string) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO bucket_settings (bucket_name, public_url) VALUES (?, ?)
		 ON CONFLICT(bucket_name) DO UPDATE SET public_url = excluded.public_url`,
		bucketName, publicUrl,
	)
	if err != nil {
		return errors.Wrap(err, "failed to upsert bucket settings")
	}

	return nil
}

func (r *SettingsRepository) BulkUpsertBucketSettings(ctx context.Context, settings []domain.BucketSettings) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx,
		`INSERT INTO bucket_settings (bucket_name, public_url) VALUES (?, ?)
		 ON CONFLICT(bucket_name) DO UPDATE SET public_url = excluded.public_url`,
	)
	if err != nil {
		return errors.Wrap(err, "failed to prepare statement")
	}
	defer stmt.Close()

	for _, s := range settings {
		if _, err := stmt.ExecContext(ctx, s.BucketName, s.PublicUrl); err != nil {
			return errors.Wrap(err, "failed to upsert bucket settings")
		}
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}
