package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/pb"
	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/repository"
)

// KudgivtServer implements the gRPC KudgivtService
type KudgivtServer struct {
	pb.UnimplementedKudgivtServiceServer
	repo *repository.KudgivtRepository
}

// NewKudgivtServer creates a new gRPC server
func NewKudgivtServer(repo *repository.KudgivtRepository) *KudgivtServer {
	return &KudgivtServer{repo: repo}
}

// CreateKudgivt creates a new kudgivt record
func (s *KudgivtServer) CreateKudgivt(ctx context.Context, req *pb.CreateKudgivtRequest) (*pb.CreateKudgivtResponse, error) {
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}
	if req.Hash == "" {
		return nil, status.Error(codes.InvalidArgument, "hash is required")
	}
	if req.Created == "" {
		return nil, status.Error(codes.InvalidArgument, "created is required")
	}
	if req.TargetDriverType == "" {
		return nil, status.Error(codes.InvalidArgument, "target_driver_type is required")
	}

	kudgivt := &repository.Kudgivt{
		OrganizationID:              req.OrganizationId,
		Hash:                        req.Hash,
		Created:                     req.Created,
		Deleted:                     ptrFromOptional(req.Deleted),
		KudguriUuid:                 ptrFromOptional(req.KudguriUuid),
		UnkouNo:                     ptrFromOptional(req.UnkouNo),
		ReadDate:                    ptrFromOptional(req.ReadDate),
		UnkouDate:                   ptrFromOptional(req.UnkouDate),
		OfficeCd:                    ptrFromOptional(req.OfficeCd),
		OfficeName:                  ptrFromOptional(req.OfficeName),
		VehicleCd:                   ptrFromOptional(req.VehicleCd),
		VehicleName:                 ptrFromOptional(req.VehicleName),
		DriverCd1:                   ptrFromOptional(req.DriverCd1),
		DriverName1:                 ptrFromOptional(req.DriverName1),
		TargetDriverType:            req.TargetDriverType,
		TargetDriverCd:              ptrFromOptional(req.TargetDriverCd),
		TargetDriverName:            ptrFromOptional(req.TargetDriverName),
		ClockInDatetime:             ptrFromOptional(req.ClockInDatetime),
		ClockOutDatetime:            ptrFromOptional(req.ClockOutDatetime),
		DepartureDatetime:           ptrFromOptional(req.DepartureDatetime),
		ReturnDatetime:              ptrFromOptional(req.ReturnDatetime),
		DepartureMeter:              ptrFromOptional(req.DepartureMeter),
		ReturnMeter:                 ptrFromOptional(req.ReturnMeter),
		TotalMileage:                ptrFromOptional(req.TotalMileage),
		DestinationCityName:         ptrFromOptional(req.DestinationCityName),
		DestinationPlaceName:        ptrFromOptional(req.DestinationPlaceName),
		ActualMileage:               ptrFromOptional(req.ActualMileage),
		LocalDriveTime:              ptrFromOptional(req.LocalDriveTime),
		ExpressDriveTime:            ptrFromOptional(req.ExpressDriveTime),
		BypassDriveTime:             ptrFromOptional(req.BypassDriveTime),
		ActualDriveTime:             ptrFromOptional(req.ActualDriveTime),
		EmptyDriveTime:              ptrFromOptional(req.EmptyDriveTime),
		Work1Time:                   ptrFromOptional(req.Work1Time),
		Work2Time:                   ptrFromOptional(req.Work2Time),
		Work3Time:                   ptrFromOptional(req.Work3Time),
		Work4Time:                   ptrFromOptional(req.Work4Time),
		Work5Time:                   ptrFromOptional(req.Work5Time),
		Work6Time:                   ptrFromOptional(req.Work6Time),
		Work7Time:                   ptrFromOptional(req.Work7Time),
		Work8Time:                   ptrFromOptional(req.Work8Time),
		Work9Time:                   ptrFromOptional(req.Work9Time),
		Work10Time:                  ptrFromOptional(req.Work10Time),
		State1Distance:              ptrFromOptional(req.State1Distance),
		State2Distance:              ptrFromOptional(req.State2Distance),
		State3Distance:              ptrFromOptional(req.State3Distance),
		State4Distance:              ptrFromOptional(req.State4Distance),
		State5Distance:              ptrFromOptional(req.State5Distance),
		State1Time:                  ptrFromOptional(req.State1Time),
		State2Time:                  ptrFromOptional(req.State2Time),
		State3Time:                  ptrFromOptional(req.State3Time),
		State4Time:                  ptrFromOptional(req.State4Time),
		State5Time:                  ptrFromOptional(req.State5Time),
		OwnMainFuel:                 ptrFromOptional(req.OwnMainFuel),
		OwnMainAdditive:             ptrFromOptional(req.OwnMainAdditive),
		OwnConsumable:               ptrFromOptional(req.OwnConsumable),
		OtherMainFuel:               ptrFromOptional(req.OtherMainFuel),
		OtherMainAdditive:           ptrFromOptional(req.OtherMainAdditive),
		OtherConsumable:             ptrFromOptional(req.OtherConsumable),
		LocalSpeedOverMax:           ptrFromOptional(req.LocalSpeedOverMax),
		LocalSpeedOverTime:          ptrFromOptional(req.LocalSpeedOverTime),
		LocalSpeedOverCount:         ptrFromOptional(req.LocalSpeedOverCount),
		ExpressSpeedOverMax:         ptrFromOptional(req.ExpressSpeedOverMax),
		ExpressSpeedOverTime:        ptrFromOptional(req.ExpressSpeedOverTime),
		ExpressSpeedOverCount:       ptrFromOptional(req.ExpressSpeedOverCount),
		DedicatedSpeedOverMax:       ptrFromOptional(req.DedicatedSpeedOverMax),
		DedicatedSpeedOverTime:      ptrFromOptional(req.DedicatedSpeedOverTime),
		DedicatedSpeedOverCount:     ptrFromOptional(req.DedicatedSpeedOverCount),
		IdlingTime:                  ptrFromOptional(req.IdlingTime),
		IdlingTimeCount:             ptrFromOptional(req.IdlingTimeCount),
		RotationOverMax:             ptrFromOptional(req.RotationOverMax),
		RotationOverCount:           ptrFromOptional(req.RotationOverCount),
		RotationOverTime:            ptrFromOptional(req.RotationOverTime),
		RapidAccelCount1:            ptrFromOptional(req.RapidAccelCount1),
		RapidAccelCount2:            ptrFromOptional(req.RapidAccelCount2),
		RapidAccelCount3:            ptrFromOptional(req.RapidAccelCount3),
		RapidAccelCount4:            ptrFromOptional(req.RapidAccelCount4),
		RapidAccelCount5:            ptrFromOptional(req.RapidAccelCount5),
		RapidAccelMax:               ptrFromOptional(req.RapidAccelMax),
		RapidAccelMaxSpeed:          ptrFromOptional(req.RapidAccelMaxSpeed),
		RapidDecelCount1:            ptrFromOptional(req.RapidDecelCount1),
		RapidDecelCount2:            ptrFromOptional(req.RapidDecelCount2),
		RapidDecelCount3:            ptrFromOptional(req.RapidDecelCount3),
		RapidDecelCount4:            ptrFromOptional(req.RapidDecelCount4),
		RapidDecelCount5:            ptrFromOptional(req.RapidDecelCount5),
		RapidDecelMax:               ptrFromOptional(req.RapidDecelMax),
		RapidDecelMaxSpeed:          ptrFromOptional(req.RapidDecelMaxSpeed),
		RapidCurveCount1:            ptrFromOptional(req.RapidCurveCount1),
		RapidCurveCount2:            ptrFromOptional(req.RapidCurveCount2),
		RapidCurveCount3:            ptrFromOptional(req.RapidCurveCount3),
		RapidCurveCount4:            ptrFromOptional(req.RapidCurveCount4),
		RapidCurveCount5:             ptrFromOptional(req.RapidCurveCount5),
		RapidCurveMax:                ptrFromOptional(req.RapidCurveMax),
		RapidCurveMaxSpeed:           ptrFromOptional(req.RapidCurveMaxSpeed),
		ContinuousDriveOverCount:     nil,
		ContinuousDriveMaxTime:       nil,
		ContinuousDriveTotalTime:     nil,
		WaveDriveCount:               nil,
		WaveDriveMaxTime:             nil,
		WaveDriveMaxSpeedDiff:        nil,
		LocalSpeedScore:              nil,
		ExpressSpeedScore:            nil,
		DedicatedSpeedScore:          nil,
		LocalDistanceScore:           nil,
		ExpressDistanceScore:         nil,
		DedicatedDistanceScore:       nil,
		RapidAccelScore:              nil,
		RapidDecelScore:              nil,
		RapidCurveScore:              nil,
		ActualLowSpeedRotationScore:  nil,
		ActualHighSpeedRotationScore: nil,
		EmptyLowSpeedRotationScore:   nil,
		EmptyHighSpeedRotationScore:  nil,
		IdlingScore:                  nil,
		ContinuousDriveScore:         nil,
		WaveDriveScore:               nil,
		SafetyScore:                  nil,
		EconomyScore:                 nil,
		TotalScore:                   nil,
	}

	result, err := s.repo.Create(ctx, kudgivt)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create kudgivt: %v", err)
	}

	return &pb.CreateKudgivtResponse{
		Kudgivt: toProtoKudgivt(result),
	}, nil
}

