package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/pb"
	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/repository"
)

// CarInspectionFilesServer implements the gRPC CarInspectionFilesService
type CarInspectionFilesServer struct {
	pb.UnimplementedCarInspectionFilesServiceServer
	repo *repository.CarInspectionFilesRepository
}

// NewCarInspectionFilesServer creates a new gRPC server
func NewCarInspectionFilesServer(repo *repository.CarInspectionFilesRepository) *CarInspectionFilesServer {
	return &CarInspectionFilesServer{repo: repo}
}

// CreateCarInspectionFile creates a new car inspection file
func (s *CarInspectionFilesServer) CreateCarInspectionFile(ctx context.Context, req *pb.CreateCarInspectionFileRequest) (*pb.CreateCarInspectionFileResponse, error) {
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

	file, err := s.repo.Create(
		ctx,
		req.OrganizationId,
		req.Type,
		req.ElectCertMgNo,
		req.ElectCertPublishdateE,
		req.ElectCertPublishdateY,
		req.ElectCertPublishdateM,
		req.ElectCertPublishdateD,
		req.Created,
		req.Modified,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create car inspection file: %v", err)
	}

	return &pb.CreateCarInspectionFileResponse{
		CarInspectionFile: toProtoCarInspectionFile(file),
	}, nil
}

// GetCarInspectionFile retrieves a car inspection file by UUID
func (s *CarInspectionFilesServer) GetCarInspectionFile(ctx context.Context, req *pb.GetCarInspectionFileRequest) (*pb.GetCarInspectionFileResponse, error) {
	if req.Uuid == "" {
		return nil, status.Error(codes.InvalidArgument, "uuid is required")
	}

	file, err := s.repo.GetByUUID(ctx, req.Uuid)
	if err != nil {
		if errors.Is(err, repository.ErrCarInspectionFileNotFound) {
			return nil, status.Error(codes.NotFound, "car inspection file not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get car inspection file: %v", err)
	}

	return &pb.GetCarInspectionFileResponse{
		CarInspectionFile: toProtoCarInspectionFile(file),
	}, nil
}

// UpdateCarInspectionFile updates an existing car inspection file
func (s *CarInspectionFilesServer) UpdateCarInspectionFile(ctx context.Context, req *pb.UpdateCarInspectionFileRequest) (*pb.UpdateCarInspectionFileResponse, error) {
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
		if errors.Is(err, repository.ErrCarInspectionFileNotFound) {
			return nil, status.Error(codes.NotFound, "car inspection file not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get car inspection file: %v", err)
	}

	// Update mutable fields from the request
	file, err := s.repo.Update(
		ctx,
		req.Uuid,
		req.Type,
		existing.ElectCertMgNo,
		existing.ElectCertPublishdateE,
		existing.ElectCertPublishdateY,
		existing.ElectCertPublishdateM,
		existing.ElectCertPublishdateD,
		req.Modified,
	)
	if err != nil {
		if errors.Is(err, repository.ErrCarInspectionFileNotFound) {
			return nil, status.Error(codes.NotFound, "car inspection file not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update car inspection file: %v", err)
	}

	return &pb.UpdateCarInspectionFileResponse{
		CarInspectionFile: toProtoCarInspectionFile(file),
	}, nil
}

// DeleteCarInspectionFile soft-deletes a car inspection file
func (s *CarInspectionFilesServer) DeleteCarInspectionFile(ctx context.Context, req *pb.DeleteCarInspectionFileRequest) (*pb.DeleteCarInspectionFileResponse, error) {
	if req.Uuid == "" {
		return nil, status.Error(codes.InvalidArgument, "uuid is required")
	}

	// Generate timestamp for soft delete
	deletedTimestamp := req.Uuid // In production, use actual timestamp

	err := s.repo.Delete(ctx, req.Uuid, deletedTimestamp)
	if err != nil {
		if errors.Is(err, repository.ErrCarInspectionFileNotFound) {
			return nil, status.Error(codes.NotFound, "car inspection file not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete car inspection file: %v", err)
	}

	return &pb.DeleteCarInspectionFileResponse{
		Success: true,
	}, nil
}

// ListCarInspectionFiles retrieves all car inspection files with pagination
func (s *CarInspectionFilesServer) ListCarInspectionFiles(ctx context.Context, req *pb.ListCarInspectionFilesRequest) (*pb.ListCarInspectionFilesResponse, error) {
	limit := int(req.PageSize)
	if limit <= 0 {
		limit = 10
	}

	offset := 0
	if req.PageToken != "" {
		// In a real implementation, decode page_token to get offset
	}

	files, err := s.repo.List(ctx, limit+1, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list car inspection files: %v", err)
	}

	var nextPageToken string
	if len(files) > limit {
		files = files[:limit]
		nextPageToken = "next"
	}

	protoFiles := make([]*pb.CarInspectionFile, len(files))
	for i, file := range files {
		protoFiles[i] = toProtoCarInspectionFile(file)
	}

	return &pb.ListCarInspectionFilesResponse{
		CarInspectionFiles: protoFiles,
		NextPageToken:      nextPageToken,
	}, nil
}

// ListCarInspectionFilesByOrganization retrieves car inspection files by organization with pagination
func (s *CarInspectionFilesServer) ListCarInspectionFilesByOrganization(ctx context.Context, req *pb.ListCarInspectionFilesByOrganizationRequest) (*pb.ListCarInspectionFilesByOrganizationResponse, error) {
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

	files, err := s.repo.ListByOrganization(ctx, req.OrganizationId, limit+1, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list car inspection files by organization: %v", err)
	}

	var nextPageToken string
	if len(files) > limit {
		files = files[:limit]
		nextPageToken = "next"
	}

	protoFiles := make([]*pb.CarInspectionFile, len(files))
	for i, file := range files {
		protoFiles[i] = toProtoCarInspectionFile(file)
	}

	return &pb.ListCarInspectionFilesByOrganizationResponse{
		CarInspectionFiles: protoFiles,
		NextPageToken:      nextPageToken,
	}, nil
}

// toProtoCarInspectionFile converts repository model to proto message
func toProtoCarInspectionFile(file *repository.CarInspectionFile) *pb.CarInspectionFile {
	proto := &pb.CarInspectionFile{
		Uuid:                  file.UUID,
		OrganizationId:        file.OrganizationID,
		Type:                  file.Type,
		ElectCertMgNo:         file.ElectCertMgNo,
		ElectCertPublishdateE: file.ElectCertPublishdateE,
		ElectCertPublishdateY: file.ElectCertPublishdateY,
		ElectCertPublishdateM: file.ElectCertPublishdateM,
		ElectCertPublishdateD: file.ElectCertPublishdateD,
		Created:               file.Created,
		Modified:              file.Modified,
	}
	if file.Deleted != nil {
		proto.Deleted = file.Deleted
	}
	return proto
}
