package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/pb"
	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/repository"
)

// DtakologsServer implements the gRPC DtakologsService
type DtakologsServer struct {
	pb.UnimplementedDtakologsServiceServer
	repo *repository.DtakologsRepository
}

// NewDtakologsServer creates a new gRPC server
func NewDtakologsServer(repo *repository.DtakologsRepository) *DtakologsServer {
	return &DtakologsServer{repo: repo}
}

// CreateDtakologs creates a new dtakologs record
func (s *DtakologsServer) CreateDtakologs(ctx context.Context, req *pb.CreateDtakologsRequest) (*pb.CreateDtakologsResponse, error) {
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}
	if req.Type == "" {
		return nil, status.Error(codes.InvalidArgument, "type is required")
	}
	if req.DataDateTime == "" {
		return nil, status.Error(codes.InvalidArgument, "data_date_time is required")
	}

	// Convert proto request to repository model
	d := &repository.Dtakologs{
		OrganizationID:               req.OrganizationId,
		Type:                         req.Type,
		AddressDispC:                 req.AddressDispC,
		AddressDispP:                 req.AddressDispP,
		AllState:                     req.AllState,
		AllStateEx:                   req.AllStateEx,
		AllStateFontColor:            req.AllStateFontColor,
		AllStateFontColorIndex:       req.AllStateFontColorIndex,
		AllStateRyoutColor:           req.AllStateRyoutColor,
		BranchCd:                     req.BranchCd,
		BranchName:                   req.BranchName,
		ComuDateTime:                 req.ComuDateTime,
		CurrentWorkCd:                req.CurrentWorkCd,
		CurrentWorkName:              req.CurrentWorkName,
		DataDateTime:                 req.DataDateTime,
		DataFilterType:               req.DataFilterType,
		DispFlag:                     req.DispFlag,
		DriverCd:                     req.DriverCd,
		DriverName:                   req.DriverName,
		EventVal:                     req.EventVal,
		GpsDirection:                 req.GpsDirection,
		GpsEnable:                    req.GpsEnable,
		GpsLatiAndLong:               req.GpsLatiAndLong,
		GpsLatitude:                  req.GpsLatitude,
		GpsLongitude:                 req.GpsLongitude,
		GpsSatelliteNum:              req.GpsSatelliteNum,
		OdoMeter:                     req.OdoMeter,
		OperationState:               req.OperationState,
		ReciveEventType:              req.ReciveEventType,
		RecivePacketType:             req.RecivePacketType,
		ReciveTypeColorName:          req.ReciveTypeColorName,
		ReciveTypeName:               req.ReciveTypeName,
		ReciveWorkCd:                 req.ReciveWorkCd,
		Revo:                         req.Revo,
		SettingTemp:                  req.SettingTemp,
		SettingTemp1:                 req.SettingTemp1,
		SettingTemp3:                 req.SettingTemp3,
		SettingTemp4:                 req.SettingTemp4,
		Speed:                        req.Speed,
		StartWorkDateTime:            req.StartWorkDateTime,
		State:                        req.State,
		State1:                       req.State1,
		State2:                       req.State2,
		State3:                       req.State3,
		StateFlag:                    req.StateFlag,
		SubDriverCd:                  req.SubDriverCd,
		Temp1:                        req.Temp1,
		Temp2:                        req.Temp2,
		Temp3:                        req.Temp3,
		Temp4:                        req.Temp4,
		TempState:                    req.TempState,
		VehicleCd:                    req.VehicleCd,
		VehicleIconColor:             req.VehicleIconColor,
		VehicleIconLabelForDatetime:  req.VehicleIconLabelForDatetime,
		VehicleIconLabelForDriver:    req.VehicleIconLabelForDriver,
		VehicleIconLabelForVehicle:   req.VehicleIconLabelForVehicle,
		VehicleName:                  req.VehicleName,
	}

	err := s.repo.Create(ctx, d)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create dtakologs: %v", err)
	}

	return &pb.CreateDtakologsResponse{
		Dtakologs: toProtoDtakologs(d),
	}, nil
}

