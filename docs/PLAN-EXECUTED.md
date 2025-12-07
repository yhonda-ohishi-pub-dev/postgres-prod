# PLAN-EXECUTED.md - 完了タスク履歴

このファイルには、plan.mdから完了したタスクを移動して記録します。

---

## 完了: 全テーブルCRUD Repository層作成 (2025-12-07)

### 概要

cloudsqlの`proto/models.proto`に定義された全テーブル（29テーブル）に対してCRUD Repositoryを作成。

### コミット履歴

| コミット | 説明 |
|---------|------|
| bf3cb64 | Repository層追加: 10テーブル分のCRUD実装と統合テスト |
| 75521ad | Repository層完成: 残り17テーブルのCRUD実装 |

### 成果物

```
pkg/repository/
  organization.go                          # CRUD + unit test
  organization_integration_test.go         # 統合テスト
  app_users.go                             # CRUD
  app_users_integration_test.go            # 統合テスト
  user_organizations.go                    # CRUD
  user_organizations_integration_test.go   # 統合テスト
  files.go                                 # CRUD
  files_integration_test.go                # 統合テスト
  flickr_photo.go                          # CRUD
  flickr_photo_integration_test.go         # 統合テスト
  cam_files.go                             # CRUD
  cam_files_integration_test.go            # 統合テスト
  cam_file_exe.go                          # CRUD
  cam_file_exe_integration_test.go         # 統合テスト
  cam_file_exe_stage.go                    # CRUD
  cam_file_exe_stage_integration_test.go   # 統合テスト
  ichiban_cars.go                          # CRUD
  ichiban_cars_integration_test.go         # 統合テスト
  dtako_cars_ichiban_cars.go               # CRUD
  dtako_cars_ichiban_cars_integration_test.go # 統合テスト
  car_inspection.go                        # CRUD
  car_inspection_files.go                  # CRUD
  car_inspection_files_a.go                # CRUD
  car_inspection_files_b.go                # CRUD
  car_inspection_deregistration.go         # CRUD
  car_inspection_deregistration_files.go   # CRUD
  car_ins_sheet_ichiban_cars.go            # CRUD
  car_ins_sheet_ichiban_cars_a.go          # CRUD
  dtakologs.go                             # CRUD
  kudguri.go                               # CRUD
  kudgcst.go                               # CRUD
  kudgfry.go                               # CRUD
  kudgful.go                               # CRUD
  kudgivt.go                               # CRUD
  kudgsir.go                               # CRUD
  uriage.go                                # CRUD
  uriage_jisha.go                          # CRUD
```

### 進捗詳細 (Repository層 29/29 完了, 統合テスト 27/27 完了)

| # | テーブル名 | Repository | Integration Test |
|---|-----------|:----------:|:----------------:|
| 1 | organizations | [x] | [x] |
| 2 | app_users | [x] | [x] |
| 3 | user_organizations | [x] | [x] |
| 4 | files | [x] | [x] |
| 5 | flickr_photo | [x] | [x] |
| 6 | cam_files | [x] | [x] |
| 7 | cam_file_exe | [x] | [x] |
| 8 | cam_file_exe_stage | [x] | [x] |
| 9 | ichiban_cars | [x] | [x] |
| 10 | dtako_cars_ichiban_cars | [x] | [x] |
| 11 | car_inspection | [x] | [x] |
| 12 | car_inspection_files | [x] | [x] |
| 13 | car_inspection_files_a | [x] | [x] |
| 14 | car_inspection_files_b | [x] | [x] |
| 15 | car_inspection_deregistration | [x] | [x] |
| 16 | car_inspection_deregistration_files | [x] | [x] |
| 17 | car_ins_sheet_ichiban_cars | [x] | [x] |
| 18 | car_ins_sheet_ichiban_cars_a | [x] | [x] |
| 19 | dtakologs | [x] | [x] |
| 20 | kudguri | [x] | [x] |
| 21 | kudgcst | [x] | [x] |
| 22 | kudgfry | [x] | [x] |
| 23 | kudgful | [x] | [x] |
| 24 | kudgivt | [x] | [x] |
| 25 | kudgsir | [x] | [x] |
| 26 | uriage | [x] | [x] |
| 27 | uriage_jisha | [x] | [x] |
| 28 | bumon_codes | [x] (既存) | - |
| 29 | bumon_code_refs | [x] (既存) | - |

