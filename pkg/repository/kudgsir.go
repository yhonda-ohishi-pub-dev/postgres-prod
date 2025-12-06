package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrKudgsirNotFound = errors.New("kudgsir not found")
)

// Kudgsir represents the database model for kudgsir table
type Kudgsir struct {
	UUID           string
	OrganizationID string
	Hash           string
	Created        string
	Deleted        *string
	KudguriUuid    *string

	// 運行基本情報
	UnkouNo  *string
	ReadDate *string

	// 営業所・車両情報
	OfficeCd    *string
	OfficeName  *string
	VehicleCd   *string
	VehicleName *string

	// ドライバー情報
	DriverCd1        *string
	DriverName1      *string
	TargetDriverType string
	TargetDriverCd   *string
	TargetDriverName *string

	// 区間情報
	StartDatetime *string
	EndDatetime   *string
	EventCd       *string
	EventName     *string

	// 距離・時間情報
	StartMileage    *string
	EndMileage      *string
	SectionTime     *string
	SectionDistance *string

	// 市区町村情報
	StartCityCd   *string
	StartCityName *string
	EndCityCd     *string
	EndCityName   *string

	// 場所情報
	StartPlaceCd   *string
	StartPlaceName *string
	EndPlaceCd     *string
	EndPlaceName   *string

	// GPS情報 (開始地点)
	StartGpsValid *string
	StartGpsLat   *string
	StartGpsLng   *string

	// GPS情報 (終了地点)
	EndGpsValid *string
	EndGpsLat   *string
	EndGpsLng   *string

	// 速度超過
	OverLimitMax *string
}

// KudgsirRepository handles database operations for kudgsir
type KudgsirRepository struct {
	db DB
}

// NewKudgsirRepository creates a new repository
func NewKudgsirRepository(pool *pgxpool.Pool) *KudgsirRepository {
	return &KudgsirRepository{db: pool}
}

// NewKudgsirRepositoryWithDB creates a repository with custom DB interface (for testing)
func NewKudgsirRepositoryWithDB(db DB) *KudgsirRepository {
	return &KudgsirRepository{db: db}
}

