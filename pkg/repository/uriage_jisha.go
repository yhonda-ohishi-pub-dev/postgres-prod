package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrUriageJishaNotFound = errors.New("uriage jisha not found")
)

// UriageJisha represents the database model
type UriageJisha struct {
	Bumon          string
	OrganizationID string
	Kingaku        *int32
	Type           *int32
	Date           string
}

// UriageJishaRepository handles database operations for uriage_jisha
type UriageJishaRepository struct {
	db DB
}

// NewUriageJishaRepository creates a new repository
func NewUriageJishaRepository(pool *pgxpool.Pool) *UriageJishaRepository {
	return &UriageJishaRepository{db: pool}
}

// NewUriageJishaRepositoryWithDB creates a repository with custom DB interface (for testing)
func NewUriageJishaRepositoryWithDB(db DB) *UriageJishaRepository {
	return &UriageJishaRepository{db: db}
}

// Create inserts a new uriage jisha
func (r *UriageJishaRepository) Create(ctx context.Context, bumon, organizationID string, kingaku, typeVal *int32, date string) (*UriageJisha, error) {
	query := `
		INSERT INTO uriage_jisha (bumon, organization_id, kingaku, type, date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING bumon, organization_id, kingaku, type, date
	`

	var uriageJisha UriageJisha
	err := r.db.QueryRow(ctx, query, bumon, organizationID, kingaku, typeVal, date).Scan(
		&uriageJisha.Bumon, &uriageJisha.OrganizationID, &uriageJisha.Kingaku, &uriageJisha.Type, &uriageJisha.Date,
	)
	if err != nil {
		return nil, err
	}

	return &uriageJisha, nil
}

// GetByPrimaryKey retrieves a uriage jisha by composite primary key (bumon, date, organization_id)
func (r *UriageJishaRepository) GetByPrimaryKey(ctx context.Context, bumon, date, organizationID string) (*UriageJisha, error) {
	query := `
		SELECT bumon, organization_id, kingaku, type, date
		FROM uriage_jisha
		WHERE bumon = $1 AND date = $2 AND organization_id = $3
	`

	var uriageJisha UriageJisha
	err := r.db.QueryRow(ctx, query, bumon, date, organizationID).Scan(
		&uriageJisha.Bumon, &uriageJisha.OrganizationID, &uriageJisha.Kingaku, &uriageJisha.Type, &uriageJisha.Date,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUriageJishaNotFound
		}
		return nil, err
	}

	return &uriageJisha, nil
}

// Update modifies an existing uriage jisha
func (r *UriageJishaRepository) Update(ctx context.Context, bumon, date, organizationID string, kingaku, typeVal *int32) (*UriageJisha, error) {
	query := `
		UPDATE uriage_jisha
		SET kingaku = $4, type = $5
		WHERE bumon = $1 AND date = $2 AND organization_id = $3
		RETURNING bumon, organization_id, kingaku, type, date
	`

	var uriageJisha UriageJisha
	err := r.db.QueryRow(ctx, query, bumon, date, organizationID, kingaku, typeVal).Scan(
		&uriageJisha.Bumon, &uriageJisha.OrganizationID, &uriageJisha.Kingaku, &uriageJisha.Type, &uriageJisha.Date,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUriageJishaNotFound
		}
		return nil, err
	}

	return &uriageJisha, nil
}

// Delete removes a uriage jisha
func (r *UriageJishaRepository) Delete(ctx context.Context, bumon, date, organizationID string) error {
	query := `
		DELETE FROM uriage_jisha
		WHERE bumon = $1 AND date = $2 AND organization_id = $3
	`

	result, err := r.db.Exec(ctx, query, bumon, date, organizationID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrUriageJishaNotFound
	}

	return nil
}

// ListByOrganization retrieves all uriage jisha entries for an organization with pagination
func (r *UriageJishaRepository) ListByOrganization(ctx context.Context, organizationID string, limit int, offset int) ([]*UriageJisha, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	query := `
		SELECT bumon, organization_id, kingaku, type, date
		FROM uriage_jisha
		WHERE organization_id = $1
		ORDER BY date DESC, bumon
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, organizationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var uriageJishas []*UriageJisha
	for rows.Next() {
		var uriageJisha UriageJisha
		err := rows.Scan(&uriageJisha.Bumon, &uriageJisha.OrganizationID, &uriageJisha.Kingaku, &uriageJisha.Type, &uriageJisha.Date)
		if err != nil {
			return nil, err
		}
		uriageJishas = append(uriageJishas, &uriageJisha)
	}

	return uriageJishas, rows.Err()
}

// List retrieves all uriage jisha entries with pagination
func (r *UriageJishaRepository) List(ctx context.Context, limit int, offset int) ([]*UriageJisha, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	query := `
		SELECT bumon, organization_id, kingaku, type, date
		FROM uriage_jisha
		ORDER BY date DESC, bumon
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var uriageJishas []*UriageJisha
	for rows.Next() {
		var uriageJisha UriageJisha
		err := rows.Scan(&uriageJisha.Bumon, &uriageJisha.OrganizationID, &uriageJisha.Kingaku, &uriageJisha.Type, &uriageJisha.Date)
		if err != nil {
			return nil, err
		}
		uriageJishas = append(uriageJishas, &uriageJisha)
	}

	return uriageJishas, rows.Err()
}
