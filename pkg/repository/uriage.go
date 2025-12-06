package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrUriageNotFound = errors.New("uriage not found")
)

// Uriage represents the database model
type Uriage struct {
	Name           string
	Bumon          string
	OrganizationID string
	Kingaku        *int32
	Type           *int32
	Cam            *int32
	Date           string
}

// UriageRepository handles database operations for uriage
type UriageRepository struct {
	db DB
}

// NewUriageRepository creates a new repository
func NewUriageRepository(pool *pgxpool.Pool) *UriageRepository {
	return &UriageRepository{db: pool}
}

// NewUriageRepositoryWithDB creates a repository with custom DB interface (for testing)
func NewUriageRepositoryWithDB(db DB) *UriageRepository {
	return &UriageRepository{db: db}
}

// Create inserts a new uriage
func (r *UriageRepository) Create(ctx context.Context, name, bumon, organizationID string, kingaku, uriageType, cam *int32, date string) (*Uriage, error) {
	query := `
		INSERT INTO uriage (name, bumon, organization_id, kingaku, type, cam, date)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING name, bumon, organization_id, kingaku, type, cam, date
	`

	var uriage Uriage
	err := r.db.QueryRow(ctx, query, name, bumon, organizationID, kingaku, uriageType, cam, date).Scan(
		&uriage.Name, &uriage.Bumon, &uriage.OrganizationID, &uriage.Kingaku, &uriage.Type, &uriage.Cam, &uriage.Date,
	)
	if err != nil {
		return nil, err
	}

	return &uriage, nil
}

// GetByPrimaryKey retrieves a uriage by composite primary key (name, bumon, date, organization_id)
func (r *UriageRepository) GetByPrimaryKey(ctx context.Context, name, bumon, date, organizationID string) (*Uriage, error) {
	query := `
		SELECT name, bumon, organization_id, kingaku, type, cam, date
		FROM uriage
		WHERE name = $1 AND bumon = $2 AND date = $3 AND organization_id = $4
	`

	var uriage Uriage
	err := r.db.QueryRow(ctx, query, name, bumon, date, organizationID).Scan(
		&uriage.Name, &uriage.Bumon, &uriage.OrganizationID, &uriage.Kingaku, &uriage.Type, &uriage.Cam, &uriage.Date,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUriageNotFound
		}
		return nil, err
	}

	return &uriage, nil
}

// Update modifies an existing uriage
func (r *UriageRepository) Update(ctx context.Context, name, bumon, date, organizationID string, kingaku, uriageType, cam *int32) (*Uriage, error) {
	query := `
		UPDATE uriage
		SET kingaku = $5, type = $6, cam = $7
		WHERE name = $1 AND bumon = $2 AND date = $3 AND organization_id = $4
		RETURNING name, bumon, organization_id, kingaku, type, cam, date
	`

	var uriage Uriage
	err := r.db.QueryRow(ctx, query, name, bumon, date, organizationID, kingaku, uriageType, cam).Scan(
		&uriage.Name, &uriage.Bumon, &uriage.OrganizationID, &uriage.Kingaku, &uriage.Type, &uriage.Cam, &uriage.Date,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUriageNotFound
		}
		return nil, err
	}

	return &uriage, nil
}

// Delete removes a uriage (hard delete)
func (r *UriageRepository) Delete(ctx context.Context, name, bumon, date, organizationID string) error {
	query := `
		DELETE FROM uriage
		WHERE name = $1 AND bumon = $2 AND date = $3 AND organization_id = $4
	`

	result, err := r.db.Exec(ctx, query, name, bumon, date, organizationID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrUriageNotFound
	}

	return nil
}

// ListByOrganization retrieves uriage entries for a specific organization with pagination
func (r *UriageRepository) ListByOrganization(ctx context.Context, organizationID string, limit int, offset int) ([]*Uriage, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	query := `
		SELECT name, bumon, organization_id, kingaku, type, cam, date
		FROM uriage
		WHERE organization_id = $1
		ORDER BY date DESC, name ASC, bumon ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, organizationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var uriages []*Uriage
	for rows.Next() {
		var uriage Uriage
		err := rows.Scan(&uriage.Name, &uriage.Bumon, &uriage.OrganizationID, &uriage.Kingaku, &uriage.Type, &uriage.Cam, &uriage.Date)
		if err != nil {
			return nil, err
		}
		uriages = append(uriages, &uriage)
	}

	return uriages, rows.Err()
}

// List retrieves all uriage entries with pagination
func (r *UriageRepository) List(ctx context.Context, limit int, offset int) ([]*Uriage, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	query := `
		SELECT name, bumon, organization_id, kingaku, type, cam, date
		FROM uriage
		ORDER BY date DESC, name ASC, bumon ASC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var uriages []*Uriage
	for rows.Next() {
		var uriage Uriage
		err := rows.Scan(&uriage.Name, &uriage.Bumon, &uriage.OrganizationID, &uriage.Kingaku, &uriage.Type, &uriage.Cam, &uriage.Date)
		if err != nil {
			return nil, err
		}
		uriages = append(uriages, &uriage)
	}

	return uriages, rows.Err()
}
