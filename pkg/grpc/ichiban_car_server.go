package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/pb"
	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/repository"
)

// IchibanCarServer implements the gRPC IchibanCarService
type IchibanCarServer struct {
	pb.UnimplementedIchibanCarServiceServer
	repo *repository.IchibanCarRepository
}

// NewIchibanCarServer creates a new gRPC server
func NewIchibanCarServer(repo *repository.IchibanCarRepository) *IchibanCarServer {
	return &IchibanCarServer{repo: repo}
}

// CreateIchibanCar creates a new ichiban car
func (s *IchibanCarServer) CreateIchibanCar(ctx context.Context, req *pb.CreateIchibanCarRequest) (*pb.CreateIchibanCarResponse, error) {
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}
	if req.Id4 == "" {
		return nil, status.Error(codes.InvalidArgument, "id4 is required")
	}
	if req.Shashu == "" {
		return nil, status.Error(codes.InvalidArgument, "shashu is required")
	}

	// Generate ID (in production, use UUID or similar)
	id := req.Id4 // Using id4 as id for simplicity

	var name, nameR *string
	if req.Name != nil {
		name = req.Name
	}
	if req.NameR != nil {
		nameR = req.NameR
	}

	var sekisai *float64
	if req.Sekisai != nil {
		sekisai = req.Sekisai
	}

	var regDate, parchDate, scrapDate, bumonCodeID, driverID *string
	if req.RegDate != nil {
		regDate = req.RegDate
	}
	if req.ParchDate != nil {
		parchDate = req.ParchDate
	}
	if req.ScrapDate != nil {
		scrapDate = req.ScrapDate
	}
	if req.BumonCodeId != nil {
		bumonCodeID = req.BumonCodeId
	}
	if req.DriverId != nil {
		driverID = req.DriverId
	}

	car, err := s.repo.Create(ctx, id, req.OrganizationId, req.Id4, req.Shashu, name, nameR, sekisai, regDate, parchDate, scrapDate, bumonCodeID, driverID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create ichiban car: %v", err)
	}

	return &pb.CreateIchibanCarResponse{
		IchibanCar: toProtoIchibanCar(car),
	}, nil
}

// GetIchibanCar retrieves an ichiban car by composite primary key (id, organization_id)
func (s *IchibanCarServer) GetIchibanCar(ctx context.Context, req *pb.GetIchibanCarRequest) (*pb.GetIchibanCarResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}

	car, err := s.repo.GetByIDAndOrg(ctx, req.Id, req.OrganizationId)
	if err != nil {
		if errors.Is(err, repository.ErrIchibanCarNotFound) {
			return nil, status.Error(codes.NotFound, "ichiban car not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get ichiban car: %v", err)
	}

	return &pb.GetIchibanCarResponse{
		IchibanCar: toProtoIchibanCar(car),
	}, nil
}

// UpdateIchibanCar updates an existing ichiban car
func (s *IchibanCarServer) UpdateIchibanCar(ctx context.Context, req *pb.UpdateIchibanCarRequest) (*pb.UpdateIchibanCarResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}
	if req.Id4 == "" {
		return nil, status.Error(codes.InvalidArgument, "id4 is required")
	}
	if req.Shashu == "" {
		return nil, status.Error(codes.InvalidArgument, "shashu is required")
	}

	var name, nameR *string
	if req.Name != nil {
		name = req.Name
	}
	if req.NameR != nil {
		nameR = req.NameR
	}

	var sekisai *float64
	if req.Sekisai != nil {
		sekisai = req.Sekisai
	}

	var regDate, parchDate, scrapDate, bumonCodeID, driverID *string
	if req.RegDate != nil {
		regDate = req.RegDate
	}
	if req.ParchDate != nil {
		parchDate = req.ParchDate
	}
	if req.ScrapDate != nil {
		scrapDate = req.ScrapDate
	}
	if req.BumonCodeId != nil {
		bumonCodeID = req.BumonCodeId
	}
	if req.DriverId != nil {
		driverID = req.DriverId
	}

	car, err := s.repo.Update(ctx, req.Id, req.OrganizationId, req.Id4, req.Shashu, name, nameR, sekisai, regDate, parchDate, scrapDate, bumonCodeID, driverID)
	if err != nil {
		if errors.Is(err, repository.ErrIchibanCarNotFound) {
			return nil, status.Error(codes.NotFound, "ichiban car not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update ichiban car: %v", err)
	}

	return &pb.UpdateIchibanCarResponse{
		IchibanCar: toProtoIchibanCar(car),
	}, nil
}

