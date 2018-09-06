package utils

import (
	"encoding/json"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var signingKey = []byte("SawyerBrooks")

// Respond func
// for handling Responses in middleware only
func Respond(w http.ResponseWriter, result *Result) {
	if result.Error != nil {
		w.Header().Set("Status", http.StatusText(result.Error.StatusCode))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(result.Error.StatusCode)

		data, _ := json.Marshal(result.Error)
		w.Write(data)
		return
	}

	if result.Success.Token != nil {
		w.Header().Set("BrewToken", *result.Success.Token)
	}
	w.Header().Set("Status", http.StatusText(result.Success.StatusCode))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(result.Success.StatusCode)

	data, _ := json.Marshal(result.Success.Data)
	w.Write(data)
}

// StringPointer func
func StringPointer(s string) *string {
	return &s
}

// GetToken func
func GetToken(email string) (*string, error) {
	claims := CustomClaims{
		email,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
			// ExpiresAt: time.Now().Add(time.Second).Unix(),
			Issuer: "Brewsite_backend",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(signingKey)
	if err != nil {
		return nil, err
	}

	return &ss, nil
}
