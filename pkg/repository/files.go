package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrFileNotFound = errors.New("file not found")
)

// File represents the database model
type File struct {
	UUID           string
	OrganizationID string
	Filename       string
	Created        string
	Deleted        string
	Type           string
	Blob           *string
}

// FileRepository handles database operations for files
type FileRepository struct {
	db DB
}

// NewFileRepository creates a new repository
func NewFileRepository(pool *pgxpool.Pool) *FileRepository {
	return &FileRepository{db: pool}
}

// NewFileRepositoryWithDB creates a repository with custom DB interface (for testing)
func NewFileRepositoryWithDB(db DB) *FileRepository {
	return &FileRepository{db: db}
}

// Create inserts a new file
func (r *FileRepository) Create(ctx context.Context, organizationID, filename, created, fileType string, blob *string) (*File, error) {
	id := uuid.New().String()

	query := `
		INSERT INTO files (uuid, organization_id, filename, created, deleted, type, blob)
		VALUES ($1, $2, $3, $4, '', $5, $6)
		RETURNING uuid, organization_id, filename, created, deleted, type, blob
	`

	var f File
	err := r.db.QueryRow(ctx, query, id, organizationID, filename, created, fileType, blob).Scan(
		&f.UUID, &f.OrganizationID, &f.Filename, &f.Created, &f.Deleted, &f.Type, &f.Blob,
	)
	if err != nil {
		return nil, err
	}

	return &f, nil
}

// GetByUUID retrieves a file by UUID
func (r *FileRepository) GetByUUID(ctx context.Context, uuid string) (*File, error) {
	query := `
		SELECT uuid, organization_id, filename, created, deleted, type, blob
		FROM files
		WHERE uuid = $1 AND deleted = ''
	`

	var f File
	err := r.db.QueryRow(ctx, query, uuid).Scan(
		&f.UUID, &f.OrganizationID, &f.Filename, &f.Created, &f.Deleted, &f.Type, &f.Blob,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrFileNotFound
		}
		return nil, err
	}

	return &f, nil
}

// Update modifies an existing file
func (r *FileRepository) Update(ctx context.Context, uuid, filename, fileType string, blob *string) (*File, error) {
	query := `
		UPDATE files
		SET filename = $2, type = $3, blob = $4
		WHERE uuid = $1 AND deleted = ''
		RETURNING uuid, organization_id, filename, created, deleted, type, blob
	`

	var f File
	err := r.db.QueryRow(ctx, query, uuid, filename, fileType, blob).Scan(
		&f.UUID, &f.OrganizationID, &f.Filename, &f.Created, &f.Deleted, &f.Type, &f.Blob,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrFileNotFound
		}
		return nil, err
	}

	return &f, nil
}

// Delete soft-deletes a file
func (r *FileRepository) Delete(ctx context.Context, uuid, deletedTimestamp string) error {
	query := `
		UPDATE files
		SET deleted = $2
		WHERE uuid = $1 AND deleted = ''
	`

	result, err := r.db.Exec(ctx, query, uuid, deletedTimestamp)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrFileNotFound
	}

	return nil
}

// ListByOrganization retrieves files by organization with pagination
func (r *FileRepository) ListByOrganization(ctx context.Context, organizationID string, limit int, offset int) ([]*File, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	query := `
		SELECT uuid, organization_id, filename, created, deleted, type, blob
		FROM files
		WHERE organization_id = $1 AND deleted = ''
		ORDER BY created DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, organizationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []*File
	for rows.Next() {
		var f File
		err := rows.Scan(&f.UUID, &f.OrganizationID, &f.Filename, &f.Created, &f.Deleted, &f.Type, &f.Blob)
		if err != nil {
			return nil, err
		}
		files = append(files, &f)
	}

	return files, rows.Err()
}

// List retrieves all files with pagination
func (r *FileRepository) List(ctx context.Context, limit int, offset int) ([]*File, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	query := `
		SELECT uuid, organization_id, filename, created, deleted, type, blob
		FROM files
		WHERE deleted = ''
		ORDER BY created DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []*File
	for rows.Next() {
		var f File
		err := rows.Scan(&f.UUID, &f.OrganizationID, &f.Filename, &f.Created, &f.Deleted, &f.Type, &f.Blob)
		if err != nil {
			return nil, err
		}
		files = append(files, &f)
	}

	return files, rows.Err()
}
