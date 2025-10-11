# Seeder

このseederは`category_masters`と`question_template_masters`テーブルにサンプルデータを挿入します。

## 実行方法

### 1. 環境変数の設定（オプション）
```bash
export DATABASE_URL="root:password@tcp(localhost:3306)/english_comp_app?charset=utf8mb4&parseTime=True&loc=Local"
```

### 2. Seederの実行
```bash
go run cmd/seeder/main.go
```

## 挿入されるデータ

### CategoryMasters
- **5件のカテゴリ**が挿入されます
- カテゴリ一覧：
  - 日常会話
  - 翻訳練習
  - 文法問題
  - ビジネス英語
  - ディスカッション
- 各カテゴリのIDはUUID形式で自動生成されます

### QuestionTemplateMasters
- **12件のサンプル問題**が挿入されます
- 問題の種類：
  - Essay（作文）- 4問
  - Translation（翻訳）- 4問
  - Fill in the blank（穴埋め）- 4問
- 難易度レベル：
  - beginner（初級）
  - intermediate（中級）
  - advanced（上級）
- 想定回答時間：5-30分
- ポイント：5-25点
- 各問題のIDはUUID形式で自動生成されます
- カテゴリIDは対応するCategoryMastersのUUIDと関連付けられます

## 注意事項

- 実行前に既存のデータは削除されます
- データベース接続が必要です
- マイグレーションが自動実行されます
