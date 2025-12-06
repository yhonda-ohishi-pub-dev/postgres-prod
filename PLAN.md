# PLAN.md - gRPCサービス拡充

## 現状

- Repository層: 29テーブル完了
- gRPCサービス: OrganizationServiceのみ
- Cloud Run: デプロイ済み、gRPC動作確認済み

## 目標

29テーブル全てにgRPCサービスを提供

## 次のステップ（予定）

### Phase 1: Proto定義追加

proto/service.protoに26個のサービス定義を追加。以下の順で実装:

#### Phase 1-1: 簡単なテーブル ✅ 完了
| テーブル | メッセージ | フィールド数 |
|---------|-----------|------------|
| app_users | AppUser | 7 |
| user_organizations | UserOrganization | 7 |
| files | File | 7 |
| flickr_photo | FlickrPhoto | 4 |

#### Phase 1-2: カメラ関連テーブル ✅ 完了
| テーブル | メッセージ | フィールド数 |
|---------|-----------|------------|
| cam_files | CamFile | 7 |
| cam_file_exe | CamFileExe | 4 |
| cam_file_exe_stage | CamFileExeStage | 3 |

#### Phase 1-3: 車両関連テーブル
| テーブル | メッセージ | フィールド数 |
|---------|-----------|------------|
| ichiban_cars | IchibanCar | 12 |
| dtako_cars_ichiban_cars | DtakoCarsIchibanCars | 3 |
| uriage | Uriage | 7 |
| uriage_jisha | UriageJisha | 5 |

#### Phase 1-4: 車検関連テーブル
| テーブル | メッセージ | フィールド数 |
|---------|-----------|------------|
| car_inspection | CarInspection | 100+ |
| car_inspection_files | CarInspectionFile | 11 |
| car_inspection_files_a | CarInspectionFilesA | 11 |
| car_inspection_files_b | CarInspectionFilesB | 11 |
| car_inspection_deregistration | CarInspectionDeregistration | 9 |
| car_inspection_deregistration_files | CarInspectionDeregistrationFiles | 4 |
| car_ins_sheet_ichiban_cars | CarInsSheetIchibanCars | 7 |
| car_ins_sheet_ichiban_cars_a | CarInsSheetIchibanCarsA | 7 |

#### Phase 1-5: KUDG系テーブル
| テーブル | メッセージ | フィールド数 |
|---------|-----------|------------|
| kudgfry | Kudgfry | 28 |
| kudguri | Kudguri | 40 |
| kudgcst | Kudgcst | 30 |
| kudgful | Kudgful | 40+ |
| kudgsir | Kudgsir | 40+ |
| kudgivt | Kudgivt | 120+ |

#### Phase 1-6: Dtakologs
| テーブル | メッセージ | フィールド数 |
|---------|-----------|------------|
| dtakologs | Dtakologs | 58 |

#### Phase 1-7: buf generate
```bash
buf generate
```

### Phase 2: gRPCサーバー実装 (26ファイル)

pkg/grpc/配下に作成するファイル一覧:

| # | ファイル名 | サービス |
|---|-----------|---------|
| 1 | app_user_server.go | AppUserService |
| 2 | user_organization_server.go | UserOrganizationService |
| 3 | file_server.go | FileService |
| 4 | flickr_photo_server.go | FlickrPhotoService |
| 5 | cam_file_server.go | CamFileService |
| 6 | cam_file_exe_server.go | CamFileExeService |
| 7 | cam_file_exe_stage_server.go | CamFileExeStageService |
| 8 | ichiban_car_server.go | IchibanCarService |
| 9 | dtako_cars_ichiban_cars_server.go | DtakoCarsIchibanCarsService |
| 10 | uriage_server.go | UriageService |
| 11 | uriage_jisha_server.go | UriageJishaService |
| 12 | car_inspection_server.go | CarInspectionService |
| 13 | car_inspection_files_server.go | CarInspectionFilesService |
| 14 | car_inspection_files_a_server.go | CarInspectionFilesAService |
| 15 | car_inspection_files_b_server.go | CarInspectionFilesBService |
| 16 | car_inspection_deregistration_server.go | CarInspectionDeregistrationService |
| 17 | car_inspection_deregistration_files_server.go | CarInspectionDeregistrationFilesService |
| 18 | car_ins_sheet_ichiban_cars_server.go | CarInsSheetIchibanCarsService |
| 19 | car_ins_sheet_ichiban_cars_a_server.go | CarInsSheetIchibanCarsAService |
| 20 | kudgfry_server.go | KudgfryService |
| 21 | kudguri_server.go | KudguriService |
| 22 | kudgcst_server.go | KudgcstService |
| 23 | kudgful_server.go | KudgfulService |
| 24 | kudgsir_server.go | KudgsirService |
| 25 | kudgivt_server.go | KudgivtService |
| 26 | dtakologs_server.go | DtakologsService |

各サーバーは以下のRPCを実装:
- Create
- Get (主キーで取得)
- Update
- Delete
- List (ページネーション付き)
- ListByOrganization (組織IDでフィルタ)

### Phase 3: main.go更新

cmd/server/main.goに全27サービスを登録

### Phase 4: ビルド・デプロイ確認

```bash
go build -o server ./cmd/server
go test ./...
gcloud builds submit --config=cloudbuild.yaml ...
```

## 完了タスク

### Phase 1-1: 簡単なテーブルのProto定義追加 (2025-12-07)

**更新ファイル:**
- proto/service.proto: 4サービス（AppUser, UserOrganization, File, FlickrPhoto）のメッセージ・RPC定義追加
- pkg/pb/service.pb.go: protoc生成
- pkg/pb/service_grpc.pb.go: protoc生成

**追加内容:**
- AppUserService: 6 RPCs (Create, Get, GetByIamEmail, Update, Delete, List)
- UserOrganizationService: 7 RPCs (Create, Get, Update, Delete, List, ListByUser, ListByOrg)
- FileService: 6 RPCs (Create, Get, Update, Delete, List, ListByOrganization)
- FlickrPhotoService: 6 RPCs (Create, Get, Update, Delete, List, ListByOrganization)

**ビルド確認:** go build 成功

### Phase 1-2: カメラ関連テーブルのProto定義追加 (2025-12-07)

**更新ファイル:**
- proto/service.proto: 3サービス（CamFile, CamFileExe, CamFileExeStage）のメッセージ・RPC定義追加
- pkg/pb/service.pb.go: protoc生成
- pkg/pb/service_grpc.pb.go: protoc生成

**追加内容:**
- CamFileService: 6 RPCs (Create, Get, Update, Delete, List, ListByOrganization)
- CamFileExeService: 6 RPCs (Create, Get, Update, Delete, List, ListByOrganization)
- CamFileExeStageService: 6 RPCs (Create, Get, Update, Delete, List, ListByOrganization)

**ビルド確認:** go build 成功

## 参考

- 完了タスクの履歴は [docs/PLAN-EXECUTED.md](docs/PLAN-EXECUTED.md) を参照
- Repository層: 29/29 完了
- 統合テスト: 27/27 完了（bumon_codes, bumon_code_refsは除外）
