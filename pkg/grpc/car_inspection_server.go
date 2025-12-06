package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/pb"
	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/repository"
)

// CarInspectionServer implements the gRPC CarInspectionService
type CarInspectionServer struct {
	pb.UnimplementedCarInspectionServiceServer
	repo *repository.CarInspectionRepository
}

// NewCarInspectionServer creates a new gRPC server
func NewCarInspectionServer(repo *repository.CarInspectionRepository) *CarInspectionServer {
	return &CarInspectionServer{repo: repo}
}

// CreateCarInspection creates a new car inspection record
func (s *CarInspectionServer) CreateCarInspection(ctx context.Context, req *pb.CreateCarInspectionRequest) (*pb.CreateCarInspectionResponse, error) {
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}
	if req.ElectCertMgNo == "" {
		return nil, status.Error(codes.InvalidArgument, "elect_cert_mg_no is required")
	}

	inspection := &repository.CarInspection{
		OrganizationID:                                              req.OrganizationId,
		CertInfoImportFileVersion:                                   req.CertInfoImportFileVersion,
		AcceptOutputNo:                                              req.AcceptOutputNo,
		FormType:                                                    req.FormType,
		ElectCertMgNo:                                               req.ElectCertMgNo,
		CarID:                                                       req.CarId,
		ElectCertPublishdateE:                                       req.ElectCertPublishdateE,
		ElectCertPublishdateY:                                       req.ElectCertPublishdateY,
		ElectCertPublishdateM:                                       req.ElectCertPublishdateM,
		ElectCertPublishdateD:                                       req.ElectCertPublishdateD,
		GrantdateE:                                                  req.GrantdateE,
		GrantdateY:                                                  req.GrantdateY,
		GrantdateM:                                                  req.GrantdateM,
		GrantdateD:                                                  req.GrantdateD,
		TranspotationBureauChiefName:                                req.TranspotationBureauChiefName,
		EntryNoCarNo:                                                req.EntryNoCarNo,
		RegGrantdateE:                                               req.RegGrantdateE,
		RegGrantdateY:                                               req.RegGrantdateY,
		RegGrantdateM:                                               req.RegGrantdateM,
		RegGrantdateD:                                               req.RegGrantdateD,
		FirstRegistDateE:                                            req.FirstRegistDateE,
		FirstRegistDateY:                                            req.FirstRegistDateY,
		FirstRegistDateM:                                            req.FirstRegistDateM,
		CarName:                                                     req.CarName,
		CarNameCode:                                                 req.CarNameCode,
		CarNo:                                                       req.CarNo,
		Model:                                                       req.Model,
		EngineModel:                                                 req.EngineModel,
		OwnerNameLowLevelChar:                                       req.OwnerNameLowLevelChar,
		OwnerNameHighLevelChar:                                      req.OwnerNameHighLevelChar,
		OwnerAddressChar:                                            req.OwnerAddressChar,
		OwnerAddressNumValue:                                        req.OwnerAddressNumValue,
		OwnerAddressCode:                                            req.OwnerAddressCode,
		UserNameLowLevelChar:                                        req.UserNameLowLevelChar,
		UserNameHighLevelChar:                                       req.UserNameHighLevelChar,
		UserAddressChar:                                             req.UserAddressChar,
		UserAddressNumValue:                                         req.UserAddressNumValue,
		UserAddressCode:                                             req.UserAddressCode,
		UseHeadqurterChar:                                           req.UseHeadqurterChar,
		UseHeadqurterNumValue:                                       req.UseHeadqurterNumValue,
		UseHeadqurterCode:                                           req.UseHeadqurterCode,
		CarKind:                                                     req.CarKind,
		Use:                                                         req.Use,
		PrivateBusiness:                                             req.PrivateBusiness,
		CarShape:                                                    req.CarShape,
		CarShapeCode:                                                req.CarShapeCode,
		NoteCap:                                                     req.NoteCap,
		Cap:                                                         req.Cap,
		NoteMaxLoadage:                                              req.NoteMaxLoadage,
		MaxLoadage:                                                  req.MaxLoadage,
		NoteCarWgt:                                                  req.NoteCarWgt,
		CarWgt:                                                      req.CarWgt,
		NoteCarTotalWgt:                                             req.NoteCarTotalWgt,
		CarTotalWgt:                                                 req.CarTotalWgt,
		NoteLength:                                                  req.NoteLength,
		Length:                                                      req.Length,
		NoteWidth:                                                   req.NoteWidth,
		Width:                                                       req.Width,
		NoteHeight:                                                  req.NoteHeight,
		Height:                                                      req.Height,
		FfAxWgt:                                                     req.FfAxWgt,
		FrAxWgt:                                                     req.FrAxWgt,
		RfAxWgt:                                                     req.RfAxWgt,
		RrAxWgt:                                                     req.RrAxWgt,
		Displacement:                                                req.Displacement,
		FuelClass:                                                   req.FuelClass,
		ModelSpecifyNo:                                              req.ModelSpecifyNo,
		ClassifyAroundNo:                                            req.ClassifyAroundNo,
		ValidPeriodExpirDateE:                                       req.ValidPeriodExpirDateE,
		ValidPeriodExpirDateY:                                       req.ValidPeriodExpirDateY,
		ValidPeriodExpirDateM:                                       req.ValidPeriodExpirDateM,
		ValidPeriodExpirDateD:                                       req.ValidPeriodExpirDateD,
		NoteInfo:                                                    req.NoteInfo,
		TwodimensionCodeInfoEntryNoCarNo:                            req.TwodimensionCodeInfoEntryNoCarNo,
		TwodimensionCodeInfoCarNo:                                   req.TwodimensionCodeInfoCarNo,
		TwodimensionCodeInfoValidPeriodExpirDate:                    req.TwodimensionCodeInfoValidPeriodExpirDate,
		TwodimensionCodeInfoModel:                                   req.TwodimensionCodeInfoModel,
		TwodimensionCodeInfoModelSpecifyNoClassifyAroundNo:          req.TwodimensionCodeInfoModelSpecifyNoClassifyAroundNo,
		TwodimensionCodeInfoCharInfo:                                req.TwodimensionCodeInfoCharInfo,
		TwodimensionCodeInfoEngineModel:                             req.TwodimensionCodeInfoEngineModel,
		TwodimensionCodeInfoCarNoStampPlace:                         req.TwodimensionCodeInfoCarNoStampPlace,
		TwodimensionCodeInfoFirstRegistDate:                         req.TwodimensionCodeInfoFirstRegistDate,
		TwodimensionCodeInfoFfAxWgt:                                 req.TwodimensionCodeInfoFfAxWgt,
		TwodimensionCodeInfoFrAxWgt:                                 req.TwodimensionCodeInfoFrAxWgt,
		TwodimensionCodeInfoRfAxWgt:                                 req.TwodimensionCodeInfoRfAxWgt,
		TwodimensionCodeInfoRrAxWgt:                                 req.TwodimensionCodeInfoRrAxWgt,
		TwodimensionCodeInfoNoiseReg:                                req.TwodimensionCodeInfoNoiseReg,
		TwodimensionCodeInfoNearNoiseReg:                            req.TwodimensionCodeInfoNearNoiseReg,
		TwodimensionCodeInfoDriveMethod:                             req.TwodimensionCodeInfoDriveMethod,
		TwodimensionCodeInfoOpacimeterMeasCar:                       req.TwodimensionCodeInfoOpacimeterMeasCar,
		TwodimensionCodeInfoNoxPmMeasMode:                           req.TwodimensionCodeInfoNoxPmMeasMode,
		TwodimensionCodeInfoNoxValue:                                req.TwodimensionCodeInfoNoxValue,
		TwodimensionCodeInfoPmValue:                                 req.TwodimensionCodeInfoPmValue,
		TwodimensionCodeInfoSafeStdDate:                             req.TwodimensionCodeInfoSafeStdDate,
		TwodimensionCodeInfoFuelClassCode:                           req.TwodimensionCodeInfoFuelClassCode,
		RegistCarLightCar:                                           req.RegistCarLightCar,
		Created:                                                     req.Created,
		Modified:                                                    req.Modified,
	}

	result, err := s.repo.Create(ctx, inspection)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create car inspection: %v", err)
	}

	return &pb.CreateCarInspectionResponse{
		CarInspection: toProtoCarInspection(result),
	}, nil
}

