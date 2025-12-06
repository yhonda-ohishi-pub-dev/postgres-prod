package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrCarInsSheetIchibanCarsNotFound = errors.New("car_ins_sheet_ichiban_cars not found")
)

// CarInsSheetIchibanCars represents the database model
type CarInsSheetIchibanCars struct {
	OrganizationID         string
	IDCars                 *string
	ElectCertMgNo          string
	ElectCertPublishdateE  string
	ElectCertPublishdateY  string
	ElectCertPublishdateM  string
	ElectCertPublishdateD  string
}

// CarInsSheetIchibanCarsRepository handles database operations for car_ins_sheet_ichiban_cars
type CarInsSheetIchibanCarsRepository struct {
	db DB
}

// NewCarInsSheetIchibanCarsRepository creates a new repository
func NewCarInsSheetIchibanCarsRepository(pool *pgxpool.Pool) *CarInsSheetIchibanCarsRepository {
	return &CarInsSheetIchibanCarsRepository{db: pool}
}

// NewCarInsSheetIchibanCarsRepositoryWithDB creates a repository with custom DB interface (for testing)
func NewCarInsSheetIchibanCarsRepositoryWithDB(db DB) *CarInsSheetIchibanCarsRepository {
	return &CarInsSheetIchibanCarsRepository{db: db}
}

// Create inserts a new car_ins_sheet_ichiban_cars record
func (r *CarInsSheetIchibanCarsRepository) Create(ctx context.Context, organizationID, electCertMgNo, electCertPublishdateE, electCertPublishdateY, electCertPublishdateM, electCertPublishdateD string, idCars *string) (*CarInsSheetIchibanCars, error) {
	query := `
		INSERT INTO car_ins_sheet_ichiban_cars (organization_id, id_cars, "ElectCertMgNo", "ElectCertPublishdateE", "ElectCertPublishdateY", "ElectCertPublishdateM", "ElectCertPublishdateD")
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING organization_id, id_cars, "ElectCertMgNo", "ElectCertPublishdateE", "ElectCertPublishdateY", "ElectCertPublishdateM", "ElectCertPublishdateD"
	`

	var record CarInsSheetIchibanCars
	err := r.db.QueryRow(ctx, query, organizationID, idCars, electCertMgNo, electCertPublishdateE, electCertPublishdateY, electCertPublishdateM, electCertPublishdateD).Scan(
		&record.OrganizationID, &record.IDCars, &record.ElectCertMgNo, &record.ElectCertPublishdateE, &record.ElectCertPublishdateY, &record.ElectCertPublishdateM, &record.ElectCertPublishdateD,
	)
	if err != nil {
		return nil, err
	}

	return &record, nil
}

// GetByPrimaryKey retrieves a car_ins_sheet_ichiban_cars record by composite primary key
func (r *CarInsSheetIchibanCarsRepository) GetByPrimaryKey(ctx context.Context, organizationID, electCertMgNo, electCertPublishdateE, electCertPublishdateY, electCertPublishdateM, electCertPublishdateD string) (*CarInsSheetIchibanCars, error) {
	query := `
		SELECT organization_id, id_cars, "ElectCertMgNo", "ElectCertPublishdateE", "ElectCertPublishdateY", "ElectCertPublishdateM", "ElectCertPublishdateD"
		FROM car_ins_sheet_ichiban_cars
		WHERE organization_id = $1 AND "ElectCertMgNo" = $2 AND "ElectCertPublishdateE" = $3 AND "ElectCertPublishdateY" = $4 AND "ElectCertPublishdateM" = $5 AND "ElectCertPublishdateD" = $6
	`

	var record CarInsSheetIchibanCars
	err := r.db.QueryRow(ctx, query, organizationID, electCertMgNo, electCertPublishdateE, electCertPublishdateY, electCertPublishdateM, electCertPublishdateD).Scan(
		&record.OrganizationID, &record.IDCars, &record.ElectCertMgNo, &record.ElectCertPublishdateE, &record.ElectCertPublishdateY, &record.ElectCertPublishdateM, &record.ElectCertPublishdateD,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCarInsSheetIchibanCarsNotFound
		}
		return nil, err
	}

	return &record, nil
}