// GetDtakologs retrieves a dtakologs record by composite primary key
// NOTE: The proto only defines organization_id and type, but the repository requires
// organization_id, data_date_time, and vehicle_cd. This is a proto/repository mismatch.
// This implementation uses the type field as a placeholder and will need correction.
func (s *DtakologsServer) GetDtakologs(ctx context.Context, req *pb.GetDtakologsRequest) (*pb.GetDtakologsResponse, error) {
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}
	if req.Type == "" {
		return nil, status.Error(codes.InvalidArgument, "type is required")
	}

	// TODO: This is a workaround - the proto needs to be updated to include data_date_time and vehicle_cd
	// For now, this will not work correctly as the repository needs the full composite key
	// Placeholder implementation that won't work without proto changes
	return nil, status.Error(codes.Unimplemented, "GetDtakologs requires proto update to include data_date_time and vehicle_cd")
}

// UpdateDtakologs updates an existing dtakologs record
func (s *DtakologsServer) UpdateDtakologs(ctx context.Context, req *pb.UpdateDtakologsRequest) (*pb.UpdateDtakologsResponse, error) {
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}
	if req.Type == "" {
		return nil, status.Error(codes.InvalidArgument, "type is required")
	}
	if req.DataDateTime == "" {
		return nil, status.Error(codes.InvalidArgument, "data_date_time is required")
	}

	// Convert proto request to repository model
	d := &repository.Dtakologs{
		OrganizationID:               req.OrganizationId,
		Type:                         req.Type,
		AddressDispC:                 req.AddressDispC,
		AddressDispP:                 req.AddressDispP,
		AllState:                     req.AllState,
		AllStateEx:                   req.AllStateEx,
		AllStateFontColor:            req.AllStateFontColor,
		AllStateFontColorIndex:       req.AllStateFontColorIndex,
		AllStateRyoutColor:           req.AllStateRyoutColor,
		BranchCd:                     req.BranchCd,
		BranchName:                   req.BranchName,
		ComuDateTime:                 req.ComuDateTime,
		CurrentWorkCd:                req.CurrentWorkCd,
		CurrentWorkName:              req.CurrentWorkName,
		DataDateTime:                 req.DataDateTime,
		DataFilterType:               req.DataFilterType,
		DispFlag:                     req.DispFlag,
		DriverCd:                     req.DriverCd,
		DriverName:                   req.DriverName,
		EventVal:                     req.EventVal,
		GpsDirection:                 req.GpsDirection,
		GpsEnable:                    req.GpsEnable,
		GpsLatiAndLong:               req.GpsLatiAndLong,
		GpsLatitude:                  req.GpsLatitude,
		GpsLongitude:                 req.GpsLongitude,
		GpsSatelliteNum:              req.GpsSatelliteNum,
		OdoMeter:                     req.OdoMeter,
		OperationState:               req.OperationState,
		ReciveEventType:              req.ReciveEventType,
		RecivePacketType:             req.RecivePacketType,
		ReciveTypeColorName:          req.ReciveTypeColorName,
		ReciveTypeName:               req.ReciveTypeName,
		ReciveWorkCd:                 req.ReciveWorkCd,
		Revo:                         req.Revo,
		SettingTemp:                  req.SettingTemp,
		SettingTemp1:                 req.SettingTemp1,
		SettingTemp3:                 req.SettingTemp3,
		SettingTemp4:                 req.SettingTemp4,
		Speed:                        req.Speed,
		StartWorkDateTime:            req.StartWorkDateTime,
		State:                        req.State,
		State1:                       req.State1,
		State2:                       req.State2,
		State3:                       req.State3,
		StateFlag:                    req.StateFlag,
		SubDriverCd:                  req.SubDriverCd,
		Temp1:                        req.Temp1,
		Temp2:                        req.Temp2,
		Temp3:                        req.Temp3,
		Temp4:                        req.Temp4,
		TempState:                    req.TempState,
		VehicleCd:                    req.VehicleCd,
		VehicleIconColor:             req.VehicleIconColor,
		VehicleIconLabelForDatetime:  req.VehicleIconLabelForDatetime,
		VehicleIconLabelForDriver:    req.VehicleIconLabelForDriver,
		VehicleIconLabelForVehicle:   req.VehicleIconLabelForVehicle,
		VehicleName:                  req.VehicleName,
	}

	err := s.repo.Update(ctx, d)
	if err != nil {
		if errors.Is(err, repository.ErrDtakologsNotFound) {
			return nil, status.Error(codes.NotFound, "dtakologs not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update dtakologs: %v", err)
	}

	return &pb.UpdateDtakologsResponse{
		Dtakologs: toProtoDtakologs(d),
	}, nil
}

