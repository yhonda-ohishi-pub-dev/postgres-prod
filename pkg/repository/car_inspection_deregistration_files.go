package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrCarInspectionDeregistrationFilesNotFound = errors.New("car inspection deregistration files not found")
)

// CarInspectionDeregistrationFiles represents the database model
type CarInspectionDeregistrationFiles struct {
	OrganizationID                              string
	CarID                                       string
	TwodimensionCodeInfoValidPeriodExpirDate    string
	FileUUID                                    string
}

// CarInspectionDeregistrationFilesRepository handles database operations for car_inspection_deregistration_files
type CarInspectionDeregistrationFilesRepository struct {
	db DB
}

// NewCarInspectionDeregistrationFilesRepository creates a new repository
func NewCarInspectionDeregistrationFilesRepository(pool *pgxpool.Pool) *CarInspectionDeregistrationFilesRepository {
	return &CarInspectionDeregistrationFilesRepository{db: pool}
}

// NewCarInspectionDeregistrationFilesRepositoryWithDB creates a repository with custom DB interface (for testing)
func NewCarInspectionDeregistrationFilesRepositoryWithDB(db DB) *CarInspectionDeregistrationFilesRepository {
	return &CarInspectionDeregistrationFilesRepository{db: db}
}

// Create inserts a new car inspection deregistration file
func (r *CarInspectionDeregistrationFilesRepository) Create(ctx context.Context, organizationID, carID, twodimensionCodeInfoValidPeriodExpirDate, fileUUID string) (*CarInspectionDeregistrationFiles, error) {
	query := `
		INSERT INTO car_inspection_deregistration_files (organization_id, "CarId", "TwodimensionCodeInfoValidPeriodExpirdate", "fileUuid")
		VALUES ($1, $2, $3, $4)
		RETURNING organization_id, "CarId", "TwodimensionCodeInfoValidPeriodExpirdate", "fileUuid"
	`

	var cidf CarInspectionDeregistrationFiles
	err := r.db.QueryRow(ctx, query, organizationID, carID, twodimensionCodeInfoValidPeriodExpirDate, fileUUID).Scan(
		&cidf.OrganizationID, &cidf.CarID, &cidf.TwodimensionCodeInfoValidPeriodExpirDate, &cidf.FileUUID,
	)
	if err != nil {
		return nil, err
	}

	return &cidf, nil
}

// GetByPrimaryKey retrieves a car inspection deregistration file by composite primary key
func (r *CarInspectionDeregistrationFilesRepository) GetByPrimaryKey(ctx context.Context, organizationID, carID, twodimensionCodeInfoValidPeriodExpirDate, fileUUID string) (*CarInspectionDeregistrationFiles, error) {
	query := `
		SELECT organization_id, "CarId", "TwodimensionCodeInfoValidPeriodExpirdate", "fileUuid"
		FROM car_inspection_deregistration_files
		WHERE organization_id = $1 AND "CarId" = $2 AND "TwodimensionCodeInfoValidPeriodExpirdate" = $3 AND "fileUuid" = $4
	`

	var cidf CarInspectionDeregistrationFiles
	err := r.db.QueryRow(ctx, query, organizationID, carID, twodimensionCodeInfoValidPeriodExpirDate, fileUUID).Scan(
		&cidf.OrganizationID, &cidf.CarID, &cidf.TwodimensionCodeInfoValidPeriodExpirDate, &cidf.FileUUID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCarInspectionDeregistrationFilesNotFound
		}
		return nil, err
	}

	return &cidf, nil
}

// Delete removes a car inspection deregistration file
func (r *CarInspectionDeregistrationFilesRepository) Delete(ctx context.Context, organizationID, carID, twodimensionCodeInfoValidPeriodExpirDate, fileUUID string) error {
	query := `
		DELETE FROM car_inspection_deregistration_files
		WHERE organization_id = $1 AND "CarId" = $2 AND "TwodimensionCodeInfoValidPeriodExpirdate" = $3 AND "fileUuid" = $4
	`

	result, err := r.db.Exec(ctx, query, organizationID, carID, twodimensionCodeInfoValidPeriodExpirDate, fileUUID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrCarInspectionDeregistrationFilesNotFound
	}

	return nil
}

// ListByOrganization retrieves all car inspection deregistration files for an organization with pagination
func (r *CarInspectionDeregistrationFilesRepository) ListByOrganization(ctx context.Context, organizationID string, limit int, offset int) ([]*CarInspectionDeregistrationFiles, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	query := `
		SELECT organization_id, "CarId", "TwodimensionCodeInfoValidPeriodExpirdate", "fileUuid"
		FROM car_inspection_deregistration_files
		WHERE organization_id = $1
		ORDER BY "CarId", "TwodimensionCodeInfoValidPeriodExpirdate", "fileUuid"
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, organizationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []*CarInspectionDeregistrationFiles
	for rows.Next() {
		var cidf CarInspectionDeregistrationFiles
		err := rows.Scan(&cidf.OrganizationID, &cidf.CarID, &cidf.TwodimensionCodeInfoValidPeriodExpirDate, &cidf.FileUUID)
		if err != nil {
			return nil, err
		}
		files = append(files, &cidf)
	}

	return files, rows.Err()
}

// ListByCarInspectionDeregistration retrieves all files for a specific car inspection deregistration with pagination
func (r *CarInspectionDeregistrationFilesRepository) ListByCarInspectionDeregistration(ctx context.Context, organizationID, carID, twodimensionCodeInfoValidPeriodExpirDate string, limit int, offset int) ([]*CarInspectionDeregistrationFiles, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	query := `
		SELECT organization_id, "CarId", "TwodimensionCodeInfoValidPeriodExpirdate", "fileUuid"
		FROM car_inspection_deregistration_files
		WHERE organization_id = $1 AND "CarId" = $2 AND "TwodimensionCodeInfoValidPeriodExpirdate" = $3
		ORDER BY "fileUuid"
		LIMIT $4 OFFSET $5
	`

	rows, err := r.db.Query(ctx, query, organizationID, carID, twodimensionCodeInfoValidPeriodExpirDate, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []*CarInspectionDeregistrationFiles
	for rows.Next() {
		var cidf CarInspectionDeregistrationFiles
		err := rows.Scan(&cidf.OrganizationID, &cidf.CarID, &cidf.TwodimensionCodeInfoValidPeriodExpirDate, &cidf.FileUUID)
		if err != nil {
			return nil, err
		}
		files = append(files, &cidf)
	}

	return files, rows.Err()
}

// List retrieves all car inspection deregistration files with pagination
func (r *CarInspectionDeregistrationFilesRepository) List(ctx context.Context, limit int, offset int) ([]*CarInspectionDeregistrationFiles, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	query := `
		SELECT organization_id, "CarId", "TwodimensionCodeInfoValidPeriodExpirdate", "fileUuid"
		FROM car_inspection_deregistration_files
		ORDER BY organization_id, "CarId", "TwodimensionCodeInfoValidPeriodExpirdate", "fileUuid"
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []*CarInspectionDeregistrationFiles
	for rows.Next() {
		var cidf CarInspectionDeregistrationFiles
		err := rows.Scan(&cidf.OrganizationID, &cidf.CarID, &cidf.TwodimensionCodeInfoValidPeriodExpirDate, &cidf.FileUUID)
		if err != nil {
			return nil, err
		}
		files = append(files, &cidf)
	}

	return files, rows.Err()
}
