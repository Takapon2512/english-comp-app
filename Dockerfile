FROM golang:1.25-alpine

WORKDIR /app

# 必要なパッケージのインストール
RUN apk add --no-cache git tzdata

# タイムゾーンの設定
ENV TZ=Asia/Tokyo

# CompileDaemonのインストール（開発環境用）
RUN go install github.com/githubnemo/CompileDaemon@latest

# 依存関係のコピーとダウンロード
COPY go.mod go.sum ./
RUN go mod download

EXPOSE 8080

# 開発環境ではCompileDaemonを使用
CMD CompileDaemon --build="go build -o main ./cmd/api" --command="./main"