### 追加コミット (2025-12-07)

| コミット | 説明 |
|---------|------|
| TBD | 統合テスト完成: 残り17テーブル分の統合テストを追加 |

### 追加作成ファイル

```
pkg/repository/
  uriage_integration_test.go
  uriage_jisha_integration_test.go
  dtakologs_integration_test.go
  kudguri_integration_test.go
  kudgcst_integration_test.go
  kudgfry_integration_test.go
  kudgful_integration_test.go
  kudgivt_integration_test.go
  kudgsir_integration_test.go
  car_inspection_files_integration_test.go
  car_inspection_files_a_integration_test.go
  car_inspection_files_b_integration_test.go
  car_inspection_deregistration_integration_test.go
  car_inspection_deregistration_files_integration_test.go
  car_ins_sheet_ichiban_cars_integration_test.go
  car_ins_sheet_ichiban_cars_a_integration_test.go
```

### 技術詳細

#### Repository層のパターン

- DB interface: テスト用モック対応
- NewXxxRepository(pool): 本番用コンストラクタ
- NewXxxRepositoryWithDB(db): テスト用コンストラクタ
- CRUD: Create, Get, Update, Delete, List
- Pagination: limit/offset (デフォルト10, 最大100)
- Soft delete: deleted フィールド (一部テーブル)
- Hard delete: uriage, uriage_jisha 等

#### 複合主キー対応

多くのテーブルで複合主キーを使用:
- car_inspection: 6列 (organization_id, ElectCertMgNo, ElectCertPublishdateE/Y/M/D)
- dtakologs: 3列 (organization_id, DataDateTime, VehicleCD)
- uriage: 4列 (name, bumon, date, organization_id)

#### DBカラム名

PostgreSQLのcase-sensitiveカラムにはダブルクォート使用:
- `"ElectCertMgNo"`, `"CarId"`, `"DataDateTime"` 等

---

## 完了: gRPCサービス拡充 Phase 1-2 (2025-12-07)

### 概要

29テーブル全てにgRPCサービスを提供するためのProto定義追加とgRPCサーバー実装。

### Phase 1: Proto定義追加

#### Phase 1-1: 簡単なテーブルのProto定義追加

**更新ファイル:**
- proto/service.proto: 4サービス（AppUser, UserOrganization, File, FlickrPhoto）のメッセージ・RPC定義追加
- pkg/pb/service.pb.go: protoc生成
- pkg/pb/service_grpc.pb.go: protoc生成

**追加内容:**
- AppUserService: 6 RPCs (Create, Get, GetByIamEmail, Update, Delete, List)
- UserOrganizationService: 7 RPCs (Create, Get, Update, Delete, List, ListByUser, ListByOrg)
- FileService: 6 RPCs (Create, Get, Update, Delete, List, ListByOrganization)
- FlickrPhotoService: 6 RPCs (Create, Get, Update, Delete, List, ListByOrganization)

#### Phase 1-2: カメラ関連テーブルのProto定義追加

**更新ファイル:**
- proto/service.proto: 3サービス（CamFile, CamFileExe, CamFileExeStage）のメッセージ・RPC定義追加
- pkg/pb/service.pb.go: protoc生成
- pkg/pb/service_grpc.pb.go: protoc生成

**追加内容:**
- CamFileService: 6 RPCs (Create, Get, Update, Delete, List, ListByOrganization)
- CamFileExeService: 6 RPCs (Create, Get, Update, Delete, List, ListByOrganization)
- CamFileExeStageService: 6 RPCs (Create, Get, Update, Delete, List, ListByOrganization)

#### Phase 1-3: 車両関連テーブルのProto定義追加

**更新ファイル:**
- proto/service.proto: 4サービス（IchibanCar, DtakoCarsIchibanCars, Uriage, UriageJisha）のメッセージ・RPC定義追加
- pkg/pb/service.pb.go: protoc生成
- pkg/pb/service_grpc.pb.go: protoc生成

**追加内容:**
- IchibanCarService: 6 RPCs (Create, Get, Update, Delete, List, ListByOrganization)
- DtakoCarsIchibanCarsService: 6 RPCs (Create, Get, Update, Delete, List, ListByOrganization)
- UriageService: 6 RPCs (Create, Get, Update, Delete, List, ListByOrganization)
- UriageJishaService: 6 RPCs (Create, Get, Update, Delete, List, ListByOrganization)

