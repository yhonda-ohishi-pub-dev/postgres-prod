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

// KudguriServer implements the gRPC KudguriService
type KudguriServer struct {
	pb.UnimplementedKudguriServiceServer
	repo *repository.KudguriRepository
}

// NewKudguriServer creates a new gRPC server
func NewKudguriServer(repo *repository.KudguriRepository) *KudguriServer {
	return &KudguriServer{repo: repo}
}

// CreateKudguri creates a new kudguri record
func (s *KudguriServer) CreateKudguri(ctx context.Context, req *pb.CreateKudguriRequest) (*pb.CreateKudguriResponse, error) {
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}
	if req.Hash == "" {
		return nil, status.Error(codes.InvalidArgument, "hash is required")
	}
	if req.UnkouNo == nil || *req.UnkouNo == "" {
		return nil, status.Error(codes.InvalidArgument, "unkou_no is required")
	}
	if req.KudguriUuid == nil || *req.KudguriUuid == "" {
		return nil, status.Error(codes.InvalidArgument, "kudguri_uuid is required")
	}
	if req.TargetDriverType == "" {
		return nil, status.Error(codes.InvalidArgument, "target_driver_type is required")
	}

	kudguri := &repository.Kudguri{
		OrganizationID:   req.OrganizationId,
		Hash:             req.Hash,
		Created:          req.Created,
		Deleted:          ptrFromOptional(req.Deleted),
		UnkouNo:          *req.UnkouNo,
		KudguriUuid:      *req.KudguriUuid,
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

	result, err := s.repo.Create(ctx, kudguri)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create kudguri: %v", err)
	}

	return &pb.CreateKudguriResponse{
		Kudguri: toProtoKudguri(result),
	}, nil
}

// GetKudguri retrieves a kudguri record by UUID
func (s *KudguriServer) GetKudguri(ctx context.Context, req *pb.GetKudguriRequest) (*pb.GetKudguriResponse, error) {
	if req.Uuid == "" {
		return nil, status.Error(codes.InvalidArgument, "uuid is required")
	}

	kudguri, err := s.repo.GetByUUID(ctx, req.Uuid)
	if err != nil {
		if errors.Is(err, repository.ErrKudguriNotFound) {
			return nil, status.Error(codes.NotFound, "kudguri not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get kudguri: %v", err)
	}

	return &pb.GetKudguriResponse{
		Kudguri: toProtoKudguri(kudguri),
	}, nil
}

// UpdateKudguri updates an existing kudguri record
func (s *KudguriServer) UpdateKudguri(ctx context.Context, req *pb.UpdateKudguriRequest) (*pb.UpdateKudguriResponse, error) {
	if req.Uuid == "" {
		return nil, status.Error(codes.InvalidArgument, "uuid is required")
	}
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}
	if req.Hash == "" {
		return nil, status.Error(codes.InvalidArgument, "hash is required")
	}
	if req.UnkouNo == nil || *req.UnkouNo == "" {
		return nil, status.Error(codes.InvalidArgument, "unkou_no is required")
	}
	if req.KudguriUuid == nil || *req.KudguriUuid == "" {
		return nil, status.Error(codes.InvalidArgument, "kudguri_uuid is required")
	}
	if req.TargetDriverType == "" {
		return nil, status.Error(codes.InvalidArgument, "target_driver_type is required")
	}

	kudguri := &repository.Kudguri{
		UUID:             req.Uuid,
		OrganizationID:   req.OrganizationId,
		Hash:             req.Hash,
		Created:          req.Created,
		Deleted:          ptrFromOptional(req.Deleted),
		UnkouNo:          *req.UnkouNo,
		KudguriUuid:      *req.KudguriUuid,
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

	result, err := s.repo.Update(ctx, kudguri)
	if err != nil {
		if errors.Is(err, repository.ErrKudguriNotFound) {
			return nil, status.Error(codes.NotFound, "kudguri not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update kudguri: %v", err)
	}

	return &pb.UpdateKudguriResponse{
		Kudguri: toProtoKudguri(result),
	}, nil
}

// DeleteKudguri soft-deletes a kudguri record
func (s *KudguriServer) DeleteKudguri(ctx context.Context, req *pb.DeleteKudguriRequest) (*pb.DeleteKudguriResponse, error) {
	if req.Uuid == "" {
		return nil, status.Error(codes.InvalidArgument, "uuid is required")
	}

	deletedAt := time.Now().Format(time.RFC3339)
	err := s.repo.Delete(ctx, req.Uuid, deletedAt)
	if err != nil {
		if errors.Is(err, repository.ErrKudguriNotFound) {
			return nil, status.Error(codes.NotFound, "kudguri not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete kudguri: %v", err)
	}

	return &pb.DeleteKudguriResponse{
		Success: true,
	}, nil
}

// ListKudguris retrieves kudguri records with pagination
func (s *KudguriServer) ListKudguris(ctx context.Context, req *pb.ListKudgurisRequest) (*pb.ListKudgurisResponse, error) {
	limit := int(req.PageSize)
	if limit <= 0 {
		limit = 10
	}

	offset := 0
	if req.PageToken != "" {
		// In production, decode page_token to get offset
	}

	kudguris, err := s.repo.List(ctx, limit+1, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list kudguris: %v", err)
	}

	var nextPageToken string
	if len(kudguris) > limit {
		kudguris = kudguris[:limit]
		nextPageToken = "next"
	}

	protoKudguris := make([]*pb.Kudguri, len(kudguris))
	for i, kudguri := range kudguris {
		protoKudguris[i] = toProtoKudguri(kudguri)
	}

	return &pb.ListKudgurisResponse{
		Kudguris:      protoKudguris,
		NextPageToken: nextPageToken,
	}, nil
}

// ListKudgurisByOrganization retrieves kudguri records for a specific organization with pagination
func (s *KudguriServer) ListKudgurisByOrganization(ctx context.Context, req *pb.ListKudgurisByOrganizationRequest) (*pb.ListKudgurisByOrganizationResponse, error) {
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

	kudguris, err := s.repo.ListByOrganization(ctx, req.OrganizationId, limit+1, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list kudguris by organization: %v", err)
	}

	var nextPageToken string
	if len(kudguris) > limit {
		kudguris = kudguris[:limit]
		nextPageToken = "next"
	}

	protoKudguris := make([]*pb.Kudguri, len(kudguris))
	for i, kudguri := range kudguris {
		protoKudguris[i] = toProtoKudguri(kudguri)
	}

	return &pb.ListKudgurisByOrganizationResponse{
		Kudguris:      protoKudguris,
		NextPageToken: nextPageToken,
	}, nil
}

// toProtoKudguri converts repository model to proto message
func toProtoKudguri(k *repository.Kudguri) *pb.Kudguri {
	unkou_no := k.UnkouNo
	kudguri_uuid := k.KudguriUuid
	return &pb.Kudguri{
		Uuid:             k.UUID,
		OrganizationId:   k.OrganizationID,
		Hash:             k.Hash,
		Created:          k.Created,
		Deleted:          optionalFromPtr(k.Deleted),
		UnkouNo:          &unkou_no,
		KudguriUuid:      &kudguri_uuid,
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
