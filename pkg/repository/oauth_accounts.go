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
	ErrOAuthAccountNotFound = errors.New("oauth account not found")
)

// OAuthAccount represents the database model
type OAuthAccount struct {
	ID             string
	AppUserID      string
	Provider       string // 'google', 'line'
	ProviderUserID string
	Email          *string
	AccessToken    *string
	RefreshToken   *string
	TokenExpiresAt *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// OAuthAccountRepository handles database operations for oauth_accounts
type OAuthAccountRepository struct {
	db DB
}

// NewOAuthAccountRepository creates a new repository
func NewOAuthAccountRepository(pool *pgxpool.Pool) *OAuthAccountRepository {
	return &OAuthAccountRepository{db: pool}
}

// NewOAuthAccountRepositoryWithDB creates a repository with custom DB interface (for testing)
func NewOAuthAccountRepositoryWithDB(db DB) *OAuthAccountRepository {
	return &OAuthAccountRepository{db: db}
}

// Create inserts a new oauth account
func (r *OAuthAccountRepository) Create(ctx context.Context, appUserID, provider, providerUserID string, email, accessToken, refreshToken *string, tokenExpiresAt *time.Time) (*OAuthAccount, error) {
	id := uuid.New().String()
	now := time.Now()

	query := `
		INSERT INTO oauth_accounts (id, app_user_id, provider, provider_user_id, email, access_token, refresh_token, token_expires_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, app_user_id, provider, provider_user_id, email, access_token, refresh_token, token_expires_at, created_at, updated_at
	`

	var account OAuthAccount
	err := r.db.QueryRow(ctx, query, id, appUserID, provider, providerUserID, email, accessToken, refreshToken, tokenExpiresAt, now, now).Scan(
		&account.ID, &account.AppUserID, &account.Provider, &account.ProviderUserID, &account.Email, &account.AccessToken, &account.RefreshToken, &account.TokenExpiresAt, &account.CreatedAt, &account.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

// GetByID retrieves an oauth account by ID
func (r *OAuthAccountRepository) GetByID(ctx context.Context, id string) (*OAuthAccount, error) {
	query := `
		SELECT id, app_user_id, provider, provider_user_id, email, access_token, refresh_token, token_expires_at, created_at, updated_at
		FROM oauth_accounts
		WHERE id = $1
	`

	var account OAuthAccount
	err := r.db.QueryRow(ctx, query, id).Scan(
		&account.ID, &account.AppUserID, &account.Provider, &account.ProviderUserID, &account.Email, &account.AccessToken, &account.RefreshToken, &account.TokenExpiresAt, &account.CreatedAt, &account.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrOAuthAccountNotFound
		}
		return nil, err
	}

	return &account, nil
}

// GetByProviderAndProviderUserID retrieves an oauth account by provider and provider_user_id
func (r *OAuthAccountRepository) GetByProviderAndProviderUserID(ctx context.Context, provider, providerUserID string) (*OAuthAccount, error) {
	query := `
		SELECT id, app_user_id, provider, provider_user_id, email, access_token, refresh_token, token_expires_at, created_at, updated_at
		FROM oauth_accounts
		WHERE provider = $1 AND provider_user_id = $2
	`

	var account OAuthAccount
	err := r.db.QueryRow(ctx, query, provider, providerUserID).Scan(
		&account.ID, &account.AppUserID, &account.Provider, &account.ProviderUserID, &account.Email, &account.AccessToken, &account.RefreshToken, &account.TokenExpiresAt, &account.CreatedAt, &account.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrOAuthAccountNotFound
		}
		return nil, err
	}

	return &account, nil
}

// ListByAppUserID retrieves all oauth accounts for an app user
func (r *OAuthAccountRepository) ListByAppUserID(ctx context.Context, appUserID string) ([]*OAuthAccount, error) {
	query := `
		SELECT id, app_user_id, provider, provider_user_id, email, access_token, refresh_token, token_expires_at, created_at, updated_at
		FROM oauth_accounts
		WHERE app_user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, appUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*OAuthAccount
	for rows.Next() {
		var account OAuthAccount
		err := rows.Scan(&account.ID, &account.AppUserID, &account.Provider, &account.ProviderUserID, &account.Email, &account.AccessToken, &account.RefreshToken, &account.TokenExpiresAt, &account.CreatedAt, &account.UpdatedAt)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, &account)
	}

	return accounts, rows.Err()
}

// UpdateTokens updates the access and refresh tokens
func (r *OAuthAccountRepository) UpdateTokens(ctx context.Context, id string, accessToken, refreshToken *string, tokenExpiresAt *time.Time) (*OAuthAccount, error) {
	query := `
		UPDATE oauth_accounts
		SET access_token = $2, refresh_token = $3, token_expires_at = $4, updated_at = $5
		WHERE id = $1
		RETURNING id, app_user_id, provider, provider_user_id, email, access_token, refresh_token, token_expires_at, created_at, updated_at
	`

	var account OAuthAccount
	err := r.db.QueryRow(ctx, query, id, accessToken, refreshToken, tokenExpiresAt, time.Now()).Scan(
		&account.ID, &account.AppUserID, &account.Provider, &account.ProviderUserID, &account.Email, &account.AccessToken, &account.RefreshToken, &account.TokenExpiresAt, &account.CreatedAt, &account.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrOAuthAccountNotFound
		}
		return nil, err
	}

	return &account, nil
}

// Delete removes an oauth account
func (r *OAuthAccountRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM oauth_accounts WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrOAuthAccountNotFound
	}

	return nil
}

// DeleteByAppUserID removes all oauth accounts for an app user
func (r *OAuthAccountRepository) DeleteByAppUserID(ctx context.Context, appUserID string) error {
	query := `DELETE FROM oauth_accounts WHERE app_user_id = $1`

	_, err := r.db.Exec(ctx, query, appUserID)
	return err
}
