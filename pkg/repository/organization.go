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
	ErrNotFound = errors.New("organization not found")
)

// Organization represents the database model
type Organization struct {
	ID        string
	Name      string
	Slug      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// OrganizationRepository handles database operations for organizations
type OrganizationRepository struct {
	db DB
}

// NewOrganizationRepository creates a new repository
func NewOrganizationRepository(pool *pgxpool.Pool) *OrganizationRepository {
	return &OrganizationRepository{db: pool}
}

// NewOrganizationRepositoryWithDB creates a repository with custom DB interface (for testing)
func NewOrganizationRepositoryWithDB(db DB) *OrganizationRepository {
	return &OrganizationRepository{db: db}
}

// Create inserts a new organization
func (r *OrganizationRepository) Create(ctx context.Context, name, slug string) (*Organization, error) {
	id := uuid.New().String()
	now := time.Now()

	query := `
		INSERT INTO organizations (id, name, slug, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, slug, created_at, updated_at, deleted_at
	`

	var org Organization
	err := r.db.QueryRow(ctx, query, id, name, slug, now, now).Scan(
		&org.ID, &org.Name, &org.Slug, &org.CreatedAt, &org.UpdatedAt, &org.DeletedAt,
	)
	if err != nil {
		return nil, err
	}

	return &org, nil
}

// GetByID retrieves an organization by ID
func (r *OrganizationRepository) GetByID(ctx context.Context, id string) (*Organization, error) {
	query := `
		SELECT id, name, slug, created_at, updated_at, deleted_at
		FROM organizations
		WHERE id = $1 AND deleted_at IS NULL
	`

	var org Organization
	err := r.db.QueryRow(ctx, query, id).Scan(
		&org.ID, &org.Name, &org.Slug, &org.CreatedAt, &org.UpdatedAt, &org.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &org, nil
}

// Update modifies an existing organization
func (r *OrganizationRepository) Update(ctx context.Context, id, name, slug string) (*Organization, error) {
	query := `
		UPDATE organizations
		SET name = $2, slug = $3, updated_at = $4
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING id, name, slug, created_at, updated_at, deleted_at
	`

	var org Organization
	err := r.db.QueryRow(ctx, query, id, name, slug, time.Now()).Scan(
		&org.ID, &org.Name, &org.Slug, &org.CreatedAt, &org.UpdatedAt, &org.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &org, nil
}

// Delete soft-deletes an organization
func (r *OrganizationRepository) Delete(ctx context.Context, id string) error {
	query := `
		UPDATE organizations
		SET deleted_at = $2, updated_at = $2
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.Exec(ctx, query, id, time.Now())
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

// List retrieves organizations with pagination
func (r *OrganizationRepository) List(ctx context.Context, limit int, offset int) ([]*Organization, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	query := `
		SELECT id, name, slug, created_at, updated_at, deleted_at
		FROM organizations
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orgs []*Organization
	for rows.Next() {
		var org Organization
		err := rows.Scan(&org.ID, &org.Name, &org.Slug, &org.CreatedAt, &org.UpdatedAt, &org.DeletedAt)
		if err != nil {
			return nil, err
		}
		orgs = append(orgs, &org)
	}

	return orgs, rows.Err()
}
