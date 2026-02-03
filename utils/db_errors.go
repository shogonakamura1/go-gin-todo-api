package utils

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

// HandleDBError DBエラーを適切なHTTPステータスコードとエラーメッセージに変換
func HandleDBError(err error) (int, string) {
	if err == nil {
		return 0, ""
	}

	// Record not found
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 404, "Resource not found"
	}

	// Unique constraint violation (PostgreSQL error code 23505)
	errStr := err.Error()
	if strings.Contains(errStr, "duplicate key") ||
		strings.Contains(errStr, "unique constraint") ||
		strings.Contains(errStr, "23505") {
		return 409, "Resource already exists"
	}

	// Foreign key constraint violation (PostgreSQL error code 23503)
	if strings.Contains(errStr, "foreign key constraint") ||
		strings.Contains(errStr, "23503") {
		return 400, "Invalid reference"
	}

	// その他のDBエラー
	return 500, "Database error"
}
