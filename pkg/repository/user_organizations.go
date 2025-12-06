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
	ErrUserOrganizationNotFound = errors.New("user organization not found")
)

// UserOrganization represents the database model
type UserOrganization struct {
	ID             string
	UserID         string
	OrganizationID string
	Role           string
	IsDefault      bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// UserOrganizationRepository handles database operations for user_organizations
type UserOrganizationRepository struct {
	db DB
}

// NewUserOrganizationRepository creates a new repository
func NewUserOrganizationRepository(pool *pgxpool.Pool) *UserOrganizationRepository {
	return &UserOrganizationRepository{db: pool}
}

// NewUserOrganizationRepositoryWithDB creates a repository with custom DB interface (for testing)
func NewUserOrganizationRepositoryWithDB(db DB) *UserOrganizationRepository {
	return &UserOrganizationRepository{db: db}
}

// Create inserts a new user organization
func (r *UserOrganizationRepository) Create(ctx context.Context, userID, organizationID, role string, isDefault bool) (*UserOrganization, error) {
	id := uuid.New().String()
	now := time.Now()

	query := `
		INSERT INTO user_organizations (id, user_id, organization_id, role, is_default, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, user_id, organization_id, role, is_default, created_at, updated_at
	`

	var uo UserOrganization
	err := r.db.QueryRow(ctx, query, id, userID, organizationID, role, isDefault, now, now).Scan(
		&uo.ID, &uo.UserID, &uo.OrganizationID, &uo.Role, &uo.IsDefault, &uo.CreatedAt, &uo.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &uo, nil
}

// GetByID retrieves a user organization by ID
func (r *UserOrganizationRepository) GetByID(ctx context.Context, id string) (*UserOrganization, error) {
	query := `
		SELECT id, user_id, organization_id, role, is_default, created_at, updated_at
		FROM user_organizations
		WHERE id = $1
	`

	var uo UserOrganization
	err := r.db.QueryRow(ctx, query, id).Scan(
		&uo.ID, &uo.UserID, &uo.OrganizationID, &uo.Role, &uo.IsDefault, &uo.CreatedAt, &uo.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserOrganizationNotFound
		}
		return nil, err
	}

	return &uo, nil
}

// GetByUserAndOrganization retrieves a user organization by user_id and organization_id
func (r *UserOrganizationRepository) GetByUserAndOrganization(ctx context.Context, userID, organizationID string) (*UserOrganization, error) {
	query := `
		SELECT id, user_id, organization_id, role, is_default, created_at, updated_at
		FROM user_organizations
		WHERE user_id = $1 AND organization_id = $2
	`

	var uo UserOrganization
	err := r.db.QueryRow(ctx, query, userID, organizationID).Scan(
		&uo.ID, &uo.UserID, &uo.OrganizationID, &uo.Role, &uo.IsDefault, &uo.CreatedAt, &uo.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserOrganizationNotFound
		}
		return nil, err
	}

	return &uo, nil
}

// ListByUserID retrieves all organizations for a user
func (r *UserOrganizationRepository) ListByUserID(ctx context.Context, userID string) ([]*UserOrganization, error) {
	query := `
		SELECT id, user_id, organization_id, role, is_default, created_at, updated_at
		FROM user_organizations
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var uos []*UserOrganization
	for rows.Next() {
		var uo UserOrganization
		err := rows.Scan(&uo.ID, &uo.UserID, &uo.OrganizationID, &uo.Role, &uo.IsDefault, &uo.CreatedAt, &uo.UpdatedAt)
		if err != nil {
			return nil, err
		}
		uos = append(uos, &uo)
	}

	return uos, rows.Err()
}

// ListByOrganizationID retrieves all users for an organization
func (r *UserOrganizationRepository) ListByOrganizationID(ctx context.Context, organizationID string) ([]*UserOrganization, error) {
	query := `
		SELECT id, user_id, organization_id, role, is_default, created_at, updated_at
		FROM user_organizations
		WHERE organization_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var uos []*UserOrganization
	for rows.Next() {
		var uo UserOrganization
		err := rows.Scan(&uo.ID, &uo.UserID, &uo.OrganizationID, &uo.Role, &uo.IsDefault, &uo.CreatedAt, &uo.UpdatedAt)
		if err != nil {
			return nil, err
		}
		uos = append(uos, &uo)
	}

	return uos, rows.Err()
}

// Update modifies an existing user organization
func (r *UserOrganizationRepository) Update(ctx context.Context, id, role string, isDefault bool) (*UserOrganization, error) {
	query := `
		UPDATE user_organizations
		SET role = $2, is_default = $3, updated_at = $4
		WHERE id = $1
		RETURNING id, user_id, organization_id, role, is_default, created_at, updated_at
	`

	var uo UserOrganization
	err := r.db.QueryRow(ctx, query, id, role, isDefault, time.Now()).Scan(
		&uo.ID, &uo.UserID, &uo.OrganizationID, &uo.Role, &uo.IsDefault, &uo.CreatedAt, &uo.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserOrganizationNotFound
		}
		return nil, err
	}

	return &uo, nil
}

// Delete removes a user organization (hard delete)
func (r *UserOrganizationRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM user_organizations WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrUserOrganizationNotFound
	}

	return nil
}

// List retrieves user organizations with pagination
func (r *UserOrganizationRepository) List(ctx context.Context, limit int, offset int) ([]*UserOrganization, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	query := `
		SELECT id, user_id, organization_id, role, is_default, created_at, updated_at
		FROM user_organizations
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var uos []*UserOrganization
	for rows.Next() {
		var uo UserOrganization
		err := rows.Scan(&uo.ID, &uo.UserID, &uo.OrganizationID, &uo.Role, &uo.IsDefault, &uo.CreatedAt, &uo.UpdatedAt)
		if err != nil {
			return nil, err
		}
		uos = append(uos, &uo)
	}

	return uos, rows.Err()
}
