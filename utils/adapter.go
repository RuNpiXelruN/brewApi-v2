package utils

import (
	"fmt"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// Adapter type
type Adapter func(http.HandlerFunc) http.HandlerFunc

// Adapt func
func Adapt(h http.HandlerFunc, adapters ...Adapter) http.HandlerFunc {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

// SayHi func
func SayHi() Adapter {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			log.Println("Hi from Middleware!")

			h.ServeHTTP(w, req)
		}
	}
}

// CheckToken func
func CheckToken() Adapter {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			var result Result
			result.Error = &Error{
				StatusCode: http.StatusUnauthorized,
				StatusText: http.StatusText(http.StatusUnauthorized),
			}

			tokenString := req.Header.Get("BrewToken")

			if tokenString == "null" || len(tokenString) < 1 {

				Respond(w, &result)
				return
			}

			token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}

				return signingKey, nil
			})

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				// save user in context
				expFloat64 := claims["exp"].(float64)
				expInt64 := int64(expFloat64)
				expTime := time.Unix(expInt64, 0)
				fmt.Println("Elapsed:", time.Since(expTime))
			} else {
				Respond(w, &result)
				return
			}
			h.ServeHTTP(w, req)
		}
	}
}
