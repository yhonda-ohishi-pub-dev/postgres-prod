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
	ErrKudgfryNotFound = errors.New("kudgfry not found")
)

// Kudgfry represents the database model
type Kudgfry struct {
	UUID                        string
	OrganizationID              string
	Hash                        string
	Created                     string
	Deleted                     *string
	KudguriUuid                 *string
	TargetDriverType            string
	UnkouNo                     *string
	UnkouDate                   *string
	ReadDate                    *string
	OfficeCd                    *string
	OfficeName                  *string
	VehicleCd                   *string
	VehicleName                 *string
	DriverCd1                   *string
	DriverName1                 *string
	DriverCd2                   *string
	DriverName2                 *string
	RelevantDatetime            *string
	RefuelInspectCategory       *string
	RefuelInspectCategoryName   *string
	RefuelInspectType           *string
	RefuelInspectTypeName       *string
	RefuelInspectKind           *string
	RefuelInspectKindName       *string
	RefillAmount                *string
	OwnOtherType                *string
	Mileage                     *string
	MeterValue                  *string
}

// KudgfryRepository handles database operations for kudgfry
type KudgfryRepository struct {
	db DB
}

// NewKudgfryRepository creates a new repository
func NewKudgfryRepository(pool *pgxpool.Pool) *KudgfryRepository {
	return &KudgfryRepository{db: pool}
}

// NewKudgfryRepositoryWithDB creates a repository with custom DB interface (for testing)
func NewKudgfryRepositoryWithDB(db DB) *KudgfryRepository {
	return &KudgfryRepository{db: db}
}

