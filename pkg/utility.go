package pkg

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"os"
	"regexp"
	"time"
)

var keyString = os.Getenv("JWT_KEY")
var jwtKey = []byte(keyString)

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

func IsValidEmail(email string) bool {
	// 使用正则表达式来验证邮箱格式
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return re.MatchString(email)
}
