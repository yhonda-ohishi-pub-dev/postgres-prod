package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrCamFileExeNotFound = errors.New("cam file exe not found")
)

// CamFileExe represents the database model
type CamFileExe struct {
	Name           string
	Cam            string
	OrganizationID string
	Stage          int32
}

// CamFileExeRepository handles database operations for cam_file_exe
type CamFileExeRepository struct {
	db DB
}

// NewCamFileExeRepository creates a new repository
func NewCamFileExeRepository(pool *pgxpool.Pool) *CamFileExeRepository {
	return &CamFileExeRepository{db: pool}
}

// NewCamFileExeRepositoryWithDB creates a repository with custom DB interface (for testing)
func NewCamFileExeRepositoryWithDB(db DB) *CamFileExeRepository {
	return &CamFileExeRepository{db: db}
}

// Create inserts a new cam file exe
func (r *CamFileExeRepository) Create(ctx context.Context, name, cam, organizationID string, stage int32) (*CamFileExe, error) {
	query := `
		INSERT INTO cam_file_exe (name, cam, organization_id, stage)
		VALUES ($1, $2, $3, $4)
		RETURNING name, cam, organization_id, stage
	`

	var camFileExe CamFileExe
	err := r.db.QueryRow(ctx, query, name, cam, organizationID, stage).Scan(
		&camFileExe.Name, &camFileExe.Cam, &camFileExe.OrganizationID, &camFileExe.Stage,
	)
	if err != nil {
		return nil, err
	}

	return &camFileExe, nil
}

// GetByKey retrieves a cam file exe by composite key (name, cam, organization_id)
func (r *CamFileExeRepository) GetByKey(ctx context.Context, name, cam, organizationID string) (*CamFileExe, error) {
	query := `
		SELECT name, cam, organization_id, stage
		FROM cam_file_exe
		WHERE name = $1 AND cam = $2 AND organization_id = $3
	`

	var camFileExe CamFileExe
	err := r.db.QueryRow(ctx, query, name, cam, organizationID).Scan(
		&camFileExe.Name, &camFileExe.Cam, &camFileExe.OrganizationID, &camFileExe.Stage,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCamFileExeNotFound
		}
		return nil, err
	}

	return &camFileExe, nil
}

// Update modifies an existing cam file exe
func (r *CamFileExeRepository) Update(ctx context.Context, name, cam, organizationID string, stage int32) (*CamFileExe, error) {
	query := `
		UPDATE cam_file_exe
		SET stage = $4
		WHERE name = $1 AND cam = $2 AND organization_id = $3
		RETURNING name, cam, organization_id, stage
	`

	var camFileExe CamFileExe
	err := r.db.QueryRow(ctx, query, name, cam, organizationID, stage).Scan(
		&camFileExe.Name, &camFileExe.Cam, &camFileExe.OrganizationID, &camFileExe.Stage,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCamFileExeNotFound
		}
		return nil, err
	}

	return &camFileExe, nil
}

// Delete removes a cam file exe
func (r *CamFileExeRepository) Delete(ctx context.Context, name, cam, organizationID string) error {
	query := `
		DELETE FROM cam_file_exe
		WHERE name = $1 AND cam = $2 AND organization_id = $3
	`

	result, err := r.db.Exec(ctx, query, name, cam, organizationID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrCamFileExeNotFound
	}

	return nil
}

// ListByOrganization retrieves all cam file exe entries for an organization with pagination
func (r *CamFileExeRepository) ListByOrganization(ctx context.Context, organizationID string, limit int, offset int) ([]*CamFileExe, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	query := `
		SELECT name, cam, organization_id, stage
		FROM cam_file_exe
		WHERE organization_id = $1
		ORDER BY name, cam
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, organizationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var camFileExes []*CamFileExe
	for rows.Next() {
		var camFileExe CamFileExe
		err := rows.Scan(&camFileExe.Name, &camFileExe.Cam, &camFileExe.OrganizationID, &camFileExe.Stage)
		if err != nil {
			return nil, err
		}
		camFileExes = append(camFileExes, &camFileExe)
	}

	return camFileExes, rows.Err()
}
