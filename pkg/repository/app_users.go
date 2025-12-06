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
	ErrAppUserNotFound = errors.New("app user not found")
)

// AppUser represents the database model
type AppUser struct {
	ID           string
	IamEmail     string
	DisplayName  string
	IsSuperadmin bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

// AppUserRepository handles database operations for app_users
type AppUserRepository struct {
	db DB
}

// NewAppUserRepository creates a new repository
func NewAppUserRepository(pool *pgxpool.Pool) *AppUserRepository {
	return &AppUserRepository{db: pool}
}

// NewAppUserRepositoryWithDB creates a repository with custom DB interface (for testing)
func NewAppUserRepositoryWithDB(db DB) *AppUserRepository {
	return &AppUserRepository{db: db}
}

// Create inserts a new app user
func (r *AppUserRepository) Create(ctx context.Context, iamEmail, displayName string, isSuperadmin bool) (*AppUser, error) {
	id := uuid.New().String()
	now := time.Now()

	query := `
		INSERT INTO app_users (id, iam_email, display_name, is_superadmin, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, iam_email, display_name, is_superadmin, created_at, updated_at, deleted_at
	`

	var user AppUser
	err := r.db.QueryRow(ctx, query, id, iamEmail, displayName, isSuperadmin, now, now).Scan(
		&user.ID, &user.IamEmail, &user.DisplayName, &user.IsSuperadmin, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetByID retrieves an app user by ID
func (r *AppUserRepository) GetByID(ctx context.Context, id string) (*AppUser, error) {
	query := `
		SELECT id, iam_email, display_name, is_superadmin, created_at, updated_at, deleted_at
		FROM app_users
		WHERE id = $1 AND deleted_at IS NULL
	`

	var user AppUser
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.IamEmail, &user.DisplayName, &user.IsSuperadmin, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrAppUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// GetByIamEmail retrieves an app user by IAM email
func (r *AppUserRepository) GetByIamEmail(ctx context.Context, iamEmail string) (*AppUser, error) {
	query := `
		SELECT id, iam_email, display_name, is_superadmin, created_at, updated_at, deleted_at
		FROM app_users
		WHERE iam_email = $1 AND deleted_at IS NULL
	`

	var user AppUser
	err := r.db.QueryRow(ctx, query, iamEmail).Scan(
		&user.ID, &user.IamEmail, &user.DisplayName, &user.IsSuperadmin, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrAppUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// Update modifies an existing app user
func (r *AppUserRepository) Update(ctx context.Context, id, displayName string, isSuperadmin bool) (*AppUser, error) {
	query := `
		UPDATE app_users
		SET display_name = $2, is_superadmin = $3, updated_at = $4
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING id, iam_email, display_name, is_superadmin, created_at, updated_at, deleted_at
	`

	var user AppUser
	err := r.db.QueryRow(ctx, query, id, displayName, isSuperadmin, time.Now()).Scan(
		&user.ID, &user.IamEmail, &user.DisplayName, &user.IsSuperadmin, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrAppUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// Delete soft-deletes an app user
func (r *AppUserRepository) Delete(ctx context.Context, id string) error {
	query := `
		UPDATE app_users
		SET deleted_at = $2, updated_at = $2
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.Exec(ctx, query, id, time.Now())
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrAppUserNotFound
	}

	return nil
}

// List retrieves app users with pagination
func (r *AppUserRepository) List(ctx context.Context, limit int, offset int) ([]*AppUser, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	query := `
		SELECT id, iam_email, display_name, is_superadmin, created_at, updated_at, deleted_at
		FROM app_users
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*AppUser
	for rows.Next() {
		var user AppUser
		err := rows.Scan(&user.ID, &user.IamEmail, &user.DisplayName, &user.IsSuperadmin, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, rows.Err()
}
