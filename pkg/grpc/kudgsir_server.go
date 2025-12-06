package grpc

import (
	"context"
	"errors"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/pb"
	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/repository"
)

// KudgsirServer implements the gRPC KudgsirService
type KudgsirServer struct {
	pb.UnimplementedKudgsirServiceServer
	repo *repository.KudgsirRepository
}

// NewKudgsirServer creates a new gRPC server
func NewKudgsirServer(repo *repository.KudgsirRepository) *KudgsirServer {
	return &KudgsirServer{repo: repo}
}

// CreateKudgsir creates a new kudgsir record
func (s *KudgsirServer) CreateKudgsir(ctx context.Context, req *pb.CreateKudgsirRequest) (*pb.CreateKudgsirResponse, error) {
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

	kudgsir := &repository.Kudgsir{
		OrganizationID:   req.OrganizationId,
		Hash:             req.Hash,
		Created:          req.Created,
		Deleted:          ptrFromOptional(req.Deleted),
		KudguriUuid:      ptrFromOptional(req.KudguriUuid),
		UnkouNo:          ptrFromOptional(req.UnkouNo),
		ReadDate:         ptrFromOptional(req.ReadDate),
		OfficeCd:         ptrFromOptional(req.OfficeCd),
		OfficeName:       ptrFromOptional(req.OfficeName),
		VehicleCd:        ptrFromOptional(req.VehicleCd),
		VehicleName:      ptrFromOptional(req.VehicleName),
		DriverCd1:        ptrFromOptional(req.DriverCd1),
		DriverName1:      ptrFromOptional(req.DriverName1),
		TargetDriverType: req.TargetDriverType,
		TargetDriverCd:   ptrFromOptional(req.TargetDriverCd),
		TargetDriverName: ptrFromOptional(req.TargetDriverName),
		StartDatetime:    ptrFromOptional(req.StartDatetime),
		EndDatetime:      ptrFromOptional(req.EndDatetime),
		EventCd:          ptrFromOptional(req.EventCd),
		EventName:        ptrFromOptional(req.EventName),
		StartMileage:     ptrFromOptional(req.StartMileage),
		EndMileage:       ptrFromOptional(req.EndMileage),
		SectionTime:      ptrFromOptional(req.SectionTime),
		SectionDistance:  ptrFromOptional(req.SectionDistance),
		StartCityCd:      ptrFromOptional(req.StartCityCd),
		StartCityName:    ptrFromOptional(req.StartCityName),
		EndCityCd:        ptrFromOptional(req.EndCityCd),
		EndCityName:      ptrFromOptional(req.EndCityName),
		StartPlaceCd:     ptrFromOptional(req.StartPlaceCd),
		StartPlaceName:   ptrFromOptional(req.StartPlaceName),
		EndPlaceCd:       ptrFromOptional(req.EndPlaceCd),
		EndPlaceName:     ptrFromOptional(req.EndPlaceName),
		StartGpsValid:    ptrFromOptional(req.StartGpsValid),
		StartGpsLat:      ptrFromOptional(req.StartGpsLat),
		StartGpsLng:      ptrFromOptional(req.StartGpsLng),
		EndGpsValid:      ptrFromOptional(req.EndGpsValid),
		EndGpsLat:        ptrFromOptional(req.EndGpsLat),
		EndGpsLng:        ptrFromOptional(req.EndGpsLng),
		OverLimitMax:     ptrFromOptional(req.OverLimitMax),
	}

	result, err := s.repo.Create(ctx, kudgsir)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create kudgsir: %v", err)
	}

	return &pb.CreateKudgsirResponse{
		Kudgsir: toProtoKudgsir(result),
	}, nil
}

