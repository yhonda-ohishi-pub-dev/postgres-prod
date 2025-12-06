package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/pb"
	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/repository"
)

// CarInspectionFilesBServer implements the gRPC CarInspectionFilesBService
type CarInspectionFilesBServer struct {
	pb.UnimplementedCarInspectionFilesBServiceServer
	repo *repository.CarInspectionFilesBRepository
}

// NewCarInspectionFilesBServer creates a new gRPC server
func NewCarInspectionFilesBServer(repo *repository.CarInspectionFilesBRepository) *CarInspectionFilesBServer {
	return &CarInspectionFilesBServer{repo: repo}
}

// CreateCarInspectionFilesB creates a new car inspection files B record
func (s *CarInspectionFilesBServer) CreateCarInspectionFilesB(ctx context.Context, req *pb.CreateCarInspectionFilesBRequest) (*pb.CreateCarInspectionFilesBResponse, error) {
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}
	if req.Type == "" {
		return nil, status.Error(codes.InvalidArgument, "type is required")
	}

	record, err := s.repo.Create(
		ctx,
		req.OrganizationId,
		req.Type,
		req.ElectCertMgNo,
		req.GrantdateE,
		req.GrantdateY,
		req.GrantdateM,
		req.GrantdateD,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create car inspection files B: %v", err)
	}

	return &pb.CreateCarInspectionFilesBResponse{
		CarInspectionFilesB: toProtoCarInspectionFilesB(record),
	}, nil
}

// GetCarInspectionFilesB retrieves a car inspection files B record by UUID
func (s *CarInspectionFilesBServer) GetCarInspectionFilesB(ctx context.Context, req *pb.GetCarInspectionFilesBRequest) (*pb.GetCarInspectionFilesBResponse, error) {
	if req.Uuid == "" {
		return nil, status.Error(codes.InvalidArgument, "uuid is required")
	}

	record, err := s.repo.GetByUUID(ctx, req.Uuid)
	if err != nil {
		if errors.Is(err, repository.ErrCarInspectionFilesBNotFound) {
			return nil, status.Error(codes.NotFound, "car inspection files B not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get car inspection files B: %v", err)
	}

	return &pb.GetCarInspectionFilesBResponse{
		CarInspectionFilesB: toProtoCarInspectionFilesB(record),
	}, nil
}

// UpdateCarInspectionFilesB updates an existing car inspection files B record
func (s *CarInspectionFilesBServer) UpdateCarInspectionFilesB(ctx context.Context, req *pb.UpdateCarInspectionFilesBRequest) (*pb.UpdateCarInspectionFilesBResponse, error) {
	if req.Uuid == "" {
		return nil, status.Error(codes.InvalidArgument, "uuid is required")
	}
	if req.Type == "" {
		return nil, status.Error(codes.InvalidArgument, "type is required")
	}
	if req.Modified == "" {
		return nil, status.Error(codes.InvalidArgument, "modified is required")
	}

	// Get the existing record first to preserve immutable fields
	existing, err := s.repo.GetByUUID(ctx, req.Uuid)
	if err != nil {
		if errors.Is(err, repository.ErrCarInspectionFilesBNotFound) {
			return nil, status.Error(codes.NotFound, "car inspection files B not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get car inspection files B: %v", err)
	}

	// Update with mutable fields (type and modified)
	record, err := s.repo.Update(
		ctx,
		req.Uuid,
		existing.OrganizationID,
		req.Type,
		existing.ElectCertMgNo,
		existing.GrantdateE,
		existing.GrantdateY,
		existing.GrantdateM,
		existing.GrantdateD,
	)
	if err != nil {
		if errors.Is(err, repository.ErrCarInspectionFilesBNotFound) {
			return nil, status.Error(codes.NotFound, "car inspection files B not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update car inspection files B: %v", err)
	}

	return &pb.UpdateCarInspectionFilesBResponse{
		CarInspectionFilesB: toProtoCarInspectionFilesB(record),
	}, nil
}

// DeleteCarInspectionFilesB soft-deletes a car inspection files B record
func (s *CarInspectionFilesBServer) DeleteCarInspectionFilesB(ctx context.Context, req *pb.DeleteCarInspectionFilesBRequest) (*pb.DeleteCarInspectionFilesBResponse, error) {
	if req.Uuid == "" {
		return nil, status.Error(codes.InvalidArgument, "uuid is required")
	}

	err := s.repo.Delete(ctx, req.Uuid)
	if err != nil {
		if errors.Is(err, repository.ErrCarInspectionFilesBNotFound) {
			return nil, status.Error(codes.NotFound, "car inspection files B not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete car inspection files B: %v", err)
	}

	return &pb.DeleteCarInspectionFilesBResponse{
		Success: true,
	}, nil
}

// ListCarInspectionFilesBs retrieves all car inspection files B records with pagination
func (s *CarInspectionFilesBServer) ListCarInspectionFilesBs(ctx context.Context, req *pb.ListCarInspectionFilesBsRequest) (*pb.ListCarInspectionFilesBsResponse, error) {
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
		return nil, status.Errorf(codes.Internal, "failed to list car inspection files B: %v", err)
	}

	var nextPageToken string
	if len(records) > limit {
		records = records[:limit]
		nextPageToken = "next"
	}

	protoRecords := make([]*pb.CarInspectionFilesB, len(records))
	for i, record := range records {
		protoRecords[i] = toProtoCarInspectionFilesB(record)
	}

	return &pb.ListCarInspectionFilesBsResponse{
		CarInspectionFilesBs: protoRecords,
		NextPageToken:        nextPageToken,
	}, nil
}

// ListCarInspectionFilesBsByOrganization retrieves car inspection files B records by organization with pagination
func (s *CarInspectionFilesBServer) ListCarInspectionFilesBsByOrganization(ctx context.Context, req *pb.ListCarInspectionFilesBsByOrganizationRequest) (*pb.ListCarInspectionFilesBsByOrganizationResponse, error) {
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
		return nil, status.Errorf(codes.Internal, "failed to list car inspection files B by organization: %v", err)
	}

	var nextPageToken string
	if len(records) > limit {
		records = records[:limit]
		nextPageToken = "next"
	}

	protoRecords := make([]*pb.CarInspectionFilesB, len(records))
	for i, record := range records {
		protoRecords[i] = toProtoCarInspectionFilesB(record)
	}

	return &pb.ListCarInspectionFilesBsByOrganizationResponse{
		CarInspectionFilesBs: protoRecords,
		NextPageToken:        nextPageToken,
	}, nil
}

// toProtoCarInspectionFilesB converts repository model to proto message
func toProtoCarInspectionFilesB(record *repository.CarInspectionFilesB) *pb.CarInspectionFilesB {
	proto := &pb.CarInspectionFilesB{
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
