package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrDtakologsNotFound = errors.New("dtakologs not found")
)

// Dtakologs represents the database model
type Dtakologs struct {
	OrganizationID               string
	Type                         string
	AddressDispC                 *string
	AddressDispP                 *string
	AllState                     *string
	AllStateEx                   *string
	AllStateFontColor            *string
	AllStateFontColorIndex       int32
	AllStateRyoutColor           string
	BranchCd                     int32
	BranchName                   string
	ComuDateTime                 *string
	CurrentWorkCd                int32
	CurrentWorkName              *string
	DataDateTime                 string
	DataFilterType               int32
	DispFlag                     int32
	DriverCd                     int32
	DriverName                   *string
	EventVal                     *string
	GpsDirection                 int32
	GpsEnable                    int32
	GpsLatiAndLong               *string
	GpsLatitude                  int32
	GpsLongitude                 int32
	GpsSatelliteNum              int32
	OdoMeter                     *string
	OperationState               int32
	ReciveEventType              int32
	RecivePacketType             int32
	ReciveTypeColorName          *string
	ReciveTypeName               *string
	ReciveWorkCd                 int32
	Revo                         int32
	SettingTemp                  string
	SettingTemp1                 string
	SettingTemp3                 string
	SettingTemp4                 string
	Speed                        float32
	StartWorkDateTime            *string
	State                        *string
	State1                       *string
	State2                       *string
	State3                       *string
	StateFlag                    string
	SubDriverCd                  int32
	Temp1                        *string
	Temp2                        *string
	Temp3                        *string
	Temp4                        *string
	TempState                    int32
	VehicleCd                    int32
	VehicleIconColor             *string
	VehicleIconLabelForDatetime  *string
	VehicleIconLabelForDriver    *string
	VehicleIconLabelForVehicle   *string
	VehicleName                  string
}

// DtakologsRepository handles database operations for dtakologs
type DtakologsRepository struct {
	db DB
}

// NewDtakologsRepository creates a new repository
func NewDtakologsRepository(pool *pgxpool.Pool) *DtakologsRepository {
	return &DtakologsRepository{db: pool}
}

// NewDtakologsRepositoryWithDB creates a repository with custom DB interface (for testing)
func NewDtakologsRepositoryWithDB(db DB) *DtakologsRepository {
	return &DtakologsRepository{db: db}
}

