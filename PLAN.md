# PLAN.md - gRPCサービス拡充

## 現状

- Repository層: 29テーブル完了 ✅
- Proto定義: 27サービス完了 ✅
- gRPCサーバー実装: 27/27 完了 ✅
- main.go: 全27サービス登録完了 ✅
- ビルド確認: go build, go vet 成功 ✅
- OAuth2認証: Google/LINE対応完了 ✅
- RLS統合テスト: 完了 ✅
- Envoyサイドカー設定: 完了 ✅
- Cloud Run: デプロイ済み、gRPC動作確認済み ⏳ (再デプロイ必要)
  - URL: https://postgres-prod-566bls5vfq-an.a.run.app

## 目標

29テーブル全てにgRPCサービスを提供

## 次のステップ（予定）

### Phase 4-4: デプロイ・動作確認

```bash
gcloud builds submit --config=cloudbuild.yaml \
  --substitutions=_CLOUDSQL_INSTANCE=postgres-prod,_DB_USER=m.tama.ramu,_SERVICE_ACCOUNT=cloudsql-sv@cloudsql-sv.iam.gserviceaccount.com
```

デプロイ後、gRPC-Webクライアントから動作確認。

---

### Phase 5: Cloud Runデプロイ確認

デプロイ後、grpcurl等で全27サービスの動作確認。

---

## 参考

- 完了タスクの履歴は [docs/PLAN-EXECUTED.md](docs/PLAN-EXECUTED.md) を参照
- Repository層: 29/29 完了
- 統合テスト: 27/27 完了（bumon_codes, bumon_code_refsは除外）
- gRPCサーバー層: 27/27 完了
- OAuth2認証: Google/LINE対応完了
