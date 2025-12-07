package repository

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrInvitationNotFound = errors.New("invitation not found")
	ErrInvitationExpired  = errors.New("invitation expired")
	ErrInvitationUsed     = errors.New("invitation already used")
)

// Invitation represents the database model
type Invitation struct {
	ID             string
	OrganizationID string
	Email          string
	Role           string
	Token          string
	InvitedBy      string
	Status         string // 'pending', 'accepted', 'expired', 'cancelled'
	ExpiresAt      time.Time
	AcceptedAt     *time.Time
	AcceptedBy     *string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// InvitationRepository handles database operations for invitations
type InvitationRepository struct {
	db DB
}

// NewInvitationRepository creates a new repository
func NewInvitationRepository(pool *pgxpool.Pool) *InvitationRepository {
	return &InvitationRepository{db: pool}
}

// NewInvitationRepositoryWithDB creates a repository with custom DB interface (for testing)
func NewInvitationRepositoryWithDB(db DB) *InvitationRepository {
	return &InvitationRepository{db: db}
}

// generateToken creates a secure random token
func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// Create inserts a new invitation
func (r *InvitationRepository) Create(ctx context.Context, organizationID, email, role, invitedBy string) (*Invitation, error) {
	id := uuid.New().String()
	now := time.Now()
	expiresAt := now.Add(7 * 24 * time.Hour) // 7 days expiry

	token, err := generateToken()
	if err != nil {
		return nil, err
	}

	if role == "" {
		role = "member"
	}

	query := `
		INSERT INTO invitations (id, organization_id, email, role, token, invited_by, status, expires_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, 'pending', $7, $8, $9)
		RETURNING id, organization_id, email, role, token, invited_by, status, expires_at, accepted_at, accepted_by, created_at, updated_at
	`

	var inv Invitation
	err = r.db.QueryRow(ctx, query, id, organizationID, email, role, token, invitedBy, expiresAt, now, now).Scan(
		&inv.ID, &inv.OrganizationID, &inv.Email, &inv.Role, &inv.Token,
		&inv.InvitedBy, &inv.Status, &inv.ExpiresAt, &inv.AcceptedAt, &inv.AcceptedBy,
		&inv.CreatedAt, &inv.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &inv, nil
}

// GetByID retrieves an invitation by ID
func (r *InvitationRepository) GetByID(ctx context.Context, id string) (*Invitation, error) {
	query := `
		SELECT id, organization_id, email, role, token, invited_by, status, expires_at, accepted_at, accepted_by, created_at, updated_at
		FROM invitations
		WHERE id = $1
	`

	var inv Invitation
	err := r.db.QueryRow(ctx, query, id).Scan(
		&inv.ID, &inv.OrganizationID, &inv.Email, &inv.Role, &inv.Token,
		&inv.InvitedBy, &inv.Status, &inv.ExpiresAt, &inv.AcceptedAt, &inv.AcceptedBy,
		&inv.CreatedAt, &inv.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvitationNotFound
		}
		return nil, err
	}

	return &inv, nil
}

// GetByToken retrieves an invitation by token
func (r *InvitationRepository) GetByToken(ctx context.Context, token string) (*Invitation, error) {
	query := `
		SELECT id, organization_id, email, role, token, invited_by, status, expires_at, accepted_at, accepted_by, created_at, updated_at
		FROM invitations
		WHERE token = $1
	`

	var inv Invitation
	err := r.db.QueryRow(ctx, query, token).Scan(
		&inv.ID, &inv.OrganizationID, &inv.Email, &inv.Role, &inv.Token,
		&inv.InvitedBy, &inv.Status, &inv.ExpiresAt, &inv.AcceptedAt, &inv.AcceptedBy,
		&inv.CreatedAt, &inv.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvitationNotFound
		}
		return nil, err
	}

	return &inv, nil
}