// Update modifies an existing car_ins_sheet_ichiban_cars record
func (r *CarInsSheetIchibanCarsRepository) Update(ctx context.Context, organizationID, electCertMgNo, electCertPublishdateE, electCertPublishdateY, electCertPublishdateM, electCertPublishdateD string, idCars *string) (*CarInsSheetIchibanCars, error) {
	query := `
		UPDATE car_ins_sheet_ichiban_cars
		SET id_cars = $7
		WHERE organization_id = $1 AND "ElectCertMgNo" = $2 AND "ElectCertPublishdateE" = $3 AND "ElectCertPublishdateY" = $4 AND "ElectCertPublishdateM" = $5 AND "ElectCertPublishdateD" = $6
		RETURNING organization_id, id_cars, "ElectCertMgNo", "ElectCertPublishdateE", "ElectCertPublishdateY", "ElectCertPublishdateM", "ElectCertPublishdateD"
	`

	var record CarInsSheetIchibanCars
	err := r.db.QueryRow(ctx, query, organizationID, electCertMgNo, electCertPublishdateE, electCertPublishdateY, electCertPublishdateM, electCertPublishdateD, idCars).Scan(
		&record.OrganizationID, &record.IDCars, &record.ElectCertMgNo, &record.ElectCertPublishdateE, &record.ElectCertPublishdateY, &record.ElectCertPublishdateM, &record.ElectCertPublishdateD,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCarInsSheetIchibanCarsNotFound
		}
		return nil, err
	}

	return &record, nil
}

// Delete hard-deletes a car_ins_sheet_ichiban_cars record (no soft delete for this table)
func (r *CarInsSheetIchibanCarsRepository) Delete(ctx context.Context, organizationID, electCertMgNo, electCertPublishdateE, electCertPublishdateY, electCertPublishdateM, electCertPublishdateD string) error {
	query := `
		DELETE FROM car_ins_sheet_ichiban_cars
		WHERE organization_id = $1 AND "ElectCertMgNo" = $2 AND "ElectCertPublishdateE" = $3 AND "ElectCertPublishdateY" = $4 AND "ElectCertPublishdateM" = $5 AND "ElectCertPublishdateD" = $6
	`

	result, err := r.db.Exec(ctx, query, organizationID, electCertMgNo, electCertPublishdateE, electCertPublishdateY, electCertPublishdateM, electCertPublishdateD)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrCarInsSheetIchibanCarsNotFound
	}

	return nil
}

// ListByOrganization retrieves car_ins_sheet_ichiban_cars records by organization with pagination
func (r *CarInsSheetIchibanCarsRepository) ListByOrganization(ctx context.Context, organizationID string, limit int, offset int) ([]*CarInsSheetIchibanCars, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	query := `
		SELECT organization_id, id_cars, "ElectCertMgNo", "ElectCertPublishdateE", "ElectCertPublishdateY", "ElectCertPublishdateM", "ElectCertPublishdateD"
		FROM car_ins_sheet_ichiban_cars
		WHERE organization_id = $1
		ORDER BY "ElectCertMgNo", "ElectCertPublishdateE", "ElectCertPublishdateY", "ElectCertPublishdateM", "ElectCertPublishdateD"
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, organizationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*CarInsSheetIchibanCars
	for rows.Next() {
		var record CarInsSheetIchibanCars
		err := rows.Scan(&record.OrganizationID, &record.IDCars, &record.ElectCertMgNo, &record.ElectCertPublishdateE, &record.ElectCertPublishdateY, &record.ElectCertPublishdateM, &record.ElectCertPublishdateD)
		if err != nil {
			return nil, err
		}
		records = append(records, &record)
	}

	return records, rows.Err()
}

// List retrieves all car_ins_sheet_ichiban_cars records with pagination
func (r *CarInsSheetIchibanCarsRepository) List(ctx context.Context, limit int, offset int) ([]*CarInsSheetIchibanCars, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	query := `
		SELECT organization_id, id_cars, "ElectCertMgNo", "ElectCertPublishdateE", "ElectCertPublishdateY", "ElectCertPublishdateM", "ElectCertPublishdateD"
		FROM car_ins_sheet_ichiban_cars
		ORDER BY organization_id, "ElectCertMgNo", "ElectCertPublishdateE", "ElectCertPublishdateY", "ElectCertPublishdateM", "ElectCertPublishdateD"
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*CarInsSheetIchibanCars
	for rows.Next() {
		var record CarInsSheetIchibanCars
		err := rows.Scan(&record.OrganizationID, &record.IDCars, &record.ElectCertMgNo, &record.ElectCertPublishdateE, &record.ElectCertPublishdateY, &record.ElectCertPublishdateM, &record.ElectCertPublishdateD)
		if err != nil {
			return nil, err
		}
		records = append(records, &record)
	}

	return records, rows.Err()
}
