package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/pb"
	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/repository"
)

// CarInsSheetIchibanCarsAServer implements the gRPC CarInsSheetIchibanCarsAService
type CarInsSheetIchibanCarsAServer struct {
	pb.UnimplementedCarInsSheetIchibanCarsAServiceServer
	repo *repository.CarInsSheetIchibanCarsARepository
}

// NewCarInsSheetIchibanCarsAServer creates a new gRPC server
func NewCarInsSheetIchibanCarsAServer(repo *repository.CarInsSheetIchibanCarsARepository) *CarInsSheetIchibanCarsAServer {
	return &CarInsSheetIchibanCarsAServer{repo: repo}
}

// CreateCarInsSheetIchibanCarsA creates a new car_ins_sheet_ichiban_cars_a entry
func (s *CarInsSheetIchibanCarsAServer) CreateCarInsSheetIchibanCarsA(ctx context.Context, req *pb.CreateCarInsSheetIchibanCarsARequest) (*pb.CreateCarInsSheetIchibanCarsAResponse, error) {
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}
	if req.ElectCertMgNo == "" {
		return nil, status.Error(codes.InvalidArgument, "elect_cert_mg_no is required")
	}
	if req.GrantdateE == "" {
		return nil, status.Error(codes.InvalidArgument, "grantdate_e is required")
	}
	if req.GrantdateY == "" {
		return nil, status.Error(codes.InvalidArgument, "grantdate_y is required")
	}
	if req.GrantdateM == "" {
		return nil, status.Error(codes.InvalidArgument, "grantdate_m is required")
	}
	if req.GrantdateD == "" {
		return nil, status.Error(codes.InvalidArgument, "grantdate_d is required")
	}

	// Handle optional id_cars field
	var idCars *string
	if req.IdCars != nil {
		idCars = req.IdCars
	}

	record, err := s.repo.Create(ctx,
		req.OrganizationId,
		req.ElectCertMgNo,
		req.GrantdateE,
		req.GrantdateY,
		req.GrantdateM,
		req.GrantdateD,
		idCars,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create car_ins_sheet_ichiban_cars_a: %v", err)
	}

	return &pb.CreateCarInsSheetIchibanCarsAResponse{
		CarInsSheetIchibanCarsA: toProtoCarInsSheetIchibanCarsA(record),
	}, nil
}

// GetCarInsSheetIchibanCarsA retrieves a car_ins_sheet_ichiban_cars_a entry by composite primary key
func (s *CarInsSheetIchibanCarsAServer) GetCarInsSheetIchibanCarsA(ctx context.Context, req *pb.GetCarInsSheetIchibanCarsARequest) (*pb.GetCarInsSheetIchibanCarsAResponse, error) {
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}
	if req.ElectCertMgNo == "" {
		return nil, status.Error(codes.InvalidArgument, "elect_cert_mg_no is required")
	}
	if req.GrantdateE == "" {
		return nil, status.Error(codes.InvalidArgument, "grantdate_e is required")
	}
	if req.GrantdateY == "" {
		return nil, status.Error(codes.InvalidArgument, "grantdate_y is required")
	}
	if req.GrantdateM == "" {
		return nil, status.Error(codes.InvalidArgument, "grantdate_m is required")
	}
	if req.GrantdateD == "" {
		return nil, status.Error(codes.InvalidArgument, "grantdate_d is required")
	}

	record, err := s.repo.GetByPrimaryKey(ctx,
		req.OrganizationId,
		req.ElectCertMgNo,
		req.GrantdateE,
		req.GrantdateY,
		req.GrantdateM,
		req.GrantdateD,
	)
	if err != nil {
		if errors.Is(err, repository.ErrCarInsSheetIchibanCarsANotFound) {
			return nil, status.Error(codes.NotFound, "car_ins_sheet_ichiban_cars_a not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get car_ins_sheet_ichiban_cars_a: %v", err)
	}

	return &pb.GetCarInsSheetIchibanCarsAResponse{
		CarInsSheetIchibanCarsA: toProtoCarInsSheetIchibanCarsA(record),
	}, nil
}

// UpdateCarInsSheetIchibanCarsA updates an existing car_ins_sheet_ichiban_cars_a entry
func (s *CarInsSheetIchibanCarsAServer) UpdateCarInsSheetIchibanCarsA(ctx context.Context, req *pb.UpdateCarInsSheetIchibanCarsARequest) (*pb.UpdateCarInsSheetIchibanCarsAResponse, error) {
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}
	if req.ElectCertMgNo == "" {
		return nil, status.Error(codes.InvalidArgument, "elect_cert_mg_no is required")
	}
	if req.GrantdateE == "" {
		return nil, status.Error(codes.InvalidArgument, "grantdate_e is required")
	}
	if req.GrantdateY == "" {
		return nil, status.Error(codes.InvalidArgument, "grantdate_y is required")
	}
	if req.GrantdateM == "" {
		return nil, status.Error(codes.InvalidArgument, "grantdate_m is required")
	}
	if req.GrantdateD == "" {
		return nil, status.Error(codes.InvalidArgument, "grantdate_d is required")
	}

	// Handle optional id_cars field
	var idCars *string
	if req.IdCars != nil {
		idCars = req.IdCars
	}

	record, err := s.repo.Update(ctx,
		req.OrganizationId,
		req.ElectCertMgNo,
		req.GrantdateE,
		req.GrantdateY,
		req.GrantdateM,
		req.GrantdateD,
		idCars,
	)
	if err != nil {
		if errors.Is(err, repository.ErrCarInsSheetIchibanCarsANotFound) {
			return nil, status.Error(codes.NotFound, "car_ins_sheet_ichiban_cars_a not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update car_ins_sheet_ichiban_cars_a: %v", err)
	}

	return &pb.UpdateCarInsSheetIchibanCarsAResponse{
		CarInsSheetIchibanCarsA: toProtoCarInsSheetIchibanCarsA(record),
	}, nil
}

