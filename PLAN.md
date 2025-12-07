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
- Cloud Run: デプロイ完了、28サービス動作確認済み ✅
  - URL: https://postgres-prod-747065218280.asia-northeast1.run.app

## 目標

29テーブル全てにgRPCサービスを提供

## 次のステップ（予定）

### Phase 5: 追加機能・改善

プロジェクトは基本機能を完了。今後の拡張案:
- フロントエンド連携 (gRPC-Web)
- 認証フロー統合テスト
- パフォーマンスチューニング
- モニタリング・ロギング強化

---

## 参考

- 完了タスクの履歴は [docs/PLAN-EXECUTED.md](docs/PLAN-EXECUTED.md) を参照
- Repository層: 29/29 完了
- 統合テスト: 27/27 完了（bumon_codes, bumon_code_refsは除外）
- gRPCサーバー層: 27/27 完了
- OAuth2認証: Google/LINE対応完了
