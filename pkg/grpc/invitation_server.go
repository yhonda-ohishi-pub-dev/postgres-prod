package grpc

import (
	"context"
	"errors"
	"fmt"
	"os"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/pb"
	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/repository"
)

// InvitationServer implements the gRPC InvitationService
type InvitationServer struct {
	pb.UnimplementedInvitationServiceServer
	invRepo  *repository.InvitationRepository
	orgRepo  *repository.OrganizationRepository
	userRepo *repository.UserOrganizationRepository
}

// NewInvitationServer creates a new gRPC server
func NewInvitationServer(
	invRepo *repository.InvitationRepository,
	orgRepo *repository.OrganizationRepository,
	userRepo *repository.UserOrganizationRepository,
) *InvitationServer {
	return &InvitationServer{
		invRepo:  invRepo,
		orgRepo:  orgRepo,
		userRepo: userRepo,
	}
}

// getInviteURL builds the invitation URL
func getInviteURL(token string) string {
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "https://localhost:3000"
	}
	return fmt.Sprintf("%s/invite/%s", frontendURL, token)
}

// CreateInvitation creates a new invitation
func (s *InvitationServer) CreateInvitation(ctx context.Context, req *pb.CreateInvitationRequest) (*pb.CreateInvitationResponse, error) {
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	// Get inviter user_id from JWT
	userID, ok := GetUserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "authentication required")
	}

	// Check if organization exists
	_, err := s.orgRepo.GetByID(ctx, req.OrganizationId)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "organization not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get organization: %v", err)
	}

	// Check if there's already a pending invitation
	existing, err := s.invRepo.GetPendingByEmailAndOrg(ctx, req.Email, req.OrganizationId)
	if err == nil && existing != nil {
		// Return existing invitation
		return &pb.CreateInvitationResponse{
			Invitation: toProtoInvitation(existing),
			InviteUrl:  getInviteURL(existing.Token),
		}, nil
	}

	// Create new invitation
	role := req.Role
	if role == "" {
		role = "member"
	}

	inv, err := s.invRepo.Create(ctx, req.OrganizationId, req.Email, role, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create invitation: %v", err)
	}

	return &pb.CreateInvitationResponse{
		Invitation: toProtoInvitation(inv),
		InviteUrl:  getInviteURL(inv.Token),
	}, nil
}

// GetInvitation retrieves an invitation by ID
func (s *InvitationServer) GetInvitation(ctx context.Context, req *pb.GetInvitationRequest) (*pb.GetInvitationResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	inv, err := s.invRepo.GetByID(ctx, req.Id)
	if err != nil {
		if errors.Is(err, repository.ErrInvitationNotFound) {
			return nil, status.Error(codes.NotFound, "invitation not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get invitation: %v", err)
	}

	return &pb.GetInvitationResponse{
		Invitation: toProtoInvitation(inv),
	}, nil
}

// GetInvitationByToken retrieves an invitation by token (for accepting)
func (s *InvitationServer) GetInvitationByToken(ctx context.Context, req *pb.GetInvitationByTokenRequest) (*pb.GetInvitationByTokenResponse, error) {
	if req.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}

	inv, err := s.invRepo.GetByToken(ctx, req.Token)
	if err != nil {
		if errors.Is(err, repository.ErrInvitationNotFound) {
			return nil, status.Error(codes.NotFound, "invitation not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get invitation: %v", err)
	}

	// Get organization details
	org, err := s.orgRepo.GetByID(ctx, inv.OrganizationID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get organization: %v", err)
	}

	return &pb.GetInvitationByTokenResponse{
		Invitation:   toProtoInvitation(inv),
		Organization: toProtoOrganization(org),
	}, nil
}

// AcceptInvitation accepts an invitation and creates user_organization
func (s *InvitationServer) AcceptInvitation(ctx context.Context, req *pb.AcceptInvitationRequest) (*pb.AcceptInvitationResponse, error) {
	if req.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}

	// Get current user from JWT
	userID, ok := GetUserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "authentication required")
	}

	// Get invitation by token
	inv, err := s.invRepo.GetByToken(ctx, req.Token)
	if err != nil {
		if errors.Is(err, repository.ErrInvitationNotFound) {
			return nil, status.Error(codes.NotFound, "invitation not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get invitation: %v", err)
	}

	// Check if already a member
	_, err = s.userRepo.GetByUserAndOrganization(ctx, userID, inv.OrganizationID)
	if err == nil {
		return nil, status.Error(codes.AlreadyExists, "already a member of this organization")
	}

	// Accept invitation
	acceptedInv, err := s.invRepo.Accept(ctx, inv.ID, userID)
	if err != nil {
		if errors.Is(err, repository.ErrInvitationExpired) {
			return nil, status.Error(codes.FailedPrecondition, "invitation expired")
		}
		if errors.Is(err, repository.ErrInvitationUsed) {
			return nil, status.Error(codes.FailedPrecondition, "invitation already used")
		}
		return nil, status.Errorf(codes.Internal, "failed to accept invitation: %v", err)
	}

	// Create user_organization
	userOrg, err := s.userRepo.Create(ctx, userID, acceptedInv.OrganizationID, acceptedInv.Role, false)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user organization: %v", err)
	}

	return &pb.AcceptInvitationResponse{
		Success:          true,
		UserOrganization: toProtoUserOrganization(userOrg),
	}, nil
}

