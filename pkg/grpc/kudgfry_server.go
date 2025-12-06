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

// KudgfryServer implements the gRPC KudgfryService
type KudgfryServer struct {
	pb.UnimplementedKudgfryServiceServer
	repo *repository.KudgfryRepository
}

// NewKudgfryServer creates a new gRPC server
func NewKudgfryServer(repo *repository.KudgfryRepository) *KudgfryServer {
	return &KudgfryServer{repo: repo}
}

// CreateKudgfry creates a new kudgfry record
func (s *KudgfryServer) CreateKudgfry(ctx context.Context, req *pb.CreateKudgfryRequest) (*pb.CreateKudgfryResponse, error) {
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}
	if req.Hash == "" {
		return nil, status.Error(codes.InvalidArgument, "hash is required")
	}
	if req.TargetDriverType == "" {
		return nil, status.Error(codes.InvalidArgument, "target_driver_type is required")
	}

	kudgfry := &repository.Kudgfry{
		OrganizationID:            req.OrganizationId,
		Hash:                      req.Hash,
		Created:                   req.Created,
		Deleted:                   ptrFromOptional(req.Deleted),
		KudguriUuid:               ptrFromOptional(req.KudguriUuid),
		TargetDriverType:          req.TargetDriverType,
		UnkouNo:                   ptrFromOptional(req.UnkouNo),
		UnkouDate:                 ptrFromOptional(req.UnkouDate),
		ReadDate:                  ptrFromOptional(req.ReadDate),
		OfficeCd:                  ptrFromOptional(req.OfficeCd),
		OfficeName:                ptrFromOptional(req.OfficeName),
		VehicleCd:                 ptrFromOptional(req.VehicleCd),
		VehicleName:               ptrFromOptional(req.VehicleName),
		DriverCd1:                 ptrFromOptional(req.DriverCd1),
		DriverName1:               ptrFromOptional(req.DriverName1),
		DriverCd2:                 ptrFromOptional(req.DriverCd2),
		DriverName2:               ptrFromOptional(req.DriverName2),
		RelevantDatetime:          ptrFromOptional(req.RelevantDatetime),
		RefuelInspectCategory:     ptrFromOptional(req.RefuelInspectCategory),
		RefuelInspectCategoryName: ptrFromOptional(req.RefuelInspectCategoryName),
		RefuelInspectType:         ptrFromOptional(req.RefuelInspectType),
		RefuelInspectTypeName:     ptrFromOptional(req.RefuelInspectTypeName),
		RefuelInspectKind:         ptrFromOptional(req.RefuelInspectKind),
		RefuelInspectKindName:     ptrFromOptional(req.RefuelInspectKindName),
		RefillAmount:              ptrFromOptional(req.RefillAmount),
		OwnOtherType:              ptrFromOptional(req.OwnOtherType),
		Mileage:                   ptrFromOptional(req.Mileage),
		MeterValue:                ptrFromOptional(req.MeterValue),
	}

	result, err := s.repo.Create(ctx, kudgfry)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create kudgfry: %v", err)
	}

	return &pb.CreateKudgfryResponse{
		Kudgfry: toProtoKudgfry(result),
	}, nil
}

