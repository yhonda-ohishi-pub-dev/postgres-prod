package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrCarInspectionFilesBNotFound = errors.New("car inspection files b not found")
)

// CarInspectionFilesB represents the database model
type CarInspectionFilesB struct {
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

// CarInspectionFilesBRepository handles database operations for car_inspection_files_b
type CarInspectionFilesBRepository struct {
	db DB
}

// NewCarInspectionFilesBRepository creates a new repository
func NewCarInspectionFilesBRepository(pool *pgxpool.Pool) *CarInspectionFilesBRepository {
	return &CarInspectionFilesBRepository{db: pool}
}

// NewCarInspectionFilesBRepositoryWithDB creates a repository with custom DB interface (for testing)
func NewCarInspectionFilesBRepositoryWithDB(db DB) *CarInspectionFilesBRepository {
	return &CarInspectionFilesBRepository{db: db}
}

// Create inserts a new car inspection file b record
func (r *CarInspectionFilesBRepository) Create(ctx context.Context, organizationID, typeVal, electCertMgNo, grantdateE, grantdateY, grantdateM, grantdateD string) (*CarInspectionFilesB, error) {
	id := uuid.New().String()
	now := time.Now().Format(time.RFC3339)

	query := `
		INSERT INTO car_inspection_files_b (uuid, organization_id, type, "ElectCertMgNo", "GrantdateE", "GrantdateY", "GrantdateM", "GrantdateD", created, modified)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING uuid, organization_id, type, "ElectCertMgNo", "GrantdateE", "GrantdateY", "GrantdateM", "GrantdateD", created, modified, deleted
	`

	var record CarInspectionFilesB
	err := r.db.QueryRow(ctx, query, id, organizationID, typeVal, electCertMgNo, grantdateE, grantdateY, grantdateM, grantdateD, now, now).Scan(
		&record.UUID, &record.OrganizationID, &record.Type, &record.ElectCertMgNo,
		&record.GrantdateE, &record.GrantdateY, &record.GrantdateM, &record.GrantdateD,
		&record.Created, &record.Modified, &record.Deleted,
	)
	if err != nil {
		return nil, err
	}

	return &record, nil
}

// GetByUUID retrieves a car inspection file b record by UUID
func (r *CarInspectionFilesBRepository) GetByUUID(ctx context.Context, uuid string) (*CarInspectionFilesB, error) {
	query := `
		SELECT uuid, organization_id, type, "ElectCertMgNo", "GrantdateE", "GrantdateY", "GrantdateM", "GrantdateD", created, modified, deleted
		FROM car_inspection_files_b
		WHERE uuid = $1 AND deleted IS NULL
	`

	var record CarInspectionFilesB
	err := r.db.QueryRow(ctx, query, uuid).Scan(
		&record.UUID, &record.OrganizationID, &record.Type, &record.ElectCertMgNo,
		&record.GrantdateE, &record.GrantdateY, &record.GrantdateM, &record.GrantdateD,
		&record.Created, &record.Modified, &record.Deleted,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCarInspectionFilesBNotFound
		}
		return nil, err
	}

	return &record, nil
}

// Update modifies an existing car inspection file b record
func (r *CarInspectionFilesBRepository) Update(ctx context.Context, uuid, organizationID, typeVal, electCertMgNo, grantdateE, grantdateY, grantdateM, grantdateD string) (*CarInspectionFilesB, error) {
	now := time.Now().Format(time.RFC3339)

	query := `
		UPDATE car_inspection_files_b
		SET organization_id = $2, type = $3, "ElectCertMgNo" = $4, "GrantdateE" = $5, "GrantdateY" = $6, "GrantdateM" = $7, "GrantdateD" = $8, modified = $9
		WHERE uuid = $1 AND deleted IS NULL
		RETURNING uuid, organization_id, type, "ElectCertMgNo", "GrantdateE", "GrantdateY", "GrantdateM", "GrantdateD", created, modified, deleted
	`

	var record CarInspectionFilesB
	err := r.db.QueryRow(ctx, query, uuid, organizationID, typeVal, electCertMgNo, grantdateE, grantdateY, grantdateM, grantdateD, now).Scan(
		&record.UUID, &record.OrganizationID, &record.Type, &record.ElectCertMgNo,
		&record.GrantdateE, &record.GrantdateY, &record.GrantdateM, &record.GrantdateD,
		&record.Created, &record.Modified, &record.Deleted,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCarInspectionFilesBNotFound
		}
		return nil, err
	}

	return &record, nil
}

// Delete soft-deletes a car inspection file b record
func (r *CarInspectionFilesBRepository) Delete(ctx context.Context, uuid string) error {
	now := time.Now().Format(time.RFC3339)

	query := `
		UPDATE car_inspection_files_b
		SET deleted = $2, modified = $2
		WHERE uuid = $1 AND deleted IS NULL
	`

	result, err := r.db.Exec(ctx, query, uuid, now)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrCarInspectionFilesBNotFound
	}

	return nil
}

// ListByOrganization retrieves car inspection files b records by organization with pagination
func (r *CarInspectionFilesBRepository) ListByOrganization(ctx context.Context, organizationID string, limit int, offset int) ([]*CarInspectionFilesB, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	query := `
		SELECT uuid, organization_id, type, "ElectCertMgNo", "GrantdateE", "GrantdateY", "GrantdateM", "GrantdateD", created, modified, deleted
		FROM car_inspection_files_b
		WHERE organization_id = $1 AND deleted IS NULL
		ORDER BY created DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, organizationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*CarInspectionFilesB
	for rows.Next() {
		var record CarInspectionFilesB
		err := rows.Scan(
			&record.UUID, &record.OrganizationID, &record.Type, &record.ElectCertMgNo,
			&record.GrantdateE, &record.GrantdateY, &record.GrantdateM, &record.GrantdateD,
			&record.Created, &record.Modified, &record.Deleted,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, &record)
	}

	return records, rows.Err()
}

// List retrieves car inspection files b records with pagination
func (r *CarInspectionFilesBRepository) List(ctx context.Context, limit int, offset int) ([]*CarInspectionFilesB, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	query := `
		SELECT uuid, organization_id, type, "ElectCertMgNo", "GrantdateE", "GrantdateY", "GrantdateM", "GrantdateD", created, modified, deleted
		FROM car_inspection_files_b
		WHERE deleted IS NULL
		ORDER BY created DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*CarInspectionFilesB
	for rows.Next() {
		var record CarInspectionFilesB
		err := rows.Scan(
			&record.UUID, &record.OrganizationID, &record.Type, &record.ElectCertMgNo,
			&record.GrantdateE, &record.GrantdateY, &record.GrantdateM, &record.GrantdateD,
			&record.Created, &record.Modified, &record.Deleted,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, &record)
	}

	return records, rows.Err()
}
