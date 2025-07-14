// utils/jwt.go
package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/entity"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string       `json:"user_id"`
	Role   string     `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint, role entity.UserRole, secret string, expirationMinutes int) (string, error) {
	fmt.Println(expirationMinutes)
	expirationTime := time.Now().Add(time.Duration(expirationMinutes) * time.Minute)
	uidStr := fmt.Sprintf("%d", userID)

	claims := &Claims{
		UserID: uidStr,
		Role:   string(role),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   uidStr, 
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", errors.New("gagal menandatangani token")
	}
	return tokenString, nil
}


func ParseToken(tokenString, secret string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("metode penandatanganan tidak valid")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("token tidak valid")
	}

	return claims, nil
}
