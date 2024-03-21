package server

import "net/http"

func respondOk(w http.ResponseWriter, body []byte) {
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func respondCreated(w http.ResponseWriter, body []byte) {
	w.WriteHeader(http.StatusCreated)
	w.Write(body)
}

func respondInternalErr(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}

func respondBadRequestErr(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(err.Error()))
}
