package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/pb"
	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/repository"
)

// CarInsSheetIchibanCarsServer implements the gRPC CarInsSheetIchibanCarsService
type CarInsSheetIchibanCarsServer struct {
	pb.UnimplementedCarInsSheetIchibanCarsServiceServer
	repo *repository.CarInsSheetIchibanCarsRepository
}

// NewCarInsSheetIchibanCarsServer creates a new gRPC server
func NewCarInsSheetIchibanCarsServer(repo *repository.CarInsSheetIchibanCarsRepository) *CarInsSheetIchibanCarsServer {
	return &CarInsSheetIchibanCarsServer{repo: repo}
}

// CreateCarInsSheetIchibanCars creates a new car_ins_sheet_ichiban_cars record
func (s *CarInsSheetIchibanCarsServer) CreateCarInsSheetIchibanCars(ctx context.Context, req *pb.CreateCarInsSheetIchibanCarsRequest) (*pb.CreateCarInsSheetIchibanCarsResponse, error) {
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}
	if req.ElectCertMgNo == "" {
		return nil, status.Error(codes.InvalidArgument, "elect_cert_mg_no is required")
	}
	if req.ElectCertPublishdateE == "" {
		return nil, status.Error(codes.InvalidArgument, "elect_cert_publishdate_e is required")
	}
	if req.ElectCertPublishdateY == "" {
		return nil, status.Error(codes.InvalidArgument, "elect_cert_publishdate_y is required")
	}
	if req.ElectCertPublishdateM == "" {
		return nil, status.Error(codes.InvalidArgument, "elect_cert_publishdate_m is required")
	}
	if req.ElectCertPublishdateD == "" {
		return nil, status.Error(codes.InvalidArgument, "elect_cert_publishdate_d is required")
	}

	// Handle optional id_cars field
	var idCars *string
	if req.IdCars != nil {
		idCars = req.IdCars
	}

	record, err := s.repo.Create(ctx,
		req.OrganizationId,
		req.ElectCertMgNo,
		req.ElectCertPublishdateE,
		req.ElectCertPublishdateY,
		req.ElectCertPublishdateM,
		req.ElectCertPublishdateD,
		idCars,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create car_ins_sheet_ichiban_cars: %v", err)
	}

	return &pb.CreateCarInsSheetIchibanCarsResponse{
		CarInsSheetIchibanCars: toProtoCarInsSheetIchibanCars(record),
	}, nil
}

// GetCarInsSheetIchibanCars retrieves a car_ins_sheet_ichiban_cars record by composite primary key
func (s *CarInsSheetIchibanCarsServer) GetCarInsSheetIchibanCars(ctx context.Context, req *pb.GetCarInsSheetIchibanCarsRequest) (*pb.GetCarInsSheetIchibanCarsResponse, error) {
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}
	if req.ElectCertMgNo == "" {
		return nil, status.Error(codes.InvalidArgument, "elect_cert_mg_no is required")
	}
	if req.ElectCertPublishdateE == "" {
		return nil, status.Error(codes.InvalidArgument, "elect_cert_publishdate_e is required")
	}
	if req.ElectCertPublishdateY == "" {
		return nil, status.Error(codes.InvalidArgument, "elect_cert_publishdate_y is required")
	}
	if req.ElectCertPublishdateM == "" {
		return nil, status.Error(codes.InvalidArgument, "elect_cert_publishdate_m is required")
	}
	if req.ElectCertPublishdateD == "" {
		return nil, status.Error(codes.InvalidArgument, "elect_cert_publishdate_d is required")
	}

	record, err := s.repo.GetByPrimaryKey(ctx,
		req.OrganizationId,
		req.ElectCertMgNo,
		req.ElectCertPublishdateE,
		req.ElectCertPublishdateY,
		req.ElectCertPublishdateM,
		req.ElectCertPublishdateD,
	)
	if err != nil {
		if errors.Is(err, repository.ErrCarInsSheetIchibanCarsNotFound) {
			return nil, status.Error(codes.NotFound, "car_ins_sheet_ichiban_cars not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get car_ins_sheet_ichiban_cars: %v", err)
	}

	return &pb.GetCarInsSheetIchibanCarsResponse{
		CarInsSheetIchibanCars: toProtoCarInsSheetIchibanCars(record),
	}, nil
}

// UpdateCarInsSheetIchibanCars updates an existing car_ins_sheet_ichiban_cars record
func (s *CarInsSheetIchibanCarsServer) UpdateCarInsSheetIchibanCars(ctx context.Context, req *pb.UpdateCarInsSheetIchibanCarsRequest) (*pb.UpdateCarInsSheetIchibanCarsResponse, error) {
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}
	if req.ElectCertMgNo == "" {
		return nil, status.Error(codes.InvalidArgument, "elect_cert_mg_no is required")
	}
	if req.ElectCertPublishdateE == "" {
		return nil, status.Error(codes.InvalidArgument, "elect_cert_publishdate_e is required")
	}
	if req.ElectCertPublishdateY == "" {
		return nil, status.Error(codes.InvalidArgument, "elect_cert_publishdate_y is required")
	}
	if req.ElectCertPublishdateM == "" {
		return nil, status.Error(codes.InvalidArgument, "elect_cert_publishdate_m is required")
	}
	if req.ElectCertPublishdateD == "" {
		return nil, status.Error(codes.InvalidArgument, "elect_cert_publishdate_d is required")
	}

	// Handle optional id_cars field
	var idCars *string
	if req.IdCars != nil {
		idCars = req.IdCars
	}

	record, err := s.repo.Update(ctx,
		req.OrganizationId,
		req.ElectCertMgNo,
		req.ElectCertPublishdateE,
		req.ElectCertPublishdateY,
		req.ElectCertPublishdateM,
		req.ElectCertPublishdateD,
		idCars,
	)
	if err != nil {
		if errors.Is(err, repository.ErrCarInsSheetIchibanCarsNotFound) {
			return nil, status.Error(codes.NotFound, "car_ins_sheet_ichiban_cars not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update car_ins_sheet_ichiban_cars: %v", err)
	}

	return &pb.UpdateCarInsSheetIchibanCarsResponse{
		CarInsSheetIchibanCars: toProtoCarInsSheetIchibanCars(record),
	}, nil
}