// GetKudgivt retrieves a kudgivt record by UUID
func (s *KudgivtServer) GetKudgivt(ctx context.Context, req *pb.GetKudgivtRequest) (*pb.GetKudgivtResponse, error) {
	if req.Uuid == "" {
		return nil, status.Error(codes.InvalidArgument, "uuid is required")
	}

	kudgivt, err := s.repo.GetByUUID(ctx, req.Uuid)
	if err != nil {
		if errors.Is(err, repository.ErrKudgivtNotFound) {
			return nil, status.Error(codes.NotFound, "kudgivt not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get kudgivt: %v", err)
	}

	return &pb.GetKudgivtResponse{
		Kudgivt: toProtoKudgivt(kudgivt),
	}, nil
}

// UpdateKudgivt updates an existing kudgivt record
func (s *KudgivtServer) UpdateKudgivt(ctx context.Context, req *pb.UpdateKudgivtRequest) (*pb.UpdateKudgivtResponse, error) {
	if req.Uuid == "" {
		return nil, status.Error(codes.InvalidArgument, "uuid is required")
	}
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}
	if req.Hash == "" {
		return nil, status.Error(codes.InvalidArgument, "hash is required")
	}
	if req.Created == "" {
		return nil, status.Error(codes.InvalidArgument, "created is required")
	}
	if req.TargetDriverType == "" {
		return nil, status.Error(codes.InvalidArgument, "target_driver_type is required")
	}

	kudgivt := &repository.Kudgivt{
		UUID:                        req.Uuid,
		OrganizationID:              req.OrganizationId,
		Hash:                        req.Hash,
		Created:                     req.Created,
		Deleted:                     ptrFromOptional(req.Deleted),
		KudguriUuid:                 ptrFromOptional(req.KudguriUuid),
		UnkouNo:                     ptrFromOptional(req.UnkouNo),
		ReadDate:                    ptrFromOptional(req.ReadDate),
		UnkouDate:                   ptrFromOptional(req.UnkouDate),
		OfficeCd:                    ptrFromOptional(req.OfficeCd),
		OfficeName:                  ptrFromOptional(req.OfficeName),
		VehicleCd:                   ptrFromOptional(req.VehicleCd),
		VehicleName:                 ptrFromOptional(req.VehicleName),
		DriverCd1:                   ptrFromOptional(req.DriverCd1),
		DriverName1:                 ptrFromOptional(req.DriverName1),
		TargetDriverType:            req.TargetDriverType,
		TargetDriverCd:              ptrFromOptional(req.TargetDriverCd),
		TargetDriverName:            ptrFromOptional(req.TargetDriverName),
		ClockInDatetime:             ptrFromOptional(req.ClockInDatetime),
		ClockOutDatetime:            ptrFromOptional(req.ClockOutDatetime),
		DepartureDatetime:           ptrFromOptional(req.DepartureDatetime),
		ReturnDatetime:              ptrFromOptional(req.ReturnDatetime),
		DepartureMeter:              ptrFromOptional(req.DepartureMeter),
		ReturnMeter:                 ptrFromOptional(req.ReturnMeter),
		TotalMileage:                ptrFromOptional(req.TotalMileage),
		DestinationCityName:         ptrFromOptional(req.DestinationCityName),
		DestinationPlaceName:        ptrFromOptional(req.DestinationPlaceName),
		ActualMileage:               ptrFromOptional(req.ActualMileage),
		LocalDriveTime:              ptrFromOptional(req.LocalDriveTime),
		ExpressDriveTime:            ptrFromOptional(req.ExpressDriveTime),
		BypassDriveTime:             ptrFromOptional(req.BypassDriveTime),
		ActualDriveTime:             ptrFromOptional(req.ActualDriveTime),
		EmptyDriveTime:              ptrFromOptional(req.EmptyDriveTime),
		Work1Time:                   ptrFromOptional(req.Work1Time),
		Work2Time:                   ptrFromOptional(req.Work2Time),
		Work3Time:                   ptrFromOptional(req.Work3Time),
		Work4Time:                   ptrFromOptional(req.Work4Time),
		Work5Time:                   ptrFromOptional(req.Work5Time),
		Work6Time:                   ptrFromOptional(req.Work6Time),
		Work7Time:                   ptrFromOptional(req.Work7Time),
		Work8Time:                   ptrFromOptional(req.Work8Time),
		Work9Time:                   ptrFromOptional(req.Work9Time),
		Work10Time:                  ptrFromOptional(req.Work10Time),
		State1Distance:              ptrFromOptional(req.State1Distance),
		State2Distance:              ptrFromOptional(req.State2Distance),
		State3Distance:              ptrFromOptional(req.State3Distance),
		State4Distance:              ptrFromOptional(req.State4Distance),
		State5Distance:              ptrFromOptional(req.State5Distance),
		State1Time:                  ptrFromOptional(req.State1Time),
		State2Time:                  ptrFromOptional(req.State2Time),
		State3Time:                  ptrFromOptional(req.State3Time),
		State4Time:                  ptrFromOptional(req.State4Time),
		State5Time:                  ptrFromOptional(req.State5Time),
		OwnMainFuel:                 ptrFromOptional(req.OwnMainFuel),
		OwnMainAdditive:             ptrFromOptional(req.OwnMainAdditive),
		OwnConsumable:               ptrFromOptional(req.OwnConsumable),
		OtherMainFuel:               ptrFromOptional(req.OtherMainFuel),
		OtherMainAdditive:           ptrFromOptional(req.OtherMainAdditive),
		OtherConsumable:             ptrFromOptional(req.OtherConsumable),
		LocalSpeedOverMax:           ptrFromOptional(req.LocalSpeedOverMax),
		LocalSpeedOverTime:          ptrFromOptional(req.LocalSpeedOverTime),
		LocalSpeedOverCount:         ptrFromOptional(req.LocalSpeedOverCount),
		ExpressSpeedOverMax:         ptrFromOptional(req.ExpressSpeedOverMax),
		ExpressSpeedOverTime:        ptrFromOptional(req.ExpressSpeedOverTime),
		ExpressSpeedOverCount:       ptrFromOptional(req.ExpressSpeedOverCount),
		DedicatedSpeedOverMax:       ptrFromOptional(req.DedicatedSpeedOverMax),
		DedicatedSpeedOverTime:      ptrFromOptional(req.DedicatedSpeedOverTime),
		DedicatedSpeedOverCount:     ptrFromOptional(req.DedicatedSpeedOverCount),
		IdlingTime:                  ptrFromOptional(req.IdlingTime),
		IdlingTimeCount:             ptrFromOptional(req.IdlingTimeCount),
		RotationOverMax:             ptrFromOptional(req.RotationOverMax),
		RotationOverCount:           ptrFromOptional(req.RotationOverCount),
		RotationOverTime:            ptrFromOptional(req.RotationOverTime),
		RapidAccelCount1:            ptrFromOptional(req.RapidAccelCount1),
		RapidAccelCount2:            ptrFromOptional(req.RapidAccelCount2),
		RapidAccelCount3:            ptrFromOptional(req.RapidAccelCount3),
		RapidAccelCount4:            ptrFromOptional(req.RapidAccelCount4),
		RapidAccelCount5:            ptrFromOptional(req.RapidAccelCount5),
		RapidAccelMax:               ptrFromOptional(req.RapidAccelMax),
		RapidAccelMaxSpeed:          ptrFromOptional(req.RapidAccelMaxSpeed),
		RapidDecelCount1:            ptrFromOptional(req.RapidDecelCount1),
		RapidDecelCount2:            ptrFromOptional(req.RapidDecelCount2),
		RapidDecelCount3:            ptrFromOptional(req.RapidDecelCount3),
		RapidDecelCount4:            ptrFromOptional(req.RapidDecelCount4),
		RapidDecelCount5:            ptrFromOptional(req.RapidDecelCount5),
		RapidDecelMax:               ptrFromOptional(req.RapidDecelMax),
		RapidDecelMaxSpeed:          ptrFromOptional(req.RapidDecelMaxSpeed),
		RapidCurveCount1:            ptrFromOptional(req.RapidCurveCount1),
		RapidCurveCount2:            ptrFromOptional(req.RapidCurveCount2),
		RapidCurveCount3:            ptrFromOptional(req.RapidCurveCount3),
		RapidCurveCount4:            ptrFromOptional(req.RapidCurveCount4),
		RapidCurveCount5:             ptrFromOptional(req.RapidCurveCount5),
		RapidCurveMax:                ptrFromOptional(req.RapidCurveMax),
		RapidCurveMaxSpeed:           ptrFromOptional(req.RapidCurveMaxSpeed),
		ContinuousDriveOverCount:     nil,
		ContinuousDriveMaxTime:       nil,
		ContinuousDriveTotalTime:     nil,
		WaveDriveCount:               nil,
		WaveDriveMaxTime:             nil,
		WaveDriveMaxSpeedDiff:        nil,
		LocalSpeedScore:              nil,
		ExpressSpeedScore:            nil,
		DedicatedSpeedScore:          nil,
		LocalDistanceScore:           nil,
		ExpressDistanceScore:         nil,
		DedicatedDistanceScore:       nil,
		RapidAccelScore:              nil,
		RapidDecelScore:              nil,
		RapidCurveScore:              nil,
		ActualLowSpeedRotationScore:  nil,
		ActualHighSpeedRotationScore: nil,
		EmptyLowSpeedRotationScore:   nil,
		EmptyHighSpeedRotationScore:  nil,
		IdlingScore:                  nil,
		ContinuousDriveScore:         nil,
		WaveDriveScore:               nil,
		SafetyScore:                  nil,
		EconomyScore:                 nil,
		TotalScore:                   nil,
	}

	result, err := s.repo.Update(ctx, kudgivt)
	if err != nil {
		if errors.Is(err, repository.ErrKudgivtNotFound) {
			return nil, status.Error(codes.NotFound, "kudgivt not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update kudgivt: %v", err)
	}

	return &pb.UpdateKudgivtResponse{
		Kudgivt: toProtoKudgivt(result),
	}, nil
}

// DeleteKudgivt soft-deletes a kudgivt record
func (s *KudgivtServer) DeleteKudgivt(ctx context.Context, req *pb.DeleteKudgivtRequest) (*pb.DeleteKudgivtResponse, error) {
	if req.Uuid == "" {
		return nil, status.Error(codes.InvalidArgument, "uuid is required")
	}

	err := s.repo.Delete(ctx, req.Uuid)
	if err != nil {
		if errors.Is(err, repository.ErrKudgivtNotFound) {
			return nil, status.Error(codes.NotFound, "kudgivt not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete kudgivt: %v", err)
	}

	return &pb.DeleteKudgivtResponse{
		Success: true,
	}, nil
}

// ListKudgivts retrieves kudgivt records with pagination
func (s *KudgivtServer) ListKudgivts(ctx context.Context, req *pb.ListKudgivtsRequest) (*pb.ListKudgivtsResponse, error) {
	limit := int(req.PageSize)
	if limit <= 0 {
		limit = 10
	}

	offset := 0
	if req.PageToken != "" {
		// In production, decode page_token to get offset
	}

	kudgivts, err := s.repo.List(ctx, limit+1, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list kudgivts: %v", err)
	}

	var nextPageToken string
	if len(kudgivts) > limit {
		kudgivts = kudgivts[:limit]
		nextPageToken = "next"
	}

	protoKudgivts := make([]*pb.Kudgivt, len(kudgivts))
	for i, kudgivt := range kudgivts {
		protoKudgivts[i] = toProtoKudgivt(kudgivt)
	}

	return &pb.ListKudgivtsResponse{
		Kudgivts:      protoKudgivts,
		NextPageToken: nextPageToken,
	}, nil
}

// ListKudgivtsByOrganization retrieves kudgivt records for a specific organization with pagination
func (s *KudgivtServer) ListKudgivtsByOrganization(ctx context.Context, req *pb.ListKudgivtsByOrganizationRequest) (*pb.ListKudgivtsByOrganizationResponse, error) {
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}

	limit := int(req.PageSize)
	if limit <= 0 {
		limit = 10
	}

	offset := 0
	if req.PageToken != "" {
		// In production, decode page_token to get offset
	}

	kudgivts, err := s.repo.ListByOrganization(ctx, req.OrganizationId, limit+1, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list kudgivts by organization: %v", err)
	}

	var nextPageToken string
	if len(kudgivts) > limit {
		kudgivts = kudgivts[:limit]
		nextPageToken = "next"
	}

	protoKudgivts := make([]*pb.Kudgivt, len(kudgivts))
	for i, kudgivt := range kudgivts {
		protoKudgivts[i] = toProtoKudgivt(kudgivt)
	}

	return &pb.ListKudgivtsByOrganizationResponse{
		Kudgivts:      protoKudgivts,
		NextPageToken: nextPageToken,
	}, nil
}

// toProtoKudgivt converts repository model to proto message
func toProtoKudgivt(k *repository.Kudgivt) *pb.Kudgivt {
	return &pb.Kudgivt{
		Uuid:                        k.UUID,
		OrganizationId:              k.OrganizationID,
		Hash:                        k.Hash,
		Created:                     k.Created,
		Deleted:                     optionalFromPtr(k.Deleted),
		KudguriUuid:                 optionalFromPtr(k.KudguriUuid),
		UnkouNo:                     optionalFromPtr(k.UnkouNo),
		ReadDate:                    optionalFromPtr(k.ReadDate),
		UnkouDate:                   optionalFromPtr(k.UnkouDate),
		OfficeCd:                    optionalFromPtr(k.OfficeCd),
		OfficeName:                  optionalFromPtr(k.OfficeName),
		VehicleCd:                   optionalFromPtr(k.VehicleCd),
		VehicleName:                 optionalFromPtr(k.VehicleName),
		DriverCd1:                   optionalFromPtr(k.DriverCd1),
		DriverName1:                 optionalFromPtr(k.DriverName1),
		TargetDriverType:            k.TargetDriverType,
		TargetDriverCd:              optionalFromPtr(k.TargetDriverCd),
		TargetDriverName:            optionalFromPtr(k.TargetDriverName),
		ClockInDatetime:             optionalFromPtr(k.ClockInDatetime),
		ClockOutDatetime:            optionalFromPtr(k.ClockOutDatetime),
		DepartureDatetime:           optionalFromPtr(k.DepartureDatetime),
		ReturnDatetime:              optionalFromPtr(k.ReturnDatetime),
		DepartureMeter:              optionalFromPtr(k.DepartureMeter),
		ReturnMeter:                 optionalFromPtr(k.ReturnMeter),
		TotalMileage:                optionalFromPtr(k.TotalMileage),
		DestinationCityName:         optionalFromPtr(k.DestinationCityName),
		DestinationPlaceName:        optionalFromPtr(k.DestinationPlaceName),
		ActualMileage:               optionalFromPtr(k.ActualMileage),
		LocalDriveTime:              optionalFromPtr(k.LocalDriveTime),
		ExpressDriveTime:            optionalFromPtr(k.ExpressDriveTime),
		BypassDriveTime:             optionalFromPtr(k.BypassDriveTime),
		ActualDriveTime:             optionalFromPtr(k.ActualDriveTime),
		EmptyDriveTime:              optionalFromPtr(k.EmptyDriveTime),
		Work1Time:                   optionalFromPtr(k.Work1Time),
		Work2Time:                   optionalFromPtr(k.Work2Time),
		Work3Time:                   optionalFromPtr(k.Work3Time),
		Work4Time:                   optionalFromPtr(k.Work4Time),
		Work5Time:                   optionalFromPtr(k.Work5Time),
		Work6Time:                   optionalFromPtr(k.Work6Time),
		Work7Time:                   optionalFromPtr(k.Work7Time),
		Work8Time:                   optionalFromPtr(k.Work8Time),
		Work9Time:                   optionalFromPtr(k.Work9Time),
		Work10Time:                  optionalFromPtr(k.Work10Time),
		State1Distance:              optionalFromPtr(k.State1Distance),
		State2Distance:              optionalFromPtr(k.State2Distance),
		State3Distance:              optionalFromPtr(k.State3Distance),
		State4Distance:              optionalFromPtr(k.State4Distance),
		State5Distance:              optionalFromPtr(k.State5Distance),
		State1Time:                  optionalFromPtr(k.State1Time),
		State2Time:                  optionalFromPtr(k.State2Time),
		State3Time:                  optionalFromPtr(k.State3Time),
		State4Time:                  optionalFromPtr(k.State4Time),
		State5Time:                  optionalFromPtr(k.State5Time),
		OwnMainFuel:                 optionalFromPtr(k.OwnMainFuel),
		OwnMainAdditive:             optionalFromPtr(k.OwnMainAdditive),
		OwnConsumable:               optionalFromPtr(k.OwnConsumable),
		OtherMainFuel:               optionalFromPtr(k.OtherMainFuel),
		OtherMainAdditive:           optionalFromPtr(k.OtherMainAdditive),
		OtherConsumable:             optionalFromPtr(k.OtherConsumable),
		LocalSpeedOverMax:           optionalFromPtr(k.LocalSpeedOverMax),
		LocalSpeedOverTime:          optionalFromPtr(k.LocalSpeedOverTime),
		LocalSpeedOverCount:         optionalFromPtr(k.LocalSpeedOverCount),
		ExpressSpeedOverMax:         optionalFromPtr(k.ExpressSpeedOverMax),
		ExpressSpeedOverTime:        optionalFromPtr(k.ExpressSpeedOverTime),
		ExpressSpeedOverCount:       optionalFromPtr(k.ExpressSpeedOverCount),
		DedicatedSpeedOverMax:       optionalFromPtr(k.DedicatedSpeedOverMax),
		DedicatedSpeedOverTime:      optionalFromPtr(k.DedicatedSpeedOverTime),
		DedicatedSpeedOverCount:     optionalFromPtr(k.DedicatedSpeedOverCount),
		IdlingTime:                  optionalFromPtr(k.IdlingTime),
		IdlingTimeCount:             optionalFromPtr(k.IdlingTimeCount),
		RotationOverMax:             optionalFromPtr(k.RotationOverMax),
		RotationOverCount:           optionalFromPtr(k.RotationOverCount),
		RotationOverTime:            optionalFromPtr(k.RotationOverTime),
		RapidAccelCount1:            optionalFromPtr(k.RapidAccelCount1),
		RapidAccelCount2:            optionalFromPtr(k.RapidAccelCount2),
		RapidAccelCount3:            optionalFromPtr(k.RapidAccelCount3),
		RapidAccelCount4:            optionalFromPtr(k.RapidAccelCount4),
		RapidAccelCount5:            optionalFromPtr(k.RapidAccelCount5),
		RapidAccelMax:               optionalFromPtr(k.RapidAccelMax),
		RapidAccelMaxSpeed:          optionalFromPtr(k.RapidAccelMaxSpeed),
		RapidDecelCount1:            optionalFromPtr(k.RapidDecelCount1),
		RapidDecelCount2:            optionalFromPtr(k.RapidDecelCount2),
		RapidDecelCount3:            optionalFromPtr(k.RapidDecelCount3),
		RapidDecelCount4:            optionalFromPtr(k.RapidDecelCount4),
		RapidDecelCount5:            optionalFromPtr(k.RapidDecelCount5),
		RapidDecelMax:               optionalFromPtr(k.RapidDecelMax),
		RapidDecelMaxSpeed:          optionalFromPtr(k.RapidDecelMaxSpeed),
		RapidCurveCount1:            optionalFromPtr(k.RapidCurveCount1),
		RapidCurveCount2:            optionalFromPtr(k.RapidCurveCount2),
		RapidCurveCount3:            optionalFromPtr(k.RapidCurveCount3),
		RapidCurveCount4:            optionalFromPtr(k.RapidCurveCount4),
		RapidCurveCount5:   optionalFromPtr(k.RapidCurveCount5),
		RapidCurveMax:      optionalFromPtr(k.RapidCurveMax),
		RapidCurveMaxSpeed: optionalFromPtr(k.RapidCurveMaxSpeed),
	}
}
