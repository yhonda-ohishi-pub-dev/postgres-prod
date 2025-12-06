package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/pb"
	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/repository"
)

// CarInspectionFilesAServer implements the gRPC CarInspectionFilesAService
type CarInspectionFilesAServer struct {
	pb.UnimplementedCarInspectionFilesAServiceServer
	repo *repository.CarInspectionFilesARepository
}

// NewCarInspectionFilesAServer creates a new gRPC server
func NewCarInspectionFilesAServer(repo *repository.CarInspectionFilesARepository) *CarInspectionFilesAServer {
	return &CarInspectionFilesAServer{repo: repo}
}

// CreateCarInspectionFilesA creates a new car inspection files A record
func (s *CarInspectionFilesAServer) CreateCarInspectionFilesA(ctx context.Context, req *pb.CreateCarInspectionFilesARequest) (*pb.CreateCarInspectionFilesAResponse, error) {
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}
	if req.Type == "" {
		return nil, status.Error(codes.InvalidArgument, "type is required")
	}
	if req.Created == "" {
		return nil, status.Error(codes.InvalidArgument, "created is required")
	}
	if req.Modified == "" {
		return nil, status.Error(codes.InvalidArgument, "modified is required")
	}

	record := &repository.CarInspectionFilesA{
		OrganizationID: req.OrganizationId,
		Type:           req.Type,
		ElectCertMgNo:  req.ElectCertMgNo,
		GrantdateE:     req.GrantdateE,
		GrantdateY:     req.GrantdateY,
		GrantdateM:     req.GrantdateM,
		GrantdateD:     req.GrantdateD,
		Created:        req.Created,
		Modified:       req.Modified,
	}

	result, err := s.repo.Create(ctx, record)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create car inspection files A: %v", err)
	}

	return &pb.CreateCarInspectionFilesAResponse{
		CarInspectionFilesA: toProtoCarInspectionFilesA(result),
	}, nil
}

