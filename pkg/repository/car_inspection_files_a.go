package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrCarInspectionFilesANotFound = errors.New("car inspection file A not found")
)

// CarInspectionFilesA represents the database model
type CarInspectionFilesA struct {
	UUID           string
	OrganizationID string
	Type           string
	ElectCertMgNo  string
	GrantdateE     string
	GrantdateY     string
	GrantdateM     string
	GrantdateD     string
	Created        string
	Modified       string
	Deleted        *string
}

// CarInspectionFilesARepository handles database operations for car_inspection_files_a
type CarInspectionFilesARepository struct {
	db DB
}

// NewCarInspectionFilesARepository creates a new repository
func NewCarInspectionFilesARepository(pool *pgxpool.Pool) *CarInspectionFilesARepository {
	return &CarInspectionFilesARepository{db: pool}
}

// NewCarInspectionFilesARepositoryWithDB creates a repository with custom DB interface (for testing)
func NewCarInspectionFilesARepositoryWithDB(db DB) *CarInspectionFilesARepository {
	return &CarInspectionFilesARepository{db: db}
}

// Create inserts a new car inspection file A record
func (r *CarInspectionFilesARepository) Create(ctx context.Context, record *CarInspectionFilesA) (*CarInspectionFilesA, error) {
	id := uuid.New().String()

	query := `
		INSERT INTO car_inspection_files_a (
			uuid, organization_id, type, "ElectCertMgNo", "GrantdateE",
			"GrantdateY", "GrantdateM", "GrantdateD", created, modified
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING uuid, organization_id, type, "ElectCertMgNo", "GrantdateE",
			"GrantdateY", "GrantdateM", "GrantdateD", created, modified, deleted
	`

	var result CarInspectionFilesA
	err := r.db.QueryRow(
		ctx, query,
		id,
		record.OrganizationID,
		record.Type,
		record.ElectCertMgNo,
		record.GrantdateE,
		record.GrantdateY,
		record.GrantdateM,
		record.GrantdateD,
		record.Created,
		record.Modified,
	).Scan(
		&result.UUID,
		&result.OrganizationID,
		&result.Type,
		&result.ElectCertMgNo,
		&result.GrantdateE,
		&result.GrantdateY,
		&result.GrantdateM,
		&result.GrantdateD,
		&result.Created,
		&result.Modified,
		&result.Deleted,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetByUUID retrieves a car inspection file A record by UUID
func (r *CarInspectionFilesARepository) GetByUUID(ctx context.Context, uuid string) (*CarInspectionFilesA, error) {
	query := `
		SELECT uuid, organization_id, type, "ElectCertMgNo", "GrantdateE",
			"GrantdateY", "GrantdateM", "GrantdateD", created, modified, deleted
		FROM car_inspection_files_a
		WHERE uuid = $1 AND deleted IS NULL
	`

	var record CarInspectionFilesA
	err := r.db.QueryRow(ctx, query, uuid).Scan(
		&record.UUID,
		&record.OrganizationID,
		&record.Type,
		&record.ElectCertMgNo,
		&record.GrantdateE,
		&record.GrantdateY,
		&record.GrantdateM,
		&record.GrantdateD,
		&record.Created,
		&record.Modified,
		&record.Deleted,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCarInspectionFilesANotFound
		}
		return nil, err
	}

	return &record, nil
}

// Update modifies an existing car inspection file A record
func (r *CarInspectionFilesARepository) Update(ctx context.Context, record *CarInspectionFilesA) (*CarInspectionFilesA, error) {
	query := `
		UPDATE car_inspection_files_a
		SET organization_id = $2, type = $3, "ElectCertMgNo" = $4, "GrantdateE" = $5,
			"GrantdateY" = $6, "GrantdateM" = $7, "GrantdateD" = $8, modified = $9
		WHERE uuid = $1 AND deleted IS NULL
		RETURNING uuid, organization_id, type, "ElectCertMgNo", "GrantdateE",
			"GrantdateY", "GrantdateM", "GrantdateD", created, modified, deleted
	`

	var result CarInspectionFilesA
	err := r.db.QueryRow(
		ctx, query,
		record.UUID,
		record.OrganizationID,
		record.Type,
		record.ElectCertMgNo,
		record.GrantdateE,
		record.GrantdateY,
		record.GrantdateM,
		record.GrantdateD,
		record.Modified,
	).Scan(
		&result.UUID,
		&result.OrganizationID,
		&result.Type,
		&result.ElectCertMgNo,
		&result.GrantdateE,
		&result.GrantdateY,
		&result.GrantdateM,
		&result.GrantdateD,
		&result.Created,
		&result.Modified,
		&result.Deleted,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCarInspectionFilesANotFound
		}
		return nil, err
	}

	return &result, nil
}

// Delete soft-deletes a car inspection file A record
func (r *CarInspectionFilesARepository) Delete(ctx context.Context, uuid string, deletedTime string) error {
	query := `
		UPDATE car_inspection_files_a
		SET deleted = $2, modified = $2
		WHERE uuid = $1 AND deleted IS NULL
	`

	result, err := r.db.Exec(ctx, query, uuid, deletedTime)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrCarInspectionFilesANotFound
	}

	return nil
}

// ListByOrganization retrieves car inspection file A records for a specific organization
func (r *CarInspectionFilesARepository) ListByOrganization(ctx context.Context, organizationID string, limit int, offset int) ([]*CarInspectionFilesA, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	query := `
		SELECT uuid, organization_id, type, "ElectCertMgNo", "GrantdateE",
			"GrantdateY", "GrantdateM", "GrantdateD", created, modified, deleted
		FROM car_inspection_files_a
		WHERE organization_id = $1 AND deleted IS NULL
		ORDER BY created DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, organizationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*CarInspectionFilesA
	for rows.Next() {
		var record CarInspectionFilesA
		err := rows.Scan(
			&record.UUID,
			&record.OrganizationID,
			&record.Type,
			&record.ElectCertMgNo,
			&record.GrantdateE,
			&record.GrantdateY,
			&record.GrantdateM,
			&record.GrantdateD,
			&record.Created,
			&record.Modified,
			&record.Deleted,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, &record)
	}

	return records, rows.Err()
}

// List retrieves car inspection file A records with pagination
func (r *CarInspectionFilesARepository) List(ctx context.Context, limit int, offset int) ([]*CarInspectionFilesA, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	query := `
		SELECT uuid, organization_id, type, "ElectCertMgNo", "GrantdateE",
			"GrantdateY", "GrantdateM", "GrantdateD", created, modified, deleted
		FROM car_inspection_files_a
		WHERE deleted IS NULL
		ORDER BY created DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*CarInspectionFilesA
	for rows.Next() {
		var record CarInspectionFilesA
		err := rows.Scan(
			&record.UUID,
			&record.OrganizationID,
			&record.Type,
			&record.ElectCertMgNo,
			&record.GrantdateE,
			&record.GrantdateY,
			&record.GrantdateM,
			&record.GrantdateD,
			&record.Created,
			&record.Modified,
			&record.Deleted,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, &record)
	}

	return records, rows.Err()
}
