package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrKudgivtNotFound = errors.New("kudgivt not found")
)

// Kudgivt represents the database model for kudgivt table
type Kudgivt struct {
	UUID           string
	OrganizationID string
	Hash           string
	Created        string
	Deleted        *string
	KudguriUuid    *string

	// 運行情報
	UnkouNo   *string
	ReadDate  *string
	UnkouDate *string

	// 事業所・車両情報
	OfficeCd    *string
	OfficeName  *string
	VehicleCd   *string
	VehicleName *string

	// 乗務員情報
	DriverCd1        *string
	DriverName1      *string
	TargetDriverType string
	TargetDriverCd   *string
	TargetDriverName *string

	// 勤務時間
	ClockInDatetime  *string
	ClockOutDatetime *string

	// 出発・帰着時刻
	DepartureDatetime *string
	ReturnDatetime    *string

	// 走行距離
	DepartureMeter *string
	ReturnMeter    *string
	TotalMileage   *string

	// 行先情報
	DestinationCityName  *string
	DestinationPlaceName *string

	// 実車距離
	ActualMileage *string

	// 運転時間
	LocalDriveTime   *string
	ExpressDriveTime *string
	BypassDriveTime  *string
	ActualDriveTime  *string
	EmptyDriveTime   *string

	// 作業時間 (Work1Time ~ Work10Time)
	Work1Time  *string
	Work2Time  *string
	Work3Time  *string
	Work4Time  *string
	Work5Time  *string
	Work6Time  *string
	Work7Time  *string
	Work8Time  *string
	Work9Time  *string
	Work10Time *string

	// 状態距離 (State1Distance ~ State5Distance)
	State1Distance *string
	State2Distance *string
	State3Distance *string
	State4Distance *string
	State5Distance *string

	// 状態時間 (State1Time ~ State5Time)
	State1Time *string
	State2Time *string
	State3Time *string
	State4Time *string
	State5Time *string

	// 自社燃料等
	OwnMainFuel     *string
	OwnMainAdditive *string
	OwnConsumable   *string

	// 他社燃料等
	OtherMainFuel     *string
	OtherMainAdditive *string
	OtherConsumable   *string

	// 速度超過（一般道）
	LocalSpeedOverMax   *string
	LocalSpeedOverTime  *string
	LocalSpeedOverCount *string

	// 速度超過（高速道）
	ExpressSpeedOverMax   *string
	ExpressSpeedOverTime  *string
	ExpressSpeedOverCount *string

	// 速度超過（専用道）
	DedicatedSpeedOverMax   *string
	DedicatedSpeedOverTime  *string
	DedicatedSpeedOverCount *string

	// アイドリング
	IdlingTime      *string
	IdlingTimeCount *string

	// 回転数超過
	RotationOverMax   *string
	RotationOverCount *string
	RotationOverTime  *string

	// 急加速
	RapidAccelCount1   *string
	RapidAccelCount2   *string
	RapidAccelCount3   *string
	RapidAccelCount4   *string
	RapidAccelCount5   *string
	RapidAccelMax      *string
	RapidAccelMaxSpeed *string

	// 急減速
	RapidDecelCount1   *string
	RapidDecelCount2   *string
	RapidDecelCount3   *string
	RapidDecelCount4   *string
	RapidDecelCount5   *string
	RapidDecelMax      *string
	RapidDecelMaxSpeed *string

	// 急ハンドル
	RapidCurveCount1   *string
	RapidCurveCount2   *string
	RapidCurveCount3   *string
	RapidCurveCount4   *string
	RapidCurveCount5   *string
	RapidCurveMax      *string
	RapidCurveMaxSpeed *string

	// 連続運転
	ContinuousDriveOverCount *string
	ContinuousDriveMaxTime   *string
	ContinuousDriveTotalTime *string

	// 波状運転
	WaveDriveCount        *string
	WaveDriveMaxTime      *string
	WaveDriveMaxSpeedDiff *string

	// スコア（速度）
	LocalSpeedScore     *string
	ExpressSpeedScore   *string
	DedicatedSpeedScore *string

	// スコア（距離）
	LocalDistanceScore     *string
	ExpressDistanceScore   *string
	DedicatedDistanceScore *string

	// スコア（急操作）
	RapidAccelScore *string
	RapidDecelScore *string
	RapidCurveScore *string

	// スコア（回転数・実車）
	ActualLowSpeedRotationScore  *string
	ActualHighSpeedRotationScore *string

	// スコア（回転数・空車）
	EmptyLowSpeedRotationScore  *string
	EmptyHighSpeedRotationScore *string

	// スコア（その他）
	IdlingScore          *string
	ContinuousDriveScore *string
	WaveDriveScore       *string

	// 総合スコア
	SafetyScore  *string
	EconomyScore *string
	TotalScore   *string
}

// KudgivtRepository handles database operations for kudgivt
type KudgivtRepository struct {
	db DB
}

// NewKudgivtRepository creates a new repository
func NewKudgivtRepository(pool *pgxpool.Pool) *KudgivtRepository {
	return &KudgivtRepository{db: pool}
}

// NewKudgivtRepositoryWithDB creates a repository with custom DB interface (for testing)
func NewKudgivtRepositoryWithDB(db DB) *KudgivtRepository {
	return &KudgivtRepository{db: db}
}

