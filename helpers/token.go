package helpers

import (
	"github.com/authentication-app-server/api-services/models"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

var hmacSampleSecret []byte

func GenerateToken(email string) (string, error) {

	if hmacSampleSecret == nil {
		hmacSampleSecret = make([]byte, 32)
	}

	issuedAt := &jwt.NumericDate{Time: time.Now()}
	expiresAt := &jwt.NumericDate{Time: time.Now().Add(24 * time.Hour)}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, models.CustomClaims{
		TokenType: "level1",
		Email:     email,
		RegisteredClaims: &jwt.RegisteredClaims{
			Issuer:    email,
			ExpiresAt: expiresAt,
			NotBefore: issuedAt,
			IssuedAt:  issuedAt,
		},
	})

	signedToken, err := token.SignedString(hmacSampleSecret)

	return signedToken, err
}
