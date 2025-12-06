package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/pb"
	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/repository"
)

// FileServer implements the gRPC FileService
type FileServer struct {
	pb.UnimplementedFileServiceServer
	repo *repository.FileRepository
}

// NewFileServer creates a new gRPC server
func NewFileServer(repo *repository.FileRepository) *FileServer {
	return &FileServer{repo: repo}
}

// CreateFile creates a new file
func (s *FileServer) CreateFile(ctx context.Context, req *pb.CreateFileRequest) (*pb.CreateFileResponse, error) {
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}
	if req.Filename == "" {
		return nil, status.Error(codes.InvalidArgument, "filename is required")
	}
	if req.Created == "" {
		return nil, status.Error(codes.InvalidArgument, "created is required")
	}
	if req.Type == "" {
		return nil, status.Error(codes.InvalidArgument, "type is required")
	}

	var blob *string
	if req.Blob != nil {
		blob = req.Blob
	}

	file, err := s.repo.Create(ctx, req.OrganizationId, req.Filename, req.Created, req.Type, blob)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create file: %v", err)
	}

	return &pb.CreateFileResponse{
		File: toProtoFile(file),
	}, nil
}

// GetFile retrieves a file by UUID
func (s *FileServer) GetFile(ctx context.Context, req *pb.GetFileRequest) (*pb.GetFileResponse, error) {
	if req.Uuid == "" {
		return nil, status.Error(codes.InvalidArgument, "uuid is required")
	}

	file, err := s.repo.GetByUUID(ctx, req.Uuid)
	if err != nil {
		if errors.Is(err, repository.ErrFileNotFound) {
			return nil, status.Error(codes.NotFound, "file not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get file: %v", err)
	}

	return &pb.GetFileResponse{
		File: toProtoFile(file),
	}, nil
}

// UpdateFile updates an existing file
func (s *FileServer) UpdateFile(ctx context.Context, req *pb.UpdateFileRequest) (*pb.UpdateFileResponse, error) {
	if req.Uuid == "" {
		return nil, status.Error(codes.InvalidArgument, "uuid is required")
	}
	if req.Filename == "" {
		return nil, status.Error(codes.InvalidArgument, "filename is required")
	}
	if req.Type == "" {
		return nil, status.Error(codes.InvalidArgument, "type is required")
	}

	var blob *string
	if req.Blob != nil {
		blob = req.Blob
	}

	file, err := s.repo.Update(ctx, req.Uuid, req.Filename, req.Type, blob)
	if err != nil {
		if errors.Is(err, repository.ErrFileNotFound) {
			return nil, status.Error(codes.NotFound, "file not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update file: %v", err)
	}

	return &pb.UpdateFileResponse{
		File: toProtoFile(file),
	}, nil
}

// DeleteFile soft-deletes a file
func (s *FileServer) DeleteFile(ctx context.Context, req *pb.DeleteFileRequest) (*pb.DeleteFileResponse, error) {
	if req.Uuid == "" {
		return nil, status.Error(codes.InvalidArgument, "uuid is required")
	}
	if req.DeletedTimestamp == "" {
		return nil, status.Error(codes.InvalidArgument, "deleted_timestamp is required")
	}

	err := s.repo.Delete(ctx, req.Uuid, req.DeletedTimestamp)
	if err != nil {
		if errors.Is(err, repository.ErrFileNotFound) {
			return nil, status.Error(codes.NotFound, "file not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete file: %v", err)
	}

	return &pb.DeleteFileResponse{
		Success: true,
	}, nil
}

// ListFiles retrieves files with pagination
func (s *FileServer) ListFiles(ctx context.Context, req *pb.ListFilesRequest) (*pb.ListFilesResponse, error) {
	limit := int(req.PageSize)
	if limit <= 0 {
		limit = 10
	}

	// Simple offset-based pagination using page_token as offset string
	offset := 0
	if req.PageToken != "" {
		// In a real implementation, decode page_token to get offset
		// For simplicity, we use 0 for empty token
	}

	files, err := s.repo.List(ctx, limit+1, offset) // +1 to check if there's next page
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list files: %v", err)
	}

	var nextPageToken string
	if len(files) > limit {
		files = files[:limit]
		nextPageToken = "next" // In production, encode the offset
	}

	protoFiles := make([]*pb.File, len(files))
	for i, file := range files {
		protoFiles[i] = toProtoFile(file)
	}

	return &pb.ListFilesResponse{
		Files:         protoFiles,
		NextPageToken: nextPageToken,
	}, nil
}

// ListFilesByOrganization retrieves files by organization with pagination
func (s *FileServer) ListFilesByOrganization(ctx context.Context, req *pb.ListFilesByOrganizationRequest) (*pb.ListFilesByOrganizationResponse, error) {
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}

	limit := int(req.PageSize)
	if limit <= 0 {
		limit = 10
	}

	// Simple offset-based pagination using page_token as offset string
	offset := 0
	if req.PageToken != "" {
		// In a real implementation, decode page_token to get offset
		// For simplicity, we use 0 for empty token
	}

	files, err := s.repo.ListByOrganization(ctx, req.OrganizationId, limit+1, offset) // +1 to check if there's next page
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list files by organization: %v", err)
	}

	var nextPageToken string
	if len(files) > limit {
		files = files[:limit]
		nextPageToken = "next" // In production, encode the offset
	}

	protoFiles := make([]*pb.File, len(files))
	for i, file := range files {
		protoFiles[i] = toProtoFile(file)
	}

	return &pb.ListFilesByOrganizationResponse{
		Files:         protoFiles,
		NextPageToken: nextPageToken,
	}, nil
}

// toProtoFile converts repository model to proto message
func toProtoFile(file *repository.File) *pb.File {
	proto := &pb.File{
		Uuid:           file.UUID,
		OrganizationId: file.OrganizationID,
		Filename:       file.Filename,
		Created:        file.Created,
		Deleted:        file.Deleted,
		Type:           file.Type,
	}
	if file.Blob != nil {
		proto.Blob = file.Blob
	}
	return proto
}