#### Phase 1-4: 車検関連テーブルのProto定義追加

**更新ファイル:**
- proto/service.proto: 8サービス（CarInspection, CarInspectionFiles, CarInspectionFilesA, CarInspectionFilesB, CarInspectionDeregistration, CarInspectionDeregistrationFiles, CarInsSheetIchibanCars, CarInsSheetIchibanCarsA）のメッセージ・RPC定義追加
- pkg/pb/service.pb.go: protoc生成 (326KB)
- pkg/pb/service_grpc.pb.go: protoc生成 (162KB)

**追加内容:**
- CarInspectionService: 6 RPCs - 98フィールド
- CarInspectionFilesService: 6 RPCs - 11フィールド
- CarInspectionFilesAService: 6 RPCs - 11フィールド
- CarInspectionFilesBService: 6 RPCs - 11フィールド
- CarInspectionDeregistrationService: 6 RPCs - 9フィールド
- CarInspectionDeregistrationFilesService: 6 RPCs - 4フィールド
- CarInsSheetIchibanCarsService: 6 RPCs - 7フィールド
- CarInsSheetIchibanCarsAService: 6 RPCs - 7フィールド

#### Phase 1-5: KUDG系テーブルのProto定義追加

**更新ファイル:**
- proto/service.proto: 6サービス（Kudgfry, Kudguri, Kudgcst, Kudgful, Kudgsir, Kudgivt）のメッセージ・RPC定義追加 (3243行に拡張)
- pkg/pb/service.pb.go: protoc生成 (1.2MB)
- pkg/pb/service_grpc.pb.go: protoc生成 (365KB)

**追加内容:**
- KudgfryService: 6 RPCs - 28フィールド
- KudguriService: 6 RPCs - 40フィールド
- KudgcstService: 6 RPCs - 30フィールド
- KudgfulService: 6 RPCs - 40フィールド
- KudgsirService: 6 RPCs - 40フィールド
- KudgivtService: 6 RPCs - 94フィールド (最大)

#### Phase 1-6: DtakologsテーブルのProto定義追加

**更新ファイル:**
- proto/service.proto: DtakologsService のメッセージ・RPC定義追加 (3483行に拡張)
- pkg/pb/service.pb.go: protoc生成 (1.2MB)
- pkg/pb/service_grpc.pb.go: protoc生成 (365KB)

**追加内容:**
- DtakologsService: 6 RPCs - 57フィールド

### Phase 2: gRPCサーバー実装 (27ファイル)

**作成ファイル (27個):**
- pkg/grpc/app_user_server.go (174行)
- pkg/grpc/user_organization_server.go (191行)
- pkg/grpc/file_server.go (216行)
- pkg/grpc/flickr_photo_server.go (197行)
- pkg/grpc/cam_file_server.go (203行)
- pkg/grpc/cam_file_exe_server.go (175行)
- pkg/grpc/cam_file_exe_stage_server.go (140行)
- pkg/grpc/ichiban_car_server.go (289行)
- pkg/grpc/dtako_cars_ichiban_cars_server.go (177行)
- pkg/grpc/uriage_server.go (248行)
- pkg/grpc/uriage_jisha_server.go (226行)
- pkg/grpc/car_inspection_server.go (510行)
- pkg/grpc/car_inspection_files_server.go (227行)
- pkg/grpc/car_inspection_files_a_server.go (233行)
- pkg/grpc/car_inspection_files_b_server.go (216行)
- pkg/grpc/car_inspection_deregistration_server.go (239行)
- pkg/grpc/car_inspection_deregistration_files_server.go (240行)
- pkg/grpc/car_ins_sheet_ichiban_cars_server.go (287行)
- pkg/grpc/car_ins_sheet_ichiban_cars_a_server.go (287行)
- pkg/grpc/kudgfry_server.go (307行)
- pkg/grpc/kudguri_server.go (328行)
- pkg/grpc/kudgcst_server.go (293行)
- pkg/grpc/kudgful_server.go (322行)
- pkg/grpc/kudgsir_server.go (317行)
- pkg/grpc/kudgivt_server.go (557行)
- pkg/grpc/dtakologs_server.go (421行)
- pkg/grpc/organization_server.go (既存)

**更新ファイル:**
- cmd/server/main.go: 全27サービス登録
- pkg/pb/service.pb.go: protoc再生成
- pkg/pb/service_grpc.pb.go: protoc再生成