// Create inserts a new kudgivt record
func (r *KudgivtRepository) Create(ctx context.Context, k *Kudgivt) (*Kudgivt, error) {
	query := `
		INSERT INTO kudgivt (
			"UUID", "OrganizationID", "Hash", "Created", "Deleted", "KudguriUuid",
			"UnkouNo", "ReadDate", "UnkouDate",
			"OfficeCd", "OfficeName", "VehicleCd", "VehicleName",
			"DriverCd1", "DriverName1", "TargetDriverType", "TargetDriverCd", "TargetDriverName",
			"ClockInDatetime", "ClockOutDatetime",
			"DepartureDatetime", "ReturnDatetime",
			"DepartureMeter", "ReturnMeter", "TotalMileage",
			"DestinationCityName", "DestinationPlaceName",
			"ActualMileage",
			"LocalDriveTime", "ExpressDriveTime", "BypassDriveTime", "ActualDriveTime", "EmptyDriveTime",
			"Work1Time", "Work2Time", "Work3Time", "Work4Time", "Work5Time",
			"Work6Time", "Work7Time", "Work8Time", "Work9Time", "Work10Time",
			"State1Distance", "State2Distance", "State3Distance", "State4Distance", "State5Distance",
			"State1Time", "State2Time", "State3Time", "State4Time", "State5Time",
			"OwnMainFuel", "OwnMainAdditive", "OwnConsumable",
			"OtherMainFuel", "OtherMainAdditive", "OtherConsumable",
			"LocalSpeedOverMax", "LocalSpeedOverTime", "LocalSpeedOverCount",
			"ExpressSpeedOverMax", "ExpressSpeedOverTime", "ExpressSpeedOverCount",
			"DedicatedSpeedOverMax", "DedicatedSpeedOverTime", "DedicatedSpeedOverCount",
			"IdlingTime", "IdlingTimeCount",
			"RotationOverMax", "RotationOverCount", "RotationOverTime",
			"RapidAccelCount1", "RapidAccelCount2", "RapidAccelCount3", "RapidAccelCount4", "RapidAccelCount5",
			"RapidAccelMax", "RapidAccelMaxSpeed",
			"RapidDecelCount1", "RapidDecelCount2", "RapidDecelCount3", "RapidDecelCount4", "RapidDecelCount5",
			"RapidDecelMax", "RapidDecelMaxSpeed",
			"RapidCurveCount1", "RapidCurveCount2", "RapidCurveCount3", "RapidCurveCount4", "RapidCurveCount5",
			"RapidCurveMax", "RapidCurveMaxSpeed",
			"ContinuousDriveOverCount", "ContinuousDriveMaxTime", "ContinuousDriveTotalTime",
			"WaveDriveCount", "WaveDriveMaxTime", "WaveDriveMaxSpeedDiff",
			"LocalSpeedScore", "ExpressSpeedScore", "DedicatedSpeedScore",
			"LocalDistanceScore", "ExpressDistanceScore", "DedicatedDistanceScore",
			"RapidAccelScore", "RapidDecelScore", "RapidCurveScore",
			"ActualLowSpeedRotationScore", "ActualHighSpeedRotationScore",
			"EmptyLowSpeedRotationScore", "EmptyHighSpeedRotationScore",
			"IdlingScore", "ContinuousDriveScore", "WaveDriveScore",
			"SafetyScore", "EconomyScore", "TotalScore"
		)
		VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9,
			$10, $11, $12, $13,
			$14, $15, $16, $17, $18,
			$19, $20,
			$21, $22,
			$23, $24, $25,
			$26, $27,
			$28,
			$29, $30, $31, $32, $33,
			$34, $35, $36, $37, $38,
			$39, $40, $41, $42, $43,
			$44, $45, $46, $47, $48,
			$49, $50, $51, $52, $53,
			$54, $55, $56,
			$57, $58, $59,
			$60, $61, $62,
			$63, $64, $65,
			$66, $67, $68,
			$69, $70,
			$71, $72, $73,
			$74, $75, $76, $77, $78,
			$79, $80,
			$81, $82, $83, $84, $85,
			$86, $87,
			$88, $89, $90, $91, $92,
			$93, $94,
			$95, $96, $97,
			$98, $99, $100,
			$101, $102, $103,
			$104, $105, $106,
			$107, $108, $109,
			$110, $111,
			$112, $113,
			$114, $115, $116,
			$117, $118, $119
		)
		RETURNING
			"UUID", "OrganizationID", "Hash", "Created", "Deleted", "KudguriUuid",
			"UnkouNo", "ReadDate", "UnkouDate",
			"OfficeCd", "OfficeName", "VehicleCd", "VehicleName",
			"DriverCd1", "DriverName1", "TargetDriverType", "TargetDriverCd", "TargetDriverName",
			"ClockInDatetime", "ClockOutDatetime",
			"DepartureDatetime", "ReturnDatetime",
			"DepartureMeter", "ReturnMeter", "TotalMileage",
			"DestinationCityName", "DestinationPlaceName",
			"ActualMileage",
			"LocalDriveTime", "ExpressDriveTime", "BypassDriveTime", "ActualDriveTime", "EmptyDriveTime",
			"Work1Time", "Work2Time", "Work3Time", "Work4Time", "Work5Time",
			"Work6Time", "Work7Time", "Work8Time", "Work9Time", "Work10Time",
			"State1Distance", "State2Distance", "State3Distance", "State4Distance", "State5Distance",
			"State1Time", "State2Time", "State3Time", "State4Time", "State5Time",
			"OwnMainFuel", "OwnMainAdditive", "OwnConsumable",
			"OtherMainFuel", "OtherMainAdditive", "OtherConsumable",
			"LocalSpeedOverMax", "LocalSpeedOverTime", "LocalSpeedOverCount",
			"ExpressSpeedOverMax", "ExpressSpeedOverTime", "ExpressSpeedOverCount",
			"DedicatedSpeedOverMax", "DedicatedSpeedOverTime", "DedicatedSpeedOverCount",
			"IdlingTime", "IdlingTimeCount",
			"RotationOverMax", "RotationOverCount", "RotationOverTime",
			"RapidAccelCount1", "RapidAccelCount2", "RapidAccelCount3", "RapidAccelCount4", "RapidAccelCount5",
			"RapidAccelMax", "RapidAccelMaxSpeed",
			"RapidDecelCount1", "RapidDecelCount2", "RapidDecelCount3", "RapidDecelCount4", "RapidDecelCount5",
			"RapidDecelMax", "RapidDecelMaxSpeed",
			"RapidCurveCount1", "RapidCurveCount2", "RapidCurveCount3", "RapidCurveCount4", "RapidCurveCount5",
			"RapidCurveMax", "RapidCurveMaxSpeed",
			"ContinuousDriveOverCount", "ContinuousDriveMaxTime", "ContinuousDriveTotalTime",
			"WaveDriveCount", "WaveDriveMaxTime", "WaveDriveMaxSpeedDiff",
			"LocalSpeedScore", "ExpressSpeedScore", "DedicatedSpeedScore",
			"LocalDistanceScore", "ExpressDistanceScore", "DedicatedDistanceScore",
			"RapidAccelScore", "RapidDecelScore", "RapidCurveScore",
			"ActualLowSpeedRotationScore", "ActualHighSpeedRotationScore",
			"EmptyLowSpeedRotationScore", "EmptyHighSpeedRotationScore",
			"IdlingScore", "ContinuousDriveScore", "WaveDriveScore",
			"SafetyScore", "EconomyScore", "TotalScore"
	`

	var result Kudgivt
	err := r.db.QueryRow(ctx, query,
		k.UUID, k.OrganizationID, k.Hash, k.Created, k.Deleted, k.KudguriUuid,
		k.UnkouNo, k.ReadDate, k.UnkouDate,
		k.OfficeCd, k.OfficeName, k.VehicleCd, k.VehicleName,
		k.DriverCd1, k.DriverName1, k.TargetDriverType, k.TargetDriverCd, k.TargetDriverName,
		k.ClockInDatetime, k.ClockOutDatetime,
		k.DepartureDatetime, k.ReturnDatetime,
		k.DepartureMeter, k.ReturnMeter, k.TotalMileage,
		k.DestinationCityName, k.DestinationPlaceName,
		k.ActualMileage,
		k.LocalDriveTime, k.ExpressDriveTime, k.BypassDriveTime, k.ActualDriveTime, k.EmptyDriveTime,
		k.Work1Time, k.Work2Time, k.Work3Time, k.Work4Time, k.Work5Time,
		k.Work6Time, k.Work7Time, k.Work8Time, k.Work9Time, k.Work10Time,
		k.State1Distance, k.State2Distance, k.State3Distance, k.State4Distance, k.State5Distance,
		k.State1Time, k.State2Time, k.State3Time, k.State4Time, k.State5Time,
		k.OwnMainFuel, k.OwnMainAdditive, k.OwnConsumable,
		k.OtherMainFuel, k.OtherMainAdditive, k.OtherConsumable,
		k.LocalSpeedOverMax, k.LocalSpeedOverTime, k.LocalSpeedOverCount,
		k.ExpressSpeedOverMax, k.ExpressSpeedOverTime, k.ExpressSpeedOverCount,
		k.DedicatedSpeedOverMax, k.DedicatedSpeedOverTime, k.DedicatedSpeedOverCount,
		k.IdlingTime, k.IdlingTimeCount,
		k.RotationOverMax, k.RotationOverCount, k.RotationOverTime,
		k.RapidAccelCount1, k.RapidAccelCount2, k.RapidAccelCount3, k.RapidAccelCount4, k.RapidAccelCount5,
		k.RapidAccelMax, k.RapidAccelMaxSpeed,
		k.RapidDecelCount1, k.RapidDecelCount2, k.RapidDecelCount3, k.RapidDecelCount4, k.RapidDecelCount5,
		k.RapidDecelMax, k.RapidDecelMaxSpeed,
		k.RapidCurveCount1, k.RapidCurveCount2, k.RapidCurveCount3, k.RapidCurveCount4, k.RapidCurveCount5,
		k.RapidCurveMax, k.RapidCurveMaxSpeed,
		k.ContinuousDriveOverCount, k.ContinuousDriveMaxTime, k.ContinuousDriveTotalTime,
		k.WaveDriveCount, k.WaveDriveMaxTime, k.WaveDriveMaxSpeedDiff,
		k.LocalSpeedScore, k.ExpressSpeedScore, k.DedicatedSpeedScore,
		k.LocalDistanceScore, k.ExpressDistanceScore, k.DedicatedDistanceScore,
		k.RapidAccelScore, k.RapidDecelScore, k.RapidCurveScore,
		k.ActualLowSpeedRotationScore, k.ActualHighSpeedRotationScore,
		k.EmptyLowSpeedRotationScore, k.EmptyHighSpeedRotationScore,
		k.IdlingScore, k.ContinuousDriveScore, k.WaveDriveScore,
		k.SafetyScore, k.EconomyScore, k.TotalScore,
	).Scan(
		&result.UUID, &result.OrganizationID, &result.Hash, &result.Created, &result.Deleted, &result.KudguriUuid,
		&result.UnkouNo, &result.ReadDate, &result.UnkouDate,
		&result.OfficeCd, &result.OfficeName, &result.VehicleCd, &result.VehicleName,
		&result.DriverCd1, &result.DriverName1, &result.TargetDriverType, &result.TargetDriverCd, &result.TargetDriverName,
		&result.ClockInDatetime, &result.ClockOutDatetime,
		&result.DepartureDatetime, &result.ReturnDatetime,
		&result.DepartureMeter, &result.ReturnMeter, &result.TotalMileage,
		&result.DestinationCityName, &result.DestinationPlaceName,
		&result.ActualMileage,
		&result.LocalDriveTime, &result.ExpressDriveTime, &result.BypassDriveTime, &result.ActualDriveTime, &result.EmptyDriveTime,
		&result.Work1Time, &result.Work2Time, &result.Work3Time, &result.Work4Time, &result.Work5Time,
		&result.Work6Time, &result.Work7Time, &result.Work8Time, &result.Work9Time, &result.Work10Time,
		&result.State1Distance, &result.State2Distance, &result.State3Distance, &result.State4Distance, &result.State5Distance,
		&result.State1Time, &result.State2Time, &result.State3Time, &result.State4Time, &result.State5Time,
		&result.OwnMainFuel, &result.OwnMainAdditive, &result.OwnConsumable,
		&result.OtherMainFuel, &result.OtherMainAdditive, &result.OtherConsumable,
		&result.LocalSpeedOverMax, &result.LocalSpeedOverTime, &result.LocalSpeedOverCount,
		&result.ExpressSpeedOverMax, &result.ExpressSpeedOverTime, &result.ExpressSpeedOverCount,
		&result.DedicatedSpeedOverMax, &result.DedicatedSpeedOverTime, &result.DedicatedSpeedOverCount,
		&result.IdlingTime, &result.IdlingTimeCount,
		&result.RotationOverMax, &result.RotationOverCount, &result.RotationOverTime,
		&result.RapidAccelCount1, &result.RapidAccelCount2, &result.RapidAccelCount3, &result.RapidAccelCount4, &result.RapidAccelCount5,
		&result.RapidAccelMax, &result.RapidAccelMaxSpeed,
		&result.RapidDecelCount1, &result.RapidDecelCount2, &result.RapidDecelCount3, &result.RapidDecelCount4, &result.RapidDecelCount5,
		&result.RapidDecelMax, &result.RapidDecelMaxSpeed,
		&result.RapidCurveCount1, &result.RapidCurveCount2, &result.RapidCurveCount3, &result.RapidCurveCount4, &result.RapidCurveCount5,
		&result.RapidCurveMax, &result.RapidCurveMaxSpeed,
		&result.ContinuousDriveOverCount, &result.ContinuousDriveMaxTime, &result.ContinuousDriveTotalTime,
		&result.WaveDriveCount, &result.WaveDriveMaxTime, &result.WaveDriveMaxSpeedDiff,
		&result.LocalSpeedScore, &result.ExpressSpeedScore, &result.DedicatedSpeedScore,
		&result.LocalDistanceScore, &result.ExpressDistanceScore, &result.DedicatedDistanceScore,
		&result.RapidAccelScore, &result.RapidDecelScore, &result.RapidCurveScore,
		&result.ActualLowSpeedRotationScore, &result.ActualHighSpeedRotationScore,
		&result.EmptyLowSpeedRotationScore, &result.EmptyHighSpeedRotationScore,
		&result.IdlingScore, &result.ContinuousDriveScore, &result.WaveDriveScore,
		&result.SafetyScore, &result.EconomyScore, &result.TotalScore,
	)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetByUUID retrieves a kudgivt record by UUID
func (r *KudgivtRepository) GetByUUID(ctx context.Context, uuid string) (*Kudgivt, error) {
	query := `
		SELECT
			"UUID", "OrganizationID", "Hash", "Created", "Deleted", "KudguriUuid",
			"UnkouNo", "ReadDate", "UnkouDate",
			"OfficeCd", "OfficeName", "VehicleCd", "VehicleName",
			"DriverCd1", "DriverName1", "TargetDriverType", "TargetDriverCd", "TargetDriverName",
			"ClockInDatetime", "ClockOutDatetime",
			"DepartureDatetime", "ReturnDatetime",
			"DepartureMeter", "ReturnMeter", "TotalMileage",
			"DestinationCityName", "DestinationPlaceName",
			"ActualMileage",
			"LocalDriveTime", "ExpressDriveTime", "BypassDriveTime", "ActualDriveTime", "EmptyDriveTime",
			"Work1Time", "Work2Time", "Work3Time", "Work4Time", "Work5Time",
			"Work6Time", "Work7Time", "Work8Time", "Work9Time", "Work10Time",
			"State1Distance", "State2Distance", "State3Distance", "State4Distance", "State5Distance",
			"State1Time", "State2Time", "State3Time", "State4Time", "State5Time",
			"OwnMainFuel", "OwnMainAdditive", "OwnConsumable",
			"OtherMainFuel", "OtherMainAdditive", "OtherConsumable",
			"LocalSpeedOverMax", "LocalSpeedOverTime", "LocalSpeedOverCount",
			"ExpressSpeedOverMax", "ExpressSpeedOverTime", "ExpressSpeedOverCount",
			"DedicatedSpeedOverMax", "DedicatedSpeedOverTime", "DedicatedSpeedOverCount",
			"IdlingTime", "IdlingTimeCount",
			"RotationOverMax", "RotationOverCount", "RotationOverTime",
			"RapidAccelCount1", "RapidAccelCount2", "RapidAccelCount3", "RapidAccelCount4", "RapidAccelCount5",
			"RapidAccelMax", "RapidAccelMaxSpeed",
			"RapidDecelCount1", "RapidDecelCount2", "RapidDecelCount3", "RapidDecelCount4", "RapidDecelCount5",
			"RapidDecelMax", "RapidDecelMaxSpeed",
			"RapidCurveCount1", "RapidCurveCount2", "RapidCurveCount3", "RapidCurveCount4", "RapidCurveCount5",
			"RapidCurveMax", "RapidCurveMaxSpeed",
			"ContinuousDriveOverCount", "ContinuousDriveMaxTime", "ContinuousDriveTotalTime",
			"WaveDriveCount", "WaveDriveMaxTime", "WaveDriveMaxSpeedDiff",
			"LocalSpeedScore", "ExpressSpeedScore", "DedicatedSpeedScore",
			"LocalDistanceScore", "ExpressDistanceScore", "DedicatedDistanceScore",
			"RapidAccelScore", "RapidDecelScore", "RapidCurveScore",
			"ActualLowSpeedRotationScore", "ActualHighSpeedRotationScore",
			"EmptyLowSpeedRotationScore", "EmptyHighSpeedRotationScore",
			"IdlingScore", "ContinuousDriveScore", "WaveDriveScore",
			"SafetyScore", "EconomyScore", "TotalScore"
		FROM kudgivt
		WHERE "UUID" = $1 AND "Deleted" IS NULL
	`

	var k Kudgivt
	err := r.db.QueryRow(ctx, query, uuid).Scan(
		&k.UUID, &k.OrganizationID, &k.Hash, &k.Created, &k.Deleted, &k.KudguriUuid,
		&k.UnkouNo, &k.ReadDate, &k.UnkouDate,
		&k.OfficeCd, &k.OfficeName, &k.VehicleCd, &k.VehicleName,
		&k.DriverCd1, &k.DriverName1, &k.TargetDriverType, &k.TargetDriverCd, &k.TargetDriverName,
		&k.ClockInDatetime, &k.ClockOutDatetime,
		&k.DepartureDatetime, &k.ReturnDatetime,
		&k.DepartureMeter, &k.ReturnMeter, &k.TotalMileage,
		&k.DestinationCityName, &k.DestinationPlaceName,
		&k.ActualMileage,
		&k.LocalDriveTime, &k.ExpressDriveTime, &k.BypassDriveTime, &k.ActualDriveTime, &k.EmptyDriveTime,
		&k.Work1Time, &k.Work2Time, &k.Work3Time, &k.Work4Time, &k.Work5Time,
		&k.Work6Time, &k.Work7Time, &k.Work8Time, &k.Work9Time, &k.Work10Time,
		&k.State1Distance, &k.State2Distance, &k.State3Distance, &k.State4Distance, &k.State5Distance,
		&k.State1Time, &k.State2Time, &k.State3Time, &k.State4Time, &k.State5Time,
		&k.OwnMainFuel, &k.OwnMainAdditive, &k.OwnConsumable,
		&k.OtherMainFuel, &k.OtherMainAdditive, &k.OtherConsumable,
		&k.LocalSpeedOverMax, &k.LocalSpeedOverTime, &k.LocalSpeedOverCount,
		&k.ExpressSpeedOverMax, &k.ExpressSpeedOverTime, &k.ExpressSpeedOverCount,
		&k.DedicatedSpeedOverMax, &k.DedicatedSpeedOverTime, &k.DedicatedSpeedOverCount,
		&k.IdlingTime, &k.IdlingTimeCount,
		&k.RotationOverMax, &k.RotationOverCount, &k.RotationOverTime,
		&k.RapidAccelCount1, &k.RapidAccelCount2, &k.RapidAccelCount3, &k.RapidAccelCount4, &k.RapidAccelCount5,
		&k.RapidAccelMax, &k.RapidAccelMaxSpeed,
		&k.RapidDecelCount1, &k.RapidDecelCount2, &k.RapidDecelCount3, &k.RapidDecelCount4, &k.RapidDecelCount5,
		&k.RapidDecelMax, &k.RapidDecelMaxSpeed,
		&k.RapidCurveCount1, &k.RapidCurveCount2, &k.RapidCurveCount3, &k.RapidCurveCount4, &k.RapidCurveCount5,
		&k.RapidCurveMax, &k.RapidCurveMaxSpeed,
		&k.ContinuousDriveOverCount, &k.ContinuousDriveMaxTime, &k.ContinuousDriveTotalTime,
		&k.WaveDriveCount, &k.WaveDriveMaxTime, &k.WaveDriveMaxSpeedDiff,
		&k.LocalSpeedScore, &k.ExpressSpeedScore, &k.DedicatedSpeedScore,
		&k.LocalDistanceScore, &k.ExpressDistanceScore, &k.DedicatedDistanceScore,
		&k.RapidAccelScore, &k.RapidDecelScore, &k.RapidCurveScore,
		&k.ActualLowSpeedRotationScore, &k.ActualHighSpeedRotationScore,
		&k.EmptyLowSpeedRotationScore, &k.EmptyHighSpeedRotationScore,
		&k.IdlingScore, &k.ContinuousDriveScore, &k.WaveDriveScore,
		&k.SafetyScore, &k.EconomyScore, &k.TotalScore,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrKudgivtNotFound
		}
		return nil, err
	}

	return &k, nil
}

