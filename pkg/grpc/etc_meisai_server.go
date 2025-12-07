package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/db"
	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/pb"
	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/repository"
)

// ETCMeisaiServer implements the gRPC ETCMeisaiService
type ETCMeisaiServer struct {
	pb.UnimplementedETCMeisaiServiceServer
	repo *repository.ETCMeisaiRepository
}

// NewETCMeisaiServer creates a new gRPC server
func NewETCMeisaiServer(repo *repository.ETCMeisaiRepository) *ETCMeisaiServer {
	return &ETCMeisaiServer{repo: repo}
}

// CreateETCMeisai creates a new ETC meisai record
func (s *ETCMeisaiServer) CreateETCMeisai(ctx context.Context, req *pb.CreateETCMeisaiRequest) (*pb.CreateETCMeisaiResponse, error) {
	if req.DateTo == nil {
		return nil, status.Error(codes.InvalidArgument, "date_to is required")
	}
	if req.DateToDate == "" {
		return nil, status.Error(codes.InvalidArgument, "date_to_date is required")
	}
	if req.IcFr == "" {
		return nil, status.Error(codes.InvalidArgument, "ic_fr is required")
	}
	if req.IcTo == "" {
		return nil, status.Error(codes.InvalidArgument, "ic_to is required")
	}
	if req.EtcNum == "" {
		return nil, status.Error(codes.InvalidArgument, "etc_num is required")
	}
	if req.Hash == "" {
		return nil, status.Error(codes.InvalidArgument, "hash is required")
	}

	// Get organization_id from context (set by RLS interceptor)
	orgID, ok := db.GetOrganizationID(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required in header")
	}

	meisai := &repository.ETCMeisai{
		DateTo:     req.DateTo.AsTime(),
		DateToDate: req.DateToDate,
		IcFr:       req.IcFr,
		IcTo:       req.IcTo,
		Price:      req.Price,
		Shashu:     req.Shashu,
		EtcNum:     req.EtcNum,
		Hash:       req.Hash,
	}

	if req.DateFr != nil {
		t := req.DateFr.AsTime()
		meisai.DateFr = &t
	}
	if req.PriceBf != nil {
		meisai.PriceBf = req.PriceBf
	}
	if req.Discount != nil {
		meisai.Discount = req.Discount
	}
	if req.CarIdNum != nil {
		meisai.CarIdNum = req.CarIdNum
	}
	if req.Detail != nil {
		meisai.Detail = req.Detail
	}
	if req.DtakoRowId != nil {
		meisai.DtakoRowId = req.DtakoRowId
	}

	result, err := s.repo.Create(ctx, orgID, meisai)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create etc_meisai: %v", err)
	}

	return &pb.CreateETCMeisaiResponse{
		EtcMeisai: toProtoETCMeisai(result),
	}, nil
}

