# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Go service for Cloud Run that connects to Cloud SQL PostgreSQL using IAM authentication. Provides gRPC CRUD services and Repository layer for 30 database tables, with OAuth2 authentication (Google/LINE), user invitation system, and Row-Level Security (RLS) for multi-tenant data isolation.

## Recent Changes

| Commit | Description |
|--------|-------------|
| 93ba7cf | Organization作成時にユーザー紐づけ自動化: JWT認証+トランザクション対応 |
| 83f7b9e | gitignore |
| 89c7442 | .gitignore: service-resolved.yaml追加（生成ファイル除外） |
| a60bdc8 | Makefile: ローカルデプロイ機能追加（Cloud Build不要） |
| cf2b3d9 | Cloud Build: ログ設定追加（CLOUD_LOGGING_ONLY） |

## Build and Run

```bash
# Using Makefile (recommended)
make deps      # Download dependencies
make build     # Build locally
make run       # Build and run
make test      # Run tests
make proto     # Generate protobuf code
make check     # Run vet, test, build

# Local deployment (no Cloud Build charges)
make deploy-local  # Docker build → push → Cloud Run deploy
make deploy-force  # Timestamp tag for forcing new revision

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
cmd/server/main.go           - Entry point, gRPC+HTTP server (h2c対応, Cloud Run compatible)
internal/config/             - Environment configuration
pkg/db/cloudsql.go           - Cloud SQL connection with IAM auth (cloudsqlconn)
pkg/http/                    - HTTP handlers
  auth_handler.go            - OAuth2 redirect endpoints (/auth/google, /auth/line, /health)
pkg/grpc/                    - gRPC server implementations (29 services)
  organization_server.go     - OrganizationService (CreateWithOwner: org+user_org同時作成)
  app_user_server.go         - AppUserService
  auth_server.go             - AuthService (OAuth2: Google/LINE)
  invitation_server.go       - InvitationService (ユーザー招待機能)
  car_inspection_server.go   - CarInspectionService (largest, 510 lines)
  interceptor.go             - RLS interceptor + JWT認証インターセプター
  kudg*_server.go            - KUDG series (6 services)
  ... and more (29 total)
pkg/auth/                    - OAuth2 authentication
  jwt.go                     - JWT token generation/validation
  google.go                  - Google OAuth2 client
  line.go                    - LINE OAuth2 client
pkg/db/rls.go                - Row-Level Security pool wrapper + トランザクション対応 (Begin/Tx)
pkg/repository/              - Database CRUD operations (30 tables)
  organization.go            - Organization repository (CreateWithOwner: トランザクションで所有者紐づけ)
  app_users.go               - AppUsers repository
  invitations.go             - Invitations repository (ユーザー招待)
  cam_files.go               - CamFiles repository
  car_inspection.go          - CarInspection repository (largest, 48KB)
  ichiban_cars.go            - IchibanCars repository
  kudg*.go                   - KUDG series repositories (6 files)
  ... and more (30 total with integration tests)
pkg/pb/                      - Generated protobuf Go code
proto/                       - Protocol buffer definitions
  service.proto              - OrganizationService definition
Makefile                     - Build automation commands
```

## gRPC Services

The service exposes gRPC on PORT (default 8080, Cloud Run compatible) with **29 services**:

| Category | Services |
|----------|----------|
| Auth | AuthService (OAuth2: Google/LINE認証), InvitationService (ユーザー招待) |
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

## HTTP Endpoints

gRPC+HTTP共存（h2c対応）により、同一ポートでHTTPエンドポイントも提供:

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/auth/google` | GET | Google OAuth2認証リダイレクト |
| `/auth/line` | GET | LINE OAuth2認証リダイレクト |
| `/health` | GET | ヘルスチェック（startup probe用） |

## Row-Level Security (RLS)

Multi-tenant data isolation using PostgreSQL RLS:

- **RLS Pool**: `pkg/db/rls.go` wraps DB connection to set `app.current_organization_id` per request
- **Transaction Support**: `Begin()` method returns `Tx` interface for transactional operations
- **gRPC Interceptor**: `pkg/grpc/interceptor.go` extracts `x-organization-id` header
- **JWT Authentication**: JWTインターセプターがAuthorizationヘッダーからuser_idを抽出しcontextに保存
- All repository operations automatically scoped to current organization

## OAuth2 Authentication

Google/LINE OAuth2 authentication with JWT tokens:

- **AuthService methods**:
  - `GetAuthURL`: Get OAuth authorization URL
  - `AuthWithGoogle` / `AuthWithLine`: Exchange auth code for JWT
  - `RefreshToken`: Refresh expired access token
  - `ValidateToken`: Validate JWT and return user info
- **oauth_accounts table**: Supports multiple OAuth providers per user

## User Invitation

Organization admins can invite users via email:

- **InvitationService methods**:
  - `CreateInvitation`: Create invitation with email and role
  - `GetInvitationByToken`: Get invitation details for accepting
  - `AcceptInvitation`: Accept and create user_organization
  - `CancelInvitation`: Cancel pending invitation
  - `ListInvitations`: List invitations for organization
  - `ResendInvitation`: Regenerate token and extend expiry

- **Invitation Flow**:
  1. Admin calls `CreateInvitation` with email and role
  2. Returns invite URL: `{FRONTEND_URL}/invite/{token}`
  3. Invited user authenticates and calls `AcceptInvitation`
  4. Creates `user_organization` record with specified role

- **invitations table**: Stores pending/accepted invitations with 7-day expiry

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
| JWT_SECRET | Secret for JWT signing | your-secret-key |
| GOOGLE_CLIENT_ID | Google OAuth2 client ID | xxx.apps.googleusercontent.com |
| GOOGLE_CLIENT_SECRET | Google OAuth2 client secret | GOCSPX-xxx |
| LINE_CHANNEL_ID | LINE OAuth2 channel ID | 123456789 |
| LINE_CHANNEL_SECRET | LINE OAuth2 channel secret | xxx |
| FRONTEND_URL | OAuth callback redirect URL | https://mtama-front.mtamaramu.com/auth/callback |

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
