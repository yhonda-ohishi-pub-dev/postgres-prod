package grpc

import (
	"context"
	"errors"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/pb"
	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/repository"
)

// KudgcstServer implements the gRPC KudgcstService
type KudgcstServer struct {
	pb.UnimplementedKudgcstServiceServer
	repo *repository.KudgcstRepository
}

// NewKudgcstServer creates a new gRPC server
func NewKudgcstServer(repo *repository.KudgcstRepository) *KudgcstServer {
	return &KudgcstServer{repo: repo}
}

// CreateKudgcst creates a new kudgcst record
func (s *KudgcstServer) CreateKudgcst(ctx context.Context, req *pb.CreateKudgcstRequest) (*pb.CreateKudgcstResponse, error) {
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}
	if req.Hash == "" {
		return nil, status.Error(codes.InvalidArgument, "hash is required")
	}
	if req.Created == "" {
		return nil, status.Error(codes.InvalidArgument, "created is required")
	}
	if req.TargetDriverType == "" {
		return nil, status.Error(codes.InvalidArgument, "target_driver_type is required")
	}

	kudgcst := &repository.Kudgcst{
		OrganizationID:       req.OrganizationId,
		Hash:                 req.Hash,
		Created:              req.Created,
		Deleted:              ptrFromOptional(req.Deleted),
		KudguriUuid:          ptrFromOptional(req.KudguriUuid),
		UnkouNo:              ptrFromOptional(req.UnkouNo),
		UnkouDate:            ptrFromOptional(req.UnkouDate),
		ReadDate:             ptrFromOptional(req.ReadDate),
		OfficeCd:             ptrFromOptional(req.OfficeCd),
		OfficeName:           ptrFromOptional(req.OfficeName),
		VehicleCd:            ptrFromOptional(req.VehicleCd),
		VehicleName:          ptrFromOptional(req.VehicleName),
		DriverCd1:            ptrFromOptional(req.DriverCd1),
		DriverName1:          ptrFromOptional(req.DriverName1),
		TargetDriverType:     req.TargetDriverType,
		StartDatetime:        ptrFromOptional(req.StartDatetime),
		EndDatetime:          ptrFromOptional(req.EndDatetime),
		FerryCompanyCd:       ptrFromOptional(req.FerryCompanyCd),
		FerryCompanyName:     ptrFromOptional(req.FerryCompanyName),
		BoardingPlaceCd:      ptrFromOptional(req.BoardingPlaceCd),
		BoardingPlaceName:    ptrFromOptional(req.BoardingPlaceName),
		TripNumber:           ptrFromOptional(req.TripNumber),
		DropoffPlaceCd:       ptrFromOptional(req.DropoffPlaceCd),
		DropoffPlaceName:     ptrFromOptional(req.DropoffPlaceName),
		SettlementType:       ptrFromOptional(req.SettlementType),
		SettlementTypeName:   ptrFromOptional(req.SettlementTypeName),
		StandardFare:         ptrFromOptional(req.StandardFare),
		ContractFare:         ptrFromOptional(req.ContractFare),
		FerryVehicleType:     nil,
		FerryVehicleTypeName: nil,
		AssumedDistance:      nil,
	}

	result, err := s.repo.Create(ctx, kudgcst)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create kudgcst: %v", err)
	}

	return &pb.CreateKudgcstResponse{
		Kudgcst: toProtoKudgcst(result),
	}, nil
}

