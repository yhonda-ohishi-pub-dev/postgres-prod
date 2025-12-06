package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrCamFileNotFound = errors.New("cam file not found")
)

// CamFile represents the database model
type CamFile struct {
	Name           string
	OrganizationID string
	Date           string
	Hour           string
	Type           string
	Cam            string
	FlickrID       *string
}

// CamFileRepository handles database operations for cam_files
type CamFileRepository struct {
	db DB
}

// NewCamFileRepository creates a new repository
func NewCamFileRepository(pool *pgxpool.Pool) *CamFileRepository {
	return &CamFileRepository{db: pool}
}

// NewCamFileRepositoryWithDB creates a repository with custom DB interface (for testing)
func NewCamFileRepositoryWithDB(db DB) *CamFileRepository {
	return &CamFileRepository{db: db}
}

// Create inserts a new cam file
func (r *CamFileRepository) Create(ctx context.Context, name, organizationID, date, hour, fileType, cam string, flickrID *string) (*CamFile, error) {
	query := `
		INSERT INTO cam_files (name, organization_id, date, hour, type, cam, flickr_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING name, organization_id, date, hour, type, cam, flickr_id
	`

	var file CamFile
	err := r.db.QueryRow(ctx, query, name, organizationID, date, hour, fileType, cam, flickrID).Scan(
		&file.Name, &file.OrganizationID, &file.Date, &file.Hour, &file.Type, &file.Cam, &file.FlickrID,
	)
	if err != nil {
		return nil, err
	}

	return &file, nil
}

// GetByNameAndOrg retrieves a cam file by composite primary key (name, organization_id)
func (r *CamFileRepository) GetByNameAndOrg(ctx context.Context, name, organizationID string) (*CamFile, error) {
	query := `
		SELECT name, organization_id, date, hour, type, cam, flickr_id
		FROM cam_files
		WHERE name = $1 AND organization_id = $2
	`

	var file CamFile
	err := r.db.QueryRow(ctx, query, name, organizationID).Scan(
		&file.Name, &file.OrganizationID, &file.Date, &file.Hour, &file.Type, &file.Cam, &file.FlickrID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCamFileNotFound
		}
		return nil, err
	}

	return &file, nil
}

// Update modifies an existing cam file
func (r *CamFileRepository) Update(ctx context.Context, name, organizationID, date, hour, fileType, cam string, flickrID *string) (*CamFile, error) {
	query := `
		UPDATE cam_files
		SET date = $3, hour = $4, type = $5, cam = $6, flickr_id = $7
		WHERE name = $1 AND organization_id = $2
		RETURNING name, organization_id, date, hour, type, cam, flickr_id
	`

	var file CamFile
	err := r.db.QueryRow(ctx, query, name, organizationID, date, hour, fileType, cam, flickrID).Scan(
		&file.Name, &file.OrganizationID, &file.Date, &file.Hour, &file.Type, &file.Cam, &file.FlickrID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCamFileNotFound
		}
		return nil, err
	}

	return &file, nil
}

// Delete removes a cam file (hard delete)
func (r *CamFileRepository) Delete(ctx context.Context, name, organizationID string) error {
	query := `
		DELETE FROM cam_files
		WHERE name = $1 AND organization_id = $2
	`

	result, err := r.db.Exec(ctx, query, name, organizationID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrCamFileNotFound
	}

	return nil
}

// ListByOrganization retrieves cam files for a specific organization with pagination
func (r *CamFileRepository) ListByOrganization(ctx context.Context, organizationID string, limit int, offset int) ([]*CamFile, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	query := `
		SELECT name, organization_id, date, hour, type, cam, flickr_id
		FROM cam_files
		WHERE organization_id = $1
		ORDER BY date DESC, hour DESC, name ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, organizationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []*CamFile
	for rows.Next() {
		var file CamFile
		err := rows.Scan(&file.Name, &file.OrganizationID, &file.Date, &file.Hour, &file.Type, &file.Cam, &file.FlickrID)
		if err != nil {
			return nil, err
		}
		files = append(files, &file)
	}

	return files, rows.Err()
}