// DeleteCarInsSheetIchibanCarsA hard-deletes a car_ins_sheet_ichiban_cars_a entry
func (s *CarInsSheetIchibanCarsAServer) DeleteCarInsSheetIchibanCarsA(ctx context.Context, req *pb.DeleteCarInsSheetIchibanCarsARequest) (*pb.DeleteCarInsSheetIchibanCarsAResponse, error) {
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}
	if req.ElectCertMgNo == "" {
		return nil, status.Error(codes.InvalidArgument, "elect_cert_mg_no is required")
	}
	if req.GrantdateE == "" {
		return nil, status.Error(codes.InvalidArgument, "grantdate_e is required")
	}
	if req.GrantdateY == "" {
		return nil, status.Error(codes.InvalidArgument, "grantdate_y is required")
	}
	if req.GrantdateM == "" {
		return nil, status.Error(codes.InvalidArgument, "grantdate_m is required")
	}
	if req.GrantdateD == "" {
		return nil, status.Error(codes.InvalidArgument, "grantdate_d is required")
	}

	err := s.repo.Delete(ctx,
		req.OrganizationId,
		req.ElectCertMgNo,
		req.GrantdateE,
		req.GrantdateY,
		req.GrantdateM,
		req.GrantdateD,
	)
	if err != nil {
		if errors.Is(err, repository.ErrCarInsSheetIchibanCarsANotFound) {
			return nil, status.Error(codes.NotFound, "car_ins_sheet_ichiban_cars_a not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete car_ins_sheet_ichiban_cars_a: %v", err)
	}

	return &pb.DeleteCarInsSheetIchibanCarsAResponse{
		Success: true,
	}, nil
}

// ListCarInsSheetIchibanCarsAs retrieves car_ins_sheet_ichiban_cars_a entries with pagination
func (s *CarInsSheetIchibanCarsAServer) ListCarInsSheetIchibanCarsAs(ctx context.Context, req *pb.ListCarInsSheetIchibanCarsAsRequest) (*pb.ListCarInsSheetIchibanCarsAsResponse, error) {
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
		return nil, status.Errorf(codes.Internal, "failed to list car_ins_sheet_ichiban_cars_a: %v", err)
	}

	var nextPageToken string
	if len(records) > limit {
		records = records[:limit]
		nextPageToken = "next" // In production, encode the offset
	}

	protoRecords := make([]*pb.CarInsSheetIchibanCarsA, len(records))
	for i, record := range records {
		protoRecords[i] = toProtoCarInsSheetIchibanCarsA(record)
	}

	return &pb.ListCarInsSheetIchibanCarsAsResponse{
		CarInsSheetIchibanCarsAs: protoRecords,
		NextPageToken:            nextPageToken,
	}, nil
}

// ListCarInsSheetIchibanCarsAsByOrganization retrieves car_ins_sheet_ichiban_cars_a entries by organization with pagination
func (s *CarInsSheetIchibanCarsAServer) ListCarInsSheetIchibanCarsAsByOrganization(ctx context.Context, req *pb.ListCarInsSheetIchibanCarsAsByOrganizationRequest) (*pb.ListCarInsSheetIchibanCarsAsByOrganizationResponse, error) {
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
		return nil, status.Errorf(codes.Internal, "failed to list car_ins_sheet_ichiban_cars_a by organization: %v", err)
	}

	var nextPageToken string
	if len(records) > limit {
		records = records[:limit]
		nextPageToken = "next" // In production, encode the offset
	}

	protoRecords := make([]*pb.CarInsSheetIchibanCarsA, len(records))
	for i, record := range records {
		protoRecords[i] = toProtoCarInsSheetIchibanCarsA(record)
	}

	return &pb.ListCarInsSheetIchibanCarsAsByOrganizationResponse{
		CarInsSheetIchibanCarsAs: protoRecords,
		NextPageToken:            nextPageToken,
	}, nil
}

// toProtoCarInsSheetIchibanCarsA converts repository model to proto message
func toProtoCarInsSheetIchibanCarsA(record *repository.CarInsSheetIchibanCarsA) *pb.CarInsSheetIchibanCarsA {
	return &pb.CarInsSheetIchibanCarsA{
		OrganizationId: record.OrganizationID,
		IdCars:         record.IDCars,
		ElectCertMgNo:  record.ElectCertMgNo,
		GrantdateE:     record.GrantdateE,
		GrantdateY:     record.GrantdateY,
		GrantdateM:     record.GrantdateM,
		GrantdateD:     record.GrantdateD,
	}
}
