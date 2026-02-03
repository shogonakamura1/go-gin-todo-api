# ビルドステージ
FROM golang:1.25-alpine AS build

# 作業ディレクトリを設定
WORKDIR /app

# 依存関係ファイルをコピー
COPY go.mod go.sum ./

# 依存関係をダウンロード（キャッシュを活用）
RUN go mod download

# アプリケーションコードをコピー
COPY . .

# アプリケーションをビルド
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# ランタイムステージ
FROM alpine:3.21

# パッケージをアップグレードしてからca-certificatesをインストール
# ca-certificatesはセキュリティパッチが頻繁に適用されるため、バージョン固定は推奨されない
# hadolint ignore=DL3018
RUN apk --no-cache upgrade && \
    apk --no-cache add ca-certificates

# セキュリティのため非rootユーザーを作成
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /app

# ビルドステージから実行ファイルをコピー
COPY --from=build /app/main .

# 非rootユーザーに切り替え
USER appuser

# ポート8080を公開
EXPOSE 8080

# アプリケーションを実行
CMD ["./main"]
