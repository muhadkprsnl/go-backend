// package utils

// import (
// 	"time"

// 	"github.com/golang-jwt/jwt/v4"
// )

// var jwtKey = []byte("your-secret-key")

//	func GenerateJWT(username string) (string, error) {
//		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
//			"username": username,
//			"exp":      time.Now().Add(time.Hour * 24).Unix(),
//		})
//		return token.SignedString(jwtKey)
//	}
package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("your-secret-key") // Ideally from env

func GenerateJWT(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
