package utils

import (
	"log"
	"net/http"
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
