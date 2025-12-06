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

// AppUserServer implements the gRPC AppUserService
type AppUserServer struct {
	pb.UnimplementedAppUserServiceServer
	repo *repository.AppUserRepository
}

// NewAppUserServer creates a new gRPC server
func NewAppUserServer(repo *repository.AppUserRepository) *AppUserServer {
	return &AppUserServer{repo: repo}
}

// CreateAppUser creates a new app user
func (s *AppUserServer) CreateAppUser(ctx context.Context, req *pb.CreateAppUserRequest) (*pb.CreateAppUserResponse, error) {
	if req.IamEmail == "" {
		return nil, status.Error(codes.InvalidArgument, "iam_email is required")
	}
	if req.DisplayName == "" {
		return nil, status.Error(codes.InvalidArgument, "display_name is required")
	}

	user, err := s.repo.Create(ctx, req.IamEmail, req.DisplayName, req.IsSuperadmin)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create app user: %v", err)
	}

	return &pb.CreateAppUserResponse{
		AppUser: toProtoAppUser(user),
	}, nil
}

// GetAppUser retrieves an app user by ID
func (s *AppUserServer) GetAppUser(ctx context.Context, req *pb.GetAppUserRequest) (*pb.GetAppUserResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	user, err := s.repo.GetByID(ctx, req.Id)
	if err != nil {
		if errors.Is(err, repository.ErrAppUserNotFound) {
			return nil, status.Error(codes.NotFound, "app user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get app user: %v", err)
	}

	return &pb.GetAppUserResponse{
		AppUser: toProtoAppUser(user),
	}, nil
}

// GetAppUserByIamEmail retrieves an app user by IAM email
func (s *AppUserServer) GetAppUserByIamEmail(ctx context.Context, req *pb.GetAppUserByIamEmailRequest) (*pb.GetAppUserByIamEmailResponse, error) {
	if req.IamEmail == "" {
		return nil, status.Error(codes.InvalidArgument, "iam_email is required")
	}

	user, err := s.repo.GetByIamEmail(ctx, req.IamEmail)
	if err != nil {
		if errors.Is(err, repository.ErrAppUserNotFound) {
			return nil, status.Error(codes.NotFound, "app user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get app user: %v", err)
	}

	return &pb.GetAppUserByIamEmailResponse{
		AppUser: toProtoAppUser(user),
	}, nil
}

// UpdateAppUser updates an existing app user
func (s *AppUserServer) UpdateAppUser(ctx context.Context, req *pb.UpdateAppUserRequest) (*pb.UpdateAppUserResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	if req.DisplayName == "" {
		return nil, status.Error(codes.InvalidArgument, "display_name is required")
	}

	user, err := s.repo.Update(ctx, req.Id, req.DisplayName, req.IsSuperadmin)
	if err != nil {
		if errors.Is(err, repository.ErrAppUserNotFound) {
			return nil, status.Error(codes.NotFound, "app user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update app user: %v", err)
	}

	return &pb.UpdateAppUserResponse{
		AppUser: toProtoAppUser(user),
	}, nil
}

// DeleteAppUser soft-deletes an app user
func (s *AppUserServer) DeleteAppUser(ctx context.Context, req *pb.DeleteAppUserRequest) (*pb.DeleteAppUserResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	err := s.repo.Delete(ctx, req.Id)
	if err != nil {
		if errors.Is(err, repository.ErrAppUserNotFound) {
			return nil, status.Error(codes.NotFound, "app user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete app user: %v", err)
	}

	return &pb.DeleteAppUserResponse{
		Success: true,
	}, nil
}

// ListAppUsers retrieves app users with pagination
func (s *AppUserServer) ListAppUsers(ctx context.Context, req *pb.ListAppUsersRequest) (*pb.ListAppUsersResponse, error) {
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

	users, err := s.repo.List(ctx, limit+1, offset) // +1 to check if there's next page
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list app users: %v", err)
	}

	var nextPageToken string
	if len(users) > limit {
		users = users[:limit]
		nextPageToken = "next" // In production, encode the offset
	}

	protoUsers := make([]*pb.AppUser, len(users))
	for i, user := range users {
		protoUsers[i] = toProtoAppUser(user)
	}

	return &pb.ListAppUsersResponse{
		AppUsers:      protoUsers,
		NextPageToken: nextPageToken,
	}, nil
}

// toProtoAppUser converts repository model to proto message
func toProtoAppUser(user *repository.AppUser) *pb.AppUser {
	proto := &pb.AppUser{
		Id:           user.ID,
		IamEmail:     user.IamEmail,
		DisplayName:  user.DisplayName,
		IsSuperadmin: user.IsSuperadmin,
		CreatedAt:    timestamppb.New(user.CreatedAt),
		UpdatedAt:    timestamppb.New(user.UpdatedAt),
	}
	if user.DeletedAt != nil {
		proto.DeletedAt = timestamppb.New(*user.DeletedAt)
	}
	return proto
}
