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

// UserOrganizationServer implements the gRPC UserOrganizationService
type UserOrganizationServer struct {
	pb.UnimplementedUserOrganizationServiceServer
	repo *repository.UserOrganizationRepository
}

// NewUserOrganizationServer creates a new gRPC server
func NewUserOrganizationServer(repo *repository.UserOrganizationRepository) *UserOrganizationServer {
	return &UserOrganizationServer{repo: repo}
}

// CreateUserOrganization creates a new user organization
func (s *UserOrganizationServer) CreateUserOrganization(ctx context.Context, req *pb.CreateUserOrganizationRequest) (*pb.CreateUserOrganizationResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}

	uo, err := s.repo.Create(ctx, req.UserId, req.OrganizationId, req.Role, req.IsDefault)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user organization: %v", err)
	}

	return &pb.CreateUserOrganizationResponse{
		UserOrganization: toProtoUserOrganization(uo),
	}, nil
}

// GetUserOrganization retrieves a user organization by ID
func (s *UserOrganizationServer) GetUserOrganization(ctx context.Context, req *pb.GetUserOrganizationRequest) (*pb.GetUserOrganizationResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	uo, err := s.repo.GetByID(ctx, req.Id)
	if err != nil {
		if errors.Is(err, repository.ErrUserOrganizationNotFound) {
			return nil, status.Error(codes.NotFound, "user organization not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get user organization: %v", err)
	}

	return &pb.GetUserOrganizationResponse{
		UserOrganization: toProtoUserOrganization(uo),
	}, nil
}

// UpdateUserOrganization updates an existing user organization
func (s *UserOrganizationServer) UpdateUserOrganization(ctx context.Context, req *pb.UpdateUserOrganizationRequest) (*pb.UpdateUserOrganizationResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	uo, err := s.repo.Update(ctx, req.Id, req.Role, req.IsDefault)
	if err != nil {
		if errors.Is(err, repository.ErrUserOrganizationNotFound) {
			return nil, status.Error(codes.NotFound, "user organization not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update user organization: %v", err)
	}

	return &pb.UpdateUserOrganizationResponse{
		UserOrganization: toProtoUserOrganization(uo),
	}, nil
}

// DeleteUserOrganization deletes a user organization
func (s *UserOrganizationServer) DeleteUserOrganization(ctx context.Context, req *pb.DeleteUserOrganizationRequest) (*pb.DeleteUserOrganizationResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	err := s.repo.Delete(ctx, req.Id)
	if err != nil {
		if errors.Is(err, repository.ErrUserOrganizationNotFound) {
			return nil, status.Error(codes.NotFound, "user organization not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete user organization: %v", err)
	}

	return &pb.DeleteUserOrganizationResponse{
		Success: true,
	}, nil
}

// ListUserOrganizations retrieves user organizations with pagination
func (s *UserOrganizationServer) ListUserOrganizations(ctx context.Context, req *pb.ListUserOrganizationsRequest) (*pb.ListUserOrganizationsResponse, error) {
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

	uos, err := s.repo.List(ctx, limit+1, offset) // +1 to check if there's next page
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list user organizations: %v", err)
	}

	var nextPageToken string
	if len(uos) > limit {
		uos = uos[:limit]
		nextPageToken = "next" // In production, encode the offset
	}

	protoUOs := make([]*pb.UserOrganization, len(uos))
	for i, uo := range uos {
		protoUOs[i] = toProtoUserOrganization(uo)
	}

	return &pb.ListUserOrganizationsResponse{
		UserOrganizations: protoUOs,
		NextPageToken:     nextPageToken,
	}, nil
}

// ListUserOrganizationsByUser retrieves all organizations for a user
func (s *UserOrganizationServer) ListUserOrganizationsByUser(ctx context.Context, req *pb.ListUserOrganizationsByUserRequest) (*pb.ListUserOrganizationsByUserResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	uos, err := s.repo.ListByUserID(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list user organizations by user: %v", err)
	}

	protoUOs := make([]*pb.UserOrganization, len(uos))
	for i, uo := range uos {
		protoUOs[i] = toProtoUserOrganization(uo)
	}

	return &pb.ListUserOrganizationsByUserResponse{
		UserOrganizations: protoUOs,
	}, nil
}

// ListUserOrganizationsByOrg retrieves all users for an organization
func (s *UserOrganizationServer) ListUserOrganizationsByOrg(ctx context.Context, req *pb.ListUserOrganizationsByOrgRequest) (*pb.ListUserOrganizationsByOrgResponse, error) {
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}

	uos, err := s.repo.ListByOrganizationID(ctx, req.OrganizationId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list user organizations by organization: %v", err)
	}

	protoUOs := make([]*pb.UserOrganization, len(uos))
	for i, uo := range uos {
		protoUOs[i] = toProtoUserOrganization(uo)
	}

	return &pb.ListUserOrganizationsByOrgResponse{
		UserOrganizations: protoUOs,
	}, nil
}

// toProtoUserOrganization converts repository model to proto message
func toProtoUserOrganization(uo *repository.UserOrganization) *pb.UserOrganization {
	return &pb.UserOrganization{
		Id:             uo.ID,
		UserId:         uo.UserID,
		OrganizationId: uo.OrganizationID,
		Role:           uo.Role,
		IsDefault:      uo.IsDefault,
		CreatedAt:      timestamppb.New(uo.CreatedAt),
		UpdatedAt:      timestamppb.New(uo.UpdatedAt),
	}
}
