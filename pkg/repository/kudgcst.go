package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrKudgcstNotFound = errors.New("kudgcst not found")
)

// Kudgcst represents the database model for kudgcst table (フェリー運賃データ)
type Kudgcst struct {
	UUID           string
	OrganizationID string
	Hash           string
	Created        string
	Deleted        *string
	KudguriUuid    *string

	// 運行基本情報
	UnkouNo   *string
	UnkouDate *string
	ReadDate  *string

	// 営業所・車両情報
	OfficeCd    *string
	OfficeName  *string
	VehicleCd   *string
	VehicleName *string

	// ドライバー情報
	DriverCd1        *string
	DriverName1      *string
	TargetDriverType string

	// フェリー運行情報
	StartDatetime *string
	EndDatetime   *string

	// フェリー会社情報
	FerryCompanyCd   *string
	FerryCompanyName *string

	// 乗船・降船情報
	BoardingPlaceCd   *string
	BoardingPlaceName *string
	TripNumber        *string
	DropoffPlaceCd    *string
	DropoffPlaceName  *string

	// 運賃情報
	SettlementType     *string
	SettlementTypeName *string
	StandardFare       *string
	ContractFare       *string

	// 車両区分
	FerryVehicleType     *string
	FerryVehicleTypeName *string

	// 距離
	AssumedDistance *string
}

// KudgcstRepository handles database operations for kudgcst
type KudgcstRepository struct {
	db DB
}

// NewKudgcstRepository creates a new repository
func NewKudgcstRepository(pool *pgxpool.Pool) *KudgcstRepository {
	return &KudgcstRepository{db: pool}
}

// NewKudgcstRepositoryWithDB creates a repository with custom DB interface (for testing)
func NewKudgcstRepositoryWithDB(db DB) *KudgcstRepository {
	return &KudgcstRepository{db: db}
}