// GetETCMeisai retrieves an ETC meisai by ID
func (s *ETCMeisaiServer) GetETCMeisai(ctx context.Context, req *pb.GetETCMeisaiRequest) (*pb.GetETCMeisaiResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	result, err := s.repo.GetByID(ctx, req.Id)
	if err != nil {
		if errors.Is(err, repository.ErrETCMeisaiNotFound) {
			return nil, status.Error(codes.NotFound, "etc_meisai not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get etc_meisai: %v", err)
	}

	return &pb.GetETCMeisaiResponse{
		EtcMeisai: toProtoETCMeisai(result),
	}, nil
}

// GetETCMeisaiByHash retrieves an ETC meisai by hash (for duplicate check)
func (s *ETCMeisaiServer) GetETCMeisaiByHash(ctx context.Context, req *pb.GetETCMeisaiByHashRequest) (*pb.GetETCMeisaiByHashResponse, error) {
	if req.Hash == "" {
		return nil, status.Error(codes.InvalidArgument, "hash is required")
	}

	result, err := s.repo.GetByHash(ctx, req.Hash)
	if err != nil {
		if errors.Is(err, repository.ErrETCMeisaiNotFound) {
			return &pb.GetETCMeisaiByHashResponse{
				Exists: false,
			}, nil
		}
		return nil, status.Errorf(codes.Internal, "failed to get etc_meisai by hash: %v", err)
	}

	return &pb.GetETCMeisaiByHashResponse{
		EtcMeisai: toProtoETCMeisai(result),
		Exists:    true,
	}, nil
}

// UpdateETCMeisai updates an ETC meisai record
func (s *ETCMeisaiServer) UpdateETCMeisai(ctx context.Context, req *pb.UpdateETCMeisaiRequest) (*pb.UpdateETCMeisaiResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	// Get existing record
	existing, err := s.repo.GetByID(ctx, req.Id)
	if err != nil {
		if errors.Is(err, repository.ErrETCMeisaiNotFound) {
			return nil, status.Error(codes.NotFound, "etc_meisai not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get etc_meisai: %v", err)
	}

	// Update fields
	if req.DateTo != nil {
		existing.DateTo = req.DateTo.AsTime()
	}
	if req.DateToDate != "" {
		existing.DateToDate = req.DateToDate
	}
	if req.IcFr != "" {
		existing.IcFr = req.IcFr
	}
	if req.IcTo != "" {
		existing.IcTo = req.IcTo
	}
	if req.DateFr != nil {
		t := req.DateFr.AsTime()
		existing.DateFr = &t
	}
	if req.PriceBf != nil {
		existing.PriceBf = req.PriceBf
	}
	if req.Discount != nil {
		existing.Discount = req.Discount
	}
	existing.Price = req.Price
	existing.Shashu = req.Shashu
	if req.CarIdNum != nil {
		existing.CarIdNum = req.CarIdNum
	}
	if req.EtcNum != "" {
		existing.EtcNum = req.EtcNum
	}
	if req.Detail != nil {
		existing.Detail = req.Detail
	}
	if req.DtakoRowId != nil {
		existing.DtakoRowId = req.DtakoRowId
	}
	if req.Hash != "" {
		existing.Hash = req.Hash
	}

	result, err := s.repo.Update(ctx, existing)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update etc_meisai: %v", err)
	}

	return &pb.UpdateETCMeisaiResponse{
		EtcMeisai: toProtoETCMeisai(result),
	}, nil
}

// DeleteETCMeisai deletes an ETC meisai record
func (s *ETCMeisaiServer) DeleteETCMeisai(ctx context.Context, req *pb.DeleteETCMeisaiRequest) (*pb.DeleteETCMeisaiResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	err := s.repo.Delete(ctx, req.Id)
	if err != nil {
		if errors.Is(err, repository.ErrETCMeisaiNotFound) {
			return nil, status.Error(codes.NotFound, "etc_meisai not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete etc_meisai: %v", err)
	}

	return &pb.DeleteETCMeisaiResponse{
		Success: true,
	}, nil
}

// ListETCMeisai lists ETC meisai with optional filters
func (s *ETCMeisaiServer) ListETCMeisai(ctx context.Context, req *pb.ListETCMeisaiRequest) (*pb.ListETCMeisaiResponse, error) {
	params := repository.ETCMeisaiListParams{
		PageSize:  int(req.PageSize),
		PageToken: req.PageToken,
	}

	if req.DateFrom != nil {
		params.DateFrom = req.DateFrom
	}
	if req.DateTo != nil {
		params.DateTo = req.DateTo
	}
	if req.EtcNum != nil {
		params.EtcNum = req.EtcNum
	}

	results, totalCount, nextPageToken, err := s.repo.List(ctx, params)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list etc_meisai: %v", err)
	}

	protoResults := make([]*pb.ETCMeisai, len(results))
	for i, m := range results {
		protoResults[i] = toProtoETCMeisai(m)
	}

	return &pb.ListETCMeisaiResponse{
		EtcMeisaiList: protoResults,
		NextPageToken: nextPageToken,
		TotalCount:    int32(totalCount),
	}, nil
}

// BulkCreateETCMeisai creates multiple ETC meisai records
func (s *ETCMeisaiServer) BulkCreateETCMeisai(ctx context.Context, req *pb.BulkCreateETCMeisaiRequest) (*pb.BulkCreateETCMeisaiResponse, error) {
	// Get organization_id from context (set by RLS interceptor)
	orgID, ok := db.GetOrganizationID(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required in header")
	}

	records := make([]*repository.ETCMeisai, len(req.Records))
	for i, r := range req.Records {
		meisai := &repository.ETCMeisai{
			DateTo:     r.DateTo.AsTime(),
			DateToDate: r.DateToDate,
			IcFr:       r.IcFr,
			IcTo:       r.IcTo,
			Price:      r.Price,
			Shashu:     r.Shashu,
			EtcNum:     r.EtcNum,
			Hash:       r.Hash,
		}

		if r.DateFr != nil {
			t := r.DateFr.AsTime()
			meisai.DateFr = &t
		}
		if r.PriceBf != nil {
			meisai.PriceBf = r.PriceBf
		}
		if r.Discount != nil {
			meisai.Discount = r.Discount
		}
		if r.CarIdNum != nil {
			meisai.CarIdNum = r.CarIdNum
		}
		if r.Detail != nil {
			meisai.Detail = r.Detail
		}
		if r.DtakoRowId != nil {
			meisai.DtakoRowId = r.DtakoRowId
		}

		records[i] = meisai
	}

	createdCount, skippedCount, errs, err := s.repo.BulkCreate(ctx, orgID, records, req.SkipDuplicates)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to bulk create etc_meisai: %v", err)
	}

	return &pb.BulkCreateETCMeisaiResponse{
		CreatedCount: int32(createdCount),
		SkippedCount: int32(skippedCount),
		Errors:       errs,
	}, nil
}

// toProtoETCMeisai converts repository model to proto message
func toProtoETCMeisai(m *repository.ETCMeisai) *pb.ETCMeisai {
	pbMeisai := &pb.ETCMeisai{
		Id:             m.ID,
		OrganizationId: m.OrganizationID,
		DateTo:         timestamppb.New(m.DateTo),
		DateToDate:     m.DateToDate,
		IcFr:           m.IcFr,
		IcTo:           m.IcTo,
		Price:          m.Price,
		Shashu:         m.Shashu,
		EtcNum:         m.EtcNum,
		Hash:           m.Hash,
		CreatedAt:      timestamppb.New(m.CreatedAt),
		UpdatedAt:      timestamppb.New(m.UpdatedAt),
	}

	if m.DateFr != nil {
		pbMeisai.DateFr = timestamppb.New(*m.DateFr)
	}
	if m.PriceBf != nil {
		pbMeisai.PriceBf = m.PriceBf
	}
	if m.Discount != nil {
		pbMeisai.Discount = m.Discount
	}
	if m.CarIdNum != nil {
		pbMeisai.CarIdNum = m.CarIdNum
	}
	if m.Detail != nil {
		pbMeisai.Detail = m.Detail
	}
	if m.DtakoRowId != nil {
		pbMeisai.DtakoRowId = m.DtakoRowId
	}

	return pbMeisai
}
