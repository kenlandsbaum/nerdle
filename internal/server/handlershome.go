package server

import (
	"net/http"
)

func (s Server) handleHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	respondOk(w, []byte(`{"message":"welcome to the app"}`))
}