// Create inserts a new kudgcst record
func (r *KudgcstRepository) Create(ctx context.Context, k *Kudgcst) (*Kudgcst, error) {
	id := uuid.New().String()
	k.UUID = id

	query := `
		INSERT INTO kudgcst (
			"UUID", "OrganizationID", "Hash", "Created", "Deleted", "KudguriUuid",
			"UnkouNo", "UnkouDate", "ReadDate",
			"OfficeCd", "OfficeName", "VehicleCd", "VehicleName",
			"DriverCd1", "DriverName1", "TargetDriverType",
			"StartDatetime", "EndDatetime",
			"FerryCompanyCd", "FerryCompanyName",
			"BoardingPlaceCd", "BoardingPlaceName", "TripNumber",
			"DropoffPlaceCd", "DropoffPlaceName",
			"SettlementType", "SettlementTypeName", "StandardFare", "ContractFare",
			"FerryVehicleType", "FerryVehicleTypeName",
			"AssumedDistance"
		)
		VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9,
			$10, $11, $12, $13,
			$14, $15, $16,
			$17, $18,
			$19, $20,
			$21, $22, $23,
			$24, $25,
			$26, $27, $28, $29,
			$30, $31,
			$32
		)
		RETURNING
			"UUID", "OrganizationID", "Hash", "Created", "Deleted", "KudguriUuid",
			"UnkouNo", "UnkouDate", "ReadDate",
			"OfficeCd", "OfficeName", "VehicleCd", "VehicleName",
			"DriverCd1", "DriverName1", "TargetDriverType",
			"StartDatetime", "EndDatetime",
			"FerryCompanyCd", "FerryCompanyName",
			"BoardingPlaceCd", "BoardingPlaceName", "TripNumber",
			"DropoffPlaceCd", "DropoffPlaceName",
			"SettlementType", "SettlementTypeName", "StandardFare", "ContractFare",
			"FerryVehicleType", "FerryVehicleTypeName",
			"AssumedDistance"
	`

	var result Kudgcst
	err := r.db.QueryRow(ctx, query,
		k.UUID, k.OrganizationID, k.Hash, k.Created, k.Deleted, k.KudguriUuid,
		k.UnkouNo, k.UnkouDate, k.ReadDate,
		k.OfficeCd, k.OfficeName, k.VehicleCd, k.VehicleName,
		k.DriverCd1, k.DriverName1, k.TargetDriverType,
		k.StartDatetime, k.EndDatetime,
		k.FerryCompanyCd, k.FerryCompanyName,
		k.BoardingPlaceCd, k.BoardingPlaceName, k.TripNumber,
		k.DropoffPlaceCd, k.DropoffPlaceName,
		k.SettlementType, k.SettlementTypeName, k.StandardFare, k.ContractFare,
		k.FerryVehicleType, k.FerryVehicleTypeName,
		k.AssumedDistance,
	).Scan(
		&result.UUID, &result.OrganizationID, &result.Hash, &result.Created, &result.Deleted, &result.KudguriUuid,
		&result.UnkouNo, &result.UnkouDate, &result.ReadDate,
		&result.OfficeCd, &result.OfficeName, &result.VehicleCd, &result.VehicleName,
		&result.DriverCd1, &result.DriverName1, &result.TargetDriverType,
		&result.StartDatetime, &result.EndDatetime,
		&result.FerryCompanyCd, &result.FerryCompanyName,
		&result.BoardingPlaceCd, &result.BoardingPlaceName, &result.TripNumber,
		&result.DropoffPlaceCd, &result.DropoffPlaceName,
		&result.SettlementType, &result.SettlementTypeName, &result.StandardFare, &result.ContractFare,
		&result.FerryVehicleType, &result.FerryVehicleTypeName,
		&result.AssumedDistance,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetByUUID retrieves a kudgcst record by UUID
func (r *KudgcstRepository) GetByUUID(ctx context.Context, uuid string) (*Kudgcst, error) {
	query := `
		SELECT
			"UUID", "OrganizationID", "Hash", "Created", "Deleted", "KudguriUuid",
			"UnkouNo", "UnkouDate", "ReadDate",
			"OfficeCd", "OfficeName", "VehicleCd", "VehicleName",
			"DriverCd1", "DriverName1", "TargetDriverType",
			"StartDatetime", "EndDatetime",
			"FerryCompanyCd", "FerryCompanyName",
			"BoardingPlaceCd", "BoardingPlaceName", "TripNumber",
			"DropoffPlaceCd", "DropoffPlaceName",
			"SettlementType", "SettlementTypeName", "StandardFare", "ContractFare",
			"FerryVehicleType", "FerryVehicleTypeName",
			"AssumedDistance"
		FROM kudgcst
		WHERE "UUID" = $1 AND "Deleted" IS NULL
	`

	var k Kudgcst
	err := r.db.QueryRow(ctx, query, uuid).Scan(
		&k.UUID, &k.OrganizationID, &k.Hash, &k.Created, &k.Deleted, &k.KudguriUuid,
		&k.UnkouNo, &k.UnkouDate, &k.ReadDate,
		&k.OfficeCd, &k.OfficeName, &k.VehicleCd, &k.VehicleName,
		&k.DriverCd1, &k.DriverName1, &k.TargetDriverType,
		&k.StartDatetime, &k.EndDatetime,
		&k.FerryCompanyCd, &k.FerryCompanyName,
		&k.BoardingPlaceCd, &k.BoardingPlaceName, &k.TripNumber,
		&k.DropoffPlaceCd, &k.DropoffPlaceName,
		&k.SettlementType, &k.SettlementTypeName, &k.StandardFare, &k.ContractFare,
		&k.FerryVehicleType, &k.FerryVehicleTypeName,
		&k.AssumedDistance,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrKudgcstNotFound
		}
		return nil, err
	}

	return &k, nil
}

// Update modifies an existing kudgcst record
func (r *KudgcstRepository) Update(ctx context.Context, k *Kudgcst) (*Kudgcst, error) {
	query := `
		UPDATE kudgcst
		SET
			"OrganizationID" = $2, "Hash" = $3, "KudguriUuid" = $4,
			"UnkouNo" = $5, "UnkouDate" = $6, "ReadDate" = $7,
			"OfficeCd" = $8, "OfficeName" = $9, "VehicleCd" = $10, "VehicleName" = $11,
			"DriverCd1" = $12, "DriverName1" = $13, "TargetDriverType" = $14,
			"StartDatetime" = $15, "EndDatetime" = $16,
			"FerryCompanyCd" = $17, "FerryCompanyName" = $18,
			"BoardingPlaceCd" = $19, "BoardingPlaceName" = $20, "TripNumber" = $21,
			"DropoffPlaceCd" = $22, "DropoffPlaceName" = $23,
			"SettlementType" = $24, "SettlementTypeName" = $25, "StandardFare" = $26, "ContractFare" = $27,
			"FerryVehicleType" = $28, "FerryVehicleTypeName" = $29,
			"AssumedDistance" = $30
		WHERE "UUID" = $1 AND "Deleted" IS NULL
		RETURNING
			"UUID", "OrganizationID", "Hash", "Created", "Deleted", "KudguriUuid",
			"UnkouNo", "UnkouDate", "ReadDate",
			"OfficeCd", "OfficeName", "VehicleCd", "VehicleName",
			"DriverCd1", "DriverName1", "TargetDriverType",
			"StartDatetime", "EndDatetime",
			"FerryCompanyCd", "FerryCompanyName",
			"BoardingPlaceCd", "BoardingPlaceName", "TripNumber",
			"DropoffPlaceCd", "DropoffPlaceName",
			"SettlementType", "SettlementTypeName", "StandardFare", "ContractFare",
			"FerryVehicleType", "FerryVehicleTypeName",
			"AssumedDistance"
	`

	var result Kudgcst
	err := r.db.QueryRow(ctx, query,
		k.UUID, k.OrganizationID, k.Hash, k.KudguriUuid,
		k.UnkouNo, k.UnkouDate, k.ReadDate,
		k.OfficeCd, k.OfficeName, k.VehicleCd, k.VehicleName,
		k.DriverCd1, k.DriverName1, k.TargetDriverType,
		k.StartDatetime, k.EndDatetime,
		k.FerryCompanyCd, k.FerryCompanyName,
		k.BoardingPlaceCd, k.BoardingPlaceName, k.TripNumber,
		k.DropoffPlaceCd, k.DropoffPlaceName,
		k.SettlementType, k.SettlementTypeName, k.StandardFare, k.ContractFare,
		k.FerryVehicleType, k.FerryVehicleTypeName,
		k.AssumedDistance,
	).Scan(
		&result.UUID, &result.OrganizationID, &result.Hash, &result.Created, &result.Deleted, &result.KudguriUuid,
		&result.UnkouNo, &result.UnkouDate, &result.ReadDate,
		&result.OfficeCd, &result.OfficeName, &result.VehicleCd, &result.VehicleName,
		&result.DriverCd1, &result.DriverName1, &result.TargetDriverType,
		&result.StartDatetime, &result.EndDatetime,
		&result.FerryCompanyCd, &result.FerryCompanyName,
		&result.BoardingPlaceCd, &result.BoardingPlaceName, &result.TripNumber,
		&result.DropoffPlaceCd, &result.DropoffPlaceName,
		&result.SettlementType, &result.SettlementTypeName, &result.StandardFare, &result.ContractFare,
		&result.FerryVehicleType, &result.FerryVehicleTypeName,
		&result.AssumedDistance,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrKudgcstNotFound
		}
		return nil, err
	}

	return &result, nil
}