// DeleteDtakologs removes a dtakologs record
// NOTE: The proto only defines organization_id and type, but the repository requires
// organization_id, data_date_time, and vehicle_cd. This is a proto/repository mismatch.
func (s *DtakologsServer) DeleteDtakologs(ctx context.Context, req *pb.DeleteDtakologsRequest) (*pb.DeleteDtakologsResponse, error) {
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}
	if req.Type == "" {
		return nil, status.Error(codes.InvalidArgument, "type is required")
	}

	// TODO: This is a workaround - the proto needs to be updated to include data_date_time and vehicle_cd
	// For now, this will not work correctly as the repository needs the full composite key
	return nil, status.Error(codes.Unimplemented, "DeleteDtakologs requires proto update to include data_date_time and vehicle_cd")
}

// ListDtakologs retrieves all dtakologs records with pagination
func (s *DtakologsServer) ListDtakologs(ctx context.Context, req *pb.ListDtakologsRequest) (*pb.ListDtakologsResponse, error) {
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
		return nil, status.Errorf(codes.Internal, "failed to list dtakologs: %v", err)
	}

	var nextPageToken string
	if len(records) > limit {
		records = records[:limit]
		nextPageToken = "next" // In production, encode the offset
	}

	protoDtakologs := make([]*pb.Dtakologs, len(records))
	for i, d := range records {
		protoDtakologs[i] = toProtoDtakologs(d)
	}

	return &pb.ListDtakologsResponse{
		Dtakologs:     protoDtakologs,
		NextPageToken: nextPageToken,
	}, nil
}

// ListDtakologsByOrganization retrieves dtakologs records for a specific organization with pagination
func (s *DtakologsServer) ListDtakologsByOrganization(ctx context.Context, req *pb.ListDtakologsByOrganizationRequest) (*pb.ListDtakologsByOrganizationResponse, error) {
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
		return nil, status.Errorf(codes.Internal, "failed to list dtakologs by organization: %v", err)
	}

	var nextPageToken string
	if len(records) > limit {
		records = records[:limit]
		nextPageToken = "next" // In production, encode the offset
	}

	protoDtakologs := make([]*pb.Dtakologs, len(records))
	for i, d := range records {
		protoDtakologs[i] = toProtoDtakologs(d)
	}

	return &pb.ListDtakologsByOrganizationResponse{
		Dtakologs:     protoDtakologs,
		NextPageToken: nextPageToken,
	}, nil
}

