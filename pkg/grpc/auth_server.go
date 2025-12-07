package grpc

import (
	"context"
	"errors"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/auth"
	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/pb"
	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/repository"
)

// AuthServer implements the gRPC AuthService
type AuthServer struct {
	pb.UnimplementedAuthServiceServer
	appUserRepo    *repository.AppUserRepository
	oauthRepo      *repository.OAuthAccountRepository
	jwtService     *auth.JWTService
	googleClient   *auth.GoogleOAuthClient
	lineClient     *auth.LineOAuthClient
}

// NewAuthServer creates a new AuthServer
func NewAuthServer(
	appUserRepo *repository.AppUserRepository,
	oauthRepo *repository.OAuthAccountRepository,
	jwtService *auth.JWTService,
	googleClient *auth.GoogleOAuthClient,
	lineClient *auth.LineOAuthClient,
) *AuthServer {
	return &AuthServer{
		appUserRepo:  appUserRepo,
		oauthRepo:    oauthRepo,
		jwtService:   jwtService,
		googleClient: googleClient,
		lineClient:   lineClient,
	}
}

// AuthWithGoogle authenticates a user with Google OAuth
func (s *AuthServer) AuthWithGoogle(ctx context.Context, req *pb.AuthWithGoogleRequest) (*pb.AuthResponse, error) {
	if req.Code == "" {
		return nil, status.Error(codes.InvalidArgument, "code is required")
	}

	// Exchange authorization code for tokens
	tokenResp, err := s.googleClient.ExchangeCode(ctx, req.Code)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "failed to exchange code: %v", err)
	}

	// Get user info from Google
	userInfo, err := s.googleClient.GetUserInfo(ctx, tokenResp.AccessToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user info: %v", err)
	}

	// Find or create user
	user, err := s.findOrCreateUser(ctx, "google", userInfo.ID, &userInfo.Email, userInfo.Name, &userInfo.Picture, tokenResp.AccessToken, tokenResp.RefreshToken, tokenResp.ExpiresIn)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to find or create user: %v", err)
	}

	// Generate JWT tokens
	return s.generateAuthResponse(user)
}

// AuthWithLine authenticates a user with LINE OAuth
func (s *AuthServer) AuthWithLine(ctx context.Context, req *pb.AuthWithLineRequest) (*pb.AuthResponse, error) {
	if req.Code == "" {
		return nil, status.Error(codes.InvalidArgument, "code is required")
	}

	// Exchange authorization code for tokens
	tokenResp, err := s.lineClient.ExchangeCode(ctx, req.Code)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "failed to exchange code: %v", err)
	}

	// Get user info from LINE
	userInfo, err := s.lineClient.GetUserInfo(ctx, tokenResp.AccessToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user info: %v", err)
	}

	// Try to get email from ID token if available
	var email *string
	if tokenResp.IDToken != "" {
		if payload, err := s.lineClient.VerifyIDToken(ctx, tokenResp.IDToken); err == nil && payload.Email != "" {
			email = &payload.Email
		}
	}

	// Find or create user
	var pictureURL *string
	if userInfo.PictureURL != "" {
		pictureURL = &userInfo.PictureURL
	}
	user, err := s.findOrCreateUser(ctx, "line", userInfo.UserID, email, userInfo.DisplayName, pictureURL, tokenResp.AccessToken, tokenResp.RefreshToken, tokenResp.ExpiresIn)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to find or create user: %v", err)
	}

	// Generate JWT tokens
	return s.generateAuthResponse(user)
}

// RefreshToken refreshes an access token
func (s *AuthServer) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.AuthResponse, error) {
	if req.RefreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh_token is required")
	}

	// Validate refresh token
	userID, err := s.jwtService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		if errors.Is(err, auth.ErrExpiredToken) {
			return nil, status.Error(codes.Unauthenticated, "refresh token has expired")
		}
		return nil, status.Error(codes.Unauthenticated, "invalid refresh token")
	}

	// Get user
	user, err := s.appUserRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrAppUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	// Generate new JWT tokens
	return s.generateAuthResponse(user)
}

