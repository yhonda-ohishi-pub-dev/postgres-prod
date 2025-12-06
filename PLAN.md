# PLAN.md - gRPCサービス拡充

## 現状

- Repository層: 29テーブル完了 ✅
- Proto定義: 27サービス完了 ✅
- gRPCサーバー実装: 27/27 完了 ✅
- main.go: 全27サービス登録完了 ✅
- ビルド確認: go build, go vet 成功 ✅
- Cloud Run: デプロイ済み、gRPC動作確認済み ⏳ (再デプロイ必要)
  - URL: https://postgres-prod-566bls5vfq-an.a.run.app
  - 現在: OrganizationServiceのみ (Phase 2完了後は全27サービス利用可能)

## 目標

29テーブル全てにgRPCサービスを提供

## 次のステップ（予定）

### Phase 3: Cloud Runデプロイ確認

```bash
gcloud builds submit --config=cloudbuild.yaml \
  --substitutions=_CLOUDSQL_INSTANCE=postgres-prod,_DB_USER=m.tama.ramu,_SERVICE_ACCOUNT=cloudsql-sv@cloudsql-sv.iam.gserviceaccount.com
```

デプロイ後、grpcurl等で全27サービスの動作確認。

---

## 完了したフェーズ

### Phase 1: Proto定義追加 ✅

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

#### Phase 1-3: 車両関連テーブル ✅ 完了
| テーブル | メッセージ | フィールド数 |
|---------|-----------|------------|
| ichiban_cars | IchibanCar | 12 |
| dtako_cars_ichiban_cars | DtakoCarsIchibanCars | 3 |
| uriage | Uriage | 7 |
| uriage_jisha | UriageJisha | 5 |

#### Phase 1-4: 車検関連テーブル ✅ 完了
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

#### Phase 1-5: KUDG系テーブル ✅ 完了
| テーブル | メッセージ | フィールド数 |
|---------|-----------|------------|
| kudgfry | Kudgfry | 28 |
| kudguri | Kudguri | 40 |
| kudgcst | Kudgcst | 30 |
| kudgful | Kudgful | 40+ |
| kudgsir | Kudgsir | 40+ |
| kudgivt | Kudgivt | 120+ |

#### Phase 1-6: Dtakologs ✅ 完了
| テーブル | メッセージ | フィールド数 |
|---------|-----------|------------|
| dtakologs | Dtakologs | 58 |

#### Phase 1-7: buf generate ✅ 完了
```bash
buf generate
```

### Phase 2: gRPCサーバー実装 (26ファイル) ✅

**完了: 2025-12-07**

全27サービスのgRPCサーバー実装完了。詳細は「Phase 2: gRPCサーバー実装 (26ファイル) - 完了 (2025-12-07)」セクション参照。

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

### Phase 1-3: 車両関連テーブルのProto定義追加 (2025-12-07)

**更新ファイル:**
- proto/service.proto: 4サービス（IchibanCar, DtakoCarsIchibanCars, Uriage, UriageJisha）のメッセージ・RPC定義追加
- pkg/pb/service.pb.go: protoc生成
- pkg/pb/service_grpc.pb.go: protoc生成

**追加内容:**
- IchibanCarService: 6 RPCs (Create, Get, Update, Delete, List, ListByOrganization)
- DtakoCarsIchibanCarsService: 6 RPCs (Create, Get, Update, Delete, List, ListByOrganization)
- UriageService: 6 RPCs (Create, Get, Update, Delete, List, ListByOrganization)
- UriageJishaService: 6 RPCs (Create, Get, Update, Delete, List, ListByOrganization)

**ビルド確認:** go build 成功

### Phase 1-4: 車検関連テーブルのProto定義追加 (2025-12-07)

**更新ファイル:**
- proto/service.proto: 8サービス（CarInspection, CarInspectionFiles, CarInspectionFilesA, CarInspectionFilesB, CarInspectionDeregistration, CarInspectionDeregistrationFiles, CarInsSheetIchibanCars, CarInsSheetIchibanCarsA）のメッセージ・RPC定義追加
- pkg/pb/service.pb.go: protoc生成 (326KB)
- pkg/pb/service_grpc.pb.go: protoc生成 (162KB)

**追加内容:**
- CarInspectionService: 6 RPCs (Create, Get, Update, Delete, List, ListByOrganization) - 98フィールド
- CarInspectionFilesService: 6 RPCs (Create, Get, Update, Delete, List, ListByOrganization) - 11フィールド
- CarInspectionFilesAService: 6 RPCs (Create, Get, Update, Delete, List, ListByOrganization) - 11フィールド
- CarInspectionFilesBService: 6 RPCs (Create, Get, Update, Delete, List, ListByOrganization) - 11フィールド
- CarInspectionDeregistrationService: 6 RPCs (Create, Get, Update, Delete, List, ListByOrganization) - 9フィールド
- CarInspectionDeregistrationFilesService: 6 RPCs (Create, Get, Update, Delete, List, ListByOrganization) - 4フィールド
- CarInsSheetIchibanCarsService: 6 RPCs (Create, Get, Update, Delete, List, ListByOrganization) - 7フィールド
- CarInsSheetIchibanCarsAService: 6 RPCs (Create, Get, Update, Delete, List, ListByOrganization) - 7フィールド

**ビルド確認:** go build 成功

### Phase 1-5: KUDG系テーブルのProto定義追加 (2025-12-07)

**更新ファイル:**
- proto/service.proto: 6サービス（Kudgfry, Kudguri, Kudgcst, Kudgful, Kudgsir, Kudgivt）のメッセージ・RPC定義追加 (3243行に拡張)
- pkg/pb/service.pb.go: protoc生成 (1.2MB)
- pkg/pb/service_grpc.pb.go: protoc生成 (365KB)

**追加内容:**
- KudgfryService: 6 RPCs (Create, Get, Update, Delete, List, ListByOrganization) - 28フィールド
- KudguriService: 6 RPCs (Create, Get, Update, Delete, List, ListByOrganization) - 40フィールド
- KudgcstService: 6 RPCs (Create, Get, Update, Delete, List, ListByOrganization) - 30フィールド
- KudgfulService: 6 RPCs (Create, Get, Update, Delete, List, ListByOrganization) - 40フィールド
- KudgsirService: 6 RPCs (Create, Get, Update, Delete, List, ListByOrganization) - 40フィールド
- KudgivtService: 6 RPCs (Create, Get, Update, Delete, List, ListByOrganization) - 94フィールド (最大)

**ビルド確認:** go build 成功

### Phase 1-6: DtakologsテーブルのProto定義追加 (2025-12-07)

**更新ファイル:**
- proto/service.proto: DtakologsService のメッセージ・RPC定義追加 (3483行に拡張)
- pkg/pb/service.pb.go: protoc生成 (1.2MB)
- pkg/pb/service_grpc.pb.go: protoc生成 (365KB)

**追加内容:**
- DtakologsService: 6 RPCs (Create, Get, Update, Delete, List, ListByOrganization) - 57フィールド

**ビルド確認:** go build 成功

### Phase 2: gRPCサーバー実装 (26ファイル) - 完了 (2025-12-07)

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

**ビルド確認:** go build, go vet 成功

## 参考

- 完了タスクの履歴は [docs/PLAN-EXECUTED.md](docs/PLAN-EXECUTED.md) を参照
- Repository層: 29/29 完了
- 統合テスト: 27/27 完了（bumon_codes, bumon_code_refsは除外）
- gRPCサーバー層: 27/27 完了
