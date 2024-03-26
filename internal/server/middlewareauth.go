package server

import (
	"errors"
	"log"
	"net/http"
)

func authenticate(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("authorization")
		path := r.URL.Path
		log.Println("path", path)
		if authHeader == "" && path != "/ui/" && path != "/ui/index.js" && path != "/ui/favicon.ico" {
			respondUnauthorizedErr(w, errors.New("you don't even have token fool"))
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
