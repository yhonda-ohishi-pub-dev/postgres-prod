# PLAN.md - gRPCサービス拡充

## 現状

- Repository層: 29テーブル完了 ✅
- Proto定義: 27サービス完了 ✅
- gRPCサーバー実装: 27/27 完了 ✅
- main.go: 全27サービス登録完了 ✅
- ビルド確認: go build, go vet 成功 ✅
- OAuth2認証: Google/LINE対応完了 ✅
- RLS統合テスト: 完了 ✅
- Cloud Run: デプロイ済み、gRPC動作確認済み ⏳ (再デプロイ必要)
  - URL: https://postgres-prod-566bls5vfq-an.a.run.app

## 目標

29テーブル全てにgRPCサービスを提供

## 次のステップ（予定）

### Phase 4: Envoyサイドカー追加（gRPC-Web対応）

ブラウザからgRPC APIに直接アクセスするため、EnvoyプロキシをCloud Runサイドカーとして構成。

#### 構成
```
Browser (gRPC-Web) → Envoy (:8080) → gRPC Server (:9090)
```

#### 4-1. Envoy設定
- [ ] `envoy.yaml` - Envoyプロキシ設定
  - gRPC-Web → gRPCプロトコル変換
  - CORS設定（開発用: `*`）
  - HTTP/1.1 → HTTP/2変換
  - ヘルスチェック対応

#### 4-2. Cloud Run設定
- [ ] `cloudbuild.yaml` - サイドカービルド対応に更新
- [ ] `service.yaml` - Cloud Runサービス定義（サイドカー構成）
  - Envoyコンテナ（port 8080、ingress）
  - gRPCサーバーコンテナ（port 9090、内部）

#### 4-3. gRPCサーバー修正
- [ ] `cmd/server/main.go` - ポート9090で起動するよう変更
- [ ] 環境変数 `PORT=9090` をサイドカー構成で設定

#### 4-4. デプロイ・動作確認
- [ ] Cloud Runへデプロイ
- [ ] gRPC-Webクライアントから動作確認

#### 備考
- Envoyイメージ: `envoyproxy/envoy:v1.28-latest`
- Cloud Runのサイドカー機能（GA）を使用

---

### Phase 5: Cloud Runデプロイ確認

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
- OAuth2認証: Google/LINE対応完了
