package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/pb"
	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/repository"
)

// CarInspectionDeregistrationServer implements the gRPC CarInspectionDeregistrationService
type CarInspectionDeregistrationServer struct {
	pb.UnimplementedCarInspectionDeregistrationServiceServer
	repo *repository.CarInspectionDeregistrationRepository
}

// NewCarInspectionDeregistrationServer creates a new gRPC server
func NewCarInspectionDeregistrationServer(repo *repository.CarInspectionDeregistrationRepository) *CarInspectionDeregistrationServer {
	return &CarInspectionDeregistrationServer{repo: repo}
}

// CreateCarInspectionDeregistration creates a new car inspection deregistration record
func (s *CarInspectionDeregistrationServer) CreateCarInspectionDeregistration(ctx context.Context, req *pb.CreateCarInspectionDeregistrationRequest) (*pb.CreateCarInspectionDeregistrationResponse, error) {
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}
	if req.CarId == "" {
		return nil, status.Error(codes.InvalidArgument, "car_id is required")
	}
	if req.TwodimensionCodeInfoValidPeriodExpirDate == "" {
		return nil, status.Error(codes.InvalidArgument, "twodimension_code_info_valid_period_expir_date is required")
	}

	record, err := s.repo.Create(ctx,
		req.OrganizationId,
		req.CarId,
		req.TwodimensionCodeInfoCarNo,
		req.CarNo,
		req.ValidPeriodExpirDateE,
		req.ValidPeriodExpirDateY,
		req.ValidPeriodExpirDateM,
		req.ValidPeriodExpirDateD,
		req.TwodimensionCodeInfoValidPeriodExpirDate,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create car inspection deregistration: %v", err)
	}

	return &pb.CreateCarInspectionDeregistrationResponse{
		CarInspectionDeregistration: toProtoCarInspectionDeregistration(record),
	}, nil
}

// GetCarInspectionDeregistration retrieves a car inspection deregistration record by composite primary key
func (s *CarInspectionDeregistrationServer) GetCarInspectionDeregistration(ctx context.Context, req *pb.GetCarInspectionDeregistrationRequest) (*pb.GetCarInspectionDeregistrationResponse, error) {
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}
	if req.CarId == "" {
		return nil, status.Error(codes.InvalidArgument, "car_id is required")
	}
	if req.TwodimensionCodeInfoValidPeriodExpirDate == "" {
		return nil, status.Error(codes.InvalidArgument, "twodimension_code_info_valid_period_expir_date is required")
	}

	record, err := s.repo.GetByPrimaryKey(ctx,
		req.OrganizationId,
		req.CarId,
		req.TwodimensionCodeInfoValidPeriodExpirDate,
	)
	if err != nil {
		if errors.Is(err, repository.ErrCarInspectionDeregistrationNotFound) {
			return nil, status.Error(codes.NotFound, "car inspection deregistration not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get car inspection deregistration: %v", err)
	}

	return &pb.GetCarInspectionDeregistrationResponse{
		CarInspectionDeregistration: toProtoCarInspectionDeregistration(record),
	}, nil
}

// UpdateCarInspectionDeregistration updates an existing car inspection deregistration record
func (s *CarInspectionDeregistrationServer) UpdateCarInspectionDeregistration(ctx context.Context, req *pb.UpdateCarInspectionDeregistrationRequest) (*pb.UpdateCarInspectionDeregistrationResponse, error) {
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}
	if req.CarId == "" {
		return nil, status.Error(codes.InvalidArgument, "car_id is required")
	}
	if req.TwodimensionCodeInfoValidPeriodExpirDate == "" {
		return nil, status.Error(codes.InvalidArgument, "twodimension_code_info_valid_period_expir_date is required")
	}

	record, err := s.repo.Update(ctx,
		req.OrganizationId,
		req.CarId,
		req.TwodimensionCodeInfoCarNo,
		req.CarNo,
		req.ValidPeriodExpirDateE,
		req.ValidPeriodExpirDateY,
		req.ValidPeriodExpirDateM,
		req.ValidPeriodExpirDateD,
		req.TwodimensionCodeInfoValidPeriodExpirDate,
	)
	if err != nil {
		if errors.Is(err, repository.ErrCarInspectionDeregistrationNotFound) {
			return nil, status.Error(codes.NotFound, "car inspection deregistration not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update car inspection deregistration: %v", err)
	}

	return &pb.UpdateCarInspectionDeregistrationResponse{
		CarInspectionDeregistration: toProtoCarInspectionDeregistration(record),
	}, nil
}

