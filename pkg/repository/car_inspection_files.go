package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrCarInspectionFileNotFound = errors.New("car inspection file not found")
)

// CarInspectionFile represents the database model
type CarInspectionFile struct {
	UUID                   string
	OrganizationID         string
	Type                   string
	ElectCertMgNo          string
	ElectCertPublishdateE  string
	ElectCertPublishdateY  string
	ElectCertPublishdateM  string
	ElectCertPublishdateD  string
	Created                string
	Modified               string
	Deleted                *string
}

// CarInspectionFilesRepository handles database operations for car inspection files
type CarInspectionFilesRepository struct {
	db DB
}

// NewCarInspectionFilesRepository creates a new repository
func NewCarInspectionFilesRepository(pool *pgxpool.Pool) *CarInspectionFilesRepository {
	return &CarInspectionFilesRepository{db: pool}
}

// NewCarInspectionFilesRepositoryWithDB creates a repository with custom DB interface (for testing)
func NewCarInspectionFilesRepositoryWithDB(db DB) *CarInspectionFilesRepository {
	return &CarInspectionFilesRepository{db: db}
}

// Create inserts a new car inspection file
func (r *CarInspectionFilesRepository) Create(ctx context.Context, organizationID, fileType, electCertMgNo, electCertPublishdateE, electCertPublishdateY, electCertPublishdateM, electCertPublishdateD, created, modified string) (*CarInspectionFile, error) {
	id := uuid.New().String()

	query := `
		INSERT INTO car_inspection_files (uuid, organization_id, type, "ElectCertMgNo", "ElectCertPublishdateE", "ElectCertPublishdateY", "ElectCertPublishdateM", "ElectCertPublishdateD", created, modified, deleted)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NULL)
		RETURNING uuid, organization_id, type, "ElectCertMgNo", "ElectCertPublishdateE", "ElectCertPublishdateY", "ElectCertPublishdateM", "ElectCertPublishdateD", created, modified, deleted
	`

	var f CarInspectionFile
	err := r.db.QueryRow(ctx, query, id, organizationID, fileType, electCertMgNo, electCertPublishdateE, electCertPublishdateY, electCertPublishdateM, electCertPublishdateD, created, modified).Scan(
		&f.UUID, &f.OrganizationID, &f.Type, &f.ElectCertMgNo, &f.ElectCertPublishdateE, &f.ElectCertPublishdateY, &f.ElectCertPublishdateM, &f.ElectCertPublishdateD, &f.Created, &f.Modified, &f.Deleted,
	)
	if err != nil {
		return nil, err
	}

	return &f, nil
}

// GetByUUID retrieves a car inspection file by UUID
func (r *CarInspectionFilesRepository) GetByUUID(ctx context.Context, uuid string) (*CarInspectionFile, error) {
	query := `
		SELECT uuid, organization_id, type, "ElectCertMgNo", "ElectCertPublishdateE", "ElectCertPublishdateY", "ElectCertPublishdateM", "ElectCertPublishdateD", created, modified, deleted
		FROM car_inspection_files
		WHERE uuid = $1 AND deleted IS NULL
	`

	var f CarInspectionFile
	err := r.db.QueryRow(ctx, query, uuid).Scan(
		&f.UUID, &f.OrganizationID, &f.Type, &f.ElectCertMgNo, &f.ElectCertPublishdateE, &f.ElectCertPublishdateY, &f.ElectCertPublishdateM, &f.ElectCertPublishdateD, &f.Created, &f.Modified, &f.Deleted,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCarInspectionFileNotFound
		}
		return nil, err
	}

	return &f, nil
}

// Update modifies an existing car inspection file
func (r *CarInspectionFilesRepository) Update(ctx context.Context, uuid, fileType, electCertMgNo, electCertPublishdateE, electCertPublishdateY, electCertPublishdateM, electCertPublishdateD, modified string) (*CarInspectionFile, error) {
	query := `
		UPDATE car_inspection_files
		SET type = $2, "ElectCertMgNo" = $3, "ElectCertPublishdateE" = $4, "ElectCertPublishdateY" = $5, "ElectCertPublishdateM" = $6, "ElectCertPublishdateD" = $7, modified = $8
		WHERE uuid = $1 AND deleted IS NULL
		RETURNING uuid, organization_id, type, "ElectCertMgNo", "ElectCertPublishdateE", "ElectCertPublishdateY", "ElectCertPublishdateM", "ElectCertPublishdateD", created, modified, deleted
	`

	var f CarInspectionFile
	err := r.db.QueryRow(ctx, query, uuid, fileType, electCertMgNo, electCertPublishdateE, electCertPublishdateY, electCertPublishdateM, electCertPublishdateD, modified).Scan(
		&f.UUID, &f.OrganizationID, &f.Type, &f.ElectCertMgNo, &f.ElectCertPublishdateE, &f.ElectCertPublishdateY, &f.ElectCertPublishdateM, &f.ElectCertPublishdateD, &f.Created, &f.Modified, &f.Deleted,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCarInspectionFileNotFound
		}
		return nil, err
	}

	return &f, nil
}

// Delete soft-deletes a car inspection file
func (r *CarInspectionFilesRepository) Delete(ctx context.Context, uuid, deletedTimestamp string) error {
	query := `
		UPDATE car_inspection_files
		SET deleted = $2
		WHERE uuid = $1 AND deleted IS NULL
	`

	result, err := r.db.Exec(ctx, query, uuid, deletedTimestamp)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrCarInspectionFileNotFound
	}

	return nil
}

// ListByOrganization retrieves car inspection files by organization with pagination
func (r *CarInspectionFilesRepository) ListByOrganization(ctx context.Context, organizationID string, limit int, offset int) ([]*CarInspectionFile, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	query := `
		SELECT uuid, organization_id, type, "ElectCertMgNo", "ElectCertPublishdateE", "ElectCertPublishdateY", "ElectCertPublishdateM", "ElectCertPublishdateD", created, modified, deleted
		FROM car_inspection_files
		WHERE organization_id = $1 AND deleted IS NULL
		ORDER BY created DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, organizationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []*CarInspectionFile
	for rows.Next() {
		var f CarInspectionFile
		err := rows.Scan(&f.UUID, &f.OrganizationID, &f.Type, &f.ElectCertMgNo, &f.ElectCertPublishdateE, &f.ElectCertPublishdateY, &f.ElectCertPublishdateM, &f.ElectCertPublishdateD, &f.Created, &f.Modified, &f.Deleted)
		if err != nil {
			return nil, err
		}
		files = append(files, &f)
	}

	return files, rows.Err()
}

// List retrieves all car inspection files with pagination
func (r *CarInspectionFilesRepository) List(ctx context.Context, limit int, offset int) ([]*CarInspectionFile, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	query := `
		SELECT uuid, organization_id, type, "ElectCertMgNo", "ElectCertPublishdateE", "ElectCertPublishdateY", "ElectCertPublishdateM", "ElectCertPublishdateD", created, modified, deleted
		FROM car_inspection_files
		WHERE deleted IS NULL
		ORDER BY created DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []*CarInspectionFile
	for rows.Next() {
		var f CarInspectionFile
		err := rows.Scan(&f.UUID, &f.OrganizationID, &f.Type, &f.ElectCertMgNo, &f.ElectCertPublishdateE, &f.ElectCertPublishdateY, &f.ElectCertPublishdateM, &f.ElectCertPublishdateD, &f.Created, &f.Modified, &f.Deleted)
		if err != nil {
			return nil, err
		}
		files = append(files, &f)
	}

	return files, rows.Err()
}