// toProtoDtakologs converts repository model to proto message
func toProtoDtakologs(d *repository.Dtakologs) *pb.Dtakologs {
	proto := &pb.Dtakologs{
		OrganizationId:              d.OrganizationID,
		Type:                        d.Type,
		AllStateFontColorIndex:      d.AllStateFontColorIndex,
		AllStateRyoutColor:          d.AllStateRyoutColor,
		BranchCd:                    d.BranchCd,
		BranchName:                  d.BranchName,
		CurrentWorkCd:               d.CurrentWorkCd,
		DataDateTime:                d.DataDateTime,
		DataFilterType:              d.DataFilterType,
		DispFlag:                    d.DispFlag,
		DriverCd:                    d.DriverCd,
		GpsDirection:                d.GpsDirection,
		GpsEnable:                   d.GpsEnable,
		GpsLatitude:                 d.GpsLatitude,
		GpsLongitude:                d.GpsLongitude,
		GpsSatelliteNum:             d.GpsSatelliteNum,
		OperationState:              d.OperationState,
		ReciveEventType:             d.ReciveEventType,
		RecivePacketType:            d.RecivePacketType,
		ReciveWorkCd:                d.ReciveWorkCd,
		Revo:                        d.Revo,
		SettingTemp:                 d.SettingTemp,
		SettingTemp1:                d.SettingTemp1,
		SettingTemp3:                d.SettingTemp3,
		SettingTemp4:                d.SettingTemp4,
		Speed:                       d.Speed,
		StateFlag:                   d.StateFlag,
		SubDriverCd:                 d.SubDriverCd,
		TempState:                   d.TempState,
		VehicleCd:                   d.VehicleCd,
		VehicleName:                 d.VehicleName,
	}

	// Handle optional fields
	if d.AddressDispC != nil {
		proto.AddressDispC = d.AddressDispC
	}
	if d.AddressDispP != nil {
		proto.AddressDispP = d.AddressDispP
	}
	if d.AllState != nil {
		proto.AllState = d.AllState
	}
	if d.AllStateEx != nil {
		proto.AllStateEx = d.AllStateEx
	}
	if d.AllStateFontColor != nil {
		proto.AllStateFontColor = d.AllStateFontColor
	}
	if d.ComuDateTime != nil {
		proto.ComuDateTime = d.ComuDateTime
	}
	if d.CurrentWorkName != nil {
		proto.CurrentWorkName = d.CurrentWorkName
	}
	if d.DriverName != nil {
		proto.DriverName = d.DriverName
	}
	if d.EventVal != nil {
		proto.EventVal = d.EventVal
	}
	if d.GpsLatiAndLong != nil {
		proto.GpsLatiAndLong = d.GpsLatiAndLong
	}
	if d.OdoMeter != nil {
		proto.OdoMeter = d.OdoMeter
	}
	if d.ReciveTypeColorName != nil {
		proto.ReciveTypeColorName = d.ReciveTypeColorName
	}
	if d.ReciveTypeName != nil {
		proto.ReciveTypeName = d.ReciveTypeName
	}
	if d.StartWorkDateTime != nil {
		proto.StartWorkDateTime = d.StartWorkDateTime
	}
	if d.State != nil {
		proto.State = d.State
	}
	if d.State1 != nil {
		proto.State1 = d.State1
	}
	if d.State2 != nil {
		proto.State2 = d.State2
	}
	if d.State3 != nil {
		proto.State3 = d.State3
	}
	if d.Temp1 != nil {
		proto.Temp1 = d.Temp1
	}
	if d.Temp2 != nil {
		proto.Temp2 = d.Temp2
	}
	if d.Temp3 != nil {
		proto.Temp3 = d.Temp3
	}
	if d.Temp4 != nil {
		proto.Temp4 = d.Temp4
	}
	if d.VehicleIconColor != nil {
		proto.VehicleIconColor = d.VehicleIconColor
	}
	if d.VehicleIconLabelForDatetime != nil {
		proto.VehicleIconLabelForDatetime = d.VehicleIconLabelForDatetime
	}
	if d.VehicleIconLabelForDriver != nil {
		proto.VehicleIconLabelForDriver = d.VehicleIconLabelForDriver
	}
	if d.VehicleIconLabelForVehicle != nil {
		proto.VehicleIconLabelForVehicle = d.VehicleIconLabelForVehicle
	}

	return proto
}
