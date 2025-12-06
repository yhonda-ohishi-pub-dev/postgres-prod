package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/pb"
	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/repository"
)

// UriageServer implements the gRPC UriageService
type UriageServer struct {
	pb.UnimplementedUriageServiceServer
	repo *repository.UriageRepository
}

// NewUriageServer creates a new gRPC server
func NewUriageServer(repo *repository.UriageRepository) *UriageServer {
	return &UriageServer{repo: repo}
}

// CreateUriage creates a new uriage entry
func (s *UriageServer) CreateUriage(ctx context.Context, req *pb.CreateUriageRequest) (*pb.CreateUriageResponse, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}
	if req.Bumon == "" {
		return nil, status.Error(codes.InvalidArgument, "bumon is required")
	}
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}
	if req.Date == "" {
		return nil, status.Error(codes.InvalidArgument, "date is required")
	}

	var kingaku, uriageType, cam *int32
	if req.Kingaku != nil {
		kingaku = req.Kingaku
	}
	if req.Type != nil {
		uriageType = req.Type
	}
	if req.Cam != nil {
		cam = req.Cam
	}

	uriage, err := s.repo.Create(ctx, req.Name, req.Bumon, req.OrganizationId, kingaku, uriageType, cam, req.Date)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create uriage: %v", err)
	}

	return &pb.CreateUriageResponse{
		Uriage: toProtoUriage(uriage),
	}, nil
}

// GetUriage retrieves a uriage by composite primary key (name, bumon, date, organization_id)
func (s *UriageServer) GetUriage(ctx context.Context, req *pb.GetUriageRequest) (*pb.GetUriageResponse, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}
	if req.Bumon == "" {
		return nil, status.Error(codes.InvalidArgument, "bumon is required")
	}
	if req.Date == "" {
		return nil, status.Error(codes.InvalidArgument, "date is required")
	}
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}

	uriage, err := s.repo.GetByPrimaryKey(ctx, req.Name, req.Bumon, req.Date, req.OrganizationId)
	if err != nil {
		if errors.Is(err, repository.ErrUriageNotFound) {
			return nil, status.Error(codes.NotFound, "uriage not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get uriage: %v", err)
	}

	return &pb.GetUriageResponse{
		Uriage: toProtoUriage(uriage),
	}, nil
}

// UpdateUriage updates an existing uriage
func (s *UriageServer) UpdateUriage(ctx context.Context, req *pb.UpdateUriageRequest) (*pb.UpdateUriageResponse, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}
	if req.Bumon == "" {
		return nil, status.Error(codes.InvalidArgument, "bumon is required")
	}
	if req.Date == "" {
		return nil, status.Error(codes.InvalidArgument, "date is required")
	}
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}

	var kingaku, uriageType, cam *int32
	if req.Kingaku != nil {
		kingaku = req.Kingaku
	}
	if req.Type != nil {
		uriageType = req.Type
	}
	if req.Cam != nil {
		cam = req.Cam
	}

	uriage, err := s.repo.Update(ctx, req.Name, req.Bumon, req.Date, req.OrganizationId, kingaku, uriageType, cam)
	if err != nil {
		if errors.Is(err, repository.ErrUriageNotFound) {
			return nil, status.Error(codes.NotFound, "uriage not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update uriage: %v", err)
	}

	return &pb.UpdateUriageResponse{
		Uriage: toProtoUriage(uriage),
	}, nil
}

// DeleteUriage deletes a uriage by composite primary key
func (s *UriageServer) DeleteUriage(ctx context.Context, req *pb.DeleteUriageRequest) (*pb.DeleteUriageResponse, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}
	if req.Bumon == "" {
		return nil, status.Error(codes.InvalidArgument, "bumon is required")
	}
	if req.Date == "" {
		return nil, status.Error(codes.InvalidArgument, "date is required")
	}
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}

	err := s.repo.Delete(ctx, req.Name, req.Bumon, req.Date, req.OrganizationId)
	if err != nil {
		if errors.Is(err, repository.ErrUriageNotFound) {
			return nil, status.Error(codes.NotFound, "uriage not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete uriage: %v", err)
	}

	return &pb.DeleteUriageResponse{
		Success: true,
	}, nil
}

// ListUriages retrieves all uriage entries with pagination
func (s *UriageServer) ListUriages(ctx context.Context, req *pb.ListUriagesRequest) (*pb.ListUriagesResponse, error) {
	limit := int(req.PageSize)
	if limit <= 0 {
		limit = 10
	}

	offset := 0
	if req.PageToken != "" {
		// In a real implementation, decode page_token to get offset
	}

	uriages, err := s.repo.List(ctx, limit+1, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list uriages: %v", err)
	}

	var nextPageToken string
	if len(uriages) > limit {
		uriages = uriages[:limit]
		nextPageToken = "next"
	}

	protoUriages := make([]*pb.Uriage, len(uriages))
	for i, uriage := range uriages {
		protoUriages[i] = toProtoUriage(uriage)
	}

	return &pb.ListUriagesResponse{
		Uriages:       protoUriages,
		NextPageToken: nextPageToken,
	}, nil
}

// ListUriagesByOrganization retrieves uriage entries for a specific organization with pagination
func (s *UriageServer) ListUriagesByOrganization(ctx context.Context, req *pb.ListUriagesByOrganizationRequest) (*pb.ListUriagesByOrganizationResponse, error) {
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

	uriages, err := s.repo.ListByOrganization(ctx, req.OrganizationId, limit+1, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list uriages by organization: %v", err)
	}

	var nextPageToken string
	if len(uriages) > limit {
		uriages = uriages[:limit]
		nextPageToken = "next"
	}

	protoUriages := make([]*pb.Uriage, len(uriages))
	for i, uriage := range uriages {
		protoUriages[i] = toProtoUriage(uriage)
	}

	return &pb.ListUriagesByOrganizationResponse{
		Uriages:       protoUriages,
		NextPageToken: nextPageToken,
	}, nil
}

// toProtoUriage converts repository model to proto message
func toProtoUriage(uriage *repository.Uriage) *pb.Uriage {
	proto := &pb.Uriage{
		Name:           uriage.Name,
		Bumon:          uriage.Bumon,
		OrganizationId: uriage.OrganizationID,
		Date:           uriage.Date,
	}

	if uriage.Kingaku != nil {
		proto.Kingaku = uriage.Kingaku
	}
	if uriage.Type != nil {
		proto.Type = uriage.Type
	}
	if uriage.Cam != nil {
		proto.Cam = uriage.Cam
	}

	return proto
}
