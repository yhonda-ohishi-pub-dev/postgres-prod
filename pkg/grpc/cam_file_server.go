package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/pb"
	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/repository"
)

// CamFileServer implements the gRPC CamFileService
type CamFileServer struct {
	pb.UnimplementedCamFileServiceServer
	repo *repository.CamFileRepository
}

// NewCamFileServer creates a new gRPC server
func NewCamFileServer(repo *repository.CamFileRepository) *CamFileServer {
	return &CamFileServer{repo: repo}
}

// CreateCamFile creates a new cam file
func (s *CamFileServer) CreateCamFile(ctx context.Context, req *pb.CreateCamFileRequest) (*pb.CreateCamFileResponse, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}
	if req.Date == "" {
		return nil, status.Error(codes.InvalidArgument, "date is required")
	}
	if req.Hour == "" {
		return nil, status.Error(codes.InvalidArgument, "hour is required")
	}
	if req.Type == "" {
		return nil, status.Error(codes.InvalidArgument, "type is required")
	}
	if req.Cam == "" {
		return nil, status.Error(codes.InvalidArgument, "cam is required")
	}

	var flickrID *string
	if req.FlickrId != nil {
		flickrID = req.FlickrId
	}

	camFile, err := s.repo.Create(ctx, req.Name, req.OrganizationId, req.Date, req.Hour, req.Type, req.Cam, flickrID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create cam file: %v", err)
	}

	return &pb.CreateCamFileResponse{
		CamFile: toProtoCamFile(camFile),
	}, nil
}

// GetCamFile retrieves a cam file by composite key (name, organization_id)
func (s *CamFileServer) GetCamFile(ctx context.Context, req *pb.GetCamFileRequest) (*pb.GetCamFileResponse, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}

	camFile, err := s.repo.GetByNameAndOrg(ctx, req.Name, req.OrganizationId)
	if err != nil {
		if errors.Is(err, repository.ErrCamFileNotFound) {
			return nil, status.Error(codes.NotFound, "cam file not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get cam file: %v", err)
	}

	return &pb.GetCamFileResponse{
		CamFile: toProtoCamFile(camFile),
	}, nil
}

// UpdateCamFile updates an existing cam file
func (s *CamFileServer) UpdateCamFile(ctx context.Context, req *pb.UpdateCamFileRequest) (*pb.UpdateCamFileResponse, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}
	if req.Date == "" {
		return nil, status.Error(codes.InvalidArgument, "date is required")
	}
	if req.Hour == "" {
		return nil, status.Error(codes.InvalidArgument, "hour is required")
	}
	if req.Type == "" {
		return nil, status.Error(codes.InvalidArgument, "type is required")
	}
	if req.Cam == "" {
		return nil, status.Error(codes.InvalidArgument, "cam is required")
	}

	var flickrID *string
	if req.FlickrId != nil {
		flickrID = req.FlickrId
	}

	camFile, err := s.repo.Update(ctx, req.Name, req.OrganizationId, req.Date, req.Hour, req.Type, req.Cam, flickrID)
	if err != nil {
		if errors.Is(err, repository.ErrCamFileNotFound) {
			return nil, status.Error(codes.NotFound, "cam file not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update cam file: %v", err)
	}

	return &pb.UpdateCamFileResponse{
		CamFile: toProtoCamFile(camFile),
	}, nil
}

// DeleteCamFile deletes a cam file
func (s *CamFileServer) DeleteCamFile(ctx context.Context, req *pb.DeleteCamFileRequest) (*pb.DeleteCamFileResponse, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}

	err := s.repo.Delete(ctx, req.Name, req.OrganizationId)
	if err != nil {
		if errors.Is(err, repository.ErrCamFileNotFound) {
			return nil, status.Error(codes.NotFound, "cam file not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete cam file: %v", err)
	}

	return &pb.DeleteCamFileResponse{
		Success: true,
	}, nil
}

// ListCamFiles retrieves cam files with pagination (not implemented - use ListCamFilesByOrganization)
func (s *CamFileServer) ListCamFiles(ctx context.Context, req *pb.ListCamFilesRequest) (*pb.ListCamFilesResponse, error) {
	return nil, status.Error(codes.Unimplemented, "use ListCamFilesByOrganization instead")
}

// ListCamFilesByOrganization retrieves cam files for a specific organization with pagination
func (s *CamFileServer) ListCamFilesByOrganization(ctx context.Context, req *pb.ListCamFilesByOrganizationRequest) (*pb.ListCamFilesByOrganizationResponse, error) {
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

	camFiles, err := s.repo.ListByOrganization(ctx, req.OrganizationId, limit+1, offset) // +1 to check if there's next page
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list cam files: %v", err)
	}

	var nextPageToken string
	if len(camFiles) > limit {
		camFiles = camFiles[:limit]
		nextPageToken = "next" // In production, encode the offset
	}

	protoCamFiles := make([]*pb.CamFile, len(camFiles))
	for i, camFile := range camFiles {
		protoCamFiles[i] = toProtoCamFile(camFile)
	}

	return &pb.ListCamFilesByOrganizationResponse{
		CamFiles:      protoCamFiles,
		NextPageToken: nextPageToken,
	}, nil
}

// toProtoCamFile converts repository model to proto message
func toProtoCamFile(camFile *repository.CamFile) *pb.CamFile {
	proto := &pb.CamFile{
		Name:           camFile.Name,
		OrganizationId: camFile.OrganizationID,
		Date:           camFile.Date,
		Hour:           camFile.Hour,
		Type:           camFile.Type,
		Cam:            camFile.Cam,
	}
	if camFile.FlickrID != nil {
		proto.FlickrId = camFile.FlickrID
	}
	return proto
}
