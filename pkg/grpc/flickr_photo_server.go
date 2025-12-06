package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/pb"
	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/repository"
)

// FlickrPhotoServer implements the gRPC FlickrPhotoService
type FlickrPhotoServer struct {
	pb.UnimplementedFlickrPhotoServiceServer
	repo *repository.FlickrPhotoRepository
}

// NewFlickrPhotoServer creates a new gRPC server
func NewFlickrPhotoServer(repo *repository.FlickrPhotoRepository) *FlickrPhotoServer {
	return &FlickrPhotoServer{repo: repo}
}

// CreateFlickrPhoto creates a new flickr photo
func (s *FlickrPhotoServer) CreateFlickrPhoto(ctx context.Context, req *pb.CreateFlickrPhotoRequest) (*pb.CreateFlickrPhotoResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	if req.OrganizationId == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}
	if req.Secret == "" {
		return nil, status.Error(codes.InvalidArgument, "secret is required")
	}
	if req.Server == "" {
		return nil, status.Error(codes.InvalidArgument, "server is required")
	}

	photo, err := s.repo.Create(ctx, req.Id, req.OrganizationId, req.Secret, req.Server)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create flickr photo: %v", err)
	}

	return &pb.CreateFlickrPhotoResponse{
		FlickrPhoto: toProtoFlickrPhoto(photo),
	}, nil
}

// GetFlickrPhoto retrieves a flickr photo by ID
func (s *FlickrPhotoServer) GetFlickrPhoto(ctx context.Context, req *pb.GetFlickrPhotoRequest) (*pb.GetFlickrPhotoResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	photo, err := s.repo.GetByID(ctx, req.Id)
	if err != nil {
		if errors.Is(err, repository.ErrFlickrPhotoNotFound) {
			return nil, status.Error(codes.NotFound, "flickr photo not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get flickr photo: %v", err)
	}

	return &pb.GetFlickrPhotoResponse{
		FlickrPhoto: toProtoFlickrPhoto(photo),
	}, nil
}

// UpdateFlickrPhoto updates an existing flickr photo
func (s *FlickrPhotoServer) UpdateFlickrPhoto(ctx context.Context, req *pb.UpdateFlickrPhotoRequest) (*pb.UpdateFlickrPhotoResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	if req.Secret == "" {
		return nil, status.Error(codes.InvalidArgument, "secret is required")
	}
	if req.Server == "" {
		return nil, status.Error(codes.InvalidArgument, "server is required")
	}

	photo, err := s.repo.Update(ctx, req.Id, req.Secret, req.Server)
	if err != nil {
		if errors.Is(err, repository.ErrFlickrPhotoNotFound) {
			return nil, status.Error(codes.NotFound, "flickr photo not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update flickr photo: %v", err)
	}

	return &pb.UpdateFlickrPhotoResponse{
		FlickrPhoto: toProtoFlickrPhoto(photo),
	}, nil
}

// DeleteFlickrPhoto hard-deletes a flickr photo
func (s *FlickrPhotoServer) DeleteFlickrPhoto(ctx context.Context, req *pb.DeleteFlickrPhotoRequest) (*pb.DeleteFlickrPhotoResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	err := s.repo.Delete(ctx, req.Id)
	if err != nil {
		if errors.Is(err, repository.ErrFlickrPhotoNotFound) {
			return nil, status.Error(codes.NotFound, "flickr photo not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete flickr photo: %v", err)
	}

	return &pb.DeleteFlickrPhotoResponse{
		Success: true,
	}, nil
}

// ListFlickrPhotos retrieves flickr photos with pagination
func (s *FlickrPhotoServer) ListFlickrPhotos(ctx context.Context, req *pb.ListFlickrPhotosRequest) (*pb.ListFlickrPhotosResponse, error) {
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

	photos, err := s.repo.List(ctx, limit+1, offset) // +1 to check if there's next page
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list flickr photos: %v", err)
	}

	var nextPageToken string
	if len(photos) > limit {
		photos = photos[:limit]
		nextPageToken = "next" // In production, encode the offset
	}

	protoPhotos := make([]*pb.FlickrPhoto, len(photos))
	for i, photo := range photos {
		protoPhotos[i] = toProtoFlickrPhoto(photo)
	}

	return &pb.ListFlickrPhotosResponse{
		FlickrPhotos:  protoPhotos,
		NextPageToken: nextPageToken,
	}, nil
}

// ListFlickrPhotosByOrganization retrieves flickr photos for a specific organization with pagination
func (s *FlickrPhotoServer) ListFlickrPhotosByOrganization(ctx context.Context, req *pb.ListFlickrPhotosByOrganizationRequest) (*pb.ListFlickrPhotosByOrganizationResponse, error) {
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

	photos, err := s.repo.ListByOrganization(ctx, req.OrganizationId, limit+1, offset) // +1 to check if there's next page
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list flickr photos by organization: %v", err)
	}

	var nextPageToken string
	if len(photos) > limit {
		photos = photos[:limit]
		nextPageToken = "next" // In production, encode the offset
	}

	protoPhotos := make([]*pb.FlickrPhoto, len(photos))
	for i, photo := range photos {
		protoPhotos[i] = toProtoFlickrPhoto(photo)
	}

	return &pb.ListFlickrPhotosByOrganizationResponse{
		FlickrPhotos:  protoPhotos,
		NextPageToken: nextPageToken,
	}, nil
}

// toProtoFlickrPhoto converts repository model to proto message
func toProtoFlickrPhoto(photo *repository.FlickrPhoto) *pb.FlickrPhoto {
	return &pb.FlickrPhoto{
		Id:             photo.ID,
		OrganizationId: photo.OrganizationID,
		Secret:         photo.Secret,
		Server:         photo.Server,
	}
}
