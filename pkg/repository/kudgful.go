package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrKudgfulNotFound = errors.New("kudgful not found")
)

// Kudgful represents the database model
type Kudgful struct {
	UUID             string
	OrganizationID   string
	Hash             string
	Created          string
	Deleted          *string
	KudguriUuid      *string
	UnkouNo          *string
	ReadDate         *string
	OfficeCd         *string
	OfficeName       *string
	VehicleCd        *string
	VehicleName      *string
	DriverCd1        *string
	DriverName1      *string
	TargetDriverType string
	TargetDriverCd   *string
	TargetDriverName *string
	StartDatetime    *string
	EndDatetime      *string
	EventCd          *string
	EventName        *string
	StartMileage     *string
	EndMileage       *string
	SectionTime      *string
	SectionDistance  *string
	StartCityCd      *string
	StartCityName    *string
	EndCityCd        *string
	EndCityName      *string
	StartPlaceCd     *string
	StartPlaceName   *string
	EndPlaceCd       *string
	EndPlaceName     *string
	StartGpsValid    *string
	StartGpsLat      *string
	StartGpsLng      *string
	EndGpsValid      *string
	EndGpsLat        *string
	EndGpsLng        *string
	OverLimitMax     *string
}

// KudgfulRepository handles database operations for kudgful
type KudgfulRepository struct {
	db DB
}

// NewKudgfulRepository creates a new repository
func NewKudgfulRepository(pool *pgxpool.Pool) *KudgfulRepository {
	return &KudgfulRepository{db: pool}
}

// NewKudgfulRepositoryWithDB creates a repository with custom DB interface (for testing)
func NewKudgfulRepositoryWithDB(db DB) *KudgfulRepository {
	return &KudgfulRepository{db: db}
}