// GetCarInspection retrieves a car inspection by composite primary key
func (s *CarInspectionServer) GetCarInspection(ctx context.Context, req *pb.GetCarInspectionRequest) (*pb.GetCarInspectionResponse, error) {
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}
	if req.ElectCertMgNo == "" {
		return nil, status.Error(codes.InvalidArgument, "elect_cert_mg_no is required")
	}
	if req.ElectCertPublishdateE == "" {
		return nil, status.Error(codes.InvalidArgument, "elect_cert_publishdate_e is required")
	}
	if req.ElectCertPublishdateY == "" {
		return nil, status.Error(codes.InvalidArgument, "elect_cert_publishdate_y is required")
	}
	if req.ElectCertPublishdateM == "" {
		return nil, status.Error(codes.InvalidArgument, "elect_cert_publishdate_m is required")
	}
	if req.ElectCertPublishdateD == "" {
		return nil, status.Error(codes.InvalidArgument, "elect_cert_publishdate_d is required")
	}

	inspection, err := s.repo.GetByPrimaryKey(ctx, req.OrganizationId, req.ElectCertMgNo, req.ElectCertPublishdateE, req.ElectCertPublishdateY, req.ElectCertPublishdateM, req.ElectCertPublishdateD)
	if err != nil {
		if errors.Is(err, repository.ErrCarInspectionNotFound) {
			return nil, status.Error(codes.NotFound, "car inspection not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get car inspection: %v", err)
	}

	return &pb.GetCarInspectionResponse{
		CarInspection: toProtoCarInspection(inspection),
	}, nil
}

