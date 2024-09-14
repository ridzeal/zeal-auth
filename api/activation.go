package api

import "net/http"

func Activation(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
