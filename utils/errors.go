package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorCode エラーコード
type ErrorCode string

const (
	ErrorCodeInvalidRequest ErrorCode = "invalid_request"
	ErrorCodeUnauthorized   ErrorCode = "unauthorized"
	ErrorCodeForbidden      ErrorCode = "forbidden"
	ErrorCodeNotFound       ErrorCode = "not_found"
	ErrorCodeConflict       ErrorCode = "conflict"
	ErrorCodeInternal       ErrorCode = "internal"
)

// ErrorResponse エラーレスポンス構造体
type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// RespondError 統一されたエラーレスポンスを返す
func RespondError(c *gin.Context, statusCode int, code ErrorCode, message string) {
	c.JSON(statusCode, ErrorResponse{
		Error: struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}{
			Code:    string(code),
			Message: message,
		},
	})
}

// RespondBadRequest 400 Bad Requestを返す
func RespondBadRequest(c *gin.Context, message string) {
	RespondError(c, http.StatusBadRequest, ErrorCodeInvalidRequest, message)
}

// RespondUnauthorized 401 Unauthorizedを返す
func RespondUnauthorized(c *gin.Context, message string) {
	RespondError(c, http.StatusUnauthorized, ErrorCodeUnauthorized, message)
}

// RespondForbidden 403 Forbiddenを返す
func RespondForbidden(c *gin.Context, message string) {
	RespondError(c, http.StatusForbidden, ErrorCodeForbidden, message)
}

// RespondNotFound 404 Not Foundを返す
func RespondNotFound(c *gin.Context, message string) {
	RespondError(c, http.StatusNotFound, ErrorCodeNotFound, message)
}

// RespondConflict 409 Conflictを返す
func RespondConflict(c *gin.Context, message string) {
	RespondError(c, http.StatusConflict, ErrorCodeConflict, message)
}

// RespondInternalError 500 Internal Server Errorを返す
func RespondInternalError(c *gin.Context, message string) {
	RespondError(c, http.StatusInternalServerError, ErrorCodeInternal, message)
}
