package server

import "net/http"

func handleError(fn HandlerFuncErr) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			switch err.Error() {
			case "this player is not playing this game":
				respondBadRequestErr(w, err)
			default:
				respondInternalErr(w, err)
			}
		}
	}
}
