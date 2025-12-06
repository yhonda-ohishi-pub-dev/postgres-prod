package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/pb"
	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/repository"
)

// CarInspectionDeregistrationFilesServer implements the gRPC CarInspectionDeregistrationFilesService
type CarInspectionDeregistrationFilesServer struct {
	pb.UnimplementedCarInspectionDeregistrationFilesServiceServer
	repo *repository.CarInspectionDeregistrationFilesRepository
}

// NewCarInspectionDeregistrationFilesServer creates a new gRPC server
func NewCarInspectionDeregistrationFilesServer(repo *repository.CarInspectionDeregistrationFilesRepository) *CarInspectionDeregistrationFilesServer {
	return &CarInspectionDeregistrationFilesServer{repo: repo}
}

// CreateCarInspectionDeregistrationFiles creates a new car inspection deregistration file
func (s *CarInspectionDeregistrationFilesServer) CreateCarInspectionDeregistrationFiles(ctx context.Context, req *pb.CreateCarInspectionDeregistrationFilesRequest) (*pb.CreateCarInspectionDeregistrationFilesResponse, error) {
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}
	if req.CarId == "" {
		return nil, status.Error(codes.InvalidArgument, "car_id is required")
	}
	if req.TwodimensionCodeInfoValidPeriodExpirDate == "" {
		return nil, status.Error(codes.InvalidArgument, "twodimension_code_info_valid_period_expir_date is required")
	}
	if req.FileUuid == "" {
		return nil, status.Error(codes.InvalidArgument, "file_uuid is required")
	}

	record, err := s.repo.Create(ctx,
		req.OrganizationId,
		req.CarId,
		req.TwodimensionCodeInfoValidPeriodExpirDate,
		req.FileUuid,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create car inspection deregistration file: %v", err)
	}

	return &pb.CreateCarInspectionDeregistrationFilesResponse{
		CarInspectionDeregistrationFiles: toProtoCarInspectionDeregistrationFiles(record),
	}, nil
}

// GetCarInspectionDeregistrationFiles retrieves a car inspection deregistration file by composite primary key
func (s *CarInspectionDeregistrationFilesServer) GetCarInspectionDeregistrationFiles(ctx context.Context, req *pb.GetCarInspectionDeregistrationFilesRequest) (*pb.GetCarInspectionDeregistrationFilesResponse, error) {
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}
	if req.CarId == "" {
		return nil, status.Error(codes.InvalidArgument, "car_id is required")
	}
	if req.TwodimensionCodeInfoValidPeriodExpirDate == "" {
		return nil, status.Error(codes.InvalidArgument, "twodimension_code_info_valid_period_expir_date is required")
	}
	if req.FileUuid == "" {
		return nil, status.Error(codes.InvalidArgument, "file_uuid is required")
	}

	record, err := s.repo.GetByPrimaryKey(ctx,
		req.OrganizationId,
		req.CarId,
		req.TwodimensionCodeInfoValidPeriodExpirDate,
		req.FileUuid,
	)
	if err != nil {
		if errors.Is(err, repository.ErrCarInspectionDeregistrationFilesNotFound) {
			return nil, status.Error(codes.NotFound, "car inspection deregistration file not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get car inspection deregistration file: %v", err)
	}

	return &pb.GetCarInspectionDeregistrationFilesResponse{
		CarInspectionDeregistrationFiles: toProtoCarInspectionDeregistrationFiles(record),
	}, nil
}

// UpdateCarInspectionDeregistrationFiles updates an existing car inspection deregistration file
func (s *CarInspectionDeregistrationFilesServer) UpdateCarInspectionDeregistrationFiles(ctx context.Context, req *pb.UpdateCarInspectionDeregistrationFilesRequest) (*pb.UpdateCarInspectionDeregistrationFilesResponse, error) {
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}
	if req.CarId == "" {
		return nil, status.Error(codes.InvalidArgument, "car_id is required")
	}
	if req.TwodimensionCodeInfoValidPeriodExpirDate == "" {
		return nil, status.Error(codes.InvalidArgument, "twodimension_code_info_valid_period_expir_date is required")
	}
	if req.FileUuid == "" {
		return nil, status.Error(codes.InvalidArgument, "file_uuid is required")
	}

	// Note: This table has no updateable fields besides the primary key,
	// so this essentially validates the record exists
	record, err := s.repo.GetByPrimaryKey(ctx,
		req.OrganizationId,
		req.CarId,
		req.TwodimensionCodeInfoValidPeriodExpirDate,
		req.FileUuid,
	)
	if err != nil {
		if errors.Is(err, repository.ErrCarInspectionDeregistrationFilesNotFound) {
			return nil, status.Error(codes.NotFound, "car inspection deregistration file not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update car inspection deregistration file: %v", err)
	}

	return &pb.UpdateCarInspectionDeregistrationFilesResponse{
		CarInspectionDeregistrationFiles: toProtoCarInspectionDeregistrationFiles(record),
	}, nil
}