// DeleteCarInspectionDeregistration hard-deletes a car inspection deregistration record
func (s *CarInspectionDeregistrationServer) DeleteCarInspectionDeregistration(ctx context.Context, req *pb.DeleteCarInspectionDeregistrationRequest) (*pb.DeleteCarInspectionDeregistrationResponse, error) {
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}
	if req.CarId == "" {
		return nil, status.Error(codes.InvalidArgument, "car_id is required")
	}
	if req.TwodimensionCodeInfoValidPeriodExpirDate == "" {
		return nil, status.Error(codes.InvalidArgument, "twodimension_code_info_valid_period_expir_date is required")
	}

	err := s.repo.Delete(ctx,
		req.OrganizationId,
		req.CarId,
		req.TwodimensionCodeInfoValidPeriodExpirDate,
	)
	if err != nil {
		if errors.Is(err, repository.ErrCarInspectionDeregistrationNotFound) {
			return nil, status.Error(codes.NotFound, "car inspection deregistration not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete car inspection deregistration: %v", err)
	}

	return &pb.DeleteCarInspectionDeregistrationResponse{
		Success: true,
	}, nil
}

// ListCarInspectionDeregistrations retrieves car inspection deregistration records with pagination
func (s *CarInspectionDeregistrationServer) ListCarInspectionDeregistrations(ctx context.Context, req *pb.ListCarInspectionDeregistrationsRequest) (*pb.ListCarInspectionDeregistrationsResponse, error) {
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

	records, err := s.repo.List(ctx, limit+1, offset) // +1 to check if there's next page
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list car inspection deregistrations: %v", err)
	}

	var nextPageToken string
	if len(records) > limit {
		records = records[:limit]
		nextPageToken = "next" // In production, encode the offset
	}

	protoRecords := make([]*pb.CarInspectionDeregistration, len(records))
	for i, record := range records {
		protoRecords[i] = toProtoCarInspectionDeregistration(record)
	}

	return &pb.ListCarInspectionDeregistrationsResponse{
		CarInspectionDeregistrations: protoRecords,
		NextPageToken:                nextPageToken,
	}, nil
}

// ListCarInspectionDeregistrationsByOrganization retrieves car inspection deregistration records by organization with pagination
func (s *CarInspectionDeregistrationServer) ListCarInspectionDeregistrationsByOrganization(ctx context.Context, req *pb.ListCarInspectionDeregistrationsByOrganizationRequest) (*pb.ListCarInspectionDeregistrationsByOrganizationResponse, error) {
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}

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

	records, err := s.repo.ListByOrganization(ctx, req.OrganizationId, limit+1, offset) // +1 to check if there's next page
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list car inspection deregistrations by organization: %v", err)
	}

	var nextPageToken string
	if len(records) > limit {
		records = records[:limit]
		nextPageToken = "next" // In production, encode the offset
	}

	protoRecords := make([]*pb.CarInspectionDeregistration, len(records))
	for i, record := range records {
		protoRecords[i] = toProtoCarInspectionDeregistration(record)
	}

	return &pb.ListCarInspectionDeregistrationsByOrganizationResponse{
		CarInspectionDeregistrations: protoRecords,
		NextPageToken:                nextPageToken,
	}, nil
}

// toProtoCarInspectionDeregistration converts repository model to proto message
func toProtoCarInspectionDeregistration(record *repository.CarInspectionDeregistration) *pb.CarInspectionDeregistration {
	return &pb.CarInspectionDeregistration{
		OrganizationId:                           record.OrganizationID,
		CarId:                                    record.CarID,
		TwodimensionCodeInfoCarNo:                record.TwodimensionCodeInfoCarNo,
		CarNo:                                    record.CarNo,
		ValidPeriodExpirDateE:                    record.ValidPeriodExpirDateE,
		ValidPeriodExpirDateY:                    record.ValidPeriodExpirDateY,
		ValidPeriodExpirDateM:                    record.ValidPeriodExpirDateM,
		ValidPeriodExpirDateD:                    record.ValidPeriodExpirDateD,
		TwodimensionCodeInfoValidPeriodExpirDate: record.TwodimensionCodeInfoValidPeriodExpirDate,
	}
}
