package utils
import(
	"fmt"
	"net/http"
	"strings"
	"github.com/golang-jwt/jwt/v4"
	

)
func ExtractRoleFromToken(r *http.Request, secretKey string) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("missing authorization header")
	}

	// The Authorization header is expected to be in the format "Bearer <token>"
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return "", fmt.Errorf("invalid authorization header format")
	}

	// Parse and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is correct
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return "", fmt.Errorf("error parsing token: %v", err)
	}

	// Extract the role from the token claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		role, ok := claims["role"].(string)
		if !ok {
			return "", fmt.Errorf("role not found in token")
		}
		return role, nil
	}

	return "", fmt.Errorf("invalid token")
}
