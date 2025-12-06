package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrFlickrPhotoNotFound = errors.New("flickr photo not found")
)

// FlickrPhoto represents the database model
type FlickrPhoto struct {
	ID             string
	OrganizationID string
	Secret         string
	Server         string
}

// FlickrPhotoRepository handles database operations for flickr_photo
type FlickrPhotoRepository struct {
	db DB
}

// NewFlickrPhotoRepository creates a new repository
func NewFlickrPhotoRepository(pool *pgxpool.Pool) *FlickrPhotoRepository {
	return &FlickrPhotoRepository{db: pool}
}

// NewFlickrPhotoRepositoryWithDB creates a repository with custom DB interface (for testing)
func NewFlickrPhotoRepositoryWithDB(db DB) *FlickrPhotoRepository {
	return &FlickrPhotoRepository{db: db}
}

// Create inserts a new flickr photo
func (r *FlickrPhotoRepository) Create(ctx context.Context, id, organizationID, secret, server string) (*FlickrPhoto, error) {
	query := `
		INSERT INTO flickr_photo (id, organization_id, secret, server)
		VALUES ($1, $2, $3, $4)
		RETURNING id, organization_id, secret, server
	`

	var photo FlickrPhoto
	err := r.db.QueryRow(ctx, query, id, organizationID, secret, server).Scan(
		&photo.ID, &photo.OrganizationID, &photo.Secret, &photo.Server,
	)
	if err != nil {
		return nil, err
	}

	return &photo, nil
}

// GetByID retrieves a flickr photo by ID
func (r *FlickrPhotoRepository) GetByID(ctx context.Context, id string) (*FlickrPhoto, error) {
	query := `
		SELECT id, organization_id, secret, server
		FROM flickr_photo
		WHERE id = $1
	`

	var photo FlickrPhoto
	err := r.db.QueryRow(ctx, query, id).Scan(
		&photo.ID, &photo.OrganizationID, &photo.Secret, &photo.Server,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrFlickrPhotoNotFound
		}
		return nil, err
	}

	return &photo, nil
}

// Update modifies an existing flickr photo
func (r *FlickrPhotoRepository) Update(ctx context.Context, id, secret, server string) (*FlickrPhoto, error) {
	query := `
		UPDATE flickr_photo
		SET secret = $2, server = $3
		WHERE id = $1
		RETURNING id, organization_id, secret, server
	`

	var photo FlickrPhoto
	err := r.db.QueryRow(ctx, query, id, secret, server).Scan(
		&photo.ID, &photo.OrganizationID, &photo.Secret, &photo.Server,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrFlickrPhotoNotFound
		}
		return nil, err
	}

	return &photo, nil
}

// Delete hard-deletes a flickr photo
func (r *FlickrPhotoRepository) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM flickr_photo
		WHERE id = $1
	`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrFlickrPhotoNotFound
	}

	return nil
}

// List retrieves flickr photos with pagination
func (r *FlickrPhotoRepository) List(ctx context.Context, limit int, offset int) ([]*FlickrPhoto, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	query := `
		SELECT id, organization_id, secret, server
		FROM flickr_photo
		ORDER BY id ASC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var photos []*FlickrPhoto
	for rows.Next() {
		var photo FlickrPhoto
		err := rows.Scan(&photo.ID, &photo.OrganizationID, &photo.Secret, &photo.Server)
		if err != nil {
			return nil, err
		}
		photos = append(photos, &photo)
	}

	return photos, rows.Err()
}

// ListByOrganization retrieves flickr photos for a specific organization with pagination
func (r *FlickrPhotoRepository) ListByOrganization(ctx context.Context, organizationID string, limit int, offset int) ([]*FlickrPhoto, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	query := `
		SELECT id, organization_id, secret, server
		FROM flickr_photo
		WHERE organization_id = $1
		ORDER BY id ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, organizationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var photos []*FlickrPhoto
	for rows.Next() {
		var photo FlickrPhoto
		err := rows.Scan(&photo.ID, &photo.OrganizationID, &photo.Secret, &photo.Server)
		if err != nil {
			return nil, err
		}
		photos = append(photos, &photo)
	}

	return photos, rows.Err()
}
