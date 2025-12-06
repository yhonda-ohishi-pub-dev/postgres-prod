package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/yhonda-ohishi-pub-dev/postgres-prod/internal/config"
	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/db"
	grpcserver "github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/grpc"
	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/pb"
	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/repository"
)

func main() {
	ctx := context.Background()

	// Load configuration
	cfg := config.Load()

	// Validate required configuration
	if cfg.InstanceConnection == "" {
		log.Fatal("Cloud SQL instance connection not configured. Set GCP_PROJECT_ID, GCP_REGION, and CLOUDSQL_INSTANCE_NAME environment variables.")
	}
	if cfg.DatabaseUser == "" {
		log.Fatal("Database user not configured. Set DB_USER environment variable (IAM email without domain for Cloud SQL IAM auth).")
	}

	// Create database connection pool
	pool, cleanup, err := db.NewPool(ctx, cfg.InstanceConnection, cfg.DatabaseUser, cfg.DatabaseName, cfg.DatabasePassword, cfg.DatabasePort)
	if err != nil {
		log.Fatalf("Failed to connect to Cloud SQL: %v", err)
	}
	defer cleanup()

	log.Printf("Connected to database: %s", cfg.DatabaseName)

	// Create repository and gRPC server
	orgRepo := repository.NewOrganizationRepository(pool)
	orgServer := grpcserver.NewOrganizationServer(orgRepo)

	// Create gRPC server with health check
	grpcServer := grpc.NewServer()
	pb.RegisterOrganizationServiceServer(grpcServer, orgServer)

	// Register health check service for Cloud Run
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	reflection.Register(grpcServer)

	// Use PORT env (Cloud Run sets this to 8080)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh

		log.Println("Shutting down gRPC server...")
		grpcServer.GracefulStop()
	}()

	log.Printf("Starting gRPC server on port %s", port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("gRPC server error: %v", err)
	}

	log.Println("Server stopped")
}