// Delete soft-deletes a kudgcst record
func (r *KudgcstRepository) Delete(ctx context.Context, uuid, deletedTimestamp string) error {
	query := `
		UPDATE kudgcst
		SET "Deleted" = $2
		WHERE "UUID" = $1 AND "Deleted" IS NULL
	`

	result, err := r.db.Exec(ctx, query, uuid, deletedTimestamp)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrKudgcstNotFound
	}

	return nil
}

// ListByOrganization retrieves kudgcst records by organization with pagination
func (r *KudgcstRepository) ListByOrganization(ctx context.Context, organizationID string, limit int, offset int) ([]*Kudgcst, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	query := `
		SELECT
			"UUID", "OrganizationID", "Hash", "Created", "Deleted", "KudguriUuid",
			"UnkouNo", "UnkouDate", "ReadDate",
			"OfficeCd", "OfficeName", "VehicleCd", "VehicleName",
			"DriverCd1", "DriverName1", "TargetDriverType",
			"StartDatetime", "EndDatetime",
			"FerryCompanyCd", "FerryCompanyName",
			"BoardingPlaceCd", "BoardingPlaceName", "TripNumber",
			"DropoffPlaceCd", "DropoffPlaceName",
			"SettlementType", "SettlementTypeName", "StandardFare", "ContractFare",
			"FerryVehicleType", "FerryVehicleTypeName",
			"AssumedDistance"
		FROM kudgcst
		WHERE "OrganizationID" = $1 AND "Deleted" IS NULL
		ORDER BY "Created" DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, organizationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*Kudgcst
	for rows.Next() {
		var k Kudgcst
		err := rows.Scan(
			&k.UUID, &k.OrganizationID, &k.Hash, &k.Created, &k.Deleted, &k.KudguriUuid,
			&k.UnkouNo, &k.UnkouDate, &k.ReadDate,
			&k.OfficeCd, &k.OfficeName, &k.VehicleCd, &k.VehicleName,
			&k.DriverCd1, &k.DriverName1, &k.TargetDriverType,
			&k.StartDatetime, &k.EndDatetime,
			&k.FerryCompanyCd, &k.FerryCompanyName,
			&k.BoardingPlaceCd, &k.BoardingPlaceName, &k.TripNumber,
			&k.DropoffPlaceCd, &k.DropoffPlaceName,
			&k.SettlementType, &k.SettlementTypeName, &k.StandardFare, &k.ContractFare,
			&k.FerryVehicleType, &k.FerryVehicleTypeName,
			&k.AssumedDistance,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, &k)
	}

	return results, rows.Err()
}

// List retrieves all kudgcst records with pagination
func (r *KudgcstRepository) List(ctx context.Context, limit int, offset int) ([]*Kudgcst, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	query := `
		SELECT
			"UUID", "OrganizationID", "Hash", "Created", "Deleted", "KudguriUuid",
			"UnkouNo", "UnkouDate", "ReadDate",
			"OfficeCd", "OfficeName", "VehicleCd", "VehicleName",
			"DriverCd1", "DriverName1", "TargetDriverType",
			"StartDatetime", "EndDatetime",
			"FerryCompanyCd", "FerryCompanyName",
			"BoardingPlaceCd", "BoardingPlaceName", "TripNumber",
			"DropoffPlaceCd", "DropoffPlaceName",
			"SettlementType", "SettlementTypeName", "StandardFare", "ContractFare",
			"FerryVehicleType", "FerryVehicleTypeName",
			"AssumedDistance"
		FROM kudgcst
		WHERE "Deleted" IS NULL
		ORDER BY "Created" DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*Kudgcst
	for rows.Next() {
		var k Kudgcst
		err := rows.Scan(
			&k.UUID, &k.OrganizationID, &k.Hash, &k.Created, &k.Deleted, &k.KudguriUuid,
			&k.UnkouNo, &k.UnkouDate, &k.ReadDate,
			&k.OfficeCd, &k.OfficeName, &k.VehicleCd, &k.VehicleName,
			&k.DriverCd1, &k.DriverName1, &k.TargetDriverType,
			&k.StartDatetime, &k.EndDatetime,
			&k.FerryCompanyCd, &k.FerryCompanyName,
			&k.BoardingPlaceCd, &k.BoardingPlaceName, &k.TripNumber,
			&k.DropoffPlaceCd, &k.DropoffPlaceName,
			&k.SettlementType, &k.SettlementTypeName, &k.StandardFare, &k.ContractFare,
			&k.FerryVehicleType, &k.FerryVehicleTypeName,
			&k.AssumedDistance,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, &k)
	}

	return results, rows.Err()
}