// Create inserts a new dtakologs record
func (r *DtakologsRepository) Create(ctx context.Context, d *Dtakologs) error {
	query := `
		INSERT INTO dtakologs (
			organization_id, __type, "AddressDispC", "AddressDispP", "AllState", "AllStateEx",
			"AllStateFontColor", "AllStateFontColorIndex", "AllStateRyoutColor", "BranchCD",
			"BranchName", "ComuDateTime", "CurrentWorkCD", "CurrentWorkName", "DataDateTime",
			"DataFilterType", "DispFlag", "DriverCD", "DriverName", "EventVal",
			"GpsDirection", "GpsEnable", "GpsLatiAndLong", "GpsLatitude", "GpsLongitude",
			"GpsSatelliteNum", "OdoMeter", "OperationState", "ReciveEventType", "RecivePacketType",
			"ReciveTypeColorName", "ReciveTypeName", "ReciveWorkCD", "Revo", "SettingTemp",
			"SettingTemp1", "SettingTemp3", "SettingTemp4", "Speed", "StartWorkDateTime",
			"State", "State1", "State2", "State3", "StateFlag", "SubDriverCD",
			"Temp1", "Temp2", "Temp3", "Temp4", "TempState", "VehicleCD",
			"VehicleIconColor", "VehicleIconLabelForDatetime", "VehicleIconLabelForDriver",
			"VehicleIconLabelForVehicle", "VehicleName"
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
			$21, $22, $23, $24, $25, $26, $27, $28, $29, $30,
			$31, $32, $33, $34, $35, $36, $37, $38, $39, $40,
			$41, $42, $43, $44, $45, $46, $47, $48, $49, $50,
			$51, $52, $53, $54, $55, $56, $57
		)
	`

	_, err := r.db.Exec(ctx, query,
		d.OrganizationID, d.Type, d.AddressDispC, d.AddressDispP, d.AllState, d.AllStateEx,
		d.AllStateFontColor, d.AllStateFontColorIndex, d.AllStateRyoutColor, d.BranchCd,
		d.BranchName, d.ComuDateTime, d.CurrentWorkCd, d.CurrentWorkName, d.DataDateTime,
		d.DataFilterType, d.DispFlag, d.DriverCd, d.DriverName, d.EventVal,
		d.GpsDirection, d.GpsEnable, d.GpsLatiAndLong, d.GpsLatitude, d.GpsLongitude,
		d.GpsSatelliteNum, d.OdoMeter, d.OperationState, d.ReciveEventType, d.RecivePacketType,
		d.ReciveTypeColorName, d.ReciveTypeName, d.ReciveWorkCd, d.Revo, d.SettingTemp,
		d.SettingTemp1, d.SettingTemp3, d.SettingTemp4, d.Speed, d.StartWorkDateTime,
		d.State, d.State1, d.State2, d.State3, d.StateFlag, d.SubDriverCd,
		d.Temp1, d.Temp2, d.Temp3, d.Temp4, d.TempState, d.VehicleCd,
		d.VehicleIconColor, d.VehicleIconLabelForDatetime, d.VehicleIconLabelForDriver,
		d.VehicleIconLabelForVehicle, d.VehicleName,
	)

	return err
}

// GetByPrimaryKey retrieves a dtakologs record by composite primary key
func (r *DtakologsRepository) GetByPrimaryKey(ctx context.Context, organizationID, dataDateTime string, vehicleCd int32) (*Dtakologs, error) {
	query := `
		SELECT
			organization_id, __type, "AddressDispC", "AddressDispP", "AllState", "AllStateEx",
			"AllStateFontColor", "AllStateFontColorIndex", "AllStateRyoutColor", "BranchCD",
			"BranchName", "ComuDateTime", "CurrentWorkCD", "CurrentWorkName", "DataDateTime",
			"DataFilterType", "DispFlag", "DriverCD", "DriverName", "EventVal",
			"GpsDirection", "GpsEnable", "GpsLatiAndLong", "GpsLatitude", "GpsLongitude",
			"GpsSatelliteNum", "OdoMeter", "OperationState", "ReciveEventType", "RecivePacketType",
			"ReciveTypeColorName", "ReciveTypeName", "ReciveWorkCD", "Revo", "SettingTemp",
			"SettingTemp1", "SettingTemp3", "SettingTemp4", "Speed", "StartWorkDateTime",
			"State", "State1", "State2", "State3", "StateFlag", "SubDriverCD",
			"Temp1", "Temp2", "Temp3", "Temp4", "TempState", "VehicleCD",
			"VehicleIconColor", "VehicleIconLabelForDatetime", "VehicleIconLabelForDriver",
			"VehicleIconLabelForVehicle", "VehicleName"
		FROM dtakologs
		WHERE organization_id = $1 AND "DataDateTime" = $2 AND "VehicleCD" = $3
	`

	var d Dtakologs
	err := r.db.QueryRow(ctx, query, organizationID, dataDateTime, vehicleCd).Scan(
		&d.OrganizationID, &d.Type, &d.AddressDispC, &d.AddressDispP, &d.AllState, &d.AllStateEx,
		&d.AllStateFontColor, &d.AllStateFontColorIndex, &d.AllStateRyoutColor, &d.BranchCd,
		&d.BranchName, &d.ComuDateTime, &d.CurrentWorkCd, &d.CurrentWorkName, &d.DataDateTime,
		&d.DataFilterType, &d.DispFlag, &d.DriverCd, &d.DriverName, &d.EventVal,
		&d.GpsDirection, &d.GpsEnable, &d.GpsLatiAndLong, &d.GpsLatitude, &d.GpsLongitude,
		&d.GpsSatelliteNum, &d.OdoMeter, &d.OperationState, &d.ReciveEventType, &d.RecivePacketType,
		&d.ReciveTypeColorName, &d.ReciveTypeName, &d.ReciveWorkCd, &d.Revo, &d.SettingTemp,
		&d.SettingTemp1, &d.SettingTemp3, &d.SettingTemp4, &d.Speed, &d.StartWorkDateTime,
		&d.State, &d.State1, &d.State2, &d.State3, &d.StateFlag, &d.SubDriverCd,
		&d.Temp1, &d.Temp2, &d.Temp3, &d.Temp4, &d.TempState, &d.VehicleCd,
		&d.VehicleIconColor, &d.VehicleIconLabelForDatetime, &d.VehicleIconLabelForDriver,
		&d.VehicleIconLabelForVehicle, &d.VehicleName,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrDtakologsNotFound
		}
		return nil, err
	}

	return &d, nil
}