// DeleteCarInsSheetIchibanCars hard-deletes a car_ins_sheet_ichiban_cars record
func (s *CarInsSheetIchibanCarsServer) DeleteCarInsSheetIchibanCars(ctx context.Context, req *pb.DeleteCarInsSheetIchibanCarsRequest) (*pb.DeleteCarInsSheetIchibanCarsResponse, error) {
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}
	if req.ElectCertMgNo == "" {
		return nil, status.Error(codes.InvalidArgument, "elect_cert_mg_no is required")
	}
	if req.ElectCertPublishdateE == "" {
		return nil, status.Error(codes.InvalidArgument, "elect_cert_publishdate_e is required")
	}
	if req.ElectCertPublishdateY == "" {
		return nil, status.Error(codes.InvalidArgument, "elect_cert_publishdate_y is required")
	}
	if req.ElectCertPublishdateM == "" {
		return nil, status.Error(codes.InvalidArgument, "elect_cert_publishdate_m is required")
	}
	if req.ElectCertPublishdateD == "" {
		return nil, status.Error(codes.InvalidArgument, "elect_cert_publishdate_d is required")
	}

	err := s.repo.Delete(ctx,
		req.OrganizationId,
		req.ElectCertMgNo,
		req.ElectCertPublishdateE,
		req.ElectCertPublishdateY,
		req.ElectCertPublishdateM,
		req.ElectCertPublishdateD,
	)
	if err != nil {
		if errors.Is(err, repository.ErrCarInsSheetIchibanCarsNotFound) {
			return nil, status.Error(codes.NotFound, "car_ins_sheet_ichiban_cars not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete car_ins_sheet_ichiban_cars: %v", err)
	}

	return &pb.DeleteCarInsSheetIchibanCarsResponse{
		Success: true,
	}, nil
}

// ListCarInsSheetIchibanCarss retrieves car_ins_sheet_ichiban_cars records with pagination
func (s *CarInsSheetIchibanCarsServer) ListCarInsSheetIchibanCarss(ctx context.Context, req *pb.ListCarInsSheetIchibanCarssRequest) (*pb.ListCarInsSheetIchibanCarssResponse, error) {
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
		return nil, status.Errorf(codes.Internal, "failed to list car_ins_sheet_ichiban_cars: %v", err)
	}

	var nextPageToken string
	if len(records) > limit {
		records = records[:limit]
		nextPageToken = "next" // In production, encode the offset
	}

	protoRecords := make([]*pb.CarInsSheetIchibanCars, len(records))
	for i, record := range records {
		protoRecords[i] = toProtoCarInsSheetIchibanCars(record)
	}

	return &pb.ListCarInsSheetIchibanCarssResponse{
		CarInsSheetIchibanCarss: protoRecords,
		NextPageToken:           nextPageToken,
	}, nil
}

// ListCarInsSheetIchibanCarssByOrganization retrieves car_ins_sheet_ichiban_cars records by organization with pagination
func (s *CarInsSheetIchibanCarsServer) ListCarInsSheetIchibanCarssByOrganization(ctx context.Context, req *pb.ListCarInsSheetIchibanCarssByOrganizationRequest) (*pb.ListCarInsSheetIchibanCarssByOrganizationResponse, error) {
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
		return nil, status.Errorf(codes.Internal, "failed to list car_ins_sheet_ichiban_cars by organization: %v", err)
	}

	var nextPageToken string
	if len(records) > limit {
		records = records[:limit]
		nextPageToken = "next" // In production, encode the offset
	}

	protoRecords := make([]*pb.CarInsSheetIchibanCars, len(records))
	for i, record := range records {
		protoRecords[i] = toProtoCarInsSheetIchibanCars(record)
	}

	return &pb.ListCarInsSheetIchibanCarssByOrganizationResponse{
		CarInsSheetIchibanCarss: protoRecords,
		NextPageToken:           nextPageToken,
	}, nil
}

// toProtoCarInsSheetIchibanCars converts repository model to proto message
func toProtoCarInsSheetIchibanCars(record *repository.CarInsSheetIchibanCars) *pb.CarInsSheetIchibanCars {
	return &pb.CarInsSheetIchibanCars{
		OrganizationId:         record.OrganizationID,
		IdCars:                 record.IDCars,
		ElectCertMgNo:          record.ElectCertMgNo,
		ElectCertPublishdateE:  record.ElectCertPublishdateE,
		ElectCertPublishdateY:  record.ElectCertPublishdateY,
		ElectCertPublishdateM:  record.ElectCertPublishdateM,
		ElectCertPublishdateD:  record.ElectCertPublishdateD,
	}
}
