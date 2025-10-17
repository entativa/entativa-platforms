package grpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// MediaServiceClient for gRPC communication with media service
type MediaServiceClient struct {
	conn   *grpc.ClientConn
	// client pb.MediaServiceClient // In production: use actual protobuf client
}

func NewMediaServiceClient(address string) (*MediaServiceClient, error) {
	// Connect to media service
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to media service: %w", err)
	}

	log.Printf("✅ Connected to media service at %s", address)

	return &MediaServiceClient{
		conn: conn,
	}, nil
}

func (c *MediaServiceClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// SaveStreamRecording - Save stream recording to media service
func (c *MediaServiceClient) SaveStreamRecording(
	ctx context.Context,
	streamID, streamerID uuid.UUID,
	recordingPath string,
	duration int,
	quality string,
) (string, error) {
	// In production: Use gRPC to call media service
	// req := &pb.SaveStreamRecordingRequest{
	//     StreamId: streamID.String(),
	//     UserId: streamerID.String(),
	//     FilePath: recordingPath,
	//     Duration: int32(duration),
	//     Quality: quality,
	// }
	// resp, err := c.client.SaveStreamRecording(ctx, req)

	// Mock implementation
	log.Printf("Saving recording for stream %s to media service", streamID)
	time.Sleep(100 * time.Millisecond) // Simulate network call

	recordingURL := fmt.Sprintf("https://cdn.socialink.com/recordings/%s.mp4", streamID)
	return recordingURL, nil
}

// GenerateThumbnail - Generate thumbnail from stream
func (c *MediaServiceClient) GenerateThumbnail(
	ctx context.Context,
	streamID uuid.UUID,
	timestamp int, // Seconds into stream
) (string, error) {
	// In production: Call media service to generate thumbnail

	thumbnailURL := fmt.Sprintf("https://cdn.socialink.com/thumbnails/%s.jpg", streamID)
	return thumbnailURL, nil
}
