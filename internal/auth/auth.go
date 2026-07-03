package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

var ErrTokenSignin = errors.New("token signing error")

type Claims struct {
	PasswordHash string `json:"password_hash"`
	jwt.RegisteredClaims
}

func GenerateToken(pass string) (string, error) {
	claims := Claims{
		PasswordHash: hashPassword(pass),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(pass))
}

func ValidateToken(tokenString, servPass string) bool {
	claims := Claims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); ok {
			return []byte(servPass), nil
		}
		return nil, ErrTokenSignin
	})
	if err != nil || !token.Valid {
		return false
	}
	return claims.PasswordHash == hashPassword(servPass)
}

func hashPassword(pass string) string {
	sum := sha256.Sum256([]byte(pass))
	return hex.EncodeToString(sum[:])
}
