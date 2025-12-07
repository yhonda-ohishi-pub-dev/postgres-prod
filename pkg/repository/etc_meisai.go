package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrETCMeisaiNotFound = errors.New("etc_meisai not found")
)

// ETCMeisai represents the database model for ETC明細
type ETCMeisai struct {
	ID             int64
	OrganizationID string
	DateFr         *time.Time // nullable
	DateTo         time.Time
	DateToDate     string // YYYY-MM-DD format
	IcFr           string
	IcTo           string
	PriceBf        *int32 // nullable
	Discount       *int32 // nullable
	Price          int32
	Shashu         int32
	CarIdNum       *int32 // nullable
	EtcNum         string
	Detail         *string // nullable
	DtakoRowId     *string // nullable
	Hash           string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// ETCMeisaiRepository handles database operations for etc_meisai
type ETCMeisaiRepository struct {
	db DB
}

// NewETCMeisaiRepository creates a new repository
func NewETCMeisaiRepository(pool *pgxpool.Pool) *ETCMeisaiRepository {
	return &ETCMeisaiRepository{db: pool}
}

// NewETCMeisaiRepositoryWithDB creates a repository with custom DB interface (for testing/RLS)
func NewETCMeisaiRepositoryWithDB(db DB) *ETCMeisaiRepository {
	return &ETCMeisaiRepository{db: db}
}

// Create inserts a new ETC meisai record
func (r *ETCMeisaiRepository) Create(ctx context.Context, organizationID string, meisai *ETCMeisai) (*ETCMeisai, error) {
	now := time.Now()

	query := `
		INSERT INTO etc_meisai (
			organization_id, date_fr, date_to, date_to_date, ic_fr, ic_to,
			price_bf, discount, price, shashu, car_id_num, etc_num,
			detail, dtako_row_id, hash, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
		RETURNING id, organization_id, date_fr, date_to, date_to_date, ic_fr, ic_to,
			price_bf, discount, price, shashu, car_id_num, etc_num,
			detail, dtako_row_id, hash, created_at, updated_at
	`

	var result ETCMeisai
	err := r.db.QueryRow(ctx, query,
		organizationID, meisai.DateFr, meisai.DateTo, meisai.DateToDate,
		meisai.IcFr, meisai.IcTo, meisai.PriceBf, meisai.Discount,
		meisai.Price, meisai.Shashu, meisai.CarIdNum, meisai.EtcNum,
		meisai.Detail, meisai.DtakoRowId, meisai.Hash, now, now,
	).Scan(
		&result.ID, &result.OrganizationID, &result.DateFr, &result.DateTo, &result.DateToDate,
		&result.IcFr, &result.IcTo, &result.PriceBf, &result.Discount,
		&result.Price, &result.Shashu, &result.CarIdNum, &result.EtcNum,
		&result.Detail, &result.DtakoRowId, &result.Hash, &result.CreatedAt, &result.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetByID retrieves an ETC meisai by ID
func (r *ETCMeisaiRepository) GetByID(ctx context.Context, id int64) (*ETCMeisai, error) {
	query := `
		SELECT id, organization_id, date_fr, date_to, date_to_date, ic_fr, ic_to,
			price_bf, discount, price, shashu, car_id_num, etc_num,
			detail, dtako_row_id, hash, created_at, updated_at
		FROM etc_meisai
		WHERE id = $1
	`

	var result ETCMeisai
	err := r.db.QueryRow(ctx, query, id).Scan(
		&result.ID, &result.OrganizationID, &result.DateFr, &result.DateTo, &result.DateToDate,
		&result.IcFr, &result.IcTo, &result.PriceBf, &result.Discount,
		&result.Price, &result.Shashu, &result.CarIdNum, &result.EtcNum,
		&result.Detail, &result.DtakoRowId, &result.Hash, &result.CreatedAt, &result.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrETCMeisaiNotFound
		}
		return nil, err
	}

	return &result, nil
}

// GetByHash retrieves an ETC meisai by hash (for duplicate check)
func (r *ETCMeisaiRepository) GetByHash(ctx context.Context, hash string) (*ETCMeisai, error) {
	query := `
		SELECT id, organization_id, date_fr, date_to, date_to_date, ic_fr, ic_to,
			price_bf, discount, price, shashu, car_id_num, etc_num,
			detail, dtako_row_id, hash, created_at, updated_at
		FROM etc_meisai
		WHERE hash = $1
	`

	var result ETCMeisai
	err := r.db.QueryRow(ctx, query, hash).Scan(
		&result.ID, &result.OrganizationID, &result.DateFr, &result.DateTo, &result.DateToDate,
		&result.IcFr, &result.IcTo, &result.PriceBf, &result.Discount,
		&result.Price, &result.Shashu, &result.CarIdNum, &result.EtcNum,
		&result.Detail, &result.DtakoRowId, &result.Hash, &result.CreatedAt, &result.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrETCMeisaiNotFound
		}
		return nil, err
	}

	return &result, nil
}

// ExistsByHash checks if a record with the given hash exists
func (r *ETCMeisaiRepository) ExistsByHash(ctx context.Context, hash string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM etc_meisai WHERE hash = $1)`

	var exists bool
	err := r.db.QueryRow(ctx, query, hash).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// Update updates an ETC meisai record
func (r *ETCMeisaiRepository) Update(ctx context.Context, meisai *ETCMeisai) (*ETCMeisai, error) {
	now := time.Now()

	query := `
		UPDATE etc_meisai SET
			date_fr = $2, date_to = $3, date_to_date = $4, ic_fr = $5, ic_to = $6,
			price_bf = $7, discount = $8, price = $9, shashu = $10, car_id_num = $11,
			etc_num = $12, detail = $13, dtako_row_id = $14, hash = $15, updated_at = $16
		WHERE id = $1
		RETURNING id, organization_id, date_fr, date_to, date_to_date, ic_fr, ic_to,
			price_bf, discount, price, shashu, car_id_num, etc_num,
			detail, dtako_row_id, hash, created_at, updated_at
	`

	var result ETCMeisai
	err := r.db.QueryRow(ctx, query,
		meisai.ID, meisai.DateFr, meisai.DateTo, meisai.DateToDate,
		meisai.IcFr, meisai.IcTo, meisai.PriceBf, meisai.Discount,
		meisai.Price, meisai.Shashu, meisai.CarIdNum, meisai.EtcNum,
		meisai.Detail, meisai.DtakoRowId, meisai.Hash, now,
	).Scan(
		&result.ID, &result.OrganizationID, &result.DateFr, &result.DateTo, &result.DateToDate,
		&result.IcFr, &result.IcTo, &result.PriceBf, &result.Discount,
		&result.Price, &result.Shashu, &result.CarIdNum, &result.EtcNum,
		&result.Detail, &result.DtakoRowId, &result.Hash, &result.CreatedAt, &result.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrETCMeisaiNotFound
		}
		return nil, err
	}

	return &result, nil
}

// Delete deletes an ETC meisai record
func (r *ETCMeisaiRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM etc_meisai WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrETCMeisaiNotFound
	}

	return nil
}

// ListParams contains parameters for listing ETC meisai
type ETCMeisaiListParams struct {
	PageSize  int
	PageToken string // base64 encoded offset
	DateFrom  *string
	DateTo    *string
	EtcNum    *string
}

// List retrieves ETC meisai with optional filters
func (r *ETCMeisaiRepository) List(ctx context.Context, params ETCMeisaiListParams) ([]*ETCMeisai, int, string, error) {
	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	if params.PageSize > 100 {
		params.PageSize = 100
	}

	offset := 0
	if params.PageToken != "" {
		// Decode page token as offset
		// For simplicity, use numeric offset directly
		// In production, consider using cursor-based pagination
	}

	// Build dynamic query
	query := `
		SELECT id, organization_id, date_fr, date_to, date_to_date, ic_fr, ic_to,
			price_bf, discount, price, shashu, car_id_num, etc_num,
			detail, dtako_row_id, hash, created_at, updated_at
		FROM etc_meisai
		WHERE 1=1
	`
	countQuery := `SELECT COUNT(*) FROM etc_meisai WHERE 1=1`

	args := []interface{}{}
	argIndex := 1

	if params.DateFrom != nil && *params.DateFrom != "" {
		query += ` AND date_to_date >= $` + string(rune('0'+argIndex))
		countQuery += ` AND date_to_date >= $` + string(rune('0'+argIndex))
		args = append(args, *params.DateFrom)
		argIndex++
	}

	if params.DateTo != nil && *params.DateTo != "" {
		query += ` AND date_to_date <= $` + string(rune('0'+argIndex))
		countQuery += ` AND date_to_date <= $` + string(rune('0'+argIndex))
		args = append(args, *params.DateTo)
		argIndex++
	}

	if params.EtcNum != nil && *params.EtcNum != "" {
		query += ` AND etc_num = $` + string(rune('0'+argIndex))
		countQuery += ` AND etc_num = $` + string(rune('0'+argIndex))
		args = append(args, *params.EtcNum)
		argIndex++
	}

	// Get total count
	var totalCount int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, "", err
	}

	// Add ordering and pagination
	query += ` ORDER BY date_to DESC, id DESC`
	query += ` LIMIT $` + string(rune('0'+argIndex)) + ` OFFSET $` + string(rune('0'+argIndex+1))
	args = append(args, params.PageSize+1, offset) // +1 to check if more exists

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, "", err
	}
	defer rows.Close()

	var results []*ETCMeisai
	for rows.Next() {
		var m ETCMeisai
		err := rows.Scan(
			&m.ID, &m.OrganizationID, &m.DateFr, &m.DateTo, &m.DateToDate,
			&m.IcFr, &m.IcTo, &m.PriceBf, &m.Discount,
			&m.Price, &m.Shashu, &m.CarIdNum, &m.EtcNum,
			&m.Detail, &m.DtakoRowId, &m.Hash, &m.CreatedAt, &m.UpdatedAt,
		)
		if err != nil {
			return nil, 0, "", err
		}
		results = append(results, &m)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, "", err
	}

	// Determine next page token
	nextPageToken := ""
	if len(results) > params.PageSize {
		results = results[:params.PageSize]
		// Use simple offset-based pagination
		nextPageToken = "more" // In production, encode actual cursor
	}

	return results, totalCount, nextPageToken, nil
}

// BulkCreate creates multiple ETC meisai records, optionally skipping duplicates
func (r *ETCMeisaiRepository) BulkCreate(ctx context.Context, organizationID string, records []*ETCMeisai, skipDuplicates bool) (int, int, []string, error) {
	createdCount := 0
	skippedCount := 0
	var errors []string

	for i, record := range records {
		if skipDuplicates {
			exists, err := r.ExistsByHash(ctx, record.Hash)
			if err != nil {
				errors = append(errors, "record "+string(rune('0'+i))+": "+err.Error())
				continue
			}
			if exists {
				skippedCount++
				continue
			}
		}

		_, err := r.Create(ctx, organizationID, record)
		if err != nil {
			errors = append(errors, "record "+string(rune('0'+i))+": "+err.Error())
			continue
		}
		createdCount++
	}

	return createdCount, skippedCount, errors, nil
}
