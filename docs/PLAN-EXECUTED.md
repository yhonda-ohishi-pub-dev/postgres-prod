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
