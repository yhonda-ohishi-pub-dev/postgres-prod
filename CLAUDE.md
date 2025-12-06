# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Go service for Cloud Run that connects to Cloud SQL PostgreSQL using IAM authentication. Provides HTTP health endpoints, gRPC CRUD services, and Repository layer for 29 database tables.

## Recent Changes

| Commit | Description |
|--------|-------------|
| 62011a7 | Proto生成コード追加: pkg/pb/proto/ |
| 19c840f | ビルドツール追加: Makefileと実行スクリプト |
| 9c09aec | Cloud Run対応: gRPC単独サーバー構成に変更 |
| 6d903dd | Phase 1-6: DtakologsテーブルのProto定義追加 |
| b5627d3 | Phase 1-5: KUDG系テーブルのProto定義追加 |

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
pkg/grpc/                    - gRPC server implementations
  organization_server.go     - OrganizationService gRPC server
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

The service exposes gRPC on PORT (default 8080, Cloud Run compatible):

- **OrganizationService**: CRUD operations for Organizations
  - `CreateOrganization`, `GetOrganization`, `UpdateOrganization`, `DeleteOrganization`, `ListOrganizations`
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
