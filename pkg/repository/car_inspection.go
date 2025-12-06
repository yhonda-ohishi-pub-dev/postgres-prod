package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrCarInspectionNotFound = errors.New("car inspection not found")
)

// CarInspection represents the database model
type CarInspection struct {
	OrganizationID                                              string
	CertInfoImportFileVersion                                   string
	AcceptOutputNo                                              string
	FormType                                                    string
	ElectCertMgNo                                               string
	CarID                                                       string
	ElectCertPublishdateE                                       string
	ElectCertPublishdateY                                       string
	ElectCertPublishdateM                                       string
	ElectCertPublishdateD                                       string
	GrantdateE                                                  string
	GrantdateY                                                  string
	GrantdateM                                                  string
	GrantdateD                                                  string
	TranspotationBureauChiefName                                string
	EntryNoCarNo                                                string
	RegGrantdateE                                               string
	RegGrantdateY                                               string
	RegGrantdateM                                               string
	RegGrantdateD                                               string
	FirstRegistDateE                                            string
	FirstRegistDateY                                            string
	FirstRegistDateM                                            string
	CarName                                                     string
	CarNameCode                                                 string
	CarNo                                                       string
	Model                                                       string
	EngineModel                                                 string
	OwnerNameLowLevelChar                                       string
	OwnerNameHighLevelChar                                      string
	OwnerAddressChar                                            string
	OwnerAddressNumValue                                        string
	OwnerAddressCode                                            string
	UserNameLowLevelChar                                        string
	UserNameHighLevelChar                                       string
	UserAddressChar                                             string
	UserAddressNumValue                                         string
	UserAddressCode                                             string
	UseHeadqurterChar                                           string
	UseHeadqurterNumValue                                       string
	UseHeadqurterCode                                           string
	CarKind                                                     string
	Use                                                         string
	PrivateBusiness                                             string
	CarShape                                                    string
	CarShapeCode                                                string
	NoteCap                                                     string
	Cap                                                         string
	NoteMaxLoadage                                              string
	MaxLoadage                                                  string
	NoteCarWgt                                                  string
	CarWgt                                                      string
	NoteCarTotalWgt                                             string
	CarTotalWgt                                                 string
	NoteLength                                                  string
	Length                                                      string
	NoteWidth                                                   string
	Width                                                       string
	NoteHeight                                                  string
	Height                                                      string
	FfAxWgt                                                     string
	FrAxWgt                                                     string
	RfAxWgt                                                     string
	RrAxWgt                                                     string
	Displacement                                                string
	FuelClass                                                   string
	ModelSpecifyNo                                              string
	ClassifyAroundNo                                            string
	ValidPeriodExpirDateE                                       string
	ValidPeriodExpirDateY                                       string
	ValidPeriodExpirDateM                                       string
	ValidPeriodExpirDateD                                       string
	NoteInfo                                                    string
	TwodimensionCodeInfoEntryNoCarNo                            string
	TwodimensionCodeInfoCarNo                                   string
	TwodimensionCodeInfoValidPeriodExpirDate                    string
	TwodimensionCodeInfoModel                                   string
	TwodimensionCodeInfoModelSpecifyNoClassifyAroundNo          string
	TwodimensionCodeInfoCharInfo                                string
	TwodimensionCodeInfoEngineModel                             string
	TwodimensionCodeInfoCarNoStampPlace                         string
	TwodimensionCodeInfoFirstRegistDate                         string
	TwodimensionCodeInfoFfAxWgt                                 string
	TwodimensionCodeInfoFrAxWgt                                 string
	TwodimensionCodeInfoRfAxWgt                                 string
	TwodimensionCodeInfoRrAxWgt                                 string
	TwodimensionCodeInfoNoiseReg                                string
	TwodimensionCodeInfoNearNoiseReg                            string
	TwodimensionCodeInfoDriveMethod                             string
	TwodimensionCodeInfoOpacimeterMeasCar                       string
	TwodimensionCodeInfoNoxPmMeasMode                           string
	TwodimensionCodeInfoNoxValue                                string
	TwodimensionCodeInfoPmValue                                 string
	TwodimensionCodeInfoSafeStdDate                             string
	TwodimensionCodeInfoFuelClassCode                           string
	RegistCarLightCar                                           string
	Created                                                     string
	Modified                                                    string
}

// CarInspectionRepository handles database operations for car_inspection
type CarInspectionRepository struct {
	db DB
}

// NewCarInspectionRepository creates a new repository
func NewCarInspectionRepository(pool *pgxpool.Pool) *CarInspectionRepository {
	return &CarInspectionRepository{db: pool}
}

// NewCarInspectionRepositoryWithDB creates a repository with custom DB interface (for testing)
func NewCarInspectionRepositoryWithDB(db DB) *CarInspectionRepository {
	return &CarInspectionRepository{db: db}
}