### コミット履歴

| コミット | 説明 |
|---------|------|
| 80824be | Phase 2完了: 全27gRPCサービス実装 |
| 752206e | RLS対応: gRPCリクエストごとにorganization_idを設定 |
| 66dbeca | ドキュメント更新: 27サービス詳細とデプロイ情報追加 |

---

## 完了: RLS統合テスト (2025-12-07)

### 概要

PostgreSQLのRow Level Security (RLS) ポリシーが正しく動作することを検証する統合テストを作成。
RLSPoolラッパーを使用して、異なるorganization間でのデータ分離を確認。

### 作成ファイル

```
pkg/db/rls_integration_test.go  # RLS統合テスト (343行)
```

### 更新ファイル

```
pkg/db/rls.go  # セッション変数名を修正 (app.organization_id → app.current_organization_id)
               # SET文の代わりにset_config()関数を使用
```

### テスト内容

1. **TestRLS_IsolationBetweenOrganizations**
   - OrgA_CanOnlySeeOwnData: organization Aのコンテキストで自分のデータのみ見えることを確認
   - OrgB_CanOnlySeeOwnData: organization Bのコンテキストで自分のデータのみ見えることを確認
   - ListQuery_ReturnsOnlyOwnData: List操作でも同様のフィルタリングが適用されることを確認

2. **TestRLS_ContextPropagation**
   - ContextWithOrganizationID: コンテキストにorganization_idが正しく保存されることを確認
   - ContextWithoutOrganizationID: organization_idなしのコンテキストが正しく処理されることを確認
   - RLSPool_RequiresOrganizationID: RLSコンテキストなしでエラーになることを確認
   - RLSPool_AcquiresConnectionWithContext: セッション変数が正しく設定されることを確認

3. **TestRLS_UpdateDeleteIsolation**
   - Update_OnlyAffectsOwnData: 他organizationのデータをUPDATEできないことを確認
   - Delete_OnlyAffectsOwnData: 他organizationのデータをDELETEできないことを確認

### 修正内容

- `SetRLSContext`関数: `SET app.organization_id = $1` → `SELECT set_config('app.current_organization_id', $1, false)`
  - PostgreSQLのSET文はパラメータ化クエリをサポートしないため、`set_config()`関数を使用
  - DBのRLSポリシーは`app.current_organization_id`を参照しているため、変数名も修正

### テスト実行方法

```bash
go test -tags=integration -v ./pkg/db/... -run TestRLS
```

### テスト結果

```
=== RUN   TestRLS_IsolationBetweenOrganizations
--- PASS: TestRLS_IsolationBetweenOrganizations (0.05s)
=== RUN   TestRLS_ContextPropagation
--- PASS: TestRLS_ContextPropagation (0.01s)
=== RUN   TestRLS_UpdateDeleteIsolation
--- PASS: TestRLS_UpdateDeleteIsolation (0.03s)
PASS
```

---

## 完了: Phase 3 OAuth2認証機能の実装 (2025-12-07)

### 概要

OAuth2認証機能（Google/LINE対応）を実装。JWTトークン発行・検証、認証プロバイダークライアント、gRPC AuthServiceを追加。

### コミット履歴

| コミット | 説明 |
|---------|------|
| ea4baf1 | OAuth2認証機能追加: Google/LINE対応 |
| 1549526 | 統合テスト修正: OAuth2対応のAPI変更に追従 |
| 5182d57 | ドキュメント更新: OAuth2/RLS機能を反映 |

### 作成ファイル

```
pkg/auth/
  jwt.go                    # JWT発行・検証
  google.go                 # Google OAuth2クライアント
  line.go                   # LINE OAuth2クライアント

pkg/repository/
  oauth_accounts.go         # oauth_accountsテーブルCRUD

pkg/grpc/
  auth_server.go            # AuthService実装
```

### 更新ファイル

```
cmd/server/main.go          # AuthService登録
CLAUDE.md                   # OAuth2/RLS機能を反映
README.md                   # ドキュメント更新
```

### 機能詳細

#### 3-1. テーブル設計・Repository
- `oauth_accounts`テーブル用Repository作成
  - Create, GetByProvider, GetByAppUserID, Update, Delete

