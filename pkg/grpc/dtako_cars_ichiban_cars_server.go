package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/pb"
	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/repository"
)

// DtakoCarsIchibanCarsServer implements the gRPC DtakoCarsIchibanCarsService
type DtakoCarsIchibanCarsServer struct {
	pb.UnimplementedDtakoCarsIchibanCarsServiceServer
	repo *repository.DtakoCarsIchibanCarsRepository
}

// NewDtakoCarsIchibanCarsServer creates a new gRPC server
func NewDtakoCarsIchibanCarsServer(repo *repository.DtakoCarsIchibanCarsRepository) *DtakoCarsIchibanCarsServer {
	return &DtakoCarsIchibanCarsServer{repo: repo}
}

// CreateDtakoCarsIchibanCars creates a new dtako_cars_ichiban_cars entry
func (s *DtakoCarsIchibanCarsServer) CreateDtakoCarsIchibanCars(ctx context.Context, req *pb.CreateDtakoCarsIchibanCarsRequest) (*pb.CreateDtakoCarsIchibanCarsResponse, error) {
	if req.IdDtako == "" {
		return nil, status.Error(codes.InvalidArgument, "id_dtako is required")
	}
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}

	var id *string
	if req.Id != nil {
		id = req.Id
	}

	entry, err := s.repo.Create(ctx, req.IdDtako, req.OrganizationId, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create dtako cars ichiban cars entry: %v", err)
	}

	return &pb.CreateDtakoCarsIchibanCarsResponse{
		DtakoCarsIchibanCars: toProtoDtakoCarsIchibanCars(entry),
	}, nil
}

// GetDtakoCarsIchibanCars retrieves an entry by composite primary key (id_dtako, organization_id)
func (s *DtakoCarsIchibanCarsServer) GetDtakoCarsIchibanCars(ctx context.Context, req *pb.GetDtakoCarsIchibanCarsRequest) (*pb.GetDtakoCarsIchibanCarsResponse, error) {
	if req.IdDtako == "" {
		return nil, status.Error(codes.InvalidArgument, "id_dtako is required")
	}
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}

	entry, err := s.repo.GetByDtakoAndOrg(ctx, req.IdDtako, req.OrganizationId)
	if err != nil {
		if errors.Is(err, repository.ErrDtakoCarsIchibanCarsNotFound) {
			return nil, status.Error(codes.NotFound, "dtako cars ichiban cars entry not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get dtako cars ichiban cars entry: %v", err)
	}

	return &pb.GetDtakoCarsIchibanCarsResponse{
		DtakoCarsIchibanCars: toProtoDtakoCarsIchibanCars(entry),
	}, nil
}

// UpdateDtakoCarsIchibanCars updates an existing entry
func (s *DtakoCarsIchibanCarsServer) UpdateDtakoCarsIchibanCars(ctx context.Context, req *pb.UpdateDtakoCarsIchibanCarsRequest) (*pb.UpdateDtakoCarsIchibanCarsResponse, error) {
	if req.IdDtako == "" {
		return nil, status.Error(codes.InvalidArgument, "id_dtako is required")
	}
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}

	var id *string
	if req.Id != nil {
		id = req.Id
	}

	entry, err := s.repo.Update(ctx, req.IdDtako, req.OrganizationId, id)
	if err != nil {
		if errors.Is(err, repository.ErrDtakoCarsIchibanCarsNotFound) {
			return nil, status.Error(codes.NotFound, "dtako cars ichiban cars entry not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update dtako cars ichiban cars entry: %v", err)
	}

	return &pb.UpdateDtakoCarsIchibanCarsResponse{
		DtakoCarsIchibanCars: toProtoDtakoCarsIchibanCars(entry),
	}, nil
}

// DeleteDtakoCarsIchibanCars deletes an entry by composite primary key
func (s *DtakoCarsIchibanCarsServer) DeleteDtakoCarsIchibanCars(ctx context.Context, req *pb.DeleteDtakoCarsIchibanCarsRequest) (*pb.DeleteDtakoCarsIchibanCarsResponse, error) {
	if req.IdDtako == "" {
		return nil, status.Error(codes.InvalidArgument, "id_dtako is required")
	}
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}

	err := s.repo.Delete(ctx, req.IdDtako, req.OrganizationId)
	if err != nil {
		if errors.Is(err, repository.ErrDtakoCarsIchibanCarsNotFound) {
			return nil, status.Error(codes.NotFound, "dtako cars ichiban cars entry not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete dtako cars ichiban cars entry: %v", err)
	}

	return &pb.DeleteDtakoCarsIchibanCarsResponse{
		Success: true,
	}, nil
}

// ListDtakoCarsIchibanCars retrieves all entries with pagination
func (s *DtakoCarsIchibanCarsServer) ListDtakoCarsIchibanCars(ctx context.Context, req *pb.ListDtakoCarsIchibanCarsRequest) (*pb.ListDtakoCarsIchibanCarsResponse, error) {
	// Note: This method is not implemented in the repository, would need to add it
	// For now, we'll return an error
	return nil, status.Error(codes.Unimplemented, "ListDtakoCarsIchibanCars is not yet implemented")
}

// ListDtakoCarsIchibanCarsByOrganization retrieves entries for a specific organization with pagination
func (s *DtakoCarsIchibanCarsServer) ListDtakoCarsIchibanCarsByOrganization(ctx context.Context, req *pb.ListDtakoCarsIchibanCarsByOrganizationRequest) (*pb.ListDtakoCarsIchibanCarsByOrganizationResponse, error) {
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

	entries, err := s.repo.ListByOrganization(ctx, req.OrganizationId, limit+1, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list dtako cars ichiban cars by organization: %v", err)
	}

	var nextPageToken string
	if len(entries) > limit {
		entries = entries[:limit]
		nextPageToken = "next"
	}

	protoEntries := make([]*pb.DtakoCarsIchibanCars, len(entries))
	for i, entry := range entries {
		protoEntries[i] = toProtoDtakoCarsIchibanCars(entry)
	}

	return &pb.ListDtakoCarsIchibanCarsByOrganizationResponse{
		DtakoCarsIchibanCars: protoEntries,
		NextPageToken:        nextPageToken,
	}, nil
}

// toProtoDtakoCarsIchibanCars converts repository model to proto message
func toProtoDtakoCarsIchibanCars(entry *repository.DtakoCarsIchibanCars) *pb.DtakoCarsIchibanCars {
	proto := &pb.DtakoCarsIchibanCars{
		IdDtako:        entry.IdDtako,
		OrganizationId: entry.OrganizationID,
	}

	if entry.Id != nil {
		proto.Id = entry.Id
	}

	return proto
}
