package grpc

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/auth"
	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/db"
)

// Context keys for user info
type contextKey string

const (
	UserIDKey       contextKey = "user_id"
	UserEmailKey    contextKey = "user_email"
	UserNameKey     contextKey = "user_name"
	IsSuperadminKey contextKey = "is_superadmin"
)

const (
	// OrganizationIDHeader is the gRPC metadata key for organization ID
	OrganizationIDHeader = "x-organization-id"
)

// skipRLSPrefixes are method prefixes that don't require x-organization-id header
var skipRLSPrefixes = []string{
	"/grpc.health.v1.Health/",
	"/grpc.reflection.v1alpha.ServerReflection/",
	"/grpc.reflection.v1.ServerReflection/",
	"/auth.AuthService/",
	"/organization.OrganizationService/",
	"/organization.AppUserService/",
	"/organization.UserOrganizationService/",
}

// shouldSkipRLS checks if the method should skip RLS enforcement
func shouldSkipRLS(fullMethod string) bool {
	for _, prefix := range skipRLSPrefixes {
		if len(fullMethod) >= len(prefix) && fullMethod[:len(prefix)] == prefix {
			return true
		}
	}
	return false
}

// RLSUnaryInterceptor extracts organization_id from gRPC metadata
// and adds it to the context for RLS enforcement.
func RLSUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Skip RLS for methods that don't require organization context
		if shouldSkipRLS(info.FullMethod) {
			return handler(ctx, req)
		}

		// Extract organization_id from metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.InvalidArgument, "missing metadata")
		}

		orgIDs := md.Get(OrganizationIDHeader)
		if len(orgIDs) == 0 {
			return nil, status.Error(codes.InvalidArgument, "missing x-organization-id header")
		}

		orgID := orgIDs[0]
		if orgID == "" {
			return nil, status.Error(codes.InvalidArgument, "x-organization-id cannot be empty")
		}

		// Add organization_id to context for RLS
		ctx = db.WithOrganizationID(ctx, orgID)

		return handler(ctx, req)
	}
}

// RLSStreamInterceptor extracts organization_id from gRPC metadata
// for streaming RPCs and adds it to the context for RLS enforcement.
func RLSStreamInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		// Skip RLS for methods that don't require organization context
		if shouldSkipRLS(info.FullMethod) {
			return handler(srv, ss)
		}

		ctx := ss.Context()

		// Extract organization_id from metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return status.Error(codes.InvalidArgument, "missing metadata")
		}

		orgIDs := md.Get(OrganizationIDHeader)
		if len(orgIDs) == 0 {
			return status.Error(codes.InvalidArgument, "missing x-organization-id header")
		}

		orgID := orgIDs[0]
		if orgID == "" {
			return status.Error(codes.InvalidArgument, "x-organization-id cannot be empty")
		}

		// Wrap stream with new context containing organization_id
		wrapped := &wrappedServerStream{
			ServerStream: ss,
			ctx:          db.WithOrganizationID(ctx, orgID),
		}

		return handler(srv, wrapped)
	}
}

// wrappedServerStream wraps grpc.ServerStream to override Context()
type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}

// skipAuthPrefixes are method prefixes that don't require JWT authentication
var skipAuthPrefixes = []string{
	"/grpc.health.v1.Health/",
	"/grpc.reflection.v1alpha.ServerReflection/",
	"/grpc.reflection.v1.ServerReflection/",
	"/auth.AuthService/",
}

// shouldSkipAuth checks if the method should skip JWT authentication
func shouldSkipAuth(fullMethod string) bool {
	for _, prefix := range skipAuthPrefixes {
		if len(fullMethod) >= len(prefix) && fullMethod[:len(prefix)] == prefix {
			return true
		}
	}
	return false
}

// JWTUnaryInterceptor validates JWT token and adds user info to context
func JWTUnaryInterceptor(jwtService *auth.JWTService) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Skip auth for public methods
		if shouldSkipAuth(info.FullMethod) {
			return handler(ctx, req)
		}

		// Extract Authorization header from metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		authHeaders := md.Get("authorization")
		if len(authHeaders) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing authorization header")
		}

		// Parse "Bearer <token>"
		authHeader := authHeaders[0]
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return nil, status.Error(codes.Unauthenticated, "invalid authorization format")
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate token
		claims, err := jwtService.ValidateAccessToken(token)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		}

		// Add user info to context
		ctx = context.WithValue(ctx, UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, UserEmailKey, claims.Email)
		ctx = context.WithValue(ctx, UserNameKey, claims.DisplayName)
		ctx = context.WithValue(ctx, IsSuperadminKey, claims.IsSuperadmin)

		return handler(ctx, req)
	}
}

// GetUserIDFromContext extracts user_id from context
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok && userID != ""
}

// GetUserEmailFromContext extracts user_email from context
func GetUserEmailFromContext(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(UserEmailKey).(string)
	return email, ok
}

// GetIsSuperadminFromContext extracts is_superadmin from context
func GetIsSuperadminFromContext(ctx context.Context) bool {
	isSuperadmin, ok := ctx.Value(IsSuperadminKey).(bool)
	return ok && isSuperadmin
}
