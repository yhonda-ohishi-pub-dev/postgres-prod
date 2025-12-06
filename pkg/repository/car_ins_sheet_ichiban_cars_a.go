package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrCarInsSheetIchibanCarsANotFound = errors.New("car ins sheet ichiban cars a entry not found")
)

// CarInsSheetIchibanCarsA represents the database model
type CarInsSheetIchibanCarsA struct {
	OrganizationID string
	IDCars         *string
	ElectCertMgNo  string
	GrantdateE     string
	GrantdateY     string
	GrantdateM     string
	GrantdateD     string
}

// CarInsSheetIchibanCarsARepository handles database operations for car_ins_sheet_ichiban_cars_a
type CarInsSheetIchibanCarsARepository struct {
	db DB
}

// NewCarInsSheetIchibanCarsARepository creates a new repository
func NewCarInsSheetIchibanCarsARepository(pool *pgxpool.Pool) *CarInsSheetIchibanCarsARepository {
	return &CarInsSheetIchibanCarsARepository{db: pool}
}

// NewCarInsSheetIchibanCarsARepositoryWithDB creates a repository with custom DB interface (for testing)
func NewCarInsSheetIchibanCarsARepositoryWithDB(db DB) *CarInsSheetIchibanCarsARepository {
	return &CarInsSheetIchibanCarsARepository{db: db}
}

// Create inserts a new car_ins_sheet_ichiban_cars_a entry
func (r *CarInsSheetIchibanCarsARepository) Create(ctx context.Context, organizationID, electCertMgNo, grantdateE, grantdateY, grantdateM, grantdateD string, idCars *string) (*CarInsSheetIchibanCarsA, error) {
	query := `
		INSERT INTO car_ins_sheet_ichiban_cars_a (organization_id, id_cars, "ElectCertMgNo", "GrantdateE", "GrantdateY", "GrantdateM", "GrantdateD")
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING organization_id, id_cars, "ElectCertMgNo", "GrantdateE", "GrantdateY", "GrantdateM", "GrantdateD"
	`

	var entry CarInsSheetIchibanCarsA
	err := r.db.QueryRow(ctx, query, organizationID, idCars, electCertMgNo, grantdateE, grantdateY, grantdateM, grantdateD).Scan(
		&entry.OrganizationID, &entry.IDCars, &entry.ElectCertMgNo, &entry.GrantdateE, &entry.GrantdateY, &entry.GrantdateM, &entry.GrantdateD,
	)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

// GetByPrimaryKey retrieves an entry by composite primary key (organization_id, ElectCertMgNo, GrantdateE, GrantdateY, GrantdateM, GrantdateD)
func (r *CarInsSheetIchibanCarsARepository) GetByPrimaryKey(ctx context.Context, organizationID, electCertMgNo, grantdateE, grantdateY, grantdateM, grantdateD string) (*CarInsSheetIchibanCarsA, error) {
	query := `
		SELECT organization_id, id_cars, "ElectCertMgNo", "GrantdateE", "GrantdateY", "GrantdateM", "GrantdateD"
		FROM car_ins_sheet_ichiban_cars_a
		WHERE organization_id = $1 AND "ElectCertMgNo" = $2 AND "GrantdateE" = $3 AND "GrantdateY" = $4 AND "GrantdateM" = $5 AND "GrantdateD" = $6
	`

	var entry CarInsSheetIchibanCarsA
	err := r.db.QueryRow(ctx, query, organizationID, electCertMgNo, grantdateE, grantdateY, grantdateM, grantdateD).Scan(
		&entry.OrganizationID, &entry.IDCars, &entry.ElectCertMgNo, &entry.GrantdateE, &entry.GrantdateY, &entry.GrantdateM, &entry.GrantdateD,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCarInsSheetIchibanCarsANotFound
		}
		return nil, err
	}

	return &entry, nil
}

