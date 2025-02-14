package middleware

import (
	"net/http"

	"github.com/Backend/go-jwt-gin-gorm/config"
	"github.com/Backend/go-jwt-gin-gorm/helpers"
	"github.com/golang-jwt/jwt/v4"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				response := map[string]string{"message": "Unauthorized"}
				helpers.ResponseJSON(w, http.StatusUnauthorized, response)
				return
			}
		}

		// mengambil token value
		tokenString := c.Value

		claims := &config.JWTClaims{}

		//parsing token JWT
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return config.JWT_KEY, nil
		})

		if err != nil {
			v, _ := err.(*jwt.ValidationError)
			switch v.Errors {
			case jwt.ValidationErrorSignatureInvalid:
				//token invalid
				response := map[string]string{"message": "Unauthorized"}
				helpers.ResponseJSON(w, http.StatusUnauthorized, response)
				return
			case jwt.ValidationErrorExpired:
				//token expired
				response := map[string]string{"message": "Unauthorized, Token Expired!!"}
				helpers.ResponseJSON(w, http.StatusUnauthorized, response)
				return
			default:
				response := map[string]string{"message": "Unauthorized"}
				helpers.ResponseJSON(w, http.StatusUnauthorized, response)
				return
			}
		}

		if !token.Valid {
			response := map[string]string{"message": "Unauthorized"}
			helpers.ResponseJSON(w, http.StatusUnauthorized, response)
			return
		}

		next.ServeHTTP(w, r)
	})
}
