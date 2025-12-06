package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/yhonda-ohishi-pub-dev/postgres-prod/internal/config"
	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/db"
	grpcserver "github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/grpc"
	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/handlers"
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
	// Uses IAM auth on Cloud Run, proxy with --auto-iam-authn locally
	pool, cleanup, err := db.NewPool(ctx, cfg.InstanceConnection, cfg.DatabaseUser, cfg.DatabaseName, cfg.DatabasePassword, cfg.DatabasePort)
	if err != nil {
		log.Fatalf("Failed to connect to Cloud SQL: %v", err)
	}
	defer cleanup()

	log.Printf("Connected to database: %s", cfg.DatabaseName)

	// Create repository and gRPC server
	orgRepo := repository.NewOrganizationRepository(pool)
	orgServer := grpcserver.NewOrganizationServer(orgRepo)

	// Create gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterOrganizationServiceServer(grpcServer, orgServer)
	reflection.Register(grpcServer)

	// Start gRPC server in background
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}
	grpcListener, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}
	go func() {
		log.Printf("Starting gRPC server on port %s", grpcPort)
		if err := grpcServer.Serve(grpcListener); err != nil {
			log.Printf("gRPC server error: %v", err)
		}
	}()

	// Setup HTTP routes
	mux := http.NewServeMux()
	mux.Handle("/health", handlers.NewHealthHandler(pool))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("postgres-prod service running"))
	})

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh

		log.Println("Shutting down servers...")

		// Gracefully stop gRPC server
		grpcServer.GracefulStop()
		log.Println("gRPC server stopped")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}
	}()

	log.Printf("Starting server on port %s", cfg.Port)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}

	log.Println("Server stopped")
}