// GetKudgcst retrieves a kudgcst record by UUID
func (s *KudgcstServer) GetKudgcst(ctx context.Context, req *pb.GetKudgcstRequest) (*pb.GetKudgcstResponse, error) {
	if req.Uuid == "" {
		return nil, status.Error(codes.InvalidArgument, "uuid is required")
	}

	kudgcst, err := s.repo.GetByUUID(ctx, req.Uuid)
	if err != nil {
		if errors.Is(err, repository.ErrKudgcstNotFound) {
			return nil, status.Error(codes.NotFound, "kudgcst not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get kudgcst: %v", err)
	}

	return &pb.GetKudgcstResponse{
		Kudgcst: toProtoKudgcst(kudgcst),
	}, nil
}

// UpdateKudgcst updates an existing kudgcst record
func (s *KudgcstServer) UpdateKudgcst(ctx context.Context, req *pb.UpdateKudgcstRequest) (*pb.UpdateKudgcstResponse, error) {
	if req.Uuid == "" {
		return nil, status.Error(codes.InvalidArgument, "uuid is required")
	}
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}
	if req.Hash == "" {
		return nil, status.Error(codes.InvalidArgument, "hash is required")
	}
	if req.TargetDriverType == "" {
		return nil, status.Error(codes.InvalidArgument, "target_driver_type is required")
	}

	// Get the existing record first to preserve created timestamp
	existing, err := s.repo.GetByUUID(ctx, req.Uuid)
	if err != nil {
		if errors.Is(err, repository.ErrKudgcstNotFound) {
			return nil, status.Error(codes.NotFound, "kudgcst not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get kudgcst: %v", err)
	}

	kudgcst := &repository.Kudgcst{
		UUID:               req.Uuid,
		OrganizationID:     req.OrganizationId,
		Hash:               req.Hash,
		Created:            req.Created,
		Deleted:            ptrFromOptional(req.Deleted),
		KudguriUuid:        ptrFromOptional(req.KudguriUuid),
		UnkouNo:            ptrFromOptional(req.UnkouNo),
		UnkouDate:          ptrFromOptional(req.UnkouDate),
		ReadDate:           ptrFromOptional(req.ReadDate),
		OfficeCd:           ptrFromOptional(req.OfficeCd),
		OfficeName:         ptrFromOptional(req.OfficeName),
		VehicleCd:          ptrFromOptional(req.VehicleCd),
		VehicleName:        ptrFromOptional(req.VehicleName),
		DriverCd1:          ptrFromOptional(req.DriverCd1),
		DriverName1:        ptrFromOptional(req.DriverName1),
		TargetDriverType:   req.TargetDriverType,
		StartDatetime:      ptrFromOptional(req.StartDatetime),
		EndDatetime:        ptrFromOptional(req.EndDatetime),
		FerryCompanyCd:     ptrFromOptional(req.FerryCompanyCd),
		FerryCompanyName:   ptrFromOptional(req.FerryCompanyName),
		BoardingPlaceCd:    ptrFromOptional(req.BoardingPlaceCd),
		BoardingPlaceName:  ptrFromOptional(req.BoardingPlaceName),
		TripNumber:         ptrFromOptional(req.TripNumber),
		DropoffPlaceCd:     ptrFromOptional(req.DropoffPlaceCd),
		DropoffPlaceName:   ptrFromOptional(req.DropoffPlaceName),
		SettlementType:     ptrFromOptional(req.SettlementType),
		SettlementTypeName: ptrFromOptional(req.SettlementTypeName),
		StandardFare:       ptrFromOptional(req.StandardFare),
		ContractFare:       ptrFromOptional(req.ContractFare),
		// Preserve immutable fields from existing record
		FerryVehicleType:     existing.FerryVehicleType,
		FerryVehicleTypeName: existing.FerryVehicleTypeName,
		AssumedDistance:      existing.AssumedDistance,
	}

	result, err := s.repo.Update(ctx, kudgcst)
	if err != nil {
		if errors.Is(err, repository.ErrKudgcstNotFound) {
			return nil, status.Error(codes.NotFound, "kudgcst not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update kudgcst: %v", err)
	}

	return &pb.UpdateKudgcstResponse{
		Kudgcst: toProtoKudgcst(result),
	}, nil
}

// DeleteKudgcst soft-deletes a kudgcst record
func (s *KudgcstServer) DeleteKudgcst(ctx context.Context, req *pb.DeleteKudgcstRequest) (*pb.DeleteKudgcstResponse, error) {
	if req.Uuid == "" {
		return nil, status.Error(codes.InvalidArgument, "uuid is required")
	}

	deletedAt := time.Now().Format(time.RFC3339)
	err := s.repo.Delete(ctx, req.Uuid, deletedAt)
	if err != nil {
		if errors.Is(err, repository.ErrKudgcstNotFound) {
			return nil, status.Error(codes.NotFound, "kudgcst not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete kudgcst: %v", err)
	}

	return &pb.DeleteKudgcstResponse{
		Success: true,
	}, nil
}

// ListKudgcsts retrieves kudgcst records with pagination
func (s *KudgcstServer) ListKudgcsts(ctx context.Context, req *pb.ListKudgcstsRequest) (*pb.ListKudgcstsResponse, error) {
	limit := int(req.PageSize)
	if limit <= 0 {
		limit = 10
	}

	offset := 0
	if req.PageToken != "" {
		// In production, decode page_token to get offset
	}

	kudgcsts, err := s.repo.List(ctx, limit+1, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list kudgcsts: %v", err)
	}

	var nextPageToken string
	if len(kudgcsts) > limit {
		kudgcsts = kudgcsts[:limit]
		nextPageToken = "next"
	}

	protoKudgcsts := make([]*pb.Kudgcst, len(kudgcsts))
	for i, kudgcst := range kudgcsts {
		protoKudgcsts[i] = toProtoKudgcst(kudgcst)
	}

	return &pb.ListKudgcstsResponse{
		Kudgcsts:      protoKudgcsts,
		NextPageToken: nextPageToken,
	}, nil
}

// ListKudgcstsByOrganization retrieves kudgcst records for a specific organization with pagination
func (s *KudgcstServer) ListKudgcstsByOrganization(ctx context.Context, req *pb.ListKudgcstsByOrganizationRequest) (*pb.ListKudgcstsByOrganizationResponse, error) {
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}

	limit := int(req.PageSize)
	if limit <= 0 {
		limit = 10
	}

	offset := 0
	if req.PageToken != "" {
		// In production, decode page_token to get offset
	}

	kudgcsts, err := s.repo.ListByOrganization(ctx, req.OrganizationId, limit+1, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list kudgcsts by organization: %v", err)
	}

	var nextPageToken string
	if len(kudgcsts) > limit {
		kudgcsts = kudgcsts[:limit]
		nextPageToken = "next"
	}

	protoKudgcsts := make([]*pb.Kudgcst, len(kudgcsts))
	for i, kudgcst := range kudgcsts {
		protoKudgcsts[i] = toProtoKudgcst(kudgcst)
	}

	return &pb.ListKudgcstsByOrganizationResponse{
		Kudgcsts:      protoKudgcsts,
		NextPageToken: nextPageToken,
	}, nil
}

// toProtoKudgcst converts repository model to proto message
func toProtoKudgcst(k *repository.Kudgcst) *pb.Kudgcst {
	return &pb.Kudgcst{
		Uuid:              k.UUID,
		OrganizationId:    k.OrganizationID,
		Hash:              k.Hash,
		Created:           k.Created,
		Deleted:           optionalFromPtr(k.Deleted),
		KudguriUuid:       optionalFromPtr(k.KudguriUuid),
		UnkouNo:           optionalFromPtr(k.UnkouNo),
		UnkouDate:         optionalFromPtr(k.UnkouDate),
		ReadDate:          optionalFromPtr(k.ReadDate),
		OfficeCd:          optionalFromPtr(k.OfficeCd),
		OfficeName:        optionalFromPtr(k.OfficeName),
		VehicleCd:         optionalFromPtr(k.VehicleCd),
		VehicleName:       optionalFromPtr(k.VehicleName),
		DriverCd1:         optionalFromPtr(k.DriverCd1),
		DriverName1:       optionalFromPtr(k.DriverName1),
		TargetDriverType:  k.TargetDriverType,
		StartDatetime:     optionalFromPtr(k.StartDatetime),
		EndDatetime:       optionalFromPtr(k.EndDatetime),
		FerryCompanyCd:    optionalFromPtr(k.FerryCompanyCd),
		FerryCompanyName:  optionalFromPtr(k.FerryCompanyName),
		BoardingPlaceCd:   optionalFromPtr(k.BoardingPlaceCd),
		BoardingPlaceName: optionalFromPtr(k.BoardingPlaceName),
		TripNumber:        optionalFromPtr(k.TripNumber),
		DropoffPlaceCd:    optionalFromPtr(k.DropoffPlaceCd),
		DropoffPlaceName:  optionalFromPtr(k.DropoffPlaceName),
		SettlementType:    optionalFromPtr(k.SettlementType),
		SettlementTypeName: optionalFromPtr(k.SettlementTypeName),
		StandardFare:      optionalFromPtr(k.StandardFare),
		ContractFare:      optionalFromPtr(k.ContractFare),
	}
}
