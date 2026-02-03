# Go Gin Todo API

Go言語とGinフレームワークを使用したTodo管理APIです。JWT認証機能を備え、ユーザーごとのTodo管理が可能です。

## 技術スタック

- **言語**: Go 1.25
- **フレームワーク**: Gin
- **データベース**: PostgreSQL
- **ORM**: GORM
- **認証**: JWT (Access Token + Refresh Token)
- **パスワードハッシュ**: bcrypt

## 機能

- ✅ ユーザー登録・ログイン
- ✅ JWT認証（Access Token + Refresh Token）
- ✅ トークンリフレッシュ機能
- ✅ ログアウト機能
- ✅ TodoのCRUD操作
- ✅ ユーザーごとのTodo管理

## セットアップ

### 前提条件

- Go 1.25以上
- PostgreSQL 16以上
- Docker & Docker Compose（オプション）

### 1. リポジトリのクローン

```bash
git clone <repository-url>
cd go-gin-todo-api
```

### 2. 依存関係のインストール

```bash
go mod download
```

### 3. 環境変数の設定

`.env.local`ファイルを作成し、以下の環境変数を設定してください：

```env
APP_ENV=dev
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=todo
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=todo
JWT_SECRET=your-secret-key-here
ACCESS_TOKEN_TTL_MIN=15
REFRESH_TOKEN_TTL_HOUR=720
```

**重要**: `JWT_SECRET`は本番環境では強力なランダム文字列に変更してください。

### 4. データベースの起動

#### Docker Composeを使用する場合

```bash
docker-compose up -d
```

#### ローカルでPostgreSQLを起動する場合

PostgreSQLを起動し、`DB_NAME`で指定したデータベースを作成してください：

```sql
CREATE DATABASE todo;
```

### 5. アプリケーションの起動

```bash
go run main.go
```

アプリケーションは`http://localhost:8080`で起動します。

## APIエンドポイント

### 認証エンドポイント（認証不要）

| メソッド | エンドポイント | 説明 |
|---------|--------------|------|
| POST | `/auth/register` | ユーザー登録 |
| POST | `/auth/login` | ログイン |
| POST | `/auth/refresh` | トークンリフレッシュ |
| POST | `/auth/logout` | ログアウト |

### 認証必須エンドポイント

すべてのリクエストに`Authorization: Bearer <access_token>`ヘッダーが必要です。

| メソッド | エンドポイント | 説明 |
|---------|--------------|------|
| GET | `/health` | ヘルスチェック |
| GET | `/me` | 現在のユーザー情報取得 |
| GET | `/todos` | Todo一覧取得 |
| POST | `/todos` | Todo作成 |
| GET | `/todos/:id` | Todo詳細取得 |
| PATCH | `/todos/:id` | Todo更新 |
| DELETE | `/todos/:id` | Todo削除 |

## API使用例

### 1. ユーザー登録

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

### 2. ログイン

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

レスポンス例：
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 900
}
```

### 3. Todo作成

```bash
curl -X POST http://localhost:8080/todos \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <access_token>" \
  -d '{
    "title": "買い物に行く"
  }'
```

### 4. Todo一覧取得

```bash
curl -X GET http://localhost:8080/todos \
  -H "Authorization: Bearer <access_token>"
```

### 5. Todo更新

```bash
curl -X PATCH http://localhost:8080/todos/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <access_token>" \
  -d '{
    "completed": true
  }'
```

### 6. トークンリフレッシュ

```bash
curl -X POST http://localhost:8080/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "<refresh_token>"
  }'
```

## プロジェクト構造

```
go-gin-todo-api/
├── main.go                 # エントリーポイント
├── database/
│   └── database.go         # データベース接続設定
├── handlers/
│   ├── auth.go             # 認証ハンドラー
│   ├── todo.go             # Todoハンドラー
│   └── user.go             # ユーザーハンドラー
├── middleware/
│   └── auth.go             # JWT認証ミドルウェア
├── models/
│   └── model.go            # データモデル定義
├── utils/
│   ├── token.go            # JWTトークン生成・検証
│   ├── password.go         # パスワードハッシュ化
│   ├── errors.go           # エラーレスポンス
│   └── db_errors.go        # DBエラーハンドリング
├── docker-compose.yml      # Docker Compose設定
├── Dockerfile              # Dockerイメージ設定
├── go.mod                  # Go依存関係
└── README.md               # このファイル
```

## データベース

アプリケーション起動時に自動的にマイグレーションが実行され、以下のテーブルが作成されます：

- `users`: ユーザー情報
- `todos`: Todo情報
- `refresh_tokens`: リフレッシュトークン管理

## 環境変数

| 変数名 | 説明 | デフォルト |
|--------|------|-----------|
| `APP_ENV` | アプリケーション環境 | `dev` |
| `PORT` | サーバーポート | `8080` |
| `DB_HOST` | データベースホスト | `localhost` |
| `DB_PORT` | データベースポート | `5432` |
| `DB_USER` | データベースユーザー名 | `postgres` |
| `DB_PASSWORD` | データベースパスワード | - |
| `DB_NAME` | データベース名 | `todo` |
| `JWT_SECRET` | JWT署名用シークレットキー | - |
| `ACCESS_TOKEN_TTL_MIN` | Access Token有効期限（分） | `15` |
| `REFRESH_TOKEN_TTL_HOUR` | Refresh Token有効期限（時間） | `720` |

## Dockerでの実行

### ビルドと起動

```bash
docker-compose up --build
```

### 停止

```bash
docker-compose down
```

### データベースのデータを保持したまま停止

```bash
docker-compose down
```

### データベースのデータも削除して停止

```bash
docker-compose down -v
```

## 開発

### テスト実行

```bash
go test ./...
```

### コードフォーマット

```bash
go fmt ./...
```

### 依存関係の更新

```bash
go get -u ./...
go mod tidy
```

## ライセンス

このプロジェクトはMITライセンスの下で公開されています。