// GetCarInspectionFilesA retrieves a car inspection files A record by UUID
func (s *CarInspectionFilesAServer) GetCarInspectionFilesA(ctx context.Context, req *pb.GetCarInspectionFilesARequest) (*pb.GetCarInspectionFilesAResponse, error) {
	if req.Uuid == "" {
		return nil, status.Error(codes.InvalidArgument, "uuid is required")
	}

	record, err := s.repo.GetByUUID(ctx, req.Uuid)
	if err != nil {
		if errors.Is(err, repository.ErrCarInspectionFilesANotFound) {
			return nil, status.Error(codes.NotFound, "car inspection files A not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get car inspection files A: %v", err)
	}

	return &pb.GetCarInspectionFilesAResponse{
		CarInspectionFilesA: toProtoCarInspectionFilesA(record),
	}, nil
}

// UpdateCarInspectionFilesA updates an existing car inspection files A record
func (s *CarInspectionFilesAServer) UpdateCarInspectionFilesA(ctx context.Context, req *pb.UpdateCarInspectionFilesARequest) (*pb.UpdateCarInspectionFilesAResponse, error) {
	if req.Uuid == "" {
		return nil, status.Error(codes.InvalidArgument, "uuid is required")
	}
	if req.Type == "" {
		return nil, status.Error(codes.InvalidArgument, "type is required")
	}
	if req.Modified == "" {
		return nil, status.Error(codes.InvalidArgument, "modified is required")
	}

	// Get the existing record first
	existing, err := s.repo.GetByUUID(ctx, req.Uuid)
	if err != nil {
		if errors.Is(err, repository.ErrCarInspectionFilesANotFound) {
			return nil, status.Error(codes.NotFound, "car inspection files A not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get car inspection files A: %v", err)
	}

	// Update mutable fields
	existing.Type = req.Type
	existing.Modified = req.Modified

	result, err := s.repo.Update(ctx, existing)
	if err != nil {
		if errors.Is(err, repository.ErrCarInspectionFilesANotFound) {
			return nil, status.Error(codes.NotFound, "car inspection files A not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update car inspection files A: %v", err)
	}

	return &pb.UpdateCarInspectionFilesAResponse{
		CarInspectionFilesA: toProtoCarInspectionFilesA(result),
	}, nil
}

// DeleteCarInspectionFilesA soft-deletes a car inspection files A record
func (s *CarInspectionFilesAServer) DeleteCarInspectionFilesA(ctx context.Context, req *pb.DeleteCarInspectionFilesARequest) (*pb.DeleteCarInspectionFilesAResponse, error) {
	if req.Uuid == "" {
		return nil, status.Error(codes.InvalidArgument, "uuid is required")
	}

	// Generate timestamp for soft delete (in production, use actual timestamp)
	deletedTime := req.Uuid

	err := s.repo.Delete(ctx, req.Uuid, deletedTime)
	if err != nil {
		if errors.Is(err, repository.ErrCarInspectionFilesANotFound) {
			return nil, status.Error(codes.NotFound, "car inspection files A not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete car inspection files A: %v", err)
	}

	return &pb.DeleteCarInspectionFilesAResponse{
		Success: true,
	}, nil
}

// ListCarInspectionFilesAs retrieves all car inspection files A records with pagination
func (s *CarInspectionFilesAServer) ListCarInspectionFilesAs(ctx context.Context, req *pb.ListCarInspectionFilesAsRequest) (*pb.ListCarInspectionFilesAsResponse, error) {
	limit := int(req.PageSize)
	if limit <= 0 {
		limit = 10
	}

	offset := 0
	if req.PageToken != "" {
		// In a real implementation, decode page_token to get offset
	}

	records, err := s.repo.List(ctx, limit+1, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list car inspection files A: %v", err)
	}

	var nextPageToken string
	if len(records) > limit {
		records = records[:limit]
		nextPageToken = "next"
	}

	protoRecords := make([]*pb.CarInspectionFilesA, len(records))
	for i, record := range records {
		protoRecords[i] = toProtoCarInspectionFilesA(record)
	}

	return &pb.ListCarInspectionFilesAsResponse{
		CarInspectionFilesAs: protoRecords,
		NextPageToken:        nextPageToken,
	}, nil
}

// ListCarInspectionFilesAsByOrganization retrieves car inspection files A records by organization with pagination
func (s *CarInspectionFilesAServer) ListCarInspectionFilesAsByOrganization(ctx context.Context, req *pb.ListCarInspectionFilesAsByOrganizationRequest) (*pb.ListCarInspectionFilesAsByOrganizationResponse, error) {
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

	records, err := s.repo.ListByOrganization(ctx, req.OrganizationId, limit+1, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list car inspection files A by organization: %v", err)
	}

	var nextPageToken string
	if len(records) > limit {
		records = records[:limit]
		nextPageToken = "next"
	}

	protoRecords := make([]*pb.CarInspectionFilesA, len(records))
	for i, record := range records {
		protoRecords[i] = toProtoCarInspectionFilesA(record)
	}

	return &pb.ListCarInspectionFilesAsByOrganizationResponse{
		CarInspectionFilesAs: protoRecords,
		NextPageToken:        nextPageToken,
	}, nil
}

// toProtoCarInspectionFilesA converts repository model to proto message
func toProtoCarInspectionFilesA(record *repository.CarInspectionFilesA) *pb.CarInspectionFilesA {
	proto := &pb.CarInspectionFilesA{
		Uuid:           record.UUID,
		OrganizationId: record.OrganizationID,
		Type:           record.Type,
		ElectCertMgNo:  record.ElectCertMgNo,
		GrantdateE:     record.GrantdateE,
		GrantdateY:     record.GrantdateY,
		GrantdateM:     record.GrantdateM,
		GrantdateD:     record.GrantdateD,
		Created:        record.Created,
		Modified:       record.Modified,
	}
	if record.Deleted != nil {
		proto.Deleted = record.Deleted
	}
	return proto
}
