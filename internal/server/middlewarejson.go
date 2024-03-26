package server

import "net/http"

func jsonContent(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