// GetKudgsir retrieves a kudgsir record by UUID
func (s *KudgsirServer) GetKudgsir(ctx context.Context, req *pb.GetKudgsirRequest) (*pb.GetKudgsirResponse, error) {
	if req.Uuid == "" {
		return nil, status.Error(codes.InvalidArgument, "uuid is required")
	}

	kudgsir, err := s.repo.GetByUUID(ctx, req.Uuid)
	if err != nil {
		if errors.Is(err, repository.ErrKudgsirNotFound) {
			return nil, status.Error(codes.NotFound, "kudgsir not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get kudgsir: %v", err)
	}

	return &pb.GetKudgsirResponse{
		Kudgsir: toProtoKudgsir(kudgsir),
	}, nil
}

// UpdateKudgsir updates an existing kudgsir record
func (s *KudgsirServer) UpdateKudgsir(ctx context.Context, req *pb.UpdateKudgsirRequest) (*pb.UpdateKudgsirResponse, error) {
	if req.Uuid == "" {
		return nil, status.Error(codes.InvalidArgument, "uuid is required")
	}
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}
	if req.Hash == "" {
		return nil, status.Error(codes.InvalidArgument, "hash is required")
	}
	if req.TargetDriverType == "" {
		return nil, status.Error(codes.InvalidArgument, "target_driver_type is required")
	}

	kudgsir := &repository.Kudgsir{
		UUID:             req.Uuid,
		OrganizationID:   req.OrganizationId,
		Hash:             req.Hash,
		KudguriUuid:      ptrFromOptional(req.KudguriUuid),
		UnkouNo:          ptrFromOptional(req.UnkouNo),
		ReadDate:         ptrFromOptional(req.ReadDate),
		OfficeCd:         ptrFromOptional(req.OfficeCd),
		OfficeName:       ptrFromOptional(req.OfficeName),
		VehicleCd:        ptrFromOptional(req.VehicleCd),
		VehicleName:      ptrFromOptional(req.VehicleName),
		DriverCd1:        ptrFromOptional(req.DriverCd1),
		DriverName1:      ptrFromOptional(req.DriverName1),
		TargetDriverType: req.TargetDriverType,
		TargetDriverCd:   ptrFromOptional(req.TargetDriverCd),
		TargetDriverName: ptrFromOptional(req.TargetDriverName),
		StartDatetime:    ptrFromOptional(req.StartDatetime),
		EndDatetime:      ptrFromOptional(req.EndDatetime),
		EventCd:          ptrFromOptional(req.EventCd),
		EventName:        ptrFromOptional(req.EventName),
		StartMileage:     ptrFromOptional(req.StartMileage),
		EndMileage:       ptrFromOptional(req.EndMileage),
		SectionTime:      ptrFromOptional(req.SectionTime),
		SectionDistance:  ptrFromOptional(req.SectionDistance),
		StartCityCd:      ptrFromOptional(req.StartCityCd),
		StartCityName:    ptrFromOptional(req.StartCityName),
		EndCityCd:        ptrFromOptional(req.EndCityCd),
		EndCityName:      ptrFromOptional(req.EndCityName),
		StartPlaceCd:     ptrFromOptional(req.StartPlaceCd),
		StartPlaceName:   ptrFromOptional(req.StartPlaceName),
		EndPlaceCd:       ptrFromOptional(req.EndPlaceCd),
		EndPlaceName:     ptrFromOptional(req.EndPlaceName),
		StartGpsValid:    ptrFromOptional(req.StartGpsValid),
		StartGpsLat:      ptrFromOptional(req.StartGpsLat),
		StartGpsLng:      ptrFromOptional(req.StartGpsLng),
		EndGpsValid:      ptrFromOptional(req.EndGpsValid),
		EndGpsLat:        ptrFromOptional(req.EndGpsLat),
		EndGpsLng:        ptrFromOptional(req.EndGpsLng),
		OverLimitMax:     ptrFromOptional(req.OverLimitMax),
	}

	result, err := s.repo.Update(ctx, kudgsir)
	if err != nil {
		if errors.Is(err, repository.ErrKudgsirNotFound) {
			return nil, status.Error(codes.NotFound, "kudgsir not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update kudgsir: %v", err)
	}

	return &pb.UpdateKudgsirResponse{
		Kudgsir: toProtoKudgsir(result),
	}, nil
}

// DeleteKudgsir soft-deletes a kudgsir record
func (s *KudgsirServer) DeleteKudgsir(ctx context.Context, req *pb.DeleteKudgsirRequest) (*pb.DeleteKudgsirResponse, error) {
	if req.Uuid == "" {
		return nil, status.Error(codes.InvalidArgument, "uuid is required")
	}

	deletedAt := time.Now().Format(time.RFC3339)
	err := s.repo.Delete(ctx, req.Uuid, deletedAt)
	if err != nil {
		if errors.Is(err, repository.ErrKudgsirNotFound) {
			return nil, status.Error(codes.NotFound, "kudgsir not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete kudgsir: %v", err)
	}

	return &pb.DeleteKudgsirResponse{
		Success: true,
	}, nil
}

// ListKudgsirs retrieves kudgsir records with pagination
func (s *KudgsirServer) ListKudgsirs(ctx context.Context, req *pb.ListKudgsirsRequest) (*pb.ListKudgsirsResponse, error) {
	limit := int(req.PageSize)
	if limit <= 0 {
		limit = 10
	}

	offset := 0
	if req.PageToken != "" {
		// In production, decode page_token to get offset
	}

	kudgsirs, err := s.repo.List(ctx, limit+1, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list kudgsirs: %v", err)
	}

	var nextPageToken string
	if len(kudgsirs) > limit {
		kudgsirs = kudgsirs[:limit]
		nextPageToken = "next"
	}

	protoKudgsirs := make([]*pb.Kudgsir, len(kudgsirs))
	for i, kudgsir := range kudgsirs {
		protoKudgsirs[i] = toProtoKudgsir(kudgsir)
	}

	return &pb.ListKudgsirsResponse{
		Kudgsirs:      protoKudgsirs,
		NextPageToken: nextPageToken,
	}, nil
}

// ListKudgsirsByOrganization retrieves kudgsir records for a specific organization with pagination
func (s *KudgsirServer) ListKudgsirsByOrganization(ctx context.Context, req *pb.ListKudgsirsByOrganizationRequest) (*pb.ListKudgsirsByOrganizationResponse, error) {
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

	kudgsirs, err := s.repo.ListByOrganization(ctx, req.OrganizationId, limit+1, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list kudgsirs by organization: %v", err)
	}

	var nextPageToken string
	if len(kudgsirs) > limit {
		kudgsirs = kudgsirs[:limit]
		nextPageToken = "next"
	}

	protoKudgsirs := make([]*pb.Kudgsir, len(kudgsirs))
	for i, kudgsir := range kudgsirs {
		protoKudgsirs[i] = toProtoKudgsir(kudgsir)
	}

	return &pb.ListKudgsirsByOrganizationResponse{
		Kudgsirs:      protoKudgsirs,
		NextPageToken: nextPageToken,
	}, nil
}

// toProtoKudgsir converts repository model to proto message
func toProtoKudgsir(k *repository.Kudgsir) *pb.Kudgsir {
	return &pb.Kudgsir{
		Uuid:             k.UUID,
		OrganizationId:   k.OrganizationID,
		Hash:             k.Hash,
		Created:          k.Created,
		Deleted:          optionalFromPtr(k.Deleted),
		KudguriUuid:      optionalFromPtr(k.KudguriUuid),
		UnkouNo:          optionalFromPtr(k.UnkouNo),
		ReadDate:         optionalFromPtr(k.ReadDate),
		OfficeCd:         optionalFromPtr(k.OfficeCd),
		OfficeName:       optionalFromPtr(k.OfficeName),
		VehicleCd:        optionalFromPtr(k.VehicleCd),
		VehicleName:      optionalFromPtr(k.VehicleName),
		DriverCd1:        optionalFromPtr(k.DriverCd1),
		DriverName1:      optionalFromPtr(k.DriverName1),
		TargetDriverType: k.TargetDriverType,
		TargetDriverCd:   optionalFromPtr(k.TargetDriverCd),
		TargetDriverName: optionalFromPtr(k.TargetDriverName),
		StartDatetime:    optionalFromPtr(k.StartDatetime),
		EndDatetime:      optionalFromPtr(k.EndDatetime),
		EventCd:          optionalFromPtr(k.EventCd),
		EventName:        optionalFromPtr(k.EventName),
		StartMileage:     optionalFromPtr(k.StartMileage),
		EndMileage:       optionalFromPtr(k.EndMileage),
		SectionTime:      optionalFromPtr(k.SectionTime),
		SectionDistance:  optionalFromPtr(k.SectionDistance),
		StartCityCd:      optionalFromPtr(k.StartCityCd),
		StartCityName:    optionalFromPtr(k.StartCityName),
		EndCityCd:        optionalFromPtr(k.EndCityCd),
		EndCityName:      optionalFromPtr(k.EndCityName),
		StartPlaceCd:     optionalFromPtr(k.StartPlaceCd),
		StartPlaceName:   optionalFromPtr(k.StartPlaceName),
		EndPlaceCd:       optionalFromPtr(k.EndPlaceCd),
		EndPlaceName:     optionalFromPtr(k.EndPlaceName),
		StartGpsValid:    optionalFromPtr(k.StartGpsValid),
		StartGpsLat:      optionalFromPtr(k.StartGpsLat),
		StartGpsLng:      optionalFromPtr(k.StartGpsLng),
		EndGpsValid:      optionalFromPtr(k.EndGpsValid),
		EndGpsLat:        optionalFromPtr(k.EndGpsLat),
		EndGpsLng:        optionalFromPtr(k.EndGpsLng),
		OverLimitMax:     optionalFromPtr(k.OverLimitMax),
	}
}
