package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/pb"
	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/repository"
)

// CamFileExeServer implements the gRPC CamFileExeService
type CamFileExeServer struct {
	pb.UnimplementedCamFileExeServiceServer
	repo *repository.CamFileExeRepository
}

// NewCamFileExeServer creates a new gRPC server
func NewCamFileExeServer(repo *repository.CamFileExeRepository) *CamFileExeServer {
	return &CamFileExeServer{repo: repo}
}

// CreateCamFileExe creates a new cam file exe
func (s *CamFileExeServer) CreateCamFileExe(ctx context.Context, req *pb.CreateCamFileExeRequest) (*pb.CreateCamFileExeResponse, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}
	if req.Cam == "" {
		return nil, status.Error(codes.InvalidArgument, "cam is required")
	}
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}

	camFileExe, err := s.repo.Create(ctx, req.Name, req.Cam, req.OrganizationId, req.Stage)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create cam file exe: %v", err)
	}

	return &pb.CreateCamFileExeResponse{
		CamFileExe: toProtoCamFileExe(camFileExe),
	}, nil
}

// GetCamFileExe retrieves a cam file exe by composite key (name, cam, organization_id)
func (s *CamFileExeServer) GetCamFileExe(ctx context.Context, req *pb.GetCamFileExeRequest) (*pb.GetCamFileExeResponse, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}
	if req.Cam == "" {
		return nil, status.Error(codes.InvalidArgument, "cam is required")
	}
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}

	camFileExe, err := s.repo.GetByKey(ctx, req.Name, req.Cam, req.OrganizationId)
	if err != nil {
		if errors.Is(err, repository.ErrCamFileExeNotFound) {
			return nil, status.Error(codes.NotFound, "cam file exe not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get cam file exe: %v", err)
	}

	return &pb.GetCamFileExeResponse{
		CamFileExe: toProtoCamFileExe(camFileExe),
	}, nil
}

// UpdateCamFileExe updates an existing cam file exe
func (s *CamFileExeServer) UpdateCamFileExe(ctx context.Context, req *pb.UpdateCamFileExeRequest) (*pb.UpdateCamFileExeResponse, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}
	if req.Cam == "" {
		return nil, status.Error(codes.InvalidArgument, "cam is required")
	}
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}

	camFileExe, err := s.repo.Update(ctx, req.Name, req.Cam, req.OrganizationId, req.Stage)
	if err != nil {
		if errors.Is(err, repository.ErrCamFileExeNotFound) {
			return nil, status.Error(codes.NotFound, "cam file exe not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update cam file exe: %v", err)
	}

	return &pb.UpdateCamFileExeResponse{
		CamFileExe: toProtoCamFileExe(camFileExe),
	}, nil
}

// DeleteCamFileExe deletes a cam file exe
func (s *CamFileExeServer) DeleteCamFileExe(ctx context.Context, req *pb.DeleteCamFileExeRequest) (*pb.DeleteCamFileExeResponse, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}
	if req.Cam == "" {
		return nil, status.Error(codes.InvalidArgument, "cam is required")
	}
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}

	err := s.repo.Delete(ctx, req.Name, req.Cam, req.OrganizationId)
	if err != nil {
		if errors.Is(err, repository.ErrCamFileExeNotFound) {
			return nil, status.Error(codes.NotFound, "cam file exe not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete cam file exe: %v", err)
	}

	return &pb.DeleteCamFileExeResponse{
		Success: true,
	}, nil
}

// ListCamFileExes retrieves cam file exes with pagination (not implemented - use ListCamFileExesByOrganization)
func (s *CamFileExeServer) ListCamFileExes(ctx context.Context, req *pb.ListCamFileExesRequest) (*pb.ListCamFileExesResponse, error) {
	return nil, status.Error(codes.Unimplemented, "use ListCamFileExesByOrganization instead")
}

// ListCamFileExesByOrganization retrieves cam file exes for a specific organization with pagination
func (s *CamFileExeServer) ListCamFileExesByOrganization(ctx context.Context, req *pb.ListCamFileExesByOrganizationRequest) (*pb.ListCamFileExesByOrganizationResponse, error) {
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

	camFileExes, err := s.repo.ListByOrganization(ctx, req.OrganizationId, limit+1, offset) // +1 to check if there's next page
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list cam file exes: %v", err)
	}

	var nextPageToken string
	if len(camFileExes) > limit {
		camFileExes = camFileExes[:limit]
		nextPageToken = "next" // In production, encode the offset
	}

	protoCamFileExes := make([]*pb.CamFileExe, len(camFileExes))
	for i, camFileExe := range camFileExes {
		protoCamFileExes[i] = toProtoCamFileExe(camFileExe)
	}

	return &pb.ListCamFileExesByOrganizationResponse{
		CamFileExes:   protoCamFileExes,
		NextPageToken: nextPageToken,
	}, nil
}

// toProtoCamFileExe converts repository model to proto message
func toProtoCamFileExe(camFileExe *repository.CamFileExe) *pb.CamFileExe {
	return &pb.CamFileExe{
		Name:           camFileExe.Name,
		Cam:            camFileExe.Cam,
		OrganizationId: camFileExe.OrganizationID,
		Stage:          camFileExe.Stage,
	}
}