// Create inserts a new kudgsir record
func (r *KudgsirRepository) Create(ctx context.Context, k *Kudgsir) (*Kudgsir, error) {
	id := uuid.New().String()
	k.UUID = id

	query := `
		INSERT INTO kudgsir (
			uuid, organization_id, hash, created, deleted, kudguri_uuid,
			unkou_no, read_date,
			office_cd, office_name, vehicle_cd, vehicle_name,
			driver_cd_1, driver_name_1, target_driver_type, target_driver_cd, target_driver_name,
			start_datetime, end_datetime, event_cd, event_name,
			start_mileage, end_mileage, section_time, section_distance,
			start_city_cd, start_city_name, end_city_cd, end_city_name,
			start_place_cd, start_place_name, end_place_cd, end_place_name,
			start_gps_valid, start_gps_lat, start_gps_lng,
			end_gps_valid, end_gps_lat, end_gps_lng,
			over_limit_max
		)
		VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8,
			$9, $10, $11, $12,
			$13, $14, $15, $16, $17,
			$18, $19, $20, $21,
			$22, $23, $24, $25,
			$26, $27, $28, $29,
			$30, $31, $32, $33,
			$34, $35, $36,
			$37, $38, $39,
			$40
		)
		RETURNING
			uuid, organization_id, hash, created, deleted, kudguri_uuid,
			unkou_no, read_date,
			office_cd, office_name, vehicle_cd, vehicle_name,
			driver_cd_1, driver_name_1, target_driver_type, target_driver_cd, target_driver_name,
			start_datetime, end_datetime, event_cd, event_name,
			start_mileage, end_mileage, section_time, section_distance,
			start_city_cd, start_city_name, end_city_cd, end_city_name,
			start_place_cd, start_place_name, end_place_cd, end_place_name,
			start_gps_valid, start_gps_lat, start_gps_lng,
			end_gps_valid, end_gps_lat, end_gps_lng,
			over_limit_max
	`

	var result Kudgsir
	err := r.db.QueryRow(ctx, query,
		k.UUID, k.OrganizationID, k.Hash, k.Created, k.Deleted, k.KudguriUuid,
		k.UnkouNo, k.ReadDate,
		k.OfficeCd, k.OfficeName, k.VehicleCd, k.VehicleName,
		k.DriverCd1, k.DriverName1, k.TargetDriverType, k.TargetDriverCd, k.TargetDriverName,
		k.StartDatetime, k.EndDatetime, k.EventCd, k.EventName,
		k.StartMileage, k.EndMileage, k.SectionTime, k.SectionDistance,
		k.StartCityCd, k.StartCityName, k.EndCityCd, k.EndCityName,
		k.StartPlaceCd, k.StartPlaceName, k.EndPlaceCd, k.EndPlaceName,
		k.StartGpsValid, k.StartGpsLat, k.StartGpsLng,
		k.EndGpsValid, k.EndGpsLat, k.EndGpsLng,
		k.OverLimitMax,
	).Scan(
		&result.UUID, &result.OrganizationID, &result.Hash, &result.Created, &result.Deleted, &result.KudguriUuid,
		&result.UnkouNo, &result.ReadDate,
		&result.OfficeCd, &result.OfficeName, &result.VehicleCd, &result.VehicleName,
		&result.DriverCd1, &result.DriverName1, &result.TargetDriverType, &result.TargetDriverCd, &result.TargetDriverName,
		&result.StartDatetime, &result.EndDatetime, &result.EventCd, &result.EventName,
		&result.StartMileage, &result.EndMileage, &result.SectionTime, &result.SectionDistance,
		&result.StartCityCd, &result.StartCityName, &result.EndCityCd, &result.EndCityName,
		&result.StartPlaceCd, &result.StartPlaceName, &result.EndPlaceCd, &result.EndPlaceName,
		&result.StartGpsValid, &result.StartGpsLat, &result.StartGpsLng,
		&result.EndGpsValid, &result.EndGpsLat, &result.EndGpsLng,
		&result.OverLimitMax,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetByUUID retrieves a kudgsir record by UUID
func (r *KudgsirRepository) GetByUUID(ctx context.Context, uuid string) (*Kudgsir, error) {
	query := `
		SELECT
			uuid, organization_id, hash, created, deleted, kudguri_uuid,
			unkou_no, read_date,
			office_cd, office_name, vehicle_cd, vehicle_name,
			driver_cd_1, driver_name_1, target_driver_type, target_driver_cd, target_driver_name,
			start_datetime, end_datetime, event_cd, event_name,
			start_mileage, end_mileage, section_time, section_distance,
			start_city_cd, start_city_name, end_city_cd, end_city_name,
			start_place_cd, start_place_name, end_place_cd, end_place_name,
			start_gps_valid, start_gps_lat, start_gps_lng,
			end_gps_valid, end_gps_lat, end_gps_lng,
			over_limit_max
		FROM kudgsir
		WHERE uuid = $1 AND deleted IS NULL
	`

	var k Kudgsir
	err := r.db.QueryRow(ctx, query, uuid).Scan(
		&k.UUID, &k.OrganizationID, &k.Hash, &k.Created, &k.Deleted, &k.KudguriUuid,
		&k.UnkouNo, &k.ReadDate,
		&k.OfficeCd, &k.OfficeName, &k.VehicleCd, &k.VehicleName,
		&k.DriverCd1, &k.DriverName1, &k.TargetDriverType, &k.TargetDriverCd, &k.TargetDriverName,
		&k.StartDatetime, &k.EndDatetime, &k.EventCd, &k.EventName,
		&k.StartMileage, &k.EndMileage, &k.SectionTime, &k.SectionDistance,
		&k.StartCityCd, &k.StartCityName, &k.EndCityCd, &k.EndCityName,
		&k.StartPlaceCd, &k.StartPlaceName, &k.EndPlaceCd, &k.EndPlaceName,
		&k.StartGpsValid, &k.StartGpsLat, &k.StartGpsLng,
		&k.EndGpsValid, &k.EndGpsLat, &k.EndGpsLng,
		&k.OverLimitMax,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrKudgsirNotFound
		}
		return nil, err
	}

	return &k, nil
}

// Update modifies an existing kudgsir record
func (r *KudgsirRepository) Update(ctx context.Context, k *Kudgsir) (*Kudgsir, error) {
	query := `
		UPDATE kudgsir
		SET
			organization_id = $2, hash = $3, kudguri_uuid = $4,
			unkou_no = $5, read_date = $6,
			office_cd = $7, office_name = $8, vehicle_cd = $9, vehicle_name = $10,
			driver_cd_1 = $11, driver_name_1 = $12, target_driver_type = $13, target_driver_cd = $14, target_driver_name = $15,
			start_datetime = $16, end_datetime = $17, event_cd = $18, event_name = $19,
			start_mileage = $20, end_mileage = $21, section_time = $22, section_distance = $23,
			start_city_cd = $24, start_city_name = $25, end_city_cd = $26, end_city_name = $27,
			start_place_cd = $28, start_place_name = $29, end_place_cd = $30, end_place_name = $31,
			start_gps_valid = $32, start_gps_lat = $33, start_gps_lng = $34,
			end_gps_valid = $35, end_gps_lat = $36, end_gps_lng = $37,
			over_limit_max = $38
		WHERE uuid = $1 AND deleted IS NULL
		RETURNING
			uuid, organization_id, hash, created, deleted, kudguri_uuid,
			unkou_no, read_date,
			office_cd, office_name, vehicle_cd, vehicle_name,
			driver_cd_1, driver_name_1, target_driver_type, target_driver_cd, target_driver_name,
			start_datetime, end_datetime, event_cd, event_name,
			start_mileage, end_mileage, section_time, section_distance,
			start_city_cd, start_city_name, end_city_cd, end_city_name,
			start_place_cd, start_place_name, end_place_cd, end_place_name,
			start_gps_valid, start_gps_lat, start_gps_lng,
			end_gps_valid, end_gps_lat, end_gps_lng,
			over_limit_max
	`

	var result Kudgsir
	err := r.db.QueryRow(ctx, query,
		k.UUID, k.OrganizationID, k.Hash, k.KudguriUuid,
		k.UnkouNo, k.ReadDate,
		k.OfficeCd, k.OfficeName, k.VehicleCd, k.VehicleName,
		k.DriverCd1, k.DriverName1, k.TargetDriverType, k.TargetDriverCd, k.TargetDriverName,
		k.StartDatetime, k.EndDatetime, k.EventCd, k.EventName,
		k.StartMileage, k.EndMileage, k.SectionTime, k.SectionDistance,
		k.StartCityCd, k.StartCityName, k.EndCityCd, k.EndCityName,
		k.StartPlaceCd, k.StartPlaceName, k.EndPlaceCd, k.EndPlaceName,
		k.StartGpsValid, k.StartGpsLat, k.StartGpsLng,
		k.EndGpsValid, k.EndGpsLat, k.EndGpsLng,
		k.OverLimitMax,
	).Scan(
		&result.UUID, &result.OrganizationID, &result.Hash, &result.Created, &result.Deleted, &result.KudguriUuid,
		&result.UnkouNo, &result.ReadDate,
		&result.OfficeCd, &result.OfficeName, &result.VehicleCd, &result.VehicleName,
		&result.DriverCd1, &result.DriverName1, &result.TargetDriverType, &result.TargetDriverCd, &result.TargetDriverName,
		&result.StartDatetime, &result.EndDatetime, &result.EventCd, &result.EventName,
		&result.StartMileage, &result.EndMileage, &result.SectionTime, &result.SectionDistance,
		&result.StartCityCd, &result.StartCityName, &result.EndCityCd, &result.EndCityName,
		&result.StartPlaceCd, &result.StartPlaceName, &result.EndPlaceCd, &result.EndPlaceName,
		&result.StartGpsValid, &result.StartGpsLat, &result.StartGpsLng,
		&result.EndGpsValid, &result.EndGpsLat, &result.EndGpsLng,
		&result.OverLimitMax,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrKudgsirNotFound
		}
		return nil, err
	}

	return &result, nil
}

