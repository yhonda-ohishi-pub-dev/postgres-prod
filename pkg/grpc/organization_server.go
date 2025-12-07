package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/pb"
	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/repository"
)

// OrganizationServer implements the gRPC OrganizationService
type OrganizationServer struct {
	pb.UnimplementedOrganizationServiceServer
	repo *repository.OrganizationRepository
}

// NewOrganizationServer creates a new gRPC server
func NewOrganizationServer(repo *repository.OrganizationRepository) *OrganizationServer {
	return &OrganizationServer{repo: repo}
}

// CreateOrganization creates a new organization and links it to the current user as owner.
// Requires JWT authentication to identify the user.
// slug is auto-generated as UUID to avoid duplicate key errors.
func (s *OrganizationServer) CreateOrganization(ctx context.Context, req *pb.CreateOrganizationRequest) (*pb.CreateOrganizationResponse, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	// Get user ID from JWT context
	userID, ok := GetUserIDFromContext(ctx)
	if !ok {
		// Fallback to old behavior if no JWT (for backwards compatibility or testing)
		org, err := s.repo.Create(ctx, req.Name)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to create organization: %v", err)
		}
		return &pb.CreateOrganizationResponse{
			Organization: toProtoOrganization(org),
		}, nil
	}

	// Create organization with owner link in a transaction
	result, err := s.repo.CreateWithOwner(ctx, req.Name, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create organization: %v", err)
	}

	return &pb.CreateOrganizationResponse{
		Organization: toProtoOrganization(result.Organization),
	}, nil
}

// GetOrganization retrieves an organization by ID
func (s *OrganizationServer) GetOrganization(ctx context.Context, req *pb.GetOrganizationRequest) (*pb.GetOrganizationResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	org, err := s.repo.GetByID(ctx, req.Id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "organization not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get organization: %v", err)
	}

	return &pb.GetOrganizationResponse{
		Organization: toProtoOrganization(org),
	}, nil
}

// UpdateOrganization updates an existing organization
func (s *OrganizationServer) UpdateOrganization(ctx context.Context, req *pb.UpdateOrganizationRequest) (*pb.UpdateOrganizationResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}
	if req.Slug == "" {
		return nil, status.Error(codes.InvalidArgument, "slug is required")
	}

	org, err := s.repo.Update(ctx, req.Id, req.Name, req.Slug)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "organization not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update organization: %v", err)
	}

	return &pb.UpdateOrganizationResponse{
		Organization: toProtoOrganization(org),
	}, nil
}

// DeleteOrganization soft-deletes an organization
func (s *OrganizationServer) DeleteOrganization(ctx context.Context, req *pb.DeleteOrganizationRequest) (*pb.DeleteOrganizationResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	err := s.repo.Delete(ctx, req.Id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "organization not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete organization: %v", err)
	}

	return &pb.DeleteOrganizationResponse{
		Success: true,
	}, nil
}

// ListOrganizations retrieves organizations with pagination
func (s *OrganizationServer) ListOrganizations(ctx context.Context, req *pb.ListOrganizationsRequest) (*pb.ListOrganizationsResponse, error) {
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

	orgs, err := s.repo.List(ctx, limit+1, offset) // +1 to check if there's next page
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list organizations: %v", err)
	}

	var nextPageToken string
	if len(orgs) > limit {
		orgs = orgs[:limit]
		nextPageToken = "next" // In production, encode the offset
	}

	protoOrgs := make([]*pb.Organization, len(orgs))
	for i, org := range orgs {
		protoOrgs[i] = toProtoOrganization(org)
	}

	return &pb.ListOrganizationsResponse{
		Organizations: protoOrgs,
		NextPageToken: nextPageToken,
	}, nil
}

// toProtoOrganization converts repository model to proto message
func toProtoOrganization(org *repository.Organization) *pb.Organization {
	proto := &pb.Organization{
		Id:        org.ID,
		Name:      org.Name,
		Slug:      org.Slug,
		CreatedAt: timestamppb.New(org.CreatedAt),
		UpdatedAt: timestamppb.New(org.UpdatedAt),
	}
	if org.DeletedAt != nil {
		proto.DeletedAt = timestamppb.New(*org.DeletedAt)
	}
	return proto
}