// GetKudgfry retrieves a kudgfry record by UUID
func (s *KudgfryServer) GetKudgfry(ctx context.Context, req *pb.GetKudgfryRequest) (*pb.GetKudgfryResponse, error) {
	if req.Uuid == "" {
		return nil, status.Error(codes.InvalidArgument, "uuid is required")
	}

	kudgfry, err := s.repo.GetByUUID(ctx, req.Uuid)
	if err != nil {
		if errors.Is(err, repository.ErrKudgfryNotFound) {
			return nil, status.Error(codes.NotFound, "kudgfry not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get kudgfry: %v", err)
	}

	return &pb.GetKudgfryResponse{
		Kudgfry: toProtoKudgfry(kudgfry),
	}, nil
}

// UpdateKudgfry updates an existing kudgfry record
func (s *KudgfryServer) UpdateKudgfry(ctx context.Context, req *pb.UpdateKudgfryRequest) (*pb.UpdateKudgfryResponse, error) {
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

	kudgfry := &repository.Kudgfry{
		UUID:                      req.Uuid,
		OrganizationID:            req.OrganizationId,
		Hash:                      req.Hash,
		KudguriUuid:               ptrFromOptional(req.KudguriUuid),
		TargetDriverType:          req.TargetDriverType,
		UnkouNo:                   ptrFromOptional(req.UnkouNo),
		UnkouDate:                 ptrFromOptional(req.UnkouDate),
		ReadDate:                  ptrFromOptional(req.ReadDate),
		OfficeCd:                  ptrFromOptional(req.OfficeCd),
		OfficeName:                ptrFromOptional(req.OfficeName),
		VehicleCd:                 ptrFromOptional(req.VehicleCd),
		VehicleName:               ptrFromOptional(req.VehicleName),
		DriverCd1:                 ptrFromOptional(req.DriverCd1),
		DriverName1:               ptrFromOptional(req.DriverName1),
		DriverCd2:                 ptrFromOptional(req.DriverCd2),
		DriverName2:               ptrFromOptional(req.DriverName2),
		RelevantDatetime:          ptrFromOptional(req.RelevantDatetime),
		RefuelInspectCategory:     ptrFromOptional(req.RefuelInspectCategory),
		RefuelInspectCategoryName: ptrFromOptional(req.RefuelInspectCategoryName),
		RefuelInspectType:         ptrFromOptional(req.RefuelInspectType),
		RefuelInspectTypeName:     ptrFromOptional(req.RefuelInspectTypeName),
		RefuelInspectKind:         ptrFromOptional(req.RefuelInspectKind),
		RefuelInspectKindName:     ptrFromOptional(req.RefuelInspectKindName),
		RefillAmount:              ptrFromOptional(req.RefillAmount),
		OwnOtherType:              ptrFromOptional(req.OwnOtherType),
		Mileage:                   ptrFromOptional(req.Mileage),
		MeterValue:                ptrFromOptional(req.MeterValue),
	}

	result, err := s.repo.Update(ctx, kudgfry)
	if err != nil {
		if errors.Is(err, repository.ErrKudgfryNotFound) {
			return nil, status.Error(codes.NotFound, "kudgfry not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update kudgfry: %v", err)
	}

	return &pb.UpdateKudgfryResponse{
		Kudgfry: toProtoKudgfry(result),
	}, nil
}

// DeleteKudgfry soft-deletes a kudgfry record
func (s *KudgfryServer) DeleteKudgfry(ctx context.Context, req *pb.DeleteKudgfryRequest) (*pb.DeleteKudgfryResponse, error) {
	if req.Uuid == "" {
		return nil, status.Error(codes.InvalidArgument, "uuid is required")
	}

	err := s.repo.Delete(ctx, req.Uuid)
	if err != nil {
		if errors.Is(err, repository.ErrKudgfryNotFound) {
			return nil, status.Error(codes.NotFound, "kudgfry not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete kudgfry: %v", err)
	}

	return &pb.DeleteKudgfryResponse{
		Success: true,
	}, nil
}

// ListKudgfrys retrieves kudgfry records with pagination
func (s *KudgfryServer) ListKudgfrys(ctx context.Context, req *pb.ListKudgfrysRequest) (*pb.ListKudgfrysResponse, error) {
	limit := int(req.PageSize)
	if limit <= 0 {
		limit = 10
	}

	offset := 0
	if req.PageToken != "" {
		// In production, decode page_token to get offset
	}

	kudgfrys, err := s.repo.List(ctx, limit+1, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list kudgfrys: %v", err)
	}

	var nextPageToken string
	if len(kudgfrys) > limit {
		kudgfrys = kudgfrys[:limit]
		nextPageToken = "next"
	}

	protoKudgfrys := make([]*pb.Kudgfry, len(kudgfrys))
	for i, kudgfry := range kudgfrys {
		protoKudgfrys[i] = toProtoKudgfry(kudgfry)
	}

	return &pb.ListKudgfrysResponse{
		Kudgfrys:      protoKudgfrys,
		NextPageToken: nextPageToken,
	}, nil
}

// ListKudgfrysByOrganization retrieves kudgfry records for a specific organization with pagination
func (s *KudgfryServer) ListKudgfrysByOrganization(ctx context.Context, req *pb.ListKudgfrysByOrganizationRequest) (*pb.ListKudgfrysByOrganizationResponse, error) {
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

	kudgfrys, err := s.repo.ListByOrganization(ctx, req.OrganizationId, limit+1, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list kudgfrys by organization: %v", err)
	}

	var nextPageToken string
	if len(kudgfrys) > limit {
		kudgfrys = kudgfrys[:limit]
		nextPageToken = "next"
	}

	protoKudgfrys := make([]*pb.Kudgfry, len(kudgfrys))
	for i, kudgfry := range kudgfrys {
		protoKudgfrys[i] = toProtoKudgfry(kudgfry)
	}

	return &pb.ListKudgfrysByOrganizationResponse{
		Kudgfrys:      protoKudgfrys,
		NextPageToken: nextPageToken,
	}, nil
}

// toProtoKudgfry converts repository model to proto message
func toProtoKudgfry(k *repository.Kudgfry) *pb.Kudgfry {
	return &pb.Kudgfry{
		Uuid:                      k.UUID,
		OrganizationId:            k.OrganizationID,
		Hash:                      k.Hash,
		Created:                   k.Created,
		Deleted:                   optionalFromPtr(k.Deleted),
		KudguriUuid:               optionalFromPtr(k.KudguriUuid),
		TargetDriverType:          k.TargetDriverType,
		UnkouNo:                   optionalFromPtr(k.UnkouNo),
		UnkouDate:                 optionalFromPtr(k.UnkouDate),
		ReadDate:                  optionalFromPtr(k.ReadDate),
		OfficeCd:                  optionalFromPtr(k.OfficeCd),
		OfficeName:                optionalFromPtr(k.OfficeName),
		VehicleCd:                 optionalFromPtr(k.VehicleCd),
		VehicleName:               optionalFromPtr(k.VehicleName),
		DriverCd1:                 optionalFromPtr(k.DriverCd1),
		DriverName1:               optionalFromPtr(k.DriverName1),
		DriverCd2:                 optionalFromPtr(k.DriverCd2),
		DriverName2:               optionalFromPtr(k.DriverName2),
		RelevantDatetime:          optionalFromPtr(k.RelevantDatetime),
		RefuelInspectCategory:     optionalFromPtr(k.RefuelInspectCategory),
		RefuelInspectCategoryName: optionalFromPtr(k.RefuelInspectCategoryName),
		RefuelInspectType:         optionalFromPtr(k.RefuelInspectType),
		RefuelInspectTypeName:     optionalFromPtr(k.RefuelInspectTypeName),
		RefuelInspectKind:         optionalFromPtr(k.RefuelInspectKind),
		RefuelInspectKindName:     optionalFromPtr(k.RefuelInspectKindName),
		RefillAmount:              optionalFromPtr(k.RefillAmount),
		OwnOtherType:              optionalFromPtr(k.OwnOtherType),
		Mileage:                   optionalFromPtr(k.Mileage),
		MeterValue:                optionalFromPtr(k.MeterValue),
	}
}

// Helper functions for optional field conversion
func ptrFromOptional(opt *string) *string {
	return opt
}

func optionalFromPtr(ptr *string) *string {
	return ptr
}

func ptrString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func stringFromPtr(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

func timeNow() string {
	return time.Now().Format(time.RFC3339)
}
