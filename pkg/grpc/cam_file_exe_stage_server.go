package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/pb"
	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/repository"
)

// CamFileExeStageServer implements the gRPC CamFileExeStageService
type CamFileExeStageServer struct {
	pb.UnimplementedCamFileExeStageServiceServer
	repo *repository.CamFileExeStageRepository
}

// NewCamFileExeStageServer creates a new gRPC server
func NewCamFileExeStageServer(repo *repository.CamFileExeStageRepository) *CamFileExeStageServer {
	return &CamFileExeStageServer{repo: repo}
}

// CreateCamFileExeStage creates a new cam file exe stage
func (s *CamFileExeStageServer) CreateCamFileExeStage(ctx context.Context, req *pb.CreateCamFileExeStageRequest) (*pb.CreateCamFileExeStageResponse, error) {
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	camStage, err := s.repo.Create(ctx, req.Stage, req.OrganizationId, req.Name)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create cam file exe stage: %v", err)
	}

	return &pb.CreateCamFileExeStageResponse{
		CamFileExeStage: toProtoCamFileExeStage(camStage),
	}, nil
}

// GetCamFileExeStage retrieves a cam file exe stage by composite key (stage, organization_id)
func (s *CamFileExeStageServer) GetCamFileExeStage(ctx context.Context, req *pb.GetCamFileExeStageRequest) (*pb.GetCamFileExeStageResponse, error) {
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}

	camStage, err := s.repo.GetByStageAndOrg(ctx, req.Stage, req.OrganizationId)
	if err != nil {
		if errors.Is(err, repository.ErrCamFileExeStageNotFound) {
			return nil, status.Error(codes.NotFound, "cam file exe stage not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get cam file exe stage: %v", err)
	}

	return &pb.GetCamFileExeStageResponse{
		CamFileExeStage: toProtoCamFileExeStage(camStage),
	}, nil
}

// UpdateCamFileExeStage updates an existing cam file exe stage
func (s *CamFileExeStageServer) UpdateCamFileExeStage(ctx context.Context, req *pb.UpdateCamFileExeStageRequest) (*pb.UpdateCamFileExeStageResponse, error) {
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	camStage, err := s.repo.Update(ctx, req.Stage, req.OrganizationId, req.Name)
	if err != nil {
		if errors.Is(err, repository.ErrCamFileExeStageNotFound) {
			return nil, status.Error(codes.NotFound, "cam file exe stage not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update cam file exe stage: %v", err)
	}

	return &pb.UpdateCamFileExeStageResponse{
		CamFileExeStage: toProtoCamFileExeStage(camStage),
	}, nil
}

// DeleteCamFileExeStage deletes a cam file exe stage
func (s *CamFileExeStageServer) DeleteCamFileExeStage(ctx context.Context, req *pb.DeleteCamFileExeStageRequest) (*pb.DeleteCamFileExeStageResponse, error) {
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}

	err := s.repo.Delete(ctx, req.Stage, req.OrganizationId)
	if err != nil {
		if errors.Is(err, repository.ErrCamFileExeStageNotFound) {
			return nil, status.Error(codes.NotFound, "cam file exe stage not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete cam file exe stage: %v", err)
	}

	return &pb.DeleteCamFileExeStageResponse{
		Success: true,
	}, nil
}

// ListCamFileExeStages retrieves cam file exe stages with pagination (not implemented - use ListCamFileExeStagesByOrganization)
func (s *CamFileExeStageServer) ListCamFileExeStages(ctx context.Context, req *pb.ListCamFileExeStagesRequest) (*pb.ListCamFileExeStagesResponse, error) {
	return nil, status.Error(codes.Unimplemented, "use ListCamFileExeStagesByOrganization instead")
}

// ListCamFileExeStagesByOrganization retrieves cam file exe stages for a specific organization
func (s *CamFileExeStageServer) ListCamFileExeStagesByOrganization(ctx context.Context, req *pb.ListCamFileExeStagesByOrganizationRequest) (*pb.ListCamFileExeStagesByOrganizationResponse, error) {
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}

	// Note: This repository method doesn't support pagination
	// It returns all stages for the organization
	camStages, err := s.repo.ListByOrganization(ctx, req.OrganizationId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list cam file exe stages: %v", err)
	}

	protoCamStages := make([]*pb.CamFileExeStage, len(camStages))
	for i, camStage := range camStages {
		protoCamStages[i] = toProtoCamFileExeStage(camStage)
	}

	return &pb.ListCamFileExeStagesByOrganizationResponse{
		CamFileExeStages: protoCamStages,
		NextPageToken:    "", // No pagination for this endpoint
	}, nil
}

// toProtoCamFileExeStage converts repository model to proto message
func toProtoCamFileExeStage(camStage *repository.CamFileExeStage) *pb.CamFileExeStage {
	return &pb.CamFileExeStage{
		Stage:          camStage.Stage,
		OrganizationId: camStage.OrganizationID,
		Name:           camStage.Name,
	}
}