// Update modifies an existing kudgivt record
func (r *KudgivtRepository) Update(ctx context.Context, k *Kudgivt) (*Kudgivt, error) {
	query := `
		UPDATE kudgivt
		SET
			"OrganizationID" = $2, "Hash" = $3, "Created" = $4, "Deleted" = $5, "KudguriUuid" = $6,
			"UnkouNo" = $7, "ReadDate" = $8, "UnkouDate" = $9,
			"OfficeCd" = $10, "OfficeName" = $11, "VehicleCd" = $12, "VehicleName" = $13,
			"DriverCd1" = $14, "DriverName1" = $15, "TargetDriverType" = $16, "TargetDriverCd" = $17, "TargetDriverName" = $18,
			"ClockInDatetime" = $19, "ClockOutDatetime" = $20,
			"DepartureDatetime" = $21, "ReturnDatetime" = $22,
			"DepartureMeter" = $23, "ReturnMeter" = $24, "TotalMileage" = $25,
			"DestinationCityName" = $26, "DestinationPlaceName" = $27,
			"ActualMileage" = $28,
			"LocalDriveTime" = $29, "ExpressDriveTime" = $30, "BypassDriveTime" = $31, "ActualDriveTime" = $32, "EmptyDriveTime" = $33,
			"Work1Time" = $34, "Work2Time" = $35, "Work3Time" = $36, "Work4Time" = $37, "Work5Time" = $38,
			"Work6Time" = $39, "Work7Time" = $40, "Work8Time" = $41, "Work9Time" = $42, "Work10Time" = $43,
			"State1Distance" = $44, "State2Distance" = $45, "State3Distance" = $46, "State4Distance" = $47, "State5Distance" = $48,
			"State1Time" = $49, "State2Time" = $50, "State3Time" = $51, "State4Time" = $52, "State5Time" = $53,
			"OwnMainFuel" = $54, "OwnMainAdditive" = $55, "OwnConsumable" = $56,
			"OtherMainFuel" = $57, "OtherMainAdditive" = $58, "OtherConsumable" = $59,
			"LocalSpeedOverMax" = $60, "LocalSpeedOverTime" = $61, "LocalSpeedOverCount" = $62,
			"ExpressSpeedOverMax" = $63, "ExpressSpeedOverTime" = $64, "ExpressSpeedOverCount" = $65,
			"DedicatedSpeedOverMax" = $66, "DedicatedSpeedOverTime" = $67, "DedicatedSpeedOverCount" = $68,
			"IdlingTime" = $69, "IdlingTimeCount" = $70,
			"RotationOverMax" = $71, "RotationOverCount" = $72, "RotationOverTime" = $73,
			"RapidAccelCount1" = $74, "RapidAccelCount2" = $75, "RapidAccelCount3" = $76, "RapidAccelCount4" = $77, "RapidAccelCount5" = $78,
			"RapidAccelMax" = $79, "RapidAccelMaxSpeed" = $80,
			"RapidDecelCount1" = $81, "RapidDecelCount2" = $82, "RapidDecelCount3" = $83, "RapidDecelCount4" = $84, "RapidDecelCount5" = $85,
			"RapidDecelMax" = $86, "RapidDecelMaxSpeed" = $87,
			"RapidCurveCount1" = $88, "RapidCurveCount2" = $89, "RapidCurveCount3" = $90, "RapidCurveCount4" = $91, "RapidCurveCount5" = $92,
			"RapidCurveMax" = $93, "RapidCurveMaxSpeed" = $94,
			"ContinuousDriveOverCount" = $95, "ContinuousDriveMaxTime" = $96, "ContinuousDriveTotalTime" = $97,
			"WaveDriveCount" = $98, "WaveDriveMaxTime" = $99, "WaveDriveMaxSpeedDiff" = $100,
			"LocalSpeedScore" = $101, "ExpressSpeedScore" = $102, "DedicatedSpeedScore" = $103,
			"LocalDistanceScore" = $104, "ExpressDistanceScore" = $105, "DedicatedDistanceScore" = $106,
			"RapidAccelScore" = $107, "RapidDecelScore" = $108, "RapidCurveScore" = $109,
			"ActualLowSpeedRotationScore" = $110, "ActualHighSpeedRotationScore" = $111,
			"EmptyLowSpeedRotationScore" = $112, "EmptyHighSpeedRotationScore" = $113,
			"IdlingScore" = $114, "ContinuousDriveScore" = $115, "WaveDriveScore" = $116,
			"SafetyScore" = $117, "EconomyScore" = $118, "TotalScore" = $119
		WHERE "UUID" = $1 AND "Deleted" IS NULL
		RETURNING
			"UUID", "OrganizationID", "Hash", "Created", "Deleted", "KudguriUuid",
			"UnkouNo", "ReadDate", "UnkouDate",
			"OfficeCd", "OfficeName", "VehicleCd", "VehicleName",
			"DriverCd1", "DriverName1", "TargetDriverType", "TargetDriverCd", "TargetDriverName",
			"ClockInDatetime", "ClockOutDatetime",
			"DepartureDatetime", "ReturnDatetime",
			"DepartureMeter", "ReturnMeter", "TotalMileage",
			"DestinationCityName", "DestinationPlaceName",
			"ActualMileage",
			"LocalDriveTime", "ExpressDriveTime", "BypassDriveTime", "ActualDriveTime", "EmptyDriveTime",
			"Work1Time", "Work2Time", "Work3Time", "Work4Time", "Work5Time",
			"Work6Time", "Work7Time", "Work8Time", "Work9Time", "Work10Time",
			"State1Distance", "State2Distance", "State3Distance", "State4Distance", "State5Distance",
			"State1Time", "State2Time", "State3Time", "State4Time", "State5Time",
			"OwnMainFuel", "OwnMainAdditive", "OwnConsumable",
			"OtherMainFuel", "OtherMainAdditive", "OtherConsumable",
			"LocalSpeedOverMax", "LocalSpeedOverTime", "LocalSpeedOverCount",
			"ExpressSpeedOverMax", "ExpressSpeedOverTime", "ExpressSpeedOverCount",
			"DedicatedSpeedOverMax", "DedicatedSpeedOverTime", "DedicatedSpeedOverCount",
			"IdlingTime", "IdlingTimeCount",
			"RotationOverMax", "RotationOverCount", "RotationOverTime",
			"RapidAccelCount1", "RapidAccelCount2", "RapidAccelCount3", "RapidAccelCount4", "RapidAccelCount5",
			"RapidAccelMax", "RapidAccelMaxSpeed",
			"RapidDecelCount1", "RapidDecelCount2", "RapidDecelCount3", "RapidDecelCount4", "RapidDecelCount5",
			"RapidDecelMax", "RapidDecelMaxSpeed",
			"RapidCurveCount1", "RapidCurveCount2", "RapidCurveCount3", "RapidCurveCount4", "RapidCurveCount5",
			"RapidCurveMax", "RapidCurveMaxSpeed",
			"ContinuousDriveOverCount", "ContinuousDriveMaxTime", "ContinuousDriveTotalTime",
			"WaveDriveCount", "WaveDriveMaxTime", "WaveDriveMaxSpeedDiff",
			"LocalSpeedScore", "ExpressSpeedScore", "DedicatedSpeedScore",
			"LocalDistanceScore", "ExpressDistanceScore", "DedicatedDistanceScore",
			"RapidAccelScore", "RapidDecelScore", "RapidCurveScore",
			"ActualLowSpeedRotationScore", "ActualHighSpeedRotationScore",
			"EmptyLowSpeedRotationScore", "EmptyHighSpeedRotationScore",
			"IdlingScore", "ContinuousDriveScore", "WaveDriveScore",
			"SafetyScore", "EconomyScore", "TotalScore"
	`

	var result Kudgivt
	err := r.db.QueryRow(ctx, query,
		k.UUID, k.OrganizationID, k.Hash, k.Created, k.Deleted, k.KudguriUuid,
		k.UnkouNo, k.ReadDate, k.UnkouDate,
		k.OfficeCd, k.OfficeName, k.VehicleCd, k.VehicleName,
		k.DriverCd1, k.DriverName1, k.TargetDriverType, k.TargetDriverCd, k.TargetDriverName,
		k.ClockInDatetime, k.ClockOutDatetime,
		k.DepartureDatetime, k.ReturnDatetime,
		k.DepartureMeter, k.ReturnMeter, k.TotalMileage,
		k.DestinationCityName, k.DestinationPlaceName,
		k.ActualMileage,
		k.LocalDriveTime, k.ExpressDriveTime, k.BypassDriveTime, k.ActualDriveTime, k.EmptyDriveTime,
		k.Work1Time, k.Work2Time, k.Work3Time, k.Work4Time, k.Work5Time,
		k.Work6Time, k.Work7Time, k.Work8Time, k.Work9Time, k.Work10Time,
		k.State1Distance, k.State2Distance, k.State3Distance, k.State4Distance, k.State5Distance,
		k.State1Time, k.State2Time, k.State3Time, k.State4Time, k.State5Time,
		k.OwnMainFuel, k.OwnMainAdditive, k.OwnConsumable,
		k.OtherMainFuel, k.OtherMainAdditive, k.OtherConsumable,
		k.LocalSpeedOverMax, k.LocalSpeedOverTime, k.LocalSpeedOverCount,
		k.ExpressSpeedOverMax, k.ExpressSpeedOverTime, k.ExpressSpeedOverCount,
		k.DedicatedSpeedOverMax, k.DedicatedSpeedOverTime, k.DedicatedSpeedOverCount,
		k.IdlingTime, k.IdlingTimeCount,
		k.RotationOverMax, k.RotationOverCount, k.RotationOverTime,
		k.RapidAccelCount1, k.RapidAccelCount2, k.RapidAccelCount3, k.RapidAccelCount4, k.RapidAccelCount5,
		k.RapidAccelMax, k.RapidAccelMaxSpeed,
		k.RapidDecelCount1, k.RapidDecelCount2, k.RapidDecelCount3, k.RapidDecelCount4, k.RapidDecelCount5,
		k.RapidDecelMax, k.RapidDecelMaxSpeed,
		k.RapidCurveCount1, k.RapidCurveCount2, k.RapidCurveCount3, k.RapidCurveCount4, k.RapidCurveCount5,
		k.RapidCurveMax, k.RapidCurveMaxSpeed,
		k.ContinuousDriveOverCount, k.ContinuousDriveMaxTime, k.ContinuousDriveTotalTime,
		k.WaveDriveCount, k.WaveDriveMaxTime, k.WaveDriveMaxSpeedDiff,
		k.LocalSpeedScore, k.ExpressSpeedScore, k.DedicatedSpeedScore,
		k.LocalDistanceScore, k.ExpressDistanceScore, k.DedicatedDistanceScore,
		k.RapidAccelScore, k.RapidDecelScore, k.RapidCurveScore,
		k.ActualLowSpeedRotationScore, k.ActualHighSpeedRotationScore,
		k.EmptyLowSpeedRotationScore, k.EmptyHighSpeedRotationScore,
		k.IdlingScore, k.ContinuousDriveScore, k.WaveDriveScore,
		k.SafetyScore, k.EconomyScore, k.TotalScore,
	).Scan(
		&result.UUID, &result.OrganizationID, &result.Hash, &result.Created, &result.Deleted, &result.KudguriUuid,
		&result.UnkouNo, &result.ReadDate, &result.UnkouDate,
		&result.OfficeCd, &result.OfficeName, &result.VehicleCd, &result.VehicleName,
		&result.DriverCd1, &result.DriverName1, &result.TargetDriverType, &result.TargetDriverCd, &result.TargetDriverName,
		&result.ClockInDatetime, &result.ClockOutDatetime,
		&result.DepartureDatetime, &result.ReturnDatetime,
		&result.DepartureMeter, &result.ReturnMeter, &result.TotalMileage,
		&result.DestinationCityName, &result.DestinationPlaceName,
		&result.ActualMileage,
		&result.LocalDriveTime, &result.ExpressDriveTime, &result.BypassDriveTime, &result.ActualDriveTime, &result.EmptyDriveTime,
		&result.Work1Time, &result.Work2Time, &result.Work3Time, &result.Work4Time, &result.Work5Time,
		&result.Work6Time, &result.Work7Time, &result.Work8Time, &result.Work9Time, &result.Work10Time,
		&result.State1Distance, &result.State2Distance, &result.State3Distance, &result.State4Distance, &result.State5Distance,
		&result.State1Time, &result.State2Time, &result.State3Time, &result.State4Time, &result.State5Time,
		&result.OwnMainFuel, &result.OwnMainAdditive, &result.OwnConsumable,
		&result.OtherMainFuel, &result.OtherMainAdditive, &result.OtherConsumable,
		&result.LocalSpeedOverMax, &result.LocalSpeedOverTime, &result.LocalSpeedOverCount,
		&result.ExpressSpeedOverMax, &result.ExpressSpeedOverTime, &result.ExpressSpeedOverCount,
		&result.DedicatedSpeedOverMax, &result.DedicatedSpeedOverTime, &result.DedicatedSpeedOverCount,
		&result.IdlingTime, &result.IdlingTimeCount,
		&result.RotationOverMax, &result.RotationOverCount, &result.RotationOverTime,
		&result.RapidAccelCount1, &result.RapidAccelCount2, &result.RapidAccelCount3, &result.RapidAccelCount4, &result.RapidAccelCount5,
		&result.RapidAccelMax, &result.RapidAccelMaxSpeed,
		&result.RapidDecelCount1, &result.RapidDecelCount2, &result.RapidDecelCount3, &result.RapidDecelCount4, &result.RapidDecelCount5,
		&result.RapidDecelMax, &result.RapidDecelMaxSpeed,
		&result.RapidCurveCount1, &result.RapidCurveCount2, &result.RapidCurveCount3, &result.RapidCurveCount4, &result.RapidCurveCount5,
		&result.RapidCurveMax, &result.RapidCurveMaxSpeed,
		&result.ContinuousDriveOverCount, &result.ContinuousDriveMaxTime, &result.ContinuousDriveTotalTime,
		&result.WaveDriveCount, &result.WaveDriveMaxTime, &result.WaveDriveMaxSpeedDiff,
		&result.LocalSpeedScore, &result.ExpressSpeedScore, &result.DedicatedSpeedScore,
		&result.LocalDistanceScore, &result.ExpressDistanceScore, &result.DedicatedDistanceScore,
		&result.RapidAccelScore, &result.RapidDecelScore, &result.RapidCurveScore,
		&result.ActualLowSpeedRotationScore, &result.ActualHighSpeedRotationScore,
		&result.EmptyLowSpeedRotationScore, &result.EmptyHighSpeedRotationScore,
		&result.IdlingScore, &result.ContinuousDriveScore, &result.WaveDriveScore,
		&result.SafetyScore, &result.EconomyScore, &result.TotalScore,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrKudgivtNotFound
		}
		return nil, err
	}

	return &result, nil
}

