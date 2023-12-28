package utils

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var SecretKey = "SECRET_TOKEN"

func GenerateToken(claims *jwt.MapClaims) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	webtoken, err := token.SignedString([]byte(SecretKey))

	if err != nil {
		return "", err
	}

	return webtoken, nil

}

func CheckPasswordHash(password, hashPassword string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))

	return err == nil

}

func HashingPassword(password string) (string, error) {
	hashedByte, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		return "", err
	}

	return string(hashedByte), nil
}

func VerifycationToken(tokenString string) (*jwt.Token, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func DecodeToken(tokenString string) (jwt.MapClaims, error) {
	token, err := VerifycationToken(tokenString)

	if err != nil {
		return nil, err
	}

	claim, isOk := token.Claims.(jwt.MapClaims)

	if isOk && token.Valid {
		return claim, nil
	}

	return nil, fmt.Errorf("invalid token")
}