// Create inserts a new car inspection
func (r *CarInspectionRepository) Create(ctx context.Context, inspection *CarInspection) (*CarInspection, error) {
	query := `
		INSERT INTO car_inspection (
			organization_id, "CertInfoImportFileVersion", "AcceptOutputNo", "FormType", "ElectCertMgNo", "CarId",
			"ElectCertPublishdateE", "ElectCertPublishdateY", "ElectCertPublishdateM", "ElectCertPublishdateD",
			"GrantdateE", "GrantdateY", "GrantdateM", "GrantdateD", "TranspotationBureauChiefName",
			"EntryNoCarNo", "RegGrantdateE", "RegGrantdateY", "RegGrantdateM", "RegGrantdateD",
			"FirstRegistDateE", "FirstRegistDateY", "FirstRegistDateM",
			"CarName", "CarNameCode", "CarNo", "Model", "EngineModel",
			"OwnerNameLowLevelChar", "OwnerNameHighLevelChar", "OwnerAddressChar", "OwnerAddressNumValue", "OwnerAddressCode",
			"UserNameLowLevelChar", "UserNameHighLevelChar", "UserAddressChar", "UserAddressNumValue", "UserAddressCode",
			"UseHeadqurterChar", "UseHeadqurterNumValue", "UseHeadqurterCode",
			"CarKind", "Use", "PrivateBusiness", "CarShape", "CarShapeCode",
			"NoteCap", "Cap", "NoteMaxLoadage", "MaxLoadage",
			"NoteCarWgt", "CarWgt", "NoteCarTotalWgt", "CarTotalWgt",
			"NoteLength", "Length", "NoteWidth", "Width", "NoteHeight", "Height",
			"FfAxWgt", "FrAxWgt", "RfAxWgt", "RrAxWgt",
			"Displacement", "FuelClass", "ModelSpecifyNo", "ClassifyAroundNo",
			"ValidPeriodExpirDateE", "ValidPeriodExpirDateY", "ValidPeriodExpirDateM", "ValidPeriodExpirDateD",
			"NoteInfo",
			"TwodimensionCodeInfoEntryNoCarNo", "TwodimensionCodeInfoCarNo", "TwodimensionCodeInfoValidPeriodExpirDate",
			"TwodimensionCodeInfoModel", "TwodimensionCodeInfoModelSpecifyNoClassifyAroundNo",
			"TwodimensionCodeInfoCharInfo", "TwodimensionCodeInfoEngineModel", "TwodimensionCodeInfoCarNoStampPlace",
			"TwodimensionCodeInfoFirstRegistDate",
			"TwodimensionCodeInfoFfAxWgt", "TwodimensionCodeInfoFrAxWgt", "TwodimensionCodeInfoRfAxWgt", "TwodimensionCodeInfoRrAxWgt",
			"TwodimensionCodeInfoNoiseReg", "TwodimensionCodeInfoNearNoiseReg",
			"TwodimensionCodeInfoDriveMethod", "TwodimensionCodeInfoOpacimeterMeasCar",
			"TwodimensionCodeInfoNoxPmMeasMode", "TwodimensionCodeInfoNoxValue", "TwodimensionCodeInfoPmValue",
			"TwodimensionCodeInfoSafeStdDate", "TwodimensionCodeInfoFuelClassCode",
			"RegistCarLightCar", "Created", "Modified"
		)
		VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10,
			$11, $12, $13, $14, $15,
			$16, $17, $18, $19, $20,
			$21, $22, $23,
			$24, $25, $26, $27, $28,
			$29, $30, $31, $32, $33,
			$34, $35, $36, $37, $38,
			$39, $40, $41,
			$42, $43, $44, $45, $46,
			$47, $48, $49, $50,
			$51, $52, $53, $54,
			$55, $56, $57, $58, $59, $60,
			$61, $62, $63, $64,
			$65, $66, $67, $68,
			$69, $70, $71, $72,
			$73,
			$74, $75, $76,
			$77, $78,
			$79, $80, $81,
			$82,
			$83, $84, $85, $86,
			$87, $88,
			$89, $90,
			$91, $92, $93,
			$94, $95,
			$96, $97, $98
		)
		RETURNING
			organization_id, "CertInfoImportFileVersion", "AcceptOutputNo", "FormType", "ElectCertMgNo", "CarId",
			"ElectCertPublishdateE", "ElectCertPublishdateY", "ElectCertPublishdateM", "ElectCertPublishdateD",
			"GrantdateE", "GrantdateY", "GrantdateM", "GrantdateD", "TranspotationBureauChiefName",
			"EntryNoCarNo", "RegGrantdateE", "RegGrantdateY", "RegGrantdateM", "RegGrantdateD",
			"FirstRegistDateE", "FirstRegistDateY", "FirstRegistDateM",
			"CarName", "CarNameCode", "CarNo", "Model", "EngineModel",
			"OwnerNameLowLevelChar", "OwnerNameHighLevelChar", "OwnerAddressChar", "OwnerAddressNumValue", "OwnerAddressCode",
			"UserNameLowLevelChar", "UserNameHighLevelChar", "UserAddressChar", "UserAddressNumValue", "UserAddressCode",
			"UseHeadqurterChar", "UseHeadqurterNumValue", "UseHeadqurterCode",
			"CarKind", "Use", "PrivateBusiness", "CarShape", "CarShapeCode",
			"NoteCap", "Cap", "NoteMaxLoadage", "MaxLoadage",
			"NoteCarWgt", "CarWgt", "NoteCarTotalWgt", "CarTotalWgt",
			"NoteLength", "Length", "NoteWidth", "Width", "NoteHeight", "Height",
			"FfAxWgt", "FrAxWgt", "RfAxWgt", "RrAxWgt",
			"Displacement", "FuelClass", "ModelSpecifyNo", "ClassifyAroundNo",
			"ValidPeriodExpirDateE", "ValidPeriodExpirDateY", "ValidPeriodExpirDateM", "ValidPeriodExpirDateD",
			"NoteInfo",
			"TwodimensionCodeInfoEntryNoCarNo", "TwodimensionCodeInfoCarNo", "TwodimensionCodeInfoValidPeriodExpirDate",
			"TwodimensionCodeInfoModel", "TwodimensionCodeInfoModelSpecifyNoClassifyAroundNo",
			"TwodimensionCodeInfoCharInfo", "TwodimensionCodeInfoEngineModel", "TwodimensionCodeInfoCarNoStampPlace",
			"TwodimensionCodeInfoFirstRegistDate",
			"TwodimensionCodeInfoFfAxWgt", "TwodimensionCodeInfoFrAxWgt", "TwodimensionCodeInfoRfAxWgt", "TwodimensionCodeInfoRrAxWgt",
			"TwodimensionCodeInfoNoiseReg", "TwodimensionCodeInfoNearNoiseReg",
			"TwodimensionCodeInfoDriveMethod", "TwodimensionCodeInfoOpacimeterMeasCar",
			"TwodimensionCodeInfoNoxPmMeasMode", "TwodimensionCodeInfoNoxValue", "TwodimensionCodeInfoPmValue",
			"TwodimensionCodeInfoSafeStdDate", "TwodimensionCodeInfoFuelClassCode",
			"RegistCarLightCar", "Created", "Modified"
	`

	var result CarInspection
	err := r.db.QueryRow(ctx, query,
		inspection.OrganizationID, inspection.CertInfoImportFileVersion, inspection.AcceptOutputNo, inspection.FormType, inspection.ElectCertMgNo, inspection.CarID,
		inspection.ElectCertPublishdateE, inspection.ElectCertPublishdateY, inspection.ElectCertPublishdateM, inspection.ElectCertPublishdateD,
		inspection.GrantdateE, inspection.GrantdateY, inspection.GrantdateM, inspection.GrantdateD, inspection.TranspotationBureauChiefName,
		inspection.EntryNoCarNo, inspection.RegGrantdateE, inspection.RegGrantdateY, inspection.RegGrantdateM, inspection.RegGrantdateD,
		inspection.FirstRegistDateE, inspection.FirstRegistDateY, inspection.FirstRegistDateM,
		inspection.CarName, inspection.CarNameCode, inspection.CarNo, inspection.Model, inspection.EngineModel,
		inspection.OwnerNameLowLevelChar, inspection.OwnerNameHighLevelChar, inspection.OwnerAddressChar, inspection.OwnerAddressNumValue, inspection.OwnerAddressCode,
		inspection.UserNameLowLevelChar, inspection.UserNameHighLevelChar, inspection.UserAddressChar, inspection.UserAddressNumValue, inspection.UserAddressCode,
		inspection.UseHeadqurterChar, inspection.UseHeadqurterNumValue, inspection.UseHeadqurterCode,
		inspection.CarKind, inspection.Use, inspection.PrivateBusiness, inspection.CarShape, inspection.CarShapeCode,
		inspection.NoteCap, inspection.Cap, inspection.NoteMaxLoadage, inspection.MaxLoadage,
		inspection.NoteCarWgt, inspection.CarWgt, inspection.NoteCarTotalWgt, inspection.CarTotalWgt,
		inspection.NoteLength, inspection.Length, inspection.NoteWidth, inspection.Width, inspection.NoteHeight, inspection.Height,
		inspection.FfAxWgt, inspection.FrAxWgt, inspection.RfAxWgt, inspection.RrAxWgt,
		inspection.Displacement, inspection.FuelClass, inspection.ModelSpecifyNo, inspection.ClassifyAroundNo,
		inspection.ValidPeriodExpirDateE, inspection.ValidPeriodExpirDateY, inspection.ValidPeriodExpirDateM, inspection.ValidPeriodExpirDateD,
		inspection.NoteInfo,
		inspection.TwodimensionCodeInfoEntryNoCarNo, inspection.TwodimensionCodeInfoCarNo, inspection.TwodimensionCodeInfoValidPeriodExpirDate,
		inspection.TwodimensionCodeInfoModel, inspection.TwodimensionCodeInfoModelSpecifyNoClassifyAroundNo,
		inspection.TwodimensionCodeInfoCharInfo, inspection.TwodimensionCodeInfoEngineModel, inspection.TwodimensionCodeInfoCarNoStampPlace,
		inspection.TwodimensionCodeInfoFirstRegistDate,
		inspection.TwodimensionCodeInfoFfAxWgt, inspection.TwodimensionCodeInfoFrAxWgt, inspection.TwodimensionCodeInfoRfAxWgt, inspection.TwodimensionCodeInfoRrAxWgt,
		inspection.TwodimensionCodeInfoNoiseReg, inspection.TwodimensionCodeInfoNearNoiseReg,
		inspection.TwodimensionCodeInfoDriveMethod, inspection.TwodimensionCodeInfoOpacimeterMeasCar,
		inspection.TwodimensionCodeInfoNoxPmMeasMode, inspection.TwodimensionCodeInfoNoxValue, inspection.TwodimensionCodeInfoPmValue,
		inspection.TwodimensionCodeInfoSafeStdDate, inspection.TwodimensionCodeInfoFuelClassCode,
		inspection.RegistCarLightCar, inspection.Created, inspection.Modified,
	).Scan(
		&result.OrganizationID, &result.CertInfoImportFileVersion, &result.AcceptOutputNo, &result.FormType, &result.ElectCertMgNo, &result.CarID,
		&result.ElectCertPublishdateE, &result.ElectCertPublishdateY, &result.ElectCertPublishdateM, &result.ElectCertPublishdateD,
		&result.GrantdateE, &result.GrantdateY, &result.GrantdateM, &result.GrantdateD, &result.TranspotationBureauChiefName,
		&result.EntryNoCarNo, &result.RegGrantdateE, &result.RegGrantdateY, &result.RegGrantdateM, &result.RegGrantdateD,
		&result.FirstRegistDateE, &result.FirstRegistDateY, &result.FirstRegistDateM,
		&result.CarName, &result.CarNameCode, &result.CarNo, &result.Model, &result.EngineModel,
		&result.OwnerNameLowLevelChar, &result.OwnerNameHighLevelChar, &result.OwnerAddressChar, &result.OwnerAddressNumValue, &result.OwnerAddressCode,
		&result.UserNameLowLevelChar, &result.UserNameHighLevelChar, &result.UserAddressChar, &result.UserAddressNumValue, &result.UserAddressCode,
		&result.UseHeadqurterChar, &result.UseHeadqurterNumValue, &result.UseHeadqurterCode,
		&result.CarKind, &result.Use, &result.PrivateBusiness, &result.CarShape, &result.CarShapeCode,
		&result.NoteCap, &result.Cap, &result.NoteMaxLoadage, &result.MaxLoadage,
		&result.NoteCarWgt, &result.CarWgt, &result.NoteCarTotalWgt, &result.CarTotalWgt,
		&result.NoteLength, &result.Length, &result.NoteWidth, &result.Width, &result.NoteHeight, &result.Height,
		&result.FfAxWgt, &result.FrAxWgt, &result.RfAxWgt, &result.RrAxWgt,
		&result.Displacement, &result.FuelClass, &result.ModelSpecifyNo, &result.ClassifyAroundNo,
		&result.ValidPeriodExpirDateE, &result.ValidPeriodExpirDateY, &result.ValidPeriodExpirDateM, &result.ValidPeriodExpirDateD,
		&result.NoteInfo,
		&result.TwodimensionCodeInfoEntryNoCarNo, &result.TwodimensionCodeInfoCarNo, &result.TwodimensionCodeInfoValidPeriodExpirDate,
		&result.TwodimensionCodeInfoModel, &result.TwodimensionCodeInfoModelSpecifyNoClassifyAroundNo,
		&result.TwodimensionCodeInfoCharInfo, &result.TwodimensionCodeInfoEngineModel, &result.TwodimensionCodeInfoCarNoStampPlace,
		&result.TwodimensionCodeInfoFirstRegistDate,
		&result.TwodimensionCodeInfoFfAxWgt, &result.TwodimensionCodeInfoFrAxWgt, &result.TwodimensionCodeInfoRfAxWgt, &result.TwodimensionCodeInfoRrAxWgt,
		&result.TwodimensionCodeInfoNoiseReg, &result.TwodimensionCodeInfoNearNoiseReg,
		&result.TwodimensionCodeInfoDriveMethod, &result.TwodimensionCodeInfoOpacimeterMeasCar,
		&result.TwodimensionCodeInfoNoxPmMeasMode, &result.TwodimensionCodeInfoNoxValue, &result.TwodimensionCodeInfoPmValue,
		&result.TwodimensionCodeInfoSafeStdDate, &result.TwodimensionCodeInfoFuelClassCode,
		&result.RegistCarLightCar, &result.Created, &result.Modified,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetByPrimaryKey retrieves a car inspection by composite primary key
func (r *CarInspectionRepository) GetByPrimaryKey(ctx context.Context, organizationID, electCertMgNo, electCertPublishdateE, electCertPublishdateY, electCertPublishdateM, electCertPublishdateD string) (*CarInspection, error) {
	query := `
		SELECT
			organization_id, "CertInfoImportFileVersion", "AcceptOutputNo", "FormType", "ElectCertMgNo", "CarId",
			"ElectCertPublishdateE", "ElectCertPublishdateY", "ElectCertPublishdateM", "ElectCertPublishdateD",
			"GrantdateE", "GrantdateY", "GrantdateM", "GrantdateD", "TranspotationBureauChiefName",
			"EntryNoCarNo", "RegGrantdateE", "RegGrantdateY", "RegGrantdateM", "RegGrantdateD",
			"FirstRegistDateE", "FirstRegistDateY", "FirstRegistDateM",
			"CarName", "CarNameCode", "CarNo", "Model", "EngineModel",
			"OwnerNameLowLevelChar", "OwnerNameHighLevelChar", "OwnerAddressChar", "OwnerAddressNumValue", "OwnerAddressCode",
			"UserNameLowLevelChar", "UserNameHighLevelChar", "UserAddressChar", "UserAddressNumValue", "UserAddressCode",
			"UseHeadqurterChar", "UseHeadqurterNumValue", "UseHeadqurterCode",
			"CarKind", "Use", "PrivateBusiness", "CarShape", "CarShapeCode",
			"NoteCap", "Cap", "NoteMaxLoadage", "MaxLoadage",
			"NoteCarWgt", "CarWgt", "NoteCarTotalWgt", "CarTotalWgt",
			"NoteLength", "Length", "NoteWidth", "Width", "NoteHeight", "Height",
			"FfAxWgt", "FrAxWgt", "RfAxWgt", "RrAxWgt",
			"Displacement", "FuelClass", "ModelSpecifyNo", "ClassifyAroundNo",
			"ValidPeriodExpirDateE", "ValidPeriodExpirDateY", "ValidPeriodExpirDateM", "ValidPeriodExpirDateD",
			"NoteInfo",
			"TwodimensionCodeInfoEntryNoCarNo", "TwodimensionCodeInfoCarNo", "TwodimensionCodeInfoValidPeriodExpirDate",
			"TwodimensionCodeInfoModel", "TwodimensionCodeInfoModelSpecifyNoClassifyAroundNo",
			"TwodimensionCodeInfoCharInfo", "TwodimensionCodeInfoEngineModel", "TwodimensionCodeInfoCarNoStampPlace",
			"TwodimensionCodeInfoFirstRegistDate",
			"TwodimensionCodeInfoFfAxWgt", "TwodimensionCodeInfoFrAxWgt", "TwodimensionCodeInfoRfAxWgt", "TwodimensionCodeInfoRrAxWgt",
			"TwodimensionCodeInfoNoiseReg", "TwodimensionCodeInfoNearNoiseReg",
			"TwodimensionCodeInfoDriveMethod", "TwodimensionCodeInfoOpacimeterMeasCar",
			"TwodimensionCodeInfoNoxPmMeasMode", "TwodimensionCodeInfoNoxValue", "TwodimensionCodeInfoPmValue",
			"TwodimensionCodeInfoSafeStdDate", "TwodimensionCodeInfoFuelClassCode",
			"RegistCarLightCar", "Created", "Modified"
		FROM car_inspection
		WHERE organization_id = $1
			AND "ElectCertMgNo" = $2
			AND "ElectCertPublishdateE" = $3
			AND "ElectCertPublishdateY" = $4
			AND "ElectCertPublishdateM" = $5
			AND "ElectCertPublishdateD" = $6
	`

	var inspection CarInspection
	err := r.db.QueryRow(ctx, query, organizationID, electCertMgNo, electCertPublishdateE, electCertPublishdateY, electCertPublishdateM, electCertPublishdateD).Scan(
		&inspection.OrganizationID, &inspection.CertInfoImportFileVersion, &inspection.AcceptOutputNo, &inspection.FormType, &inspection.ElectCertMgNo, &inspection.CarID,
		&inspection.ElectCertPublishdateE, &inspection.ElectCertPublishdateY, &inspection.ElectCertPublishdateM, &inspection.ElectCertPublishdateD,
		&inspection.GrantdateE, &inspection.GrantdateY, &inspection.GrantdateM, &inspection.GrantdateD, &inspection.TranspotationBureauChiefName,
		&inspection.EntryNoCarNo, &inspection.RegGrantdateE, &inspection.RegGrantdateY, &inspection.RegGrantdateM, &inspection.RegGrantdateD,
		&inspection.FirstRegistDateE, &inspection.FirstRegistDateY, &inspection.FirstRegistDateM,
		&inspection.CarName, &inspection.CarNameCode, &inspection.CarNo, &inspection.Model, &inspection.EngineModel,
		&inspection.OwnerNameLowLevelChar, &inspection.OwnerNameHighLevelChar, &inspection.OwnerAddressChar, &inspection.OwnerAddressNumValue, &inspection.OwnerAddressCode,
		&inspection.UserNameLowLevelChar, &inspection.UserNameHighLevelChar, &inspection.UserAddressChar, &inspection.UserAddressNumValue, &inspection.UserAddressCode,
		&inspection.UseHeadqurterChar, &inspection.UseHeadqurterNumValue, &inspection.UseHeadqurterCode,
		&inspection.CarKind, &inspection.Use, &inspection.PrivateBusiness, &inspection.CarShape, &inspection.CarShapeCode,
		&inspection.NoteCap, &inspection.Cap, &inspection.NoteMaxLoadage, &inspection.MaxLoadage,
		&inspection.NoteCarWgt, &inspection.CarWgt, &inspection.NoteCarTotalWgt, &inspection.CarTotalWgt,
		&inspection.NoteLength, &inspection.Length, &inspection.NoteWidth, &inspection.Width, &inspection.NoteHeight, &inspection.Height,
		&inspection.FfAxWgt, &inspection.FrAxWgt, &inspection.RfAxWgt, &inspection.RrAxWgt,
		&inspection.Displacement, &inspection.FuelClass, &inspection.ModelSpecifyNo, &inspection.ClassifyAroundNo,
		&inspection.ValidPeriodExpirDateE, &inspection.ValidPeriodExpirDateY, &inspection.ValidPeriodExpirDateM, &inspection.ValidPeriodExpirDateD,
		&inspection.NoteInfo,
		&inspection.TwodimensionCodeInfoEntryNoCarNo, &inspection.TwodimensionCodeInfoCarNo, &inspection.TwodimensionCodeInfoValidPeriodExpirDate,
		&inspection.TwodimensionCodeInfoModel, &inspection.TwodimensionCodeInfoModelSpecifyNoClassifyAroundNo,
		&inspection.TwodimensionCodeInfoCharInfo, &inspection.TwodimensionCodeInfoEngineModel, &inspection.TwodimensionCodeInfoCarNoStampPlace,
		&inspection.TwodimensionCodeInfoFirstRegistDate,
		&inspection.TwodimensionCodeInfoFfAxWgt, &inspection.TwodimensionCodeInfoFrAxWgt, &inspection.TwodimensionCodeInfoRfAxWgt, &inspection.TwodimensionCodeInfoRrAxWgt,
		&inspection.TwodimensionCodeInfoNoiseReg, &inspection.TwodimensionCodeInfoNearNoiseReg,
		&inspection.TwodimensionCodeInfoDriveMethod, &inspection.TwodimensionCodeInfoOpacimeterMeasCar,
		&inspection.TwodimensionCodeInfoNoxPmMeasMode, &inspection.TwodimensionCodeInfoNoxValue, &inspection.TwodimensionCodeInfoPmValue,
		&inspection.TwodimensionCodeInfoSafeStdDate, &inspection.TwodimensionCodeInfoFuelClassCode,
		&inspection.RegistCarLightCar, &inspection.Created, &inspection.Modified,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCarInspectionNotFound
		}
		return nil, err
	}

	return &inspection, nil
}

// Update modifies an existing car inspection
func (r *CarInspectionRepository) Update(ctx context.Context, inspection *CarInspection) (*CarInspection, error) {
	query := `
		UPDATE car_inspection
		SET
			"CertInfoImportFileVersion" = $7, "AcceptOutputNo" = $8, "FormType" = $9, "CarId" = $10,
			"GrantdateE" = $11, "GrantdateY" = $12, "GrantdateM" = $13, "GrantdateD" = $14, "TranspotationBureauChiefName" = $15,
			"EntryNoCarNo" = $16, "RegGrantdateE" = $17, "RegGrantdateY" = $18, "RegGrantdateM" = $19, "RegGrantdateD" = $20,
			"FirstRegistDateE" = $21, "FirstRegistDateY" = $22, "FirstRegistDateM" = $23,
			"CarName" = $24, "CarNameCode" = $25, "CarNo" = $26, "Model" = $27, "EngineModel" = $28,
			"OwnerNameLowLevelChar" = $29, "OwnerNameHighLevelChar" = $30, "OwnerAddressChar" = $31, "OwnerAddressNumValue" = $32, "OwnerAddressCode" = $33,
			"UserNameLowLevelChar" = $34, "UserNameHighLevelChar" = $35, "UserAddressChar" = $36, "UserAddressNumValue" = $37, "UserAddressCode" = $38,
			"UseHeadqurterChar" = $39, "UseHeadqurterNumValue" = $40, "UseHeadqurterCode" = $41,
			"CarKind" = $42, "Use" = $43, "PrivateBusiness" = $44, "CarShape" = $45, "CarShapeCode" = $46,
			"NoteCap" = $47, "Cap" = $48, "NoteMaxLoadage" = $49, "MaxLoadage" = $50,
			"NoteCarWgt" = $51, "CarWgt" = $52, "NoteCarTotalWgt" = $53, "CarTotalWgt" = $54,
			"NoteLength" = $55, "Length" = $56, "NoteWidth" = $57, "Width" = $58, "NoteHeight" = $59, "Height" = $60,
			"FfAxWgt" = $61, "FrAxWgt" = $62, "RfAxWgt" = $63, "RrAxWgt" = $64,
			"Displacement" = $65, "FuelClass" = $66, "ModelSpecifyNo" = $67, "ClassifyAroundNo" = $68,
			"ValidPeriodExpirDateE" = $69, "ValidPeriodExpirDateY" = $70, "ValidPeriodExpirDateM" = $71, "ValidPeriodExpirDateD" = $72,
			"NoteInfo" = $73,
			"TwodimensionCodeInfoEntryNoCarNo" = $74, "TwodimensionCodeInfoCarNo" = $75, "TwodimensionCodeInfoValidPeriodExpirDate" = $76,
			"TwodimensionCodeInfoModel" = $77, "TwodimensionCodeInfoModelSpecifyNoClassifyAroundNo" = $78,
			"TwodimensionCodeInfoCharInfo" = $79, "TwodimensionCodeInfoEngineModel" = $80, "TwodimensionCodeInfoCarNoStampPlace" = $81,
			"TwodimensionCodeInfoFirstRegistDate" = $82,
			"TwodimensionCodeInfoFfAxWgt" = $83, "TwodimensionCodeInfoFrAxWgt" = $84, "TwodimensionCodeInfoRfAxWgt" = $85, "TwodimensionCodeInfoRrAxWgt" = $86,
			"TwodimensionCodeInfoNoiseReg" = $87, "TwodimensionCodeInfoNearNoiseReg" = $88,
			"TwodimensionCodeInfoDriveMethod" = $89, "TwodimensionCodeInfoOpacimeterMeasCar" = $90,
			"TwodimensionCodeInfoNoxPmMeasMode" = $91, "TwodimensionCodeInfoNoxValue" = $92, "TwodimensionCodeInfoPmValue" = $93,
			"TwodimensionCodeInfoSafeStdDate" = $94, "TwodimensionCodeInfoFuelClassCode" = $95,
			"RegistCarLightCar" = $96, "Created" = $97, "Modified" = $98
		WHERE organization_id = $1
			AND "ElectCertMgNo" = $2
			AND "ElectCertPublishdateE" = $3
			AND "ElectCertPublishdateY" = $4
			AND "ElectCertPublishdateM" = $5
			AND "ElectCertPublishdateD" = $6
		RETURNING
			organization_id, "CertInfoImportFileVersion", "AcceptOutputNo", "FormType", "ElectCertMgNo", "CarId",
			"ElectCertPublishdateE", "ElectCertPublishdateY", "ElectCertPublishdateM", "ElectCertPublishdateD",
			"GrantdateE", "GrantdateY", "GrantdateM", "GrantdateD", "TranspotationBureauChiefName",
			"EntryNoCarNo", "RegGrantdateE", "RegGrantdateY", "RegGrantdateM", "RegGrantdateD",
			"FirstRegistDateE", "FirstRegistDateY", "FirstRegistDateM",
			"CarName", "CarNameCode", "CarNo", "Model", "EngineModel",
			"OwnerNameLowLevelChar", "OwnerNameHighLevelChar", "OwnerAddressChar", "OwnerAddressNumValue", "OwnerAddressCode",
			"UserNameLowLevelChar", "UserNameHighLevelChar", "UserAddressChar", "UserAddressNumValue", "UserAddressCode",
			"UseHeadqurterChar", "UseHeadqurterNumValue", "UseHeadqurterCode",
			"CarKind", "Use", "PrivateBusiness", "CarShape", "CarShapeCode",
			"NoteCap", "Cap", "NoteMaxLoadage", "MaxLoadage",
			"NoteCarWgt", "CarWgt", "NoteCarTotalWgt", "CarTotalWgt",
			"NoteLength", "Length", "NoteWidth", "Width", "NoteHeight", "Height",
			"FfAxWgt", "FrAxWgt", "RfAxWgt", "RrAxWgt",
			"Displacement", "FuelClass", "ModelSpecifyNo", "ClassifyAroundNo",
			"ValidPeriodExpirDateE", "ValidPeriodExpirDateY", "ValidPeriodExpirDateM", "ValidPeriodExpirDateD",
			"NoteInfo",
			"TwodimensionCodeInfoEntryNoCarNo", "TwodimensionCodeInfoCarNo", "TwodimensionCodeInfoValidPeriodExpirDate",
			"TwodimensionCodeInfoModel", "TwodimensionCodeInfoModelSpecifyNoClassifyAroundNo",
			"TwodimensionCodeInfoCharInfo", "TwodimensionCodeInfoEngineModel", "TwodimensionCodeInfoCarNoStampPlace",
			"TwodimensionCodeInfoFirstRegistDate",
			"TwodimensionCodeInfoFfAxWgt", "TwodimensionCodeInfoFrAxWgt", "TwodimensionCodeInfoRfAxWgt", "TwodimensionCodeInfoRrAxWgt",
			"TwodimensionCodeInfoNoiseReg", "TwodimensionCodeInfoNearNoiseReg",
			"TwodimensionCodeInfoDriveMethod", "TwodimensionCodeInfoOpacimeterMeasCar",
			"TwodimensionCodeInfoNoxPmMeasMode", "TwodimensionCodeInfoNoxValue", "TwodimensionCodeInfoPmValue",
			"TwodimensionCodeInfoSafeStdDate", "TwodimensionCodeInfoFuelClassCode",
			"RegistCarLightCar", "Created", "Modified"
	`

	var result CarInspection
	err := r.db.QueryRow(ctx, query,
		inspection.OrganizationID, inspection.ElectCertMgNo, inspection.ElectCertPublishdateE, inspection.ElectCertPublishdateY, inspection.ElectCertPublishdateM, inspection.ElectCertPublishdateD,
		inspection.CertInfoImportFileVersion, inspection.AcceptOutputNo, inspection.FormType, inspection.CarID,
		inspection.GrantdateE, inspection.GrantdateY, inspection.GrantdateM, inspection.GrantdateD, inspection.TranspotationBureauChiefName,
		inspection.EntryNoCarNo, inspection.RegGrantdateE, inspection.RegGrantdateY, inspection.RegGrantdateM, inspection.RegGrantdateD,
		inspection.FirstRegistDateE, inspection.FirstRegistDateY, inspection.FirstRegistDateM,
		inspection.CarName, inspection.CarNameCode, inspection.CarNo, inspection.Model, inspection.EngineModel,
		inspection.OwnerNameLowLevelChar, inspection.OwnerNameHighLevelChar, inspection.OwnerAddressChar, inspection.OwnerAddressNumValue, inspection.OwnerAddressCode,
		inspection.UserNameLowLevelChar, inspection.UserNameHighLevelChar, inspection.UserAddressChar, inspection.UserAddressNumValue, inspection.UserAddressCode,
		inspection.UseHeadqurterChar, inspection.UseHeadqurterNumValue, inspection.UseHeadqurterCode,
		inspection.CarKind, inspection.Use, inspection.PrivateBusiness, inspection.CarShape, inspection.CarShapeCode,
		inspection.NoteCap, inspection.Cap, inspection.NoteMaxLoadage, inspection.MaxLoadage,
		inspection.NoteCarWgt, inspection.CarWgt, inspection.NoteCarTotalWgt, inspection.CarTotalWgt,
		inspection.NoteLength, inspection.Length, inspection.NoteWidth, inspection.Width, inspection.NoteHeight, inspection.Height,
		inspection.FfAxWgt, inspection.FrAxWgt, inspection.RfAxWgt, inspection.RrAxWgt,
		inspection.Displacement, inspection.FuelClass, inspection.ModelSpecifyNo, inspection.ClassifyAroundNo,
		inspection.ValidPeriodExpirDateE, inspection.ValidPeriodExpirDateY, inspection.ValidPeriodExpirDateM, inspection.ValidPeriodExpirDateD,
		inspection.NoteInfo,
		inspection.TwodimensionCodeInfoEntryNoCarNo, inspection.TwodimensionCodeInfoCarNo, inspection.TwodimensionCodeInfoValidPeriodExpirDate,
		inspection.TwodimensionCodeInfoModel, inspection.TwodimensionCodeInfoModelSpecifyNoClassifyAroundNo,
		inspection.TwodimensionCodeInfoCharInfo, inspection.TwodimensionCodeInfoEngineModel, inspection.TwodimensionCodeInfoCarNoStampPlace,
		inspection.TwodimensionCodeInfoFirstRegistDate,
		inspection.TwodimensionCodeInfoFfAxWgt, inspection.TwodimensionCodeInfoFrAxWgt, inspection.TwodimensionCodeInfoRfAxWgt, inspection.TwodimensionCodeInfoRrAxWgt,
		inspection.TwodimensionCodeInfoNoiseReg, inspection.TwodimensionCodeInfoNearNoiseReg,
		inspection.TwodimensionCodeInfoDriveMethod, inspection.TwodimensionCodeInfoOpacimeterMeasCar,
		inspection.TwodimensionCodeInfoNoxPmMeasMode, inspection.TwodimensionCodeInfoNoxValue, inspection.TwodimensionCodeInfoPmValue,
		inspection.TwodimensionCodeInfoSafeStdDate, inspection.TwodimensionCodeInfoFuelClassCode,
		inspection.RegistCarLightCar, inspection.Created, inspection.Modified,
	).Scan(
		&result.OrganizationID, &result.CertInfoImportFileVersion, &result.AcceptOutputNo, &result.FormType, &result.ElectCertMgNo, &result.CarID,
		&result.ElectCertPublishdateE, &result.ElectCertPublishdateY, &result.ElectCertPublishdateM, &result.ElectCertPublishdateD,
		&result.GrantdateE, &result.GrantdateY, &result.GrantdateM, &result.GrantdateD, &result.TranspotationBureauChiefName,
		&result.EntryNoCarNo, &result.RegGrantdateE, &result.RegGrantdateY, &result.RegGrantdateM, &result.RegGrantdateD,
		&result.FirstRegistDateE, &result.FirstRegistDateY, &result.FirstRegistDateM,
		&result.CarName, &result.CarNameCode, &result.CarNo, &result.Model, &result.EngineModel,
		&result.OwnerNameLowLevelChar, &result.OwnerNameHighLevelChar, &result.OwnerAddressChar, &result.OwnerAddressNumValue, &result.OwnerAddressCode,
		&result.UserNameLowLevelChar, &result.UserNameHighLevelChar, &result.UserAddressChar, &result.UserAddressNumValue, &result.UserAddressCode,
		&result.UseHeadqurterChar, &result.UseHeadqurterNumValue, &result.UseHeadqurterCode,
		&result.CarKind, &result.Use, &result.PrivateBusiness, &result.CarShape, &result.CarShapeCode,
		&result.NoteCap, &result.Cap, &result.NoteMaxLoadage, &result.MaxLoadage,
		&result.NoteCarWgt, &result.CarWgt, &result.NoteCarTotalWgt, &result.CarTotalWgt,
		&result.NoteLength, &result.Length, &result.NoteWidth, &result.Width, &result.NoteHeight, &result.Height,
		&result.FfAxWgt, &result.FrAxWgt, &result.RfAxWgt, &result.RrAxWgt,
		&result.Displacement, &result.FuelClass, &result.ModelSpecifyNo, &result.ClassifyAroundNo,
		&result.ValidPeriodExpirDateE, &result.ValidPeriodExpirDateY, &result.ValidPeriodExpirDateM, &result.ValidPeriodExpirDateD,
		&result.NoteInfo,
		&result.TwodimensionCodeInfoEntryNoCarNo, &result.TwodimensionCodeInfoCarNo, &result.TwodimensionCodeInfoValidPeriodExpirDate,
		&result.TwodimensionCodeInfoModel, &result.TwodimensionCodeInfoModelSpecifyNoClassifyAroundNo,
		&result.TwodimensionCodeInfoCharInfo, &result.TwodimensionCodeInfoEngineModel, &result.TwodimensionCodeInfoCarNoStampPlace,
		&result.TwodimensionCodeInfoFirstRegistDate,
		&result.TwodimensionCodeInfoFfAxWgt, &result.TwodimensionCodeInfoFrAxWgt, &result.TwodimensionCodeInfoRfAxWgt, &result.TwodimensionCodeInfoRrAxWgt,
		&result.TwodimensionCodeInfoNoiseReg, &result.TwodimensionCodeInfoNearNoiseReg,
		&result.TwodimensionCodeInfoDriveMethod, &result.TwodimensionCodeInfoOpacimeterMeasCar,
		&result.TwodimensionCodeInfoNoxPmMeasMode, &result.TwodimensionCodeInfoNoxValue, &result.TwodimensionCodeInfoPmValue,
		&result.TwodimensionCodeInfoSafeStdDate, &result.TwodimensionCodeInfoFuelClassCode,
		&result.RegistCarLightCar, &result.Created, &result.Modified,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCarInspectionNotFound
		}
		return nil, err
	}

	return &result, nil
}

