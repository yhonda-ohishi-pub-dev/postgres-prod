package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrDtakoCarsIchibanCarsNotFound = errors.New("dtako cars ichiban cars entry not found")
)

// DtakoCarsIchibanCars represents the database model
type DtakoCarsIchibanCars struct {
	IdDtako        string
	OrganizationID string
	Id             *string
}

// DtakoCarsIchibanCarsRepository handles database operations for dtako_cars_ichiban_cars
type DtakoCarsIchibanCarsRepository struct {
	db DB
}

// NewDtakoCarsIchibanCarsRepository creates a new repository
func NewDtakoCarsIchibanCarsRepository(pool *pgxpool.Pool) *DtakoCarsIchibanCarsRepository {
	return &DtakoCarsIchibanCarsRepository{db: pool}
}

// NewDtakoCarsIchibanCarsRepositoryWithDB creates a repository with custom DB interface (for testing)
func NewDtakoCarsIchibanCarsRepositoryWithDB(db DB) *DtakoCarsIchibanCarsRepository {
	return &DtakoCarsIchibanCarsRepository{db: db}
}

// Create inserts a new dtako_cars_ichiban_cars entry
func (r *DtakoCarsIchibanCarsRepository) Create(ctx context.Context, idDtako, organizationID string, id *string) (*DtakoCarsIchibanCars, error) {
	query := `
		INSERT INTO dtako_cars_ichiban_cars (id_dtako, organization_id, id)
		VALUES ($1, $2, $3)
		RETURNING id_dtako, organization_id, id
	`

	var entry DtakoCarsIchibanCars
	err := r.db.QueryRow(ctx, query, idDtako, organizationID, id).Scan(
		&entry.IdDtako, &entry.OrganizationID, &entry.Id,
	)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

// GetByDtakoAndOrg retrieves an entry by composite primary key (id_dtako, organization_id)
func (r *DtakoCarsIchibanCarsRepository) GetByDtakoAndOrg(ctx context.Context, idDtako, organizationID string) (*DtakoCarsIchibanCars, error) {
	query := `
		SELECT id_dtako, organization_id, id
		FROM dtako_cars_ichiban_cars
		WHERE id_dtako = $1 AND organization_id = $2
	`

	var entry DtakoCarsIchibanCars
	err := r.db.QueryRow(ctx, query, idDtako, organizationID).Scan(
		&entry.IdDtako, &entry.OrganizationID, &entry.Id,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrDtakoCarsIchibanCarsNotFound
		}
		return nil, err
	}

	return &entry, nil
}

// Update modifies an existing entry
func (r *DtakoCarsIchibanCarsRepository) Update(ctx context.Context, idDtako, organizationID string, id *string) (*DtakoCarsIchibanCars, error) {
	query := `
		UPDATE dtako_cars_ichiban_cars
		SET id = $3
		WHERE id_dtako = $1 AND organization_id = $2
		RETURNING id_dtako, organization_id, id
	`

	var entry DtakoCarsIchibanCars
	err := r.db.QueryRow(ctx, query, idDtako, organizationID, id).Scan(
		&entry.IdDtako, &entry.OrganizationID, &entry.Id,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrDtakoCarsIchibanCarsNotFound
		}
		return nil, err
	}

	return &entry, nil
}

// Delete removes an entry by composite primary key
func (r *DtakoCarsIchibanCarsRepository) Delete(ctx context.Context, idDtako, organizationID string) error {
	query := `
		DELETE FROM dtako_cars_ichiban_cars
		WHERE id_dtako = $1 AND organization_id = $2
	`

	result, err := r.db.Exec(ctx, query, idDtako, organizationID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrDtakoCarsIchibanCarsNotFound
	}

	return nil
}

// ListByOrganization retrieves entries for a specific organization with pagination
func (r *DtakoCarsIchibanCarsRepository) ListByOrganization(ctx context.Context, organizationID string, limit int, offset int) ([]*DtakoCarsIchibanCars, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	query := `
		SELECT id_dtako, organization_id, id
		FROM dtako_cars_ichiban_cars
		WHERE organization_id = $1
		ORDER BY id_dtako ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, organizationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*DtakoCarsIchibanCars
	for rows.Next() {
		var entry DtakoCarsIchibanCars
		err := rows.Scan(&entry.IdDtako, &entry.OrganizationID, &entry.Id)
		if err != nil {
			return nil, err
		}
		entries = append(entries, &entry)
	}

	return entries, rows.Err()
}