// DeleteIchibanCar hard-deletes an ichiban car
func (s *IchibanCarServer) DeleteIchibanCar(ctx context.Context, req *pb.DeleteIchibanCarRequest) (*pb.DeleteIchibanCarResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}

	err := s.repo.Delete(ctx, req.Id, req.OrganizationId)
	if err != nil {
		if errors.Is(err, repository.ErrIchibanCarNotFound) {
			return nil, status.Error(codes.NotFound, "ichiban car not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete ichiban car: %v", err)
	}

	return &pb.DeleteIchibanCarResponse{
		Success: true,
	}, nil
}

// ListIchibanCars retrieves all ichiban cars with pagination
func (s *IchibanCarServer) ListIchibanCars(ctx context.Context, req *pb.ListIchibanCarsRequest) (*pb.ListIchibanCarsResponse, error) {
	limit := int(req.PageSize)
	if limit <= 0 {
		limit = 10
	}

	offset := 0
	if req.PageToken != "" {
		// In a real implementation, decode page_token to get offset
	}

	cars, err := s.repo.List(ctx, limit+1, offset) // +1 to check if there's next page
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list ichiban cars: %v", err)
	}

	var nextPageToken string
	if len(cars) > limit {
		cars = cars[:limit]
		nextPageToken = "next"
	}

	protoCars := make([]*pb.IchibanCar, len(cars))
	for i, car := range cars {
		protoCars[i] = toProtoIchibanCar(car)
	}

	return &pb.ListIchibanCarsResponse{
		IchibanCars:   protoCars,
		NextPageToken: nextPageToken,
	}, nil
}

// ListIchibanCarsByOrganization retrieves ichiban cars for a specific organization with pagination
func (s *IchibanCarServer) ListIchibanCarsByOrganization(ctx context.Context, req *pb.ListIchibanCarsByOrganizationRequest) (*pb.ListIchibanCarsByOrganizationResponse, error) {
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

	cars, err := s.repo.ListByOrganization(ctx, req.OrganizationId, limit+1, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list ichiban cars by organization: %v", err)
	}

	var nextPageToken string
	if len(cars) > limit {
		cars = cars[:limit]
		nextPageToken = "next"
	}

	protoCars := make([]*pb.IchibanCar, len(cars))
	for i, car := range cars {
		protoCars[i] = toProtoIchibanCar(car)
	}

	return &pb.ListIchibanCarsByOrganizationResponse{
		IchibanCars:   protoCars,
		NextPageToken: nextPageToken,
	}, nil
}

// toProtoIchibanCar converts repository model to proto message
func toProtoIchibanCar(car *repository.IchibanCar) *pb.IchibanCar {
	proto := &pb.IchibanCar{
		Id:             car.ID,
		OrganizationId: car.OrganizationID,
		Id4:            car.ID4,
		Shashu:         car.Shashu,
	}

	if car.Name != nil {
		proto.Name = car.Name
	}
	if car.NameR != nil {
		proto.NameR = car.NameR
	}
	if car.Sekisai != nil {
		proto.Sekisai = car.Sekisai
	}
	if car.RegDate != nil {
		proto.RegDate = car.RegDate
	}
	if car.ParchDate != nil {
		proto.ParchDate = car.ParchDate
	}
	if car.ScrapDate != nil {
		proto.ScrapDate = car.ScrapDate
	}
	if car.BumonCodeID != nil {
		proto.BumonCodeId = car.BumonCodeID
	}
	if car.DriverID != nil {
		proto.DriverId = car.DriverID
	}

	return proto
}
