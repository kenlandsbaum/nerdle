package server

import (
	"net/http"
)

func (s Server) handleHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	respondOk(w, []byte(`{"message":"welcome to the app"}`))
}

// func (s Server) handleSlow(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()
// 	processTime := time.Duration(3) * time.Second

// 	select {
// 	case <-ctx.Done():
// 		return

// 	case <-time.After(processTime):
// 	}
// 	respondOk(w, []byte("done"))
// }
