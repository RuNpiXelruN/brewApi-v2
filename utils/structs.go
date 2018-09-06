package utils

import jwt "github.com/dgrijalva/jwt-go"

// Result type
type Result struct {
	Success *Success `json:"success"`
	Error   *Error   `json:"error"`
}

// Success type
type Success struct {
	StatusCode int         `json:"status_code"`
	Data       interface{} `json:"data"`
	Token      *string     `json:"token"`
}

// Error type
type Error struct {
	StatusCode int    `json:"status_code"`
	StatusText string `json:"error_text"`
}

// CustomClaims struct
type CustomClaims struct {
	Content interface{} `json:"content"`
	jwt.StandardClaims
}