// Update modifies an existing entry
func (r *CarInsSheetIchibanCarsARepository) Update(ctx context.Context, organizationID, electCertMgNo, grantdateE, grantdateY, grantdateM, grantdateD string, idCars *string) (*CarInsSheetIchibanCarsA, error) {
	query := `
		UPDATE car_ins_sheet_ichiban_cars_a
		SET id_cars = $7
		WHERE organization_id = $1 AND "ElectCertMgNo" = $2 AND "GrantdateE" = $3 AND "GrantdateY" = $4 AND "GrantdateM" = $5 AND "GrantdateD" = $6
		RETURNING organization_id, id_cars, "ElectCertMgNo", "GrantdateE", "GrantdateY", "GrantdateM", "GrantdateD"
	`

	var entry CarInsSheetIchibanCarsA
	err := r.db.QueryRow(ctx, query, organizationID, electCertMgNo, grantdateE, grantdateY, grantdateM, grantdateD, idCars).Scan(
		&entry.OrganizationID, &entry.IDCars, &entry.ElectCertMgNo, &entry.GrantdateE, &entry.GrantdateY, &entry.GrantdateM, &entry.GrantdateD,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCarInsSheetIchibanCarsANotFound
		}
		return nil, err
	}

	return &entry, nil
}

// Delete removes an entry by composite primary key
func (r *CarInsSheetIchibanCarsARepository) Delete(ctx context.Context, organizationID, electCertMgNo, grantdateE, grantdateY, grantdateM, grantdateD string) error {
	query := `
		DELETE FROM car_ins_sheet_ichiban_cars_a
		WHERE organization_id = $1 AND "ElectCertMgNo" = $2 AND "GrantdateE" = $3 AND "GrantdateY" = $4 AND "GrantdateM" = $5 AND "GrantdateD" = $6
	`

	result, err := r.db.Exec(ctx, query, organizationID, electCertMgNo, grantdateE, grantdateY, grantdateM, grantdateD)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrCarInsSheetIchibanCarsANotFound
	}

	return nil
}

// ListByOrganization retrieves entries for a specific organization with pagination
func (r *CarInsSheetIchibanCarsARepository) ListByOrganization(ctx context.Context, organizationID string, limit int, offset int) ([]*CarInsSheetIchibanCarsA, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	query := `
		SELECT organization_id, id_cars, "ElectCertMgNo", "GrantdateE", "GrantdateY", "GrantdateM", "GrantdateD"
		FROM car_ins_sheet_ichiban_cars_a
		WHERE organization_id = $1
		ORDER BY "ElectCertMgNo", "GrantdateE", "GrantdateY", "GrantdateM", "GrantdateD"
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, organizationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*CarInsSheetIchibanCarsA
	for rows.Next() {
		var entry CarInsSheetIchibanCarsA
		err := rows.Scan(&entry.OrganizationID, &entry.IDCars, &entry.ElectCertMgNo, &entry.GrantdateE, &entry.GrantdateY, &entry.GrantdateM, &entry.GrantdateD)
		if err != nil {
			return nil, err
		}
		entries = append(entries, &entry)
	}

	return entries, rows.Err()
}

// List retrieves all entries with pagination
func (r *CarInsSheetIchibanCarsARepository) List(ctx context.Context, limit int, offset int) ([]*CarInsSheetIchibanCarsA, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	query := `
		SELECT organization_id, id_cars, "ElectCertMgNo", "GrantdateE", "GrantdateY", "GrantdateM", "GrantdateD"
		FROM car_ins_sheet_ichiban_cars_a
		ORDER BY organization_id, "ElectCertMgNo", "GrantdateE", "GrantdateY", "GrantdateM", "GrantdateD"
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*CarInsSheetIchibanCarsA
	for rows.Next() {
		var entry CarInsSheetIchibanCarsA
		err := rows.Scan(&entry.OrganizationID, &entry.IDCars, &entry.ElectCertMgNo, &entry.GrantdateE, &entry.GrantdateY, &entry.GrantdateM, &entry.GrantdateD)
		if err != nil {
			return nil, err
		}
		entries = append(entries, &entry)
	}

	return entries, rows.Err()
}
