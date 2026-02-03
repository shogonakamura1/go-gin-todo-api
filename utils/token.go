package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims JWTのクレーム構造
type JWTClaims struct {
	UserID uint `json:"sub"`
	jwt.RegisteredClaims
}

// GenerateAccessToken アクセストークンを生成します
func GenerateAccessToken(userID uint) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET is not set")
	}

	ttlMinutes := 15 // デフォルト値
	if ttlStr := os.Getenv("ACCESS_TOKEN_TTL_MIN"); ttlStr != "" {
		if parsed, err := strconv.Atoi(ttlStr); err == nil {
			ttlMinutes = parsed
		}
	}

	now := time.Now()
	expiresAt := now.Add(time.Duration(ttlMinutes) * time.Minute)

	claims := JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateAccessToken アクセストークンを検証し、user_idを返します
func ValidateAccessToken(tokenString string) (uint, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return 0, fmt.Errorf("JWT_SECRET is not set")
	}

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims.UserID, nil
	}

	return 0, fmt.Errorf("invalid token")
}

// GenerateRefreshToken ランダムなリフレッシュトークンを生成します
func GenerateRefreshToken() (string, error) {
	bytes := make([]byte, 32) // 256ビット
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// HashRefreshToken リフレッシュトークンをSHA256でハッシュ化します
func HashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// GetRefreshTokenTTL リフレッシュトークンの有効期限を取得します
func GetRefreshTokenTTL() time.Duration {
	ttlHours := 720 // デフォルト値（30日）
	if ttlStr := os.Getenv("REFRESH_TOKEN_TTL_HOUR"); ttlStr != "" {
		if parsed, err := strconv.Atoi(ttlStr); err == nil {
			ttlHours = parsed
		}
	}
	return time.Duration(ttlHours) * time.Hour
}