#### 3-2. 認証ロジック
- `pkg/auth/jwt.go` - JWT発行・検証
- `pkg/auth/google.go` - Google OAuth2クライアント
- `pkg/auth/line.go` - LINE OAuth2クライアント

#### 3-3. AuthService (gRPC)
- `GetAuthURL` - OAuth認可URLを取得
- `AuthWithGoogle` - Google認証コードをJWTに交換
- `AuthWithLine` - LINE認証コードをJWTに交換
- `RefreshToken` - リフレッシュトークンでアクセストークン更新
- `ValidateToken` - JWTを検証しユーザー情報を返却

#### 3-4. 環境変数
- JWT_SECRET
- GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET
- LINE_CHANNEL_ID, LINE_CHANNEL_SECRET

---

## 完了: Phase 4 Envoyサイドカー追加 (2025-12-07)

### 概要

ブラウザからgRPC APIに直接アクセスするため、EnvoyプロキシをCloud Runサイドカーとして構成。gRPC-WebからgRPCへのプロトコル変換を実現。

### 構成

```
Browser (gRPC-Web) → Envoy (:8080) → gRPC Server (:9090)
```

### 作成ファイル

```
envoy.yaml           # Envoyプロキシ設定 (gRPC-Web変換、CORS)
Dockerfile.envoy     # Envoyイメージビルド用
service.yaml         # Cloud Runサービス定義 (サイドカー構成)
```

### 更新ファイル

```
cloudbuild.yaml      # サイドカービルド対応に更新
```

### 技術詳細

#### Envoy設定 (envoy.yaml)
- gRPC-Web → gRPCプロトコル変換 (envoy.filters.http.grpc_web)
- CORS設定（開発用: `*`）
- HTTP/1.1 → HTTP/2変換
- バックエンドは127.0.0.1:9090 (gRPCサーバー)

#### Cloud Run サイドカー構成 (service.yaml)
- Envoyコンテナ: port 8080 (ingress)
- gRPCサーバーコンテナ: port 9090 (内部)
- container-dependencies: gRPC-serverはEnvoyに依存

#### cloudbuild.yaml
- gRPCサーバーイメージビルド
- Envoyサイドカーイメージビルド (Dockerfile.envoy)
- service.yamlによるデプロイ

---

## 完了: Phase 4-4 Cloud Runデプロイ・動作確認 (2025-12-07)

### 概要

Cloud Buildを使用してgRPCサーバー + Envoyサイドカーをデプロイし、28サービスの動作を確認。

### デプロイ結果

**サービスURL:** https://postgres-prod-747065218280.asia-northeast1.run.app

**デプロイコマンド:**
```bash
gcloud builds submit --config=cloudbuild.yaml \
  --substitutions=_CLOUDSQL_INSTANCE=postgres-prod,_DB_USER=747065218280-compute@developer,_SERVICE_ACCOUNT=747065218280-compute@developer.gserviceaccount.com
```

### 動作確認

```bash
# サービス一覧取得
grpcurl -H "x-organization-id: test-org" postgres-prod-747065218280.asia-northeast1.run.app:443 list
```

**確認済みサービス (28個 + Health + Reflection):**
- grpc.health.v1.Health
- grpc.reflection.v1.ServerReflection
- organization.AuthService
- organization.OrganizationService
- organization.AppUserService
- organization.UserOrganizationService
- organization.FileService
- organization.FlickrPhotoService
- organization.CamFileService
- organization.CamFileExeService
- organization.CamFileExeStageService
- organization.IchibanCarService
- organization.DtakoCarsIchibanCarsService
- organization.UriageService
- organization.UriageJishaService
- organization.CarInspectionService
- organization.CarInspectionFilesService
- organization.CarInspectionFilesAService
- organization.CarInspectionFilesBService
- organization.CarInspectionDeregistrationService
- organization.CarInspectionDeregistrationFilesService
- organization.CarInsSheetIchibanCarsService
- organization.CarInsSheetIchibanCarsAService
- organization.KudgfryService
- organization.KudguriService
- organization.KudgcstService
- organization.KudgfulService
- organization.KudgsirService
- organization.KudgivtService
- organization.DtakologsService

### 備考

- Cloud BuildのサービスアカウントにIAM Policy設定権限が不足していたため、`run.services.setIamPolicy`のエラーが発生したが、デプロイ自体は成功
- デフォルトのCompute Service Account (`747065218280-compute@developer.gserviceaccount.com`) を使用