// DeleteCarInspectionDeregistrationFiles hard-deletes a car inspection deregistration file
func (s *CarInspectionDeregistrationFilesServer) DeleteCarInspectionDeregistrationFiles(ctx context.Context, req *pb.DeleteCarInspectionDeregistrationFilesRequest) (*pb.DeleteCarInspectionDeregistrationFilesResponse, error) {
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}
	if req.CarId == "" {
		return nil, status.Error(codes.InvalidArgument, "car_id is required")
	}
	if req.TwodimensionCodeInfoValidPeriodExpirDate == "" {
		return nil, status.Error(codes.InvalidArgument, "twodimension_code_info_valid_period_expir_date is required")
	}
	if req.FileUuid == "" {
		return nil, status.Error(codes.InvalidArgument, "file_uuid is required")
	}

	err := s.repo.Delete(ctx,
		req.OrganizationId,
		req.CarId,
		req.TwodimensionCodeInfoValidPeriodExpirDate,
		req.FileUuid,
	)
	if err != nil {
		if errors.Is(err, repository.ErrCarInspectionDeregistrationFilesNotFound) {
			return nil, status.Error(codes.NotFound, "car inspection deregistration file not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete car inspection deregistration file: %v", err)
	}

	return &pb.DeleteCarInspectionDeregistrationFilesResponse{
		Success: true,
	}, nil
}

// ListCarInspectionDeregistrationFiless retrieves car inspection deregistration files with pagination
func (s *CarInspectionDeregistrationFilesServer) ListCarInspectionDeregistrationFiless(ctx context.Context, req *pb.ListCarInspectionDeregistrationFilessRequest) (*pb.ListCarInspectionDeregistrationFilessResponse, error) {
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
		return nil, status.Errorf(codes.Internal, "failed to list car inspection deregistration files: %v", err)
	}

	var nextPageToken string
	if len(records) > limit {
		records = records[:limit]
		nextPageToken = "next" // In production, encode the offset
	}

	protoRecords := make([]*pb.CarInspectionDeregistrationFiles, len(records))
	for i, record := range records {
		protoRecords[i] = toProtoCarInspectionDeregistrationFiles(record)
	}

	return &pb.ListCarInspectionDeregistrationFilessResponse{
		CarInspectionDeregistrationFiless: protoRecords,
		NextPageToken:                     nextPageToken,
	}, nil
}

// ListCarInspectionDeregistrationFilessByOrganization retrieves car inspection deregistration files by organization with pagination
func (s *CarInspectionDeregistrationFilesServer) ListCarInspectionDeregistrationFilessByOrganization(ctx context.Context, req *pb.ListCarInspectionDeregistrationFilessByOrganizationRequest) (*pb.ListCarInspectionDeregistrationFilessByOrganizationResponse, error) {
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
		return nil, status.Errorf(codes.Internal, "failed to list car inspection deregistration files by organization: %v", err)
	}

	var nextPageToken string
	if len(records) > limit {
		records = records[:limit]
		nextPageToken = "next" // In production, encode the offset
	}

	protoRecords := make([]*pb.CarInspectionDeregistrationFiles, len(records))
	for i, record := range records {
		protoRecords[i] = toProtoCarInspectionDeregistrationFiles(record)
	}

	return &pb.ListCarInspectionDeregistrationFilessByOrganizationResponse{
		CarInspectionDeregistrationFiless: protoRecords,
		NextPageToken:                     nextPageToken,
	}, nil
}

// toProtoCarInspectionDeregistrationFiles converts repository model to proto message
func toProtoCarInspectionDeregistrationFiles(record *repository.CarInspectionDeregistrationFiles) *pb.CarInspectionDeregistrationFiles {
	return &pb.CarInspectionDeregistrationFiles{
		OrganizationId:                          record.OrganizationID,
		CarId:                                   record.CarID,
		TwodimensionCodeInfoValidPeriodExpirDate: record.TwodimensionCodeInfoValidPeriodExpirDate,
		FileUuid:                                 record.FileUUID,
	}
}
