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

	// Create repositories
	orgRepo := repository.NewOrganizationRepository(pool)
	appUserRepo := repository.NewAppUserRepository(pool)
	userOrgRepo := repository.NewUserOrganizationRepository(pool)
	fileRepo := repository.NewFileRepository(pool)
	flickrPhotoRepo := repository.NewFlickrPhotoRepository(pool)
	camFileRepo := repository.NewCamFileRepository(pool)
	camFileExeRepo := repository.NewCamFileExeRepository(pool)
	camFileExeStageRepo := repository.NewCamFileExeStageRepository(pool)
	ichibanCarRepo := repository.NewIchibanCarRepository(pool)
	dtakoCarsIchibanCarsRepo := repository.NewDtakoCarsIchibanCarsRepository(pool)
	uriageRepo := repository.NewUriageRepository(pool)
	uriageJishaRepo := repository.NewUriageJishaRepository(pool)
	carInspectionRepo := repository.NewCarInspectionRepository(pool)
	carInspectionFilesRepo := repository.NewCarInspectionFilesRepository(pool)
	carInspectionFilesARepo := repository.NewCarInspectionFilesARepository(pool)
	carInspectionFilesBRepo := repository.NewCarInspectionFilesBRepository(pool)
	carInspectionDeregistrationRepo := repository.NewCarInspectionDeregistrationRepository(pool)
	carInspectionDeregistrationFilesRepo := repository.NewCarInspectionDeregistrationFilesRepository(pool)
	carInsSheetIchibanCarsRepo := repository.NewCarInsSheetIchibanCarsRepository(pool)
	carInsSheetIchibanCarsARepo := repository.NewCarInsSheetIchibanCarsARepository(pool)
	kudgfryRepo := repository.NewKudgfryRepository(pool)
	kudguriRepo := repository.NewKudguriRepository(pool)
	kudgcstRepo := repository.NewKudgcstRepository(pool)
	kudgfulRepo := repository.NewKudgfulRepository(pool)
	kudgsirRepo := repository.NewKudgsirRepository(pool)
	kudgivtRepo := repository.NewKudgivtRepository(pool)
	dtakologsRepo := repository.NewDtakologsRepository(pool)

	// Create gRPC servers
	orgServer := grpcserver.NewOrganizationServer(orgRepo)
	appUserServer := grpcserver.NewAppUserServer(appUserRepo)
	userOrgServer := grpcserver.NewUserOrganizationServer(userOrgRepo)
	fileServer := grpcserver.NewFileServer(fileRepo)
	flickrPhotoServer := grpcserver.NewFlickrPhotoServer(flickrPhotoRepo)
	camFileServer := grpcserver.NewCamFileServer(camFileRepo)
	camFileExeServer := grpcserver.NewCamFileExeServer(camFileExeRepo)
	camFileExeStageServer := grpcserver.NewCamFileExeStageServer(camFileExeStageRepo)
	ichibanCarServer := grpcserver.NewIchibanCarServer(ichibanCarRepo)
	dtakoCarsIchibanCarsServer := grpcserver.NewDtakoCarsIchibanCarsServer(dtakoCarsIchibanCarsRepo)
	uriageServer := grpcserver.NewUriageServer(uriageRepo)
	uriageJishaServer := grpcserver.NewUriageJishaServer(uriageJishaRepo)
	carInspectionServer := grpcserver.NewCarInspectionServer(carInspectionRepo)
	carInspectionFilesServer := grpcserver.NewCarInspectionFilesServer(carInspectionFilesRepo)
	carInspectionFilesAServer := grpcserver.NewCarInspectionFilesAServer(carInspectionFilesARepo)
	carInspectionFilesBServer := grpcserver.NewCarInspectionFilesBServer(carInspectionFilesBRepo)
	carInspectionDeregistrationServer := grpcserver.NewCarInspectionDeregistrationServer(carInspectionDeregistrationRepo)
	carInspectionDeregistrationFilesServer := grpcserver.NewCarInspectionDeregistrationFilesServer(carInspectionDeregistrationFilesRepo)
	carInsSheetIchibanCarsServer := grpcserver.NewCarInsSheetIchibanCarsServer(carInsSheetIchibanCarsRepo)
	carInsSheetIchibanCarsAServer := grpcserver.NewCarInsSheetIchibanCarsAServer(carInsSheetIchibanCarsARepo)
	kudgfryServer := grpcserver.NewKudgfryServer(kudgfryRepo)
	kudguriServer := grpcserver.NewKudguriServer(kudguriRepo)
	kudgcstServer := grpcserver.NewKudgcstServer(kudgcstRepo)
	kudgfulServer := grpcserver.NewKudgfulServer(kudgfulRepo)
	kudgsirServer := grpcserver.NewKudgsirServer(kudgsirRepo)
	kudgivtServer := grpcserver.NewKudgivtServer(kudgivtRepo)
	dtakologsServer := grpcserver.NewDtakologsServer(dtakologsRepo)

	// Create gRPC server with health check
	grpcServer := grpc.NewServer()

	// Register all services
	pb.RegisterOrganizationServiceServer(grpcServer, orgServer)
	pb.RegisterAppUserServiceServer(grpcServer, appUserServer)
	pb.RegisterUserOrganizationServiceServer(grpcServer, userOrgServer)
	pb.RegisterFileServiceServer(grpcServer, fileServer)
	pb.RegisterFlickrPhotoServiceServer(grpcServer, flickrPhotoServer)
	pb.RegisterCamFileServiceServer(grpcServer, camFileServer)
	pb.RegisterCamFileExeServiceServer(grpcServer, camFileExeServer)
	pb.RegisterCamFileExeStageServiceServer(grpcServer, camFileExeStageServer)
	pb.RegisterIchibanCarServiceServer(grpcServer, ichibanCarServer)
	pb.RegisterDtakoCarsIchibanCarsServiceServer(grpcServer, dtakoCarsIchibanCarsServer)
	pb.RegisterUriageServiceServer(grpcServer, uriageServer)
	pb.RegisterUriageJishaServiceServer(grpcServer, uriageJishaServer)
	pb.RegisterCarInspectionServiceServer(grpcServer, carInspectionServer)
	pb.RegisterCarInspectionFilesServiceServer(grpcServer, carInspectionFilesServer)
	pb.RegisterCarInspectionFilesAServiceServer(grpcServer, carInspectionFilesAServer)
	pb.RegisterCarInspectionFilesBServiceServer(grpcServer, carInspectionFilesBServer)
	pb.RegisterCarInspectionDeregistrationServiceServer(grpcServer, carInspectionDeregistrationServer)
	pb.RegisterCarInspectionDeregistrationFilesServiceServer(grpcServer, carInspectionDeregistrationFilesServer)
	pb.RegisterCarInsSheetIchibanCarsServiceServer(grpcServer, carInsSheetIchibanCarsServer)
	pb.RegisterCarInsSheetIchibanCarsAServiceServer(grpcServer, carInsSheetIchibanCarsAServer)
	pb.RegisterKudgfryServiceServer(grpcServer, kudgfryServer)
	pb.RegisterKudguriServiceServer(grpcServer, kudguriServer)
	pb.RegisterKudgcstServiceServer(grpcServer, kudgcstServer)
	pb.RegisterKudgfulServiceServer(grpcServer, kudgfulServer)
	pb.RegisterKudgsirServiceServer(grpcServer, kudgsirServer)
	pb.RegisterKudgivtServiceServer(grpcServer, kudgivtServer)
	pb.RegisterDtakologsServiceServer(grpcServer, dtakologsServer)

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
