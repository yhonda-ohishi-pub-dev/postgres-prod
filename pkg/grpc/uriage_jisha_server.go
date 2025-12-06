package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/pb"
	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/repository"
)

// UriageJishaServer implements the gRPC UriageJishaService
type UriageJishaServer struct {
	pb.UnimplementedUriageJishaServiceServer
	repo *repository.UriageJishaRepository
}

// NewUriageJishaServer creates a new gRPC server
func NewUriageJishaServer(repo *repository.UriageJishaRepository) *UriageJishaServer {
	return &UriageJishaServer{repo: repo}
}

// CreateUriageJisha creates a new uriage jisha entry
func (s *UriageJishaServer) CreateUriageJisha(ctx context.Context, req *pb.CreateUriageJishaRequest) (*pb.CreateUriageJishaResponse, error) {
	if req.Bumon == "" {
		return nil, status.Error(codes.InvalidArgument, "bumon is required")
	}
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}
	if req.Date == "" {
		return nil, status.Error(codes.InvalidArgument, "date is required")
	}

	var kingaku, typeVal *int32
	if req.Kingaku != nil {
		kingaku = req.Kingaku
	}
	if req.Type != nil {
		typeVal = req.Type
	}

	uriageJisha, err := s.repo.Create(ctx, req.Bumon, req.OrganizationId, kingaku, typeVal, req.Date)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create uriage jisha: %v", err)
	}

	return &pb.CreateUriageJishaResponse{
		UriageJisha: toProtoUriageJisha(uriageJisha),
	}, nil
}

// GetUriageJisha retrieves a uriage jisha by composite primary key (bumon, date, organization_id)
func (s *UriageJishaServer) GetUriageJisha(ctx context.Context, req *pb.GetUriageJishaRequest) (*pb.GetUriageJishaResponse, error) {
	if req.Bumon == "" {
		return nil, status.Error(codes.InvalidArgument, "bumon is required")
	}
	if req.Date == "" {
		return nil, status.Error(codes.InvalidArgument, "date is required")
	}
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}

	uriageJisha, err := s.repo.GetByPrimaryKey(ctx, req.Bumon, req.Date, req.OrganizationId)
	if err != nil {
		if errors.Is(err, repository.ErrUriageJishaNotFound) {
			return nil, status.Error(codes.NotFound, "uriage jisha not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get uriage jisha: %v", err)
	}

	return &pb.GetUriageJishaResponse{
		UriageJisha: toProtoUriageJisha(uriageJisha),
	}, nil
}

// UpdateUriageJisha updates an existing uriage jisha
func (s *UriageJishaServer) UpdateUriageJisha(ctx context.Context, req *pb.UpdateUriageJishaRequest) (*pb.UpdateUriageJishaResponse, error) {
	if req.Bumon == "" {
		return nil, status.Error(codes.InvalidArgument, "bumon is required")
	}
	if req.Date == "" {
		return nil, status.Error(codes.InvalidArgument, "date is required")
	}
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}

	var kingaku, typeVal *int32
	if req.Kingaku != nil {
		kingaku = req.Kingaku
	}
	if req.Type != nil {
		typeVal = req.Type
	}

	uriageJisha, err := s.repo.Update(ctx, req.Bumon, req.Date, req.OrganizationId, kingaku, typeVal)
	if err != nil {
		if errors.Is(err, repository.ErrUriageJishaNotFound) {
			return nil, status.Error(codes.NotFound, "uriage jisha not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update uriage jisha: %v", err)
	}

	return &pb.UpdateUriageJishaResponse{
		UriageJisha: toProtoUriageJisha(uriageJisha),
	}, nil
}

// DeleteUriageJisha deletes a uriage jisha by composite primary key
func (s *UriageJishaServer) DeleteUriageJisha(ctx context.Context, req *pb.DeleteUriageJishaRequest) (*pb.DeleteUriageJishaResponse, error) {
	if req.Bumon == "" {
		return nil, status.Error(codes.InvalidArgument, "bumon is required")
	}
	if req.Date == "" {
		return nil, status.Error(codes.InvalidArgument, "date is required")
	}
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}

	err := s.repo.Delete(ctx, req.Bumon, req.Date, req.OrganizationId)
	if err != nil {
		if errors.Is(err, repository.ErrUriageJishaNotFound) {
			return nil, status.Error(codes.NotFound, "uriage jisha not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete uriage jisha: %v", err)
	}

	return &pb.DeleteUriageJishaResponse{
		Success: true,
	}, nil
}

// ListUriageJishas retrieves all uriage jisha entries with pagination
func (s *UriageJishaServer) ListUriageJishas(ctx context.Context, req *pb.ListUriageJishasRequest) (*pb.ListUriageJishasResponse, error) {
	limit := int(req.PageSize)
	if limit <= 0 {
		limit = 10
	}

	offset := 0
	if req.PageToken != "" {
		// In a real implementation, decode page_token to get offset
	}

	uriageJishas, err := s.repo.List(ctx, limit+1, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list uriage jishas: %v", err)
	}

	var nextPageToken string
	if len(uriageJishas) > limit {
		uriageJishas = uriageJishas[:limit]
		nextPageToken = "next"
	}

	protoUriageJishas := make([]*pb.UriageJisha, len(uriageJishas))
	for i, uriageJisha := range uriageJishas {
		protoUriageJishas[i] = toProtoUriageJisha(uriageJisha)
	}

	return &pb.ListUriageJishasResponse{
		UriageJishas:  protoUriageJishas,
		NextPageToken: nextPageToken,
	}, nil
}

// ListUriageJishasByOrganization retrieves uriage jisha entries for a specific organization with pagination
func (s *UriageJishaServer) ListUriageJishasByOrganization(ctx context.Context, req *pb.ListUriageJishasByOrganizationRequest) (*pb.ListUriageJishasByOrganizationResponse, error) {
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}

	limit := int(req.PageSize)
	if limit <= 0 {
		limit = 10
	}

	offset := 0
	if req.PageToken != "" {
		// In a real implementation, decode page_token to get offset
	}

	uriageJishas, err := s.repo.ListByOrganization(ctx, req.OrganizationId, limit+1, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list uriage jishas by organization: %v", err)
	}

	var nextPageToken string
	if len(uriageJishas) > limit {
		uriageJishas = uriageJishas[:limit]
		nextPageToken = "next"
	}

	protoUriageJishas := make([]*pb.UriageJisha, len(uriageJishas))
	for i, uriageJisha := range uriageJishas {
		protoUriageJishas[i] = toProtoUriageJisha(uriageJisha)
	}

	return &pb.ListUriageJishasByOrganizationResponse{
		UriageJishas:  protoUriageJishas,
		NextPageToken: nextPageToken,
	}, nil
}

// toProtoUriageJisha converts repository model to proto message
func toProtoUriageJisha(uriageJisha *repository.UriageJisha) *pb.UriageJisha {
	proto := &pb.UriageJisha{
		Bumon:          uriageJisha.Bumon,
		OrganizationId: uriageJisha.OrganizationID,
		Date:           uriageJisha.Date,
	}

	if uriageJisha.Kingaku != nil {
		proto.Kingaku = uriageJisha.Kingaku
	}
	if uriageJisha.Type != nil {
		proto.Type = uriageJisha.Type
	}

	return proto
}
