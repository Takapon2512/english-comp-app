# e-comp バックエンドAPI

## 概要
e-compは、初心者向けにLLMを使って英作文を学べるアプリケーションのバックエンドAPIです。

## 技術スタック
- Go 1.22
- MySQL 8.0
- Docker & Docker Compose

## 開発環境のセットアップ

### 必要条件
- Docker
- Docker Compose
- Make（オプション）

### 環境構築手順

1. リポジトリのクローン
```bash
git clone [your-repository-url]
cd english-app
```

2. 環境変数の設定
```bash
cp .env.example .env
# .envファイルを編集して必要な環境変数を設定
```

3. Dockerコンテナの起動
```bash
docker-compose up -d
```

4. APIサーバーの起動確認
```bash
curl http://localhost:8080/health
```

### 開発用コマンド

- サーバーの起動
```bash
docker-compose up -d
```

- サーバーの停止
```bash
docker-compose down
```

- ログの確認
```bash
docker-compose logs -f app
```

- データベースの接続
```bash
docker-compose exec db mysql -u english_app -p english_app
```

## プロジェクト構造
```
.
├── cmd/
│   └── api/
│       └── main.go      # アプリケーションのエントリーポイント
├── internal/
│   ├── auth/           # 認証関連
│   ├── handler/        # HTTPハンドラー
│   ├── middleware/     # ミドルウェア
│   ├── model/          # データモデル
│   └── repository/     # データベースアクセス
├── docker/
│   └── mysql/
│       └── init/       # DBの初期化スクリプト
├── Dockerfile          # アプリケーションのDockerfile
├── docker-compose.yml  # Docker Compose設定
├── go.mod             # Goの依存関係
└── go.sum             # Goの依存関係のチェックサム
```

## APIエンドポイント

### 認証関連
- POST /api/v1/auth/signup - ユーザー登録
- POST /api/v1/auth/login - ログイン
- POST /api/v1/auth/logout - ログアウト

### ヘルスチェック
- GET /health - サーバーの状態確認

## 貢献方法
1. このリポジトリをフォーク
2. 新しいブランチを作成 (`git checkout -b feature/amazing-feature`)
3. 変更をコミット (`git commit -m 'Add some amazing feature'`)
4. ブランチをプッシュ (`git push origin feature/amazing-feature`)
5. プルリクエストを作成

## ライセンス
[ライセンスを記載]