// UpdateCarInspection updates an existing car inspection
func (s *CarInspectionServer) UpdateCarInspection(ctx context.Context, req *pb.UpdateCarInspectionRequest) (*pb.UpdateCarInspectionResponse, error) {
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}
	if req.ElectCertMgNo == "" {
		return nil, status.Error(codes.InvalidArgument, "elect_cert_mg_no is required")
	}
	if req.Modified == "" {
		return nil, status.Error(codes.InvalidArgument, "modified is required")
	}

	// Get the existing record first to preserve the created timestamp
	existing, err := s.repo.GetByPrimaryKey(ctx, req.OrganizationId, req.ElectCertMgNo, req.ElectCertPublishdateE, req.ElectCertPublishdateY, req.ElectCertPublishdateM, req.ElectCertPublishdateD)
	if err != nil {
		if errors.Is(err, repository.ErrCarInspectionNotFound) {
			return nil, status.Error(codes.NotFound, "car inspection not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get car inspection: %v", err)
	}

	inspection := &repository.CarInspection{
		OrganizationID:                                              req.OrganizationId,
		CertInfoImportFileVersion:                                   req.CertInfoImportFileVersion,
		AcceptOutputNo:                                              req.AcceptOutputNo,
		FormType:                                                    req.FormType,
		ElectCertMgNo:                                               req.ElectCertMgNo,
		CarID:                                                       req.CarId,
		ElectCertPublishdateE:                                       req.ElectCertPublishdateE,
		ElectCertPublishdateY:                                       req.ElectCertPublishdateY,
		ElectCertPublishdateM:                                       req.ElectCertPublishdateM,
		ElectCertPublishdateD:                                       req.ElectCertPublishdateD,
		GrantdateE:                                                  req.GrantdateE,
		GrantdateY:                                                  req.GrantdateY,
		GrantdateM:                                                  req.GrantdateM,
		GrantdateD:                                                  req.GrantdateD,
		TranspotationBureauChiefName:                                req.TranspotationBureauChiefName,
		EntryNoCarNo:                                                req.EntryNoCarNo,
		RegGrantdateE:                                               req.RegGrantdateE,
		RegGrantdateY:                                               req.RegGrantdateY,
		RegGrantdateM:                                               req.RegGrantdateM,
		RegGrantdateD:                                               req.RegGrantdateD,
		FirstRegistDateE:                                            req.FirstRegistDateE,
		FirstRegistDateY:                                            req.FirstRegistDateY,
		FirstRegistDateM:                                            req.FirstRegistDateM,
		CarName:                                                     req.CarName,
		CarNameCode:                                                 req.CarNameCode,
		CarNo:                                                       req.CarNo,
		Model:                                                       req.Model,
		EngineModel:                                                 req.EngineModel,
		OwnerNameLowLevelChar:                                       req.OwnerNameLowLevelChar,
		OwnerNameHighLevelChar:                                      req.OwnerNameHighLevelChar,
		OwnerAddressChar:                                            req.OwnerAddressChar,
		OwnerAddressNumValue:                                        req.OwnerAddressNumValue,
		OwnerAddressCode:                                            req.OwnerAddressCode,
		UserNameLowLevelChar:                                        req.UserNameLowLevelChar,
		UserNameHighLevelChar:                                       req.UserNameHighLevelChar,
		UserAddressChar:                                             req.UserAddressChar,
		UserAddressNumValue:                                         req.UserAddressNumValue,
		UserAddressCode:                                             req.UserAddressCode,
		UseHeadqurterChar:                                           req.UseHeadqurterChar,
		UseHeadqurterNumValue:                                       req.UseHeadqurterNumValue,
		UseHeadqurterCode:                                           req.UseHeadqurterCode,
		CarKind:                                                     req.CarKind,
		Use:                                                         req.Use,
		PrivateBusiness:                                             req.PrivateBusiness,
		CarShape:                                                    req.CarShape,
		CarShapeCode:                                                req.CarShapeCode,
		NoteCap:                                                     req.NoteCap,
		Cap:                                                         req.Cap,
		NoteMaxLoadage:                                              req.NoteMaxLoadage,
		MaxLoadage:                                                  req.MaxLoadage,
		NoteCarWgt:                                                  req.NoteCarWgt,
		CarWgt:                                                      req.CarWgt,
		NoteCarTotalWgt:                                             req.NoteCarTotalWgt,
		CarTotalWgt:                                                 req.CarTotalWgt,
		NoteLength:                                                  req.NoteLength,
		Length:                                                      req.Length,
		NoteWidth:                                                   req.NoteWidth,
		Width:                                                       req.Width,
		NoteHeight:                                                  req.NoteHeight,
		Height:                                                      req.Height,
		FfAxWgt:                                                     req.FfAxWgt,
		FrAxWgt:                                                     req.FrAxWgt,
		RfAxWgt:                                                     req.RfAxWgt,
		RrAxWgt:                                                     req.RrAxWgt,
		Displacement:                                                req.Displacement,
		FuelClass:                                                   req.FuelClass,
		ModelSpecifyNo:                                              req.ModelSpecifyNo,
		ClassifyAroundNo:                                            req.ClassifyAroundNo,
		ValidPeriodExpirDateE:                                       req.ValidPeriodExpirDateE,
		ValidPeriodExpirDateY:                                       req.ValidPeriodExpirDateY,
		ValidPeriodExpirDateM:                                       req.ValidPeriodExpirDateM,
		ValidPeriodExpirDateD:                                       req.ValidPeriodExpirDateD,
		NoteInfo:                                                    req.NoteInfo,
		TwodimensionCodeInfoEntryNoCarNo:                            req.TwodimensionCodeInfoEntryNoCarNo,
		TwodimensionCodeInfoCarNo:                                   req.TwodimensionCodeInfoCarNo,
		TwodimensionCodeInfoValidPeriodExpirDate:                    req.TwodimensionCodeInfoValidPeriodExpirDate,
		TwodimensionCodeInfoModel:                                   req.TwodimensionCodeInfoModel,
		TwodimensionCodeInfoModelSpecifyNoClassifyAroundNo:          req.TwodimensionCodeInfoModelSpecifyNoClassifyAroundNo,
		TwodimensionCodeInfoCharInfo:                                req.TwodimensionCodeInfoCharInfo,
		TwodimensionCodeInfoEngineModel:                             req.TwodimensionCodeInfoEngineModel,
		TwodimensionCodeInfoCarNoStampPlace:                         req.TwodimensionCodeInfoCarNoStampPlace,
		TwodimensionCodeInfoFirstRegistDate:                         req.TwodimensionCodeInfoFirstRegistDate,
		TwodimensionCodeInfoFfAxWgt:                                 req.TwodimensionCodeInfoFfAxWgt,
		TwodimensionCodeInfoFrAxWgt:                                 req.TwodimensionCodeInfoFrAxWgt,
		TwodimensionCodeInfoRfAxWgt:                                 req.TwodimensionCodeInfoRfAxWgt,
		TwodimensionCodeInfoRrAxWgt:                                 req.TwodimensionCodeInfoRrAxWgt,
		TwodimensionCodeInfoNoiseReg:                                req.TwodimensionCodeInfoNoiseReg,
		TwodimensionCodeInfoNearNoiseReg:                            req.TwodimensionCodeInfoNearNoiseReg,
		TwodimensionCodeInfoDriveMethod:                             req.TwodimensionCodeInfoDriveMethod,
		TwodimensionCodeInfoOpacimeterMeasCar:                       req.TwodimensionCodeInfoOpacimeterMeasCar,
		TwodimensionCodeInfoNoxPmMeasMode:                           req.TwodimensionCodeInfoNoxPmMeasMode,
		TwodimensionCodeInfoNoxValue:                                req.TwodimensionCodeInfoNoxValue,
		TwodimensionCodeInfoPmValue:                                 req.TwodimensionCodeInfoPmValue,
		TwodimensionCodeInfoSafeStdDate:                             req.TwodimensionCodeInfoSafeStdDate,
		TwodimensionCodeInfoFuelClassCode:                           req.TwodimensionCodeInfoFuelClassCode,
		RegistCarLightCar:                                           req.RegistCarLightCar,
		Created:                                                     existing.Created,
		Modified:                                                    req.Modified,
	}

	result, err := s.repo.Update(ctx, inspection)
	if err != nil {
		if errors.Is(err, repository.ErrCarInspectionNotFound) {
			return nil, status.Error(codes.NotFound, "car inspection not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update car inspection: %v", err)
	}

	return &pb.UpdateCarInspectionResponse{
		CarInspection: toProtoCarInspection(result),
	}, nil
}

// DeleteCarInspection hard-deletes a car inspection
func (s *CarInspectionServer) DeleteCarInspection(ctx context.Context, req *pb.DeleteCarInspectionRequest) (*pb.DeleteCarInspectionResponse, error) {
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}
	if req.ElectCertMgNo == "" {
		return nil, status.Error(codes.InvalidArgument, "elect_cert_mg_no is required")
	}
	if req.ElectCertPublishdateE == "" {
		return nil, status.Error(codes.InvalidArgument, "elect_cert_publishdate_e is required")
	}
	if req.ElectCertPublishdateY == "" {
		return nil, status.Error(codes.InvalidArgument, "elect_cert_publishdate_y is required")
	}
	if req.ElectCertPublishdateM == "" {
		return nil, status.Error(codes.InvalidArgument, "elect_cert_publishdate_m is required")
	}
	if req.ElectCertPublishdateD == "" {
		return nil, status.Error(codes.InvalidArgument, "elect_cert_publishdate_d is required")
	}

	err := s.repo.Delete(ctx, req.OrganizationId, req.ElectCertMgNo, req.ElectCertPublishdateE, req.ElectCertPublishdateY, req.ElectCertPublishdateM, req.ElectCertPublishdateD)
	if err != nil {
		if errors.Is(err, repository.ErrCarInspectionNotFound) {
			return nil, status.Error(codes.NotFound, "car inspection not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete car inspection: %v", err)
	}

	return &pb.DeleteCarInspectionResponse{
		Success: true,
	}, nil
}

// ListCarInspections retrieves all car inspections with pagination
func (s *CarInspectionServer) ListCarInspections(ctx context.Context, req *pb.ListCarInspectionsRequest) (*pb.ListCarInspectionsResponse, error) {
	limit := int(req.PageSize)
	if limit <= 0 {
		limit = 10
	}

	offset := 0
	if req.PageToken != "" {
		// In a real implementation, decode page_token to get offset
	}

	inspections, err := s.repo.List(ctx, limit+1, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list car inspections: %v", err)
	}

	var nextPageToken string
	if len(inspections) > limit {
		inspections = inspections[:limit]
		nextPageToken = "next"
	}

	protoInspections := make([]*pb.CarInspection, len(inspections))
	for i, inspection := range inspections {
		protoInspections[i] = toProtoCarInspection(inspection)
	}

	return &pb.ListCarInspectionsResponse{
		CarInspections: protoInspections,
		NextPageToken:  nextPageToken,
	}, nil
}

// ListCarInspectionsByOrganization retrieves car inspections by organization with pagination
func (s *CarInspectionServer) ListCarInspectionsByOrganization(ctx context.Context, req *pb.ListCarInspectionsByOrganizationRequest) (*pb.ListCarInspectionsByOrganizationResponse, error) {
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}

	limit := int(req.PageSize)
	if limit <= 0 {
		limit = 10
	}

	offset := 0
	if req.PageToken != "" {
		// In a real implementation, decode page_token to get offset
	}

	inspections, err := s.repo.ListByOrganization(ctx, req.OrganizationId, limit+1, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list car inspections by organization: %v", err)
	}

	var nextPageToken string
	if len(inspections) > limit {
		inspections = inspections[:limit]
		nextPageToken = "next"
	}

	protoInspections := make([]*pb.CarInspection, len(inspections))
	for i, inspection := range inspections {
		protoInspections[i] = toProtoCarInspection(inspection)
	}

	return &pb.ListCarInspectionsByOrganizationResponse{
		CarInspections: protoInspections,
		NextPageToken:  nextPageToken,
	}, nil
}