// CancelInvitation cancels an invitation
func (s *InvitationServer) CancelInvitation(ctx context.Context, req *pb.CancelInvitationRequest) (*pb.CancelInvitationResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	err := s.invRepo.Cancel(ctx, req.Id)
	if err != nil {
		if errors.Is(err, repository.ErrInvitationNotFound) {
			return nil, status.Error(codes.NotFound, "invitation not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to cancel invitation: %v", err)
	}

	return &pb.CancelInvitationResponse{
		Success: true,
	}, nil
}

// ListInvitations lists invitations for an organization
func (s *InvitationServer) ListInvitations(ctx context.Context, req *pb.ListInvitationsRequest) (*pb.ListInvitationsResponse, error) {
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}

	invitations, err := s.invRepo.List(ctx, req.OrganizationId, req.Status, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list invitations: %v", err)
	}

	protoInvitations := make([]*pb.Invitation, len(invitations))
	for i, inv := range invitations {
		protoInvitations[i] = toProtoInvitation(inv)
	}

	return &pb.ListInvitationsResponse{
		Invitations: protoInvitations,
	}, nil
}

// ResendInvitation regenerates token and extends expiry
func (s *InvitationServer) ResendInvitation(ctx context.Context, req *pb.ResendInvitationRequest) (*pb.ResendInvitationResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	inv, err := s.invRepo.Resend(ctx, req.Id)
	if err != nil {
		if errors.Is(err, repository.ErrInvitationNotFound) {
			return nil, status.Error(codes.NotFound, "invitation not found or already accepted")
		}
		return nil, status.Errorf(codes.Internal, "failed to resend invitation: %v", err)
	}

	return &pb.ResendInvitationResponse{
		Invitation: toProtoInvitation(inv),
		InviteUrl:  getInviteURL(inv.Token),
	}, nil
}

// toProtoInvitation converts repository model to proto message
func toProtoInvitation(inv *repository.Invitation) *pb.Invitation {
	pbInv := &pb.Invitation{
		Id:             inv.ID,
		OrganizationId: inv.OrganizationID,
		Email:          inv.Email,
		Role:           inv.Role,
		Token:          inv.Token,
		InvitedBy:      inv.InvitedBy,
		Status:         inv.Status,
		ExpiresAt:      timestamppb.New(inv.ExpiresAt),
		CreatedAt:      timestamppb.New(inv.CreatedAt),
		UpdatedAt:      timestamppb.New(inv.UpdatedAt),
	}

	if inv.AcceptedAt != nil {
		pbInv.AcceptedAt = timestamppb.New(*inv.AcceptedAt)
	}
	if inv.AcceptedBy != nil {
		pbInv.AcceptedBy = inv.AcceptedBy
	}

	return pbInv
}

