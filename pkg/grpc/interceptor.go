package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/db"
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