// Update modifies an existing dtakologs record
func (r *DtakologsRepository) Update(ctx context.Context, d *Dtakologs) error {
	query := `
		UPDATE dtakologs
		SET
			__type = $4, "AddressDispC" = $5, "AddressDispP" = $6, "AllState" = $7, "AllStateEx" = $8,
			"AllStateFontColor" = $9, "AllStateFontColorIndex" = $10, "AllStateRyoutColor" = $11, "BranchCD" = $12,
			"BranchName" = $13, "ComuDateTime" = $14, "CurrentWorkCD" = $15, "CurrentWorkName" = $16,
			"DataFilterType" = $17, "DispFlag" = $18, "DriverCD" = $19, "DriverName" = $20,
			"EventVal" = $21, "GpsDirection" = $22, "GpsEnable" = $23, "GpsLatiAndLong" = $24,
			"GpsLatitude" = $25, "GpsLongitude" = $26, "GpsSatelliteNum" = $27, "OdoMeter" = $28,
			"OperationState" = $29, "ReciveEventType" = $30, "RecivePacketType" = $31, "ReciveTypeColorName" = $32,
			"ReciveTypeName" = $33, "ReciveWorkCD" = $34, "Revo" = $35, "SettingTemp" = $36,
			"SettingTemp1" = $37, "SettingTemp3" = $38, "SettingTemp4" = $39, "Speed" = $40,
			"StartWorkDateTime" = $41, "State" = $42, "State1" = $43, "State2" = $44,
			"State3" = $45, "StateFlag" = $46, "SubDriverCD" = $47, "Temp1" = $48,
			"Temp2" = $49, "Temp3" = $50, "Temp4" = $51, "TempState" = $52,
			"VehicleIconColor" = $53, "VehicleIconLabelForDatetime" = $54, "VehicleIconLabelForDriver" = $55,
			"VehicleIconLabelForVehicle" = $56, "VehicleName" = $57
		WHERE organization_id = $1 AND "DataDateTime" = $2 AND "VehicleCD" = $3
	`

	result, err := r.db.Exec(ctx, query,
		d.OrganizationID, d.DataDateTime, d.VehicleCd,
		d.Type, d.AddressDispC, d.AddressDispP, d.AllState, d.AllStateEx,
		d.AllStateFontColor, d.AllStateFontColorIndex, d.AllStateRyoutColor, d.BranchCd,
		d.BranchName, d.ComuDateTime, d.CurrentWorkCd, d.CurrentWorkName,
		d.DataFilterType, d.DispFlag, d.DriverCd, d.DriverName,
		d.EventVal, d.GpsDirection, d.GpsEnable, d.GpsLatiAndLong,
		d.GpsLatitude, d.GpsLongitude, d.GpsSatelliteNum, d.OdoMeter,
		d.OperationState, d.ReciveEventType, d.RecivePacketType, d.ReciveTypeColorName,
		d.ReciveTypeName, d.ReciveWorkCd, d.Revo, d.SettingTemp,
		d.SettingTemp1, d.SettingTemp3, d.SettingTemp4, d.Speed,
		d.StartWorkDateTime, d.State, d.State1, d.State2,
		d.State3, d.StateFlag, d.SubDriverCd, d.Temp1,
		d.Temp2, d.Temp3, d.Temp4, d.TempState,
		d.VehicleIconColor, d.VehicleIconLabelForDatetime, d.VehicleIconLabelForDriver,
		d.VehicleIconLabelForVehicle, d.VehicleName,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrDtakologsNotFound
	}

	return nil
}