// Create inserts a new kudgful record
func (r *KudgfulRepository) Create(ctx context.Context, kudgful *Kudgful) (*Kudgful, error) {
	if kudgful.UUID == "" {
		kudgful.UUID = uuid.New().String()
	}

	query := `
		INSERT INTO kudgful (
			"uuid", "OrganizationID", "Hash", "Created", "Deleted", "KudguriUuid",
			"UnkouNo", "ReadDate", "OfficeCd", "OfficeName", "VehicleCd", "VehicleName",
			"DriverCd1", "DriverName1", "TargetDriverType", "TargetDriverCd", "TargetDriverName",
			"StartDatetime", "EndDatetime", "EventCd", "EventName", "StartMileage", "EndMileage",
			"SectionTime", "SectionDistance", "StartCityCd", "StartCityName", "EndCityCd", "EndCityName",
			"StartPlaceCd", "StartPlaceName", "EndPlaceCd", "EndPlaceName",
			"StartGpsValid", "StartGpsLat", "StartGpsLng", "EndGpsValid", "EndGpsLat", "EndGpsLng",
			"OverLimitMax"
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
			$21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40
		)
		RETURNING
			"uuid", "OrganizationID", "Hash", "Created", "Deleted", "KudguriUuid",
			"UnkouNo", "ReadDate", "OfficeCd", "OfficeName", "VehicleCd", "VehicleName",
			"DriverCd1", "DriverName1", "TargetDriverType", "TargetDriverCd", "TargetDriverName",
			"StartDatetime", "EndDatetime", "EventCd", "EventName", "StartMileage", "EndMileage",
			"SectionTime", "SectionDistance", "StartCityCd", "StartCityName", "EndCityCd", "EndCityName",
			"StartPlaceCd", "StartPlaceName", "EndPlaceCd", "EndPlaceName",
			"StartGpsValid", "StartGpsLat", "StartGpsLng", "EndGpsValid", "EndGpsLat", "EndGpsLng",
			"OverLimitMax"
	`

	var result Kudgful
	err := r.db.QueryRow(ctx, query,
		kudgful.UUID, kudgful.OrganizationID, kudgful.Hash, kudgful.Created, kudgful.Deleted, kudgful.KudguriUuid,
		kudgful.UnkouNo, kudgful.ReadDate, kudgful.OfficeCd, kudgful.OfficeName, kudgful.VehicleCd, kudgful.VehicleName,
		kudgful.DriverCd1, kudgful.DriverName1, kudgful.TargetDriverType, kudgful.TargetDriverCd, kudgful.TargetDriverName,
		kudgful.StartDatetime, kudgful.EndDatetime, kudgful.EventCd, kudgful.EventName, kudgful.StartMileage, kudgful.EndMileage,
		kudgful.SectionTime, kudgful.SectionDistance, kudgful.StartCityCd, kudgful.StartCityName, kudgful.EndCityCd, kudgful.EndCityName,
		kudgful.StartPlaceCd, kudgful.StartPlaceName, kudgful.EndPlaceCd, kudgful.EndPlaceName,
		kudgful.StartGpsValid, kudgful.StartGpsLat, kudgful.StartGpsLng, kudgful.EndGpsValid, kudgful.EndGpsLat, kudgful.EndGpsLng,
		kudgful.OverLimitMax,
	).Scan(
		&result.UUID, &result.OrganizationID, &result.Hash, &result.Created, &result.Deleted, &result.KudguriUuid,
		&result.UnkouNo, &result.ReadDate, &result.OfficeCd, &result.OfficeName, &result.VehicleCd, &result.VehicleName,
		&result.DriverCd1, &result.DriverName1, &result.TargetDriverType, &result.TargetDriverCd, &result.TargetDriverName,
		&result.StartDatetime, &result.EndDatetime, &result.EventCd, &result.EventName, &result.StartMileage, &result.EndMileage,
		&result.SectionTime, &result.SectionDistance, &result.StartCityCd, &result.StartCityName, &result.EndCityCd, &result.EndCityName,
		&result.StartPlaceCd, &result.StartPlaceName, &result.EndPlaceCd, &result.EndPlaceName,
		&result.StartGpsValid, &result.StartGpsLat, &result.StartGpsLng, &result.EndGpsValid, &result.EndGpsLat, &result.EndGpsLng,
		&result.OverLimitMax,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetByUUID retrieves a kudgful record by UUID
func (r *KudgfulRepository) GetByUUID(ctx context.Context, id string) (*Kudgful, error) {
	query := `
		SELECT
			"uuid", "OrganizationID", "Hash", "Created", "Deleted", "KudguriUuid",
			"UnkouNo", "ReadDate", "OfficeCd", "OfficeName", "VehicleCd", "VehicleName",
			"DriverCd1", "DriverName1", "TargetDriverType", "TargetDriverCd", "TargetDriverName",
			"StartDatetime", "EndDatetime", "EventCd", "EventName", "StartMileage", "EndMileage",
			"SectionTime", "SectionDistance", "StartCityCd", "StartCityName", "EndCityCd", "EndCityName",
			"StartPlaceCd", "StartPlaceName", "EndPlaceCd", "EndPlaceName",
			"StartGpsValid", "StartGpsLat", "StartGpsLng", "EndGpsValid", "EndGpsLat", "EndGpsLng",
			"OverLimitMax"
		FROM kudgful
		WHERE "uuid" = $1 AND "Deleted" IS NULL
	`

	var result Kudgful
	err := r.db.QueryRow(ctx, query, id).Scan(
		&result.UUID, &result.OrganizationID, &result.Hash, &result.Created, &result.Deleted, &result.KudguriUuid,
		&result.UnkouNo, &result.ReadDate, &result.OfficeCd, &result.OfficeName, &result.VehicleCd, &result.VehicleName,
		&result.DriverCd1, &result.DriverName1, &result.TargetDriverType, &result.TargetDriverCd, &result.TargetDriverName,
		&result.StartDatetime, &result.EndDatetime, &result.EventCd, &result.EventName, &result.StartMileage, &result.EndMileage,
		&result.SectionTime, &result.SectionDistance, &result.StartCityCd, &result.StartCityName, &result.EndCityCd, &result.EndCityName,
		&result.StartPlaceCd, &result.StartPlaceName, &result.EndPlaceCd, &result.EndPlaceName,
		&result.StartGpsValid, &result.StartGpsLat, &result.StartGpsLng, &result.EndGpsValid, &result.EndGpsLat, &result.EndGpsLng,
		&result.OverLimitMax,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrKudgfulNotFound
		}
		return nil, err
	}

	return &result, nil
}

// Update modifies an existing kudgful record
func (r *KudgfulRepository) Update(ctx context.Context, kudgful *Kudgful) (*Kudgful, error) {
	query := `
		UPDATE kudgful
		SET
			"OrganizationID" = $2, "Hash" = $3, "Created" = $4, "Deleted" = $5, "KudguriUuid" = $6,
			"UnkouNo" = $7, "ReadDate" = $8, "OfficeCd" = $9, "OfficeName" = $10, "VehicleCd" = $11, "VehicleName" = $12,
			"DriverCd1" = $13, "DriverName1" = $14, "TargetDriverType" = $15, "TargetDriverCd" = $16, "TargetDriverName" = $17,
			"StartDatetime" = $18, "EndDatetime" = $19, "EventCd" = $20, "EventName" = $21, "StartMileage" = $22, "EndMileage" = $23,
			"SectionTime" = $24, "SectionDistance" = $25, "StartCityCd" = $26, "StartCityName" = $27, "EndCityCd" = $28, "EndCityName" = $29,
			"StartPlaceCd" = $30, "StartPlaceName" = $31, "EndPlaceCd" = $32, "EndPlaceName" = $33,
			"StartGpsValid" = $34, "StartGpsLat" = $35, "StartGpsLng" = $36, "EndGpsValid" = $37, "EndGpsLat" = $38, "EndGpsLng" = $39,
			"OverLimitMax" = $40
		WHERE "uuid" = $1 AND "Deleted" IS NULL
		RETURNING
			"uuid", "OrganizationID", "Hash", "Created", "Deleted", "KudguriUuid",
			"UnkouNo", "ReadDate", "OfficeCd", "OfficeName", "VehicleCd", "VehicleName",
			"DriverCd1", "DriverName1", "TargetDriverType", "TargetDriverCd", "TargetDriverName",
			"StartDatetime", "EndDatetime", "EventCd", "EventName", "StartMileage", "EndMileage",
			"SectionTime", "SectionDistance", "StartCityCd", "StartCityName", "EndCityCd", "EndCityName",
			"StartPlaceCd", "StartPlaceName", "EndPlaceCd", "EndPlaceName",
			"StartGpsValid", "StartGpsLat", "StartGpsLng", "EndGpsValid", "EndGpsLat", "EndGpsLng",
			"OverLimitMax"
	`

	var result Kudgful
	err := r.db.QueryRow(ctx, query,
		kudgful.UUID, kudgful.OrganizationID, kudgful.Hash, kudgful.Created, kudgful.Deleted, kudgful.KudguriUuid,
		kudgful.UnkouNo, kudgful.ReadDate, kudgful.OfficeCd, kudgful.OfficeName, kudgful.VehicleCd, kudgful.VehicleName,
		kudgful.DriverCd1, kudgful.DriverName1, kudgful.TargetDriverType, kudgful.TargetDriverCd, kudgful.TargetDriverName,
		kudgful.StartDatetime, kudgful.EndDatetime, kudgful.EventCd, kudgful.EventName, kudgful.StartMileage, kudgful.EndMileage,
		kudgful.SectionTime, kudgful.SectionDistance, kudgful.StartCityCd, kudgful.StartCityName, kudgful.EndCityCd, kudgful.EndCityName,
		kudgful.StartPlaceCd, kudgful.StartPlaceName, kudgful.EndPlaceCd, kudgful.EndPlaceName,
		kudgful.StartGpsValid, kudgful.StartGpsLat, kudgful.StartGpsLng, kudgful.EndGpsValid, kudgful.EndGpsLat, kudgful.EndGpsLng,
		kudgful.OverLimitMax,
	).Scan(
		&result.UUID, &result.OrganizationID, &result.Hash, &result.Created, &result.Deleted, &result.KudguriUuid,
		&result.UnkouNo, &result.ReadDate, &result.OfficeCd, &result.OfficeName, &result.VehicleCd, &result.VehicleName,
		&result.DriverCd1, &result.DriverName1, &result.TargetDriverType, &result.TargetDriverCd, &result.TargetDriverName,
		&result.StartDatetime, &result.EndDatetime, &result.EventCd, &result.EventName, &result.StartMileage, &result.EndMileage,
		&result.SectionTime, &result.SectionDistance, &result.StartCityCd, &result.StartCityName, &result.EndCityCd, &result.EndCityName,
		&result.StartPlaceCd, &result.StartPlaceName, &result.EndPlaceCd, &result.EndPlaceName,
		&result.StartGpsValid, &result.StartGpsLat, &result.StartGpsLng, &result.EndGpsValid, &result.EndGpsLat, &result.EndGpsLng,
		&result.OverLimitMax,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrKudgfulNotFound
		}
		return nil, err
	}

	return &result, nil
}

