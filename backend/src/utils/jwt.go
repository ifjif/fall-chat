package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrSigningMethod = errors.New("unexpected signing method")
	ErrInvalidToken  = errors.New("invalid token")
	JWT_SECRET       = []byte("fall-chat-jwt")
)

const (
	TOKEN_ACCESS_TYPE               = "access"
	TOKEN_FILE_ACCESS_TYPE          = "file_access"
	TOKEN_REFRESH_TYPE              = "refresh"
	TOKEN_ACCESS_EXPIRE_DURATION    = 30 * time.Minute
	TOKEN_REFRESH_EXPIRE_DURATION   = 7 * 24 * time.Hour
	TOKEN_FILE_ACCESS_TYPE_DURATION = 5 * time.Minute
)

type CustomClaims struct {
	UserID               string `json:"user_id"`
	Username             string `json:"username"`
	DeviceID             string `json:"device_id"`
	JTI                  string `json:"jti"`  // 唯一标识
	Type                 string `json:"type"` // access / refresh / file_access
	jwt.RegisteredClaims        // 标准声明
}

// 长效token
func GenerateRefreshToken(userID, username, deviceID string) (string, string, int64, error) {
	jti := uuid.NewString()
	token, expiresIn, err := GenerateJWT(TOKEN_REFRESH_TYPE, jti, userID, username, deviceID, TOKEN_REFRESH_EXPIRE_DURATION)

	return jti, token, expiresIn, err
}

// 文件访问 token
func GenerateFileAccessToken(userID, username, filename string) (string, int64, error) {
	token, expiresIn, err := GenerateJWT(TOKEN_FILE_ACCESS_TYPE, filename, userID, username, "", TOKEN_FILE_ACCESS_TYPE_DURATION)
	return token, expiresIn, err
}

// 短效token
func GenerateAccessToken(userID, username, deviceID string) (string, int64, error) {
	token, expiresIn, err := GenerateJWT(TOKEN_ACCESS_TYPE, "", userID, username, deviceID, TOKEN_ACCESS_EXPIRE_DURATION)
	return token, expiresIn, err
}

// 返回token和 过期秒数
func GenerateJWT(type_, jti, userID, username, deviceID string, duration time.Duration) (string, int64, error) {
	expireTime := time.Now().Add(duration)

	claims := CustomClaims{
		UserID:   userID,
		Username: username,
		DeviceID: deviceID,
		Type:     type_,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()), // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()), // 生效时间
			Issuer:    "fall-chat",                    // 签发者
		},
	}

	// hs256 算法创建 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名
	tokenString, err := token.SignedString(JWT_SECRET)
	if err != nil {
		return "", 0, err
	}

	// 距过期时间的秒数
	expiresIn := int64(time.Until(expireTime).Seconds())

	return tokenString, expiresIn, nil
}

// 解析并验证jwt
func ParseJWT(tokenString string) (*CustomClaims, error) {
	claims := &CustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		// 确保签名算法是 HMAC,防止混淆攻击
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrSigningMethod
		}

		return JWT_SECRET, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