// Create inserts a new kudgfry record
func (r *KudgfryRepository) Create(ctx context.Context, kudgfry *Kudgfry) (*Kudgfry, error) {
	if kudgfry.UUID == "" {
		kudgfry.UUID = uuid.New().String()
	}
	if kudgfry.Created == "" {
		kudgfry.Created = time.Now().Format(time.RFC3339)
	}

	query := `
		INSERT INTO kudgfry (
			uuid, "OrganizationID", "Hash", "Created", "Deleted", "KudguriUuid",
			"TargetDriverType", "UnkouNo", "UnkouDate", "ReadDate",
			"OfficeCd", "OfficeName", "VehicleCd", "VehicleName",
			"DriverCd1", "DriverName1", "DriverCd2", "DriverName2",
			"RelevantDatetime", "RefuelInspectCategory", "RefuelInspectCategoryName",
			"RefuelInspectType", "RefuelInspectTypeName", "RefuelInspectKind",
			"RefuelInspectKindName", "RefillAmount", "OwnOtherType",
			mileage, "MeterValue"
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16, $17, $18,
			$19, $20, $21, $22, $23, $24, $25, $26,
			$27, $28, $29
		)
		RETURNING
			uuid, "OrganizationID", "Hash", "Created", "Deleted", "KudguriUuid",
			"TargetDriverType", "UnkouNo", "UnkouDate", "ReadDate",
			"OfficeCd", "OfficeName", "VehicleCd", "VehicleName",
			"DriverCd1", "DriverName1", "DriverCd2", "DriverName2",
			"RelevantDatetime", "RefuelInspectCategory", "RefuelInspectCategoryName",
			"RefuelInspectType", "RefuelInspectTypeName", "RefuelInspectKind",
			"RefuelInspectKindName", "RefillAmount", "OwnOtherType",
			mileage, "MeterValue"
	`

	var result Kudgfry
	err := r.db.QueryRow(ctx, query,
		kudgfry.UUID, kudgfry.OrganizationID, kudgfry.Hash, kudgfry.Created,
		kudgfry.Deleted, kudgfry.KudguriUuid, kudgfry.TargetDriverType,
		kudgfry.UnkouNo, kudgfry.UnkouDate, kudgfry.ReadDate,
		kudgfry.OfficeCd, kudgfry.OfficeName, kudgfry.VehicleCd, kudgfry.VehicleName,
		kudgfry.DriverCd1, kudgfry.DriverName1, kudgfry.DriverCd2, kudgfry.DriverName2,
		kudgfry.RelevantDatetime, kudgfry.RefuelInspectCategory, kudgfry.RefuelInspectCategoryName,
		kudgfry.RefuelInspectType, kudgfry.RefuelInspectTypeName, kudgfry.RefuelInspectKind,
		kudgfry.RefuelInspectKindName, kudgfry.RefillAmount, kudgfry.OwnOtherType,
		kudgfry.Mileage, kudgfry.MeterValue,
	).Scan(
		&result.UUID, &result.OrganizationID, &result.Hash, &result.Created,
		&result.Deleted, &result.KudguriUuid, &result.TargetDriverType,
		&result.UnkouNo, &result.UnkouDate, &result.ReadDate,
		&result.OfficeCd, &result.OfficeName, &result.VehicleCd, &result.VehicleName,
		&result.DriverCd1, &result.DriverName1, &result.DriverCd2, &result.DriverName2,
		&result.RelevantDatetime, &result.RefuelInspectCategory, &result.RefuelInspectCategoryName,
		&result.RefuelInspectType, &result.RefuelInspectTypeName, &result.RefuelInspectKind,
		&result.RefuelInspectKindName, &result.RefillAmount, &result.OwnOtherType,
		&result.Mileage, &result.MeterValue,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetByUUID retrieves a kudgfry record by UUID
func (r *KudgfryRepository) GetByUUID(ctx context.Context, uuid string) (*Kudgfry, error) {
	query := `
		SELECT
			uuid, "OrganizationID", "Hash", "Created", "Deleted", "KudguriUuid",
			"TargetDriverType", "UnkouNo", "UnkouDate", "ReadDate",
			"OfficeCd", "OfficeName", "VehicleCd", "VehicleName",
			"DriverCd1", "DriverName1", "DriverCd2", "DriverName2",
			"RelevantDatetime", "RefuelInspectCategory", "RefuelInspectCategoryName",
			"RefuelInspectType", "RefuelInspectTypeName", "RefuelInspectKind",
			"RefuelInspectKindName", "RefillAmount", "OwnOtherType",
			mileage, "MeterValue"
		FROM kudgfry
		WHERE uuid = $1 AND "Deleted" IS NULL
	`

	var kudgfry Kudgfry
	err := r.db.QueryRow(ctx, query, uuid).Scan(
		&kudgfry.UUID, &kudgfry.OrganizationID, &kudgfry.Hash, &kudgfry.Created,
		&kudgfry.Deleted, &kudgfry.KudguriUuid, &kudgfry.TargetDriverType,
		&kudgfry.UnkouNo, &kudgfry.UnkouDate, &kudgfry.ReadDate,
		&kudgfry.OfficeCd, &kudgfry.OfficeName, &kudgfry.VehicleCd, &kudgfry.VehicleName,
		&kudgfry.DriverCd1, &kudgfry.DriverName1, &kudgfry.DriverCd2, &kudgfry.DriverName2,
		&kudgfry.RelevantDatetime, &kudgfry.RefuelInspectCategory, &kudgfry.RefuelInspectCategoryName,
		&kudgfry.RefuelInspectType, &kudgfry.RefuelInspectTypeName, &kudgfry.RefuelInspectKind,
		&kudgfry.RefuelInspectKindName, &kudgfry.RefillAmount, &kudgfry.OwnOtherType,
		&kudgfry.Mileage, &kudgfry.MeterValue,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrKudgfryNotFound
		}
		return nil, err
	}

	return &kudgfry, nil
}

// Update modifies an existing kudgfry record
func (r *KudgfryRepository) Update(ctx context.Context, kudgfry *Kudgfry) (*Kudgfry, error) {
	query := `
		UPDATE kudgfry
		SET
			"OrganizationID" = $2, "Hash" = $3, "KudguriUuid" = $4,
			"TargetDriverType" = $5, "UnkouNo" = $6, "UnkouDate" = $7, "ReadDate" = $8,
			"OfficeCd" = $9, "OfficeName" = $10, "VehicleCd" = $11, "VehicleName" = $12,
			"DriverCd1" = $13, "DriverName1" = $14, "DriverCd2" = $15, "DriverName2" = $16,
			"RelevantDatetime" = $17, "RefuelInspectCategory" = $18, "RefuelInspectCategoryName" = $19,
			"RefuelInspectType" = $20, "RefuelInspectTypeName" = $21, "RefuelInspectKind" = $22,
			"RefuelInspectKindName" = $23, "RefillAmount" = $24, "OwnOtherType" = $25,
			mileage = $26, "MeterValue" = $27
		WHERE uuid = $1 AND "Deleted" IS NULL
		RETURNING
			uuid, "OrganizationID", "Hash", "Created", "Deleted", "KudguriUuid",
			"TargetDriverType", "UnkouNo", "UnkouDate", "ReadDate",
			"OfficeCd", "OfficeName", "VehicleCd", "VehicleName",
			"DriverCd1", "DriverName1", "DriverCd2", "DriverName2",
			"RelevantDatetime", "RefuelInspectCategory", "RefuelInspectCategoryName",
			"RefuelInspectType", "RefuelInspectTypeName", "RefuelInspectKind",
			"RefuelInspectKindName", "RefillAmount", "OwnOtherType",
			mileage, "MeterValue"
	`

	var result Kudgfry
	err := r.db.QueryRow(ctx, query,
		kudgfry.UUID, kudgfry.OrganizationID, kudgfry.Hash,
		kudgfry.KudguriUuid, kudgfry.TargetDriverType,
		kudgfry.UnkouNo, kudgfry.UnkouDate, kudgfry.ReadDate,
		kudgfry.OfficeCd, kudgfry.OfficeName, kudgfry.VehicleCd, kudgfry.VehicleName,
		kudgfry.DriverCd1, kudgfry.DriverName1, kudgfry.DriverCd2, kudgfry.DriverName2,
		kudgfry.RelevantDatetime, kudgfry.RefuelInspectCategory, kudgfry.RefuelInspectCategoryName,
		kudgfry.RefuelInspectType, kudgfry.RefuelInspectTypeName, kudgfry.RefuelInspectKind,
		kudgfry.RefuelInspectKindName, kudgfry.RefillAmount, kudgfry.OwnOtherType,
		kudgfry.Mileage, kudgfry.MeterValue,
	).Scan(
		&result.UUID, &result.OrganizationID, &result.Hash, &result.Created,
		&result.Deleted, &result.KudguriUuid, &result.TargetDriverType,
		&result.UnkouNo, &result.UnkouDate, &result.ReadDate,
		&result.OfficeCd, &result.OfficeName, &result.VehicleCd, &result.VehicleName,
		&result.DriverCd1, &result.DriverName1, &result.DriverCd2, &result.DriverName2,
		&result.RelevantDatetime, &result.RefuelInspectCategory, &result.RefuelInspectCategoryName,
		&result.RefuelInspectType, &result.RefuelInspectTypeName, &result.RefuelInspectKind,
		&result.RefuelInspectKindName, &result.RefillAmount, &result.OwnOtherType,
		&result.Mileage, &result.MeterValue,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrKudgfryNotFound
		}
		return nil, err
	}

	return &result, nil
}

