package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrKudguriNotFound = errors.New("kudguri not found")
)

// Kudguri represents the database model for kudguri table
type Kudguri struct {
	UUID           string
	OrganizationID string
	Hash           string
	Created        string
	Deleted        *string
	UnkouNo        string
	KudguriUuid    string
	ReadDate       *string
	OfficeCd       *string
	OfficeName     *string
	VehicleCd      *string
	VehicleName    *string
	DriverCd1      *string
	DriverName1    *string
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

// KudguriRepository handles database operations for kudguri
type KudguriRepository struct {
	db DB
}

// NewKudguriRepository creates a new repository
func NewKudguriRepository(pool *pgxpool.Pool) *KudguriRepository {
	return &KudguriRepository{db: pool}
}

// NewKudguriRepositoryWithDB creates a repository with custom DB interface (for testing)
func NewKudguriRepositoryWithDB(db DB) *KudguriRepository {
	return &KudguriRepository{db: db}
}

// Create inserts a new kudguri record
func (r *KudguriRepository) Create(ctx context.Context, k *Kudguri) (*Kudguri, error) {
	if k.UUID == "" {
		k.UUID = uuid.New().String()
	}

	query := `
		INSERT INTO kudguri (
			uuid, "OrganizationID", "Hash", "Created", "Deleted",
			"unkouNo", "kudguriUuid", "ReadDate", "OfficeCd", "OfficeName",
			"VehicleCd", "VehicleName", "DriverCd1", "DriverName1", "TargetDriverType",
			"TargetDriverCd", "TargetDriverName", "StartDatetime", "EndDatetime", "EventCd",
			"EventName", "StartMileage", "EndMileage", "SectionTime", "SectionDistance",
			"StartCityCd", "StartCityName", "EndCityCd", "EndCityName", "StartPlaceCd",
			"StartPlaceName", "EndPlaceCd", "EndPlaceName", "StartGpsValid", "StartGpsLat",
			"StartGpsLng", "EndGpsValid", "EndGpsLat", "EndGpsLng", "OverLimitMax"
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
			$21, $22, $23, $24, $25, $26, $27, $28, $29, $30,
			$31, $32, $33, $34, $35, $36, $37, $38, $39, $40
		)
		RETURNING uuid, "OrganizationID", "Hash", "Created", "Deleted",
			"unkouNo", "kudguriUuid", "ReadDate", "OfficeCd", "OfficeName",
			"VehicleCd", "VehicleName", "DriverCd1", "DriverName1", "TargetDriverType",
			"TargetDriverCd", "TargetDriverName", "StartDatetime", "EndDatetime", "EventCd",
			"EventName", "StartMileage", "EndMileage", "SectionTime", "SectionDistance",
			"StartCityCd", "StartCityName", "EndCityCd", "EndCityName", "StartPlaceCd",
			"StartPlaceName", "EndPlaceCd", "EndPlaceName", "StartGpsValid", "StartGpsLat",
			"StartGpsLng", "EndGpsValid", "EndGpsLat", "EndGpsLng", "OverLimitMax"
	`

	var result Kudguri
	err := r.db.QueryRow(ctx, query,
		k.UUID, k.OrganizationID, k.Hash, k.Created, k.Deleted,
		k.UnkouNo, k.KudguriUuid, k.ReadDate, k.OfficeCd, k.OfficeName,
		k.VehicleCd, k.VehicleName, k.DriverCd1, k.DriverName1, k.TargetDriverType,
		k.TargetDriverCd, k.TargetDriverName, k.StartDatetime, k.EndDatetime, k.EventCd,
		k.EventName, k.StartMileage, k.EndMileage, k.SectionTime, k.SectionDistance,
		k.StartCityCd, k.StartCityName, k.EndCityCd, k.EndCityName, k.StartPlaceCd,
		k.StartPlaceName, k.EndPlaceCd, k.EndPlaceName, k.StartGpsValid, k.StartGpsLat,
		k.StartGpsLng, k.EndGpsValid, k.EndGpsLat, k.EndGpsLng, k.OverLimitMax,
	).Scan(
		&result.UUID, &result.OrganizationID, &result.Hash, &result.Created, &result.Deleted,
		&result.UnkouNo, &result.KudguriUuid, &result.ReadDate, &result.OfficeCd, &result.OfficeName,
		&result.VehicleCd, &result.VehicleName, &result.DriverCd1, &result.DriverName1, &result.TargetDriverType,
		&result.TargetDriverCd, &result.TargetDriverName, &result.StartDatetime, &result.EndDatetime, &result.EventCd,
		&result.EventName, &result.StartMileage, &result.EndMileage, &result.SectionTime, &result.SectionDistance,
		&result.StartCityCd, &result.StartCityName, &result.EndCityCd, &result.EndCityName, &result.StartPlaceCd,
		&result.StartPlaceName, &result.EndPlaceCd, &result.EndPlaceName, &result.StartGpsValid, &result.StartGpsLat,
		&result.StartGpsLng, &result.EndGpsValid, &result.EndGpsLat, &result.EndGpsLng, &result.OverLimitMax,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetByUUID retrieves a kudguri record by UUID
func (r *KudguriRepository) GetByUUID(ctx context.Context, uuid string) (*Kudguri, error) {
	query := `
		SELECT uuid, "OrganizationID", "Hash", "Created", "Deleted",
			"unkouNo", "kudguriUuid", "ReadDate", "OfficeCd", "OfficeName",
			"VehicleCd", "VehicleName", "DriverCd1", "DriverName1", "TargetDriverType",
			"TargetDriverCd", "TargetDriverName", "StartDatetime", "EndDatetime", "EventCd",
			"EventName", "StartMileage", "EndMileage", "SectionTime", "SectionDistance",
			"StartCityCd", "StartCityName", "EndCityCd", "EndCityName", "StartPlaceCd",
			"StartPlaceName", "EndPlaceCd", "EndPlaceName", "StartGpsValid", "StartGpsLat",
			"StartGpsLng", "EndGpsValid", "EndGpsLat", "EndGpsLng", "OverLimitMax"
		FROM kudguri
		WHERE uuid = $1 AND "Deleted" IS NULL
	`

	var k Kudguri
	err := r.db.QueryRow(ctx, query, uuid).Scan(
		&k.UUID, &k.OrganizationID, &k.Hash, &k.Created, &k.Deleted,
		&k.UnkouNo, &k.KudguriUuid, &k.ReadDate, &k.OfficeCd, &k.OfficeName,
		&k.VehicleCd, &k.VehicleName, &k.DriverCd1, &k.DriverName1, &k.TargetDriverType,
		&k.TargetDriverCd, &k.TargetDriverName, &k.StartDatetime, &k.EndDatetime, &k.EventCd,
		&k.EventName, &k.StartMileage, &k.EndMileage, &k.SectionTime, &k.SectionDistance,
		&k.StartCityCd, &k.StartCityName, &k.EndCityCd, &k.EndCityName, &k.StartPlaceCd,
		&k.StartPlaceName, &k.EndPlaceCd, &k.EndPlaceName, &k.StartGpsValid, &k.StartGpsLat,
		&k.StartGpsLng, &k.EndGpsValid, &k.EndGpsLat, &k.EndGpsLng, &k.OverLimitMax,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrKudguriNotFound
		}
		return nil, err
	}

	return &k, nil
}

// Update modifies an existing kudguri record
func (r *KudguriRepository) Update(ctx context.Context, k *Kudguri) (*Kudguri, error) {
	query := `
		UPDATE kudguri
		SET "OrganizationID" = $2, "Hash" = $3, "Created" = $4, "Deleted" = $5,
			"unkouNo" = $6, "kudguriUuid" = $7, "ReadDate" = $8, "OfficeCd" = $9, "OfficeName" = $10,
			"VehicleCd" = $11, "VehicleName" = $12, "DriverCd1" = $13, "DriverName1" = $14, "TargetDriverType" = $15,
			"TargetDriverCd" = $16, "TargetDriverName" = $17, "StartDatetime" = $18, "EndDatetime" = $19, "EventCd" = $20,
			"EventName" = $21, "StartMileage" = $22, "EndMileage" = $23, "SectionTime" = $24, "SectionDistance" = $25,
			"StartCityCd" = $26, "StartCityName" = $27, "EndCityCd" = $28, "EndCityName" = $29, "StartPlaceCd" = $30,
			"StartPlaceName" = $31, "EndPlaceCd" = $32, "EndPlaceName" = $33, "StartGpsValid" = $34, "StartGpsLat" = $35,
			"StartGpsLng" = $36, "EndGpsValid" = $37, "EndGpsLat" = $38, "EndGpsLng" = $39, "OverLimitMax" = $40
		WHERE uuid = $1 AND "Deleted" IS NULL
		RETURNING uuid, "OrganizationID", "Hash", "Created", "Deleted",
			"unkouNo", "kudguriUuid", "ReadDate", "OfficeCd", "OfficeName",
			"VehicleCd", "VehicleName", "DriverCd1", "DriverName1", "TargetDriverType",
			"TargetDriverCd", "TargetDriverName", "StartDatetime", "EndDatetime", "EventCd",
			"EventName", "StartMileage", "EndMileage", "SectionTime", "SectionDistance",
			"StartCityCd", "StartCityName", "EndCityCd", "EndCityName", "StartPlaceCd",
			"StartPlaceName", "EndPlaceCd", "EndPlaceName", "StartGpsValid", "StartGpsLat",
			"StartGpsLng", "EndGpsValid", "EndGpsLat", "EndGpsLng", "OverLimitMax"
	`

	var result Kudguri
	err := r.db.QueryRow(ctx, query,
		k.UUID, k.OrganizationID, k.Hash, k.Created, k.Deleted,
		k.UnkouNo, k.KudguriUuid, k.ReadDate, k.OfficeCd, k.OfficeName,
		k.VehicleCd, k.VehicleName, k.DriverCd1, k.DriverName1, k.TargetDriverType,
		k.TargetDriverCd, k.TargetDriverName, k.StartDatetime, k.EndDatetime, k.EventCd,
		k.EventName, k.StartMileage, k.EndMileage, k.SectionTime, k.SectionDistance,
		k.StartCityCd, k.StartCityName, k.EndCityCd, k.EndCityName, k.StartPlaceCd,
		k.StartPlaceName, k.EndPlaceCd, k.EndPlaceName, k.StartGpsValid, k.StartGpsLat,
		k.StartGpsLng, k.EndGpsValid, k.EndGpsLat, k.EndGpsLng, k.OverLimitMax,
	).Scan(
		&result.UUID, &result.OrganizationID, &result.Hash, &result.Created, &result.Deleted,
		&result.UnkouNo, &result.KudguriUuid, &result.ReadDate, &result.OfficeCd, &result.OfficeName,
		&result.VehicleCd, &result.VehicleName, &result.DriverCd1, &result.DriverName1, &result.TargetDriverType,
		&result.TargetDriverCd, &result.TargetDriverName, &result.StartDatetime, &result.EndDatetime, &result.EventCd,
		&result.EventName, &result.StartMileage, &result.EndMileage, &result.SectionTime, &result.SectionDistance,
		&result.StartCityCd, &result.StartCityName, &result.EndCityCd, &result.EndCityName, &result.StartPlaceCd,
		&result.StartPlaceName, &result.EndPlaceCd, &result.EndPlaceName, &result.StartGpsValid, &result.StartGpsLat,
		&result.StartGpsLng, &result.EndGpsValid, &result.EndGpsLat, &result.EndGpsLng, &result.OverLimitMax,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrKudguriNotFound
		}
		return nil, err
	}

	return &result, nil
}

// Delete soft-deletes a kudguri record
func (r *KudguriRepository) Delete(ctx context.Context, uuid, deletedTimestamp string) error {
	query := `
		UPDATE kudguri
		SET "Deleted" = $2
		WHERE uuid = $1 AND "Deleted" IS NULL
	`

	result, err := r.db.Exec(ctx, query, uuid, deletedTimestamp)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrKudguriNotFound
	}

	return nil
}

