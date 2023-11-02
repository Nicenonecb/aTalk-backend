package pkg

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"time"
)

var jwtKey = []byte("your_secret_key") // 选择一个强密钥

type Claims struct {
	UserID uuid.UUID `json:"UserID"`
	jwt.StandardClaims
}

func GenerateToken(UserID uuid.UUID) (string, error) {
	expirationTime := time.Now().Add(100 * time.Hour)
	claims := &Claims{
		UserID: UserID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ValidateToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if !tkn.Valid {
		return nil, jwt.NewValidationError("Token has expired or is not valid", jwt.ValidationErrorExpired)
	}
	return claims, err
}
