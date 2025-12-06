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

// KudgfulServer implements the gRPC KudgfulService
type KudgfulServer struct {
	pb.UnimplementedKudgfulServiceServer
	repo *repository.KudgfulRepository
}

// NewKudgfulServer creates a new gRPC server
func NewKudgfulServer(repo *repository.KudgfulRepository) *KudgfulServer {
	return &KudgfulServer{repo: repo}
}

// CreateKudgful creates a new kudgful record
func (s *KudgfulServer) CreateKudgful(ctx context.Context, req *pb.CreateKudgfulRequest) (*pb.CreateKudgfulResponse, error) {
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

	kudgful := &repository.Kudgful{
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

	result, err := s.repo.Create(ctx, kudgful)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create kudgful: %v", err)
	}

	return &pb.CreateKudgfulResponse{
		Kudgful: toProtoKudgful(result),
	}, nil
}

// GetKudgful retrieves a kudgful record by UUID
func (s *KudgfulServer) GetKudgful(ctx context.Context, req *pb.GetKudgfulRequest) (*pb.GetKudgfulResponse, error) {
	if req.Uuid == "" {
		return nil, status.Error(codes.InvalidArgument, "uuid is required")
	}

	kudgful, err := s.repo.GetByUUID(ctx, req.Uuid)
	if err != nil {
		if errors.Is(err, repository.ErrKudgfulNotFound) {
			return nil, status.Error(codes.NotFound, "kudgful not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get kudgful: %v", err)
	}

	return &pb.GetKudgfulResponse{
		Kudgful: toProtoKudgful(kudgful),
	}, nil
}

// UpdateKudgful updates an existing kudgful record
func (s *KudgfulServer) UpdateKudgful(ctx context.Context, req *pb.UpdateKudgfulRequest) (*pb.UpdateKudgfulResponse, error) {
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

	kudgful := &repository.Kudgful{
		UUID:             req.Uuid,
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

	result, err := s.repo.Update(ctx, kudgful)
	if err != nil {
		if errors.Is(err, repository.ErrKudgfulNotFound) {
			return nil, status.Error(codes.NotFound, "kudgful not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update kudgful: %v", err)
	}

	return &pb.UpdateKudgfulResponse{
		Kudgful: toProtoKudgful(result),
	}, nil
}

// DeleteKudgful soft-deletes a kudgful record
func (s *KudgfulServer) DeleteKudgful(ctx context.Context, req *pb.DeleteKudgfulRequest) (*pb.DeleteKudgfulResponse, error) {
	if req.Uuid == "" {
		return nil, status.Error(codes.InvalidArgument, "uuid is required")
	}

	deletedAt := time.Now().Format(time.RFC3339)
	err := s.repo.Delete(ctx, req.Uuid, deletedAt)
	if err != nil {
		if errors.Is(err, repository.ErrKudgfulNotFound) {
			return nil, status.Error(codes.NotFound, "kudgful not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete kudgful: %v", err)
	}

	return &pb.DeleteKudgfulResponse{
		Success: true,
	}, nil
}

// ListKudgfuls retrieves kudgful records with pagination
func (s *KudgfulServer) ListKudgfuls(ctx context.Context, req *pb.ListKudgfulsRequest) (*pb.ListKudgfulsResponse, error) {
	limit := int(req.PageSize)
	if limit <= 0 {
		limit = 10
	}

	offset := 0
	if req.PageToken != "" {
		// In production, decode page_token to get offset
	}

	kudgfuls, err := s.repo.List(ctx, limit+1, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list kudgfuls: %v", err)
	}

	var nextPageToken string
	if len(kudgfuls) > limit {
		kudgfuls = kudgfuls[:limit]
		nextPageToken = "next"
	}

	protoKudgfuls := make([]*pb.Kudgful, len(kudgfuls))
	for i, kudgful := range kudgfuls {
		protoKudgfuls[i] = toProtoKudgful(kudgful)
	}

	return &pb.ListKudgfulsResponse{
		Kudgfuls:      protoKudgfuls,
		NextPageToken: nextPageToken,
	}, nil
}

// ListKudgfulsByOrganization retrieves kudgful records for a specific organization with pagination
func (s *KudgfulServer) ListKudgfulsByOrganization(ctx context.Context, req *pb.ListKudgfulsByOrganizationRequest) (*pb.ListKudgfulsByOrganizationResponse, error) {
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

	kudgfuls, err := s.repo.ListByOrganization(ctx, req.OrganizationId, limit+1, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list kudgfuls by organization: %v", err)
	}

	var nextPageToken string
	if len(kudgfuls) > limit {
		kudgfuls = kudgfuls[:limit]
		nextPageToken = "next"
	}

	protoKudgfuls := make([]*pb.Kudgful, len(kudgfuls))
	for i, kudgful := range kudgfuls {
		protoKudgfuls[i] = toProtoKudgful(kudgful)
	}

	return &pb.ListKudgfulsByOrganizationResponse{
		Kudgfuls:      protoKudgfuls,
		NextPageToken: nextPageToken,
	}, nil
}

// toProtoKudgful converts repository model to proto message
func toProtoKudgful(k *repository.Kudgful) *pb.Kudgful {
	return &pb.Kudgful{
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