// Delete soft-deletes a kudgful record
func (r *KudgfulRepository) Delete(ctx context.Context, id string, deletedTimestamp string) error {
	query := `
		UPDATE kudgful
		SET "Deleted" = $2
		WHERE "uuid" = $1 AND "Deleted" IS NULL
	`

	result, err := r.db.Exec(ctx, query, id, deletedTimestamp)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrKudgfulNotFound
	}

	return nil
}

// ListByOrganization retrieves kudgful records by organization with pagination
func (r *KudgfulRepository) ListByOrganization(ctx context.Context, organizationID string, limit int, offset int) ([]*Kudgful, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	query := `
		SELECT
			"uuid", "OrganizationID", "Hash", "Created", "Deleted", "KudguriUuid",
			"UnkouNo", "ReadDate", "OfficeCd", "OfficeName", "VehicleCd", "VehicleName",
			"DriverCd1", "DriverName1", "TargetDriverType", "TargetDriverCd", "TargetDriverName",
			"StartDatetime", "EndDatetime", "EventCd", "EventName", "StartMileage", "EndMileage",
			"SectionTime", "SectionDistance", "StartCityCd", "StartCityName", "EndCityCd", "EndCityName",
			"StartPlaceCd", "StartPlaceName", "EndPlaceCd", "EndPlaceName",
			"StartGpsValid", "StartGpsLat", "StartGpsLng", "EndGpsValid", "EndGpsLat", "EndGpsLng",
			"OverLimitMax"
		FROM kudgful
		WHERE "OrganizationID" = $1 AND "Deleted" IS NULL
		ORDER BY "Created" DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, organizationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*Kudgful
	for rows.Next() {
		var result Kudgful
		err := rows.Scan(
			&result.UUID, &result.OrganizationID, &result.Hash, &result.Created, &result.Deleted, &result.KudguriUuid,
			&result.UnkouNo, &result.ReadDate, &result.OfficeCd, &result.OfficeName, &result.VehicleCd, &result.VehicleName,
			&result.DriverCd1, &result.DriverName1, &result.TargetDriverType, &result.TargetDriverCd, &result.TargetDriverName,
			&result.StartDatetime, &result.EndDatetime, &result.EventCd, &result.EventName, &result.StartMileage, &result.EndMileage,
			&result.SectionTime, &result.SectionDistance, &result.StartCityCd, &result.StartCityName, &result.EndCityCd, &result.EndCityName,
			&result.StartPlaceCd, &result.StartPlaceName, &result.EndPlaceCd, &result.EndPlaceName,
			&result.StartGpsValid, &result.StartGpsLat, &result.StartGpsLng, &result.EndGpsValid, &result.EndGpsLat, &result.EndGpsLng,
			&result.OverLimitMax,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, &result)
	}

	return results, rows.Err()
}

// List retrieves kudgful records with pagination
func (r *KudgfulRepository) List(ctx context.Context, limit int, offset int) ([]*Kudgful, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	query := `
		SELECT
			"uuid", "OrganizationID", "Hash", "Created", "Deleted", "KudguriUuid",
			"UnkouNo", "ReadDate", "OfficeCd", "OfficeName", "VehicleCd", "VehicleName",
			"DriverCd1", "DriverName1", "TargetDriverType", "TargetDriverCd", "TargetDriverName",
			"StartDatetime", "EndDatetime", "EventCd", "EventName", "StartMileage", "EndMileage",
			"SectionTime", "SectionDistance", "StartCityCd", "StartCityName", "EndCityCd", "EndCityName",
			"StartPlaceCd", "StartPlaceName", "EndPlaceCd", "EndPlaceName",
			"StartGpsValid", "StartGpsLat", "StartGpsLng", "EndGpsValid", "EndGpsLat", "EndGpsLng",
			"OverLimitMax"
		FROM kudgful
		WHERE "Deleted" IS NULL
		ORDER BY "Created" DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*Kudgful
	for rows.Next() {
		var result Kudgful
		err := rows.Scan(
			&result.UUID, &result.OrganizationID, &result.Hash, &result.Created, &result.Deleted, &result.KudguriUuid,
			&result.UnkouNo, &result.ReadDate, &result.OfficeCd, &result.OfficeName, &result.VehicleCd, &result.VehicleName,
			&result.DriverCd1, &result.DriverName1, &result.TargetDriverType, &result.TargetDriverCd, &result.TargetDriverName,
			&result.StartDatetime, &result.EndDatetime, &result.EventCd, &result.EventName, &result.StartMileage, &result.EndMileage,
			&result.SectionTime, &result.SectionDistance, &result.StartCityCd, &result.StartCityName, &result.EndCityCd, &result.EndCityName,
			&result.StartPlaceCd, &result.StartPlaceName, &result.EndPlaceCd, &result.EndPlaceName,
			&result.StartGpsValid, &result.StartGpsLat, &result.StartGpsLng, &result.EndGpsValid, &result.EndGpsLat, &result.EndGpsLng,
			&result.OverLimitMax,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, &result)
	}

	return results, rows.Err()
}
