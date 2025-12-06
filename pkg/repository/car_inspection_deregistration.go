package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrCarInspectionDeregistrationNotFound = errors.New("car inspection deregistration not found")
)

// CarInspectionDeregistration represents the database model
type CarInspectionDeregistration struct {
	OrganizationID                             string
	CarID                                      string
	TwodimensionCodeInfoCarNo                  string
	CarNo                                      string
	ValidPeriodExpirDateE                      string
	ValidPeriodExpirDateY                      string
	ValidPeriodExpirDateM                      string
	ValidPeriodExpirDateD                      string
	TwodimensionCodeInfoValidPeriodExpirDate   string
}

// CarInspectionDeregistrationRepository handles database operations for car_inspection_deregistration
type CarInspectionDeregistrationRepository struct {
	db DB
}

// NewCarInspectionDeregistrationRepository creates a new repository
func NewCarInspectionDeregistrationRepository(pool *pgxpool.Pool) *CarInspectionDeregistrationRepository {
	return &CarInspectionDeregistrationRepository{db: pool}
}

// NewCarInspectionDeregistrationRepositoryWithDB creates a repository with custom DB interface (for testing)
func NewCarInspectionDeregistrationRepositoryWithDB(db DB) *CarInspectionDeregistrationRepository {
	return &CarInspectionDeregistrationRepository{db: db}
}

// Create inserts a new car inspection deregistration record
func (r *CarInspectionDeregistrationRepository) Create(ctx context.Context, organizationID, carID, twodimensionCodeInfoCarNo, carNo, validPeriodExpirDateE, validPeriodExpirDateY, validPeriodExpirDateM, validPeriodExpirDateD, twodimensionCodeInfoValidPeriodExpirDate string) (*CarInspectionDeregistration, error) {
	query := `
		INSERT INTO car_inspection_deregistration (organization_id, "CarId", "TwodimensionCodeInfoCarNo", "CarNo", "ValidPeriodExpirdateE", "ValidPeriodExpirdateY", "ValidPeriodExpirdateM", "ValidPeriodExpirdateD", "TwodimensionCodeInfoValidPeriodExpirdate")
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING organization_id, "CarId", "TwodimensionCodeInfoCarNo", "CarNo", "ValidPeriodExpirdateE", "ValidPeriodExpirdateY", "ValidPeriodExpirdateM", "ValidPeriodExpirdateD", "TwodimensionCodeInfoValidPeriodExpirdate"
	`

	var record CarInspectionDeregistration
	err := r.db.QueryRow(ctx, query, organizationID, carID, twodimensionCodeInfoCarNo, carNo, validPeriodExpirDateE, validPeriodExpirDateY, validPeriodExpirDateM, validPeriodExpirDateD, twodimensionCodeInfoValidPeriodExpirDate).Scan(
		&record.OrganizationID, &record.CarID, &record.TwodimensionCodeInfoCarNo, &record.CarNo, &record.ValidPeriodExpirDateE, &record.ValidPeriodExpirDateY, &record.ValidPeriodExpirDateM, &record.ValidPeriodExpirDateD, &record.TwodimensionCodeInfoValidPeriodExpirDate,
	)
	if err != nil {
		return nil, err
	}

	return &record, nil
}

// GetByPrimaryKey retrieves a car inspection deregistration record by composite primary key (organization_id, CarId, TwodimensionCodeInfoValidPeriodExpirdate)
func (r *CarInspectionDeregistrationRepository) GetByPrimaryKey(ctx context.Context, organizationID, carID, twodimensionCodeInfoValidPeriodExpirDate string) (*CarInspectionDeregistration, error) {
	query := `
		SELECT organization_id, "CarId", "TwodimensionCodeInfoCarNo", "CarNo", "ValidPeriodExpirdateE", "ValidPeriodExpirdateY", "ValidPeriodExpirdateM", "ValidPeriodExpirdateD", "TwodimensionCodeInfoValidPeriodExpirdate"
		FROM car_inspection_deregistration
		WHERE organization_id = $1 AND "CarId" = $2 AND "TwodimensionCodeInfoValidPeriodExpirdate" = $3
	`

	var record CarInspectionDeregistration
	err := r.db.QueryRow(ctx, query, organizationID, carID, twodimensionCodeInfoValidPeriodExpirDate).Scan(
		&record.OrganizationID, &record.CarID, &record.TwodimensionCodeInfoCarNo, &record.CarNo, &record.ValidPeriodExpirDateE, &record.ValidPeriodExpirDateY, &record.ValidPeriodExpirDateM, &record.ValidPeriodExpirDateD, &record.TwodimensionCodeInfoValidPeriodExpirDate,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCarInspectionDeregistrationNotFound
		}
		return nil, err
	}

	return &record, nil
}