// Delete soft-deletes a kudgivt record
func (r *KudgivtRepository) Delete(ctx context.Context, uuid string) error {
	query := `
		UPDATE kudgivt
		SET "Deleted" = NOW()::text
		WHERE "UUID" = $1 AND "Deleted" IS NULL
	`

	result, err := r.db.Exec(ctx, query, uuid)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrKudgivtNotFound
	}

	return nil
}

// ListByOrganization retrieves kudgivt records for a specific organization with pagination
func (r *KudgivtRepository) ListByOrganization(ctx context.Context, organizationID string, limit int, offset int) ([]*Kudgivt, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	query := `
		SELECT
			"UUID", "OrganizationID", "Hash", "Created", "Deleted", "KudguriUuid",
			"UnkouNo", "ReadDate", "UnkouDate",
			"OfficeCd", "OfficeName", "VehicleCd", "VehicleName",
			"DriverCd1", "DriverName1", "TargetDriverType", "TargetDriverCd", "TargetDriverName",
			"ClockInDatetime", "ClockOutDatetime",
			"DepartureDatetime", "ReturnDatetime",
			"DepartureMeter", "ReturnMeter", "TotalMileage",
			"DestinationCityName", "DestinationPlaceName",
			"ActualMileage",
			"LocalDriveTime", "ExpressDriveTime", "BypassDriveTime", "ActualDriveTime", "EmptyDriveTime",
			"Work1Time", "Work2Time", "Work3Time", "Work4Time", "Work5Time",
			"Work6Time", "Work7Time", "Work8Time", "Work9Time", "Work10Time",
			"State1Distance", "State2Distance", "State3Distance", "State4Distance", "State5Distance",
			"State1Time", "State2Time", "State3Time", "State4Time", "State5Time",
			"OwnMainFuel", "OwnMainAdditive", "OwnConsumable",
			"OtherMainFuel", "OtherMainAdditive", "OtherConsumable",
			"LocalSpeedOverMax", "LocalSpeedOverTime", "LocalSpeedOverCount",
			"ExpressSpeedOverMax", "ExpressSpeedOverTime", "ExpressSpeedOverCount",
			"DedicatedSpeedOverMax", "DedicatedSpeedOverTime", "DedicatedSpeedOverCount",
			"IdlingTime", "IdlingTimeCount",
			"RotationOverMax", "RotationOverCount", "RotationOverTime",
			"RapidAccelCount1", "RapidAccelCount2", "RapidAccelCount3", "RapidAccelCount4", "RapidAccelCount5",
			"RapidAccelMax", "RapidAccelMaxSpeed",
			"RapidDecelCount1", "RapidDecelCount2", "RapidDecelCount3", "RapidDecelCount4", "RapidDecelCount5",
			"RapidDecelMax", "RapidDecelMaxSpeed",
			"RapidCurveCount1", "RapidCurveCount2", "RapidCurveCount3", "RapidCurveCount4", "RapidCurveCount5",
			"RapidCurveMax", "RapidCurveMaxSpeed",
			"ContinuousDriveOverCount", "ContinuousDriveMaxTime", "ContinuousDriveTotalTime",
			"WaveDriveCount", "WaveDriveMaxTime", "WaveDriveMaxSpeedDiff",
			"LocalSpeedScore", "ExpressSpeedScore", "DedicatedSpeedScore",
			"LocalDistanceScore", "ExpressDistanceScore", "DedicatedDistanceScore",
			"RapidAccelScore", "RapidDecelScore", "RapidCurveScore",
			"ActualLowSpeedRotationScore", "ActualHighSpeedRotationScore",
			"EmptyLowSpeedRotationScore", "EmptyHighSpeedRotationScore",
			"IdlingScore", "ContinuousDriveScore", "WaveDriveScore",
			"SafetyScore", "EconomyScore", "TotalScore"
		FROM kudgivt
		WHERE "OrganizationID" = $1 AND "Deleted" IS NULL
		ORDER BY "Created" DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, organizationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*Kudgivt
	for rows.Next() {
		var k Kudgivt
		err := rows.Scan(
			&k.UUID, &k.OrganizationID, &k.Hash, &k.Created, &k.Deleted, &k.KudguriUuid,
			&k.UnkouNo, &k.ReadDate, &k.UnkouDate,
			&k.OfficeCd, &k.OfficeName, &k.VehicleCd, &k.VehicleName,
			&k.DriverCd1, &k.DriverName1, &k.TargetDriverType, &k.TargetDriverCd, &k.TargetDriverName,
			&k.ClockInDatetime, &k.ClockOutDatetime,
			&k.DepartureDatetime, &k.ReturnDatetime,
			&k.DepartureMeter, &k.ReturnMeter, &k.TotalMileage,
			&k.DestinationCityName, &k.DestinationPlaceName,
			&k.ActualMileage,
			&k.LocalDriveTime, &k.ExpressDriveTime, &k.BypassDriveTime, &k.ActualDriveTime, &k.EmptyDriveTime,
			&k.Work1Time, &k.Work2Time, &k.Work3Time, &k.Work4Time, &k.Work5Time,
			&k.Work6Time, &k.Work7Time, &k.Work8Time, &k.Work9Time, &k.Work10Time,
			&k.State1Distance, &k.State2Distance, &k.State3Distance, &k.State4Distance, &k.State5Distance,
			&k.State1Time, &k.State2Time, &k.State3Time, &k.State4Time, &k.State5Time,
			&k.OwnMainFuel, &k.OwnMainAdditive, &k.OwnConsumable,
			&k.OtherMainFuel, &k.OtherMainAdditive, &k.OtherConsumable,
			&k.LocalSpeedOverMax, &k.LocalSpeedOverTime, &k.LocalSpeedOverCount,
			&k.ExpressSpeedOverMax, &k.ExpressSpeedOverTime, &k.ExpressSpeedOverCount,
			&k.DedicatedSpeedOverMax, &k.DedicatedSpeedOverTime, &k.DedicatedSpeedOverCount,
			&k.IdlingTime, &k.IdlingTimeCount,
			&k.RotationOverMax, &k.RotationOverCount, &k.RotationOverTime,
			&k.RapidAccelCount1, &k.RapidAccelCount2, &k.RapidAccelCount3, &k.RapidAccelCount4, &k.RapidAccelCount5,
			&k.RapidAccelMax, &k.RapidAccelMaxSpeed,
			&k.RapidDecelCount1, &k.RapidDecelCount2, &k.RapidDecelCount3, &k.RapidDecelCount4, &k.RapidDecelCount5,
			&k.RapidDecelMax, &k.RapidDecelMaxSpeed,
			&k.RapidCurveCount1, &k.RapidCurveCount2, &k.RapidCurveCount3, &k.RapidCurveCount4, &k.RapidCurveCount5,
			&k.RapidCurveMax, &k.RapidCurveMaxSpeed,
			&k.ContinuousDriveOverCount, &k.ContinuousDriveMaxTime, &k.ContinuousDriveTotalTime,
			&k.WaveDriveCount, &k.WaveDriveMaxTime, &k.WaveDriveMaxSpeedDiff,
			&k.LocalSpeedScore, &k.ExpressSpeedScore, &k.DedicatedSpeedScore,
			&k.LocalDistanceScore, &k.ExpressDistanceScore, &k.DedicatedDistanceScore,
			&k.RapidAccelScore, &k.RapidDecelScore, &k.RapidCurveScore,
			&k.ActualLowSpeedRotationScore, &k.ActualHighSpeedRotationScore,
			&k.EmptyLowSpeedRotationScore, &k.EmptyHighSpeedRotationScore,
			&k.IdlingScore, &k.ContinuousDriveScore, &k.WaveDriveScore,
			&k.SafetyScore, &k.EconomyScore, &k.TotalScore,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, &k)
	}

	return results, rows.Err()
}

// List retrieves all kudgivt records with pagination
func (r *KudgivtRepository) List(ctx context.Context, limit int, offset int) ([]*Kudgivt, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	query := `
		SELECT
			"UUID", "OrganizationID", "Hash", "Created", "Deleted", "KudguriUuid",
			"UnkouNo", "ReadDate", "UnkouDate",
			"OfficeCd", "OfficeName", "VehicleCd", "VehicleName",
			"DriverCd1", "DriverName1", "TargetDriverType", "TargetDriverCd", "TargetDriverName",
			"ClockInDatetime", "ClockOutDatetime",
			"DepartureDatetime", "ReturnDatetime",
			"DepartureMeter", "ReturnMeter", "TotalMileage",
			"DestinationCityName", "DestinationPlaceName",
			"ActualMileage",
			"LocalDriveTime", "ExpressDriveTime", "BypassDriveTime", "ActualDriveTime", "EmptyDriveTime",
			"Work1Time", "Work2Time", "Work3Time", "Work4Time", "Work5Time",
			"Work6Time", "Work7Time", "Work8Time", "Work9Time", "Work10Time",
			"State1Distance", "State2Distance", "State3Distance", "State4Distance", "State5Distance",
			"State1Time", "State2Time", "State3Time", "State4Time", "State5Time",
			"OwnMainFuel", "OwnMainAdditive", "OwnConsumable",
			"OtherMainFuel", "OtherMainAdditive", "OtherConsumable",
			"LocalSpeedOverMax", "LocalSpeedOverTime", "LocalSpeedOverCount",
			"ExpressSpeedOverMax", "ExpressSpeedOverTime", "ExpressSpeedOverCount",
			"DedicatedSpeedOverMax", "DedicatedSpeedOverTime", "DedicatedSpeedOverCount",
			"IdlingTime", "IdlingTimeCount",
			"RotationOverMax", "RotationOverCount", "RotationOverTime",
			"RapidAccelCount1", "RapidAccelCount2", "RapidAccelCount3", "RapidAccelCount4", "RapidAccelCount5",
			"RapidAccelMax", "RapidAccelMaxSpeed",
			"RapidDecelCount1", "RapidDecelCount2", "RapidDecelCount3", "RapidDecelCount4", "RapidDecelCount5",
			"RapidDecelMax", "RapidDecelMaxSpeed",
			"RapidCurveCount1", "RapidCurveCount2", "RapidCurveCount3", "RapidCurveCount4", "RapidCurveCount5",
			"RapidCurveMax", "RapidCurveMaxSpeed",
			"ContinuousDriveOverCount", "ContinuousDriveMaxTime", "ContinuousDriveTotalTime",
			"WaveDriveCount", "WaveDriveMaxTime", "WaveDriveMaxSpeedDiff",
			"LocalSpeedScore", "ExpressSpeedScore", "DedicatedSpeedScore",
			"LocalDistanceScore", "ExpressDistanceScore", "DedicatedDistanceScore",
			"RapidAccelScore", "RapidDecelScore", "RapidCurveScore",
			"ActualLowSpeedRotationScore", "ActualHighSpeedRotationScore",
			"EmptyLowSpeedRotationScore", "EmptyHighSpeedRotationScore",
			"IdlingScore", "ContinuousDriveScore", "WaveDriveScore",
			"SafetyScore", "EconomyScore", "TotalScore"
		FROM kudgivt
		WHERE "Deleted" IS NULL
		ORDER BY "Created" DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*Kudgivt
	for rows.Next() {
		var k Kudgivt
		err := rows.Scan(
			&k.UUID, &k.OrganizationID, &k.Hash, &k.Created, &k.Deleted, &k.KudguriUuid,
			&k.UnkouNo, &k.ReadDate, &k.UnkouDate,
			&k.OfficeCd, &k.OfficeName, &k.VehicleCd, &k.VehicleName,
			&k.DriverCd1, &k.DriverName1, &k.TargetDriverType, &k.TargetDriverCd, &k.TargetDriverName,
			&k.ClockInDatetime, &k.ClockOutDatetime,
			&k.DepartureDatetime, &k.ReturnDatetime,
			&k.DepartureMeter, &k.ReturnMeter, &k.TotalMileage,
			&k.DestinationCityName, &k.DestinationPlaceName,
			&k.ActualMileage,
			&k.LocalDriveTime, &k.ExpressDriveTime, &k.BypassDriveTime, &k.ActualDriveTime, &k.EmptyDriveTime,
			&k.Work1Time, &k.Work2Time, &k.Work3Time, &k.Work4Time, &k.Work5Time,
			&k.Work6Time, &k.Work7Time, &k.Work8Time, &k.Work9Time, &k.Work10Time,
			&k.State1Distance, &k.State2Distance, &k.State3Distance, &k.State4Distance, &k.State5Distance,
			&k.State1Time, &k.State2Time, &k.State3Time, &k.State4Time, &k.State5Time,
			&k.OwnMainFuel, &k.OwnMainAdditive, &k.OwnConsumable,
			&k.OtherMainFuel, &k.OtherMainAdditive, &k.OtherConsumable,
			&k.LocalSpeedOverMax, &k.LocalSpeedOverTime, &k.LocalSpeedOverCount,
			&k.ExpressSpeedOverMax, &k.ExpressSpeedOverTime, &k.ExpressSpeedOverCount,
			&k.DedicatedSpeedOverMax, &k.DedicatedSpeedOverTime, &k.DedicatedSpeedOverCount,
			&k.IdlingTime, &k.IdlingTimeCount,
			&k.RotationOverMax, &k.RotationOverCount, &k.RotationOverTime,
			&k.RapidAccelCount1, &k.RapidAccelCount2, &k.RapidAccelCount3, &k.RapidAccelCount4, &k.RapidAccelCount5,
			&k.RapidAccelMax, &k.RapidAccelMaxSpeed,
			&k.RapidDecelCount1, &k.RapidDecelCount2, &k.RapidDecelCount3, &k.RapidDecelCount4, &k.RapidDecelCount5,
			&k.RapidDecelMax, &k.RapidDecelMaxSpeed,
			&k.RapidCurveCount1, &k.RapidCurveCount2, &k.RapidCurveCount3, &k.RapidCurveCount4, &k.RapidCurveCount5,
			&k.RapidCurveMax, &k.RapidCurveMaxSpeed,
			&k.ContinuousDriveOverCount, &k.ContinuousDriveMaxTime, &k.ContinuousDriveTotalTime,
			&k.WaveDriveCount, &k.WaveDriveMaxTime, &k.WaveDriveMaxSpeedDiff,
			&k.LocalSpeedScore, &k.ExpressSpeedScore, &k.DedicatedSpeedScore,
			&k.LocalDistanceScore, &k.ExpressDistanceScore, &k.DedicatedDistanceScore,
			&k.RapidAccelScore, &k.RapidDecelScore, &k.RapidCurveScore,
			&k.ActualLowSpeedRotationScore, &k.ActualHighSpeedRotationScore,
			&k.EmptyLowSpeedRotationScore, &k.EmptyHighSpeedRotationScore,
			&k.IdlingScore, &k.ContinuousDriveScore, &k.WaveDriveScore,
			&k.SafetyScore, &k.EconomyScore, &k.TotalScore,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, &k)
	}

	return results, rows.Err()
}