// Delete soft-deletes a kudgfry record
func (r *KudgfryRepository) Delete(ctx context.Context, uuid string) error {
	deletedAt := time.Now().Format(time.RFC3339)
	query := `
		UPDATE kudgfry
		SET "Deleted" = $2
		WHERE uuid = $1 AND "Deleted" IS NULL
	`

	result, err := r.db.Exec(ctx, query, uuid, deletedAt)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrKudgfryNotFound
	}

	return nil
}

// ListByOrganization retrieves kudgfry records for a specific organization with pagination
func (r *KudgfryRepository) ListByOrganization(ctx context.Context, organizationID string, limit int, offset int) ([]*Kudgfry, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	query := `
		SELECT
			uuid, "OrganizationID", "Hash", "Created", "Deleted", "KudguriUuid",
			"TargetDriverType", "UnkouNo", "UnkouDate", "ReadDate",
			"OfficeCd", "OfficeName", "VehicleCd", "VehicleName",
			"DriverCd1", "DriverName1", "DriverCd2", "DriverName2",
			"RelevantDatetime", "RefuelInspectCategory", "RefuelInspectCategoryName",
			"RefuelInspectType", "RefuelInspectTypeName", "RefuelInspectKind",
			"RefuelInspectKindName", "RefillAmount", "OwnOtherType",
			mileage, "MeterValue"
		FROM kudgfry
		WHERE "OrganizationID" = $1 AND "Deleted" IS NULL
		ORDER BY "Created" DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, organizationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var kudgfrys []*Kudgfry
	for rows.Next() {
		var kudgfry Kudgfry
		err := rows.Scan(
			&kudgfry.UUID, &kudgfry.OrganizationID, &kudgfry.Hash, &kudgfry.Created,
			&kudgfry.Deleted, &kudgfry.KudguriUuid, &kudgfry.TargetDriverType,
			&kudgfry.UnkouNo, &kudgfry.UnkouDate, &kudgfry.ReadDate,
			&kudgfry.OfficeCd, &kudgfry.OfficeName, &kudgfry.VehicleCd, &kudgfry.VehicleName,
			&kudgfry.DriverCd1, &kudgfry.DriverName1, &kudgfry.DriverCd2, &kudgfry.DriverName2,
			&kudgfry.RelevantDatetime, &kudgfry.RefuelInspectCategory, &kudgfry.RefuelInspectCategoryName,
			&kudgfry.RefuelInspectType, &kudgfry.RefuelInspectTypeName, &kudgfry.RefuelInspectKind,
			&kudgfry.RefuelInspectKindName, &kudgfry.RefillAmount, &kudgfry.OwnOtherType,
			&kudgfry.Mileage, &kudgfry.MeterValue,
		)
		if err != nil {
			return nil, err
		}
		kudgfrys = append(kudgfrys, &kudgfry)
	}

	return kudgfrys, rows.Err()
}

// List retrieves all kudgfry records with pagination
func (r *KudgfryRepository) List(ctx context.Context, limit int, offset int) ([]*Kudgfry, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	query := `
		SELECT
			uuid, "OrganizationID", "Hash", "Created", "Deleted", "KudguriUuid",
			"TargetDriverType", "UnkouNo", "UnkouDate", "ReadDate",
			"OfficeCd", "OfficeName", "VehicleCd", "VehicleName",
			"DriverCd1", "DriverName1", "DriverCd2", "DriverName2",
			"RelevantDatetime", "RefuelInspectCategory", "RefuelInspectCategoryName",
			"RefuelInspectType", "RefuelInspectTypeName", "RefuelInspectKind",
			"RefuelInspectKindName", "RefillAmount", "OwnOtherType",
			mileage, "MeterValue"
		FROM kudgfry
		WHERE "Deleted" IS NULL
		ORDER BY "Created" DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var kudgfrys []*Kudgfry
	for rows.Next() {
		var kudgfry Kudgfry
		err := rows.Scan(
			&kudgfry.UUID, &kudgfry.OrganizationID, &kudgfry.Hash, &kudgfry.Created,
			&kudgfry.Deleted, &kudgfry.KudguriUuid, &kudgfry.TargetDriverType,
			&kudgfry.UnkouNo, &kudgfry.UnkouDate, &kudgfry.ReadDate,
			&kudgfry.OfficeCd, &kudgfry.OfficeName, &kudgfry.VehicleCd, &kudgfry.VehicleName,
			&kudgfry.DriverCd1, &kudgfry.DriverName1, &kudgfry.DriverCd2, &kudgfry.DriverName2,
			&kudgfry.RelevantDatetime, &kudgfry.RefuelInspectCategory, &kudgfry.RefuelInspectCategoryName,
			&kudgfry.RefuelInspectType, &kudgfry.RefuelInspectTypeName, &kudgfry.RefuelInspectKind,
			&kudgfry.RefuelInspectKindName, &kudgfry.RefillAmount, &kudgfry.OwnOtherType,
			&kudgfry.Mileage, &kudgfry.MeterValue,
		)
		if err != nil {
			return nil, err
		}
		kudgfrys = append(kudgfrys, &kudgfry)
	}

	return kudgfrys, rows.Err()
}
