package utils

import (
	"log"
	"net/http"
)

type Adapter func(http.HandlerFunc) http.HandlerFunc

func Adapt(h http.HandlerFunc, adapters ...Adapter) http.HandlerFunc {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

func SayHi() Adapter {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			log.Println("Hi from Middleware!")

			h.ServeHTTP(w, req)
		}
	}
}
