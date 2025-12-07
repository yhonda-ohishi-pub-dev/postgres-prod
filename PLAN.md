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

### Phase 3: OAuth2認証機能の実装

userをOAuth2対応させ、LINE/Googleログインを実装する。

#### 3-1. テーブル設計・Repository
- [ ] `oauth_accounts`テーブル作成（マイグレーションSQL）
  - provider, provider_user_id, app_user_id, tokens等
- [ ] `pkg/repository/oauth_accounts.go` 新規作成
  - Create, GetByProvider, GetByAppUserID, Update, Delete

#### 3-2. 認証ロジック
- [ ] `pkg/auth/jwt.go` - JWT発行・検証
- [ ] `pkg/auth/google.go` - Google OAuth2クライアント
- [ ] `pkg/auth/line.go` - LINE OAuth2クライアント

#### 3-3. Proto定義・gRPCサービス
- [ ] `proto/auth.proto` - AuthService定義
  - AuthWithGoogle(code) → AuthResponse(jwt, refresh_token, user)
  - AuthWithLine(code) → AuthResponse
  - RefreshToken(refresh_token) → AuthResponse
- [ ] `pkg/grpc/auth_server.go` - AuthService実装

#### 3-4. 統合・テスト
- [ ] main.goにAuthService登録
- [ ] 環境変数追加（GOOGLE_CLIENT_ID, LINE_CHANNEL_ID等）
- [ ] 統合テスト作成

#### 備考
- app_usersテーブルの変更は別リポジトリで対応（指示待ち）
- クライアント: Web（将来的にモバイル対応）

---

### Phase 4: Cloud Runデプロイ確認

```bash
gcloud builds submit --config=cloudbuild.yaml \
  --substitutions=_CLOUDSQL_INSTANCE=postgres-prod,_DB_USER=m.tama.ramu,_SERVICE_ACCOUNT=cloudsql-sv@cloudsql-sv.iam.gserviceaccount.com
```

デプロイ後、grpcurl等で全27サービスの動作確認。

---

## 参考

- 完了タスクの履歴は [docs/PLAN-EXECUTED.md](docs/PLAN-EXECUTED.md) を参照
- Repository層: 29/29 完了
- 統合テスト: 27/27 完了（bumon_codes, bumon_code_refsは除外）
- gRPCサーバー層: 27/27 完了
