package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/yhonda-ohishi-pub-dev/postgres-prod/internal/config"
	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/auth"
	"github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/db"
	grpcserver "github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/grpc"
	httphandler "github.com/yhonda-ohishi-pub-dev/postgres-prod/pkg/http"
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

	// Create RLS-aware pool wrapper
	rlsPool := db.NewRLSPool(pool)

	// Create repositories with RLS pool (auto-sets app.organization_id per request)
	orgRepo := repository.NewOrganizationRepositoryWithDB(rlsPool)
	appUserRepo := repository.NewAppUserRepositoryWithDB(rlsPool)
	userOrgRepo := repository.NewUserOrganizationRepositoryWithDB(rlsPool)
	fileRepo := repository.NewFileRepositoryWithDB(rlsPool)
	flickrPhotoRepo := repository.NewFlickrPhotoRepositoryWithDB(rlsPool)
	camFileRepo := repository.NewCamFileRepositoryWithDB(rlsPool)
	camFileExeRepo := repository.NewCamFileExeRepositoryWithDB(rlsPool)
	camFileExeStageRepo := repository.NewCamFileExeStageRepositoryWithDB(rlsPool)
	ichibanCarRepo := repository.NewIchibanCarRepositoryWithDB(rlsPool)
	dtakoCarsIchibanCarsRepo := repository.NewDtakoCarsIchibanCarsRepositoryWithDB(rlsPool)
	uriageRepo := repository.NewUriageRepositoryWithDB(rlsPool)
	uriageJishaRepo := repository.NewUriageJishaRepositoryWithDB(rlsPool)
	carInspectionRepo := repository.NewCarInspectionRepositoryWithDB(rlsPool)
	carInspectionFilesRepo := repository.NewCarInspectionFilesRepositoryWithDB(rlsPool)
	carInspectionFilesARepo := repository.NewCarInspectionFilesARepositoryWithDB(rlsPool)
	carInspectionFilesBRepo := repository.NewCarInspectionFilesBRepositoryWithDB(rlsPool)
	carInspectionDeregistrationRepo := repository.NewCarInspectionDeregistrationRepositoryWithDB(rlsPool)
	carInspectionDeregistrationFilesRepo := repository.NewCarInspectionDeregistrationFilesRepositoryWithDB(rlsPool)
	carInsSheetIchibanCarsRepo := repository.NewCarInsSheetIchibanCarsRepositoryWithDB(rlsPool)
	carInsSheetIchibanCarsARepo := repository.NewCarInsSheetIchibanCarsARepositoryWithDB(rlsPool)
	kudgfryRepo := repository.NewKudgfryRepositoryWithDB(rlsPool)
	kudguriRepo := repository.NewKudguriRepositoryWithDB(rlsPool)
	kudgcstRepo := repository.NewKudgcstRepositoryWithDB(rlsPool)
	kudgfulRepo := repository.NewKudgfulRepositoryWithDB(rlsPool)
	kudgsirRepo := repository.NewKudgsirRepositoryWithDB(rlsPool)
	kudgivtRepo := repository.NewKudgivtRepositoryWithDB(rlsPool)
	dtakologsRepo := repository.NewDtakologsRepositoryWithDB(rlsPool)
	oauthAccountRepo := repository.NewOAuthAccountRepositoryWithDB(rlsPool)
	invitationRepo := repository.NewInvitationRepositoryWithDB(rlsPool)
	etcMeisaiRepo := repository.NewETCMeisaiRepositoryWithDB(rlsPool)

	// Create auth services
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-secret-change-in-production"
		log.Println("WARNING: JWT_SECRET not set, using default secret")
	}
	jwtService := auth.NewJWTService(jwtSecret, 15*time.Minute, 7*24*time.Hour)

	googleClient := auth.NewGoogleOAuthClient(
		os.Getenv("GOOGLE_CLIENT_ID"),
		os.Getenv("GOOGLE_CLIENT_SECRET"),
		os.Getenv("GOOGLE_REDIRECT_URI"),
	)
	lineClient := auth.NewLineOAuthClient(
		os.Getenv("LINE_CHANNEL_ID"),
		os.Getenv("LINE_CHANNEL_SECRET"),
		os.Getenv("LINE_REDIRECT_URI"),
	)

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
	authServer := grpcserver.NewAuthServer(appUserRepo, oauthAccountRepo, jwtService, googleClient, lineClient)
	invitationServer := grpcserver.NewInvitationServer(invitationRepo, orgRepo, userOrgRepo)
	etcMeisaiServer := grpcserver.NewETCMeisaiServer(etcMeisaiRepo)

	// Create HTTP auth handler
	authHandler := httphandler.NewAuthHandler(googleClient, lineClient, jwtService, appUserRepo, oauthAccountRepo, cfg.FrontendURL)

	// Create gRPC server with health check, JWT auth, and RLS interceptor
	// JWT interceptor runs first (validates token and sets user context),
	// then RLS interceptor (sets organization context)
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcserver.JWTUnaryInterceptor(jwtService),
			grpcserver.RLSUnaryInterceptor(),
		),
		grpc.StreamInterceptor(grpcserver.RLSStreamInterceptor()),
	)

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
	pb.RegisterAuthServiceServer(grpcServer, authServer)
	pb.RegisterInvitationServiceServer(grpcServer, invitationServer)
	pb.RegisterETCMeisaiServiceServer(grpcServer, etcMeisaiServer)

	// Register health check service for Cloud Run
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	reflection.Register(grpcServer)

	// Use PORT env (Cloud Run sets this to 9090 for gRPC in sidecar setup)
	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}

	// Create HTTP mux for auth endpoints
	httpMux := http.NewServeMux()
	authHandler.RegisterRoutes(httpMux)

	// Create a handler that routes based on content-type and HTTP/2
	// gRPC uses HTTP/2 with content-type starting with "application/grpc"
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")
		// gRPC requests: HTTP/2 with application/grpc content type
		if r.ProtoMajor == 2 && strings.HasPrefix(contentType, "application/grpc") {
			grpcServer.ServeHTTP(w, r)
			return
		}
		httpMux.ServeHTTP(w, r)
	})

	// Use h2c to support HTTP/2 cleartext (required for gRPC without TLS)
	h2s := &http2.Server{}
	h2cHandler := h2c.NewHandler(handler, h2s)

	httpServer := &http.Server{
		Addr:    ":" + port,
		Handler: h2cHandler,
	}

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh

		log.Println("Shutting down servers...")
		grpcServer.GracefulStop()
		httpServer.Shutdown(context.Background())
	}()

	log.Printf("Starting server on port %s (HTTP + gRPC with h2c)", port)
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}

	log.Println("Server stopped")
}
