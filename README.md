# postgres-prod

Cloud Run上で動作するGoサービス。Cloud SQL PostgreSQLにIAM認証で接続し、gRPC APIを提供します。

## Features

- **Cloud SQL IAM認証**: パスワード不要のセキュアな接続
- **gRPC API**: 27サービス（29テーブル対応）のCRUD API
- **Repository層**: 29テーブル分のCRUD実装（統合テスト完備）
- **Cloud Run対応**: 本番環境ですぐにデプロイ可能
- **ローカル開発対応**: Cloud SQL Auth Proxyでの開発をサポート

## Quick Start

### Prerequisites

- Go 1.21+
- Google Cloud SDK
- Cloud SQL Auth Proxy (ローカル開発用)
- buf (Protocol Buffers生成用)

### Build

```bash
# Makefile を使用（推奨）
make deps      # 依存関係をダウンロード
make build     # ビルド
make run       # ビルドして実行
make test      # テスト実行
make proto     # Protobufコード生成
make check     # vet, test, build を実行

# 手動コマンド
go mod tidy
go build -o server ./cmd/server
```

### Run Locally

Cloud SQL Auth Proxyを起動:
```bash
cloud-sql-proxy --auto-iam-authn PROJECT:REGION:INSTANCE
```

サーバーを起動:
```bash
export GCP_PROJECT_ID=your-project
export GCP_REGION=asia-northeast1
export CLOUDSQL_INSTANCE_NAME=your-instance
export DB_NAME=postgres
export DB_USER=your-iam-user
./server
```

## API

### gRPC (port 8080, Cloud Run compatible)

27サービスが利用可能。各サービスは標準CRUD操作（Create, Get, Update, Delete, List）を提供:

| カテゴリ | サービス |
|----------|----------|
| Core | OrganizationService, AppUserService, UserOrganizationService, FileService |
| Media | FlickrPhotoService, CamFileService, CamFileExeService, CamFileExeStageService |
| Vehicle | IchibanCarService, DtakoCarsIchibanCarsService, UriageService, UriageJishaService |
| Inspection | CarInspectionService, CarInspectionFilesService, CarInspectionFilesAService, CarInspectionFilesBService |
| Deregistration | CarInspectionDeregistrationService, CarInspectionDeregistrationFilesService |
| Sheet | CarInsSheetIchibanCarsService, CarInsSheetIchibanCarsAService |
| KUDG | KudgfryService, KudguriService, KudgcstService, KudgfulService, KudgsirService, KudgivtService |
| Logs | DtakologsService |

**Health Check**
- gRPC Health Check Protocol（Cloud Run のスタートアップ/ライブネスプローブ用）

## Project Structure

```
cmd/server/main.go       - エントリーポイント（27サービス登録）
internal/config/         - 環境設定
pkg/
  db/cloudsql.go         - Cloud SQL接続（IAM認証）
  grpc/                  - gRPCサーバー実装（27サービス）
  handlers/              - HTTPハンドラー
  pb/                    - 生成されたProtobufコード
  repository/            - データベース操作（29テーブル分のCRUD + 統合テスト）
proto/service.proto      - gRPCサービス定義
```

## Deployment

### Cloud Build

```bash
gcloud builds submit --config=cloudbuild.yaml \
  --substitutions=_CLOUDSQL_INSTANCE=instance,_DB_USER=user,_SERVICE_ACCOUNT=sa@project.iam.gserviceaccount.com
```

### Requirements

- Cloud SQLインスタンスで `cloudsql.iam_authentication` フラグを有効化
- サービスアカウントに `roles/cloudsql.instanceUser` と `roles/cloudsql.client` を付与
- データベースにIAMユーザーを作成: `CREATE USER "sa@project.iam" WITH LOGIN`

## Proto Generation

```bash
buf generate
```

## Environment Variables

| Variable | Description |
|----------|-------------|
| GCP_PROJECT_ID | GCPプロジェクトID |
| GCP_REGION | Cloud SQLリージョン |
| CLOUDSQL_INSTANCE_NAME | Cloud SQLインスタンス名 |
| DB_NAME | データベース名 |
| DB_USER | IAMユーザー名 |
| PORT | gRPCサーバーポート (default: 8080, Cloud Runが設定) |

## License

Private