// Delete hard-deletes a car inspection
func (r *CarInspectionRepository) Delete(ctx context.Context, organizationID, electCertMgNo, electCertPublishdateE, electCertPublishdateY, electCertPublishdateM, electCertPublishdateD string) error {
	query := `
		DELETE FROM car_inspection
		WHERE organization_id = $1
			AND "ElectCertMgNo" = $2
			AND "ElectCertPublishdateE" = $3
			AND "ElectCertPublishdateY" = $4
			AND "ElectCertPublishdateM" = $5
			AND "ElectCertPublishdateD" = $6
	`

	result, err := r.db.Exec(ctx, query, organizationID, electCertMgNo, electCertPublishdateE, electCertPublishdateY, electCertPublishdateM, electCertPublishdateD)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrCarInspectionNotFound
	}

	return nil
}

// ListByOrganization retrieves car inspections by organization with pagination
func (r *CarInspectionRepository) ListByOrganization(ctx context.Context, organizationID string, limit int, offset int) ([]*CarInspection, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	query := `
		SELECT
			organization_id, "CertInfoImportFileVersion", "AcceptOutputNo", "FormType", "ElectCertMgNo", "CarId",
			"ElectCertPublishdateE", "ElectCertPublishdateY", "ElectCertPublishdateM", "ElectCertPublishdateD",
			"GrantdateE", "GrantdateY", "GrantdateM", "GrantdateD", "TranspotationBureauChiefName",
			"EntryNoCarNo", "RegGrantdateE", "RegGrantdateY", "RegGrantdateM", "RegGrantdateD",
			"FirstRegistDateE", "FirstRegistDateY", "FirstRegistDateM",
			"CarName", "CarNameCode", "CarNo", "Model", "EngineModel",
			"OwnerNameLowLevelChar", "OwnerNameHighLevelChar", "OwnerAddressChar", "OwnerAddressNumValue", "OwnerAddressCode",
			"UserNameLowLevelChar", "UserNameHighLevelChar", "UserAddressChar", "UserAddressNumValue", "UserAddressCode",
			"UseHeadqurterChar", "UseHeadqurterNumValue", "UseHeadqurterCode",
			"CarKind", "Use", "PrivateBusiness", "CarShape", "CarShapeCode",
			"NoteCap", "Cap", "NoteMaxLoadage", "MaxLoadage",
			"NoteCarWgt", "CarWgt", "NoteCarTotalWgt", "CarTotalWgt",
			"NoteLength", "Length", "NoteWidth", "Width", "NoteHeight", "Height",
			"FfAxWgt", "FrAxWgt", "RfAxWgt", "RrAxWgt",
			"Displacement", "FuelClass", "ModelSpecifyNo", "ClassifyAroundNo",
			"ValidPeriodExpirDateE", "ValidPeriodExpirDateY", "ValidPeriodExpirDateM", "ValidPeriodExpirDateD",
			"NoteInfo",
			"TwodimensionCodeInfoEntryNoCarNo", "TwodimensionCodeInfoCarNo", "TwodimensionCodeInfoValidPeriodExpirDate",
			"TwodimensionCodeInfoModel", "TwodimensionCodeInfoModelSpecifyNoClassifyAroundNo",
			"TwodimensionCodeInfoCharInfo", "TwodimensionCodeInfoEngineModel", "TwodimensionCodeInfoCarNoStampPlace",
			"TwodimensionCodeInfoFirstRegistDate",
			"TwodimensionCodeInfoFfAxWgt", "TwodimensionCodeInfoFrAxWgt", "TwodimensionCodeInfoRfAxWgt", "TwodimensionCodeInfoRrAxWgt",
			"TwodimensionCodeInfoNoiseReg", "TwodimensionCodeInfoNearNoiseReg",
			"TwodimensionCodeInfoDriveMethod", "TwodimensionCodeInfoOpacimeterMeasCar",
			"TwodimensionCodeInfoNoxPmMeasMode", "TwodimensionCodeInfoNoxValue", "TwodimensionCodeInfoPmValue",
			"TwodimensionCodeInfoSafeStdDate", "TwodimensionCodeInfoFuelClassCode",
			"RegistCarLightCar", "Created", "Modified"
		FROM car_inspection
		WHERE organization_id = $1
		ORDER BY "ElectCertMgNo", "ElectCertPublishdateE", "ElectCertPublishdateY", "ElectCertPublishdateM", "ElectCertPublishdateD"
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, organizationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inspections []*CarInspection
	for rows.Next() {
		var inspection CarInspection
		err := rows.Scan(
			&inspection.OrganizationID, &inspection.CertInfoImportFileVersion, &inspection.AcceptOutputNo, &inspection.FormType, &inspection.ElectCertMgNo, &inspection.CarID,
			&inspection.ElectCertPublishdateE, &inspection.ElectCertPublishdateY, &inspection.ElectCertPublishdateM, &inspection.ElectCertPublishdateD,
			&inspection.GrantdateE, &inspection.GrantdateY, &inspection.GrantdateM, &inspection.GrantdateD, &inspection.TranspotationBureauChiefName,
			&inspection.EntryNoCarNo, &inspection.RegGrantdateE, &inspection.RegGrantdateY, &inspection.RegGrantdateM, &inspection.RegGrantdateD,
			&inspection.FirstRegistDateE, &inspection.FirstRegistDateY, &inspection.FirstRegistDateM,
			&inspection.CarName, &inspection.CarNameCode, &inspection.CarNo, &inspection.Model, &inspection.EngineModel,
			&inspection.OwnerNameLowLevelChar, &inspection.OwnerNameHighLevelChar, &inspection.OwnerAddressChar, &inspection.OwnerAddressNumValue, &inspection.OwnerAddressCode,
			&inspection.UserNameLowLevelChar, &inspection.UserNameHighLevelChar, &inspection.UserAddressChar, &inspection.UserAddressNumValue, &inspection.UserAddressCode,
			&inspection.UseHeadqurterChar, &inspection.UseHeadqurterNumValue, &inspection.UseHeadqurterCode,
			&inspection.CarKind, &inspection.Use, &inspection.PrivateBusiness, &inspection.CarShape, &inspection.CarShapeCode,
			&inspection.NoteCap, &inspection.Cap, &inspection.NoteMaxLoadage, &inspection.MaxLoadage,
			&inspection.NoteCarWgt, &inspection.CarWgt, &inspection.NoteCarTotalWgt, &inspection.CarTotalWgt,
			&inspection.NoteLength, &inspection.Length, &inspection.NoteWidth, &inspection.Width, &inspection.NoteHeight, &inspection.Height,
			&inspection.FfAxWgt, &inspection.FrAxWgt, &inspection.RfAxWgt, &inspection.RrAxWgt,
			&inspection.Displacement, &inspection.FuelClass, &inspection.ModelSpecifyNo, &inspection.ClassifyAroundNo,
			&inspection.ValidPeriodExpirDateE, &inspection.ValidPeriodExpirDateY, &inspection.ValidPeriodExpirDateM, &inspection.ValidPeriodExpirDateD,
			&inspection.NoteInfo,
			&inspection.TwodimensionCodeInfoEntryNoCarNo, &inspection.TwodimensionCodeInfoCarNo, &inspection.TwodimensionCodeInfoValidPeriodExpirDate,
			&inspection.TwodimensionCodeInfoModel, &inspection.TwodimensionCodeInfoModelSpecifyNoClassifyAroundNo,
			&inspection.TwodimensionCodeInfoCharInfo, &inspection.TwodimensionCodeInfoEngineModel, &inspection.TwodimensionCodeInfoCarNoStampPlace,
			&inspection.TwodimensionCodeInfoFirstRegistDate,
			&inspection.TwodimensionCodeInfoFfAxWgt, &inspection.TwodimensionCodeInfoFrAxWgt, &inspection.TwodimensionCodeInfoRfAxWgt, &inspection.TwodimensionCodeInfoRrAxWgt,
			&inspection.TwodimensionCodeInfoNoiseReg, &inspection.TwodimensionCodeInfoNearNoiseReg,
			&inspection.TwodimensionCodeInfoDriveMethod, &inspection.TwodimensionCodeInfoOpacimeterMeasCar,
			&inspection.TwodimensionCodeInfoNoxPmMeasMode, &inspection.TwodimensionCodeInfoNoxValue, &inspection.TwodimensionCodeInfoPmValue,
			&inspection.TwodimensionCodeInfoSafeStdDate, &inspection.TwodimensionCodeInfoFuelClassCode,
			&inspection.RegistCarLightCar, &inspection.Created, &inspection.Modified,
		)
		if err != nil {
			return nil, err
		}
		inspections = append(inspections, &inspection)
	}

	return inspections, rows.Err()
}

// List retrieves all car inspections with pagination
func (r *CarInspectionRepository) List(ctx context.Context, limit int, offset int) ([]*CarInspection, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	query := `
		SELECT
			organization_id, "CertInfoImportFileVersion", "AcceptOutputNo", "FormType", "ElectCertMgNo", "CarId",
			"ElectCertPublishdateE", "ElectCertPublishdateY", "ElectCertPublishdateM", "ElectCertPublishdateD",
			"GrantdateE", "GrantdateY", "GrantdateM", "GrantdateD", "TranspotationBureauChiefName",
			"EntryNoCarNo", "RegGrantdateE", "RegGrantdateY", "RegGrantdateM", "RegGrantdateD",
			"FirstRegistDateE", "FirstRegistDateY", "FirstRegistDateM",
			"CarName", "CarNameCode", "CarNo", "Model", "EngineModel",
			"OwnerNameLowLevelChar", "OwnerNameHighLevelChar", "OwnerAddressChar", "OwnerAddressNumValue", "OwnerAddressCode",
			"UserNameLowLevelChar", "UserNameHighLevelChar", "UserAddressChar", "UserAddressNumValue", "UserAddressCode",
			"UseHeadqurterChar", "UseHeadqurterNumValue", "UseHeadqurterCode",
			"CarKind", "Use", "PrivateBusiness", "CarShape", "CarShapeCode",
			"NoteCap", "Cap", "NoteMaxLoadage", "MaxLoadage",
			"NoteCarWgt", "CarWgt", "NoteCarTotalWgt", "CarTotalWgt",
			"NoteLength", "Length", "NoteWidth", "Width", "NoteHeight", "Height",
			"FfAxWgt", "FrAxWgt", "RfAxWgt", "RrAxWgt",
			"Displacement", "FuelClass", "ModelSpecifyNo", "ClassifyAroundNo",
			"ValidPeriodExpirDateE", "ValidPeriodExpirDateY", "ValidPeriodExpirDateM", "ValidPeriodExpirDateD",
			"NoteInfo",
			"TwodimensionCodeInfoEntryNoCarNo", "TwodimensionCodeInfoCarNo", "TwodimensionCodeInfoValidPeriodExpirDate",
			"TwodimensionCodeInfoModel", "TwodimensionCodeInfoModelSpecifyNoClassifyAroundNo",
			"TwodimensionCodeInfoCharInfo", "TwodimensionCodeInfoEngineModel", "TwodimensionCodeInfoCarNoStampPlace",
			"TwodimensionCodeInfoFirstRegistDate",
			"TwodimensionCodeInfoFfAxWgt", "TwodimensionCodeInfoFrAxWgt", "TwodimensionCodeInfoRfAxWgt", "TwodimensionCodeInfoRrAxWgt",
			"TwodimensionCodeInfoNoiseReg", "TwodimensionCodeInfoNearNoiseReg",
			"TwodimensionCodeInfoDriveMethod", "TwodimensionCodeInfoOpacimeterMeasCar",
			"TwodimensionCodeInfoNoxPmMeasMode", "TwodimensionCodeInfoNoxValue", "TwodimensionCodeInfoPmValue",
			"TwodimensionCodeInfoSafeStdDate", "TwodimensionCodeInfoFuelClassCode",
			"RegistCarLightCar", "Created", "Modified"
		FROM car_inspection
		ORDER BY organization_id, "ElectCertMgNo", "ElectCertPublishdateE", "ElectCertPublishdateY", "ElectCertPublishdateM", "ElectCertPublishdateD"
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inspections []*CarInspection
	for rows.Next() {
		var inspection CarInspection
		err := rows.Scan(
			&inspection.OrganizationID, &inspection.CertInfoImportFileVersion, &inspection.AcceptOutputNo, &inspection.FormType, &inspection.ElectCertMgNo, &inspection.CarID,
			&inspection.ElectCertPublishdateE, &inspection.ElectCertPublishdateY, &inspection.ElectCertPublishdateM, &inspection.ElectCertPublishdateD,
			&inspection.GrantdateE, &inspection.GrantdateY, &inspection.GrantdateM, &inspection.GrantdateD, &inspection.TranspotationBureauChiefName,
			&inspection.EntryNoCarNo, &inspection.RegGrantdateE, &inspection.RegGrantdateY, &inspection.RegGrantdateM, &inspection.RegGrantdateD,
			&inspection.FirstRegistDateE, &inspection.FirstRegistDateY, &inspection.FirstRegistDateM,
			&inspection.CarName, &inspection.CarNameCode, &inspection.CarNo, &inspection.Model, &inspection.EngineModel,
			&inspection.OwnerNameLowLevelChar, &inspection.OwnerNameHighLevelChar, &inspection.OwnerAddressChar, &inspection.OwnerAddressNumValue, &inspection.OwnerAddressCode,
			&inspection.UserNameLowLevelChar, &inspection.UserNameHighLevelChar, &inspection.UserAddressChar, &inspection.UserAddressNumValue, &inspection.UserAddressCode,
			&inspection.UseHeadqurterChar, &inspection.UseHeadqurterNumValue, &inspection.UseHeadqurterCode,
			&inspection.CarKind, &inspection.Use, &inspection.PrivateBusiness, &inspection.CarShape, &inspection.CarShapeCode,
			&inspection.NoteCap, &inspection.Cap, &inspection.NoteMaxLoadage, &inspection.MaxLoadage,
			&inspection.NoteCarWgt, &inspection.CarWgt, &inspection.NoteCarTotalWgt, &inspection.CarTotalWgt,
			&inspection.NoteLength, &inspection.Length, &inspection.NoteWidth, &inspection.Width, &inspection.NoteHeight, &inspection.Height,
			&inspection.FfAxWgt, &inspection.FrAxWgt, &inspection.RfAxWgt, &inspection.RrAxWgt,
			&inspection.Displacement, &inspection.FuelClass, &inspection.ModelSpecifyNo, &inspection.ClassifyAroundNo,
			&inspection.ValidPeriodExpirDateE, &inspection.ValidPeriodExpirDateY, &inspection.ValidPeriodExpirDateM, &inspection.ValidPeriodExpirDateD,
			&inspection.NoteInfo,
			&inspection.TwodimensionCodeInfoEntryNoCarNo, &inspection.TwodimensionCodeInfoCarNo, &inspection.TwodimensionCodeInfoValidPeriodExpirDate,
			&inspection.TwodimensionCodeInfoModel, &inspection.TwodimensionCodeInfoModelSpecifyNoClassifyAroundNo,
			&inspection.TwodimensionCodeInfoCharInfo, &inspection.TwodimensionCodeInfoEngineModel, &inspection.TwodimensionCodeInfoCarNoStampPlace,
			&inspection.TwodimensionCodeInfoFirstRegistDate,
			&inspection.TwodimensionCodeInfoFfAxWgt, &inspection.TwodimensionCodeInfoFrAxWgt, &inspection.TwodimensionCodeInfoRfAxWgt, &inspection.TwodimensionCodeInfoRrAxWgt,
			&inspection.TwodimensionCodeInfoNoiseReg, &inspection.TwodimensionCodeInfoNearNoiseReg,
			&inspection.TwodimensionCodeInfoDriveMethod, &inspection.TwodimensionCodeInfoOpacimeterMeasCar,
			&inspection.TwodimensionCodeInfoNoxPmMeasMode, &inspection.TwodimensionCodeInfoNoxValue, &inspection.TwodimensionCodeInfoPmValue,
			&inspection.TwodimensionCodeInfoSafeStdDate, &inspection.TwodimensionCodeInfoFuelClassCode,
			&inspection.RegistCarLightCar, &inspection.Created, &inspection.Modified,
		)
		if err != nil {
			return nil, err
		}
		inspections = append(inspections, &inspection)
	}

	return inspections, rows.Err()
}