// Accept marks an invitation as accepted
func (r *InvitationRepository) Accept(ctx context.Context, id, acceptedBy string) (*Invitation, error) {
	now := time.Now()

	// First check status and expiry
	inv, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if inv.Status != "pending" {
		return nil, ErrInvitationUsed
	}

	if inv.ExpiresAt.Before(now) {
		// Update status to expired
		_, _ = r.db.Exec(ctx, `UPDATE invitations SET status = 'expired', updated_at = $1 WHERE id = $2`, now, id)
		return nil, ErrInvitationExpired
	}

	query := `
		UPDATE invitations
		SET status = 'accepted', accepted_at = $2, accepted_by = $3, updated_at = $4
		WHERE id = $1
		RETURNING id, organization_id, email, role, token, invited_by, status, expires_at, accepted_at, accepted_by, created_at, updated_at
	`

	err = r.db.QueryRow(ctx, query, id, now, acceptedBy, now).Scan(
		&inv.ID, &inv.OrganizationID, &inv.Email, &inv.Role, &inv.Token,
		&inv.InvitedBy, &inv.Status, &inv.ExpiresAt, &inv.AcceptedAt, &inv.AcceptedBy,
		&inv.CreatedAt, &inv.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return inv, nil
}

// Cancel marks an invitation as cancelled
func (r *InvitationRepository) Cancel(ctx context.Context, id string) error {
	query := `UPDATE invitations SET status = 'cancelled', updated_at = $2 WHERE id = $1 AND status = 'pending'`

	result, err := r.db.Exec(ctx, query, id, time.Now())
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrInvitationNotFound
	}

	return nil
}

// List retrieves invitations for an organization with optional status filter
func (r *InvitationRepository) List(ctx context.Context, organizationID, status string, limit, offset int) ([]*Invitation, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	var query string
	var args []interface{}

	if status != "" {
		query = `
			SELECT id, organization_id, email, role, token, invited_by, status, expires_at, accepted_at, accepted_by, created_at, updated_at
			FROM invitations
			WHERE organization_id = $1 AND status = $2
			ORDER BY created_at DESC
			LIMIT $3 OFFSET $4
		`
		args = []interface{}{organizationID, status, limit, offset}
	} else {
		query = `
			SELECT id, organization_id, email, role, token, invited_by, status, expires_at, accepted_at, accepted_by, created_at, updated_at
			FROM invitations
			WHERE organization_id = $1
			ORDER BY created_at DESC
			LIMIT $2 OFFSET $3
		`
		args = []interface{}{organizationID, limit, offset}
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invitations []*Invitation
	for rows.Next() {
		var inv Invitation
		err := rows.Scan(
			&inv.ID, &inv.OrganizationID, &inv.Email, &inv.Role, &inv.Token,
			&inv.InvitedBy, &inv.Status, &inv.ExpiresAt, &inv.AcceptedAt, &inv.AcceptedBy,
			&inv.CreatedAt, &inv.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		invitations = append(invitations, &inv)
	}

	return invitations, rows.Err()
}

// Resend regenerates the token and extends the expiry
func (r *InvitationRepository) Resend(ctx context.Context, id string) (*Invitation, error) {
	now := time.Now()
	expiresAt := now.Add(7 * 24 * time.Hour)

	token, err := generateToken()
	if err != nil {
		return nil, err
	}

	query := `
		UPDATE invitations
		SET token = $2, expires_at = $3, status = 'pending', updated_at = $4
		WHERE id = $1 AND status IN ('pending', 'expired')
		RETURNING id, organization_id, email, role, token, invited_by, status, expires_at, accepted_at, accepted_by, created_at, updated_at
	`

	var inv Invitation
	err = r.db.QueryRow(ctx, query, id, token, expiresAt, now).Scan(
		&inv.ID, &inv.OrganizationID, &inv.Email, &inv.Role, &inv.Token,
		&inv.InvitedBy, &inv.Status, &inv.ExpiresAt, &inv.AcceptedAt, &inv.AcceptedBy,
		&inv.CreatedAt, &inv.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvitationNotFound
		}
		return nil, err
	}

	return &inv, nil
}

// GetPendingByEmailAndOrg checks if there's already a pending invitation for email+org
func (r *InvitationRepository) GetPendingByEmailAndOrg(ctx context.Context, email, organizationID string) (*Invitation, error) {
	query := `
		SELECT id, organization_id, email, role, token, invited_by, status, expires_at, accepted_at, accepted_by, created_at, updated_at
		FROM invitations
		WHERE email = $1 AND organization_id = $2 AND status = 'pending'
	`

	var inv Invitation
	err := r.db.QueryRow(ctx, query, email, organizationID).Scan(
		&inv.ID, &inv.OrganizationID, &inv.Email, &inv.Role, &inv.Token,
		&inv.InvitedBy, &inv.Status, &inv.ExpiresAt, &inv.AcceptedAt, &inv.AcceptedBy,
		&inv.CreatedAt, &inv.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvitationNotFound
		}
		return nil, err
	}

	return &inv, nil
}
