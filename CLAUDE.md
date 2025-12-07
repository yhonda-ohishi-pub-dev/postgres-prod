# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Go service for Cloud Run that connects to Cloud SQL PostgreSQL using IAM authentication. Provides HTTP health endpoints, gRPC CRUD services, and Repository layer for 29 database tables.

## Recent Changes

| Commit | Description |
|--------|-------------|
| 80824be | Phase 2完了: 全27gRPCサービス実装 |
| 92608b4 | ドキュメント更新: gRPC単独構成とMakefileを反映 |
| 62011a7 | Proto生成コード追加: pkg/pb/proto/ |
| 19c840f | ビルドツール追加: Makefileと実行スクリプト |
| 9c09aec | Cloud Run対応: gRPC単独サーバー構成に変更 |

## Build and Run

```bash
# Using Makefile (recommended)
make deps      # Download dependencies
make build     # Build locally
make run       # Build and run
make test      # Run tests
make proto     # Generate protobuf code
make check     # Run vet, test, build

# Manual commands
go mod tidy
go build -o server ./cmd/server

# Run locally (requires Cloud SQL Auth Proxy for local development)
export GCP_PROJECT_ID=your-project
export GCP_REGION=asia-northeast1
export CLOUDSQL_INSTANCE_NAME=your-instance
export DB_NAME=postgres
export DB_USER=your-iam-user
./server

# Build Docker image
docker build -t postgres-prod .

# Deploy via Cloud Build
gcloud builds submit --config=cloudbuild.yaml \
  --substitutions=_CLOUDSQL_INSTANCE=instance-name,_DB_USER=iam-user,_SERVICE_ACCOUNT=sa@project.iam.gserviceaccount.com
```

## Architecture

```
cmd/server/main.go           - Entry point, gRPC-only server (Cloud Run compatible)
internal/config/             - Environment configuration
pkg/db/cloudsql.go           - Cloud SQL connection with IAM auth (cloudsqlconn)
pkg/grpc/                    - gRPC server implementations (27 services)
  organization_server.go     - OrganizationService
  app_user_server.go         - AppUserService
  car_inspection_server.go   - CarInspectionService (largest, 510 lines)
  kudg*_server.go            - KUDG series (6 services)
  ... and more (27 total)
pkg/repository/              - Database CRUD operations (29 tables)
  organization.go            - Organization repository (with interface for mocking)
  app_users.go               - AppUsers repository
  cam_files.go               - CamFiles repository
  car_inspection.go          - CarInspection repository (largest, 48KB)
  ichiban_cars.go            - IchibanCars repository
  kudg*.go                   - KUDG series repositories (6 files)
  ... and more (29 total with integration tests)
pkg/pb/                      - Generated protobuf Go code
proto/                       - Protocol buffer definitions
  service.proto              - OrganizationService definition
Makefile                     - Build automation commands
```

## gRPC Services

The service exposes gRPC on PORT (default 8080, Cloud Run compatible) with **27 services**:

| Category | Services |
|----------|----------|
| Core | OrganizationService, AppUserService, UserOrganizationService, FileService |
| Media | FlickrPhotoService, CamFileService, CamFileExeService, CamFileExeStageService |
| Vehicle | IchibanCarService, DtakoCarsIchibanCarsService, UriageService, UriageJishaService |
| Inspection | CarInspectionService, CarInspectionFilesService, CarInspectionFilesAService, CarInspectionFilesBService |
| Deregistration | CarInspectionDeregistrationService, CarInspectionDeregistrationFilesService |
| Sheet | CarInsSheetIchibanCarsService, CarInsSheetIchibanCarsAService |
| KUDG | KudgfryService, KudguriService, KudgcstService, KudgfulService, KudgsirService, KudgivtService |
| Logs | DtakologsService |

Each service provides standard CRUD operations: `Create`, `Get`, `Update`, `Delete`, `List`

- **Health Check**: gRPC Health Check Protocol for Cloud Run startup/liveness probes

## Cloud SQL IAM Authentication

This service uses automatic IAM database authentication via `cloud.google.com/go/cloudsqlconn` with `WithIAMAuthN()`. No passwords required.

Requirements:
- Cloud SQL instance must have `cloudsql.iam_authentication` flag enabled
- Service account must have `roles/cloudsql.instanceUser` and `roles/cloudsql.client`
- Database user must be created as IAM user: `CREATE USER "sa@project.iam" WITH LOGIN`
- Set `DB_USER` to the IAM principal name (email without domain suffix for service accounts)

## Proto Definitions

Service definitions are in `proto/service.proto`. Generate Go code with:

```bash
buf generate
```

Import the generated package:

```go
import "postgres-prod/pkg/pb"
```

## Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| GCP_PROJECT_ID | GCP project ID | my-project |
| GCP_REGION | Cloud SQL region | asia-northeast1 |
| CLOUDSQL_INSTANCE_NAME | Cloud SQL instance name | my-instance |
| DB_NAME | Database name | postgres |
| DB_USER | IAM user (without @domain) | my-sa |
| PORT | gRPC server port (Cloud Run sets this) | 8080 |

## Deployment

### Deploy to Cloud Run (Direct)

```bash
gcloud run deploy postgres-prod \
  --image=asia-northeast1-docker.pkg.dev/cloudsql-sv/postgres-prod/postgres-prod:latest \
  --region=asia-northeast1 \
  --platform=managed \
  --allow-unauthenticated \
  --add-cloudsql-instances=cloudsql-sv:asia-northeast1:postgres-prod \
  --set-env-vars=GCP_PROJECT_ID=cloudsql-sv,GCP_REGION=asia-northeast1,CLOUDSQL_INSTANCE_NAME=postgres-prod,DB_NAME=myapp,DB_USER=747065218280-compute@developer \
  --use-http2 \
  --project=cloudsql-sv
```

### Service URL

```
https://postgres-prod-747065218280.asia-northeast1.run.app
```

## gRPC Reflection

Verify deployed services using grpcurl:

```bash
# List all services
grpcurl postgres-prod-747065218280.asia-northeast1.run.app:443 list

# Describe a service
grpcurl postgres-prod-747065218280.asia-northeast1.run.app:443 describe organization.OrganizationService

# Call a method (example: ListOrganizations)
grpcurl postgres-prod-747065218280.asia-northeast1.run.app:443 organization.OrganizationService/ListOrganizations
```