// Delete soft-deletes a kudgsir record
func (r *KudgsirRepository) Delete(ctx context.Context, uuid, deletedTimestamp string) error {
	query := `
		UPDATE kudgsir
		SET deleted = $2
		WHERE uuid = $1 AND deleted IS NULL
	`

	result, err := r.db.Exec(ctx, query, uuid, deletedTimestamp)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrKudgsirNotFound
	}

	return nil
}

// ListByOrganization retrieves kudgsir records by organization with pagination
func (r *KudgsirRepository) ListByOrganization(ctx context.Context, organizationID string, limit int, offset int) ([]*Kudgsir, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	query := `
		SELECT
			uuid, organization_id, hash, created, deleted, kudguri_uuid,
			unkou_no, read_date,
			office_cd, office_name, vehicle_cd, vehicle_name,
			driver_cd_1, driver_name_1, target_driver_type, target_driver_cd, target_driver_name,
			start_datetime, end_datetime, event_cd, event_name,
			start_mileage, end_mileage, section_time, section_distance,
			start_city_cd, start_city_name, end_city_cd, end_city_name,
			start_place_cd, start_place_name, end_place_cd, end_place_name,
			start_gps_valid, start_gps_lat, start_gps_lng,
			end_gps_valid, end_gps_lat, end_gps_lng,
			over_limit_max
		FROM kudgsir
		WHERE organization_id = $1 AND deleted IS NULL
		ORDER BY created DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, organizationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*Kudgsir
	for rows.Next() {
		var k Kudgsir
		err := rows.Scan(
			&k.UUID, &k.OrganizationID, &k.Hash, &k.Created, &k.Deleted, &k.KudguriUuid,
			&k.UnkouNo, &k.ReadDate,
			&k.OfficeCd, &k.OfficeName, &k.VehicleCd, &k.VehicleName,
			&k.DriverCd1, &k.DriverName1, &k.TargetDriverType, &k.TargetDriverCd, &k.TargetDriverName,
			&k.StartDatetime, &k.EndDatetime, &k.EventCd, &k.EventName,
			&k.StartMileage, &k.EndMileage, &k.SectionTime, &k.SectionDistance,
			&k.StartCityCd, &k.StartCityName, &k.EndCityCd, &k.EndCityName,
			&k.StartPlaceCd, &k.StartPlaceName, &k.EndPlaceCd, &k.EndPlaceName,
			&k.StartGpsValid, &k.StartGpsLat, &k.StartGpsLng,
			&k.EndGpsValid, &k.EndGpsLat, &k.EndGpsLng,
			&k.OverLimitMax,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, &k)
	}

	return results, rows.Err()
}

// List retrieves all kudgsir records with pagination
func (r *KudgsirRepository) List(ctx context.Context, limit int, offset int) ([]*Kudgsir, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	query := `
		SELECT
			uuid, organization_id, hash, created, deleted, kudguri_uuid,
			unkou_no, read_date,
			office_cd, office_name, vehicle_cd, vehicle_name,
			driver_cd_1, driver_name_1, target_driver_type, target_driver_cd, target_driver_name,
			start_datetime, end_datetime, event_cd, event_name,
			start_mileage, end_mileage, section_time, section_distance,
			start_city_cd, start_city_name, end_city_cd, end_city_name,
			start_place_cd, start_place_name, end_place_cd, end_place_name,
			start_gps_valid, start_gps_lat, start_gps_lng,
			end_gps_valid, end_gps_lat, end_gps_lng,
			over_limit_max
		FROM kudgsir
		WHERE deleted IS NULL
		ORDER BY created DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*Kudgsir
	for rows.Next() {
		var k Kudgsir
		err := rows.Scan(
			&k.UUID, &k.OrganizationID, &k.Hash, &k.Created, &k.Deleted, &k.KudguriUuid,
			&k.UnkouNo, &k.ReadDate,
			&k.OfficeCd, &k.OfficeName, &k.VehicleCd, &k.VehicleName,
			&k.DriverCd1, &k.DriverName1, &k.TargetDriverType, &k.TargetDriverCd, &k.TargetDriverName,
			&k.StartDatetime, &k.EndDatetime, &k.EventCd, &k.EventName,
			&k.StartMileage, &k.EndMileage, &k.SectionTime, &k.SectionDistance,
			&k.StartCityCd, &k.StartCityName, &k.EndCityCd, &k.EndCityName,
			&k.StartPlaceCd, &k.StartPlaceName, &k.EndPlaceCd, &k.EndPlaceName,
			&k.StartGpsValid, &k.StartGpsLat, &k.StartGpsLng,
			&k.EndGpsValid, &k.EndGpsLat, &k.EndGpsLng,
			&k.OverLimitMax,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, &k)
	}

	return results, rows.Err()
}
