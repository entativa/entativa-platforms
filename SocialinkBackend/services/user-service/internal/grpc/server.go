package grpc

import (
	"context"
	"fmt"
	"log"
	"net"

	"socialink/user-service/internal/repository"
	"socialink/user-service/internal/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	grpcServer       *grpc.Server
	userRepo         *repository.UserRepository
	authService      *service.AuthService
	twoFactorService *service.TwoFactorService
	port             string
}

// NewGRPCServer creates a new gRPC server
func NewGRPCServer(
	port string,
	userRepo *repository.UserRepository,
	authService *service.AuthService,
	twoFactorService *service.TwoFactorService,
) *Server {
	grpcServer := grpc.NewServer()

	server := &Server{
		grpcServer:       grpcServer,
		userRepo:         userRepo,
		authService:      authService,
		twoFactorService: twoFactorService,
		port:             port,
	}

	// Register services
	// pb.RegisterUserServiceServer(grpcServer, server)

	// Enable reflection for debugging (remove in production)
	reflection.Register(grpcServer)

	return server
}

// Start starts the gRPC server
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %w", s.port, err)
	}

	log.Printf("🔌 gRPC server listening on port %s", s.port)

	if err := s.grpcServer.Serve(listener); err != nil {
		return fmt.Errorf("failed to serve gRPC: %w", err)
	}

	return nil
}

// Stop gracefully stops the gRPC server
func (s *Server) Stop() {
	log.Println("Stopping gRPC server...")
	s.grpcServer.GracefulStop()
}

// Example gRPC method implementation
func (s *Server) GetUser(ctx context.Context, req interface{}) (interface{}, error) {
	// Implementation here when proto is compiled
	return nil, nil
}