// ListByOrganization retrieves kudguri records for a specific organization with pagination
func (r *KudguriRepository) ListByOrganization(ctx context.Context, organizationID string, limit int, offset int) ([]*Kudguri, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	query := `
		SELECT uuid, "OrganizationID", "Hash", "Created", "Deleted",
			"unkouNo", "kudguriUuid", "ReadDate", "OfficeCd", "OfficeName",
			"VehicleCd", "VehicleName", "DriverCd1", "DriverName1", "TargetDriverType",
			"TargetDriverCd", "TargetDriverName", "StartDatetime", "EndDatetime", "EventCd",
			"EventName", "StartMileage", "EndMileage", "SectionTime", "SectionDistance",
			"StartCityCd", "StartCityName", "EndCityCd", "EndCityName", "StartPlaceCd",
			"StartPlaceName", "EndPlaceCd", "EndPlaceName", "StartGpsValid", "StartGpsLat",
			"StartGpsLng", "EndGpsValid", "EndGpsLat", "EndGpsLng", "OverLimitMax"
		FROM kudguri
		WHERE "OrganizationID" = $1 AND "Deleted" IS NULL
		ORDER BY "Created" DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, organizationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var kudguriList []*Kudguri
	for rows.Next() {
		var k Kudguri
		err := rows.Scan(
			&k.UUID, &k.OrganizationID, &k.Hash, &k.Created, &k.Deleted,
			&k.UnkouNo, &k.KudguriUuid, &k.ReadDate, &k.OfficeCd, &k.OfficeName,
			&k.VehicleCd, &k.VehicleName, &k.DriverCd1, &k.DriverName1, &k.TargetDriverType,
			&k.TargetDriverCd, &k.TargetDriverName, &k.StartDatetime, &k.EndDatetime, &k.EventCd,
			&k.EventName, &k.StartMileage, &k.EndMileage, &k.SectionTime, &k.SectionDistance,
			&k.StartCityCd, &k.StartCityName, &k.EndCityCd, &k.EndCityName, &k.StartPlaceCd,
			&k.StartPlaceName, &k.EndPlaceCd, &k.EndPlaceName, &k.StartGpsValid, &k.StartGpsLat,
			&k.StartGpsLng, &k.EndGpsValid, &k.EndGpsLat, &k.EndGpsLng, &k.OverLimitMax,
		)
		if err != nil {
			return nil, err
		}
		kudguriList = append(kudguriList, &k)
	}

	return kudguriList, rows.Err()
}

// List retrieves all kudguri records with pagination
func (r *KudguriRepository) List(ctx context.Context, limit int, offset int) ([]*Kudguri, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	query := `
		SELECT uuid, "OrganizationID", "Hash", "Created", "Deleted",
			"unkouNo", "kudguriUuid", "ReadDate", "OfficeCd", "OfficeName",
			"VehicleCd", "VehicleName", "DriverCd1", "DriverName1", "TargetDriverType",
			"TargetDriverCd", "TargetDriverName", "StartDatetime", "EndDatetime", "EventCd",
			"EventName", "StartMileage", "EndMileage", "SectionTime", "SectionDistance",
			"StartCityCd", "StartCityName", "EndCityCd", "EndCityName", "StartPlaceCd",
			"StartPlaceName", "EndPlaceCd", "EndPlaceName", "StartGpsValid", "StartGpsLat",
			"StartGpsLng", "EndGpsValid", "EndGpsLat", "EndGpsLng", "OverLimitMax"
		FROM kudguri
		WHERE "Deleted" IS NULL
		ORDER BY "Created" DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var kudguriList []*Kudguri
	for rows.Next() {
		var k Kudguri
		err := rows.Scan(
			&k.UUID, &k.OrganizationID, &k.Hash, &k.Created, &k.Deleted,
			&k.UnkouNo, &k.KudguriUuid, &k.ReadDate, &k.OfficeCd, &k.OfficeName,
			&k.VehicleCd, &k.VehicleName, &k.DriverCd1, &k.DriverName1, &k.TargetDriverType,
			&k.TargetDriverCd, &k.TargetDriverName, &k.StartDatetime, &k.EndDatetime, &k.EventCd,
			&k.EventName, &k.StartMileage, &k.EndMileage, &k.SectionTime, &k.SectionDistance,
			&k.StartCityCd, &k.StartCityName, &k.EndCityCd, &k.EndCityName, &k.StartPlaceCd,
			&k.StartPlaceName, &k.EndPlaceCd, &k.EndPlaceName, &k.StartGpsValid, &k.StartGpsLat,
			&k.StartGpsLng, &k.EndGpsValid, &k.EndGpsLat, &k.EndGpsLng, &k.OverLimitMax,
		)
		if err != nil {
			return nil, err
		}
		kudguriList = append(kudguriList, &k)
	}

	return kudguriList, rows.Err()
}
