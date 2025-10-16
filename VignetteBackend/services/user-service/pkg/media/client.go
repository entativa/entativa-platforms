package media

import (
	"context"
	"fmt"
	"io"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "socialink/user-service/proto/media"
)

// Client is a gRPC client for the media service
type Client struct {
	conn   *grpc.ClientConn
	client pb.MediaServiceClient
}

// NewClient creates a new media service client
func NewClient(address string) (*Client, error) {
	conn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(100*1024*1024), // 100MB
			grpc.MaxCallSendMsgSize(100*1024*1024), // 100MB
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to media service: %w", err)
	}

	return &Client{
		conn:   conn,
		client: pb.NewMediaServiceClient(conn),
	}, nil
}

// Close closes the gRPC connection
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// UploadMedia uploads a media file
func (c *Client) UploadMedia(ctx context.Context, req *pb.UploadMediaRequest) (*pb.UploadMediaResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	return c.client.UploadMedia(ctx, req)
}

// UploadMediaStream uploads a large media file in chunks
func (c *Client) UploadMediaStream(ctx context.Context, chunks []*pb.UploadChunk) (*pb.UploadMediaResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	stream, err := c.client.UploadMediaStream(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to open stream: %w", err)
	}

	for _, chunk := range chunks {
		if err := stream.Send(chunk); err != nil {
			return nil, fmt.Errorf("failed to send chunk: %w", err)
		}
	}

	return stream.CloseAndRecv()
}

// GetMedia retrieves media information
func (c *Client) GetMedia(ctx context.Context, mediaID string) (*pb.MediaResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return c.client.GetMedia(ctx, &pb.GetMediaRequest{
		MediaId: mediaID,
	})
}

// DeleteMedia deletes a media file
func (c *Client) DeleteMedia(ctx context.Context, mediaID, userID string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := c.client.DeleteMedia(ctx, &pb.DeleteMediaRequest{
		MediaId: mediaID,
		UserId:  userID,
	})
	return err
}

// GetSignedUrl generates a signed URL for media download
func (c *Client) GetSignedUrl(ctx context.Context, mediaID string, expirySeconds int32) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	resp, err := c.client.GetSignedUrl(ctx, &pb.GetSignedUrlRequest{
		MediaId:       mediaID,
		ExpirySeconds: expirySeconds,
	})
	if err != nil {
		return "", err
	}

	return resp.SignedUrl, nil
}

// BatchGetMedia retrieves multiple media files
func (c *Client) BatchGetMedia(ctx context.Context, mediaIDs []string) ([]*pb.MediaResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	resp, err := c.client.BatchGetMedia(ctx, &pb.BatchGetMediaRequest{
		MediaIds: mediaIDs,
	})
	if err != nil {
		return nil, err
	}

	return resp.Media, nil
}

// BatchDeleteMedia deletes multiple media files
func (c *Client) BatchDeleteMedia(ctx context.Context, mediaIDs []string, userID string) (int32, []string, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	resp, err := c.client.BatchDeleteMedia(ctx, &pb.BatchDeleteMediaRequest{
		MediaIds: mediaIDs,
		UserId:   userID,
	})
	if err != nil {
		return 0, nil, err
	}

	return resp.DeletedCount, resp.FailedIds, nil
}

// Helper function to create upload request
func CreateUploadRequest(data []byte, filename, contentType, userID string, purpose pb.MediaPurpose) *pb.UploadMediaRequest {
	return &pb.UploadMediaRequest{
		Data:        data,
		Filename:    filename,
		ContentType: contentType,
		UserId:      userID,
		Purpose:     purpose,
		ProcessingOptions: &pb.ProcessingOptions{
			GenerateThumbnails: true,
			GenerateBlurhash:   true,
			Quality:            92,
		},
	}
}

// Helper function to create chunks for streaming upload
func CreateChunks(data []byte, filename, contentType, userID string, purpose pb.MediaPurpose, chunkSize int) []*pb.UploadChunk {
	if chunkSize <= 0 {
		chunkSize = 1024 * 1024 // 1MB default
	}

	uploadID := fmt.Sprintf("upload_%d", time.Now().UnixNano())
	totalChunks := (len(data) + chunkSize - 1) / chunkSize
	chunks := make([]*pb.UploadChunk, 0, totalChunks)

	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}

		chunk := &pb.UploadChunk{
			UploadId:    uploadID,
			ChunkNumber: int32(i / chunkSize),
			Data:        data[i:end],
			IsLastChunk: end == len(data),
			Filename:    filename,
			ContentType: contentType,
			UserId:      userID,
			Purpose:     purpose,
		}

		chunks = append(chunks, chunk)
	}

	return chunks
}

// UploadProfilePicture is a convenience method for profile picture uploads
func (c *Client) UploadProfilePicture(ctx context.Context, data []byte, filename, userID string) (*pb.UploadMediaResponse, error) {
	req := CreateUploadRequest(data, filename, "image/jpeg", userID, pb.MediaPurpose_PROFILE_PICTURE)
	
	// For profile pictures, ensure thumbnails are generated
	if req.ProcessingOptions == nil {
		req.ProcessingOptions = &pb.ProcessingOptions{}
	}
	req.ProcessingOptions.GenerateThumbnails = true
	req.ProcessingOptions.ThumbnailSizes = []*pb.ThumbnailSize{
		{Name: "thumb", Width: 150, Height: 150},
		{Name: "medium", Width: 300, Height: 300},
		{Name: "large", Width: 600, Height: 600},
	}
	req.ProcessingOptions.MaxWidth = 2048
	req.ProcessingOptions.MaxHeight = 2048

	return c.UploadMedia(ctx, req)
}

// UploadCoverPhoto is a convenience method for cover photo uploads
func (c *Client) UploadCoverPhoto(ctx context.Context, data []byte, filename, userID string) (*pb.UploadMediaResponse, error) {
	req := CreateUploadRequest(data, filename, "image/jpeg", userID, pb.MediaPurpose_COVER_PHOTO)
	
	// For cover photos, use wider thumbnails
	if req.ProcessingOptions == nil {
		req.ProcessingOptions = &pb.ProcessingOptions{}
	}
	req.ProcessingOptions.GenerateThumbnails = true
	req.ProcessingOptions.ThumbnailSizes = []*pb.ThumbnailSize{
		{Name: "thumb", Width: 400, Height: 150},
		{Name: "medium", Width: 800, Height: 300},
		{Name: "large", Width: 1600, Height: 600},
	}
	req.ProcessingOptions.MaxWidth = 1920
	req.ProcessingOptions.MaxHeight = 1080

	return c.UploadMedia(ctx, req)
}
