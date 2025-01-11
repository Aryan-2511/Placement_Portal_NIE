package utils

import (
	// "errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte(os.Getenv("jwt_secret_key"))

// GenerateToken generates a JWT for the given email, role, and name
func GenerateToken(email, role, name string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"role":  role,
		"name":  name,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
    if err != nil {
    return "", err
    }

 	return tokenString, nil
}

// ValidateToken validates the given JWT and returns the claims if valid
// ValidateToken validates a JWT and extracts claims
// func ValidateToken(tokenString, secretKey string) (map[string]interface{}, error) {
// 	// Parse the token
// 	// token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 	// 	// Ensure the token uses HMAC signing method
// 	// 	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 	// 		return nil, errors.New("unexpected signing method")
// 	// 	}
// 	// 	return []byte(secretKey), nil
// 	// })
// 	// if err != nil {
// 	// 	return nil, err
// 	// }

// 	// // Extract claims
// 	// if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 	// 	// Check expiration
// 	// 	if exp, ok := claims["exp"].(float64); ok {
// 	// 		if time.Now().Unix() > int64(exp) {
// 	// 			return nil, errors.New("token has expired")
// 	// 		}
// 	// 	} else {
// 	// 		return nil, errors.New("exp claim is missing in token")
// 	// 	}

// 	// 	// Ensure role exists in claims
// 	// 	if _, ok := claims["role"]; !ok {
// 	// 		return nil, errors.New("role claim is missing in token")
// 	// 	}

// 	// 	// Return claims as a map
// 	// 	return claims, nil
// 	// }

// 	// return nil, errors.New("invalid token")
	
// }
func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	   return jwtSecret, nil
	})
   
	if err != nil {
	   return err
	}
   
	if !token.Valid {
	   return fmt.Errorf("invalid token")
	}
   
	return nil
 }