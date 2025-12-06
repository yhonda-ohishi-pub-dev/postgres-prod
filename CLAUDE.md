# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Go service for Cloud Run that connects to Cloud SQL PostgreSQL using IAM authentication. Provides both HTTP health endpoints and gRPC CRUD services for Organizations.

## Recent Changes

| Commit | Description |
|--------|-------------|
| b4e3ab7 | Cloud Runデプロイ対応: gRPCサーバー追加とセキュリティ改善 |
| e42653d | Organizations gRPC CRUDサービスの実装 |
| 5402165 | Cloud SQL PostgreSQL接続サービスの初期実装 |

## Build and Run

```bash
# Download dependencies
go mod tidy

# Build locally
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
cmd/server/main.go           - Entry point, HTTP + gRPC server setup
internal/config/             - Environment configuration
pkg/db/cloudsql.go           - Cloud SQL connection with IAM auth (cloudsqlconn)
pkg/handlers/                - HTTP handlers (health check)
pkg/grpc/                    - gRPC server implementations
  organization_server.go     - OrganizationService gRPC server
pkg/repository/              - Database CRUD operations
  organization.go            - Organization repository (with interface for mocking)
pkg/pb/                      - Generated protobuf Go code
proto/                       - Protocol buffer definitions
  service.proto              - OrganizationService definition
```

## gRPC Services

The service exposes gRPC on port 50051:

- **OrganizationService**: CRUD operations for Organizations
  - `CreateOrganization`, `GetOrganization`, `UpdateOrganization`, `DeleteOrganization`, `ListOrganizations`

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
| PORT | HTTP port | 8080 |
| GRPC_PORT | gRPC port | 50051 |
