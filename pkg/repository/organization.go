package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/db"
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
	db      DB
	rlsPool *db.RLSPool
}

// NewOrganizationRepository creates a new repository
func NewOrganizationRepository(pool *pgxpool.Pool) *OrganizationRepository {
	return &OrganizationRepository{db: pool}
}

// NewOrganizationRepositoryWithDB creates a repository with custom DB interface (for testing)
func NewOrganizationRepositoryWithDB(d DB) *OrganizationRepository {
	// Check if the DB is an RLSPool to enable transaction support
	if rlsPool, ok := d.(*db.RLSPool); ok {
		return &OrganizationRepository{db: d, rlsPool: rlsPool}
	}
	return &OrganizationRepository{db: d}
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

// CreateWithOwnerResult contains both the organization and user_organization created
type CreateWithOwnerResult struct {
	Organization     *Organization
	UserOrganization *UserOrganization
}

// CreateWithOwner creates an organization and links it to the user as owner in a single transaction.
// The user is assigned the "owner" role and this organization is set as their default.
func (r *OrganizationRepository) CreateWithOwner(ctx context.Context, name, slug, userID string) (*CreateWithOwnerResult, error) {
	if r.rlsPool == nil {
		return nil, errors.New("transaction support requires RLSPool")
	}

	tx, err := r.rlsPool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	now := time.Now()
	orgID := uuid.New().String()
	userOrgID := uuid.New().String()

	// Create organization
	orgQuery := `
		INSERT INTO organizations (id, name, slug, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, slug, created_at, updated_at, deleted_at
	`
	var org Organization
	err = tx.QueryRow(ctx, orgQuery, orgID, name, slug, now, now).Scan(
		&org.ID, &org.Name, &org.Slug, &org.CreatedAt, &org.UpdatedAt, &org.DeletedAt,
	)
	if err != nil {
		return nil, err
	}

	// Create user_organization link with owner role
	userOrgQuery := `
		INSERT INTO user_organizations (id, user_id, organization_id, role, is_default, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, user_id, organization_id, role, is_default, created_at, updated_at
	`
	var userOrg UserOrganization
	err = tx.QueryRow(ctx, userOrgQuery, userOrgID, userID, orgID, "owner", true, now, now).Scan(
		&userOrg.ID, &userOrg.UserID, &userOrg.OrganizationID, &userOrg.Role, &userOrg.IsDefault, &userOrg.CreatedAt, &userOrg.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &CreateWithOwnerResult{
		Organization:     &org,
		UserOrganization: &userOrg,
	}, nil
}
