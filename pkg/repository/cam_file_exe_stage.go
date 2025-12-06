package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrCamFileExeStageNotFound = errors.New("cam file exe stage not found")
)

// CamFileExeStage represents the database model
type CamFileExeStage struct {
	Stage          int32
	OrganizationID string
	Name           string
}

// CamFileExeStageRepository handles database operations for cam_file_exe_stage
type CamFileExeStageRepository struct {
	db DB
}

// NewCamFileExeStageRepository creates a new repository
func NewCamFileExeStageRepository(pool *pgxpool.Pool) *CamFileExeStageRepository {
	return &CamFileExeStageRepository{db: pool}
}

// NewCamFileExeStageRepositoryWithDB creates a repository with custom DB interface (for testing)
func NewCamFileExeStageRepositoryWithDB(db DB) *CamFileExeStageRepository {
	return &CamFileExeStageRepository{db: db}
}

// Create inserts a new cam file exe stage
func (r *CamFileExeStageRepository) Create(ctx context.Context, stage int32, organizationID, name string) (*CamFileExeStage, error) {
	query := `
		INSERT INTO cam_file_exe_stage (stage, organization_id, name)
		VALUES ($1, $2, $3)
		RETURNING stage, organization_id, name
	`

	var camStage CamFileExeStage
	err := r.db.QueryRow(ctx, query, stage, organizationID, name).Scan(
		&camStage.Stage, &camStage.OrganizationID, &camStage.Name,
	)
	if err != nil {
		return nil, err
	}

	return &camStage, nil
}

// GetByStageAndOrg retrieves a cam file exe stage by composite key (stage, organization_id)
func (r *CamFileExeStageRepository) GetByStageAndOrg(ctx context.Context, stage int32, organizationID string) (*CamFileExeStage, error) {
	query := `
		SELECT stage, organization_id, name
		FROM cam_file_exe_stage
		WHERE stage = $1 AND organization_id = $2
	`

	var camStage CamFileExeStage
	err := r.db.QueryRow(ctx, query, stage, organizationID).Scan(
		&camStage.Stage, &camStage.OrganizationID, &camStage.Name,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCamFileExeStageNotFound
		}
		return nil, err
	}

	return &camStage, nil
}

// Update modifies an existing cam file exe stage
func (r *CamFileExeStageRepository) Update(ctx context.Context, stage int32, organizationID, name string) (*CamFileExeStage, error) {
	query := `
		UPDATE cam_file_exe_stage
		SET name = $3
		WHERE stage = $1 AND organization_id = $2
		RETURNING stage, organization_id, name
	`

	var camStage CamFileExeStage
	err := r.db.QueryRow(ctx, query, stage, organizationID, name).Scan(
		&camStage.Stage, &camStage.OrganizationID, &camStage.Name,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCamFileExeStageNotFound
		}
		return nil, err
	}

	return &camStage, nil
}

// Delete removes a cam file exe stage
func (r *CamFileExeStageRepository) Delete(ctx context.Context, stage int32, organizationID string) error {
	query := `
		DELETE FROM cam_file_exe_stage
		WHERE stage = $1 AND organization_id = $2
	`

	result, err := r.db.Exec(ctx, query, stage, organizationID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrCamFileExeStageNotFound
	}

	return nil
}

// ListByOrganization retrieves all cam file exe stages for an organization
func (r *CamFileExeStageRepository) ListByOrganization(ctx context.Context, organizationID string) ([]*CamFileExeStage, error) {
	query := `
		SELECT stage, organization_id, name
		FROM cam_file_exe_stage
		WHERE organization_id = $1
		ORDER BY stage ASC
	`

	rows, err := r.db.Query(ctx, query, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stages []*CamFileExeStage
	for rows.Next() {
		var stage CamFileExeStage
		err := rows.Scan(&stage.Stage, &stage.OrganizationID, &stage.Name)
		if err != nil {
			return nil, err
		}
		stages = append(stages, &stage)
	}

	return stages, rows.Err()
}