// GetAuthURL returns the OAuth authorization URL
func (s *AuthServer) GetAuthURL(ctx context.Context, req *pb.GetAuthURLRequest) (*pb.GetAuthURLResponse, error) {
	if req.Provider == "" {
		return nil, status.Error(codes.InvalidArgument, "provider is required")
	}
	if req.State == "" {
		return nil, status.Error(codes.InvalidArgument, "state is required")
	}

	var url string
	switch req.Provider {
	case "google":
		url = s.googleClient.GetAuthURL(req.State)
	case "line":
		url = s.lineClient.GetAuthURL(req.State)
	default:
		return nil, status.Errorf(codes.InvalidArgument, "unsupported provider: %s", req.Provider)
	}

	return &pb.GetAuthURLResponse{Url: url}, nil
}

// ValidateToken validates an access token
func (s *AuthServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	if req.AccessToken == "" {
		return nil, status.Error(codes.InvalidArgument, "access_token is required")
	}

	claims, err := s.jwtService.ValidateAccessToken(req.AccessToken)
	if err != nil {
		return &pb.ValidateTokenResponse{Valid: false}, nil
	}

	// Get user from database
	user, err := s.appUserRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return &pb.ValidateTokenResponse{Valid: false}, nil
	}

	return &pb.ValidateTokenResponse{
		Valid: true,
		User:  toAuthProtoAppUser(user),
	}, nil
}

// findOrCreateUser finds an existing user by OAuth account or creates a new one
func (s *AuthServer) findOrCreateUser(ctx context.Context, provider, providerUserID string, email *string, displayName string, avatarURL *string, accessToken, refreshToken string, expiresIn int) (*repository.AppUser, error) {
	// Try to find existing OAuth account
	oauthAccount, err := s.oauthRepo.GetByProviderAndProviderUserID(ctx, provider, providerUserID)
	if err == nil {
		// Found existing OAuth account, get the user
		user, err := s.appUserRepo.GetByID(ctx, oauthAccount.AppUserID)
		if err != nil {
			return nil, err
		}

		// Update tokens
		tokenExpiresAt := time.Now().Add(time.Duration(expiresIn) * time.Second)
		_, err = s.oauthRepo.UpdateTokens(ctx, oauthAccount.ID, &accessToken, &refreshToken, &tokenExpiresAt)
		if err != nil {
			return nil, err
		}

		return user, nil
	}

	if !errors.Is(err, repository.ErrOAuthAccountNotFound) {
		return nil, err
	}

	// No existing OAuth account, create new user and OAuth account
	user, err := s.appUserRepo.Create(ctx, email, displayName, avatarURL, false)
	if err != nil {
		return nil, err
	}

	// Create OAuth account
	tokenExpiresAt := time.Now().Add(time.Duration(expiresIn) * time.Second)
	_, err = s.oauthRepo.Create(ctx, user.ID, provider, providerUserID, email, &accessToken, &refreshToken, &tokenExpiresAt)
	if err != nil {
		// Rollback user creation
		_ = s.appUserRepo.Delete(ctx, user.ID)
		return nil, err
	}

	return user, nil
}

// generateAuthResponse generates an AuthResponse with JWT tokens
func (s *AuthServer) generateAuthResponse(user *repository.AppUser) (*pb.AuthResponse, error) {
	var email string
	if user.Email != nil {
		email = *user.Email
	}

	tokenPair, err := s.jwtService.GenerateTokenPair(user.ID, email, user.DisplayName, user.IsSuperadmin)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate tokens: %v", err)
	}

	return &pb.AuthResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
		User:         toAuthProtoAppUser(user),
	}, nil
}

// toAuthProtoAppUser converts repository model to proto message
func toAuthProtoAppUser(user *repository.AppUser) *pb.AppUser {
	proto := &pb.AppUser{
		Id:           user.ID,
		Email:        user.Email,
		DisplayName:  user.DisplayName,
		AvatarUrl:    user.AvatarURL,
		IsSuperadmin: user.IsSuperadmin,
		CreatedAt:    timestamppb.New(user.CreatedAt),
		UpdatedAt:    timestamppb.New(user.UpdatedAt),
	}
	if user.DeletedAt != nil {
		proto.DeletedAt = timestamppb.New(*user.DeletedAt)
	}
	return proto
}