// Delete removes a dtakologs record by composite primary key
func (r *DtakologsRepository) Delete(ctx context.Context, organizationID, dataDateTime string, vehicleCd int32) error {
	query := `
		DELETE FROM dtakologs
		WHERE organization_id = $1 AND "DataDateTime" = $2 AND "VehicleCD" = $3
	`

	result, err := r.db.Exec(ctx, query, organizationID, dataDateTime, vehicleCd)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrDtakologsNotFound
	}

	return nil
}

// ListByOrganization retrieves dtakologs records for a specific organization with pagination
func (r *DtakologsRepository) ListByOrganization(ctx context.Context, organizationID string, limit, offset int) ([]*Dtakologs, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	query := `
		SELECT
			organization_id, __type, "AddressDispC", "AddressDispP", "AllState", "AllStateEx",
			"AllStateFontColor", "AllStateFontColorIndex", "AllStateRyoutColor", "BranchCD",
			"BranchName", "ComuDateTime", "CurrentWorkCD", "CurrentWorkName", "DataDateTime",
			"DataFilterType", "DispFlag", "DriverCD", "DriverName", "EventVal",
			"GpsDirection", "GpsEnable", "GpsLatiAndLong", "GpsLatitude", "GpsLongitude",
			"GpsSatelliteNum", "OdoMeter", "OperationState", "ReciveEventType", "RecivePacketType",
			"ReciveTypeColorName", "ReciveTypeName", "ReciveWorkCD", "Revo", "SettingTemp",
			"SettingTemp1", "SettingTemp3", "SettingTemp4", "Speed", "StartWorkDateTime",
			"State", "State1", "State2", "State3", "StateFlag", "SubDriverCD",
			"Temp1", "Temp2", "Temp3", "Temp4", "TempState", "VehicleCD",
			"VehicleIconColor", "VehicleIconLabelForDatetime", "VehicleIconLabelForDriver",
			"VehicleIconLabelForVehicle", "VehicleName"
		FROM dtakologs
		WHERE organization_id = $1
		ORDER BY "DataDateTime" DESC, "VehicleCD" ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, organizationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*Dtakologs
	for rows.Next() {
		var d Dtakologs
		err := rows.Scan(
			&d.OrganizationID, &d.Type, &d.AddressDispC, &d.AddressDispP, &d.AllState, &d.AllStateEx,
			&d.AllStateFontColor, &d.AllStateFontColorIndex, &d.AllStateRyoutColor, &d.BranchCd,
			&d.BranchName, &d.ComuDateTime, &d.CurrentWorkCd, &d.CurrentWorkName, &d.DataDateTime,
			&d.DataFilterType, &d.DispFlag, &d.DriverCd, &d.DriverName, &d.EventVal,
			&d.GpsDirection, &d.GpsEnable, &d.GpsLatiAndLong, &d.GpsLatitude, &d.GpsLongitude,
			&d.GpsSatelliteNum, &d.OdoMeter, &d.OperationState, &d.ReciveEventType, &d.RecivePacketType,
			&d.ReciveTypeColorName, &d.ReciveTypeName, &d.ReciveWorkCd, &d.Revo, &d.SettingTemp,
			&d.SettingTemp1, &d.SettingTemp3, &d.SettingTemp4, &d.Speed, &d.StartWorkDateTime,
			&d.State, &d.State1, &d.State2, &d.State3, &d.StateFlag, &d.SubDriverCd,
			&d.Temp1, &d.Temp2, &d.Temp3, &d.Temp4, &d.TempState, &d.VehicleCd,
			&d.VehicleIconColor, &d.VehicleIconLabelForDatetime, &d.VehicleIconLabelForDriver,
			&d.VehicleIconLabelForVehicle, &d.VehicleName,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, &d)
	}

	return records, rows.Err()
}

// List retrieves dtakologs records with pagination across all organizations
func (r *DtakologsRepository) List(ctx context.Context, limit, offset int) ([]*Dtakologs, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	query := `
		SELECT
			organization_id, __type, "AddressDispC", "AddressDispP", "AllState", "AllStateEx",
			"AllStateFontColor", "AllStateFontColorIndex", "AllStateRyoutColor", "BranchCD",
			"BranchName", "ComuDateTime", "CurrentWorkCD", "CurrentWorkName", "DataDateTime",
			"DataFilterType", "DispFlag", "DriverCD", "DriverName", "EventVal",
			"GpsDirection", "GpsEnable", "GpsLatiAndLong", "GpsLatitude", "GpsLongitude",
			"GpsSatelliteNum", "OdoMeter", "OperationState", "ReciveEventType", "RecivePacketType",
			"ReciveTypeColorName", "ReciveTypeName", "ReciveWorkCD", "Revo", "SettingTemp",
			"SettingTemp1", "SettingTemp3", "SettingTemp4", "Speed", "StartWorkDateTime",
			"State", "State1", "State2", "State3", "StateFlag", "SubDriverCD",
			"Temp1", "Temp2", "Temp3", "Temp4", "TempState", "VehicleCD",
			"VehicleIconColor", "VehicleIconLabelForDatetime", "VehicleIconLabelForDriver",
			"VehicleIconLabelForVehicle", "VehicleName"
		FROM dtakologs
		ORDER BY organization_id ASC, "DataDateTime" DESC, "VehicleCD" ASC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*Dtakologs
	for rows.Next() {
		var d Dtakologs
		err := rows.Scan(
			&d.OrganizationID, &d.Type, &d.AddressDispC, &d.AddressDispP, &d.AllState, &d.AllStateEx,
			&d.AllStateFontColor, &d.AllStateFontColorIndex, &d.AllStateRyoutColor, &d.BranchCd,
			&d.BranchName, &d.ComuDateTime, &d.CurrentWorkCd, &d.CurrentWorkName, &d.DataDateTime,
			&d.DataFilterType, &d.DispFlag, &d.DriverCd, &d.DriverName, &d.EventVal,
			&d.GpsDirection, &d.GpsEnable, &d.GpsLatiAndLong, &d.GpsLatitude, &d.GpsLongitude,
			&d.GpsSatelliteNum, &d.OdoMeter, &d.OperationState, &d.ReciveEventType, &d.RecivePacketType,
			&d.ReciveTypeColorName, &d.ReciveTypeName, &d.ReciveWorkCd, &d.Revo, &d.SettingTemp,
			&d.SettingTemp1, &d.SettingTemp3, &d.SettingTemp4, &d.Speed, &d.StartWorkDateTime,
			&d.State, &d.State1, &d.State2, &d.State3, &d.StateFlag, &d.SubDriverCd,
			&d.Temp1, &d.Temp2, &d.Temp3, &d.Temp4, &d.TempState, &d.VehicleCd,
			&d.VehicleIconColor, &d.VehicleIconLabelForDatetime, &d.VehicleIconLabelForDriver,
			&d.VehicleIconLabelForVehicle, &d.VehicleName,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, &d)
	}

	return records, rows.Err()
}