// toProtoCarInspection converts repository model to proto message
func toProtoCarInspection(inspection *repository.CarInspection) *pb.CarInspection {
	return &pb.CarInspection{
		OrganizationId:                                            inspection.OrganizationID,
		CertInfoImportFileVersion:                                 inspection.CertInfoImportFileVersion,
		AcceptOutputNo:                                            inspection.AcceptOutputNo,
		FormType:                                                  inspection.FormType,
		ElectCertMgNo:                                             inspection.ElectCertMgNo,
		CarId:                                                     inspection.CarID,
		ElectCertPublishdateE:                                     inspection.ElectCertPublishdateE,
		ElectCertPublishdateY:                                     inspection.ElectCertPublishdateY,
		ElectCertPublishdateM:                                     inspection.ElectCertPublishdateM,
		ElectCertPublishdateD:                                     inspection.ElectCertPublishdateD,
		GrantdateE:                                                inspection.GrantdateE,
		GrantdateY:                                                inspection.GrantdateY,
		GrantdateM:                                                inspection.GrantdateM,
		GrantdateD:                                                inspection.GrantdateD,
		TranspotationBureauChiefName:                              inspection.TranspotationBureauChiefName,
		EntryNoCarNo:                                              inspection.EntryNoCarNo,
		RegGrantdateE:                                             inspection.RegGrantdateE,
		RegGrantdateY:                                             inspection.RegGrantdateY,
		RegGrantdateM:                                             inspection.RegGrantdateM,
		RegGrantdateD:                                             inspection.RegGrantdateD,
		FirstRegistDateE:                                          inspection.FirstRegistDateE,
		FirstRegistDateY:                                          inspection.FirstRegistDateY,
		FirstRegistDateM:                                          inspection.FirstRegistDateM,
		CarName:                                                   inspection.CarName,
		CarNameCode:                                               inspection.CarNameCode,
		CarNo:                                                     inspection.CarNo,
		Model:                                                     inspection.Model,
		EngineModel:                                               inspection.EngineModel,
		OwnerNameLowLevelChar:                                     inspection.OwnerNameLowLevelChar,
		OwnerNameHighLevelChar:                                    inspection.OwnerNameHighLevelChar,
		OwnerAddressChar:                                          inspection.OwnerAddressChar,
		OwnerAddressNumValue:                                      inspection.OwnerAddressNumValue,
		OwnerAddressCode:                                          inspection.OwnerAddressCode,
		UserNameLowLevelChar:                                      inspection.UserNameLowLevelChar,
		UserNameHighLevelChar:                                     inspection.UserNameHighLevelChar,
		UserAddressChar:                                           inspection.UserAddressChar,
		UserAddressNumValue:                                       inspection.UserAddressNumValue,
		UserAddressCode:                                           inspection.UserAddressCode,
		UseHeadqurterChar:                                         inspection.UseHeadqurterChar,
		UseHeadqurterNumValue:                                     inspection.UseHeadqurterNumValue,
		UseHeadqurterCode:                                         inspection.UseHeadqurterCode,
		CarKind:                                                   inspection.CarKind,
		Use:                                                       inspection.Use,
		PrivateBusiness:                                           inspection.PrivateBusiness,
		CarShape:                                                  inspection.CarShape,
		CarShapeCode:                                              inspection.CarShapeCode,
		NoteCap:                                                   inspection.NoteCap,
		Cap:                                                       inspection.Cap,
		NoteMaxLoadage:                                            inspection.NoteMaxLoadage,
		MaxLoadage:                                                inspection.MaxLoadage,
		NoteCarWgt:                                                inspection.NoteCarWgt,
		CarWgt:                                                    inspection.CarWgt,
		NoteCarTotalWgt:                                           inspection.NoteCarTotalWgt,
		CarTotalWgt:                                               inspection.CarTotalWgt,
		NoteLength:                                                inspection.NoteLength,
		Length:                                                    inspection.Length,
		NoteWidth:                                                 inspection.NoteWidth,
		Width:                                                     inspection.Width,
		NoteHeight:                                                inspection.NoteHeight,
		Height:                                                    inspection.Height,
		FfAxWgt:                                                   inspection.FfAxWgt,
		FrAxWgt:                                                   inspection.FrAxWgt,
		RfAxWgt:                                                   inspection.RfAxWgt,
		RrAxWgt:                                                   inspection.RrAxWgt,
		Displacement:                                              inspection.Displacement,
		FuelClass:                                                 inspection.FuelClass,
		ModelSpecifyNo:                                            inspection.ModelSpecifyNo,
		ClassifyAroundNo:                                          inspection.ClassifyAroundNo,
		ValidPeriodExpirDateE:                                     inspection.ValidPeriodExpirDateE,
		ValidPeriodExpirDateY:                                     inspection.ValidPeriodExpirDateY,
		ValidPeriodExpirDateM:                                     inspection.ValidPeriodExpirDateM,
		ValidPeriodExpirDateD:                                     inspection.ValidPeriodExpirDateD,
		NoteInfo:                                                  inspection.NoteInfo,
		TwodimensionCodeInfoEntryNoCarNo:                          inspection.TwodimensionCodeInfoEntryNoCarNo,
		TwodimensionCodeInfoCarNo:                                 inspection.TwodimensionCodeInfoCarNo,
		TwodimensionCodeInfoValidPeriodExpirDate:                  inspection.TwodimensionCodeInfoValidPeriodExpirDate,
		TwodimensionCodeInfoModel:                                 inspection.TwodimensionCodeInfoModel,
		TwodimensionCodeInfoModelSpecifyNoClassifyAroundNo:        inspection.TwodimensionCodeInfoModelSpecifyNoClassifyAroundNo,
		TwodimensionCodeInfoCharInfo:                              inspection.TwodimensionCodeInfoCharInfo,
		TwodimensionCodeInfoEngineModel:                           inspection.TwodimensionCodeInfoEngineModel,
		TwodimensionCodeInfoCarNoStampPlace:                       inspection.TwodimensionCodeInfoCarNoStampPlace,
		TwodimensionCodeInfoFirstRegistDate:                       inspection.TwodimensionCodeInfoFirstRegistDate,
		TwodimensionCodeInfoFfAxWgt:                               inspection.TwodimensionCodeInfoFfAxWgt,
		TwodimensionCodeInfoFrAxWgt:                               inspection.TwodimensionCodeInfoFrAxWgt,
		TwodimensionCodeInfoRfAxWgt:                               inspection.TwodimensionCodeInfoRfAxWgt,
		TwodimensionCodeInfoRrAxWgt:                               inspection.TwodimensionCodeInfoRrAxWgt,
		TwodimensionCodeInfoNoiseReg:                              inspection.TwodimensionCodeInfoNoiseReg,
		TwodimensionCodeInfoNearNoiseReg:                          inspection.TwodimensionCodeInfoNearNoiseReg,
		TwodimensionCodeInfoDriveMethod:                           inspection.TwodimensionCodeInfoDriveMethod,
		TwodimensionCodeInfoOpacimeterMeasCar:                     inspection.TwodimensionCodeInfoOpacimeterMeasCar,
		TwodimensionCodeInfoNoxPmMeasMode:                         inspection.TwodimensionCodeInfoNoxPmMeasMode,
		TwodimensionCodeInfoNoxValue:                              inspection.TwodimensionCodeInfoNoxValue,
		TwodimensionCodeInfoPmValue:                               inspection.TwodimensionCodeInfoPmValue,
		TwodimensionCodeInfoSafeStdDate:                           inspection.TwodimensionCodeInfoSafeStdDate,
		TwodimensionCodeInfoFuelClassCode:                         inspection.TwodimensionCodeInfoFuelClassCode,
		RegistCarLightCar:                                         inspection.RegistCarLightCar,
		Created:                                                   inspection.Created,
		Modified:                                                  inspection.Modified,
	}
}
