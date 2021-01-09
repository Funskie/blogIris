package middlewares

import (
	"net/http"

	"github.com/Funskie/blogIris/api/auth"
	"github.com/Funskie/blogIris/api/responses"
)

func SetMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

func SetMiddlewareAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := auth.TokenValid(r)
		if err != nil {
			responses.Error(w, http.StatusUnauthorized, err)
			return
		}
		next(w, r)
	}
}