// Update modifies an existing car inspection deregistration record
func (r *CarInspectionDeregistrationRepository) Update(ctx context.Context, organizationID, carID, twodimensionCodeInfoCarNo, carNo, validPeriodExpirDateE, validPeriodExpirDateY, validPeriodExpirDateM, validPeriodExpirDateD, twodimensionCodeInfoValidPeriodExpirDate string) (*CarInspectionDeregistration, error) {
	query := `
		UPDATE car_inspection_deregistration
		SET "TwodimensionCodeInfoCarNo" = $3, "CarNo" = $4, "ValidPeriodExpirdateE" = $5, "ValidPeriodExpirdateY" = $6, "ValidPeriodExpirdateM" = $7, "ValidPeriodExpirdateD" = $8
		WHERE organization_id = $1 AND "CarId" = $2 AND "TwodimensionCodeInfoValidPeriodExpirdate" = $9
		RETURNING organization_id, "CarId", "TwodimensionCodeInfoCarNo", "CarNo", "ValidPeriodExpirdateE", "ValidPeriodExpirdateY", "ValidPeriodExpirdateM", "ValidPeriodExpirdateD", "TwodimensionCodeInfoValidPeriodExpirdate"
	`

	var record CarInspectionDeregistration
	err := r.db.QueryRow(ctx, query, organizationID, carID, twodimensionCodeInfoCarNo, carNo, validPeriodExpirDateE, validPeriodExpirDateY, validPeriodExpirDateM, validPeriodExpirDateD, twodimensionCodeInfoValidPeriodExpirDate).Scan(
		&record.OrganizationID, &record.CarID, &record.TwodimensionCodeInfoCarNo, &record.CarNo, &record.ValidPeriodExpirDateE, &record.ValidPeriodExpirDateY, &record.ValidPeriodExpirDateM, &record.ValidPeriodExpirDateD, &record.TwodimensionCodeInfoValidPeriodExpirDate,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCarInspectionDeregistrationNotFound
		}
		return nil, err
	}

	return &record, nil
}

// Delete hard-deletes a car inspection deregistration record (no soft delete for this table)
func (r *CarInspectionDeregistrationRepository) Delete(ctx context.Context, organizationID, carID, twodimensionCodeInfoValidPeriodExpirDate string) error {
	query := `
		DELETE FROM car_inspection_deregistration
		WHERE organization_id = $1 AND "CarId" = $2 AND "TwodimensionCodeInfoValidPeriodExpirdate" = $3
	`

	result, err := r.db.Exec(ctx, query, organizationID, carID, twodimensionCodeInfoValidPeriodExpirDate)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrCarInspectionDeregistrationNotFound
	}

	return nil
}

// ListByOrganization retrieves car inspection deregistration records by organization with pagination
func (r *CarInspectionDeregistrationRepository) ListByOrganization(ctx context.Context, organizationID string, limit int, offset int) ([]*CarInspectionDeregistration, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	query := `
		SELECT organization_id, "CarId", "TwodimensionCodeInfoCarNo", "CarNo", "ValidPeriodExpirdateE", "ValidPeriodExpirdateY", "ValidPeriodExpirdateM", "ValidPeriodExpirdateD", "TwodimensionCodeInfoValidPeriodExpirdate"
		FROM car_inspection_deregistration
		WHERE organization_id = $1
		ORDER BY "CarId", "TwodimensionCodeInfoValidPeriodExpirdate"
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, organizationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*CarInspectionDeregistration
	for rows.Next() {
		var record CarInspectionDeregistration
		err := rows.Scan(&record.OrganizationID, &record.CarID, &record.TwodimensionCodeInfoCarNo, &record.CarNo, &record.ValidPeriodExpirDateE, &record.ValidPeriodExpirDateY, &record.ValidPeriodExpirDateM, &record.ValidPeriodExpirDateD, &record.TwodimensionCodeInfoValidPeriodExpirDate)
		if err != nil {
			return nil, err
		}
		records = append(records, &record)
	}

	return records, rows.Err()
}

// List retrieves all car inspection deregistration records with pagination
func (r *CarInspectionDeregistrationRepository) List(ctx context.Context, limit int, offset int) ([]*CarInspectionDeregistration, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	query := `
		SELECT organization_id, "CarId", "TwodimensionCodeInfoCarNo", "CarNo", "ValidPeriodExpirdateE", "ValidPeriodExpirdateY", "ValidPeriodExpirdateM", "ValidPeriodExpirdateD", "TwodimensionCodeInfoValidPeriodExpirdate"
		FROM car_inspection_deregistration
		ORDER BY organization_id, "CarId", "TwodimensionCodeInfoValidPeriodExpirdate"
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*CarInspectionDeregistration
	for rows.Next() {
		var record CarInspectionDeregistration
		err := rows.Scan(&record.OrganizationID, &record.CarID, &record.TwodimensionCodeInfoCarNo, &record.CarNo, &record.ValidPeriodExpirDateE, &record.ValidPeriodExpirDateY, &record.ValidPeriodExpirDateM, &record.ValidPeriodExpirDateD, &record.TwodimensionCodeInfoValidPeriodExpirDate)
		if err != nil {
			return nil, err
		}
		records = append(records, &record)
	}

	return records, rows.Err()
}
