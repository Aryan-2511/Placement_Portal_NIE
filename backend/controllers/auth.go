package controllers

import (
	// "Github.com/Aryan-2511/Placement_NIE/utils"
	// "net/http"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword generates a secure hash from plain text password
// Uses bcrypt with default cost factor
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}


// func AuthMiddleware(next http.Handler) http.Handler{
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
// 		tokenString := r.Header.Get("Authorization")
// 		if tokenString == ""{
// 			http.Error(w,"Authorization token required",http.StatusUnauthorized)
// 			return
// 		}
// 		claims,err := utils.ValidateToken(tokenString)
// 		if err!=nil{
// 			http.Error(w, "Invalid token", http.StatusUnauthorized)
// 			return
// 		}
// 		r.Header.Set("User",claims["USN"].(string))
// 		r.Header.Set("Role",claims["role"].(string))

// 		next.ServeHTTP(w,r)
// 	})
// }